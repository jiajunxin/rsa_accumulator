package proof

import (
	"math/big"
	"sync"
)

var (
	big0 = big.NewInt(0)
	big1 = big.NewInt(1)
	big2 = big.NewInt(2)
	big3 = big.NewInt(3)
	big5 = big.NewInt(5)
	big7 = big.NewInt(7)
	big8 = big.NewInt(8)
	// sync pool for big integers, lease GC and improve performance
	iPool = sync.Pool{
		New: func() interface{} { return new(big.Int) },
	}
)

func log2(n *big.Int) int {
	return n.BitLen() - 1
}
