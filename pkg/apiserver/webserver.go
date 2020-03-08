package apiserver

import (
	"fmt"
	"log"
	"net/http"
	"smh-apiengine/pkg/devicecontrol"
	"time"

	"github.com/gorilla/mux"
)
type ServerConfig struct {
	Protocol string
	Address  string
	Port     int
	Token    string
	TLSCert  string
	TLSKey   string
}

type server struct {
	router *mux.Router
	config ServerConfig
	server *http.Server
	dataProvider devicecontrol.DeviceControl
}

func newServer(serverConfig ServerConfig, router *mux.Router, control devicecontrol.DeviceControl) *server {
	server := &server{
		router: router,
		config: serverConfig,
		server: &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		dataProvider: control,
	}

	server.routes()

	return server
}

func (s *server) routes() {
	s.router.NotFoundHandler = http.HandlerFunc(s.handleNotFound)
	s.router.Use(s.headersMiddleware)
	s.router.Use(s.authTokenMiddleware)

	// Run routes
	s.router.HandleFunc("/run/command/{commandId}", s.handleRunCommand)
	s.router.HandleFunc("/run/scenario/{scenarioId}", s.handleRunScenario)
	s.router.HandleFunc("/run/intent", s.handleRunIntent).Methods("POST")

	// Api routes
	s.router.HandleFunc("/controls", s.handleControls)
}

func (s *server) logProcess()  {
	log.Printf("%s api server is listening requests on %s:%d\n",
		s.config.Protocol, s.config.Address, s.config.Port)

	if s.config.Token != "" {
		log.Printf("Requests should use the following token: \"%s\"\n", s.config.Token)
	}
}

// ServeHTTP runs http server
func ServeHTTP(serverConfig ServerConfig, control devicecontrol.DeviceControl) {
	srv := newServer(serverConfig, mux.NewRouter(), control)

	srv.logProcess()
	log.Println(srv.server.ListenAndServe())
}

// ServeHTTPS runs https server using provided TLS certificate and key
func ServeHTTPS(serverConfig ServerConfig, control devicecontrol.DeviceControl) {
	srv := newServer(serverConfig, mux.NewRouter(), control)

	srv.logProcess()
	log.Println(srv.server.ListenAndServeTLS(serverConfig.TLSCert, serverConfig.TLSKey))
}
