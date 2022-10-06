package dihash

import (
	"crypto/sha256"
	"github.com/jiajunxin/rsa_accumulator/param"
	"math/big"
	"math/rand"
	"sync"
)

const (
	deltaValString = "30731438344250145947882657666206403727243332864808664054575262055190442942812700108124167942976" +
		"653745028212341196692947492080562974589240558404052155436479139607283861572110186639866316589725954212169900" +
		"473106847592072353357762907262662369230376196184226071545259316873351199416881666739376881925207433619609913" +
		"435128355340248285568061176332195286623104126482371089555666194830543043595601648501184952472930075767818065" +
		"617175977748228906417030406830990961578747315754348300610520710090878042950122953510395835606916522592211024" +
		"941845938097013497415239566963754154588561352876059012472806373183052035005766579987123343"
)

var (
	// delta is a 2000 bits large integer
	delta     *big.Int
	onceDelta sync.Once
	// max256 is set to 2^256 after init
	max256     *big.Int
	onceMax256 sync.Once
)

func Delta() *big.Int {
	onceDelta.Do(func() {
		delta = new(big.Int)
		if _, ok := delta.SetString(deltaValString, 10); !ok {
			panic("failed to set delta")
		}
	})
	return delta
}

func Max256() *big.Int {
	onceMax256.Do(func() {
		max256 = new(big.Int)
		max256.Lsh(param.Big1, 256)
		max256.Sub(max256, param.Big1)
	})
	return max256
}

// DIHash returns the Delta + Sha256(input)
func DIHash(input []byte) *big.Int {
	h := sha256.New()
	h.Write(input)
	hashBytes := h.Sum(nil)
	temp := new(big.Int).SetBytes(hashBytes)
	ret := new(big.Int).Add(Delta(), temp)
	return ret
}

// Get2048Rnd returns the Sha256(input||0) || Sha256(input||1) || ... Sha256(input||7)
func get2048Rnd(rnd *rand.Rand) *big.Int {
	h := sha256.New()
	var (
		ranNum    = new(big.Int)
		hashJoint []byte
	)

	for i := 0; i < 8; i++ {
		ranNum.Rand(rnd, Max256())
		tempBytes := append(ranNum.Bytes(), byte(i))
		h.Write(tempBytes)
		hashTemp := h.Sum(nil)
		hashJoint = append(hashJoint, hashTemp...)
	}

	ret := new(big.Int).SetBytes(hashJoint)
	return ret
}
