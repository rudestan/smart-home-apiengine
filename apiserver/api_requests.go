package apiserver

import (
	"apiengine/devicecontrol"
	"io"
	"log"
	"net/http"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request)  {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	_, ioErr := io.WriteString(w, NewErrorResponse("Resource not found"))

	if ioErr != nil {
		log.Println(ioErr)
	}
}

// ActionRunIntent api action that accepts alexa request JSON and tries to execute matched scenario or command
func ActionGroups(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")

	if !isRequestAuthenticated(authToken, w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)

	_, ioErr := io.WriteString(w, NewSuccessResponse("groups", devicecontrol.AllGroups()))

	if ioErr != nil {
		log.Println(ioErr)
	}
}
