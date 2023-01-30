package experiments

import (
	"fmt"
	"time"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

// TestBasiczkRSA test a naive case of zero-knowledge RSA accumulator
func TestBasiczkRSA() {
	setSize := 65536 // 2 ^ 16 65536
	fmt.Println("Test set size = ", setSize)
	fmt.Println("GenRepresentatives with MultiDIHashFromPoseidon")
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	startingTime := time.Now().UTC()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	// generate a zero-knowledge RSA accumulator
	r := accumulator.GenRandomizer()
	randomizedbase := AccumulateNew(setup.G, r, setup.N)
	// calculate the exponentation
	exp := accumulator.SetProductRecursiveFast(rep)
	accumulator.AccumulateNew(randomizedbase, exp, setup.N)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generating a zero-knowledge RSA accumulator with set size = %d, takes [%.3f] Seconds \n", setSize, duration.Seconds())

	// startingTime = time.Now().UTC()
	// maxLen := setSize * 256 / bits.UintSize
	// table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	// duration = time.Now().UTC().Sub(startingTime)
	// fmt.Printf("Running PreCompute Takes [%.3f] Seconds \n", duration.Seconds())
	// startingTime = time.Now().UTC()
	// accumulator.ProveMembershipParallelWithTable(setup.G, setup.N, rep, 2, table)
	// duration = time.Now().UTC().Sub(startingTime)
	// fmt.Printf("Running ProveMembershipParallelWithTable Takes [%.3f] Seconds \n", duration.Seconds())
}
