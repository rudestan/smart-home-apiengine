package main

import (
    "apiengine/apiserver"
    "apiengine/devicecontrol"
    "errors"
    "fmt"
    "github.com/urfave/cli/v2"
    "log"
    "os"
)

const (
    defaultProtocol = "http"
    defaultAddress = "127.0.0.1"
    defaultPort = 8787
)

type serverConfig struct {
    Protocol string
    Address string
    Port int
    TLSCert string
    TLSKey string
}

func main() {
    var configFile string
    var logFile string
    var srvConfig serverConfig

    app := &cli.App{
        Name: "Broadlink API Engine",
        Description: "A Web server that serves an incoming requests and controls Broadlink devices using API. Uses a " +
            "user created JSON for pre-configured device commands, scenarios and devices. Supports an incoming requests" +
            "from Alexa API web servers.",
        Usage: "an app and web server for controlling Broadlink devices",
        UsageText: "apiengine [global options] command [command options]",
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
                Name:    "serve",
                Usage:   "Starts web server",
                UsageText: "serve [options]",
                Description: fmt.Sprintf("Starts web server using default parameters (protocol: %d, host:%s:%d) " +
                    "or from provided arguemnts", defaultPort, defaultAddress, defaultPort),
                Flags: []cli.Flag {
                    &cli.StringFlag{
                        Name:        "proto",
                        Value:       defaultProtocol,
                        Usage:       "Protocol for web server to run (values: \"http\", \"https\")",
                        Destination: &srvConfig.Protocol,
                        Aliases: []string{"pr"},
                    },
                    &cli.StringFlag{
                        Name:        "address",
                        Value:       defaultAddress,
                        Usage:       "Ip address for web server",
                        Destination: &srvConfig.Address,
                        Aliases: []string{"a"},
                    },
                    &cli.IntFlag{
                        Name:        "port",
                        Value:       defaultPort,
                        Usage:       "Port for web server",
                        Destination: &srvConfig.Port,
                        Aliases: []string{"p"},
                    },
                    &cli.StringFlag{
                        Name:        "tls-cert",
                        Usage:       "TLS Certificate file path (only in case https protocol is used)",
                        Destination: &srvConfig.TLSCert,
                        Aliases: []string{"tc"},
                    },
                    &cli.StringFlag{
                        Name:        "tls-key",
                        Usage:       "TLS Key file path (only in case https protocol is used)",
                        Destination: &srvConfig.TLSKey,
                        Aliases: []string{"tk"},
                    },
                },
                Action: func(c *cli.Context) error {
                    err := devicecontrol.Init(configFile)
                    if err != nil {
                        return err
                    }

                    return runServer(srvConfig)
                },
            },
            {
                Name:    "run",
                Usage:   "Runs command or scenario with provided id",
                UsageText: "apiengine run [type: \"scenario\" or \"cmd\"] [id of the command]",
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

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}

func runServer(serverConfig serverConfig) error  {
    if serverConfig.Protocol == "https" {
        if serverConfig.TLSCert == "" || serverConfig.TLSKey == "" {
            return errors.New("TLS Certificate and Key files are required when using https protocol")
        }
    }

    switch serverConfig.Protocol {
    case "http":
        apiserver.ServeHttp(fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port))
    case "https":
        apiserver.ServeHttps(
            fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port),
            serverConfig.TLSCert,
            serverConfig.TLSKey)
    }

    return nil
}

func setLogOutputToFile(fileName string) error  {
    logFile, err := os.Create(fileName)
    if err != nil {
        return err
    }

    log.SetOutput(logFile)

    return nil
}
