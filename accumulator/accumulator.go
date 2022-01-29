package accumulator

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"

	"github.com/rsa_accumulator/dihash"
)

func init() {
	_ = Max2048.Lsh(one, 2048)
	_ = Max2048.Sub(Max2048, one)
}

// TrustedSetup returns a pointer to AccumulatorSetup with 2048 bits key length
func TrustedSetup() *AccumulatorSetup {
	var ret AccumulatorSetup
	ret.N.SetString(N2048String, 10)
	ret.G.SetString(G2048String, 10)
	return &ret
}

func Accumulate(g, power, n *big.Int) *big.Int {
	var ret big.Int
	ret.Exp(g, power, n)
	return &ret
}

// This function can be done in a smarter way, but not the current version
func PreCompute() []big.Int {
	return preCompute(PreComputeSize)
}

func preCompute(preComputeSize int) []big.Int {
	trustedSetup := *TrustedSetup()

	ret := make([]big.Int, preComputeSize)

	ret[0] = trustedSetup.G
	ret[1] = *Accumulate(&trustedSetup.G, dihash.Delta, &trustedSetup.N)
	for i := 2; i < preComputeSize; i++ {
		ret[i] = *Accumulate(&ret[i-1], dihash.Delta, &trustedSetup.N)
	}
	return ret
}

func getSafePrime() *big.Int {
	ranNum, _ := crand.Prime(crand.Reader, securityPara/2)
	var temp big.Int
	flag := false
	for !flag {
		temp.Mul(ranNum, two)
		temp.Add(&temp, one)
		flag = temp.ProbablyPrime(securityParaInBits / 2)
		if !flag {
			ranNum, _ = crand.Prime(crand.Reader, securityPara)
		}
	}
	return &temp
}

func getRanQR(p, q *big.Int) *big.Int {
	rng := rand.New(rand.NewSource(123456))
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
