package main

// fase17.go - THE DISPLACEMENT FIELD. Her order: define it without circularity,
// simulate it, compare against GUE-pure and GUE+density, and destroy its phase
// as the kill-control.
//
// THE RULE, derived - ZERO free parameters (her sections 4 and 8):
// the explicit formula's counting fluctuation is
//     S(g) = -(1/pi) * sum_{n<=97, Lambda>0} Lambda(n) sin(g log n) / (sqrt(n) log n)
// and a zero pinned by N_smooth(g) + S(g) = n - 1/2 sits displaced from its
// smooth position by
//     dg(g) = -S(g) / rhobar(g),      rhobar = log(g/2pi) / 2pi.
// The amplitude is fixed BY THE FORMULA; the truncation n <= 97 matches the
// echo's own periods (declared in Phase XVI). Nothing is read from the real
// curve: not the crossing, not the amplitude, not the slope, not the residual.
//
// PREDICTION, pre-registered before running: if the hypothesis is right, the
// displaced arm beats the density arm in amplitude, moves the crossing from
// 0.70 toward the real 0.863, steepens the slope beyond -0.6, produces the
// sin-sin structure (T4), and DIES when the component phases are randomized.
// If it needs calibration to work, it fails by her section 18.
//
// ARMS (identical Phase XV/XVI pipeline, 60 shuffles, 3 seeds each):
//   a) GUE pure                          [Phase XVI reference]
//   b) GUE + density modulation          [Phase XVI reference]
//   c) GUE + DISPLACEMENT (positions shifted by dg)
//   d) GUE + density + displacement      [declared: double-counts the gradient]
//   e) displacement with RANDOMIZED component phases phi_n (kill control:
//      same amplitude distribution, coherence destroyed)
//   f) robustness: amplitude x0.5 and x2 on arm c (robustness, never selection)

import (
	"fmt"
	"math"
)

var nLam = []int{2, 3, 4, 5, 7, 8, 9, 11, 13, 16, 17, 19, 23, 25, 27, 29, 31, 32, 37, 41, 43, 47, 49, 53, 59, 61, 64, 67, 71, 73, 79, 81, 83, 89, 97}

// sTrunc: the truncated counting fluctuation S(g); fases = optional random
// phase per component (the kill control).
func sTrunc(g float64, fases []float64) float64 {
	s := 0.0
	for i, n := range nLam {
		lp := math.Log(float64(basePrimo(n)))
		ln := math.Log(float64(n))
		ph := g * ln
		if fases != nil {
			ph += fases[i]
		}
		s += lp * math.Sin(ph) / (math.Sqrt(float64(n)) * ln)
	}
	return -s / math.Pi
}

// desplazar shifts every position by A * dg(g) = -A * S(g)/rhobar(g).
func desplazar(g []float64, A float64, fases []float64) []float64 {
	out := make([]float64, len(g))
	for i, x := range g {
		rho := math.Log(x/(2*math.Pi)) / (2 * math.Pi)
		out[i] = x - A*sTrunc(x, fases)/rho
	}
	return out
}

// descomponer: the Phase XVI four-term decomposition, returning the summed
// selection-sin term T4 per bin (the one the hypothesis must explain).
func descomponer(g []float64, Tp []float64) (T4 []float64) {
	ps := paresDe(g)
	nb := int((zHi - zLo) / zW)
	T4 = make([]float64, nb)
	Sall := make([]float64, len(Tp))
	for ti, T := range Tp {
		for _, p := range ps {
			Sall[ti] += math.Sin(p.base * T)
		}
		Sall[ti] /= float64(len(ps))
	}
	miembros := make([][]int, nb)
	for i, p := range ps {
		s := sDe(p, true)
		if s < zLo || s >= zHi {
			continue
		}
		miembros[int((s-zLo)/zW)] = append(miembros[int((s-zLo)/zW)], i)
	}
	for b := 0; b < nb; b++ {
		if len(miembros[b]) == 0 {
			continue
		}
		for ti, T := range Tp {
			var sC, sG float64
			for _, i := range miembros[b] {
				sC += math.Sin(ps[i].base * T)
				sG += math.Sin(ps[i].gap * T / 2)
			}
			n := float64(len(miembros[b]))
			T4[b] += 2 * ((sC/n - Sall[ti]) * (sG / n)) / float64(len(Tp))
		}
	}
	return
}

func resumen(E []float64) (amp float64) {
	for _, v := range E {
		if math.Abs(v) > amp {
			amp = math.Abs(v)
		}
	}
	return
}

