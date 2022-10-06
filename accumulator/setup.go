package accumulator

import "math/big"

// TrustedSetup returns a pointer to AccumulatorSetup with 2048 bits key length
func TrustedSetup() *Setup {
	n := new(big.Int)
	if _, ok := n.SetString(N2048String, 10); !ok {
		panic("Error in TrustedSetup, N2048String is not a valid number")
	}
	g := new(big.Int)
	if _, ok := g.SetString(G2048String, 10); !ok {
		panic("Error in TrustedSetup, G2048String is not a valid number")
	}
	return &Setup{
		N: n,
		G: g,
	}
}
