package devicecontrol

import (
	"errors"
	"fmt"
	"smh-apiengine/pkg/alexakit"
	"strings"
)

// HandleAlexaRequest tries to find the command and device for the alexa request execution. In case of execution
// failure, for example because the device has changed the ip address, retries to discover the devices again and
// execute command. If the execution was successful, updates the device's data save it into config json file
func (deviceControl *DeviceControl) HandleAlexaRequest(reqIntent alexakit.SimpleIntent) error {
	scenario, err := deviceControl.config.findScenario(reqIntent)

	if err == nil {
		if len(scenario.Sequence) > 0 {
			return deviceControl.ExecScenarioFullCycle(scenario)
		}

		return fmt.Errorf("scenario \"%s\" has no sequence items", scenario.Name)
	}

/*	cmd, err := deviceControl.config.findCommand(reqIntent)

	if err != nil {
		return err
	}

	err = deviceControl.ExecCommandFullCycle(cmd)

	if err != nil {
		return err
	}*/

	return nil
}

// NewSimpleRequestIntent create SimpleRequestIntent struct from full AlexaRequest
func (deviceControl *DeviceControl) NewSimpleRequestIntent(request alexakit.AlexaRequest) (alexakit.SimpleIntent, error) {
	var simpleRequestIntent alexakit.SimpleIntent
	intent := request.Request.Intent

	if len(intent.Slots) == 0 {
		return simpleRequestIntent, errors.New("no intents found in the request")
	}

	if _, ok := deviceControl.config.Intents[intent.Name]; !ok {
		return alexakit.SimpleIntent{}, errors.New("intent not supported")
	}

	targetSlots := deviceControl.config.Intents[intent.Name].Slots
	requestSlots := map[string]alexakit.SimpleSlot{}

	for _, slot := range intent.Slots {
		value, err := deviceControl.config.searchSlotValueWithSynonyms(targetSlots, slot)

		if err != nil {
			return simpleRequestIntent, errors.New("can not find the original supported slot value")
		}

		requestSlots[slot.Name] = alexakit.SimpleSlot{Name: slot.Name, Value: value}
	}

	return alexakit.SimpleIntent{
		Name:  intent.Name,
		Slots: requestSlots,
	}, nil
}

// searchSlotValueWithSynonyms searches for corresponding slot value, defined in intents json config. Compares not
// only name but also synonyms if any defined
func (c *Config) searchSlotValueWithSynonyms(targetSlots map[string]slot, searchSlot alexakit.Slot) (string, error) {
	if _, ok := targetSlots[searchSlot.Name]; !ok {
		return "", errors.New("slot not found")
	}

	values := targetSlots[searchSlot.Name].Values

	for _, value := range values {
		if strings.EqualFold(value.Name, searchSlot.Value) || c.iContains(value.Synonyms, searchSlot.Value) {
			return value.Name, nil
		}
	}

	return "", errors.New("no supported value found")
}

func (c *Config) findCommand(reqIntent alexakit.SimpleIntent) (Command, error) {
	for _, command := range c.Commands {
		if len(c.Intents) == 0 {
			continue
		}

		if c.hasMatchedIntent(command.Intents, reqIntent) {
			return command, nil
		}
	}

	return Command{}, fmt.Errorf("command not found. Searched for: %s", reqIntent)
}

func (c *Config) findScenario(reqIntent alexakit.SimpleIntent) (Scenario, error) {
	for _, scenario := range c.Scenarios {
		if len(scenario.Intents) == 0 {
			continue
		}

		if c.hasMatchedIntent(scenario.Intents, reqIntent) {
			return scenario, nil
		}
	}

	return Scenario{}, fmt.Errorf("scenario not found. Searched for: %s", reqIntent)
}

func (c *Config) hasMatchedIntent(intents []CommandIntent, reqIntent alexakit.SimpleIntent) bool {
	for i := 0; i < len(intents); i++ {
		if intents[i].Name != reqIntent.Name {
			continue
		}

		if c.hasMatchedSlots(intents[i].Slots, reqIntent.Slots) {
			return true
		}
	}

	return false
}

func (c *Config) hasMatchedSlots(sourceSlots map[string]CommandSlot, lookupSlots map[string]alexakit.SimpleSlot) bool {
	for _, lookupSlot := range lookupSlots {
		if slot, ok := sourceSlots[lookupSlot.Name]; ok {
			if !strings.EqualFold(slot.Value, lookupSlot.Value) {
				return false
			}
		} else {
			return false
		}
	}

	return true
}

func (c *Config) iContains(arr []string, val string) bool {
	for i := 0; i < len(arr); i++ {
		if strings.EqualFold(arr[i], val) {
			return true
		}
	}

	return false
}
