package liqpay

import "fmt"

type AntiFraudError string

const (
	AntiFraudLimitExceeded       AntiFraudError = "limit"   // Exceeded limit on amount or number of client payments
	AntiFraudFraudDetected       AntiFraudError = "frod"    // Transaction identified as atypical/risky according to Bank's Anti-Fraud rules
	AntiFraudDeclinedTransaction AntiFraudError = "decline" // Transaction identified as atypical/risky according to Bank's Anti-Fraud system
)

type NonFinancialError string

const (
	NonFinancialAuthorizationRequired         NonFinancialError = "err_auth"                    // Authorization required
	NonFinancialCacheTimeElapsed              NonFinancialError = "err_cache"                   // Data storage time elapsed for this operation
	NonFinancialUserNotFound                  NonFinancialError = "user_not_found"              // User not found
	NonFinancialSMSSendFailed                 NonFinancialError = "err_sms_send"                // Failed to send SMS
	NonFinancialSMSOTPIncorrect               NonFinancialError = "err_sms_otp"                 // SMS password entered incorrectly
	NonFinancialShopBlocked                   NonFinancialError = "shop_blocked"                // Shop is blocked
	NonFinancialShopNotActive                 NonFinancialError = "shop_not_active"             // Shop is not active
	NonFinancialInvalidSignature              NonFinancialError = "invalid_signature"           // Invalid request signature
	NonFinancialOrderIDEmpty                  NonFinancialError = "order_id_empty"              // Empty order_id passed
	NonFinancialShopNotAgent                  NonFinancialError = "err_shop_not_agent"          // You are not an agent for the specified shop
	NonFinancialCardNotFound                  NonFinancialError = "err_card_def_notfound"       // Card for receiving payments not found in wallet
	NonFinancialNoCardToken                   NonFinancialError = "err_no_card_token"           // User has no card with such card_token
	NonFinancialCardLiqpayDefault             NonFinancialError = "err_card_liqpay_def"         // Specify another card
	NonFinancialInvalidCardType               NonFinancialError = "err_card_type"               // Invalid card type
	NonFinancialInvalidCardCountry            NonFinancialError = "err_card_country"            // Specify another card
	NonFinancialAmountBelowLimit              NonFinancialError = "err_limit_amount"            // Transfer amount less than or greater than specified limit
	NonFinancialPaymentAmountLimit            NonFinancialError = "err_payment_amount_limit"    // Transfer amount less than or greater than specified limit
	NonFinancialAmountLimitExceeded           NonFinancialError = "amount_limit"                // Exceeded amount limit
	NonFinancialPaymentSenderCard             NonFinancialError = "payment_err_sender_card"     // Specify another sender card
	NonFinancialPaymentProcessing             NonFinancialError = "payment_processing"          // Payment is being processed
	NonFinancialPaymentDiscountNotFound       NonFinancialError = "err_payment_discount"        // Discount for this payment not found
	NonFinancialWalletLoadFailed              NonFinancialError = "err_wallet"                  // Failed to load wallet
	NonFinancialVerifyCodeRequired            NonFinancialError = "err_get_verify_code"         // Card verification required
	NonFinancialIncorrectVerifyCode           NonFinancialError = "err_verify_code"             // Incorrect verification code
	NonFinancialAdditionalInfoRequired        NonFinancialError = "wait_info"                   // Additional information is expected, try again later
	NonFinancialInvalidRequestPath            NonFinancialError = "err_path"                    // Invalid request address
	NonFinancialCashPaymentAcquirerNotAllowed NonFinancialError = "err_payment_cash_acq"        // Payment cannot be made in this shop
	NonFinancialSplitAmountMismatch           NonFinancialError = "err_split_amount"            // Split payment amounts do not match payment amount
	NonFinancialReceiverCardNotSet            NonFinancialError = "err_card_receiver_def"       // Recipient has not set up card to receive payments
	NonFinancialPaymentStatusError            NonFinancialError = "payment_err_status"          // Incorrect payment status
	NonFinancialPublicKeyNotFound             NonFinancialError = "public_key_not_found"        // Public key not found
	NonFinancialPaymentNotFound               NonFinancialError = "payment_not_found"           // Payment not found
	NonFinancialPaymentNotSubscribed          NonFinancialError = "payment_not_subscribed"      // Payment is not regular
	NonFinancialWrongAmountCurrency           NonFinancialError = "wrong_amount_currency"       // Payment currency does not match debit currency
	NonFinancialAmountHoldError               NonFinancialError = "err_amount_hold"             // Amount cannot exceed payment amount
	NonFinancialAccessError                   NonFinancialError = "err_access"                  // Access error
	NonFinancialDuplicateOrderID              NonFinancialError = "order_id_duplicate"          // Such order_id already exists
	NonFinancialAccountBlocked                NonFinancialError = "err_blocked"                 // Account access closed
	NonFinancialParameterEmpty                NonFinancialError = "err_empty"                   // Parameter not filled
	NonFinancialPhoneParameterEmpty           NonFinancialError = "err_empty_phone"             // Phone parameter not filled
	NonFinancialParameterMissing              NonFinancialError = "err_missing"                 // Parameter not passed
	NonFinancialParameterIncorrect            NonFinancialError = "err_wrong"                   // Parameter specified incorrectly
	NonFinancialIncorrectCurrency             NonFinancialError = "err_wrong_currency"          // Incorrect currency specified. Use: USD, UAH, EUR
	NonFinancialInvalidPhoneNumber            NonFinancialError = "err_phone"                   // Incorrect phone number entered
	NonFinancialInvalidCardNumber             NonFinancialError = "err_card"                    // Incorrect card number specified
	NonFinancialCardBINNotFound               NonFinancialError = "err_card_bin"                // Card BIN not found
	NonFinancialTerminalNotFound              NonFinancialError = "err_terminal_notfound"       // Terminal not found
	NonFinancialCommissionNotFound            NonFinancialError = "err_commission_notfound"     // Commission not found
	NonFinancialPaymentCreationFailed         NonFinancialError = "err_payment_create"          // Failed to create payment
	NonFinancialMPIVerificationFailed         NonFinancialError = "err_mpi"                     // Failed to verify card
	NonFinancialCurrencyNotAllowed            NonFinancialError = "err_currency_is_not_allowed" // Currency not allowed
	NonFinancialOperationIncomplete           NonFinancialError = "err_look"                    // Operation incomplete
	NonFinancialModsEmpty                     NonFinancialError = "err_mods_empty"              // Operation incomplete
	NonFinancialPaymentTypeError              NonFinancialError = "payment_err_type"            // Incorrect payment type
	NonFinancialPaymentCurrencyError          NonFinancialError = "err_payment_currency"        // Card or transfer currency not allowed
	NonFinancialExchangeRateNotFound          NonFinancialError = "err_payment_exchangerates"   // Failed to find corresponding exchange rate
	NonFinancialInvalidRequestSignature       NonFinancialError = "err_signature"               // Invalid request signature
	NonFinancialAPIActionParameterMissing     NonFinancialError = "err_api_action"              // Action parameter not passed
	NonFinancialAPICallbackParameterMissing   NonFinancialError = "err_api_callback"            // Callback parameter not passed
	NonFinancialAPIIPForbidden                NonFinancialError = "err_api_ip"                  // API call from this IP address is forbidden in this merchant
	NonFinancialPhoneConfirmationExpired      NonFinancialError = "expired_phone"               // Payment confirmation deadline by entering phone number expired
	NonFinancialThreeDSecureExpired           NonFinancialError = "expired_3ds"                 // 3DS client verification deadline expired
	NonFinancialOTPConfirmationExpired        NonFinancialError = "expired_otp"                 // Payment confirmation deadline by OTP password expired
	NonFinancialCVVConfirmationExpired        NonFinancialError = "expired_cvv"                 // Payment confirmation deadline by entering CVV code expired
	NonFinancialPrivat24Expired               NonFinancialError = "expired_p24"                 // Privat24 card selection deadline expired
	NonFinancialSenderDataExpired             NonFinancialError = "expired_sender"              // Sender data collection deadline expired
	NonFinancialPINConfirmationExpired        NonFinancialError = "expired_pin"                 // Payment confirmation deadline by card PIN expired
	NonFinancialIVRConfirmationExpired        NonFinancialError = "expired_ivr"                 // Payment confirmation deadline by IVR call expired
	NonFinancialCaptchaConfirmationExpired    NonFinancialError = "expired_captcha"             // Payment confirmation deadline by captcha expired
	NonFinancialPasswordConfirmationExpired   NonFinancialError = "expired_password"            // Payment confirmation deadline by Privat24 password expired
	NonFinancialSenderAppConfirmationExpired  NonFinancialError = "expired_senderapp"           // Payment confirmation deadline by Privat24 form expired
	NonFinancialPreparedTransactionExpired    NonFinancialError = "expired_prepared"            // Deadline for completing created payment expired
	NonFinancialMasterPassExpired             NonFinancialError = "expired_mp"                  // Deadline for completing payment in MasterPass wallet expired
	NonFinancialQRCodeExpired                 NonFinancialError = "expired_qr"                  // Deadline for confirming payment by scanning QR code expired
	NonFinancialCardNot3DSupported            NonFinancialError = "5"                           // Card does not support 3DSecure
)

