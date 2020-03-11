package devicecontrol

// lock some custom lock implementation for time consuming operations. Unfortunately standard mutex is not really
// meets the requirements since goroutines waiting until mutex is unlocked and then continues execution. This is not
// our case since the routines must be executed anyway without waiting but the long running process (e.g. execution
// of the command on unavailable device or device discovering) should be blocked from execution more than once
type lock struct {
	locked bool
}

// Lock sets lock to true
func (deviceControl *DeviceControl) Lock() {
	deviceControl.lock.locked = true
}

// Unlock sets lock to false
func (deviceControl *DeviceControl) Unlock() {
	deviceControl.lock.locked = false
}

// Locked returns the state of lock
func (deviceControl *DeviceControl) Locked() bool {
	return deviceControl.lock.locked
}
