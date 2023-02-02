package experiments

import (
	"fmt"
	"math/big"
	"math/bits"
	"sync"
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
func TestRSAMembershipPreComputeMultiDIParallel(setSize int, limit int) {
	fmt.Println("Test set size = ", setSize)
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
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r1, setup.N, rep[:setSize], limit, table, c1)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r2, setup.N, rep[setSize:2*setSize], limit, table, c2)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r3, setup.N, rep[2*setSize:], limit, table, c3)
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

// Test the time to pre-compute all the membership proofs of one RSA accumulator, for different set size, with single core
func TestRSAMembershipPreComputeDIParallel(setSize int, limit int) {
	fmt.Println("Test set size = ", setSize)
	fmt.Println("Core limit = 2^", limit)
	fmt.Println("GenRepresentatives with DIHashFromPoseidon")
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()

	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	// generate a zero-knowledge RSA accumulator
	r1 := accumulator.GenRandomizer()

	maxLen := setSize * 1024 / bits.UintSize
	startingTime := time.Now().UTC()
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running PreComputeTable Takes [%.3f] Seconds \n", duration.Seconds())

	startingTime = time.Now().UTC()
	proofs := accumulator.ProveMembershipParallelWithTableWithRandomizer(setup.G, r1, setup.N, rep[:setSize], limit, table)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallelWithTableWithRandomizer for three RSA accumulators Takes [%.3f] Seconds \n", duration.Seconds())
	startingTime = time.Now().UTC()
	func() {
		tempProof := proofs[0]
		_ = tempProof.BitLen() // this line is simply used to allow accessing tempProof
	}()
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Online phase to get one membership proof, Takes [%d] Nanoseconds \n", duration.Nanoseconds())
}

// Test the time to pre-compute all the membership proofs of one RSA accumulator, for different set size, with single core
func TestPreComputeTableSize(setSize int) {
	//setSize := 65536 // 2 ^ 16 65536
	fmt.Println("Test set size = ", setSize)
	setup := *accumulator.TrustedSetup()

	maxLen := setSize * 256 / bits.UintSize
	startingTime := time.Now().UTC()
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running PreComputeTable Takes [%.3f] Seconds \n", duration.Seconds())
	fmt.Println("The table size = ", table.TableSize, "rows, ", bits.UintSize, " columns, each element size = ", bits.UintSize)
	fmt.Println("Totally ", table.TableSize*bits.UintSize*bits.UintSize/8, "bytes")
}

func TestDifferentMembershipForDI() {
	TestRSAMembershipPreComputeDIParallel(16384, 0) //2^14, 1 core
	TestRSAMembershipPreComputeDIParallel(16384, 2) //2^14, 4 cores
	TestRSAMembershipPreComputeDIParallel(16384, 4) //2^14, 16 cores

	TestRSAMembershipPreCompute(16384)                   //2^14, 1 core
	TestRSAMembershipPreComputeMultiDIParallel(16384, 0) //2^14, 3 cores
	TestRSAMembershipPreComputeMultiDIParallel(16384, 2) //2^14, 12 cores

	TestRSAMembershipPreComputeDIParallel(65536, 0) //2^16, 1 core
	TestRSAMembershipPreComputeDIParallel(65536, 2) //2^16, 4 cores
	TestRSAMembershipPreComputeDIParallel(65536, 4) //2^16, 16 cores

	TestRSAMembershipPreCompute(65536)                   //2^16, 1 core
	TestRSAMembershipPreComputeMultiDIParallel(65536, 0) //2^16, 3 cores
	TestRSAMembershipPreComputeMultiDIParallel(65536, 2) //2^16, 12 cores

	TestRSAMembershipPreComputeDIParallel(262144, 0) //2^18, 1 core
	TestRSAMembershipPreComputeDIParallel(262144, 2) //2^18, 4 cores
	TestRSAMembershipPreComputeDIParallel(262144, 4) //2^18, 16 cores

	TestRSAMembershipPreCompute(262144)                   //2^18, 1 core
	TestRSAMembershipPreComputeMultiDIParallel(262144, 0) //2^18, 3 cores
	TestRSAMembershipPreComputeMultiDIParallel(262144, 2) //2^18, 12 cores
}

