package rmqproc

import (
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const rmq_dsn = "devstan.com:5672"
const rmq_login = "rmqadmin"
const rmq_password = "XhDfE35SlFD"
const rmq_exchange = "alexa_sync"
const rmq_routing_key = "alexa.response.json"

const progress_text = "Ok."

func OutPutAlexaProgressResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("{\"version\": \"1.0\",\"response\": {\"outputSpeech\": {\"type\": \"PlainText\",\"text\": \"%s\"},\"card\": {\"type\": \"Simple\",\"title\": \"The card\",\"content\": \"Visual card\"}}}", progress_text)))
}

func AlexaProxy(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusOK)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("- Recieved payload, pushing to the RMQ...")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", rmq_login, rmq_password, rmq_dsn))

	displayError(err,"Failed to connect!")
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("channel.open: %s", err)
	}
	defer ch.Close()

	publishToRmq(ch, fmt.Sprintf("%s", reqBody))

	OutPutAlexaProgressResponse(w)
}

func PublishPayloadToRmq(payload string) error  {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/", rmq_login, rmq_password, rmq_dsn))

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Println("failed to close the connection")
		}
	}()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	defer func() {
		err = ch.Close()
		if err != nil {
			log.Println("failed to close the channel")
		}
	}()

	err = publishToRmq(ch, payload)

	if err != nil {
		return err
	}

	return nil
}

func publishToRmq(ch *amqp.Channel, jsonMessage string) error  {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(jsonMessage),
		Expiration:   "50000",
	}

	err := ch.Publish(rmq_exchange, rmq_routing_key, false, false, msg)

	if err != nil {
		return err
	}
	fmt.Printf("- Payload pushed to the RMQ, please check %s exchange", rmq_exchange)

	return nil
}

func ListenHttpAndPublish()  {
	r := mux.NewRouter()

	r.HandleFunc("/alexaproxy", AlexaProxy).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8844",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("- Listening to the incoming POST requests...")

	fmt.Println(srv.ListenAndServe())
}
