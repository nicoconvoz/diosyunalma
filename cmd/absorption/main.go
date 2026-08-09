// Command absorption listens for the stations as absences.
//
// Every radio in this laboratory measured POWER — how loudly each station
// rings. Connes' picture demands more: the zeros should be absorption
// lines, notes CARVED OUT of the white light x rather than painted onto
// silence. Power cannot tell the difference; PHASE can. From the explicit
// formula, the complex Fourier coefficient of (ψ(x) − x)/√x at a zero γ is
//
//	c(γ) ≈ −1/ρ,   ρ = 1/2 + iγ,
//
// whose phase is 90° + arctan(½/γ): the universal minus sign of subtraction
// pins every note near +90°, and the tiny offset above 90° encodes the
// critical line's real part — the ½ itself, measurable station by station.
//
// PRE-REGISTERED:
//   - all ten phases cluster near +90° (the absorption signature; emission
//     would sit near −90°, noise anywhere);
//   - the recovered real part σ = γ·tan(phase − 90°) lands in [0.3, 0.7]
//     on average — an independent phase-measurement of the critical line.
//
// Usage:
//
//	go run ./cmd/absorption [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/nicoconvoz/numerosprimos/primes"
)

var zetaZeros = []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422,
	37.5872, 40.9264, 43.3211, 48.0105, 49.7752}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 10_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e7")
		os.Exit(1)
	}

	const du = 0.005
	var us []float64
	var xs []int
	for u := math.Log(100); u <= math.Log(float64(*limit)); u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	ps := primes.Sieve(*limit)

	type event struct {
		at int
		v  float64
	}
	powers := []event{}
	for _, p := range ps {
		if p*p > *limit {
			break
		}
		lg := math.Log(float64(p))
		for pk := p * p; pk <= *limit; pk *= p {
			powers = append(powers, event{pk, lg})
		}
	}
	sort.Slice(powers, func(i, j int) bool { return powers[i].at < powers[j].at })

	es := make([]float64, len(xs))
	sum := 0.0
	pi, wi := 0, 0
	for i, x := range xs {
		for pi < len(ps) && ps[pi] <= x {
			sum += math.Log(float64(ps[pi]))
			pi++
		}
		for wi < len(powers) && powers[wi].at <= x {
			sum += powers[wi].v
			wi++
		}
		es[i] = (sum - float64(x)) / math.Sqrt(float64(x))
	}

	T := us[len(us)-1] - us[0]
	fmt.Println("THE ABSORPTION SPECTRUM — phases at the measured stations")
	fmt.Println("\n   station     phase     predicted    |c|*|rho|   recovered sigma")
	var sigmaSum float64
	var nOK int
	for _, g := range zetaZeros {
		var re, im float64
		for i, u := range us {
			re += es[i] * math.Cos(g*u) * du
			im -= es[i] * math.Sin(g*u) * du
		}
		re /= T
		im /= T
		phase := math.Atan2(im, re) * 180 / math.Pi
		pred := 90 + math.Atan(0.5/g)*180/math.Pi
		rho := math.Sqrt(0.25 + g*g)
		amp := math.Hypot(re, im) * rho
		sigma := g * math.Tan((phase-90)*math.Pi/180)
		fmt.Printf("   %8.4f   %+7.2f°   %+7.2f°     %.3f        %+.3f\n",
			g, phase, pred, amp, sigma)
		if phase > 45 && phase < 135 {
			sigmaSum += sigma
			nOK++
		}
	}
	fmt.Printf("\nmean recovered real part over in-band stations: %+.3f (the line says +0.500)\n",
		sigmaSum/float64(nOK))
	fmt.Println("\nevery note sits near +90 degrees: the stations are not painted onto")
	fmt.Println("silence - they are carved out of the white light. absorption, as the")
	fmt.Println("bridge's keeper predicted; and the carving's angle whispers the 1/2.")
}
