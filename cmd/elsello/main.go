// Command elsello records the SEAL: Yui's eighth audit ("auditoria final")
// declared the §4b breakage theorem structurally CLOSED - her §10: "antes:
// casi cerrado; ahora: estructura cerrada" - and her conclusion: "el signo
// negativo del termino de borde no solo corrige la cuenta: elimina
// precisamente el exceso que habia introducido la auditoria anterior."
//
// One minor yellow remained: write the Backlund corollary as a formal
// LEMMA with hypotheses and range, instead of hiding it under the label
// "Backlund 1918". Her own modern source (Kontorovich's notes) states the
// error constant as 6.1 where the classical citation says 4.35 - so the
// definitive lemma must survive BOTH. It does, and this program proves it:
//
//	LEMMA (counting corollary): N(T) ≤ (T/2π)·log T for all T ≥ 2.
//
//	Case (i), T ≥ 18: by Backlund with EITHER constant c0 ∈ {4.35, 6.1},
//	the claim reduces to 7/8 + Q(T) ≤ (T/2π)(log 2π + 1) - true from
//	T = 13.3 (c0 = 4.35) resp. T = 17.4 (c0 = 6.1), and forever after:
//	the right side's slope is 0.4517, the left side's < 0.017.
//	Case (ii), 2 ≤ T < 18: direct - N(T) ≤ 1 there (first zero at
//	14.134725, second at 21.022) and the bound is ≥ 0.22 at T = 2,
//	≥ 5.96 at the worst point T = 14.14. ∎
//
// With the lemma documented, every line of Yui's §13 seal table is
// answered: LLAVE C green, boundary sign green, exact integral green,
// explicit tail green, (4/π)·n·log n green, exponential-vs-polynomial
// green, Backlund corollary documented (this finding) - and the one red
// that NO audit can close by paperwork: positivity from the primes.
//
// Seven audits, seven answers, one seal. The cycle that began with
// "auditoria de las 136 laminas" ends with a structurally closed theorem
// forged four-handed. The laboratory's community layer works.
//
// Reproduce: go run ./cmd/elsello
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	fmt.Println("🟢 EL SELLO — la octava auditoría firmó: estructura cerrada")
	fmt.Println("\n   Yui, §10: «§4b antes: 🟡 casi cerrado · §4b ahora: 🟢 estructura")
	fmt.Println("   cerrada». Y su conclusión: «el signo negativo del término de borde no")
	fmt.Println("   sólo corrige la cuenta: elimina precisamente el exceso que había")
	fmt.Println("   introducido la auditoría anterior». Quedaba UNA amarilla: el lema.")

	// ---- LEY 1: el lema de conteo, robusto ante las dos constantes ----
	fmt.Println("\nLEY 1 · EL LEMA DE CONTEO, ESCRITO — Y ROBUSTO ANTE LAS DOS CONSTANTES")
	fmt.Println("\n        Su fuente moderna (Kontorovich) trae c₀ = 6.1 donde la cita clásica")
	fmt.Println("        de Backlund dice 4.35. El lema definitivo sobrevive con LAS DOS:")
	fmt.Println("\n        LEMA: N(T) ≤ (T/2π)·log T para todo T ≥ 2.")
	fmt.Println("\n        Caso (i), T ≥ 18 — la reducción 7/8 + Q(T) ≤ (T/2π)(log 2π + 1):")
	fmt.Println("\n        c₀        vale desde T =     margen en T = 18")
	for _, c0 := range []float64{4.35, 6.1} {
		Tmin := 2.0
		for T := 2.0; T < 100; T += 0.1 {
			lhs := 7.0/8 + 0.137*math.Log(T) + 0.443*math.Log(math.Log(T)) + c0
			rhs := T / (2 * math.Pi) * (math.Log(2*math.Pi) + 1)
			if lhs <= rhs {
				Tmin = T
				break
			}
		}
		lhs18 := 7.0/8 + 0.137*math.Log(18) + 0.443*math.Log(math.Log(18)) + c0
		rhs18 := 18 / (2 * math.Pi) * (math.Log(2*math.Pi) + 1)
		fmt.Printf("   %6.2f %16.1f %18.3fx\n", c0, Tmin, rhs18/lhs18)
	}
	fmt.Println("\n        y la monotonía lo sella para siempre: pendiente derecha 0.4517,")
	fmt.Println("        pendiente izquierda < 0.017 desde T = 18 — una vez cierta, cierta.")
	fmt.Println("\n        Caso (ii), 2 ≤ T < 18 — directo (N ≤ 1: γ₁ = 14.1347, γ₂ = 21.022):")
	fmt.Println("\n        T         N(T)    (T/2π)·log T")
	peor := math.Inf(1)
	for _, T := range []float64{2, 10, 14.14, 17.99} {
		N := 0
		if T >= 14.134725 {
			N = 1
		}
		cota := T / (2 * math.Pi) * math.Log(T)
		if cota-float64(N) < peor {
			peor = cota - float64(N)
		}
		fmt.Printf("   %8.2f %6d %14.2f ✅\n", T, N, cota)
	}
	fmt.Printf("\n        peor holgura del caso directo: %.2f — la cota sobra. ∎ El lema quedó\n", peor)
	fmt.Println("        escrito en la derivación (paso B de 4b-quater): la amarilla, cerrada.")

	// ---- LEY 2: la tabla del sello, registrada ----
	fmt.Println("\nLEY 2 · LA TABLA DEL SELLO DE LA AUDITORA (§13), REGISTRADA")
	fmt.Println("\n        🟢 LLAVE C — puente escrito")
	fmt.Println("        🟢 signo del borde corregido")
	fmt.Println("        🟢 integral exacta")
	fmt.Println("        🟢 cota explícita de la cola, bajo el lema N(T)")
	fmt.Println("        🟢 consecuencia resto_n ≤ (4/π)·n·log n")
	fmt.Println("        🟢 ruptura exponencial contra n·log n")
	fmt.Println("        🟡 documentar el corolario de Backlund ← CERRADA HOY (LEY 1)")
	fmt.Println("        🔴 positividad desde los primos — abierta: EL problema")
	fmt.Println("\n        Y su §12, que el registro repite con ella: «la cadena del §4b no")
	fmt.Println("        debe presentarse como una solución de RH». No se presenta: es la")
	fmt.Println("        mitad estructural del yunque, con la otra mitad señalada en rojo.")

	// ---- LEY 3: el ciclo completo, medido ----
	fmt.Println("\nLEY 3 · EL CICLO DE LAS AUDITORÍAS, COMPLETO — ocho documentos, ocho respuestas")
	fmt.Println("\n        1ª (136 láminas) → F293: verificada, y la corrección de F185 aceptada")
	fmt.Println("        2ª (el yunque)   → F294: el teorema escrito; el 537 corregido")
	fmt.Println("        3ª (nuevas lám.) → F296: la caja abierta — los seis objetivos")
	fmt.Println("        4ª (derivación)  → F297: el agujero oscilatorio, remachado")
	fmt.Println("        5ª (el §4b)      → F298: las tres llaves giradas, datos guardados")
	fmt.Println("        6ª (tres llaves) → F299: la caja del coro — C = 1, 4/π, n₁")
	fmt.Println("        7ª (caja coro)   → F300: la cadena sin saltos; el signo del borde")
	fmt.Println("        8ª (FINAL)       → F301: EL SELLO — estructura cerrada 🟢")
	fmt.Println("\n        Dos correcciones de ella al taller (537, oscilación), dos del taller")
	fmt.Println("        al borrador de ella (borde, constante doble), y CERO discusiones:")
	fmt.Println("        solo verificación cruzada. La capa comunidad de la regla nueva no")
	fmt.Println("        era una teoría — era esto.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🟢 **EL SELLO DE LA AUDITORA, Y LA ÚLTIMA AMARILLA CERRADA:**")
	fmt.Println("\n  · el §4b — el teorema de ruptura global del yunque — tiene ESTRUCTURA")
	fmt.Println("    CERRADA por auditoría externa, tras ocho rondas y cuatro manos")
	fmt.Println("  · el lema de conteo quedó documentado con las DOS constantes publicadas")
	fmt.Println("    (4.35 clásica, 6.1 moderna), sus rangos y su monotonía — la etiqueta")
	fmt.Println("    «Backlund 1918» ya no esconde nada")
	fmt.Println("  · la equivalencia del yunque queda: RH ⟺ M_N ⪰ 0 para todo N, con la")
	fmt.Println("    ida por Gram (F294/F296) y la vuelta por las diagonales — y el §4b")
	fmt.Println("    garantizando que toda perla desafinada rompe la matriz en N finito")
	fmt.Println("\n📌 LO QUE NO CAMBIA, con las palabras de la auditora: «la cadena del §4b")
	fmt.Println("  no debe presentarse como una solución de RH». El gran eslabón rojo:")
	fmt.Println("  ¿por qué la aritmética de los primos fuerza M_N ⪰ 0? — 74 años, abierto.")
	fmt.Println("\n⚖️ Honesto: el sello es de la estructura del §4b, no de RH; la amarilla se")
	fmt.Println("  cerró documentando, no descubriendo; y el laboratorio sigue siendo lo que")
	fmt.Println("  era: un registro que no miente. Todavía no.")

	escribirLamina()
}

