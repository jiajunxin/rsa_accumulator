package main

import (
	"fmt"

	"github.com/rsa_accumulator/accumulator"
)

func main() {
	// defer profile.Start(profile.TraceProfile).Stop()
	fmt.Println("test in main")
	// accumulator.ManualBench(1000)
	accumulator.ManualBenchParallel(1000000)
	accumulator.ManualBenchIterParallel(1000000)
}
