package liqpay

import "time"

type Action string

const (
	ActionPay             Action = "pay"              // Default payment
	ActionHold            Action = "hold"             // Amount of hold on sender's account
	ActionSubscribe       Action = "subscribe"        // Create subscription
	ActionSubscribeUpdate Action = "subscribe_update" // Update subscription
	ActionUnsubscribe     Action = "unsubscribe"      // Unsubscribe
	ActionStatus          Action = "status"           // Payment status
	ActionPayDonate       Action = "paydonate"        // Donation
	ActionPaySplit        Action = "paysplit"         // Splitting payments
	ActionAuth            Action = "auth"             // Card preauth
	ActionRegular         Action = "regular"          // Regular payment
	ActionRefund          Action = "refund"           // Refund payment
	ActionInvoiceSend     Action = "invoice_send"     // Send invoice
	ActionInvoiceCancel   Action = "invoice_cancel"   // Cancel invoice
)

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyEUR Currency = "EUR"
	CurrencyUAH Currency = "UAH"
)

type Language string

const (
	LanguageUK Language = "uk"
	LanguageEN Language = "en"
)

type PayType string

const (
	PayTypeApplePay       PayType = "apay"        // Pay with Apple Pay
	PayTypeGooglePay      PayType = "gpay"        // Pay with Google Pay
	PayTypeCard           PayType = "card"        // Pay with card
	PayTypePrivat24       PayType = "privat24"    // Pay with a Privat24 account
	PayTypeMomentPart     PayType = "moment_part" // Pay with installments
	PayTypePayPart        PayType = "paypart"     // Payment by parts
	PayTypeCash           PayType = "cash"        // Pay with cash
	PayTypeInvoice        PayType = "invoice"     // Pay with an invoice
	PayTypeQRCodeScanning PayType = "qr"          // Pay through QR code scanning
)

type Status string

const (
	StatusError    Status = "error"    // Failed payment. Data is incorrect
	StatusFailure  Status = "failure"  // Failed payment
	StatusReversed Status = "reversed" // Payment refunded
	StatusSuccess  Status = "success"  // Successful payment
)

type Item struct {
	Amount float64 `json:"amount"` // Quantity/volume
	Cost   string  `json:"cost"`   // The cost of all units of the specified product in the receipt (number of units * unit cost)
	ID     string  `json:"id"`     // Item ID. You can get it in the Liqpay account - SCR - Kasa - Goods
	Price  string  `json:"price"`  // Unit cost of goods
}

type RROInfo struct {
	Items          []Item   `json:"items,omitempty"`           // Data about products for which payment is performed
	DeliveryEmails []string `json:"delivery_emails,omitempty"` // List of e-mails to which receipts should be sent after fiscalization

}

type CheckoutRequest struct {
	Action      Action    `json:"action"`                 // Transaction type
	Amount      float64   `json:"amount"`                 // Payment amount
	Currency    Currency  `json:"currency"`               // Payment currency
	Description string    `json:"description"`            // Payment description
	OrderID     string    `json:"order_id"`               // Unique purchase ID in your shop. Maximum length is 255 symbols
	RROInfo     RROInfo   `json:"rro_info,omitempty"`     // Data for fiscalization
	ExpiredDate string    `json:"expired_date,omitempty"` // Date and time until which customer is able to pay invoice by UTC. Should be sent in the following format 2016-04-24 00:00:00
	Language    Language  `json:"language,omitempty"`     // Customer's language
	PayTypes    []PayType `json:"pay_types,omitempty"`    // Parameter that gets the methods of payments that displayed on checkout. If the parameter is not passed, shop settings will be applied, Checkout tab
	ResultURL   string    `json:"result_url,omitempty"`   // URL of your shop where the buyer would be redirected after completion of the purchase. Maximum length 510 symbols
	ServerURL   string    `json:"server_url,omitempty"`   // URL API in your store for notifications of payment status change (server -> server). Maximum length is 510 symbols
	VerifyCode  string    `json:"verifycode,omitempty"`   // Possible value Y. Dynamic verification code is generated and going back to Callback. Also generated code will be transferred to verification transactions for displaying in statement by client's card. Works for action = auth
}

