package amqp

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

const defaultExp = "50000" // 50 sec.

func (proc *Rmq) Publish(payload string) error {
	conn, err := proc.connect()

	defer func() {
		err = conn.Close()
		if err != nil {
			log.Println("failed to close the connection")
		}
	}()

	return proc.publishToRmq(conn, payload)
}

func (proc *Rmq) publishToRmq(conn *amqp.Connection, jsonMessage string) error  {
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

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(jsonMessage),
		Expiration:   defaultExp,
	}

	err = ch.Publish(
		proc.config.Exchange,
		proc.config.RoutingKey,
		false,
		false,
		msg)

	if err != nil {
		return err
	}

	return nil
}
