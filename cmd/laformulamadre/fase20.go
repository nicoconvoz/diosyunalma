package main

// fase20.go - THE TRUNCATION TEST. Her strict order: ADD NOTHING. The pure
// field model (the rigid picket pinned by the counting equation, arm F of
// Phase XVIII) with ONE experimental variable: the truncation of S.
//
//     N_smooth(g) + S_N(g) = k - 1/2
//     S_N(g) = -(1/pi) sum_{n<=N, Lambda>0} Lambda(n) sin(g log n)/(sqrt(n) log n)
//
// Rungs DECLARED BEFORE RUNNING: N in {97, 997, 9973, 99991}. Everything else
// frozen: same N_smooth, same gamma range, same E(s) pipeline, same T-terms,
// same crossing and slope procedure, same spacing statistic.
//
// The model is DETERMINISTIC (no seed): stability is checked by re-running the
// shuffle-null estimation with a second seed (declared; her section 10 allows
// documenting this limitation).
//
// The pinning walks FORWARD to the first upward crossing of the counting
// equation (Berry's reading when N_smooth + S is locally non-monotone), then
// bisects: no Newton jumps, fully reproducible.
//
// The three convergences sought (her section 7): amplitude 0.285 -> 0.346,
// crossing 0.94 -> 0.862, spacing sd staying at 0.42. Plus T4 per rung.
// If all three converge: STOP AND AUDIT, do not interpret (her section 16).

import (
	"fmt"
	"math"
)

type compS struct{ ln, w float64 }

// componentes builds Lambda-weighted components for all prime powers <= tope.
func componentes(tope int) []compS {
	esP := criba20(tope)
	var cs []compS
	for p := 2; p <= tope; p++ {
		if !esP[p] {
			continue
		}
		lp := math.Log(float64(p))
		for q := p; q <= tope; q *= p {
			lq := math.Log(float64(q))
			cs = append(cs, compS{lq, lp / (math.Sqrt(float64(q)) * lq)})
			if q > tope/p {
				break
			}
		}
	}
	return cs
}

func criba20(tope int) []bool {
	c := make([]bool, tope+1)
	for i := 2; i <= tope; i++ {
		c[i] = true
	}
	for i := 2; i*i <= tope; i++ {
		if c[i] {
			for j := i * i; j <= tope; j += i {
				c[j] = false
			}
		}
	}
	return c
}

func sDeComp(g float64, cs []compS) float64 {
	s := 0.0
	for _, c := range cs {
		s += c.w * math.Sin(g*c.ln)
	}
	return -s / math.Pi
}

// picketPuro pins the deterministic picket u_k = u_0 + k by marching to the
// first upward crossing, then bisecting.
func picketPuro(tope float64, cs []compS) []float64 {
	f := func(g, u float64) float64 { return theta(g)/math.Pi + 1 + sDeComp(g, cs) - u }
	u := theta(30)/math.Pi + 1
	g := 30.0
	var out []float64
	for {
		u += 1
		a := g + 1e-6
		fa := f(a, u)
		paso := 0.02
		b := a
		fb := fa
		for fb < 0 {
			b += paso
			fb = f(b, u)
			if b > tope+5 {
				break
			}
		}
		if b > tope+5 {
			break
		}
		lo, hi := b-paso, b
		for i := 0; i < 45; i++ {
			m := (lo + hi) / 2
			if f(m, u) < 0 {
				lo = m
			} else {
				hi = m
			}
		}
		g = (lo + hi) / 2
		if g > tope {
			break
		}
		out = append(out, g)
	}
	return out
}