type StatusRequest struct {
	Action  Action `json:"action"`   // Transaction type
	OrderID string `json:"order_id"` // Unique purchase ID in your shop. Maximum length is 255 symbols
}

type StatusResponse struct {
	AcqID              int      `json:"acq_id"`              // Acquirer ID
	Action             Action   `json:"action"`              // Transaction type: pay, hold, paysplit, subscribe, paydonate, auth, regular
	AgentCommission    float64  `json:"agent_commission"`    // Agent commission in payment currency
	Amount             float64  `json:"amount"`              // Payment amount
	AmountBonus        float64  `json:"amount_bonus"`        // Payer bonus amount in payment currency debit
	AmountCredit       float64  `json:"amount_credit"`       // Payment amount for credit in currency of currency_credit
	AmountDebit        float64  `json:"amount_debit"`        // Payment amount for debit in currency of currency_debit
	AuthCodeCredit     string   `json:"authcode_credit"`     // Authorization code for transaction of credit
	AuthCodeDebit      string   `json:"authcode_debit"`      // Authorization code for transaction of debit
	BonusProcent       float64  `json:"bonus_procent"`       // Discount rate in percent
	BonusType          string   `json:"bonus_type"`          // Bonus type: bonusplus, discount_club, personal, promo
	CardToken          string   `json:"card_token"`          // Sender's card token
	CommissionCredit   float64  `json:"commission_credit"`   // Commission from the receiver in currency_credit
	CommissionDebit    float64  `json:"commission_debit"`    // Commission from the sender in currency_debit
	CreateDate         int64    `json:"create_date"`         // Date of payment creation
	Currency           Currency `json:"currency"`            // Payment currency
	CurrencyCredit     string   `json:"currency_credit"`     // Transaction currency of credit
	CurrencyDebit      string   `json:"currency_debit"`      // Transaction currency of debit
	Description        string   `json:"description"`         // Payment description
	EndDate            int64    `json:"end_date"`            // Date of payment edition/end
	Info               string   `json:"info"`                // Additional payment information
	IP                 string   `json:"ip"`                  // Sender's IP address
	Is3DS              bool     `json:"is_3ds"`              // True if transaction passed with 3DS, false otherwise
	LiqpayOrderID      string   `json:"liqpay_order_id"`     // Payment order_id in LiqPay system
	MomentPart         string   `json:"moment_part"`         // Payment indication in parts
	MPIECI             string   `json:"mpi_eci"`             // MPI ECI: 5 - transaction passed with 3DS, 6 - issuer of payer card doesn't support 3d Secure, 7 - operation passed without 3d Secure
	OrderID            string   `json:"order_id"`            // Order_id payment
	PaymentID          int      `json:"payment_id"`          // Payment id in LiqPay system
	Paytype            string   `json:"paytype"`             // Method of payment: card, privat24, moment_part, cash, invoice, qr
	PublicKey          string   `json:"public_key"`          // Shop public key
	ReceiverCommission float64  `json:"receiver_commission"` // Receiver commission in payment currency
	RRNCredit          string   `json:"rrn_credit"`          // Unique transaction ID in authorization and settlement system of issuer bank for credit
	RRNDebit           string   `json:"rrn_debit"`           // Unique transaction ID in authorization and settlement system of issuer bank for debit
	SenderBonus        float64  `json:"sender_bonus"`        // Sender's bonus in the payment currency
	SenderCardBank     string   `json:"sender_card_bank"`    // Sender's card bank
	SenderCardCountry  int      `json:"sender_card_country"` // Sender's card country (digital ISO 3166-1 code)
	SenderCardMask2    string   `json:"sender_card_mask2"`   // Sender's card mask2
	SenderCardType     string   `json:"sender_card_type"`    // Sender's card type (MC/Visa)
	SenderCommission   float64  `json:"sender_commission"`   // Commission from the sender in the payment currency
	SenderPhone        string   `json:"sender_phone"`        // Sender's phone number
	Status             Status   `json:"status"`              // Payment status
}

