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
	VerifyLoginEndpoint  = RequestLoginEndpoint + "/verify?email=%s&token=%s"
)

// NewEmptyTokenConfig creates a new ClientConfig with an empty Token
func NewEmptyTokenConfig() *ClientConfig {
	return &ClientConfig{
		Endpoint: BaseEndpoint,
		Token:    "",
	}
}

// NewConfig creates a new ClientConfig with token as its Token
func NewConfig(token string) *ClientConfig {
	return &ClientConfig{
		Endpoint: BaseEndpoint,
		Token:    token,
	}
}
