package zkmultiswap

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/jiajunxin/rsa_accumulator/accumulator"
	fiatshamir "github.com/jiajunxin/rsa_accumulator/fiat-shamir"
)

const (
	// BitLength is the bit length of the user ID, balnace and epoch number. It can be 32, 64 or any valid number within the field
	BitLength = 32
	// CurrentEpochNum is used for *test purpose* only. It should be larger than the test set size and all OriginalUpdEpoch
	CurrentEpochNum = 1000000
	// OriginalSum is used for *test purpose* only. It should be larger than 0 and the updated balance should also be positive
	OriginalSum = 10000

	// KeyPathPrefix denotes the path to store the circuit and keys. fileName = KeyPathPrefix + "_" + strconv.FormatInt(int64(size), 10) + different names
	KeyPathPrefix = "zkmultiswap"
)

// UpdateSet32 is one set for the prover with uint32 for CurrentEpochNum,
type UpdateSet32 struct {
	ChallengeL1      big.Int
	ChallengeL2      big.Int
	RemainderR1      big.Int
	RemainderR2      big.Int
	CurrentEpochNum  uint32
	DeltaModL1       big.Int
	DeltaModL2       big.Int
	Randomizer1      big.Int
	Randomizer2      big.Int
	OriginalSum      uint32
	UpdatedSum       uint32
	UserID           []uint32
	OriginalBalances []uint32
	OriginalHashes   []big.Int
	OriginalUpdEpoch []uint32
	UpdatedBalances  []uint32
}

// PublicInfo is the public information part of UpdateSet32
type PublicInfo struct {
	ChallengeL1     big.Int
	ChallengeL2     big.Int
	RemainderR1     big.Int
	RemainderR2     big.Int
	CurrentEpochNum uint32
	DeltaModL1      big.Int
	DeltaModL2      big.Int
}

// IsValid returns true only if the input is valid for multiSwap
func (input *UpdateSet32) IsValid() bool {
	if len(input.UserID) < 2 {
		return false
	}
	if len(input.UserID) != len(input.OriginalBalances) {
		return false
	}
	if len(input.UserID) != len(input.OriginalHashes) {
		return false
	}
	if len(input.UserID) != len(input.OriginalUpdEpoch) {
		return false
	}
	if len(input.UserID) != len(input.UpdatedBalances) {
		return false
	}
	return true
}

func getRandomAcc(setup *accumulator.Setup) *big.Int {
	var ret big.Int
	rand := accumulator.GenRandomizer()
	ret.Exp(setup.G, rand, setup.N)
	return &ret
}

// SetupTranscript should takes in all public information regarding the MultiSwap
func SetupTranscript(setup *accumulator.Setup, accOld, accMid, accNew *big.Int, CurrentEpochNum uint32) *fiatshamir.Transcript {
	transcript := fiatshamir.InitTranscript([]string{setup.G.String(), setup.N.String()}, fiatshamir.Max252)
	transcript.Append(strconv.Itoa(int(CurrentEpochNum)))
	return transcript
}

