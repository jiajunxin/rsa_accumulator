package precompute

import (
	"math/big"
)

const tableSize = 512

var big1 = big.NewInt(1)

// Table is the precomputing table
type Table struct {
	base        *big.Int
	n           *big.Int
	maxBitLen   int
	numElements uint64
	stepSize    uint
	table       []*big.Int
}

// NewTable creates a new precomputing table
func NewTable(base, n, elementUpperBound *big.Int, numElements uint64) *Table {
	t := &Table{
		base:        base,
		n:           n,
		numElements: numElements,
	}
	t.table = make([]*big.Int, tableSize)
	t.maxBitLen = elementUpperBound.BitLen() * int(numElements)
	t.stepSize = uint(t.maxBitLen / tableSize)
	opt := new(big.Int).Lsh(big1, t.stepSize)
	t.table[0] = new(big.Int).Set(base)
	//fmt.Println("table[ 0 ] = ", t.table[0])
	for i := uint(1); i < tableSize; i++ {
		t.table[i] = new(big.Int).Exp(t.table[i-1], opt, n)
		//fmt.Println("table[", i, "] = ", t.table[i])
	}
	return t
}

// Compute computes the result of base^x mod n with specified number of goroutines
func (t *Table) Compute(x *big.Int, numRoutine int) *big.Int {
	xBitLen := x.BitLen()
	steps := xBitLen / int(t.stepSize)
	batchStep := steps / numRoutine
	if batchStep == 0 {
		batchStep = 1
	}
	batchSize := batchStep * int(t.stepSize)
	resChan := make(chan *big.Int, numRoutine)
	for i := 0; i < numRoutine; i++ {
		startBitLen := i * batchSize
		var endBitLen int
		if i == numRoutine-1 {
			endBitLen = xBitLen
		} else {
			endBitLen = startBitLen + batchSize
		}
		go routineCompute(t.table[batchStep*i], x, t.n, startBitLen, endBitLen, resChan)
	}
	res := big.NewInt(1)
	for i := 0; i < numRoutine; i++ {
		tmp := <-resChan
		res.Mul(res, tmp)
		res.Mod(res, t.n)
	}
	return res
}

func routineCompute(base, x, n *big.Int, startBitLen, endBitLen int, resChan chan *big.Int) {
	pow := big.NewInt(0)
	for i := endBitLen - 1; i >= startBitLen; i-- {
		pow.Lsh(pow, 1)
		if x.Bit(i) == 1 {
			pow.Add(pow, big1)
		}
	}
	resChan <- new(big.Int).Exp(base, pow, n)
}
