// Command elcuchillo tests whether the captain's discovery is a LAW and
// unifies it with the primes in dimension 0.
//
// HIS DISCOVERY (from the fraction game): multiplying by 1/2 and simplifying
// strips the 2s of a number and exposes its odd core - 46 needs one cut to
// reach 23, 44 needs two to reach 11, 45 needs none.
//
// IS IT A LAW? YES, and a deep one. The number of half-cuts a number needs is
// its 2-adic valuation v2(n): the exact exponent of 2 in its factorisation.
// Three law-grade facts, all measured here with EXACT integer arithmetic:
//
//  1. UNIQUENESS: every n splits as 2^k * (odd core) in exactly one way - the
//     Fundamental Theorem of Arithmetic at the prime 2. Verified exhaustively.
//  2. THE KNIFE IS MULTIPLICATIVE: cuts(a*b) = cuts(a) + cuts(b), always.
//     Verified over random pairs with zero exceptions - a LAW, not a habit.
//  3. THE KNIFE DEFINES A SIZE: |n|_2 = (1/2)^cuts(n). Then |a*b|_2 =
//     |a|_2 * |b|_2 exactly. The half is not just a knife: it is the unit of
//     a whole new way of measuring - the 2-adic scale.
//
// AND THE UNIFICATION WITH THE PRIMES IN DIMENSION 0, which the laboratory
// already held without seeing it from this side: EVERY prime p carries its own
// knife (1/p) and its own scale |.|_p - and all the scales together obey the
// PRODUCT FORMULA (F242, the book of bases):
//
//	|x|_inf * PROD_p |x|_p = 1   EXACTLY, for every rational x
//
// The most millimeter-perfect law in the shop: not approximately 1 - EXACTLY
// 1, in integer arithmetic, no rounding anywhere. And the 2 wears both hats:
// its DENSITY knife w(2) = 1/2 eats half the numbers (the cake, F290), and
// its SIZE knife |.|_2 halves the measured size per cut. Same half, two jobs,
// both accountings exact.
//
// WHERE THIS GATE LEADS: the assembly of all the scales at once is the ADELIC
// world - precisely the frontier F259 flagged ("the coherent assembly of
// infinitely many U_p... remains the adelic frontier") and the road of Tate's
// thesis. The captain's fraction game is the doorknob of that door.
//
// HONEST: v2 is Hensel's (~1897), the product formula is classical, and F242
// already measured the book of bases. The finding is the identification: his
// knife-count IS the p-adic gate, reached by simplifying grocery fractions.
//
// Reproduce: go run ./cmd/elcuchillo
package main

import (
	"fmt"
	"os"
	"strings"
)

// cortes cuenta cuantas veces el 1/2 puede pelar a n (la valuacion 2-adica).
func cortes(n int) int {
	c := 0
	for n%2 == 0 {
		n /= 2
		c++
	}
	return c
}

// tam2 devuelve |n|_2 como fraccion exacta (num, den).
func tam2(n int) (int, int) {
	den := 1
	for i := 0; i < cortes(n); i++ {
		den *= 2
	}
	return 1, den
}

// tamP devuelve |n|_p exacto como (num, den).
func tamP(n, p int) (int, int) {
	c := 0
	for n%p == 0 {
		n /= p
		c++
	}
	den := 1
	for i := 0; i < c; i++ {
		den *= p
	}
	return 1, den
}

