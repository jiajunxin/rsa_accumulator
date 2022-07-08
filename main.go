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
	r, _ := new(big.Int).SetString("24235325234562345123453242432454789165454564654564698", 10)
	prover := proof.NewRPProver(r, big.NewInt(12), pp)
	commitX, err := prover.CommitX()
	if err != nil {
		panic(err)
	}
	commitment, err := prover.ComposeCommitment()
	if err != nil {
		panic(err)
	}
	fmt.Println("commitment done")
	verifier := proof.NewRPVerifier(prover.C, pp)
	verifier.Commitment = commitment
	e, err := verifier.Challenge()
	if err != nil {
		panic(err)
	}
	fmt.Printf("e: %s\n", e)
	response, err := prover.Response(e)
	if err != nil {
		panic(err)
	}
	res := verifier.Verify(e, commitX, response)
	fmt.Printf("res: %t\n", res)
}
