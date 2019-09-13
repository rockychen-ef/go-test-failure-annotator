package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_passing_a_gojunit_format_report_including_test_failures_returns_test_failure_details(test *testing.T) {
	svc := TestResultParseService{}

	result, err := svc.Parse("../fixture/test_report_gojunit_f.xml")

	assert.NoError(test, err)
	assert.Equal(test, 2, len(result))

	assert.Equal(test, "TestList", result[0].Name)
	assert.Equal(test, "handler/user_handler_test.go", result[0].File)
	assert.Equal(test, 53, result[0].Line)
	assert.Contains(test, result[0].Reason, "Not equal:")

	assert.Equal(test, "TestSave_Create", result[1].Name)
	assert.Equal(test, "repository/user_repo_test.go", result[1].File)
	assert.Equal(test, 81, result[1].Line)
	assert.Contains(test, result[1].Reason, "Not equal:")
}

func Test_passing_a_gojunit_format_report_without_test_failures_returns_a_zero_test_failure_array(test *testing.T) {
	svc := TestResultParseService{}

	result, err := svc.Parse("../fixture/test_report_gojunit_s.xml")

	assert.NoError(test, err)
	assert.Equal(test, 0, len(result))
}

func Test_passing_a_gotestsum_junit_format_report_including_test_failures_returns_test_failure_details(test *testing.T) {
	svc := TestResultParseService{}

	result, err := svc.Parse("../fixture/test_report_gotestsum_f.xml")

	assert.NoError(test, err)
	assert.Equal(test, 2, len(result))

	assert.Equal(test, "TestList", result[0].Name)
	assert.Equal(test, "handler/user_handler_test.go", result[0].File)
	assert.Equal(test, 53, result[0].Line)
	assert.Contains(test, result[0].Reason, "Not equal:")

	assert.Equal(test, "TestSave_Create", result[1].Name)
	assert.Equal(test, "repository/user_repo_test.go", result[1].File)
	assert.Equal(test, 81, result[1].Line)
	assert.Contains(test, result[1].Reason, "Not equal:")
}

func Test_passing_a_gotestsum_junit_format_report_without_test_failures_returns_a_zero_test_failure_array(test *testing.T) {
	svc := TestResultParseService{}

	result, err := svc.Parse("../fixture/test_report_gotestsum_s.xml")

	assert.NoError(test, err)
	assert.Equal(test, 0, len(result))
}
