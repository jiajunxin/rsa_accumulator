package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/jiajunxin/rsa_accumulator/zkmultiswap"
)

func main() {
	//zkmultiswap.TestMimc()
	zkmultiswap.SetupZkMultiswap(1000)
	runtime.GC()
	startingTime := time.Now().UTC()
	zkmultiswap.Prove()
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generating a SNARK proof for set size = %d, takes [%.3f] Seconds \n", 100, duration.Seconds())
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
