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
			want:    nil,
			wantErr: true,
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

func TestVerify(t *testing.T) {
	type args struct {
		target *big.Int
		w1     *big.Int
		w2     *big.Int
		w3     *big.Int
		w4     *big.Int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_verify_success_4",
			args: args{
				target: big.NewInt(4),
				w1:     big.NewInt(2),
				w2:     big.NewInt(0),
				w3:     big.NewInt(0),
				w4:     big.NewInt(0),
			},
			want: true,
		},
		{
			name: "test_verify_success_35955023",
			args: args{
				target: big.NewInt(35955023),
				w1:     big.NewInt(2323),
				w2:     big.NewInt(5454),
				w3:     big.NewInt(893),
				w4:     big.NewInt(123),
			},
			want: true,
		},
		{
			name: "test_verify_fail_35955024",
			args: args{
				target: big.NewInt(35955024),
				w1:     big.NewInt(2323),
				w2:     big.NewInt(5454),
				w3:     big.NewInt(893),
				w4:     big.NewInt(123),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := [squareNum]*big.Int{
				tt.args.w1,
				tt.args.w2,
				tt.args.w3,
				tt.args.w4,
			}
			if got := Verify(tt.args.target, fs); got != tt.want {
				t.Errorf("Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLagrangeFourSquares(t *testing.T) {
	type args struct {
		n *big.Int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test_4",
			args: args{
				n: big.NewInt(4),
			},
			wantErr: false,
		},
		{
			name: "test_35955023",
			args: args{
				n: big.NewInt(35955023),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LagrangeFourSquares(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("LagrangeFourSquares() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !Verify(tt.args.n, got) {
				t.Errorf("LagrangeFourSquares() verify failed, got: %v != %v", got, tt.args.n)
			}
		})
	}
}
