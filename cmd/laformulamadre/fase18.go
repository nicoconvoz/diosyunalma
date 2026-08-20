package main

// fase18.go - REPULSION + FIELD, SOLVED TOGETHER. Her hypothesis: in the real
// zeros the GUE repulsion and the explicit-formula field determine the final
// positions JOINTLY, not in sequence.
//
// THE SELF-CONSISTENT CONSTRUCTION, zero new knobs (her sections 5-6):
// real zeros are pinned by the FULL counting function, N_smooth(g)+S(g) = k-1/2.
// So the model does exactly that to GUE material: draw unfolded GUE points
//     u_k = u_{k-1} + w_k,   w_k ~ Wigner(GUE surmise, mean 1)
// and pin each final position g_k as the solution of
//     N_smooth(g_k) + S(g_k) = u_k
// with S the SAME truncated fluctuation of Phases XVI-XVII (n <= 97). The field
// enters INSIDE the pinning: where S rises, positions shift AND local gaps
// compress simultaneously and consistently - sizes and positions emerge
// together, which is precisely her Phase XVIII hypothesis. No amplitude, no
// dial: the formula fixes everything.
// (Her section 5's other family, the Coulomb gas, needs a temperature and a
// relaxation schedule - free knobs - so by her section 6 the knobless
// implementation is run and the gas is deferred, declared, not hidden.)
//
// PREDICTION, pre-registered: if self-consistency is the missing piece, the
// coupled arm beats BOTH references (density 0.124, posterior displacement
// 0.082) in amplitude, grows T4 toward the real 0.267, moves the crossing
// toward 0.863, keeps Wigner-like spacings, and dies under phase destruction.
//
// ARMS (identical pipeline, 60 shuffles, 3 seeds):
//   R)  repulsion alone: pin with S = 0            [must recover GUE-pure]
//   F)  field alone: rigid picket u_k = k-1/2 + S  [no repulsion]
//   AC) COUPLED, coherent                          [the hypothesis]
//   AF) COUPLED, phases randomized                 [the kill control]
//   AP) a-posteriori displacement (Phase XVII arm) [the decisive comparison]
//   AD) density-only (Phase XVI arm)               [reference]

import (
	"fmt"
	"math"
)

// clavar solves N_smooth(g) + S(g) = u by Newton with a bisection guard.
func clavar(u, gIni float64, fases []float64, campo bool) float64 {
	g := gIni
	for it := 0; it < 40; it++ {
		s := 0.0
		if campo {
			s = sTrunc(g, fases)
		}
		f := theta(g)/math.Pi + 1 + s - u
		rho := math.Log(g/(2*math.Pi)) / (2 * math.Pi)
		if math.Abs(f) < 1e-10 {
			break
		}
		paso := f / rho
		if paso > 2 {
			paso = 2
		}
		if paso < -2 {
			paso = -2
		}
		g -= paso
		if g < 15 {
			g = 15
		}
	}
	return g
}

// autoconsistente builds the coupled spectrum: GUE unfolded points pinned by
// the full counting function. repulsion=false gives the rigid picket (field
// alone); campo=false gives repulsion alone.
func autoconsistente(tope float64, d *dado, repulsion, campo bool, fases []float64) []float64 {
	u := theta(30)/math.Pi + 1
	g := 30.0
	var out []float64
	for {
		w := 1.0
		if repulsion {
			w = wigner(d)
		}
		u += w
		g = clavar(u, g+w/(math.Log(g/(2*math.Pi))/(2*math.Pi)), fases, campo)
		if g > tope {
			break
		}
		out = append(out, g)
	}
	return out
}

