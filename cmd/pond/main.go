// Command pond drops the rock into two ponds and listens.
//
// The flash behind it: ripples that bounce off the walls and return against
// themselves look like chaos, yet hold a complex harmony. Wave chaos has a
// standard laboratory: the Bunimovich stadium billiard, whose classical
// trajectories are chaotic, against a circular pond, whose trajectories are
// regular. The modes of each pond are the eigenvalues of −∇²ψ = Eψ with
// ψ = 0 on the shore, computed here by finite differences on desymmetrized
// (quarter) domains of equal area, with banded LDLᵀ Sturm counting.
//
// PRE-REGISTERED:
//   - the chaotic pond's unfolded spacings must REPEL — variance near the
//     GOE value 0.27, small gaps near 18% — the harmony hidden in the chaos;
//   - the circular pond must relax toward Poisson (variance near 1,
//     small gaps near 39%) — order, paradoxically, gives indifference;
//   - NEITHER pond may match the stations' GUE (variance 0.178): water
//     respects time reversal, so its repulsion is the GOE kind. The primes'
//     stations repel in the GUE class — the fingerprint of a hidden system
//     that BREAKS time reversal, as Berry–Keating predicts.
//
// Usage:
//
//	go run ./cmd/pond
package main

import (
	"fmt"
	"math"
)

const (
	h      = 0.025
	levels = 130
	skip   = 15
)

// domain builds the interior-point index map for a quarter shape.
type domain struct {
	name   string
	inside func(x, y float64) bool
	xMax   float64
	yMax   float64
}

func stats(sp []float64) (mean, variance, small float64) {
	for _, s := range sp {
		mean += s
	}
	mean /= float64(len(sp))
	n := 0
	for _, s := range sp {
		variance += (s - mean) * (s - mean)
		if s < 0.5*mean {
			n++
		}
	}
	variance /= float64(len(sp)) * mean * mean
	small = float64(n) / float64(len(sp))
	return
}

func modes(dm domain) []float64 {
	// enumerate interior points column-major (x outer, y inner) so the
	// banded matrix's bandwidth is one column's height.
	type pt struct{ i, j int }
	idx := map[pt]int{}
	pts := []pt{}
	nx := int(dm.xMax/h) + 2
	ny := int(dm.yMax/h) + 2
	for i := 1; i < nx; i++ {
		for j := 1; j < ny; j++ {
			x, y := float64(i)*h, float64(j)*h
			if dm.inside(x, y) {
				idx[pt{i, j}] = len(pts)
				pts = append(pts, pt{i, j})
			}
		}
	}
	n := len(pts)
	band := 0
	for k, p := range pts {
		for _, nb := range []pt{{p.i - 1, p.j}, {p.i, p.j - 1}} {
			if m, ok := idx[nb]; ok {
				if k-m > band {
					band = k - m
				}
			}
		}
	}

	// banded symmetric storage: a[k][0] = diagonal, a[k][d] = entry (k, k-d).
	a := make([][]float64, n)
	inv := 1 / (h * h)
	for k, p := range pts {
		a[k] = make([]float64, band+1)
		a[k][0] = 4 * inv
		for _, nb := range []pt{{p.i - 1, p.j}, {p.i, p.j - 1}} {
			if m, ok := idx[nb]; ok {
				a[k][k-m] = -inv
			}
		}
	}

	// Sturm count via banded LDL^T factorization of A - sigma*I.
	l := make([][]float64, n)
	d := make([]float64, n)
	for k := range l {
		l[k] = make([]float64, band+1)
	}
	count := func(sigma float64) int {
		neg := 0
		for j := 0; j < n; j++ {
			lo := j - band
			if lo < 0 {
				lo = 0
			}
			dj := a[j][0] - sigma
			for k := lo; k < j; k++ {
				dj -= l[j][j-k] * l[j][j-k] * d[k]
			}
			if dj == 0 {
				dj = 1e-30
			}
			d[j] = dj
			if dj < 0 {
				neg++
			}
			hi := j + band
			if hi >= n {
				hi = n - 1
			}
			for i := j + 1; i <= hi; i++ {
				s := 0.0
				if i-j <= band {
					s = a[i][i-j]
				}
				lo2 := i - band
				if lo2 < 0 {
					lo2 = 0
				}
				for k := lo2; k < j; k++ {
					s -= l[i][i-k] * l[j][j-k] * d[k]
				}
				l[i][i-j] = s / dj
			}
		}
		return neg
	}

	out := make([]float64, levels)
	for k := 1; k <= levels; k++ {
		lo, hi := 0.0, 4000.0
		for hi-lo > 1e-3 {
			mid := (lo + hi) / 2
			if count(mid) >= k {
				hi = mid
			} else {
				lo = mid
			}
		}
		out[k-1] = hi
	}
	return out
}

func unfoldedSpacings(es []float64) []float64 {
	// Weyl-form least squares N(E) = c0 + c1*E + c2*sqrt(E).
	var m [3][4]float64
	for k, e := range es {
		row := [3]float64{1, e, math.Sqrt(e)}
		y := float64(k + 1)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				m[i][j] += row[i] * row[j]
			}
			m[i][3] += row[i] * y
		}
	}
	for c := 0; c < 3; c++ {
		p := m[c][c]
		for j := c; j < 4; j++ {
			m[c][j] /= p
		}
		for r := 0; r < 3; r++ {
			if r != c {
				f := m[r][c]
				for j := c; j < 4; j++ {
					m[r][j] -= f * m[c][j]
				}
			}
		}
	}
	nb := func(e float64) float64 {
		return m[0][3] + m[1][3]*e + m[2][3]*math.Sqrt(e)
	}
	sp := []float64{}
	for k := skip; k+1 < len(es); k++ {
		sp = append(sp, nb(es[k+1])-nb(es[k]))
	}
	return sp
}

func main() {
	fmt.Println("THE POND — the rock dropped into chaos and into order")

	area := 1.8 + math.Pi/4
	rCirc := math.Sqrt(4 * area / math.Pi)
	ponds := []domain{
		{"stadium (chaotic)", func(x, y float64) bool {
			if y >= 1 {
				return false
			}
			if x <= 1.8 {
				return x > 0
			}
			dx := x - 1.8
			return dx*dx+y*y < 1
		}, 2.8, 1.0},
		{"circle (regular)", func(x, y float64) bool {
			return x*x+y*y < rCirc*rCirc
		}, rCirc, rCirc},
	}

	for _, p := range ponds {
		es := modes(p)
		sp := unfoldedSpacings(es)
		mean, v, small := stats(sp)
		fmt.Printf("\n%s — %d modes, %d spacings kept\n", p.name, levels, len(sp))
		fmt.Printf("  mean %.3f   var %.3f   frac < half-mean %.1f%%\n", mean, v, 100*small)
	}

	fmt.Println("\nreference rhythms:")
	fmt.Println("  GOE (waves in a pond, time-reversal kept):  var 0.273   small-gap 18%")
	fmt.Println("  GUE (the stations, time-reversal broken):   var 0.178   small-gap 11%")
	fmt.Println("  Poisson (order, indifference):              var 1.000   small-gap 39%")
	fmt.Println("\nthe rock's ripples hold harmony in their chaos - and the primes'")
	fmt.Println("harmony is one shade stricter: their hidden pond breaks time reversal.")
}
