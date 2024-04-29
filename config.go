package liqpay

const (
	ServerServerURL   = "https://www.liqpay.ua/api/request"
	ClientServerURL   = "https://www.liqpay.ua/api/3/checkout"
	CurrentAPIVersion = "3"
)

type Config struct {
	PrivateKey string
	PublicKey  string
	Debug      bool
}

func NewConfig(publicKey, privateKey string, debugMode bool) *Config {
	return &Config{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Debug:      debugMode,
	}
}
