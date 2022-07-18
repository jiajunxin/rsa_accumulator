package proof

import "math/big"

const (
	lenG = 10
)

var (
	// proof of exponentiation security parameter, lambda
	epSecurityParam = big.NewInt(128)
	// bound B
	epB *big.Int
)

func init() {
	// calculate bound B for proof of exponentiation, B > (2^(2*lambda))*|G|
	epB = new(big.Int)
	epB.Exp(big.NewInt(2), big.NewInt(2*epSecurityParam.Int64()), nil)
	epB.Mul(epB, big.NewInt(int64(lenG)))
}

// ExpProver is the prover for proof of exponentiation
type ExpProver struct {
}

// ExpVerifier is the verifier for proof of exponentiation
type ExpVerifier struct {
}