func fase18() {
	fmt.Println("🧵⚛️ FASE XVIII — REPULSIÓN + CAMPO, RESUELTOS JUNTOS")
	fmt.Println("   construcción sin perillas: puntos GUE desplegados u_k clavados por la ecuación")
	fmt.Println("   de conteo COMPLETA N_liso(γ) + S(γ) = u_k — el campo entra ADENTRO del clavado:")
	fmt.Println("   tamaños y posiciones emergen juntos, que es exactamente la hipótesis.")
	fmt.Println("   predicción pre-registrada: el acoplado supera a densidad (0,124) y a posteriori")
	fmt.Println("   (0,082), T4 crece hacia 0,267, el cruce corre hacia 0,863, y muere sin fase.")

	var Tp []float64
	for _, p := range []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} {
		Tp = append(Tp, math.Log(float64(p)))
	}
	g := cerosPaso(4000, 0.02)
	d := &dado{s: 20260825}
	Ereal, crReal, pendReal := pipeline(g, Tp, d, 120)
	T4real := descomponer(g, Tp)

	type brazo18 struct {
		nom           string
		E, T4         []float64
		cruces, pends []float64
		amp           float64
		sMed, sDes    float64
	}
	estad := func(sp []float64) (m, sd float64) {
		var ss []float64
		for _, p := range paresDe(sp) {
			ss = append(ss, sDe(p, true))
		}
		return media(ss), desvio(ss)
	}
	corre := func(nom string, gen func(*dado) []float64, sem int) brazo18 {
		br := brazo18{nom: nom}
		for r := 0; r < sem; r++ {
			sp := gen(d)
			E, c, p := pipeline(sp, Tp, d, 60)
			if r == 0 {
				br.E = E
				br.T4 = descomponer(sp, Tp)
				br.sMed, br.sDes = estad(sp)
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
	bR := corre("R)  repulsión sola (S=0)", func(dd *dado) []float64 {
		return autoconsistente(4000, dd, true, false, nil)
	}, 3)
	bF := corre("F)  campo solo (sin repulsión)", func(dd *dado) []float64 {
		return autoconsistente(4000, dd, false, true, nil)
	}, 1)
	bAC := corre("AC) ACOPLADO coherente", func(dd *dado) []float64 {
		return autoconsistente(4000, dd, true, true, nil)
	}, 3)
	bAF := corre("AF) acoplado, fase DESTRUIDA", func(dd *dado) []float64 {
		f := make([]float64, len(nLam))
		for i := range f {
			f[i] = 2 * math.Pi * dd.u()
		}
		return autoconsistente(4000, dd, true, true, f)
	}, 3)
	bAP := corre("AP) desplazamiento A POSTERIORI", func(dd *dado) []float64 {
		return desplazar(surrogate(4000, nil, dd), 1, nil)
	}, 3)
	bAD := corre("AD) densidad sola (Fase XVI)", func(dd *dado) []float64 {
		return surrogate(4000, epsilonPrimos, dd)
	}, 3)

	fmt.Println("\n§1 · LA TABLA (su §9)")
	fmt.Printf("   %-34s %9s %16s %18s %8s %8s\n", "modelo", "amplitud", "cruce s*", "pendiente", "s med", "s desv")
	fmt.Printf("   %-34s %9.3f %16.3f %18.3f %8s %8s\n", "REAL (3474 ceros)", resumen(Ereal), crReal, pendReal, "1.00", "0.42")
	for _, br := range []brazo18{bR, bF, bAC, bAF, bAP, bAD} {
		fmt.Printf("   %-34s %9.3f %16s %18s %8.2f %8.2f\n", br.nom, br.amp, listaF(br.cruces), listaF(br.pends), br.sMed, br.sDes)
	}
	fmt.Println("   (s med/desv: estadística de espaciados del brazo, su §12 — Wigner da ~1,00/0,42)")

	fmt.Println("\n§2 · LA PRUEBA DECISIVA (su §8): ¿aparece el factor 3 de T4?")
	fmt.Printf("   %-11s %9s %10s %10s %10s\n", "bin s", "T4 real", "T4 acopl.", "T4 a-post.", "T4 densid.")
	nb := int((zHi - zLo) / zW)
	for b := 0; b < nb; b++ {
		fmt.Printf("   %.2f–%.2f %+9.4f %+10.4f %+10.4f %+10.4f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW, T4real[b], bAC.T4[b], bAP.T4[b], bAD.T4[b])
	}

	fmt.Println("\n§3 · E(s) ACOPLADO contra REAL, bin por bin, y el residuo")
	rmax := 0.0
	for b := 0; b < nb; b++ {
		r := Ereal[b] - bAC.E[b]
		if math.Abs(r) > rmax {
			rmax = math.Abs(r)
		}
		fmt.Printf("   %.2f–%.2f  real %+8.4f  acoplado %+8.4f  residuo %+8.4f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW, Ereal[b], bAC.E[b], r)
	}
	fmt.Printf("   residuo máximo del acoplado: %.4f (densidad: 0,312 · a posteriori: 0,300)\n", rmax)

	fmt.Println("\n§4 · SIMULTÁNEO contra A POSTERIORI (su §10, la comparación que decide)")
	fmt.Printf("   amplitud: acoplado %.3f · a posteriori %.3f · densidad %.3f\n", bAC.amp, bAP.amp, bAD.amp)
	fmt.Printf("   fase destruida: %.3f — ¿muere? (GUE puro anda en ~0,03)\n", bAF.amp)

	dibujar18(Ereal, bAC.E, bAP.E, bAF.E, crReal, bAC.cruces)
}