type FinancialError int

const (
	FinancialGeneralError                  FinancialError = 90   // General error during processing
	FinancialInvalidTokenMerchant          FinancialError = 101  // Token created not by this merchant
	FinancialInactiveToken                 FinancialError = 102  // Sent token is not active
	FinancialMaxPurchaseAmountExceeded     FinancialError = 103  // Maximum purchase amount reached for the token
	FinancialTransactionLimitExceeded      FinancialError = 104  // Transaction limit for the token exceeded
	FinancialUnsupportedCard               FinancialError = 105  // Card is not supported
	FinancialPreauthorizationNotAllowed    FinancialError = 106  // Merchant not allowed to preauthorize
	FinancialAcquirerDoesNotSupport3DS     FinancialError = 107  // Such token does not exist
	FinancialTokenNotFound                 FinancialError = 108  // Such token does not exist
	FinancialTokenDoesNotExist             FinancialError = 109  // IP attempt limit exceeded
	FinancialIPAttemptsLimitExceeded       FinancialError = 110  // Session expired
	FinancialSessionExpired                FinancialError = 111  // Card branch blocked
	FinancialCardBranchBlocked             FinancialError = 112  // Daily card branch limit reached
	FinancialDailyCardBranchLimitExceeded  FinancialError = 113  // Temporary restriction on P2P payments from PB cards to cards of foreign banks
	FinancialP2PBlocked                    FinancialError = 113  // Daily limit for using the card reached
	FinancialDailyTransactionLimitExceeded FinancialError = 2903 // Such order_id already exists
	FinancialDuplicateOrderID              FinancialError = 2915 // Payments to this country are forbidden
	FinancialPaymentCountryForbidden       FinancialError = 3914 // Card expiration date expired
	FinancialCardExpirationExpired         FinancialError = 9851 // Incorrect card number
	FinancialInvalidCardNumber             FinancialError = 9852 // Payment declined. Try again later
	FinancialPaymentDeclined               FinancialError = 9854 // Card does not support this type of transaction
	FinancialUnsupportedTransactionType    FinancialError = 9855 // Card does not support this type of transaction
)

// APIError represents an error returned by the LiqPay API.
type APIError struct {
	Status string `json:"status"`
	Code   string `json:"err_code"`
	Desc   string `json:"err_description"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("status: %s, code: %s, description: %s", e.Status, e.Code, e.Desc)
}

// ConvertToAPIError converts an error to *APIError type if possible.
func ConvertToAPIError(err error) (*APIError, error) {
	apiErr, ok := err.(*APIError)
	if !ok {
		return nil, fmt.Errorf("failed to convert error to APIError: %w", err)
	}
	return apiErr, nil
}

// ErrorRefersToAPI checks if the error refers to an APIError.
func ErrorRefersToAPI(err error) bool {
	if _, ok := err.(*APIError); ok {
		return true
	}
	return false
}
