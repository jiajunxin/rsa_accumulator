package zkmultiswap

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/poseidon"
	"github.com/consensys/gnark/std/math/bits"
)

// CircuitExp y == x**e mod N
// only the bitSize least significant bits of e are used
type CircuitExp struct {
	// tagging a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:",public"`
	Y frontend.Variable `gnark:",public"`
	N frontend.Variable `gnark:",public"`

	E frontend.Variable
}

// Define declares the circuit's constraints
// y == x**e
func (circuit *CircuitExp) Define(api frontend.API) error {

	// number of bits of exponent
	const bitSize = 255

	// specify constraints
	output := frontend.Variable(1)
	binary := bits.ToBinary(api, circuit.E, bits.WithNbDigits(bitSize))

	for i := 0; i < len(binary); i++ {
		if i != 0 {
			output = api.MulModP(output, output, circuit.N)
		}
		multiply := api.MulModP(output, circuit.X, circuit.N)
		output = api.Select(binary[len(binary)-1-i], multiply, output)
	}

	api.AssertIsEqual(circuit.Y, output)

	return nil
}

// CircuitRSAInside is the Zk-MultiSwap circuit for gnark.
// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The difference between CircuitRSAInside and Circuit is that:
// CircuitRSAInside checks two PoKE protocol inside SNARK
type CircuitRSAInside struct {
	// struct tag on a variable is optional
	// default uses variable name and secret visibility.
	ChallengeL1     frontend.Variable `gnark:",public"` // a prime challenge number L1
	ChallengeL2     frontend.Variable `gnark:",public"` // a prime challenge number L2
	RemainderR1     frontend.Variable `gnark:",public"` // a remainder R1
	RemainderR2     frontend.Variable `gnark:",public"` // a remainder R2
	CurrentEpochNum frontend.Variable `gnark:",public"` // current epoch number
	// Delta (2^1024) should be able to fixed as public parameters, however, gnark still cannot support big Int for now
	// we use the following two public input to replace the Delta
	// This because Delta + Hash(x) mod L = (Delta mod L + Hash(x) mod L) mod L
	DeltaModL1 frontend.Variable `gnark:",public"` // 2^1024 mod L1
	DeltaModL2 frontend.Variable `gnark:",public"` // 2^1024 mod L2
	Acc0       frontend.Variable `gnark:",public"` // the accumulator before update
	AccMid     frontend.Variable `gnark:",public"` // the intermediate accumulator
	Acc1       frontend.Variable `gnark:",public"` // the accumulator after update
	N          frontend.Variable `gnark:",public"` // the RSA modulus for the hidden order group
	//------------------------------private witness below--------------------------------------
	Randomizer1      frontend.Variable   // Used to randomize the removed set
	Randomizer2      frontend.Variable   // Used to randomize the inserted set
	OriginalSum      frontend.Variable   // original sum of balances for all users
	UpdatedSum       frontend.Variable   // updated sum of balances for all users
	UserID           []frontend.Variable // list of user IDs to be updated
	OriginalBalances []frontend.Variable // list of user balances before update
	OriginalHashes   []frontend.Variable // list of user hasher before update
	OriginalUpdEpoch []frontend.Variable // list of user updated epoch number before update
	UpdatedBalances  []frontend.Variable // list of user balances after update
	Q1               frontend.Variable   // PoKE proof for AccMid to Acc0, Q1^(ChallengeL1)AccMid^(RemainderR1)=Acc0
	Q2               frontend.Variable   // PoKE proof for AccMid to Acc1, Q2^(ChallengeL2)AccMid^(RemainderR2)=Acc1
}

