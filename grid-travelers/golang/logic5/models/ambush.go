package models

import (
	"time"
)

// AmbushSeizeRequestType is a type of request
// that is sent by the tile server
// to inform the ambush to store its position.
// This is done solely for animation purposes
type AmbushSeizeRequestType struct {
	changeTimestamp time.Time
}

type AmbushSeizeRequestChannel struct {
	Channel chan AmbushSeizeRequestType
}

type AmbushSeizeRequestChannels struct {
	Channels map[Position]*AmbushSeizeRequestChannel
}

func RunAmbush(ambush *Traveler, ambushSeizeRequests *AmbushSeizeRequestChannel) {
	for {
		// Wait for the ambush to be seized
		// and reprint the ambush symbol
		ambushSeizeRequest := <-ambushSeizeRequests.Channel
		_ = ambush.addTraceAt(ambushSeizeRequest.changeTimestamp)
	}
}
