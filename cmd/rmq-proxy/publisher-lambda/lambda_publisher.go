package main

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"smh-apiengine/pkg/rmqproc"
)

type MyEvent struct {
	Name string `json:"What is your name?"`
	Age int     `json:"How old are you?"`
}

type Response struct {
	Message string `json:"message"`
}

func HandleLambdaEvent(event MyEvent) (Response, error) {
	payload, err := json.Marshal(event)

	if err != nil {
		return Response{Message: "Failed to encode JSON"}, err
	}

	err = rmqproc.PublishPayloadToRmq(string(payload))

	if err != nil {
		return Response{Message: "Failed to Publish payload!"}, err
	}

	return Response{Message: "Payload published, check RMQ"}, nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
