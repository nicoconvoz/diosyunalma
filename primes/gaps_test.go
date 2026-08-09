package primes

import (
	"reflect"
	"testing"
)

func TestGaps(t *testing.T) {
	tests := []struct {
		name string
		seq  []int
		want []int
	}{
		{name: "nil sequence has no gaps", seq: nil, want: []int{}},
		{name: "single element has no gaps", seq: []int{7}, want: []int{}},
		{name: "two elements yield one gap", seq: []int{2, 3}, want: []int{1}},
		{
			name: "primes up to thirty",
			seq:  []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29},
			want: []int{1, 2, 2, 4, 2, 4, 2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Gaps(tt.seq)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Gaps(%v) = %v, want %v", tt.seq, got, tt.want)
			}
		})
	}
}

// TestGapsSumsToSpan pins the invariant that every downstream measurement
// leans on: the gaps must reconstruct the distance between the endpoints.
func TestGapsSumsToSpan(t *testing.T) {
	seq := Sieve(1000)
	got := Gaps(seq)

	if len(got) != len(seq)-1 {
		t.Fatalf("len(Gaps) = %d, want %d", len(got), len(seq)-1)
	}

	sum := 0
	for _, g := range got {
		sum += g
	}
	if want := seq[len(seq)-1] - seq[0]; sum != want {
		t.Errorf("sum of gaps = %d, want %d", sum, want)
	}
}

// TestGapsAreEvenAfterTheFirst records the arithmetic fact that makes the
// palindrome parity law possible: every prime past 2 is odd, so every gap past
// the first is even.
func TestGapsAreEvenAfterTheFirst(t *testing.T) {
	got := Gaps(Sieve(10_000))

	if got[0] != 1 {
		t.Fatalf("first gap = %d, want 1 (from 2 to 3)", got[0])
	}
	for i, g := range got[1:] {
		if g%2 != 0 {
			t.Errorf("gap at index %d is %d, want even", i+1, g)
		}
	}
}
