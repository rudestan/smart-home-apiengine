package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"smh-apiengine/pkg/devicecontrol"

	"github.com/urfave/cli/v2"
)

var (
	ErrArguments = errors.New("not all required arguments provided")
	ErrInvalidRunType = errors.New("invalid run type specified. Must be either \"scenario\" or \"cmd\"")
)

func main() {
	var configFile string
	var logFile string

	execName, err := os.Executable()

	if err != nil {
		execName = "run"
	}

	app := &cli.App{
		Name:        "Smart Home Broadlink API Engine Runner App",
		Description: "Application runs commands and scenarios from the configuration JSON file",
		Usage:       "an app for running commands and scenarios on Broadlink devices",
		UsageText:   fmt.Sprintf("%s [global options] [type: \"scenario\" or \"cmd\"] [id of the command]", path.Base(execName)),
		HideHelp:    false,
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
				EnvVars:	 []string{"SMH_RUNNER_LOG_FILE"},
			},
		},
		Action: func(c *cli.Context) error {
			if c.NArg() != 2 {
				return ErrArguments
			}

			runType := c.Args().First()

			if runType != "cmd" && runType != "scenario" {
				return ErrInvalidRunType
			}

			config, err := devicecontrol.NewConfiguration(configFile)

			if err != nil {
				return err
			}

			deviceControl := devicecontrol.NewDeviceControl(&config)
			id := c.Args().Get(1)

			switch runType {
			case "cmd":
				cmd := deviceControl.FindCommandByID(id)
				if cmd == nil {
					return errors.New("command not found")
				}

				return deviceControl.ExecCommandFullCycle(*cmd)
			case "scenario":
				scenario, err := deviceControl.FindScenarioByName(id)
				if err != nil {
					return err
				}

				return deviceControl.ExecScenarioFullCycle(scenario)
			}

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

	err = app.Run(os.Args)
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