// GenTestSet generates a set of values for test purpose.
func GenTestSet(setsize uint32, setup *accumulator.Setup) *UpdateSet32 {
	var ret UpdateSet32
	ret.UserID = make([]uint32, setsize)
	ret.OriginalBalances = make([]uint32, setsize)
	ret.OriginalUpdEpoch = make([]uint32, setsize)
	ret.OriginalHashes = make([]big.Int, setsize)
	ret.UpdatedBalances = make([]uint32, setsize)

	ret.CurrentEpochNum = CurrentEpochNum
	for i := uint32(0); i < setsize; i++ {
		j := i*2 + 1      // no special meaning for j, just need some non-repeating positive integers
		ret.UserID[i] = j // we need to arrange user IDs in accending order for checking them efficiently
		ret.OriginalBalances[i] = j
		ret.OriginalUpdEpoch[i] = 10
		ret.OriginalHashes[i].SetInt64(int64(j))
		ret.UpdatedBalances[i] = j
	}
	ret.OriginalSum = OriginalSum
	ret.UpdatedSum = OriginalSum // UpdatedSum can be any valid positive numbers, but we are testing the case UpdatedSum = OriginalSum for simplicity

	// get slice of elements removed and inserted
	removeSet := make([]*big.Int, setsize)
	insertSet := make([]*big.Int, setsize)

	var poseidonhash *fr.Element // this is the Poseidon part of the DI hash. We use this to build the hash chain. The original DI hash is to long to directly input into Poseidon hash
	for i := uint32(0); i < setsize; i++ {
		poseidonhash, removeSet[i] = accumulator.PoseidonAndDIHash(accumulator.ElementFromUint32(ret.UserID[i]), accumulator.ElementFromUint32(ret.OriginalBalances[i]),
			accumulator.ElementFromUint32(ret.OriginalUpdEpoch[i]), accumulator.ElementFromString(ret.OriginalHashes[i].String()))
		//fmt.Println("poseidonhash i = ", poseidonhash.String())

		insertSet[i] = accumulator.DIHashPoseidon(accumulator.ElementFromUint32(ret.UserID[i]), accumulator.ElementFromUint32(ret.UpdatedBalances[i]),
			accumulator.ElementFromUint32(ret.CurrentEpochNum), poseidonhash)
	}
	prod1 := accumulator.SetProductRecursiveFast(removeSet)
	prod2 := accumulator.SetProductRecursiveFast(insertSet)

	// Randomizers are FIXED!!! for test purpose
	ret.Randomizer1 = *big.NewInt(200)
	ret.Randomizer2 = *big.NewInt(300)
	var tempInt big.Int
	// because gnark cannot support 2048-bits large integers, we are using the product of 8 255-bits random numbers to replace one large RSA-domain randomizer.
	for i := 0; i < 8; i++ {
		tempHash := poseidon.Poseidon(accumulator.ElementFromBigInt(&ret.Randomizer1), accumulator.ElementFromUint32(uint32(i)))
		tempHash.ToBigIntRegular(&tempInt)
		prod1.Mul(prod1, &tempInt)

		tempHash = poseidon.Poseidon(accumulator.ElementFromBigInt(&ret.Randomizer2), accumulator.ElementFromUint32(uint32(i)))
		tempHash.ToBigIntRegular(&tempInt)
		prod2.Mul(prod2, &tempInt)
	}

	// get accumulators
	accMid := getRandomAcc(setup)
	var accOld, accNew big.Int
	accOld.Exp(accMid, prod1, setup.N)
	accNew.Exp(accMid, prod2, setup.N)

	// get challenge
	transcript := SetupTranscript(setup, &accOld, accMid, &accNew, ret.CurrentEpochNum)
	challengeL1 := transcript.GetChallengeAndAppendTranscript()
	challengeL2 := transcript.GetChallengeAndAppendTranscript()

	// get remainder
	remainderR1 := big.NewInt(1)
	remainderR2 := big.NewInt(1)
	remainderR1.Mod(prod1, challengeL1)
	remainderR2.Mod(prod2, challengeL2)

	ret.ChallengeL1 = *challengeL1
	ret.ChallengeL2 = *challengeL2
	ret.RemainderR1 = *remainderR1
	ret.RemainderR2 = *remainderR2
	var deltaModL1, deltaModL2 big.Int
	deltaModL1.Mod(accumulator.Min1024, challengeL1)
	deltaModL2.Mod(accumulator.Min1024, challengeL2)
	ret.DeltaModL1 = deltaModL1
	ret.DeltaModL2 = deltaModL2

	if !ret.IsValid() {
		panic("error in GenTestSet, the generated test set is invalid")
	}
	return &ret
}

// PublicPart returns a new UpdateSet32 with same public part and hidden part 0
func (input *UpdateSet32) PublicPart() *PublicInfo {
	var ret PublicInfo
	ret.ChallengeL1 = input.ChallengeL1
	ret.ChallengeL2 = input.ChallengeL2
	ret.RemainderR1 = input.RemainderR1
	ret.RemainderR2 = input.RemainderR2
	ret.CurrentEpochNum = input.CurrentEpochNum
	ret.DeltaModL1 = input.DeltaModL1
	ret.DeltaModL2 = input.DeltaModL2
	return &ret
}

