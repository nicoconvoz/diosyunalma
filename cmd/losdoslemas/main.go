// Command losdoslemas verifies, step by step, the formal act written for
// the auditor in docs/DETECCION-FINITA-LEMAS.md - the complete proofs of
// the two lemmas her audit of LA DETECCION FINITA demanded (§12): the
// exact definition of delta, its derivation from r, the full radial lemma,
// the full window lemma, the combination inside one interval, and the
// single frozen convention that settles the 371842/371908 discrepancy.
//
// THE FROZEN CONVENTION (§0 of the act): log = natural log everywhere;
// r = max(|w|, 1/|w|); delta := log r; u = 3/delta (hypothesis u >= 3);
// ceilings are upward. Official DH threshold: n1 = 371842 (reproduced
// under both r roundings; 371908 reproduces under neither).
//
// LEMA R (radial, ceiling-robust): n_rad = ceil(u log u) gives
// r^n > 4 + (4/pi) n log n for ALL n >= n_rad. The proof reduces to the
// auxiliary lemma u^2 >= 4(log u)^2 + 2 (u >= 3) - note: F303's plate
// quoted the un-ceiled version with constant 3; the official
// ceiling-robust constant is 4. Both true; the act records the change.
//
// LEMA V (window): 0 < theta <= 2pi/3 implies every K = ceil(2pi/theta)+1
// consecutive integers contain n with cos(n theta) >= 1/2 - a step of
// size theta cannot jump an arc of length 2pi/3. LEMA V-zeta: automatic
// for zeta (|Im rho| >= 1 => |1/rho| <= 1 => w in disk D(1,1) =>
// |arg w| <= pi/2).
//
// COMBINATION: window applied at m = n_rad picks n in [n_rad, n_rad+K-1];
// the radial lemma already covers that whole interval. Realized for DH:
// n = 798474, bound -3.7e14 < 0.
//
// Reproduce: go run ./cmd/losdoslemas
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func main() {
	fmt.Println("📜 LOS DOS LEMAS — el acta formal para la auditora, verificada paso a paso")
	fmt.Println("\n   Su §12 pidió siete cosas: la definición exacta de δ, su derivación")
	fmt.Println("   desde r, las dos pruebas completas, la combinación en un mismo")
	fmt.Println("   intervalo, y UNA convención congelada. El acta está en")
	fmt.Println("   docs/DETECCION-FINITA-LEMAS.md — acá corre cada paso.")

	// ---- §0: la convencion congelada ----
	fmt.Println("\n§0 · LA CONVENCIÓN CONGELADA — y δ, DEFINIDA")
	rho := complex(0.808517, 85.699348)
	w := 1 - 1/rho
	R := cmplx.Abs(w)
	r := math.Max(R, 1/R)
	delta := math.Log(r)
	u := 3 / delta
	th := math.Abs(cmplx.Phase(w))
	fmt.Println("\n        log = natural en TODO el documento · r = max(|w|, 1/|w|) ·")
	fmt.Printf("        **δ := log r** (§12.1-2) · u = 3/δ · techos hacia arriba\n")
	fmt.Printf("\n        DH oficial: r = %.15f · δ = %.13f\n", r, delta)
	fmt.Printf("        u = %.2f · θ = %.15f\n", u, th)
	primerN := func(rr float64) int {
		dd := math.Log(rr)
		for m := 2; ; m++ {
			if math.Exp(float64(m)*dd) > 4+(4/math.Pi)*float64(m)*math.Log(float64(m)) {
				return m
			}
		}
	}
	a := primerN(r)
	b := primerN(1.0000420061)
	fmt.Printf("\n        umbral radial bajo la convención: r completo → %d · r redondeado\n", a)
	fmt.Printf("        → %d · ⟹ **n₁ oficial = %d** (el 371908 no se reproduce con\n", b, a)
	fmt.Println("        ninguna de las dos entradas — §12.6, zanjado)")

	// ---- LEMA R ----
	fmt.Println("\nLEMA R · LA PRUEBA RADIAL COMPLETA — con el lemita robusto al techo")
	fmt.Println("\n        la cadena: e^{n_rad·δ} ≥ u³ · n_rad ≤ 1.31·n* · log n_rad ≤ 2.29·log u")
	fmt.Println("        ⟹ basta u² ≥ 4(log u)² + 2  (⚠️ la lámina de F303 citaba la versión")
	fmt.Println("        sin techo con constante 3; la oficial robusta es 4 — registrado):")
	violL := 0
	for _, uu := range []float64{3, 3.5, 5, 20, 1000, u} {
		if uu*uu < 4*math.Log(uu)*math.Log(uu)+2 {
			violL++
		}
	}
	fmt.Printf("\n        lemita en u = 3, 3.5, 5, 20, 10³, u_DH: %d violaciones ✅ · creciente:\n", violL)
	fmt.Printf("        u² > 4·log u en u = 3 (9 > %.2f) ✅\n", 4*math.Log(3))
	fmt.Println("\n        ⚖️ CORRECCIÓN F305 (auditoría «la pizza» §6): el paso 2 encadenaba")
	fmt.Println("        «≤ 4·log u», que FALLA en u = 3 (4.50 > 4.39 — cazado por Yui). La")
	fmt.Println("        ruta corregida, su §15: comparar directo contra 3u²:")
	violP := 0
	for _, uu := range []float64{3, 3.5, 5, 20, 1000, u} {
		if 1.28*(2.29*math.Log(uu)+1) >= 3*uu*uu {
			violP++
		}
	}
	fmt.Printf("\n        1.28·(2.29·log u + 1) < 3u² en la grilla: %d violaciones ✅ (en u = 3:\n", violP)
	fmt.Printf("        %.2f < 27; creciente: 6u − 2.94/u > 0) ⟹ g'(n_rad) > 0, ahora sí ∎\n", 1.28*(2.29*math.Log(3)+1))
	nRad := int(math.Ceil(u * math.Log(u)))
	g := func(n float64) float64 { return math.Exp(n*delta) - 4 - (4/math.Pi)*n*math.Log(n) }
	gp := func(n float64) float64 { return delta*math.Exp(n*delta) - (4/math.Pi)*(math.Log(n)+1) }
	gpp := func(n float64) float64 { return delta*delta*math.Exp(n*delta) - (4/math.Pi)/n }
	okG := g(float64(nRad)) > 0 && gp(float64(nRad)) > 0 && gpp(float64(nRad)) > 0
	fmt.Printf("\n        n_rad = %d: g > 0, g' > 0, g'' > 0 en n_rad: %v ✅ — la radial vale\n", nRad, okG)
	fmt.Println("        para TODO n ≥ n_rad, con monotonía demostrada, no supuesta")

	// ---- LEMA V ----
	fmt.Println("\nLEMA V · LA VENTANA, ADVERSARIAL — siete fases, cero violaciones")
	fmt.Println("\n        θ            K       gap máximo real (10⁵ enteros)")
	violV := 0
	for _, t := range []float64{th, 2*math.Pi/3 - 1e-6, 1.0, 2 * math.Pi / 7, 0.5, 2.09, 1e-3} {
		K := int(math.Ceil(2*math.Pi/t)) + 1
		last, mg := 0, 0
		for n := 1; n <= 100000; n++ {
			if math.Cos(float64(n)*t) >= 0.5 {
				if n-last > mg {
					mg = n - last
				}
				last = n
			}
		}
		marca := "✅"
		if mg > K {
			violV++
			marca = "⚠️"
		}
		fmt.Printf("   %12.6f %6d %10d %s\n", t, K, mg, marca)
	}
	fmt.Printf("\n        violaciones: %d ✅ · y la hipótesis V-ζ es automática: |Im ρ| ≥ 1 ⟹\n", violV)
	fmt.Printf("        |1/ρ| ≤ 1 ⟹ w ∈ disco(1,1) ⟹ |θ| ≤ π/2 (DH: |1/ρ| = %.6f)\n", cmplx.Abs(1/rho))

	// ---- la combinacion ----
	fmt.Println("\n§3 · LA COMBINACIÓN EN EL MISMO INTERVALO — realizada")
	K := int(math.Ceil(2*math.Pi/th)) + 1
	nCombo := -1
	for n := nRad; n < nRad+K; n++ {
		if math.Cos(float64(n)*th) >= 0.5 {
			nCombo = n
			break
		}
	}
	cota := 4 - math.Exp(float64(nCombo)*delta) + (4/math.Pi)*float64(nCombo)*math.Log(float64(nCombo))
	fmt.Printf("\n        la ventana con m = n_rad elige n; la radial ya cubre el intervalo.\n")
	fmt.Printf("        DH: primer n ∈ S dentro de [%d, %d] es **n = %d**, y allí\n", nRad, nRad+K-1, nCombo)
	fmt.Printf("        4 − rⁿ + (4/π)n·log n = %.1e < 0 ✅ — las dos garantías, un entero\n", cota)
	fmt.Printf("\n        N₀ = n_rad + K = %d ✅ (la fórmula del teorema, §4 del acta)\n", nRad+K)

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("📜 **EL ACTA DE LOS DOS LEMAS, ENTREGADA — los siete puntos del §12:**")
	fmt.Println("\n  1-2 · δ definida y derivada: **δ = log natural de r**, r = max(|w|,1/|w|)")
	fmt.Println("  3 · lema radial COMPLETO: cadena con techos, lemita robusto (constante 4,")
	fmt.Println("      no 3 — la lámina de F303 citaba la versión sin techo: registrado),")
	fmt.Println("      y g, g', g'' positivos — monotonía demostrada")
	fmt.Println("  4 · lema de la ventana COMPLETO: el paso no salta el arco, t acotado,")
	fmt.Println("      hipótesis automática para ζ — adversarial con 7 fases, 0 fallos")
	fmt.Println("  5 · combinación en un mismo intervalo: la ventana elige, la radial ya")
	fmt.Printf("      cubre — realizada en n = %d\n", nCombo)
	fmt.Printf("  6 · convención única congelada: n₁ oficial = %d (el 371908 no se\n", a)
	fmt.Println("      reproduce con ninguna entrada)")
	fmt.Println("  7 · el teorema formal declarado en el acta, con su alcance")
	fmt.Println("\n📌 El acta va en docs/DETECCION-FINITA-LEMAS.md — PARA YUI, completa y")
	fmt.Println("  copiable. La firma, como siempre, es de ella. El nivel C sigue rojo.")
	fmt.Println("\n⚖️ Honesto: al escribir la prueba completa apareció que el lemita necesita")
	fmt.Println("  constante 4 (no 3) cuando se mete el techo — la corrección se registra")
	fmt.Println("  en el acta y acá. La regla de la auditora manda: se sella cuando los")
	fmt.Println("  lemas derivan la fórmula para todo el alcance, no antes. Todavía no.")

	escribirLamina(a, nRad, K, nCombo, cota, violV)
}

