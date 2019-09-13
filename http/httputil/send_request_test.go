package httputil

import (
	"bytes"
	"elb2c/gh-action/http/httpconst"
	"elb2c/gh-action/testutil"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sending_a_request_with_an_URL_a_method_and_headers_returns_a_HTTP_status_code_headers_and_a_body(test *testing.T) {
	reqURL := "http://test.local/api/1"
	reqMethod := httpconst.MethodPost
	reqHeder := make(http.Header)
	reqHeder.Add("Content-Type", "application/x-www-form-urlencoded")
	reqHeder.Add("Accept", "application/json")

	respCode := 200
	respBody := []byte("response body")
	respHeader := make(http.Header)
	respHeader.Add("X-EF-EFID", "EFID token")
	respHeader.Add("X-EF-ACCESS", "Access token")
	respHeader.Add("X-EF-CORRELATION-ID", "Correlation ID")

	client := testutil.NewTestHTTPClient(func(req *http.Request) *http.Response {
		assert.Equal(test, reqURL, req.URL.String())
		assert.Equal(test, reqMethod, req.Method)
		assert.Equal(test, reqHeder, req.Header)

		return &http.Response{
			StatusCode: respCode,
			Header:     respHeader,
			Body:       ioutil.NopCloser(bytes.NewBufferString("response body")),
		}
	})

	result, err := SendRequest(client, reqURL, reqMethod, reqHeder, nil)

	assert.NoError(test, err)

	assert.Equal(test, respCode, result.StatusCode)
	assert.Equal(test, respHeader, result.Header)
	assert.Equal(test, respBody, result.Body)
}
