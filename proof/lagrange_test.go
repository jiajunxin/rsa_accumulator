package proof

import (
	"math/big"
	"reflect"
	"testing"
)

var (
	big2Pow20 big.Int
	big2Pow32 big.Int
)

func setup() {
	var ok bool
	if _, ok = big2Pow20.SetString("1048576", 10); !ok {
		panic("failed to set big2Pow20")
	}
	if _, ok = big2Pow32.SetString("4294967296", 10); !ok {
		panic("failed to set big2Pow32")
	}
}

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

func Test_preCompute(t *testing.T) {
	type args struct {
		n *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr bool
	}{
		{
			name: "test_8",
			args: args{
				n: big.NewInt(8),
			},
			want:    big.NewInt(6),
			wantErr: false,
		},
		{
			name: "test_2^20",
			args: args{
				n: &big2Pow20,
			},
			want:    big.NewInt(9699690),
			wantErr: false,
		},
		{
			name: "test_2^32",
			args: args{
				n: &big2Pow32,
			},
			want:    big.NewInt(200560490130),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			got, err := preCompute(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("preCompute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("preCompute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
