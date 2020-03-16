package amqp

import (
	"bytes"
	"log"
	"net/http"
)

type Handler struct {
	EndPoint string
}

// handle initializes AlexaRequest struct with all intents and slots received in json message payload.
// Then it creates simplified filtered struct and performs execution with device control package.
func (h *Handler) handle(req string) {
	err := h.postToApi(req)

	if err != nil {
		log.Println(err)
	}
}

func (h* Handler) postToApi(req string) error {
	resp, err := http.Post(h.EndPoint, "application/json", bytes.NewBufferString(req))

	defer func() {
		if resp != nil && resp.Body != nil {
			err = resp.Body.Close()

			if err != nil {
				log.Println("error closing response body: ", err)
			}
		}
	}()

	if err != nil {
		return err
	}

	log.Println(resp)

	return nil
}
