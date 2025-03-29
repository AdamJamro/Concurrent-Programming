package models

import (
	"fmt"
	"grid-travelers-v1/config"
	"time"
)

type Traveler struct {
	id            int
	pos           position
	symbol        rune
	noOfSteps     int
	timestamp     time.Duration     // store creation time
	traceSequence TraceSequenceType // store all history
	generator     *config.Generator
}

type TravelerData struct {
	Id        int
	Symbol    rune
	Generator *config.Generator
}

// MakeTravelers defines how travelers are being made
func MakeTravelers() [config.NoOfTravelers]Traveler {
	generators := config.MakeGenerators(config.NoOfTravelers)
	travelers := [config.NoOfTravelers]Traveler{}

	for i := range config.NoOfTravelers {
		travelers[i].InitializeTraveler(
			TravelerData{Id: i, Symbol: rune('a' + i), Generator: generators[i]},
		)
	}

	return travelers
}

func (t *Traveler) InitializeTraveler(data TravelerData) {
	t.id = data.Id
	t.generator = data.Generator
	t.symbol = data.Symbol

	randomPosition := position{x: t.generator.Intn(config.GridWidth), y: t.generator.Intn(config.GridHeight)}
	t.pos = randomPosition
	t.noOfSteps = t.generator.Intn(config.MaxSteps-config.MinSteps) + config.MinSteps
	t.timestamp = nowInNanoseconds()
	errorStatus := t.addTrace()
	if errorStatus != nil {
		fmt.Println("")
	}

}

// addTrace saves current state and pushes it to the trace sequence
func (t *Traveler) addTrace() error {
	trace := TraceType{
		timeStamp: nowInNanoseconds() - t.timestamp,
		id:        t.id,
		pos:       t.pos,
		symbol:    t.symbol,
	}
	return t.traceSequence.add(trace)
}

func (t *Traveler) getPosition() position {
	return position{x: t.pos.x, y: t.pos.y}
}

func (t *Traveler) MakeRandomMove() {
	moveType := MoveType(t.generator.Intn(4))
	t.Move(moveType)
}

func (t *Traveler) SaveState() error {
	return t.addTrace()
}

// Move operates on board of a 2D torus topology
// with (config.GridHeight x config.GridWidth) dimensions
func (t *Traveler) Move(m MoveType) {
	switch m {
	case up:
		t.pos.y = (t.pos.y + 1) % config.GridHeight
	case down:
		t.pos.y = (t.pos.y - 1 + config.GridHeight) % config.GridHeight
	case left:
		t.pos.x = (t.pos.x - 1 + config.GridWidth) % config.GridWidth
	case right:
		t.pos.x = (t.pos.x + 1) % config.GridWidth
	default: // exemplary error handling
		fmt.Print("Error: Invalid move on Traveler ", t.id)
	}
}

func (t *Traveler) Delay() {
	time.Sleep(time.Duration(t.generator.Intn(config.MaxDelay-config.MinDelay)+config.MinDelay) * time.Nanosecond)
}

func (t *Traveler) GetNoOfSteps() int {
	return t.noOfSteps
}

func (t *Traveler) GetId() int {
	return t.id
}

func (t *Traveler) PrintReport() {
	t.traceSequence.PrintTrace()
}

func nowInNanoseconds() time.Duration {
	return time.Duration(time.Now().UnixNano()) * time.Nanosecond
}
