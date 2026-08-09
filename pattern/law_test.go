package pattern

import (
	"math"
	"testing"
)

// TestLawViolationsCountsTheForbiddenConfiguration exercises the rule directly:
// a palindromic window may hold an odd number of gaps that are not multiples of
// 3, or none at all. An even non-zero count is the forbidden case.
func TestLawViolationsCountsTheForbiddenConfiguration(t *testing.T) {
	tests := []struct {
		name string
		gaps []int
		k    int
		want int
	}{
		{name: "two equal multiples of three are allowed", gaps: []int{6, 6}, k: 2, want: 0},
		{name: "two equal non-multiples are forbidden", gaps: []int{2, 2}, k: 2, want: 1},
		{name: "one flip at the centre is allowed", gaps: []int{6, 4, 6}, k: 3, want: 0},
		{name: "two flips on the outside are forbidden", gaps: []int{4, 6, 4}, k: 3, want: 1},
		{name: "three flips are allowed", gaps: []int{4, 4, 4}, k: 3, want: 0},
		{name: "non-palindromes are not judged", gaps: []int{2, 4}, k: 2, want: 0},
		{name: "window larger than the input", gaps: []int{6, 6}, k: 9, want: 0},
		{name: "window below two", gaps: []int{6, 6}, k: 1, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LawViolations(tt.gaps, tt.k); got != tt.want {
				t.Errorf("LawViolations(%v, %d) = %d, want %d", tt.gaps, tt.k, got, tt.want)
			}
		})
	}
}

// TestLawViolationsCountsEveryOffendingWindow makes sure overlapping windows are
// each judged rather than the scan stopping at the first hit.
func TestLawViolationsCountsEveryOffendingWindow(t *testing.T) {
	gaps := []int{2, 2, 2, 2}

	if got := LawViolations(gaps, 2); got != 3 {
		t.Errorf("LawViolations(%v, 2) = %d, want 3", gaps, got)
	}
}

// TestSingularBoostIsExactlyOneForPureTwoThreeGaps is a regression test.
//
// The first implementation of the prime factorisation behind this failed to
// strip the factors of 2 and 3 before collecting the product over primes above
// 3, so a gap of 6 was treated as having 6 itself as a factor and returned
// 1.0667 instead of 1. Every boost computed in that pass was wrong.
func TestSingularBoostIsExactlyOneForPureTwoThreeGaps(t *testing.T) {
	for _, d := range []int{2, 3, 4, 6, 8, 12, 18, 24, 36, 48, 54, 72, 96, 108, 162} {
		if got := SingularBoost(d); got != 1.0 {
			t.Errorf("SingularBoost(%d) = %v, want exactly 1", d, got)
		}
	}
}

func TestSingularBoostAppliesOneFactorPerPrimeAboveThree(t *testing.T) {
	// (q-2)^2 / ((q-1)(q-3)) per distinct prime q > 3 dividing d.
	five := 9.0 / (4.0 * 2.0)   // 1.125
	seven := 25.0 / (6.0 * 4.0) // 1.041666...

	tests := []struct {
		d    int
		want float64
	}{
		{d: 30, want: five},          // 2·3·5
		{d: 60, want: five},          // 2²·3·5
		{d: 150, want: five},         // repeated 5 counts once
		{d: 42, want: seven},         // 2·3·7
		{d: 210, want: five * seven}, // 2·3·5·7
	}

	for _, tt := range tests {
		got := SingularBoost(tt.d)
		if math.Abs(got-tt.want) > 1e-12 {
			t.Errorf("SingularBoost(%d) = %v, want %v", tt.d, got, tt.want)
		}
	}
}

// TestSingularBoostGrowsWithSmallPrimeFactors states the qualitative claim the
// measurement rests on: a gap carrying more small primes is favoured more.
func TestSingularBoostGrowsWithSmallPrimeFactors(t *testing.T) {
	if SingularBoost(30) <= SingularBoost(6) {
		t.Error("a gap divisible by 5 must be boosted above a pure 2-3 gap")
	}
	if SingularBoost(30) <= SingularBoost(42) {
		t.Error("the boost from 5 must exceed the boost from 7")
	}
	if SingularBoost(210) <= SingularBoost(30) {
		t.Error("two extra prime factors must beat one")
	}
}

func TestSingularBoostHandlesDegenerateInput(t *testing.T) {
	for _, d := range []int{0, 1, -6} {
		if got := SingularBoost(d); got != 1.0 {
			t.Errorf("SingularBoost(%d) = %v, want 1", d, got)
		}
	}
}
