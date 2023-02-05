package experiments

import (
	"fmt"
	"math/big"
	"math/bits"
	"sync"
	"time"

	"github.com/jiajunxin/multiexp"
	"github.com/jiajunxin/rsa_accumulator/accumulator"
	"github.com/remyoudompheng/bigfft"
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

func TestPreComputeMultiDIParallelRepeatedTogetherWithSNARK(setSize int) {
	setup := *accumulator.TrustedSetup()
	maxLen := setSize * 256 / bits.UintSize //256 comes from the length of each multiDI hash
	//tables := make([]*multiexp.PreTable, 8)
	fmt.Println("TestPreComputeMultiDIParallelRepeated, Test set size = ", setSize)
	fmt.Println("Generating precomputation tables")
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	var wg sync.WaitGroup
	fmt.Println("32 trial: run PreComputeMultiDIParallel with 3 cores for 32 sets of", setSize, " elements")
	startingTime := time.Now().UTC()
	repeatNum := 32
	wg.Add(repeatNum)
	for i := 0; i < repeatNum; i++ {
		go func(i int) {
			defer wg.Done()
			preComputeMultiDIParallel(setSize, 0, table)
		}(i)
	}
	wg.Wait()
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running the 32 th trial Takes [%.3f] Seconds \n", duration.Seconds())
}

func TestDifferentGroupingSize(setSize int) {
	max := 262144 //2^18
	setup := *accumulator.TrustedSetup()
	maxLen := setSize * 256 / bits.UintSize //256 comes from the length of each multiDI hash
	//tables := make([]*multiexp.PreTable, 8)
	fmt.Println("TestDifferentGroupingSize, Test set size = ", setSize)
	repeatNum := max / setSize
	fmt.Println("To have same number of users, repeatNum = ", repeatNum)
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	var wg sync.WaitGroup
	fmt.Println("Start timer for precomputation of membership proofs")
	startingTime := time.Now().UTC()
	wg.Add(repeatNum)
	for i := 0; i < repeatNum; i++ {
		go func(i int) {
			defer wg.Done()
			preComputeMultiDIParallel(setSize, 0, table)
		}(i)
	}
	wg.Wait()
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running the trial Takes [%.3f] Seconds \n", duration.Seconds())
}

func TestDifferentGroupSize() {
	TestDifferentGroupingSize(1024)   //2^10
	TestDifferentGroupingSize(4096)   //2^12
	TestDifferentGroupingSize(16384)  //2^14
	TestDifferentGroupingSize(65536)  //2^16
	TestDifferentGroupingSize(262144) //2^18
}

func TestDifferentPrecomputationTableSize() {
	TestPreComputeTableSize(1024)    //2^10
	TestPreComputeTableSize(4096)    //2^12
	TestPreComputeTableSize(16384)   //2^14
	TestPreComputeTableSize(65536)   //2^16
	TestPreComputeTableSize(262144)  //2^18
	TestPreComputeTableSize(1048576) //2^20
}

func PoKE(base, exp, newAcc, N *big.Int) {
	l := accumulator.HashToPrime(append(newAcc.Bytes(), base.Bytes()...))
	//fmt.Println("primeChallenge = ", l.String())
	remainder := big.NewInt(1)
	quotient := big.NewInt(1)
	quotient, remainder = quotient.DivMod(exp, l, remainder)
	_ = accumulator.AccumulateNew(base, quotient, N)
	//fmt.Println("Q = ", Q.String())
	//fmt.Println("r = ", remainder.String())
	// AccTest1 := accumulator.AccumulateNew(Q, l, N)
	// fmt.Println("Q^l = ", AccTest1.String())
	// AccTest2 := accumulator.AccumulateNew(base, remainder, N)
	// fmt.Println("g^r = ", AccTest2.String())
	// AccTest3 := AccTest1.Mul(AccTest1, AccTest2)
	// AccTest3.Mod(AccTest3, setup.N)
	// fmt.Println("Accumulator_test3 = ", AccTest3.String()) //AccTest3 should be the same as the AccOld
}

// RemovedSet is generated to test performance.
// The specific code here is to keep consistant with our SNARK experiments.
func TestNotusSingleThread(setSize, updatedSetSize int) {
	fmt.Println("Test code for Notus system. DO NOT use in production.")
	fmt.Println("Random numbers are fixed for test purpose!")
	fmt.Println("Test code with Single thread.")
	fmt.Println("Test set size = ", setSize)
	fmt.Println("Test updated set size = ", updatedSetSize)

	var currentEpoch int64
	currentEpoch = 500
	// generate the RemovedSet and insertedSet
	removed1 := make([]*big.Int, updatedSetSize)
	removed2 := make([]*big.Int, updatedSetSize)
	removed3 := make([]*big.Int, updatedSetSize)
	insert1 := make([]*big.Int, updatedSetSize)
	insert2 := make([]*big.Int, updatedSetSize)
	insert3 := make([]*big.Int, updatedSetSize)
	if updatedSetSize < 1 {
		panic("invalid updatedSetSize")
	}
	listID := make([]uint32, updatedSetSize)
	listValueOriginal := make([]uint32, updatedSetSize)
	listValueUpdated := make([]uint32, updatedSetSize)
	listLastUpdatedEpoch := make([]uint32, updatedSetSize)
	listPrevHash := make([]big.Int, updatedSetSize)
	var tempHashInput, tempValLeftShifted big.Int
	for i := 0; i < updatedSetSize; i++ {
		j := i*2 + 1
		listID[i] = uint32(j)
		listValueOriginal[i] = uint32(j)
		listLastUpdatedEpoch[i] = 10
		listPrevHash[i].SetInt64(int64(j))
		// input each into Poseidon Hash
		tempHashInput.Lsh(big.NewInt(int64(listID[i])), 64)
		tempValLeftShifted.Lsh(big.NewInt(int64(listValueOriginal[i])), 32)
		tempHashInput.Add(&tempHashInput, &tempValLeftShifted)
		tempHashInput.Add(&tempHashInput, big.NewInt(int64(listLastUpdatedEpoch[i])))
		removed1[i] = accumulator.PoseidonWith2Inputs([]*big.Int{&tempHashInput, &listPrevHash[i]})
		removed2[i] = accumulator.UniversalHashToInt(removed1[i])
		removed3[i] = accumulator.UniversalHashToInt(removed2[i])

		listValueUpdated[i] = uint32(j) // the updated value is same as the original value, which is allowed, for the simplicity of testing
		tempHashInput.Lsh(big.NewInt(int64(listID[i])), 64)
		tempValLeftShifted.Lsh(big.NewInt(int64(listValueUpdated[i])), 32)
		tempHashInput.Add(&tempHashInput, &tempValLeftShifted)
		tempHashInput.Add(&tempHashInput, big.NewInt(int64(currentEpoch)))
		insert1[i] = accumulator.PoseidonWith2Inputs([]*big.Int{&tempHashInput, removed1[i]})
		insert2[i] = accumulator.UniversalHashToInt(insert1[i])
		insert3[i] = accumulator.UniversalHashToInt(insert2[i])
	}

	unchangedSet := accumulator.GenBenchSet(setSize - updatedSetSize)
	rep := accumulator.GenRepresentatives(unchangedSet, accumulator.MultiDIHashFromPoseidon)

	var unchanged1, unchanged2, unchanged3 []*big.Int
	unchanged1 = rep[:setSize-updatedSetSize]
	unchanged2 = rep[setSize-updatedSetSize : 2*(setSize-updatedSetSize)]
	unchanged3 = rep[2*(setSize-updatedSetSize):]

	// This is also fro test purpose only.
	// We use Hash of tau as the random source, generate 6 different 2048 bits random numbers
	// Each 2048 bits random number is composed by 8 256 bits random number, therefore, we
	// need 48 256 bits random numbers.
	var tau, temp big.Int
	tau.SetString("13790045313639893950773977216617751241425918462119445469315488891147110571211", 10)
	poseidonHashResult := accumulator.PoseidonWith2Inputs([]*big.Int{&tau, &tau})
	tempRandomList := make([]*big.Int, 48)
	tempRandomList[0] = poseidonHashResult
	for i := 1; i < 48; i++ {
		tempRandomList[i] = new(big.Int)
		tempRandomList[i] = accumulator.UniversalHashToInt(tempRandomList[i-1])
	}
	var ranRem1, ranRem2, ranRem3, ranIns1, ranIns2, ranIns3 big.Int
	var leftShiftBits uint
	for i := 0; i < 8; i++ {
		leftShiftBits = 256 * (7 - 1)
		temp.Lsh(tempRandomList[i], leftShiftBits)
		ranRem1.Add(&ranRem1, &temp)
	}
	for i := 8; i < 16; i++ {
		leftShiftBits = 256 * (15 - 1)
		temp.Lsh(tempRandomList[i], leftShiftBits)
		ranRem2.Add(&ranRem2, &temp)
	}
	for i := 16; i < 24; i++ {
		leftShiftBits = 256 * (23 - 1)
		temp.Lsh(tempRandomList[i], leftShiftBits)
		ranRem3.Add(&ranRem3, &temp)
	}
	for i := 24; i < 32; i++ {
		leftShiftBits = 256 * (31 - 1)
		temp.Lsh(tempRandomList[i], leftShiftBits)
		ranIns1.Add(&ranIns1, &temp)
	}
	for i := 32; i < 40; i++ {
		leftShiftBits = 256 * (39 - 1)
		temp.Lsh(tempRandomList[i], leftShiftBits)
		ranIns2.Add(&ranIns2, &temp)
	}
	for i := 40; i < 48; i++ {
		leftShiftBits = 256 * (47 - 1)
		temp.Lsh(tempRandomList[i], leftShiftBits)
		ranIns3.Add(&ranIns3, &temp)
	}

	var original1, original2, original3 []*big.Int
	original1 = append(unchanged1, removed1...)
	original2 = append(unchanged2, removed2...)
	original3 = append(unchanged3, removed3...)
	originalProd1 := accumulator.SetProductRecursiveFast(original1)
	originalProd2 := accumulator.SetProductRecursiveFast(original2)
	originalProd3 := accumulator.SetProductRecursiveFast(original3)
	originalProd1 = bigfft.Mul(originalProd1, &ranRem1)
	originalProd2 = bigfft.Mul(originalProd2, &ranRem2)
	originalProd3 = bigfft.Mul(originalProd3, &ranRem3)

	setup := *accumulator.TrustedSetup()
	maxLen := setSize * 256 / bits.UintSize
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)

	// generate original zero-knowledge RSA accumulators
	accOri1 := multiexp.ExpParallel(setup.G, originalProd1, setup.N, table, 1, 0)
	accOri2 := multiexp.ExpParallel(setup.G, originalProd2, setup.N, table, 1, 0)
	accOri3 := multiexp.ExpParallel(setup.G, originalProd3, setup.N, table, 1, 0)

	fmt.Println("Precomputation and original RSA accumulators setup. Start to zero-knowledge MultiSwap")
	totalTime := time.Now().UTC()
	startingTime := time.Now().UTC()
	fmt.Println("Generate Acc_mid1,2,3")
	remProd1 := accumulator.SetProductRecursiveFast(removed1)
	remProd1 = bigfft.Mul(remProd1, &ranRem1)
	remProd2 := accumulator.SetProductRecursiveFast(removed2)
	remProd2 = bigfft.Mul(remProd2, &ranRem2)
	remProd3 := accumulator.SetProductRecursiveFast(removed3)
	remProd3 = bigfft.Mul(remProd3, &ranRem3)
	var accmidProd1, accmidProd2, accmidProd3 big.Int
	accmidProd1.Div(originalProd1, remProd1)
	accmidProd2.Div(originalProd2, remProd2)
	accmidProd3.Div(originalProd3, remProd3)

	accMid1 := multiexp.ExpParallel(setup.G, originalProd1, setup.N, table, 1, 0)
	accMid2 := multiexp.ExpParallel(setup.G, originalProd2, setup.N, table, 1, 0)
	accMid3 := multiexp.ExpParallel(setup.G, originalProd3, setup.N, table, 1, 0)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running Generate Acc_mid1,2,3 Takes [%.3f] Seconds \n", duration.Seconds())
	startingTime = time.Now().UTC()
	fmt.Println("Generate three zkPoKE")
	PoKE(accMid1, remProd1, accOri1, setup.N)
	PoKE(accMid2, remProd2, accOri2, setup.N)
	PoKE(accMid3, remProd3, accOri3, setup.N)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running Generate three zkPoKE Takes [%.3f] Seconds \n", duration.Seconds())
	fmt.Println("Generate Updated accumulators")
	startingTime = time.Now().UTC()
	insProd1 := accumulator.SetProductRecursiveFast(insert1)
	insProd1 = bigfft.Mul(insProd1, &ranIns1)
	insProd2 := accumulator.SetProductRecursiveFast(insert2)
	insProd2 = bigfft.Mul(insProd2, &ranIns2)
	insProd3 := accumulator.SetProductRecursiveFast(insert3)
	insProd3 = bigfft.Mul(insProd3, &ranIns3)
	accUpd1 := accumulator.AccumulateNew(accMid1, insProd1, setup.N)
	accUpd2 := accumulator.AccumulateNew(accMid2, insProd2, setup.N)
	accUpd3 := accumulator.AccumulateNew(accMid3, insProd3, setup.N)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running Generate Updated accumulators Takes [%.3f] Seconds \n", duration.Seconds())
	fmt.Println("Generate three zkPoKE")
	startingTime = time.Now().UTC()
	PoKE(accMid1, insProd1, accUpd1, setup.N)
	PoKE(accMid2, insProd2, accUpd2, setup.N)
	PoKE(accMid3, insProd3, accUpd3, setup.N)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running Generate three zkPoKE Takes [%.3f] Seconds \n", duration.Seconds())
	fmt.Println("Generate membership proofs for the three accumulators")
	startingTime = time.Now().UTC()

	newSet1 := append(unchanged1[:], insert1...)
	newSet2 := append(unchanged2[:], insert2...)
	newSet3 := append(unchanged3[:], insert3...)
	proofs1 := accumulator.ProveMembershipParallelWithTableWithRandomizer(setup.G, &ranIns1, setup.N, newSet1[:], 0, table)
	proofs2 := accumulator.ProveMembershipParallelWithTableWithRandomizer(setup.G, &ranIns2, setup.N, newSet2[:], 0, table)
	proofs3 := accumulator.ProveMembershipParallelWithTableWithRandomizer(setup.G, &ranIns3, setup.N, newSet3[:], 0, table)

	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running Generate membership proofs Takes [%.3f] Seconds \n", duration.Seconds())

	duration = time.Now().UTC().Sub(totalTime)
	fmt.Printf("Running full process Takes [%.3f] Seconds \n", duration.Seconds())
	startingTime = time.Now().UTC()
	func() {
		tempProof := proofs1[0]
		_ = tempProof.BitLen()
		tempProof = proofs2[0]
		_ = tempProof.BitLen()
		tempProof = proofs3[0]
		_ = tempProof.BitLen() // this line is simply used to allow accessing tempProof
	}()
}
