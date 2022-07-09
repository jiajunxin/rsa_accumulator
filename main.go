package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/rsa_accumulator/proof"
)

func main() {
	//setup := accumulator.TrustedSetup()
	//h, err := rand.Int(rand.Reader, setup.N)
	//if err != nil {
	//	panic(err)
	//}
	//
	//pp := proof.NewPublicParameters(setup.N, setup.G, h)
	//r, err := rand.Int(rand.Reader, setup.N)
	//if err != nil {
	//	panic(err)
	//}
	//prover := proof.NewRPProver(r, big.NewInt(12), pp)
	//commitX, err := prover.CommitX()
	//if err != nil {
	//	panic(err)
	//}
	//commitment, err := prover.ComposeCommitment()
	//if err != nil {
	//	panic(err)
	//}
	//verifier := proof.NewRPVerifier(pp)
	//verifier.SetC(prover.C)
	//verifier.SetCommitment(commitment)
	//verifier.SetCommitX(commitX)
	//response, err := prover.Response()
	//if err != nil {
	//	panic(err)
	//}
	//res := verifier.Verify(response)
	//if res {
	//	fmt.Println("argument accepted")
	//} else {
	//	fmt.Println("argument rejected")
	//}
	target := new(big.Int)
	target.Exp(big.NewInt(2), big.NewInt(600), nil)
	target.Sub(target, big.NewInt(1))
	fmt.Println("bit len: ", target.BitLen())
	start := time.Now()
	fs, err := proof.LagrangeFourSquares(target)
	if err != nil {
		panic(err)
	}
	fmt.Println(time.Since(start))
	verify := proof.Verify(target, fs)
	fmt.Println(fs)
	if verify {
		fmt.Println("verify success")
	} else {
		fmt.Println("verify failed")
	}
}