type RefundRequest struct {
	Action  Action `json:"action"`   // Transaction type
	Amount  string `json:"amount"`   // Payment amount. For example: 5, 7.34
	OrderID string `json:"order_id"` // Unique purchase ID in your shop. Maximum length is 255 symbols
}

type RefundResponse struct {
	Action    Action `json:"action"`     // Transaction type
	PaymentID int64  `json:"payment_id"` // Payment id in LiqPay system
	Status    string `json:"status"`     // Payment status
}

type SubscribePeriod string

const (
	SubscribePeriodDaily   SubscribePeriod = "day"   // Daily
	SubscribePeriodWeekly  SubscribePeriod = "week"  // Weekly
	SubscribePeriodMonthly SubscribePeriod = "month" // Monthly
	SubscribePeriodYearly  SubscribePeriod = "year"  // Yearly
)

type SubscriptionRequest struct {
	Action             Action          `json:"action"`                          // Action to perform, e.g., "subscribe"
	Amount             float64         `json:"amount"`                          // Payment amount. For example: 5, 7.34
	Card               string          `json:"card"`                            // Card number of the payer
	CardCVV            string          `json:"card_cvv"`                        // CVV/CVV2
	CardExpMonth       string          `json:"card_exp_month"`                  // Expiry month of the payer's card. For example: 08
	CardExpYear        string          `json:"card_exp_year"`                   // Expiry year of the payer's card. For example: 19
	Currency           Currency        `json:"currency"`                        // Payment currency. Possible values: USD, EUR, UAH
	Description        string          `json:"description"`                     // Payment description
	IP                 string          `json:"ip"`                              // Client IP
	OrderID            string          `json:"order_id"`                        // Unique purchase ID in your shop. Maximum length is 255 symbols
	Phone              string          `json:"phone"`                           // Payer's mobile phone.
	Language           string          `json:"language,omitempty"`              // Customer's language uk, en
	Prepare            string          `json:"prepare,omitempty"`               // Preliminary preparation of the payment
	RecurringByToken   string          `json:"recurringbytoken,omitempty"`      // Generate payer card_token
	Recurring          bool            `json:"recurring,omitempty"`             // Token recurring payment flag
	ServerURL          string          `json:"server_url,omitempty"`            // URL API in your store for notifications of payment status change
	Subscribe          string          `json:"subscribe,omitempty"`             // Regular payment
	SubscribeDateStart string          `json:"subscribe_date_start,omitempty"`  // Date of the first payment
	SubscribePeriod    SubscribePeriod `json:"subscribe_periodicity,omitempty"` // Period of payments
	SenderAddress      string          `json:"sender_address,omitempty"`        // Sender's address
	SenderCity         string          `json:"sender_city,omitempty"`           // Sender's city
	SenderCountryCode  string          `json:"sender_country_code,omitempty"`   // Country code of the sender
	SenderFirstName    string          `json:"sender_first_name,omitempty"`     // Sender's first name
	SenderLastName     string          `json:"sender_last_name,omitempty"`      // Sender's last name
	SenderPostalCode   string          `json:"sender_postal_code,omitempty"`    // Sender's zip code
	Customer           string          `json:"customer,omitempty"`              // Unique customer ID in your shop
	DAE                string          `json:"dae,omitempty"`                   // Detail Addenda
	Info               string          `json:"info,omitempty"`                  // Information to add details to payment
	ProductCategory    string          `json:"product_category,omitempty"`      // Product category in your shop
	ProductDescription string          `json:"product_description,omitempty"`   // Product description in your shop
	ProductName        string          `json:"product_name,omitempty"`          // Product name in your shop
	ProductURL         string          `json:"product_url,omitempty"`           // Product page address
}

