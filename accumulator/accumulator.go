package accumulator

import (
	"fmt"
	"math/big"
	"time"
)

// AccAndProve generates the accumulator with all the memberships precomputed
func AccAndProve(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	startingTime := time.Now().UTC()
	rep := HashEncode(set, encodeType)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running HashEncode Takes [%.3f] Seconds \n",
		duration.Seconds())

	proofs := ProveMembership(setup.G, setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// AccAndProveIter iteratively generates the accumulator with all the memberships precomputed
func AccAndProveIter(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	rep := HashEncode(set, encodeType)

	proofs := ProveMembershipIter(*setup.G, setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// ProveMembership uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlog(n))
func ProveMembership(base, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) <= 2 {
		return handleSmallSet(base, N, set)
	}
	// if len(set) <= 1024 {
	// 	return set
	// }
	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftBase := accumulateNew(base, N, set[len(set)/2:])
	rightBase := accumulateNew(base, N, set[0:len(set)/2])
	proofs := ProveMembership(leftBase, N, set[0:len(set)/2])
	proofs = append(proofs, ProveMembership(rightBase, N, set[len(set)/2:])...)
	return proofs
}

func handleSmallSet(base, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) == 1 {
		return []*big.Int{base}
	}
	// set length = 2
	return []*big.Int{
		AccumulateNew(base, set[1], N),
		AccumulateNew(base, set[0], N),
	}
}

// ProofNode is the linked-list node for iterating proofs
type proofNode struct {
	left  int // left index of proofs
	right int // right index of proofs
	proof *big.Int
	next  *proofNode
}

// ProveMembershipIter uses divide-and-conquer method to pre-compute the all membership proofs iteratively
func ProveMembershipIter(base big.Int, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) <= 0 {
		return nil
	}
	var (
		header = &proofNode{
			right: len(set),
			proof: &base,
		}
		iter            = header
		iterNotFinished = true
	)

	for iterNotFinished {
		iterNotFinished = false
		iter = header
		for iter != nil {
			if iter.right-iter.left <= 1 {
				iter = iter.next
				continue
			}
			iter = insertNewNode(iter, N, set)
			iterNotFinished = true
		}
	}

	proofs := make([]*big.Int, 0, len(set))
	for iter = header; iter != nil; iter = iter.next {
		proofs = append(proofs, iter.proof)
	}
	return proofs
}

func insertNewNode(iter *proofNode, N *big.Int, set []*big.Int) *proofNode {
	left := iter.left
	right := iter.right
	mid := left + (right-left)/2
	newProofNode := &proofNode{
		left:  mid,
		right: right,
		proof: accumulateNew(iter.proof, N, set[left:mid]),
		next:  iter.next,
	}
	iter.left = left
	iter.right = mid
	iter.proof = accumulate(iter.proof, N, set[mid:right])
	iter.next = newProofNode
	return newProofNode.next
}

// AccumulateNew calculates g^{power} mod N
func AccumulateNew(g, power, N *big.Int) *big.Int {
	ret := new(big.Int).Set(g)
	ret.Exp(g, power, N)
	return ret
}

func accumulate(g, N *big.Int, set []*big.Int) *big.Int {
	for _, v := range set {
		g.Exp(g, v, N)
	}
	return g
}

func accumulateNew(g, N *big.Int, set []*big.Int) *big.Int {
	acc := new(big.Int).Set(g)
	return accumulate(acc, N, set)
}
