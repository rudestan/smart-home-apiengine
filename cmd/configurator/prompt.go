package main

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"smh-apiengine/pkg/devicecontrol"
	"strconv"
)

type commandItem struct {
	Name string
	ID string
	DeviceName string
}

type scenarioItem struct {
	ID string
	Name string
	CmdCount int
}

func promptEnterName(subject string) (string, error)  {
	validate := func(input string) error {
		if len(input) > 0 {
			return nil
		}
		return errors.New(subject + " must not be empty")
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Enter %s name", subject),
		Validate: validate,
	}

	input, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return input, nil
}

func promptEnterInt(subject string) (int, error)  {
	validate := func(input string) error {
		_, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("Invalid number")
		}

		return nil
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("Enter %s", subject),
		Validate: validate,
	}

	input, err := prompt.Run()

	if err != nil {
		return -1, err
	}

	val, err := strconv.Atoi(input)

	if err != nil {
		return -1, err
	}

	return val, nil
}

func selectSimplePrompt(label string, answers []string) (string, error)  {
	prompt := promptui.Select{
		Label: label,
		Items: answers,
		Templates: simpleTemplate,
		Size: 10,
	}

	_, choice, err := prompt.Run()

	if err != nil {
		return "", err
	}

	return choice, nil
}

func selectChooseCommand(config *devicecontrol.Config, extraItems []commandItem, size int) (string, error)  {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .Name }} (device: {{ .DeviceName | red }})",
		Active:   "\U00002705 {{ .Name | yellow }} (device: {{ .DeviceName | red }})",

		Inactive: "  {{ .Name | cyan }} (device: {{ .DeviceName | red }})",
		Selected: "\U00002705 {{ .Name | red | cyan }} (device: {{ .DeviceName | red }})",
	}

	var items []commandItem

	for _, command := range config.Commands {
		device := config.FindDeviceById(command.DeviceID)
		deviceName := "unknown"

		if device != nil {
			deviceName = device.Name
		}

		items = append(items, commandItem{
			Name:        command.Name,
			ID: command.ID,
			DeviceName:  deviceName,
		})
	}
	if extraItems != nil {
		for _, extraItem := range extraItems {
			items = append(items, extraItem)
		}
	}

	prompt := promptui.Select{
		Label: "Please select a command to add to scenario",
		Items: items,
		Templates: templates,
		Size: size,
	}

	idx, _, err := prompt.Run()

	if err != nil {
		return "", err
	}

	if len(items) > 0 {
		return items[idx].ID, nil
	}

	return "", nil
}

func selectChooseScenario(config *devicecontrol.Config, extraItems []scenarioItem, size int) (string, error)  {
	templates := &promptui.SelectTemplates{
		Label:    "{{ .Name }} {{ .ID }} (commands: {{ .CmdCount | red }})",
		Active:   "\U00002705 {{ .Name | yellow }} {{ .ID }} (commands: {{ .CmdCount | red }})",

		Inactive: "  {{ .Name | cyan }} {{ .ID }} (commands: {{ .CmdCount | red }})",
		Selected: "\U00002705 {{ .Name | red | cyan }} {{ .ID }} (commands: {{ .CmdCount | red }})",
	}

	var items []scenarioItem

	for _, scenario := range config.Scenarios {
		items = append(items, scenarioItem{
			ID: scenario.ID,
			Name:        scenario.Name,
			CmdCount: len(scenario.Sequence),
		})
	}
    if extraItems != nil {
		for _, extraItem := range extraItems {
			items = append(items, extraItem)
		}
	}

	prompt := promptui.Select{
		Label: "Please select scenario to add to control item",
		Items: items,
		Templates: templates,
		Size: size,
	}

	idx, _, err := prompt.Run()

	if err != nil {
		return "", err
	}

	if len(items) > 0 {
		return items[idx].ID, nil
	}

	return "", nil
}
