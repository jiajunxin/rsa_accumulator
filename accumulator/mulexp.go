package accumulator

import (
	"math/big"
)

// SimpleExp should calculate g^x mod n.
// It is implemented here to campare with golang's official Exp and MultiExp
func SimpleExp(g, x, n *big.Int) *big.Int {
	if g.Cmp(big1) <= 0 || n.Cmp(big1) <= 0 || x.Cmp(big1) < 0 {
		panic("invalid input for function SimpleExp")
	}
	// change x to its binary representation
	//binaryX := x.Bytes()
	bitLen := x.BitLen()
	//fmt.Println("BitLen = ", bitLen)
	bits := x.Bits() // bits is a slice of uint32
	//fmt.Println("Bitslice len = ", len(bits))
	var mask uint
	var gCopy, output big.Int
	gCopy.Set(g)
	output.SetInt64(1)
	for i := 0; i < bitLen; i++ {
		chunk := i / 64
		for j := 0; j < 64; j++ {
			mask = uint(1 << j)
			if (uint(bits[chunk]) & mask) == mask {
				output.Mul(&output, &gCopy)
				output.Mod(&output, n)
			}
			gCopy.Mul(&gCopy, &gCopy)
			gCopy.Mod(&gCopy, n)
		}
		i += 64
	}
	_ = big.DoubleExp(g, g, g, g)
	return &output
}

// GCB calculates the greatest common binaries of a and b.
// For example, if a = 1011 (binary) and b = 1100,
// the return will be of 1000(binary)
func GCB(a, b *big.Int) *big.Int {
	bitStringA := a.Bits()
	bitStringB := b.Bits()

	var maxBitLen int
	if len(bitStringA) > len(bitStringB) {
		maxBitLen = len(bitStringB)
	} else {
		maxBitLen = len(bitStringA)
	}
	bitStingsRet := make([]big.Word, maxBitLen)
	for i := 0; i < maxBitLen; i++ {
		bitStingsRet[i] = CommonBits(bitStringA[i], bitStringB[i])
		bitStringA[i] = bitStringA[i] - bitStingsRet[i]
		bitStringB[i] = bitStringB[i] - bitStingsRet[i]
	}
	var ret big.Int
	ret.SetBits(bitStingsRet)
	return &ret
}

// CommonBits calculates the greatest common binaries of a and b when they are uint.
func CommonBits(a, b big.Word) big.Word {
	var ret uint
	ret = 0
	var mask uint
	for i := 0; i < 32; i++ {
		mask = uint(1 << i)
		if ((uint(a) & mask) == mask) && ((uint(b) & mask) == mask) {
			//fmt.Println("i == ", i, "mask = ", mask)
			ret = uint(ret) | mask
		}
	}

	return big.Word(ret)
}
