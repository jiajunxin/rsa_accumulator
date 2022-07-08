package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/rsa_accumulator/accumulator"
	"github.com/rsa_accumulator/proof"
)

func main() {
	setup := accumulator.TrustedSetup()
	h, err := rand.Int(rand.Reader, setup.N)
	if err != nil {
		panic(err)
	}

	pp := proof.NewPublicParameters(setup.N, setup.G, h)
	r, err := rand.Int(rand.Reader, setup.N)
	if err != nil {
		panic(err)
	}
	prover := proof.NewRPProver(r, big.NewInt(12), pp)
	commitX, err := prover.CommitX()
	if err != nil {
		panic(err)
	}
	commitment, err := prover.ComposeCommitment()
	if err != nil {
		panic(err)
	}
	verifier := proof.NewRPVerifier(prover.C, pp)
	verifier.SetCommitment(commitment)
	e, err := verifier.Challenge()
	if err != nil {
		panic(err)
	}
	response, err := prover.Response(e)
	if err != nil {
		panic(err)
	}
	res := verifier.Verify(e, commitX, response)
	if res {
		fmt.Println("argument accepted")
	} else {
		fmt.Println("argument rejected")
	}
}
