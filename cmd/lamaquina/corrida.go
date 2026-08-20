package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	fmt.Println("🔥 LA MÁQUINA — Fase III: de primo a órbita")
	fmt.Println("   no una lista que imite los ceros: una máquina que no necesite conocerlos.")

	const t0, t1 = 100.0, 1000.0
	const topeT = 3.2
	verdad := cerosVerdaderos(t0, t1, 0.02)

	// -----------------------------------------------------------------------
	// §1 · R6: the input list, declared before anything is tested
	// -----------------------------------------------------------------------
	fmt.Println("\n§1 · R6 — entradas declaradas ANTES de probar nada")
	fmt.Println("     admisibles usadas: Λ(n) (o sea, los primos). Nada más.")
	fmt.Println("     los γₙ aparecen SÓLO como regla de medir, nunca dentro de una definición.")

	// -----------------------------------------------------------------------
	// §2 · THE CONCATENATION OBSTRUCTION - the new no-go, with its hypotheses
	// -----------------------------------------------------------------------
	fmt.Println("\n§2 · LA OBSTRUCCIÓN DE CONCATENACIÓN — un no-go nuevo, con sus hipótesis exactas")
	fmt.Println("     HIPÓTESIS: (i) las órbitas primitivas del sistema tienen longitudes {log p};")
	fmt.Println("                (ii) dos órbitas cerradas que comparten un punto se pueden CONCATENAR")
	fmt.Println("                     en otra órbita cerrada (cierto en todo grafo cuántico conexo y en")
	fmt.Println("                     todo flujo con un punto recurrente común).")
	fmt.Println("     CONSECUENCIA: el espectro de longitudes es cerrado bajo la suma, o sea contiene")
	fmt.Println("                   log a + log b = log(ab) para TODO par. Contiene log n para todo n.")
	fmt.Println("     PERO la fórmula explícita le da peso Λ(n) = 0 a todo n que no sea potencia de primo.")

	const N = 4096
	var forzadas, permitidas, prohibidas int
	var ejemplos []string
	for n := 2; n <= N; n++ {
		l := lambda(n)
		// a length log n is FORCED on a concatenating system as soon as n factors
		// into at least two prime factors counted with multiplicity
		if n > 1 {
			forzadas++
			if l > 0 {
				permitidas++
			} else {
				prohibidas++
				if len(ejemplos) < 8 {
					ejemplos = append(ejemplos, fmt.Sprintf("log %d", n))
				}
			}
		}
	}
	fmt.Printf("     hasta n = %d: el sistema que concatena está OBLIGADO a tener %d longitudes;\n", N, forzadas)
	fmt.Printf("     la fórmula explícita sólo permite %d (las potencias de primo) y PROHÍBE %d — el %.1f%%.\n",
		permitidas, prohibidas, 100*float64(prohibidas)/float64(forzadas))
	fmt.Printf("     las primeras prohibidas: %v …\n", ejemplos)
	fmt.Println("     ⟹ NO-GO: ningún sistema cuyas órbitas se puedan concatenar puede tener el espectro")
	fmt.Println("       de longitudes de la fórmula explícita. O las órbitas NO se tocan, o las")
	fmt.Println("       concatenaciones se CANCELAN exactamente — y cancelar exactamente Λ(n) = 0 para")
	fmt.Println("       todo compuesto es una condición extraordinariamente fuerte, no un detalle técnico.")
	fmt.Println("     (esto explica de otra manera por qué el candidato tiene que ser raro: en un sistema")
	fmt.Println("      dinámico normal las órbitas se concatenan, y acá no pueden.)")

	// -----------------------------------------------------------------------
	// §3 · CANDIDATE A: the disjoint ladders - the only escape from §2
	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · CANDIDATO A — LAS ESCALERAS DISJUNTAS (la única salida de §2)")
	fmt.Println("     una órbita por primo, sin tocarse: el espectro es la UNIÓN de las escaleras 2πm/log p.")
	ps := primos(400)
	var niv []float64
	for _, p := range ps {
		lp := math.Log(float64(p))
		for m := 1; ; m++ {
			e := 2 * math.Pi * float64(m) / lp
			if e > t1 {
				break
			}
			if e >= t0 {
				niv = append(niv, e)
			}
		}
	}
	sort.Float64s(niv)
	fmt.Printf("     con los %d primos hasta 400: %d niveles en [%.0f,%.0f] contra %d ceros verdaderos\n",
		len(ps), len(niv), t0, t1, len(verdad))
	// unfold with its OWN mean density (constant), which is what a fixed union gives
	dens := float64(len(niv)) / (t1 - t0)
	sp := make([]float64, 0, len(niv))
	for i := 0; i+1 < len(niv); i++ {
		sp = append(sp, (niv[i+1]-niv[i])*dens)
	}
	var per []float64
	for n := 2; n <= 4096; n++ {
		if l := math.Log(float64(n)); lambda(n) > 0 && l <= topeT && l > 0.3 {
			per = append(per, l)
		}
	}
	spV := make([]float64, 0, len(verdad))
	for i := 0; i+1 < len(verdad); i++ {
		spV = append(spV, suave(verdad[i+1])-suave(verdad[i]))
	}
	fmt.Printf("     CONTEO      : densidad CONSTANTE %.4f/unidad — los ceros piden (1/2π)ln(T/2π), que crece ✗\n", dens)
	fmt.Printf("     CORRELACIÓN : Σ²(10) = %.4f  ·  los ceros dan %.4f  ·  espaciado mínimo %.2e ✗\n",
		sigma2(sp, 10, 600), sigma2(spV, 10, 600), mini(sp))
	fmt.Printf("     ARITMÉTICA  : eco en k·log p = %.3f (los ceros: medido en F350 = 182.07)\n",
		murcielago(niv, per, topeT))
	fmt.Println("     ⟹ las escaleras disjuntas SÍ tienen aritmética genuina (cada escalera ES un primo)")
	fmt.Println("       pero su densidad es constante y sus niveles se amontonan: no hay repulsión, hay")
	fmt.Println("       superposición de relojes independientes. Muere en CONTEO y en CORRELACIONES.")

	// -----------------------------------------------------------------------
	// §4 · THE WEIGHT, DERIVED - what the explicit formula forces on the dynamics
	// -----------------------------------------------------------------------
	fmt.Println("\n§4 · EL PESO, DERIVADO (no impuesto) — lo que la fórmula explícita le exige a la dinámica")
	fmt.Println("     una órbita ESTABLE (un círculo) aporta peso ℓ y NO se amortigua con las repeticiones:")
	fmt.Println("     la fórmula explícita sí se amortigua, como p^(−m/2). Un círculo no puede ser un primo.")
	fmt.Println("     una órbita HIPERBÓLICA de longitud ℓ y exponente λ aporta  ℓ/(2·senh(mλ/2)),")
	fmt.Println("     que decae en m con TASA λ. La fórmula explícita decae con tasa log p. Igualando TASAS:")
	fmt.Printf("     %6s %5s %16s %16s %12s\n", "p", "m", "tasa Selberg", "tasa Weil", "cociente")
	for _, pp := range []int{2, 7, 97} {
		lp := math.Log(float64(pp))
		for _, m := range []int{2, 8, 32} {
			// decay rate in m of each amplitude, with lambda set to log p
			ampS := lp / (2 * math.Sinh(float64(m)*lp/2))
			tasaS := -2 * math.Log(ampS/lp) / float64(m)
			tasaW := -2 * math.Log(2*lp*math.Pow(float64(pp), -float64(m)/2)/(2*lp)) / float64(m)
			fmt.Printf("     %6d %5d %16.8f %16.8f %12.6f\n", pp, m, tasaS, tasaW, tasaS/tasaW)
		}
	}
	fmt.Println("     ⟹ las dos tasas coinciden cuando λ = log p — o sea, cuando el exponente de")
	fmt.Println("       inestabilidad de la órbita ES SU PROPIA LONGITUD. En unidades donde el tiempo es la")
	fmt.Println("       longitud, eso es EXPONENTE DE LYAPUNOV EXACTAMENTE 1, que es el flujo de xp (ẋ = x, ṗ = −p).")
	fmt.Println("     Ese es el puente entre las dos mitades: Berry-Keating aporta la inestabilidad correcta")
	fmt.Println("     y los primos aportan las longitudes. La máquina necesita las dos a la vez.")
	fmt.Println("     ⚠ CORRECCIÓN NUESTRA, en el mismo turno: primero despejamos λ de la amplitud COMPLETA")
	fmt.Println("       y dijimos que tendía a log p. Es falso: eso da λ = log p − (2·log 2)/m, y el cociente")
	fmt.Println("       λ/log p ni siquiera es monótono (2 da 1.000, 11 da 0.630, 97 da 0.714). Lo que la")
	fmt.Println("       fórmula fuerza es la TASA de decaimiento, no la constante — y esa sí da λ = log p exacto.")
	fmt.Println("       La constante que sobra es justamente el factor del §5.")

	// -----------------------------------------------------------------------
	// §5 · THE SIGN, and the residue after matching amplitudes
	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · EL SIGNO Y EL RESTO — lo que queda cuando las amplitudes se emparejan")
	fmt.Println("     con λ = ℓ = log p exacto, la razón entre lo que pide Weil y lo que da Selberg es")
	fmt.Printf("     %6s %4s %16s %16s %14s\n", "p", "m", "Selberg (+)", "Weil (−)", "razón")
	for _, p := range []int{2, 3, 5, 7, 11, 97} {
		lp := math.Log(float64(p))
		for _, m := range []int{1, 2} {
			sel := lp / (2 * math.Sinh(float64(m)*lp/2))
			wei := -2 * lp * math.Pow(float64(p), -float64(m)/2)
			fmt.Printf("     %6d %4d %16.8f %16.8f %14.6f\n", p, m, sel, wei, wei/sel)
		}
	}
	fmt.Println("     la razón es exactamente −2·(1 − p^(−m)): el signo NEGATIVO no sale de emparejar")
	fmt.Println("     amplitudes — sobrevive a todo el ajuste. Hay que DERIVARLO de otra cosa")
	fmt.Println("     (orientación, índice, grado cohomológico, o la estructura de absorción de Connes).")
	fmt.Println("     Y el factor (1 − p^(−m)) es el «problema asintótico» de Sierra, medido: la")
	fmt.Println("     amortiguación de Selberg y la de Weil coinciden sólo cuando m → ∞.")

	// -----------------------------------------------------------------------
	// §6 · THE EVALUATION MATRIX she asked for
	// -----------------------------------------------------------------------
	fmt.Println("\n§6 · LA MATRIZ DE EVALUACIÓN (su §9), llenada")
	fmt.Printf("     %-26s %-8s %-11s %-11s %-12s %-9s %s\n", "candidato", "R6", "operador", "conteo", "correlac.", "eco", "muere en")
	fmt.Printf("     %-26s %-8s %-11s %-11s %-12s %-9s %s\n", "inversión aritmética (F350)", "limpio", "NO", "sí", "sí", "tautológico", "no es operador")
	fmt.Printf("     %-26s %-8s %-11s %-11s %-12s %-9s %s\n", "grafo/flujo que concatena", "limpio", "sí", "—", "—", "IMPOSIBLE", "§2: longitudes prohibidas")
	fmt.Printf("     %-26s %-8s %-11s %-11s %-12s %-9s %s\n", "escaleras disjuntas", "limpio", "sí", "NO", "NO", "sí", "conteo y correlaciones")
	fmt.Printf("     %-26s %-8s %-11s %-11s %-12s %-9s %s\n", "Berry-Keating 2011", "limpio", "sí", "sí", "?", "NO", "sin aritmética")
	fmt.Println("     ⟹ el hueco es exactamente uno: un operador con la inestabilidad de xp cuyas órbitas")
	fmt.Println("       cerradas sean los primos SIN poder concatenarse. Eso es lo que la Fase IV debe buscar.")

	dibujar(Res{
		Verdad: verdad, Niveles: niv, Periodos: per, T0: t0, T1: t1,
		Forzadas: forzadas, Permitidas: permitidas, Prohibidas: prohibidas, TopeN: N,
		DensA: dens, S2a: sigma2(sp, 10, 600), S2v: sigma2(spV, 10, 600),
		EcoA: murcielago(niv, per, topeT), MinA: mini(sp), Primos: len(ps),
	})
}

func mini(v []float64) float64 {
	m := v[0]
	for _, x := range v {
		if x < m {
			m = x
		}
	}
	return m
}

// Res carries the measured numbers to the plate.
type Res struct {
	Verdad, Niveles, Periodos        []float64
	T0, T1                           float64
	Forzadas, Permitidas, Prohibidas int
	TopeN                            int
	DensA, S2a, S2v, EcoA, MinA      float64
	Primos                           int
}
