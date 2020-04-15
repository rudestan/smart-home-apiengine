package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"path"
)

func main() {
	var logFile string
	var configFile string

	execName, err := os.Executable()

	if err != nil {
		execName = "configure"
	}

	app := &cli.App{
		Name: "Smart Home Broadlink API Engine Configuration App",
		Description: "Application allows to discover Broadlink devices, configure commands and save the configuration " +
			"to JSON file. This configuration file can be used to run the commands on corresponding Broadlink devices " +
			"(e.g. RMP3 Pro, S2, S3, SC1 etc.) as well as run web server app to be able to server the " +
			"incoming requests and execute the commands.",
		Usage:     "an app for configuring Broadlink devices",
		UsageText: fmt.Sprintf("%s [global options] command [command options]", path.Base(execName)),
		HideHelp:  false,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log",
				Usage:       "Log file for logs output",
				Destination: &logFile,
				Aliases:     []string{"l"},
			},
			&cli.PathFlag{
				Name:        "config",
				Usage:       "Path to JSON configuration with commands and devices",
				Destination: &configFile,
				Aliases:     []string{"c"},
				EnvVars:	 []string{"SMH_CONFIG"},
				Required:    true,
			},
		},
		Commands: []*cli.Command{
			{
				Name:        "discover",
				Usage:       "Discovers the devices",
				Action: func(c *cli.Context) error {
					return CmdDiscover(configFile)
				},
			},
			{
				Name:        "add_commands",
				Usage:       "Adds commands",
				Action: func(c *cli.Context) error {
					return CmdAddCommands(configFile)
				},
			},
			{
				Name:        "add_scenarios",
				Usage:       "Adds scenarios",
				Action: func(c *cli.Context) error {
					return CmdAddScenarios(configFile)
				},
			},
			{
				Name:        "add_controls",
				Usage:       "Adds controls",
				Action: func(c *cli.Context) error {
					return CmdAddControls(configFile)
				},
			},
			{
				Name:        "run",
				Usage:       "Runs command or scenario",
				Action: func(c *cli.Context) error {
					return CmdRun(configFile)
				},
			},
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
