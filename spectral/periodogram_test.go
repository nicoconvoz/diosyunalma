package spectral

import (
	"math"
	"testing"
)

func grid(n int, du float64) []float64 {
	u := make([]float64, n)
	for i := range u {
		u[i] = float64(i) * du
	}
	return u
}

func argmax(v []float64) int {
	best := 0
	for i, x := range v {
		if x > v[best] {
			best = i
		}
	}
	return best
}

// TestPeriodogramFindsAPureTone is the core contract: a clean cosine must put
// its power at its own frequency, not somewhere nearby.
func TestPeriodogramFindsAPureTone(t *testing.T) {
	const f0 = 14.0
	u := grid(2000, 0.01) // 20 units of u — plenty of cycles at f0
	y := make([]float64, len(u))
	for i, ui := range u {
		y[i] = math.Cos(f0 * ui)
	}

	freqs := make([]float64, 0)
	for f := 5.0; f <= 30.0; f += 0.01 {
		freqs = append(freqs, f)
	}

	power := Periodogram(u, y, freqs)
	found := freqs[argmax(power)]

	if math.Abs(found-f0) > 0.05 {
		t.Errorf("peak at %v, want %v", found, f0)
	}
}

// TestPeriodogramSeparatesTwoTonesByAmplitude checks that a stronger tone
// carries more power than a weaker one, and both are located.
func TestPeriodogramSeparatesTwoTonesByAmplitude(t *testing.T) {
	u := grid(4000, 0.005)
	y := make([]float64, len(u))
	for i, ui := range u {
		y[i] = 3*math.Cos(10*ui) + 1*math.Cos(20*ui+0.7)
	}

	freqs := []float64{10, 20}
	power := Periodogram(u, y, freqs)

	if power[0] <= power[1] {
		t.Errorf("power at f=10 (%v) should exceed power at f=20 (%v)", power[0], power[1])
	}
	if power[1] < power[0]/20 {
		t.Errorf("the weaker tone vanished: %v vs %v", power[1], power[0])
	}
}

// TestPeriodogramIgnoresAConstant pins the mean removal: a flat series has no
// periodic content, whatever its offset.
func TestPeriodogramIgnoresAConstant(t *testing.T) {
	u := grid(500, 0.02)
	y := make([]float64, len(u))
	for i := range y {
		y[i] = 42.0
	}

	power := Periodogram(u, y, []float64{5, 14, 25})
	for i, p := range power {
		if p > 1e-18 {
			t.Errorf("power[%d] = %v on a constant series, want ~0", i, p)
		}
	}
}

func TestPeriodogramHandlesDegenerateInput(t *testing.T) {
	if got := Periodogram(nil, nil, []float64{1}); got != nil {
		t.Errorf("empty series = %v, want nil", got)
	}
	if got := Periodogram([]float64{1, 2}, []float64{1}, []float64{1}); got != nil {
		t.Errorf("mismatched lengths = %v, want nil", got)
	}
	if got := Periodogram([]float64{1}, []float64{1}, nil); len(got) != 0 {
		t.Errorf("no frequencies = %v, want empty", got)
	}
}
