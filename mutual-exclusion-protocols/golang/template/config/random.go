package config

////////////////////////////////////
//
// define the source of randomness
//
////////////////////////////////////

import (
	"math/rand"
	"time"
)

type Generator struct {
	rng *rand.Rand
}

func getSeed(id int) rand.Source {
	unixTime := time.Now().UnixNano()

	return rand.New(rand.NewSource(
		unixTime + (unixTime % (100 * int64(id+1)) * int64(id+1))),
	)
	//return 42
}

// MakeGenerators creates an array of independent sources of randomness
func MakeGenerators() [NoOfProcesses]*Generator {
	var generators [NoOfProcesses]*Generator
	for i := range NoOfProcesses {
		generators[i] = new(Generator)
		generators[i].rng = rand.New(getSeed(i))
	}
	return generators
}

func NewGenerator() *Generator {
	generator := new(Generator)
	generator.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
	return generator
}

func (g Generator) Intn(bound int) int {
	return g.rng.Intn(bound)
}
