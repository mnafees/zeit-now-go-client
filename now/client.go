package now

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"runtime"
)

// Client struct
type Client struct {
	URL        *url.URL
	HTTPClient *http.Client
	Config     ClientConfig
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
	reqBody := map[string]string{
		"email":     authEmail,
		"tokenName": authTokenName,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(reqBody)
	// Request a login - https://zeit.co/docs/api#endpoints/authentication/request-a-login
	req, err := http.NewRequest(http.MethodPost, RequestLoginEndpoint, buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", userAgent)
	log.Println("Requesting a login ...")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(resp.Body)
	var requestLoginResponse map[string]string
	err = decoder.Decode(&requestLoginResponse)
	if err != nil {
		return "", err
	}
	// Verify a login - https://zeit.co/docs/api#endpoints/authentication/verify-login
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf(VerifyLoginEndpoint, authEmail, requestLoginResponse["token"]), nil)
	req.Header.Set("User-Agent", userAgent)
	log.Printf("Verifying login with token: %s\n", requestLoginResponse["token"])
	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	for resp.StatusCode != 200 {
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return "", err
		}
	}
	var verifyLoginResponse map[string]string
	decoder = json.NewDecoder(resp.Body)
	decoder.Decode(&verifyLoginResponse)
	return verifyLoginResponse["token"], nil
}
