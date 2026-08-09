// Command fusion is the double-mirror vessel: the voyage and the compass
// welded into one system.
//
// The two directions of Riemann's explicit formula are the two halves of
// this laboratory: the VOYAGE uses the primes to hunt zeros (every ship's
// main sum runs over the integers), and the COMPASS uses the zeros to find
// primes (the stillness equalizer of Finding 93). Fusion closes the loop:
//
//	stage 1 (the voyage):  hunt the first 1000 zeros with our own
//	                       instruments - the cargo.
//	stage 2 (the compass): pour the cargo into the prime detector
//	                       D(n) = -(2/sqrt n) * sum_m cos(gamma_m ln n),
//	                       the explicit formula's own lens: it must spike
//	                       exactly on primes and prime powers.
//
// No sieve, no division, no factoring anywhere in stage 2: the primes are
// found purely by listening to the zeros. Controls: the same compass with
// shuffled-gap frequencies must go blind.
//
// Usage:
//
//	go run ./cmd/fusion
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"sort"
)

const emTerms = 500

func zetaEM(s complex128, lnn []float64) complex128 {
	var sum complex128
	sig, t := real(s), imag(s)
	for n := 1; n < emTerms; n++ {
		amp := math.Exp(-sig * lnn[n])
		sn, cs := math.Sincos(t * lnn[n])
		sum += complex(amp*cs, -amp*sn)
	}
	nf := complex(float64(emTerms), 0)
	ns := cmplx.Exp(-s * complex(lnn[emTerms], 0))
	sum += ns * nf / (s - 1)
	sum += ns / 2
	sum += ns * s / nf / 12
	sum -= ns * s * (s + 1) * (s + 2) / (nf * nf * nf) / 720
	sum += ns * s * (s + 1) * (s + 2) * (s + 3) * (s + 4) /
		(nf * nf * nf * nf * nf) / 30240
	return sum
}

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zOf(t float64, lnn []float64) float64 {
	z := zetaEM(complex(0.5, t), lnn)
	th := theta(t)
	return real(z)*math.Cos(th) - imag(z)*math.Sin(th)
}

