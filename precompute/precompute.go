package precompute

import (
	"fmt"
	"math/big"
)

const tableSize = 16

var big1 = big.NewInt(1)

type Table struct {
	base        *big.Int
	n           *big.Int
	maxBitLen   int
	numElements uint64
	stepSize    uint
	table       []*big.Int
}

func NewTable(base, n, elementUpperBound *big.Int, numElements uint64) *Table {
	t := &Table{
		base:        base,
		n:           n,
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
	t.stepSize = uint(t.maxBitLen / tableSize)
	opt := new(big.Int)
	for i := uint(0); i < tableSize; i++ {
		opt.Lsh(big1, t.stepSize*i)
		t.table[i] = new(big.Int).Exp(base, opt, n)
		//fmt.Println("table[", i, "] = ", t.table[i])
	}
	return t
}

func (t *Table) Compute(x *big.Int, numRoutine int) *big.Int {
	xBitLen := x.BitLen()
	steps := xBitLen / int(t.stepSize)
	batchStep := steps / numRoutine
	if batchStep == 0 {
		batchStep = 1
	}
	batchSize := batchStep * int(t.stepSize)
	resChan := make(chan *big.Int)
	for i := 0; i < numRoutine; i++ {
		go routineCompute(t.table[batchStep*i], x, t.n, i*batchSize, batchSize, resChan)
	}
	res := big.NewInt(1)
	for i := 0; i < numRoutine; i++ {
		tmp := <-resChan
		fmt.Println("tmp = ", tmp)
		res.Mul(res, tmp)
		res.Mod(res, t.n)
	}
	return res
}

func routineCompute(base, x, n *big.Int, startBitLen, batchSize int, resChan chan *big.Int) {
	pow := big.NewInt(0)
	endBitLen := startBitLen + batchSize
	if endBitLen > x.BitLen() {
		endBitLen = x.BitLen()
	}
	for i := startBitLen; i < endBitLen; i++ {
		if x.Bit(i) == 1 {
			pow.Add(pow, big1)
		}
		pow.Lsh(pow, 1)
	}
	fmt.Println("pow = ", pow)
	fmt.Println(pow.BitLen())
	fmt.Println("base = ", base)
	res := new(big.Int).Exp(base, pow, n)
	fmt.Println("res = ", res)
	resChan <- res
}
