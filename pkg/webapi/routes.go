package webapi

import (
	"net/http"
	"smh-apiengine/pkg/devicecontrol"
	s "smh-apiengine/pkg/webserver"
	"github.com/gorilla/mux"
)

type ApiRouteHandlers struct {
	dataProvider *devicecontrol.DeviceControl
	middleWare *s.Middleware
}

func NewApiRouteHandlers(config *s.ServerConfig, deviceControl *devicecontrol.DeviceControl) ApiRouteHandlers  {
	middleware := s.NewMiddleware(config)
	return ApiRouteHandlers{dataProvider: deviceControl, middleWare:middleware}
}

func (apiHandlers *ApiRouteHandlers) Init(r *mux.Router)  {
	r.NotFoundHandler = http.HandlerFunc(apiHandlers.HandleNotFound)
	r.Use(apiHandlers.middleWare.HeadersMiddleware)
	r.Use(apiHandlers.middleWare.AuthTokenMiddleware)

	// Run routes
	r.HandleFunc("/run/command/{commandId}", apiHandlers.handleRunCommand)
	r.HandleFunc("/run/scenario/{scenarioId}", apiHandlers.handleRunScenario)
	r.HandleFunc("/run/intent", apiHandlers.handleRunIntent).Methods("POST")

	// Api routes
	r.HandleFunc("/controls", apiHandlers.handleControls)
}