func fase20() {
	fmt.Println("🪞🔬 FASE XX — LA PRUEBA DE LA TRUNCACIÓN: más fórmula, sin más ruido")
	fmt.Println("   NO SE AGREGA NADA. Única variable: el tope N de la suma de S.")
	fmt.Println("   escalones declarados antes de correr: N ∈ {97, 997, 9973, 99991}")
	fmt.Println("   convergencias buscadas (su §7): amplitud 0,285→0,346 · cruce 0,94→0,862 ·")
	fmt.Println("   espaciados quietos en 0,42 · y T4 por escalón. Si las tres convergen:")
	fmt.Println("   FRENAR Y AUDITAR, no interpretar (su §16).")

	var Tp []float64
	for _, p := range []int{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} {
		Tp = append(Tp, math.Log(float64(p)))
	}
	g := cerosPaso(4000, 0.02)
	d := &dado{s: 20260827}
	Ereal, crReal, pendReal := pipeline(g, Tp, d, 120)
	T4real := descomponer(g, Tp)
	estad := func(sp []float64) float64 {
		var ss []float64
		for _, p := range paresDe(sp) {
			ss = append(ss, sDe(p, true))
		}
		return desvio(ss)
	}

	nb := int((zHi - zLo) / zW)
	fmt.Printf("\n   %-10s %9s %10s %11s %8s %9s %9s\n",
		"N", "amplitud", "cruce s*", "pendiente", "s desv", "T4 bin1", "resid máx")
	fmt.Printf("   %-10s %9.3f %10.3f %11.3f %8s %9.3f %9s\n",
		"REAL", resumen(Ereal), crReal, pendReal, "0.42", T4real[0], "—")

	type filon struct {
		N              int
		E, T4          []float64
		cr, pend, sdev float64
		amp, rmax      float64
	}
	var filas []filon
	for _, N := range []int{97, 997, 9973, 99991} {
		cs := componentes(N)
		sp := picketPuro(4000, cs)
		E, cr, pend := pipeline(sp, Tp, d, 60)
		T4 := descomponer(sp, Tp)
		rmax := 0.0
		for b := 0; b < nb; b++ {
			if r := math.Abs(Ereal[b] - E[b]); r > rmax {
				rmax = r
			}
		}
		fl := filon{N, E, T4, cr, pend, estad(sp), resumen(E), rmax}
		filas = append(filas, fl)
		fmt.Printf("   %-10d %9.3f %10.3f %11.3f %8.2f %9.3f %9.4f   (%d componentes, %d puntos)\n",
			N, fl.amp, fl.cr, fl.pend, fl.sdev, fl.T4[0], fl.rmax, len(cs), len(sp))
	}

	fmt.Println("\n§2 · E(s) POR ESCALÓN, superpuestas (bin a bin)")
	fmt.Printf("   %-11s %8s %8s %8s %8s %8s\n", "bin s", "N=97", "N=997", "N=9973", "N=99991", "REAL")
	for b := 0; b < nb; b++ {
		fmt.Printf("   %.2f–%.2f %+8.3f %+8.3f %+8.3f %+8.3f %+8.3f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW,
			filas[0].E[b], filas[1].E[b], filas[2].E[b], filas[3].E[b], Ereal[b])
	}

	fmt.Println("\n§3 · T4 POR ESCALÓN (¿más términos → más T4 coherente?)")
	fmt.Printf("   %-11s %8s %8s %8s %8s %8s\n", "bin s", "N=97", "N=997", "N=9973", "N=99991", "REAL")
	for b := 0; b < nb; b += 3 {
		fmt.Printf("   %.2f–%.2f %+8.3f %+8.3f %+8.3f %+8.3f %+8.3f\n",
			zLo+float64(b)*zW, zLo+float64(b+1)*zW,
			filas[0].T4[b], filas[1].T4[b], filas[2].T4[b], filas[3].T4[b], T4real[b])
	}

	// stability: the model is deterministic; re-estimate nulls with a 2nd seed
	fmt.Println("\n§4 · ESTABILIDAD (su §10): el modelo es DETERMINISTA — la única fuente")
	fmt.Println("   aleatoria es la estimación del nulo. Se repite el escalón N=9973 con otra")
	fmt.Println("   semilla de nulos:")
	d2 := &dado{s: 777777}
	spR := picketPuro(4000, componentes(9973))
	E2, cr2, pend2 := pipeline(spR, Tp, d2, 60)
	fmt.Printf("   semilla A: amplitud %.3f · cruce %.3f · pendiente %.3f\n", filas[2].amp, filas[2].cr, filas[2].pend)
	fmt.Printf("   semilla B: amplitud %.3f · cruce %.3f · pendiente %.3f\n", resumen(E2), cr2, pend2)

	fmt.Println("\n§5 · LECTURA con su tabla de decisión del §18 — la dictan las columnas:")
	fmt.Println("   amplitud hacia 0,346 · cruce hacia 0,862 · s desv quieto en 0,42 · T4 hacia 0,267")
	fmt.Println("   Si las tres convergen: FRENAR Y AUDITAR. Si no: documentar dónde falla.")

	dibujar20(Ereal, filas[0].E, filas[2].E, filas[3].E, crReal)
}
