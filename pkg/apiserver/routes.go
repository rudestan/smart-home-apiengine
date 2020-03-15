package apiserver

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (s *Server) Init(r *mux.Router)  {
	r.NotFoundHandler = http.HandlerFunc(s.handleNotFound)
	r.Use(HeadersMiddleware)
	r.Use(s.authTokenMiddleware)

	// Run routes
	r.HandleFunc("/run/command/{commandId}", s.handleRunCommand)
	r.HandleFunc("/run/scenario/{scenarioId}", s.handleRunScenario)
	r.HandleFunc("/run/intent", s.handleRunIntent).Methods("POST")

	// Api routes
	r.HandleFunc("/controls", s.handleControls)
}
