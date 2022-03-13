package main

import (
	"fmt"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/rsa_accumulator/accumulator"
)

func main() {
	fmt.Println("start test in main")
	min := big.NewInt(0)
	// we set max to 0x7FFFFFFF instead of 0xFFFFFFFF to suit the xjSNARK, it seems there is some problem in there
	max := big.NewInt(0x7FFFFFFF)
	setsize := 16
	var err error

	sortedList1 := make([]*big.Int, setsize)
	for i := 0; i < setsize; i++ {
		sortedList1[i] = big.NewInt(int64(i))
	}

	// generate product of the input sorted list
	poseidonHashResult := make([]*big.Int, setsize+1)
	tempHashInput := make([]*big.Int, 2)
	tempHashInput[0] = min
	tempHashInput[1] = sortedList1[0]

	poseidonHashResult[0], err = poseidon.Hash(tempHashInput)

	if err != nil {
		// not expecting error from a Hash function
		panic(err)
	}
	for i := 1; i < setsize; i++ {
		tempHashInput[0] = sortedList1[i-1]
		tempHashInput[1] = sortedList1[i]
		poseidonHashResult[i], err = poseidon.Hash(tempHashInput)
		if err != nil {
			// not expecting error from a Hash function
			panic(err)
		}
	}
	tempHashInput[0] = sortedList1[setsize-1]
	tempHashInput[1] = max
	poseidonHashResult[setsize], err = poseidon.Hash(tempHashInput)

	if err != nil {
		// not expecting error from a Hash function
		panic(err)
	}

	var l1 big.Int
	l1.SetString("75117285383387635827127513071317", 10)
	rem1 := big.NewInt(1)
	tempDIHash := big.NewInt(1)
	tempMod := big.NewInt(1)
	tempPrint := big.NewInt(1)
	for i := 0; i < setsize+1; i++ {
		tempDIHash.Add(accumulator.Min2048, poseidonHashResult[i])
		tempPrint.Set(tempDIHash)
		tempMod = tempMod.Mod(tempDIHash, &l1)
		rem1.Mul(rem1, tempMod)
		rem1.Mod(rem1, &l1)
	}
	fmt.Println("rem1 = ", rem1)

	// build the accumulator and calculate the prime challenge
	setup := accumulator.TrustedSetup()
	prod := big.NewInt(1)
	for i := 0; i < setsize+1; i++ {
		prod.Mul(prod, poseidonHashResult[i])
	}
	AccOld := accumulator.AccumulateNew(setup.G, prod, setup.N)
	primeChallenge := accumulator.HashToPrime(AccOld.Bytes())
	fmt.Println("Accumulator_old = ", AccOld.String())
	fmt.Println("primeChallenge = ", primeChallenge)

}
