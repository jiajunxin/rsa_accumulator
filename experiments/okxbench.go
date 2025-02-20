package experiments

import (
	"fmt"
	"math/big"
	"math/bits"
	"runtime"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/jiajunxin/multiexp"
	"github.com/remyoudompheng/bigfft"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
	"github.com/jiajunxin/rsa_accumulator/zkmultiswap"
)

// OKXBenchParallel benchmark OKX circuit with at most 32 cores
func OKXBenchParallel(setsize, updatedSetSize uint32) {
	if !isCircuitExist(updatedSetSize) {
		fmt.Println("Circuit haven't been compiled for testSetSize = ", updatedSetSize, ". Start compiling.")
		startingTime := time.Now().UTC()
		zkmultiswap.SetupZkMultiswap(updatedSetSize)
		duration := time.Now().UTC().Sub(startingTime)
		fmt.Printf("Generating a SNARK circuit for set size = %d, takes [%.3f] Seconds \n", updatedSetSize, duration.Seconds())
		runtime.GC()
	} else {
		fmt.Println("Circuit have already been compiled for test purpose.")
	}
	//--------------------------------------------generating accumulator--------------------------------
	var ret zkmultiswap.UpdateSet32
	ret.UserID = make([]uint32, updatedSetSize)
	ret.OriginalBalances = make([]uint32, updatedSetSize)
	ret.OriginalUpdEpoch = make([]uint32, updatedSetSize)
	ret.OriginalHashes = make([]big.Int, updatedSetSize)
	ret.UpdatedBalances = make([]uint32, updatedSetSize)

	ret.CurrentEpochNum = zkmultiswap.CurrentEpochNum
	for i := uint32(0); i < updatedSetSize; i++ {
		j := i*2 + 1      // no special meaning for j, just need some non-repeating positive integers
		ret.UserID[i] = j // we need to arrange user IDs in accending order for checking them efficiently
		ret.OriginalBalances[i] = j
		ret.OriginalUpdEpoch[i] = 10
		ret.OriginalHashes[i].SetInt64(int64(j))
		ret.UpdatedBalances[i] = j
	}
	ret.OriginalSum = zkmultiswap.OriginalSum
	ret.UpdatedSum = zkmultiswap.OriginalSum // UpdatedSum can be any valid positive numbers, but we are testing the case UpdatedSum = OriginalSum for simplicity

	// get slice of elements removed and inserted
	removeSet := make([]*big.Int, updatedSetSize)
	insertSet := make([]*big.Int, updatedSetSize)
	startingTime := time.Now().UTC()
	for i := uint32(0); i < updatedSetSize; i++ {
		var poseidonhash *fr.Element // this is the Poseidon part of the DI hash. We use this to build the hash chain. The original DI hash is to long to directly input into Poseidon hash
		poseidonhash, removeSet[i] = accumulator.PoseidonAndDIHash(accumulator.ElementFromUint32(ret.UserID[i]), accumulator.ElementFromUint32(ret.OriginalBalances[i]),
			accumulator.ElementFromUint32(ret.OriginalUpdEpoch[i]), accumulator.ElementFromString(ret.OriginalHashes[i].String()))
		//fmt.Println("poseidonhash i = ", poseidonhash.String())

		insertSet[i] = accumulator.DIHashPoseidon(accumulator.ElementFromUint32(ret.UserID[i]), accumulator.ElementFromUint32(ret.UpdatedBalances[i]),
			accumulator.ElementFromUint32(ret.CurrentEpochNum), poseidonhash)
	}

	// Randomizers are FIXED!!! for test purpose
	ret.Randomizer1 = *big.NewInt(200)
	ret.Randomizer2 = *big.NewInt(300)
	var removedRanProd, insertedRanProd big.Int
	removedRanProd.SetInt64(1)
	insertedRanProd.SetInt64(1)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generate DI representatives Takes [%.3f] Seconds \n", duration.Seconds())

	// get accumulators
	setup := *accumulator.TrustedSetup()
	maxLen := setsize * 1025 / bits.UintSize
	fmt.Println("Generating precomputing tables")
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, int(maxLen))

	unchangedSet := accumulator.GenBenchSet(int(setsize - updatedSetSize))
	unchanged := accumulator.GenRepresentatives(unchangedSet, accumulator.DIHashFromPoseidon)

	original := append(unchanged, removeSet...)
	originalProd := accumulator.SetProductRecursiveFast(original)
	originalProd = bigfft.Mul(originalProd, &removedRanProd)
	runtime.GC()
	startingTime = time.Now().UTC()
	_, err := multiexp.ExpParallelAuto(setup.G, originalProd, setup.N, table)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generate subset accumulator Takes [%.4f] Seconds \n", duration.Seconds())
	if err != nil {
		panic(err)
	}
	runtime.GC()
	startingTime = time.Now().UTC()
	_, err = multiexp.ExpParallelAuto(setup.G, originalProd, setup.N, table)
	// accOri := multiexp.ExpParallel(setup.G, originalProd, setup.N, table, limit, 0)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generate subset accumulator Takes [%.4f] Seconds \n", duration.Seconds())
	if err != nil {
		panic(err)
	}
	// fmt.Println("accOri = ", accOri.String())
	//fmt.Println("accMid = ", accMid.String())
	//--------------------------------------------finish generating accumulator--------------------------------
	table = nil
	runtime.GC()
}
