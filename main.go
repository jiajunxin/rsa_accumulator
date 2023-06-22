package main

import (
	"fmt"

	"github.com/jiajunxin/rsa_accumulator/experiments"
)

func main() {

	//zkmultiswap.TestMultiSwap()
	fmt.Println("TestBasicZKrsa")
	experiments.TestBasicZKrsa()
	fmt.Println("TestDifferentGroupSize")
	experiments.TestDifferentGroupSize()
	fmt.Println("TestDifferentMembership")
	experiments.TestDifferentMembership()
	fmt.Println("TestDifferentMembershipForDI")
	experiments.TestDifferentMembershipForDI()
	fmt.Println("TestDifferentPrecomputationTableSize")
	experiments.TestDifferentPrecomputationTableSize()
	fmt.Println("TestMembership")
	experiments.TestMembership()
	fmt.Println("TestPoKE")
	experiments.TestPoKE()
	fmt.Println("TestNotusSingleThread")
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
