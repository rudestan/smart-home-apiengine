package directpublisher

import (
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"smh-apiengine/pkg/alexakit"
	"smh-apiengine/pkg/webserver"
)

func (dp *DirectPublisher) Router() *mux.Router  {
	return dp.router
}

// Init implements routes webservers
func (dp *DirectPublisher) InitRoutes()  {
	headersMiddleware := webserver.HeadersMiddleware{}
	dp.router.Use(headersMiddleware.Middleware)
	dp.router.HandleFunc("/alexaproxy", dp.handleAlexaRequest).Methods("POST")
}

func (dp *DirectPublisher) handleAlexaRequest(w http.ResponseWriter, r *http.Request)  {
	w.WriteHeader(http.StatusOK)

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Recieved payload, pushing to the RMQ...")

	err = dp.rmq.Publish(string(reqBody))

	if err != nil {
		log.Println(err)

		return
	}

	dp.AlexaTextResponse(w)
}

func (dp *DirectPublisher) AlexaTextResponse(w http.ResponseWriter) {
	alexaResponse := alexakit.NewPlainTextSpeechResponse(alexakit.SpeechTextConfirmation)
	responseJson, err := alexaResponse.ToJson()

	if err != nil {
		log.Println(err)

		return
	}

	_, err = w.Write([]byte(responseJson))

	if err != nil {
		log.Println(err)
	}
}
