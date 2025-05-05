package travelers_simulator

import (
	"fmt"
	"grid-travelers-v4/config"
	"grid-travelers-v4/models"
	"sync"
)

func runTraveler(traveler *models.Traveler, semaphores models.GridFieldSemaphores) {
	for range traveler.GetNoOfSteps() {
		delay := traveler.Delay()
		if success := traveler.MakeRandomMove(config.MaxDelay-delay, semaphores); !success {
			//fmt.Println("Timeout reached for traveler no. ", traveler.GetId())
			traveler.RuneSymbolToLowerCase()
			_ = traveler.LogState() // discard errors
			break
		}
	}
}

func RunSimulation() {
	// synchronization
	gridSemaphores := models.GridFieldSemaphores{}
	gridSemaphores.InitGridFields()

	// initialize travelers
	travelers := models.MakeTravelers(gridSemaphores)
	wg := sync.WaitGroup{}
	//fmt.Println("Travelers has been initialized")

	//fmt.Println("timestamp | id | x | y | id-symbol")
	fmt.Println("-1 " + " 15 15")

	// run travelers
	for i := range config.NoOfTravelers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runTraveler(&travelers[id], gridSemaphores)
		}(i)
	}

	wg.Wait()
	for i := range config.NoOfTravelers {
		travelers[i].PrintReport()
	}

	//fmt.Println("Simulation stops.")
}
