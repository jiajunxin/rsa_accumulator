package accumulator

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
)

func init() {
	_ = Max2048.Lsh(one, 2048)
	_ = Max2048.Sub(Max2048, one)
}

// TrustedSetup returns a pointer to AccumulatorSetup with 2048 bits key length
func TrustedSetup() *AccumulatorSetup {
	var ret AccumulatorSetup
	ret.P.SetString(P2048String, 10)
	ret.Q.SetString(Q2048String, 10)
	ret.N.SetString(N2048String, 10)
	ret.G.SetString(G2048String, 10)
	return &ret
}

func getSafePrime() *big.Int {
	ranNum, _ := crand.Prime(crand.Reader, securityPara/2)
	fmt.Println("test 0.6")
	var temp big.Int
	flag := false
	for !flag {
		temp.Mul(ranNum, two)
		temp.Add(&temp, one)
		//fmt.Println("test 0.7")
		flag = temp.ProbablyPrime(securityParaInBits)
		//fmt.Println("test 0.8")
		if !flag {
			ranNum, _ = crand.Prime(crand.Reader, securityPara)
		}
	}
	return &temp
}

func getRanQR(p, q *big.Int) *big.Int {
	rng := rand.New(rand.NewSource(crsNum))
	var N big.Int
	N.Mul(p, q)
	var ranNum big.Int
	ranNum.Rand(rng, &N)

	flag := false
	for !flag {
		flag = isQR(&ranNum, p, q)
		if !flag {
			ranNum.Rand(rng, &N)
		}
	}
	return &ranNum
}

func isQR(input, p, q *big.Int) bool {
	if big.Jacobi(input, p) == 1 && big.Jacobi(input, q) == 1 {
		return true
	}
	return false
}
