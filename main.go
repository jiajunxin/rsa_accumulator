package main

import (
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/rsa_accumulator/accumulator"
)

func GenSortedListSet(inputList []*big.Int) []*big.Int {
	var err error
	setsize := len(inputList)
	min := big.NewInt(0)
	// we set max to 0x7FFFFFFF instead of 0xFFFFFFFF to suit the xjSNARK, it seems there is some problem in there
	max := big.NewInt(0x7FFFFFFF)

	poseidonHashResult := make([]*big.Int, setsize+1)
	tempHashInput := make([]*big.Int, 2)
	tempHashInput[0] = min
	tempHashInput[1] = inputList[0]

	poseidonHashResult[0], err = poseidon.Hash(tempHashInput)
	if err != nil {
		// not expecting error from a Hash function
		panic(err)
	}

	for i := 1; i < setsize; i++ {
		tempHashInput[0] = inputList[i-1]
		tempHashInput[1] = inputList[i]
		poseidonHashResult[i], err = poseidon.Hash(tempHashInput)
		if err != nil {
			// not expecting error from a Hash function
			panic(err)
		}
	}

	tempHashInput[0] = inputList[setsize-1]
	tempHashInput[1] = max
	poseidonHashResult[setsize], err = poseidon.Hash(tempHashInput)
	if err != nil {
		// not expecting error from a Hash function
		panic(err)
	}
	return poseidonHashResult
}

// genUpdateListSet is simply used for testing.
func genUpdateListSet(inputList []*big.Int) []*big.Int {
	var err error
	setsize := len(inputList)

	poseidonHashResult := make([]*big.Int, setsize-1)
	tempHashInput := make([]*big.Int, 2)

	for i := 0; i < setsize-1; i++ {
		tempHashInput[0] = inputList[i]
		tempHashInput[1] = inputList[i+1]
		poseidonHashResult[i], err = poseidon.Hash(tempHashInput)
		if err != nil {
			// not expecting error from a Hash function
			panic(err)
		}
	}
	return poseidonHashResult
}

func CalSetProd(inputSet []*big.Int) *big.Int {
	setsize := len(inputSet)
	prod := big.NewInt(1)
	for i := 0; i < setsize; i++ {
		prod.Mul(prod, inputSet[i])
	}
	return prod
}

func main() {
	fmt.Println("start test in main")
	setup := accumulator.TrustedSetup()
	oldSetSize := 10000
	delSetSize := 16
	//	addSetSize := 32

	// generate Accumulator_old
	oldSortedList := make([]*big.Int, oldSetSize)
	// i start from 1 to avoid 0 as input. 0 is already reserved for min value.
	for i := 0; i < oldSetSize; i++ {
		oldSortedList[i] = big.NewInt(int64(1 + i*2))
	}
	oldSortedListSet := GenSortedListSet(oldSortedList)
	fmt.Println("sortedListSet1.len = ", len(oldSortedListSet))
	// get the product of the set
	prodOldSet := CalSetProd(oldSortedListSet)
	AccOld := accumulator.AccumulateNew(setup.G, prodOldSet, setup.N)

	fmt.Println("Accumulator_old = ", AccOld.String())
	// Generate delete set
	// Note that delList do not have to be continuous, we use this set just for example.
	delList := make([]*big.Int, delSetSize)
	// delList start from 1 to avoid 0 as input. 0 is already reserved for min value.
	for i := 0; i < delSetSize; i++ {
		delList[i] = big.NewInt(int64(1 + i*2))
	}
	// delListSet should be a subset of oldSortedListSet
	delListSet := genUpdateListSet(delList)
	prodDelSet := CalSetProd(delListSet)
	var prodMidSet big.Int
	prodMidSet.Div(prodOldSet, prodDelSet)
	AccMid := accumulator.AccumulateNew(setup.G, &prodMidSet, setup.N)
	fmt.Println("Accumulator_mid = ", AccMid.String())

	AccTest := accumulator.AccumulateNew(AccMid, prodDelSet, setup.N)
	fmt.Println("AccTest = ", AccTest.String()) //AccTest should be the same as the AccOld

	l1 := accumulator.HashToPrime(append(AccOld.Bytes(), AccMid.Bytes()...))
	fmt.Println("primeChallenge = ", l1.String())

	r1 := big.NewInt(1)
	q1 := big.NewInt(1)
	q1, r1 = q1.DivMod(prodDelSet, l1, r1)
	Q1 := accumulator.AccumulateNew(AccMid, q1, setup.N)
	fmt.Println("Q1 = ", Q1.String())
	fmt.Println("r1 = ", r1.String())

	AccTest1 := accumulator.AccumulateNew(Q1, l1, setup.N)
	AccTest2 := accumulator.AccumulateNew(AccMid, r1, setup.N)
	AccTest3 := AccTest1.Mul(AccTest1, AccTest2)
	AccTest3.Mod(AccTest3, setup.N)
	fmt.Println("Accumulator_test = ", AccTest3.String()) //AccTest3 should be the same as the AccOld

	// Generate delete set
	// Note that delList do not have to be continuous, we use this set just for example.
	// addList := make([]*big.Int, addSetSize)
	// // addList start from 1 to avoid 0 as input. 0 is already reserved for min value.
	// for i := 0; i < addSetSize; i++ {
	// 	addList[i] = big.NewInt(int64(1 + i))
	// }
	// // delList should be a subset of addList
	// addListSet := genUpdateListSet(addList)
	// prodAddSet := CalSetProd(addListSet)

}
