package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/rsa_accumulator/accumulator"
	pf "github.com/rsa_accumulator/proof"
)

func main() {
	setup := accumulator.TrustedSetup()
	h, err := rand.Int(rand.Reader, setup.N)
	if err != nil {
		panic(err)
	}
	pp := pf.NewPublicParameters(setup.N, setup.G, h)
	r, err := rand.Int(rand.Reader, setup.N)
	if err != nil {
		panic(err)
	}
	x := new(big.Int)
	x.Exp(big.NewInt(2), big.NewInt(100), nil)
	x.Sub(x, big.NewInt(1))

	prover := pf.NewRPProver(r, x, pp)
	proof, err := prover.Prove()
	if err != nil {
		panic(err)
	}
	verifier := pf.NewRPVerifier(pp)
	isAccepted := verifier.Verify(proof)

	if isAccepted {
		fmt.Println("argument accepted")
	} else {
		fmt.Println("argument rejected")
	}
}
