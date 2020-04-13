package main

import (
	"errors"
	"smh-apiengine/pkg/devicecontrol"
)

func CmdRun(configFile string) error {
	config, err := devicecontrol.NewConfiguration(configFile)
	if err != nil {
		return err
	}

	deviceControl := devicecontrol.NewDeviceControl(&config)

	for {
		elementType, err := selectSimplePrompt(
			"What do you want to lunch?",
			[]string{"Scenario", "Command", "Exit"})
		if err != nil {
			return err
		}

		if elementType == "Exit" {
			return nil
		}

		for {
			var elementId string

			if elementType == "Command" {
				elementId, err = selectChooseCommand(
					&config,
					[]commandItem{
						{"\U0001F448 Back", "Back", "-"},
						{"\U0001F5A5 Exit", "Exit", "-"}},
					15)
			} else {
				elementId, err = selectChooseScenario(
					&config,
					[]scenarioItem{
						{ID: "Back", Name: "\U0001F448", CmdCount: 0},
						{ID: "Exit", Name: "\U0001F5A5", CmdCount: 0}},
					15)
			}
			if err != nil {
				return err
			}

			if elementId == "Back" {
				break
			}

			if elementId == "Exit" {
				return nil
			}

			if elementType == "Command" {
				err = execCommandById(&deviceControl, &config, elementId)
			} else {
				err = execScenarioById(&deviceControl, &config, elementId)
			}

			if err != nil {
				return err
			}
		}
	}
}

func execCommandById(dc *devicecontrol.DeviceControl, config *devicecontrol.Config, cmdId string) error {
	command := config.FindCommandByID(cmdId)

	if command == nil {
		return errors.New("no command found")
	}

	return dc.ExecCommand(command)
}

func execScenarioById(dc *devicecontrol.DeviceControl, config *devicecontrol.Config, scId string) error {
	scenario := config.FindScenarioByID(scId)

	if scenario == nil {
		return errors.New("no scenario found")
	}

	return dc.ExecScenario(scenario)
}
