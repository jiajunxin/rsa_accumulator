package complex

import "math/big"

const (
	roundingDelta = 0.4999999
)

var (
	// big integer
	big1    = big.NewInt(1)
	bigNeg1 = big.NewInt(-1)
	big2    = big.NewInt(2)
	big4    = big.NewInt(4)

	// big float
	big2f = big.NewFloat(2)
)
