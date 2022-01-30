package accumulator

import (
	"math/big"
)

func init() {
	_ = Min2048.Lsh(one, 2047)
}

// TrustedSetup returns a pointer to AccumulatorSetup with 2048 bits key length
func TrustedSetup() *AccumulatorSetup {
	var ret AccumulatorSetup
	ret.N.SetString(N2048String, 10)
	ret.G.SetString(G2048String, 10)
	return &ret
}

func GenRepersentatives(set []string, encodeType EncodeType) []big.Int {
	switch encodeType {
	case HashToPrimeFromSha256:
		return genRepWithHashToPrimeFromSHA256(set)
	case DIHashFromPoseidon:
		return genRepWithDIHashFromPoseidon(set)
	default:
		return genRepWithHashToPrimeFromSHA256(set)
	}
}

// AccAndProve generates the accumulator with all the memberships precomputed
func AccAndProve(set []string, encodeType EncodeType, setup *AccumulatorSetup) (*big.Int, []big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembership(&setup.G, &setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := Accumulate(&proofs[0], &rep[0], &setup.N)

	return acc, proofs
}

// ProveMembership uses divide-and-conqure method to pre-compute the all membership proofs in time O(nlogn)
func ProveMembership(base, N *big.Int, set []big.Int) []big.Int {
	if len(set) == 1 {
		ret := make([]big.Int, 1)
		ret[0] = *base
		return ret
	}
	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftBase := *accumulate(set[len(set)/2:], base, N)
	rightBase := *accumulate(set[0:len(set)/2], base, N)
	proofs := ProveMembership(&leftBase, N, set[0:len(set)/2])
	proofs = append(proofs, ProveMembership(&rightBase, N, set[len(set)/2:])...)
	return proofs
}

func accumulate(set []big.Int, g, N *big.Int) *big.Int {
	var acc big.Int
	acc.Set(g)
	for _, v := range set {
		acc.Exp(&acc, &v, N)
	}
	return &acc
}

func Accumulate(g, power, N *big.Int) *big.Int {
	var ret big.Int
	ret.Exp(g, power, N)
	return &ret
}
