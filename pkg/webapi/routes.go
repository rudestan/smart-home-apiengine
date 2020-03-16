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
	router *mux.Router
}

func NewApiRouteHandlers(config *s.ServerConfig, deviceControl *devicecontrol.DeviceControl) *ApiRouteHandlers  {
	return &ApiRouteHandlers{
		dataProvider: deviceControl,
		middleWare:s.NewMiddleware(config),
		router:mux.NewRouter()}
}

func (apiHandlers *ApiRouteHandlers) Router() *mux.Router  {
	return apiHandlers.router
}

func (apiHandlers *ApiRouteHandlers) InitRoutes()  {
	apiHandlers.router.NotFoundHandler = http.HandlerFunc(apiHandlers.HandleNotFound)
	apiHandlers.router.Use(apiHandlers.middleWare.HeadersMiddleware)
	apiHandlers.router.Use(apiHandlers.middleWare.AuthTokenMiddleware)

	// Run routes
	apiHandlers.router.HandleFunc("/run/command/{commandId}", apiHandlers.handleRunCommand)
	apiHandlers.router.HandleFunc("/run/scenario/{scenarioId}", apiHandlers.handleRunScenario)
	apiHandlers.router.HandleFunc("/run/intent", apiHandlers.handleRunIntent).Methods("POST")

	// Api routes
	apiHandlers.router.HandleFunc("/controls", apiHandlers.handleControls)
}
