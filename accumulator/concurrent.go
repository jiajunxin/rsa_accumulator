package accumulator

import (
	"math"
	"math/big"
	"sync"
)

const (
	numWorkerPowerOfTwo = 3
)

func AccAndProveParallel(set []string, encodeType EncodeType, setup *AccumulatorSetup) (*big.Int, []*big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembershipParallel(setup.G, &setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], &setup.N)

	return acc, proofs
}

// ProveMembershipIter uses divide-and-conquer method to pre-compute the all membership proofs iteratively
func ProveMembershipParallel(base big.Int, N *big.Int, set []*big.Int) []*big.Int {
	numWorkers := int(math.Pow(2, numWorkerPowerOfTwo))
	if len(set) <= numWorkers*2 {
		return ProveMembershipIter(base, N, set)
	}
	var (
		header *proofNode = &proofNode{
			right: len(set),
			proof: &base,
		}
		iter *proofNode = header
	)

	for i := 0; i < numWorkerPowerOfTwo; i++ {
		iter = header
		for iter != nil {
			left := iter.left
			right := iter.right
			mid := right - (right-left)/2
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
			iter = newProofNode.next
		}
	}

	type receiver struct {
		left   int
		right  int
		proofs []*big.Int
	}
	receivers := make(chan receiver, numWorkers)
	wg := &sync.WaitGroup{}
	iter = header
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(iter *proofNode) {
			defer wg.Done()
			receivers <- receiver{
				left:   iter.left,
				right:  iter.right,
				proofs: proveMembershipIter(*iter.proof, N, set, iter.left, iter.right),
			}
		}(iter)
		iter = iter.next
	}
	wg.Wait()
	close(receivers)

	proofs := make([]*big.Int, len(set))
	for receiver := range receivers {
		copy(proofs[receiver.left:receiver.right], receiver.proofs)
	}
	return proofs
}

func proveMembershipIter(base big.Int, N *big.Int, set []*big.Int, left, right int) []*big.Int {
	if len(set) <= 0 {
		return nil
	}
	var (
		header *proofNode = &proofNode{
			left:  left,
			right: right,
			proof: &base,
		}
		iter       *proofNode = header
		finishFlag bool       = true
	)

	for finishFlag {
		finishFlag = false
		iter = header
		for iter != nil {
			left := iter.left
			right := iter.right
			if right-left <= 1 {
				iter = iter.next
				continue
			}
			mid := right - (right-left)/2
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
			iter = newProofNode.next
			finishFlag = true
		}
	}

	proofs := make([]*big.Int, 0, len(set))
	for iter = header; iter != nil; iter = iter.next {
		proofs = append(proofs, iter.proof)
	}
	return proofs
}
