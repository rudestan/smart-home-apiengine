package apiserver

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func newServer(serverAddr string) *http.Server {
	r := mux.NewRouter()

	r.HandleFunc("/run/command/{commandId}", ActionRunCommand)
	r.HandleFunc("/run/scenario/{scenarioId}", ActionRunScenario)
	r.
		HandleFunc("/run/intent", ActionRunIntent).
		Methods("POST")

	return &http.Server{
		Handler:      r,
		Addr:         serverAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
}

// ServeHttp runs http server
func ServeHttp(serverAddr string)  {
	srv := newServer(serverAddr)
	log.Printf("http api server is listening requests on %s", serverAddr)
	log.Println(srv.ListenAndServe())
}

// ServeHttps runs https server using provided TLS certificate and key
func ServeHttps(serverAddr string, certFile string, keyFile string)  {
	srv := newServer(serverAddr)
	log.Printf("http api is server serving requests on %s", serverAddr)
	log.Println(srv.ListenAndServeTLS(certFile, keyFile))
}
