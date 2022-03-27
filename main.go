package main

import (
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/rsa_accumulator/accumulator"
)

func diHash(input []*big.Int) *big.Int {
	if len(input) != 2 {
		panic("diHash requires 2 inputs")
	}
	ret, err := poseidon.Hash(input)
	if err != nil {
		panic(err)
	}
	ret.Add(ret, accumulator.Min2048)
	return ret
}

func GenSortedListSet(inputList []*big.Int) []*big.Int {
	setsize := len(inputList)
	min := big.NewInt(0)
	// we set max to 0x7FFFFFFF instead of 0xFFFFFFFF to suit the xjSNARK, it seems there is some problem in there
	max := big.NewInt(0x7FFFFFFF)

	poseidonHashResult := make([]*big.Int, setsize+1)
	tempHashInput := make([]*big.Int, 2)
	tempHashInput[0] = min
	tempHashInput[1] = inputList[0]

	poseidonHashResult[0] = diHash(tempHashInput)

	for i := 1; i < setsize; i++ {
		tempHashInput[0] = inputList[i-1]
		tempHashInput[1] = inputList[i]
		poseidonHashResult[i] = diHash(tempHashInput)
	}

	tempHashInput[0] = inputList[setsize-1]
	tempHashInput[1] = max
	poseidonHashResult[setsize] = diHash(tempHashInput)

	return poseidonHashResult
}

// genUpdateListSet is simply used for testing.
func genUpdateListSet(inputList []*big.Int) []*big.Int {
	setsize := len(inputList)

	poseidonHashResult := make([]*big.Int, setsize-1)
	tempHashInput := make([]*big.Int, 2)

	for i := 0; i < setsize-1; i++ {
		tempHashInput[0] = inputList[i]
		tempHashInput[1] = inputList[i+1]
		poseidonHashResult[i] = diHash(tempHashInput)
		// if i == setsize-2 {
		// 	fmt.Println("Hash = ", poseidonHashResult[i])
		// }
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
	// AccTest3 := AccTest1.Mul(AccTest1, AccTest2)
	// AccTest3.Mod(AccTest3, setup.N)
	// fmt.Println("Accumulator_test3 = ", AccTest3.String()) //AccTest3 should be the same as the AccOld
}

func main() {
	fmt.Println("start test in main")
	setup := accumulator.TrustedSetup()
	oldSetSize := 10000
	delSetSize := 1000
	addSetSize := 2 * delSetSize

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
	fmt.Println("PoKE1 ")
	PoKE(AccMid, prodDelSet, AccOld, setup.N)

	// Generate add set
	// Note that add do not have to be continuous, we use this set just for example.
	// Delete set have to be covered by the add set
	addList := make([]*big.Int, addSetSize)
	// addList start from 1 to avoid 0 as input. 0 is already reserved for min value.
	for i := 0; i < addSetSize; i++ {
		addList[i] = big.NewInt(int64(1 + i))
	}
	// delList should be a subset of addList
	addListSet := genUpdateListSet(addList)
	prodAddSet := CalSetProd(addListSet)
	// var prodNewSet big.Int
	// prodNewSet.Mul(&prodMidSet, prodAddSet)
	AccNew := accumulator.AccumulateNew(AccMid, prodAddSet, setup.N)
	fmt.Println("Accumulator_New = ", AccNew.String())
	fmt.Println("PoKE2 ")
	PoKE(AccMid, prodAddSet, AccNew, setup.N)

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
