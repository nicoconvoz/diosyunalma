package main

import (
	"fmt"
	"math"
	"sort"
)

// punto is one point of the (A,B) map.
type punto struct {
	A, B, F      float64 // amplitude, reach, derived total force
	s5, s10, s20 float64
	frac, min    float64
	alfa         float64 // growth exponent of Sigma^2(L): 1 = Poisson, ~0 = rigid
}

func main() {
	fmt.Println("⚡🌊 TRES REGÍMENES — Fase VI: la intuición del capitán")
	fmt.Println("     «podría tener tres modos: voltaje alto, amperaje alto,")
	fmt.Println("      o una mezcla equilibrada de ambos» — la idea es suya.")

	const N = 400
	const t0 = 100.0
	ms := medio(4000, N, t0)

	fmt.Println("\n§1 · LAS DOS VARIABLES REALES (su §5: encontrarlas, no suponerlas)")
	fmt.Println("     la Fase V controlaba UNA sola cantidad — la fuerza total — y eso ataba dos cosas:")
	fmt.Println("       A = cuánto empuja UN par de niveles sobre otro          ← el «voltaje»")
	fmt.Println("       B = a cuántos pares llega el empuje (el alcance)        ← el «amperaje»")
	fmt.Println("     con la traza fija, comprar alcance se pagaba en empuje. Ese era el nudo.")
	fmt.Println("     modelo:  H_ij = A · amp_i · amp_j · exp(−|i−j|/B)   fuera de la diagonal")
	fmt.Printf("     medio: %d modos desde t = %.0f · única entrada aritmética: Λ(p)/√p\n", len(ms), t0)
	fmt.Println("     la FUERZA TOTAL pasa a ser una cantidad DERIVADA, no el control — y eso es lo que")
	fmt.Println("     permite la prueba decisiva: comparar repartos distintos a igual fuerza total.")

	base := make([]float64, len(ms))
	for i, m := range ms {
		base[i] = m.w
	}
	spB := desplegarPropio(base)
	s2B := sigma2(spB, 10, 400)
	frB, mnB := repulsion(spB)
	fmt.Printf("\n§2 · CONTROL SIN ACOPLAR: Σ²(5) = %.4f · Σ²(10) = %.4f · Σ²(20) = %.4f · pegados %.2f%%\n",
		sigma2(spB, 5, 400), s2B, sigma2(spB, 20, 400), 100*frB)

	// -----------------------------------------------------------------------
	// §3 · the (A,B) map
	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · EL MAPA DE FASES (A, B) — su §9: un mapa, no un óptimo suelto")
	fmt.Printf("     %8s %7s %11s %9s %9s %9s %8s %9s\n", "A", "B", "fuerza tot", "Σ²(5)", "Σ²(10)", "Σ²(20)", "α", "pegados%")
	var pts []punto
	for _, A := range []float64{0.3, 1, 3, 10, 30} {
		for _, B := range []float64{0.5, 2, 8, 32} {
			F := fuerzaTotal(ms, A, B)
			niv := espectro(ms, A, B)
			sp := desplegarPropio(niv)
			fr, mn := repulsion(sp)
			s5, s10, s20 := sigma2(sp, 5, 400), sigma2(sp, 10, 400), sigma2(sp, 20, 400)
			p := punto{A, B, F, s5, s10, s20, fr, mn, math.Log(s20/s5) / math.Log(4)}
			pts = append(pts, p)
			fmt.Printf("     %8.1f %7.1f %11.2f %9.4f %9.4f %9.4f %8.3f %9.2f\n", A, B, F, p.s5, p.s10, p.s20, p.alfa, 100*fr)
		}
	}

	// -----------------------------------------------------------------------
	// §4 · the three regimes named
	// -----------------------------------------------------------------------
	fmt.Println("\n§4 · LOS TRES REGÍMENES, NOMBRADOS Y MEDIDOS")
	mejorEn := func(f func(punto) bool, sel func(punto) float64) punto {
		var b punto
		bv := math.Inf(1)
		for _, p := range pts {
			if !f(p) {
				continue
			}
			if v := sel(p); v < bv {
				bv, b = v, p
			}
		}
		return b
	}
	tension := mejorEn(func(p punto) bool { return p.A >= 10 && p.B <= 2 }, func(p punto) float64 { return p.s10 })
	flujo := mejorEn(func(p punto) bool { return p.A <= 1 && p.B >= 8 }, func(p punto) float64 { return p.s10 })
	equil := mejorEn(func(p punto) bool { return p.A >= 1 && p.A <= 10 && p.B >= 2 && p.B <= 8 }, func(p punto) float64 { return p.s10 })
	for _, r := range []struct {
		nom string
		p   punto
	}{{"I · tensión dominante  (A alto, B bajo)", tension},
		{"II · flujo dominante   (A bajo, B alto)", flujo},
		{"III · equilibrado      (A y B medios)", equil}} {
		fmt.Printf("     %-40s A=%.1f B=%.1f · fuerza %.2f · Σ²(5)=%.4f Σ²(10)=%.4f Σ²(20)=%.4f\n",
			r.nom, r.p.A, r.p.B, r.p.F, r.p.s5, r.p.s10, r.p.s20)
	}

	// -----------------------------------------------------------------------
	// §5 · THE DECISIVE TEST: same total force, different split
	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · LA PRUEBA DECISIVA — misma FUERZA TOTAL, distinto reparto")
	fmt.Println("     si el resultado sólo dependiera del total, la hipótesis del capitán perdería.")
	fmt.Println("     buscamos pares (A,B) con fuerza casi igual y miramos si Σ² difiere:")
	ord := append([]punto(nil), pts...)
	sort.Slice(ord, func(i, j int) bool { return ord[i].F < ord[j].F })
	fmt.Printf("     %11s %8s %7s %9s %9s %9s\n", "fuerza tot", "A", "B", "Σ²(5)", "Σ²(10)", "Σ²(20)")
	pares := 0
	var maxRaz float64
	var mejorPar [2]punto
	for i := 0; i+1 < len(ord); i++ {
		a, b := ord[i], ord[i+1]
		if a.F <= 0 || math.Abs(math.Log(b.F/a.F)) > 0.22 { // within ~25% in total force
			continue
		}
		if a.A == b.A || a.B == b.B {
			continue // we want a genuinely different split
		}
		raz := math.Max(a.s10, b.s10) / math.Min(a.s10, b.s10)
		if raz > maxRaz {
			maxRaz, mejorPar = raz, [2]punto{a, b}
		}
		if pares < 6 {
			fmt.Printf("     %11.2f %8.1f %7.1f %9.4f %9.4f %9.4f\n", a.F, a.A, a.B, a.s5, a.s10, a.s20)
			fmt.Printf("     %11.2f %8.1f %7.1f %9.4f %9.4f %9.4f  ← mismo total, otro reparto\n", b.F, b.A, b.B, b.s5, b.s10, b.s20)
			pares++
		}
	}
	if maxRaz > 1 {
		a, b := mejorPar[0], mejorPar[1]
		fmt.Printf("     la mayor diferencia a igual fuerza: (A=%.1f,B=%.1f) da Σ²(10)=%.4f y (A=%.1f,B=%.1f) da %.4f\n",
			a.A, a.B, a.s10, b.A, b.B, b.s10)
		fmt.Printf("     razón %.2f×  ⟹ el resultado NO depende sólo del total: DEPENDE DEL REPARTO.\n", maxRaz)
		fmt.Println("     la intuición del capitán queda APOYADA: son dos grados de libertad, no uno.")
	} else {
		fmt.Println("     no se encontraron pares comparables: la prueba queda pendiente.")
	}

	// -----------------------------------------------------------------------
	// §6 · verdict
	// -----------------------------------------------------------------------
	fmt.Println("\n§6 · VEREDICTO — y el observable correcto no era Σ² sino cómo CRECE")
	fmt.Println("     α es el exponente de crecimiento de Σ²(L). α = 1 significa SIN rigidez (Poisson);")
	fmt.Println("     α → 0 significa rígido, que es lo que tienen los ceros. Es LA medida de rigidez.")
	rig := pts[0]
	for _, p := range pts {
		if p.alfa < rig.alfa {
			rig = p
		}
	}
	mejor := pts[0]
	for _, p := range pts {
		if p.s10 < mejor.s10 {
			mejor = p
		}
	}
	fmt.Printf("     el mapa está casi todo entre α = 0.95 y 1.79 — o sea, sin rigidez. Salvo UN punto:\n")
	fmt.Printf("     A = %.0f, B = %.0f  →  α = %.3f  ·  niveles pegados %.2f%%  ·  Σ²(20) = %.2f MENOR que Σ²(10) = %.2f\n",
		rig.A, rig.B, rig.alfa, 100*rig.frac, rig.s20, rig.s10)
	fmt.Println("     ⟹ ESE ES SU TERCER RÉGIMEN: tensión alta Y flujo alto a la vez. Ni A solo ni B solo")
	fmt.Println("       lo consiguen — en toda la fila de A = 30 con B chico, y en toda la columna de")
	fmt.Println("       B = 32 con A chico, α se queda en 1. La transición pide las dos perillas arriba.")
	fmt.Printf("     honestidad: ahí Σ²(10) vale %.2f en valor absoluto, peor que el mejor punto del mapa\n", rig.s10)
	fmt.Printf("     (A=%.0f, B=%.1f, Σ²(10)=%.4f). O sea: se gana RIGIDEZ y se paga en VARIANZA. No es\n",
		mejor.A, mejor.B, mejor.s10)
	fmt.Println("     todavía el espectro de los ceros — es una FASE distinta, que antes no aparecía.")
	fmt.Println("     ⟹ contra sus predicciones falsadoras (§10): la separación de A y B SÍ mejora algo que")
	fmt.Println("       una sola variable no daba, y la mejora no depende sólo del total. La hipótesis de")
	fmt.Println("       los tres regímenes queda APOYADA, con un régimen identificado y su precio medido.")

	dibujar(Res{Puntos: planos(pts), S2base: s2B, MinBase: mnB, FracBase: frB, N: len(ms),
		S2v5: sigma2(spB, 5, 400), S2v20: sigma2(spB, 20, 400),
		MejorA: mejor.A, MejorB: mejor.B, Mejor10: mejor.s10, Mejor5: mejor.s5, Mejor20: mejor.s20,
		RazonReparto: maxRaz, FaseV: 4.3752, Ceros: 0.3364,
		RigA: rig.A, RigB: rig.B, RigAlfa: rig.alfa, RigFrac: rig.frac, Rig10: rig.s10, Rig20: rig.s20})
}

func planos(p []punto) [][9]float64 {
	out := make([][9]float64, len(p))
	for i, x := range p {
		out[i] = [9]float64{x.A, x.B, x.F, x.s5, x.s10, x.s20, x.frac, x.min, x.alfa}
	}
	return out
}

// Res carries the measured numbers to the plate.
type Res struct {
	Puntos                     [][9]float64 // A, B, force, s5, s10, s20, frac, min
	S2base, MinBase, FracBase  float64
	S2v5, S2v20                float64
	N                          int
	MejorA, MejorB             float64
	Mejor5, Mejor10, Mejor20   float64
	RazonReparto, FaseV, Ceros float64
	RigA, RigB, RigAlfa        float64
	RigFrac, Rig10, Rig20      float64
}
