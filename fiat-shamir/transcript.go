package fiatshamir

import (
	"crypto/sha256"
	"fmt"
	"math/big"
)

// Based on the Miller-Robin test, the probability to have a non-prime probability is less than 1/(securityParaHash*4)
const securityParameter = 30

var Min253 big.Int

func init() {
	Min253.SetInt64(1)
	_ = Min253.Lsh(&Min253, 252)
}

type Transcript struct {
	info []string
}

// Print outputs the info in the transcript
func (transcript *Transcript) Print() {
	fmt.Println("The transcript has ", len(transcript.info), "strings as info.")
	for i := range transcript.info {
		fmt.Println("Info[", i, "] = ", transcript.info[i])
	}
}

// InitTranscript inits a transcript with the input strings
func InitTranscript(input []string) *Transcript {
	var ret Transcript
	// we need a deep copy to make sure the transcript will not be changed
	ret.info = append(ret.info, input...)
	return &ret
}

// Append add new info into the transcript
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

func SetFieldElement(input []byte) *big.Int {
	var ret big.Int
	ret.SetBytes(input)
	if ret.Cmp(&Min253) != 0 {
		ret.Mod(&ret, &Min253)
	}
	return &ret
}

// HashToPrime takes the input into Poseidon and take the hash output to input repeatedly until we hit a prime number
func HashToPrime(input []string) *big.Int {
	h := sha256.New()
	for i := 0; i < len(input); i++ {
		_, err := h.Write([]byte(input[i]))
		if err != nil {
			panic(err)
		}
	}
	hashTemp := h.Sum(nil)
	ret := SetFieldElement(hashTemp)
	flag := false
	for !flag {
		flag = ret.ProbablyPrime(securityParameter)
		if !flag {
			h.Reset()
			_, err := h.Write(ret.Bytes())
			if err != nil {
				panic(err)
			}
			hashTemp = h.Sum(nil)
			ret = SetFieldElement(hashTemp)
		}
	}
	return ret
}
