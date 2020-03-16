package webapi

import (
	"io"
	"log"
	"net/http"
	"smh-apiengine/pkg/webserver"
)

// handleNotFound used for not found responses
func (apiHandlers *ApiRouteHandlers) HandleNotFound(w http.ResponseWriter, r *http.Request) {
	_, ioErr := io.WriteString(w, webserver.NewErrorResponse("Resource not found"))

	if ioErr != nil {
		log.Println(ioErr)
	}
}

// handleControls api handler that accepts alexa request JSON and tries to execute matched scenario or command
func (apiHandlers *ApiRouteHandlers) handleControls(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, ioErr := io.WriteString(w, webserver.NewSuccessResponse("controls", apiHandlers.dataProvider.AllControls()))

	if ioErr != nil {
		log.Println(ioErr)
	}
}
