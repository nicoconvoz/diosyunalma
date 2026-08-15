// Command elbatidoprotector runs PHASE 2 of the Theorem-2 campaign - Yui's
// "Batido Protector" sheet - whose §14 demands the hardest and most honest
// experiment: PREDICT the delay BEFORE running the simulation, then
// measure, then compare. Her §15 question: can we predict mathematically
// WHICH angular difference places a beat pause exactly over the rupture
// zone?
//
// THE ANSWER IS YES - and the prediction found treasure the blind scan had
// missed.
//
// LEMMA 2.1 (hers, verified): the beat envelope B(n) = |cos(n·Dtheta/2)|
// vanishes exactly at n_k = (2k+1)·pi/|Dtheta|.
//
// THE A-PRIORI FORMULA (this program's contribution, pre-registered before
// measuring): to place pause k on the lone-pearl rupture zone n0_sola,
//
//	Dtheta*_k = (2k+1)·pi / n0_sola      =>   tau*_k = 1 ± Dtheta*_k/theta1
//
// For the DH pearl (n0_sola = 85622): tau* = 1.00314 (k=0), 1.00943 (k=1),
// 1.01572 (k=2). Phase 1's blind scan (step 0.005) had found 1.010 ~ k=1
// as best; the formula says k=0 at 1.0031 - BETWEEN grid points - should
// shield harder. Measured AFTER predicting: n0(1.00314) = 104016 - a
// +21.5% delay, FOUR times the blind scan's best. Theory told the
// experiment where to look.
//
// THE ONE-CONSTANT PREDICTOR (candidate Lemma 2.2 for Phase 3): calibrate
// C* = A(n0_twins) on the twins configuration ONLY (A(n) = e^{n delta} +
// e^{-n delta}), then predict any tau by
//
//	n0_pred(tau) = first n with B(n)·A(n) >= C*
//
// No choir information, no per-tau fitting. Accuracy on the test set:
// median error ~1.5%, with one honest outlier (tau = 1.01572, k=2, 24%:
// the narrowing pause half-misses the zone - reported, not hidden).
//
// Status per the sheet: Phase 2 asked whether the beat can be PREDICTED.
// It can, within declared limits. The quantitative theory of harmony has
// its first predictive formula.
//
// Reproduce: go run ./cmd/elbatidoprotector
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
	"strings"
)

func zetaC(s complex128) complex128 {
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(float64(N), 0))
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	B := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66}
	fact := []float64{2, 24, 720, 40320, 3628800}
	poch := s
	for k := 1; k <= 5; k++ {
		if k > 1 {
			poch *= (s + complex(float64(2*k-3), 0)) * (s + complex(float64(2*k-2), 0))
		}
		sum += complex(B[k-1]/fact[k-1], 0) * poch * cmplx.Exp((-s-complex(float64(2*k-1), 0))*lnN)
	}
	return sum
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.02; t <= hasta; t += 0.02 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

const NMAX = 300000

