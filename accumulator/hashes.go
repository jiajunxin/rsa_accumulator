package accumulator

import (
	"crypto/sha256"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/poseidon"
)

// ElementFromBigInt returns an element in BN256 generated from BigInt
func ElementFromBigInt(v *big.Int) *fr.Element {
	var e fr.Element
	e.SetBigInt(v)
	return &e
}

// ElementFromString returns an element in BN256 generated from string of decimal integers
func ElementFromString(v string) *fr.Element {
	n, success := new(big.Int).SetString(v, 10)
	if !success {
		panic("Error parsing hex number")
	}
	var e fr.Element
	e.SetBigInt(n)
	return &e
}

// ElementFromUint32 returns an element in BN256 generated from uint32
func ElementFromUint32(v uint32) *fr.Element {
	var e fr.Element
	e.SetInt64(int64(v))
	return &e
}

// HashToPrime takes the input into Sha256 and take the hash output to input repeatedly until we hit a prime number
func HashToPrime(input []byte) *big.Int {
	var ret big.Int
	h := sha256.New()
	_, err := h.Write(input)
	if err != nil {
		panic(err)
	}
	hashTemp := h.Sum(nil)
	ret.SetBytes(hashTemp)
	flag := false
	for !flag {
		flag = ret.ProbablyPrime(securityParaHashToPrime)
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

// SHA256ToInt calculates the input with Sha256 and change it to big.Int
func SHA256ToInt(input []byte) *big.Int {
	var ret big.Int
	h := sha256.New()
	_, err := h.Write(input)
	if err != nil {
		panic(err)
	}
	hashTemp := h.Sum(nil)
	ret.SetBytes(hashTemp)
	return &ret
}

// PoseidonWith2Inputs inputs 2 big.Int and generate a Poseidon hash result.
func PoseidonWith2Inputs(inputs []*big.Int) *big.Int {
	if len(inputs) != 2 {
		panic("PoseidonWith2Inputs requires 2 inputs")
	}
	fieldElement := poseidon.Poseidon(ElementFromBigInt(inputs[0]), (ElementFromBigInt(inputs[1])))
	var ret big.Int
	fieldElement.ToBigIntRegular(&ret)
	return &ret
}

// UniversalHashToInt calculates output = input * A + B mod P
func UniversalHashToInt(input *big.Int) *big.Int {
	var ret big.Int
	ret.Mul(input, A)
	ret.Mod(&ret, P)
	ret.Add(&ret, B)
	ret.Mod(&ret, P)
	if ret.Bit(0) == 0 {
		ret.Add(&ret, big1)
	}
	return &ret
}

// DIHashPoseidon generates DI hash with Poseidon hash
func DIHashPoseidon(input ...*fr.Element) *big.Int {
	ret := new(big.Int)
	temp := poseidon.Poseidon(input...)
	temp.ToBigIntRegular(ret)
	ret.Add(ret, Min1024)
	return ret
}

// PoseidonAndDIHash returns the Poseidon Hash result together with DI hash result
func PoseidonAndDIHash(input ...*fr.Element) (*fr.Element, *big.Int) {
	ret := new(big.Int)
	temp := poseidon.Poseidon(input...)
	temp.ToBigIntRegular(ret)
	ret.Add(ret, Min1024)
	return temp, ret
}
