package models

import (
	"fmt"
	"grid-travelers-v4/config"
	"time"
	"unicode"
)

type Traveler struct {
	id            int
	pos           Position
	symbol        rune
	noOfSteps     int
	instanceType  int               // either wildlife, grid traveler (or later an ambush)
	timestamp     time.Time         // store creation time
	traceSequence TraceSequenceType // store all history
	generator     *config.Generator
}

type TravelerData struct {
	Id           int
	Symbol       rune
	Generator    *config.Generator
	Position     Position
	Alive        bool
	InstanceType int
}

func MakeTraveler(id int, x int, y int, instanceType int, channels *TileRequestChannels) *Traveler {
	traveler := Traveler{}

	initialPosition := Position{x, y}

	if errorMsg := traveler.InitializeTraveler(
		TravelerData{
			Id:           id,
			Symbol:       rune('A' + id),
			Generator:    config.NewGenerator(),
			Position:     initialPosition,
			Alive:        true,
			InstanceType: instanceType,
		},
		channels,
	); errorMsg != nil {
		fmt.Println(errorMsg)
	}

	return &traveler
}

// MakeTravelers defines how travelers are being made
func MakeTravelers(tileRequestChannels *TileRequestChannels) [config.NoOfTravelers]Traveler {
	generators := config.MakeGenerators()
	travelers := [config.NoOfTravelers]Traveler{}

	initialPositions := generators[0].GenerateRandomInitialPositions()

	for i := range config.NoOfTravelers {
		if errorMsg := travelers[i].InitializeTraveler(
			TravelerData{
				Id:        i,
				Symbol:    rune('A' + i),
				Generator: generators[i],
				Position: Position{
					initialPositions[i].X,
					initialPositions[i].Y,
				},
				Alive:        true,
				InstanceType: GridTraveler,
			},
			tileRequestChannels,
		); errorMsg != nil {
			fmt.Println(errorMsg)
		}
	}
	return travelers
}

func (t *Traveler) InitializeTraveler(data TravelerData, requestChannels *TileRequestChannels) error {
	t.id = data.Id
	t.generator = data.Generator
	t.symbol = data.Symbol
	t.instanceType = data.InstanceType
	t.timestamp = time.Now()
	initialPosition := data.Position

	someArbitraryTimeout := time.Duration(2) * time.Second

	timeout := time.After(someArbitraryTimeout)
	responseChannel := make(chan ResponseType)
	select {
	case requestChannels.tiles[initialPosition.y][initialPosition.x].channel <- RequestType{
		requester: Requester{
			pos:  Position{x: initialPosition.x, y: initialPosition.y},
			id:   t.id,
			kind: t.instanceType,
		},
		destination: initialPosition,
		timestamp:   time.Now(),
		timeout:     someArbitraryTimeout,
		response:    responseChannel,
	}:
		response := <-responseChannel
		if response.ResponseCode != ResponseCodeSuccess {
			return fmt.Errorf("error: traveler %d Unable to acquire initial position", t.id)
		} else {
			fmt.Println(t.id, "Successfully acquired initial position,", time.Now().Nanosecond())
			t.pos = initialPosition
			t.addTrace()
		}

	case <-timeout:
		return fmt.Errorf("error: Timeout while acquiring semaphore for initial position")
	}

	t.noOfSteps = t.generator.Intn(config.MaxSteps-config.MinSteps) + config.MinSteps

	// this now happens when request is handled:
	//errorStatus := t.addTrace()
	//if errorStatus != nil {
	//	return fmt.Errorf("initial addTrace() call failed ")
	//}
	return nil
}

// addTrace saves current state and pushes it to the trace sequence
func (t *Traveler) addTrace() error {
	trace := TraceType{
		timeStamp: time.Now(),
		id:        t.id,
		pos:       t.pos,
		symbol:    t.symbol,
	}
	//if unicode.IsLower(t.symbol) {
	//	fmt.Println("traveler no.", t.id, "is lower case")
	//}
	return t.traceSequence.add(trace)
}

// addTrace saves current state and pushes it to the trace sequence
func (t *Traveler) addTraceAt(timestamp time.Time) error {
	trace := TraceType{
		timeStamp: timestamp,
		id:        t.id,
		pos:       t.pos,
		symbol:    t.symbol,
	}
	return t.traceSequence.add(trace)
}

func (t *Traveler) GetPosition() Position {
	return Position{x: t.pos.x, y: t.pos.y}
}

// MakeRandomMove unused (used in precious logic systems)
func (t *Traveler) MakeRandomMove(stepTimeout time.Duration, channels *TileRequestChannels) ResponseType {
	moveType := MoveType(t.generator.Intn(4))
	return t.Move(moveType, stepTimeout, channels)
}

func (t *Traveler) LogState() error {
	return t.addTrace()
}

