package impl

import (
	"mutex-protocols/backery/config"
	"sync"
	"time"
)

type Resources struct {
	MaxTicket int32
	Choosing  *VolatileBoolArray
	Ticket    *VolatileIntArray
}

type VolatileIntArray struct {
	mutex *[config.NoOfProcesses]sync.Mutex
	data  *[config.NoOfProcesses]int32
}
type VolatileIntArrayInterface interface {
	Get(index int) int32
	Set(index int, value int32)
}

func (v *VolatileIntArray) Get(index int) int32 {
	v.mutex[index].Lock()
	defer v.mutex[index].Unlock()
	return v.data[index]
}
func (v *VolatileIntArray) Set(index int, value int32) {
	v.mutex[index].Lock()
	defer v.mutex[index].Unlock()
	v.data[index] = value
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
		MaxTicket: 0,
		Choosing: &VolatileBoolArray{
			data:  &[config.NoOfProcesses]bool{},
			mutex: &[config.NoOfProcesses]sync.Mutex{},
		},
		Ticket: &VolatileIntArray{
			data:  &[config.NoOfProcesses]int32{},
			mutex: &[config.NoOfProcesses]sync.Mutex{},
		},
	}
	for i := 0; i < config.NoOfProcesses; i++ {
		resources.Ticket.Set(i, 0)
		resources.Choosing.Set(i, false)
	}
	return resources
}

func (v *VolatileIntArray) GetMax() int32 {
	// this is a weaker type of mutual exclusion
	// it may slower down the program
	// but won't cause deadlock or starvation
	// we could opt out of this and sync the whole operation
	// by acquiring the mutex for the entire max value search
	currentMax := int32(0)
	for i := 0; i < len(v.data); i++ {
		val := v.Get(i)
		if val > currentMax {
			currentMax = val
		}
	}
	return currentMax
}

func entryProtocol(processId int, res *Resources) {
	res.Choosing.Set(processId, true)
	newTicket := res.Ticket.GetMax() + 1
	res.Ticket.Set(processId, newTicket)
	res.Choosing.Set(processId, false)
	for i := 0; i < config.NoOfProcesses; i++ {
		if i != processId {
			for res.Choosing.Get(i) == true { /*wait*/
			}
			for res.Ticket.Get(i) != 0 && res.Ticket.Get(i) <= res.Ticket.Get(processId) &&
				(res.Ticket.Get(i) != res.Ticket.Get(processId) || i <= processId) { /*wait*/
			}
		}
	}
	time.Sleep(10 * time.Millisecond)
}

func exitProtocol(processId int, res *Resources) {
	// util we don't change any resource we can treat this as a footer of critical section
	// thus we can safely update MaxTicket
	res.MaxTicket = max(res.MaxTicket, res.Ticket.Get(processId))

	// now we perform the exit protocol
	res.Ticket.Set(processId, 0)
	time.Sleep(10 * time.Millisecond)
}
