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

// AccAndProveIter iteratively generates the accumulator with all the memberships precomputed
func AccAndProveIter(set []string, encodeType EncodeType, setup *AccumulatorSetup) (*big.Int, []big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembershipIter(&setup.G, &setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := Accumulate(&proofs[0], &rep[0], &setup.N)

	return acc, proofs
}

// ProveMembership uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlogn)
func ProveMembership(base, N *big.Int, set []big.Int) []big.Int {
	if len(set) == 1 {
		ret := make([]big.Int, 1)
		ret[0] = *base
		return ret
	}
	if len(set) == 2 {
		ret := make([]big.Int, 2)
		ret[0] = *Accumulate(base, &set[1], N)
		ret[1] = *Accumulate(base, &set[0], N)
		return ret
	}
	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftBase := *accumulate(set[len(set)/2:], base, N)
	rightBase := *accumulate(set[0:len(set)/2], base, N)
	proofs := ProveMembership(&leftBase, N, set[0:len(set)/2])
	proofs = append(proofs, ProveMembership(&rightBase, N, set[len(set)/2:])...)
	return proofs
}

// ProofIterator is the linked-list node for iterating proofs
type proofIterator struct {
	left  int // left index of proofs
	right int // right index of proofs
	proof big.Int
	next  *proofIterator
}

// ProveMembershipIter uses divide-and-conquer method to pre-compute the all membership proofs iteratively
func ProveMembershipIter(base, N *big.Int, set []big.Int) []big.Int {
	var (
		dummy *proofIterator = &proofIterator{
			next: &proofIterator{
				left:  0,
				right: len(set),
				proof: *base,
				next:  nil,
			},
		}
		prev       *proofIterator = dummy
		iter       *proofIterator = dummy.next
		finishFlag bool
	)
	for {
		for iter != nil {
			left := iter.left
			right := iter.right
			if right-left <= 1 {
				iter = iter.next
				prev = prev.next
				continue
			}
			mid := right - (right-left)/2
			acc := iter.proof
			secondNewProof := &proofIterator{
				left:  mid,
				right: right,
				proof: *accumulate(set[left:mid], &acc, N),
				next:  iter.next,
			}
			firstNewProof := &proofIterator{
				left:  left,
				right: mid,
				proof: *accumulate(set[mid:right], &acc, N),
				next:  secondNewProof,
			}
			prev.next = firstNewProof
			iter = iter.next
			prev = secondNewProof
			finishFlag = true
		}
		if !finishFlag {
			break
		}
		finishFlag = false
		prev = dummy
		iter = dummy.next
	}
	proofs := make([]big.Int, len(set))
	for i := 0; i < len(set); i++ {
		proofs[i] = dummy.next.proof
		dummy = dummy.next
	}
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
