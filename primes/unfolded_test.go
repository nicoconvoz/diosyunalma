package primes

import (
	"math"
	"testing"
)

// TestUnfoldedMatchesTheHandFormula pins the rescaling: each gap is divided by
// the local mean gap ln(p), which is the unfolding Finding 7 used.
func TestUnfoldedMatchesTheHandFormula(t *testing.T) {
	got := Unfolded([]int{3, 5, 11})
	want := []float64{2 / math.Log(3), 6 / math.Log(5)}

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if math.Abs(got[i]-want[i]) > 1e-14 {
			t.Errorf("spacing %d = %v, want %v", i, got[i], want[i])
		}
	}
}

// TestUnfoldedMeanIsNearOne is the property that makes the rescaling an
// unfolding at all: measured in units of the local mean gap, the average
// spacing must sit at 1.
func TestUnfoldedMeanIsNearOne(t *testing.T) {
	s := Unfolded(Sieve(1_000_000))

	sum := 0.0
	for _, v := range s {
		sum += v
	}
	mean := sum / float64(len(s))

	if mean < 0.97 || mean > 1.03 {
		t.Errorf("mean unfolded spacing = %v, want within 3%% of 1", mean)
	}
}

func TestUnfoldedHandlesDegenerateInput(t *testing.T) {
	if got := Unfolded(nil); len(got) != 0 {
		t.Errorf("nil = %v, want empty", got)
	}
	if got := Unfolded([]int{7}); len(got) != 0 {
		t.Errorf("single value = %v, want empty", got)
	}
}
