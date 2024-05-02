package liqpay

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type Client interface {
	CreateCheckout(req *CheckoutRequest) (string, error)

	CreateSubscription(req *SubscriptionRequest) (string, error)
	UpdateSubscription(req *EditSubscriptionRequest) (*SubscriptionResponse, error)
	RemoveSubscription(orderID string) (*SubscriptionResponse, error)

	CreateInvoice(req *InvoiceRequest) (*InvoiceResponse, error)
	CancelInvoice(orderID string) (*CancelInvoiceResponse, error)

	Status(orderID string) (*StatusResponse, error)
	Refund(orderID string, amount string) (*RefundResponse, error)

	ValidateCallback(data string, signature string) error
}

type client struct {
	config     *Config
	httpClient *http.Client
}

// NewClient creates a new LiqPay client with the provided configuration and HTTP client.
func NewClient(config *Config, httpClient *http.Client) Client {
	var httpC = http.DefaultClient

	if httpClient != nil {
		httpC = httpClient
	}

	httpC.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &client{
		config:     config,
		httpClient: httpC,
	}
}

// sign generates a signature for the given data using the client's private key.
func (c client) sign(data []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(c.config.PrivateKey))
	hasher.Write(data)
	hasher.Write([]byte(c.config.PrivateKey))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

// encode encodes the payload to base64 format.
func (с client) encode(payload any) (string, error) {
	obj, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("liqpay client: failed to encode request body: %w", err)
	}
	return base64.StdEncoding.EncodeToString(obj), nil
}

// injectMissingKeys injects missing keys (version, public_key) into the payload.
func (c client) injectMissingKeys(payload any) (map[string]interface{}, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &data); err != nil {
		return nil, err
	}

	if data["version"] == nil || data["version"] == "" {
		data["version"] = CurrentAPIVersion
	}

	if data["public_key"] == nil || data["public_key"] == "" {
		data["public_key"] = c.config.PublicKey
	}

	return data, nil
}

// sendClientRequest sends a client-server request to LiqPay API.
func (c client) sendClientRequest(payload any) (*http.Response, error) {
	injectedPayload, err := c.injectMissingKeys(payload)
	if err != nil {
		return nil, fmt.Errorf("liqpay client: failed to inject missing keys: %w", err)
	}

	encodedJSON, err := c.encode(injectedPayload)
	if err != nil {
		return nil, fmt.Errorf("liqpay client: failed to encode payload: %w", err)
	}
	signature := c.sign([]byte(encodedJSON))

	formData := url.Values{
		"data":      {encodedJSON},
		"signature": {signature},
	}

	resp, err := c.httpClient.PostForm(ClientServerURL, formData)
	if err != nil {
		return nil, fmt.Errorf("liqpay client: failed to parse liqpay form: %w", err)
	}
	defer resp.Body.Close()

	return resp, nil
}

// getClientRedirectURL extracts the redirect URL from the HTTP response.
func (c client) getClientRedirectURL(resp *http.Response) (string, error) {
	var location string

	if resp.StatusCode == http.StatusFound {
		location = resp.Header.Get("Location")
		if location == "" {
			return "", errors.New("liqpay client: redirection response missing Location header")
		}

		return location, nil
	}

	return "", fmt.Errorf("redirect not found")
}

// prepareServerRequest prepares a server-server HTTP request to LiqPay API.
func (c client) prepareServerRequest(payload any) (*http.Request, error) {
	injectedPayload, err := c.injectMissingKeys(payload)
	if err != nil {
		return nil, fmt.Errorf("liqpay client: failed to inject missing keys: %w", err)
	}

	encodedJSON, err := c.encode(injectedPayload)
	if err != nil {
		return nil, fmt.Errorf("liqpay client: failed to encode payload: %w", err)
	}
	signature := c.sign([]byte(encodedJSON))

	formData := url.Values{
		"data":      {encodedJSON},
		"signature": {signature},
	}

	reqBody := bytes.NewBufferString(formData.Encode())
	req, err := http.NewRequest(http.MethodPost, ServerServerURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("liqpay client: failed to create new http request: %w", err)
	}

	return req, nil
}