type SubscriptionResponse struct {
	AcqID              int64     `json:"acq_id"`              // Acquirer ID
	Action             Action    `json:"action"`              // Transaction type
	AgentCommission    float64   `json:"agent_commission"`    // Agent commission in payment currency
	Amount             float64   `json:"amount"`              // Payment amount
	AmountBonus        float64   `json:"amount_bonus"`        // Payer bonus amount in payment currency debit
	AmountCredit       float64   `json:"amount_credit"`       // Payment amount for credit in currency of currency_credit
	AmountDebit        float64   `json:"amount_debit"`        // Payment amount for debit in currency of currency_debit
	CardToken          string    `json:"card_token"`          // Sender's card token
	CommissionCredit   float64   `json:"commission_credit"`   // Commission from the receiver in currency_credit
	CommissionDebit    float64   `json:"commission_debit"`    // Commission from the sender in currency_debit
	CreateDate         time.Time `json:"create_date"`         // Date of payment creation
	Currency           Currency  `json:"currency"`            // Payment currency
	CurrencyCredit     string    `json:"currency_credit"`     // Transaction currency of credit
	CurrencyDebit      string    `json:"currency_debit"`      // Transaction currency of debit
	Description        string    `json:"description"`         // Payment description
	EndDate            string    `json:"end_date"`            // Date of payment edition/end
	Is3DS              bool      `json:"is_3ds"`              // Whether the transaction passed with 3DS
	LiqpayOrderID      string    `json:"liqpay_order_id"`     // Payment order_id in LiqPay system
	MPIECI             int64     `json:"mpi_eci"`             // MPI ECI value
	OrderID            string    `json:"order_id"`            // Order_id payment
	PaymentID          int64     `json:"payment_id"`          // Payment id in LiqPay system
	PayType            string    `json:"paytype"`             // Methods of payment
	PublicKey          string    `json:"public_key"`          // Shop public key
	ReceiverCommission float64   `json:"receiver_commission"` // Receiver commission in payment currency
	SenderBonus        float64   `json:"sender_bonus"`        // Sender's bonus in the payment currency
	SenderCardBank     string    `json:"sender_card_bank"`    // Sender's card bank
	SenderCardCountry  string    `json:"sender_card_country"` // Sender's card country
	SenderCardMask2    string    `json:"sender_card_mask2"`   // Sender's card
	SenderCardType     string    `json:"sender_card_type"`    // Sender's card type MC/Visa
	SenderCommission   float64   `json:"sender_commission"`   // Commission from the sender in the payment currency
	SenderPhone        string    `json:"sender_phone"`        // Sender's phone number
	Status             Status    `json:"status"`              // Payment status
	TransactionID      int64     `json:"transaction_id"`      // Id transactions in the LiqPay system
	Version            int       `json:"version"`             // Version API
}

type EditSubscriptionRequest struct {
	Action      Action  `json:"action"`      // Action to be performed, in this case, 'subscribe_update'
	Amount      float64 `json:"amount"`      // Payment amount. For example: 5, 7.34
	Currency    string  `json:"currency"`    // Payment currency. Possible values: USD, EUR, UAH
	Description string  `json:"description"` // Payment description
	OrderID     string  `json:"order_id"`    // Unique purchase ID in your system
}

type UnsubscribeRequest struct {
	Action  Action `json:"action"`   // Transaction type
	OrderID string `json:"order_id"` // Unique purchase ID in your shop. Maximum length is 255 symbols
}

type InvoiceItem struct {
	Amount float64 `json:"amount"` // Price per unit
	Count  int     `json:"count"`  // Number of units
	Unit   string  `json:"unit"`   // Units of measurement
	Name   string  `json:"name"`   // Name of the product or service
}

