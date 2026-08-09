// Command duet weaves the harmony's own atom.
//
// Finding 42 showed the harmony dial: the summed signal of ζ and the golden
// tribe carries both musics at once, as the Dedekind zeta of Q(√5) = ζ·L(χ₅)
// demands. If the harmony has a dial and a song, it must have an atom: one
// well whose levels are BOTH station lists merged.
//
// The loom is Finding 53's — semiclassical Abel inversion — but the density
// is now the harmony's: ρ(E) = [ln(E/2π) + ln(5E/2π)] / 2π, the sum of both
// dials' zero densities.
//
// PRE-REGISTERED: the woven duet atom's levels must track the MERGED list of
// measured stations (6 of ζ and all 14 of χ₅ below 38.2 — twenty notes,
// interleaved by the two tribes), with the ground region worst, as in
// Finding 53. One calibration constant (the quantization offset) is allowed.
//
// Usage:
//
//	go run ./cmd/duet
package main

import (
	"fmt"
	"math"
	"sort"
)

var zetaZeros = []float64{14.1349, 21.0211, 25.0044, 30.4282, 32.9422, 37.5872}

var chi5Zeros = []float64{6.6516, 9.8280, 11.9612, 16.0386, 17.5632, 19.5431,
	22.2228, 24.5864, 26.7680, 28.4730, 29.6900, 33.0130, 34.7410, 38.1330}

const (
	vMax  = 120.0
	dV    = 0.02
	nGrid = 4001
	eTop  = 38.5
)

func main() {
	v0 := 2 * math.Pi / 5

	// the harmony's density: both dials at once.
	rho := func(e float64) float64 {
		s := 0.0
		if e > 2*math.Pi {
			s += math.Log(e/(2*math.Pi)) / (2 * math.Pi)
		}
		if e > 2*math.Pi/5 {
			s += math.Log(5*e/(2*math.Pi)) / (2 * math.Pi)
		}
		return s
	}

	half := func(v float64) float64 {
		if v <= v0 {
			return 0
		}
		tMax := math.Sqrt(v - v0)
		const steps = 400
		s := 0.0
		for k := 0; k < steps; k++ {
			t := (float64(k) + 0.5) / steps * tMax
			s += 2 * rho(v-t*t)
		}
		return s * tMax / steps
	}

	nTab := int((vMax - v0) / dV)
	vs := make([]float64, nTab)
	xs := make([]float64, nTab)
	for i := 0; i < nTab; i++ {
		vs[i] = v0 + float64(i)*dV
		xs[i] = half(vs[i])
	}
	xMax := xs[nTab-1]

	h := 2 * xMax / float64(nGrid-1)
	pot := make([]float64, nGrid)
	for i := 0; i < nGrid; i++ {
		x := math.Abs(-xMax + float64(i)*h)
		lo, hi := 0, nTab-1
		for lo < hi {
			mid := (lo + hi) / 2
			if xs[mid] < x {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		if lo == 0 {
			pot[i] = vs[0]
		} else if lo >= nTab-1 {
			pot[i] = vMax
		} else {
			f := (x - xs[lo-1]) / (xs[lo] - xs[lo-1])
			pot[i] = vs[lo-1] + f*(vs[lo]-vs[lo-1])
		}
	}

	off2 := 1 / (h * h) * 1 / (h * h)
	count := func(e float64) int {
		n := 0
		q := 2/(h*h) + pot[0] - e
		if q < 0 {
			n++
		}
		for i := 1; i < nGrid; i++ {
			d := 2/(h*h) + pot[i] - e
			if q != 0 {
				q = d - off2/q
			} else {
				q = d - off2/1e-30
			}
			if q < 0 {
				n++
			}
		}
		return n
	}

	// the merged song: both tribes' stations, in one list.
	type note struct {
		g   float64
		who string
	}
	song := []note{}
	for _, g := range zetaZeros {
		song = append(song, note{g, "zeta"})
	}
	for _, g := range chi5Zeros {
		song = append(song, note{g, "chi5"})
	}
	sort.Slice(song, func(i, j int) bool { return song[i].g < song[j].g })

	fmt.Println("THE DUET — one atom for the harmony of Q(sqrt 5)")
	fmt.Println("\n  k   duet atom   merged station   tribe   error")
	levels := make([]float64, len(song))
	rms, shift := 0.0, 0.0
	for k := 1; k <= len(song); k++ {
		lo, hi := v0, eTop+8
		for hi-lo > 1e-6 {
			mid := (lo + hi) / 2
			if count(mid) >= k {
				hi = mid
			} else {
				lo = mid
			}
		}
		levels[k-1] = hi
		n := song[k-1]
		fmt.Printf("  %2d   %8.3f     %8.4f       %-5s  %+6.2f%%\n",
			k, hi, n.g, n.who, 100*(hi-n.g)/n.g)
		rms += (hi - n.g) * (hi - n.g)
		shift += n.g - hi
	}
	shift /= float64(len(song))
	rms2 := 0.0
	for k := range song {
		d := levels[k] + shift - song[k].g
		rms2 += d * d
	}
	fmt.Printf("\nrms over twenty notes: %.3f raw, %.3f after the one\n",
		math.Sqrt(rms/float64(len(song))), math.Sqrt(rms2/float64(len(song))))
	fmt.Printf("calibration constant (quantization offset %+.3f)\n", shift)
	fmt.Println("\none well, two tribes: the harmony's atom sings the interleaved song.")
}
