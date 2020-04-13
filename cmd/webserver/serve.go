package main

import (
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
	"smh-apiengine/pkg/devicecontrol"
	"smh-apiengine/pkg/webserver"
	"time"
)

const (
	defaultProtocol = "http"
	defaultAddress  = "127.0.0.1"
	defaultPort     = 8787
	defaultStartRetires = 5
	defaultStartRetryInterval = 3 // seconds
)

func main() {
	var configFile string
	var logFile string
	var srvConfig webserver.ServerConfig

	execName, err := os.Executable()

	if err != nil {
		execName = "webserver"
	}

	app := &cli.App{
		Name: "Smart Home Broadlink API Server",
		Description: "A Web server that serves an incoming requests and controls Broadlink devices using API. Uses a " +
			"user created JSON for pre-configured device commands, scenarios and devices. Supports an incoming requests" +
			"from Alexa API web servers.",
		Usage:     "an app and web server for controlling Broadlink devices",
		UsageText: fmt.Sprintf("%s [global options] command [command options]", path.Base(execName)),
		HideHelp:  false,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       "Path to JSON configuration with commands and devices",
				Destination: &configFile,
				Aliases:     []string{"c"},
				EnvVars:	 []string{"SMH_CONFIG"},
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "log",
				Usage:       "Log file for logs output",
				Destination: &logFile,
				Aliases:     []string{"l"},
				EnvVars:	 []string{"SMH_SERVER_LOG_FILE"},
			},
			&cli.StringFlag{
				Name:        "proto",
				Value:       defaultProtocol,
				Usage:       "Protocol for web server to run (values: \"http\", \"https\")",
				Destination: &srvConfig.Protocol,
				Aliases:     []string{"r"},
				EnvVars:	 []string{"SMH_SERVER_PROTO"},
			},
			&cli.StringFlag{
				Name:        "address",
				Value:       defaultAddress,
				Usage:       "Ip address for web server",
				Destination: &srvConfig.Address,
				Aliases:     []string{"a"},
				EnvVars:	 []string{"SMH_SERVER_IP_ADDRESS"},
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       defaultPort,
				Usage:       "Port for web server",
				Destination: &srvConfig.Port,
				Aliases:     []string{"p"},
				EnvVars:	 []string{"SMH_SERVER_PORT"},
			},
			&cli.StringFlag{
				Name:        "tls-cert",
				Usage:       "TLS Certificate file path (only when https protocol is used)",
				Destination: &srvConfig.TLSCert,
				Aliases:     []string{"s"},
				EnvVars:	 []string{"SMH_SERVER_TLS_CERT"},
			},
			&cli.StringFlag{
				Name:        "tls-key",
				Usage:       "TLS Key file path (only when https protocol is used)",
				Destination: &srvConfig.TLSKey,
				Aliases:     []string{"k"},
				EnvVars:	 []string{"SMH_SERVER_TLS_KEY"},
			},
			&cli.StringFlag{
				Name:        "token",
				Usage:       "Authorization token bearer for the requests",
				Destination: &srvConfig.Token,
				Aliases:     []string{"t"},
				EnvVars:	 []string{"SMH_SERVER_AUTH_TOKEN"},
			},
		},
		Action: func(c *cli.Context) error {
			if logFile != "" {
				err := setLogOutputToFile(logFile)
				if err != nil {
					log.Println("Failed to set the log output to file: ", err)
				}
			}

			config, err := devicecontrol.NewConfiguration(configFile)

			if err != nil {
				return err
			}

			deviceControl := devicecontrol.NewDeviceControl(&config)

			return runServer(&srvConfig, &deviceControl)
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(serverConfig *webserver.ServerConfig, deviceControl *devicecontrol.DeviceControl) error {
	if serverConfig.Protocol == "https" {
		if serverConfig.TLSCert == "" || serverConfig.TLSKey == "" {
			return errors.New("TLS Certificate and Key files are required when using https protocol")
		}
	}

	apiRouteHandlers := webserver.NewApiRouteHandlers(serverConfig, deviceControl)
	server := webserver.NewServer(serverConfig, apiRouteHandlers)

	var err error

	for i := 0; i < defaultStartRetires; i++ {
		if serverConfig.Protocol == "https" {
			err = server.ServeHTTPS()
		} else {
			err = server.ServeHTTP()
		}

		if err != nil {
			log.Printf("Failed: %s; attempt %d of %d ...\n", err.Error(), i + 1, defaultStartRetires)
			time.Sleep(time.Second * defaultStartRetryInterval)
		}
	}

	return nil
}

func setLogOutputToFile(fileName string) error {
	logFile, err := os.Create(fileName)
	if err != nil {
		return err
	}

	log.SetOutput(logFile)

	return nil
}
