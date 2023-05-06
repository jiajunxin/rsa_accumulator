package fiatshamir

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon"
)

// Based on the Miller-Robin test, the probability to have a non-prime probability is less than 1/(securityParaHash*4)
const securityParameter = 30

type Transcript struct {
	info []string
}

// InitTranscript inits a transcript with the input strings
func InitTranscript(input []string) *Transcript {
	var ret Transcript
	// we need a deep copy to make sure the transcript will not be changed
	copy(ret.info, input)
	return &ret
}

func (oldTranscript *Transcript) Append(newInfo string) {
	oldTranscript.info = append(oldTranscript.info, newInfo)
}

func (oldTranscript *Transcript) AppendSlice(newInfo []string) {
	oldTranscript.info = append(oldTranscript.info, newInfo...)
}

func (oldTranscript *Transcript) GetChallengeAndAppendTranscript() *big.Int {
	var ret big.Int
	ret.Set(HashToPrime(oldTranscript.info))
	oldTranscript.Append(ret.String())
	return &ret
}

func ElementFromString(v string) *fr.Element {
	n, success := new(big.Int).SetString(v, 10)
	if !success {
		panic("Error parsing hex number")
	}
	var e fr.Element
	e.SetBigInt(n)
	return &e
}

// HashToPrime takes the input into Poseidon and take the hash output to input repeatedly until we hit a prime number
func HashToPrime(input []string) *big.Int {
	var ret big.Int
	elementSlice := make([]*fr.Element, len(input))
	for i := range input {
		elementSlice[i] = new(fr.Element)
		elementSlice[i] = ElementFromString(input[i])
	}
	temp := poseidon.Poseidon(elementSlice...)
	temp.ToBigInt(&ret)
	flag := false
	for !flag {
		flag = ret.ProbablyPrime(securityParameter)
		if !flag {
			temp = poseidon.Poseidon(temp)
			temp.ToBigInt(&ret)
		}
	}
	return &ret
	// h := sha256.New()
	// for i := 0; i < len(input); i++ {
	// 	_, err := h.Write([]byte(input[i]))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// hashTemp := h.Sum(nil)
	// ret.SetBytes(hashTemp)
	// flag := false
	// for !flag {
	// 	flag = ret.ProbablyPrime(securityParameter)
	// 	if !flag {
	// 		h.Reset()
	// 		_, err := h.Write(hashTemp)
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		hashTemp = h.Sum(nil)
	// 		ret.SetBytes(hashTemp)
	// 	}
	// }
}