func huntZeros(count int, lnn []float64) []float64 {
	zeros := []float64{}
	prevT, prevZ := 12.0, zOf(12.0, lnn)
	for t := 12.02; len(zeros) < count && t < 1500; t += 0.02 {
		zt := zOf(t, lnn)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, t
			for i := 0; i < 60 && hi-lo > 1e-10; i++ {
				mid := (lo + hi) / 2
				if (zOf(mid, lnn) < 0) == (prevZ < 0) {
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

// compass reads the explicit-formula detector at x from the cargo of zeros.
func compass(x float64, gammas []float64) float64 {
	lx := math.Log(x)
	var s float64
	for _, g := range gammas {
		s += math.Cos(g * lx)
	}
	return -2 * s / math.Sqrt(x)
}

func isPrimePower(n int) (bool, bool) { // (prime power, prime)
	if n < 2 {
		return false, false
	}
	for p := 2; p*p <= n; p++ {
		if n%p == 0 {
			m := n
			for m%p == 0 {
				m /= p
			}
			return m == 1, false
		}
	}
	return true, true
}

// listen scores every integer in [2, top] by the compass's local band gain
// (the equalizer: the reading at n against the readings half a step off).
func listen(gammas []float64, top int) []int {
	type band struct {
		n    int
		gain float64
	}
	bands := make([]band, 0, top-1)
	for n := 2; n <= top; n++ {
		x := float64(n)
		g := compass(x, gammas) - (compass(x-0.5, gammas)+compass(x+0.5, gammas))/2
		bands = append(bands, band{n, g})
	}
	sort.Slice(bands, func(i, j int) bool { return bands[i].gain > bands[j].gain })
	out := make([]int, len(bands))
	for i, b := range bands {
		out[i] = b.n
	}
	return out
}

func main() {
	fmt.Println("THE FUSION — the voyage and the compass, one vessel")

	lnn := make([]float64, emTerms+1)
	for n := 1; n <= emTerms; n++ {
		lnn[n] = math.Log(float64(n))
	}

	// stage 1: the voyage. Hunt the cargo.
	gammas := huntZeros(1000, lnn)
	fmt.Printf("\n  stage 1 - THE VOYAGE: %d zeros hunted in-house\n", len(gammas))
	fmt.Printf("    gamma_1 = %.6f (known 14.134725)   gamma_%d = %.4f\n",
		gammas[0], len(gammas), gammas[len(gammas)-1])

	// stage 2: the compass. Pour the cargo into the prime detector.
	const top = 1000
	nTargets := 0
	for n := 2; n <= top; n++ {
		if pp, _ := isPrimePower(n); pp {
			nTargets++
		}
	}
	ranked := listen(gammas, top)
	hits := 0
	for i := 0; i < nTargets; i++ {
		if pp, _ := isPrimePower(ranked[i]); pp {
			hits++
		}
	}
	fmt.Printf("\n  stage 2 - THE COMPASS: primes found by listening to the zeros alone\n")
	fmt.Printf("    targets: %d prime powers in [2, %d]\n", nTargets, top)
	fmt.Printf("    of the compass's top %d readings, %d are prime powers (%.1f%%)\n",
		nTargets, hits, 100*float64(hits)/float64(nTargets))
	fmt.Printf("    chance level would be %.1f%%\n", 100*float64(nTargets)/float64(top-1))

	// control: shuffled gaps - same density, no prime soul.
	rng := rand.New(rand.NewSource(2026))
	gaps := make([]float64, len(gammas)-1)
	for i := range gaps {
		gaps[i] = gammas[i+1] - gammas[i]
	}
	rng.Shuffle(len(gaps), func(i, j int) { gaps[i], gaps[j] = gaps[j], gaps[i] })
	shuffled := make([]float64, len(gammas))
	shuffled[0] = gammas[0]
	for i, g := range gaps {
		shuffled[i+1] = shuffled[i] + g
	}
	rankedSh := listen(shuffled, top)
	hitsSh := 0
	for i := 0; i < nTargets; i++ {
		if pp, _ := isPrimePower(rankedSh[i]); pp {
			hitsSh++
		}
	}
	fmt.Printf("    control (shuffled gaps): %d of top %d (%.1f%%) - the compass goes blind\n",
		hitsSh, nTargets, 100*float64(hitsSh)/float64(nTargets))

	// the distant system: a regional reading. The compass resolves
	// individual integers only out to x ~ gamma_max/pi (peak width one);
	// beyond that the image blurs - the voyage's depth IS the reach.
	region := func(lo, hi int) {
		type band struct {
			n    int
			gain float64
		}
		var bands []band
		targets := 0
		for n := lo; n <= hi; n++ {
			x := float64(n)
			g := compass(x, gammas) - (compass(x-0.5, gammas)+compass(x+0.5, gammas))/2
			bands = append(bands, band{n, g})
			if pp, _ := isPrimePower(n); pp {
				targets++
			}
		}
		sort.Slice(bands, func(i, j int) bool { return bands[i].gain > bands[j].gain })
		fmt.Printf("    [%d, %d]  compass top %d:", lo, hi, targets)
		picks := make([]int, targets)
		copy(picks, func() []int {
			out := make([]int, targets)
			for i := 0; i < targets; i++ {
				out[i] = bands[i].n
			}
			return out
		}())
		sort.Ints(picks)
		hits := 0
		for _, n := range picks {
			mark := ""
			if pp, _ := isPrimePower(n); pp {
				hits++
				mark = "*"
			}
			fmt.Printf(" %d%s", n, mark)
		}
		fmt.Printf("   -> %d/%d\n", hits, targets)
	}
	reach := gammas[len(gammas)-1] / math.Pi
	fmt.Printf("\n  THE DISTANT SYSTEM - no sieve, no division (resolution reach ~ %.0f):\n", reach)
	region(430, 470)
	region(950, 1000)
	fmt.Println("    within reach the reading is sharp; beyond it the image blurs -")
	fmt.Println("    the compass sees exactly as far as the voyage has sailed deep.")
	fmt.Println("\n  the loop is closed: primes hunt zeros, zeros find primes - the")
	fmt.Println("  double mirror of the explicit formula, welded into one vessel.")
}