// Define declares the circuit constraints
func (circuit CircuitRSAInside) Define(api frontend.API) error {
	api.ToBinary(circuit.Randomizer1, BitLength)
	api.ToBinary(circuit.Randomizer2, BitLength)
	api.AssertIsLess(circuit.DeltaModL1, circuit.ChallengeL1)
	api.AssertIsLess(circuit.DeltaModL2, circuit.ChallengeL2)

	api.AssertIsEqual(len(circuit.UserID), len(circuit.OriginalBalances))
	api.AssertIsEqual(len(circuit.UserID), len(circuit.OriginalHashes))
	api.AssertIsEqual(len(circuit.UserID), len(circuit.OriginalUpdEpoch))
	api.AssertIsEqual(len(circuit.UserID), len(circuit.UpdatedBalances))
	//check input are in the correct range
	api.AssertIsLess(circuit.RemainderR1, circuit.ChallengeL1)
	api.AssertIsLess(circuit.RemainderR2, circuit.ChallengeL2)
	// ToBinary not only returns the binary, but additionally checks if the binary representation is same as the input,
	// which means the input can be represented with the bit-length
	api.ToBinary(circuit.CurrentEpochNum, BitLength)
	api.ToBinary(circuit.OriginalSum, BitLength)
	api.ToBinary(circuit.UpdatedSum, BitLength)

	// checking PoKE proof for AccMid to Acc0, Q1^(ChallengeL1)AccMid^(RemainderR1)=Acc0
	// and PoKE proof for AccMid to Acc1, Q2^(ChallengeL2)AccMid^(RemainderR2)=Acc1
	api.AssertIsLess(circuit.Q1, circuit.N)
	api.AssertIsLess(circuit.Q2, circuit.N)
	api.AssertIsLess(circuit.Acc0, circuit.N)
	api.AssertIsLess(circuit.AccMid, circuit.N)
	api.AssertIsLess(circuit.Acc1, circuit.N)

	const challengeBitLen = 255

	Q1L1 := frontend.Variable(1)
	L1Binary := bits.ToBinary(api, circuit.ChallengeL1, bits.WithNbDigits(challengeBitLen))

	for i := 0; i < len(L1Binary); i++ {
		if i != 0 {
			Q1L1 = api.MulModP(Q1L1, Q1L1, circuit.N)
		}
		multiplier := api.MulModP(Q1L1, circuit.Q1, circuit.N)
		Q1L1 = api.Select(L1Binary[len(L1Binary)-1-i], multiplier, Q1L1)
	}
	api.AssertIsEqual(Q1L1, circuit.Acc0)

	// end of checking the PoKE proof

	// check we do not have repeating IDs and IDs in correct range
	for i := 0; i < len(circuit.UserID)-1; i++ {
		api.AssertIsLess(circuit.UserID[i], circuit.UserID[i+1])
	}
	//api.ToBinary(circuit.UserID[len(circuit.UserID)-1], BitLength)

	for i := 0; i < len(circuit.UserID); i++ {
		api.ToBinary(circuit.OriginalBalances[i], BitLength)
		api.AssertIsLess(circuit.OriginalUpdEpoch[i], circuit.CurrentEpochNum)
		api.ToBinary(circuit.UpdatedBalances[i], BitLength)
	}

	var remainder1, remainder2 frontend.Variable = 1, 1
	tempSum := circuit.OriginalSum
	tempSum = api.Sub(tempSum, circuit.OriginalBalances[0])
	tempSum = api.Add(tempSum, circuit.UpdatedBalances[0])
	for i := 0; i < len(circuit.UserID); i++ {
		tempHash0 := poseidon.Poseidon(api, circuit.UserID[i], circuit.OriginalBalances[i], circuit.OriginalUpdEpoch[i], circuit.OriginalHashes[i])
		//api.Println(tempHash0)
		tempHash1 := api.Add(tempHash0, circuit.DeltaModL1)
		remainder1 = api.MulModP(remainder1, tempHash1, circuit.ChallengeL1)

		// Check HashChain
		tempHash2 := poseidon.Poseidon(api, circuit.UserID[i], circuit.UpdatedBalances[i], circuit.CurrentEpochNum, tempHash0)
		tempHash2 = api.Add(tempHash2, circuit.DeltaModL2)
		remainder2 = api.MulModP(remainder2, tempHash2, circuit.ChallengeL2)

		tempSum = api.Sub(tempSum, circuit.OriginalBalances[i])
		tempSum = api.Add(tempSum, circuit.UpdatedBalances[i])
	}
	// because gnark cannot support 2048-bits large integers, we are using the product of 8 255-bits random numbers to replace one large RSA-domain randomizer.
	for i := 0; i < 8; i++ {
		tempHash := poseidon.Poseidon(api, circuit.Randomizer1, i)
		remainder1 = api.MulModP(remainder1, tempHash, circuit.ChallengeL1)
		//api.Println(tempHash)
		tempHash = poseidon.Poseidon(api, circuit.Randomizer2, i)
		remainder2 = api.MulModP(remainder2, tempHash, circuit.ChallengeL2)
	}
	api.AssertIsEqual(remainder1, circuit.RemainderR1)
	api.AssertIsEqual(remainder2, circuit.RemainderR2)
	api.AssertIsEqual(tempSum, circuit.UpdatedSum)

	return nil
}

