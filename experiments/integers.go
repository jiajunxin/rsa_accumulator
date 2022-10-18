package experiments

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

// TestFirstLayerPercentage tests the first layer of divide-and-conquer
func TestFirstLayerPercentage() {
	setSize := 10000
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	startingTime := time.Now().UTC()
	accumulator.ProveMembershipParallel(setup.G, setup.N, rep, 2)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallel Takes [%.3f] Seconds \n",
		duration.Seconds())
}

func TestMembership() {
	setSize := 1000000
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	startingTime := time.Now().UTC()
	prod := accumulator.SetProductRecursiveFast(rep)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running SetProductRecursiveFast Takes [%.3f] Seconds \n",
		duration.Seconds())

	startingTime = time.Now().UTC()
	accumulator.AccumulateNew(setup.G, prod, setup.N)
	//accumulator.ProveMembershipParallel(setup.G, setup.N, rep, 2)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccumulateNew Takes [%.3f] Seconds \n",
		duration.Seconds())
}

// TestProduct test different ways to
func TestProduct() {
	setSize := 1000000
	set := accumulator.GenBenchSet(setSize)
	var prod big.Int
	//rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	startingTime := time.Now().UTC()
	//prod = *accumulator.SetProduct2(rep)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallel Takes [%.3f] Seconds \n",
		duration.Seconds())
	fmt.Println("product length is", prod.BitLen())
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	startingTime = time.Now().UTC()
	prod = *accumulator.SetProductRecursive(rep)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallel Takes [%.3f] Seconds \n",
		duration.Seconds())
	fmt.Println("product length is", prod.BitLen())
}

// TestProduct test different ways to
func TestProduct2() {
	setSize := 1000000
	set := accumulator.GenBenchSet(setSize)
	var prod big.Int
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	startingTime := time.Now().UTC()
	prod = *accumulator.SetProductRecursive(rep)
	endingTime := time.Now().UTC()
	duration := endingTime.Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallel Takes [%.3f] Seconds \n",
		duration.Seconds())
	fmt.Println("product length is", prod.BitLen())
	var temp big.Int
	startingTime = time.Now().UTC()
	for i := 0; i < 100; i++ {
		temp.Div(&prod, rep[0])
	}
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Divide for 100 times Takes [%.3f] Seconds \n",
		duration.Seconds())
	fmt.Println("product length is", prod.BitLen())
}

// TestProduct test different ways to
func TestProduct3() {
	setSize := 10000
	set := accumulator.GenBenchSet(setSize)
	var prod big.Int
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	startingTime := time.Now().UTC()
	prod = *accumulator.SetProductRecursiveFast(rep)
	endingTime := time.Now().UTC()
	duration := endingTime.Sub(startingTime)
	fmt.Printf("Running SetProductRecursiveFast Takes [%.3f] Seconds \n",
		duration.Seconds())
	fmt.Println("product length is", prod.BitLen())
	startingTime = time.Now().UTC()
	prod = *accumulator.SetProductRecursive(rep)
	endingTime = time.Now().UTC()
	duration = endingTime.Sub(startingTime)
	fmt.Printf("Running SetProductRecursive Takes [%.3f] Seconds \n",
		duration.Seconds())
	fmt.Println("product length is", prod.BitLen())
}

func genDIMin(size int) []*big.Int {
	ret := make([]*big.Int, size)
	for i := 0; i < size; i++ {
		ret[i] = accumulator.Min2048
	}
	return ret
}

func genDIMax(size int) []*big.Int {
	ret := make([]*big.Int, size)
	var min257 big.Int
	min257.SetInt64(1)
	min257.Lsh(&min257, 256)
	for i := 0; i < size; i++ {
		ret[i] = new(big.Int)
		ret[i].Add(accumulator.Min2048, &min257)
	}
	//fmt.Println("2048 = ", accumulator.Min2048.String())
	// fmt.Println("257 = ", min257.String())
	// fmt.Println("set1[0] = ", ret[0].String())
	return ret
}

