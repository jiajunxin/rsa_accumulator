package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
	"github.com/jiajunxin/rsa_accumulator/precompute"
)

func testPreCompute() {
	setSize := 1000
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	startingTime := time.Now().UTC()
	prod := accumulator.SetProductRecursiveFast(rep)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running SetProductRecursive Takes [%.3f] Seconds \n",
		duration.Seconds())

	startingTime = time.Now().UTC()
	originalResult := accumulator.AccumulateNew(setup.G, prod, setup.N)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccumulateNew Takes [%.3f] Seconds \n",
		duration.Seconds())

	startingTime = time.Now().UTC()
	table := precompute.GenPreTable(setup.G, setup.N, setSize*2048, 1024)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running GenPreTable Takes [%.3f] Seconds \n",
		duration.Seconds())
	//precompute.PrintTable(table)
	fmt.Println(" ")
	fmt.Println(" ")
	fmt.Println(" ")
	startingTime = time.Now().UTC()
	result := precompute.ComputeFromTableParallel(table, prod, setup.N)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running ComputeFromTable Takes [%.3f] Seconds \n",
		duration.Seconds())

	if result.Cmp(originalResult) != 0 {
		fmt.Println("wrong result")
	}

	elementUpperBound := new(big.Int).Lsh(big.NewInt(1), 2048)
	elementUpperBound.Sub(elementUpperBound, big.NewInt(1))
	startingTime = time.Now().UTC()
	table1 := precompute.NewTable(setup.G, setup.N, elementUpperBound, uint64(setSize), 1024)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running NewTable Takes [%.3f] Seconds \n", duration.Seconds())

	startingTime = time.Now().UTC()
	result1 := table1.Compute(prod, 8)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running ComputeFromTable Takes [%.3f] Seconds \n", duration.Seconds())

	if result1.Cmp(originalResult) != 0 {
		fmt.Println("wrong result")
	}
}

func main() {

	testPreCompute()

}
