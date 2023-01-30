package accumulator

import (
	"crypto/sha256"
	"math/big"
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

// output = input * A + B mod P
func UniversalHashToInt(input *big.Int) *big.Int {
	var ret big.Int
	ret.Mul(input, A)
	ret.Mod(&ret, P)
	ret.Add(&ret, B)
	ret.Mod(&ret, P)
	return &ret
}
