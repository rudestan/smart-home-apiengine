package amqp

import (
	"fmt"
	"github.com/spf13/cast"
	"github.com/streadway/amqp"
	"log"
	"reflect"
)

// Consume creates a new RMQ connection and starts listener to the preselect queue. Upon receiving
// a message handlerFunc will be executed with the received message payload.
func (proc* Rmq) Consume(handler interface{}) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", proc.config.Login, proc.config.Password, proc.config.Host, proc.config.Port))

	displayError(err, "failed to connect")

	defer func() {
		err := conn.Close()

		if err != nil {
			panic(err)
		}
	}()

	ch, q := proc.openChannelAndQueue(conn)
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

			handlerValue := reflect.ValueOf(handler)

			switch handlerValue.Elem().Interface().(type) {
			case MessageHandler:
				handler.(MessageHandler).handle(cast.ToString(d.Body))
			default:
				log.Println("Wrong handler type provided!")
				return
			}
		}
	}()

	log.Println("RMQ consumer started")

	<-consumer
}

func (proc *Rmq) openChannelAndQueue(conn *amqp.Connection) (*amqp.Channel, amqp.Queue) {
	ch, err := conn.Channel()

	displayError(err, "could not create a channel")

	err = ch.ExchangeDeclare(
		proc.config.Exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	displayError(err, "can not declare the exchange")

	q, err := ch.QueueDeclare(
		proc.config.Queue,
		true,
		false,
		false,
		false,
		nil,
	)
	displayError(err, "failed to declare the queue")

	err = ch.QueueBind(
		q.Name,
		proc.config.RoutingKey,
		proc.config.Exchange,
		false,
		nil,
	)
	displayError(err, "failed to bind the queue")

	return ch, q
}
