package accumulator

import (
	"math"
	"math/big"
	"sync"
)

const (
	numWorkerPowerOfTwo = 3
)

// AccAndProveParallel concurrently generates the accumulator with all the memberships precomputed
func AccAndProveParallel(set []string, encodeType EncodeType, setup *AccumulatorSetup) (*big.Int, []*big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembershipParallel(setup.G, &setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], &setup.N)

	return acc, proofs
}

type parallelReceiver struct {
	left   int
	right  int
	proofs []*big.Int
}

// ProveMembershipParallel uses divide-and-conquer method to pre-compute the all membership proofs iteratively
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
			iter = insertNewProofNode(iter, N, set)
		}
	}

	receivers := make(chan parallelReceiver, numWorkers)
	wg := &sync.WaitGroup{}
	wg.Add(numWorkers)
	iter = header
	for i := 0; i < numWorkers; i++ {
		go func(iter *proofNode) {
			defer wg.Done()
			receivers <- parallelReceiver{
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
			if iter.right-iter.left <= 1 {
				iter = iter.next
				continue
			}
			iter = insertNewProofNode(iter, N, set)
			finishFlag = true
		}
	}

	proofs := make([]*big.Int, 0, len(set))
	for iter = header; iter != nil; iter = iter.next {
		proofs = append(proofs, iter.proof)
	}
	return proofs
}
