package accumulator

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"
)

// AccAndProveParallel generates the accumulator with all the memberships precomputed in parallel
func AccAndProveParallel(set []string, encodeType EncodeType, setup *AccumulatorSetup) (*big.Int, []*big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembershipParallel(setup.G, setup.N, rep, 4)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// AccAndProveParallel concurrently generates the accumulator with all the memberships precomputed
func AccAndProveIterParallel(set []string, encodeType EncodeType,
	setup *AccumulatorSetup) (*big.Int, []*big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembershipIterParallel(*setup.G, setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// ProveMembershipParallel uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlogn)
// It uses at most O(2^limit) Goroutines
func ProveMembershipParallel(base, N *big.Int, set []*big.Int, limit uint16) []*big.Int {
	if 0 == limit {
		return ProveMembership(base, N, set)
	}
	limit--
	fmt.Println("limit = ", limit)
	if len(set) == 0 {
		fmt.Println("Errorwwwwwwwwwwwwwwwwwwww")
	}
	if len(set) <= 2 {
		return handleSmallSet(base, N, set)
	}

	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftBase := *accumulate(base, N, set[len(set)/2:])
	rightBase := *accumulate(base, N, set[0:len(set)/2])
	//leftBase, rightBase := calBaseParallel(base, N, set)
	c3 := make(chan []*big.Int)
	c4 := make(chan []*big.Int)
	go proveMembershipWithChan(&leftBase, N, set[0:len(set)/2], limit, c3)
	go proveMembershipWithChan(&rightBase, N, set[len(set)/2:], limit, c4)
	proofs1 := <-c3
	proofs2 := <-c4

	proofs1 = append(proofs1, proofs2...)
	return proofs1
}

// proveMembership uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlogn)
func proveMembershipWithChan(base, N *big.Int, set []*big.Int, limit uint16, c chan []*big.Int) {
	if limit == 0 {
		c <- ProveMembership(base, N, set)
		close(c)
		return
	}
	limit--
	fmt.Println("limit = ", limit)
	if len(set) == 0 {
		fmt.Println("wwwwwwwwwwwwwwwwwwww")
	}
	if len(set) <= 2 {
		c <- handleSmallSet(base, N, set)
		close(c)
		return
	}

	// c1 := make(chan *big.Int)
	// c2 := make(chan *big.Int)
	// go accumulateWithChan(set[len(set)/2:], base, N, c1)
	// go accumulateWithChan(set[0:len(set)/2], base, N, c2)
	leftBase, rightBase := calBaseParallel(base, N, set)
	c3 := make(chan []*big.Int)
	c4 := make(chan []*big.Int)
	go proveMembershipWithChan(leftBase, N, set[0:len(set)/2], limit, c3)
	go proveMembershipWithChan(rightBase, N, set[len(set)/2:], limit, c4)
	proofs1 := <-c3
	proofs2 := <-c4
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

func accumulateWithChan(set []*big.Int, g, N *big.Int, c chan *big.Int) {
	var acc big.Int
	acc.Set(g)
	for _, v := range set {
		acc.Exp(&acc, v, N)
	}
	c <- &acc
	close(c)
}

type parallelReceiver struct {
	left   int
	right  int
	proofs []*big.Int
}

// ProveMembershipParallel uses divide-and-conquer method to pre-compute the all membership proofs
// iteratively and concurrently
func ProveMembershipIterParallel(base big.Int, N *big.Int, set []*big.Int) []*big.Int {
	numWorkers, numWorkerPowerOfTwo := calNumWorkers()
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
			iter = insertNewProofNodeParallel(iter, N, set)
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

func calNumWorkers() (int, int) {
	numWorkersPowerOfTwo := 0
	numWorkers := 1
	numCPUs := runtime.NumCPU()
	for numWorkers < numCPUs {
		numWorkersPowerOfTwo++
		numWorkers *= 2
	}
	return numWorkers / 2, numWorkersPowerOfTwo - 1
}

func insertNewProofNodeParallel(iter *proofNode, N *big.Int, set []*big.Int) *proofNode {
	left := iter.left
	right := iter.right
	mid := left + (right-left)/2
	newProofNodeChan := make(chan *big.Int)
	iterChan := make(chan *big.Int)
	go func() {
		newProofNodeChan <- accumulateNew(iter.proof, N, set[left:mid])
	}()
	go func() {
		iterChan <- accumulateNew(iter.proof, N, set[mid:right])
	}()
	newProofNode := &proofNode{
		left:  mid,
		right: right,
		proof: <-newProofNodeChan,
		next:  iter.next,
	}
	iter.left = left
	iter.right = mid
	iter.proof = <-iterChan
	iter.next = newProofNode
	return newProofNode.next
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
			iter = insertNewProofNodeParallel(iter, N, set)
			finishFlag = true
		}
	}

	proofs := make([]*big.Int, 0, len(set))
	for iter = header; iter != nil; iter = iter.next {
		proofs = append(proofs, iter.proof)
	}
	return proofs
}
