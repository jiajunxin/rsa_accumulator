package proof

import (
	"math/big"
	"reflect"
	"testing"
)

func Test_isPerfectSquare(t *testing.T) {
	type args struct {
		n *big.Int
	}
	tests := []struct {
		name     string
		args     args
		wantVal  *big.Int
		wantBool bool
	}{
		{
			name: "test_is_perfect_square",
			args: args{
				n: big.NewInt(4),
			},
			wantVal:  big.NewInt(2),
			wantBool: true,
		},
		{
			name: "test_is_not_perfect_square",
			args: args{
				n: big.NewInt(5),
			},
			wantVal:  big.NewInt(2),
			wantBool: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := isPerfectSquare(tt.args.n)
			if !reflect.DeepEqual(got, tt.wantVal) {
				t.Errorf("isPerfectSquare() got = %v, want %v", got, tt.wantVal)
			}
			if got1 != tt.wantBool {
				t.Errorf("isPerfectSquare() got1 = %v, want %v", got1, tt.wantBool)
			}
		})
	}
}
