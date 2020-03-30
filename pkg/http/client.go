package http

import "net/http"

// Client represents a http.client
type Client interface {
	Get(url string) (resp *http.Response, err error)
	Do(req *http.Request) (*http.Response, error)
}
