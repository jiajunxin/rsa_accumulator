package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
	"github.com/jiajunxin/rsa_accumulator/zkmultiswap"
)

func main() {

	testSetSize := uint32(10)
	startingTime := time.Now().UTC()
	zkmultiswap.SetupZkMultiswap(testSetSize)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generating a SNARK circuit for set size = %d, takes [%.3f] Seconds \n", testSetSize, duration.Seconds())
	runtime.GC()
	testSet := zkmultiswap.GenTestSet(testSetSize, accumulator.TrustedSetup())
	startingTime = time.Now().UTC()
	proof, publicWitness, err := zkmultiswap.Prove(testSet)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generating a SNARK proof for set size = %d, takes [%.3f] Seconds \n", testSetSize, duration.Seconds())
	if err != nil {
		fmt.Println("Error during Prove")
		panic(err)
	}

	flag := zkmultiswap.Verify(proof, testSetSize, publicWitness)
	if flag {
		fmt.Println("Verification passed")
	}
	//zkmultiswap.TestMultiSwap()
}

// experiments.TestBasicZKrsa()
// experiments.TestDifferentGroupSize()
// experiments.TestDifferentMembership()
// experiments.TestDifferentMembershipForDI()
// experiments.TestDifferentPrecomputationTableSize()
// experiments.TestFirstLayerPercentage()
// experiments.TestMembership()
// experiments.TestProduct()
// experiments.TestProduct2()
// experiments.TestProduct3()
// experiments.TestRange()
// experiments.TestPoKE()
