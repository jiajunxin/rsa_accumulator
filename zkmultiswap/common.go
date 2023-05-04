package zkmultiswap

import (
	"fmt"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

const (
	// BitLength is the bit length of the user ID, balnace and epoch number. It can be 32, 64 or any valid number within the field
	BitLength = 32

	keyPathPrefix = "zkmultiswap"
)

func ElementFromString(v string) *fr.Element {
	n, success := new(big.Int).SetString(v, 10)
	if !success {
		panic("Error parsing hex number")
	}
	var e fr.Element
	e.SetBigInt(n)
	return &e
}

func ElementFromUint32(v uint32) *fr.Element {
	var e fr.Element
	e.SetInt64(int64(v))
	return &e
}

// Set32 is one set for the prover with uint32 for CurrentEpochNum,
type UpdateSet32 struct {
	ChallengeL1      *big.Int
	ChallengeL2      *big.Int
	RemainderR1      *big.Int
	RemainderR2      *big.Int
	Randomizer       *big.Int
	CurrentEpochNum  uint32
	OriginalSum      uint32
	UpdatedSum       uint32
	UserID           []uint32
	OriginalBalances []uint32
	OriginalHashes   []*big.Int
	OriginalUpdEpoch []uint32
	UpdatedBalances  []uint32
}

// func GenTestSet(setsize uint32) *UpdateSet32 {
// 	var ret UpdateSet32

// 	ret.CurrentEpochNum = 500
// 	for i := uint32(0); i < setsize; i++ {
// 		j := i*2 + 1
// 		ret.UserID[i] = j
// 		ret.OriginalBalances[i] = j
// 		ret.OriginalUpdEpoch[i] = 10
// 		ret.OriginalHashes[i].SetInt64(int64(j))
// 		ret.UpdatedBalances[i] = j
// 	}

// 	// get challenge

// 	poseidonHasher := poseidon.Poseidon(ElementFromUint32(ret.UserID[0]))
// 	return &ret
// }

// TestMultiSwap is temporarily used for test purpose
func TestMultiSwap() {
	fmt.Println("Start TestMultiSwap")
	testSetSize := uint32(100)
	SetupZkMultiswap(testSetSize)

	proof, publicWitness, err := Prove()
	if err != nil {
		panic(err)
	}

	flag := Verify(proof, testSetSize, publicWitness)
	if flag {
		fmt.Println("Verification passed")
	}
	fmt.Println("Verification failed")
}
