package primes

import (
	"math"
	"testing"
)

// TestPsiAtOnHandComputedValues checks Chebyshev's psi against values worked
// out from the definition: psi(x) sums ln p over every prime POWER up to x.
//
//	psi(10) = 3·ln2 + 2·ln3 + ln5 + ln7      (2,4,8 | 3,9 | 5 | 7)
//	psi(30) = 4·ln2 + 3·ln3 + 2·ln5 + ln(7·11·13·17·19·23·29)
func TestPsiAtOnHandComputedValues(t *testing.T) {
	ln := math.Log

	psi10 := 3*ln(2) + 2*ln(3) + ln(5) + ln(7)
	psi30 := 4*ln(2) + 3*ln(3) + 2*ln(5) +
		ln(7) + ln(11) + ln(13) + ln(17) + ln(19) + ln(23) + ln(29)

	got := PsiAt(30, []int{1, 2, 10, 29, 30})
	want := []float64{0, ln(2), psi10, psi30, psi30}

	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if math.Abs(got[i]-want[i]) > 1e-12 {
			t.Errorf("psi(%v) = %v, want %v", []int{1, 2, 10, 29, 30}[i], got[i], want[i])
		}
	}
}

// TestPsiAtIsMonotone pins that psi never decreases: each sample point can only
// accumulate more prime powers.
func TestPsiAtIsMonotone(t *testing.T) {
	xs := []int{10, 100, 1_000, 10_000, 100_000}
	got := PsiAt(100_000, xs)

	for i := 1; i < len(got); i++ {
		if got[i] < got[i-1] {
			t.Fatalf("psi decreased from %v to %v between %d and %d",
				got[i-1], got[i], xs[i-1], xs[i])
		}
	}
}

// TestPsiAtTracksX records the prime number theorem in miniature: psi(x)/x
// approaches 1, and at x = 10^5 it is already within a few percent.
func TestPsiAtTracksX(t *testing.T) {
	const x = 100_000
	got := PsiAt(x, []int{x})[0]

	ratio := got / float64(x)
	if ratio < 0.98 || ratio > 1.02 {
		t.Errorf("psi(%d)/%d = %v, want within 2%% of 1", x, x, ratio)
	}
}

func TestPsiAtHandlesDegenerateInput(t *testing.T) {
	if got := PsiAt(100, nil); len(got) != 0 {
		t.Errorf("no sample points = %v, want empty", got)
	}
	got := PsiAt(10, []int{0, 1})
	for i, v := range got {
		if v != 0 {
			t.Errorf("psi below 2 at index %d = %v, want 0", i, v)
		}
	}
}