type InvoiceRequest struct {
	Action        Action        `json:"action"`                   // Action type, e.g., "invoice_send"
	Amount        float64       `json:"amount"`                   // Payment amount. For example: 5, 7.34
	Currency      Currency      `json:"currency"`                 // Payment currency. Possible values: USD, EUR, UAH
	Description   string        `json:"description"`              // Payment description
	Email         string        `json:"email"`                    // Customer's e-mail to send invoice (phone or email required parameters for transmission)
	OrderID       string        `json:"order_id"`                 // Unique purchase ID in your shop. Maximum length is 255 symbols
	Phone         string        `json:"phone"`                    // The phone number to which the invoice will be sent as a push notification to the Privat24 mobile application (phone or email required parameters for transmission)
	ActionPayment string        `json:"action_payment,omitempty"` // Transaction type. Possible values: pay, hold, subscribe, paydonate
	ExpiredDate   string        `json:"expired_date,omitempty"`   // Date and time until which customer is able to pay invoice by UTC. Should be sent in the following format 2016-04-24 00:00:00
	Goods         []InvoiceItem `json:"goods,omitempty"`          // Optional list of goods
	Language      Language      `json:"language,omitempty"`       // Customer's language uk, en
	ResultURL     string        `json:"result_url,omitempty"`     // URL of your shop where the buyer would be redirected after completion of the purchase. Maximum length 510 symbols
	ServerURL     string        `json:"server_url,omitempty"`     // URL API in your store for notifications of payment status change (server -> server). Maximum length is 510 symbols
}

type InvoiceResponse struct {
	Action        Action   `json:"action"`          // Transaction type. Possible values: pay, hold, paysplit, subscribe, paydonate, auth, regular
	Amount        float64  `json:"amount"`          // Payment amount
	Currency      Currency `json:"currency"`        // Payment currency
	Description   string   `json:"description"`     // Payment description
	Href          string   `json:"href"`            // Link to invoice
	ID            int      `json:"id"`              // Payment id in LiqPay system
	OrderID       string   `json:"order_id"`        // Order_id payment
	ReceiverType  string   `json:"receiver_type"`   // Receive channel type
	ReceiverValue string   `json:"receiver_value"`  // The value obtained in the parameter receiver_type
	Status        string   `json:"status"`          // Payment status. Possible values: error, failure, success, invoice_wait, token
	Token         string   `json:"token,omitempty"` // Payment token
}

type CancelInvoiceRequest struct {
	Action  Action `json:"action"`   // Transaction type
	OrderID string `json:"order_id"` // Unique purchase ID in your shop. Maximum length is 255 symbols
}

type CancelInvoiceResult string

const (
	CancelInvoiceResultOK    = "ok"
	CancelInvoiceResultError = "error"
)

type CancelInvoiceResponse struct {
	InvoiceID int64               `json:"invoice_id"` // Unique identifier of the invoice
	Result    CancelInvoiceResult `json:"order_id"`   // The result of a request ok or error
}

