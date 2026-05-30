package tsubo_test

import (
	"net/http"
	"testing"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		baseURL   string
		options   []tsubo.Option
		assertion func(t *testing.T, c *tsubo.Client)
	}{
		{
			name:    "default client",
			baseURL: "https://example.com",
			assertion: func(t *testing.T, c *tsubo.Client) {
				t.Helper()

				assert.Equal(t, "https://example.com", c.BaseURL())
				assert.NotNil(t, c.HTTPClient())
				assert.Equal(t, "tsubo-client", c.UserAgent())
			},
		},
		{
			name:    "with user agent",
			baseURL: "https://example.com",
			options: []tsubo.Option{
				tsubo.WithUserAgent("custom-agent"),
			},
			assertion: func(t *testing.T, c *tsubo.Client) {
				t.Helper()

				assert.Equal(t, "custom-agent", c.UserAgent())
			},
		},
		{
			name:    "with http client",
			baseURL: "https://example.com",
			options: []tsubo.Option{
				tsubo.WithHTTPClient(&http.Client{}),
			},
			assertion: func(t *testing.T, c *tsubo.Client) {
				t.Helper()

				assert.NotNil(t, c.HTTPClient())
			},
		},
		{
			name:    "with nil http client",
			baseURL: "https://example.com",
			options: []tsubo.Option{
				tsubo.WithHTTPClient(nil),
			},
			assertion: func(t *testing.T, c *tsubo.Client) {
				t.Helper()

				assert.NotNil(t, c.HTTPClient())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := tsubo.NewClient(tt.baseURL, tt.options...)

			require.NotNil(t, c)

			tt.assertion(t, c)
		})
	}
}

func TestClientSettersAndGetters(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		run  func(t *testing.T)
	}{
		{
			name: "base url",
			run: func(t *testing.T) {
				t.Helper()

				c := tsubo.NewClient("https://example.com")

				c.SetBaseURL("https://example.org")

				assert.Equal(t, "https://example.org", c.BaseURL())
			},
		},
		{
			name: "http client",
			run: func(t *testing.T) {
				t.Helper()

				c := tsubo.NewClient("https://example.com")
				httpClient := &http.Client{}

				c.SetHTTPClient(httpClient)

				assert.Same(t, httpClient, c.HTTPClient())
			},
		},
		{
			name: "user agent",
			run: func(t *testing.T) {
				t.Helper()

				c := tsubo.NewClient("https://example.com")

				c.SetUserAgent("custom-agent")

				assert.Equal(t, "custom-agent", c.UserAgent())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.run(t)
		})
	}
}

func TestClientDo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		userAgent string
	}{
		{
			name:      "default user agent",
			userAgent: "tsubo-client",
		},
		{
			name:      "custom user agent",
			userAgent: "custom-agent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var gotUserAgent string

			httpClient := &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					gotUserAgent = req.Header.Get("User-Agent")

					return &http.Response{
						StatusCode: http.StatusOK,
						Header:     make(http.Header),
						Body:       http.NoBody,
						Request:    req,
					}, nil
				}),
			}

			c := tsubo.NewClient(
				"https://example.com",
				tsubo.WithHTTPClient(httpClient),
				tsubo.WithUserAgent(tt.userAgent),
			)

			req, err := http.NewRequest(
				http.MethodGet,
				"https://example.com",
				nil,
			)
			require.NoError(t, err)

			resp, err := c.Do(req)
			require.NoError(t, err)
			require.NotNil(t, resp)

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Equal(t, tt.userAgent, gotUserAgent)
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
