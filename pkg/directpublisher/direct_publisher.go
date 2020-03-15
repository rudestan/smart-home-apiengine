package directpublisher

import (
	"smh-apiengine/pkg/amqp"
)

type DirectPublisher struct {
	rmq *amqp.Rmq
}

func NewDirectPublisher(rmqConfig *amqp.Config) *DirectPublisher  {
	rmq := amqp.NewRmq(rmqConfig)

	return &DirectPublisher{rmq:rmq}
}
