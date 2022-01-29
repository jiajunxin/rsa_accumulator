package accumulator

import (
	"math/big"

	"github.com/rsa_accumulator/dihash"
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

func Accumulate(g, power, n *big.Int) *big.Int {
	var ret big.Int
	ret.Exp(g, power, n)
	return &ret
}

// This function can be done in a smarter way, but not the current version
func PreCompute() []big.Int {
	return preCompute(PreComputeSize)
}

func preCompute(preComputeSize int) []big.Int {
	trustedSetup := *TrustedSetup()

	ret := make([]big.Int, preComputeSize)

	ret[0] = trustedSetup.G
	ret[1] = *Accumulate(&trustedSetup.G, dihash.Delta, &trustedSetup.N)
	for i := 2; i < preComputeSize; i++ {
		ret[i] = *Accumulate(&ret[i-1], dihash.Delta, &trustedSetup.N)
	}
	return ret
}
