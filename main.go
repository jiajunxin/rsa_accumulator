package main

import "github.com/jiajunxin/rsa_accumulator/experiments"

func main() {

	//zkmultiswap.TestMultiSwap()
	experiments.TestBasicZKrsa()
	experiments.TestDifferentGroupSize()
	experiments.TestDifferentMembership()
	experiments.TestDifferentMembershipForDI()
	experiments.TestDifferentPrecomputationTableSize()
	experiments.TestFirstLayerPercentage()
	experiments.TestMembership()
	experiments.TestProduct()
	experiments.TestProduct2()
	experiments.TestProduct3()
	experiments.TestRange()
	experiments.TestPoKE()
	experiments.TestNotusSingleThread(1024, 100)
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
