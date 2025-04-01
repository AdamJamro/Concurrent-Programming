package travelers_simulator

import (
	"fmt"
	"grid-travelers-v2/config"
	"grid-travelers-v2/models"
	"sync"
)

func runTraveler(traveler *models.Traveler, semaphores models.GridFieldSemaphores) {
	for range traveler.GetNoOfSteps() {
		delay := traveler.Delay()
		if success := traveler.MakeRandomMove(config.MaxDelay-delay, semaphores); !success {
			//fmt.Println("Timeout reached for traveler no. ", traveler.GetId())
			traveler.RuneSymbolToLowerCase()
			traveler.SaveState()
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
