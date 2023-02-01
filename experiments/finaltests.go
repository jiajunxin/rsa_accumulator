package experiments

import (
	"fmt"
	"math/big"
	"math/bits"
	"time"

	"github.com/jiajunxin/multiexp"
	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

// TestBasiczkRSA test a naive case of zero-knowledge RSA accumulator
func TestBasiczkRSA() {
	setSize := 65536 // 2 ^ 16 65536
	fmt.Println("Test set size = ", setSize)
	fmt.Println("GenRepresentatives with MultiDIHashFromPoseidon")
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	startingTime := time.Now().UTC()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	// generate a zero-knowledge RSA accumulator
	r := accumulator.GenRandomizer()
	randomizedbase := AccumulateNew(setup.G, r, setup.N)
	// calculate the exponentation
	exp := accumulator.SetProductRecursiveFast(rep)
	accumulator.AccumulateNew(randomizedbase, exp, setup.N)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generating a zero-knowledge RSA accumulator with set size = %d, takes [%.3f] Seconds \n", setSize, duration.Seconds())

	// startingTime = time.Now().UTC()
	// maxLen := setSize * 256 / bits.UintSize
	// table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	// duration = time.Now().UTC().Sub(startingTime)
	// fmt.Printf("Running PreCompute Takes [%.3f] Seconds \n", duration.Seconds())
	// startingTime = time.Now().UTC()
	// accumulator.ProveMembershipParallelWithTable(setup.G, setup.N, rep, 2, table)
	// duration = time.Now().UTC().Sub(startingTime)
	// fmt.Printf("Running ProveMembershipParallelWithTable Takes [%.3f] Seconds \n", duration.Seconds())
}

// Test the time to pre-compute all the membership proofs of one RSA accumulator, for different set size, with single core
func TestRSAMembershipPreCompute(setSize int) {
	//setSize := 65536 // 2 ^ 16 65536
	fmt.Println("Test set size = ", setSize)
	fmt.Println("GenRepresentatives with MultiDIHashFromPoseidon")
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()

	rep := accumulator.GenRepresentatives(set, accumulator.MultiDIHashFromPoseidon)
	// generate a zero-knowledge RSA accumulator
	r1 := accumulator.GenRandomizer()
	r2 := accumulator.GenRandomizer()
	r3 := accumulator.GenRandomizer()

	maxLen := setSize * 256 / bits.UintSize
	startingTime := time.Now().UTC()
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running PreComputeTable Takes [%.3f] Seconds \n", duration.Seconds())

	startingTime = time.Now().UTC()
	proofs1 := accumulator.ProveMembershipParallelWithTableWithRandomizer(setup.G, r1, setup.N, rep[:setSize], 0, table)
	proofs2 := accumulator.ProveMembershipParallelWithTableWithRandomizer(setup.G, r2, setup.N, rep[setSize:2*setSize], 0, table)
	proofs3 := accumulator.ProveMembershipParallelWithTableWithRandomizer(setup.G, r3, setup.N, rep[2*setSize:], 0, table)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallelWithTableWithRandomizer with single core for three RSA accumulators Takes [%.3f] Seconds \n", duration.Seconds())
	startingTime = time.Now().UTC()
	func() {
		tempProof := proofs1[0]
		_ = tempProof.BitLen()
		tempProof = proofs2[0]
		_ = tempProof.BitLen()
		tempProof = proofs3[0]
		_ = tempProof.BitLen() // this line is simply used to allow accessing tempProof
	}()
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Online phase to get one membership proof, Takes [%d] Nanoseconds \n", duration.Nanoseconds())
}

func TestDifferentMembership() {
	TestRSAMembershipPreCompute(1024)    //2^10
	TestRSAMembershipPreCompute(4096)    //2^12
	TestRSAMembershipPreCompute(16384)   //2^14
	TestRSAMembershipPreCompute(65536)   //2^16
	TestRSAMembershipPreCompute(262144)  //2^18
	TestRSAMembershipPreCompute(1048576) //2^20
}

// Test the time to pre-compute all the membership proofs of one RSA accumulator, for different set size, with single core
func TestRSAMembershipPreComputeMultiDIParallel(setSize int) {
	//setSize := 65536 // 2 ^ 16 65536
	fmt.Println("Test set size = ", setSize)
	fmt.Println("GenRepresentatives with MultiDIHashFromPoseidon")
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()

	rep := accumulator.GenRepresentatives(set, accumulator.MultiDIHashFromPoseidon)
	// generate a zero-knowledge RSA accumulator
	r1 := accumulator.GenRandomizer()
	r2 := accumulator.GenRandomizer()
	r3 := accumulator.GenRandomizer()

	maxLen := setSize * 256 / bits.UintSize
	startingTime := time.Now().UTC()
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running PreComputeTable Takes [%.3f] Seconds \n", duration.Seconds())

	c1 := make(chan []*big.Int)
	c2 := make(chan []*big.Int)
	c3 := make(chan []*big.Int)
	startingTime = time.Now().UTC()
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r1, setup.N, rep[:setSize], 2, table, c1)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r2, setup.N, rep[setSize:2*setSize], 2, table, c2)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r3, setup.N, rep[2*setSize:], 2, table, c3)
	proofs1 := <-c1
	proofs2 := <-c2
	proofs3 := <-c3
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallelWithTableWithRandomizer with 12 cores for three RSA accumulators Takes [%.3f] Seconds \n", duration.Seconds())
	startingTime = time.Now().UTC()
	func() {
		tempProof := proofs1[0]
		_ = tempProof.BitLen()
		tempProof = proofs2[0]
		_ = tempProof.BitLen()
		tempProof = proofs3[0]
		_ = tempProof.BitLen() // this line is simply used to allow accessing tempProof
	}()
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Online phase to get one membership proof, Takes [%d] Nanoseconds \n", duration.Nanoseconds())
}
