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
	CreateCheckoutPage(req *CheckoutRequest) (string, error)
	GetPaymentStatus(orderID string) (*PaymentStatusResponse, error)
	RefundPayment(orderID string, amount string) (*RefundResponse, error)
}

type client struct {
	config     *Config
	httpClient *http.Client
}

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

func (c client) sign(data []byte) string {
	hasher := sha1.New()
	hasher.Write([]byte(c.config.PrivateKey))
	hasher.Write(data)
	hasher.Write([]byte(c.config.PrivateKey))
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
}

func (с client) encode(payload any) (string, error) {
	obj, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to encode request body: %w", err)
	}
	return base64.StdEncoding.EncodeToString(obj), nil
}

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

func (c client) sendClientRequest(payload any) (*http.Response, error) {
	injectedPayload, err := c.injectMissingKeys(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to inject missing keys: %w", err)
	}

	encodedJSON, err := c.encode(injectedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}
	signature := c.sign([]byte(encodedJSON))

	formData := url.Values{
		"data":      {encodedJSON},
		"signature": {signature},
	}

	resp, err := c.httpClient.PostForm(ClientServerURL, formData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse liqpay form: %w", err)
	}
	defer resp.Body.Close()

	return resp, nil
}

func (c client) getClientRedirectURL(resp *http.Response) (string, error) {
	var location string

	if resp.StatusCode == http.StatusFound {
		location = resp.Header.Get("Location")
		if location == "" {
			return "", errors.New("redirection response missing Location header")
		}

		return location, nil
	}

	return "", fmt.Errorf("redirect not found")
}

func (c client) prepareServerRequest(payload any) (*http.Request, error) {
	injectedPayload, err := c.injectMissingKeys(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to inject missing keys: %w", err)
	}

	encodedJSON, err := c.encode(injectedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}
	signature := c.sign([]byte(encodedJSON))

	formData := url.Values{
		"data":      {encodedJSON},
		"signature": {signature},
	}

	reqBody := bytes.NewBufferString(formData.Encode())
	req, err := http.NewRequest(http.MethodPost, ServerServerURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create new http request: %w", err)
	}

	return req, nil
}

func (c client) sendServerRequest(req *http.Request, v any) error {
	if c.config.Debug {
		log.Printf("[LiqPay Request] method: %s, url: %s\n", req.Method, req.URL.String())
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}

	if v == nil {
		return nil
	}

	jsonResp, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	if c.config.Debug {
		log.Printf("[LiqPay Response] %s", string(jsonResp))
	}

	if err := json.Unmarshal(jsonResp, v); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if res["status"] == "error" || res["status"] == "failure" || res["result"] == "error" {
		errResp := &APIError{
			Status: res["status"].(string),
			Code:   res["err_code"].(string),
			Desc:   res["err_description"].(string),
		}

		if c.config.Debug {
			log.Printf("[LiqPay Request Error] status: %s, status_code: %d, code: %s, description: %s\n",
				errResp.Status, resp.StatusCode, errResp.Code, errResp.Desc)
		}

		return errResp
	}

	return nil
}

func (c client) CreateCheckoutPage(data *CheckoutRequest) (string, error) {
	data.Action = CheckoutActionPay
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

func (c client) GetPaymentStatus(orderID string) (*PaymentStatusResponse, error) {
	data := &PaymentStatusRequest{Action: "status", OrderID: orderID}
	req, err := c.prepareServerRequest(data)
	if err != nil {
		return nil, err
	}

	v := &PaymentStatusResponse{}
	err = c.sendServerRequest(req, v)
	switch {
	case err != nil && ErrorRefersToAPI(err):
		return v, err
	case err != nil:
		return nil, err
	}

	return v, nil
}

func (c client) RefundPayment(orderID string, amount string) (*RefundResponse, error) {
	data := &RefundRequest{Action: "refund", OrderID: orderID, Amount: amount}
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
