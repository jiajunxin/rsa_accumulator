package main

import (
	"fmt"
	"math/big"

	"github.com/rsa_accumulator/accumulator"
)

func main() {
	fmt.Println("test in main")
	var testObject accumulator.AccumulatorSetup
	testObject = *accumulator.Init()

	var N big.Int
	N.Mul(&testObject.P, &testObject.Q)

	fmt.Println("p = ", testObject.P.String())
	fmt.Println("q = ", testObject.Q.String())
	fmt.Println("N = ", testObject.N.String())
	fmt.Println("G = ", testObject.G.String())
}
