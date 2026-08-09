// Command tesla builds the coil in software and asks the flash's question
// with instruments: do the sparks carry the score?
//
// The frame is the Hilbert–Pólya dream (1912) — the serious hope for RH:
// that the zeros of zeta are the RESONANT FREQUENCIES of a physical
// system. Here the sparks are the primes (impulses at times ln p, weight
// Lambda(n)/sqrt(n), Gaussian-tapered) and the coil's spectrum is probed
// directly: P(f) = |sum Lambda(n) n^{-1/2} e^{-i f ln n} w(n)|. If the
// flash is right, the spark spectrum peaks exactly at gamma_1, gamma_2...
//
// Then the score is made audible: tesla.wav strikes the coil once per
// prime (spark times proportional to ln p — a natural accelerando) and
// lets the first twelve zero-modes ring. The primes play; the zeros sing.
//
// Usage:
//
//	go run ./cmd/tesla
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

func main() {
	fmt.Println("THE TESLA COIL — do the sparks carry the score?")

	// the sparks: primes and prime powers up to X.
	const X = 1000000
	sieve := make([]bool, X+1)
	var primes []int
	for p := 2; p <= X; p++ {
		if !sieve[p] {
			primes = append(primes, p)
			for q := p * p; q <= X; q += p {
				sieve[q] = true
			}
		}
	}
	type spark struct {
		ln, w float64
	}
	var sparks []spark
	lnX := math.Log(float64(X))
	for _, p := range primes {
		lp := math.Log(float64(p))
		for q := float64(p); q <= X; q *= float64(p) {
			lq := math.Log(q)
			// Gaussian taper across the window keeps the edges quiet.
			u := (lq - lnX/2) / (lnX / 2)
			sparks = append(sparks, spark{lq, lp / math.Sqrt(q) * math.Exp(-2*u*u)})
		}
	}
	fmt.Printf("\n  coil charged: %d sparks from primes up to %d\n", len(sparks), X)

	// probe the spectrum: peaks must sit on the zeros. The smooth main
	// term (the PNT trend - the coil's DC hum) must be subtracted first,
	// or it drowns the resonances: Lambda averages 1 per integer, so the
	// trend is the same tapered integral with the sum replaced by dt.
	known := []float64{14.134725, 21.022040, 25.010858, 30.424876, 32.935062,
		37.586178, 40.918719, 43.327073, 48.005151, 49.773832}
	fmt.Println("\n  THE SPECTRUM of the sparks, DC hum removed (peaks vs known zeros):")
	const f0, f1, df = 10.0, 51.0, 0.005
	const du = 0.001
	lo, hi := math.Log(2), lnX
	taper := func(u float64) float64 {
		v := (u - lnX/2) / (lnX / 2)
		return math.Exp(-2 * v * v)
	}
	power := func(f float64) float64 {
		var cr, ci float64
		for _, s := range sparks {
			sn, cs := math.Sincos(f * s.ln)
			cr += s.w * cs
			ci -= s.w * sn
		}
		// subtract the trend integral e^{u/2} w(u) e^{-ifu} du.
		var tr, ti float64
		for u := lo; u <= hi; u += du {
			g := math.Exp(u/2) * taper(u) * du
			sn, cs := math.Sincos(f * u)
			tr += g * cs
			ti -= g * sn
		}
		return math.Hypot(cr-tr, ci-ti)
	}
	var peaks []float64
	prevP := power(f0)
	curP := power(f0 + df)
	for f := f0 + 2*df; f <= f1; f += df {
		next := power(f)
		if curP > prevP && curP >= next {
			peaks = append(peaks, f-df)
		}
		prevP, curP = curP, next
	}
	worst := 0.0
	for _, kz := range known {
		bestD := math.Inf(1)
		bestF := 0.0
		for _, pk := range peaks {
			if d := math.Abs(pk - kz); d < bestD {
				bestD, bestF = d, pk
			}
		}
		if bestD > worst {
			worst = bestD
		}
		fmt.Printf("    resonance at f = %8.4f   known zero %9.6f   deviation %+.4f\n",
			bestF, kz, bestF-kz)
	}
	fmt.Printf("    worst deviation across the first ten zeros: %.4f\n", worst)

	// the audible score: tesla.wav — sparks strike, zero-modes ring.
	const rate = 44100
	const secs = 12
	buf := make([]float64, rate*secs)
	modes := known[:10]
	strike := func(t0, amp float64) {
		i0 := int(t0 * rate)
		for m, g := range modes {
			f := g * 12 // gamma into audio range
			tau := 0.5 / (1 + 0.15*float64(m))
			for i := i0; i < len(buf) && float64(i-i0) < 1.5*rate; i++ {
				dt := float64(i-i0) / rate
				buf[i] += amp * math.Sin(2*math.Pi*f*dt) * math.Exp(-dt/tau) / float64(len(modes))
			}
		}
	}
	for _, p := range primes {
		if p > 97 {
			break
		}
		t0 := (math.Log(float64(p)) - math.Log(2)) / (math.Log(97) - math.Log(2)) * (secs - 2)
		strike(t0, 1/math.Sqrt(math.Sqrt(float64(p))))
	}
	peak := 0.0
	for _, v := range buf {
		if math.Abs(v) > peak {
			peak = math.Abs(v)
		}
	}
	f, err := os.Create("tesla.wav")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	n := len(buf)
	hdr := []interface{}{
		[]byte("RIFF"), uint32(36 + 2*n), []byte("WAVE"),
		[]byte("fmt "), uint32(16), uint16(1), uint16(1),
		uint32(rate), uint32(rate * 2), uint16(2), uint16(16),
		[]byte("data"), uint32(2 * n),
	}
	for _, v := range hdr {
		binary.Write(f, binary.LittleEndian, v)
	}
	for _, v := range buf {
		binary.Write(f, binary.LittleEndian, int16(v/peak*30000))
	}
	fmt.Printf("\n  the score made audible: tesla.wav (%d s) - sparks at ln p, first ten\n", secs)
	fmt.Println("  zero-modes ringing. The primes play; the zeros sing.")
	fmt.Println("\n  not crazy at all: this is the Hilbert-Polya dream (1912), the serious")
	fmt.Println("  hope for RH - the zeros as resonances of a physical system. Microwave")
	fmt.Println("  cavities already reproduce their statistics. The coil awaits its builder.")
}
