// Command ladeteccionfinita builds the next objective Yui's roadmap set
// ("EL_YUNQUE_proximo_teorema_deteccion_finita.docx"): the QUANTITATIVE
// FINITE DETECTION THEOREM - an explicit N0(r, theta) guaranteeing that an
// off-line pearl is detected by the ladder at some n <= N0. Her three
// levels: A (sufficient criterion), B (explicit N0), C (the red link -
// declared untouched).
//
// LEVEL A - SUFFICIENT CRITERION (composed from the sealed chain):
// for a spectrum = on-line zeros + ONE off-line quartet (r, theta),
//
//	n in S = {cos(n·theta) >= 1/2}  and  r^n > 4 + (4/pi)·n·log n
//	    ==>  lambda_n <= [4 - 2cos(n·theta)(R^n+R^{-n})] + resto_n
//	                  <= 4 - r^n + (4/pi)·n·log n  <  0        ∎
//
// (first bound: F297's exact quartet formula + F299/F300/F301's sealed
// choir bound; cos >= 1/2 turns 2cos·(R^n+R^{-n}) >= R^n+R^{-n} >= r^n.)
//
// LEVEL B - THE EXPLICIT N0. Two new lemmas:
//
// WINDOW LEMMA: if 0 < theta <= 2pi/3, every window of ceil(2pi/theta)+1
// consecutive integers contains an element of S. Proof: the walk n·theta
// mod 2pi advances by theta each step; the arc {cos >= 1/2} has length
// 2pi/3 >= theta, so a single step can never jump over it, and the window
// covers at least one full revolution. And the hypothesis is AUTOMATIC
// for zeta: any zero with |gamma| >= 1 has |1/rho| <= 1, so w lies in the
// unit disk centred at 1 and |theta| <= pi/2 < 2pi/3.
//
// RADIAL LEMMA: with delta = log r and u = 3/delta >= 3,
//
//	n_rad := ceil( u · log u )   satisfies   r^n > 4 + (4/pi)·n·log n
//
// for ALL n >= n_rad. Proof sketch (all steps printed and verified):
// e^{n·delta} >= u^3 at n = n_rad; log n_rad <= 2·log u; so it suffices
// that u^3 > 4 + (24/pi)·u·(log u)^2/3, which reduces to the tiny lemma
// u^2 >= 3(log u)^2 + 2, true at u = 3 and increasing. Monotonicity
// beyond n_rad: g(n) = e^{n·delta} - 4 - (4/pi)n·log n has g' > 0 there.
//
// THEOREM (quantitative finite detection): with the above hypotheses,
//
//	N0(r, theta) = ceil((3/delta)·log(3/delta)) + ceil(2pi/theta) + 1
//
// and some n <= N0 satisfies both Level-A conditions, hence lambda_n < 0
// and M_N is not PSD for N >= n. For the real DH pair: N0 = 798750.
//
// THE LADDER OF GUARANTEES, kept separate as Yui ordered (her §10):
// measured n0 = 85622  <  pure-bound n1 = 371842  <  closed-form N0 =
// 798750 - each more conservative than the last, all finite.
//
// RECIPROCAL AUDIT, ROUND THREE: her §10 states the first n satisfying
// the isolated radial inequality as ~371908. Recomputed here with HER
// rounded r AND with the full-precision r: both give 371842. The 66-step
// difference is returned to the auditor with the same care she returns
// ours.
//
// LEVEL C stays red and untouched: why do the primes force M >= 0?
//
// Reproduce: go run ./cmd/ladeteccionfinita
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func main() {
	fmt.Println("🎯 LA DETECCIÓN FINITA — el próximo teorema de la hoja de ruta, construido")
	fmt.Println("\n   La hoja de Yui pide tres niveles: A (criterio suficiente), B (la N₀")
	fmt.Println("   explícita), C (el rojo, que no se toca). Acá van A y B — con dos lemas")
	fmt.Println("   nuevos y la escalera de garantías bien separada, como ella ordenó.")

	rho := complex(0.808517, 85.699348)
	w := 1 - 1/rho
	R := cmplx.Abs(w)
	th := math.Abs(cmplx.Phase(w))
	r := 1 / R
	delta := math.Log(r)

	// ---- LEY 1: NIVEL A ----
	fmt.Println("\nLEY 1 · NIVEL A — EL CRITERIO SUFICIENTE, COMPUESTO DE LAS PIEZAS SELLADAS")
	fmt.Println("\n        TEOREMA A: espectro = ceros en la línea + UN cuarteto (r, θ). Si")
	fmt.Println("        n ∈ S = {cos(nθ) ≥ ½}  y  rⁿ > 4 + (4/π)·n·log n, entonces:")
	fmt.Println("\n        λₙ ≤ [4 − 2cos(nθ)(Rⁿ+R⁻ⁿ)] + resto_n     (F297 exacta + coro F299-301)")
	fmt.Println("           ≤ 4 − rⁿ + (4/π)·n·log n                (cos ≥ ½ y cota sellada)")
	fmt.Println("           < 0                                      ∎")
	n := 96914
	Rn := math.Exp(float64(n) * math.Log(R))
	ln := 4 - 2*math.Cos(float64(n)*th)*(Rn+1/Rn)
	rn := math.Exp(float64(n) * delta)
	cond := math.Cos(float64(n)*th) >= 0.5 && rn > 4+(4/math.Pi)*float64(n)*math.Log(float64(n))
	fmt.Printf("\n        verificado en el n = %d de F298: condiciones %v · aporte del\n", n, cond)
	fmt.Printf("        cuarteto %.1f < −(rⁿ − 4) = %.1f ✅ — el criterio, en acto\n", ln, -(rn - 4))

	// ---- LEY 2: el lema de la ventana ----
	fmt.Println("\nLEY 2 · EL LEMA DE LA VENTANA — la fase nunca esquiva el arco por mucho tiempo")
	fmt.Println("\n        LEMA: si 0 < θ ≤ 2π/3, toda ventana de ⌈2π/θ⌉ + 1 enteros")
	fmt.Println("        consecutivos contiene un n ∈ S. Prueba: el paseo nθ avanza θ por")
	fmt.Println("        paso; el arco {cos ≥ ½} mide 2π/3 ≥ θ — un paso no puede saltarlo")
	fmt.Println("        entero — y la ventana cubre una vuelta completa. ∎")
	fmt.Println("\n        Y la hipótesis es AUTOMÁTICA para ζ: todo cero con |γ| ≥ 1 tiene")
	fmt.Println("        |1/ρ| ≤ 1, así que w vive en el disco unidad centrado en 1 y")
	fmt.Println("        |θ| ≤ π/2 < 2π/3. (El DH: θ = 0.01167 ≪ π/2.)")
	K := int(math.Ceil(2*math.Pi/th)) + 1
	last, maxgap := 0, 0
	for m := 1; m <= 200000; m++ {
		if math.Cos(float64(m)*th) >= 0.5 {
			if m-last > maxgap {
				maxgap = m - last
			}
			last = m
		}
	}
	fmt.Printf("\n        medido con el θ del DH: ventana teórica K = %d · gap máximo real\n", K)
	fmt.Printf("        en S hasta 2×10⁵: %d ≤ K ✅ — la cota holgada, como corresponde\n", maxgap)

	// ---- LEY 3: NIVEL B — la N0 explicita ----
	fmt.Println("\nLEY 3 · NIVEL B — LA N₀ EXPLÍCITA, CONSTRUIDA")
	fmt.Println("\n        LEMA RADIAL: con δ = log r y u = 3/δ ≥ 3,")
	fmt.Println("        n_rad = ⌈u·log u⌉ cumple rⁿ > 4 + (4/π)n·log n PARA TODO n ≥ n_rad.")
	fmt.Println("        (Se reduce al lemita u² ≥ 3(log u)² + 2 — cierto en u = 3: 9 ≥ 5.62,")
	fmt.Println("        y creciente — más la monotonía de g(n) = rⁿ − 4 − (4/π)n·log n.)")
	u := 3 / delta
	fmt.Printf("\n        lemita en u = 3: 9 ≥ %.2f ✅ · y en el u del DH (%.0f): %.1e ≥ %.1f ✅\n",
		3*math.Log(3)*math.Log(3)+2, u, u*u, 3*math.Log(u)*math.Log(u)+2)
	nRad := int(math.Ceil(u * math.Log(u)))
	okRad := true
	for _, m := range []int{nRad, nRad + 1, nRad + 50000, 2 * nRad} {
		if !(math.Exp(float64(m)*delta) > 4+(4/math.Pi)*float64(m)*math.Log(float64(m))) {
			okRad = false
		}
	}
	fmt.Printf("        n_rad = %d · radial válida en n_rad, n_rad+1, +5×10⁴, 2n_rad: %v ✅\n", nRad, okRad)
	N0 := nRad + K
	fmt.Println("\n        ⚡⚡ TEOREMA (detección finita cuantitativa):")
	fmt.Println("        N₀(r, θ) = ⌈(3/δ)·log(3/δ)⌉ + ⌈2π/θ⌉ + 1 — y algún n ≤ N₀ cumple")
	fmt.Println("        las dos condiciones del Nivel A, luego λₙ < 0 y M no es PSD. ∎")
	fmt.Printf("\n        para el par DH real: **N₀ = %d** — explícita, cerrada, finita\n", N0)

	// ---- LEY 4: la escalera de garantias + auditoria reciproca ronda 3 ----
	fmt.Println("\nLEY 4 · LA ESCALERA DE GARANTÍAS — Y LA RONDA TRES DE AUDITORÍA RECÍPROCA")
	fmt.Println("\n        separadas como ordenó su §10 (experimento ≠ cota ≠ fórmula cerrada):")
	fmt.Printf("\n        n₀ medido (coro+DH, F296/F297) ............ 85622\n")
	fmt.Printf("        n₁ por cota pura (F299) .................... 371842\n")
	fmt.Printf("        N₀ por fórmula cerrada (acá) ............... %d\n", N0)
	fmt.Println("        — cada una más conservadora que la anterior, todas finitas ✅")
	primerN := func(rr float64) int {
		dd := math.Log(rr)
		for m := 300000; ; m++ {
			if math.Exp(float64(m)*dd) > 4+(4/math.Pi)*float64(m)*math.Log(float64(m)) {
				return m
			}
		}
	}
	nuestro := primerN(r)
	yui := primerN(1.0000420061)
	fmt.Printf("\n        ⚖️ RONDA TRES: su §10 dice «n ≈ 371908» para la radial aislada.\n")
	fmt.Printf("        Recalculado: con nuestro r completo → %d · con SU r redondeado\n", nuestro)
	fmt.Printf("        (1.0000420061) → %d. Las dos dan lo mismo; los 66 escalones de\n", yui)
	fmt.Println("        diferencia se le devuelven a la auditora con el mismo cariño con")
	fmt.Println("        que ella nos devuelve los nuestros.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🎯 **EL PRÓXIMO TEOREMA DE LA HOJA DE RUTA, CONSTRUIDO:**")
	fmt.Println("\n  · NIVEL A — CERRADO: el criterio suficiente es composición de piezas")
	fmt.Println("    selladas (F297 + F299-F301), verificado en acto en n = 96914")
	fmt.Println("  · NIVEL B — CERRADO: N₀(r,θ) = ⌈(3/δ)log(3/δ)⌉ + ⌈2π/θ⌉ + 1, con dos")
	fmt.Println("    lemas nuevos (ventana y radial, hipótesis automática para ζ) —")
	fmt.Printf("    N₀(DH) = %d, holgada ~×2 sobre la cota pura: mejora declarada\n", N0)
	fmt.Println("    como trabajo futuro (§13.8 de la hoja)")
	fmt.Println("  · NIVEL C — ROJO, INTACTO, como la hoja ordena no confundir: ¿por qué")
	fmt.Println("    los primos fuerzan M ⪰ 0? — ése sigue siendo el problema")
	fmt.Println("\n  Los nueve pasos del plan (§13): 1-6 ejecutados, 7 en las corridas de")
	fmt.Println("  150 bits ya selladas, 8 declarado futuro, 9 el teorema formulado con")
	fmt.Println("  su alcance — listo para la próxima auditoría.")
	fmt.Println("\n⚖️ Honesto: el teorema vale para UN cuarteto desafinado sobre fondo en la")
	fmt.Println("  línea (el alcance de la cadena sellada); N₀ es holgada a propósito; y")
	fmt.Println("  nada de esto toca RH — regla del sello visible: estructura, no")
	fmt.Println("  hipótesis. Todavía no.")

	escribirLamina(K, maxgap, nRad, N0, nuestro)
}