func fase17() {
	fmt.Println("🧵🌊 FASE XVII — EL CAMPO DE DESPLAZAMIENTO: construir una causa que pueda fallar")
	fmt.Println("   regla derivada, CERO parámetros libres: δγ(γ) = −S(γ)/ρ̄(γ), S truncada en n ≤ 97")
	fmt.Println("   predicción pre-registrada: c) supera a densidad en amplitud, corre el cruce hacia")
	fmt.Println("   0,863, empina la pendiente, genera T4 — y MUERE al aleatorizar las fases (e).")

	var Tp []float64
	for _, p := range []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} {
		Tp = append(Tp, math.Log(float64(p)))
	}
	g := cerosPaso(4000, 0.02)
	d := &dado{s: 20260824}
	fmt.Printf("   %d ceros reales de referencia · pipeline de Fase XV intacto\n", len(g))

	Ereal, crReal, pendReal := pipeline(g, Tp, d, 120)
	T4real := descomponer(g, Tp)

	type brazo17 struct {
		nom           string
		E, T4         []float64
		cruces, pends []float64
		amp           float64
	}
	corre := func(nom string, gen func(*dado) []float64, sem int) brazo17 {
		br := brazo17{nom: nom}
		for r := 0; r < sem; r++ {
			sp := gen(d)
			E, c, p := pipeline(sp, Tp, d, 60)
			if r == 0 {
				br.E = E
				br.T4 = descomponer(sp, Tp)
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

	fmt.Println("\n   corriendo los brazos (3 semillas cada uno)…")
	bA := corre("a) GUE puro", func(dd *dado) []float64 { return surrogate(4000, nil, dd) }, 3)
	bB := corre("b) GUE + densidad", func(dd *dado) []float64 { return surrogate(4000, epsilonPrimos, dd) }, 3)
	bC := corre("c) GUE + DESPLAZAMIENTO", func(dd *dado) []float64 {
		return desplazar(surrogate(4000, nil, dd), 1, nil)
	}, 3)
	bD := corre("d) densidad + desplazamiento", func(dd *dado) []float64 {
		return desplazar(surrogate(4000, epsilonPrimos, dd), 1, nil)
	}, 3)
	bE := corre("e) fase DESTRUIDA (φ al azar)", func(dd *dado) []float64 {
		f := make([]float64, len(nLam))
		for i := range f {
			f[i] = 2 * math.Pi * dd.u()
		}
		return desplazar(surrogate(4000, nil, dd), 1, f)
	}, 3)
	bC5 := corre("f) desplazamiento ×0,5", func(dd *dado) []float64 {
		return desplazar(surrogate(4000, nil, dd), 0.5, nil)
	}, 2)
	bC2 := corre("f) desplazamiento ×2", func(dd *dado) []float64 {
		return desplazar(surrogate(4000, nil, dd), 2, nil)
	}, 2)

	fmt.Println("\n§1 · LA TABLA QUE IMPORTA (su §7)")
	fmt.Printf("   %-32s %9s %18s %18s\n", "modelo", "amplitud", "cruce s*", "pendiente")
	fmt.Printf("   %-32s %9.3f %18.3f %18.3f\n", "REAL (3474 ceros)", resumen(Ereal), crReal, pendReal)
	for _, br := range []brazo17{bA, bB, bC, bD, bE, bC5, bC2} {
		fmt.Printf("   %-32s %9.3f %18s %18s\n", br.nom, br.amp, listaF(br.cruces), listaF(br.pends))
	}

	fmt.Println("\n§2 · LA PRUEBA CLAVE (su §9): ¿el desplazamiento genera T4?")
	fmt.Printf("   %-11s %9s %9s %9s\n", "bin s", "T4 real", "T4 despl.", "T4 densid.")
	nb := int((zHi - zLo) / zW)
	for b := 0; b < nb; b++ {
		fmt.Printf("   %.2f–%.2f %+9.4f %+9.4f %+9.4f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW, T4real[b], bC.T4[b], bB.T4[b])
	}

	fmt.Println("\n§3 · RESIDUO respecto del REAL, brazo por brazo (máx |real − modelo|)")
	for _, br := range []brazo17{bB, bC, bD} {
		rmax := 0.0
		for b := range Ereal {
			if r := math.Abs(Ereal[b] - br.E[b]); r > rmax {
				rmax = r
			}
		}
		fmt.Printf("   %-32s residuo máx %.4f\n", br.nom, rmax)
	}

	fmt.Println("\n§4 · VEREDICTO contra sus criterios del §17–18 — lo dicta la tabla:")
	fmt.Println("   éxito = c) mejora a b) en amplitud/cruce/pendiente Y e) pierde la señal.")
	fmt.Println("   fracaso = c) no mejora, o la señal sobrevive a la fase destruida.")

	dibujar17(Ereal, bB.E, bC.E, bE.E, crReal, bC.cruces, bE.amp)
}

func listaF(v []float64) string {
	s := ""
	for i, x := range v {
		if i > 0 {
			s += "/"
		}
		s += fmt.Sprintf("%.2f", x)
	}
	return s
}
