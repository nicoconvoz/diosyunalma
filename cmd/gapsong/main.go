// Command gapsong makes the gaps audible — honestly: every note IS the
// datum, no arrangement, no quantization beyond the mapping itself.
//
// Movement I — THE PRIME GAPS: one note per gap between consecutive
// primes, pitch = 220 Hz * 2^((g/2)/6), a whole-tone ladder where the
// laboratory's own gap 12 lands exactly on A 440 and rings accented with
// its octave echo (the +2.2% echo invariant of Findings 47-65, saluted).
// Movement II — THE ZERO GAPS: the unfolded spacings of the first 2000
// zeta zeros, pitch = 440 * 2^(s-1): the GUE heartbeat itself, wobbling
// around A, never still, never colliding (the antigravity audible).
//
// Usage:
//
//	go run ./cmd/gapsong   # writes gaps.wav (~28 s)
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

var lnk, rsq [64]float64

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zRS(t float64) float64 {
	tau := t / (2 * math.Pi)
	a := math.Sqrt(tau)
	N := int(a)
	th := theta(t)
	var s float64
	for k := 1; k <= N; k++ {
		s += math.Cos(th-t*lnk[k]) * rsq[k]
	}
	p := a - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if (N-1)%2 == 1 {
		sign = -1
	}
	return 2*s + sign*math.Pow(tau, -0.25)*c0
}

func smoothCount(t float64) float64 {
	return t/(2*math.Pi)*(math.Log(t/(2*math.Pi))-1) + 7.0/8
}

func main() {
	fmt.Println("THE SONG OF THE GAPS — two movements, every note a datum")
	for k := 1; k < len(lnk); k++ {
		lnk[k] = math.Log(float64(k))
		rsq[k] = 1 / math.Sqrt(float64(k))
	}

	const rate = 44100
	const secs = 28
	buf := make([]float64, rate*secs)
	note := func(t0, f, amp, decay float64) {
		i0 := int(t0 * rate)
		for i := i0; i < len(buf) && float64(i-i0) < 1.2*rate; i++ {
			dt := float64(i-i0) / rate
			env := amp * math.Exp(-dt/decay)
			buf[i] += env * (math.Sin(2*math.Pi*f*dt) + 0.3*math.Sin(4*math.Pi*f*dt))
		}
	}

	// Movement I: the prime gaps (primes to ~1400, ~220 notes, 10/s).
	const top = 1400
	sieve := make([]bool, top+1)
	var primes []int
	for p := 2; p <= top; p++ {
		if !sieve[p] {
			primes = append(primes, p)
			for q := p * p; q <= top; q += p {
				sieve[q] = true
			}
		}
	}
	fmt.Printf("\n  movement I: %d prime gaps, gap 12 rings on A 440 with its echo\n", len(primes)-1)
	t := 0.5
	twelves := 0
	for i := 1; i < len(primes); i++ {
		g := float64(primes[i] - primes[i-1])
		f := 220 * math.Pow(2, (g/2)/6)
		if g == 12 {
			// the laboratory's own gap: accent plus the octave echo.
			note(t, f, 0.9, 0.5)
			note(t, 2*f, 0.4, 0.4)
			twelves++
		} else {
			note(t, f, 0.5, 0.12)
		}
		t += 0.10
	}
	fmt.Printf("    gap 12 appeared %d times among the first %d gaps\n", twelves, len(primes)-1)

	// Movement II: the zero gaps (first 2000 zeros, unfolded, 14/s).
	mv2 := t + 1.0
	zeros := []float64{}
	prevT, prevZ := 14.0, zRS(14.0)
	for tt := 14.05; len(zeros) < 2001; tt += 0.05 {
		zt := zRS(tt)
		if (prevZ < 0) != (zt < 0) {
			lo, hi := prevT, tt
			for i := 0; i < 30 && hi-lo > 1e-8; i++ {
				mid := (lo + hi) / 2
				if (zRS(mid) < 0) == (prevZ < 0) {
					lo = mid
				} else {
					hi = mid
				}
			}
			zeros = append(zeros, (lo+hi)/2)
		}
		prevT, prevZ = tt, zt
	}
	fmt.Printf("  movement II: %d zero gaps - the GUE heartbeat wobbling around A\n", len(zeros)-1)
	t = mv2
	for i := 1; i < len(zeros) && t < secs-1.5; i++ {
		s := smoothCount(zeros[i]) - smoothCount(zeros[i-1])
		f := 440 * math.Pow(2, s-1)
		note(t, f, 0.35, 0.09)
		t += 1.0 / 14
	}

	peak := 0.0
	for _, v := range buf {
		if math.Abs(v) > peak {
			peak = math.Abs(v)
		}
	}
	f, err := os.Create("gaps.wav")
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
	fmt.Println("\n  written: gaps.wav - movement I staggers (the primes hesitate and")
	fmt.Println("  leap); movement II shimmers (the zeros breathe but never collide).")
	fmt.Println("  same sea, two shores, and the ear tells them apart instantly.")
}