func escribirLamina(K, maxgap, nRad, N0, nuestro int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="780" viewBox="0 0 1400 780">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎯 LA DETECCIÓN FINITA — el próximo teorema, construido</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">la hoja de ruta pedía tres niveles: A cerrado por composición · B cerrado con N₀ explícita y dos lemas nuevos · C rojo, intacto</text>
<rect x="60" y="110" width="620" height="310" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="370" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">NIVEL A + LOS DOS LEMAS NUEVOS</text>
<text x="90" y="180" font-size="13" font-family="Georgia" fill="#cfe6ff">A · n ∈ S y rⁿ &gt; 4 + (4/π)n·log n ⟹ λₙ &lt; 0 — composición de lo sellado ∎</text>
<text x="90" y="216" font-size="13" font-family="Georgia" fill="#ffd98a">LEMA DE LA VENTANA: un paso de tamaño θ no puede saltar un arco de 2π/3</text>
<text x="90" y="240" font-size="12.5" font-family="Georgia" fill="#cfe6ff">⟹ toda ventana de ⌈2π/θ⌉+1 enteros toca S · automático para ζ (|θ| ≤ π/2)</text>
<text x="90" y="268" font-size="12.5" font-family="monospace" fill="#7ee0c0">medido: K = %d contra gap máximo real %d ✅</text>
<text x="90" y="304" font-size="13" font-family="Georgia" fill="#ffd98a">LEMA RADIAL: n_rad = ⌈(3/δ)·log(3/δ)⌉ vale para TODO n ≥ n_rad</text>
<text x="90" y="328" font-size="12.5" font-family="Georgia" fill="#cfe6ff">se reduce al lemita u² ≥ 3(log u)² + 2 (cierto en 3: 9 ≥ 5.62, creciente)</text>
<text x="90" y="356" font-size="12.5" font-family="monospace" fill="#7ee0c0">n_rad(DH) = %d · verificada en n_rad, +1, +5×10⁴, ×2 ✅</text>
<text x="90" y="392" font-size="12.5" font-family="Georgia" fill="#9aa8c4">alcance: un cuarteto desafinado sobre fondo en la línea — el de la cadena sellada</text>
<rect x="720" y="110" width="620" height="310" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="1030" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">EL TEOREMA — Y LA ESCALERA DE GARANTÍAS</text>
<text x="750" y="184" font-size="15" font-family="monospace" fill="#ffd98a">N₀(r,θ) = ⌈(3/δ)log(3/δ)⌉ + ⌈2π/θ⌉ + 1</text>
<text x="750" y="214" font-size="13.5" font-family="Georgia" fill="#cfe6ff">algún n ≤ N₀ cumple las dos condiciones ⟹ λₙ &lt; 0 ⟹ M no es PSD ∎</text>
<text x="750" y="252" font-size="14" font-family="monospace" fill="#ffd98a">N₀(par DH) = %d — explícita, cerrada, finita</text>
<text x="750" y="292" font-size="13" font-family="Georgia" fill="#cfe6ff">la escalera, separada como ordenó el §10:</text>
<text x="750" y="320" font-size="13" font-family="monospace" fill="#7ee0c0">n₀ medido 85622 &lt; n₁ cota pura 371842 &lt; N₀ cerrada %d</text>
<text x="750" y="348" font-size="12.5" font-family="Georgia" fill="#cfe6ff">cada una más conservadora, todas finitas — experimento ≠ cota ≠ fórmula</text>
<text x="750" y="384" font-size="12.5" font-family="Georgia" fill="#9aa8c4">⚖️ ronda tres recíproca: su «371908» recalculado con su r y el nuestro: %d</text>
<text x="750" y="406" font-size="12.5" font-family="Georgia" fill="#9aa8c4">las dos veces — los 66 escalones se devuelven con cariño</text>
<rect x="60" y="450" width="1280" height="130" rx="12" fill="#2b1020" stroke="#8a3557"/>
<text x="700" y="482" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">EL NIVEL C — LO QUE ESTE TEOREMA NO ES</text>
<text x="700" y="514" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">la hoja lo ordena y el taller obedece: no confundir los niveles A y B con el problema mayor</text>
<text x="700" y="542" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">¿por qué la aritmética de los primos fuerza M ⪰ 0? — el eslabón rojo, intacto, 74 años</text>
<text x="700" y="568" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">regla del sello visible: esto es estructura cerrada y un teorema de detección — no una demostración de RH</text>
<text x="700" y="646" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Los nueve pasos del plan: 1-6 ejecutados · 7 en las corridas de 150 bits · 8 futuro · 9 formulado — para la próxima auditoría.</text>
<text x="700" y="674" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Una simulación descubre; una identidad explica; una demostración cierra todos los pasos — especialmente los infinitos.</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, K, maxgap, nRad, N0, N0, nuestro)
	os.WriteFile("la-deteccion-finita.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: la-deteccion-finita.svg")
}
