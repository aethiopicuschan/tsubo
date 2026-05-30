package tsubo

import "net/http"

// Client is a client for the tsubo API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

// Option is a function that configures a Client.
type Option func(*Client)

// NewClient creates a new Client with the given base URL and options.
func NewClient(baseURL string, options ...Option) (c *Client) {
	c = &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		userAgent:  "tsubo-client",
	}
	for _, option := range options {
		option(c)
	}
	return
}

// WithUserAgent returns an Option that sets the User-Agent header of the Client.
func WithUserAgent(userAgent string) func(*Client) {
	return func(c *Client) {
		c.SetUserAgent(userAgent)
	}
}

// WithHTTPClient returns an Option that sets the HTTP client of the Client.
func WithHTTPClient(httpClient *http.Client) func(*Client) {
	return func(c *Client) {
		c.SetHTTPClient(httpClient)
	}
}

// Getter for baseURL
func (c *Client) BaseURL() string {
	return c.baseURL
}

// Getter for httpClient
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// Getter for userAgent
func (c *Client) UserAgent() string {
	return c.userAgent
}

// Setter for baseURL
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// Setter for httpClient
func (c *Client) SetHTTPClient(httpClient *http.Client) {
	// If httpClient is nil, use the default http.Client
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	c.httpClient = httpClient
}

// Setter for userAgent
func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

// Do sends an HTTP request and returns an HTTP response, setting the User-Agent header.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", c.UserAgent())
	return c.HTTPClient().Do(req)
}
