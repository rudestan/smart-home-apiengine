package directpublisher

import "github.com/gorilla/mux"

func (dp *DirectPublisher) Router() *mux.Router  {
	return dp.router
}

// Init implements routes webservers
func (dp *DirectPublisher) InitRoutes()  {
	dp.router.Use(dp.middleWare.HeadersMiddleware)
	dp.router.HandleFunc("/alexaproxy", dp.handleAlexaRequest).Methods("POST")
}
