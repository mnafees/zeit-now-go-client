package now

// ClientConfig struct
type ClientConfig struct {
	Endpoint string
	Token    string
}

// Base consts
const (
	BaseEndpoint         = "https://api.zeit.co"
	RequestLoginEndpoint = BaseEndpoint + "/now/registration"
	VerifyLoginEndpoint  = RequestLoginEndpoint + "?email=%s&token=%s"
)

// NewConfig creates a new ClientConfig
func NewConfig() *ClientConfig {
	return &ClientConfig{
		Endpoint: BaseEndpoint,
		Token:    "",
	}
}
