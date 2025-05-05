package config

import "time"

///////////////////////////////////////////
//
// define the parameters of the simulation
//
///////////////////////////////////////////

// grid dimensions

const GridWidth = 10
const GridHeight = 10

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
const NoOfAmbushes = 10

// GridFieldCapacity DEPRECATED
const GridFieldCapacity = 1
