package riemann

import (
	"math"
	"testing"
)

// TestZeroCountAgainstTheKnownCensus checks Riemann–von Mangoldt against the
// actual zero counts: there are exactly 29 zeros below height 100 and 10 below
// height 50. The formula omits the small fluctuation term S(T), so the second
// value carries a wider tolerance.
func TestZeroCountAgainstTheKnownCensus(t *testing.T) {
	if got := ZeroCount(100); math.Abs(got-29) > 0.2 {
		t.Errorf("ZeroCount(100) = %v, want 29 within 0.2", got)
	}
	if got := ZeroCount(50); math.Abs(got-10) > 0.7 {
		t.Errorf("ZeroCount(50) = %v, want 10 within 0.7", got)
	}
}

// TestZeroCountBeforeTheFirstZeroIsSmall pins the low end: below γ₁ = 14.13
// the true count is zero, and the smooth formula must stay under one.
func TestZeroCountBeforeTheFirstZeroIsSmall(t *testing.T) {
	if got := ZeroCount(14.0); got < -0.5 || got > 1 {
		t.Errorf("ZeroCount(14) = %v, want between -0.5 and 1", got)
	}
}

func TestZeroCountIsMonotoneAboveTheFirstZero(t *testing.T) {
	prev := ZeroCount(15)
	for x := 16.0; x <= 200; x++ {
		cur := ZeroCount(x)
		if cur < prev {
			t.Fatalf("ZeroCount decreased between %v and %v: %v -> %v", x-1, x, prev, cur)
		}
		prev = cur
	}
}

func TestZeroCountHandlesDegenerateInput(t *testing.T) {
	if got := ZeroCount(0); got != 0 {
		t.Errorf("ZeroCount(0) = %v, want 0", got)
	}
	if got := ZeroCount(-5); got != 0 {
		t.Errorf("ZeroCount(-5) = %v, want 0", got)
	}
}
