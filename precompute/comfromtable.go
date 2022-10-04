package precompute

import (
	"fmt"
	"math/big"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

var (
	big0 = big.NewInt(0)
	big1 = big.NewInt(1)
	big2 = big.NewInt(2)
)

// PreTable only allows to pre-compute a power of 2 for the base g.
// base[0] should be g and n[0] is 0.
// base[i] = g^{2^{n[i]}} mod N
type PreTable struct {
	base []*big.Int
	n    []int
}

func GenPreTable(base, N *big.Int, bitLen, tableSize int) *PreTable {
	var table PreTable
	table.base = make([]*big.Int, tableSize)
	table.n = make([]int, tableSize)

	stepSize := bitLen / tableSize

	table.base[0] = new(big.Int)
	table.base[0].Set(base)
	table.n[0] = 0

	for i := 1; i < tableSize; i++ {
		table.n[i] = table.n[i-1] + stepSize
		var power big.Int
		// need to optimize
		power.Exp(big2, big.NewInt(int64(table.n[i])), nil)
		table.base[i] = accumulator.AccumulateNew(base, &power, N)
	}

	return &table
}

func ComputeFromTable(table *PreTable, x, N *big.Int) *big.Int {
	// Todo: more checks for the validity of the table
	if len(table.base) != len(table.n) {
		panic("invalid pre-compute table, unbalanced")
	}
	if len(table.base) < 1 {
		panic("invalid pre-compute table, too small")
	}
	if x.Cmp(big0) < 1 {
		panic("invalid x, negative")
	}

	// Now, we divide x according to the n
	// We first find out how many sub part can we separate x according to the table
	length := x.BitLen()
	var iCounter int
	var xCopy big.Int
	xCopy.Set(x)
	fmt.Println("len(table.n) = ", len(table.n))
	fmt.Println("len(x) = ", length)
	for iCounter = 0; iCounter < len(table.n); iCounter++ {
		if table.n[iCounter] >= length {
			break
		}
	}
	if iCounter == len(table.n) {
		iCounter--
	}

	fmt.Println("iCounter = ", iCounter)
	subX := make([]big.Int, iCounter+2)

	for i := 1; i < iCounter+1; i++ {
		var modulo big.Int
		modulo.Exp(big2, big.NewInt(int64(table.n[i]-table.n[i-1])), nil)
		subX[i-1].Mod(&xCopy, &modulo)
		fmt.Println("subX[i-1] = ", subX[i-1].String())
		xCopy.Rsh(&xCopy, uint(table.n[i]-table.n[i-1]))
	}
	subX[iCounter].Set(&xCopy)
	fmt.Println("subX[iCounter+1] = ", subX[iCounter+1].String())

	// --------------start of test code------------------------------
	//the following code tests the if the separation for x is correct or not
	// var prod, temp big.Int
	// for i := 0; i < iCounter+1; i++ {
	// 	if i == 0 {
	// 		prod.Add(&prod, &subX[i])
	// 		continue
	// 	}
	// 	temp.Exp(big2, big.NewInt(int64(table.n[i])), nil)
	// 	temp.Mul(&temp, &subX[i])
	// 	prod.Add(&prod, &temp)
	// }
	// fmt.Println("prod = ", prod.String())
	// fmt.Println("x = ", x.String())
	// --------------end of test code------------------------------

	// The next part can be paralleled
	var prod big.Int
	prod.SetInt64(1)
	for i := 0; i < iCounter+1; i++ {
		if i == 0 {
			temp := accumulator.AccumulateNew(table.base[0], &subX[0], N)
			prod.Mul(&prod, temp)
			continue
		}
		temp := accumulator.AccumulateNew(table.base[i], &subX[i], N)
		prod.Mul(&prod, temp)
		prod.Mod(&prod, N)
	}
	return &prod
}
