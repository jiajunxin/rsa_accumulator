package benchmark

import (
	"math/big"
	"sync"
	"testing"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

const testSize = 1000

var (
	onceSet            sync.Once
	onceSetup          sync.Once
	onceRepresentation sync.Once
)

func getBenchSet() []string {
	var set []string
	onceSet.Do(func() {
		set = accumulator.GenBenchSet(testSize)
	})
	return set
}

func getSetup() *accumulator.Setup {
	var setup *accumulator.Setup
	onceSetup.Do(func() {
		setup = accumulator.TrustedSetup()
	})
	return setup
}

func getRepresentations() []*big.Int {
	var rep []*big.Int
	onceRepresentation.Do(func() {
		set := getBenchSet()
		rep = accumulator.GenRepresentatives(set, accumulator.HashToPrimeFromSha256)
	})
	return rep
}

func BenchmarkAccumulatorFirstLayer(b *testing.B) {
	//set := getBenchSet()
	//setup := getSetup()

}
