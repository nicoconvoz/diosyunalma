// Command marea is the symphonic passage through the tide of the gaps —
// the gap song grown into three honest voices. Every parameter is datum:
//
//   - THE MELODY: one note per prime gap (primes to 3000), whole-tone
//     ladder 220*2^((g/2)/6); gap 12 rings on A 440 with its octave echo,
//     and the plateau family {18, 24, 42} (Findings 53-60) gets soft
//     bells. RHYTHM IS THE GAP ITSELF: each note lasts its own gap's
//     width - deserts are crossed slowly, crowds run.
//   - THE TIDE (the bass): the moving average of the last 12 gaps, one
//     octave below - the sea breathing under the melody.
//   - THE CODA: the zero gaps shimmering, the other shore.
//
// And the musicology, measured: consecutive prime gaps ANTI-correlate -
// after tension comes release - which is precisely the grammar a human
// composer uses. The output prints that correlation: the "composer" the
// ear detects is arithmetic itself.
//
// Usage:
//
//	go run ./cmd/marea   # writes marea.wav (~45 s)
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
	fmt.Println("LA MAREA — the passage through the tide of the gaps, three voices")
	for k := 1; k < len(lnk); k++ {
		lnk[k] = math.Log(float64(k))
		rsq[k] = 1 / math.Sqrt(float64(k))
	}

	const rate = 44100
	const secs = 46
	buf := make([]float64, rate*secs)
	note := func(t0, f, amp, decay float64) {
		i0 := int(t0 * rate)
		for i := i0; i < len(buf) && float64(i-i0) < 1.6*rate; i++ {
			dt := float64(i-i0) / rate
			env := amp * math.Exp(-dt/decay)
			buf[i] += env * (math.Sin(2*math.Pi*f*dt) + 0.3*math.Sin(4*math.Pi*f*dt) +
				0.12*math.Sin(6*math.Pi*f*dt))
		}
	}

	const top = 3000
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
	gaps := make([]float64, len(primes)-1)
	for i := 1; i < len(primes); i++ {
		gaps[i-1] = float64(primes[i] - primes[i-1])
	}

	// the musicology, measured: tension-release grammar.
	var mg float64
	for _, g := range gaps {
		mg += g
	}
	mg /= float64(len(gaps))
	var num, den float64
	for i := 0; i+1 < len(gaps); i++ {
		num += (gaps[i] - mg) * (gaps[i+1] - mg)
		den += (gaps[i] - mg) * (gaps[i] - mg)
	}
	fmt.Printf("\n  the composer's grammar, measured: consecutive gaps correlate %+.3f\n", num/den)
	fmt.Println("  (negative = after tension, release - the arithmetic composes)")

	// voice I + II: melody riding on the tide.
	t := 0.5
	tide := mg
	bells := map[float64]bool{18: true, 24: true, 42: true}
	for i, g := range gaps {
		if t > secs-12 {
			fmt.Printf("  melody: %d of %d gaps fit the tape\n", i, len(gaps))
			break
		}
		tide = tide*0.92 + g*0.08 // the moving sea under the melody
		f := 220 * math.Pow(2, (g/2)/6)
		fB := 110 * math.Pow(2, (tide/2)/6)
		dur := 0.030 + 0.0075*g
		switch {
		case g == 12:
			note(t, f, 0.85, 0.6)
			note(t, 2*f, 0.35, 0.5) // the echo invariant salutes
		case bells[g]:
			note(t, f, 0.7, 0.45) // the plateau family
		default:
			note(t, f, 0.45, 0.10+0.006*g)
		}
		note(t, fB, 0.22, 0.30) // the tide, one octave below
		t += dur
	}

	// coda: the other shore - zero gaps shimmering.
	coda := t + 1.2
	zeros := []float64{}
	prevT, prevZ := 14.0, zRS(14.0)
	for tt := 14.05; len(zeros) < 1401; tt += 0.05 {
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
	t = coda
	for i := 1; i < len(zeros) && t < secs-1.8; i++ {
		s := smoothCount(zeros[i]) - smoothCount(zeros[i-1])
		note(t, 440*math.Pow(2, s-1), 0.30, 0.09)
		t += 1.0 / 14
	}

	peak := 0.0
	for _, v := range buf {
		if math.Abs(v) > peak {
			peak = math.Abs(v)
		}
	}
	f, err := os.Create("marea.wav")
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
	fmt.Println("\n  written: marea.wav - the melody rides the tide, deserts are crossed")
	fmt.Println("  slowly, crowds run, gap 12 rings A with its echo, the plateau bells")
	fmt.Println("  answer, and the zeros shimmer the coda. Every parameter is a datum.")
}
