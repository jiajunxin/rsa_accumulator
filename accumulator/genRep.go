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
	//fmt.Println("set size = ", len(set))
	var err error
	for i := range set {
		ret[i], err = poseidon.HashBytes([]byte(set[i]))
		if err != nil {
			panic(err)
		}
		ret[i].Add(ret[i], Min1024)
		if ret[i].Bit(0) == 0 {
			ret[i].Add(ret[i], big1)
		}
		//fmt.Println("Inside DI hash, ret = ", ret[i].String())
	}
	// fmt.Println("Inside DI hash, loop ret ")
	// for i := range ret {
	// 		fmt.Println("i = ", i, "ret[i] = ", ret[i].String())
	// }
	return ret
}

// This function is a specific Multi-DI Hash with 80 bits of security.
// The first 255 bits are the output of Poseidon Hash,
// the second 255 bits are inputing the first part into a Universal Hash
// the third 255 bits are inputing the second part into the same Universal Hash
func genRepWithMultiDIHashFromPoseidon(set []string) []*big.Int {
	ret := make([]*big.Int, 3*len(set))
	setSize := len(set)
	for i := range set {
		ret[i] = new(big.Int)
		temp, err := poseidon.HashBytes([]byte(set[i]))
		if err != nil {
			panic(err)
		}
		ret[i].Add(ret[i], temp)
		if ret[i].Bit(0) == 0 {
			ret[i].Add(ret[i], big1)
		}
	}
	for i := 0; i < setSize; i++ {
		ret[i+setSize] = new(big.Int)
		ret[i+setSize*2] = new(big.Int)
		ret[i+setSize] = UniversalHashToInt(ret[i])
		ret[i+setSize*2] = UniversalHashToInt(ret[i+setSize])
	}
	return ret
}
