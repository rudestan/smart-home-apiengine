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

// ActionRunIntent api action that accepts alexa request JSON and tries to execute matched scenario or command
func ActionRunIntent(w http.ResponseWriter, r *http.Request)  {
    logRequest(r)
    w.Header().Set("Content-Type", "application/json")

    if !isRequestAuthenticated(authToken, w, r) {
        return
    }

    w.WriteHeader(http.StatusOK)

    alexaRequestIntent, err := alexakit.NewAlexaRequestIntent(r)

    if err != nil {
        _, ioErr := io.WriteString(w, NewErrorResponse("Failed to accept POST body of alexa intent"))

        if ioErr != nil {
            log.Println(ioErr)
        }

        log.Println(err)

        return
    }

    simpleAlexaIntent, err := devicecontrol.NewSimpleRequestIntent(alexaRequestIntent)

    if err != nil {
        _, ioErr := io.WriteString(w, NewErrorResponse("Failed to create a simple alexa request intent"))

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

    io.WriteString(w, NewSuccessResponse("intent executed", nil))
}

// ActionRunCommand api action that accepts command id and tries to execute matched command
func ActionRunCommand(w http.ResponseWriter, r *http.Request)  {
    logRequest(r)
    w.Header().Set("Content-Type", "application/json")

    if !isRequestAuthenticated(authToken, w, r) {
        return
    }

    w.WriteHeader(http.StatusOK)

    vars := mux.Vars(r)
    commandId := vars["commandId"]

    cmd, err := devicecontrol.FindCommandById(commandId)

    if err != nil {
        _, ioErr := io.WriteString(w, NewErrorResponse(
            fmt.Sprintf("Command with id %s was not found", commandId)))

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

    _, err = io.WriteString(w, NewSuccessResponse("command executed", nil))

    if err != nil {
        log.Println(err)
    }
}

// ActionRunScenario api action that accepts scenario id and tries to execute matched scenario
func ActionRunScenario(w http.ResponseWriter, r *http.Request)  {
    logRequest(r)
    w.Header().Set("Content-Type", "application/json")

    if !isRequestAuthenticated(authToken, w, r) {
        return
    }

    w.WriteHeader(http.StatusOK)

    vars := mux.Vars(r)
    scenarioId := vars["scenarioId"]
    scenario, err := devicecontrol.FindScenarioByName(scenarioId)

    if err != nil {
        _, ioErr := io.WriteString(w, NewErrorResponse(
            fmt.Sprintf("Scenario with id %s was not found", scenarioId)))

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

    _, err = io.WriteString(w, NewSuccessResponse("scenario executed", nil))

    if err != nil {
        log.Println(err)
    }
}
