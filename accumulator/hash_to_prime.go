package accumulator

import (
	"crypto/sha256"
	"math/big"
)

// HashToPrime takes the input into Sha256 and take the hash output to input repeatedly until we hit a prime number
func HashToPrime(input []byte) *big.Int {
	ret := new(big.Int)
	h := sha256.New()
	h.Write(input)
	hashTemp := h.Sum(nil)
	ret.SetBytes(hashTemp)
	for {
		isPrime := ret.ProbablyPrime(securityParaHashToPrime)
		if isPrime {
			break
		}
		h.Reset()
		h.Write(hashTemp)
		hashTemp = h.Sum(nil)
		ret.SetBytes(hashTemp)
	}
	return ret
}

// SHA256ToInt calculates the input with Sha256 and change it to big.Int
func SHA256ToInt(input []byte) *big.Int {
	h := sha256.New()
	h.Write(input)
	hashTemp := h.Sum(nil)
	ret := new(big.Int).SetBytes(hashTemp)
	return ret
}
