package alexakit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Slot struct represents amazon alexa slot element
type Slot struct {
	Name               string `json:"name"`
	Value              string `json:"value"`
	ConfirmationStatus string `json:"confirmationStatus"`
	Source             string `json:"source"`
}

// Intent struct represents alexa intent element and contains slots
type Intent struct {
	Name               string          `json:"name"`
	ConfirmationStatus string          `json:"confirmationStatus"`
	Slots              map[string]Slot `json:"slots"`
}

// Request struct is the root amazon alexa request element with all the data
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

// NewAlexaRequestIntent creates AlexaRequest struct with Intent from received http request
func NewAlexaRequestIntent(r *http.Request) (AlexaRequest, error) {
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

func (r *AlexaRequest) ToJson() (string, error) {
	content, err := json.Marshal(r)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
