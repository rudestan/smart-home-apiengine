package directpublisher

import (
	"smh-apiengine/pkg/amqp"
	"smh-apiengine/pkg/webserver"
)

type DirectPublisher struct {
	rmq *amqp.Rmq
	middleWare *webserver.Middleware
}

func NewDirectPublisher(rmqConfig *amqp.Config, config *webserver.ServerConfig) *DirectPublisher  {
	rmq := amqp.NewRmq(rmqConfig)
	middleWare := webserver.NewMiddleware(config)

	return &DirectPublisher{rmq:rmq, middleWare:middleWare}
}
