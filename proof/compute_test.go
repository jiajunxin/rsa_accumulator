package proof

import (
	"math/big"
	"reflect"
	"testing"
)

func Test_log2(t *testing.T) {
	type args struct {
		n *big.Int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test_0",
			args: args{
				n: big.NewInt(0),
			},
			want: -1,
		},
		{
			name: "test_1",
			args: args{
				n: big.NewInt(1),
			},
			want: 0,
		},
		{
			name: "test_32",
			args: args{
				n: big.NewInt(32),
			},
			want: 5,
		},
		{
			name: "test_35",
			args: args{
				n: big.NewInt(35),
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := log2(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("log2() = %v, want %v", got, tt.want)
			}
		})
	}
}
