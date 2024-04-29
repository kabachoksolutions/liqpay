package liqpay

type CheckoutType string

const (
	CheckoutTypePay        CheckoutType = "pay"         // Direct debit from the card
	CheckoutTypeHold       CheckoutType = "hold"        // Blocking funds on the client's card as part of a two-stage payment
	CheckoutTypeSubscribe  CheckoutType = "subscribe"   // Subscription registration
	CheckoutTypePayDonate  CheckoutType = "paydonate"   // Accepting donations with an arbitrary amount
	CheckoutTypeAuth       CheckoutType = "auth"        // Card pre-authorization
	CheckoutTypeSplitRules CheckoutType = "split_rules" // Split payment to several recipients
	CheckoutTypeApplePay   CheckoutType = "apay"        // Payment with Apple Pay
	CheckoutTypeGooglePay  CheckoutType = "gpay"        // Payment with Google Pay
)

type Action string

const (
	CheckoutActionPay       Action = "pay"       // Default payment
	CheckoutActionHold      Action = "hold"      // Amount of hold on sender's account
	CheckoutActionSubscribe Action = "subscribe" // Regular payment
	CheckoutActionPayDonate Action = "paydonate" // Donation
)

type CheckoutCurrency string

const (
	CheckoutCurrencyUSD CheckoutCurrency = "USD"
	CheckoutCurrencyEUR CheckoutCurrency = "EUR"
	CheckoutCurrencyUAH CheckoutCurrency = "UAH"
)

type CheckoutCustomerLanguage string

const (
	CheckoutCustomerLanguageUK CheckoutCustomerLanguage = "uk"
	CheckoutCustomerLanguageEN CheckoutCustomerLanguage = "en"
)

type CheckoutPayType string

const (
	CheckoutPayTypeApplePay       CheckoutPayType = "apay"        // Pay with Apple Pay
	CheckoutPayTypeGooglePay      CheckoutPayType = "gpay"        // Pay with Google Pay
	CheckoutPayTypeCard           CheckoutPayType = "card"        // Pay with card
	CheckoutPayTypePrivat24       CheckoutPayType = "privat24"    // Pay with a Privat24 account
	CheckoutPayTypeMomentPart     CheckoutPayType = "moment_part" // Pay with installments
	CheckoutPayTypePayPart        CheckoutPayType = "paypart"     // Payment by parts
	CheckoutPayTypeCash           CheckoutPayType = "cash"        // Pay with cash
	CheckoutPayTypeInvoice        CheckoutPayType = "invoice"     // Pay with an invoice
	CheckoutPayTypeQRCodeScanning CheckoutPayType = "qr"          // Pay through QR code scanning
)

type PaymentStatus string

const (
	PaymentStatusError    PaymentStatus = "error"    // Failed payment. Data is incorrect
	PaymentStatusFailure  PaymentStatus = "failure"  // Failed payment
	PaymentStatusReversed PaymentStatus = "reversed" // Payment refunded
	PaymentStatusSuccess  PaymentStatus = "success"  // Successful payment
)

type Item struct {
	Amount string `json:"amount"` // Quantity/volume
	Cost   string `json:"cost"`   // The cost of all units of the specified product in the receipt (number of units * unit cost)
	ID     string `json:"id"`     // Item ID. You can get it in the Liqpay account - SCR - Kasa - Goods
	Price  string `json:"price"`  // Unit cost of goods
}

type RROInfo struct {
	Items          []Item   `json:"items,omitempty"`           // Data about products for which payment is performed
	DeliveryEmails []string `json:"delivery_emails,omitempty"` // List of e-mails to which receipts should be sent after fiscalization

}

type CheckoutRequest struct {
	Action      Action                    `json:"action"`                 // Transaction type
	Amount      string                    `json:"amount"`                 // Payment amount
	Currency    CheckoutCurrency          `json:"currency"`               // Payment currency
	Description string                    `json:"description"`            // Payment description
	OrderID     string                    `json:"order_id"`               // Unique purchase ID in your shop. Maximum length is 255 symbols
	RROInfo     *RROInfo                  `json:"rro_info,omitempty"`     // Data for fiscalization
	ExpiredDate *string                   `json:"expired_date,omitempty"` // Date and time until which customer is able to pay invoice by UTC. Should be sent in the following format 2016-04-24 00:00:00
	Language    *CheckoutCustomerLanguage `json:"language,omitempty"`     // Customer's language
	PayTypes    *CheckoutPayType          `json:"pay_types,omitempty"`    // Parameter that gets the methods of payments that displayed on checkout. If the parameter is not passed, shop settings will be applied, Checkout tab
	ResultURL   *string                   `json:"result_url,omitempty"`   // URL of your shop where the buyer would be redirected after completion of the purchase. Maximum length 510 symbols
	ServerURL   *string                   `json:"server_url,omitempty"`   // URL API in your store for notifications of payment status change (server -> server). Maximum length is 510 symbols
	VerifyCode  *string                   `json:"verifycode,omitempty"`   // Possible value Y. Dynamic verification code is generated and going back to Callback. Also generated code will be transferred to verification transactions for displaying in statement by client's card. Works for action = auth
}

