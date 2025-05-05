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
func MakeGenerators(len int) [NoOfTravelers]*Generator {
	var generators [NoOfTravelers]*Generator
	for i := range len {
		generators[i] = new(Generator)
		generators[i].rng = rand.New(getSeed(i))
	}
	return generators
}

func (g Generator) Intn(bound int) int {
	return g.rng.Intn(bound)
}

func (g Generator) GenerateRandomInitialPositions() [NoOfTravelers]struct{ X, Y int } {
	cols := GridWidth
	rows := GridHeight

	// use shuffle for uniform distribution

	positions := [NoOfTravelers]struct{ X, Y int }{}
	permutation := g.rng.Perm(cols * rows)
	for i := range positions {
		codedPosition := permutation[i]
		decodedPosition := struct{ X, Y int }{X: codedPosition % cols, Y: codedPosition / cols}
		positions[i] = struct{ X, Y int }{X: decodedPosition.X, Y: decodedPosition.Y}
	}
	return positions
}
