package fiatshamir

import (
	"crypto/sha256"
	"math/big"
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

// HashToPrime takes the input into Sha256 and take the hash output to input repeatedly until we hit a prime number
func HashToPrime(input []string) *big.Int {
	var ret big.Int
	h := sha256.New()
	for i := 0; i < len(input); i++ {
		_, err := h.Write([]byte(input[i]))
		if err != nil {
			panic(err)
		}
	}
	hashTemp := h.Sum(nil)
	ret.SetBytes(hashTemp)
	flag := false
	for !flag {
		flag = ret.ProbablyPrime(securityParameter)
		if !flag {
			h.Reset()
			_, err := h.Write(hashTemp)
			if err != nil {
				panic(err)
			}
			hashTemp = h.Sum(nil)
			ret.SetBytes(hashTemp)
		}
	}
	return &ret
}
