package control

import (
	"math"
	"testing"
)

func TestEvaluateComputesTheDecoyStatistics(t *testing.T) {
	// A sampler alternating 8 and 12 has mean 10 and population stddev 2.
	values := []int{8, 12, 8, 12}
	i := 0
	sample := func() int {
		v := values[i%len(values)]
		i++
		return v
	}

	got := Evaluate(20, len(values), sample)

	if got.Observed != 20 {
		t.Errorf("Observed = %d, want 20", got.Observed)
	}
	if math.Abs(got.Mean-10) > 1e-9 {
		t.Errorf("Mean = %v, want 10", got.Mean)
	}
	if math.Abs(got.StdDev-2) > 1e-9 {
		t.Errorf("StdDev = %v, want 2", got.StdDev)
	}
	if math.Abs(got.Ratio-2) > 1e-9 {
		t.Errorf("Ratio = %v, want 2", got.Ratio)
	}
	if math.Abs(got.ZScore-5) > 1e-9 {
		t.Errorf("ZScore = %v, want 5", got.ZScore)
	}
}

// TestEvaluateStaysFiniteWhenDecoysNeverVary is the guard that matters in
// practice: a constant decoy makes the standard deviation zero, and a naive
// z-score would divide by it and report Inf or NaN as if it were a discovery.
func TestEvaluateStaysFiniteWhenDecoysNeverVary(t *testing.T) {
	got := Evaluate(50, 10, func() int { return 7 })

	if math.IsInf(got.ZScore, 0) || math.IsNaN(got.ZScore) {
		t.Fatalf("ZScore = %v, want a finite value", got.ZScore)
	}
	if got.ZScore != 0 {
		t.Errorf("ZScore = %v, want 0 when decoys never vary", got.ZScore)
	}
	if got.Significant() {
		t.Error("Significant() = true with zero spread; a flat control proves nothing")
	}
}

func TestEvaluateWithoutTrialsClaimsNothing(t *testing.T) {
	got := Evaluate(100, 0, func() int { return 1 })

	if got.ZScore != 0 || got.Ratio != 0 || got.Significant() {
		t.Errorf("Evaluate with 0 trials = %+v, want a null result", got)
	}
}

func TestSignificantNeedsFiveSigmaInEitherDirection(t *testing.T) {
	tests := []struct {
		name string
		z    float64
		want bool
	}{
		{name: "noise", z: 1.5, want: false},
		{name: "suggestive but not enough", z: 4.9, want: false},
		{name: "clear excess", z: 44, want: true},
		{name: "clear deficit", z: -44, want: true},
		{name: "exactly at the line is not past it", z: 5, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{ZScore: tt.z, StdDev: 1}
			if got := r.Significant(); got != tt.want {
				t.Errorf("Result{ZScore: %v}.Significant() = %v, want %v", tt.z, got, tt.want)
			}
		})
	}
}
