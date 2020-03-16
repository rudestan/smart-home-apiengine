package webserver

import (
	"log"
	"net/http"
)

type Middleware struct {
	config *ServerConfig
}

func NewMiddleware(config *ServerConfig) *Middleware  {
	return &Middleware{config:config}
}

func (m* Middleware) HeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: \"%s\", from: %s", r.RequestURI, r.RemoteAddr)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
