package zkmultiswap

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

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
func Prove(input *UpdateSet32) (*groth16.Proof, error) {
	fmt.Println("Start Proving")
	fileName := KeyPathPrefix + "_" + strconv.FormatInt(int64(len(input.UserID)), 10)
	startingTime := time.Now().UTC()
	pk, err := groth16.ReadSegmentProveKey(fileName)
	if err != nil {
		fmt.Println("error while ReadSegmentProveKey")
		return nil, err
	}
	r1cs, err := groth16.LoadR1CSFromFile(fileName)
	if err != nil {
		fmt.Println("error while LoadR1CSFromFile")
		return nil, err
	}
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Loading a SNARK circuit and proving key for set size = %d, takes [%.3f] Seconds \n", len(input.UserID), duration.Seconds())

	assignment := AssignCircuit(input)
	witness, err := frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		fmt.Println("error while AssignCircuit")
		return nil, err
	}
	runtime.GC()
	startingTime = time.Now().UTC()
	proof, err := groth16.ProveRoll(r1cs, pk[0], pk[1], witness, fileName, backend.IgnoreSolverError()) // backend.IgnoreSolverError() can be used for testing
	if err != nil {
		fmt.Println("error while ProveRoll")
		return nil, err
	}
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Generating a SNARK proof for set size = %d, takes [%.3f] Seconds \n", len(input.UserID), duration.Seconds())
	return &proof, nil
}

// VerifyPublicWitness returns true is the public witness is valid for zkMultiSwap
func VerifyPublicWitness(publicWitness *witness.Witness, publicInfo *PublicInfo) bool {
	startingTime := time.Now().UTC()
	assignment2 := AssignCircuitHelper(publicInfo)
	publicWitness2, err := frontend.NewWitness(assignment2, ecc.BN254, frontend.PublicOnly())
	if err != nil {
		fmt.Println("Error generating NewWitness")
		return false
	}
	if !reflect.DeepEqual(publicWitness.Vector, publicWitness2.Vector) {
		fmt.Println("Verification failed for publicWitness")
		duration := time.Now().UTC().Sub(startingTime)
		fmt.Printf("Checking publicWitness using reflect takes [%.3f] Seconds \n", duration.Seconds())
		return false
	}
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Checking publicWitness using reflect takes [%.3f] Seconds \n", duration.Seconds())
	return true
}

// GenPublicWitness generates the publicWitness based on publicInfo
func GenPublicWitness(publicInfo *PublicInfo) *witness.Witness {
	assignment := AssignCircuitHelper(publicInfo)
	publicWitness, err := frontend.NewWitness(assignment, ecc.BN254, frontend.PublicOnly())
	if err != nil {
		fmt.Println("Error generating NewWitness in GenPublicWitness")
		return nil
	}
	return publicWitness
}

// Verify is used to check a Groth16 proof and public witness for the zkMultiSwap
func Verify(proof *groth16.Proof, setsize uint32, publicInfo *PublicInfo) bool {
	fileName := KeyPathPrefix + "_" + strconv.FormatInt(int64(setsize), 10)
	vk, err := LoadVerifyingKey(fileName)
	if err != nil {
		panic("r1cs init error")
	}
	runtime.GC()
	startingTime := time.Now().UTC()
	publicWitness := GenPublicWitness(publicInfo)
	if publicWitness == nil {
		return false
	}
	err = groth16.Verify(*proof, vk, publicWitness)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Verifying a SNARK proof for set size = %d, takes [%.3f] Seconds \n", setsize, duration.Seconds())
	if err != nil {
		fmt.Println("verify error = ", err)
		return false
	}
	return true
}