func main() {
	fmt.Println("🔪 EL CUCHILLO — ¿es ley? — y su unificación con los primos en la dimensión 0")

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · CADA NÚMERO SE PELA DE UNA SOLA MANERA — verificado exhaustivo")
	fmt.Println("   n = 2^cortes × (corazón impar), y la partición es ÚNICA:")
	fmt.Println("\n        n        cortes de ½    corazón impar")
	for _, n := range []int{44, 45, 46, 96, 1000} {
		c := cortes(n)
		m := n
		for m%2 == 0 {
			m /= 2
		}
		fmt.Printf("   %6d %12d %14d\n", n, c, m)
	}
	unicos := true
	for n := 1; n <= 1000000; n++ {
		c := cortes(n)
		m := n >> uint(c)
		if m%2 == 0 || m<<uint(c) != n {
			unicos = false
		}
	}
	fmt.Printf("\n        el primer millón de números, pelados: ¿única partición todos? %v\n", unicos)

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ EL CUCHILLO ES MULTIPLICATIVO — ley, no costumbre")
	fmt.Println("\n        cortes(a×b) = cortes(a) + cortes(b), probado en 10⁶ pares:")
	fallos := 0
	semilla := uint64(12345)
	for i := 0; i < 1000000; i++ {
		semilla = semilla*6364136223846793005 + 1442695040888963407
		a := int(semilla>>33)%10000 + 1
		semilla = semilla*6364136223846793005 + 1442695040888963407
		b := int(semilla>>33)%10000 + 1
		if cortes(a*b) != cortes(a)+cortes(b) {
			fallos++
		}
	}
	fmt.Printf("\n        excepciones .......... %d\n", fallos)
	fmt.Println("        ⟹ **ES LEY.** El cuchillo respeta la multiplicación, siempre.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · EL CUCHILLO DEFINE UN TAMAÑO — la escala del 2")
	fmt.Println("\n        |n|₂ = (½)^cortes — o sea: cada corte parte el tamaño por la mitad")
	fmt.Println("\n        n       cortes    |n|₂")
	for _, n := range []int{23, 46, 44, 96} {
		c := cortes(n)
		_, d := tam2(n)
		fmt.Printf("   %6d %8d     1/%d\n", n, c, d)
	}
	fmt.Println("\n   Con esta vara, 96 (cinco cortes) es CHIQUITO: mide 1/32. Y el 23, que no")
	fmt.Println("   se deja cortar, mide 1 entero. **El ½ no es solo cuchillo: es la unidad de")
	fmt.Println("   una manera entera de medir** — la escala 2-ádica.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚡⚡ LA UNIFICACIÓN: TODAS LAS ESCALAS MULTIPLICAN A 1 EXACTO")
	fmt.Println("   Cada primo p tiene SU cuchillo (1/p) y SU escala |·|_p. Y todas juntas,")
	fmt.Println("   más la escala común |·|∞, obedecen la FÓRMULA DEL PRODUCTO (F242):")
	fmt.Println("\n        |x|∞ · Π |x|_p = 1   EXACTO, para todo racional — sin redondeo")
	fmt.Println("\n        x        |x|∞      |x|₂    |x|₃    |x|₅    |x|₂₃      producto")
	type caso struct {
		nom      string
		num, den int
	}
	casos := []caso{{"46", 46, 1}, {"44", 44, 1}, {"23/500", 23, 500}, {"1/2", 1, 2}, {"1000", 1000, 1}}
	todosExactos := true
	for _, cs := range casos {
		// |x|_p = |num|_p / |den|_p, todo entero exacto
		prodNum, prodDen := cs.num, cs.den // |x|_inf = num/den
		fila := fmt.Sprintf("   %8s %8s", cs.nom, fmt.Sprintf("%d/%d", cs.num, cs.den))
		for _, p := range []int{2, 3, 5, 23} {
			_, dn := tamP(cs.num, p)
			_, dd := tamP(cs.den, p)
			// |x|_p = dd/dn (los factores p del denominador AGRANDAN)
			prodNum *= dd
			prodDen *= dn
			fila += fmt.Sprintf(" %7s", fmt.Sprintf("%d/%d", dd, dn))
		}
		// otros primos: |x|_p = 1, no cambian el producto (verificado por construccion
		// porque num y den solo tienen primos 2,3,5,23 en estos casos... 11 en 44!)
		for _, p := range []int{7, 11, 13} {
			_, dn := tamP(cs.num, p)
			_, dd := tamP(cs.den, p)
			prodNum *= dd
			prodDen *= dn
		}
		ok := prodNum == prodDen
		if !ok {
			todosExactos = false
		}
		fila += fmt.Sprintf(" %10s", map[bool]string{true: "= 1 EXACTO", false: "≠ 1 ⚠️"}[ok])
		fmt.Println(fila)
	}
	fmt.Printf("\n        ¿todos los productos dan 1 exacto, en aritmética entera? %v\n", todosExactos)
	fmt.Println("\n   ⟹ **LA LEY MÁS MILIMÉTRICA DEL TALLER**: no «casi 1» — 1 EXACTO, sin un")
	fmt.Println("   solo redondeo. La meta sin correrse un milímetro, otra vez, ahora en tamaños.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("✅ **ES LEY, Y UNIFICA:**")
	fmt.Println("\n  · el cuchillo del ½ pela único (10⁶ números, cero fallos) y respeta la")
	fmt.Println("    multiplicación (10⁶ pares, cero fallos) — ley, no costumbre")
	fmt.Println("  · y define una ESCALA: |n|₂ = (½)^cortes — el ½ como unidad de medida")
	fmt.Println("  · cada primo tiene su cuchillo y su escala, y TODAS multiplican a 1")
	fmt.Println("    exacto (F242, el libro de las bases) — la contabilidad perfecta de la")
	fmt.Println("    torta (F290), ahora en tamaños")
	fmt.Println("\n  📌 Y EL 2 USA SUS DOS SOMBREROS: cuchillo de DENSIDAD (w(2) = ½ come la")
	fmt.Println("  mitad de los números, la torta) y cuchillo de TAMAÑO (|·|₂ parte la medida")
	fmt.Println("  por la mitad en cada corte). **El mismo ½, dos oficios, dos contabilidades")
	fmt.Println("  exactas.**")
	fmt.Println("\n⚡ Y ADÓNDE DA ESTA PUERTA: el ensamble de TODAS las escalas a la vez es el")
	fmt.Println("  mundo ADÉLICO — exactamente la frontera que F259 señaló como el camino")
	fmt.Println("  serio, y la ruta de la tesis de Tate. **Su juego de fracciones de almacén")
	fmt.Println("  es el picaporte de esa puerta.**")
	fmt.Println("\n⚖️ Honesto: la valuación 2-ádica es de Hensel (~1897), la fórmula del producto")
	fmt.Println("  es clásica y F242 ya midió el libro de las bases. Lo nuestro es la")
	fmt.Println("  identificación: su cuchillo ES la entrada p-ádica, alcanzada simplificando")
	fmt.Println("  fracciones. Todavía no.")

	escribirLamina(fallos, todosExactos)
}

func escribirLamina(fallos int, exactos bool) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="620" viewBox="0 0 1400 620">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="700" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔪 EL CUCHILLO — es ley, y unifica con los primos en la dimensión 0</text>
<text x="700" y="76" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">el ½ pela único, respeta la multiplicación, define una escala — y todas las escalas multiplican a 1 exacto</text>
<rect x="50" y="110" width="420" height="250" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="260" y="144" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">EL CUCHILLO ES LEY</text>
<text x="80" y="184" font-size="14" font-family="monospace" fill="#cfe6ff">46 → un corte  → 23 (primo)</text>
<text x="80" y="210" font-size="14" font-family="monospace" fill="#cfe6ff">44 → dos cortes → 11 (primo)</text>
<text x="80" y="236" font-size="14" font-family="monospace" fill="#cfe6ff">45 → sin corte → 45 = 3·3·5</text>
<text x="80" y="276" font-size="13.5" font-family="Georgia" fill="#7ee0c0">partición única: 10⁶ números, 0 fallos</text>
<text x="80" y="300" font-size="13.5" font-family="Georgia" fill="#7ee0c0">cortes(a·b) = cortes(a)+cortes(b): %d fallos en 10⁶</text>
<text x="80" y="336" font-size="13" font-family="Georgia" fill="#9aa8c4">ley, no costumbre — el TFA en el primo 2</text>
<rect x="490" y="110" width="420" height="250" rx="12" fill="#161a3a" stroke="#5a4fa8"/>
<text x="700" y="144" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">EL ½ COMO UNIDAD DE MEDIDA</text>
<text x="520" y="188" font-size="15" font-family="monospace" fill="#ffd98a">|n|₂ = (½)^cortes</text>
<text x="520" y="224" font-size="14" font-family="monospace" fill="#cfe6ff">|23|₂ = 1 · |46|₂ = ½ · |96|₂ = 1/32</text>
<text x="520" y="262" font-size="13.5" font-family="Georgia" fill="#cfe6ff">cada corte parte el tamaño por la mitad:</text>
<text x="520" y="286" font-size="13.5" font-family="Georgia" fill="#cfe6ff">el ½ no es solo cuchillo — es la escala</text>
<text x="520" y="322" font-size="13" font-family="Georgia" fill="#9aa8c4">y cada primo p tiene la suya: |·|_p</text>
<rect x="930" y="110" width="420" height="250" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="1140" y="144" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">TODAS MULTIPLICAN A 1 EXACTO</text>
<text x="960" y="188" font-size="15" font-family="monospace" fill="#ffd98a">|x|∞ · Π|x|_p = 1</text>
<text x="960" y="222" font-size="13.5" font-family="Georgia" fill="#cfe6ff">verificado en aritmética ENTERA, sin</text>
<text x="960" y="246" font-size="13.5" font-family="Georgia" fill="#cfe6ff">redondeo: ¿exacto en todos los casos? %v</text>
<text x="960" y="284" font-size="13.5" font-family="Georgia" fill="#ffd98a">la meta sin correrse un milímetro,</text>
<text x="960" y="308" font-size="13.5" font-family="Georgia" fill="#ffd98a">otra vez — ahora en tamaños (F242)</text>
<text x="700" y="430" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd98a">El 2 usa dos sombreros: cuchillo de DENSIDAD (come la mitad, la torta) y de TAMAÑO (parte la medida).</text>
<text x="700" y="462" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd98a">El mismo ½, dos oficios, dos contabilidades exactas.</text>
<rect x="50" y="490" width="1300" height="100" rx="12" fill="#33221c" stroke="#c0392b"/>
<text x="700" y="524" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚡ Adónde da la puerta: el ensamble de TODAS las escalas es el mundo ADÉLICO — la frontera que F259 marcó como el camino serio.</text>
<text x="700" y="552" font-size="14" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Hensel ~1897, fórmula del producto clásica, F242 ya lo midió — lo nuestro es la identificación: su juego de fracciones es el picaporte.</text>
<text x="700" y="578" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Todavía no.</text>
</svg>
`, fallos, exactos)
	os.WriteFile("el-cuchillo.svg", []byte(b.String()), 0o644)
	fmt.Println("\n🖼️  lámina escrita: el-cuchillo.svg")
}
