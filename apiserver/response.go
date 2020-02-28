package apiserver

import "encoding/json"

type ApiResultResponse struct {
    Result string `json:"result"`
    Message string `json:"message"`
}

func NewResponse(result string, msg string) string  {
    resp := ApiResultResponse{Result:result, Message:msg}
    jsonResp, err := json.Marshal(resp)

    if err != nil {
        return "{\"result\":\"error\",\"message\":\"internal error\"}"
    }

    return string(jsonResp)
}
