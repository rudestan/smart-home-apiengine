package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/urfave/cli/v2"
)

func main() {
	var logFile string

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
		},
		Commands: []*cli.Command{
			{
				Name:        "create",
				Usage:       "Creates a an empty configuration JSON file",
				UsageText:   fmt.Sprintf("%s run [type: \"scenario\" or \"cmd\"] [id of the command]", execName),
				Description: "Creates an empty configuration JSON file in provided path",
				Action: func(c *cli.Context) error {
					if c.NArg() != 1 {
						return errors.New("not all required arguments provided")
					}

					configFilepath := c.Args().First()

					log.Println(configFilepath)

					return nil
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
