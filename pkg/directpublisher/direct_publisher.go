package directpublisher

import (
	"smh-apiengine/pkg/amqp"

	"github.com/gorilla/mux"
)

type DirectPublisher struct {
	rmq *amqp.Rmq
	router *mux.Router
}

func NewDirectPublisher(rmqConfig *amqp.Config) *DirectPublisher  {
	return &DirectPublisher{
		rmq:amqp.NewRmq(rmqConfig),
		router:mux.NewRouter()}
}
