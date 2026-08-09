// Command deepen harvests deeper stations for the small dials.
//
// The tutti's songbook (Finding 58) is capped at 15.15 because the dials of
// moduli 3, 4, 7, 8, 11 and 13 were only ever harvested six stations deep.
// This run re-tunes each of those radios with a wider frequency range and a
// deeper take — twelve stations per dial — so the tutti's ceiling can rise.
//
// The first six stations of each dial were verified against the published
// tables in Findings 41–44; the deeper ones harvested here are the
// laboratory's own measurements, verified only by their stability.
//
// Usage:
//
//	go run ./cmd/deepen [-limit N]
package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/nicoconvoz/diosyunalma/primes"
	"github.com/nicoconvoz/diosyunalma/spectral"
)

func modPow(base, exp, mod int) int {
	r := 1
	base %= mod
	for exp > 0 {
		if exp&1 == 1 {
			r = r * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return r
}

func character(q int) func(int) float64 {
	switch q {
	case 4:
		return func(n int) float64 {
			switch n % 4 {
			case 1:
				return 1
			case 3:
				return -1
			}
			return 0
		}
	case 8:
		return func(n int) float64 {
			switch n % 8 {
			case 1, 7:
				return 1
			case 3, 5:
				return -1
			}
			return 0
		}
	}
	half := (q - 1) / 2
	table := make([]float64, q)
	for n := 1; n < q; n++ {
		switch modPow(n, half, q) {
		case 1:
			table[n] = 1
		case q - 1:
			table[n] = -1
		}
	}
	return func(n int) float64 { return table[n%q] }
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

func main() {
	limit := flag.Int("limit", 100_000_000, "sieve primes up to this value")
	flag.Parse()
	if *limit < 10_000_000 {
		fmt.Fprintln(os.Stderr, "limit must be at least 1e7")
		os.Exit(1)
	}

	const du = 0.005
	us, xs := []float64{}, []int{}
	for u := math.Log(100); u <= math.Log(float64(*limit)); u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	freqs := []float64{}
	for f := 1.0; f <= 34.0; f += 0.005 {
		freqs = append(freqs, f)
	}
	ps := primes.Sieve(*limit)

	for _, q := range []int{3, 4, 7, 8, 11, 13} {
		w := character(q)

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

		e := make([]float64, len(xs))
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
			e[i] = sum / math.Sqrt(float64(x))
		}
		n := float64(len(e) - 1)
		for i := range e {
			e[i] *= 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
		}

		pw := spectral.Periodogram(us, e, freqs)
		fmt.Printf("dial mod %d:", q)
		for _, k := range topPeaks(freqs, pw, 12) {
			fmt.Printf(" %.4f", k.f)
		}
		fmt.Println()
	}
}
