package main

// fase19.go - THE RIGIDITY STEP ("el escalon con dientes"). Phase XVIII proved
// the coupling matters but its GUE material was an iid Wigner-gap random walk:
// diffusive, no long-range rigidity, so it dilutes the field's coherence and
// broadens spacings (0.58 vs the real 0.42). Her blessing: replace the material
// with GENUINE long-range-rigid GUE and repeat the identical judgment.
//
// THE MATERIAL, zero knobs: the Hermite beta=2 tridiagonal ensemble
// (Dumitriu-Edelman) whose eigenvalues are EXACTLY GUE:
//     diagonal  d_i ~ N(0, 1)
//     off-diag  c_i ~ chi(2*(N-i)) / sqrt(2)
// Eigenvalues by tridiagonal QL with implicit shifts; unfolded with the
// analytic semicircle CDF; only the BULK (|t| < 0.75) is used, away from the
// edges. The unfolded u_k have unit mean spacing, Wigner spacings AND the log
// number variance - the rigidity the walk lacked. Then the SAME pinning as
// Phase XVIII: N_smooth(g) + S(g) = u_k.
//
// PREDICTION, pre-registered: rigid material preserves the field's coherence -
// amplitude rises from 0.138 toward the real 0.349, spacings stay near 0.42,
// T4 grows past 0.116, the crossing stays near 0.86-0.96, and the phase
// destruction still kills it. If the rigid arm does NOT beat the walk arm, the
// rigidity hypothesis fails and is declared dead.

import (
	"fmt"
	"math"
	"sort"
)

func (d *dado) normal() float64 {
	u1 := d.u()
	if u1 < 1e-300 {
		u1 = 1e-300
	}
	return math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*d.u())
}

func (d *dado) chi(k int) float64 {
	if k <= 0 {
		return 0
	}
	if k <= 100 {
		s := 0.0
		for i := 0; i < k; i++ {
			x := d.normal()
			s += x * x
		}
		return math.Sqrt(s)
	}
	// large-k: chi_k ~ sqrt(k - 1/2) + N(0, 1/2)
	return math.Sqrt(float64(k)-0.5) + d.normal()/math.Sqrt2
}

// tql: eigenvalues of a symmetric tridiagonal matrix (diagonal a, off-diag b),
// QL with implicit shifts. Standard algorithm, values only.
func tql(a, b []float64) []float64 {
	n := len(a)
	d2 := append([]float64(nil), a...)
	e := make([]float64, n)
	copy(e, b)
	for l := 0; l < n; l++ {
		iter := 0
		for {
			m := l
			for ; m < n-1; m++ {
				dd := math.Abs(d2[m]) + math.Abs(d2[m+1])
				if math.Abs(e[m]) <= 1e-14*dd {
					break
				}
			}
			if m == l {
				break
			}
			iter++
			if iter > 50 {
				break
			}
			g := (d2[l+1] - d2[l]) / (2 * e[l])
			r := math.Hypot(g, 1)
			sg := r
			if g < 0 {
				sg = -r
			}
			g = d2[m] - d2[l] + e[l]/(g+sg)
			s, c := 1.0, 1.0
			p := 0.0
			for i := m - 1; i >= l; i-- {
				f := s * e[i]
				bb := c * e[i]
				r = math.Hypot(f, g)
				e[i+1] = r
				if r == 0 {
					d2[i+1] -= p
					e[m] = 0
					break
				}
				s = f / r
				c = g / r
				g = d2[i+1] - p
				r = (d2[i]-g)*s + 2*c*bb
				p = s * r
				d2[i+1] = g + p
				g = c*r - bb
			}
			if r == 0 && m-1 >= l {
				continue
			}
			d2[l] -= p
			e[l] = g
			e[m] = 0
		}
	}
	sort.Float64s(d2)
	return d2
}

