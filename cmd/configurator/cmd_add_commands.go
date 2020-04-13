package main

import (
	"errors"
	"fmt"
	"smh-apiengine/pkg/devicecontrol"
)

func CmdAddCommands(configFile string) error {
	config, err := devicecontrol.NewConfiguration(configFile)

	if err != nil {
		return err
	}

	deviceControl := devicecontrol.NewDeviceControl(&config)
	deviceMac := selectDevicePrompt(deviceControl.GetDevices())
	device := config.FindDeviceById(deviceMac)

	if device == nil {
		return errors.New("device not found")
	}

	for {
		fmt.Printf("Adding a command for the device [%s]\n", device.Name)
		fmt.Print("- Existing commands: -\n\n")
		for _, existingCmd := range config.Commands {
			if existingCmd.DeviceID != device.Mac {
				continue
			}

			fmt.Println("\U000027A4  " + existingCmd.Name)
		}
		fmt.Println("----------------------")

		var cmd devicecontrol.Command

		if device.DeviceCategory == devicecontrol.DevicePowerSwitch {
			powerSwitchCommand, err := selectSimplePrompt(
				"Please select command type for Power Switch",
				[]string{"On", "Off", "Exit"})

			if err != nil {
				return err
			}

			if powerSwitchCommand == "Exit" {
				return nil
			}
			cmdName, err := promptEnterName("command")

			if err != nil {
				return err
			}

			cmd = deviceControl.NewCommandForPowerSwitch(device, cmdName, powerSwitchCommand)
		} else {
			startLearning, err := selectSimplePrompt(
				"Enter learning mode for the blaster device? Select Yes once you are ready.",
				[]string{"Yes", "Exit"})
			if err != nil {
				return err
			}

			if startLearning == "Exit" {
				return nil
			}

			fmt.Println("Point your remote control to the device and press the button that you want to learn.")

			cmdCode, err := deviceControl.LearnCommand(device.Mac)
			if err != nil {
				return err
			}

			cmdName, err := promptEnterName("command")

			if err != nil {
				return err
			}

			cmd = deviceControl.NewCommand(device, cmdName, cmdCode)
		}

		deviceControl.AddCommand(cmd)
		fmt.Println("Command added")

		addNew, err := selectSimplePrompt("Add new command?", []string{"Yes", "Exit", "Exit and save"})

		if err != nil {
			return err
		}

		if addNew == "Exit and save" {
			err = config.SaveConfiguration(configFile)

			if err != nil {
				return err
			}
			fmt.Println("Configuration saved")

			return nil
		}

		if addNew == "Exit" {
			return nil
		}
	}
}