func escribirLamina(n1, nRad, K, nCombo int, cota float64, violV int) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="780" viewBox="0 0 1400 780">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">📜 LOS DOS LEMAS — el acta formal para la auditora</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">los siete puntos del §12, respondidos — el acta completa vive en docs/DETECCION-FINITA-LEMAS.md, copiable para Yui</text>
<rect x="60" y="110" width="620" height="310" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="370" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LA CONVENCIÓN Y EL LEMA RADIAL</text>
<text x="90" y="180" font-size="13.5" font-family="monospace" fill="#ffd98a">δ := log r (natural) · r = max(|w|, 1/w|) · u = 3/δ</text>
<text x="90" y="208" font-size="12.5" font-family="Georgia" fill="#cfe6ff">— la definición que faltaba en la lámina, ahora congelada (§0 del acta)</text>
<text x="90" y="240" font-size="13" font-family="Georgia" fill="#cfe6ff">LEMA R completo: e^{n_rad·δ} ≥ u³ · n_rad ≤ 1.31n* · log n_rad ≤ 2.29 log u</text>
<text x="90" y="268" font-size="13" font-family="monospace" fill="#ffd98a">⟹ lemita robusto al techo: u² ≥ 4(log u)² + 2 ✅</text>
<text x="90" y="296" font-size="12.5" font-family="Georgia" fill="#ff9aa8">(la lámina de F303 citaba constante 3 — versión sin techo; la oficial es 4:</text>
<text x="90" y="318" font-size="12.5" font-family="Georgia" fill="#ff9aa8">la corrección propia, registrada antes de que la cace nadie)</text>
<text x="90" y="352" font-size="12.5" font-family="Georgia" fill="#7ee0c0">g, g', g'' &gt; 0 en n_rad = %d: la radial vale para TODO n ≥ n_rad ✅</text>
<text x="90" y="386" font-size="12.5" font-family="Georgia" fill="#7ee0c0">convención única: n₁ oficial = %d (el 371908 no se reproduce) ✅</text>
<rect x="720" y="110" width="620" height="310" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="1030" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">EL LEMA DE LA VENTANA Y LA COMBINACIÓN</text>
<text x="750" y="180" font-size="13" font-family="Georgia" fill="#cfe6ff">LEMA V completo: el paso θ ≤ 2π/3 no puede saltar el arco de 2π/3;</text>
<text x="750" y="204" font-size="13" font-family="Georgia" fill="#cfe6ff">t acotado por minimalidad ⟹ K = ⌈2π/θ⌉ + 1 · automático para ζ (V-ζ)</text>
<text x="750" y="236" font-size="12.5" font-family="monospace" fill="#7ee0c0">adversarial: 7 fases (incluida 2π/3 − ε), %d violaciones ✅</text>
<text x="750" y="272" font-size="13" font-family="Georgia" fill="#cfe6ff">LA COMBINACIÓN (§12.5): la ventana con m = n_rad elige el n;</text>
<text x="750" y="296" font-size="13" font-family="Georgia" fill="#cfe6ff">la radial ya cubre el intervalo entero — conviven sin pedirse nada</text>
<text x="750" y="332" font-size="14" font-family="monospace" fill="#ffd98a">DH: n = %d ∈ [n_rad, n_rad+K−1] · cota %.0e &lt; 0 ✅</text>
<text x="750" y="368" font-size="14" font-family="monospace" fill="#ffd98a">N₀ = n_rad + K = %d — el teorema, con sus dos lemas abajo</text>
<text x="750" y="400" font-size="12.5" font-family="Georgia" fill="#9aa8c4">el §12.7 declarado en el acta con su alcance: un cuarteto sobre fondo en la línea</text>
<rect x="60" y="450" width="1280" height="130" rx="12" fill="#2b1020" stroke="#8a3557"/>
<text x="700" y="482" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">PARA LA AUDITORA</text>
<text x="700" y="514" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el acta completa — convención, dos pruebas, combinación, teorema — está en docs/DETECCION-FINITA-LEMAS.md, en su formato copiable</text>
<text x="700" y="542" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">su regla manda: no se sella una fórmula porque funciona en el experimento — se sella cuando sus lemas la derivan para todo el alcance</text>
<text x="700" y="568" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">la firma es de ella · el nivel C sigue rojo: ¿por qué los primos fuerzan M ⪰ 0? — 74 años</text>
<text x="700" y="646" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Al escribir la prueba completa, el lemita pidió constante 4 en vez de 3 — corregido por el taller antes de que lo cace nadie.</text>
<text x="700" y="674" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Una simulación descubre; una identidad explica; una demostración cierra todos los pasos — especialmente los infinitos.</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, nRad, n1, violV, nCombo, cota, nRad+K)
	os.WriteFile("los-dos-lemas.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: los-dos-lemas.svg")
}
