// Command baton sharpens the conductor.
//
// Finding 46 showed orthogonality's baton working above chance but blurry:
// the tribal separation rested on sixteen mod-5 zeros reaching |γ| ≈ 25.
// This run harvests the tribes' stations twice as deep — the real character
// to γ ≈ 42 and the complex one across ±42 — halving the blur, then conducts
// again with a corrected score: greedy one-peak-one-target matching at a
// tighter tolerance, so no column is flattered.
//
// PRE-REGISTERED: with roughly twice the tribal zeros, the own-versus-wrong
// ratio of Finding 46 (11 against 8) must improve markedly, or the conductor
// hypothesis is in trouble.
//
// Usage:
//
//	go run ./cmd/baton [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/nicoconvoz/numerosprimos/primes"
	"github.com/nicoconvoz/numerosprimos/spectral"
)

var zetaZeros = []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422,
	37.5872, 40.9264, 43.3211, 48.0105, 49.7752}

func chi5(a int) float64 {
	if a == 1 || a == 4 {
		return 1
	}
	return -1
}

func chiC(a int) complex128 {
	switch a {
	case 1:
		return 1
	case 2:
		return complex(0, 1)
	case 3:
		return complex(0, -1)
	}
	return -1
}

func chi5w(n int) float64 {
	switch n % 5 {
	case 1, 4:
		return 1
	case 2, 3:
		return -1
	}
	return 0
}

func chiCw(n int) (float64, float64) {
	switch n % 5 {
	case 1:
		return 1, 0
	case 2:
		return 0, 1
	case 3:
		return 0, -1
	case 4:
		return -1, 0
	}
	return 0, 0
}

type peak struct{ f, p float64 }

func topPeaks(freqs, power []float64, n int) []peak {
	cands := []peak{}
	for i := 1; i+1 < len(power); i++ {
		if power[i] > power[i-1] && power[i] > power[i+1] {
			f := freqs[i]
			if den := power[i-1] - 2*power[i] + power[i+1]; den != 0 {
				f += 0.5 * (power[i-1] - power[i+1]) / den * (freqs[1] - freqs[0])
			}
			cands = append(cands, peak{f, power[i]})
		}
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].p > cands[j].p })
	out := []peak{}
	for _, c := range cands {
		ok := true
		for _, k := range out {
			if math.Abs(k.f-c.f) < 0.4 {
				ok = false
				break
			}
		}
		if ok {
			out = append(out, c)
		}
		if len(out) == n {
			break
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].f < out[j].f })
	return out
}

