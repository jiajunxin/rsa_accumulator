package complex

import "math/big"

// roundFloat rounds the given big float to the nearest big integer
func roundFloat(f *big.Float) *big.Int {
	delta := big.NewFloat(roundingDelta)
	if f.Sign() < 0 {
		delta.Neg(delta)
	}
	f.Add(f, delta)
	res := new(big.Int)
	f.Int(res)
	return res
}
