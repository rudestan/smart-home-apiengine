package main

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"smh-apiengine/pkg/devicecontrol"
)

var simpleTemplate = &promptui.SelectTemplates{
	Label:    "{{ . }}",
	Active:   "\U00002705 {{ . | yellow }}",
	Inactive: "  {{ . | cyan }}",
	Selected: "\U00002705 {{ . | red | cyan }}",
}

func CmdAddControls(configFile string) error {
	config, err := devicecontrol.NewConfiguration(configFile)
	if err != nil {
		return err
	}

	for {
		fmt.Println("Adding controls")
		fmt.Print("- Existing controls: -\n\n")
		for _, existingControl := range config.Controls {
			fmt.Printf(
				"\U000027A4  %s (items: %d)\n",
				existingControl.Name,
				len(existingControl.Items))
		}
		fmt.Println("----------------------")

		controlName, err := promptEnterName("control (Light control, TV, Media etc.)")
		if err != nil {
			return err
		}

		controlIcon, err := promptEnterName("control icon (e.g. light, tv etc.)")
		if err != nil {
			return err
		}

		fmt.Printf("Creating a \"%s\" control", controlName)

		control := config.NewControl(controlName, controlIcon)

		for {
			controlItemName, err := promptEnterName("control item name (e.g. Lamp, Power TV, etc.)")
			if err != nil {
				return err
			}

			controlItemIcon, err := promptEnterName("control item icon (e.g. power, stop, play, eject etc.)")
			if err != nil {
				return err
			}

			fmt.Printf("Creating a \"%s\" control item\n", controlItemName)

			controlItem := config.NewControlItem(controlItemName, controlItemIcon)

			fmt.Println("Now let's add some elements to the control item (commands, scenarios that will be executed)")

			for {
				elementType, err := selectSimplePrompt(
					"What you would like to add?",
					[]string{"Scenario", "Command", "Finish"})

				if err != nil {
					return err
				}

				if elementType == "Finish" {
					break
				}

				// adding a command
				var entity devicecontrol.Entity

				if elementType == "Command" {
					commandId, err := selectChooseCommand(&config, nil, 5)
					if err != nil {
						return err
					}

					if commandId == "Exit" {
						break
					}

					cmd := config.FindCommandByID(commandId)
					if cmd == nil {
						return errors.New("command not found")
					}

					device := config.FindDeviceById(cmd.DeviceID)
					if device == nil {
						return errors.New("device not found")
					}

					state := ""
					if device.SupportsPowerSwitch() {
						state, err = getStateOn()
						if err != nil {
							return err
						}
					}

					entity = config.NewControlItemCommandEntity(commandId, state)
				} else {
					scenarioId, err := selectChooseScenario(&config, nil, 5)

					if err != nil {
						return err
					}

					if scenarioId == "Exit" {
						break
					}

					stateOn, err := getStateOn()
					if err != nil {
						return err
					}

					entity = config.NewControlItemScenarioEntity(scenarioId, stateOn)
				}

				controlItem.AddControlItemStateEntity(entity)

			}

			control.AddControlItem(controlItem)

			answer, err := selectSimplePrompt("Would you like to add more control items?", []string{"Yes", "No"})
			if answer == "No" {
				break
			}
		}

		config.AddControl(control)

		controlAnswer, err := selectSimplePrompt(
			"What would you like to do next?",
			[]string{"Add another control", "Save and exit", "Exit"})

		if controlAnswer == "Save and exit" {
			err = config.SaveConfiguration(configFile)

			if err != nil {
				return err
			}

			return nil
		}

		if controlAnswer == "Exit" {
			return nil
		}
	}
}

func getStateOn() (string, error)  {
	state, err := selectSimplePrompt("Which type of power state is this?", []string{"Stateless", "On", "Off"})

	if err != nil {
		return "", err
	}

	switch state {
	case "On": return devicecontrol.StateOn, nil
	case "Off": return devicecontrol.StateOff, nil
	}

	return devicecontrol.StateNA, nil
}
