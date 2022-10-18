package accumulator

import (
	"math/big"

	"github.com/iden3/go-iden3-crypto/poseidon"
)

func genRepWithHashToPrimeFromSHA256(set []string) []*big.Int {
	ret := make([]*big.Int, len(set))
	for i, v := range set {
		ret[i] = HashToPrime([]byte(v))
	}
	return ret
}

func genRepWithDIHashFromPoseidon(set []string) []*big.Int {
	ret := make([]*big.Int, len(set))
	for i := range set {
		ret[i] = Min2048
		temp, err := poseidon.HashBytes([]byte(set[i]))
		if err != nil {
			panic(err)
		}
		ret[i].Add(ret[i], temp)
	}
	return ret
}

func genRepWithMultiDIHashFromPoseidon(set []string) []*big.Int {
	ret := make([]*big.Int, len(set))
	var err error
	for i := range set {
		ret[i], err = poseidon.HashBytes([]byte(set[i]))
		if err != nil {
			panic(err)
		}
	}
	return ret
}
