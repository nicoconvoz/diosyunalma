package primes

import (
	"reflect"
	"testing"
)

func TestFrom(t *testing.T) {
	seq := []int{2, 3, 5, 7, 11, 13}

	tests := []struct {
		name string
		min  int
		want []int
	}{
		{name: "drops the primes the mod-3 walk cannot cover", min: 5, want: []int{5, 7, 11, 13}},
		{name: "keeps everything when min is below the first", min: 2, want: []int{2, 3, 5, 7, 11, 13}},
		{name: "boundary is inclusive", min: 13, want: []int{13}},
		{name: "past the end yields empty", min: 99, want: []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := From(seq, tt.min); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("From(%v, %d) = %v, want %v", seq, tt.min, got, tt.want)
			}
		})
	}
}

func TestFromDoesNotAliasTheInput(t *testing.T) {
	seq := []int{2, 3, 5, 7}
	got := From(seq, 5)
	got[0] = 999

	if seq[2] != 5 {
		t.Errorf("From aliased its input: seq = %v", seq)
	}
}

// TestFlipsMarksTheGapsThatChangeResidue records the fact the whole mod-3
// analysis rests on: a gap divisible by the modulus leaves the residue where it
// was, any other gap moves it.
func TestFlipsMarksTheGapsThatChangeResidue(t *testing.T) {
	tests := []struct {
		name    string
		gaps    []int
		modulus int
		want    []int
	}{
		{
			name:    "multiples of three stay put",
			gaps:    []int{6, 4, 12, 2},
			modulus: 3,
			want:    []int{0, 1, 0, 1},
		},
		{
			name:    "every gap past the first is even, so mod 2 never flips",
			gaps:    []int{2, 4, 6, 8},
			modulus: 2,
			want:    []int{0, 0, 0, 0},
		},
		{
			name:    "empty input",
			gaps:    []int{},
			modulus: 3,
			want:    []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flips(tt.gaps, tt.modulus)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flips(%v, %d) = %v, want %v", tt.gaps, tt.modulus, got, tt.want)
			}
		})
	}
}

func TestFlipsRejectsAModulusBelowTwo(t *testing.T) {
	for _, m := range []int{1, 0, -3} {
		if got := Flips([]int{6, 4}, m); len(got) != 0 {
			t.Errorf("Flips(gaps, %d) = %v, want empty", m, got)
		}
	}
}

// TestTransitionMatrixCountsResidueSteps checks the counting against a hand
// worked example: 5,7,11,13 have residues 2,1,2,1 mod 3, giving the steps
// 2->1, 1->2 and 2->1.
func TestTransitionMatrixCountsResidueSteps(t *testing.T) {
	got := TransitionMatrix([]int{5, 7, 11, 13}, 3)

	want := [][]int{
		{0, 0, 0},
		{0, 0, 1},
		{0, 2, 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TransitionMatrix = %v, want %v", got, want)
	}
}

func TestTransitionMatrixTotalsMatchTheStepCount(t *testing.T) {
	seq := From(Sieve(10_000), 5)
	got := TransitionMatrix(seq, 3)

	total := 0
	for _, row := range got {
		for _, c := range row {
			total += c
		}
	}
	if want := len(seq) - 1; total != want {
		t.Errorf("transitions total = %d, want %d", total, want)
	}

	// Primes above 3 never land in residue 0, so that row and column stay empty.
	for i := range got {
		if got[0][i] != 0 || got[i][0] != 0 {
			t.Errorf("residue 0 is populated at index %d: %v", i, got)
		}
	}
}

func TestTransitionMatrixHandlesDegenerateInput(t *testing.T) {
	if got := TransitionMatrix([]int{7}, 3); len(got) != 3 {
		t.Errorf("single element should still return a %dx%d matrix, got %v", 3, 3, got)
	}
	if got := TransitionMatrix([]int{5, 7}, 1); got != nil {
		t.Errorf("TransitionMatrix with modulus 1 = %v, want nil", got)
	}
}
