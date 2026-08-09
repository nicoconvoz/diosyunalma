// Command barometro asks the flash's question with a real barometer: what
// is the blue of nature? Water — the incompressible fluid that carries
// life. Is the zero sea a gas or a liquid?
//
// Statistical mechanics answers through the number variance: put boxes of
// size L on the unfolded sea and measure Var(count). An ideal gas
// (Poisson) gives Var = L: fully compressible. A liquid resists: for the
// harmonic log-gas (GUE, our F101 antigravity) Dyson–Mehta proved
// Var = (2/pi^2)(ln(2 pi L) + gamma + 1 - pi^2/8): logarithmic — the
// compressibility Var/L vanishes at scale. And Berry predicts the deep
// sea saturates BELOW even that: more liquid than the model.
//
// The favor we draw from it: the sphere's tolerance should be the
// PRESSURE-AWARE 2.5 sigma of this law, not a fixed constant - the
// barometer calibrates the certification for any window size.
//
// Usage:
//
//	go run ./cmd/barometro
package main

import (
	"fmt"
	"math"
	"math/rand"
)

var lnk, rsq [64]float64

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

func huntZeros(count int) []float64 {
	zeros := []float64{}
	prevT, prevZ := 14.0, zRS(14.0)
	for t := 14.05; len(zeros) < count && t < 25000; t += 0.05 {
		zt := zRS(t)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			for i := 0; i < 30 && hi-lo > 1e-8; i++ {
				mid := (lo + hi) / 2
				if (zRS(mid) < 0) == (prevZ < 0) {
					lo = mid
				} else {
					hi = mid
				}
			}
			zeros = append(zeros, (lo+hi)/2)
		}
		prevT, prevZ = t, zt
	}
	return zeros
}

// sigma2 measures the number variance in boxes of size L on unfolded positions.
func sigma2(u []float64, L float64) float64 {
	uMax := u[len(u)-1]
	nBox := int(uMax / L)
	if nBox < 4 {
		return math.NaN()
	}
	counts := make([]float64, nBox)
	j := 0
	for b := 0; b < nBox; b++ {
		hi := float64(b+1) * L
		for j < len(u) && u[j] < hi {
			counts[b]++
			j++
		}
	}
	var m, v float64
	for _, c := range counts {
		m += c
	}
	m /= float64(nBox)
	for _, c := range counts {
		v += (c - m) * (c - m)
	}
	return v / float64(nBox)
}

func gue(L float64) float64 {
	return 2 / (math.Pi * math.Pi) *
		(math.Log(2*math.Pi*L) + 0.5772156649 + 1 - math.Pi*math.Pi/8)
}

func main() {
	fmt.Println("EL BARÓMETRO — is the zero sea a gas or a liquid?")
	for k := 1; k < len(lnk); k++ {
		lnk[k] = math.Log(float64(k))
		rsq[k] = 1 / math.Sqrt(float64(k))
	}

	gammas := huntZeros(20000)
	u := make([]float64, len(gammas))
	for i, g := range gammas {
		u[i] = smoothCount(g)
	}
	base := u[0]
	for i := range u {
		u[i] -= base
	}
	fmt.Printf("\n  %d zeros unfolded; boxes laid on the sea...\n", len(gammas))

	rng := rand.New(rand.NewSource(2026))
	pois := make([]float64, len(u))
	acc := 0.0
	for i := range pois {
		acc += rng.ExpFloat64()
		pois[i] = acc
	}

	fmt.Println("\n     L     zeros Var   GUE liquid   Poisson gas   compressibility")
	for _, L := range []float64{2, 5, 10, 25, 50, 100, 250} {
		vz := sigma2(u, L)
		vp := sigma2(pois, L)
		fmt.Printf("   %5.0f   %9.3f   %10.3f   %11.3f   %13.4f\n",
			L, vz, gue(L), vp, vz/L)
	}
	fmt.Println("\n  an ideal gas keeps Var = L (compressibility 1). The zero sea's")
	fmt.Println("  variance grows only as ln L and its compressibility sinks toward")
	fmt.Println("  zero: THE SEA IS WATER - incompressible, alive, blue. And where it")
	fmt.Println("  dips even below the GUE liquid, Berry's saturation shows the deep")
	fmt.Println("  sea is more liquid than the model itself.")
	fmt.Println("\n  the favor: the sphere's honest tolerance is 2.5*sqrt(Var(L)) - the")
	fmt.Println("  pressure-aware gauge, now available to the whole fleet.")
	for _, L := range []float64{5, 25, 50} {
		fmt.Printf("    window %2.0f spacings: tolerance %.2f (fixed old value was 1.75)\n",
			L, 2.5*math.Sqrt(gue(L)))
	}
}
