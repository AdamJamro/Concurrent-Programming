package travelers_simulator

import (
	"fmt"
	"grid-travelers-v3/config"
	"grid-travelers-v3/models"
	"sync"
)

func runTraveler(traveler *models.Traveler, semaphores models.GridFieldSemaphores) {
	for range traveler.GetNoOfSteps() {
		delay := traveler.Delay()
		if success := traveler.MakeMove(config.MaxDelay-delay, semaphores); !success {
			//fmt.Println("Timeout reached for traveler no. ", traveler.GetId())
			traveler.RuneSymbolToLowerCase()
			_ = traveler.SaveState() // do not handle errors here
			break
		}
		if errorStatus := traveler.SaveState(); errorStatus != nil {
			fmt.Println(errorStatus, "\nOverflow on saving state for traveler no. ", traveler.GetId())
		}
	}
}

func RunSimulation() {
	fmt.Println("Simulation commences...")

	// synchronization
	gridSemaphores := models.GridFieldSemaphores{}
	gridSemaphores.InitGridFields()

	// initialize travelers
	travelers := models.MakeTravelers(gridSemaphores)
	wg := sync.WaitGroup{}
	fmt.Println("Travelers has been initialized")

	fmt.Println("timestamp | id | x | y | id-symbol")

	// run travelers
	for i := range config.NoOfTravelers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runTraveler(&travelers[id], gridSemaphores)
			travelers[id].PrintReport()
		}(i)
	}

	wg.Wait()
	fmt.Println("Simulation stops.")
}
