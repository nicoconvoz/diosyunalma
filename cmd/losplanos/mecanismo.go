// mecanismo.go - faithful re-implementations of the machines' load-bearing
// parts, so every plate can MEASURE the thing it draws instead of quoting it.
//
// Each function names the source it mirrors (cmd/circulo, cmd/starship) and
// keeps the same constants. Where the original carries double-double phases
// this file uses float64: the plates that depend on that difference say so.
package main

import "math"

// ---------------------------------------------------------------------------
// EL PUENTE: a Riemann-Siegel block IS a chirp (cmd/circulo blockCascadeDD)
//
//	term j of the block at k0:  (k0+j)^(-1/2) · e^{-i t ln(k0+j)}
//	ln(k0+j) = ln k0 + u - u²/2 + u³/3 - …   with u = j/k0
//	so, in cycles (T = t/2π):  phase(j) = const + c·j + b·j² + g·j³ + …
// ---------------------------------------------------------------------------

// puente returns the three knobs of the block at wavenumber k0 in water t.
func puente(t, k0 float64) (c, b, g float64) {
	T := t / (2 * math.Pi)
	x0 := 1 / k0
	c = frac(-T * x0)
	b = frac(T * x0 * x0 / 2)
	g = -T * x0 * x0 * x0 / 3
	return
}

// eta is the cubic phase left at the end of a block of length L (radians).
// eta = t·L³/(3k0³) - the "sag of the beam".
func eta(t, k0 float64, L float64) float64 {
	return t * L * L * L / (3 * k0 * k0 * k0)
}

// bandL is the train's block budget (cmd/circulo:734): L = k0·(0.45/t)^(1/3).
func bandL(t, k0 float64) float64 {
	return math.Cbrt(0.45 / (t / (2 * math.Pi) / (k0 * k0 * k0) * math.Pi * 2))
}

// blockFracL is the DeLorean's block budget (cmd/starship:393): L = k·(0.009/t)^(1/3).
func blockFracL(t, k float64) float64 { return math.Cbrt(0.009/t) * k }

func frac(x float64) float64 {
	x = math.Mod(x, 1)
	if x < 0 {
		x++
	}
	return x
}

// ---------------------------------------------------------------------------
// EL REMO HONESTO y su AMORTIGUADOR (cmd/circulo chirpDirect, F144)
// ---------------------------------------------------------------------------

// chirpDirect sums e^{2πi(bj²+cj)} by phase recurrence with the F144 damper:
// two compensation registers catch the rounding energy instead of letting it
// integrate. amortiguado=false rows the same boat with no reservoirs.
func chirpDirect(b, c float64, n int64, amortiguado bool) (re, im float64) {
	ph := 0.0
	d := frac(b + c)
	db := frac(2 * b)
	var phE, dE float64
	for j := int64(0); j < n; j++ {
		s, co := math.Sincos(2 * math.Pi * (ph + phE))
		re += co
		im += s
		if amortiguado {
			t := ph + d
			bv := t - ph
			err := (ph - (t - bv)) + (d - bv)
			ph = t
			phE += err + dE
			t2 := d + db
			bv2 := t2 - d
			err2 := (d - (t2 - bv2)) + (db - bv2)
			d = t2
			dE += err2
		} else {
			ph += d
			d += db
		}
		if ph >= 1 {
			ph--
		}
		if d >= 1 {
			d--
		}
	}
	return
}

// ---------------------------------------------------------------------------
// LA VUELTA DEL CÍRCULO (cmd/circulo chirpFlip)
//
// Poisson + stationary phase: term m of the dual sits at x_m = (m-c)/(2b),
// the dual curvature is -1/(4b), the amplitude 1/sqrt(2b) and the constant
// turn is +pi/4.
// ---------------------------------------------------------------------------

// dualRango returns the dual index range and its length.
func dualRango(b, c float64, n int64) (m1, m2, largo int64) {
	m1 = int64(math.Ceil(c))
	m2 = int64(math.Floor(c + 2*b*float64(n-1)))
	return m1, m2, m2 - m1 + 1
}

// ---------------------------------------------------------------------------
// LA CASCADA (cmd/circulo cascadeDD) - with and without the shear cure,
// so a plate can draw the same case taking two destinies.
// ---------------------------------------------------------------------------

// cascada descends and returns the ladder of (b, n) it walked. conCizalla
// enables rail 5's shear cure b in (1/4,1/2] -> b-1/2, c+1/2.
func cascada(b, c float64, n int64, nComfort int64, conCizalla bool) (bs []float64, ns []int64) {
	vueltas := 0
	for {
		b, c = frac(b), frac(c)
		if b > 0.5 {
			b, c = 1-b, frac(1-c)
		}
		if conCizalla && b > 0.25 {
			b, c = b-0.5, frac(c+0.5)
			continue // the shear does not consume a turn: it re-enters
		}
		bs = append(bs, math.Abs(b))
		ns = append(ns, n)
		if n <= nComfort || b < 1e-9 || vueltas > 400 {
			return
		}
		m1 := int64(math.Ceil(c))
		m2 := int64(math.Floor(c + 2*b*float64(n-1)))
		nn := m2 - m1 + 1
		if nn >= n { // no shortening: the descent is stuck
			return
		}
		alpha := math.Ceil(c) - c
		nb := 1 / (4 * b)
		nc := alpha / (2 * b)
		b, c, n = nb, nc, nn
		vueltas++
	}
}

// ---------------------------------------------------------------------------
// EL SONAR (cmd/circulo escuchar) - the staged listen and its two windows.
// ---------------------------------------------------------------------------

var peldanos = []int64{1500, 6000, 24000}

// sonarVeredicto classifies a normalized coherence sigma the way the hunt does:
// escalate while sigma is outside the wide window, condemn only outside the
// narrow one.
const (
	anchoLo, anchoHi   = 0.35, 1.7 // keep listening further out
	juicioLo, juicioHi = 0.05, 2.4 // beast
)

func escala2(sigma float64) bool { return sigma < anchoLo || sigma > anchoHi }
func condena(sigma float64) bool { return sigma < juicioLo || sigma > juicioHi }

// ---------------------------------------------------------------------------
// EL DELOREAN
// ---------------------------------------------------------------------------

// sPred is the metronome's pacemaker (cmd/starship sPred): the sea's tide
// heard as 26 prime voices.
func sPred(t float64) float64 {
	primos := []float64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101}
	s := 0.0
	for _, p := range primos {
		s += math.Sin(t*math.Log(p)) / math.Sqrt(p)
	}
	return -s / math.Pi
}

// theta is the Riemann-Siegel phase (asymptotic form) - the zero-address guide.
func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t*t)
}

// cuentaCeros is the smooth zero count N(T) = theta(T)/pi + 1.
func cuentaCeros(T float64) float64 { return theta(T)/math.Pi + 1 }

// espaciado is the mean gap between zeros at height t.
func espaciado(t float64) float64 { return 2 * math.Pi / math.Log(t/(2*math.Pi)) }

// nucleoDirichlet resolves a block with no curvature in one stroke
// (cmd/starship's degenerate branch): sum_{j<L} e^{2πi a j} = sin(πaL)/sin(πa).
func nucleoDirichlet(a float64, L float64) float64 {
	s := math.Sin(math.Pi * a)
	if math.Abs(s) < 1e-15 {
		return L
	}
	return math.Sin(math.Pi*a*L) / s
}