// materialRigido: unfolded bulk GUE eigenvalues with unit mean spacing.
func materialRigido(cuantos int, d *dado) []float64 {
	N := cuantos*2 + 800 // bulk fraction ~ enough
	a := make([]float64, N)
	b := make([]float64, N-1)
	for i := 0; i < N; i++ {
		a[i] = d.normal()
	}
	for i := 0; i < N-1; i++ {
		b[i] = d.chi(2*(N-1-i)) / math.Sqrt2
	}
	ev := tql(a, b)
	// semicircle unfolding: t = lambda/(2 sqrt(N)); F(t) in [0,1]
	var u []float64
	for _, l := range ev {
		t := l / (2 * math.Sqrt(float64(N)))
		if t < -0.75 || t > 0.75 {
			continue
		}
		F := 0.5 + (t*math.Sqrt(1-t*t)+math.Asin(t))/math.Pi
		u = append(u, float64(N)*F)
	}
	sort.Float64s(u)
	if len(u) > cuantos {
		u = u[:cuantos]
	}
	return u
}

// clavarSerie pins an unfolded series to the full counting equation.
func clavarSerie(u []float64, fases []float64, campo bool) []float64 {
	off := theta(30)/math.Pi + 1 - u[0]
	g := 30.0
	out := make([]float64, 0, len(u))
	for _, uk := range u {
		g = clavar(uk+off, g+1/(math.Log(g/(2*math.Pi))/(2*math.Pi)), fases, campo)
		out = append(out, g)
	}
	return out
}

// varConteo: number variance of the unfolded series in windows of length L.
func varConteo(u []float64, L float64) float64 {
	var cs []float64
	for x := u[0]; x+L < u[len(u)-1]; x += L {
		n := 0
		for _, v := range u { // fine at this size
			if v >= x && v < x+L {
				n++
			}
		}
		cs = append(cs, float64(n))
	}
	m := media(cs)
	s := 0.0
	for _, c := range cs {
		s += (c - m) * (c - m)
	}
	return s / float64(len(cs)-1)
}

