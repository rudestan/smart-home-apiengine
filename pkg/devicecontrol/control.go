package devicecontrol

import (
	"github.com/rudestan/broadlinkrm"
	"github.com/spf13/cast"
	"log"
)

type DeviceControl struct {
	config *Config
	broadlink broadlinkrm.Broadlink
	lock *lock
}

func NewDeviceControl(configFile string) (DeviceControl, error)  {
	config, err := NewConfiguration(configFile)

	if err != nil {
		return DeviceControl{}, err
	}

	deviceControl := DeviceControl{
		config:    &config,
		broadlink: broadlinkrm.NewBroadlink(),
		lock: 	   &lock{locked:false},
	}

	if len(config.Devices) > 0 {
		deviceControl.initDevices()
	}

	return deviceControl, nil
}

// FindCommandByID finds Command structure by provided id or error if there is no Command found
func (deviceControl *DeviceControl) FindCommandByID(id string) (Command, error) {
	return deviceControl.config.findCommandByID(id)
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
