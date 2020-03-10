package devicecontrol

import (
	"errors"
	"log"

	"github.com/rudestan/broadlinkrm"
	"github.com/spf13/cast"
)


var errorTryingExecute = errors.New("devicecontrol is trying to execute command")

type lock struct {
	locked bool
}

type DeviceControl struct {
	config Config
	broadlink broadlinkrm.Broadlink
	lock *lock
}

func NewDeviceControl(configFile string) (DeviceControl, error)  {
	config, err := NewConfiguration(configFile)

	if err != nil {
		return DeviceControl{}, err
	}

	deviceControl := DeviceControl{
		config:    config,
		broadlink: broadlinkrm.NewBroadlink(),
		lock: 	   &lock{locked:false},
	}

	if len(config.Devices) > 0 {
		deviceControl.initDevices()
	}

	return deviceControl, nil
}

func (deviceControl *DeviceControl) Lock()  {
	deviceControl.lock.locked = true
}

func (deviceControl *DeviceControl) Unlock()  {
	deviceControl.lock.locked = false
}

func (deviceControl *DeviceControl) Locked() bool  {
	return deviceControl.lock.locked
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

// ExecScenarioFullCycle executes scenario full cycle with commands one after another, including the delay
func (deviceControl *DeviceControl) ExecScenarioFullCycle(scenario Scenario) error {
	log.Printf("Executing scenario \"%s\" with %d sequence items", scenario.Name, len(scenario.Sequence))

	for _, sequenceItem := range scenario.Sequence {
		log.Printf("Executing sequence item \"%s\"", sequenceItem.Name)

/*		cmd, err := deviceControl.config.findCommandByID(sequenceItem.Name)

		if err != nil {
			return err
		}

		err = deviceControl.ExecCommandFullCycle(cmd)

		if err != nil {
			return err
		}

		if sequenceItem.Delay > 0 {
			log.Printf("Sleeping %d seconds\n", sequenceItem.Delay)
			time.Sleep(time.Second * time.Duration(sequenceItem.Delay))
		}*/
	}

	return nil
}

// ExecCommandFullCycle executes the command in full cycle with retry and discover, as well as updating and saving
// the device data
//func (deviceControl *DeviceControl) ExecCommandFullCycle(command Command, errorChan chan error) error {
func (deviceControl *DeviceControl) ExecCommandFullCycle(command Command, errorChan chan error) error {
	device, err := deviceControl.config.findDeviceByMac(command.DeviceID)

	if err != nil {
		errorChan <- err

		return err
	}

	err = deviceControl.ExecCommandWithRetryAndDiscover(device, command, errorChan)

	if err != nil {
		return err
	}

	err = deviceControl.updateAndSaveMatchedDiscoveredDevice(device)

	if err != nil {
		return err
	}

	return nil
}

// ExecCommandWithRetryAndDiscover executes the command on passed device, in case of failure calls the execution
// with device discovering
func (deviceControl *DeviceControl) ExecCommandWithRetryAndDiscover(device *Device, command Command, errorChan chan error) error {
	if deviceControl.Locked() {
		errorChan <- errorTryingExecute

		return errorTryingExecute
	}

	log.Printf("Executing a command on device: %s (%s, %s)\n", device.Name, device.IP, device.Mac)

	err := deviceControl.execCommand(device.Mac, command.Code, errorChan)

	if err == nil {
		return err
	}

	log.Printf("Failed to execute a command, will retry with discovering: %s\n", err)

	return deviceControl.ExecCommandWithDiscover(device, command, errorChan)
}

func (deviceControl *DeviceControl) execCommand(mac string, code string, errorChan chan error) error  {
	deviceControl.Lock()
	defer deviceControl.Unlock()

	errorChan <- nil

	return deviceControl.broadlink.Execute(mac, code)
}

// ExecCommandWithDiscover executes the command on passed device
func (deviceControl *DeviceControl) ExecCommandWithDiscover(device *Device, command Command, execChan chan error) error {
	deviceControl.Lock()
	defer deviceControl.Unlock()

	log.Printf("Discovering the devices")

	execChan <- nil

	deviceControl.broadlink = broadlinkrm.NewBroadlink()
	err := deviceControl.broadlink.Discover()

	if err != nil {
		return err
	}

	return deviceControl.broadlink.Execute(device.Mac, command.Code)
}

func  (deviceControl *DeviceControl) updateAndSaveMatchedDiscoveredDevice(device *Device) error {
	deviceInfo, err := deviceControl.broadlink.GetDeviceInfo(device.Mac)

	if err != nil {
		return err
	}

	device.ID = deviceInfo.Id
	device.IP = deviceInfo.Ip
	device.Mac = deviceInfo.Mac
	device.DeviceType = deviceInfo.DeviceType
	device.Key = deviceInfo.Key

	return deviceControl.config.saveUpdated()
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
