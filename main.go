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
	target.SetString("409642780809890809878787847837483748374839283748329747837482973874893784732", 10)
	fmt.Println(target.BitLen())
	start := time.Now()
	res, err := proof.LagrangeFourSquares(target)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start))
	fmt.Println(res)
	if proof.Verify(target, res.W1, res.W2, res.W3, res.W4) {
		fmt.Println("verify success")
	} else {
		fmt.Println("verify failed")
	}
	//fmt.Println("end test in main")
}
