package pattern

import "testing"

func TestPalindromesOnHandCheckedWindows(t *testing.T) {
	// Worked out by hand from [1 2 2 4 2 4]:
	//   k=2 -> only [2 2]
	//   k=3 -> [2 4 2] and [4 2 4]
	//   k=4 -> none
	gaps := []int{1, 2, 2, 4, 2, 4}

	tests := []struct {
		name string
		k    int
		want int
	}{
		{name: "k below two is not a window", k: 1, want: 0},
		{name: "k zero is not a window", k: 0, want: 0},
		{name: "negative k is not a window", k: -3, want: 0},
		{name: "pairs of equal adjacent gaps", k: 2, want: 1},
		{name: "triples mirrored around a centre", k: 3, want: 2},
		{name: "quadruples", k: 4, want: 0},
		{name: "window longer than the input", k: 99, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Palindromes(gaps, tt.k); got != tt.want {
				t.Errorf("Palindromes(%v, %d) = %d, want %d", gaps, tt.k, got, tt.want)
			}
		})
	}
}

// TestPalindromesFindsTheOnlyBalancedPrimeBelowThirty ties the detector to a
// verifiable arithmetic fact: 3-5-7 is the single place below 30 where a prime
// sits exactly halfway between its neighbours, because a repeated gap that is
// not a multiple of 3 forces one of the three terms to be divisible by 3.
func TestPalindromesFindsTheOnlyBalancedPrimeBelowThirty(t *testing.T) {
	gaps := []int{1, 2, 2, 4, 2, 4, 2, 4, 6} // primes up to 29

	if got := Palindromes(gaps, 2); got != 1 {
		t.Errorf("Palindromes(primes<30, 2) = %d, want 1 (only 3-5-7)", got)
	}
}

// TestPalindromesEveryWindowOfEqualGaps is the degenerate case: a constant
// sequence is a palindrome at every position and every length.
func TestPalindromesEveryWindowOfEqualGaps(t *testing.T) {
	gaps := []int{6, 6, 6, 6, 6}

	for k := 2; k <= len(gaps); k++ {
		want := len(gaps) - k + 1
		if got := Palindromes(gaps, k); got != want {
			t.Errorf("Palindromes(constant, %d) = %d, want %d", k, got, want)
		}
	}
}

// TestPalindromesIgnoresDirection states the property the whole experiment
// rests on: a window and its reverse must be scored identically.
func TestPalindromesIgnoresDirection(t *testing.T) {
	gaps := []int{2, 4, 6, 4, 2, 8, 2, 4, 6}

	reversed := make([]int, len(gaps))
	for i, g := range gaps {
		reversed[len(gaps)-1-i] = g
	}

	for k := 2; k <= 5; k++ {
		forward := Palindromes(gaps, k)
		backward := Palindromes(reversed, k)
		if forward != backward {
			t.Errorf("k=%d: forward = %d, backward = %d; must match", k, forward, backward)
		}
	}
}