// Move operates on board of a 2D torus topology
// with (config.GridHeight x config.GridWidth) dimensions
// returns passed on responseCode of the request handler defined in synchronization.go
func (t *Traveler) Move(m MoveType, stepTimeout time.Duration, requestChannels *TileRequestChannels) ResponseType {
	destinationPos := Position{x: t.pos.x, y: t.pos.y}
	switch m {
	case Down:
		destinationPos.y = (destinationPos.y + 1) % config.GridHeight
	case Up:
		destinationPos.y = (destinationPos.y - 1 + config.GridHeight) % config.GridHeight
	case Left:
		destinationPos.x = (destinationPos.x - 1 + config.GridWidth) % config.GridWidth
	case Right:
		destinationPos.x = (destinationPos.x + 1) % config.GridWidth
	default: // exemplary error handling
		fmt.Print("Error: Invalid move on Traveler ", t.id)
	}

	timeout := time.After(stepTimeout)

	responseChannel := make(chan ResponseType)
	for {
		timeReference := time.Now()
		select {
		case requestChannels.tiles[t.pos.y][t.pos.x].channel <- RequestType{
			requester: Requester{
				pos:  Position{x: t.pos.x, y: t.pos.y},
				id:   t.id,
				kind: t.instanceType,
			},
			destination: destinationPos,
			timestamp:   timeReference,
			timeout:     stepTimeout,
			response:    responseChannel,
		}:
			response := <-responseChannel
			if response.ResponseCode == ResponseCodeSuccess {
				fmt.Println("Traveler no.", t.id, "moved from", t.pos, "to", destinationPos, "at", response.ChangeTimestamp.Nanosecond())
				t.pos.y = destinationPos.y
				t.pos.x = destinationPos.x
				t.addTraceAt(response.ChangeTimestamp)
			} else if response.ResponseCode == ResponseCodeError {
				// try again
				fmt.Println("Traveler no.", t.id, "got error whilst moving onto (", destinationPos, ")", "and is retrying...")
				time.Sleep(50 * time.Microsecond)
				continue
			}
			return response
		case <-timeout:
			fmt.Println("in t.Move() Timeout reached for traveler no. ", t.id)
			return ResponseType{
				ResponseCode:    ResponseCodeTimeout,
				ChangeTimestamp: timeReference,
			}
		}
	}
}

func (t *Traveler) MoveOutsideTheBoard(tileChannel *RequestChannelType) {
	t.pos.x = config.GridWidth
	t.pos.y = config.GridHeight

	tileChannel.channel <- RequestType{
		requester: Requester{
			pos:  Position{x: t.pos.x, y: t.pos.y},
			id:   t.id,
			kind: t.instanceType,
		},
		destination: Position{x: config.GridWidth, y: config.GridHeight},
		timestamp:   time.Now(),
		timeout:     10 * time.Second,
		response:    nil,
	}

	_ = t.addTrace() // discard overflow error
	fmt.Println("Traveler no.", t.id, "moved outside the board")
}

func (t *Traveler) Delay() time.Duration {
	delayTime := time.Duration(t.generator.Intn(config.MaxDelay-config.MinDelay)+config.MinDelay) * time.Nanosecond
	time.Sleep(delayTime)
	return config.MaxDelay - delayTime // returns how quickly a traveler needs to proceed
}

func (t *Traveler) GetDelayDuration() time.Duration {
	delayTime := time.Duration(t.generator.Intn(config.MaxDelay-config.MinDelay)+config.MinDelay) * time.Nanosecond
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

func (t *Traveler) getNeighbours(omitX int, omitY int) []Position {
	var positions []Position
	var newPos Position

	// up
	newPos = t.GetPosition()
	newPos.y = (newPos.y + 1) % config.GridHeight
	if newPos.y != omitY {
		positions = append(positions, newPos)
	}

	// down
	newPos = t.GetPosition()
	newPos.y = (newPos.y - 1 + config.GridHeight) % config.GridHeight
	if newPos.y != omitY {
		positions = append(positions, newPos)
	}

	// left
	newPos = t.GetPosition()
	newPos.x = (newPos.x - 1 + config.GridWidth) % config.GridWidth
	if newPos.x != omitX {
		positions = append(positions, newPos)
	}

	// right
	newPos = t.GetPosition()
	newPos.x = (newPos.x + 1) % config.GridWidth
	if newPos.x != omitX {
		positions = append(positions, newPos)
	}

	return positions
}

func getNeighbours(pos Position, omitX int, omitY int) []Position {
	var positions []Position
	var newPos Position

	// up
	newPos = Position{
		x: pos.x,
		y: pos.y,
	}
	newPos.y = (newPos.y + 1) % config.GridHeight
	if newPos.y != omitY {
		positions = append(positions, newPos)
	}

	// down
	newPos = Position{
		x: pos.x,
		y: pos.y,
	}
	newPos.y = (newPos.y - 1 + config.GridHeight) % config.GridHeight
	if newPos.y != omitY {
		positions = append(positions, newPos)
	}

	// left
	newPos = Position{
		x: pos.x,
		y: pos.y,
	}
	newPos.x = (newPos.x - 1 + config.GridWidth) % config.GridWidth
	if newPos.x != omitX {
		positions = append(positions, newPos)
	}

	// right
	newPos = Position{
		x: pos.x,
		y: pos.y,
	}
	newPos.x = (newPos.x + 1) % config.GridWidth
	if newPos.x != omitX {
		positions = append(positions, newPos)
	}

	return positions
}

func (t *Traveler) Kill() {
	t.RuneSymbolToLowerCase()
	_ = t.addTrace() // discard overflow error
}

func (t *Traveler) SetPosition(x int, y int) {
	t.pos.x = x
	t.pos.y = y
}

func (t *Traveler) KillAt(timestamp time.Time) {
	t.RuneSymbolToLowerCase()
	_ = t.addTraceAt(timestamp) // discard overflow error
}

func nowInNanoseconds() time.Duration {
	return time.Duration(time.Now().UnixNano())
}
