package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"smh-apiengine/pkg/apiserver"
	"smh-apiengine/pkg/devicecontrol"

	"github.com/urfave/cli/v2"
)

const (
	defaultProtocol = "http"
	defaultAddress  = "127.0.0.1"
	defaultPort     = 8787
)

type serverConfig struct {
	Protocol string
	Address  string
	Port     int
	TLSCert  string
	TLSKey   string
}

func main() {
	var configFile string
	var logFile string
	var srvConfig serverConfig
	var authToken string

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
				Required:    true,
			},
			&cli.StringFlag{
				Name:        "log",
				Usage:       "Log file for logs output",
				Destination: &logFile,
				Aliases:     []string{"l"},
			},
			&cli.StringFlag{
				Name:        "proto",
				Value:       defaultProtocol,
				Usage:       "Protocol for web server to run (values: \"http\", \"https\")",
				Destination: &srvConfig.Protocol,
				Aliases:     []string{"r"},
			},
			&cli.StringFlag{
				Name:        "address",
				Value:       defaultAddress,
				Usage:       "Ip address for web server",
				Destination: &srvConfig.Address,
				Aliases:     []string{"a"},
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       defaultPort,
				Usage:       "Port for web server",
				Destination: &srvConfig.Port,
				Aliases:     []string{"p"},
			},
			&cli.StringFlag{
				Name:        "tls-cert",
				Usage:       "TLS Certificate file path (only in case https protocol is used)",
				Destination: &srvConfig.TLSCert,
				Aliases:     []string{"s"},
			},
			&cli.StringFlag{
				Name:        "tls-key",
				Usage:       "TLS Key file path (only in case https protocol is used)",
				Destination: &srvConfig.TLSKey,
				Aliases:     []string{"k"},
			},
			&cli.StringFlag{
				Name:        "token",
				Usage:       "Authorization token bearer for the requests",
				Aliases:     []string{"t"},
				Destination: &authToken,
			},
		},
		Action: func(c *cli.Context) error {
			if logFile != "" {
				err := setLogOutputToFile(logFile)
				if err != nil {
					log.Println("Failed to set the log output to file: ", err)
				}
			}

			err := devicecontrol.Init(configFile)
			if err != nil {
				return err
			}

			return runServer(srvConfig, authToken)
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(serverConfig serverConfig, authToken string) error {
	if serverConfig.Protocol == "https" {
		if serverConfig.TLSCert == "" || serverConfig.TLSKey == "" {
			return errors.New("TLS Certificate and Key files are required when using https protocol")
		}
	}

	switch serverConfig.Protocol {
	case "http":
		apiserver.ServeHttp(fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port), authToken)
	case "https":
		apiserver.ServeHttps(
			fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port),
			authToken,
			serverConfig.TLSCert,
			serverConfig.TLSKey)
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
