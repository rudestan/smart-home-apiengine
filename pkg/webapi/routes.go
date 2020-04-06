package webapi

import (
	"net/http"
	"smh-apiengine/pkg/devicecontrol"
	s "smh-apiengine/pkg/webserver"
	"github.com/gorilla/mux"
)

type ApiRouteHandlers struct {
	dataProvider *devicecontrol.DeviceControl
	middleware []mux.MiddlewareFunc
	router *mux.Router
}

func NewApiRouteHandlers(config *s.ServerConfig, deviceControl *devicecontrol.DeviceControl) *ApiRouteHandlers  {
	middleware := []mux.MiddlewareFunc{
		s.AuthMiddleware{Token:config.Token}.Middleware,
		s.HeadersMiddleware{}.Middleware}

	return &ApiRouteHandlers{
		dataProvider: deviceControl,
		middleware:middleware,
		router:mux.NewRouter()}
}

func (apiHandlers *ApiRouteHandlers) Router() *mux.Router  {
	return apiHandlers.router
}

func (apiHandlers *ApiRouteHandlers) InitRoutes()  {
	apiHandlers.router.NotFoundHandler = http.HandlerFunc(apiHandlers.HandleNotFound)

	for _, middlewareFunc := range apiHandlers.middleware {
		apiHandlers.router.Use(middlewareFunc)
	}

	// Run routes
	apiHandlers.router.HandleFunc("/run/command/{commandId}", apiHandlers.handleRunCommand)
	apiHandlers.router.HandleFunc("/run/scenario/{scenarioId}", apiHandlers.handleRunScenario)
	apiHandlers.router.HandleFunc("/run/intent", apiHandlers.handleRunIntent).Methods("POST")

	// Api routes
	apiHandlers.router.HandleFunc("/controls", apiHandlers.handleControls)
}
