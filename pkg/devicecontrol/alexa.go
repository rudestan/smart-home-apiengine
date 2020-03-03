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
func HandleAlexaRequest(reqIntent alexakit.SimpleIntent) error {
    scenario, err := findScenario(reqIntent)

    if err == nil {
        if len(scenario.Sequence) > 0 {
            return ExecScenarioFullCycle(scenario)
        } else {
            return fmt.Errorf("scenario \"%s\" has no sequence items", scenario.Name)
        }
    }

    cmd, err := findCommand(reqIntent)

    if err != nil {
        return err
    }

    err = ExecCommandFullCycle(cmd)

    if err != nil {
        return err
    }

    return nil
}

// NewSimpleRequestIntent create SimpleRequestIntent struct from full AlexaRequest
func NewSimpleRequestIntent(request alexakit.AlexaRequest) (alexakit.SimpleIntent, error) {
    var simpleRequestIntent alexakit.SimpleIntent
    intent := request.Request.Intent

    if len(intent.Slots) == 0 {
        return simpleRequestIntent, errors.New("no intents found in the request")
    }

    if _, ok := config.Intents[intent.Name]; !ok {
        return alexakit.SimpleIntent{}, errors.New("intent not supported")
    }

    targetSlots := config.Intents[intent.Name].Slots
    requestSlots := map[string]alexakit.SimpleSlot{}

    for _, slot := range intent.Slots {
        value, err := searchSlotValueWithSynonyms(targetSlots, slot)

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
func searchSlotValueWithSynonyms(targetSlots map[string]Slot, searchSlot alexakit.Slot) (string, error) {
    if _, ok := targetSlots[searchSlot.Name]; !ok {
        return "", errors.New("slot not found")
    }

    values := targetSlots[searchSlot.Name].Values

    for _, value := range values {
        if strings.EqualFold(value.Name, searchSlot.Value) || iContains(value.Synonyms, searchSlot.Value) {
            return value.Name, nil
        }
    }

    return "", errors.New("no supported value found")
}

func findCommand(reqIntent alexakit.SimpleIntent) (Command, error) {
    for _, command := range config.Commands {
        if len(command.Intents) == 0 {
            continue
        }

        if hasMatchedIntent(command.Intents, reqIntent) {
            return command, nil
        }
    }

    return Command{}, fmt.Errorf("command not found. Searched for: %s", reqIntent)
}

func findScenario(reqIntent alexakit.SimpleIntent) (Scenario, error) {
    for _, scenario := range config.Scenarios {
        if len(scenario.Intents) == 0 {
            continue
        }

        if hasMatchedIntent(scenario.Intents, reqIntent) {
            return scenario, nil
        }
    }

    return Scenario{}, fmt.Errorf("scenario not found. Searched for: %s", reqIntent)
}

func hasMatchedIntent(intents []CommandIntent, reqIntent alexakit.SimpleIntent) bool {
    for i := 0; i < len(intents); i++ {
        if intents[i].Name != reqIntent.Name {
            continue
        }

        if hasMatchedSlots(intents[i].Slots, reqIntent.Slots) {
            return true
        }
    }

    return false
}

func hasMatchedSlots(sourceSlots map[string]CommandSlot, lookupSlots map[string]alexakit.SimpleSlot) bool {
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

func iContains(arr []string, val string) bool {
    for i := 0; i < len(arr); i++ {
        if strings.EqualFold(arr[i], val) {
            return true
        }
    }

    return false
}
