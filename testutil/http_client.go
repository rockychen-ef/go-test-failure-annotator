package testutil

import (
	"net/http"
)

// NewTestHTTPClient ...
func NewTestHTTPClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

// RoundTripFunc ...
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .,..
func (fn RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req), nil
}
