package accumulator

import (
	"math/big"
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
	//Todo: generate DI from Poseidon Hash
	for i := range set {
		ret[i] = Min2048
		//Todo ret[i] = ret[i] + PoseidonHash(v)
	}
	return ret
}
