package zkmultiswap

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	bnPoseidon "github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon"
	"github.com/consensys/gnark/test"
	iden3Poseidon "github.com/iden3/go-iden3-crypto/poseidon"
	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

// compare return 0 if input1 == input2
func compare(input1 *big.Int, input2 []byte) int {
	input1bytes := input1.Bytes()
	return bytes.Compare(input1bytes, input2)
}

func TestPoseidonHash(t *testing.T) {
	inputs := "anything"
	result1, err := iden3Poseidon.HashBytes([]byte(inputs))
	if err != nil {
		t.Errorf(err.Error())
	}

	poseidonHasher := bnPoseidon.NewPoseidon()
	_, err = poseidonHasher.Write([]byte(inputs))
	if err != nil {
		t.Errorf(err.Error())
	}
	result2 := poseidonHasher.Sum(nil)

	if compare(result1, result2) != 0 {
		fmt.Println("result1 = ", result1.String())
		fmt.Println("result2 = ", result2)
	}
}

func TestZkMultiSwap(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit, witness Circuit
	testSetSize := uint32(30)

	circuit = *InitCircuitWithSize(testSetSize)

	testSet := GenTestSet(testSetSize, accumulator.TrustedSetup())
	witness = *InitCircuit(testSet)
	assert.SolvingSucceeded(&circuit, &witness, test.WithCurves(ecc.BN254))
}
