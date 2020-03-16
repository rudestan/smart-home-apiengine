package directpublisher

import (
	"smh-apiengine/pkg/webserver"

	"github.com/gorilla/mux"
)

func (dp *DirectPublisher) Router() *mux.Router  {
	return dp.router
}

// Init implements routes webservers
func (dp *DirectPublisher) InitRoutes()  {
	dp.router.Use(webserver.HeadersMiddleware)
	dp.router.HandleFunc("/alexaproxy", dp.handleAlexaRequest).Methods("POST")
}