type PaymentStatusRequest struct {
	Action  string `json:"action"`
	OrderID string `json:"order_id"`
}

type PaymentStatusResponse struct {
	AcqID              int           `json:"acq_id"`              // Acquirer ID
	Action             string        `json:"action"`              // Transaction type: pay, hold, paysplit, subscribe, paydonate, auth, regular
	AgentCommission    float64       `json:"agent_commission"`    // Agent commission in payment currency
	Amount             float64       `json:"amount"`              // Payment amount
	AmountBonus        float64       `json:"amount_bonus"`        // Payer bonus amount in payment currency debit
	AmountCredit       float64       `json:"amount_credit"`       // Payment amount for credit in currency of currency_credit
	AmountDebit        float64       `json:"amount_debit"`        // Payment amount for debit in currency of currency_debit
	AuthCodeCredit     string        `json:"authcode_credit"`     // Authorization code for transaction of credit
	AuthCodeDebit      string        `json:"authcode_debit"`      // Authorization code for transaction of debit
	BonusProcent       float64       `json:"bonus_procent"`       // Discount rate in percent
	BonusType          string        `json:"bonus_type"`          // Bonus type: bonusplus, discount_club, personal, promo
	CardToken          string        `json:"card_token"`          // Sender's card token
	CommissionCredit   float64       `json:"commission_credit"`   // Commission from the receiver in currency_credit
	CommissionDebit    float64       `json:"commission_debit"`    // Commission from the sender in currency_debit
	CreateDate         string        `json:"create_date"`         // Date of payment creation
	Currency           string        `json:"currency"`            // Payment currency
	CurrencyCredit     string        `json:"currency_credit"`     // Transaction currency of credit
	CurrencyDebit      string        `json:"currency_debit"`      // Transaction currency of debit
	Description        string        `json:"description"`         // Payment description
	EndDate            string        `json:"end_date"`            // Date of payment edition/end
	Info               string        `json:"info"`                // Additional payment information
	IP                 string        `json:"ip"`                  // Sender's IP address
	Is3DS              bool          `json:"is_3ds"`              // True if transaction passed with 3DS, false otherwise
	LiqpayOrderID      string        `json:"liqpay_order_id"`     // Payment order_id in LiqPay system
	MomentPart         string        `json:"moment_part"`         // Payment indication in parts
	MPIECI             int           `json:"mpi_eci"`             // MPI ECI: 5 - transaction passed with 3DS, 6 - issuer of payer card doesn't support 3d Secure, 7 - operation passed without 3d Secure
	OrderID            string        `json:"order_id"`            // Order_id payment
	PaymentID          int           `json:"payment_id"`          // Payment id in LiqPay system
	Paytype            string        `json:"paytype"`             // Method of payment: card, privat24, moment_part, cash, invoice, qr
	PublicKey          string        `json:"public_key"`          // Shop public key
	ReceiverCommission float64       `json:"receiver_commission"` // Receiver commission in payment currency
	RRNCredit          string        `json:"rrn_credit"`          // Unique transaction ID in authorization and settlement system of issuer bank for credit
	RRNDebit           string        `json:"rrn_debit"`           // Unique transaction ID in authorization and settlement system of issuer bank for debit
	SenderBonus        float64       `json:"sender_bonus"`        // Sender's bonus in the payment currency
	SenderCardBank     string        `json:"sender_card_bank"`    // Sender's card bank
	SenderCardCountry  string        `json:"sender_card_country"` // Sender's card country (digital ISO 3166-1 code)
	SenderCardMask2    string        `json:"sender_card_mask2"`   // Sender's card mask2
	SenderCardType     string        `json:"sender_card_type"`    // Sender's card type (MC/Visa)
	SenderCommission   float64       `json:"sender_commission"`   // Commission from the sender in the payment currency
	SenderPhone        string        `json:"sender_phone"`        // Sender's phone number
	Status             PaymentStatus `json:"status"`              // Payment status
}

type RefundRequest struct {
	Action  string `json:"action"`   // Transaction type
	Amount  string `json:"amount"`   // Payment amount. For example: 5, 7.34
	OrderID string `json:"order_id"` // Unique purchase ID in your shop. Maximum length is 255 symbols
}

type RefundResponse struct {
	Action    Action `json:"action"`     // Transaction type
	PaymentID int64  `json:"payment_id"` // Payment id in LiqPay system
	Status    string `json:"status"`     // Payment status
}
