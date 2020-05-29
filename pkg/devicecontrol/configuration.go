package devicecontrol

import (
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	DevicePowerSwitch = "power_switch"
	DeviceBlaster = "blaster"
)

const (
	PowerSwitchOnCmdName = "on"
	PowerSwitchOffCmdName = "off"
	PowerSwitchOnCmd = "01"
	PowerSwitchOffCmd = "00"
)

const (
	StateOn = "on"
	StateOff = "off"
	StateNA = "na"
)

const (
	NsUUIDCommand = "smart-home:command:"
	NsUUIDScenario = "smart-home:scenario:"
	NsUUIDElement = "smart-home:element:"
	NsUUIDControl = "smart-home:control:"
)

const (
	ElementTypeCommand = "command"
	ElementTypeScenario = "scenario"
)

// Device struct that stores all required for usage data
type Device struct {
	Name       string `json:"name"`
	IP         string `json:"ip"`
	Mac        string `json:"mac"`
	Key        string `json:"key"`
	ID         string `json:"id"`
	DeviceType string `json:"device_type"`
	DeviceCategory string `json:"device_category"`
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
	ID 		 string 		 `json:"id"`
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
	CommandId string `json:"command_id"`
	Delay int    `json:"delay"`
}

// Scenario struct contains data that allows to execute multiple commands sequence using some triggering intents
// scenario can be used as one control item that triggers N command sequence with certain delay
type Scenario struct {
	ID 		 string 		 `json:"id"`
	Name     string          `json:"name"`
	Sequence []SequenceItem  `json:"sequence"`
	Intents  []CommandIntent `json:"intents"`
}

// Control struct is used for creating a virtual remote control with items (buttons)
type Control struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"` // icon for control e.g. in tabs
	Items map[string]*ControlItem `json:"items"`
}

// ControlItem struct represents some virtual control button
type ControlItem struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	StateEntities []Entity `json:"state_entities"`
	activeState string
}

// Entity struct that represents an element from the configuration, either command or scenario
type Entity struct {
	ID string `json:"id"`
	Target string `json:"target"`
	Type string `json:"type"`
	State string `json:"state"`
}

// Config struct is the root struct that defines a device control struct
type Config struct {
	Devices   map[string]*Device  `json:"devices"`
	Intents   map[string]Intent   `json:"intents"`
	Commands  map[string]Command  `json:"commands"`
	Scenarios map[string]Scenario `json:"scenarios"`
	Controls  map[string]Control  `json:"controls"`
	Schedule  map[string]ScheduleItem `json:"schedule"`
	fileName  string
	sync.Mutex
}

// ScheduleItem struct represents the schedule item that can be executed at certain times or some interval etc.
type ScheduleItem struct {
	ExecutionTimes map[string]string `json:"execution_times"`
	Entity Entity `json:"entity"`
}

func (c *Config) AddControl(ctrl Control)  {
	if c.Controls == nil {
		c.Controls = make(map[string]Control)
	}

	c.Controls[ctrl.ID] = ctrl
}

func (c *Config) NewControl(name string, icon string) Control  {
	controlId := c.getUUIDV5(NsUUIDControl, name)

	return Control{
		ID:    controlId.String(),
		Name:  name,
		Icon:  icon,
		Items: nil,
	}
}

func (ctrl *Control) AddControlItem(ci *ControlItem)  {
	if ctrl.Items == nil {
		ctrl.Items = make(map[string]*ControlItem)
	}

	ctrl.Items[ci.ID] = ci
}

func (ci *ControlItem) AddControlItemStateEntity(et Entity)  {
	ci.StateEntities = append(ci.StateEntities, et)
}

func (ci *ControlItem) FindEntityByState(state string) *Entity {
	for _, entity := range ci.StateEntities {
		if entity.State == state {
			return &entity
		}
	}

	return nil
}

func (ci *ControlItem) FindNextStateEntity() *Entity {
	if len(ci.StateEntities) == 0 {
		return nil
	}

	if ci.activeState == "" || ci.activeState == "na" {
		return &ci.StateEntities[0]
	}

	for idx, entity := range ci.StateEntities {
		if ci.activeState == entity.State {
			if idx+1 < len(ci.StateEntities) {
				return &ci.StateEntities[idx+1]
			} else if idx+1 >= len(ci.StateEntities) {
				return &ci.StateEntities[0]
			}
		}
	}

	return nil
}

func (c *Config) NewControlItem(name string, icon string) *ControlItem  {
	elementId := uuid.NewV4()

	return &ControlItem{
		ID:       elementId.String(),
		Name:     name,
		Icon:     icon,
		StateEntities: nil,
	}
}

func (c *Config) NewControlItemCommandEntity(target string, state string) Entity  {
	return c.newControlItemStateEntity(target, ElementTypeCommand, state)
}

func (c *Config) NewControlItemScenarioEntity(target string, state string) Entity  {
	return c.newControlItemStateEntity(target, ElementTypeScenario, state)
}

func (c *Config) getUUIDV5(ns string, name string) uuid.UUID  {
	nsUUID := uuid.NewV5(uuid.UUID{}, ns)

	return uuid.NewV5(nsUUID, name)
}

func (c *Config) newControlItemStateEntity(target string, elementType string, state string) Entity  {
	elementId := uuid.NewV4()

	return Entity{
		ID:      elementId.String(),
		Target:  target,
		Type:    elementType,
		State: 	 state,
	}
}

func (d *Device) SupportsPowerSwitch() bool  {
	return d.DeviceCategory == DevicePowerSwitch
}

func (c *Config) findDeviceByMac(deviceMac string) (*Device, error) {
	for _, device := range c.Devices {
		if device.Mac == deviceMac {
			return device, nil
		}
	}

	return nil, errors.New("no device found")
}

func (c *Config) FindDeviceById(id string) *Device {
	if device, ok := c.Devices[id]; ok {
		return device
	}

	return nil
}

func (c *Config) FindCommandByID(id string) *Command {
	if cmd, ok := c.Commands[id]; ok {
		return &cmd
	}

	return nil
}

func (c *Config) FindControlItemByID(id string) *ControlItem {
	for _, control := range c.Controls {
		if controlItem, ok := control.Items[id]; ok {
			return controlItem
		}
	}

	return nil
}

func (c *Config) FindScenarioByID(id string) *Scenario {
	if scenario, ok := c.Scenarios[id]; ok {
		return &scenario
	}

	return nil
}

func (c *Config) findScenarioByName(name string) (Scenario, error) {
	if cmd, ok := c.Scenarios[name]; ok {
		return cmd, nil
	}

	return Scenario{}, fmt.Errorf("scenario \"%s\" not found", name)
}

func (s *Scenario) AddSequenceItem(item SequenceItem)  {
	s.Sequence = append(s.Sequence, item)
}

// NewConfiguration loads the json configuration from provided file
func NewConfiguration(fileName string) (Config, error) {
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

	var config Config

	err = json.Unmarshal(contents, &config)

	if err != nil {
		return Config{}, err
	}

	config.fileName = fileName

	return config, nil
}

// saveConfiguration saves the configuration to the provided filename
func (c *Config) SaveConfiguration(fileName string) error {
	c.Lock()
	defer c.Unlock()

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

	data, err := json.MarshalIndent(c, "", "    ")

	if err != nil {
		return errors.New("failed to save config")
	}

	err = ioutil.WriteFile(fileName, data, mode)

	return err
}
