package alexakit

import "encoding/json"

const (
	version = "1.0"
	OutputSpeechTypePlainText = "PlainText"
)

type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Response struct {
	OutputSpeech OutputSpeech `json:"outputSpeech"`
}

type AlexaResponse struct {
	Version string `json:"version"`
	Response Response `json:"response"`
}

func NewPlainTextSpeechResponse(speechText string) AlexaResponse {
	return AlexaResponse{
		Version:  version,
		Response: Response{
			OutputSpeech: OutputSpeech{
				Type: OutputSpeechTypePlainText,
				Text: speechText,
			},
		},
	}
}

func (r *AlexaResponse) ToJson() (string, error) {
	content, err := json.Marshal(r)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
