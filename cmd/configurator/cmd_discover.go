package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"smh-apiengine/pkg/devicecontrol"

	"github.com/manifoldco/promptui"
	"github.com/rudestan/broadlinkrm"
)

type deviceItem struct {
	Name string
	Ip string
	Mac string
	Status string
}

func CmdDiscover(configFile string) error {
	config, err := devicecontrol.NewConfiguration(configFile)

	if err != nil {
		if os.IsNotExist(err) {
			createConfig := promptCreateConfig()

			if !createConfig {
				return nil
			}

			err = config.SaveConfiguration(configFile)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	deviceControl := devicecontrol.NewDeviceControl(&config)

	fmt.Println("Discovering, please wait...")

	err = deviceControl.Discover(false)
	if err != nil {
		return err
	}

	devicesList := deviceControl.GetDiscoveredDevices()
	if len(devicesList) == 0 {
		fmt.Println("No devices found!")
		return nil
	}

	for {
		deviceInfo := selectDiscoveredDeviceAdd(&deviceControl, devicesList)

		switch deviceInfo.Name {
		case "Exit": return nil
		case "Save changes":
			err = config.SaveConfiguration(configFile)
			if err != nil {
				log.Println("Failed to save the config!")
			} else {
				log.Println("Changes saved")
			}
			break
		default:
			if deviceInfo.Mac == "-" {
				return errors.New("unknown choice")
			}

			deviceName := promptAddDevice(deviceInfo.Name)
			err = deviceControl.AddOrUpdateDiscoveredDevice(deviceName, deviceInfo.Mac)
			if err != nil {
				log.Println("Failed to add/update device!")
			} else {
				log.Printf("Device \"%s\" added!\n", deviceName)
			}

			break
		}
	}
}

func selectDevicePrompt(devicesList map[string]*devicecontrol.Device) string  {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U00002705 {{ .Name | yellow }} ({{ .Ip | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Ip | red }})",
		Selected: "\U00002705 {{ .Name | red | cyan }}",
	}

	var items []deviceItem

	for _, deviceInfo := range devicesList {
		items = append(items, deviceItem{
			Name:   deviceInfo.Name,
			Ip:     deviceInfo.IP,
			Mac:     deviceInfo.Mac,
		})
	}
	items = append(items, deviceItem{Name:   "Exit"})

	prompt := promptui.Select{
		Label: "Please select the device for which you want to add a command?",
		Items: items,
		Templates: templates,
		Size: 5,
	}

	idx, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	selectedDevice := items[idx]

	return selectedDevice.Mac
}

func promptAddDevice(defaultName string) string  {
	validate := func(input string) error {
		if len(input) > 0 {
			return nil
		}
		return errors.New("device name must not be empty")
	}

	prompt := promptui.Prompt{
		Label:     "Enter device name",
		Default:   defaultName,
		Validate:validate,
	}

	choice, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return choice
}

func selectDiscoveredDeviceAdd(deviceControl *devicecontrol.DeviceControl, devicesList map[string]broadlinkrm.DeviceInfo) deviceItem  {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U00002705 {{ .Name | yellow }} ({{ .Ip | red }}) - {{ .Status | green }}",
		Inactive: "  {{ .Name | cyan }} ({{ .Ip | red }}) - {{ .Status | green }}",
		Selected: "\U00002705 {{ .Name | red | cyan }}",
	}

	var items []deviceItem

	for _, deviceInfo := range devicesList {
		deviceStatus := deviceControl.AnalyzeDevice(deviceInfo)
		items = append(items, deviceItem{
			Name:   deviceInfo.Name,
			Ip:     deviceInfo.Ip,
			Mac:     deviceInfo.Mac,
			Status: deviceStatus,
		})
	}
	items = append(items, deviceItem{
		Name:   "Exit",
		Ip:     "",
		Mac:    "",
		Status: "",
	})
	items = append(items, deviceItem{
		Name:   "Save changes",
		Ip:     "",
		Mac:    "",
		Status: "",
	})

	prompt := promptui.Select{
		Label: "Please select which device to add/update?",
		Items: items,
		Templates: templates,
		Size: 5,
	}

	idx, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return deviceItem{}
	}

	return items[idx]
}

func promptCreateConfig() bool  {
	prompt := promptui.Select{
		Label: "Config file does not exist, should we create a new one?",
		Items: []string{"Yes", "No"},
	}

	_, answer, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return false
	}

	if answer == "No" {
		return false
	}

	return true
}
