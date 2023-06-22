package zkmultiswap

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/poseidon"
)

// Circuit is the Zk-MultiSwap circuit for gnark.
// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
type Circuit struct {
	// struct tag on a variable is optional
	// default uses variable name and secret visibility.
	ChallengeL1     frontend.Variable `gnark:",public"` // a prime challenge number L1
	ChallengeL2     frontend.Variable `gnark:",public"` // a prime challenge number L2
	RemainderR1     frontend.Variable `gnark:",public"` // a remainder R1
	RemainderR2     frontend.Variable `gnark:",public"` // a remainder R2
	CurrentEpochNum frontend.Variable `gnark:",public"` // current epoch number
	//------------------------------private witness below--------------------------------------
	Randomizer       frontend.Variable   // Used to randomize the witness for commit-and-prove, reserved for future
	OriginalSum      frontend.Variable   // original sum of balances for all users
	UpdatedSum       frontend.Variable   // updated sum of balances for all users
	UserID           []frontend.Variable // list of user IDs to be updated
	OriginalBalances []frontend.Variable // list of user balances before update
	OriginalHashes   []frontend.Variable // list of user hasher before update
	OriginalUpdEpoch []frontend.Variable // list of user updated epoch number before update
	UpdatedBalances  []frontend.Variable // list of user balances after update
}

// Define declares the circuit constraints
func (circuit Circuit) Define(api frontend.API) error {
	// To be fixed!
	api.ToBinary(circuit.Randomizer, BitLength)

	api.AssertIsEqual(len(circuit.UserID), len(circuit.OriginalBalances))
	api.AssertIsEqual(len(circuit.UserID), len(circuit.OriginalHashes))
	api.AssertIsEqual(len(circuit.UserID), len(circuit.OriginalUpdEpoch))
	api.AssertIsEqual(len(circuit.UserID), len(circuit.UpdatedBalances))
	//check input are in the correct range
	api.AssertIsLess(circuit.RemainderR1, circuit.ChallengeL1)
	api.AssertIsLess(circuit.RemainderR2, circuit.ChallengeL2)
	// ToBinary not only returns the binary, but additionaly checks if the binary representation is same as the input,
	// which means the input can be represented with the bit-length
	api.ToBinary(circuit.CurrentEpochNum, BitLength)
	api.ToBinary(circuit.OriginalSum, BitLength)
	api.ToBinary(circuit.UpdatedSum, BitLength)

	// check we do not have repeating IDs and IDs in correct range
	for i := 0; i < len(circuit.UserID)-1; i++ {
		api.AssertIsLess(circuit.UserID[i], circuit.UserID[i+1])
	}
	//api.ToBinary(circuit.UserID[len(circuit.UserID)-1], BitLength)

	for i := 0; i < len(circuit.UserID); i++ {
		api.ToBinary(circuit.OriginalBalances[i], BitLength)
		//api.AssertIsLess(circuit.OriginalHashes[i], api.Curve().Info().Fp.Modulus)
		api.AssertIsLess(circuit.OriginalUpdEpoch[i], circuit.CurrentEpochNum)
		api.ToBinary(circuit.UpdatedBalances[i], BitLength)
	}

	remainder1 := poseidon.Poseidon(api, circuit.UserID[0], circuit.OriginalBalances[0], circuit.OriginalUpdEpoch[0], circuit.OriginalHashes[0])
	//	api.Println(remainder1)
	remainder2 := poseidon.Poseidon(api, circuit.UserID[0], circuit.UpdatedBalances[0], circuit.CurrentEpochNum, remainder1)
	tempSum := circuit.OriginalSum
	api.Sub(tempSum, circuit.OriginalBalances[0])
	api.Add(tempSum, circuit.UpdatedBalances[0])
	for i := 1; i < len(circuit.UserID); i++ {
		tempHash := poseidon.Poseidon(api, circuit.UserID[i], circuit.OriginalBalances[i], circuit.OriginalUpdEpoch[i], circuit.OriginalHashes[i])
		remainder1 = api.MulModP(remainder1, tempHash, circuit.ChallengeL1)

		tempHash2 := poseidon.Poseidon(api, circuit.UserID[i], circuit.UpdatedBalances[i], circuit.CurrentEpochNum, tempHash)
		remainder2 = api.MulModP(remainder2, tempHash2, circuit.ChallengeL2)

		api.Sub(tempSum, circuit.OriginalBalances[i])
		api.Add(tempSum, circuit.UpdatedBalances[i])
	}
	api.AssertIsEqual(remainder1, circuit.RemainderR1)
	api.AssertIsEqual(remainder2, circuit.RemainderR2)
	api.AssertIsEqual(tempSum, circuit.UpdatedSum)

	return nil
}

// InitCircuitWithSize init a circuit with challenges, OriginalHashes and CurrentEpochNum value 1, all other values 0. Use for test purpose only.
func InitCircuitWithSize(size uint32) *Circuit {
	var circuit Circuit
	circuit.ChallengeL1 = 1
	circuit.ChallengeL2 = 1
	circuit.RemainderR1 = 0
	circuit.RemainderR2 = 0
	circuit.CurrentEpochNum = 1
	circuit.OriginalSum = 1
	circuit.UpdatedSum = 1
	circuit.Randomizer = 1

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

// AssignCircuit assign a circuit with UpdateSet32 values.
func AssignCircuit(input *UpdateSet32) *Circuit {
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
	circuit.OriginalSum = input.OriginalSum
	circuit.UpdatedSum = input.UpdatedSum
	circuit.Randomizer = input.Randomizer

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

// AssignCircuitHelper assign a circuit with PublicInfo values.
func AssignCircuitHelper(input *PublicInfo) *Circuit {
	circuit := InitCircuitWithSize(1)
	circuit.ChallengeL1 = input.ChallengeL1
	circuit.ChallengeL2 = input.ChallengeL2
	circuit.RemainderR1 = input.RemainderR1
	circuit.RemainderR2 = input.RemainderR2
	circuit.CurrentEpochNum = input.CurrentEpochNum

	return circuit
}
