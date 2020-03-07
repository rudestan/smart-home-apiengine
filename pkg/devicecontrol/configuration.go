package devicecontrol

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Device struct that stores all required for usage data
type Device struct {
	Name       string `json:"name"`
	IP         string `json:"ip"`
	Mac        string `json:"mac"`
	Key        string `json:"key"`
	ID         string `json:"id"`
	DeviceType string `json:"device_type"`
	Enabled    bool   `json:"enabled"`
}

// Intent struct contains Alexa Intent's name and related slots
type Intent struct {
	Name  string          `json:"name"`
	Slots map[string]Slot `json:"slots"`
}

// Slot struct contains name of the slot and possible values
type Slot struct {
	Name   string               `json:"name"`
	Values map[string]SlotValue `json:"values"`
}

// SlotValue struct contains the name of the slot value and an array of possible synonyms
type SlotValue struct {
	Name     string   `json:"name"`
	Synonyms []string `json:"synonyms"`
}

// Command struct contains all the data required for the execution as well as related intents that can trigger it
type Command struct {
	DeviceID string          `json:"device_id"`
	Name     string          `json:"name"`
	Code     string          `json:"code"`
	Intents  []CommandIntent `json:"intents"`
}

// CommandIntent simplified version of Intent that contains the name of intent and slots array
type CommandIntent struct {
	Name  string                 `json:"name"`
	Slots map[string]CommandSlot `json:"slots"`
}

// CommandSlot simplified version of Slot that does not contains synonyms
type CommandSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// SequenceItem struct contains the command name as a reference to a command and delay to the next execution in seconds
type SequenceItem struct {
	Name  string `json:"name"`
	Delay int    `json:"delay"`
}

// Scenario struct contains data that allows to execute multiple commands sequence using some triggering intents
type Scenario struct {
	Name     string          `json:"name"`
	Sequence []SequenceItem  `json:"sequence"`
	Intents  []CommandIntent `json:"intents"`
}

// Group struct holds a collection of ids of control items such as command and/or scenario. Groups can be used to
// organize items for example "Light in the Living room" with commands to all light devices in the living room,
// "TV" with commands related to TV etc.
type Group struct {
	Name      string              `json:"name"`
	Commands  []map[string]string `json:"commands"`
	Scenarios []string            `json:"scenarios"`
}

// Item a struct that represents control item
type Item struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	State string `json:"state"`
} 

// Control struct is used for organising commands and scenarios in control groups (e.g. remote control)
type Control struct {
	Name string `json:"name"`
	Items map[string][]Item `json:"items"`
} 

// Config struct is the root struct that defines a device control struct
type Config struct {
	Devices   map[string]*Device  `json:"devices"`
	Intents   map[string]Intent   `json:"intents"`
	Commands  map[string]Command  `json:"commands"`
	Scenarios map[string]Scenario `json:"scenarios"`
	Groups    map[string]Group    `json:"groups"`
	Controls  map[string]Control  `json:"controls"`
	fileName  string
}

func (c *Config) findDeviceByMac(deviceMac string) (*Device, error) {
	for _, device := range c.Devices {
		if device.Mac == deviceMac {
			return device, nil
		}
	}

	return nil, errors.New("no device found")
}

func (c *Config) findCommandByID(id string) (Command, error) {
	if cmd, ok := c.Commands[id]; ok {
		return cmd, nil
	}

	return Command{}, fmt.Errorf("command \"%s\" not found", id)
}

func (c *Config) findScenarioByName(name string) (Scenario, error) {
	if cmd, ok := c.Scenarios[name]; ok {
		return cmd, nil
	}

	return Scenario{}, fmt.Errorf("scenario \"%s\" not found", name)
}

// LoadConfiguration loads the json configuration from provided file
func LoadConfiguration(fileName string) (Config, error) {
	jsonFile, err := os.Open(fileName)

	if err != nil {
		return Config{}, err
	}

	defer func() {
		err := jsonFile.Close()

		if err != nil {
			log.Println("Unable to close the config JSON file")
		}
	}()

	contents, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return Config{}, err
	}

	err = json.Unmarshal(contents, &config)

	if err != nil {
		return Config{}, err
	}

	config.fileName = fileName

	return config, nil
}

// SaveConfiguration saves the configuration to the provided filename
func SaveConfiguration(fileName string) error {
	fileInfo, err := os.Stat(fileName)
	mode := os.FileMode(0666)
	var jsonFile *os.File

	if err != nil {
		if os.IsNotExist(err) {
			jsonFile, err = os.Create(fileName)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		mode = fileInfo.Mode()
		jsonFile, err = os.Open(fileName)

		if err != nil {
			return fmt.Errorf("can not read file \"%s\"", fileName)
		}
	}

	defer func() {
		err := jsonFile.Close()

		if err != nil {
			log.Println("Unable to close the config JSON file")
		}
	}()

	data, err := json.MarshalIndent(config, "", "    ")

	if err != nil {
		return errors.New("failed to save config")
	}

	err = ioutil.WriteFile(fileName, data, mode)

	return err
}
