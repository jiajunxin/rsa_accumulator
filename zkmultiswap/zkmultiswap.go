package zkmultiswap

import (
	"fmt"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/poseidon"
)

// gnark is a zk-SNARK library written in Go. Circuits are regular structs.
// The inputs must be of type frontend.Variable and make up the witness.
// The witness has a
//   - secret part --> known to the prover only
//   - public part --> known to the prover and the verifier
type Circuit struct {
	Secret frontend.Variable // pre-image of the hash secret known to the prover only
	Hash   frontend.Variable `gnark:",public"` // hash of the secret known to all
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (circuit *Circuit) Define(api frontend.API) error {
	hashOutput := poseidon.Poseidon(api, circuit.Secret, circuit.Secret)
	api.AssertIsEqual(circuit.Hash, hashOutput)
	return nil
}

// TestMultiSwap
func TestMultiSwap() {
	fmt.Println("Start TestMultiSwap")
	// compiles our circuit into a R1CS
	var circuit Circuit
	fmt.Println("Start Compiling")
	ccs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Finish Compiling")
	fmt.Printf("ccs: %v\n", ccs)

	// groth16 zkSNARK: Setup
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		fmt.Println(err)
	}

	// witness definition
	assignment := Circuit{Secret: 3, Hash: 35}
	witness, err := frontend.NewWitness(&assignment, ecc.BN254)
	if err != nil {
		fmt.Println(err)
	}
	publicWitness, err := witness.Public()
	if err != nil {
		fmt.Println(err)
	}

	// groth16: Prove & Verify
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		fmt.Println(err)
	}
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		fmt.Println(err)
	}
}
