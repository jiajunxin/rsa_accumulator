package zkmultiswap

import (
	"fmt"
	"os"
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
	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				runtime.GC()
			}
		}
	}()
	r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, circuit)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish Compiling")
	fmt.Println("Number of constrains: ", r1cs.GetNbConstraints())

	fileName := keyPathPrefix + "_" + strconv.FormatInt(int64(size), 10)
	err = groth16.SetupLazyWithDump(r1cs, fileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finish Setup")
}

// AssignWitness to do
func AssignWitness() *Circuit {
	//var ret ZKMultiSwapCircuit
	return InitCircuitWithSize(1000)
}

// Prove is used to generate a Groth16 proof and public witness for the zkMultiSwap
func Prove() (*groth16.Proof, *witness.Witness, error) {
	fmt.Println("Start Proving")
	go func() {
		for {
			select {
			case <-time.After(time.Second * 10):
				runtime.GC()
			}
		}
	}()
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

func VerifyPublicWitness(*witness.Witness) bool {
	//Todo.
	return true
}

// Verify is used to check a Groth16 proof and public witness for the zkMultiSwap
func Verify(proof *groth16.Proof, setsize uint32, publicWitness *witness.Witness) bool {
	fileName := keyPathPrefix + "_" + strconv.FormatInt(int64(setsize), 10)
	vk, err := LoadVerifyingKey(fileName)
	if err != nil {
		panic("r1cs init error")
	}
	runtime.GC()
	if !VerifyPublicWitness(publicWitness) {
		return false
	}

	err = groth16.Verify(*proof, vk, publicWitness)
	return err != nil
}