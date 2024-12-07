package errors

import "encoding/json"

type CustomError struct {
	Message map[string]interface{}
}

func (e *CustomError) Error() string {
	jsonMessage, err := json.Marshal(e.Message)
	if err != nil {
		return "failed to marshal error message"
	}
	return string(jsonMessage)
}

func NewCustomError(message map[string]interface{}) error {
	return &CustomError{Message: message}
}