// InitCircuitRSAWithSize init a circuit with challenges, OriginalHashes and CurrentEpochNum value 1, all other values 0. Use for test purpose only.
func InitCircuitRSAWithSize(size uint32) *Circuit {
	var circuit Circuit
	circuit.ChallengeL1 = 1
	circuit.ChallengeL2 = 1
	circuit.RemainderR1 = 0
	circuit.RemainderR2 = 0
	circuit.CurrentEpochNum = 1
	circuit.DeltaModL1 = 0
	circuit.DeltaModL2 = 0
	circuit.OriginalSum = 1
	circuit.UpdatedSum = 1
	circuit.Randomizer1 = 1
	circuit.Randomizer2 = 1

	circuit.UserID = make([]frontend.Variable, size)
	circuit.OriginalBalances = make([]frontend.Variable, size)
	circuit.OriginalHashes = make([]frontend.Variable, size)
	circuit.OriginalUpdEpoch = make([]frontend.Variable, size)
	circuit.UpdatedBalances = make([]frontend.Variable, size)
	for i := uint32(0); i < size; i++ {
		circuit.UserID[i] = i
		circuit.OriginalBalances[i] = 0
		circuit.OriginalHashes[i] = 1
		circuit.OriginalUpdEpoch[i] = 0
		circuit.UpdatedBalances[i] = 0
	}
	return &circuit
}

// AssignCircuitRSA assign a circuit with UpdateSet32 values.
func AssignCircuitRSA(input *UpdateSet32) *Circuit {
	if !input.IsValid() {
		panic("error in InitCircuit, the input set is invalid")
	}
	var circuit Circuit
	size := len(input.OriginalBalances)
	circuit.ChallengeL1 = input.ChallengeL1
	circuit.ChallengeL2 = input.ChallengeL2
	circuit.RemainderR1 = input.RemainderR1
	circuit.RemainderR2 = input.RemainderR2
	circuit.CurrentEpochNum = input.CurrentEpochNum
	circuit.DeltaModL1 = input.DeltaModL1
	circuit.DeltaModL2 = input.DeltaModL2
	circuit.OriginalSum = input.OriginalSum
	circuit.UpdatedSum = input.UpdatedSum
	circuit.Randomizer1 = input.Randomizer1
	circuit.Randomizer2 = input.Randomizer2

	circuit.UserID = make([]frontend.Variable, size)
	circuit.OriginalBalances = make([]frontend.Variable, size)
	circuit.OriginalHashes = make([]frontend.Variable, size)
	circuit.OriginalUpdEpoch = make([]frontend.Variable, size)
	circuit.UpdatedBalances = make([]frontend.Variable, size)
	for i := 0; i < size; i++ {
		circuit.UserID[i] = input.UserID[i]
		circuit.OriginalBalances[i] = input.OriginalBalances[i]
		circuit.OriginalHashes[i] = input.OriginalHashes[i]
		circuit.OriginalUpdEpoch[i] = input.OriginalUpdEpoch[i]
		circuit.UpdatedBalances[i] = input.UpdatedBalances[i]
	}
	return &circuit
}

// AssignCircuitRSAHelper assign a circuit with PublicInfo values.
func AssignCircuitRSAHelper(input *PublicInfo) *Circuit {
	circuit := InitCircuitWithSize(1)
	circuit.ChallengeL1 = input.ChallengeL1
	circuit.ChallengeL2 = input.ChallengeL2
	circuit.RemainderR1 = input.RemainderR1
	circuit.RemainderR2 = input.RemainderR2
	circuit.CurrentEpochNum = input.CurrentEpochNum
	circuit.DeltaModL1 = input.DeltaModL1
	circuit.DeltaModL2 = input.DeltaModL2

	return circuit
}
