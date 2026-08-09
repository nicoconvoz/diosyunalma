package pattern

import "testing"

func TestLagEqualityOnHandCheckedInput(t *testing.T) {
	// [1 2 2 4 2 4]
	//   lag 1 -> only positions 1,2 hold equal values
	//   lag 2 -> (g2,g4)=(2,2) and (g3,g5)=(4,4)
	//   lag 3 -> (g0,g3)=(1,4) (g1,g4)=(2,2) (g2,g5)=(2,4)
	gaps := []int{1, 2, 2, 4, 2, 4}

	tests := []struct {
		name string
		lag  int
		want int
	}{
		{name: "lag zero compares nothing", lag: 0, want: 0},
		{name: "negative lag compares nothing", lag: -2, want: 0},
		{name: "adjacent", lag: 1, want: 1},
		{name: "two apart", lag: 2, want: 2},
		{name: "three apart", lag: 3, want: 1},
		{name: "lag beyond the input", lag: 6, want: 0},
		{name: "lag far beyond the input", lag: 99, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LagEquality(gaps, tt.lag); got != tt.want {
				t.Errorf("LagEquality(%v, %d) = %d, want %d", gaps, tt.lag, got, tt.want)
			}
		})
	}
}

// TestLagEqualityAgreesWithShortPalindromes pins the two detectors together.
// A palindrome of 2 gaps is exactly an adjacent equal pair, and a palindrome of
// 3 gaps is exactly a pair equal across one element. If these ever disagree,
// one of the two detectors is wrong and every measurement built on them is
// suspect.
func TestLagEqualityAgreesWithShortPalindromes(t *testing.T) {
	inputs := [][]int{
		{1, 2, 2, 4, 2, 4, 2, 4, 6},
		{6, 6, 6, 6, 6},
		{2, 4, 6, 4, 2, 8, 2, 4, 6, 10, 2, 2},
		{},
		{4},
	}

	for _, gaps := range inputs {
		if lag, pal := LagEquality(gaps, 1), Palindromes(gaps, 2); lag != pal {
			t.Errorf("%v: LagEquality(.,1) = %d but Palindromes(.,2) = %d", gaps, lag, pal)
		}
		if lag, pal := LagEquality(gaps, 2), Palindromes(gaps, 3); lag != pal {
			t.Errorf("%v: LagEquality(.,2) = %d but Palindromes(.,3) = %d", gaps, lag, pal)
		}
	}
}

// TestLagEqualityCountsEveryPairOnConstantInput fixes the upper bound: with all
// values equal, every comparable position matches.
func TestLagEqualityCountsEveryPairOnConstantInput(t *testing.T) {
	gaps := []int{6, 6, 6, 6, 6, 6}

	for lag := 1; lag < len(gaps); lag++ {
		want := len(gaps) - lag
		if got := LagEquality(gaps, lag); got != want {
			t.Errorf("LagEquality(constant, %d) = %d, want %d", lag, got, want)
		}
	}
}
