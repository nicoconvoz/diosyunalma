package primes

import "testing"

func TestSetKnowsWhichNumbersArePrime(t *testing.T) {
	s := NewSet(20)

	for _, p := range []int{2, 3, 5, 7, 11, 13, 17, 19} {
		if !s.Has(p) {
			t.Errorf("Has(%d) = false, want true", p)
		}
	}
	for _, n := range []int{0, 1, 4, 6, 8, 9, 10, 12, 15, 16, 18, 20} {
		if s.Has(n) {
			t.Errorf("Has(%d) = true, want false", n)
		}
	}
}

func TestSetIsSafeOutsideItsRange(t *testing.T) {
	s := NewSet(20)

	for _, n := range []int{-1, -100, 21, 1000} {
		if s.Has(n) {
			t.Errorf("Has(%d) = true, want false for out-of-range input", n)
		}
	}
}

func TestSetAgreesWithSieve(t *testing.T) {
	const limit = 10_000
	s := NewSet(limit)

	for _, p := range Sieve(limit) {
		if !s.Has(p) {
			t.Fatalf("Sieve says %d is prime, Set disagrees", p)
		}
	}

	count := 0
	for n := 0; n <= limit; n++ {
		if s.Has(n) {
			count++
		}
	}
	if want := len(Sieve(limit)); count != want {
		t.Errorf("Set holds %d primes, Sieve found %d", count, want)
	}
}

// TestProgressionTriplesCountsAllPrimeTriples measures the quantity
// Hardy–Littlewood actually predicts: p, p+d and p+2d all prime, with no
// requirement that they be consecutive. The difference between this and the
// consecutive count is the whole open question.
func TestProgressionTriplesCountsAllPrimeTriples(t *testing.T) {
	s := NewSet(20) // primes: 2 3 5 7 11 13 17 19

	tests := []struct {
		name string
		d    int
		want int
	}{
		{name: "d=2 finds only 3,5,7", d: 2, want: 1},
		{name: "d=4 finds only 3,7,11", d: 4, want: 1},
		{name: "d=6 finds 5,11,17 and 7,13,19", d: 6, want: 2},
		{name: "d beyond the range finds nothing", d: 50, want: 0},
		{name: "d below one is not a progression", d: 0, want: 0},
		{name: "negative d is not a progression", d: -6, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProgressionTriples(s, tt.d); got != tt.want {
				t.Errorf("ProgressionTriples(d=%d) = %d, want %d", tt.d, got, tt.want)
			}
		})
	}
}

// TestProgressionTriplesNeverUndercountsConsecutiveOnes pins the containment
// the comparison rests on: every pair of equal consecutive gaps is also an
// arithmetic triple, so the unconstrained count can only be larger.
func TestProgressionTriplesNeverUndercountsConsecutiveOnes(t *testing.T) {
	const limit = 100_000
	s := NewSet(limit)
	gaps := Gaps(From(Sieve(limit), 5))

	for _, d := range []int{6, 12, 18, 24, 30} {
		consecutive := 0
		for i := 1; i < len(gaps); i++ {
			if gaps[i-1] == d && gaps[i] == d {
				consecutive++
			}
		}
		all := ProgressionTriples(s, d)

		if all < consecutive {
			t.Errorf("d=%d: %d triples but %d consecutive pairs; the unconstrained count cannot be smaller",
				d, all, consecutive)
		}
	}
}
