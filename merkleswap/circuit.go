package merkleswap

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/poseidon"
)

// TreeDepth denotes the depth of Merkle in the test circuit.
const TreeDepth = 28

// Circuit is the Merkle tree-based MultiSwap circuit for gnark.
// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The Circuit is only used for test purpose to measure the overhead of MerkleSwap
type Circuit struct {
	// default uses variable name and secret visibility.
	TestInputs []frontend.Variable `gnark:",public"` // test values
	//------------------------------private witness below--------------------------------------
	TestOutputs []frontend.Variable // test values
}

// Define declares the circuit constraints
func (circuit Circuit) Define(api frontend.API) error {
	for i := 0; i < len(circuit.TestInputs); i++ {
		tempHash0 := poseidon.Poseidon(api, circuit.TestInputs[i], circuit.TestInputs[i])
		for j := 0; j < TreeDepth; j++ {
			tempHash0 = poseidon.Poseidon(api, circuit.TestInputs[i], tempHash0)
		}
		api.AssertIsEqual(tempHash0, circuit.TestOutputs[i])
	}
	return nil
}
