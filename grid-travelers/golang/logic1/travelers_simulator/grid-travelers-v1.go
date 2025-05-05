package travelers_simulator

import (
	"fmt"
	"grid-travelers-v1/config"
	"grid-travelers-v1/models"
	"sync"
)

func runTraveler(traveler *models.Traveler) {
	for range traveler.GetNoOfSteps() {
		traveler.Delay()
		traveler.MakeRandomMove()
		if errorStatus := traveler.SaveState(); errorStatus != nil {
			fmt.Println(errorStatus, "\nOverflow on saving state for traveler no. ", traveler.GetId())
		}
	}
}

func RunSimulation() {
	//fmt.Println("Simulation commences...")

	// initialize travelers
	travelers := models.MakeTravelers()
	wg := sync.WaitGroup{}
	//fmt.Println("Travelers has been initialized")

	//fmt.Println("timestamp | id | x | y | id-symbol")
	fmt.Println("-1 15 15 15")

	// run travelers
	for i := range config.NoOfTravelers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			runTraveler(&travelers[id])
			travelers[id].PrintReport()
		}(i)
	}

	wg.Wait()
	//fmt.Println("Simulation stops.")
}
