package apiserver

import (
	"io"
	"log"
	"net/http"
	"smh-apiengine/pkg/devicecontrol"
)

// handleNotFound used for not found responses
func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	_, ioErr := io.WriteString(w, newErrorResponse("Resource not found"))

	if ioErr != nil {
		log.Println(ioErr)
	}
}

// handleControls api handler that accepts alexa request JSON and tries to execute matched scenario or command
func (s *Server) handleControls(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	dc := s.dataProvider.(devicecontrol.DeviceControl)
	_, ioErr := io.WriteString(w, newSuccessResponse("controls", dc.AllControls()))

	if ioErr != nil {
		log.Println(ioErr)
	}
}
