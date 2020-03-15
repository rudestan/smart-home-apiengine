package main

import (
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"smh-apiengine/pkg/amqp"
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
const apiEndpoint = "http://localhost:8787/run/intent"

func main() {
	var rmqConfig amqp.Config
	var msgHandler amqp.Handler
	var logFile string

	app := &cli.App{
		Name: "Smart home RMQ Proxy",
		Description: "Smart home RabbitMQ Proxy app consumes messages from the queue and posts the payload to " +
			"certain configured endpoint",
		Usage:     "Application for consuming RMQ messages and posting their payload to some endpoint",
		UsageText: "rmqproxy [global options]",
		HideHelp:  false,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       rmqHost,
				Usage:       "RabbitMQ Host",
				Destination: &rmqConfig.Host,
				Aliases:     []string{"t"},
			},
			&cli.IntFlag{
				Name:        "port",
				Value:       rmqPort,
				Usage:       "RabbitMQ Host",
				Destination: &rmqConfig.Port,
				Aliases:     []string{"p"},
			},
			&cli.StringFlag{
				Name:        "login",
				Value:       rmqLogin,
				Usage:       "RabbitMQ Login",
				Destination: &rmqConfig.Login,
				Aliases:     []string{"l"},
			},
			&cli.StringFlag{
				Name:        "password",
				Value:       rmqPassword,
				Usage:       "RabbitMQ Password",
				Destination: &rmqConfig.Password,
				Aliases:     []string{"s"},
			},
			&cli.StringFlag{
				Name:        "exchange",
				Value:       rmqExchange,
				Usage:       "RabbitMQ Exchange name",
				Destination: &rmqConfig.Exchange,
				Aliases:     []string{"e"},
			},
			&cli.StringFlag{
				Name:        "queue",
				Value:       rmqQueue,
				Usage:       "RabbitMQ Queue name",
				Destination: &rmqConfig.Queue,
				Aliases:     []string{"q"},
			},
			&cli.StringFlag{
				Name:        "rkey",
				Value:       rmqRoutingKey,
				Usage:       "RabbitMQ Queue name",
				Destination: &rmqConfig.RoutingKey,
				Aliases:     []string{"r"},
			},
			&cli.StringFlag{
				Name:        "endpoint",
				Value:       apiEndpoint,
				Usage:       "Endpoint where to post the payload using POST method",
				Destination: &msgHandler.EndPoint,
				Aliases:     []string{"u"},
			},
			&cli.StringFlag{
				Name:        "log",
				Usage:       "Log file for logs output",
				Destination: &logFile,
			},
		},
		Action: func(context *cli.Context) error {
			rmqProc := amqp.NewRmq(&rmqConfig)
			rmqProc.Consume(msgHandler)

			return nil
		},
		Before: func(context *cli.Context) error {
			if logFile != "" {
				err := setLogOutputToFile(logFile)
				if err != nil {
					log.Println("Failed to set the log output to file: ", err)
				}
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func setLogOutputToFile(fileName string) error {
	logFile, err := os.Create(fileName)

	if err != nil {
		return err
	}

	log.SetOutput(logFile)

	return nil
}
