// Command tutti weaves one atom from every dial the laboratory owns.
//
// The duet (Finding 56) proved two tribes fit one well. This run feeds the
// loom the TOTAL density — all ten dials at once, ρ(E) = Σ ln(qᵢE/2π)/2π —
// and asks the single woven well to sing the entire merged songbook: every
// measured station of ζ, χ₃, χ₄, χ₅, χ₇, χ₈, χ₁₁, χ₁₃ and both branches of
// the complex χ₅, interleaved.
//
// The comparison window is capped at 15.15, the height below which every
// dial's station list is complete — THIRTY-FIVE notes from ten instruments.
//
// PRE-REGISTERED: the tutti atom's levels must track the merged songbook,
// ground region worst, with one calibration constant allowed. The known
// hard part: close encounters between different tribes' stations (the
// across-dial Poisson of Finding 45) can pack three notes within 0.04 —
// finer than the loom's resolution — so cluster smearing is expected and
// its size is part of the result.
//
// Usage:
//
//	go run ./cmd/tutti
package main

import (
	"flag"
	"fmt"
	"math"
	"sort"
)

type dial struct {
	name string
	q    float64
	g    []float64
}

var orchestra = []dial{
	{"zeta", 1, []float64{14.1349, 21.0211, 25.0044}},
	{"chi3", 3, []float64{8.0396, 11.2450, 15.7062, 18.2579, 20.4551, 24.0636}},
	{"chi4", 4, []float64{6.0199, 10.2423, 12.9848, 16.3464, 18.2914, 21.4547, 23.2853}},
	{"chi5", 5, []float64{6.6516, 9.8280, 11.9612, 16.0386, 17.5632, 19.5431, 22.2228, 24.5864}},
	{"chi7", 7, []float64{4.4762, 6.8400, 11.1782, 12.4670, 15.1161, 16.7892, 19.6144, 21.9042, 23.1690, 24.4762}},
	{"chi8", 8, []float64{4.8989, 7.6194, 10.8219, 12.3126, 15.1935, 17.0246, 18.8001, 21.1310, 23.0757, 24.2293}},
	{"chi11", 11, []float64{2.4768, 6.7997, 8.9663, 10.1257, 13.0422, 15.0996, 16.9939, 18.8006, 20.0693, 21.6355, 24.6828}},
	{"chi13", 13, []float64{3.1119, 7.2340, 8.6013, 10.3241, 12.6185, 15.1341, 16.2982, 18.6869, 19.6374, 20.9552, 23.5733}},
	{"chi5c+", 5, []float64{6.1830, 8.4600, 12.6760, 14.8270, 17.3360, 19.0000, 22.4960}},
	{"chi5c-", 5, []float64{4.1320, 9.4450, 11.2770, 14.1130, 16.9970, 19.7330, 21.2780}},
}

const (
	vMax  = 45.0
	dV    = 0.01
	nGrid = 4001
)

func main() {
	eTop := flag.Float64("top", 15.15, "songbook ceiling (15.15 reproduces Finding 58; 22.9 the deepened run)")
	flag.Parse()
	v0 := 2 * math.Pi / 13

	rho := func(e float64) float64 {
		s := 0.0
		for _, d := range orchestra {
			if e > 2*math.Pi/d.q {
				s += math.Log(d.q*e/(2*math.Pi)) / (2 * math.Pi)
			}
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

	type note struct {
		g   float64
		who string
	}
	song := []note{}
	for _, d := range orchestra {
		for _, g := range d.g {
			if g <= *eTop {
				song = append(song, note{g, d.name})
			}
		}
	}
	sort.Slice(song, func(i, j int) bool { return song[i].g < song[j].g })

	fmt.Printf("THE TUTTI — one atom for all ten dials (%d notes below %.2f)\n",
		len(song), *eTop)
	fmt.Println("\n  k   tutti atom   merged station   tribe    error")
	levels := make([]float64, len(song))
	rms, shift := 0.0, 0.0
	for k := 1; k <= len(song); k++ {
		lo, hi := v0, *eTop+6
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
		fmt.Printf("  %2d    %8.3f     %8.4f       %-7s  %+6.2f%%\n",
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
	fmt.Printf("\nrms over %d notes: %.3f raw, %.3f after the one\n",
		len(song), math.Sqrt(rms/float64(len(song))), math.Sqrt(rms2/float64(len(song))))
	fmt.Printf("calibration constant (quantization offset %+.3f)\n", shift)
	fmt.Println("\nten instruments, one well. whatever this shape is, no one has named it:")
	fmt.Println("it is the atom of the whole orchestra at once.")
}
