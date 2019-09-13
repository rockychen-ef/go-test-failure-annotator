package checkrun

import (
	"elb2c/gh-action/api"
	"elb2c/gh-action/http/httpconst"
	"elb2c/gh-action/http/httputil"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	nameOfCheckRun     = "Test failure annotator"
	checkRunDateFormat = "2006-01-02T15:04:05Z"
)

func makeHeaders(token string) (header http.Header) {
	header = make(http.Header)
	header.Add(httpconst.HeaderContentType, httpconst.MediaTypeApplicationJSON)
	header.Add(httpconst.HeaderAccept, httpconst.MediaTypeGitHubAntiopePreviewJSON)
	header.Add(httpconst.HeaderAuth, "Bearer "+token)
	header.Add(httpconst.HeaderUserAgent, nameOfCheckRun)

	return header
}

func submit(client *http.Client, URL string, method string,
	header http.Header, body io.Reader, successStatusCode int) (map[string]interface{}, error) {

	if client == nil {
		return nil, errors.New("HTTP client must not be nil")
	}

	if URL == "" {
		return nil, errors.New("URL must not be empty")
	}

	if method == "" {
		return nil, errors.New("Method must not be empty")
	}

	if header == nil {
		return nil, errors.New("Header must not be nil")
	}

	if body == nil {
		return nil, errors.New("Request body must not be nil")
	}

	resp, err := httputil.SendRequest(client, URL, method, header, body)
	if err != nil {
		return nil, fmt.Errorf("Failed to hit the API of check run. method: %s, url: %s, reason: %s", method, URL, err.Error())
	}

	log.Printf("HTTP response of check run. HTTP status code: %d, Response body: %s\n", resp.StatusCode, string(resp.Body))
	if resp.StatusCode != successStatusCode {
		return nil, &api.Error{
			StatusCode:   resp.StatusCode,
			ResponseBody: resp.Body,
			Message:      fmt.Sprintf("The API of check run doesn't return %d HTTP status code", successStatusCode),
		}
	}

	var respMap map[string]interface{}
	if err := json.Unmarshal(resp.Body, &respMap); err != nil {
		return nil, fmt.Errorf("Http status code: %d. Failed to unmarshal the response body of check run. method: %s, URL: %s, reason: %s",
			resp.StatusCode, method, URL, err.Error())
	}

	return respMap, nil
}
