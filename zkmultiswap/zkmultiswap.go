package zkmultiswap

import (
	"fmt"
	"math/big"
	"os"
	"runtime"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

const keyPathPrefix = "zkmultiswap"

// LoadVerifyingKey load the verification key from the filepath
func LoadVerifyingKey(filepath string) (verifyingKey groth16.VerifyingKey, err error) {
	verifyingKey = groth16.NewVerifyingKey(ecc.BN254)
	f, _ := os.Open(filepath + ".vk.save")
	_, err = verifyingKey.ReadFrom(f)
	if err != nil {
		return verifyingKey, fmt.Errorf("read file error")
	}
	err = f.Close()
	if err != nil {
		return verifyingKey, fmt.Errorf("close file error")
	}
	return verifyingKey, nil
}

func elementFromString(v string) *fr.Element {
	n, success := new(big.Int).SetString(v, 10)
	if !success {
		panic("Error parsing hex number")
	}
	var e fr.Element
	e.SetBigInt(n)
	return &e
}

// SetupZkMultiswap generates the circuit and public/verification keys with Groth16
// "keyPathPrefix".pk* are for public keys, "keyPathPrefix".ccs* are for r1cs, "keyPathPrefix".vk,save is for verification keys
func SetupZkMultiswap(size uint32) {
	// compiles our circuit into a R1CS
	circuit := InitCircuitWithSize(size)
	fmt.Println("Start Compiling")
	r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, circuit)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish Compiling")
	runtime.GC()
	fmt.Println("Number of constrains: ", r1cs.GetNbConstraints())

	err = groth16.SetupLazyWithDump(r1cs, keyPathPrefix)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish Setup")
}

// AssignWitness to do
func AssignWitness() *ZKMultiSwapCircuit {
	//var ret ZKMultiSwapCircuit
	return InitCircuitWithSize(1000)
}

// Prove is used to generate a Groth16 proof and public witness for the zkMultiSwap
func Prove() (*groth16.Proof, *witness.Witness, error) {
	fmt.Println("Start Proving")
	pk, err := groth16.ReadSegmentProveKey(keyPathPrefix)
	if err != nil {
		return nil, nil, err
	}
	runtime.GC()
	r1cs, err := groth16.LoadR1CSFromFile(keyPathPrefix)
	if err != nil {
		return nil, nil, err
	}

	// Todo: witness to be input
	//outputs := elementFromString("17517277496620338529366114881698763424837036587329561912313499393581702161864")
	assignment := AssignWitness()

	witness, err := frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		return nil, nil, err
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return nil, nil, err
	}
	proof, err := groth16.ProveRoll(r1cs, pk[0], pk[1], witness, keyPathPrefix, backend.IgnoreSolverError())
	if err != nil {
		return nil, nil, err
	}
	return &proof, publicWitness, nil
}

// Verify is used to check a Groth16 proof and public witness for the zkMultiSwap
func Verify(proof *groth16.Proof, publicWitness *witness.Witness) bool {
	vk, err := LoadVerifyingKey(keyPathPrefix)
	if err != nil {
		panic("r1cs init error")
	}
	runtime.GC()

	err = groth16.Verify(*proof, vk, publicWitness)
	return err != nil
}

// TestMultiSwap is temporarily used for test purpose
func TestMultiSwap() {
	fmt.Println("Start TestMultiSwap")
	SetupZkMultiswap(100)

	proof, publicWitness, err := Prove()
	if err != nil {
		panic(err)
	}

	flag := Verify(proof, publicWitness)
	if flag {
		fmt.Println("Verification passed")
	}
	fmt.Println("Verification failed")
}
