package models

import (
	"fmt"
	"grid-travelers-v3/config"
	"time"
	"unicode"
)

type Traveler struct {
	id            int
	pos           position
	symbol        rune
	direction     MoveType
	noOfSteps     int
	timestamp     time.Duration     // store creation time
	traceSequence TraceSequenceType // store all history
	generator     *config.Generator
}

type TravelerData struct {
	Id        int
	Symbol    rune
	Generator *config.Generator
	Position  position
}

func InitialPosition(index int) position {
	return position{
		x: index,
		y: index,
	}
}

// MakeTravelers defines how travelers are being made
func MakeTravelers(semaphores GridFieldSemaphores) [config.NoOfTravelers]Traveler {
	generators := config.MakeGenerators(config.NoOfTravelers)
	travelers := [config.NoOfTravelers]Traveler{}

	for i := range config.NoOfTravelers {
		if errorMsg := travelers[i].InitializeTraveler(
			TravelerData{Id: i, Symbol: rune('A' + i), Generator: generators[i], Position: InitialPosition(i)},
			semaphores,
		); errorMsg != nil {
			fmt.Println(errorMsg)
		}
	}

	return travelers
}

func (t *Traveler) InitializeTraveler(data TravelerData, semaphores GridFieldSemaphores) error {
	t.id = data.Id
	t.generator = data.Generator
	t.symbol = data.Symbol

	// acquire the semaphore for the initial position
	timeout := time.After(time.Duration(5) * time.Second)
	select {
	case <-semaphores.at(data.Position.x, data.Position.y):
		break
	case <-timeout:
		return fmt.Errorf("error: Timeout while acquiring semaphore for initial position")
	}
	t.pos = data.Position
	// move up or down
	travelDirection := t.generator.Intn(2)
	if t.id%2 == 0 {
		// or move left or right
		travelDirection += 2
	}
	t.direction = MoveType(travelDirection)
	t.noOfSteps = t.generator.Intn(config.MaxSteps-config.MinSteps) + config.MinSteps
	t.timestamp = nowInNanoseconds()
	errorStatus := t.addTrace()
	if errorStatus != nil {
		return fmt.Errorf("initial addTrace() call failed ")
	}
	return nil
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

// MakeRandomMove unused (used in precious logic systems)
func (t *Traveler) MakeRandomMove(stepTimeout time.Duration, semaphores GridFieldSemaphores) bool {
	moveType := MoveType(t.generator.Intn(4))
	return t.Move(moveType, stepTimeout, semaphores)
}

func (t *Traveler) MakeMove(stepTimeout time.Duration, semaphores GridFieldSemaphores) bool {
	return t.Move(t.direction, stepTimeout, semaphores)
}

func (t *Traveler) SaveState() error {
	return t.addTrace()
}

// Move operates on board of a 2D torus topology
// with (config.GridHeight x config.GridWidth) dimensions
// returns true if move was successful, false otherwise
func (t *Traveler) Move(m MoveType, stepTimeout time.Duration, semaphores GridFieldSemaphores) bool {
	newPos := position{x: t.pos.x, y: t.pos.y}
	switch m {
	case up:
		newPos.y = (newPos.y + 1) % config.GridHeight
	case down:
		newPos.y = (newPos.y - 1 + config.GridHeight) % config.GridHeight
	case left:
		newPos.x = (newPos.x - 1 + config.GridWidth) % config.GridWidth
	case right:
		newPos.x = (newPos.x + 1) % config.GridWidth
	default: // exemplary error handling
		fmt.Print("Error: Invalid move on Traveler ", t.id)
	}
	timeout := time.After(stepTimeout)

	select {
	case <-semaphores.at(newPos.x, newPos.y):
		semaphores.at(t.pos.x, t.pos.y) <- struct{}{}
		t.pos = newPos
		return true
	case <-timeout:
		return false
	}
}

func (t *Traveler) Delay() time.Duration {
	delayTime := time.Duration(t.generator.Intn(config.MaxDelay-config.MinDelay)+config.MinDelay) * time.Nanosecond
	time.Sleep(delayTime)
	return delayTime
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

func (t *Traveler) RuneSymbolToLowerCase() {
	t.symbol = unicode.ToLower(rune(t.symbol))
}

func nowInNanoseconds() time.Duration {
	return time.Duration(time.Now().UnixNano()) * time.Nanosecond
}
