package primes

import (
	"math"
	"testing"
)

// TestPsiSegmentedMatchesPsiAt pins the segmented evaluation to the direct
// one on a range where both fit in memory. Any drift between them would
// poison every measurement made through the telescope.
func TestPsiSegmentedMatchesPsiAt(t *testing.T) {
	xs := []int{2, 10, 97, 1000, 65_536, 999_983, 1_000_000}

	direct := PsiAt(1_000_000, xs)
	segmented := PsiSegmented(1_000_000, xs)

	if len(segmented) != len(direct) {
		t.Fatalf("len = %d, want %d", len(segmented), len(direct))
	}
	for i := range direct {
		diff := math.Abs(segmented[i] - direct[i])
		if diff > 1e-6 {
			t.Errorf("psi(%d): segmented %v vs direct %v (diff %v)",
				xs[i], segmented[i], direct[i], diff)
		}
	}
}

// TestPsiSegmentedCrossesSegmentBoundaries uses a tiny segment size indirectly:
// sample points straddling powers of two exercise the boundary bookkeeping.
func TestPsiSegmentedCrossesSegmentBoundaries(t *testing.T) {
	xs := []int{}
	for x := 100; x <= 100_000; x += 997 {
		xs = append(xs, x)
	}

	direct := PsiAt(100_000, xs)
	segmented := PsiSegmented(100_000, xs)

	for i := range direct {
		if math.Abs(segmented[i]-direct[i]) > 1e-7 {
			t.Fatalf("psi(%d): segmented %v vs direct %v", xs[i], segmented[i], direct[i])
		}
	}
}

func TestPsiSegmentedHandlesDegenerateInput(t *testing.T) {
	if got := PsiSegmented(1000, nil); len(got) != 0 {
		t.Errorf("no sample points = %v, want empty", got)
	}
	got := PsiSegmented(100, []int{0, 1})
	for i, v := range got {
		if v != 0 {
			t.Errorf("psi below 2 at index %d = %v, want 0", i, v)
		}
	}
}
