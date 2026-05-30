package tsubo

import "net/http"

// HTTPClient is an interface that represents the ability to perform HTTP requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
