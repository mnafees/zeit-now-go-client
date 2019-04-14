package now

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"
)

// Client struct
type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Config     ClientConfig
}

type requestLoginReqBody struct {
	email     string
	tokenName string
}

type requestLoginRespBody struct {
	token        string
	securityCode string
}

type verifyLoginRespBody struct {
	token string
}

var userAgent = fmt.Sprintf("NowGoClient/%s", runtime.Version())

// NewClient creates a new client to interact with the Zeit API
func NewClient(config ClientConfig) (*Client, error) {
	parsedURL, err := url.ParseRequestURI(config.Endpoint)
	if err != nil {
		return nil, err
	}
	return &Client{
		URL:        parsedURL,
		HTTPClient: &http.Client{},
		Config:     config,
	}, nil
}

func (c *Client) url(endpoint string) string {
	u := *c.URL
	u.Path = path.Join(c.URL.Path, endpoint)
	return u.String()
}

// FetchAuthToken gets a new auth token
func (c *Client) FetchAuthToken(authEmail string, authTokenName string) (string, error) {
	requestLoginBodyJSON, _ := json.Marshal(&requestLoginReqBody{
		email:     authEmail,
		tokenName: authTokenName,
	})
	// Request a login - https://zeit.co/docs/api#endpoints/authentication/request-a-login
	req, err := http.NewRequest(http.MethodPost, RequestLoginEndpoint, bytes.NewBuffer(requestLoginBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", nil
	}
	decoder := json.NewDecoder(resp.Body)
	var requestLoginResponse requestLoginRespBody
	decoder.Decode(&requestLoginResponse)
	// Verify a login - https://zeit.co/docs/api#endpoints/authentication/verify-login
	var authToken string
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf(VerifyLoginEndpoint, authEmail, requestLoginResponse.token), nil)
	req.Header.Set("User-Agent", userAgent)
	for authToken == "" {
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return "", nil
		}
		var verifyLoginResponse verifyLoginRespBody
		err = decoder.Decode(&requestLoginResponse)
		if err == nil {
			authToken = verifyLoginResponse.token
		} else {
			time.Sleep(time.Second * 1)
		}
	}
	return authToken, nil
}
