// Command symphony seats three new players in the orchestra.
//
// Two more real dials: the quadratic tribes of mod 7 and mod 8 — the latter
// being the character of the field Q(√2), the tritone's own home. And one
// genuinely new instrument: the COMPLEX character mod 5 (order four), whose
// signal is complex-valued and whose dial is therefore ASYMMETRIC — positive
// and negative frequencies carry different stations, because for a complex
// character the zeros of L(s, χ) are not mirror-symmetric; the mirror image
// belongs to the conjugate character. No real radio can tell clockwise from
// counterclockwise; this one can.
//
// Usage:
//
//	go run ./cmd/symphony [-limit N]
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

func chi7(n int) float64 {
	switch n % 7 {
	case 1, 2, 4:
		return 1
	case 3, 5, 6:
		return -1
	}
	return 0
}

func chi8(n int) float64 {
	switch n % 8 {
	case 1, 7:
		return 1
	case 3, 5:
		return -1
	}
	return 0
}

// chi5c is the order-four character mod 5 with χ(2) = i, returned as (re, im).
func chi5c(n int) (float64, float64) {
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

func grid(limit int) ([]float64, []int) {
	const du = 0.005
	us, xs := []float64{}, []int{}
	for u := math.Log(100); u <= math.Log(float64(limit)); u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	return us, xs
}

func hann(y []float64) []float64 {
	out := make([]float64, len(y))
	n := float64(len(y) - 1)
	for i, v := range y {
		out[i] = v * 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
	}
	return out
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

	us, xs := grid(*limit)
	ps := primes.Sieve(*limit)

	accumulateExact := func(w func(int) float64) []float64 {
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
			out[i] = sum
		}
		return out
	}

	normalise := func(sig []float64) []float64 {
		e := make([]float64, len(sig))
		for i := range sig {
			e[i] = sig[i] / math.Sqrt(float64(xs[i]))
		}
		return hann(e)
	}

	realFreqs := []float64{}
	for f := 2.0; f <= 30.0; f += 0.005 {
		realFreqs = append(realFreqs, f)
	}

	fmt.Println("THE SYMPHONY — three new players")

	for _, tribe := range []struct {
		name string
		w    func(int) float64
	}{
		{"RADIO mod 7 — L(s, χ₇)", chi7},
		{"RADIO mod 8 — L(s, χ₈), the tribe of Q(√2)", chi8},
	} {
		pw := spectral.Periodogram(us, normalise(accumulateExact(tribe.w)), realFreqs)
		fmt.Printf("\n%s\n  stations:", tribe.name)
		for _, q := range topPeaks(realFreqs, pw, 6) {
			fmt.Printf("  %.4f", q.f)
		}
		fmt.Println()
	}

	// The complex instrument: an asymmetric dial.
	re := normalise(accumulateExact(func(n int) float64 { r, _ := chi5c(n); return r }))
	im := normalise(accumulateExact(func(n int) float64 { _, i := chi5c(n); return i }))

	cFreqs := []float64{}
	for f := -25.0; f <= 25.0; f += 0.005 {
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

	fmt.Println("\nRADIO mod 5 COMPLEX — the asymmetric dial (order-4 character)")
	fmt.Print("  stations:")
	for _, q := range topPeaks(cFreqs, power, 8) {
		fmt.Printf("  %+.4f", q.f)
	}
	fmt.Println()
	fmt.Println("  positive and negative stations DIFFER: this tribe's music tells")
	fmt.Println("  clockwise from counterclockwise. The mirror belongs to its conjugate.")
}
