// Command tresregimenes tests the captain's Phase VI intuition (Auditoria/39),
// in his own words:
//
//	"Podria tener tres modos: voltaje alto, amperaje alto o una mezcla
//	 equilibrada de ambos."
//
// The auditor's §5 forbids taking the metaphor literally: we must find the real
// mathematical variables, if they exist. They do, and Phase V had them tied
// together without noticing.
//
// WHAT PHASE V ACTUALLY CONTROLLED. There the perturbation was a sum of K rank
// one channels with the TOTAL force held fixed. That single constraint locks two
// different things at once:
//
//	A - how hard one pair of levels pushes on another (the captain's voltage)
//	B - how many pairs are connected at all, i.e. the reach (his amperage)
//
// With the trace fixed, buying reach was paid for in push, which is exactly the
// trade-off Phase V measured and could not escape. THE INTUITION IS THE
// DIAGNOSIS: one knob where there are two.
//
// THE MODEL. Off the diagonal,
//
//	H_ij = A * amp_i * amp_j * exp(-|i-j|/B)
//
// A is the amplitude, B the reach in units of level spacings, and the arithmetic
// enters only through amp = Lambda(p)/sqrt(p). The total force is now a DERIVED
// quantity, not the control - which is what makes the decisive test possible:
//
//	compare points with the SAME total force and different (A,B). If the
//	measurement only depends on the total, the hypothesis loses. If it depends
//	on how the force is split, the captain is right.
//
// R6: the only arithmetic input is Lambda(n). The gamma_n are read at the end,
// as a ruler, and never to choose A, B or the window.
//
// Reproduce: go run ./cmd/tresregimenes
package main

import (
	"math"
	"sort"
)

func primos(tope int) []int {
	criba := make([]bool, tope+1)
	var ps []int
	for p := 2; p <= tope; p++ {
		if criba[p] {
			continue
		}
		ps = append(ps, p)
		for m := 2 * p; m <= tope; m += p {
			criba[m] = true
		}
	}
	return ps
}

type modo struct {
	w, amp float64
	p      int
}

func medio(topeP, n int, t0 float64) []modo {
	var ms []modo
	for _, p := range primos(topeP) {
		lp := math.Log(float64(p))
		amp := lp / math.Sqrt(float64(p))
		for k := 1; ; k++ {
			w := 2 * math.Pi * float64(k) / lp
			if w > t0*40 {
				break
			}
			if w >= t0 {
				ms = append(ms, modo{w, amp, p})
			}
		}
	}
	sort.Slice(ms, func(a, b int) bool { return ms[a].w < ms[b].w })
	if len(ms) > n {
		ms = ms[:n]
	}
	return ms
}

// ---------------------------------------------------------------------------
// The two knobs, at last separated
// ---------------------------------------------------------------------------

// fuerzaTotal is the Frobenius norm of the off-diagonal perturbation: the single
// quantity Phase V held fixed. Here it is DERIVED from (A,B), so two different
// splits can be compared at equal total.
func fuerzaTotal(ms []modo, A, B float64) float64 {
	s := 0.0
	for i := range ms {
		for j := range ms {
			if i == j {
				continue
			}
			x := A * ms[i].amp * ms[j].amp * math.Exp(-math.Abs(float64(i-j))/B)
			s += x * x
		}
	}
	return math.Sqrt(s)
}

func espectro(ms []modo, A, B float64) []float64 {
	n := len(ms)
	H := make([][]float64, n)
	for i := range H {
		H[i] = make([]float64, n)
	}
	for i := range ms {
		H[i][i] = ms[i].w
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := float64(j - i)
			if d > 12*B { // the tail is numerically dead well before this
				break
			}
			x := A * ms[i].amp * ms[j].amp * math.Exp(-d/B)
			H[i][j] += x
			H[j][i] += x
		}
	}
	ev := jacobi(H, 20)
	// the states pushed outside the band are collective doorways, not the
	// medium's response - the trap Phase V fell into on its first run
	lo, hi := ms[0].w, ms[len(ms)-1].w
	dentro := ev[:0]
	for _, e := range ev {
		if e >= lo && e <= hi {
			dentro = append(dentro, e)
		}
	}
	return dentro
}

func jacobi(a [][]float64, barridos int) []float64 {
	n := len(a)
	for sw := 0; sw < barridos; sw++ {
		off := 0.0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				off += a[i][j] * a[i][j]
			}
		}
		if off < 1e-20 {
			break
		}
		for p := 0; p < n-1; p++ {
			for q := p + 1; q < n; q++ {
				if math.Abs(a[p][q]) < 1e-16 {
					continue
				}
				th := (a[q][q] - a[p][p]) / (2 * a[p][q])
				t := math.Copysign(1, th) / (math.Abs(th) + math.Sqrt(th*th+1))
				c := 1 / math.Sqrt(t*t+1)
				s := t * c
				for k := 0; k < n; k++ {
					akp, akq := a[k][p], a[k][q]
					a[k][p] = c*akp - s*akq
					a[k][q] = s*akp + c*akq
				}
				for k := 0; k < n; k++ {
					apk, aqk := a[p][k], a[q][k]
					a[p][k] = c*apk - s*aqk
					a[q][k] = s*apk + c*aqk
				}
			}
		}
	}
	ev := make([]float64, n)
	for i := range ev {
		ev[i] = a[i][i]
	}
	sort.Float64s(ev)
	return ev
}

// ---------------------------------------------------------------------------
// The bench
// ---------------------------------------------------------------------------

func media(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x
	}
	return s / float64(len(v))
}

func desplegarPropio(niv []float64) []float64 {
	if len(niv) < 3 {
		return nil
	}
	d := float64(len(niv)-1) / (niv[len(niv)-1] - niv[0])
	sp := make([]float64, 0, len(niv))
	for i := 0; i+1 < len(niv); i++ {
		sp = append(sp, (niv[i+1]-niv[i])*d)
	}
	return sp
}

func sigma2(sp []float64, L float64, pasos int) float64 {
	x := make([]float64, len(sp)+1)
	for i, s := range sp {
		x[i+1] = x[i] + s
	}
	total := x[len(x)-1]
	if total <= L*1.2 {
		return math.NaN()
	}
	var ns []float64
	paso := (total - L) / float64(pasos)
	for k := 0; k < pasos; k++ {
		a := float64(k) * paso
		lo := sort.SearchFloat64s(x, a)
		hi := sort.SearchFloat64s(x, a+L)
		ns = append(ns, float64(hi-lo))
	}
	m := media(ns)
	v := 0.0
	for _, y := range ns {
		v += (y - m) * (y - m)
	}
	return v / float64(len(ns))
}

func repulsion(sp []float64) (frac, minimo float64) {
	minimo = math.Inf(1)
	n := 0
	for _, s := range sp {
		if s < minimo {
			minimo = s
		}
		if s < 0.1 {
			n++
		}
	}
	return float64(n) / float64(len(sp)), minimo
}
