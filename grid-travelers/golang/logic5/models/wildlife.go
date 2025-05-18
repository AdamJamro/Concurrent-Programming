package models

import (
	"context"
	"fmt"
	"grid-travelers-v5/config"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Counter is a struct that holds our atomic counter.  Using a struct
type Counter struct {
	value int64 // The counter value.  Must be int64 for atomic ops.
}

func NewCounter() *Counter {
	return &Counter{value: 0}
}

func (c *Counter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

func (c *Counter) Decrement() {
	atomic.AddInt64(&c.value, -1)
}

// Get atomically reads the current value of the counter.
func (c *Counter) Get() int64 {
	return atomic.LoadInt64(&c.value)
}

func RunWildlifeSpawner(requestChannels *TileRequestChannels, forcedMoveChannels *ForcedMoveChannels, ctx context.Context, reportChannel chan *TraceSequenceType, wg *sync.WaitGroup) {
	gen := config.NewGenerator()
	wildLifeSemaphore := make(chan struct{}, config.MaxWildlifeSpawn)
	wildlifeCounter := NewCounter()
	symbolSequence := NewCounter()
	lastSpawnTrialSuccessful := false

	for i := 0; i < config.MaxWildlifeSpawn; i++ {
		wildLifeSemaphore <- struct{}{}
	}

	for {
		wildlifeCounter.Increment()
		index := gen.Intn(config.GridHeight * config.GridWidth)
		x := index % config.GridWidth
		y := index / config.GridWidth
		id := int(wildlifeCounter.Get())
		traceSequence := TraceSequenceType{
			len:  0,
			data: [config.MaxSteps + 1]TraceType{},
		}
		wildlife := &Traveler{
			id: id + config.NoOfTravelers,
			pos: Position{
				x: x,
				y: y,
			},
			symbol:        rune('0' + id - 1),
			instanceType:  Wildlife,
			generator:     config.NewGenerator(),
			traceSequence: traceSequence,
		}

		select {
		case <-ctx.Done():
			fmt.Println("Spawner finished")
			return
		case <-wildLifeSemaphore:
			if lastSpawnTrialSuccessful {
				time.Sleep(config.SpawnRate)
			}
			lastSpawnTrialSuccessful = false

			// spawn
			arbitraryTimeout := time.After(time.Second)
			tileResponse := make(chan ResponseType)
			select {
			case requestChannels.tiles[y][x].channel <- RequestType{
				requester: Requester{
					pos:  Position{x: x, y: y},
					id:   wildlife.id,
					kind: wildlife.instanceType,
				},
				destination: Position{x: x, y: y},
				timestamp:   time.Now(),
				timeout:     time.Second,
				response:    tileResponse,
			}:
				response := <-tileResponse
				if response.ResponseCode == ResponseCodeAmbushed {
					wildlife.addTrace()
					time.Sleep(100)
					wildlife.symbol = '*'
					wildlife.addTrace()
					time.Sleep(config.AmbushTimeout)
					wildlife.MoveOutsideTheBoard(&requestChannels.tiles[y][x])
				}
				if response.ResponseCode != ResponseCodeSuccess {
					fmt.Println("Spawn failed, wildlife id:", wildlife.id, "pos:", wildlife.pos)
					wildlifeCounter.Decrement()
					wildLifeSemaphore <- struct{}{}
					continue
				}
			case <-arbitraryTimeout:
				// timeout
				wildlifeCounter.Decrement()
				wildLifeSemaphore <- struct{}{}
				continue
			}

			fmt.Println("Spawned wildlife no.", wildlife.id, " pos:", wildlife.pos)
			wg.Add(1)

			wildlife.symbol = rune('0' + (symbolSequence.Get() % 10))
			symbolSequence.Increment()
			lastSpawnTrialSuccessful = true // REMEMBER WE SUCCEEDED

			// simulate
			go func() {
				defer wg.Done()
				wildlife.addTrace()
				runWildLife(wildlife, requestChannels, forcedMoveChannels, ctx)
				wildlifeCounter.Decrement()
				wildLifeSemaphore <- struct{}{}
				fmt.Println("Wildlife no,", wildlife.id, " finished, trying to print report...")
				reportChannel <- &wildlife.traceSequence
			}()
		}
	}
}

func runWildLife(wildlife *Traveler, requestChannels *TileRequestChannels, forcedMoveChannels *ForcedMoveChannels, ctx context.Context) {
	extinct := time.After(time.Duration((config.MaxWildLifetime-config.MinWildLifetime)*rand.Float64() + config.MinWildLifetime))

	for {
		moveTimeout := time.After(wildlife.GetDelayDuration() * 2)
		select {
		case <-extinct:
			wildlife.MoveOutsideTheBoard(&requestChannels.tiles[wildlife.pos.y][wildlife.pos.x])
			fmt.Println("wildlife,", wildlife.id, " died")
			return
		case <-ctx.Done():
			fmt.Println("Wildlife finished")
			wildlife.MoveOutsideTheBoard(&requestChannels.tiles[wildlife.pos.y][wildlife.pos.x])
			fmt.Println("wildlife,", wildlife.id, " died")
			return
		case forcedMove := <-forcedMoveChannels.tiles[wildlife.pos.y][wildlife.pos.x].channel:
			wildlife.pos = forcedMove.destination
			wildlife.addTraceAt(forcedMove.timestamp)
			fmt.Println("Wildlife no.", wildlife.id, "gave way and moved to", wildlife.pos, "at", forcedMove.timestamp.Nanosecond())
			forcedMove.response <- ResponseType{
				ResponseCode: ResponseCodeSuccess,
			}
		case <-moveTimeout:
			//fmt.Println("Wildlife running")
			response := wildlife.MakeRandomMove(time.Duration(config.MinTimeout/2), requestChannels)
			if response.ResponseCode == ResponseCodeAmbushed {
				return
			}
			// ignore timeouts and errors in case of wildlife
			// this is an arbitrary decision
		}
	}
}
