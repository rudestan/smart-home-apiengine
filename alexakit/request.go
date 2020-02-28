package alexakit

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
)

type Slot struct {
    Name               string `json:"name"`
    Value              string `json:"value"`
    ConfirmationStatus string `json:"confirmationStatus"`
    Source             string `json:"source"`
}

type Intent struct {
    Name               string          `json:"name"`
    ConfirmationStatus string          `json:"confirmationStatus"`
    Slots              map[string]Slot `json:"slots"`
}

type Request struct {
    Type      string `json:"type"`
    RequestID string `json:"requestId"`
    TimeStamp string `json:"timestamp"`
    Locale    string `json:"locale"`
    Intent    Intent `json:"intent"`
}

// AlexaRequest struct represent the json structure of the request from alexa api
type AlexaRequest struct {
    Version string  `json:"version"`
    Request Request `json:"request"`
}

// SimpleSlot struct that represent simplified version of the alexa slot data
type SimpleSlot struct {
    Name  string
    Value string
}

// SimpleIntent struct that represents simplified version of the alexa intent with slots map
type SimpleIntent struct {
    Name  string
    Slots map[string]SimpleSlot
}

func (ar *AlexaRequest) GetRequestType() string  {
    return ar.Request.Type
}

func NewAlexaRequestIntent(r *http.Request) (AlexaRequest, error)  {
    var alexaRequestIntent AlexaRequest

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        return AlexaRequest{}, err
    }

    err = json.Unmarshal(body, &alexaRequestIntent)
    if err != nil {
        return AlexaRequest{}, err
    }

    return alexaRequestIntent, nil
}
