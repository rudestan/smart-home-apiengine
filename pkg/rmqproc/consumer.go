package rmqproc

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/streadway/amqp"
	"log"
)

type RmqConfig struct {
	Host       string
	Port       int
	Login      string
	Password   string
	Exchange   string
	Queue      string
	RoutingKey string
}

type MessageHandler interface {
	handle(req string)
}

func displayError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Consume creates a new RMQ connection and starts listener to the preselect queue. Upon receiving
// a message handlerFunc will be executed with the received message payload.
func Consume(config RmqConfig, handler interface{}) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", config.Login, config.Password, config.Host, config.Port))

	displayError(err, "failed to connect")

	defer func() {
		err := conn.Close()

		if err != nil {
			panic(err)
		}
	}()

	ch, q := initRmq(conn, config)
	defer func() {
		err := ch.Close()

		if err != nil {
			panic(err)
		}
	}()

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	displayError(err, "failed to register a consumer")

	consumer := make(chan bool)

	go func() {
		for d := range msgs {
			log.Println("Message received")
			msgHandler := handler.(MessageHandler)
			msgHandler.handle(cast.ToString(d.Body))
		}
	}()

	log.Println("RMQ consumer started")

	<-consumer
}

func initRmq(conn *amqp.Connection, config RmqConfig) (*amqp.Channel, amqp.Queue) {
	ch, err := conn.Channel()

	displayError(err, "could not create a channel")

	err = ch.ExchangeDeclare(
		config.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	displayError(err, "can not declare the exchange")

	q, err := ch.QueueDeclare(
		config.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	displayError(err, "failed to declare the queue")

	err = ch.QueueBind(
		q.Name,
		config.RoutingKey,
		config.Exchange,
		false,
		nil,
	)
	displayError(err, "failed to bind the queue")

	return ch, q
}
