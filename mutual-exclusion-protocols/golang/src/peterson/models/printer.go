package models

import (
	"fmt"
	"mutex-protocols/peterson/config"
	"time"
)

type Position struct {
	x int
	y int
}

type TraceType struct {
	timestamp time.Time
	id        int
	pos       Position
	symbol    rune
}

type TraceArray [config.MaxSteps + 1]TraceType

type TraceSequenceType struct {
	len  int
	data TraceArray
}

func (t *TraceSequenceType) PrintTrace() {
	for i := 0; i < t.len; i++ {
		fmt.Printf(
			"%d %d %d %d %c\n",
			t.data[i].timestamp.UnixNano(),
			t.data[i].id,
			t.data[i].pos.x,
			t.data[i].pos.y,
			t.data[i].symbol,
		)
	}
}

func (t *TraceSequenceType) add(trace TraceType) error {
	if t.len > config.MaxSteps {
		return fmt.Errorf("error: TraceSequenceType is full")
	}
	t.data[t.len] = trace
	t.len++
	return nil
}
