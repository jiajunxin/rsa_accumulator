package solidity

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/examples/cubic"
	"github.com/consensys/gnark/frontend"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/jiajunxin/rsa_accumulator/accumulator"
	"github.com/jiajunxin/rsa_accumulator/zkmultiswap"
	"github.com/stretchr/testify/suite"
)

const KeyPathPrefix = "zkmultiswap"
const testSetSize = 5 //for Gorth16, the testSetSize should not affect the gas cost

type ExportSolidityTestSuiteGroth16 struct {
	suite.Suite

	// backend
	backend *backends.SimulatedBackend

	// verifier contract
	verifierContract *Verifier

	// groth16 gnark objects
	vk      groth16.VerifyingKey
	pk      []groth16.ProvingKey
	circuit cubic.Circuit
	r1cs    frontend.CompiledConstraintSystem

	address common.Address
}

func TestRunExportSolidityTestSuiteGroth16(t *testing.T) {
	suite.Run(t, new(ExportSolidityTestSuiteGroth16))
}

func (t *ExportSolidityTestSuiteGroth16) SetupTest() {

	const gasLimit uint64 = 4712388

	// setup simulated backend
	key, _ := crypto.GenerateKey()
	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	t.NoError(err, "init keyed transactor")

	genesis := map[common.Address]core.GenesisAccount{
		auth.From: {Balance: big.NewInt(1000000000000000000)}, // 1 Eth
	}
	t.backend = backends.NewSimulatedBackend(genesis, gasLimit)

	// deploy verifier contract
	addr, _, v, err := DeployVerifier(auth, t.backend)
	t.address = addr
	t.NoError(err, "deploy verifier contract failed")
	t.verifierContract = v
	t.backend.Commit()

	t.r1cs, err = groth16.LoadR1CSFromFile(KeyPathPrefix)
	if err != nil {
		fmt.Println("error while LoadR1CSFromFile")
	}
	t.NoError(err, "compiling R1CS failed")

	// read proving and verifying keys
	t.pk, err = groth16.ReadSegmentProveKey(KeyPathPrefix)
	if err != nil {
		fmt.Println("error while ReadSegmentProveKey")
	}
	t.vk = groth16.NewVerifyingKey(ecc.BN254)
	{
		f, _ := os.Open(KeyPathPrefix + ".vk.save")
		_, err = t.vk.ReadFrom(f)
		f.Close()
		t.NoError(err, "reading verifying key failed")
	}

}

func (t *ExportSolidityTestSuiteGroth16) TestVerifyProof() {

	// create a valid proof
	testSet := zkmultiswap.GenTestSet(testSetSize, accumulator.TrustedSetup())
	assignment := zkmultiswap.AssignCircuit(testSet)
	witness, err := frontend.NewWitness(assignment, ecc.BN254)
	if err != nil {
		fmt.Println("error while AssignCircuit")
	}

	// prove
	proof, err := groth16.ProveRoll(t.r1cs, t.pk[0], t.pk[1], witness, KeyPathPrefix, backend.IgnoreSolverError()) // backend.IgnoreSolverError() can be used for testing
	t.NoError(err, "proving failed")

	// ensure gnark (Go) code verifies it
	publicInfo := testSet.PublicPart()
	publicWitness := zkmultiswap.GenPublicWitness(publicInfo)
	if publicWitness == nil {
		fmt.Println("error while publicWitness")
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
	//err = groth16.Verify(proof, verifyingKey, publicWitness)
	//t.NoError(err, "verifying failed")

	// get proof bytes
	const fpSize = 4 * 8
	var buf bytes.Buffer
	_, err = proof.WriteRawTo(&buf)
	if err != nil {
		fmt.Println(err.Error())
	}
	proofBytes := buf.Bytes()

	// solidity contract inputs
	var (
		a     [2]*big.Int
		b     [2][2]*big.Int
		c     [2]*big.Int
		input [7]*big.Int
	)
	for i := 0; i < 7; i++ {
		input[i] = new(big.Int)
	}

	// proof.Ar, proof.Bs, proof.Krs
	a[0] = new(big.Int).SetBytes(proofBytes[fpSize*0 : fpSize*1])
	a[1] = new(big.Int).SetBytes(proofBytes[fpSize*1 : fpSize*2])
	b[0][0] = new(big.Int).SetBytes(proofBytes[fpSize*2 : fpSize*3])
	b[0][1] = new(big.Int).SetBytes(proofBytes[fpSize*3 : fpSize*4])
	b[1][0] = new(big.Int).SetBytes(proofBytes[fpSize*4 : fpSize*5])
	b[1][1] = new(big.Int).SetBytes(proofBytes[fpSize*5 : fpSize*6])
	c[0] = new(big.Int).SetBytes(proofBytes[fpSize*6 : fpSize*7])
	c[1] = new(big.Int).SetBytes(proofBytes[fpSize*7 : fpSize*8])

	// public witness
	input[0] = &publicInfo.ChallengeL1
	input[1] = &publicInfo.ChallengeL2
	input[2] = &publicInfo.RemainderR1
	input[3] = &publicInfo.RemainderR2
	input[4].SetInt64(int64(publicInfo.CurrentEpochNum))
	input[5] = &publicInfo.DeltaModL1
	input[6] = &publicInfo.DeltaModL2

	//------
	snarkInput := make([]interface{}, 0)
	snarkInput = append(snarkInput, a)
	snarkInput = append(snarkInput, b)
	snarkInput = append(snarkInput, c)
	snarkInput = append(snarkInput, input)
	//snarkInput = {a, b, c, input}
	parsed, err := VerifierMetaData.GetAbi()
	data, err := parsed.Pack("verifyProof", snarkInput...)
	if err != nil {
		panic(err)
	}
	msg := ethereum.CallMsg{From: bind.CallOpts{}.From, To: &t.address, Data: data}
	gasLimit, err := t.backend.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("Failed to estimate gas needed: %v", err)
	}
	fmt.Println("Gas Limit:", gasLimit)
	//------

	// call the contract
	res, err := t.verifierContract.VerifyProof(&bind.CallOpts{}, a, b, c, input)
	if t.NoError(err, "calling verifier on chain gave error") {
		t.True(res, "calling verifier on chain didn't succeed")
	}

	// (wrong) public witness
	input[0] = new(big.Int).SetUint64(42)

	// call the contract should fail
	res, err = t.verifierContract.VerifyProof(&bind.CallOpts{}, a, b, c, input)
	if t.NoError(err, "calling verifier on chain gave error") {
		t.False(res, "calling verifier on chain succeed, and shouldn't have")
	}
}
