package apiserver

import (
	"io"
	"log"
	"net/http"
)

// handleNotFound used for not found responses
func (s *server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	_, ioErr := io.WriteString(w, newErrorResponse("Resource not found"))

	if ioErr != nil {
		log.Println(ioErr)
	}
}

// handleControls api handler that accepts alexa request JSON and tries to execute matched scenario or command
func (s *server) handleControls(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, ioErr := io.WriteString(w, newSuccessResponse("controls", s.dataProvider.AllControls()))

	if ioErr != nil {
		log.Println(ioErr)
	}
}
