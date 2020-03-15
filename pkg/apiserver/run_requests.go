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
func (s *Server) handleRunIntent(w http.ResponseWriter, r *http.Request) {
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

	dc := s.dataProvider.(devicecontrol.DeviceControl)
	simpleAlexaIntent, err := dc.NewSimpleRequestIntent(alexaRequestIntent)

	if err != nil {
		_, ioErr := io.WriteString(w, newErrorResponse("Failed to create a simple alexa request intent"))

		if ioErr != nil {
			log.Println(ioErr)
		}

		log.Println(err)

		return
	}

	go func() {
		err = dc.HandleAlexaRequest(simpleAlexaIntent)

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
func (s *Server) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	commandID := vars["commandId"]
	dc := s.dataProvider.(devicecontrol.DeviceControl)
	cmd, err := dc.FindCommandByID(commandID)

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
		err = dc.ExecCommandFullCycle(cmd)

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
func (s *Server) handleRunScenario(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	scenarioID := vars["scenarioId"]
	dc := s.dataProvider.(devicecontrol.DeviceControl)
	scenario, err := dc.FindScenarioByName(scenarioID)

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
		err = dc.ExecScenarioFullCycle(scenario)

		if err != nil {
			log.Println(err)
		}
	}()

	_, err = io.WriteString(w, newSuccessResponse("scenario executed", nil))

	if err != nil {
		log.Println(err)
	}
}
