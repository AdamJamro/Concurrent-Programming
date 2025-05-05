package travelers_simulator

import (
	"context"
	"fmt"
	"grid-travelers-v4/config"
	"grid-travelers-v4/models"
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
		case models.ResponseCodeNoPath:
			fmt.Println("runTraveler: ResponseCodeNoPath moving traveler no. ", traveler.GetId(), "onto (", traveler.GetPosition(), ")")
			// try again
			continue
		case models.ResponseCodeError:
			// should never happen!!!
			fmt.Println("runTraveler traveler no. ", traveler.GetId(), " got error whilst moving onto (", traveler.GetPosition(), ")")
			traveler.KillAt(response.ChangeTimestamp)
			return
		}
	}
}

func RunSimulation() {
	// synchronization
	tileServersCtx, cancelTileServers := context.WithCancel(context.Background())
	var tileRequestChannels models.TileRequestChannels
	tileRequestChannels.InitTileRequestChannels()
	var forcedMoveChannels models.ForcedMoveChannels
	forcedMoveChannels.InitForcedMoveChannels()
	models.InitTileServers(&tileRequestChannels, &forcedMoveChannels, tileServersCtx)

	// initialize travelers
	travelers := models.MakeTravelers(&tileRequestChannels)
	wg := sync.WaitGroup{}
	//fmt.Println("Travelers has been initialized")

	fmt.Printf("-1 %d %d %d\n", config.NoOfTravelers+config.MaxWildlifeSpawn+1, config.GridWidth, config.GridHeight)

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

	for i := range config.NoOfTravelers {
		travelers[i].PrintReport()
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
	//printerFinish <- struct{}{}
	//fmt.Println("Simulation stops.")
}

func RunTests() {
	// playground config

	tileServersCtx, cancelTileServers := context.WithCancel(context.Background())
	var tileRequestChannels models.TileRequestChannels
	tileRequestChannels.InitTileRequestChannels()
	var forcedMoveChannels models.ForcedMoveChannels
	forcedMoveChannels.InitForcedMoveChannels()
	models.InitTileServers(&tileRequestChannels, &forcedMoveChannels, tileServersCtx)

	// scenario #1

	var travelers []*models.Traveler
	var wildlife []*models.Traveler
	travelers = append(travelers, models.MakeTraveler(0, 0, 0, models.GridTraveler, &tileRequestChannels))
	travelers = append(travelers, models.MakeTraveler(1, 0, 1, models.GridTraveler, &tileRequestChannels))
	wildlife = append(wildlife, models.MakeTraveler(2, 1, 1, models.Wildlife, &tileRequestChannels))
	timeout := time.Second

	fmt.Printf("-1 %d %d %d\n", 20, config.GridWidth, config.GridHeight)

	//travelers[0].Move(models.Down, timeout, &tileRequestChannels)
	time.Sleep(500)
	travelers[1].Move(models.Right, timeout, &tileRequestChannels)

	for _, traveler := range travelers {
		traveler.PrintReport()
	}
	for _, wildlife := range wildlife {
		wildlife.PrintReport()
	}
	cancelTileServers()

}