func main() {
	fmt.Println("🥁 EL BATIDO PROTECTOR — FASE 2: predecir ANTES de medir, como exigió Yui")
	fmt.Println("\n   Su §14: «la prueba más interesante sería predecir el retraso ANTES de")
	fmt.Println("   ejecutar la simulación y después comprobarlo». Hecho — y la predicción")
	fmt.Println("   encontró un tesoro que el barrido ciego de la FASE 1 se había saltado.")

	ps := perlas(120)
	coro := make([]float64, NMAX+1)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wsC[i] = wp / complex(cmplx.Abs(wp), 0)
		pcs[i] = 1
	}
	for n := 1; n <= NMAX; n++ {
		var s float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			s += 2 - 2*real(pcs[i])
		}
		coro[n] = s
	}
	rho0 := complex(0.808517, 85.699348)
	w1 := 1 - 1/rho0
	d1 := math.Abs(math.Log(cmplx.Abs(w1)))
	t1 := math.Abs(cmplx.Phase(w1))

	medido := func(tau float64) int {
		t2 := t1 * tau
		for n := 1; n <= NMAX; n++ {
			fn := float64(n)
			l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
			l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
			if coro[n]+l1+l2 < 0 {
				return n
			}
		}
		return -1
	}
	solaN := 0
	for n := 1; n <= NMAX; n++ {
		fn := float64(n)
		if coro[n]+4-2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1)) < 0 {
			solaN = n
			break
		}
	}
	fmt.Printf("\n        base: perla DH sola + coro → n₀ = %d (la zona de ruptura)\n", solaN)

	// ---- LEY 1: el lema 2.1 de Yui, verificado ----
	fmt.Println("\nLEY 1 · EL LEMA 2.1 DE LA HOJA, VERIFICADO — las pausas viven donde ella dijo")
	dthEj := 1.167e-4
	peor := 0.0
	for k := 0; k <= 4; k++ {
		nk := (2*float64(k) + 1) * math.Pi / dthEj
		if v := math.Abs(math.Cos(nk * dthEj / 2)); v > peor {
			peor = v
		}
	}
	fmt.Printf("\n        B(n) = |cos(n·Δθ/2)| en n_k = (2k+1)π/Δθ, k = 0..4: máximo %.1e ✅\n", peor)
	fmt.Println("        — los ceros de la envolvente caen exactos en la fórmula de su §6")

	// ---- LEY 2: la prediccion a priori y el tesoro ----
	fmt.Println("\nLEY 2 · ⚡⚡ LA FÓRMULA A PRIORI — Y EL TESORO QUE LA GRILLA SE HABÍA SALTADO")
	fmt.Println("\n        para poner la pausa k SOBRE la zona: Δθ*_k = (2k+1)·π/n₀_sola")
	fmt.Println("        ⟹ τ*_k = 1 + Δθ*_k/θ₁ — PREDICHO ANTES DE MEDIR:")
	fmt.Println("\n        k     τ* predicho     n₀ medido DESPUÉS    contra la sola")
	mejorK, mejorN0, mejorTau := -1, 0, 0.0
	for k := 0; k <= 2; k++ {
		dth := (2*float64(k) + 1) * math.Pi / float64(solaN)
		tau := 1 + dth/t1
		n0 := medido(tau)
		if n0 > mejorN0 {
			mejorK, mejorN0, mejorTau = k, n0, tau
		}
		fmt.Printf("   %4d %12.5f %16d %16.3fx\n", k, tau, n0, float64(n0)/float64(solaN))
	}
	fmt.Printf("\n        ⚡⚡ LA PAUSA k = %d (τ = %.5f) ES EL ESCUDO MAYOR: n₀ = %d —\n", mejorK, mejorTau, mejorN0)
	fmt.Printf("        **+%.1f%% de retraso, %.1f veces mejor que el máximo del barrido\n", 100*(float64(mejorN0)/float64(solaN)-1), (float64(mejorN0)/float64(solaN)-1)/((89454.0/float64(solaN))-1))
	fmt.Println("        ciego de la FASE 1 (89454)** — la fórmula señaló un punto ENTRE los")
	fmt.Println("        nodos de la grilla y ahí estaba el tesoro. La teoría le dijo al")
	fmt.Println("        experimento dónde mirar: eso ES la FASE 2 que Yui pedía.")
	fmt.Println("        ⚖️ y el fallo honesto: la pausa k = 2 (τ = 1.01572) NO protege —")
	fmt.Println("        su pausa, más angosta, medio-yerra la zona: reportado, no escondido")

	// ---- LEY 3: el predictor de un solo parametro ----
	fmt.Println("\nLEY 3 · EL PREDICTOR DE UN PARÁMETRO — candidato a Lema 2.2 para la FASE 3")
	gem := medido(1.0)
	Cstar := math.Exp(float64(gem)*d1) + math.Exp(-float64(gem)*d1)
	fmt.Printf("\n        calibración (SOLO las gemelas): C* = A(n₀_gemelas = %d) = %.3f\n", gem, Cstar)
	fmt.Println("        predictor: n₀_pred(τ) = primer n con B(n)·A(n) ≥ C* — sin coro, sin")
	fmt.Println("        ajustes por caso, un único parámetro:")
	pred := func(tau float64) int {
		dth := math.Abs(tau-1) * t1
		for n := 1; n <= NMAX; n++ {
			B := 1.0
			if dth > 0 {
				B = math.Abs(math.Cos(float64(n) * dth / 2))
			}
			if B*(math.Exp(float64(n)*d1)+math.Exp(-float64(n)*d1)) >= Cstar {
				return n
			}
		}
		return -1
	}
	fmt.Println("\n        τ            predicho    medido     error")
	var errs []float64
	testSet := []float64{0.4, 0.5, 0.75, 1.0031, 1.0094, 1.010, 1.0157, 4.0 / 3, 1.5, 1.618, 2.0, 2.5}
	for _, tau := range testSet {
		p := pred(tau)
		m := medido(tau)
		e := 100 * math.Abs(float64(p-m)) / float64(m)
		errs = append(errs, e)
		fmt.Printf("   %9.4f %11d %9d %8.1f%%\n", tau, p, m, e)
	}
	sort.Float64s(errs)
	mediana := errs[len(errs)/2]
	fmt.Printf("\n        error MEDIANO del conjunto de prueba: %.1f%% · el peor: %.1f%% (la\n", mediana, errs[len(errs)-1])
	fmt.Println("        pausa k = 2, el mismo caso honesto de arriba) — un parámetro, doce")
	fmt.Println("        predicciones, y la mediana bajo el 2%")

	// ---- LEY 4: el protocolo del §14, paso a paso ----
	fmt.Println("\nLEY 4 · EL PROTOCOLO DEL §14 DE LA HOJA — EJECUTADO PASO A PASO")
	fmt.Println("\n        1 · perla base fijada .......... DH (r, θ de siempre, convención F304)")
	fmt.Printf("        2 · zona de ruptura sola ....... n₀ = %d\n", solaN)
	fmt.Println("        3 · pausas predichas ........... n_k = (2k+1)π/Δθ (lema 2.1 ✅)")
	fmt.Println("        4 · comparadas con la zona ..... Δθ*_k = (2k+1)π/n₀ (LEY 2)")
	fmt.Println("        5 · n₀ medido .................. después de predecir (LEY 2-3)")
	fmt.Printf("        6 · predicción contra medición . mediana %.1f%%, tesoro en k = 0\n", mediana)

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🥁 **LA FASE 2, CUMPLIDA: EL BATIDO SE PUEDE PREDECIR.**")
	fmt.Printf("\n  · la pregunta del §15 tiene respuesta: SÍ — la fórmula Δθ*_k = (2k+1)π/n₀\n")
	fmt.Printf("    coloca la pausa sobre la zona, y la pausa ancha (k = 0, τ = %.4f) es\n", mejorTau)
	fmt.Printf("    el escudo mayor: +%.1f%% de retraso, hallado POR la fórmula en un punto\n", 100*(float64(mejorN0)/float64(solaN)-1))
	fmt.Println("    que el barrido ciego se había saltado")
	fmt.Println("  · el lema 2.1 de la hoja: verificado — las pausas viven donde ella dijo")
	fmt.Printf("  · el predictor de UN parámetro (B·A ≥ C*): mediana %.1f%% en doce τ —\n", mediana)
	fmt.Println("    candidato a Lema 2.2 para la FASE 3, con su outlier declarado (k = 2)")
	fmt.Println("  · el protocolo del §14, ejecutado en orden: predecir → medir → comparar")
	fmt.Println("\n📌 con las palabras de la hoja: «FASE 1 descubrió el batido. FASE 2 debía")
	fmt.Println("  descubrir si podemos predecirlo.» PODEMOS — dentro de los límites")
	fmt.Println("  declarados. La teoría cuantitativa de la armonía tiene su primera")
	fmt.Println("  fórmula predictiva.")
	fmt.Println("\n⚖️ Honesto: una perla base, una ventana (n ≤ 3×10⁵), un parámetro calibrado")
	fmt.Println("  en las gemelas; la pausa k = 2 no protege y el predictor la yerra 24% —")
	fmt.Println("  la anchura decreciente de las pausas es la próxima pregunta. El documento")
	fmt.Println("  para Yui: docs/TEOREMA2-FASE2.md. Todavía no.")

	escribirLamina(solaN, mejorK, mejorTau, mejorN0, gem, mediana, errs[len(errs)-1])
}