func TestRange() {
	setSize := 1000000

	var prodUpper, prodLower, difference big.Int

	set1 := genDIMax(setSize)
	set2 := genDIMin(setSize)

	prodUpper = *accumulator.SetProductRecursive(set1)
	prodLower = *accumulator.SetProductRecursive(set2)

	difference.Sub(&prodUpper, &prodLower)
	fmt.Println("Bit Length of the range is:", difference.BitLen())
	//fmt.Println("The range is:", difference.String())
	fmt.Println("Bit Length of the lower is:", prodLower.BitLen())
	//fmt.Println("The Lower is:", prodLower.String())
	fmt.Println("Bit Length of the upper is:", prodUpper.BitLen())
	//fmt.Println("The Upper is:", prodUpper.String())
}

func TestBreakDI(seed, bitLength, setSize int64) {
	rng := rand.New(rand.NewSource(seed))
	set := make([]*big.Int, setSize)

	upperBound := big.NewInt(1)
	upperBound = upperBound.Lsh(upperBound, uint(bitLength))
	fmt.Println("SetSize = ", setSize)
	fmt.Println("BitLength of each element = ", bitLength)
	for i := range set {
		set[i] = new(big.Int)
		set[i].Rand(rng, upperBound)
	}

	prod := accumulator.SetProductParallel(set, 2)

	fmt.Println("Bit length of the prod = ", prod.BitLen())
	var bigzero, bigone, counter, temp big.Int
	bigzero.SetInt64(0)
	bigone.SetInt64(1)
	counter.SetInt64(0)

	var ranElement, remainder *big.Int
	ranElement = new(big.Int)
	for {
		ranElement.Rand(rng, upperBound)
		remainder = temp.Mod(prod, ranElement)
		//fmt.Println("Bit length of the prof = ", prod.BitLen())
		//fmt.Println("Bit length of the ranEle = ", ranElement.BitLen())
		if remainder.Cmp(&bigzero) == 0 {
			fmt.Println("random element bit length ", ranElement.BitLen())
			//fmt.Println("remainder bit length ", remainder.BitLen())
			fmt.Println("The random seed is:", seed, ", the counter bit length =", counter.BitLen())
			return
		}
		counter.Add(&counter, &bigone)
	}
}

func TestDI() {
	// TestBreakDI(999, 100)
	// TestBreakDI(99, 100)
	// TestBreakDI(9, 100)
	// TestBreakDI(8, 100)

	// TestBreakDI(999, 128)
	// TestBreakDI(99, 128)
	// TestBreakDI(9, 128)
	// TestBreakDI(8, 128)
	// TestBreakDI(999, 130, 10000000)
	// TestBreakDI(99, 130, 10000000)
	// TestBreakDI(9, 130, 10000000)
	// TestBreakDI(8, 130, 10000000)

	// TestBreakDI(999, 140, 10000000)
	// TestBreakDI(99, 140, 10000000)
	// TestBreakDI(9, 140, 10000000)
	// TestBreakDI(8, 140, 10000000)

	// TestBreakDI(999, 150, 10000000)
	// TestBreakDI(99, 150, 10000000)
	// TestBreakDI(9, 150, 10000000)
	// TestBreakDI(8, 150, 10000000)

	// TestBreakDI(999, 160, 10000000)
	// TestBreakDI(99, 160, 10000000)
	// TestBreakDI(9, 160, 10000000)
	// TestBreakDI(8, 160, 10000000)

	// TestBreakDI(999, 192, 10000000)
	// TestBreakDI(99, 192, 10000000)
	// TestBreakDI(9, 192, 10000000)
	// TestBreakDI(8, 192, 10000000)

	TestBreakDI(999, 256, 100000000)
	TestBreakDI(99, 256, 100000000)
	TestBreakDI(9, 256, 100000000)
	TestBreakDI(8, 256, 100000000)
}
