package main

import (
	"fmt"

	"github.com/rsa_accumulator/accumulator"
)

func main() {
	fmt.Println("start test in main")
	testSizes := []int{1000, 10000, 1000000}
	for _, size := range testSizes {
		fmt.Printf("test size: %d\n", size)
		accumulator.ManualBench(size)
		accumulator.ManualBenchParallel(size)
		accumulator.ManualBenchIterParallel(size)
		fmt.Println()
	}
	// {
	// 	set := accumulator.GenBenchSet(1000)
	// 	setup := *accumulator.TrustedSetup()
	// 	rep := accumulator.GenRepersentatives(set, accumulator.DIHashFromPoseidon)
	// 	defer profile.Start(profile.TraceProfile).Stop()
	// 	accumulator.ProveMembershipIterParallel(*setup.G, setup.N, rep)
	// }
	fmt.Println("end test in main")
}
