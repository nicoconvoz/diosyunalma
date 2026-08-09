// Command entrelazado decodes the quantum-entanglement flash and measures
// it: can the GPS's local data - computed entirely at home - touch the
// state of the distant sea without sailing there?
//
// The entangled pairs exist trivially and perfectly: every zero at
// 1/2+i*gamma has its mirror twin at 1/2-i*gamma (the two beams force
// it); measure one, know the other exactly, at any distance.
//
// The deeper entanglement is the pacemaker's: the 26 prime voices are
// computable AT HOME for any coordinate whatsoever. Here hundreds of
// windows are sailed for real (cheap water), the actual tide of each is
// measured, and the correlation with the home-computed forecast is the
// ENTANGLEMENT STRENGTH r: r^2 of the remote sea's state is known
// without travel. Control: shuffled pairing must decohere to r ~ 0.
//
// Usage:
//
//	go run ./cmd/entrelazado
package main

import (
	"fmt"
	"math"
	"math/rand"
)

var lnk, rsq [512]float64

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zRS(t float64) float64 {
	tau := t / (2 * math.Pi)
	a := math.Sqrt(tau)
	N := int(a)
	th := theta(t)
	var s float64
	for k := 1; k <= N; k++ {
		s += math.Cos(th-t*lnk[k]) * rsq[k]
	}
	p := a - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (N-1)%2 == 1 {
		sign = -1
	}
	return 2*s + sign*math.Pow(tau, -0.25)*c0
}

func smoothCount(t float64) float64 {
	return t/(2*math.Pi)*(math.Log(t/(2*math.Pi))-1) + 7.0/8
}

// countZeros counts sign changes of Z in [a, b] with a fine step.
func countZeros(a, b float64) int {
	n := 0
	step := 0.02
	prev := zRS(a)
	for t := a + step; t <= b; t += step {
		z := zRS(t)
		if (prev < 0) != (z < 0) {
			n++
		}
		prev = z
	}
	return n
}

// sPred: the pacemaker - the home-computed voice of the 26 first primes.
func sPred(t float64) float64 {
	s := 0.0
	for _, p := range []float64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37,
		41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101} {
		s -= math.Sin(math.Mod(t*math.Log(p), 2*math.Pi)) / math.Sqrt(p)
	}
	return s / math.Pi
}

func pearson(x, y []float64) float64 {
	n := float64(len(x))
	var sx, sy, sxx, syy, sxy float64
	for i := range x {
		sx += x[i]
		sy += y[i]
		sxx += x[i] * x[i]
		syy += y[i] * y[i]
		sxy += x[i] * y[i]
	}
	return (n*sxy - sx*sy) / math.Sqrt((n*sxx-sx*sx)*(n*syy-sy*sy))
}

func main() {
	fmt.Println("EL ENTRELAZAMIENTO — touching the distant sea from home")
	for k := 1; k < len(lnk); k++ {
		lnk[k] = math.Log(float64(k))
		rsq[k] = 1 / math.Sqrt(float64(k))
	}

	const windows = 400
	local := make([]float64, windows)  // computed at home, no sailing
	remote := make([]float64, windows) // measured by sailing the window
	rng := rand.New(rand.NewSource(2026))
	for i := 0; i < windows; i++ {
		t0 := 100000 + rng.Float64()*900000
		spacing := 2 * math.Pi / math.Log(t0/(2*math.Pi))
		span := 5 * spacing
		local[i] = sPred(t0+span) - sPred(t0)
		remote[i] = float64(countZeros(t0, t0+span)) - (smoothCount(t0+span) - smoothCount(t0))
	}
	r := pearson(local, remote)

	// decoherence control: shuffle the pairing.
	shuf := make([]float64, windows)
	copy(shuf, local)
	rng.Shuffle(windows, func(i, j int) { shuf[i], shuf[j] = shuf[j], shuf[i] })
	rc := pearson(shuf, remote)

	fmt.Printf("\n  %d windows sailed in cheap water; forecasts computed at home:\n", windows)
	fmt.Printf("    entanglement strength  r  = %.3f   (r^2 = %.0f%% of the remote state known from home)\n", r, 100*r*r)
	fmt.Printf("    decoherence control (shuffled pairing): r = %+.3f\n", rc)
	fmt.Printf("    F108's M-sigma prediction: the 26 voices carry ~2/3-5/6 of the variance -> r ~ 0.82-0.91\n")

	fmt.Println("\n  the decoded flash: the GPS coordinates and the primes were born of one")
	fmt.Println("  wavefunction - the explicit formula - like entangled particles born of")
	fmt.Println("  one event. Study the local tool (26 voices, computable for ANY")
	fmt.Println("  coordinate, even 1e24) and you touch the distant number's sea without")
	fmt.Println("  travelling. And the perfect pairs: every zero and its mirror twin at")
	fmt.Println("  1/2 - i*gamma - measure one, know the other, at any distance, exactly.")
}
