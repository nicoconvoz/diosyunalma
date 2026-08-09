// Command radios turns on every dial at once.
//
// Each residue tribe has a mother function — a Dirichlet L-function — and its
// imbalance oscillates at that function's zeros. Three tribal radios play
// here: mod 3, mod 4 and mod 5 (the golden tribe, already verified against
// LMFDB in Finding 41). And one more: THE HARMONY DIAL.
//
// The harmony is a classical theorem: ζ(s)·L(s,χ₅) is the Dedekind zeta
// function of Q(√5) — the number world whose integers are built from the
// golden ratio. Its zeros are the UNION of both station lists, so the summed
// signal ψ(x) + ψ(x,χ₅) − x, normalised, must show zeta's stations AND the
// golden tribe's stations on one dial. Pre-registered before looking.
//
// Usage:
//
//	go run ./cmd/radios [-limit N]
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

func chi3(n int) float64 {
	switch n % 3 {
	case 1:
		return 1
	case 2:
		return -1
	}
	return 0
}

func chi4(n int) float64 {
	switch n % 4 {
	case 1:
		return 1
	case 3:
		return -1
	}
	return 0
}

func chi5(n int) float64 {
	switch n % 5 {
	case 1, 4:
		return 1
	case 2, 3:
		return -1
	}
	return 0
}

// psiWeighted accumulates Σ w(p^k)·ln p at ascending sample points.
func psiWeighted(ps []int, limit int, xs []int, w func(int) float64) []float64 {
	type event struct {
		at int
		v  float64
	}
	powers := []event{}
	for _, p := range ps {
		if p*p > limit {
			break
		}
		lg := math.Log(float64(p))
		for pk := p * p; pk <= limit; pk *= p {
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

type peak struct{ f, p float64 }

func topPeaks(freqs, power []float64, n int) []peak {
	cands := []peak{}
	for i := 1; i+1 < len(power); i++ {
		if power[i] > power[i-1] && power[i] > power[i+1] {
			f := freqs[i]
			if den := power[i-1] - 2*power[i] + power[i+1]; den != 0 {
				f += 0.5 * (power[i-1] - power[i+1]) / den * 0.005
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
	uMin, uMax := math.Log(100), math.Log(float64(*limit))
	us, xs := []float64{}, []int{}
	for u := uMin; u <= uMax; u += du {
		us = append(us, u)
		xs = append(xs, int(math.Round(math.Exp(u))))
	}
	freqs := []float64{}
	for f := 2.0; f <= 30.0; f += 0.005 {
		freqs = append(freqs, f)
	}

	ps := primes.Sieve(*limit)
	hann := func(y []float64) []float64 {
		out := make([]float64, len(y))
		n := float64(len(y) - 1)
		for i, v := range y {
			out[i] = v * 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/n))
		}
		return out
	}

	tune := func(name string, sig []float64) []peak {
		e := make([]float64, len(sig))
		for i := range sig {
			e[i] = sig[i] / math.Sqrt(float64(xs[i]))
		}
		pw := spectral.Periodogram(us, hann(e), freqs)
		pk := topPeaks(freqs, pw, 6)
		fmt.Printf("\n%s\n  stations:", name)
		for _, q := range pk {
			fmt.Printf("  %.4f", q.f)
		}
		fmt.Println()
		return pk
	}

	fmt.Println("ALL RADIOS ON — one sieve, four dials")

	tune("RADIO mod 3 — L(s, χ₃)", psiWeighted(ps, *limit, xs, chi3))
	tune("RADIO mod 4 — L(s, χ₄)", psiWeighted(ps, *limit, xs, chi4))
	tune("RADIO mod 5 — L(s, χ₅)  [verified 8/8 in Finding 41]",
		psiWeighted(ps, *limit, xs, chi5))

	// The harmony dial: ψ + ψ(χ₅) − x is the prime-ideal count of Z[φ],
	// the golden field, minus its main term.
	plain := psiWeighted(ps, *limit, xs, func(int) float64 { return 1 })
	gold := psiWeighted(ps, *limit, xs, chi5)
	harmony := make([]float64, len(xs))
	for i := range xs {
		harmony[i] = plain[i] + gold[i] - float64(xs[i])
	}
	hp := tune("RADIO 4 — THE HARMONY: Dedekind zeta of Q(√5), the golden field", harmony)

	zetaS := []float64{14.1347, 21.0220, 25.0109}
	goldS := []float64{6.6485, 9.8314, 11.9588, 16.0338, 17.5670}
	both := 0
	for _, q := range hp {
		for _, z := range append(append([]float64{}, zetaS...), goldS...) {
			if math.Abs(q.f-z) < 0.15 {
				both++
				break
			}
		}
	}
	fmt.Printf("\nharmony stations matching the union of both known lists: %d of %d\n",
		both, len(hp))
	fmt.Println("pre-registered: the harmony dial carries BOTH musics at once.")
}
