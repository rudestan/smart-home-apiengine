package directpublisher

import (
	"smh-apiengine/pkg/amqp"
	"smh-apiengine/pkg/webserver"

	"github.com/gorilla/mux"
)

type DirectPublisher struct {
	rmq *amqp.Rmq
	router *mux.Router
	middleWare *webserver.Middleware
}

func NewDirectPublisher(rmqConfig *amqp.Config, config *webserver.ServerConfig) *DirectPublisher  {
	return &DirectPublisher{
		rmq:amqp.NewRmq(rmqConfig),
		middleWare:webserver.NewMiddleware(config),
		router:mux.NewRouter()}
}
