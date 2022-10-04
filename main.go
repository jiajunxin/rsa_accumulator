package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
	"github.com/jiajunxin/rsa_accumulator/precompute"
)

func testPreCompute() {
	setSize := 10000
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	startingTime := time.Now().UTC()
	prod := accumulator.SetProductRecursiveFast(rep)
	var duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running SetProductRecursive Takes [%.3f] Seconds \n",
		duration.Seconds())

	elementUpperBound := new(big.Int).Lsh(big.NewInt(1), 2047)
	startingTime = time.Now().UTC()
	t := precompute.NewTable(setup.G, setup.N, elementUpperBound, uint64(setSize), 1024)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running NewTable Takes [%.3f] Seconds \n", duration.Seconds())

	startingTime = time.Now().UTC()
	t.Compute(prod, 8)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running ComputeFromTable Takes [%.3f] Seconds \n", duration.Seconds())
}

func main() {

	testPreCompute()

}
