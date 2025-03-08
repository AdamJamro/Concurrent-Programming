////////////////////////////////////////////////////
// implement simulator's source of random numbers //
////////////////////////////////////////////////////

package configuration

import (
	common "factorySimulator/commonModels"
	"math/rand"
	"time"
)

// getRandom...() use the Rng and getSeed() to generate random values
func getSeed() rand.Source {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
	//return 42
}

var Rng = rand.New(getSeed())

func GetRandomIntArgument() int {
	return Rng.Intn(Bound)
}

func GetRandomOperation() *common.Operation {
	// race condition may happen but since there's only one boss we don't care
	var randomIndex = Rng.Intn(len(Operations))
	return &Operations[randomIndex]
}

func GetBossDelay() time.Duration {
	// race condition may happen but since there's only one boss we don't care
	return time.Duration(rand.Intn(BossDelayUpperBound)) * time.Second
}
