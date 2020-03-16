package main

import (
	"log"
	"smh-apiengine/pkg/alexakit"
	"smh-apiengine/pkg/amqp"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleLambdaEvent(alexaRequest alexakit.AlexaRequest) (alexakit.AlexaResponse, error) {
	payload, err := alexaRequest.ToJson()

	if err != nil {
		log.Println("Failed to encode JSON")

		return alexakit.NewPlainTextSpeechResponse(alexakit.SpeechTextFailed), err
	}

	rmqConfig := alexakit.NewConfigFromEnv()
	rmq := amqp.NewRmq(rmqConfig)

	err = rmq.Publish(payload)

	if err != nil {
		return alexakit.NewPlainTextSpeechResponse(alexakit.SpeechTextFailed), err
	}

	speechResponse := alexakit.NewPlainTextSpeechResponse(alexakit.SpeechTextConfirmation)

	return speechResponse, nil
}


func main() {
	lambda.Start(HandleLambdaEvent)
}
