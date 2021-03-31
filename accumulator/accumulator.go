package accumulator

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
)

const securityPara = 2048
const securityParaInBits = 128
const crs = "HKUST2021" //used as the seed for generating random numbers
const crsNum = 100      //used as the seed for generating random numbers

var one = big.NewInt(1)
var two = big.NewInt(2)
var Max2048 = big.NewInt(0)

func init() {
	_ = Max2048.Lsh(one, 2048)
	_ = Max2048.Sub(Max2048, one)
}

type AccumulatorSetup struct {
	p big.Int
	q big.Int
	N big.Int
	g big.Int //generator in QR_N
}

func Init() *AccumulatorSetup {
	var p, q, N, g big.Int

	crand.Read([]byte(crs))
	p = *getSafePrime()
	q = *getSafePrime()
	N.Mul(&p, &q)
	g = *getRanQR(&p, &q)

	var ret AccumulatorSetup
	ret.p = p
	ret.q = q
	ret.N = N
	ret.g = g
	return &ret
}

func getSafePrime() *big.Int {
	ranNum, _ := crand.Prime(crand.Reader, securityPara)

	var temp big.Int
	flag := false
	for !flag {
		temp.Sub(ranNum, one)
		temp.Div(&temp, two)
		flag = temp.ProbablyPrime(securityParaInBits)
		if !flag {
			ranNum, _ = crand.Prime(crand.Reader, securityPara)
		}
	}
	return ranNum
}

func getRanQR(p, q *big.Int) *big.Int {
	rng := rand.New(rand.NewSource(crsNum))
	var ranNum *big.Int
	ranNum.Rand(rng, Max2048)

	flag := false
	for !flag {
		flag = isQR(ranNum, p, q)
	}
	return ranNum
}

func isQR(input, p, q *big.Int) bool {
	if big.Jacobi(input, p) == 1 && big.Jacobi(input, q) == 1 {
		return true
	}
	return false
}