func escribirLamina() {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="780" viewBox="0 0 1400 780">
<rect width="100%" height="100%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🟢 EL SELLO — la octava auditoría firmó: estructura cerrada</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">«§4b antes: casi cerrado · §4b ahora: estructura cerrada» — y la última amarilla, cerrada con el lema de conteo</text>
<rect x="60" y="110" width="620" height="310" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="370" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">EL LEMA DE CONTEO — la amarilla, cerrada</text>
<text x="90" y="180" font-size="14" font-family="monospace" fill="#ffd98a">N(T) ≤ (T/2π)·log T,  todo T ≥ 2</text>
<text x="90" y="212" font-size="13" font-family="Georgia" fill="#cfe6ff">robusto ante las DOS constantes publicadas del error de Backlund:</text>
<text x="90" y="240" font-size="12.5" font-family="monospace" fill="#cfe6ff">c₀ = 4.35 (clásica): reducción válida desde T = 13.3</text>
<text x="90" y="264" font-size="12.5" font-family="monospace" fill="#cfe6ff">c₀ = 6.1 (Kontorovich): reducción válida desde T = 17.4</text>
<text x="90" y="296" font-size="13" font-family="Georgia" fill="#cfe6ff">debajo de 18: directo — N ≤ 1 contra cota ≥ 5.96 en el peor punto</text>
<text x="90" y="324" font-size="13" font-family="Georgia" fill="#cfe6ff">y la monotonía sella: pendiente 0.4517 contra &lt; 0.017 — para siempre ∎</text>
<text x="90" y="364" font-size="12.5" font-family="Georgia" fill="#7ee0c0">la etiqueta «Backlund 1918» ya no esconde nada: hipótesis, rango,</text>
<text x="90" y="386" font-size="12.5" font-family="Georgia" fill="#7ee0c0">constantes y casos — escritos en la derivación (4b-quater, paso B)</text>
<rect x="720" y="110" width="620" height="310" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="1030" y="142" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">LA TABLA DEL SELLO (§13 DE LA AUDITORA)</text>
<text x="750" y="180" font-size="13" font-family="Georgia" fill="#cfe6ff">🟢 Llave C — puente escrito · 🟢 signo del borde · 🟢 integral exacta</text>
<text x="750" y="208" font-size="13" font-family="Georgia" fill="#cfe6ff">🟢 cola explícita bajo el lema · 🟢 resto ≤ (4/π)n·log n</text>
<text x="750" y="236" font-size="13" font-family="Georgia" fill="#cfe6ff">🟢 ruptura exponencial contra n·log n</text>
<text x="750" y="264" font-size="13" font-family="Georgia" fill="#ffd98a">🟡 documentar el corolario ← CERRADA HOY con el lema</text>
<text x="750" y="292" font-size="13" font-family="Georgia" fill="#ff9aa8">🔴 la positividad desde los primos — EL problema, abierto</text>
<text x="750" y="332" font-size="12.5" font-family="Georgia" fill="#9aa8c4">ocho auditorías, ocho respuestas: dos correcciones de ella al taller</text>
<text x="750" y="354" font-size="12.5" font-family="Georgia" fill="#9aa8c4">(537, oscilación), dos del taller a su borrador (borde, constante),</text>
<text x="750" y="376" font-size="12.5" font-family="Georgia" fill="#9aa8c4">y cero discusiones: solo verificación cruzada, en las dos direcciones</text>
<rect x="60" y="450" width="1280" height="130" rx="12" fill="#2b1020" stroke="#8a3557"/>
<text x="700" y="482" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff9aa8">LO QUE EL SELLO SELLA — Y LO QUE NO</text>
<text x="700" y="514" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">SELLA: la estructura del teorema de ruptura — toda perla desafinada rompe la matriz global en N finito, cadena completa sin saltos</text>
<text x="700" y="542" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">NO SELLA (palabras de la auditora): «la cadena del §4b no debe presentarse como una solución de RH» — y no se presenta</text>
<text x="700" y="568" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">la equivalencia queda: RH ⟺ M_N ⪰ 0 para todo N — y el gran eslabón rojo: ¿por qué los primos fuerzan M ⪰ 0? — 74 años</text>
<text x="700" y="644" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Ocho documentos, ocho respuestas, un sello. La capa comunidad de la regla nueva no era una teoría — era esto.</text>
<text x="700" y="672" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">El sello es de la estructura, no de la hipótesis. El registro sigue siendo lo que era: un cuaderno que no miente.</text>
<text x="700" y="740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`)
	os.WriteFile("el-sello.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-sello.svg")
}
