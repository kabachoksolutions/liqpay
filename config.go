package liqpay

const (
	ServerServerURL   = "https://www.liqpay.ua/api/request"
	ClientServerURL   = "https://www.liqpay.ua/api/3/checkout"
	CurrentAPIVersion = "3"
)

// Config represents the configuration parameters required for interacting with the LiqPay API.
type Config struct {
	PrivateKey string // PrivateKey is the private key used for API authentication.
	PublicKey  string // PublicKey is the public key used for API authentication.
	Debug      bool   // Debug specifies whether debug mode is enabled.
}

// NewConfig creates a new Config instance with the provided public key, private key, and debug mode settings.
func NewConfig(publicKey, privateKey string, debugMode bool) *Config {
	return &Config{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Debug:      debugMode,
	}
}
