package impl

import (
	"fmt"
	"mutex-protocols/peterson/config"
	"mutex-protocols/peterson/models"
	"sync"
)

func runProcess(process *models.Process, commonResources *Resources) {
	// assume process is already in LocalSection state
	for step := 0; step < (process.GetNoOfSteps() / 4); step++ {
		process.Delay()

		process.ChangeState(config.EntryProtocol)
		entryProtocol(process.GetId(), commonResources)

		process.ChangeState(config.CriticalSection)
		process.Delay()

		process.ChangeState(config.ExitProtocol)
		exitProtocol(process.GetId(), commonResources)

		process.ChangeState(config.LocalSection)
	}
	process.Kill() // for animation purposes
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
