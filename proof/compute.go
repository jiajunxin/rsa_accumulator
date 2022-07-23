package proof

import (
	"math/big"
)

var (
	big0 = big.NewInt(0)
	big1 = big.NewInt(1)
	big2 = big.NewInt(2)
	big4 = big.NewInt(4)
	big8 = big.NewInt(8)
)

func log2(n *big.Int) int {
	return n.BitLen() - 1
}
