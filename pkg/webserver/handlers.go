package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math"
	"net/http"
	"smh-apiengine/pkg/alexakit"
	"smh-apiengine/pkg/devicecontrol"
	"time"
)

type ApiRouteHandlers struct {
	dataProvider *devicecontrol.DeviceControl
	middleware []mux.MiddlewareFunc
	router *mux.Router
	routesInited time.Time
}

func NewApiRouteHandlers(config *ServerConfig, deviceControl *devicecontrol.DeviceControl) *ApiRouteHandlers  {
	authMiddleware := AuthMiddleware{Token: config.Token}
	headersMiddleware := HeadersMiddleware{}
	middleware := []mux.MiddlewareFunc{
		authMiddleware.Middleware,
		headersMiddleware.Middleware}

	return &ApiRouteHandlers{
		dataProvider: deviceControl,
		middleware:middleware,
		router:mux.NewRouter(),
		routesInited: time.Now()}
}

func (apiHandlers *ApiRouteHandlers) Router() *mux.Router {
	return apiHandlers.router
}

func (apiHandlers *ApiRouteHandlers) InitRoutes()  {
	apiHandlers.router.NotFoundHandler = http.HandlerFunc(apiHandlers.handleNotFound)

	for _, middlewareFunc := range apiHandlers.middleware {
		apiHandlers.router.Use(middlewareFunc)
	}

	apiHandlers.router.HandleFunc("/uptime", apiHandlers.handleUptime)

	// Run routes
	apiHandlers.router.HandleFunc("/run/command/{commandId}", apiHandlers.handleRunCommand)
	apiHandlers.router.HandleFunc("/run/scenario/{scenarioId}", apiHandlers.handleRunScenario)
	apiHandlers.router.HandleFunc("/run/intent", apiHandlers.handleRunIntent).Methods("POST")

	// Api routes
	apiHandlers.router.HandleFunc("/controls", apiHandlers.handleControls)
	apiHandlers.router.HandleFunc("/device/state", apiHandlers.handleWebsocketDeviceState)
}

// handleNotFound used for not found responses
func (apiHandlers *ApiRouteHandlers) handleNotFound(w http.ResponseWriter, r *http.Request) {
	_, ioErr := io.WriteString(w, NewErrorResponse("Resource not found"))

	if ioErr != nil {
		log.Println(ioErr)
	}
}

// handlePing simple handler that returns a so called "pong" response
func (apiHandlers *ApiRouteHandlers) handleUptime(w http.ResponseWriter, r *http.Request) {
	type uptime struct {
		Since time.Duration `json:"uptime"`
		StartedOn string `json:"started_on"`
	}

	days := 0
	uptimeData := uptime{
		Since: time.Since(apiHandlers.routesInited),
		StartedOn: fmt.Sprintf(
			"%d.%d.%d at %d:%d",
			apiHandlers.routesInited.Day(),
			apiHandlers.routesInited.Month(),
			apiHandlers.routesInited.Year(),
			apiHandlers.routesInited.Hour(),
			apiHandlers.routesInited.Minute())}

	if uptimeData.Since.Hours() >= 24 {
		days = int(math.Round(uptimeData.Since.Hours() / 24))
	}

	_, ioErr := io.WriteString(w, NewSuccessResponse(
		fmt.Sprintf(
			"Uptime %d day(s), %.f hour(s), %.f minute(s), %.f second(s)",
			days,
			uptimeData.Since.Hours(),
			uptimeData.Since.Minutes(),
			uptimeData.Since.Seconds()),
			uptimeData))
	if ioErr != nil {
		log.Println(ioErr)
	}
}

// handleControls api handler that accepts alexa request JSON and tries to execute matched scenario or command
func (apiHandlers *ApiRouteHandlers) handleControls(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	_, ioErr := io.WriteString(w, NewSuccessResponse("controls", apiHandlers.dataProvider.AllControls()))

	if ioErr != nil {
		log.Println(ioErr)
	}
}

// handleRunIntent api action that accepts alexa request JSON and tries to execute matched scenario or command
func (apiHandlers *ApiRouteHandlers) handleRunIntent(w http.ResponseWriter, r *http.Request) {
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

	simpleAlexaIntent, err := apiHandlers.dataProvider.NewSimpleRequestIntent(alexaRequestIntent)

	if err != nil {
		_, ioErr := io.WriteString(w, NewErrorResponse("Failed to create a simple alexa request intent"))

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

	_, err = io.WriteString(w, NewSuccessResponse("intent executed", nil))

	if err != nil {
		log.Println(err)
	}
}

// handleRunCommand api action that accepts command id and tries to execute matched command
func (apiHandlers *ApiRouteHandlers) handleRunCommand(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	commandID := vars["commandId"]
	cmd := apiHandlers.dataProvider.FindCommandByID(commandID)

	if cmd == nil {
		_, ioErr := io.WriteString(w, NewErrorResponse(
			fmt.Sprintf("Command with id %s was not found", commandID)))

		if ioErr != nil {
			log.Println(ioErr)
		}

		return
	}

	go func() {
		err := apiHandlers.dataProvider.ExecCommandFullCycle(*cmd)

		if err != nil {
			log.Println(err)
		}
	}()

	_, err := io.WriteString(w, NewSuccessResponse("command executed", nil))

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
		_, ioErr := io.WriteString(w, NewErrorResponse(
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

	_, err = io.WriteString(w, NewSuccessResponse("scenario executed", nil))

	if err != nil {
		log.Println(err)
	}
}

func (apiHandlers *ApiRouteHandlers) handleWebsocketDeviceState(w http.ResponseWriter, r *http.Request)  {
	conn, _, _, err := ws.UpgradeHTTP(r, w)

	if err != nil {
		log.Println(err)

		return
	}

	go func() {
		defer func() {
			err = conn.Close()
			if err != nil {
				log.Println("Failed to close the connection")
			}
		}()

		for {
			if len(apiHandlers.dataProvider.DeviceStateBuffer) == 0 {
				continue
			}

			for time, bufferItem := range apiHandlers.dataProvider.DeviceStateBuffer {
				jsonStr, err := json.Marshal(bufferItem)

				if err != nil {
					log.Println(err)
					continue
				}

				err = wsutil.WriteServerText(conn, jsonStr)

				if err != nil {
					log.Println(err)
				}

				delete(apiHandlers.dataProvider.DeviceStateBuffer, time)
			}
		}
	}()
}