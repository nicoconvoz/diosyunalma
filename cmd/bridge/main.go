// Command bridge walks the span between Li's criterion and Hilbert–Pólya,
// carrying this laboratory's own measured zeros across it.
//
// THE BRIDGE, in one line: |1 − 1/ρ| = 1 exactly when Re ρ = 1/2. The Möbius
// map sends the critical line to the unit circle, a self-adjoint operator to a
// unitary one (Cayley), and Li's coefficients to traces: λₙ = Tr(I − Uⁿ),
// each conjugate pair contributing 2(1 − cos nθ) ≥ 0. Li positivity is what
// the world looks like when the Hilbert–Pólya operator exists.
//
// THE MEASUREMENT THIS ADDS. The zeta hunt measured the heights γ of ten
// zeros but never their real parts β — and β is what the Riemann Hypothesis
// is about. The explicit formula makes each zero's contribution to
// (ψ(x) − x)/√x an oscillation of amplitude x^(β−1/2): constant when β = 1/2,
// drifting like e^((β−1/2)u) otherwise. Splitting the sampled range into two
// halves and comparing each zero's spectral power measures β directly:
//
//	β̂ = 1/2 + ln(P₂/P₁) / (2·Δū)
//
// PRE-REGISTERED: if RH holds, every β̂ lands near 0.5 and every Möbius image
// lands on the unit circle.
//
// Usage:
//
//	go run ./cmd/bridge [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"math/cmplx"
	"os"

	"github.com/nicoconvoz/diosyunalma/primes"
	"github.com/nicoconvoz/diosyunalma/riemann"
	"github.com/nicoconvoz/diosyunalma/spectral"
)

// measuredGammas are this laboratory's own zero heights, measured by cmd/zeta
// at 10^9 (Finding 26).
var measuredGammas = []float64{
	14.1349, 21.0211, 25.0044, 30.4282, 32.9422,
	37.5872, 40.9264, 43.3211, 48.0105, 49.7752,
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()

	if *limit < 1_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e6")
		os.Exit(1)
	}

	// The same sampled series as the zeta hunt.
	const du = 0.005
	uMin, uMax := math.Log(100), math.Log(float64(*limit))
	us := []float64{}
	xs := []int{}
	for u := uMin; u <= uMax; u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	psi := primes.PsiAt(*limit, xs)
	e := make([]float64, len(psi))
	for i := range psi {
		x := float64(xs[i])
		e[i] = (psi[i] - x) / math.Sqrt(x)
	}

	// Split into halves; measure each zero's power in each.
	half := len(us) / 2
	u1, y1 := us[:half], hann(e[:half])
	u2, y2 := us[half:], hann(e[half:])
	deltaU := mean(u2) - mean(u1)

	fmt.Printf("samples: %d + %d, half-centres %.2f and %.2f, Δū = %.2f\n",
		half, len(us)-half, mean(u1), mean(u2), deltaU)

	fmt.Println("\nTHE REAL PARTS, MEASURED — the content of the hypothesis, zero by zero")
	fmt.Printf("%-4s %-12s %-12s %-12s %-10s %s\n",
		"#", "γ measured", "P first half", "P second", "β̂", "|1−1/ρ̂| (1 = on the circle)")

	zeros := []complex128{}
	sumBeta := 0.0
	for i, g := range measuredGammas {
		p1 := spectral.Periodogram(u1, y1, []float64{g})[0]
		p2 := spectral.Periodogram(u2, y2, []float64{g})[0]
		beta := 0.5 + math.Log(p2/p1)/(4*deltaU) // power is amplitude², hence 4

		rho := complex(beta, g)
		zeros = append(zeros, rho, cmplx.Conj(rho))
		sumBeta += beta

		fmt.Printf("%-4d %-12.4f %-12.4g %-12.4g %-10.4f %.5f\n",
			i+1, g, p1, p2, beta, cmplx.Abs(1-1/rho))
	}
	fmt.Printf("\nmean β̂ = %.4f    (RH says 0.5)\n", sumBeta/float64(len(measuredGammas)))

	// Partial Li coefficients from the measured spectrum.
	fmt.Println("\nPARTIAL LI COEFFICIENTS — λₙ = Tr(I − Uⁿ) over the measured zeros")
	fmt.Printf("%-8s %-14s %s\n", "n", "λₙ (measured)", "λₙ if forced off-line (β=0.75 quadruple added)")

	offLine := append(append([]complex128{}, zeros...),
		complex(0.75, 14.1349), complex(0.75, -14.1349),
		complex(0.25, 14.1349), complex(0.25, -14.1349))

	worstReal, worstOff := math.Inf(1), math.Inf(1)
	for _, n := range []int{1, 2, 3, 5, 8, 13, 21, 34, 55, 100, 300, 1000, 3000, 8000} {
		lr := riemann.LiSum(zeros, n)
		lo := riemann.LiSum(offLine, n)
		if lr < worstReal {
			worstReal = lr
		}
		if lo < worstOff {
			worstOff = lo
		}
		fmt.Printf("%-8d %-14.4f %.4f\n", n, lr, lo)
	}

	fmt.Printf("\nminimum λₙ, measured zeros    : %.4f\n", worstReal)
	fmt.Printf("minimum λₙ, with off-line zero: %.4f   <- the crash Li's criterion detects\n", worstOff)
}

func hann(y []float64) []float64 {
	out := make([]float64, len(y))
	n := float64(len(y) - 1)
	for i, v := range y {
		out[i] = v * 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
	}
	return out
}

func mean(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}
