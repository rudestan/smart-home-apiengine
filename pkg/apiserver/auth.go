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

func (s *server) isTokenValid(authHeader string) bool {
	if s.token == "" {
		return true
	}

	if len(authHeader) > 0 && strings.HasPrefix(authHeader, bearerPrefix) {
		return s.token == strings.TrimPrefix(authHeader, bearerPrefix)
	}

	return false
}

func (s *server) isRequestAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	if s.isTokenValid(r.Header.Get(headerAuthorization)) {
		return true
	}

	log.Printf("Wrong token provided: %s\n", r.Header.Get(headerAuthorization))
	w.WriteHeader(http.StatusForbidden)

	_, ioErr := io.WriteString(w, newErrorResponse("Wrong token provided!"))
	if ioErr != nil {
		log.Println(ioErr)
	}

	return false
}