func fase19() {
	fmt.Println("🧵🦷 FASE XIX — LA RIGIDEZ: el escalón con dientes")
	fmt.Println("   material nuevo, cero perillas: autovalores GUE GENUINOS (ensamble tridiagonal")
	fmt.Println("   de Hermite β=2), desplegados por la semicircular, clavados por la MISMA")
	fmt.Println("   ecuación de conteo de la Fase XVIII. Predicción pre-registrada: la rigidez")
	fmt.Println("   preserva la coherencia — amplitud sube de 0,138 hacia 0,349, espaciados cerca")
	fmt.Println("   de 0,42, T4 crece, y la fase destruida sigue matando. Si el rígido NO supera")
	fmt.Println("   al paseo, la hipótesis de la rigidez muere y se declara.")

	var Tp []float64
	for _, p := range []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} {
		Tp = append(Tp, math.Log(float64(p)))
	}
	g := cerosPaso(4000, 0.02)
	d := &dado{s: 20260826}
	Ereal, crReal, pendReal := pipeline(g, Tp, d, 120)
	T4real := descomponer(g, Tp)

	// the material check itself, before any field: rigidity measured
	uR := materialRigido(3500, d)
	uW := func() []float64 { // the walk, for the same check
		u := make([]float64, 3500)
		x := 0.0
		for i := range u {
			x += wigner(d)
			u[i] = x
		}
		return u
	}()
	fmt.Printf("\n§0 · EL MATERIAL, verificado antes del campo (varianza de conteo en L=10):\n")
	fmt.Printf("   paseo iid (Fase XVIII): %.2f · GUE rígido: %.2f · GUE teórico: ~0,59 · Poisson: 10\n",
		varConteo(uW, 10), varConteo(uR, 10))

	type brazo19 struct {
		nom           string
		E, T4         []float64
		cruces, pends []float64
		amp           float64
		sDes          float64
	}
	estad := func(sp []float64) float64 {
		var ss []float64
		for _, p := range paresDe(sp) {
			ss = append(ss, sDe(p, true))
		}
		return desvio(ss)
	}
	corre := func(nom string, gen func(*dado) []float64, sem int) brazo19 {
		br := brazo19{nom: nom}
		for r := 0; r < sem; r++ {
			sp := gen(d)
			E, c, p := pipeline(sp, Tp, d, 60)
			if r == 0 {
				br.E = E
				br.T4 = descomponer(sp, Tp)
				br.sDes = estad(sp)
			} else {
				for b := range br.E {
					br.E[b] += E[b]
				}
			}
			br.cruces = append(br.cruces, c)
			br.pends = append(br.pends, p)
		}
		for b := range br.E {
			br.E[b] /= float64(sem)
		}
		br.amp = resumen(br.E)
		return br
	}

	fmt.Println("\n   corriendo los brazos…")
	bR0 := corre("R0) rígido, SIN campo", func(dd *dado) []float64 {
		return clavarSerie(materialRigido(3500, dd), nil, false)
	}, 3)
	bRC := corre("RC) RÍGIDO + CAMPO coherente", func(dd *dado) []float64 {
		return clavarSerie(materialRigido(3500, dd), nil, true)
	}, 3)
	bRF := corre("RF) rígido + campo, fase rota", func(dd *dado) []float64 {
		f := make([]float64, len(nLam))
		for i := range f {
			f[i] = 2 * math.Pi * dd.u()
		}
		return clavarSerie(materialRigido(3500, dd), f, true)
	}, 3)
	bAC := corre("AC) paseo + campo (Fase XVIII)", func(dd *dado) []float64 {
		return autoconsistente(4000, dd, true, true, nil)
	}, 3)

	fmt.Println("\n§1 · LA TABLA")
	fmt.Printf("   %-34s %9s %16s %18s %8s\n", "modelo", "amplitud", "cruce s*", "pendiente", "s desv")
	fmt.Printf("   %-34s %9.3f %16.3f %18.3f %8s\n", "REAL (3474 ceros)", resumen(Ereal), crReal, pendReal, "0.42")
	for _, br := range []brazo19{bR0, bRC, bRF, bAC} {
		fmt.Printf("   %-34s %9.3f %16s %18s %8.2f\n", br.nom, br.amp, listaF(br.cruces), listaF(br.pends), br.sDes)
	}

	fmt.Println("\n§2 · T4, la prueba decisiva")
	fmt.Printf("   %-11s %9s %10s %10s\n", "bin s", "T4 real", "T4 rígido", "T4 paseo")
	nb := int((zHi - zLo) / zW)
	for b := 0; b < nb; b++ {
		fmt.Printf("   %.2f–%.2f %+9.4f %+10.4f %+10.4f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW, T4real[b], bRC.T4[b], bAC.T4[b])
	}

	fmt.Println("\n§3 · E(s) RÍGIDO+CAMPO contra REAL, y el residuo")
	rmax := 0.0
	for b := 0; b < nb; b++ {
		r := Ereal[b] - bRC.E[b]
		if math.Abs(r) > rmax {
			rmax = math.Abs(r)
		}
		fmt.Printf("   %.2f–%.2f  real %+8.4f  rígido %+8.4f  residuo %+8.4f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW, Ereal[b], bRC.E[b], r)
	}
	fmt.Printf("   residuo máximo: %.4f (XVIII acoplado: 0,250 · XVI densidad: 0,312)\n", rmax)

	fmt.Println("\n§4 · EL VEREDICTO: rígido contra paseo — ¿la rigidez era el ingrediente?")
	fmt.Printf("   amplitud: rígido %.3f · paseo %.3f · real %.3f\n", bRC.amp, bAC.amp, resumen(Ereal))
	fmt.Printf("   espaciados: rígido %.2f · paseo %.2f · real 0,42\n", bRC.sDes, bAC.sDes)
	fmt.Printf("   fase rota: %.3f — ¿sigue muriendo?\n", bRF.amp)

	dibujar19(Ereal, bRC.E, bAC.E, bRF.E, crReal, bRC.cruces)
}
