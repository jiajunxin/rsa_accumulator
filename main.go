package main

import (
	"fmt"
	"github.com/jiajunxin/rsa_accumulator/experiments"
	"math/big"
	"time"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
	"github.com/jiajunxin/rsa_accumulator/precompute"
)

func testPreCompute() {
	setSize := 100000
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)

	elementUpperBound := new(big.Int).Lsh(big.NewInt(1), 2047)
	startingTime := time.Now().UTC()
	t := precompute.NewTable(setup.G, setup.N, elementUpperBound, uint64(setSize), 102400)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running NewTable Takes [%.3f] Seconds \n", duration.Seconds())

	startingTime = time.Now().UTC()
	experiments.ProveMembershipParallel2(t, setup.G, setup.N, rep, 4, 16)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallel2 Takes [%.3f] Seconds \n", duration.Seconds())
}

func main() {

	testPreCompute()

}
