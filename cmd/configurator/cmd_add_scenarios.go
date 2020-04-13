package main

import (
	"fmt"
	"smh-apiengine/pkg/devicecontrol"
)

func CmdAddScenarios(configFile string) error  {
	config, err := devicecontrol.NewConfiguration(configFile)

	if err != nil {
		return err
	}

	deviceControl := devicecontrol.NewDeviceControl(&config)

	for {
		fmt.Println("Adding scenarios")
		fmt.Print("- Existing scenarios: -\n\n")
		for _, existingScenario := range config.Scenarios {
			fmt.Printf(
				"\U000027A4  %s (sequence items: %d)\n",
				existingScenario.Name,
				len(existingScenario.Sequence))
		}
		fmt.Println("----------------------")

		choice, err := promptEnterName("scenario")
		if err != nil {
			return err
		}

		fmt.Printf("Creating a \"%s\" scenario", choice)

		scenario := deviceControl.NewScenario(choice)

		for {
			cmdId, err := selectChooseCommand(&config, []commandItem{{
				Name:       "Finish adding",
				ID:         "Exit",
				DeviceName: "-",
			}}, 5)

			if err != nil {
				return err
			}

			if cmdId == "Exit" {
				break
			}

			delay, err := promptEnterInt("delay in seconds after command")
			if err != nil {
				return err
			}

			sequenceItem := deviceControl.NewSequenceItem(cmdId, delay)
			scenario.AddSequenceItem(sequenceItem)

			fmt.Println("Command added")
		}

		deviceControl.AddScenario(scenario)
		fmt.Println("Scenario added")

		saveOrExit, err := selectSimplePrompt(
			"Create another scenario or exit?",
			[]string{"Create New", "Save and exit", "Exit"})

		if err != nil {
			return err
		}

		if saveOrExit == "Exit" {
			return nil
		}

		if saveOrExit == "Save and exit" {
			err = config.SaveConfiguration(configFile)

			if err != nil {
				return err
			}

			return nil
		}
	}
}
