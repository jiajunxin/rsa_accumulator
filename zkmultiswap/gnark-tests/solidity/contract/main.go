package main

import (
	"fmt"
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"github.com/jiajunxin/rsa_accumulator/zkmultiswap"
)

const KeyPathPrefix = "zkmultiswap"
const testSetSize = 5

func main() {
	err := generateGroth16()
	if err != nil {
		log.Fatal("groth16 error:", err)
	}
}

func generateGroth16() error {
	//var circuit zkmultiswap.Circuit
	circuit := zkmultiswap.InitCircuitWithSize(testSetSize)

	r1cs, err := frontend.Compile(ecc.BN254, r1cs.NewBuilder, circuit)
	if err != nil {
		return err
	}

	err = groth16.SetupLazyWithDump(r1cs, KeyPathPrefix)
	if err != nil {
		return err
	}
	verifyingKey := groth16.NewVerifyingKey(ecc.BN254)
	f, _ := os.Open(KeyPathPrefix + ".vk.save")
	_, err = verifyingKey.ReadFrom(f)
	if err != nil {
		fmt.Println("read file error")
	}
	err = f.Close()
	if err != nil {
		fmt.Println("close file error")
	}
	// _, vk, err := groth16.Setup(r1cs)
	{
		f, err := os.Create("Notuscontract_g16.sol")
		if err != nil {
			return err
		}
		err = verifyingKey.ExportSolidity(f)
		if err != nil {
			return err
		}
	}

	return nil
}