func preComputeMultiDIParallel(setSize int, limit int, table *multiexp.PreTable) {
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()

	rep := accumulator.GenRepresentatives(set, accumulator.MultiDIHashFromPoseidon)
	// generate a zero-knowledge RSA accumulator
	r1 := accumulator.GenRandomizer()
	r2 := accumulator.GenRandomizer()
	r3 := accumulator.GenRandomizer()

	c1 := make(chan []*big.Int)
	c2 := make(chan []*big.Int)
	c3 := make(chan []*big.Int)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r1, setup.N, rep[:setSize], limit, table, c1)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r2, setup.N, rep[setSize:2*setSize], limit, table, c2)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r3, setup.N, rep[2*setSize:], limit, table, c3)
	proofs1 := <-c1
	proofs2 := <-c2
	proofs3 := <-c3
	func() {
		tempProof := proofs1[0]
		_ = tempProof.BitLen()
		tempProof = proofs2[0]
		_ = tempProof.BitLen()
		tempProof = proofs3[0]
		_ = tempProof.BitLen() // this line is simply used to allow accessing tempProof
	}()
}

// We only test till 5 sets with 12*5 threads because in our test environment we have 64 threads.
// You may need to adjust the parameters based on your own test environments.
func TestPreComputeMultiDIParallelRepeated() {
	setSize := 65536 //2^16, 12 cores
	setup := *accumulator.TrustedSetup()
	maxLen := setSize * 256 / bits.UintSize //256 comes from the length of each multiDI hash
	tables := make([]*multiexp.PreTable, 5)
	for i := 0; i < 5; i++ {
		tables[i] = multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	}
	fmt.Println("TestPreComputeMultiDIParallelRepeated, Test set size = ", setSize)
	fmt.Println("First trial: run PreComputeMultiDIParallel with 12 cores for 1 set of", setSize, " elements")
	var wg sync.WaitGroup
	startingTime := time.Now().UTC()
	repeatNum := 1
	wg.Add(repeatNum)
	for i := 0; i < repeatNum; i++ {
		defer wg.Done()
		preComputeMultiDIParallel(65536, 2, tables[i]) //2^16, 12 cores
	}
	wg.Wait()
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running First trial Takes [%.3f] Seconds \n", duration.Seconds())

	fmt.Println("Second trial: run PreComputeMultiDIParallel with 12*2 cores for 2 sets of", setSize, " elements")
	startingTime = time.Now().UTC()
	repeatNum = 2
	wg.Add(repeatNum)
	for i := 0; i < repeatNum; i++ {
		defer wg.Done()
		preComputeMultiDIParallel(65536, 2, tables[i]) //2^16, 12 cores
	}
	wg.Wait()
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running second trial Takes [%.3f] Seconds \n", duration.Seconds())

	fmt.Println("Third trial: run PreComputeMultiDIParallel with 12*3 cores for 3 sets of", setSize, " elements")
	startingTime = time.Now().UTC()
	repeatNum = 3
	wg.Add(repeatNum)
	for i := 0; i < repeatNum; i++ {
		defer wg.Done()
		preComputeMultiDIParallel(65536, 2, tables[i]) //2^16, 12 cores
	}
	wg.Wait()
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running third trial Takes [%.3f] Seconds \n", duration.Seconds())

	fmt.Println("Fourth trial: run PreComputeMultiDIParallel with 12*4 cores for 4 sets of", setSize, " elements")
	startingTime = time.Now().UTC()
	repeatNum = 4
	wg.Add(repeatNum)
	for i := 0; i < repeatNum; i++ {
		defer wg.Done()
		preComputeMultiDIParallel(65536, 2, tables[i]) //2^16, 12 cores
	}
	wg.Wait()
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running fourth trial Takes [%.3f] Seconds \n", duration.Seconds())

	fmt.Println("Fifth trial: run PreComputeMultiDIParallel with 12*5 cores for 5 sets of", setSize, " elements")
	startingTime = time.Now().UTC()
	repeatNum = 5
	wg.Add(repeatNum)
	for i := 0; i < repeatNum; i++ {
		defer wg.Done()
		preComputeMultiDIParallel(65536, 2, tables[i]) //2^16, 12 cores
	}
	wg.Wait()
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running fifth trial Takes [%.3f] Seconds \n", duration.Seconds())
}

func TestDifferentPrecomputationTableSize() {
	TestPreComputeTableSize(1024)    //2^10
	TestPreComputeTableSize(4096)    //2^12
	TestPreComputeTableSize(16384)   //2^14
	TestPreComputeTableSize(65536)   //2^16
	TestPreComputeTableSize(262144)  //2^18
	TestPreComputeTableSize(1048576) //2^20
}

// Test the time to pre-compute all the membership proofs of one RSA accumulator, for different set size, with single core
func TestNotusParallel(setSize int, limit int) {
	fmt.Println("Test set size = ", setSize)
	fmt.Println("Core limit = 2^", limit)
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
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r1, setup.N, rep[:setSize], limit, table, c1)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r2, setup.N, rep[setSize:2*setSize], limit, table, c2)
	go accumulator.ProveMembershipParallelWithTableWithRandomizerWithChan(setup.G, r3, setup.N, rep[2*setSize:], limit, table, c3)
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
