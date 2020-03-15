package directpublisher

import (
	"github.com/gorilla/mux"
)

func (dp *DirectPublisher) Init(r *mux.Router)  {
	r.HandleFunc("/alexaproxy", dp.handleAlexaRequest).Methods("POST")
}
