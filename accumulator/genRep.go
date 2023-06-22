package accumulator

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon"
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
		ret[i] = new(big.Int)
		temp := poseidon.Poseidon(ElementFromString(set[i]))
		temp.ToBigIntRegular(ret[i])
		ret[i].Add(ret[i], Min1024)
	}
	return ret
}
