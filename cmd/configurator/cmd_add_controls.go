package main

import (
	"errors"
	"fmt"
	"smh-apiengine/pkg/devicecontrol"

	"github.com/manifoldco/promptui"
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

			controlItemStateful := false
			controlItemStatefulAnswer, err := selectSimplePrompt(
				"Does control item support sates (all commands must use the same device)?",
				[]string{"Yes", "No"})
			if err != nil {
				return err
			}

			if controlItemStatefulAnswer == "Yes" {
				controlItemStateful = true
			}

			fmt.Printf("Creating a \"%s\" control item\n", controlItemName)

			controlItem := config.NewControlItem(controlItemName, controlItemIcon, controlItemStateful)

			fmt.Println("Now let's add some elements to the control item (commands, scenarios that will be executed)")

			for {
				elementType := "Command"

				if (!controlItemStateful) {
					elementType, err = selectSimplePrompt(
						"What you would like to add?",
						[]string{"Scenario", "Command"})

					if err != nil {
						return err
					}
				}

				// adding a command
				var element devicecontrol.Element

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

					state := devicecontrol.StateOn

					if controlItemStateful && device.SupportsPowerSwitch() {
						state, err = getState()
						if err != nil {
							return err
						}
					}

					element = config.NewControlItemCommandElement(commandId, state)
				} else {
					scenarioId, err := selectChooseScenario(&config, nil, 5)

					if err != nil {
						return err
					}

					if scenarioId == "Exit" {
						break
					}

					element = config.NewControlItemScenarioElement(scenarioId, devicecontrol.StateOn)
				}

				controlItem.AddControlItemElement(element)

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

func getState() (string, error)  {
	powerState, err := selectSimplePrompt(
		"Which type of power state is this?",
		[]string{
			devicecontrol.StateOn,
			devicecontrol.StateOff})

	if err != nil {
		return devicecontrol.StateOff, err
	}

	return powerState, nil
}
