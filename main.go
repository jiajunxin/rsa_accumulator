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
	target.SetString("127521675156751657516572156572165721657216721671276427624672167216572416574216724165777547", 10)
	fmt.Println(target.BitLen())
	start := time.Now()
	fmt.Println(proof.LagrangeFourSquaresPollack(target))
	fmt.Println(time.Since(start))
	//fmt.Println("end test in main")
}
