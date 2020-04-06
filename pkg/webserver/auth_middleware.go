package webserver

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

// ResultResponse api response for messages without payload
type ResultResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
}

type AuthMiddleware struct {
	Token string
}

func (am *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if am.Token == "" || am.isTokenValid(am.Token, r.Header.Get(headerAuthorization)) {
			next.ServeHTTP(w, r)

			return
		}

		log.Printf("Wrong token provided: %s\n", r.Header.Get(headerAuthorization))
		w.WriteHeader(http.StatusForbidden)

		_, ioErr := io.WriteString(w, NewErrorResponse("Wrong token provided!"))
		if ioErr != nil {
			log.Println(ioErr)
		}
	})
}

func (am *AuthMiddleware) isTokenValid(token string, authHeader string) bool {
	if len(authHeader) > 0 && strings.HasPrefix(authHeader, bearerPrefix) {
		return token == strings.TrimPrefix(authHeader, bearerPrefix)
	}

	return false
}
