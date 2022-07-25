package main

import (
	"crypto/rand"
	"math/big"

	"github.com/rsa_accumulator/accumulator"
	"github.com/rsa_accumulator/proof"
)

func main() {
	n := new(big.Int)
	n.SetString(accumulator.N2048String, 10)
	g := new(big.Int)
	g.SetString(accumulator.G2048String, 10)
	h := new(big.Int)
	h.SetString(accumulator.H2048String, 10)
	r, err := rand.Int(rand.Reader, n)
	handleError(err)
	pp := proof.NewPublicParameters(n, g, h)
	x := big.NewInt(100)
	prover := proof.NewRPProver(pp, r, x)
	pf, err := prover.Prove()
	handleError(err)
	verifier := proof.NewRPVerifier(pp)
	ok := verifier.Verify(pf)
	if !ok {
		panic("verification failed")
	} else {
		println("verification succeeded")
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
