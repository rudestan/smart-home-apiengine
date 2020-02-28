package apiserver

import (
    "apiengine/alexakit"
    "apiengine/devicecontrol"
    "fmt"
    "github.com/gorilla/mux"
    "io"
    "log"
    "net/http"
)

func logRequest(r *http.Request)  {
    log.Printf("Request: \"%s\", from: %s", r.RequestURI, r.RemoteAddr)
}

// ActionRunIntent api action that accepts alexa request JSON and tries to execute matched scenario or command
func ActionRunIntent(w http.ResponseWriter, r *http.Request)  {
    logRequest(r)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    alexaRequestIntent, err := alexakit.NewAlexaRequestIntent(r)

    if err != nil {
        _, ioErr := io.WriteString(w, NewResponse("error", "Failed to accept POST body of alexa intent"))

        if ioErr != nil {
            log.Println(ioErr)
        }

        log.Println(err)

        return
    }

    simpleAlexaIntent, err := devicecontrol.NewSimpleRequestIntent(alexaRequestIntent)

    if err != nil {
        _, ioErr := io.WriteString(w,
            NewResponse("error", "Failed to create a simple alexa request intent"))

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

    io.WriteString(w, NewResponse("success", "intent executed"))
}

// ActionRunCommand api action that accepts command id and tries to execute matched command
func ActionRunCommand(w http.ResponseWriter, r *http.Request)  {
    logRequest(r)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    vars := mux.Vars(r)
    commandId := vars["commandId"]

    cmd, err := devicecontrol.FindCommandById(commandId)

    if err != nil {
        _, ioErr := io.WriteString(w, NewResponse("error",
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

    _, err = io.WriteString(w, NewResponse("success", "command executed"))

    if err != nil {
        log.Println(err)
    }
}

// ActionRunScenario api action that accepts scenario id and tries to execute matched scenario
func ActionRunScenario(w http.ResponseWriter, r *http.Request)  {
    logRequest(r)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    vars := mux.Vars(r)
    scenarioId := vars["scenarioId"]
    scenario, err := devicecontrol.FindScenarioByName(scenarioId)

    if err != nil {
        _, ioErr := io.WriteString(w, NewResponse("error",
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

    _, err = io.WriteString(w, NewResponse("success", "scenario executed"))

    if err != nil {
        log.Println(err)
    }
}
