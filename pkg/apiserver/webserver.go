package apiserver

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

type ServerRoutes interface {
	Init(r *mux.Router)
}

type Server struct {
	router *mux.Router
	config ServerConfig
	server *http.Server
	dataProvider interface{}
}

func NewServer(serverConfig ServerConfig, router *mux.Router, routes ServerRoutes, dataProvider interface{}) *Server {
	server := &Server{
		router: router,
		config: serverConfig,
		server: &http.Server{
			Handler:      router,
			Addr:         fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		dataProvider: dataProvider,
	}

	routes.Init(router)

	return server
}

// ServeHTTP runs http server
func ServeHTTP(s *Server) {
	s.logProcess()
	log.Println(s.server.ListenAndServe())
}

// ServeHTTPS runs https server using provided TLS certificate and key
func ServeHTTPS(s *Server) {
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