package apiserver

import (
    "encoding/json"
    "log"
)

const (
    responseSuccess = "success"
    responseError = "error"
)

type ApiResultResponseWithPayload struct {
    Result string `json:"result"`
    Message string `json:"message"`
    Payload interface{} `json:"payload"`
}

type ApiResultResponse struct {
    Result string `json:"result"`
    Message string `json:"message"`
}

func NewResponse(result string, msg string, payload interface{}) string  {
    var resp interface{}

    if payload != nil {
        resp = ApiResultResponseWithPayload{Result:result, Message:msg, Payload:payload}
    } else {
        resp = ApiResultResponse{Result:result, Message:msg}
    }

    jsonResp, err := json.Marshal(resp)

    if err != nil {
        log.Println("failed to build the response")

        return "{\"result\":\"error\",\"message\":\"internal error\"}"
    }

    return string(jsonResp)
}

func NewSuccessResponse(msg string, payload interface{}) string  {
    return NewResponse(responseSuccess, msg, payload)
}

func NewErrorResponse(msg string) string  {
    return NewResponse(responseError, msg, nil)
}
