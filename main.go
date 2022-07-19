package main

import (
	"crypto/rand"
	"fmt"
	"github.com/rsa_accumulator/accumulator"
	pf "github.com/rsa_accumulator/proof"
	"math/big"
)

func main() {
	setup := accumulator.TrustedSetup()
	h, err := rand.Int(rand.Reader, setup.N)
	if err != nil {
		panic(err)
	}
	pp := pf.NewPublicParameters(setup.N, setup.G, h)
	prover := pf.NewExpProver(pp)
	verifier := pf.NewExpVerifier(pp)
	u := big.NewInt(213123)
	x := big.NewInt(123)
	w := new(big.Int).Exp(u, x, nil)
	commitment, err := prover.Commit(u, w, x)
	if err != nil {
		panic(err)
	}
	verifier.SetCommitment(commitment)
	challenge, err := verifier.Challenge()
	if err != nil {
		panic(err)
	}
	response, err := prover.Response(challenge)
	if err != nil {
		panic(err)
	}
	ok, err := verifier.VerifyResponse(u, w, response)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic("verification failed")
	}
	fmt.Println("Verification succeeded")
	//r, err := rand.Int(rand.Reader, setup.N)
	//if err != nil {
	//	panic(err)
	//}
	//x := new(big.Int)
	//x.Exp(big.NewInt(2), big.NewInt(100), nil)
	//x.Sub(x, big.NewInt(1))
	//
	//prover := pf.NewRPProver(pp, r, x)
	//proof, err := prover.Commit()
	//if err != nil {
	//	panic(err)
	//}
	//verifier := pf.NewRPVerifier(pp)
	//isAccepted := verifier.Verify(proof)
	//
	//if isAccepted {
	//	fmt.Println("argument accepted")
	//} else {
	//	fmt.Println("argument rejected")
	//}
}
