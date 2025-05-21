package impl

import (
	"mutex-protocols/peterson/config"
	"sync"
	"sync/atomic"
	"time"
)

type Resources struct {
	Choosing *VolatileBoolArray
	Last     int32 // either 0 or 1
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
		Last: int32(0), // 0 or 1, arbitrary
	}
	for i := 0; i < config.NoOfProcesses; i++ {
		resources.Choosing.Set(i, false)
	}
	return resources
}

func entryProtocol(processId int, res *Resources) {
	pid := int32(processId)

	res.Choosing.Set(processId, true)
	atomic.StoreInt32(&res.Last, pid)
	for res.Choosing.Get(1-processId) == true && atomic.LoadInt32(&res.Last) == pid {
		// wait
	}
	time.Sleep(5 * time.Millisecond) // animation purposes only
}

func exitProtocol(processId int, res *Resources) {
	res.Choosing.Set(processId, false)
	time.Sleep(5 * time.Millisecond) // animation purposes only
}
