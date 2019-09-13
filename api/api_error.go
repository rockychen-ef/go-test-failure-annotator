package api

import (
	"encoding/json"
	"errors"
)

// Error an error object used when calling 3rd party APIs failure
type Error struct {
	StatusCode   int
	ResponseBody []byte
	Message      string
}

func (self *Error) Error() string {
	return self.Message
}

// ToJSONMap convert ResponseBody to JSON maps
func (self *Error) ToJSONMap() (map[string]interface{}, error) {
	if self.ResponseBody == nil {
		return nil, errors.New("ResponseBody must not be nil")
	}

	var result map[string]interface{}
	json.Unmarshal(self.ResponseBody, &result)

	return result, nil
}
