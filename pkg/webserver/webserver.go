package webserver

import (
	"fmt"
	"log"
	"net/http"
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

type RouteHandlers interface {
	Router() *mux.Router
	InitRoutes()
}

type server struct {
	config *ServerConfig
	server *http.Server
}

func NewServer(serverConfig *ServerConfig, routing interface{}) *server {
	rHandlers := routing.(RouteHandlers)

	rHandlers.InitRoutes()

	server := &server{
		config: serverConfig,
		server: &http.Server{
			Handler:      rHandlers.Router(),
			Addr:         fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}

	return server
}

// ServeHTTP runs http server
func (s *server) ServeHTTP() error {
	s.logProcess()
	return s.server.ListenAndServe()
}

// ServeHTTPS runs https server using provided TLS certificate and key
func (s *server) ServeHTTPS() error {
	s.logProcess()
	return s.server.ListenAndServeTLS(s.config.TLSCert, s.config.TLSKey)
}

func (s *server) logProcess()  {
	log.Printf("%s server is listening requests on %s:%d\n",
		s.config.Protocol, s.config.Address, s.config.Port)

	if s.config.Token != "" {
		log.Printf("Requests should use the following token: \"%s\"\n", s.config.Token)
	}
}
