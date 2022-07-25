package proof

import (
	"math/big"
)

const log2E = 1.44269504089

var (
	bigNeg1 = big.NewInt(-1)
	big0    = big.NewInt(0)
	big1    = big.NewInt(1)
	big2    = big.NewInt(2)
	big3    = big.NewInt(3)
	big4    = big.NewInt(4)
	big8    = big.NewInt(8)
)

func log2(n *big.Int) int {
	return n.BitLen() - 1
}

func log(n *big.Int) int {
	log2 := log2(n)
	return int(float64(log2) / log2E)
}
