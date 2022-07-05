package proof

import "math/big"

func isPerfectSquare(n *big.Int) (*big.Int, bool) {
	sqrt := new(big.Int).Sqrt(n)
	return sqrt, new(big.Int).Mul(sqrt, sqrt).Cmp(n) == 0
}

func euclideanStep(a, b *big.Int) (*big.Int, error) {
	if a.Cmp(b) == -1 {
		a, b = b, a
	}
	q := new(big.Int).Mod(b, a)
	return q, nil
}
