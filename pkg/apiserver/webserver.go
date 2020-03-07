package apiserver

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type server struct {
	router *mux.Router
	token string
	server *http.Server
}

func newServer(serverAddr string, router *mux.Router, token string) *server {
	server := &server{
		router: router,
		token: token,
		server: &http.Server{
			Handler:      router,
			Addr:         serverAddr,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
	}

	server.routes()

	return server
}

func (s *server) routes()  {
	s.router.NotFoundHandler = http.HandlerFunc(s.handleNotFound)

	// Run routes
	s.router.HandleFunc("/run/command/{commandId}", s.handleRunCommand)
	s.router.HandleFunc("/run/scenario/{scenarioId}", s.handleRunScenario)
	s.router.HandleFunc("/run/intent", s.handleRunIntent).Methods("POST")

	// Api routes
	s.router.HandleFunc("/controls", s.handleControls)
}

// ServeHTTP runs http server
func ServeHTTP(serverAddr string, token string) {
	srv := newServer(serverAddr, mux.NewRouter(), token)

	log.Printf("http api server is listening requests on %s\n", serverAddr)

	if token != "" {
		log.Printf("Requests should use the following token: \"%s\"\n", token)
	}

	log.Println(srv.server.ListenAndServe())
}

// ServeHTTPS runs https server using provided TLS certificate and key
func ServeHTTPS(serverAddr string, token string, certFile string, keyFile string) {
	srv := newServer(serverAddr, mux.NewRouter(), token)

	log.Printf("http api is server serving requests on %s\n", serverAddr)

	if token != "" {
		log.Printf("Requests should use the following token: \"%s\"\n", token)
	}

	log.Println(srv.server.ListenAndServeTLS(certFile, keyFile))
}

func logRequest(r *http.Request) {
	log.Printf("Request: \"%s\", from: %s", r.RequestURI, r.RemoteAddr)
}
