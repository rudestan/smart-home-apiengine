package directpublisher

import "github.com/gorilla/mux"

// Init implements routes webservers
func (dp *DirectPublisher) Init(r *mux.Router)  {
	r.Use(dp.middleWare.HeadersMiddleware)
	r.HandleFunc("/alexaproxy", dp.handleAlexaRequest).Methods("POST")
}
