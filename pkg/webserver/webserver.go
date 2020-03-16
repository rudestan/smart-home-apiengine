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
	Init(r *mux.Router)
}

type Server struct {
	router *mux.Router
	config *ServerConfig
	server *http.Server
}

func NewServer(serverConfig *ServerConfig, router *mux.Router, handlers RouteHandlers) *Server {
	server := &Server{
		router: router,
		config: serverConfig,
		server: &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}

	handlers.Init(router)

	return server
}

// ServeHTTP runs http server
func (s *Server) ServeHTTP() {
	s.logProcess()
	log.Println(s.server.ListenAndServe())
}

// ServeHTTPS runs https server using provided TLS certificate and key
func (s *Server) ServeHTTPS() {
	s.logProcess()
	log.Println(s.server.ListenAndServeTLS(s.config.TLSCert, s.config.TLSKey))
}

func (s *Server) logProcess()  {
	log.Printf("%s server is listening requests on %s:%d\n",
		s.config.Protocol, s.config.Address, s.config.Port)

	if s.config.Token != "" {
		log.Printf("Requests should use the following token: \"%s\"\n", s.config.Token)
	}
}
