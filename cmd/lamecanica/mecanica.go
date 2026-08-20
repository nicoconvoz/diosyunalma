package main

// mecanica.go - EL TRINQUETE TEMPLADO.
//
// ONE real number per site (exactly Phase VIII's state count), ONE continuous
// dial b, ONE bit eps, TWO derived constants (s_c, eta). Against the eleven
// independent knobs the auditor's section 6 forbade: one dial.
//
// The eleven words collapse into four objects:
//   {hardness, resistance, toughness, brittleness} = one scalar function U read
//       as curvature, maximum, area, and discontinuity at the maximum
//   {elasticity, plasticity}                       = which basin of that same U
//   {tension, compression}                         = the sign of one state variable
//   {fatigue}                                      = a property of the COUPLED field,
//       deliberately NOT constructed, so it is a prediction that may return null
//   {shear, torsion}                               = a bond connection and its curl,
//       declared out of scope with proofs (see gauge.go)

import "math"

// material is the medium plus its one state variable.
type material struct {
	n, kmax int
	h0      []float64 // the Phase VIII static field, inherited
	x       []float64 // THE state: one real per site. x(0) = 0 in every arm.
	amp2    []float64 // amp_i^2
	bf      []float64 // bf[d] = f(d)^2, so B_ij = bf[|i-j|]*amp2_i*amp2_j
	sig0    []float64 // sigma_i(0), frozen: makes the virgin medium an exact fixed point
	sc      float64   // YIELD STRESS, derived: stdev of the load the medium already applies
	b       float64   // the one dial: plastic step in log-h is 2b
	kappa   float64   // U''(0) = sc*pi/b
	eta     float64   // DERIVED step size
	eps     float64   // +1 shedding (closes under load) or -1 softening (opens)
	fuera   bool      // guard tripped: FUERA DE POZO, a report and never a clamp
	pasos   int
	// scratch buffers: the map runs hundreds of thousands of steps, and the
	// dynamics must not spend its time allocating.
	bhb, bq, bg []float64
}

// nuevoMaterial freezes every derived constant at t = 0. No spectrum is read here
// or anywhere inside the evolution: grep for espectro in this file returns nothing.
func nuevoMaterial(ms []modo, h0 []float64, c func(int) float64, kmax int, b, eps float64) *material {
	n := len(ms)
	m := &material{n: n, kmax: kmax, h0: normalizar1(h0), b: b, eps: eps}
	m.x = make([]float64, n)
	m.bhb = make([]float64, n)
	m.bq = make([]float64, n)
	m.bg = make([]float64, n)
	m.amp2 = make([]float64, n)
	for i := range ms {
		m.amp2[i] = ms[i].amp * ms[i].amp
	}
	m.bf = make([]float64, kmax+1)
	for d := 1; d <= kmax; d++ {
		v := c(d)
		m.bf[d] = v * v // f^2: the material law is SIGN-BLIND, so the sign channel
	} //               touches the operator only and the two channels separate cleanly
	m.sig0 = m.tension0()
	m.sc = desvio(m.sig0)
	m.kappa = m.sc * math.Pi / b
	mu := m.radioJacobiano()
	hb := m.hbar()
	m.eta = 0.2 / (m.kappa + maxi(hb)*mu)
	return m
}

// hbar is the mean-1 COPY of the permeability. It is used to build the operator and
// the load; the material state x is never renormalised, so no normalisation drift
// can manufacture memory (control 12).
func (m *material) hbar() []float64 {
	h := m.bhb
	s := 0.0
	for i := range h {
		h[i] = m.h0[i] * math.Exp(m.x[i]) // positivity with NO clamp
		s += h[i]
	}
	f := float64(m.n) / s
	for i := range h {
		h[i] *= f
	}
	return h
}

// carga is the LOCAL LOAD: q_i = hbar_i * sum_j B_ij hbar_j, in O(n*kmax).
// No eigenvalue is involved. Returns q and its total.
func (m *material) carga(hb []float64) ([]float64, float64) {
	g := m.bg
	for i := range g {
		g[i] = m.amp2[i] * hb[i]
	}
	q := m.bq
	Q := 0.0
	for i := 0; i < m.n; i++ {
		s := 0.0
		for d := 1; d <= m.kmax; d++ {
			if i-d >= 0 {
				s += m.bf[d] * g[i-d]
			}
			if i+d < m.n {
				s += m.bf[d] * g[i+d]
			}
		}
		q[i] = hb[i] * m.amp2[i] * s
		Q += q[i]
	}
	return q, Q
}

// tension0 is sigma_i = N*q_i/Q - 1. Two identities hold by construction and are
// asserted in the run: sum_i sigma_i = 0 exactly, and sigma is invariant under
// h -> c*h, so the F = 30 convention CANNOT steer the material.
func (m *material) tension0() []float64 {
	hb := m.hbar()
	q, Q := m.carga(hb)
	s := make([]float64, m.n) // copied out: q lives in a shared buffer
	for i := range s {
		s[i] = float64(m.n)*q[i]/Q - 1
	}
	return s
}

