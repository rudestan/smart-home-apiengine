package apiserver

import (
	"io"
	"log"
	"net/http"
	"smh-apiengine/pkg/devicecontrol"
)

// NotFoundHandler used for not found responses
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_, ioErr := io.WriteString(w, NewErrorResponse("Resource not found"))

	if ioErr != nil {
		log.Println(ioErr)
	}
}

// ActionControls api action that accepts alexa request JSON and tries to execute matched scenario or command
func ActionControls(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !isRequestAuthenticated(authToken, w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)

	_, ioErr := io.WriteString(w, NewSuccessResponse("controls", devicecontrol.AllControls()))

	if ioErr != nil {
		log.Println(ioErr)
	}
}
