package main

import (
	"github.com/jiajunxin/rsa_accumulator/experiments"
)

func main() {
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
}
