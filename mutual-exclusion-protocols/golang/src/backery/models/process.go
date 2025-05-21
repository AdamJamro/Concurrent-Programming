package models

import (
	"fmt"
	"mutex-protocols/backery/config"
	"time"
	"unicode"
)

type Process struct {
	id            int
	pos           Position
	symbol        rune
	noOfSteps     int
	timestamp     time.Time         // store creation time
	traceSequence TraceSequenceType // store all history
	generator     *config.Generator
}

type ProcessData struct {
	Id        int
	Symbol    rune
	Generator *config.Generator
	Position  Position
}

func MakeProcesses() *[config.NoOfProcesses]Process {
	generators := config.MakeGenerators()
	processes := [config.NoOfProcesses]Process{}

	for i, _ := range processes {
		if error := processes[i].InitializeProcess(
			ProcessData{
				Id:        i,
				Symbol:    rune('A' + i),
				Generator: generators[i],
				Position:  Position{x: i, y: int(config.LocalSection)},
			},
		); error != nil {
			fmt.Printf("Error initializing process #%d\n", i)
		}
	}

	return &processes
}

func (t *Process) InitializeProcess(data ProcessData) error {
	t.id = data.Id
	t.generator = data.Generator
	t.symbol = data.Symbol
	t.noOfSteps = t.generator.Intn(config.MaxSteps-config.MinSteps+1) + config.MinSteps
	t.pos = data.Position
	t.timestamp = time.Now()
	return t.addTraceAt(t.timestamp)
}

func (t *Process) addTrace() error {
	return t.addTraceAt(time.Now())
}

// addTrace saves current state and pushes it to the trace sequence at a specific timestamp
func (t *Process) addTraceAt(timestamp time.Time) error {
	trace := TraceType{
		timestamp: timestamp,
		id:        t.id,
		pos:       t.pos,
		symbol:    t.symbol,
	}
	return t.traceSequence.add(trace)
}

func (t *Process) GetPosition() Position {
	return Position{x: t.pos.x, y: t.pos.y}
}

// ChangeState unused (used in precious logic systems)
func (t *Process) ChangeState(state config.ProcessState) {
	t.pos.y = int(state)
	changeTime := time.Now()
	t.addTraceAt(changeTime)
}

func (t *Process) LogState() error {
	return t.addTrace()
}

func (t *Process) Delay() time.Duration {
	delayTime := time.Duration(t.generator.Intn(config.MaxDelay-config.MinDelay)+config.MinDelay) * time.Nanosecond
	time.Sleep(delayTime)
	return config.MaxDelay - delayTime // returns how quickly a traveler needs to proceed
}

func (t *Process) GetDelayDuration() time.Duration {
	delayTime := time.Duration(t.generator.Intn(config.MaxDelay-config.MinDelay)+config.MinDelay) * time.Nanosecond
	return delayTime
}

func (t *Process) GetNoOfSteps() int {
	return t.noOfSteps
}

func (t *Process) GetId() int {
	return t.id
}

func (t *Process) PrintReport() {
	t.traceSequence.PrintTrace()
}

func (t *Process) RuneSymbolToLowerCase() {
	t.symbol = unicode.ToLower(rune(t.symbol))
}

func (t *Process) KillAt(timestamp time.Time) {
	t.RuneSymbolToLowerCase()
	_ = t.addTraceAt(timestamp) // discard overflow error
}

func (t *Process) Kill() {
	t.RuneSymbolToLowerCase()
	_ = t.addTrace() // discard overflow error
}

func (t *Process) MoveOutsideTheBoard() {
	t.pos.x = config.GridWidth
	t.pos.y = config.GridHeight
	_ = t.addTrace() // discard overflow error
}

func (t *Process) SetPosition(x int, y int) {
	t.pos.x = x
	t.pos.y = y
}
