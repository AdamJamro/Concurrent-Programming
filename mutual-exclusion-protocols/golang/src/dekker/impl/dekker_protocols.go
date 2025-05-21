package impl

import (
	"mutex-protocols/dekker/config"
	"sync"
	"sync/atomic"
	"time"
)

type Resources struct {
	Choosing *VolatileBoolArray
	Turn     int32 // either 0 or 1
}

type VolatileBoolArray struct {
	mutex *[config.NoOfProcesses]sync.Mutex
	data  *[config.NoOfProcesses]bool
}
type VolatileBoolArrayInterface interface {
	Get(index int) bool
	Set(index int, value bool)
}

func (v *VolatileBoolArray) Get(index int) bool {
	v.mutex[index].Lock()
	defer v.mutex[index].Unlock()
	return v.data[index]
}
func (v *VolatileBoolArray) Set(index int, value bool) {
	v.mutex[index].Lock()
	defer v.mutex[index].Unlock()
	v.data[index] = value
}

func createResources() *Resources {
	resources := &Resources{
		Choosing: &VolatileBoolArray{
			data:  &[config.NoOfProcesses]bool{},
			mutex: &[config.NoOfProcesses]sync.Mutex{},
		},
		Turn: int32(0), // 0 or 1, arbitrary
	}
	for i := 0; i < config.NoOfProcesses; i++ {
		resources.Choosing.Set(i, false)
	}
	return resources
}

func entryProtocol(processId int, res *Resources) {
	pid := int32(processId)

	res.Choosing.Set(processId, true)
	for res.Choosing.Get(1-processId) == true {
		if atomic.LoadInt32(&res.Turn) != pid {
			res.Choosing.Set(processId, false)
			for atomic.LoadInt32(&res.Turn) != pid {
				// wait
			}
			res.Choosing.Set(processId, true)
		}
	}
	time.Sleep(5 * time.Millisecond) // animation purposes only
}

func exitProtocol(processId int, res *Resources) {
	// prepare value as if we still were inside synchronized block
	opponentId := int32(1 - processId)

	// perform the exit protocol
	res.Choosing.Set(processId, false)
	atomic.StoreInt32(&res.Turn, opponentId)
	time.Sleep(5 * time.Millisecond) // animation purposes only
}
