package accumulator

import (
	"context"
	"fmt"
	"math/big"
	"runtime"
	"time"
)

// AccAndProveParallel recursively generates the accumulator with all the memberships precomputed in parallel
func AccAndProveParallel(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	startingTime := time.Now().UTC()
	rep := GenRepresentatives(set, encodeType)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running GenRepresentatives Takes [%.3f] Seconds \n",
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
	startingTime := time.Now().UTC()
	rep := GenRepresentatives(set, encodeType)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running GenRepresentatives Takes [%.3f] Seconds \n",
		duration.Seconds())
	proofs := ProveMembershipIterParallel(*setup.G, setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// ProveMembershipParallel uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlog(n))
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
	startingTime := time.Now().UTC()
	leftBase, rightBase := calBaseParallel(base, N, set)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallel for the first layer with 2 cores Takes [%.3f] Seconds \n",
		duration.Seconds())
	c1 := make(chan []*big.Int)
	c2 := make(chan []*big.Int)
	go proveMembershipWithChan(leftBase, N, set[0:len(set)/2], limit, c1)
	go proveMembershipWithChan(rightBase, N, set[len(set)/2:], limit, c2)
	proofs1 := <-c1
	proofs2 := <-c2

	proofs1 = append(proofs1, proofs2...)
	return proofs1
}

// proveMembership uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlog(n))
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

	// if len(set) <= 1024 {
	// 	c <- set[:]
	// 	//c <- handleSmallSet(base, N, set)
	// 	close(c)
	// 	return
	// }

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
	initNodeChan := make(chan *proofNode, numWorkers)
	go initMembershipProofs(&base, N, set, 0, len(set),
		numWorkerPowerOfTwo, 0, initNodeChan)

	receivers := make(chan parallelReceiver, numWorkers)
	var cnt int
	for node := range initNodeChan {
		go func(node *proofNode) {
			receivers <- parallelReceiver{
				left:   node.left,
				right:  node.right,
				proofs: proveMembershipIter(*node.proof, N, set, node.left, node.right),
			}
		}(node)
		cnt++
		if cnt == numWorkers {
			close(initNodeChan)
		}
	}

	proofChan := make(chan []*big.Int)
	go func() {
		var cnt int
		proofs := make([]*big.Int, len(set))
		for r := range receivers {
			copy(proofs[r.left:r.right], r.proofs)
			cnt++
			if cnt != numWorkers {
				continue
			}
			close(receivers)
			proofChan <- proofs
			close(proofChan)
			return
		}
	}()

	return <-proofChan
}

func calNumWorkers() (int, int) {
	numWorkersPowerOfTwo := 0
	numWorkers := 1
	numCPUs := runtime.NumCPU()
	for numWorkers <= numCPUs {
		numWorkersPowerOfTwo++
		numWorkers *= 2
	}
	fmt.Printf("CPU Number: %d, Number of Workers: %d\n", numCPUs, numWorkers/2)
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
		header = &proofNode{
			left:  left,
			right: right,
			proof: &base,
		}
		iter       = header
		finishFlag = true
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sendChan := make(chan sendParam)
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
			iter = insertNewProofNodeWithChan(iter, N, set, sendChan, iterChan)
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
	left, right, powerOfTwo, depth int, initNodeChan chan *proofNode) {
	if depth > powerOfTwo {
		return
	}
	if depth == powerOfTwo {
		initNodeChan <- &proofNode{
			left:  left,
			right: right,
			proof: base,
		}
		return
	}
	mid := left + (right-left)/2
	go func() {
		proof1 := accumulateNew(base, N, set[left:mid])
		go initMembershipProofs(proof1, N, set, mid, right,
			powerOfTwo, depth+1, initNodeChan)
	}()
	go func() {
		proof2 := accumulateNew(base, N, set[mid:right])
		go initMembershipProofs(proof2, N, set, left, mid,
			powerOfTwo, depth+1, initNodeChan)
	}()
}

func insertNewProofNodeWithChan(iter *proofNode, N *big.Int, set []*big.Int,
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
