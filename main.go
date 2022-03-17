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

	// AccTest := accumulator.AccumulateNew(AccMid, prodDelSet, setup.N)
	// fmt.Println("AccTest = ", AccTest.String()) //AccTest should be the same as the AccOld

	l1 := accumulator.HashToPrime(append(AccOld.Bytes(), AccMid.Bytes()...))
	fmt.Println("primeChallenge = ", l1.String())

	//～～～～～～～～～～～～～～～～～～～～～～～～～～～～～～～～
	var rTemp big.Int
	rTemp.SetString("1", 10)
	var rTemp2 big.Int
	for i := 0; i < len(delListSet); i++ {
		rTemp2.Mod(delListSet[i], l1)
		fmt.Println("rTemp2 = ", rTemp2.String())
		rTemp.Mul(&rTemp, &rTemp2)
		fmt.Println("rTemp * rTemp2 = ", rTemp.String())
		rTemp.Mod(&rTemp, l1)
		fmt.Println("rTemp = ", rTemp.String())
	}

	//fmt.Println("rTemp = ", rTemp.String())
	//～～～～～～～～～～～～～～～～～～～～～～～～～～～～～～～～

	var temp big.Int
	temp.SetString("16158503035655503650357438344334975980222051334857742016065172713762327569433945446598600705761456731844358980460949009747059779575245460547544076193224141560315438683650498045875098875194826053398028819192033784138396109321309878080919047169238085235290822926018152521443787945770532904303776199561965192760957166694834171210342487393282284747428088017663161029038902829665513096354230157075129296432088558362971801859230928678799175576150822952201848806616643615613562842355410104862578550863465661734839271290328348967522998634176499319122440250970673193049291436530055113939486576928778008576099804946423629969146", 10)
	temp.Mod(&temp, l1)
	fmt.Println("temp = ", temp.String())
	r1 := big.NewInt(1)
	q1 := big.NewInt(1)
	q1, r1 = q1.DivMod(prodDelSet, l1, r1)
	Q1 := accumulator.AccumulateNew(AccMid, q1, setup.N)
	fmt.Println("Q1 = ", Q1.String())
	fmt.Println("r1 = ", r1.String())

	// AccTest1 := accumulator.AccumulateNew(Q1, l1, setup.N)
	// AccTest2 := accumulator.AccumulateNew(AccMid, r1, setup.N)
	// AccTest3 := AccTest1.Mul(AccTest1, AccTest2)
	// AccTest3.Mod(AccTest3, setup.N)
	// fmt.Println("Accumulator_test = ", AccTest3.String()) //AccTest3 should be the same as the AccOld

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
