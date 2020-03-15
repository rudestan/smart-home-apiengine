package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
	"smh-apiengine/pkg/apiserver"
	"smh-apiengine/pkg/directpublisher"
	"smh-apiengine/pkg/amqp"
)

const (
	defaultProtocol = "http"
	defaultAddress  = "127.0.0.1"
	defaultPort     = 8844
)

const (
	rmqHost       = "localhost"
	rmqPort       = 5672
	rmqLogin      = "guest"
	rmqPassword   = "guest"
	rmqExchange   = "alexa_sync"
	rmqQueue      = "alexa.responses"
	rmqRoutingKey = "alexa.response.json"
)

func main() {
	var logFile string
	var srvConfig apiserver.ServerConfig
	var rmqConfig amqp.Config

	execName, err := os.Executable()

	if err != nil {
		execName = "rmq-direct-publisher"
	}

	app := &cli.App{
		Name: "Smart Home Alexa Request Direct RMQ Publisher",
		Description: "A Web server that listens for incoming Alexa requests and publishes them to the RMQ",
		Usage:     "an app and web server for publishing the requests to the RMQ",
		UsageText: fmt.Sprintf("%s [global options] command [command options]", path.Base(execName)),
		HideHelp:  false,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log",
				Usage:       "Log file for logs output",
				Destination: &logFile,
				Aliases:     []string{"l"},
				EnvVars:	 []string{"RMQ_DIRECT_PUBLISHER_LOG_FILE"},
			},
			&cli.StringFlag{
				Name:        "proto",
				Value:       defaultProtocol,
				Usage:       "Protocol for web server to run (values: \"http\", \"https\")",
				Destination: &srvConfig.Protocol,
				Aliases:     []string{"r"},
				EnvVars:	 []string{"RMQ_DIRECT_PUBLISHER_SERVER_PROTO"},
			},
			&cli.StringFlag{
				Name:        "address",
				Value:       defaultAddress,
				Usage:       "Ip address for web server",
				Destination: &srvConfig.Address,
				Aliases:     []string{"a"},
				EnvVars:	 []string{"RMQ_DIRECT_PUBLISHER_SERVER_IP_ADDRESS"},
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       defaultPort,
				Usage:       "Port for web server",
				Destination: &srvConfig.Port,
				Aliases:     []string{"p"},
				EnvVars:	 []string{"RMQ_DIRECT_PUBLISHER_SERVER_PORT"},
			},
			&cli.StringFlag{
				Name:        "tls-cert",
				Usage:       "TLS Certificate file path (only when https protocol is used)",
				Destination: &srvConfig.TLSCert,
				Aliases:     []string{"c"},
				EnvVars:	 []string{"RMQ_DIRECT_PUBLISHER_SERVER_TLS_CERT"},
			},
			&cli.StringFlag{
				Name:        "tls-key",
				Usage:       "TLS Key file path (only when https protocol is used)",
				Destination: &srvConfig.TLSKey,
				Aliases:     []string{"k"},
				EnvVars:	 []string{"RMQ_DIRECT_PUBLISHER_SERVER_TLS_KEY"},
			},
			&cli.StringFlag{
				Name:        "rmqhost",
				Value:       rmqHost,
				Usage:       "RabbitMQ Host",
				Destination: &rmqConfig.Host,
				Aliases:     []string{"t"},
			},
			&cli.IntFlag{
				Name:        "rmqport",
				Value:       rmqPort,
				Usage:       "RabbitMQ Port",
				Destination: &rmqConfig.Port,
				Aliases:     []string{"o"},
			},
			&cli.StringFlag{
				Name:        "rmqlogin",
				Value:       rmqLogin,
				Usage:       "RabbitMQ Login",
				Destination: &rmqConfig.Login,
				Aliases:     []string{"i"},
			},
			&cli.StringFlag{
				Name:        "rmqpassword",
				Value:       rmqPassword,
				Usage:       "RabbitMQ Password",
				Destination: &rmqConfig.Password,
				Aliases:     []string{"s"},
			},
			&cli.StringFlag{
				Name:        "rmqexchange",
				Value:       rmqExchange,
				Usage:       "RabbitMQ Exchange name",
				Destination: &rmqConfig.Exchange,
				Aliases:     []string{"e"},
			},
			&cli.StringFlag{
				Name:        "rmqqueue",
				Value:       rmqQueue,
				Usage:       "RabbitMQ Queue name",
				Destination: &rmqConfig.Queue,
				Aliases:     []string{"q"},
			},
			&cli.StringFlag{
				Name:        "rmqrkey",
				Value:       rmqRoutingKey,
				Usage:       "RabbitMQ Queue name",
				Destination: &rmqConfig.RoutingKey,
				Aliases:     []string{"n"},
			},
		},
		Action: func(c *cli.Context) error {
			if logFile != "" {
				err := setLogOutputToFile(logFile)
				if err != nil {
					log.Println("Failed to set the log output to file: ", err)
				}
			}

			return runServer(srvConfig, rmqConfig)
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(serverConfig apiserver.ServerConfig, rmqConfig amqp.Config) error {
	if serverConfig.Protocol == "https" {
		if serverConfig.TLSCert == "" || serverConfig.TLSKey == "" {
			return errors.New("TLS Certificate and Key files are required when using https protocol")
		}
	}

	directPublisher := directpublisher.NewDirectPublisher(&rmqConfig)
	server := apiserver.NewServer(serverConfig, mux.NewRouter(), directPublisher, nil)
	switch serverConfig.Protocol {
	case "http":
		apiserver.ServeHTTP(server)
	case "https":
		apiserver.ServeHTTPS(server)
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

