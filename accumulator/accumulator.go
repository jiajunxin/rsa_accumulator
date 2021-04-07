package accumulator

import (
	crand "crypto/rand"
	"fmt"
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
	ret.P.SetString(P2048String, 10)
	ret.Q.SetString(Q2048String, 10)
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
	newBase := *Accumulate(&trustedSetup.G, dihash.Delta, &trustedSetup.N)

	ret := make([]big.Int, preComputeSize)

	ret[0] = newBase
	var temp big.Int
	for i := 1; i < preComputeSize; i++ {
		temp.SetUint64(uint64(i + 1))
		ret[i] = *Accumulate(&newBase, &temp, &trustedSetup.N)
		//fmt.Println("i = ", i)
	}
	return ret
}

// AccumulateSetWirhPreCompute accumulates the input elements using DIHash with pre-compute
func AccumulateSetWirhPreCompute(inputSet []Element, bases []big.Int) *big.Int {
	trustedSetup := *TrustedSetup()
	var ret big.Int
	setSize := len(inputSet)
	fmt.Println("set size = ", setSize)
	setWindowValue := SetWindowValue(inputSet)
	fmt.Println("setWindowValue size = ", len(setWindowValue))
	fmt.Println("set window value 0 = ", setWindowValue[0].String())
	//fmt.Println("set window value 1 = ", setWindowValue[1].String())
	//fmt.Println("set window value 2 = ", setWindowValue[2].String())
	fmt.Println("test in AccumulateSetWirhPreCompute 0")
	var temp big.Int
	ret.Set(one)
	for i := 0; i < setSize; i++ {
		temp = *Accumulate(&bases[setSize-i-1], &setWindowValue[i], &trustedSetup.N)
		ret.Mul(&ret, &temp)
		ret.Mod(&ret, &trustedSetup.N)
	}
	fmt.Println("test in AccumulateSetWirhPreCompute 1")
	temp = *Accumulate(&trustedSetup.G, &setWindowValue[setSize], &trustedSetup.N)
	ret.Mul(&ret, &temp)
	ret.Mod(&ret, &trustedSetup.N)

	return &ret
}

// AccumulateSetWirhPreCompute accumulates the input elements using DIHash with pre-compute
func AccumulateSetWithoutPreCompute(inputSet []Element) *big.Int {
	trustedSetup := *TrustedSetup()
	var ret big.Int
	setSize := len(inputSet)
	hashValues := make([]big.Int, setSize)

	{
		//sha256 values
		a := SHA256ToInt(inputSet[0])
		b := SHA256ToInt(inputSet[1])
		var temp big.Int
		temp.Add(a, b)
		fmt.Println("a + b = ", temp.String())
		temp.Mul(a, b)
		fmt.Println("a * b = ", temp.String())
	}

	for i := 0; i < setSize; i++ {
		hashValues[i] = *dihash.DIHash(inputSet[i])
	}
	fmt.Println("di hashValues 0 = ", hashValues[0])
	//fmt.Println("di hashValues 1 = ", hashValues[1])

	var tempGenerator big.Int
	tempGenerator.Set(&trustedSetup.G)

	for i := 0; i < setSize; i++ {
		tempGenerator = *Accumulate(&tempGenerator, &hashValues[i], &trustedSetup.N)
	}
	ret.Set(&tempGenerator)

	return &ret
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
