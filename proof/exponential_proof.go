package proof

import "math/big"

var (
	epSecurityPara = big.NewInt(128)
)

// ExpProver is the prover for proof of exponentiation
type ExpProver struct {
}

// ExpVerifier is the verifier for proof of exponentiation
type ExpVerifier struct {
}
