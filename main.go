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

	//zkmultiswap.TestMultiSwap()
	// fmt.Println("TestBasicZKrsa")
	// experiments.TestBasicZKrsa()
	// fmt.Println("TestDifferentGroupSize")

	// experiments.TestDifferentMembership()

	// fmt.Println("TestDifferentPrecomputationTableSize")
	// experiments.TestDifferentPrecomputationTableSize()
	// fmt.Println("TestMembership")
	// experiments.TestMembership()
	// fmt.Println("TestPoKE")
	// experiments.TestPoKE()
	// fmt.Println("TestNotusSingleThread")
	// experiments.TestNotusSingleThread(1024, 100)
}
