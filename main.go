package main

import (
	"fmt"

	"math/big"

	"github.com/rsa_accumulator/proof"
)

func main() {
	fmt.Println("start test in main")
	//testSizes := []int{1000}
	//for _, size := range testSizes {
	//	fmt.Printf("test size: %d\n", size)
	//	//accumulator.ManualBench(size)
	//	accumulator.ManualBenchZKAcc(size)
	//	fmt.Println()
	//}

	// {
	// 	set := accumulator.GenBenchSet(1000)
	// 	setup := *accumulator.TrustedSetup()
	// 	rep := accumulator.GenRepersentatives(set, accumulator.DIHashFromPoseidon)
	// 	defer profile.Start(profile.TraceProfile).Stop()
	// 	accumulator.ProveMembershipIterParallel(*setup.G, setup.N, rep)
	// }

	res, err := proof.Lagrange(big.NewInt(8))
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	//fmt.Println("end test in main")
}
