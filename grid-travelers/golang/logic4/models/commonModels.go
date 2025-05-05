package models

import (
	"grid-travelers-v4/config"
	"time"
)

type MoveType int

const (
	Up MoveType = iota
	Down
	Left
	Right
)

type Position struct {
	x int
	y int
}

type TraceType struct {
	timeStamp time.Time
	id        int
	pos       Position
	symbol    rune
}

type TraceArray [config.MaxSteps + 1]TraceType

type TraceSequenceType struct {
	len  int
	data TraceArray
}
