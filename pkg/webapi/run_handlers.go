package webapi

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"smh-apiengine/pkg/alexakit"
	"smh-apiengine/pkg/webserver"

	"github.com/gorilla/mux"
)

// handleRunIntent api action that accepts alexa request JSON and tries to execute matched scenario or command
func (apiHandlers *ApiRouteHandlers) handleRunIntent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	alexaRequestIntent, err := alexakit.NewAlexaRequestIntent(r)

	if err != nil {
		_, ioErr := io.WriteString(w, webserver.NewErrorResponse("Failed to accept POST body of alexa intent"))

		if ioErr != nil {
			log.Println(ioErr)
		}

		log.Println(err)

		return
	}

	simpleAlexaIntent, err := apiHandlers.dataProvider.NewSimpleRequestIntent(alexaRequestIntent)

	if err != nil {
		_, ioErr := io.WriteString(w, webserver.NewErrorResponse("Failed to create a simple alexa request intent"))

		if ioErr != nil {
			log.Println(ioErr)
		}

		log.Println(err)

		return
	}

	go func() {
		err = apiHandlers.dataProvider.HandleAlexaRequest(simpleAlexaIntent)

		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.WriteString(w, webserver.NewSuccessResponse("intent executed", nil))

	if err != nil {
		log.Println(err)
	}
}

// handleRunCommand api action that accepts command id and tries to execute matched command
func (apiHandlers *ApiRouteHandlers) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	commandID := vars["commandId"]
	cmd, err := apiHandlers.dataProvider.FindCommandByID(commandID)

	if err != nil {
		_, ioErr := io.WriteString(w, webserver.NewErrorResponse(
			fmt.Sprintf("Command with id %s was not found", commandID)))

		if ioErr != nil {
			log.Println(ioErr)
		}

		log.Println(err)

		return
	}

	go func() {
		err = apiHandlers.dataProvider.ExecCommandFullCycle(cmd)

		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.WriteString(w, webserver.NewSuccessResponse("command executed", nil))

	if err != nil {
		log.Println(err)
	}
}

// handleRunScenario api action that accepts scenario id and tries to execute matched scenario
func (apiHandlers *ApiRouteHandlers) handleRunScenario(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	scenarioID := vars["scenarioId"]
	scenario, err := apiHandlers.dataProvider.FindScenarioByName(scenarioID)

	if err != nil {
		_, ioErr := io.WriteString(w, webserver.NewErrorResponse(
			fmt.Sprintf("Scenario with id %s was not found", scenarioID)))

		if ioErr != nil {
			log.Println(err)
		}

		log.Println(err)

		return
	}

	go func() {
		err = apiHandlers.dataProvider.ExecScenarioFullCycle(scenario)

		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.WriteString(w, webserver.NewSuccessResponse("scenario executed", nil))

	if err != nil {
		log.Println(err)
	}
}
