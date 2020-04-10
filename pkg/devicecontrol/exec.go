package devicecontrol

import (
	"github.com/rudestan/broadlinkrm"
	"log"
	"time"
)

// ExecScenarioFullCycle executes scenario full cycle with commands one after another, including the delay
func (deviceControl *DeviceControl) ExecScenarioFullCycle(scenario Scenario) error {
	log.Printf("Executing scenario \"%s\" with %d sequence items", scenario.Name, len(scenario.Sequence))

	for _, sequenceItem := range scenario.Sequence {
		log.Printf("Executing sequence item \"%s\"", sequenceItem.Name)

		cmd, err := deviceControl.config.findCommandByID(sequenceItem.Name)
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
		}
	}

	return nil
}

// ExecCommandFullCycle executes the command in full cycle with retry and discover, as well as updating and saving
// the device data
func (deviceControl *DeviceControl) ExecCommandFullCycle(command Command) error {
	device, err := deviceControl.config.findDeviceByMac(command.DeviceID)
	if err != nil {
		return err
	}

	err = deviceControl.execCommandWithRetryAndDiscover(device, command)
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
func (deviceControl *DeviceControl) execCommandWithRetryAndDiscover(device *Device, command Command) error {
	err := deviceControl.execCommand(device, command.Code)

	if err == nil {
		return nil
	}

	log.Printf("Failed trying with discovering: %s\n", err)

	err = deviceControl.discover()

	if err != nil {
		return err
	}

	log.Printf("Retrying execution on device: %s (%s, %s)\n", device.Name, device.IP, device.Mac)

	return deviceControl.execCommand(device, command.Code)
}

// discover function discovers the devices, this operation is time consuming and should be executed only once (for
// example from some goroutine).
func (deviceControl *DeviceControl) discover() error {
	if deviceControl.lock.Locked() {
		log.Println("device control is locked, can not discover")
		return nil
	}

	deviceControl.lock.Lock()
	defer deviceControl.lock.Unlock()

	deviceControl.broadlink = broadlinkrm.NewBroadlink()
	err := deviceControl.broadlink.Discover()

	if err != nil {
		return err
	}

	return nil
}

// execCommand executes command on the device, should not be during lock. The operation can be time consuming in
// case the device is not available on the network. It will fail on timeout.
func (deviceControl *DeviceControl) execCommand(device *Device, code string) error  {
	if deviceControl.lock.Locked() {
		log.Println("device control is locked, can not execute command")
		return nil
	}

	deviceControl.lock.Lock()
	defer deviceControl.lock.Unlock()

	err := deviceControl.broadlink.Execute(device.Mac, code)

	if err != nil {
		return err
	}

	err = deviceControl.getPowerState(device)

	if err != nil {
		log.Println(err)
	}

	return nil
}

func (deviceControl *DeviceControl) getPowerState(device *Device) error  {
	powerState, err := deviceControl.broadlink.GetPowerState(device.Mac)

	if err != nil {
		return err
	}

	var status string

	if powerState {
		status = "on"
	} else {
		status = "off"
	}

	deviceControl.DeviceLogs = append(deviceControl.DeviceLogs, "State changed: " + device.Name + " power state: " + status)

	return nil
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

	return deviceControl.config.saveConfiguration(deviceControl.config.fileName)
}
