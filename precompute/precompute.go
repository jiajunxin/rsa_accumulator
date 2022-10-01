package precompute

import "math/big"

const tableSize = 1024

var big1 = big.NewInt(1)

type Table struct {
	base        *big.Int
	maxBitLen   int
	numElements uint64
	table       []*big.Int
}

func NewTable(base, n, elementUpperBound *big.Int, numElements uint64) *Table {
	t := &Table{
		base:        base,
		numElements: numElements,
	}
	t.table = make([]*big.Int, tableSize)
	// calculate the maximum possible value of the product of the elements
	prodMax := big.NewInt(1)
	// TODO: optimize the multiplication
	for i := uint64(0); i < numElements; i++ {
		prodMax.Mul(prodMax, elementUpperBound)
	}

	t.maxBitLen = prodMax.BitLen()
	step := uint(t.maxBitLen / tableSize)
	opt := new(big.Int)
	for i := uint(0); i < tableSize; i++ {
		opt.Lsh(big1, step*i)
		t.table[i] = new(big.Int).Exp(base, opt, n)
	}
	return t
}
