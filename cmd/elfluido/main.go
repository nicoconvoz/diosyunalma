// Command elfluido tests the captain's Phase IV intuition (Auditoria/37), in
// his own words:
//
//	"Es como una onda en un fluido: algunas grandes, otras pequeñas, otras
//	 medianas, todas resonando y creando una melodía única."
//
// The idea is his; the auditor formalised it as a research route. The change of
// language is the point: Phase III closed every family whose closed orbits can
// be CONCATENATED, and in a medium there are no orbits to concatenate - there
// are MODES that couple. The question stops being "how do we forbid gluing two
// orbits" and becomes "can gluing simply not be the fundamental operation".
//
// THE MINIMAL MEDIUM. Each prime power p^m is an excitation whose natural scale
// is log p: it contributes a ladder of frequencies 2*pi*k/log p. Uncoupled,
// those ladders are Phase III's disjoint loops, which keep the arithmetic and
// lose the correlations. Here they are all coupled to ONE common field mode -
// the fluid - with the strength the explicit formula itself assigns to that
// excitation, Lambda(p^m)/sqrt(p^m) = log p * p^(-m/2). Nothing is fitted.
//
// That coupling is rank one, so the coupled spectrum is EXACT: the levels are
// the roots of the secular equation
//
//	1 = g * sum_i v_i^2 / (E - w_i)
//
// one root strictly between each pair of consecutive uncoupled levels. No
// eigensolver, no truncation error, and the whole family in g is available.
//
// R6: the only arithmetic input is Lambda(n). The gamma_n are read only as a
// ruler, after the spectrum exists.
//
// Reproduce: go run ./cmd/elfluido
package main

import (
	"math"
	"sort"
)

// ---------------------------------------------------------------------------
// The excitations: one per prime power, with the weight the explicit formula
// gives it. This is the whole arithmetic input.
// ---------------------------------------------------------------------------

type excitacion struct {
	p, m int
	esc  float64 // log p: the characteristic scale
	amp  float64 // log p * p^(-m/2): the weight Weil assigns to this excitation
}

func excitaciones(topeP, topeM int) []excitacion {
	criba := make([]bool, topeP+1)
	var out []excitacion
	for p := 2; p <= topeP; p++ {
		if criba[p] {
			continue
		}
		for m := 2 * p; m <= topeP; m += p {
			criba[m] = true
		}
		lp := math.Log(float64(p))
		for m := 1; m <= topeM; m++ {
			out = append(out, excitacion{p, m, lp, lp * math.Pow(float64(p), -float64(m)/2)})
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// The medium: uncoupled ladders, and the exact rank-one coupled spectrum
// ---------------------------------------------------------------------------

type modo struct {
	w   float64 // frequency
	v   float64 // coupling amplitude to the common field
	p   int
	esc float64
}

// ladrillos builds the uncoupled modes in [t0,t1]: excitation (p,m) puts a
// ladder of frequencies 2*pi*k/log p, and carries its Weil weight.
func ladrillos(ex []excitacion, t0, t1 float64) []modo {
	var ms []modo
	for _, e := range ex {
		if e.m > 1 {
			continue // the ladder of p already contains its harmonics
		}
		for k := 1; ; k++ {
			w := 2 * math.Pi * float64(k) / e.esc
			if w > t1 {
				break
			}
			if w >= t0 {
				ms = append(ms, modo{w, e.amp, e.p, e.esc})
			}
		}
	}
	sort.Slice(ms, func(a, b int) bool { return ms[a].w < ms[b].w })
	return ms
}

// acoplar solves the secular equation of the rank-one coupling exactly. For
// g > 0 there is exactly one root in each open interval (w_i, w_{i+1}), and the
// last root sits above the top level.
func acoplar(ms []modo, g float64) []float64 {
	if g == 0 {
		out := make([]float64, len(ms))
		for i, m := range ms {
			out[i] = m.w
		}
		return out
	}
	f := func(E float64) float64 {
		s := 0.0
		for _, m := range ms {
			s += m.v * m.v / (E - m.w)
		}
		return g*s - 1
	}
	out := make([]float64, 0, len(ms))
	for i := 0; i+1 < len(ms); i++ {
		lo, hi := ms[i].w, ms[i+1].w
		// f goes from +inf just above w_i to -inf just below w_{i+1}
		a := lo + (hi-lo)*1e-12
		b := hi - (hi-lo)*1e-12
		fa := f(a)
		fb := f(b)
		if fa*fb > 0 {
			out = append(out, (lo+hi)/2) // degenerate interval: keep the level
			continue
		}
		for it := 0; it < 90; it++ {
			mm := (a + b) / 2
			if f(mm)*fa > 0 {
				a = mm
			} else {
				b = mm
			}
		}
		out = append(out, (a+b)/2)
	}
	return out
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

func desvio(v []float64) float64 {
	m := media(v)
	s := 0.0
	for _, x := range v {
		s += (x - m) * (x - m)
	}
	return math.Sqrt(s / float64(len(v)-1))
}

// desplegarPropio unfolds a spectrum by its OWN mean density - the honest way
// when the candidate does not claim the Riemann counting law.
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

// repulsion counts how much probability sits at very small spacings: a medium
// with genuine collective interaction empties the origin, a superposition of
// independent clocks does not.
func repulsion(sp []float64) (frac float64, minimo float64) {
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

func eco(niv []float64, T float64) float64 {
	M := len(niv)
	s := 0.0
	for i, g := range niv {
		w := 0.5 - 0.5*math.Cos(2*math.Pi*(float64(i)+0.5)/float64(M))
		s += w * math.Cos(g*T)
	}
	return 2 * s / float64(M)
}

func murcielago(niv, per []float64, tope float64) float64 {
	const ancho = 0.012
	var dentro, fuera float64
	nIn, nOut := 0, 0
	for T := 0.35; T <= tope; T += 0.0009 {
		cerca := false
		for _, p := range per {
			if math.Abs(T-p) < ancho {
				cerca = true
				break
			}
		}
		v := math.Abs(eco(niv, T))
		if cerca {
			dentro += v
			nIn++
		} else {
			fuera += v
			nOut++
		}
	}
	return (dentro / float64(nIn)) / (fuera / float64(nOut))
}

// ---------------------------------------------------------------------------
// The ruler: the true zeros, read only after the candidate exists
// ---------------------------------------------------------------------------

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t*t)
}

func suave(T float64) float64 { return theta(T)/math.Pi + 1 }

func zetaZ(t float64) float64 {
	th := theta(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		s += math.Cos(th-t*math.Log(fn)) / math.Sqrt(fn)
	}
	s *= 2
	p := u - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sign := 1.0
	if N%2 == 0 {
		sign = -1
	}
	return s + sign*math.Pow(2*math.Pi/t, 0.25)*c0
}

func cerosVerdaderos(t0, t1, paso float64) []float64 {
	var g []float64
	a, za := t0, zetaZ(t0)
	for b := t0 + paso; b <= t1; b += paso {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 60; i++ {
				m := (lo + hi) / 2
				zm := zetaZ(m)
				if (zlo < 0) != (zm < 0) {
					hi = m
				} else {
					lo, zlo = m, zm
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}
