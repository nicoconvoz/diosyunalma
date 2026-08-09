package primes

import (
	"reflect"
	"testing"
)

func TestSieveSmallValues(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		want  []int
	}{
		{name: "negative limit yields nothing", limit: -5, want: []int{}},
		{name: "zero yields nothing", limit: 0, want: []int{}},
		{name: "one is not prime", limit: 1, want: []int{}},
		{name: "two is the first prime", limit: 2, want: []int{2}},
		{name: "primes up to ten", limit: 10, want: []int{2, 3, 5, 7}},
		{
			name:  "primes up to thirty",
			limit: 30,
			want:  []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sieve(tt.limit)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Sieve(%d) = %v, want %v", tt.limit, got, tt.want)
			}
		})
	}
}

// TestSieveMatchesKnownPrimeCounts checks the output against published values of
// the prime-counting function pi(x). These are independently verified constants,
// not values produced by this implementation.
func TestSieveMatchesKnownPrimeCounts(t *testing.T) {
	tests := []struct {
		limit int
		want  int
	}{
		{limit: 10, want: 4},
		{limit: 100, want: 25},
		{limit: 1_000, want: 168},
		{limit: 10_000, want: 1_229},
		{limit: 100_000, want: 9_592},
		{limit: 1_000_000, want: 78_498},
	}

	for _, tt := range tests {
		got := len(Sieve(tt.limit))
		if got != tt.want {
			t.Errorf("pi(%d) = %d, want %d", tt.limit, got, tt.want)
		}
	}
}

// TestSieveIsSortedAndUnique guards the ordering contract that every downstream
// analysis (gaps, spirals, density) depends on.
func TestSieveIsSortedAndUnique(t *testing.T) {
	got := Sieve(10_000)

	for i := 1; i < len(got); i++ {
		if got[i] <= got[i-1] {
			t.Fatalf("output not strictly increasing at index %d: %d then %d", i, got[i-1], got[i])
		}
	}
}

// TestSieveBoundaryIsInclusive pins down that the limit itself is included when
// it happens to be prime — an off-by-one here would silently corrupt every count.
func TestSieveBoundaryIsInclusive(t *testing.T) {
	got := Sieve(7)
	want := []int{2, 3, 5, 7}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sieve(7) = %v, want %v", got, want)
	}
}