// sendServerRequest sends a server-server request to LiqPay API.
func (c client) sendServerRequest(req *http.Request, v any) error {
	if c.config.Debug {
		log.Printf("[LIQPAY DEBUG] Request method: %s, url: %s\n", req.Method, req.URL.String())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("liqpay client: request failed: %w", err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("liqpay client: failed to decode json: %w", err)
	}

	if v == nil {
		return nil
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("liqpay client: failed to marshal response: %w", err)
	}

	if c.config.Debug {
		log.Printf("[LIQPAY DEBUG] Response: %s", string(jsonResp))
	}

	if err := json.Unmarshal(jsonResp, v); err != nil {
		return fmt.Errorf("liqpay client: failed to unmarshal response: %w", err)
	}

	if res["status"] == "error" || res["status"] == "failure" || res["result"] == "error" {
		errResp := &APIError{
			Status: res["status"].(string),
			Code:   res["err_code"].(string),
			Desc:   res["err_description"].(string),
		}

		if c.config.Debug {
			log.Printf("[LIQPAY DEBUG] Error status: %s, status_code: %d, code: %s, description: %s\n",
				errResp.Status, resp.StatusCode, errResp.Code, errResp.Desc)
		}

		return errResp
	}

	return nil
}

// CreateCheckout creates a new checkout link.
func (c client) CreateCheckout(data *CheckoutRequest) (string, error) {
	data.Action = ActionPay

	resp, err := c.sendClientRequest(data)
	if err != nil {
		return "", err
	}

	link, err := c.getClientRedirectURL(resp)
	if err != nil {
		return "", err
	}

	return link, nil
}

// CreateSubscription creates a new subscription link.
func (c client) CreateSubscription(data *SubscriptionRequest) (string, error) {
	data.Action = ActionSubscribe
	data.Subscribe = "1"

	resp, err := c.sendClientRequest(data)
	if err != nil {
		return "", err
	}

	link, err := c.getClientRedirectURL(resp)
	if err != nil {
		return "", err
	}

	return link, nil
}

// UpdateSubscription updates an existing subscription.
func (c client) UpdateSubscription(data *EditSubscriptionRequest) (*SubscriptionResponse, error) {
	data.Action = ActionSubscribeUpdate

	req, err := c.prepareServerRequest(data)
	if err != nil {
		return nil, err
	}

	v := &SubscriptionResponse{}
	err = c.sendServerRequest(req, v)
	switch {
	case err != nil && ErrorRefersToAPI(err):
		return v, err
	case err != nil:
		return nil, err
	}

	return v, nil
}

// RemoveSubscription removes a subscription.
func (c client) RemoveSubscription(orderID string) (*SubscriptionResponse, error) {
	data := &UnsubscribeRequest{Action: ActionUnsubscribe, OrderID: orderID}

	req, err := c.prepareServerRequest(data)
	if err != nil {
		return nil, err
	}

	v := &SubscriptionResponse{}
	err = c.sendServerRequest(req, v)
	switch {
	case err != nil && ErrorRefersToAPI(err):
		return v, err
	case err != nil:
		return nil, err
	}

	return v, nil
}

// CreateInvoice creates a new invoice.
func (c client) CreateInvoice(data *InvoiceRequest) (*InvoiceResponse, error) {
	data.Action = ActionInvoiceSend

	req, err := c.prepareServerRequest(data)
	if err != nil {
		return nil, err
	}

	v := &InvoiceResponse{}
	err = c.sendServerRequest(req, v)
	switch {
	case err != nil && ErrorRefersToAPI(err):
		return v, err
	case err != nil:
		return nil, err
	}

	return v, nil
}

// CancelInvoice cancels an invoice.
func (c client) CancelInvoice(orderID string) (*CancelInvoiceResponse, error) {
	data := &CancelInvoiceRequest{Action: ActionInvoiceCancel, OrderID: orderID}

	req, err := c.prepareServerRequest(data)
	if err != nil {
		return nil, err
	}

	v := &CancelInvoiceResponse{}
	err = c.sendServerRequest(req, v)
	switch {
	case err != nil && ErrorRefersToAPI(err):
		return v, err
	case err != nil:
		return nil, err
	}

	return v, nil
}

// Status retrieves the status of an order.
func (c client) Status(orderID string) (*StatusResponse, error) {
	data := &StatusRequest{Action: ActionStatus, OrderID: orderID}

	req, err := c.prepareServerRequest(data)
	if err != nil {
		return nil, err
	}

	v := &StatusResponse{}
	err = c.sendServerRequest(req, v)
	switch {
	case err != nil && ErrorRefersToAPI(err):
		return v, err
	case err != nil:
		return nil, err
	}

	return v, nil
}

// Refund processes a refund for an order.
func (c client) Refund(orderID string, amount string) (*RefundResponse, error) {
	data := &RefundRequest{Action: ActionRefund, OrderID: orderID, Amount: amount}

	req, err := c.prepareServerRequest(data)
	if err != nil {
		return nil, err
	}

	v := &RefundResponse{}
	err = c.sendServerRequest(req, v)
	switch {
	case err != nil && ErrorRefersToAPI(err):
		return v, err
	case err != nil:
		return nil, err
	}

	return v, nil
}

// ValidateCallback validates the callback data and signature received from LiqPay.
func (c client) ValidateCallback(data string, signature string) error {
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("liqpay client: failed to decode data: %w", err)
	}

	expectedSignature := c.sign(decodedData)
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("liqpay client: failed to decode signature: %w", err)
	}

	if !bytes.Equal(decodedSignature, []byte(expectedSignature)) {
		return errors.New("liqpay client: callback signature verification failed")
	}

	return nil
}
