package devicecontrol

import (
    "github.com/rudestan/broadlinkrm"
    "github.com/spf13/cast"
    "log"
    "time"
)

var config Config
var broadlink broadlinkrm.Broadlink

// Init func initializes the device control json configuration, creates new broadlink struct and manually adds
// the devices to it from the config
func Init(configFile string) error {
    config, err := LoadConfiguration(configFile)
    if err != nil {
        return err
    }

    broadlink = broadlinkrm.NewBroadlink()

    if len(config.Devices) > 0 {
        initDevices(config.Devices, broadlink)
    }

    return nil
}

// FindCommandById finds Command structure by provided id or error if there is no Command found
func FindCommandById(id string) (Command, error)  {
    return config.findCommandById(id)
}

// FindCommandById finds Scenario structure by provided name or error if there is no Scenario found
func FindScenarioByName(name string) (Scenario, error)  {
    return config.findScenarioByName(name)
}

// ExecScenarioFullCycle executes scenario full cycle with commands one after another, including the delay
func ExecScenarioFullCycle(scenario Scenario) error {
    log.Printf("Executing scenario \"%s\" with %d sequence items", scenario.Name, len(scenario.Sequence))

    for _, sequenceItem := range scenario.Sequence {
        log.Printf("Executing sequence item \"%s\"", sequenceItem.Name)

        cmd, err := config.findCommandById(sequenceItem.Name)

        if err != nil {
            return err
        }

        err = ExecCommandFullCycle(cmd)

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
func ExecCommandFullCycle(command Command) error {
    device, err := config.findDeviceByMac(command.DeviceID)

    if err != nil {
        return err
    }

    err = ExecCommandWithRetryAndDiscover(device, command)

    if err != nil {
        return err
    }

    err = updateAndSaveMatchedDiscoveredDevice(device)

    if err != nil {
        return err
    }

    return nil
}

// ExecCommandWithRetryAndDiscover executes the command on passed device, in case of failure calls the execution
// with device discovering
func ExecCommandWithRetryAndDiscover(device *Device, command Command) error {
    log.Printf("Executing a command on device: %s (%s, %s)\n", device.Name, device.IP, device.Mac)
    err := broadlink.Execute(device.Mac, command.Code)

    if err == nil {
        return err
    }

    log.Printf("Failed to execute a command, will retry with discovering: %s\n", err)

    return ExecCommandWithDiscover(device, command)
}

// ExecCommandWithDiscover executes the command on passed device
func ExecCommandWithDiscover(device *Device, command Command) error {
    broadlink = broadlinkrm.NewBroadlink()
    err := broadlink.Discover()

    if err != nil {
        return err
    }

    return broadlink.Execute(device.Mac, command.Code)
}

func updateAndSaveMatchedDiscoveredDevice(device *Device) error {
    deviceInfo, err := broadlink.GetDeviceInfo(device.Mac)

    if err != nil {
        return err
    }

    device.ID = deviceInfo.Id
    device.IP = deviceInfo.Ip
    device.Mac = deviceInfo.Mac
    device.DeviceType = deviceInfo.DeviceType
    device.Key = deviceInfo.Key

    return SaveConfiguration(config.fileName)
}

func initDevices(devices map[string]*Device, broadlink broadlinkrm.Broadlink) {
    for _, deviceConfig := range devices {
        if !deviceConfig.Enabled {
            log.Printf("The device with ip: %s is disbled. Skipping\n", deviceConfig.IP)
            continue
        }

        err := broadlink.AddManualDevice(
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
