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

func Test_CreationAPI_missing_HTTP_client_returns_an_error(test *testing.T) {
	api := CreationAPI{}

	_, err := api.Create()

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "HTTP client must not be nil")
}

func Test_CreationAPI_missing_base_URL_returns_an_error(test *testing.T) {
	api := CreationAPI{
		client: &http.Client{},
	}

	_, err := api.Create()

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "Base URL must not be empty")
}

func Test_CreationAPI_missing_config_returns_an_error(test *testing.T) {
	api := CreationAPI{
		client:  &http.Client{},
		baseURL: "http://test.local",
	}

	_, err := api.Create()

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "Config must not be nil")
}

func Test_CreationAPI_missing_GitHub_config_returns_an_error(test *testing.T) {
	api := CreationAPI{
		client:  &http.Client{},
		baseURL: "http://test.local",
		config:  &config.Config{},
	}

	_, err := api.Create()

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "GitHub config must not be empty")
}

func Test_CreationAPI_missing_GitHub_repository_returns_an_error(test *testing.T) {
	api := CreationAPI{
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

	_, err := api.Create()

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "GitHub repository must not be empty")
}

func Test_CreationAPI_missing_GitHub_token_returns_an_error(test *testing.T) {
	api := CreationAPI{
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

	_, err := api.Create()

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "GitHub token must not be empty")
}

func Test_CreationAPI_missing_GitHub_SHA_returns_an_error(test *testing.T) {
	api := CreationAPI{
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

	_, err := api.Create()

	assert.Error(test, err)
	assert.Equal(test, err.Error(), "GitHub SHA must not be empty")
}

func Test_CreationAPI_passing_valid_HTTP_client_base_URL_and_GitHub_configs_returns_an_ID(test *testing.T) {
	client := testutil.NewTestHTTPClient(func(req *http.Request) *http.Response {
		assert.Equal(test, "http://test.local/repos/octocat/Hello-World/check-runs", req.URL.String())
		assert.Equal(test, "POST", req.Method)
		// Verify headers
		assert.Equal(test, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(test, "application/vnd.github.antiope-preview+json", req.Header.Get("Accept"))
		assert.Equal(test, "Bearer token", req.Header.Get("Authorization"))
		assert.Equal(test, "Test failure annotator", req.Header.Get("User-Agent"))
		// Verify request body
		var reqBody CreationRequestBody
		json.Unmarshal(testutil.ToBytes(req.Body), &reqBody)
		assert.Equal(test, "Test failure annotator", reqBody.Name)
		assert.Equal(test, "sha", reqBody.SHA)
		assert.Equal(test, "in_progress", reqBody.Status)
		assert.NotEmpty(test, reqBody.StartedAt)

		return &http.Response{
			StatusCode: 201,
			Header:     make(http.Header),
			Body: ioutil.NopCloser(bytes.NewBufferString(`{
				"id": 4,
				"head_sha": "ce587453ced02b1526dfb4cb910479d431683101",
				"node_id": "MDg6Q2hlY2tSdW40",
				"external_id": "42",
				"url": "https://api.github.com/repos/github/hello-world/check-runs/4",
				"html_url": "http://github.com/github/hello-world/runs/4",
				"details_url": "https://example.com",
				"status": "in_progress",
				"conclusion": null,
				"started_at": "2018-05-04T01:14:52Z",
				"completed_at": null,
				"output": {
				  "title": "Mighty Readme Report",
				  "summary": "",
				  "text": ""
				},
				"name": "mighty_readme",
				"check_suite": {
				  "id": 5
				},
				"app": {
				  "id": 1,
				  "node_id": "MDExOkludGVncmF0aW9uMQ==",
				  "owner": {
					"login": "github",
					"id": 1,
					"node_id": "MDEyOk9yZ2FuaXphdGlvbjE=",
					"url": "https://api.github.com/orgs/github",
					"repos_url": "https://api.github.com/orgs/github/repos",
					"events_url": "https://api.github.com/orgs/github/events",
					"hooks_url": "https://api.github.com/orgs/github/hooks",
					"issues_url": "https://api.github.com/orgs/github/issues",
					"members_url": "https://api.github.com/orgs/github/members{/member}",
					"public_members_url": "https://api.github.com/orgs/github/public_members{/member}",
					"avatar_url": "https://github.com/images/error/octocat_happy.gif",
					"description": "A great organization"
				  },
				  "name": "Super CI",
				  "description": "",
				  "external_url": "https://example.com",
				  "html_url": "https://github.com/apps/super-ci",
				  "created_at": "2017-07-08T16:18:44-04:00",
				  "updated_at": "2017-07-08T16:18:44-04:00"
				},
				"pull_requests": [
				  {
					"url": "https://api.github.com/repos/github/hello-world/pulls/1",
					"id": 1934,
					"number": 3956,
					"head": {
					  "ref": "say-hello",
					  "sha": "3dca65fa3e8d4b3da3f3d056c59aee1c50f41390",
					  "repo": {
						"id": 526,
						"url": "https://api.github.com/repos/github/hello-world",
						"name": "hello-world"
					  }
					},
					"base": {
					  "ref": "master",
					  "sha": "e7fdf7640066d71ad16a86fbcbb9c6a10a18af4f",
					  "repo": {
						"id": 526,
						"url": "https://api.github.com/repos/github/hello-world",
						"name": "hello-world"
					  }
					}
				  }
				]
			  }`)),
		}
	})
	api := CreationAPI{
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

	result, err := api.Create()

	assert.NoError(test, err)
	assert.Equal(test, 4, result)
}

func Test_CreationAPI_assuming_remote_API_returning_non_201_status_code_returns_an_error(test *testing.T) {
	client := testutil.NewTestHTTPClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 500,
			Header:     make(http.Header),
			Body:       ioutil.NopCloser(bytes.NewBufferString("internal server error")),
		}
	})
	api := CreationAPI{
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

	_, err := api.Create()

	assert.Error(test, err)
}