// radioJacobiano is the spectral radius of the closed-form stress Jacobian
//
//	M_ik = N*[ (delta_ik q_i + B_ik hbar_i hbar_k)/Q - 2 q_i q_k/Q^2 ]
//
// by power iteration, matrix-free. Used ONLY to fix the step size eta.
func (m *material) radioJacobiano() float64 {
	hb := append([]float64(nil), m.hbar()...)
	qs, Q := m.carga(hb)
	q := append([]float64(nil), qs...)
	v := make([]float64, m.n)
	for i := range v {
		v[i] = 1
	}
	lam := 0.0
	for it := 0; it < 20; it++ {
		g := make([]float64, m.n)
		for k := range g {
			g[k] = m.amp2[k] * hb[k] * v[k]
		}
		qv := 0.0
		for k := range q {
			qv += q[k] * v[k]
		}
		w := make([]float64, m.n)
		for i := 0; i < m.n; i++ {
			s := 0.0
			for d := 1; d <= m.kmax; d++ {
				if i-d >= 0 {
					s += m.bf[d] * g[i-d]
				}
				if i+d < m.n {
					s += m.bf[d] * g[i+d]
				}
			}
			w[i] = float64(m.n) * ((q[i]*v[i]+hb[i]*m.amp2[i]*s)/Q - 2*q[i]*qv/(Q*Q))
		}
		nr := 0.0
		for _, y := range w {
			nr += y * y
		}
		nr = math.Sqrt(nr)
		if nr == 0 {
			return 0
		}
		lam = nr
		for i := range v {
			v[i] = w[i] / nr
		}
	}
	return lam
}

// paso is THE MAP - one line, the whole mechanism:
//
//	x_i(t+1) = x_i(t) + eta*[ hbar_i*(P_i - eps*tau_i) - sc*sin(pi*x_i/b) ]
//
// At P = 0 and x = 0 the increment is identically zero (tau(0) = 0, sin 0 = 0), so
// the virgin medium is an EXACT fixed point and no algorithmic drift is possible.
func (m *material) paso(P []float64) {
	hb := m.hbar()
	q, Q := m.carga(hb)
	lim := 5 * m.b
	for i := 0; i < m.n; i++ {
		tau := (float64(m.n)*q[i]/Q - 1) - m.sig0[i]
		p := 0.0
		if P != nil {
			p = P[i]
		}
		m.x[i] += m.eta * (hb[i]*(p-m.eps*tau) - m.sc*math.Sin(math.Pi*m.x[i]/m.b))
		if math.Abs(m.x[i]) > lim {
			m.fuera = true // a REPORT, never a clamp
		}
	}
	m.pasos++
}

func (m *material) correr(P []float64, pasos int) {
	for t := 0; t < pasos; t++ {
		m.paso(P)
	}
}

// --- derived readings, no new parameters ------------------------------------

// residuo is the total permanent deformation.
func (m *material) residuo() float64 {
	s := 0.0
	for _, v := range m.x {
		s += math.Abs(v)
	}
	return s
}

// cedidos counts sites that left their original well.
func (m *material) cedidos() int {
	k := 0
	for _, v := range m.x {
		if math.Abs(v) > m.b/2 {
			k++
		}
	}
	return k
}

// maxAbs is the guard reading.
func (m *material) maxAbs() float64 {
	s := 0.0
	for _, v := range m.x {
		if a := math.Abs(v); a > s {
			s = a
		}
	}
	return s
}

// cuantizacion is the SHARPEST falsifiable signature of this mechanism: every
// yielded site must sit within 0.05 of an integer number of wells. A smooth
// residual would mean the effect is NOT this mechanism.
func (m *material) cuantizacion() (peor float64, cuantos int) {
	for _, v := range m.x {
		if math.Abs(v) > m.b/2 {
			cuantos++
			r := v / (2 * m.b)
			d := math.Abs(r - math.Round(r))
			if d > peor {
				peor = d
			}
		}
	}
	return
}

// contiguos reports whether the yielded set is a crack (contiguous) or a set of
// independent failures (scattered). Calling scattered failures a crack is the fake.
func (m *material) contiguos() (bloquesN, mayor int) {
	run := 0
	for i := 0; i < m.n; i++ {
		if math.Abs(m.x[i]) > m.b/2 {
			run++
			if run == 1 {
				bloquesN++
			}
			if run > mayor {
				mayor = run
			}
		} else {
			run = 0
		}
	}
	return
}

// --- perturbation profiles: the EXPERIMENT, not the model -------------------

func perfilPunto(n int, i0 int, a float64) []float64 {
	p := make([]float64, n)
	p[i0] = a
	return p
}

func perfilBulto(n, i0 int, a float64) []float64 {
	p := make([]float64, n)
	for i := range p {
		d := float64(i-i0) / 10
		p[i] = a * math.Exp(-d*d)
	}
	return p
}

func perfilUniforme(n int, a float64) []float64 {
	p := make([]float64, n)
	for i := range p {
		p[i] = a
	}
	return p
}

func maxi(v []float64) float64 {
	m := v[0]
	for _, x := range v {
		if x > m {
			m = x
		}
	}
	return m
}
