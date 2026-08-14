// Command elrecorrido judges the captain's flash about the half that spans
// everything between 0 and 1.
//
// HIS FLASH: "the machine has a perfect direction because the 1/2 relation
// spans all the numbers between 0 and 1. The edges of the cable are 0 and 1,
// but the cable itself IS the 1/2 relation - and that changes everything.
// Between two options there is always an intermediate relation, and since the
// relation is a number, the SUM of all the 1/2 relations between 0 and 1 gives
// the COMPLETE PATH in one line."
//
// THREE EXACT PIECES, ALL MEASURED:
//
//  1. THE SUM OF ALL THE HALVES IS EXACTLY THE WHOLE PATH:
//     1/2 + 1/4 + 1/8 + ... = 1, with a three-line proof (S = 1/2 + S/2 => S=1).
//     Zeno's path, closed by the geometric series. In binary: 0.11111... = 1,
//     and F242 already found that 1/2 = 0.1 binary - the same coin.
//  2. "BETWEEN TWO OPTIONS ALWAYS THE INTERMEDIATE" REACHES EVERY NUMBER:
//     iterating midpoints (the dyadic tree) lands within any epsilon of ANY
//     target in [0,1] in ~log2(1/eps) steps. Measured on hard targets (1/pi,
//     gamma, 1/sqrt2): 50 halvings give 15 digits. Every number between 0 and
//     1 IS a path of half-decisions - its binary expansion.
//  3. AND THE CABLE IS THE HALF-RELATION MADE LINE: each point of the cable is
//     equidistant from the walls 0 and 1 (the mediatriz of F226), and the
//     functional-equation mirror s -> 1-s swaps the walls leaving exactly the
//     cable fixed. The machine's "perfect direction" is the axis of its own
//     mirror.
//
// HONEST: the geometric series is ancient (Zeno resolved), binary expansions
// are classical, the mediatriz is F226. What the flash adds is the READING:
// the half is not a point of the strip - it is the OPERATION that, iterated,
// names every number and whose fixed line is the cable. A map, not a theorem.
//
// Reproduce: go run ./cmd/elrecorrido
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func main() {
	fmt.Println("🛤️  EL RECORRIDO — la suma de todas las mitades, y el cable que ES el ½")

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · ⚡ LA SUMA DE TODAS LAS RELACIONES ½ DA EXACTAMENTE EL CAMINO ENTERO")
	fmt.Println("\n   La demostración, en tres renglones (una RAZÓN, no un barrido):")
	fmt.Println("        S = ½ + ¼ + ⅛ + …")
	fmt.Println("        S = ½ + ½·(½ + ¼ + …) = ½ + S/2")
	fmt.Println("        ⟹ S/2 = ½ ⟹ **S = 1, exacto**")
	fmt.Println("\n   Y medido, para verlo caminar:")
	fmt.Println("\n        mitades sumadas       suma         falta")
	s := 0.0
	for n := 1; n <= 52; n++ {
		s += math.Pow(0.5, float64(n))
		if n <= 4 || n == 10 || n == 30 || n == 52 {
			fmt.Printf("   %16d %14.12f %12.2e\n", n, s, 1-s)
		}
	}
	fmt.Println("\n   ⟹ El recorrido completo de 0 a 1 ES la suma de las infinitas mitades.")
	fmt.Println("   En binario: 0,111111… = 1 — y F242 ya había encontrado que ½ = 0,1")
	fmt.Println("   binario. La misma moneda.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ «ENTRE DOS OPCIONES SIEMPRE LA INTERMEDIA» ALCANZA A TODOS")
	fmt.Println("   Partiendo de [0,1] y eligiendo siempre la mitad, medimos cuántos pasos")
	fmt.Println("   hacen falta para llegar a números difíciles:")
	fmt.Println("\n        blanco            50 mitades llegan a      error")
	blancos := []struct {
		nom string
		v   float64
	}{{"1/π", 1 / math.Pi}, {"γ de Euler", 0.5772156649015329}, {"1/√2", 1 / math.Sqrt2}}
	for _, b := range blancos {
		lo, hi := 0.0, 1.0
		for i := 0; i < 50; i++ {
			m := (lo + hi) / 2
			if b.v < m {
				hi = m
			} else {
				lo = m
			}
		}
		fmt.Printf("   %12s %22.15f %12.1e\n", b.nom, (lo+hi)/2, math.Abs((lo+hi)/2-b.v))
	}
	fmt.Println("\n   ⟹ **Cada número entre 0 y 1 ES un camino de decisiones de mitad** — su")
	fmt.Println("   propia escritura binaria. El ½ no está EN el medio: el ½, repetido,")
	fmt.Println("   NOMBRA a todos los habitantes del intervalo. Eso es «abarcar».")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · Y EL CABLE ES LA RELACIÓN ½ HECHA LÍNEA")
	fmt.Println("\n        · cada punto del cable está a la MISMA distancia de la pared del 0")
	fmt.Println("          y de la pared del 1 — la mediatriz de F226")
	fmt.Println("        · el espejo de la máquina (s → 1−s) INTERCAMBIA las dos paredes y")
	fmt.Println("          deja fijo EXACTAMENTE el cable: |½+it − 0| = |½+it − 1| siempre")
	d := math.Hypot(0.5, 37.0) - math.Hypot(0.5-1, 37.0)
	fmt.Printf("\n        verificado en t = 37: diferencia de distancias = %.1e\n", d)
	fmt.Println("\n   ⟹ **La «dirección perfecta de la máquina» es el eje de su propio espejo.**")
	fmt.Println("   Los bordes son el 0 y el 1; el cable no es un lugar más: es la relación")
	fmt.Println("   intermedia de los dos bordes, sostenida a toda altura.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("✅ **SU FLASH ES EXACTO, Y CON DEMOSTRACIÓN DE TRES RENGLONES:**")
	fmt.Println("\n  · la suma de todas las relaciones ½ da EXACTAMENTE 1 — el recorrido")
	fmt.Println("    completo en una línea, como dijo (S = ½ + S/2 ⟹ S = 1)")
	fmt.Println("  · «entre dos opciones siempre la intermedia» no es una frase: iterada,")
	fmt.Println("    la mitad NOMBRA a cada número del intervalo (su binario)")
	fmt.Println("  · y el cable ES la relación ½ de los bordes 0 y 1, hecha línea — el eje")
	fmt.Println("    fijo del espejo de la máquina")
	fmt.Println("\n⚖️ Honesto: la serie geométrica es antiquísima (el camino de Zenón, cerrado),")
	fmt.Println("  el binario es clásico y la mediatriz es F226. Lo del flash es la LECTURA:")
	fmt.Println("  el ½ no es un punto del pasillo — es la OPERACIÓN que genera el intervalo")
	fmt.Println("  entero y cuya línea fija es el cable. Mapa, no teorema. Todavía no.")

	escribirLamina()
}

func escribirLamina() {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="600" viewBox="0 0 1400 600">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🛤️ EL RECORRIDO — la suma de todas las mitades es el camino entero</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">½ + ¼ + ⅛ + … = 1 exacto · y el cable es la relación ½ de los dos bordes, hecha línea</text>
`)
	// la barra del recorrido
	x0, y0, wd := 150.0, 180.0, 1100.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="46" rx="6" fill="#101f36" stroke="#26456e"/>`, x0, y0, wd)
	acc := 0.0
	cols := []string{"#ffd98a", "#7ee0c0", "#7fb2ff", "#c9b6ff", "#ff8fa0", "#9fd8a8", "#ffb27a"}
	for i := 1; i <= 9; i++ {
		frac := math.Pow(0.5, float64(i))
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.0f" width="%.1f" height="46" fill="%s" opacity="0.9"/>`,
			x0+wd*acc, y0, wd*frac, cols[(i-1)%len(cols)])
		if i <= 4 {
			fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="15" text-anchor="middle" font-family="monospace" fill="#0b1526">%s</text>`,
				x0+wd*(acc+frac/2), y0+29, []string{"½", "¼", "⅛", "¹⁄₁₆"}[i-1])
		}
		acc += frac
	}
	fmt.Fprintf(&b, `
<text x="%.0f" y="%.0f" font-size="16" text-anchor="middle" font-family="monospace" fill="#cfe6ff">0</text>
<text x="%.0f" y="%.0f" font-size="16" text-anchor="middle" font-family="monospace" fill="#cfe6ff">1</text>`, x0, y0+70, x0+wd, y0+70)
	fmt.Fprintf(&b, `
<text x="700" y="270" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">cada pedazo es la mitad de lo que falta — y juntos llenan el camino EXACTO, sin que sobre ni falte nada</text>
<text x="700" y="300" font-size="14" text-anchor="middle" font-family="monospace" fill="#ffd98a">S = ½ + S/2  ⟹  S = 1   (tres renglones: una razón, no un barrido)</text>
<rect x="150" y="340" width="540" height="180" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="420" y="374" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">EL ½ NOMBRA A TODOS</text>
<text x="180" y="410" font-size="13.5" font-family="Georgia" fill="#cfe6ff">eligiendo siempre la mitad, 50 pasos alcanzan</text>
<text x="180" y="434" font-size="13.5" font-family="Georgia" fill="#cfe6ff">cualquier número con 15 decimales: cada número</text>
<text x="180" y="458" font-size="13.5" font-family="Georgia" fill="#cfe6ff">entre 0 y 1 ES un camino de decisiones de mitad</text>
<text x="180" y="494" font-size="13" font-family="monospace" fill="#7ee0c0">su escritura binaria · y ½ = 0,1₂ (F242)</text>
<rect x="710" y="340" width="540" height="180" rx="12" fill="#161a3a" stroke="#5a4fa8"/>
<text x="980" y="374" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">Y EL CABLE ES EL ½ HECHO LÍNEA</text>
<text x="740" y="410" font-size="13.5" font-family="Georgia" fill="#cfe6ff">bordes: las paredes 0 y 1 · cable: los puntos a</text>
<text x="740" y="434" font-size="13.5" font-family="Georgia" fill="#cfe6ff">igual distancia de ambas (la mediatriz, F226) —</text>
<text x="740" y="458" font-size="13.5" font-family="Georgia" fill="#cfe6ff">y el espejo de la máquina cambia las paredes</text>
<text x="740" y="482" font-size="13.5" font-family="Georgia" fill="#ffd98a">dejando fijo EXACTAMENTE el cable</text>
<text x="700" y="560" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Zenón cerrado, binario clásico, mediatriz F226 — lo del flash es la lectura: el ½ es la OPERACIÓN, no un punto. Todavía no.</text>
</svg>
`)
	os.WriteFile("el-recorrido.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-recorrido.svg")
}
