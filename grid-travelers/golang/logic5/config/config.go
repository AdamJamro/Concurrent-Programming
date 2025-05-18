package config

import "time"

///////////////////////////////////////////
//
// define the parameters of the simulation
//
///////////////////////////////////////////

// grid dimensions

const GridWidth = 15
const GridHeight = 15

// movement

const NoOfTravelers = 15

const MinSteps = 10
const MaxSteps = 15

const MinDelay = 1000000                 //nanoseconds
const MaxDelay = 5000000                 //nanoseconds
const MinTimeout = 50 * time.Microsecond //nanoseconds

// wild travelers
const millisecond = 1000000
const MaxWildlifeSpawn = 9
const SpawnRate = (MaxDelay-MinDelay)*0.3 + MinDelay
const MinWildLifetime = 80 * millisecond
const MaxWildLifetime = 150 * millisecond

// ambushes
const AmbushTimeout = 10 * millisecond
const NoOfAmbushes = 15

// GridFieldCapacity DEPRECATED
const GridFieldCapacity = 1

func ValidateConfiguration() {
	if GridWidth <= 0 {
		panic("GridWidth must be greater than 0")
	}
	if GridHeight <= 0 {
		panic("GridHeight must be greater than 0")
	}
	if NoOfTravelers <= 0 {
		panic("NoOfTravelers must be greater than 0")
	}
	if MinSteps <= 0 {
		panic("MinSteps must be greater than 0")
	}
	if MaxSteps <= 0 {
		panic("MaxSteps must be greater than 0")
	}
	if MinSteps >= MaxSteps {
		panic("MinSteps must be strictly less than MaxSteps")
	}
	if GridHeight*GridWidth < NoOfTravelers+NoOfAmbushes {
		panic("Board must be greater than NoOfTravelers and NoOfAmbushes in order to contain them")
	}
}