func isCircuitExist(testSetSize uint32) bool {
	fileName := KeyPathPrefix + "_" + strconv.FormatInt(int64(testSetSize), 10) + ".ccs.save"
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

// TestMultiSwap is temporarily used for test purpose
func TestMultiSwap(testSetSize uint32) {
	if !isCircuitExist(testSetSize) {
		fmt.Println("Circuit haven't been compiled for testSetSize = ", testSetSize, ". Start compiling.")
		startingTime := time.Now().UTC()
		SetupZkMultiswap(testSetSize)
		duration := time.Now().UTC().Sub(startingTime)
		fmt.Printf("Generating a SNARK circuit for set size = %d, takes [%.3f] Seconds \n", testSetSize, duration.Seconds())
		runtime.GC()
	} else {
		fmt.Println("Circuit have already been compiled for test purpose.")
	}
	testSet := GenTestSet(testSetSize, accumulator.TrustedSetup())
	publicInfo := testSet.PublicPart()
	proof, err := Prove(testSet)
	if err != nil {
		fmt.Println("Error during Prove")
		panic(err)
	}
	runtime.GC()

	flag := Verify(proof, testSetSize, publicInfo)
	if flag {
		fmt.Println("Verification passed")
		return
	}
	fmt.Println("Verification failed")
}

// TestMultiSwapAndOutputSmartContract outputs a Solidity smart contract to verify the SNARK
func TestMultiSwapAndOutputSmartContract(testSetSize uint32) {
	if !isCircuitExist(testSetSize) {
		fmt.Println("Circuit haven't been compiled for testSetSize = ", testSetSize, ". Start compiling.")
		startingTime := time.Now().UTC()
		SetupZkMultiswap(testSetSize)
		duration := time.Now().UTC().Sub(startingTime)
		fmt.Printf("Generating a SNARK circuit for set size = %d, takes [%.3f] Seconds \n", testSetSize, duration.Seconds())
		runtime.GC()
	} else {
		fmt.Println("Circuit have already been compiled for test purpose.")
	}
	fileName := KeyPathPrefix + "_" + strconv.FormatInt(int64(testSetSize), 10)
	vk, err := LoadVerifyingKey(fileName)
	if err != nil {
		panic(err)
	}
	f, err := os.Create("contract_g16.sol")
	if err != nil {
		panic(err)
	}
	err = vk.ExportSolidity(f)
	if err != nil {
		panic(err)
	}
}

func TestMultiSwapAndOutputSmartContract2(testSetSize uint32) error {
	var circuit Circuit
	circuit = *InitCircuitWithSize(testSetSize)
	r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		return err
	}
	SetupZkMultiswap(testSetSize)
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return err
	}
	{
		f, err := os.Create("Notus.g16.vk")
		if err != nil {
			return err
		}
		_, err = vk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}
	{
		f, err := os.Create("Notus.g16.pk")
		if err != nil {
			return err
		}
		_, err = pk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}
	{
		f, err := os.Create("Notuscontract_g16.sol")
		if err != nil {
			return err
		}
		err = vk.ExportSolidity(f)
		if err != nil {
			return err
		}
	}

	testSet := GenTestSet(testSetSize, accumulator.TrustedSetup())
	publicInfo := testSet.PublicPart()
	proof, err := Prove(testSet)
	if err != nil {
		fmt.Println("Error during Prove, error = ", err.Error())
		panic(err)
	}

	flag := Verify(proof, testSetSize, publicInfo)
	if !flag {
		fmt.Println("Verification failed")
		return nil
	}
	fmt.Println("Verification passed")

	// get proof bytes
	const fpSize = 4 * 8
	var buf bytes.Buffer
	(*proof).WriteRawTo(&buf)
	proofBytes := buf.Bytes()
	// solidity contract inputs
	var (
		a     [2]*big.Int
		b     [2][2]*big.Int
		c     [2]*big.Int
		input [7]*big.Int
	)
	for i := 0; i < 7; i++ {
		input[i] = new(big.Int)
	}

	// proof.Ar, proof.Bs, proof.Krs
	a[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	a[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	b[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	b[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	b[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	b[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	c[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	c[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])
	input[0] = &publicInfo.ChallengeL1
	input[1] = &publicInfo.ChallengeL2
	input[2].SetInt64(int64(publicInfo.CurrentEpochNum))
	input[3] = &publicInfo.DeltaModL1
	input[4] = &publicInfo.DeltaModL2
	input[5] = &publicInfo.RemainderR1
	input[6] = &publicInfo.RemainderR2

	fmt.Println("a[0] = ", a[0].String())
	fmt.Println("a[1] = ", a[1].String())
	fmt.Println("b[0][0] = ", b[0][0].String())
	fmt.Println("b[0][1] = ", b[0][1].String())
	fmt.Println("b[1][0] = ", b[1][0].String())
	fmt.Println("b[1][1] = ", b[1][1].String())
	fmt.Println("c[0] = ", c[0].String())
	fmt.Println("c[1] = ", c[1].String())
	for i := 0; i < 7; i++ {
		fmt.Println("input[", i, "] = ", input[i].String())
	}
	return nil
}
