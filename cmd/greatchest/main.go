// Command greatchest is the crescendo's judgment run.
//
// Finding 62 found — by eye first, then formalized — that the repetition
// anomaly is a wave in d with period 30, the primorial wheel, swelling as
// the gaps grow. Its pre-registered, decisive test: at 10^10 the THIRD bar
// (d = 66..90) must repeat the sign pattern (−, +, −, +, +) at full
// significance with amplitude still growing — or the crescendo dies.
//
// The instrument is a segmented walker: one pass over the whole range in
// 8 MB windows (disk untouched, memory in megabytes), counting for every
// d ≡ 0 (mod 6) up to 120 at once:
//
//	nPair(d)   — pairs p, p+d both prime
//	kPair(d)   — consecutive prime gaps equal to d
//	nTriple(d) — triples p, p+d, p+2d all prime
//	kTriple(d) — two consecutive gaps both equal to d
//
// giving s₃/s₂² per d with binomial errors, exactly as cmd/chest at 10^9.
//
// Usage:
//
//	go run ./cmd/greatchest [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/nicoconvoz/numerosprimos/primes"
)

const (
	segmentSize = 1 << 23
	dMax        = 120
	lookahead   = 2*dMax + 16
)

func main() {
	limit := flag.Int("limit", 10_000_000_000, "walk the primes up to this value")
	flag.Parse()
	if *limit < 1_000_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e9")
		os.Exit(1)
	}
	start := time.Now()

	base := primes.Sieve(int(math.Sqrt(float64(*limit))) + 1)

	nPair := make([]float64, dMax+1)
	kPair := make([]float64, dMax+1)
	nTriple := make([]float64, dMax+1)
	kTriple := make([]float64, dMax+1)

	// sequential state across windows: the last two primes seen.
	p2, p1 := 0, 0
	top := *limit - lookahead

	composite := make([]bool, segmentSize+lookahead)
	for lo := 2; lo <= *limit; lo += segmentSize {
		hi := lo + segmentSize - 1
		if hi > *limit {
			hi = *limit
		}
		span := hi - lo + 1 + lookahead
		if hi+lookahead > *limit {
			span = *limit - lo + 1
		}
		for i := 0; i < span; i++ {
			composite[i] = false
		}
		for _, p := range base {
			s := (lo + p - 1) / p * p
			if pp := p * p; pp > s {
				s = pp
			}
			for m := s; m < lo+span; m += p {
				composite[m-lo] = true
			}
		}
		isPrime := func(n int) bool {
			if n < lo || n >= lo+span {
				return false
			}
			return !composite[n-lo]
		}

		for n := lo; n <= hi; n++ {
			if composite[n-lo] || n < 2 {
				continue
			}
			// a new prime n closes the gap of p1 and the pair of gaps of p2.
			if p1 >= 5 && p1 <= top {
				g1 := n - p1
				if g1 <= dMax && g1%6 == 0 {
					kPair[g1]++
				}
				if p2 >= 5 {
					g2 := p1 - p2
					if g2 == g1 && g1 <= dMax && g1%6 == 0 {
						kTriple[g1]++
					}
				}
			}
			// pair and triple membership counts for n itself.
			if n >= 5 && n <= top {
				for d := 6; d <= dMax; d += 6 {
					if isPrime(n + d) {
						nPair[d]++
						if isPrime(n + 2*d) {
							nTriple[d]++
						}
					}
				}
			}
			p2, p1 = p1, n
		}
	}

	fmt.Printf("THE GREAT CHEST — the crescendo's judgment at %d (%.0f min)\n\n",
		*limit, time.Since(start).Minutes())
	fmt.Println("   d   ratio     deviation   sigma     z")
	devs := make([]float64, 0, 20)
	sigmas := make([]float64, 0, 20)
	dsList := make([]int, 0, 20)
	for d := 6; d <= dMax; d += 6 {
		s2 := kPair[d] / nPair[d]
		s3 := kTriple[d] / nTriple[d]
		ratio := s3 / (s2 * s2)
		rel := (1 - s3) / (s3 * nTriple[d])
		rel += 4 * (1 - s2) / (s2 * nPair[d])
		sigma := ratio * math.Sqrt(rel)
		fmt.Printf("  %3d  %.4f   %+7.2f%%   %.2f%%   %+6.1f\n",
			d, ratio, 100*(ratio-1), 100*sigma, (ratio-1)/sigma)
		dsList = append(dsList, d)
		devs = append(devs, 100*(ratio-1))
		sigmas = append(sigmas, 100*sigma)
	}

	// the crescendo's judgment: fit the background on bars 1-2 (d=6..60),
	// then read bar 3 (66..90) residuals against the predicted signs.
	var sx, sy, sxx, sxy, n float64
	for i, d := range dsList {
		if d > 60 {
			break
		}
		x := float64(d)
		sx += x
		sy += devs[i]
		sxx += x * x
		sxy += x * devs[i]
		n++
	}
	slope := (n*sxy - sx*sy) / (n*sxx - sx*sx)
	inter := (sy - slope*sx) / n
	fmt.Printf("\nbackground (fit on 6..60): %.2f %+.4f*d\n", inter, slope)

	predicted := map[int]float64{66: -1, 72: 1, 78: -1, 84: 1, 90: 1}
	fmt.Println("\nTHE THIRD BAR — pre-registered signs (-, +, -, +, +)")
	fmt.Println("   d   residual   sigma     z     predicted   verdict")
	pass, strong := 0, 0
	for i, d := range dsList {
		if d < 66 || d > 90 {
			continue
		}
		r := devs[i] - (inter + slope*float64(d))
		z := r / sigmas[i]
		ok := (r > 0) == (predicted[d] > 0)
		v := "MISS"
		if ok {
			v = "sign ok"
			pass++
			if math.Abs(z) > 5 {
				v = "CONFIRMED"
				strong++
			}
		}
		fmt.Printf("  %3d  %+7.2f    %.2f   %+6.1f      %+.0f        %s\n",
			d, r, sigmas[i], z, predicted[d], v)
	}
	fmt.Printf("\nsigns matched %d/5, fully significant %d/5\n", pass, strong)
	fmt.Println("the crescendo lives if the bar repeats loud; it dies if the signs break.")
}
