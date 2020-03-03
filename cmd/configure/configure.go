package main

import (
    "errors"
    "fmt"
    "log"
    "os"
    "smh-apiengine/pkg/devicecontrol"

    "github.com/urfave/cli/v2"
)

func main() {
    var configFile string
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
        Usage: "an app for configuring and controlling Broadlink devices",
        UsageText: fmt.Sprintf("%s [global options] command [command options]", execName),
        HideHelp: false,
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
                Aliases: []string{"l"},
            },
        },
        Commands: []*cli.Command{
            {
                Name:    "run",
                Usage:   "Runs command or scenario with provided id",
                UsageText: fmt.Sprintf("%s run [type: \"scenario\" or \"cmd\"] [id of the command]", execName),
                Description: "Runs command or scenario with provided id",
                Action: func(c *cli.Context) error {
                    if c.NArg() != 2 {
                        return errors.New("not all required arguments provided")
                    }

                    runType := c.Args().First()

                    if runType != "cmd" && runType != "scenario" {
                        return errors.New("invalid run type specified. Must be either \"scenario\" or \"cmd\"")
                    }

                    err := devicecontrol.Init(configFile)
                    if err != nil {
                        return err
                    }

                    id := c.Args().Get(1)

                    switch runType {
                    case "cmd":
                        cmd, err := devicecontrol.FindCommandById(id)
                        if err != nil {
                            return err
                        }

                        return devicecontrol.ExecCommandFullCycle(cmd)
                    case "scenario":
                        scenario, err := devicecontrol.FindScenarioByName(id)
                        if err != nil {
                            return err
                        }

                        return devicecontrol.ExecScenarioFullCycle(scenario)
                    }

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

func setLogOutputToFile(fileName string) error  {
    logFile, err := os.Create(fileName)
    if err != nil {
        return err
    }

    log.SetOutput(logFile)

    return nil
}
