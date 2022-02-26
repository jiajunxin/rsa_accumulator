package accumulator

import (
	"fmt"
	"math/big"
)

func init() {
	_ = Min2048.Lsh(one, securityPara-1)
}

// TrustedSetup returns a pointer to AccumulatorSetup with 2048 bits key length
func TrustedSetup() *Setup {
	var ret Setup
	ret.N.SetString(N2048String, 10)
	ret.G.SetString(G2048String, 10)
	return &ret
}

// GenRepersentatives generates different representatives that can be inputted into RSA accumulator
func GenRepersentatives(set []string, encodeType EncodeType) []*big.Int {
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
func AccAndProve(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembership(&setup.G, &setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := Accumulate(proofs[0], rep[0], &setup.N)

	return acc, proofs
}

// AccAndProveParallel generates the accumulator with all the memberships precomputed in parallel
func AccAndProveParallel(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembershipParallel(&setup.G, &setup.N, rep, 4)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := Accumulate(proofs[0], rep[0], &setup.N)

	return acc, proofs
}

// AccAndProveIter iteratively generates the accumulator with all the memberships precomputed
func AccAndProveIter(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembershipIter(&setup.G, &setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := Accumulate(proofs[0], rep[0], &setup.N)

	return acc, proofs
}

// ProveMembership uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlogn)
func ProveMembership(base, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) <= 2 {
		return handleSmallSet(base, N, set)
	}
	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftBase := *accumulate(set[len(set)/2:], base, N)
	rightBase := *accumulate(set[0:len(set)/2], base, N)
	proofs := ProveMembership(&leftBase, N, set[0:len(set)/2])
	proofs = append(proofs, ProveMembership(&rightBase, N, set[len(set)/2:])...)
	return proofs
}

// ProveMembershipParallel uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlogn)
// It uses at most O(2^limit) Goroutines
func ProveMembershipParallel(base, N *big.Int, set []*big.Int, limit uint16) []*big.Int {
	if 0 == limit {
		return ProveMembership(base, N, set)
	}
	limit--
	if len(set) <= 2 {
		return handleSmallSet(base, N, set)
	}

	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftBase, rightBase := calBaseParallel(base, N, set)
	c1 := make(chan []*big.Int)
	c2 := make(chan []*big.Int)
	go proveMembershipWithChan(leftBase, N, set[0:len(set)/2], limit, c1)
	go proveMembershipWithChan(rightBase, N, set[len(set)/2:], limit, c2)
	proofs1 := <-c1
	proofs2 := <-c2

	proofs1 = append(proofs1, proofs2...)
	return proofs1
}

// proveMembership uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlogn)
func proveMembershipWithChan(base, N *big.Int, set []*big.Int, limit uint16, c chan []*big.Int) {
	if 0 == limit {
		c <- ProveMembership(base, N, set)
		close(c)
		return
	}
	limit--

	if len(set) <= 2 {
		c <- handleSmallSet(base, N, set)
		close(c)
		return
	}

	leftBase := *accumulate(set[len(set)/2:], base, N)
	rightBase := *accumulate(set[0:len(set)/2], base, N)
	c1 := make(chan []*big.Int)
	c2 := make(chan []*big.Int)
	go proveMembershipWithChan(&leftBase, N, set[0:len(set)/2], limit, c1)
	go proveMembershipWithChan(&rightBase, N, set[len(set)/2:], limit, c2)
	proofs1 := <-c1
	proofs2 := <-c2
	proofs1 = append(proofs1, proofs2...)
	c <- proofs1
	close(c)
}

func calBaseParallel(base, N *big.Int, set []*big.Int) (*big.Int, *big.Int) {
	// the left part of proof need to accumulate the right part of the set, vice versa.
	c1 := make(chan *big.Int)
	c2 := make(chan *big.Int)
	go accumulateWithChan(set[len(set)/2:], base, N, c1)
	go accumulateWithChan(set[0:len(set)/2], base, N, c2)
	leftBase, rightBase := <-c1, <-c2
	return leftBase, rightBase
}

// ProofIterator is the linked-list node for iterating proofs
type proofIterator struct {
	left  int // left index of proofs
	right int // right index of proofs
	proof big.Int
	next  *proofIterator
}

// ProveMembershipIter uses divide-and-conquer method to pre-compute the all membership proofs iteratively
func ProveMembershipIter(base, N *big.Int, set []*big.Int) []*big.Int {
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
	proofs := make([]*big.Int, len(set))
	for i := 0; i < len(set); i++ {
		proofs[i] = &dummy.next.proof
		dummy = dummy.next
	}
	return proofs
}

func accumulate(set []*big.Int, g, N *big.Int) *big.Int {
	var acc big.Int
	acc.Set(g)
	for _, v := range set {
		acc.Exp(&acc, v, N)
	}
	return &acc
}

func accumulateWithChan(set []*big.Int, g, N *big.Int, c chan *big.Int) {
	var acc big.Int
	acc.Set(g)
	for _, v := range set {
		acc.Exp(&acc, v, N)
	}
	c <- &acc
	close(c)
}

func handleSmallSet(base, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) == 1 {
		ret := make([]*big.Int, 1)
		ret[0] = base
		return ret
	}
	if len(set) == 2 {
		ret := make([]*big.Int, 2)
		ret[0] = Accumulate(base, set[1], N)
		ret[1] = Accumulate(base, set[0], N)
		return ret
	}
	// Should never reach here
	fmt.Println("Error in handleSmallSet, set size =", len(set))
	panic("Error in handleSmallSet, set size")
}

// Accumulate calculates g^{power} mod N
func Accumulate(g, power, N *big.Int) *big.Int {
	var ret big.Int
	ret.Exp(g, power, N)
	return &ret
}
