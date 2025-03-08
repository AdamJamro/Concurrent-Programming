//////////////////////////////////////////////////
// define constant parameters of the simulation //
//////////////////////////////////////////////////

package configuration

import (
	"time"
)

var IsVerboseModeOn = true

// Bound for random integer arguments (input of operations)
const Bound = 4274

const WorkerDelay = 6 * time.Second
const ClientDelay = 8 * time.Second
const BossDelayUpperBound = 4 // in seconds

const NumOfWorkers = 5
const NumOfClients = 3

// SizeOfQueue maximum number of staged tasks waiting to be fulfilled
const SizeOfQueue = 2

// SizeOfTaskStorage maximum capacity of storage
const SizeOfTaskStorage = 5
