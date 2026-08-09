// Command radio2 is the repaired instrument. The first attempt read the
// spectrum at single frequencies, where a periodogram has 100% relative error
// by construction; both columns bounced and proved nothing except that the
// instrument was broken. The repair is textbook: Welch averaging — split the
// signal into many segments, average their spectra per band.
//
// The flash under test: the primes' static should be missing its BASS. Every
// measured autocorrelation of gap values is negative (lag 1: −0.035, decaying
// slowly), which means long waves are suppressed: the primes keep a budget,
// overspending on one gap is repaid by the next ones. White noise has no such
// memory; the shuffled control must come out flat.
package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/nicoconvoz/numerosprimos/control"
	"github.com/nicoconvoz/numerosprimos/primes"
	"github.com/nicoconvoz/numerosprimos/spectral"
)

const (
	segLen   = 2048
	segCount = 512
)

func centred(gaps []int) []float64 {
	mean := 0.0
	for _, g := range gaps {
		mean += float64(g)
	}
	mean /= float64(len(gaps))
	out := make([]float64, len(gaps))
	for i, g := range gaps {
		out[i] = float64(g) - mean
	}
	return out
}

// welch averages per-segment periodograms: the standard cure for the
// single-shot periodogram's built-in 100% variance.
func welch(x []float64, freqs []float64) []float64 {
	us := make([]float64, segLen)
	for i := range us {
		us[i] = float64(i)
	}
	sum := make([]float64, len(freqs))
	for s := 0; s < segCount; s++ {
		seg := x[s*segLen : (s+1)*segLen]
		p := spectral.Periodogram(us, seg, freqs)
		for i := range sum {
			sum[i] += p[i]
		}
	}
	for i := range sum {
		sum[i] /= float64(segCount)
	}
	return sum
}

func main() {
	rng := rand.New(rand.NewSource(2026))
	gaps := primes.Gaps(primes.From(primes.Sieve(100_000_000), 5))
	window := gaps[len(gaps)/2 : len(gaps)/2+segLen*segCount]
	x := centred(window)

	fmt.Printf("signal: %d gaps, %d segments of %d\n", len(x), segCount, segLen)

	// The prediction curve from the measured correlations, nothing fitted.
	maxLag := 100
	rho := make([]float64, maxLag+1)
	var den float64
	for _, v := range x {
		den += v * v
	}
	for lag := 1; lag <= maxLag; lag++ {
		var num float64
		for i := 0; i+lag < len(x); i++ {
			num += x[i] * x[i+lag]
		}
		rho[lag] = num / den
	}

	freqs := []float64{}
	for f := 0.02; f <= 0.49; f += 0.03 {
		freqs = append(freqs, 2*math.Pi*f)
	}

	real := welch(x, freqs)
	decoy := welch(centred(control.ShuffleGaps(window, rng)), freqs)

	norm := func(v []float64) []float64 {
		m := 0.0
		for _, p := range v {
			m += p
		}
		m /= float64(len(v))
		out := make([]float64, len(v))
		for i, p := range v {
			out[i] = p / m
		}
		return out
	}
	rN, dN := norm(real), norm(decoy)

	fmt.Println("\nTHE DIAL, REPAIRED — 512-segment average (1.00 = white)")
	fmt.Printf("%-8s %-10s %-10s %s\n", "freq", "PRIMES", "shuffled", "predicted from rho_k")
	for i, w := range freqs {
		pred := 1.0
		for lag := 1; lag <= maxLag; lag++ {
			pred += 2 * rho[lag] * math.Cos(w*float64(lag))
		}
		fmt.Printf("%-8.2f %-10.3f %-10.3f %.3f\n", w/(2*math.Pi), rN[i], dN[i], pred)
	}

	fmt.Printf("\nbass (f=0.02): primes %.3f vs shuffled %.3f\n", rN[0], dN[0])
	fmt.Printf("treble (f=0.47): primes %.3f vs shuffled %.3f\n", rN[len(rN)-1], dN[len(dN)-1])
	fmt.Println("\nmissing bass in the primes with a flat decoy = the static is rigid:")
	fmt.Println("the primes keep a budget; long waves are cancelled by the bookkeeping.")
}
