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
	// generate product of the input sorted list
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
	setsize := 16

	// generate Accumulator_old
	oldSetsize := 10000
	sortedList0 := make([]*big.Int, oldSetsize)
	for i := 0; i < oldSetsize; i++ {
		sortedList0[i] = big.NewInt(int64(i))
	}
	sortedListSet1 := GenSortedListSet(sortedList0)
	fmt.Println("sortedListSet1.len = ", len(sortedListSet1))
	// get the product of the set
	prod1 := CalSetProd(sortedListSet1)

	setup := accumulator.TrustedSetup()
	AccOld := accumulator.AccumulateNew(setup.G, prod1, setup.N)

	fmt.Println("Accumulator_old = ", AccOld.String())
	// Generate test set
	// Suppopse we want to remove the first setsize elements from sortedListSet1
	// Note that sortedList1 do not have to be continuous, we use this set just for example.
	sortedListSet2 := sortedListSet1[setsize:]
	fmt.Println("sortedListSet2.len = ", len(sortedListSet2))
	prod2 := CalSetProd(sortedListSet2)

	AccMid := accumulator.AccumulateNew(setup.G, prod2, setup.N)

	fmt.Println("Accumulator_mid = ", AccMid.String())

	r1 := big.NewInt(1)
	q1, r1 := prod1.DivMod(prod1, prod2, r1)
	//fmt.Println("q1 = ", q1.String())
	fmt.Println("r1 = ", r1.String())

	AccTest := accumulator.AccumulateNew(AccMid, q1, setup.N)
	fmt.Println("Accumulator_test = ", AccTest.String())

	l1 := accumulator.HashToPrime(append(AccOld.Bytes(), AccMid.Bytes()...))
	fmt.Println("primeChallenge = ", l1.String())

	// prod1 is the product of all the hash result of sortedList0_1

	// calculate Q s.t. q1*l1 + r1 = prod1
	// r1 := big.NewInt(1)
	// q1, r1 := prod1.DivMod(prod1, l1, r1)
	// Q1 := accumulator.AccumulateNew(setup.G, q1, setup.N)
	// fmt.Println("Q1 = ", Q1)
	// fmt.Println("r1 = ", r1)

	// rem1 := big.NewInt(1)
	// tempDIHash := big.NewInt(1)
	// tempMod := big.NewInt(1)
	// tempPrint := big.NewInt(1)
	// for i := 0; i < setsize+1; i++ {
	// 	tempDIHash.Add(accumulator.Min2048, poseidonHashResult[i])
	// 	tempPrint.Set(tempDIHash)
	// 	tempMod = tempMod.Mod(tempDIHash, l1)
	// 	rem1.Mul(rem1, tempMod)
	// 	rem1.Mod(rem1, l1)
	// }

}
