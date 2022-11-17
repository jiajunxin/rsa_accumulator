package accumulator

import (
	"fmt"
	"math/big"
	"time"

	"github.com/remyoudompheng/bigfft"
)

func init() {
	_ = Min2048.Lsh(big1, RSABitLength-1)
}

// TrustedSetup returns a pointer to AccumulatorSetup with 2048 bits key length
func TrustedSetup() *Setup {
	ret := &Setup{
		N: &big.Int{},
		G: &big.Int{},
	}
	ret.N.SetString(N2048String, 10)
	ret.G.SetString(G2048String, 10)
	return ret
}

// GenRepresentatives generates different representatives that can be inputted into RSA accumulator
func GenRepresentatives(set []string, encodeType EncodeType) []*big.Int {
	switch encodeType {
	case HashToPrimeFromSha256:
		return genRepWithHashToPrimeFromSHA256(set)
	case DIHashFromPoseidon:
		return genRepWithDIHashFromPoseidon(set)
	case MultiDIHashFromPoseidon:
		return genRepWithMultiDIHashFromPoseidon(set)
	default:
		return genRepWithHashToPrimeFromSHA256(set)
	}
}

// AccAndProve generates the accumulator with all the memberships precomputed
func AccAndProve(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	startingTime := time.Now().UTC()
	rep := GenRepresentatives(set, encodeType)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running GenRepresentatives Takes [%.3f] Seconds \n",
		duration.Seconds())

	proofs := ProveMembership(setup.G, setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// AccAndProveIter iteratively generates the accumulator with all the memberships precomputed
func AccAndProveIter(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	rep := GenRepresentatives(set, encodeType)

	proofs := ProveMembershipIter(*setup.G, setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}

// ProveMembership uses divide-and-conquer method to pre-compute the all membership proofs in time O(nlog(n))
func ProveMembership(base, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) <= 4 {
		return handleSmallSet(base, N, set)
	}
	// if len(set) <= 1024 {
	// 	return set
	// }
	// the left part of proof need to accumulate the right part of the set, vice versa.
	leftProd := SetProductRecursiveFast(set[len(set)/2:])
	rightProd := SetProductRecursiveFast(set[0 : len(set)/2])
	bases := big.DoubleExp(base, leftProd, rightProd, N)
	// leftBase := accumulateNew(base, N, set[len(set)/2:])
	// rightBase := accumulateNew(base, N, set[0:len(set)/2])
	proofs := ProveMembership(bases[0], N, set[0:len(set)/2])
	proofs = append(proofs, ProveMembership(bases[1], N, set[len(set)/2:])...)
	return proofs
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
		iter       = header
		finishFlag = true
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

func insertNewProofNode(iter *proofNode, N *big.Int, set []*big.Int) *proofNode {
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

func handleSmallSet(base, N *big.Int, set []*big.Int) []*big.Int {
	if len(set) == 4 {
		// suppose the set is x0, x1, x2, x3, the membership for x0 is base^{x1x2x3}
		x0x1Prod := bigfft.Mul(set[0], set[1])
		x2x3Prod := bigfft.Mul(set[2], set[3])
		x0Prod := bigfft.Mul(x2x3Prod, set[1])
		x1Prod := bigfft.Mul(x2x3Prod, set[0])
		x2Prod := bigfft.Mul(x0x1Prod, set[3])
		x3Prod := bigfft.Mul(x0x1Prod, set[2])
		return big.FourFoldExp(base, N, []*big.Int{x0Prod, x1Prod, x2Prod, x3Prod})
	}
	if len(set) == 3 {
		// suppose the set is x0, x1, x2, the membership for x0 is base^{x1x2}
		x0Prod := bigfft.Mul(set[1], set[2])
		x1Prod := bigfft.Mul(set[0], set[2])
		x2Prod := bigfft.Mul(set[0], set[1])
		return append(big.DoubleExp(base, x0Prod, x1Prod, N), AccumulateNew(base, x2Prod, N))
	}
	if len(set) == 2 {
		return big.DoubleExp(base, set[1], set[0], N)
	}
	if len(set) == 1 {
		ret := make([]*big.Int, 1)
		ret[0] = base
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
