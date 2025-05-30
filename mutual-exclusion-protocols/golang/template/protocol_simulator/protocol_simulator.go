package protocol_simulator

import (
	"fmt"
	"mutex-protocols-template/config"
	"mutex-protocols-template/models"
	"sync"
	"time"
)

// Resources
type Resources struct {
	// TODO
}

func entryProtocol(commonResources *Resources) {
	time.Sleep(time.Millisecond) // FOR ANIMATION PURPOSES
	// TODO
}

func exitProtocol(commonResources *Resources) {
	time.Sleep(time.Millisecond) // FOR ANIMATION PURPOSES
	// TODO
}

func createResources() *Resources {
	// TODO
	return &Resources{}
}

func runProcess(process *models.Process, commonResources *Resources) {
	// assume process is in LocalSection state
	for step := 0; step < (process.GetNoOfSteps() / 4); step++ {
		process.Delay()

		process.ChangeState(config.EntryProtocol)
		entryProtocol(commonResources)

		process.ChangeState(config.CriticalSection)
		process.Delay()

		process.ChangeState(config.ExitProtocol)
		exitProtocol(commonResources)

		process.ChangeState(config.LocalSection)
	}
}

func RunSimulation() {
	// initialize processes at the LocalSection each
	// and store them for animation purposes
	processes := models.MakeProcesses()
	processCommenceChannel := make(chan struct{})
	processWaitGroup := sync.WaitGroup{}

	commonResources := createResources()

	for i, _ := range processes {
		processWaitGroup.Add(1)
		go func(proc *models.Process, res *Resources) {
			defer processWaitGroup.Done()
			processCommenceChannel <- struct{}{}
			runProcess(proc, res)
		}(&processes[i], commonResources)
	}

	for range processes {
		// maybe, just maybe, this makes it more fair to compete at the start
		<-processCommenceChannel
	}

	processWaitGroup.Wait()
	fmt.Printf("-1  %d  %d  %d ", config.NoOfProcesses, config.GridWidth, config.GridHeight)
	for i := 0; i < config.NoOfSections; i++ {
		fmt.Printf("%s;", config.SectionLabels[i])
	}
	fmt.Printf("NO_EXTRA_LABELS;\n")

	for _, process := range processes {
		process.PrintReport()
	}
}
