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

	return nil
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
	return nil
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
