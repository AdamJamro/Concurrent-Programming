package models

import (
	"grid-travelers-v2/config"
	"time"
)

type MoveType int

const (
	up MoveType = iota
	down
	left
	right
)

type position struct {
	x int
	y int
}

type TraceType struct {
	timeStamp time.Duration
	id        int
	pos       position
	symbol    rune
}

type TraceArray [config.MaxSteps]TraceType

type TraceSequenceType struct {
	len  int
	data TraceArray
}
