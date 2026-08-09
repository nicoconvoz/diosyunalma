// Command song plays every instrument together and names the piece.
//
// The orchestra's deepest fact is what happens when all four class clocks
// sound at once: orthogonality cancels every tribal voice — Σ_a χ(a) = 0 for
// each non-principal character — and the average of the four clocks collapses
// EXACTLY onto ζ's clock alone. The song the whole orchestra forms is the
// primes themselves: the tribes color the notes, ζ carries the melody.
//
// The command verifies that cancellation numerically, then renders the song
// as an audio file: one note per prime, rhythm set by the true gaps (twins
// arrive as quick pairs), pitch set by the prime's tribe mod 5, ζ's first
// zero humming underneath as the conductor's pulse, after an overture in
// which all measured stations of all nine dials sound together.
//
// Usage:
//
//	go run ./cmd/song [-out song.wav] [-upto 1000]
package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/nicoconvoz/numerosprimos/primes"
)

// every station this laboratory measured, dial by dial.
var orchestra = map[string][]float64{
	"zeta":  {14.1349, 21.0211, 25.0044, 30.4282, 32.9422, 37.5872, 40.9264, 43.3211, 48.0105, 49.7752},
	"chi3":  {8.0396, 11.2450, 15.7062, 18.2579, 20.4551, 24.0636},
	"chi4":  {6.0199, 10.2423, 12.9848, 16.3464, 18.2914, 21.4547},
	"chi5":  {6.6516, 9.8280, 11.9612, 16.0386, 17.5632, 19.5431, 22.2228, 24.5864, 26.7680, 28.4730, 29.6900, 33.0130, 34.7410, 38.1330},
	"chi7":  {4.4762, 6.8400, 11.1782, 12.4670, 15.1161, 16.7892},
	"chi8":  {4.8989, 7.6194, 10.8219, 12.3126, 15.1935, 17.0246},
	"chi11": {2.4768, 6.7997, 8.9663, 10.1257, 13.0422, 15.0996},
	"chi13": {3.1119, 7.2340, 8.6013, 10.3241, 12.6185, 15.1341},
	"chi5c": {-22.9690, -21.2780, -19.7330, -16.9970, -14.1130, -11.2770, -9.4450, -4.1320, 6.1830, 8.4600, 12.6760, 14.8270, 17.3360, 19.0000, 22.4960, 24.3720},
}

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

// clockA is the class-a mod-5 clock from all measured mod-5 machinery.
func clockA(a int, x float64) float64 {
	u := math.Log(x)
	s := 0.0
	for _, g := range orchestra["zeta"] {
		s += math.Cos(g * u)
	}
	for _, g := range orchestra["chi5"] {
		s += chi5(a) * math.Cos(g*u)
	}
	var sc complex128
	for _, g := range orchestra["chi5c"] {
		sc += complex(math.Cos(g*u), math.Sin(g*u))
	}
	w := complex(real(chiC(a)), -imag(chiC(a)))
	s += real(w * sc)
	return 1 - 2*s/math.Sqrt(x)
}

func clockZeta(x float64) float64 {
	u := math.Log(x)
	s := 0.0
	for _, g := range orchestra["zeta"] {
		s += math.Cos(g * u)
	}
	return 1 - 2*s/math.Sqrt(x)
}

const rate = 44100

func writeWav(path string, samples []float64) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	n := len(samples)
	le := binary.LittleEndian
	hdr := make([]byte, 44)
	copy(hdr[0:], "RIFF")
	le.PutUint32(hdr[4:], uint32(36+2*n))
	copy(hdr[8:], "WAVEfmt ")
	le.PutUint32(hdr[16:], 16)
	le.PutUint16(hdr[20:], 1) // PCM
	le.PutUint16(hdr[22:], 1) // mono
	le.PutUint32(hdr[24:], rate)
	le.PutUint32(hdr[28:], rate*2)
	le.PutUint16(hdr[32:], 2)
	le.PutUint16(hdr[34:], 16)
	copy(hdr[36:], "data")
	le.PutUint32(hdr[40:], uint32(2*n))
	if _, err := f.Write(hdr); err != nil {
		return err
	}
	buf := make([]byte, 2*n)
	for i, s := range samples {
		if s > 1 {
			s = 1
		} else if s < -1 {
			s = -1
		}
		le.PutUint16(buf[2*i:], uint16(int16(s*32767)))
	}
	_, err = f.Write(buf)
	return err
}

