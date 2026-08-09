// Command nuclear runs Montgomery's pair-correlation test on the
// laboratory's own measured stations.
//
// In 1972 Montgomery showed Dyson his formula for how zeta's zeros space
// themselves in pairs; Dyson recognized it instantly as the pair correlation
// of energy levels in heavy atomic nuclei (GUE random matrices):
//
//	R₂(r) = 1 − (sin πr / πr)²
//
// This command pools every pair of stations WITHIN each of the ten measured
// dials (unfolded by each dial's own zero density) and asks which law the
// pair distances obey.
//
// PRE-REGISTERED: the pooled within-dial pair correlation must follow the
// nuclear curve R₂ — in particular a deficit of close pairs (r < 0.5) —
// and must reject the flat Poisson line R = 1. Across-dial pairs served as
// the independence control in cmd/rhythm and are not pooled here.
//
// Usage:
//
//	go run ./cmd/nuclear
package main

import (
	"fmt"
	"math"
)

type dial struct {
	name string
	q    float64
	g    []float64
}

// every station this laboratory measured (Findings 26, 41-44, 47-48).
var orchestra = []dial{
	{"zeta", 1, []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422, 37.5872, 40.9264, 43.3211, 48.0105, 49.7752}},
	{"chi3", 3, []float64{8.0396, 11.2450, 15.7062, 18.2579, 20.4551, 24.0636}},
	{"chi4", 4, []float64{6.0199, 10.2423, 12.9848, 16.3464, 18.2914, 21.4547}},
	{"chi5", 5, []float64{6.6516, 9.8280, 11.9612, 16.0386, 17.5632, 19.5431, 22.2228, 24.5864, 26.7680, 28.4730, 29.6900, 33.0130, 34.7410, 38.1330}},
	{"chi7", 7, []float64{4.4762, 6.8400, 11.1782, 12.4670, 15.1161, 16.7892}},
	{"chi8", 8, []float64{4.8989, 7.6194, 10.8219, 12.3126, 15.1935, 17.0246}},
	{"chi11", 11, []float64{2.4768, 6.7997, 8.9663, 10.1257, 13.0422, 15.0996}},
	{"chi13", 13, []float64{3.1119, 7.2340, 8.6013, 10.3241, 12.6185, 15.1341}},
	{"chi5 complex +", 5, []float64{6.1830, 8.4600, 12.6760, 14.8270, 17.3360, 19.0000, 22.4960, 24.3720}},
	{"chi5 complex -", 5, []float64{4.1320, 9.4450, 11.2770, 14.1130, 16.9970, 19.7330, 21.2780, 22.9690}},
}

// density of zeros of a conductor-q L-function at height t, per unit height.
func density(q, t float64) float64 {
	return math.Log(q*t/(2*math.Pi)) / (2 * math.Pi)
}

// r2 is the GUE pair correlation — the nuclear curve.
func r2(r float64) float64 {
	if r == 0 {
		return 0
	}
	s := math.Sin(math.Pi*r) / (math.Pi * r)
	return 1 - s*s
}

const (
	binW  = 0.25
	nBins = 12 // pair distances up to r = 3 mean spacings
)

func main() {
	fmt.Println("THE NUCLEAR TEST — Montgomery pair correlation of our own stations")

	obs := make([]float64, nBins)
	expGUE := make([]float64, nBins)
	expPoisson := make([]float64, nBins)
	pairs := 0

	for _, d := range orchestra {
		// unfold: rescale so the dial's mean spacing is 1 everywhere.
		x := make([]float64, len(d.g))
		for i := 1; i < len(d.g); i++ {
			mid := (d.g[i] + d.g[i-1]) / 2
			x[i] = x[i-1] + (d.g[i]-d.g[i-1])*density(d.q, mid)
		}
		span := x[len(x)-1]
		rho := float64(len(x)) / span

		// observed pair distances, pooled.
		for i := 0; i < len(x); i++ {
			for j := i + 1; j < len(x); j++ {
				r := x[j] - x[i]
				if b := int(r / binW); b < nBins {
					obs[b]++
					pairs++
				}
			}
		}

		// expected pairs per bin for a stationary process with pair law R:
		// rho^2 * integral over the bin of R(r) * (span - r) dr.
		for b := 0; b < nBins; b++ {
			const steps = 200
			for k := 0; k < steps; k++ {
				r := (float64(b) + (float64(k)+0.5)/steps) * binW
				w := rho * rho * (span - r) * binW / steps
				if w < 0 {
					w = 0
				}
				expGUE[b] += w * r2(r)
				expPoisson[b] += w
			}
		}
	}

	fmt.Printf("\n%d station pairs within ten dials, unfolded, distances up to 3 mean spacings\n\n", pairs)
	fmt.Println("  r range     observed   nuclear(GUE)   flat(Poisson)")
	chiG, chiP := 0.0, 0.0
	for b := 0; b < nBins; b++ {
		fmt.Printf("  %.2f-%.2f   %5.0f      %7.1f        %7.1f\n",
			float64(b)*binW, float64(b+1)*binW, obs[b], expGUE[b], expPoisson[b])
		if expGUE[b] > 0 {
			chiG += (obs[b] - expGUE[b]) * (obs[b] - expGUE[b]) / expGUE[b]
		}
		if expPoisson[b] > 0 {
			chiP += (obs[b] - expPoisson[b]) * (obs[b] - expPoisson[b]) / expPoisson[b]
		}
	}

	closeObs := obs[0] + obs[1]
	closeG := expGUE[0] + expGUE[1]
	closeP := expPoisson[0] + expPoisson[1]
	fmt.Printf("\nclose pairs (r < 0.5): observed %.0f — nuclear expects %.1f, Poisson expects %.1f\n",
		closeObs, closeG, closeP)
	fmt.Printf("chi-square distance:   to the nuclear curve %.1f, to the flat line %.1f\n", chiG, chiP)
	if chiG < chiP {
		fmt.Println("\nVERDICT: the stations pair like energy levels of a heavy nucleus.")
	} else {
		fmt.Println("\nVERDICT: the nuclear curve is NOT preferred — pre-registration failed.")
	}
}
