package devicecontrol

// lock some custom lock implementation for time consuming operations. Unfortunately standard mutex is not really
// meets the requirements since goroutines waiting until mutex is unlocked and then continues execution. This is not
// our case since the routines must be executed anyway without waiting but the long running process (e.g. execution
// of the command on unavailable device or device discovering) should be blocked from execution more than once
import (
	"runtime"
	"sync/atomic"
)

const (
	unlocked int32 = iota
	locked
)

type spinLock struct {
	state int32
}
func (lock *spinLock) Lock() {
	for !atomic.CompareAndSwapInt32(&lock.state, unlocked, locked) {
		runtime.Gosched()
	}
}
func (lock *spinLock) Unlock() {
	for !atomic.CompareAndSwapInt32(&lock.state, locked, unlocked) {
		runtime.Gosched()
	}
}

func (lock *spinLock) Locked() bool  {
	return lock.state == locked
}
