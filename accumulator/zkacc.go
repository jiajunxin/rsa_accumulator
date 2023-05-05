package accumulator

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"time"
)

// GenRandomizer outputs random number uniformly between 0 to 2^2047
func GenRandomizer() *big.Int {
	ranNum, err := crand.Int(crand.Reader, Min2048)
	if err != nil {
		panic(err)
	}
	return ranNum
}

// ZKAccumulate generates one accumulator which is zero-knowledge
func ZKAccumulate(set []string, encodeType EncodeType, setup *Setup) (*big.Int, []*big.Int) {
	startingTime := time.Now().UTC()
	rep := GenRepresentatives(set, encodeType)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running GenRepresentatives Takes [%.3f] Seconds \n",
		duration.Seconds())

	r := GenRandomizer()
	base := AccumulateNew(setup.G, r, setup.N)

	proofs := ProveMembership(base, setup.N, rep)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation
	acc := AccumulateNew(proofs[0], rep[0], setup.N)

	return acc, proofs
}