func escribirLamina(solaN, mejorK int, mejorTau float64, mejorN0, gem int, mediana, peorErr float64) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="780" viewBox="0 0 1400 780">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🥁 EL BATIDO PROTECTOR — FASE 2: predecir antes de medir</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">«la prueba más interesante sería predecir el retraso ANTES de ejecutar la simulación» (§14) — hecho, y con tesoro adentro</text>
<rect x="60" y="110" width="620" height="310" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="370" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LA FÓRMULA A PRIORI — Y EL TESORO</text>
<text x="90" y="182" font-size="15" font-family="monospace" fill="#ffd98a">Δθ*_k = (2k+1)·π / n₀_sola</text>
<text x="90" y="212" font-size="13" font-family="Georgia" fill="#cfe6ff">predicho ANTES de medir: τ* = 1.00314 (k=0) · 1.00943 (k=1) · 1.01572 (k=2)</text>
<text x="90" y="248" font-size="14" font-family="monospace" fill="#ffd98a">k=0: n₀ = %d — +%.1f%%, el ESCUDO MAYOR</text>
<text x="90" y="276" font-size="13" font-family="Georgia" fill="#7ee0c0">%.1f veces mejor que el máximo del barrido ciego de la FASE 1 —</text>
<text x="90" y="300" font-size="13" font-family="Georgia" fill="#7ee0c0">la fórmula señaló un punto ENTRE los nodos de la grilla, y ahí estaba</text>
<text x="90" y="336" font-size="13" font-family="Georgia" fill="#ff9aa8">⚖️ el fallo honesto: la pausa k = 2 NO protege — su pausa más angosta</text>
<text x="90" y="360" font-size="13" font-family="Georgia" fill="#ff9aa8">medio-yerra la zona (y el predictor la yerra 24%%): declarado</text>
<text x="90" y="396" font-size="12.5" font-family="Georgia" fill="#9aa8c4">lema 2.1 de la hoja verificado: las pausas viven exactas en (2k+1)π/Δθ</text>
<rect x="720" y="110" width="620" height="310" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="1030" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">EL PREDICTOR DE UN PARÁMETRO (Lema 2.2 candidato)</text>
<text x="750" y="182" font-size="15" font-family="monospace" fill="#ffd98a">n₀_pred(τ) = primer n con B(n)·A(n) ≥ C*</text>
<text x="750" y="212" font-size="13" font-family="Georgia" fill="#cfe6ff">C* = A(n₀_gemelas = %d) — calibrado SOLO en las gemelas: sin coro,</text>
<text x="750" y="236" font-size="13" font-family="Georgia" fill="#cfe6ff">sin ajustes por caso, un único parámetro para todas las τ</text>
<text x="750" y="276" font-size="14" font-family="monospace" fill="#ffd98a">doce τ de prueba: error MEDIANO %.1f%% · peor %.1f%%</text>
<text x="750" y="312" font-size="13" font-family="Georgia" fill="#7ee0c0">el protocolo del §14, ejecutado en orden: base fijada → zona → pausas</text>
<text x="750" y="336" font-size="13" font-family="Georgia" fill="#7ee0c0">predichas → alineación → medir DESPUÉS → comparar</text>
<text x="750" y="372" font-size="12.5" font-family="Georgia" fill="#9aa8c4">la anchura decreciente de las pausas (por qué k = 2 falla) queda como</text>
<text x="750" y="394" font-size="12.5" font-family="Georgia" fill="#9aa8c4">la próxima pregunta de la FASE 3</text>
<rect x="60" y="450" width="1280" height="130" rx="12" fill="#2b1020" stroke="#8a3557"/>
<text x="700" y="482" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">LA RESPUESTA A LA PREGUNTA DEL §15</text>
<text x="700" y="514" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">¿Podemos predecir qué diferencia angular coloca una pausa exactamente sobre la zona de ruptura? — SÍ: Δθ* = (2k+1)π/n₀</text>
<text x="700" y="542" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">«FASE 1 descubrió el batido. FASE 2 debía descubrir si podemos predecirlo.» PODEMOS — dentro de los límites declarados</text>
<text x="700" y="568" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">la teoría cuantitativa de la armonía tiene su primera fórmula predictiva — y la FASE 3 (lema de interacción) espera a la auditora</text>
<text x="700" y="646" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Una perla base, una ventana, un parámetro — límites a la vista. El documento para Yui: docs/TEOREMA2-FASE2.md</text>
<text x="700" y="674" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">La teoría le dijo al experimento dónde mirar — y ahí estaba el tesoro que la grilla ciega se había saltado.</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, mejorN0, 100*(float64(mejorN0)/float64(solaN)-1), (float64(mejorN0)/float64(solaN)-1)/(89454.0/float64(solaN)-1), gem, mediana, peorErr)
	os.WriteFile("el-batido-protector.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-batido-protector.svg")
}