type Callback struct {
	AcqID              int     `json:"acq_id"`              // ID of the acquirer
	Action             Action  `json:"action"`              // Type of operation: pay, hold, paysplit, subscribe, auth, regular
	AgentCommission    float64 `json:"agent_commission"`    // Agent commission in payment currency
	Amount             float64 `json:"amount"`              // Payment amount
	AmountBonus        float64 `json:"amount_bonus"`        // Sender's bonus in payment currency (debit)
	AmountCredit       float64 `json:"amount_credit"`       // Amount of credit transaction in currency_credit
	AmountDebit        float64 `json:"amount_debit"`        // Amount of debit transaction in currency_debit
	AuthcodeCredit     string  `json:"authcode_credit"`     // Authorization code for credit transaction
	AuthcodeDebit      string  `json:"authcode_debit"`      // Authorization code for debit transaction
	CardToken          string  `json:"card_token"`          // Sender's card token
	CommissionCredit   float64 `json:"commission_credit"`   // Receiver's commission in currency_credit
	CommissionDebit    float64 `json:"commission_debit"`    // Sender's commission in currency_debit
	CompletionDate     string  `json:"completion_date"`     // Date of funds debit
	CreateDate         string  `json:"create_date"`         // Payment creation date
	Currency           string  `json:"currency"`            // Payment currency
	CurrencyCredit     string  `json:"currency_credit"`     // Currency of credit transaction
	CurrencyDebit      string  `json:"currency_debit"`      // Currency of debit transaction
	Customer           string  `json:"customer"`            // Unique identifier of the customer on merchant's site
	Description        string  `json:"description"`         // Payment comment
	EndDate            string  `json:"end_date"`            // End/change date of payment
	ErrCode            string  `json:"err_code"`            // Error code
	ErrDescription     string  `json:"err_description"`     // Error description
	Info               string  `json:"info"`                // Additional information about the payment
	IP                 string  `json:"ip"`                  // Sender's IP address
	Is3DS              bool    `json:"is_3ds"`              // Indicates if the transaction passed 3DS verification
	LiqpayOrderID      string  `json:"liqpay_order_id"`     // Payment order_id in LiqPay system
	MpiEci             int     `json:"mpi_eci"`             // MPI ECI value
	OrderID            string  `json:"order_id"`            // Payment order_id
	PaymentID          int     `json:"payment_id"`          // Payment ID in LiqPay system
	Paytype            string  `json:"paytype"`             // Payment method: card, privat24, masterpass, moment_part, cash, invoice, qr
	PublicKey          string  `json:"public_key"`          // Merchant's public key
	ReceiverCommission float64 `json:"receiver_commission"` // Receiver's commission in payment currency
	RedirectTo         string  `json:"redirect_to"`         // Link to redirect the client for 3DS verification
	RefundDateLast     string  `json:"refund_date_last"`    // Last refund date for the payment
	RRNCredit          string  `json:"rrn_credit"`          // Unique transaction number in issuer and acquiring bank's system (credit)
	RRNDebit           string  `json:"rrn_debit"`           // Unique transaction number in issuer and acquiring bank's system (debit)
	SenderBonus        float64 `json:"sender_bonus"`        // Sender's bonus in payment currency
	SenderCardBank     string  `json:"sender_card_bank"`    // Sender's card bank
	SenderCardCountry  string  `json:"sender_card_country"` // Sender's card country ISO 3166-1 code
	SenderCardMask2    string  `json:"sender_card_mask2"`   // Sender's card mask
	SenderCardType     string  `json:"sender_card_type"`    // Sender's card type (MC/Visa)
	SenderCommission   float64 `json:"sender_commission"`   // Sender's commission in payment currency
	SenderFirstName    string  `json:"sender_first_name"`   // Sender's first name
	SenderLastName     string  `json:"sender_last_name"`    // Sender's last name
	SenderPhone        string  `json:"sender_phone"`        // Sender's phone number
	Status             string  `json:"status"`              // Payment status
	WaitReserveStatus  string  `json:"wait_reserve_status"` // Additional payment status indicating that the current payment is reserved for refund
	Token              string  `json:"token"`               // Payment token
	Type               string  `json:"type"`                // Payment type
	Version            int     `json:"version"`             // API version
	ErrErc             string  `json:"err_erc"`             // Error code
	ProductCategory    string  `json:"product_category"`    // Product category
	ProductDescription string  `json:"product_description"` // Product description
	ProductName        string  `json:"product_name"`        // Product name
	ProductURL         string  `json:"product_url"`         // Product page URL
	RefundAmount       float64 `json:"refund_amount"`       // Refund amount
	Verifycode         string  `json:"verifycode"`          // Verification code
}
