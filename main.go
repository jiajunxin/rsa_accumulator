package main

import (
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
	for i := 0; i < 100; i++ {
		pp := proof.NewPublicParameters(n, g, h)
		u := big.NewInt(123)
		x := big.NewInt(3)
		w := new(big.Int)
		w.Exp(u, x, nil)
		prover := proof.NewExpProver(pp)
		verifier := proof.NewExpVerifier(pp)
		pf, err := prover.Prove(u, x)
		handleError(err)
		ok, err := verifier.Verify(pf, u, w)
		handleError(err)
		if !ok {
			panic("verification failed")
		}
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
