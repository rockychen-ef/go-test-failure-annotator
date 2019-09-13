package httputil

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func SendRequest(client *http.Client, URL string, method string, header http.Header, reqBody io.Reader) (*Response, error) {
	if client == nil {
		return nil, errors.New("HTTP client must not be nil")
	}

	if URL == "" {
		return nil, errors.New("Target URL must not be empty")
	}

	if method == "" {
		return nil, errors.New("HTTP method must not be empty")
	}

	if header == nil {
		return nil, errors.New("HTTP request header must not be nil")
	}

	req, err := http.NewRequest(method, URL, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header = header

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Process response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       respBody,
	}, nil
}
