package apiserver

import (
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	headerAuthorization = "Authorization"
	bearerPrefix        = "Bearer "
)

func (s *Server) authTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.config.Token == "" || s.isTokenValid(r.Header.Get(headerAuthorization)) {
			next.ServeHTTP(w, r)

			return
		}

		log.Printf("Wrong token provided: %s\n", r.Header.Get(headerAuthorization))
		w.WriteHeader(http.StatusForbidden)

		_, ioErr := io.WriteString(w, newErrorResponse("Wrong token provided!"))
		if ioErr != nil {
			log.Println(ioErr)
		}
	})
}

func (s *Server) isTokenValid(authHeader string) bool {
	if len(authHeader) > 0 && strings.HasPrefix(authHeader, bearerPrefix) {
		return s.config.Token == strings.TrimPrefix(authHeader, bearerPrefix)
	}

	return false
}
