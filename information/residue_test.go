package information

import (
	"math"
	"testing"
)

func TestResidueEntropyOnKnownDistributions(t *testing.T) {
	tests := []struct {
		name    string
		values  []int
		modulus int
		want    float64
	}{
		{
			name:    "uniform residues carry full entropy",
			values:  []int{5, 6, 7, 8, 9, 10, 11, 12},
			modulus: 4,
			want:    1,
		},
		{
			name:    "a single residue class carries none",
			values:  []int{7, 14, 21, 28},
			modulus: 7,
			want:    0,
		},
		{
			name:    "odd primes mod 2 collapse to one class",
			values:  []int{3, 5, 7, 11, 13, 17, 19},
			modulus: 2,
			want:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResidueEntropy(tt.values, tt.modulus)
			if math.Abs(got-tt.want) > 1e-12 {
				t.Errorf("ResidueEntropy = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestResidueEntropySkipsValuesUpToTheModulus pins the convention Finding 4
// used: the primes dividing or equal to the modulus are structural exceptions,
// not samples of the walk.
func TestResidueEntropySkipsValuesUpToTheModulus(t *testing.T) {
	with := ResidueEntropy([]int{2, 3, 5, 7, 11, 13}, 6)
	without := ResidueEntropy([]int{7, 11, 13}, 6)

	if math.Abs(with-without) > 1e-12 {
		t.Errorf("values <= modulus changed the result: %v vs %v", with, without)
	}
}

func TestResidueEntropyHandlesDegenerateInput(t *testing.T) {
	if got := ResidueEntropy(nil, 6); got != 0 {
		t.Errorf("no values = %v, want 0", got)
	}
	if got := ResidueEntropy([]int{7, 11}, 1); got != 0 {
		t.Errorf("modulus 1 = %v, want 0", got)
	}
	if got := ResidueEntropy([]int{2, 3, 5}, 10); got != 0 {
		t.Errorf("all values below modulus = %v, want 0", got)
	}
}
