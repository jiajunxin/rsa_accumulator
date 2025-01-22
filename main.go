package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/jiajunxin/rsa_accumulator/experiments"
	"github.com/jiajunxin/rsa_accumulator/zkmultiswap"
)

const (
	twoTo14 = 16384
	twoTo15 = 32768
	twoTo16 = 65536
	twoTo17 = 131072
	twoTo18 = 262144
	twoTo19 = 524288
)

func testMembershipproof() {
	// test Membership proof Verification and proof size
	experiments.TestMembershipVerify()
	size := 1000
	startingTime := time.Now().UTC()
	proofs := make([]int, size)
	for i := 0; i < size; i++ {
		proofs[i] = 1
	}
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallelWithTableWithRandomizer with single core for an RSA accumulator Takes [%.3f] Seconds \n", duration.Seconds())
	startingTime = time.Now().UTC()
	tempProof := proofs[0]
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Println("t = ", tempProof)
	fmt.Printf("Online phase to get one membership proof (retrieve it from memory), Takes [%d] Nanoseconds \n", duration.Nanoseconds())
}

func testbasicprocess() {
	experiments.TestPoKE()
	zkmultiswap.TestMultiSwap(10)
	// outputs a Solidity smart
	zkmultiswap.TestMultiSwapAndOutputSmartContract(10)
	err := zkmultiswap.TestMultiSwapAndOutputSmartContract2(10)
	if err != nil {
		panic(err)
	}
}

func singleCoreComponentProfiler(updateRates int) {
	// The code below measures and reports the proportion of execution time consumed by each component of the system when run on a single processor core.
	fmt.Println("TestNotusMultiSwap for ", twoTo14, " users and ", uint32(twoTo14/updateRates), " updates.")
	experiments.TestNotusMultiSwap(uint32(twoTo14), uint32(twoTo14/updateRates))
	fmt.Println("TestNotusMultiSwap for ", twoTo15, " users and ", uint32(twoTo15/updateRates), " updates.")
	experiments.TestNotusMultiSwap(uint32(twoTo15), uint32(twoTo15/updateRates))
	fmt.Println("TestNotusMultiSwap for ", twoTo16, " users and ", uint32(twoTo16/updateRates), " updates.")
	experiments.TestNotusMultiSwap(uint32(twoTo16), uint32(twoTo16/updateRates))
}

func testNotusParallel(updateRates int) {
	fmt.Println("Test Notus With Different Size")
	experiments.TestMultiSwapWithDifferentSize()
	runtime.GC()
	fmt.Println("TestNotusMultiSwap for ", twoTo15, " users and ", uint32(twoTo15/updateRates), " updates in parallel")
	fmt.Println("Assuming 32 group partations for 2^20 user, each group of 2^15 users uses 32 cores, total 1024 cores")
	experiments.TestNotusParallel(uint32(twoTo15), uint32(twoTo15/updateRates))
	experiments.TestRSAMembershipPreComputeDIParallel(32768, 5) //2^15, 32 cores
	experiments.TestRSAMembershipPreComputeDIParallel(65536, 5) //2^16, 32 cores
	runtime.GC()
	fmt.Println("TestNotusMultiSwap for ", twoTo16, " users and ", uint32(twoTo16/updateRates), " updates in parallel")
	fmt.Println("Assuming 16 group partations for 2^20 user, each group of 2^16 users uses 32 cores, total 512 cores")
	experiments.TestNotusParallel(uint32(twoTo16), uint32(twoTo16/updateRates))
	runtime.GC()
	fmt.Println("TestNotusMultiSwap for ", twoTo17, " users and ", uint32(twoTo17/updateRates), " updates in parallel")
	fmt.Println("Assuming 8 group partations for 2^20 user, each group of 2^17 users uses 32 cores, total 256 cores")
	experiments.TestNotusParallel(uint32(twoTo17), uint32(twoTo17/updateRates))
	runtime.GC()
	fmt.Println("TestNotusMultiSwap for ", twoTo18, " users and ", uint32(twoTo18/updateRates), " updates in parallel")
	fmt.Println("Assuming 4 group partations for 2^20 user, each group of 2^18 users uses 32 cores, total 128 cores")
	experiments.TestNotusParallel(uint32(twoTo18), uint32(twoTo18/updateRates))
	runtime.GC()
	fmt.Println("TestNotusMultiSwap for ", twoTo19, " users and ", uint32(twoTo19/updateRates), " updates in parallel")
	fmt.Println("Assuming 2 group partations for 2^20 user, each group of 2^19 users uses 32 cores, total 64 cores")
	experiments.TestNotusParallel(uint32(twoTo19), uint32(twoTo19/updateRates))
	runtime.GC()
}

func main() {
	var totalSize uint32 = twoTo14
	var updatedSetSize uint32 = 1 << 10
	powNumCores := 3
	println("Total Size: ", totalSize, " Updated Set Size: ", updatedSetSize, " Number of Cores: ", 1<<powNumCores)
	experiments.OKXBenchParallel(totalSize, updatedSetSize, powNumCores)
	//println("Test for ", twoTo15, " total users, with 1024 transaction in this batch, with 2^", 3, " cores")
	// experiments.TestRSASubsetParallel(twoTo15, 1024, 3)
	//println("Test for ", twoTo16, " total users, with 1024 transaction in this batch, with 2^", 3, " cores")
	// experiments.TestRSASubsetParallel(twoTo16, 1024, 3)
	//println("Test for ", twoTo17, " total users, with 1024 transaction in this batch, with 2^", 3, " cores")
	// experiments.TestRSASubsetParallel(twoTo17, 1024, 3)
	//println("Test for ", twoTo18, " total users, with 1024 transaction in this batch, with 2^", 3, " cores")
	// experiments.TestRSASubsetParallel(twoTo18, 1024, 3)
}
