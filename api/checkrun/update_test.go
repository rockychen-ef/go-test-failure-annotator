package checkrun

import (
	"bytes"
	"elb2c/gh-action/config"
	"elb2c/gh-action/testutil"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateAPI_passing_invalid_check_ID_returns_an_error(test *testing.T) {
	api := UpdateAPI{}
	annotations := make([]Annotation, 0)

	err := api.Update(0, annotations)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "Invalid check ID")
}

func Test_UpdateAPI_passing_nil_annotation_array_returns_an_error(test *testing.T) {
	api := UpdateAPI{}

	err := api.Update(1, nil)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "Annotation array must not be nil")
}

func Test_UpdateAPI_missing_HTTP_client_returns_an_error(test *testing.T) {
	api := UpdateAPI{}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "HTTP client must not be nil")
}

func Test_UpdateAPI_missing_base_URL_returns_an_error(test *testing.T) {
	api := UpdateAPI{
		client: &http.Client{},
	}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "Base URL must not be empty")
}

func Test_UpdateAPI_missing_config_returns_an_error(test *testing.T) {
	api := UpdateAPI{
		client:  &http.Client{},
		baseURL: "http://test.local",
	}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "Config must not be nil")
}

func Test_UpdateAPI_missing_GitHub_config_returns_an_error(test *testing.T) {
	api := UpdateAPI{
		client:  &http.Client{},
		baseURL: "http://test.local",
		config:  &config.Config{},
	}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "GitHub config must not be empty")
}

func Test_UpdateAPI_missing_GitHub_repository_returns_an_error(test *testing.T) {
	api := UpdateAPI{
		client:  &http.Client{},
		baseURL: "http://test.local",
		config: &config.Config{
			GitHub: config.GitHub{
				Repository: "",
				Token:      "token",
				SHA:        "sha",
			},
		},
	}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "GitHub repository must not be empty")
}

func Test_UpdateAPI_missing_GitHub_token_returns_an_error(test *testing.T) {
	api := UpdateAPI{
		client:  &http.Client{},
		baseURL: "http://test.local",
		config: &config.Config{
			GitHub: config.GitHub{
				Repository: "repository",
				Token:      "",
				SHA:        "sha",
			},
		},
	}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "GitHub token must not be empty")
}

func Test_UpdateAPI_missing_GitHub_SHA_returns_an_error(test *testing.T) {
	api := UpdateAPI{
		client:  &http.Client{},
		baseURL: "http://test.local",
		config: &config.Config{
			GitHub: config.GitHub{
				Repository: "repository",
				Token:      "token",
				SHA:        "",
			},
		},
	}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "GitHub SHA must not be empty")
}

func Test_UpdateAPI_passing_valid_HTTP_client_base_URL_and_GitHub_configs_returns_no_error(test *testing.T) {
	client := testutil.NewTestHTTPClient(func(req *http.Request) *http.Response {
		assert.Equal(test, "http://test.local/repos/octocat/Hello-World/check-runs/1", req.URL.String())
		assert.Equal(test, "PATCH", req.Method)
		// Verify headers
		assert.Equal(test, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(test, "application/vnd.github.antiope-preview+json", req.Header.Get("Accept"))
		assert.Equal(test, "Bearer token", req.Header.Get("Authorization"))
		assert.Equal(test, "Test failure annotator", req.Header.Get("User-Agent"))
		// Verify request body
		var reqBody UpdateRequestBody
		json.Unmarshal(testutil.ToBytes(req.Body), &reqBody)
		assert.Equal(test, "Test failure annotator", reqBody.Name)
		assert.Equal(test, "sha", reqBody.SHA)
		assert.Equal(test, "completed", reqBody.Status)
		assert.NotEmpty(test, reqBody.CompletedAt)
		assert.Equal(test, "failure", reqBody.Conclusion)
		assert.Equal(test, "Test failure details", reqBody.Output.Title)
		assert.Equal(test, "1 test failure(s) found", reqBody.Output.Summary)
		assert.Equal(test, "Test_passing_errors_returns_an_error", reqBody.Output.Annotations[0].Title)
		assert.Equal(test, "api/api_handler.go", reqBody.Output.Annotations[0].Path)
		assert.Equal(test, 1, reqBody.Output.Annotations[0].StartLine)
		assert.Equal(test, 3, reqBody.Output.Annotations[0].EndLine)
		assert.Equal(test, "failure", reqBody.Output.Annotations[0].Level)
		assert.Equal(test, "Because you passed errors", reqBody.Output.Annotations[0].Message)

		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		}
	})
	api := UpdateAPI{
		client:  client,
		baseURL: "http://test.local",
		config: &config.Config{
			GitHub: config.GitHub{
				Repository: "octocat/Hello-World",
				Token:      "token",
				SHA:        "sha",
			},
		},
	}
	annotations := make([]Annotation, 0)
	annotation := Annotation{
		Title:     "Test_passing_errors_returns_an_error",
		Path:      "api/api_handler.go",
		StartLine: 1,
		EndLine:   3,
		Level:     "failure",
		Message:   "Because you passed errors",
	}
	annotations = append(annotations, annotation)

	err := api.Update(1, annotations)

	assert.NoError(test, err)
}

func Test_UpdateAPI_assuming_remote_API_returning_non_200_status_code_returns_an_error(test *testing.T) {
	client := testutil.NewTestHTTPClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 500,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(bytes.NewBufferString("internal server error")),
		}
	})
	api := UpdateAPI{
		client:  client,
		baseURL: "http://test.local",
		config: &config.Config{
			GitHub: config.GitHub{
				Repository: "octocat/Hello-World",
				Token:      "token",
				SHA:        "sha",
			},
		},
	}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.Error(test, err)
}

func Test_UpdateAPI_passing_zero_annotation_returns_no_error(test *testing.T) {
	client := testutil.NewTestHTTPClient(func(req *http.Request) *http.Response {
		// Verify request body
		var reqBody UpdateRequestBody
		json.Unmarshal(testutil.ToBytes(req.Body), &reqBody)
		assert.Equal(test, "success", reqBody.Conclusion)
		assert.Equal(test, "0 test failure(s) found", reqBody.Output.Summary)

		return &http.Response{
			StatusCode: 200,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(bytes.NewBufferString("{}")),
		}
	})
	api := UpdateAPI{
		client:  client,
		baseURL: "http://test.local",
		config: &config.Config{
			GitHub: config.GitHub{
				Repository: "octocat/Hello-World",
				Token:      "token",
				SHA:        "sha",
			},
		},
	}
	annotations := make([]Annotation, 0)

	err := api.Update(1, annotations)

	assert.NoError(test, err)
}
