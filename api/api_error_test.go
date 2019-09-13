package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calling_the_Error_method_returns_an_error_message(test *testing.T) {
	apiErr := Error{
		Message: "error!",
	}

	result := apiErr.Error()

	assert.Equal(test, apiErr.Message, result)
}

func Test_calling_the_ToJSONMap_method_with_response_body_returns_JSON_maps(test *testing.T) {
	apiErr := Error{
		ResponseBody: []byte(`{
			"error_description": "Invalid access token.",
			"error": "invalid_token"
		  }`),
	}

	result, err := apiErr.ToJSONMap()

	assert.NoError(test, err)
	assert.Equal(test, "Invalid access token.", result["error_description"].(string))
	assert.Equal(test, "invalid_token", result["error"].(string))
}

func Test_calling_the_ToJSONMap_method_when_response_body_is_nil_returns_an_error(test *testing.T) {
	apiErr := Error{}

	result, err := apiErr.ToJSONMap()

	assert.Nil(test, result)
	assert.Equal(test, "ResponseBody must not be nil", err.Error())
}
