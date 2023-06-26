package merkleswap

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

// KeyPathPrefix denotes the file name for Merkle MultiSwap circuits
const KeyPathPrefix = "merkleswap"

// InitCircuitWithSize init a circuit with challenges, OriginalHashes and CurrentEpochNum value 1, all other values 0. Use for test purpose only.
func InitCircuitWithSize(size uint32) *Circuit {
	var circuit Circuit
	circuit.TestInputs = make([]frontend.Variable, size)
	circuit.TestOutputs = make([]frontend.Variable, size)
	for i := uint32(0); i < size; i++ {
		circuit.TestInputs[i] = 0
		circuit.TestOutputs[i] = 0
	}
	return &circuit
}

// AssignCircuit assign a circuit with UpdateSet32 values.
func AssignCircuit(size uint32) *Circuit {
	var circuit Circuit
	helper := make([]uint32, size)
	circuit.TestInputs = make([]frontend.Variable, size)
	circuit.TestOutputs = make([]frontend.Variable, size)
	for i := uint32(0); i < size; i++ {
		circuit.TestInputs[i] = i
		helper[i] = i
	}
	for i := uint32(0); i < size; i++ {
		tempHash0 := poseidon.Poseidon(accumulator.ElementFromUint32(helper[i]), accumulator.ElementFromUint32(helper[i]))
		for j := 0; j < TreeDepth; j++ {
			tempHash0 = poseidon.Poseidon(accumulator.ElementFromUint32(helper[i]), tempHash0)
		}
		circuit.TestOutputs[i] = tempHash0
	}
	return &circuit
}

//func genTestSet(size uint32)

// TestMerkleMultiSwap is temporarily used for test purpose
func TestMerkleMultiSwap(testSetSize uint32) {
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

	proof, err := Prove(testSetSize)
	if err != nil {
		fmt.Println("Error during Prove")
		panic(err)
	}
	runtime.GC()

	flag := Verify(proof, testSetSize)
	if flag {
		fmt.Println("Verification passed")
		return
	}
	fmt.Println("Verification failed")
}

func isCircuitExist(testSetSize uint32) bool {
	fileName := KeyPathPrefix + "_" + strconv.FormatInt(int64(testSetSize), 10) + ".ccs.save"
	_, err := os.Stat(fileName)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

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

// SetupZkMultiswap generates the circuit and public/verification keys with Groth16
// "keyPathPrefix".pk* are for public keys, "keyPathPrefix".ccs* are for r1cs, "keyPathPrefix".vk,save is for verification keys
func SetupZkMultiswap(size uint32) {
	// compiles our circuit into a R1CS
	circuit := InitCircuitWithSize(size)
	fmt.Println("Start Compiling")
	r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, circuit) //, frontend.IgnoreUnconstrainedInputs()
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish Compiling")
	fmt.Println("Number of constrains: ", r1cs.GetNbConstraints())

	fileName := KeyPathPrefix + "_" + strconv.FormatInt(int64(size), 10)
	err = groth16.SetupLazyWithDump(r1cs, fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish Setup")
}

// Prove is used to generate a Groth16 proof and public witness for the zkMultiSwap
func Prove(testSetSize uint32) (*groth16.Proof, error) {
	fmt.Println("Start Proving")
	fileName := KeyPathPrefix + "_" + strconv.FormatInt(int64(testSetSize), 10)
	startingTime := time.Now().UTC()
	pk, err := groth16.ReadSegmentProveKey(fileName)
	if err != nil {
		fmt.Println("error while ReadSegmentProveKey")
		return nil, err
	}
	runtime.GC()
	r1cs, err := groth16.LoadR1CSFromFile(fileName)
	if err != nil {
		fmt.Println("error while LoadR1CSFromFile")
		return nil, err
	}
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Loading a SNARK circuit and proving key for set size = %d, takes [%.3f] Seconds \n", testSetSize, duration.Seconds())

	assignment := AssignCircuit(testSetSize)
	witness, err := frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		fmt.Println("error while AssignCircuit")
		return nil, err
	}

	startingTime = time.Now().UTC()
	proof, err := groth16.ProveRoll(r1cs, pk[0], pk[1], witness, fileName, backend.IgnoreSolverError()) // backend.IgnoreSolverError() can be used for testing
	if err != nil {
		fmt.Println("error while ProveRoll")
		return nil, err
	}
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generating a SNARK proof Merkle tree MultiSwap for set size = %d, takes [%.3f] Seconds \n", testSetSize, duration.Seconds())
	return &proof, nil
}

// Verify is used to check a Groth16 proof and public witness for the zkMultiSwap
func Verify(proof *groth16.Proof, setsize uint32) bool {
	fileName := KeyPathPrefix + "_" + strconv.FormatInt(int64(setsize), 10)
	vk, err := LoadVerifyingKey(fileName)
	if err != nil {
		panic("r1cs init error")
	}
	runtime.GC()

	assignment := AssignCircuit(setsize)
	publicWitness, err := frontend.NewWitness(assignment, ecc.BN254, frontend.PublicOnly())
	if err != nil {
		fmt.Println("Error generating NewWitness in GenPublicWitness")
		return false
	}
	if publicWitness == nil {
		return false
	}
	startingTime := time.Now().UTC()
	err = groth16.Verify(*proof, vk, publicWitness)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Verifying a SNARK proof for set size = %d, takes [%.3f] Seconds \n", setsize, duration.Seconds())
	if err != nil {
		fmt.Println("verify error = ", err)
		return false
	}
	return true
}
