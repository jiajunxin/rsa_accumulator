package accumulator

import (
	"math/big"
	"testing"
)

func TestGcd(t *testing.T) {
	var testObject AccumulatorSetup
	testObject = *TrustedSetup()

	var N big.Int
	N.Mul(&testObject.P, &testObject.Q)

	tmp := N.Cmp(&testObject.N)
	if tmp != 0 {
		t.Errorf("the N is not calculated correctly")
	}
	flag := isQR(&testObject.G, &testObject.P, &testObject.Q)
	if flag == false {
		t.Errorf("G is QR")
	}
}
