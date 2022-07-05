package proof

import (
	"math/big"
	"reflect"
	"testing"
)

func Test_log10(t *testing.T) {
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
			name: "log10_0",
			args: args{
				n: big.NewInt(0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "log10_1",
			args: args{
				n: big.NewInt(1),
			},
			want:    big.NewInt(0),
			wantErr: false,
		},
		{
			name: "log10_10",
			args: args{
				n: big.NewInt(10),
			},
			want:    big.NewInt(1),
			wantErr: false,
		},
		{
			name: "log10_20",
			args: args{
				n: big.NewInt(20),
			},
			want:    big.NewInt(1),
			wantErr: false,
		},
		{
			name: "log10_10000",
			args: args{
				n: big.NewInt(10000),
			},
			want:    big.NewInt(4),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := log10(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("log10() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("log10() got = %v, want %v", got, tt.want)
			}
		})
	}
}