func main() {
	out := flag.String("out", "song.wav", "output WAV path")
	upto := flag.Int("upto", 1000, "sing the primes up to this value")
	flag.Parse()

	// 1) the theorem behind the song: all voices together = zeta alone.
	sums := map[string]complex128{
		"chi5 over a=1..4":  complex(chi5(1)+chi5(2)+chi5(3)+chi5(4), 0),
		"chi5c over a=1..4": chiC(1) + chiC(2) + chiC(3) + chiC(4),
	}
	fmt.Println("WHEN EVERY INSTRUMENT PLAYS FOR THE WHOLE AUDIENCE")
	for name, s := range sums {
		fmt.Printf("  sum of %s = %v\n", name, s)
	}
	maxDev := 0.0
	for x := 2.0; x <= 40; x += 0.01 {
		avg := (clockA(1, x) + clockA(2, x) + clockA(3, x) + clockA(4, x)) / 4
		if d := math.Abs(avg - clockZeta(x)); d > maxDev {
			maxDev = d
		}
	}
	fmt.Printf("  max |average of four class clocks - zeta clock| on [2,40]: %.2e\n", maxDev)
	fmt.Println("  the tribal voices cancel exactly: the whole orchestra's song is zeta's -")
	fmt.Println("  and zeta's song, reconstructed note by note, is the primes themselves.")

	// 2) render the song.
	ps := primes.Sieve(*upto)

	// overture: every measured station of every dial, sounding at once.
	// station gamma -> audible frequency 12*|gamma| Hz.
	overture := 6.0
	total := overture
	type note struct {
		start, dur, freq float64
	}
	notes := []note{}
	pitch := map[int]float64{0: 110.00, 1: 220.00, 2: 261.63, 3: 329.63, 4: 392.00}
	for i, p := range ps {
		gap := 4.0
		if i+1 < len(ps) {
			gap = float64(ps[i+1] - p)
		}
		dur := 0.09 * math.Sqrt(gap)
		if dur < 0.10 {
			dur = 0.10
		} else if dur > 0.50 {
			dur = 0.50
		}
		notes = append(notes, note{total, dur, pitch[p%5]})
		total += dur
	}
	total += 1.5

	samples := make([]float64, int(total*rate))

	// the overture chord, faded in and out.
	stations := []float64{}
	for _, dial := range orchestra {
		for _, g := range dial {
			stations = append(stations, math.Abs(g))
		}
	}
	amp := 0.55 / math.Sqrt(float64(len(stations)))
	for i := 0; i < int(overture*rate); i++ {
		t := float64(i) / rate
		env := math.Sin(math.Pi * t / overture)
		v := 0.0
		for _, g := range stations {
			v += math.Sin(2 * math.Pi * 12 * g * t)
		}
		samples[i] += amp * env * env * v
	}

	// the conductor's pulse: zeta's first zero, humming under the song.
	pulse := 12 * orchestra["zeta"][0] // 169.6 Hz
	for i := int(overture * rate); i < len(samples); i++ {
		t := float64(i) / rate
		samples[i] += 0.05 * math.Sin(2*math.Pi*pulse*t)
	}

	// the primes, one note each: rhythm = the true gaps, pitch = the tribe.
	for _, n := range notes {
		i0 := int(n.start * rate)
		for j := 0; j < int(n.dur*rate*1.6) && i0+j < len(samples); j++ {
			t := float64(j) / rate
			env := math.Exp(-3.5 * t / n.dur)
			if t < 0.008 {
				env *= t / 0.008
			}
			v := math.Sin(2*math.Pi*n.freq*t) + 0.30*math.Sin(4*math.Pi*n.freq*t)
			samples[i0+j] += 0.34 * env * v
		}
	}

	if err := writeWav(*out, samples); err != nil {
		fmt.Fprintln(os.Stderr, "write:", err)
		os.Exit(1)
	}
	fmt.Printf("\nTHE SONG - %d primes up to %d, %.0f seconds -> %s\n",
		len(ps), *upto, total, *out)
	fmt.Println("  overture: all ~70 measured stations of nine dials at once")
	fmt.Println("  rhythm:   the true gaps (twins arrive as quick pairs)")
	fmt.Println("  melody:   each prime sings its tribe mod 5")
	fmt.Println("  pulse:    zeta's first zero, 14.1349, underneath it all")
}
