package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/rsa_accumulator/proof"
)

func main() {
	fmt.Println("start test in main")
	//testSizes := []int{1000}
	//for _, size := range testSizes {
	//	fmt.Printf("test size: %d\n", size)
	//   accumulator.ManualBench(size)
	//	accumulator.ManualBenchZKAcc(size)
	//	fmt.Println()
	//}

	{
		// set := accumulator.GenBenchSet(1000)
		// setup := *accumulator.TrustedSetup()
		// rep := accumulator.GenRepersentatives(set, accumulator.DIHashFromPoseidon)
		// 	defer profile.Start(profile.TraceProfile).Stop()
		// accumulator.ProveMembershipIterParallel(*setup.G, setup.N, rep)
	}

	target := new(big.Int)
	target.SetString("4096427", 10)
	fmt.Println(target.BitLen())
	start := time.Now()
	fmt.Println(proof.LagrangeFourSquares(target))
	fmt.Println(time.Since(start))
	//fmt.Println("end test in main")
}
