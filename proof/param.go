package proof

import "math/big"

const (
	// security parameter for range proof and proof of exponentiation
	securityParam = 128
)

// PublicParameters holds public parameters initialized during the setup procedure
type PublicParameters struct {
	N *big.Int
	G *big.Int
	H *big.Int
}

// NewPublicParameters generates a new public parameter configuration
func NewPublicParameters(n, g, h *big.Int) *PublicParameters {
	return &PublicParameters{
		N: n,
		G: g,
		H: h,
	}
}
