package accumulator

import "testing"

func TestAccAndProveIterParallel(t *testing.T) {

	type args struct {
		set        []string
		encodeType EncodeType
		setup      *Setup
	}
	tests := []struct {
		name         string
		args         args
		idx          int
		wantProofLen int
	}{
		{
			name: "set_size_16",
			args: args{
				set:        GenTestSet(16),
				encodeType: HashToPrimeFromSha256,
				setup:      TrustedSetup(),
			},
			idx:          5,
			wantProofLen: 16,
		},
		{
			name: "set_size_17",
			args: args{
				set:        GenTestSet(17),
				encodeType: HashToPrimeFromSha256,
				setup:      TrustedSetup(),
			},
			idx:          7,
			wantProofLen: 17,
		},
		{
			name: "set_size_254",
			args: args{
				set:        GenTestSet(254),
				encodeType: HashToPrimeFromSha256,
				setup:      TrustedSetup(),
			},
			idx:          253,
			wantProofLen: 254,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, proofs := AccAndProveIterParallel(tt.args.set, tt.args.encodeType, tt.args.setup)
			acc1, acc2 := genAccts(tt.args.set, tt.args.setup, proofs, tt.idx)
			if len(proofs) != tt.wantProofLen {
				t.Errorf("AccAndProveParallel() got proof len = %v, want %v", len(proofs), tt.wantProofLen)
				return
			}
			if acc.Cmp(acc2) != 0 {
				t.Errorf("proofs generated are not consistent acc = %v, acc2 %v", acc, acc2)
				return
			}
			if acc1.Cmp(acc2) != 0 {
				t.Errorf("proofs generated are not consistent acc1 = %v, acc2 %v", acc1, acc2)
				return
			}
		})
	}
}
