package main

import (
	"fmt"

	"github.com/jiajunxin/rsa_accumulator/experiments"
)

func main() {

	// updateRates denotes the percentage of updates in the total users, i.e. number of updates = users/updateRates
	updateRates := 1024

	twoTo14 := 16384
	twoTo15 := 32768
	twoTo16 := 65536
	twoTo17 := 131072
	twoTo18 := 262144
	twoTo19 := 524288
	twoTo20 := 1048576

	experiments.TestMembershipVerify()

	// the following code is to test the time percentage for each part of the system in a single core.
	fmt.Println("TestNotusMultiSwap for ", twoTo14, " users and ", uint32(twoTo14/updateRates), " updates.")
	experiments.TestNotusMultiSwap(uint32(twoTo14), uint32(twoTo14/updateRates))
	fmt.Println("TestNotusMultiSwap for ", twoTo15, " users and ", uint32(twoTo15/updateRates), " updates.")
	experiments.TestNotusMultiSwap(uint32(twoTo14), uint32(twoTo14/updateRates))
	fmt.Println("TestNotusMultiSwap for ", twoTo16, " users and ", uint32(twoTo16/updateRates), " updates.")
	experiments.TestNotusMultiSwap(uint32(twoTo14), uint32(twoTo14/updateRates))

	fmt.Println("Test Membership precomputation under different group size")
	experiments.TestDifferentMembershipForDISingleThread()

	fmt.Println("Test MultiSwap With Different Size")
	experiments.TestMultiSwapWithDifferentSize()

	fmt.Println("TestNotusMultiSwap for ", twoTo15, " users and ", uint32(twoTo15/updateRates), " updates in parallel")
	fmt.Println("Assuming 32 group partations for 2^20 user, each group of 2^15 users uses 32 cores, total 1024 cores")
	experiments.TestNotusParallel(uint32(twoTo15), uint32(twoTo15/updateRates))

	fmt.Println("TestNotusMultiSwap for ", twoTo16, " users and ", uint32(twoTo16/updateRates), " updates in parallel")
	fmt.Println("Assuming 16 group partations for 2^20 user, each group of 2^16 users uses 32 cores, total 512 cores")
	experiments.TestNotusParallel(uint32(twoTo16), uint32(twoTo16/updateRates))

	fmt.Println("TestNotusMultiSwap for ", twoTo17, " users and ", uint32(twoTo17/updateRates), " updates in parallel")
	fmt.Println("Assuming 8 group partations for 2^20 user, each group of 2^17 users uses 32 cores, total 256 cores")
	experiments.TestNotusParallel(uint32(twoTo17), uint32(twoTo17/updateRates))

	fmt.Println("TestNotusMultiSwap for ", twoTo18, " users and ", uint32(twoTo18/updateRates), " updates in parallel")
	fmt.Println("Assuming 4 group partations for 2^20 user, each group of 2^18 users uses 32 cores, total 128 cores")
	experiments.TestNotusParallel(uint32(twoTo18), uint32(twoTo18/updateRates))

	fmt.Println("TestNotusMultiSwap for ", twoTo19, " users and ", uint32(twoTo19/updateRates), " updates in parallel")
	fmt.Println("Assuming 4 group partations for 2^19 user, each group of 2^19 users uses 32 cores, total 64 cores")
	experiments.TestNotusParallel(uint32(twoTo19), uint32(twoTo19/updateRates))

	fmt.Println("TestNotusMultiSwap for ", twoTo20, " users and ", uint32(twoTo20/updateRates), " updates in parallel")
	fmt.Println("Assuming 1 group partation for 2^20 user, each group of 2^20 users uses 32 cores, total 32 cores")
	experiments.TestNotusParallel(uint32(twoTo20), uint32(twoTo20/updateRates))
}
