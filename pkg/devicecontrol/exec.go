package devicecontrol

import (
	"errors"
	"fmt"
	"log"
	"github.com/rudestan/broadlinkrm"
	"time"
)

const (
	stateOff = "off"
	stateOn = "on"
)

type DeviceState struct {
	Id string `json:"id"`
	State string `json:"state"`
}

// ExecScenarioFullCycle executes scenario full cycle with commands one after another, including the delay
func (deviceControl *DeviceControl) ExecScenarioFullCycle(scenario Scenario) error {
	log.Printf("Executing scenario \"%s\" with %d sequence items", scenario.Name, len(scenario.Sequence))

	for _, sequenceItem := range scenario.Sequence {
		log.Printf("Executing sequence item \"%s\"", sequenceItem.CommandId)

		cmd := deviceControl.config.FindCommandByID(sequenceItem.CommandId)
		if cmd == nil {
			return errors.New("command not found")
		}

		err := deviceControl.ExecCommandFullCycle(*cmd)
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

// ExecControlItem executes the command in full cycle with retry and discover, as well as updating and saving
// the device data
func (deviceControl *DeviceControl) ExecControlItem(controlItem *ControlItem, state string) error {
	var stateEntity *Entity

	if state != "" {
		stateEntity = controlItem.FindEntityByState(state)

		if stateEntity == nil {
			return errors.New(fmt.Sprintf("Can not find Entity with %s state", state))
		}
	} else {
		stateEntity = controlItem.FindNextStateEntity()

		if stateEntity == nil {
			return errors.New("Can not find Entity with next state")
		}
	}

	switch stateEntity.Type {
	case ElementTypeCommand:
		cmd := deviceControl.config.FindCommandByID(stateEntity.Target)
		if cmd == nil {
			return errors.New("command not found")
		}

		err := deviceControl.ExecCommand(cmd)
		if err != nil {
			return err
		}
	case ElementTypeScenario:
		scenario := deviceControl.config.FindScenarioByID(stateEntity.Target)
		if scenario == nil {
			return errors.New("scenario not found")
		}

		err := deviceControl.ExecScenario(scenario)
		if err != nil {
			return err
		}
	default:
		return errors.New("unknown element type")
	}

	controlItem.activeState = stateEntity.State

	return nil
}

// ExecCommandWithRetryAndDiscover executes the command on passed device, in case of failure calls the execution
// with device discovering
func (deviceControl *DeviceControl) execCommandWithRetryAndDiscover(device *Device, command Command) error {
	err := deviceControl.ExecCommand(&command)

	if err == nil {
		return nil
	}

	log.Printf("Failed trying with discovering: %s\n", err)

	err = deviceControl.Discover(true)

	if err != nil {
		return err
	}

	log.Printf("Retrying execution on device: %s (%s, %s)\n", device.Name, device.IP, device.Mac)

	return deviceControl.ExecCommand(&command)
}

// discover function discovers the devices, this operation is time consuming and should be executed only once (for
// example from some goroutine).
func (deviceControl *DeviceControl) Discover(debug bool) error {
	if deviceControl.lock.Locked() {
		log.Println("device control is locked, can not discover")
		return nil
	}

	deviceControl.lock.Lock()
	defer deviceControl.lock.Unlock()

	deviceControl.broadlink = broadlinkrm.NewBroadlink()

	if !debug {
		deviceControl.broadlink.DebugOff()
	}

	err := deviceControl.broadlink.Discover()

	if err != nil {
		return err
	}

	return nil
}

// ExecScenarioFullCycle executes scenario full cycle with commands one after another, including the delay
func (deviceControl *DeviceControl) ExecScenario(scenario *Scenario) error {
	log.Printf("Executing scenario \"%s\" with %d sequence items", scenario.Name, len(scenario.Sequence))

	for _, sequenceItem := range scenario.Sequence {
		log.Printf("Executing sequence item \"%s\"", sequenceItem.CommandId)

		cmd := deviceControl.config.FindCommandByID(sequenceItem.CommandId)
		if cmd == nil {
			return errors.New("command not found")
		}

		err := deviceControl.ExecCommand(cmd)
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

// execCommand executes command on the device, should not be during lock. The operation can be time consuming in
// case the device is not available on the network. It will fail on timeout.
func (deviceControl *DeviceControl) ExecCommand(command *Command) error  {
	if deviceControl.lock.Locked() {
		log.Println("device control is locked, can not execute command")
		return nil
	}

	deviceControl.lock.Lock()
	defer deviceControl.lock.Unlock()

	device := deviceControl.config.FindDeviceById(command.DeviceID)

	if device == nil {
		return errors.New(fmt.Sprintf("No device with id %s found", command.DeviceID))
	}

	return deviceControl.broadlink.Execute(device.Mac, command.Code)
}

func (deviceControl *DeviceControl) getPowerState(device *Device) error  {
	powerState, err := deviceControl.broadlink.GetPowerState(device.Mac)

	if err != nil {
		return err
	}

	state := stateOff

	if powerState {
		state = stateOn
	}

	now := time.Now().String()
	deviceControl.DeviceStateBuffer[now] = DeviceState{
		Id:    device.ID,
		State: state,
	}

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

	return deviceControl.config.SaveConfiguration(deviceControl.config.fileName)
}
