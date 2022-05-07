package main

import (
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/rsa_accumulator/accumulator"
)

func CalSetProd(inputSet []*big.Int) *big.Int {
	setsize := len(inputSet)
	prod := big.NewInt(1)
	tempSet := make([]*big.Int, 1)
	for i := 0; i < setsize; i++ {
		tempSet[0] = inputSet[i]
		temp, err := poseidon.Hash(tempSet)
		if err != nil {
			panic(err)
		}
		prod.Mul(prod, temp)
	}
	return prod
}

// PoKE is an implementation of the PoKE algorithm in "Batching Techniques for Accumulators with Applications to IOPs and Stateless Blockchains"
func PoKE(base, exp, newAcc, N *big.Int) {
	l := accumulator.HashToPrime(append(newAcc.Bytes(), base.Bytes()...))
	fmt.Println("primeChallenge = ", l.String())
	remainder := big.NewInt(1)
	quotient := big.NewInt(1)
	quotient, remainder = quotient.DivMod(exp, l, remainder)
	Q := accumulator.AccumulateNew(base, quotient, N)
	fmt.Println("Q = ", Q.String())
	fmt.Println("r = ", remainder.String())
	AccTest1 := accumulator.AccumulateNew(Q, l, N)
	fmt.Println("Q^l = ", AccTest1.String())
	AccTest2 := accumulator.AccumulateNew(base, remainder, N)
	fmt.Println("g^r = ", AccTest2.String())
	AccTest3 := AccTest1.Mul(AccTest1, AccTest2)
	AccTest3.Mod(AccTest3, N)
	fmt.Println("Accumulator_test3 = ", AccTest3.String()) //AccTest3 should be the same as the AccOld
}

func main() {
	fmt.Println("start test in main")
	setup := accumulator.TrustedSetup()
	setSize := 1000
	// generate Accumulatorset
	inputSet := make([]*big.Int, setSize)
	// i start from 1 to avoid 0 as input. 0 is already reserved for min value.
	for i := 0; i < setSize; i++ {
		inputSet[i] = big.NewInt(int64(10000 + i*2))
	}

	// get the product of the set
	prodOfSet := CalSetProd(inputSet)
	Acc, r := accumulator.ZKAccExp(setup.G, prodOfSet, setup.N)
	prodOfSet.Mul(prodOfSet, r)

	fmt.Println("\nAccumulator = ", Acc.String())

	fmt.Println("PoKE1 ")
	PoKE(setup.G, prodOfSet, Acc, setup.N)

	// l := accumulator.HashToPrime(append(AccNew.Bytes(), AccMid.Bytes()...))
	// r := big.NewInt(1)
	// for i, v := range addListSet {
	// 	var temp big.Int
	// 	temp.Mod(v, l)
	// 	//fmt.Println("the", i, "th element mod l = ", temp.String())
	// 	var temp2 big.Int
	// 	temp2.Mul(r, &temp)
	// 	fmt.Println("the", i, "th element = ", temp2.String())
	// 	r.Mod(&temp2, l)
	// 	fmt.Println("the", i, "th r-element = ", r.String())
	// }
}
