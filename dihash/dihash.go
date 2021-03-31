package dihash

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/rand"
)

// Delta is a 2000 bits large integer
var Delta = big.NewInt(0)

// Max256 is set to 2^256 after init
var Max256 = big.NewInt(0)

var one = big.NewInt(1)

func init() {
	_, flag := Delta.SetString("30731438344250145947882657666206403727243332864808664054575262055190442942812700108124167942976653745028212341196692947492080562974589240558404052155436479139607283861572110186639866316589725954212169900473106847592072353357762907262662369230376196184226071545259316873351199416881666739376881925207433619609913435128355340248285568061176332195286623104126482371089555666194830543043595601648501184952472930075767818065617175977748228906417030406830990961578747315754348300610520710090878042950122953510395835606916522592211024941845938097013497415239566963754154588561352876059012472806373183052035005766579987123343", 10)
	if !flag {
		fmt.Printf("flag is false, init of DI Hash failed")
		panic(0)
	}
	_ = Max256.Lsh(one, 256)
	_ = Max256.Sub(Max256, one)
	//fmt.Println("the Delta = ", Delta.String())
	//fmt.Println("the Max256 in bits = ", Max256.Text(2))
}

// GetDIHash returns the Delta + Sha256(input) with the last bit 1
func GetDIHash(rnd *rand.Rand) *big.Int {
	h := sha256.New()
	var ranNum, temp, ret big.Int
	ranNum.Rand(rnd, Max256)
	h.Write(ranNum.Bytes())
	hashTemp := h.Sum(nil)
	temp.SetBytes(hashTemp)
	_ = temp.SetBit(&temp, 0, 0)
	_ = ret.Add(Delta, &temp)
	return &ret
}

// Get2048Rnd returns the Sha256(input||0) || Sha256(input||1) || ... Sha256(input||7)
func Get2048Rnd(rnd *rand.Rand) *big.Int {
	h := sha256.New()
	var ranNum, ret big.Int
	var hashJoint []byte

	for i := 0; i < 8; i++ {
		ranNum.Rand(rnd, Max256)
		tempBytes := append(ranNum.Bytes(), byte(i))
		h.Write(tempBytes)
		hashTemp := h.Sum(nil)
		hashJoint = append(hashJoint, hashTemp...)
	}

	ret.SetBytes(hashJoint)
	return &ret
}
