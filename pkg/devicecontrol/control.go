package devicecontrol

import (
	"github.com/satori/go.uuid"
	"github.com/spf13/cast"
	"log"
	"smh-apiengine/pkg/broadlinkrm"
	"strings"
)

type DeviceControl struct {
	config *Config
	broadlink broadlinkrm.Broadlink
	lock *spinLock
	DeviceStateBuffer map[string]DeviceState
}

func NewDeviceControl(config *Config) DeviceControl  {
	deviceControl := DeviceControl{
		config:    config,
		broadlink: broadlinkrm.NewBroadlink(),
		lock: 	   &spinLock{},
		DeviceStateBuffer: make(map[string]DeviceState),
	}

	if len(config.Devices) > 0 {
		deviceControl.initDevices()
	}

	return deviceControl
}

func (deviceControl *DeviceControl) AddCommand(cmd Command) {
	if deviceControl.config.Commands == nil {
		deviceControl.config.Commands = make(map[string]Command)
	}

	deviceControl.config.Commands[cmd.ID] = cmd
}

func (deviceControl *DeviceControl) NewCommandForPowerSwitch(device *Device, name string, commandType string) Command {
	if strings.ToLower(commandType) == PowerSwitchOnCmdName {
		return deviceControl.NewCommand(device, name, PowerSwitchOnCmd)
	}

	return deviceControl.NewCommand(device, name, PowerSwitchOffCmd)
}

func (deviceControl *DeviceControl) NewCommand(device *Device, name string, code string) Command {
	cmdUUID := deviceControl.getUUIDV5(NsUUIDCommand + device.Mac, name)

	return Command{
		ID: 	  cmdUUID.String(),
		DeviceID: device.Mac,
		Name:     name,
		Code:     code,
		Intents:  nil,
	}
}

func (deviceControl *DeviceControl) NewSequenceItem(cmdId string, delay int) SequenceItem {
	return SequenceItem{
		CommandId: cmdId,
		Delay:     delay,
	}
}

func (deviceControl *DeviceControl) NewScenario(name string) Scenario {
	scenarioUUID := deviceControl.getUUIDV5(NsUUIDScenario, name)

	return Scenario{
		ID:       scenarioUUID.String(),
		Name:     name,
		Sequence: nil,
		Intents:  nil,
	}
}

func (deviceControl *DeviceControl) AddScenario(scenario Scenario) {
	if deviceControl.config.Scenarios == nil {
		deviceControl.config.Scenarios = make(map[string]Scenario)
	}

	deviceControl.config.Scenarios[scenario.ID] = scenario
}

func (deviceControl *DeviceControl) getUUIDV5(ns string, name string) uuid.UUID  {
	nsUUID := uuid.NewV5(uuid.UUID{}, ns)

	return uuid.NewV5(nsUUID, name)
}

func (deviceControl *DeviceControl) LearnCommand(deviceMac string) (string, error)  {
	return deviceControl.broadlink.Learn(deviceMac)
}

func (deviceControl *DeviceControl) GetDiscoveredDevices() map[string]broadlinkrm.DeviceInfo {
	return deviceControl.broadlink.GetDeviceInfoList()
}

func (deviceControl *DeviceControl) IsKnownDevice(deviceInfo broadlinkrm.DeviceInfo) bool {
	if _, ok := deviceControl.config.Devices[deviceInfo.Mac]; !ok {
		return false
	}

	return true
}

func (deviceControl *DeviceControl) AnalyzeDevice(deviceInfo broadlinkrm.DeviceInfo) string {
	if !deviceControl.IsKnownDevice(deviceInfo) {
		return "new device"
	}

	device := deviceControl.config.Devices[deviceInfo.Mac]

	if device.IP != deviceInfo.Ip {
		return "IP not matching"
	}

	return "existing [" + device.Name + "]"
}

func (deviceControl *DeviceControl) GetDevices() map[string]*Device {
	return deviceControl.config.Devices
}

// FindCommandByID finds Command structure by provided id or error if there is no Command found
func (deviceControl *DeviceControl) FindCommandByID(id string) *Command {
	return deviceControl.config.FindCommandByID(id)
}

// FindCommandByID finds Command structure by provided id or error if there is no Command found
func (deviceControl *DeviceControl) FindControlItemByID(id string) *ControlItem {
	return deviceControl.config.FindControlItemByID(id)
}

// FindScenarioByName finds Scenario structure by provided name or error if there is no Scenario found
func (deviceControl *DeviceControl) FindScenarioByName(name string) (Scenario, error) {
	return deviceControl.config.findScenarioByName(name)
}

// AllControls returns controls from config
func (deviceControl *DeviceControl) AllControls() map[string]Control {
	return deviceControl.config.Controls
}

func (deviceControl *DeviceControl) initDevices() {
	for _, deviceConfig := range deviceControl.config.Devices {
		if !deviceConfig.Enabled {
			log.Printf("The device with ip: %s is disbled. Skipping\n", deviceConfig.IP)
			continue
		}

		err := deviceControl.broadlink.AddManualDevice(
			deviceConfig.IP,
			deviceConfig.Mac,
			deviceConfig.Key,
			deviceConfig.ID,
			cast.ToInt(deviceConfig.DeviceType))

		if err != nil {
			log.Printf("Failed to add the device with ip: %s\n", deviceConfig.IP)
		} else {
			log.Printf("Device \"%s\" added\n", deviceConfig.Name)
		}
	}
}

func (deviceControl *DeviceControl) AddOrUpdateDiscoveredDevice(name string, mac string) error {
	deviceInfo, err := deviceControl.broadlink.GetDeviceInfo(mac)

	if err != nil {
		return err
	}

	if deviceControl.config.Devices == nil {
		deviceControl.config.Devices = make(map[string]*Device)
	}

	deviceCategory := DeviceBlaster

	if deviceInfo.SupportsPower {
		deviceCategory = DevicePowerSwitch
	}

	deviceControl.config.Devices[deviceInfo.Mac] = &Device{
		Name:           name,
		IP:             deviceInfo.Ip,
		Mac:            deviceInfo.Mac,
		Key:            deviceInfo.Key,
		ID:             deviceInfo.Id,
		DeviceType:     deviceInfo.DeviceType,
		DeviceCategory: deviceCategory,
		Enabled:        true,
	}

	return nil
}
