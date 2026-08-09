package riemann

import (
	"math"
	"testing"
)

// TestPrimeClockMatchesTheHandFormula pins the evaluation to the truncated
// explicit formula: with a single zero the clock must read exactly
// 1 − 2·cos(γ·ln x)/√x.
func TestPrimeClockMatchesTheHandFormula(t *testing.T) {
	const gamma, x = 14.134725, 10.0
	want := 1 - 2*math.Cos(gamma*math.Log(x))/math.Sqrt(x)

	got := PrimeClock([]float64{gamma}, x)
	if math.Abs(got-want) > 1e-14 {
		t.Errorf("PrimeClock = %v, want %v", got, want)
	}
}

// TestPrimeClockWithNoHandsReadsTheSmoothDensity fixes the baseline: no zeros
// means no oscillation, and the clock face shows only the smooth part.
func TestPrimeClockWithNoHandsReadsTheSmoothDensity(t *testing.T) {
	if got := PrimeClock(nil, 50); got != 1 {
		t.Errorf("PrimeClock(no zeros) = %v, want 1", got)
	}
}

// TestPrimeClockHandsAdd checks linearity: two hands contribute the sum of
// their separate oscillations.
func TestPrimeClockHandsAdd(t *testing.T) {
	const x = 7.5
	a, b := 14.134725, 21.022040

	sum := PrimeClock([]float64{a}, x) + PrimeClock([]float64{b}, x) - 1
	got := PrimeClock([]float64{a, b}, x)

	if math.Abs(got-sum) > 1e-12 {
		t.Errorf("two hands = %v, want %v", got, sum)
	}
}

// TestPrimeClockAlignsOnThePrimePowers is the behavioural claim the sundial
// rests on. Ten hands give a resolution of about x·π/50, so the assertion
// stays where that resolves: at small x the clock must read clearly higher on
// prime powers than between them. (An earlier version asserted x = 29 against
// x = 30 and failed — correctly: at x ≈ 30 the ten-hand resolution is ~1.9, so
// neighbouring integers blur. The test was wrong, not the formula.)
func TestPrimeClockAlignsOnThePrimePowers(t *testing.T) {
	trueZeros := []float64{
		14.134725, 21.022040, 25.010858, 30.424876, 32.935062,
		37.586178, 40.918719, 43.327073, 48.005151, 49.773832,
	}

	primePowers := []float64{2, 3, 4, 5, 7, 8, 9, 11, 13}
	between := []float64{2.5, 3.5, 4.5, 6, 7.5, 8.5, 10, 12, 14.5}

	on, off := 0.0, 0.0
	for _, x := range primePowers {
		on += PrimeClock(trueZeros, x)
	}
	for _, x := range between {
		off += PrimeClock(trueZeros, x)
	}
	on /= float64(len(primePowers))
	off /= float64(len(between))

	if on < off+2 {
		t.Errorf("clock mean on prime powers = %v, between = %v; want a clear separation", on, off)
	}

	// Two sharp pointwise cases well inside the resolution.
	if PrimeClock(trueZeros, 2) <= PrimeClock(trueZeros, 2.5) {
		t.Error("clock at 2 should exceed clock at 2.5")
	}
	if PrimeClock(trueZeros, 13) <= PrimeClock(trueZeros, 12) {
		t.Error("clock at 13 should exceed clock at 12")
	}
}

func TestPrimeClockHandlesDegenerateInput(t *testing.T) {
	if got := PrimeClock([]float64{14.13}, 1); got != 0 {
		t.Errorf("x=1 = %v, want 0 (ln 1 = 0 makes the clock meaningless)", got)
	}
	if got := PrimeClock([]float64{14.13}, -5); got != 0 {
		t.Errorf("negative x = %v, want 0", got)
	}
}
