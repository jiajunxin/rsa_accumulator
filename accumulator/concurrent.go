package accumulator

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"
)

// AccAndProveParallel recursively generates the accumulator with all the memberships precomputed in parallel
func AccAndProveParallel(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	startingTime := time.Now().UTC()
	rep := GenRepersentatives(set, encodeType)
	endingTime := time.Now().UTC()
	var duration time.Duration = endingTime.Sub(startingTime)
	fmt.Printf("Running GenRepersentatives Takes [%.3f] Seconds \n",
		duration.Seconds())
	numWorkers, _ := calNumWorkers()
	proofs := ProveMembershipParallel(setup.G, setup.N, rep, numWorkers)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// AccAndProveIterParallel iteratively and concurrently generates the accumulator with all the memberships precomputed
func AccAndProveIterParallel(set []string, encodeType EncodeType,
	setup *Setup) (*big.Int, []*big.Int) {
	rep := GenRepersentatives(set, encodeType)

	proofs := ProveMembershipIterParallel(*setup.G, setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// ProveMembershipParallel uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlogn)
// It uses at most O(2^limit) Goroutines
func ProveMembershipParallel(base, N *big.Int, set []*big.Int, limit int) []*big.Int {
	if limit == 0 {
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
func proveMembershipWithChan(base, N *big.Int, set []*big.Int, limit int, c chan []*big.Int) {
	if limit == 0 {
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

	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftBase, rightBase := calBaseParallel(base, N, set)
	c1 := make(chan []*big.Int)
	c2 := make(chan []*big.Int)
	go proveMembershipWithChan(leftBase, N, set[0:len(set)/2], limit, c1)
	go proveMembershipWithChan(rightBase, N, set[len(set)/2:], limit, c2)
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

// ProveMembershipIterParallel uses divide-and-conquer method to pre-compute the all membership proofs
// iteratively and concurrently
func ProveMembershipIterParallel(base big.Int, N *big.Int, set []*big.Int) []*big.Int {
	numWorkers, numWorkerPowerOfTwo := calNumWorkers()
	if len(set) <= numWorkers*2 {
		return ProveMembershipIter(base, N, set)
	}
	initChans := make([]chan *proofNode, numWorkers)
	indexes := make([]int, numWorkers)
	for i := 0; i < numWorkers; i++ {
		initChans[i] = make(chan *proofNode)
		indexes[i] = i
	}
	initMembershipProofs(&base, N, set, 0, len(set),
		numWorkerPowerOfTwo, 0, indexes, initChans)
	header := <-initChans[0]
	iter := header
	for i := 1; i < numWorkers; i++ {
		iter.next = <-initChans[i]
		iter = iter.next
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

type sendParam struct {
	left  int
	right int
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sendChan := make(chan sendParam, 5)
	iterChan := make(chan *big.Int)
	go func() {
		for {
			select {
			case send := <-sendChan:
				iterChan <- accumulateNew(iter.proof, N, set[send.left:send.right])
			case <-ctx.Done():
				return
			}
		}
	}()
	for finishFlag {
		finishFlag = false
		iter = header
		for iter != nil {
			if iter.right-iter.left <= 1 {
				iter = iter.next
				continue
			}
			iter = insertNewProofNodeParallelWithChan(iter, N, set, sendChan, iterChan)
			finishFlag = true
		}
	}

	proofs := make([]*big.Int, 0, len(set))
	for iter = header; iter != nil; iter = iter.next {
		proofs = append(proofs, iter.proof)
	}
	return proofs
}

func initMembershipProofs(base, N *big.Int, set []*big.Int,
	left, right, numWorkerPowerOfTwo, depth int, indexes []int, initChans []chan *proofNode) {
	if depth > numWorkerPowerOfTwo {
		return
	}
	if depth == numWorkerPowerOfTwo {
		initChans[indexes[0]] <- &proofNode{
			left:  left,
			right: right,
			proof: base,
		}
		close(initChans[indexes[0]])
		return
	}
	mid := left + (right-left)/2
	idxMid := len(indexes) / 2
	go func() {
		proof1 := accumulateNew(base, N, set[left:mid])
		go initMembershipProofs(proof1, N, set, mid, right,
			numWorkerPowerOfTwo, depth+1, indexes[idxMid:], initChans)
	}()
	go func() {
		proof2 := accumulateNew(base, N, set[mid:right])
		go initMembershipProofs(proof2, N, set, left, mid,
			numWorkerPowerOfTwo, depth+1, indexes[:idxMid], initChans)
	}()
}

func insertNewProofNodeParallelWithChan(iter *proofNode, N *big.Int, set []*big.Int,
	sendChan chan<- sendParam, iterChan <-chan *big.Int) *proofNode {
	left := iter.left
	right := iter.right
	mid := left + (right-left)/2
	sendChan <- sendParam{left: mid, right: right}
	newProofNode := &proofNode{
		left:  mid,
		right: right,
		proof: accumulateNew(iter.proof, N, set[left:mid]),
		next:  iter.next,
	}
	iter.left = left
	iter.right = mid
	iter.proof = <-iterChan
	iter.next = newProofNode
	return newProofNode.next
}
