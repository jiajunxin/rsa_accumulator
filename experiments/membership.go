package experiments

import (
	"fmt"
	"math/big"
	"runtime"
	"time"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

// AccAndProveParallel recursively generates the accumulator with all the memberships precomputed in parallel
func AccAndProveParallel(set []string, encodeType accumulator.EncodeType, setup *accumulator.Setup) (*big.Int, []*big.Int) {
	startingTime := time.Now().UTC()
	rep := accumulator.GenRepresentatives(set, encodeType)
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

// ProveMembership uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlog(n))
func ProveMembership(base, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) <= 2 {
		return handleSmallSet(base, N, set)
	}
	// if len(set) <= 1024 {
	// 	return set
	// }
	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftBase := *accumulateNew(base, N, set[len(set)/2:])
	rightBase := *accumulateNew(base, N, set[0:len(set)/2])
	proofs := ProveMembership(&leftBase, N, set[0:len(set)/2])
	proofs = append(proofs, ProveMembership(&rightBase, N, set[len(set)/2:])...)
	return proofs
}

func handleSmallSet(base, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) == 1 {
		ret := make([]*big.Int, 1)
		ret[0] = base
		return ret
	}
	if len(set) == 2 {
		ret := make([]*big.Int, 2)
		ret[0] = AccumulateNew(base, set[1], N)
		ret[1] = AccumulateNew(base, set[0], N)
		return ret
	}
	// Should never reach here
	fmt.Println("Error in handleSmallSet, set size =", len(set))
	panic("Error in handleSmallSet, set size")
}

// AccumulateNew calculates g^{power} mod N
func AccumulateNew(g, power, N *big.Int) *big.Int {
	ret := &big.Int{}
	ret.Set(g)
	ret.Exp(g, power, N)
	return ret
}

func accumulate(g, N *big.Int, set []*big.Int) *big.Int {
	for _, v := range set {
		g.Exp(g, v, N)
	}
	return g
}

// AccumulateParallel is a test function for Parallelly accumulating elements
// func AccumulateParallel(g, N *big.Int, set []*big.Int) *big.Int {
// 	// test function. Just parallel for 4 cores.
// 	var prod big.Int
// 	prod.SetInt64(1)
// 	for _, v := range set {
// 		prod.Mul(&prod, v)
// 	}
// 	bitLength := prod.BitLen()
// 	// find the decimal for the bit length
// 	g.Exp(g, &prod, N)
// 	return g
// }

func accumulateNew(g, N *big.Int, set []*big.Int) *big.Int {
	acc := &big.Int{}
	acc.Set(g)
	return accumulate(acc, N, set)
}
