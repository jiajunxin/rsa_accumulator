package accumulator

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

// HashEncode generates different representatives that can be inputted into RSA accumulator
func HashEncode(set []string, encodeType EncodeType) []*big.Int {
	if encodeType == EncodeTypeSHA256HashToPrime {
		return sha256HashToPrime(set)
	}
	if encodeType == EncodeTypePoseidonDIHash {
		return poseidonDIHash(set)
	}
	return sha256HashToPrime(set)
}

func sha256HashToPrime(set []string) []*big.Int {
	ret := make([]*big.Int, len(set))
	for i, v := range set {
		ret[i] = HashToPrime([]byte(v))
	}
	return ret
}

func poseidonDIHash(set []string) []*big.Int {
	var (
		ret = make([]*big.Int, len(set))
		opt = new(big.Int)
		err error
	)
	for i := range set {
		ret[i] = Min2048()
		opt, err = poseidon.HashBytes([]byte(set[i]))
		if err != nil {
			panic(err)
		}
		ret[i].Add(ret[i], opt)
	}
	return ret
}
