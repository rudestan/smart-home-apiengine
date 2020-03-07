package apiserver

import (
	"io"
	"log"
	"net/http"
	"smh-apiengine/pkg/devicecontrol"
)

// handleNotFound used for not found responses
func (s *server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_, ioErr := io.WriteString(w, NewErrorResponse("Resource not found"))

	if ioErr != nil {
		log.Println(ioErr)
	}
}

// handleControls api handler that accepts alexa request JSON and tries to execute matched scenario or command
func (s *server) handleControls(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !isRequestAuthenticated(s.token, w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)

	_, ioErr := io.WriteString(w, NewSuccessResponse("controls", devicecontrol.AllControls()))

	if ioErr != nil {
		log.Println(ioErr)
	}
}
