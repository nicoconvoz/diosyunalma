package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	fmt.Println("🌊 EL FLUIDO — Fase IV: la intuición del capitán, puesta a prueba")
	fmt.Println("   «algunas ondas grandes, otras pequeñas, otras medianas, todas resonando")
	fmt.Println("    y creando una melodía única» — la idea es suya; acá se la mide.")

	const t0, t1 = 100.0, 1000.0
	const topeT = 3.2
	verdad := cerosVerdaderos(t0, t1, 0.02)
	spV := make([]float64, 0, len(verdad))
	for i := 0; i+1 < len(verdad); i++ {
		spV = append(spV, suave(verdad[i+1])-suave(verdad[i]))
	}
	s2V := sigma2(spV, 10, 600)
	fracV, minV := repulsion(spV)

	var per []float64
	for n := 2; n <= 4096; n++ {
		l := math.Log(float64(n))
		if l > topeT || l <= 0.3 {
			continue
		}
		// prime powers only
		q, e := n, 0
		for d := 2; d*d <= q; d++ {
			if q%d == 0 {
				for q%d == 0 {
					q /= d
					e++
				}
				break
			}
		}
		if q == 1 || e == 0 {
			per = append(per, l)
		}
	}
	fmt.Printf("\n§0 · la regla · %d ceros verdaderos · Σ²(10) = %.4f · espaciados bajo 0,1: %.2f%% · mínimo %.4f\n",
		len(verdad), s2V, 100*fracV, minV)

	// -----------------------------------------------------------------------
	// §1 · R6 declared first
	// -----------------------------------------------------------------------
	fmt.Println("\n§1 · R6 — el medio, declarado antes de mirar nada")
	fmt.Println("     UN primo p = UNA excitación de escala log p, que aporta la escalera 2πk/log p.")
	fmt.Println("     TODAS las excitaciones se acoplan a UN mismo modo común (el fluido), con la fuerza")
	fmt.Println("     que la propia fórmula explícita le da a esa excitación: Λ(p)/√p = log p · p^(−1/2).")
	fmt.Println("     entradas: Λ(n) y nada más. Ningún parámetro ajustado a los γₙ. El acoplamiento es")
	fmt.Println("     de rango uno, así que el espectro acoplado es EXACTO (raíces de la secular).")

	ex := excitaciones(400, 1)
	ms := ladrillos(ex, t0, t1)
	fmt.Printf("     %d excitaciones (primos ≤ 400) → %d modos del medio en [%.0f,%.0f]\n", len(ex), len(ms), t0, t1)

	// -----------------------------------------------------------------------
	// §2 · does the coupling do anything at all? the ladder of g
	// -----------------------------------------------------------------------
	fmt.Println("\n§2 · ¿EL MEDIO HACE ALGO? — subiendo el acoplamiento desde cero")
	fmt.Printf("     %10s %12s %12s %14s %12s\n", "g", "Σ²(10)", "mín. espac.", "bajo 0,1 (%)", "eco k·log p")
	var filas []fila
	for _, g := range []float64{0, 0.01, 0.1, 1, 10, 100} {
		niv := acoplar(ms, g)
		sp := desplegarPropio(niv)
		s2 := sigma2(sp, 10, 600)
		fr, mn := repulsion(sp)
		ec := murcielago(niv, per, topeT)
		fmt.Printf("     %10.2f %12.4f %12.2e %14.2f %12.3f\n", g, s2, mn, 100*fr, ec)
		filas = append(filas, fila{g, s2, mn, fr, ec})
	}
	fmt.Printf("     los ceros verdaderos, para comparar: Σ² = %.4f · mínimo %.4f · bajo 0,1 = %.2f%%\n",
		s2V, minV, 100*fracV)

	// -----------------------------------------------------------------------
	// §3 · the auditor's own experiment: add excitations one by one
	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · EL EXPERIMENTO QUE PIDIÓ LA AUDITORA (§8) — excitaciones de a una")
	fmt.Println("     ¿mejoran las correlaciones de forma COLECTIVA, o sólo se superponen espectros?")
	fmt.Printf("     %8s %8s %14s %14s %14s\n", "primos", "modos", "Σ² sin acoplar", "Σ² acoplado", "mejora")
	var evol [][5]float64
	for _, np := range []int{2, 4, 8, 16, 32, 78} {
		exN := excitaciones(1000, 1)
		if np < len(exN) {
			exN = exN[:np]
		}
		msN := ladrillos(exN, t0, t1)
		if len(msN) < 40 {
			continue
		}
		s0 := sigma2(desplegarPropio(acoplar(msN, 0)), 10, 600)
		s1 := sigma2(desplegarPropio(acoplar(msN, 10)), 10, 600)
		fmt.Printf("     %8d %8d %14.4f %14.4f %14.4f\n", len(exN), len(msN), s0, s1, s0-s1)
		evol = append(evol, [5]float64{float64(len(exN)), float64(len(msN)), s0, s1, s0 - s1})
	}

	// -----------------------------------------------------------------------
	// §4 · WHY: what a rank-one medium can and cannot do
	// -----------------------------------------------------------------------
	fmt.Println("\n§4 · POR QUÉ — lo que un medio de UN SOLO modo común puede y no puede hacer")
	niv := acoplar(ms, 10)
	sp := desplegarPropio(niv)
	fr, mn := repulsion(sp)
	fmt.Printf("     con g = 10: Σ²(10) = %.4f (sin acoplar %.4f · los ceros %.4f)\n",
		sigma2(sp, 10, 600), filas[0].s2, s2V)
	fmt.Printf("     espaciado mínimo %.2e contra %.2e sin acoplar — y %.4f en los ceros\n", mn, filas[0].mn, minV)
	fmt.Printf("     espaciados bajo 0,1: %.2f%% contra %.2f%% sin acoplar — los ceros: %.2f%%\n",
		100*fr, 100*filas[0].fr, 100*fracV)
	fmt.Println("     el acoplamiento de rango uno INTERCALA: mete exactamente una raíz entre cada par de")
	fmt.Println("     niveles vecinos. Empuja, pero no puede empujar más allá de sus vecinos: el")
	fmt.Println("     entrelazado de Cauchy se lo prohíbe. Un solo modo común no alcanza para fabricar")
	fmt.Println("     repulsión de largo alcance — mueve cada nivel dentro de su propia celda.")
	fmt.Println("     ⟹ la intuición del fluido NO muere acá: lo que muere es la versión de UN SOLO modo.")
	fmt.Println("       Lo que el experimento dice es que el medio necesita MUCHOS modos comunes —")
	fmt.Println("       o sea, un acoplamiento de rango alto — para que las excitaciones se sientan")
	fmt.Println("       entre sí más allá de la vecina inmediata. Eso es una predicción medible.")

	// -----------------------------------------------------------------------
	// §5 · the honest verdict against R6 and the auditor's success criterion
	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · VEREDICTO CONTRA SU CRITERIO DE ÉXITO (§16)")
	fmt.Println("     «éxito NO es que la simulación se parezca a los ceros; éxito sería que el medio")
	fmt.Println("      genere una estructura que no estaba codificada en la salida».")
	fmt.Printf("     R6            : LIMPIO — sólo Λ(n); ningún parámetro tocado contra los γₙ\n")
	fmt.Printf("     eco aritmético: %.1f — pero es TAUTOLÓGICO: las escalas log p están en la definición\n", filas[len(filas)-1].ec)
	fmt.Printf("     correlaciones : %.4f contra %.4f de los ceros — el medio de un modo NO las produce\n",
		sigma2(sp, 10, 600), s2V)
	fmt.Println("     conteo        : densidad constante — el medio no respira; ese frente sigue abierto")
	fmt.Println("     ⟹ resultado NEGATIVO y ÚTIL: el fluido de un solo modo común queda descartado,")
	fmt.Println("       con la razón exacta (entrelazado de Cauchy), y la hipótesis del capitán queda")
	fmt.Println("       viva en su versión fuerte: un medio con muchos modos, que es lo próximo a medir.")

	dibujar(R{
		Verdad: verdad, S2v: s2V, FracV: fracV, MinV: minV, Modos: len(ms), Exc: len(ex),
		Filas: filasPlanas(filas), Evol: evol, T0: t0, T1: t1,
		NivSin: acoplar(ms, 0), NivCon: niv,
	})
}

// fila is one coupling strength of the medium.
type fila struct{ g, s2, mn, fr, ec float64 }

func filasPlanas(f []fila) [][5]float64 {
	out := make([][5]float64, len(f))
	for i, x := range f {
		out[i] = [5]float64{x.g, x.s2, x.mn, x.fr, x.ec}
	}
	return out
}

var _ = sort.Float64s

// R carries the measured numbers to the plate.
type R struct {
	Verdad           []float64
	S2v, FracV, MinV float64
	Modos, Exc       int
	Filas            [][5]float64 // g, sigma2, min spacing, frac below 0.1, echo
	Evol             [][5]float64 // primes, modes, sigma2 uncoupled, coupled, gain
	T0, T1           float64
	NivSin, NivCon   []float64
}
