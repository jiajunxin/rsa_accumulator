package main

import (
	"fmt"

	"github.com/pkg/profile"

	"github.com/rsa_accumulator/accumulator"
)

func main() {
	defer profile.Start(profile.MemProfile).Stop()
	fmt.Println("test in main")
	// accumulator.ManualBench(10000)
	accumulator.ManualBenchIter(10000)
}
