package accumulator

import (
	"crypto/sha256"
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

// HashToPrime takes the input into Sha256 and take the hash output to input repeatedly until we hit a prime number
func HashToPrime(input []byte) *big.Int {
	var ret big.Int
	h := sha256.New()
	h.Write(input)
	hashTemp := h.Sum(nil)
	ret.SetBytes(hashTemp)
	flag := false
	for !flag {
		flag = ret.ProbablyPrime(securityParaHashToPrime)
		if !flag {
			h.Reset()
			h.Write(hashTemp)
			hashTemp = h.Sum(nil)
			ret.SetBytes(hashTemp)
		}
	}
	return &ret
}

// SHA256ToInt calculates the input with Sha256 and change it to big.Int
func SHA256ToInt(input []byte) *big.Int {
	var ret big.Int
	h := sha256.New()
	h.Write(input)
	hashTemp := h.Sum(nil)
	ret.SetBytes(hashTemp)
	return &ret
}

func PoseidonWith2Inputs(inputs []*big.Int) *big.Int {
	if len(inputs) != 2 {
		panic("PoseidonWith2Inputs requires 2 inputs")
	}
	ret, err := poseidon.Hash(inputs)
	if err != nil {
		panic(err)
	}
	return ret
}

// UniversalHashToInt calculates output = input * A + B mod P
func UniversalHashToInt(input *big.Int) *big.Int {
	var ret big.Int
	ret.Mul(input, A)
	ret.Mod(&ret, P)
	ret.Add(&ret, B)
	ret.Mod(&ret, P)
	if ret.Bit(0) == 0 {
		ret.Add(&ret, big1)
	}
	return &ret
}
