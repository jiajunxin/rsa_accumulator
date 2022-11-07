package accumulator

import (
	"math/big"
)

// SimpleExp should calculate g^x mod n.
// It is implemented here to campare with golang's official Exp and MultiExp
func SimpleExp(g, x, n *big.Int) *big.Int {
	if g.Cmp(big1) <= 0 || n.Cmp(big1) <= 0 || x.Cmp(big1) < 0 {
		panic("invalid input for function SimpleExp")
	}
	// change x to its binary representation
	binaryX := x.Bits()

}
