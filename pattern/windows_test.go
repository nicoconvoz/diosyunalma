package pattern

import (
	"reflect"
	"testing"
)

func TestWindowsReturnsTheMatchingWindowsInOrder(t *testing.T) {
	// [6 4 6 2 6] — the k=3 palindromes are [6 4 6] at 0 and [6 2 6] at 2.
	gaps := []int{6, 4, 6, 2, 6}
	keepAll := func([]int) bool { return true }

	got := Windows(gaps, 3, keepAll)
	want := [][]int{{6, 4, 6}, {6, 2, 6}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Windows(%v, 3) = %v, want %v", gaps, got, want)
	}
}

func TestWindowsHonoursTheFilter(t *testing.T) {
	gaps := []int{6, 4, 6, 2, 6}
	centreIsTwo := func(w []int) bool { return w[len(w)/2] == 2 }

	got := Windows(gaps, 3, centreIsTwo)
	want := [][]int{{6, 2, 6}}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Windows filtered = %v, want %v", got, want)
	}
}

// TestWindowsAgreesWithPalindromesWith pins the collector to the counter. If
// they ever disagree, every pattern tally built on Windows is suspect.
func TestWindowsAgreesWithPalindromesWith(t *testing.T) {
	inputs := [][]int{
		{1, 2, 2, 4, 2, 4, 2, 4, 6},
		{6, 6, 6, 6, 6},
		{2, 4, 6, 4, 2, 8, 2, 4, 6, 10, 2, 2},
		{},
	}
	keep := func(w []int) bool { return w[0]%4 == 2 }

	for _, gaps := range inputs {
		for k := 2; k <= 5; k++ {
			collected := len(Windows(gaps, k, keep))
			counted := PalindromesWith(gaps, k, keep)
			if collected != counted {
				t.Errorf("%v k=%d: Windows found %d but PalindromesWith counted %d",
					gaps, k, collected, counted)
			}
		}
	}
}

// TestWindowsReturnsCopies guards the tally use case: mutating a returned
// window must not corrupt the source sequence or later calls.
func TestWindowsReturnsCopies(t *testing.T) {
	gaps := []int{6, 6, 6}
	keep := func([]int) bool { return true }

	first := Windows(gaps, 2, keep)
	first[0][0] = 999

	if gaps[0] != 6 {
		t.Fatalf("mutating a result corrupted the input: %v", gaps)
	}
	if again := Windows(gaps, 2, keep); again[0][0] != 6 {
		t.Errorf("a later call sees the mutation: %v", again)
	}
}

func TestWindowsHandlesDegenerateInput(t *testing.T) {
	keep := func([]int) bool { return true }

	if got := Windows([]int{6, 6}, 1, keep); len(got) != 0 {
		t.Errorf("k below two = %v, want empty", got)
	}
	if got := Windows([]int{6, 6}, 9, keep); len(got) != 0 {
		t.Errorf("k beyond the input = %v, want empty", got)
	}
	if got := Windows([]int{6, 6}, 2, nil); len(got) != 0 {
		t.Errorf("nil filter = %v, want empty", got)
	}
}
