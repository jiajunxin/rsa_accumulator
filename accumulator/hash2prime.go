package accumulator

import (
	"crypto/sha256"
	"math/big"
)

func HashToPrime(input []byte) *big.Int {
	var ret big.Int
	h := sha256.New()
	h.Write([]byte(input))
	hashTemp := h.Sum(nil)
	ret.SetBytes(hashTemp)
	flag := false
	for !flag {
		flag = ret.ProbablyPrime(securityParaInBits / 2)
		if !flag {
			h.Reset()
			h.Write(hashTemp)
			hashTemp = h.Sum(nil)
			ret.SetBytes(hashTemp)
		}
	}
	return &ret
}
