package apiserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"smh-apiengine/pkg/alexakit"
	"smh-apiengine/pkg/devicecontrol"

	"github.com/gorilla/mux"
)

// handleRunIntent api action that accepts alexa request JSON and tries to execute matched scenario or command
func (s *server) handleRunIntent(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !s.isRequestAuthenticated(w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)

	alexaRequestIntent, err := alexakit.NewAlexaRequestIntent(r)

	if err != nil {
		_, ioErr := io.WriteString(w, newErrorResponse("Failed to accept POST body of alexa intent"))

		if ioErr != nil {
			log.Println(ioErr)
		}

		log.Println(err)

		return
	}

	simpleAlexaIntent, err := devicecontrol.NewSimpleRequestIntent(alexaRequestIntent)

	if err != nil {
		_, ioErr := io.WriteString(w, newErrorResponse("Failed to create a simple alexa request intent"))

		if ioErr != nil {
			log.Println(ioErr)
		}

		log.Println(err)

		return
	}

	go func() {
		err = devicecontrol.HandleAlexaRequest(simpleAlexaIntent)

		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.WriteString(w, newSuccessResponse("intent executed", nil))

	if err != nil {
		log.Println(err)
	}
}

// handleRunCommand api action that accepts command id and tries to execute matched command
func (s *server) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !s.isRequestAuthenticated(w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	commandID := vars["commandId"]

	cmd, err := devicecontrol.FindCommandByID(commandID)

	if err != nil {
		_, ioErr := io.WriteString(w, newErrorResponse(
			fmt.Sprintf("Command with id %s was not found", commandID)))

		if ioErr != nil {
			log.Println(ioErr)
		}

		log.Println(err)

		return
	}

	go func() {
		err = devicecontrol.ExecCommandFullCycle(cmd)

		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.WriteString(w, newSuccessResponse("command executed", nil))

	if err != nil {
		log.Println(err)
	}
}

// handleRunScenario api action that accepts scenario id and tries to execute matched scenario
func (s *server) handleRunScenario(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if !s.isRequestAuthenticated(w, r) {
		return
	}

	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	scenarioID := vars["scenarioId"]
	scenario, err := devicecontrol.FindScenarioByName(scenarioID)

	if err != nil {
		_, ioErr := io.WriteString(w, newErrorResponse(
			fmt.Sprintf("Scenario with id %s was not found", scenarioID)))

		if ioErr != nil {
			log.Println(err)
		}

		log.Println(err)

		return
	}

	go func() {
		err = devicecontrol.ExecScenarioFullCycle(scenario)

		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.WriteString(w, newSuccessResponse("scenario executed", nil))

	if err != nil {
		log.Println(err)
	}
}
