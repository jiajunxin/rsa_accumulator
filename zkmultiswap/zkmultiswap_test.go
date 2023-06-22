package zkmultiswap

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/test"
	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

// compare return 0 if input1 == input2
func compare(input1 *big.Int, input2 []byte) int {
	input1bytes := input1.Bytes()
	return bytes.Compare(input1bytes, input2)
}

// func TestPoseidonHash(t *testing.T) {
// 	inputs := "anything"
// 	result1, err := iden3Poseidon.HashBytes([]byte(inputs))
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}

// 	poseidonHasher := bnPoseidon.NewPoseidon()
// 	_, err = poseidonHasher.Write([]byte(inputs))
// 	if err != nil {
// 		t.Errorf(err.Error())
// 	}
// 	result2 := poseidonHasher.Sum(nil)

// 	if compare(result1, result2) != 0 {
// 		fmt.Println("result1 = ", result1.String())
// 		fmt.Println("result2 = ", result2)
// 	}
// }

func TestPublicWitness(t *testing.T) {
	testSetSize := uint32(10)
	testSet := GenTestSet(testSetSize, accumulator.TrustedSetup())

	assignment := AssignCircuit(testSet)
	witness, err := frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		t.Errorf(err.Error())
	}
	publicWitness, err := witness.Public()
	if err != nil {
		fmt.Println("error while generating public witness")
		t.Errorf(err.Error())
	}

	publicPart := testSet.PublicPart()
	assignment2 := AssignCircuitHelper(publicPart)
	publicWitness2, err := frontend.NewWitness(assignment2, ecc.BN254, frontend.PublicOnly())
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(publicWitness.Vector, publicWitness2.Vector) {
		t.Errorf("public witness and public witness build from public info are not equal")
	}
	// test cases that should not be equal
	assignment.CurrentEpochNum = 666
	witness, err = frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		t.Errorf(err.Error())
	}
	result2 := poseidonHasher.Sum(nil)

// 	if compare(result1, result2) != 0 {
// 		fmt.Println("result1 = ", result1.String())
// 		fmt.Println("result2 = ", result2)
// 	}
// }

func TestPublicWitness(t *testing.T) {
	testSetSize := uint32(10)
	testSet := GenTestSet(testSetSize, accumulator.TrustedSetup())

	assignment := AssignCircuit(testSet)
	witness, err := frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		t.Errorf(err.Error())
	}
	publicWitness, err := witness.Public()
	if err != nil {
		fmt.Println("error while generating public witness")
		t.Errorf(err.Error())
	}

	publicPart := testSet.PublicPart()
	assignment2 := AssignCircuitHelper(publicPart)
	publicWitness2, err := frontend.NewWitness(assignment2, ecc.BN254, frontend.PublicOnly())
	if err != nil {
		t.Errorf(err.Error())
	}
	if !reflect.DeepEqual(publicWitness.Vector, publicWitness2.Vector) {
		t.Errorf("public witness and public witness build from public info are not equal")
	}
	// test cases that should not be equal
	assignment.CurrentEpochNum = 666
	witness, err = frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		t.Errorf(err.Error())
	}
	publicWitness, err = witness.Public()
	if err != nil {
		fmt.Println("error while generating public witness")
		t.Errorf(err.Error())
	}
	if reflect.DeepEqual(publicWitness.Vector, publicWitness2.Vector) {
		t.Errorf("public witness and public witness should not be equal")
	}
}

func TestZkMultiSwap(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit, witness Circuit
	testSetSize := uint32(30)

	circuit = *InitCircuitWithSize(testSetSize)

	testSet := GenTestSet(testSetSize, accumulator.TrustedSetup())
	witness = *AssignCircuit(testSet)
	assert.SolvingSucceeded(&circuit, &witness, test.WithCurves(ecc.BN254))
}

func TestZkMultiSwapFailCases(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit, witness Circuit
	testSetSize := uint32(10)

	circuit = *InitCircuitWithSize(testSetSize)

	testSet := GenTestSet(testSetSize, accumulator.TrustedSetup())
	witness = *AssignCircuit(testSet)

	// case for incorrect update of sum
	witness.UpdatedBalances[0] = 5
	assert.SolvingFailed(&circuit, &witness, test.WithCurves(ecc.BN254))
	//-------------------
	witness = *AssignCircuit(testSet)
	witness.OriginalSum = 10
	assert.SolvingFailed(&circuit, &witness, test.WithCurves(ecc.BN254))
	//-------------------
	witness = *AssignCircuit(testSet)
	witness.UpdatedSum = 10
	assert.SolvingFailed(&circuit, &witness, test.WithCurves(ecc.BN254))
	//-------------------
	witness = *AssignCircuit(testSet)
	witness.OriginalBalances[0] = 10
	assert.SolvingFailed(&circuit, &witness, test.WithCurves(ecc.BN254))

	// case for incorrect range of user balance
	witness = *AssignCircuit(testSet)
	witness.OriginalBalances[0] = 1000000000
	witness.UpdatedBalances[0] = 1000000000
	assert.SolvingFailed(&circuit, &witness, test.WithCurves(ecc.BN254))

	// case for incorrect remainders
	witness = *AssignCircuit(testSet)
	witness.RemainderR1 = testSet.ChallengeL1
	assert.SolvingFailed(&circuit, &witness, test.WithCurves(ecc.BN254))
	//-------------------
	witness = *AssignCircuit(testSet)
	witness.RemainderR2 = testSet.ChallengeL2
	assert.SolvingFailed(&circuit, &witness, test.WithCurves(ecc.BN254))
}
