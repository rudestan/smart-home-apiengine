package apiserver

import (
	"encoding/json"
	"log"
)

const (
	responseSuccess = "success"
	responseError   = "error"
)

// APIResultResponseWithPayload api response with some additional payload
type APIResultResponseWithPayload struct {
	Result  string      `json:"result"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

// APIResultResponse api response for messages without payload
type APIResultResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

// newResponse creates new response and returns marshalled json string
func newResponse(result string, msg string, payload interface{}) string {
	var resp interface{}

	if payload != nil {
		resp = APIResultResponseWithPayload{Result: result, Message: msg, Payload: payload}
	} else {
		resp = APIResultResponse{Result: result, Message: msg}
	}

	jsonResp, err := json.Marshal(resp)

	if err != nil {
		log.Println("failed to build the response")

		return "{\"result\":\"error\",\"message\":\"internal error\"}"
	}

	return string(jsonResp)
}

// newSuccessResponse creates success response
func newSuccessResponse(msg string, payload interface{}) string {
	return newResponse(responseSuccess, msg, payload)
}

// newErrorResponse creates an error response
func newErrorResponse(msg string) string {
	return newResponse(responseError, msg, nil)
}
