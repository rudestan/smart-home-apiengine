package apiserver

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var authToken string

func newServer(serverAddr string) *http.Server {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	// Run routes
	r.HandleFunc("/run/command/{commandId}", ActionRunCommand)
	r.HandleFunc("/run/scenario/{scenarioId}", ActionRunScenario)
	r.HandleFunc("/run/intent", ActionRunIntent).Methods("POST")

	// Api routes
	r.HandleFunc("/controls", ActionControls)

	return &http.Server{
		Handler:      r,
		Addr:         serverAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

// ServeHTTP runs http server
func ServeHTTP(serverAddr string, token string) {
	srv := newServer(serverAddr)
	authToken = token
	log.Printf("http api server is listening requests on %s\n", serverAddr)

	if token != "" {
		log.Printf("Requests should use the following token: \"%s\"\n", token)
	}

	log.Println(srv.ListenAndServe())
}

// ServeHTTPS runs https server using provided TLS certificate and key
func ServeHTTPS(serverAddr string, token string, certFile string, keyFile string) {
	srv := newServer(serverAddr)
	authToken = token
	log.Printf("http api is server serving requests on %s\n", serverAddr)

	if token != "" {
		log.Printf("Requests should use the following token: \"%s\"\n", token)
	}

	log.Println(srv.ListenAndServeTLS(certFile, keyFile))
}

func logRequest(r *http.Request) {
	log.Printf("Request: \"%s\", from: %s", r.RequestURI, r.RemoteAddr)
}
