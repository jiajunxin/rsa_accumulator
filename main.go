package main

import (
	"fmt"

	"github.com/rsa_accumulator/accumulator"
)

func main() {
	// defer profile.Start(profile.TraceProfile).Stop()
	testSize := 1000
	fmt.Printf("test in main, test size: %d\n", testSize)
	// accumulator.ManualBench(testSize)
	accumulator.ManualBenchParallel(testSize)
	accumulator.ManualBenchIterParallel(testSize)
}