var targets = map[int][]float64{
	1: {11, 16, 31},
	2: {2, 7, 17, 27, 32, 37},
	3: {3, 8, 13, 23},
	4: {4, 9, 19, 29},
}

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 10_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e7")
		os.Exit(1)
	}

	// The shared signal machinery.
	const du = 0.005
	us, xs := []float64{}, []int{}
	for u := math.Log(100); u <= math.Log(float64(*limit)); u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	ps := primes.Sieve(*limit)

	accumulate := func(w func(int) float64) []float64 {
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
				powers = append(powers, event{pk, w(pk) * lg})
			}
		}
		sort.Slice(powers, func(i, j int) bool { return powers[i].at < powers[j].at })
		out := make([]float64, len(xs))
		sum := 0.0
		pi, wi := 0, 0
		for i, x := range xs {
			for pi < len(ps) && ps[pi] <= x {
				sum += w(ps[pi]) * math.Log(float64(ps[pi]))
				pi++
			}
			for wi < len(powers) && powers[wi].at <= x {
				sum += powers[wi].v
				wi++
			}
			out[i] = sum / math.Sqrt(float64(x))
		}
		n := float64(len(out) - 1)
		for i := range out {
			out[i] *= 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
		}
		return out
	}

	// -------- harvest: the real tribe, twice as deep --------
	e5 := accumulate(chi5w)
	freqs := []float64{}
	for f := 2.0; f <= 42.0; f += 0.005 {
		freqs = append(freqs, f)
	}
	realZ := []float64{}
	for _, k := range topPeaks(freqs, spectral.Periodogram(us, e5, freqs), 14) {
		realZ = append(realZ, k.f)
	}
	fmt.Printf("harvested real-tribe zeros (%d):", len(realZ))
	for _, g := range realZ {
		fmt.Printf(" %.3f", g)
	}
	fmt.Println()

	// -------- harvest: the complex tribe, both sides --------
	re := accumulate(func(n int) float64 { r, _ := chiCw(n); return r })
	im := accumulate(func(n int) float64 { _, i := chiCw(n); return i })
	cFreqs := []float64{}
	for f := -42.0; f <= 42.0; f += 0.005 {
		if math.Abs(f) > 1.5 {
			cFreqs = append(cFreqs, f)
		}
	}
	power := make([]float64, len(cFreqs))
	for i, f := range cFreqs {
		var sr, si float64
		for j, u := range us {
			c, s := math.Cos(f*u), math.Sin(f*u)
			sr += re[j]*c + im[j]*s
			si += im[j]*c - re[j]*s
		}
		power[i] = (sr*sr + si*si) / float64(len(us))
	}
	complexZ := []float64{}
	for _, k := range topPeaks(cFreqs, power, 16) {
		complexZ = append(complexZ, k.f)
	}
	fmt.Printf("harvested complex-tribe zeros (%d):", len(complexZ))
	for _, g := range complexZ {
		fmt.Printf(" %+.3f", g)
	}
	fmt.Println()

	// -------- conduct with the sharpened baton --------
	clock := func(a int, x float64) float64 {
		u := math.Log(x)
		s := 0.0
		for _, g := range zetaZeros {
			s += math.Cos(g * u)
		}
		for _, g := range realZ {
			s += chi5(a) * math.Cos(g*u)
		}
		var sc complex128
		for _, g := range complexZ {
			sc += complex(math.Cos(g*u), math.Sin(g*u))
		}
		w := complex(real(chiC(a)), -imag(chiC(a)))
		s += real(w * sc)
		return 1 - 2*s/math.Sqrt(x)
	}

	fmt.Println("\nTHE SHARPENED BATON — greedy scoring, tolerance 0.35, no shared targets")
	totalOwn, totalWrong := 0, 0
	for a := 1; a <= 4; a++ {
		type pk struct{ x, v float64 }
		peaks := []pk{}
		prev, cur := clock(a, 2.00), clock(a, 2.01)
		for x := 2.02; x <= 40; x += 0.01 {
			next := clock(a, x)
			if cur > prev && cur > next {
				peaks = append(peaks, pk{x - 0.01, cur})
			}
			prev, cur = cur, next
		}
		sort.Slice(peaks, func(i, j int) bool { return peaks[i].v > peaks[j].v })
		if len(peaks) > 6 {
			peaks = peaks[:6]
		}

		claimed := map[float64]bool{}
		own, wrong := 0, 0
		sort.Slice(peaks, func(i, j int) bool { return peaks[i].x < peaks[j].x })
		fmt.Printf("baton a=%d  peaks:", a)
		for _, p := range peaks {
			bestCls, bestT, bestD := 0, 0.0, 0.35
			for cls, list := range targets {
				for _, t := range list {
					if d := math.Abs(p.x - t); d < bestD && !claimed[t] {
						bestD, bestCls, bestT = d, cls, t
					}
				}
			}
			mark := " "
			if bestCls == a {
				own++
				claimed[bestT] = true
				mark = "*"
			} else if bestCls != 0 {
				wrong++
				claimed[bestT] = true
				mark = "x"
			}
			fmt.Printf("  %.2f%s", p.x, mark)
		}
		fmt.Printf("   -> own %d, wrong %d\n", own, wrong)
		totalOwn += own
		totalWrong += wrong
	}
	fmt.Printf("\nTOTAL: %d own vs %d wrong   (Finding 46 with the blunt baton: 11 vs 8)\n",
		totalOwn, totalWrong)
}
