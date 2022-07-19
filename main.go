package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/rsa_accumulator/proof"
)

func main() {
	//accumulator.TrustedSetupForQRN()
	//setup := accumulator.TrustedSetup()
	//h, err := rand.Int(rand.Reader, setup.N)
	//if err != nil {
	//	panic(err)
	//}
	//pp := pf.NewPublicParameters(setup.N, setup.G, h)
	//prover := pf.NewExpProver(pp)
	//verifier := pf.NewExpVerifier(pp)
	//u := big.NewInt(213123)
	//target := big.NewInt(123)
	//w := new(big.Int).Exp(u, target, nil)
	//commitment, err := prover.Commit(u, w, target)
	//if err != nil {
	//	panic(err)
	//}
	//verifier.SetCommitment(commitment)
	//challenge, err := verifier.Challenge()
	//if err != nil {
	//	panic(err)
	//}
	//response, err := prover.Response(challenge)
	//if err != nil {
	//	panic(err)
	//}
	//ok, err := verifier.VerifyResponse(u, w, response)
	//if err != nil {
	//	panic(err)
	//}
	//if !ok {
	//	panic("verification failed")
	//}
	//fmt.Println("Verification succeeded")

	//r, err := rand.Int(rand.Reader, setup.N)
	//if err != nil {
	//	panic(err)
	//}
	//target := new(big.Int)
	//target.Exp(big.NewInt(2), big.NewInt(100), nil)
	//target.Sub(target, big.NewInt(1))
	//
	//prover := pf.NewRPProver(pp, r, target)
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
	const logDir = "test.log"
	f, err := os.OpenFile(logDir, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	for i := 0; i < 10; i++ {
		_, err = f.WriteString(time.Now().String() + "\n")
		if err != nil {
			panic(err)
		}
		randLmt := new(big.Int).Exp(big.NewInt(2), big.NewInt(1792), nil)
		target, err := rand.Int(rand.Reader, randLmt)
		if err != nil {
			panic(err)
		}
		fmt.Println(target)
		_, err = f.WriteString(target.String() + "\n")
		if err != nil {
			panic(err)
		}
		start := time.Now()
		fs, err := proof.LagrangeFourSquares(target)
		if err != nil {
			panic(err)
		}
		currTime := time.Now()
		fmt.Println(currTime.Sub(start))
		_, err = f.WriteString(time.Since(start).String() + "\n")
		if err != nil {
			panic(err)
		}
		ok := proof.Verify(target, fs)
		if ok {
			fmt.Println("verification succeeded")
			_, err := f.WriteString("verification succeeded\n")
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("verification failed")
			_, err := f.WriteString("verification failed\n")
			if err != nil {
				panic(err)
			}
		}
	}
}
