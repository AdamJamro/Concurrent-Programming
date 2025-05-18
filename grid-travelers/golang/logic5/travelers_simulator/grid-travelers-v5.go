package travelers_simulator

import (
	"context"
	"fmt"
	"grid-travelers-v5/config"
	"grid-travelers-v5/models"
	"sync"
	"time"
)

func runTraveler(traveler *models.Traveler, channels *models.TileRequestChannels) {
	for range traveler.GetNoOfSteps() {
		//timestamp := time.Now()
		delay := traveler.Delay()

		originalPos := traveler.GetPosition()

		response := traveler.MakeRandomMove(config.MinTimeout+delay, channels)

		switch response.ResponseCode {
		case models.ResponseCodeSuccess:
			fmt.Println("runTraveler moved traveler no. ", traveler.GetId(), "from", originalPos, "to", traveler.GetPosition())
			//traveler = response.instance
			//t.traceSequence = response.instance.traceSequence
			//t.pos = response.instance.pos
			continue
		case models.ResponseCodeTimeout:
			// assume deadlock
			fmt.Println("runTraveler Timeout reached for traveler no. ", traveler.GetId())
			traveler.KillAt(response.ChangeTimestamp)
			return
		case models.ResponseCodeNoPath: // deprecated (could be used to avoid timeout caused by wildlife but was not implemented)
			fmt.Println("runTraveler: ResponseCodeNoPath moving traveler no. ", traveler.GetId(), "onto (", traveler.GetPosition(), ")")
			// try again
			continue
		case models.ResponseCodeAmbushed:
			return // end the walk, traveler.Move() handles this case autonomously
		case models.ResponseCodeError:
			// should never happen!!!
			// MakeRandomMove() and Move() do not return this code
			fmt.Println("runTraveler traveler no. ", traveler.GetId(), " got error whilst moving onto (", traveler.GetPosition(), ")")
			traveler.KillAt(response.ChangeTimestamp)
			return
		}
	}
}

func RunSimulation() {
	// result script ids of travelers are going to be
	// 1) grid travelers (0, 1, 2, ..., noOfTravelers-1)
	// 2) wildlife (noOfTravelers, noOfTravelers+1, ..., noOfTravelers+maxWildlifeSpawn-1)
	// 3) ambushes (noOfTravelers+maxWildlifeSpawn, noOfTravelers+maxWildlifeSpawn+1, ..., noOfTravelers+maxWildlifeSpawn+noOfAmbushes-1)

	// synchronization
	tileServersCtx, cancelTileServers := context.WithCancel(context.Background())
	var tileRequestChannels models.TileRequestChannels
	tileRequestChannels.InitTileRequestChannels()
	var forcedMoveChannels models.ForcedMoveChannels
	forcedMoveChannels.InitForcedMoveChannels()
	var ambushSeizeRequestChannels models.AmbushSeizeRequestChannels = models.AmbushSeizeRequestChannels{
		Channels: make(map[models.Position]*models.AmbushSeizeRequestChannel),
	}

	// initialize tile servers
	models.InitTileServers(&tileRequestChannels, &forcedMoveChannels, &ambushSeizeRequestChannels, tileServersCtx)

	// initialize travelers i.e.
	// grid travelers & ambushes only,
	// since wildlife is going to be randomly spawned later
	travelers, ambushes := models.MakeTravelers(&tileRequestChannels, &ambushSeizeRequestChannels)
	wg := sync.WaitGroup{}
	//fmt.Println("Travelers has been initialized")

	// run travelers
	for i := range config.NoOfTravelers {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Println("Spawned Traveler no", i)
			runTraveler(&travelers[id], &tileRequestChannels)
			fmt.Println("Died Traveler no", i)
		}(i)
	}
	for i := range config.NoOfAmbushes {
		go func(id int) {
			ambushSeizeRequestChannel := ambushSeizeRequestChannels.Channels[ambushes[id].GetPosition()]
			models.RunAmbush(&ambushes[id], ambushSeizeRequestChannel)
		}(i)
	}

	spawnCtx, cancelSpawner := context.WithCancel(context.Background())
	//run wildlife
	wildlifeReportChannel := make(chan *models.TraceSequenceType)
	wildlifeWaitGroup := sync.WaitGroup{}
	go models.RunWildlifeSpawner(&tileRequestChannels, &forcedMoveChannels, spawnCtx, wildlifeReportChannel, &wildlifeWaitGroup)

	fmt.Println("DEBUH Waiting for travelers to end")
	wg.Wait()
	fmt.Println("DEBUH trying to cancel spawner")
	cancelSpawner()
	// wait for all tile servers to finish
	time.Sleep(time.Second)

	fmt.Printf("-1 %d %d %d\n", config.NoOfTravelers+config.MaxWildlifeSpawn+config.NoOfAmbushes+1, config.GridWidth, config.GridHeight)
	for _, traveler := range travelers {
		traveler.PrintReport()
	}
	for _, ambush := range ambushes {
		ambush.PrintReport()
	}
	// print wildlife traces
	printerFinish := make(chan struct{})
	go func() {
		isEmpty := false
		for !isEmpty {
			emptyTimeout := time.After(time.Second * 2)
			select {
			case traceSequence := <-wildlifeReportChannel:
				traceSequence.PrintTrace()
			case <-emptyTimeout:
				isEmpty = true
			}
		}
		printerFinish <- struct{}{}
		//for {
		//	select {
		//	case traceSequence := <-wildlifeReportChannel:
		//		traceSequence.PrintTrace()
		//	case <-printerFinish:
		//		break
		//	}
		//}
	}()
	wildlifeWaitGroup.Wait()
	cancelTileServers()
	<-printerFinish
	//fmt.Println("Simulation stops.")
}
