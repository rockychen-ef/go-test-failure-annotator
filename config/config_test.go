package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setting_required_env_vars_returns_values(test *testing.T) {
	os.Setenv("TEST_RESULT", "/tmp/result.json")
	os.Setenv("GITHUB_API_URL", "https://api.url")
	os.Setenv("GITHUB_SHA", "sha")
	os.Setenv("GITHUB_TOKEN", "token")
	os.Setenv("GITHUB_WORKSPACE", "workspace")
	os.Setenv("GITHUB_REPOSITORY", "repository")

	result, err := Load()

	assert.NoError(test, err)
	assert.Equal(test, "workspace/tmp/result.json", result.TestResult())
	assert.Equal(test, "https://api.url", result.GitHub.URL)
	assert.Equal(test, "sha", result.GitHub.SHA)
	assert.Equal(test, "token", result.GitHub.Token)
	assert.Equal(test, "workspace", result.GitHub.Workspace)
	assert.Equal(test, "repository", result.GitHub.Repository)

	// clean up
	os.Unsetenv("TEST_RESULT")
	os.Unsetenv("GITHUB_API_URL")
	os.Unsetenv("GITHUB_SHA")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_WORKSPACE")
	os.Unsetenv("GITHUB_REPOSITORY")
}

func Test_missimg_required_env_variables_returns_an_errors(test *testing.T) {
	_, err := Load()
	assert.Error(test, err)
}

func Test_setting_valid_GitHub_repository_when_getting_owner_returns_owner(test *testing.T) {
	os.Setenv("TEST_RESULT", "/tmp/result.json")
	os.Setenv("GITHUB_API_URL", "https://api.url")
	os.Setenv("GITHUB_SHA", "sha")
	os.Setenv("GITHUB_TOKEN", "token")
	os.Setenv("GITHUB_WORKSPACE", "workspace")
	os.Setenv("GITHUB_REPOSITORY", "octocat/Hello-World")

	cfg, _ := Load()
	result, err := cfg.Owner()

	assert.NoError(test, err)
	assert.Equal(test, "octocat", result)

	os.Unsetenv("TEST_RESULT")
	os.Unsetenv("GITHUB_API_URL")
	os.Unsetenv("GITHUB_SHA")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_WORKSPACE")
	os.Unsetenv("GITHUB_REPOSITORY")
}

func Test_setting_unexpected_format_of_GitHub_repository_when_getting_owner_returns_an_error(test *testing.T) {
	os.Setenv("TEST_RESULT", "/tmp/result.json")
	os.Setenv("GITHUB_API_URL", "https://api.url")
	os.Setenv("GITHUB_SHA", "sha")
	os.Setenv("GITHUB_TOKEN", "token")
	os.Setenv("GITHUB_WORKSPACE", "workspace")
	os.Setenv("GITHUB_REPOSITORY", "Hello-World")

	cfg, _ := Load()
	_, err := cfg.Owner()

	assert.Error(test, err)
	assert.Equal(test, "Invalid environment variable 'GITHUB_REPOSITORY'. should be 'owner/repository' instead of 'Hello-World'", err.Error())

	os.Unsetenv("TEST_RESULT")
	os.Unsetenv("GITHUB_API_URL")
	os.Unsetenv("GITHUB_SHA")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_WORKSPACE")
	os.Unsetenv("GITHUB_REPOSITORY")
}

func Test_setting_valid_GitHub_repository_when_getting_repo_returns_repo(test *testing.T) {
	os.Setenv("TEST_RESULT", "/tmp/result.json")
	os.Setenv("GITHUB_API_URL", "https://api.url")
	os.Setenv("GITHUB_SHA", "sha")
	os.Setenv("GITHUB_TOKEN", "token")
	os.Setenv("GITHUB_WORKSPACE", "workspace")
	os.Setenv("GITHUB_REPOSITORY", "octocat/Hello-World")

	cfg, _ := Load()
	result, err := cfg.Repo()

	assert.NoError(test, err)
	assert.Equal(test, "Hello-World", result)

	os.Unsetenv("TEST_RESULT")
	os.Unsetenv("GITHUB_API_URL")
	os.Unsetenv("GITHUB_SHA")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_WORKSPACE")
	os.Unsetenv("GITHUB_REPOSITORY")
}

func Test_setting_unexpected_format_of_GitHub_repository_when_getting_repo_returns_an_error(test *testing.T) {
	os.Setenv("TEST_RESULT", "/tmp/result.json")
	os.Setenv("GITHUB_API_URL", "https://api.url")
	os.Setenv("GITHUB_SHA", "sha")
	os.Setenv("GITHUB_TOKEN", "token")
	os.Setenv("GITHUB_WORKSPACE", "workspace")
	os.Setenv("GITHUB_REPOSITORY", "a/b/c")

	cfg, _ := Load()
	_, err := cfg.Repo()

	assert.Error(test, err)
	assert.Equal(test, "Invalid environment variable 'GITHUB_REPOSITORY'. should be 'owner/repository' instead of 'a/b/c'", err.Error())

	os.Unsetenv("TEST_RESULT")
	os.Unsetenv("GITHUB_API_URL")
	os.Unsetenv("GITHUB_SHA")
	os.Unsetenv("GITHUB_TOKEN")
	os.Unsetenv("GITHUB_WORKSPACE")
	os.Unsetenv("GITHUB_REPOSITORY")
}
