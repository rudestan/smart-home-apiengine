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

func isValidToken(token string, authHeader string) bool {
	if token == "" {
		return true
	}

	if len(authHeader) > 0 && strings.HasPrefix(authHeader, bearerPrefix) {
		return token == strings.TrimPrefix(authHeader, bearerPrefix)
	}

	return false
}

func isRequestAuthenticated(token string, w http.ResponseWriter, r *http.Request) bool {
	if isValidToken(token, r.Header.Get(headerAuthorization)) {
		return true
	}

	log.Printf("Wrong token provided: %s\n", r.Header.Get(headerAuthorization))
	w.WriteHeader(http.StatusForbidden)

	_, ioErr := io.WriteString(w, NewErrorResponse("Wrong token provided!"))
	if ioErr != nil {
		log.Println(ioErr)
	}

	return false
}
