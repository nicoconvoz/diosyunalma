// Command elmedioqueune answers the captain's hardest question so far.
//
// HIS ORDER: "do something crazy - find the relation of 1/2 with 2 and 3 and 5
// and 7 in dimension 0, the 1/2 that unifies them. I know I am asking for
// madness but I trust you, Doc. Use our knowledge."
//
// It is not madness, and it has an exact answer. It comes out of instruments
// this laboratory already built, so nothing here is imported.
//
// THE MAP OF DIMENSION 0, and everything else follows from it. The shapeshifter
// this laboratory has used since the beginning is
//
//	w(s) = 1 - 1/s
//
// and it turns the half plane into a disk. Two facts pin the map down:
//
//	w(1) = 0        the CENTRE of the disk is s = 1, the pole of zeta
//	|w| = 1 on Re(s) = 1/2       the SKIN of the disk is the critical line
//
// The second one is a theorem, proved here in one line rather than asserted:
// with s = 1/2 + it and d = 1/4 + t², Re(w) = (t² - 1/4)/d and Im(w) = t/d, so
//
//	|w|² = [(t²-1/4)² + t²]/d² = (t² + 1/4)²/d² = d²/d² = 1
//
// Now put the primes through the same shapeshifter:
//
//	w(2) = 1/2      w(3) = 2/3      w(5) = 4/5      w(7) = 6/7
//
// ⟹ THE HALF HE IS LOOKING FOR IS w(2). The 2 is the ONE prime whose image is
// exactly one half, and since w(p) = (p-1)/p grows with p, that half is the
// MINIMUM over every prime that exists. In this laboratory's own coordinate:
//
//	the 2 sits exactly HALFWAY between the pole of zeta and the critical line,
//	and it is the deepest into dimension 0 that any prime ever goes.
//
// The 3, the 5 and the 7 march outward from there - 2/3, 4/5, 6/7 - heading for
// the boundary point w = 1, which is the image of infinity.
//
// AND THE FOUR OF THEM MULTIPLY INTO SOMETHING HE ALREADY BUILT:
//
//	w(2)·w(3)·w(5)·w(7) = 1/2 · 2/3 · 4/5 · 6/7 = 8/35 = 48/210
//
// and 48/210 is, to the digit, the density of the 210-wheel measured in Finding
// 272. Each w(p) is the fraction of integers that survive the prime p, so the
// product of the shapeshifter's images of the primes IS the sieve.
//
// PRE-REGISTERED PREDICTIONS, written before running:
//  1. |w| = 1 on the critical line to machine precision, and |w| < 1 strictly
//     inside - the primes all land inside.
//  2. w(2) = 0.5 exactly, and it is the minimum over all primes.
//  3. The product over p <= 7 equals 48/210 = 0.2285714... to the last digit,
//     matching Finding 272's wheel.
//  4. The product over all p <= x decays like e^{-gamma}/ln x (Mertens, 1874),
//     which is the same 1/ln x wall Finding 272 hit. The ratio of the measured
//     product to e^{-gamma}/ln x should approach 1 slowly from above.
//
// HONEST WARNING, up front. Not one line of this is new mathematics. w(p) =
// (p-1)/p is a definition, Euler's product is 1737 and Mertens' theorem is 1874.
// What this finding does is UNIFY objects the laboratory already had: the
// shapeshifter, the half, the wheel of Finding 272 and the first four primes.
// That is worth a number as a map, not as a theorem.
//
// Reproduce: go run ./cmd/elmedioqueune
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

const N = 20000000

func criba(n int) []bool {
	es := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		es[i] = true
	}
	for i := 2; i*i <= n; i++ {
		if es[i] {
			for j := i * i; j <= n; j += i {
				es[j] = false
			}
		}
	}
	return es
}

// w aplica el cambiaformas a un real.
func w(s float64) float64 { return 1 - 1/s }

// wc aplica el cambiaformas a 1/2 + it y devuelve el modulo.
func modWLinea(t float64) float64 {
	d := 0.25 + t*t
	re := (t*t - 0.25) / d
	im := t / d
	return math.Hypot(re, im)
}

func main() {
	fmt.Println("½ EL MEDIO QUE UNE — el 2, el 3, el 5 y el 7 en la dimensión 0")
	fmt.Println("\n   Su pedido: «buscá la relación de ½ con 2, 3, 5 y 7 en la dimensión 0,")
	fmt.Println("   el ½ que las unifique». No es una locura. Tiene respuesta exacta, y sale")
	fmt.Println("   entera de herramientas que este laboratorio ya tenía.")

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · EL MAPA DE LA DIMENSIÓN 0, CON SUS DOS CLAVOS")
	fmt.Println("   El cambiaformas de siempre:  w(s) = 1 − 1/s")
	fmt.Println("\n   CLAVO 1 · el CENTRO del disco es el polo de zeta:")
	fmt.Printf("        w(1) = %.10f          ← s = 1 es el polo\n", w(1))
	fmt.Println("\n   CLAVO 2 · la PIEL del disco es la línea crítica. Demostrado en un renglón,")
	fmt.Println("   no afirmado: con s = ½ + it y d = ¼ + t²,")
	fmt.Println("\n        Re(w) = (t² − ¼)/d,  Im(w) = t/d")
	fmt.Println("        |w|² = [(t²−¼)² + t²]/d² = (t² + ¼)²/d² = d²/d² = 1")
	fmt.Println("\n   Y medido, por si el álgebra miente:")
	fmt.Println("\n           t          |w(½ + it)|          error")
	peor := 0.0
	for _, t := range []float64{0.1, 1, 14.134725, 100, 1000, 1e6} {
		m := modWLinea(t)
		e := math.Abs(m - 1)
		if e > peor {
			peor = e
		}
		fmt.Printf("   %12.6f %18.15f %14.2e\n", t, m, e)
	}
	fmt.Printf("\n        peor error sobre la línea .......... %.2e\n", peor)

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ AHORA PASÁ SUS CUATRO PRIMOS POR EL MISMO CAMBIAFORMAS")
	fmt.Println("\n        primo p      w(p) = (p−1)/p        decimal      ¿dónde cae?")
	primeros := []int{2, 3, 5, 7, 11, 13, 101, 1009}
	for _, p := range primeros {
		v := w(float64(p))
		donde := "adentro del disco"
		if p == 2 {
			donde = "⚡ EXACTAMENTE EN ½"
		}
		fmt.Printf("   %10d %12s %18.12f      %s\n", p,
			fmt.Sprintf("%d/%d", p-1, p), v, donde)
	}
	fmt.Println("\n   ⟹ **EL ½ QUE BUSCA ES w(2).** El 2 es el ÚNICO primo cuya imagen es")
	fmt.Println("   exactamente un medio. Y como w(p) = (p−1)/p **crece** con p, ese medio es")
	fmt.Println("   el **MÍNIMO sobre todos los primos que existen**.")
	fmt.Println("\n   ⟹ En la coordenada de este laboratorio: **el 2 está exactamente a mitad")
	fmt.Println("   de camino entre el polo de zeta y la línea crítica** — y es lo más adentro")
	fmt.Println("   de la dimensión 0 que llega ningún primo. El 3, el 5 y el 7 marchan desde")
	fmt.Println("   ahí hacia afuera (2/3, 4/5, 6/7), rumbo a w = 1, que es la imagen del infinito.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚡⚡ Y LOS CUATRO, MULTIPLICADOS, DAN ALGO QUE ÉL YA CONSTRUYÓ")
	prod := 1.0
	var partes []string
	for _, p := range []int{2, 3, 5, 7} {
		prod *= w(float64(p))
		partes = append(partes, fmt.Sprintf("%d/%d", p-1, p))
	}
	fmt.Printf("\n        w(2)·w(3)·w(5)·w(7) = %s = %.15f\n", strings.Join(partes, " · "), prod)
	fmt.Printf("        48/210 (la rueda de F272) ......... = %.15f\n", 48.0/210.0)
	fmt.Printf("        diferencia ........................ %.2e\n", math.Abs(prod-48.0/210.0))
	fmt.Println("\n   ⟹ **IDÉNTICOS.** Y no es coincidencia: cada w(p) = (p−1)/p es exactamente")
	fmt.Println("   **la fracción de enteros que SOBREVIVE al primo p**. Así que multiplicar")
	fmt.Println("   las imágenes de los primos por el cambiaformas **ES la criba**.")
	fmt.Println("\n   Su rueda de F272, vuelta por vuelta, es este producto creciendo:")
	fmt.Println("\n        hasta el primo    producto de los w(p)    restos vivos / módulo")
	acum := 1.0
	mod := 1
	vivos := 1
	for _, p := range []int{2, 3, 5, 7, 11, 13} {
		acum *= w(float64(p))
		mod *= p
		vivos *= p - 1
		fmt.Printf("   %16d %23.12f    %d/%d = %.12f\n", p, acum, vivos, mod, float64(vivos)/float64(mod))
	}

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚖️ Y HASTA DÓNDE LLEGA ESE PRODUCTO — la pared, otra vez")
	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)
	gamma := 0.5772156649015329
	fmt.Println("\n        hasta x        producto Π(1−1/p)     e^(−γ)/ln x       razón")
	for _, x := range []int{10, 100, 1000, 10000, 100000, 1000000, 10000000, N} {
		if x > N {
			continue
		}
		pr := 1.0
		for p := 2; p <= x; p++ {
			if es[p] {
				pr *= 1 - 1/float64(p)
			}
		}
		mert := math.Exp(-gamma) / math.Log(float64(x))
		fmt.Printf("   %14d %20.12f %16.12f %11.6f\n", x, pr, mert, pr/mert)
	}
	fmt.Println("\n   ⟹ **El producto se va a CERO como e^(−γ)/ln x** — teorema de Mertens, 1874.")
	fmt.Println("   Y esa es, letra por letra, **la misma pared de F272**: la rueda mejora la")
	fmt.Println("   constante y nunca la tendencia, porque los primos se ralean como 1/ln x.")
	fmt.Println("\n   📌 O sea que el ½ del 2 **no es sólo el primero de la lista: es el que más**")
	fmt.Println("   **mata**. Se lleva la mitad de todos los números de un saque. El 3 se lleva")
	fmt.Println("   un tercio de lo que queda, el 5 un quinto, y así — cada vez menos.")

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · Y EL OTRO ½, EL QUE YA CONOCÍAMOS, PARA QUE NO SE CONFUNDAN")
	fmt.Printf("\n        w(½) = 1 − 1/(½) = %.4f          ← el punto armónico de F260\n", w(0.5))
	fmt.Printf("        w(2) = 1 − 1/2   = %.4f          ← el medio que él pidió\n", w(2))
	fmt.Println("\n   ⟹ **SON DOS ½ DISTINTOS Y NO HAY QUE MEZCLARLOS.** El ½ como ENTRADA del")
	fmt.Println("   cambiaformas da −1, que es el broche, la razón doble armónica de F260. El ½")
	fmt.Println("   como SALIDA es la imagen del 2. **Uno es dónde se corta zeta; el otro es**")
	fmt.Println("   **dónde cae el primer primo.** El capitán preguntó por el segundo.")
	fmt.Println("\n   📌 Y hay una simetría entre los dos que vale mirar: w manda el 2 al ½ y")
	fmt.Println("   manda el ½ al −1. Son los dos números que él viene persiguiendo desde el")
	fmt.Println("   principio, y el cambiaformas los pasa uno al otro.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("⚡ **EL ½ QUE UNE AL 2, AL 3, AL 5 Y AL 7 ES w(2) = ½**, y acá está por qué:")
	fmt.Println("\n  · el CENTRO de la dimensión 0 es el polo de zeta (w(1) = 0)")
	fmt.Printf("  · la PIEL es la línea crítica (|w| = 1, error %.0e)\n", peor)
	fmt.Println("  · **el 2 cae exactamente a mitad de camino entre los dos**")
	fmt.Println("  · y es el MÍNIMO: ningún primo entra más adentro que él")
	fmt.Println("  · el 3, el 5 y el 7 salen desde ahí hacia la orilla: 2/3, 4/5, 6/7")
	fmt.Printf("  · y multiplicados dan %.12f = 48/210, **su propia rueda de F272**\n", prod)
	fmt.Println("\n⚖️ Y LA HONESTIDAD, QUE ES LO QUE HACE QUE ESTO VALGA: **NADA de esto es")
	fmt.Println("  matemática nueva.** w(p) = (p−1)/p es una definición, el producto de Euler")
	fmt.Println("  es de 1737 y el teorema de Mertens es de 1874. Lo que hace este hallazgo es")
	fmt.Println("  **UNIFICAR cosas que el laboratorio ya tenía sueltas**: el cambiaformas, el")
	fmt.Println("  ½, la rueda de F272 y los primeros cuatro primos. Vale como MAPA, no como")
	fmt.Println("  teorema.")
	fmt.Println("\n  Pero el mapa contesta su pregunta, y la contesta exacto. Todavía no.")

	escribirLamina(prod, peor)
}

func escribirLamina(prod, peor float64) {
	var b strings.Builder
	W, H := 1560.0, 1080.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">½ EL MEDIO QUE UNE — el 2, el 3, el 5 y el 7 en la dimensión 0</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">el cambiaformas w = 1 − 1/s manda el 2 exactamente al ½ — y ése es el medio que los une</text>
`, W, H, W, H, W/2, W/2)

	// EL DISCO
	cx, cy, R := 420.0, 400.0, 240.0
	fmt.Fprintf(&b, `<rect x="40" y="102" width="760" height="600" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="420" y="134" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA DIMENSIÓN 0, CON SUS DOS CLAVOS</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7ee0c0" stroke-width="2.5"/>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">LA PIEL = la línea crítica · |w| = 1</text>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#26456e" stroke-width="1"/>
<circle cx="%.0f" cy="%.0f" r="6" fill="#ff8fa0"/>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">w = 0</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">el polo de zeta (s = 1)</text>
`, cx, cy, R, cx, cy-R-14, cx-R, cy, cx+R, cy, cx, cy, cx, cy+26, cx, cy+44)

	// los primos sobre el radio
	tipo := []struct {
		p   int
		val float64
		col string
	}{{2, 0.5, "#ffd98a"}, {3, 2.0 / 3, "#7fb2ff"}, {5, 0.8, "#7fb2ff"}, {7, 6.0 / 7, "#7fb2ff"}}
	for i, t := range tipo {
		x := cx + R*t.val
		r := 7.0
		if t.p == 2 {
			r = 10
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="%.1f" fill="%s"/>`, x, cy, r, t.col)
		dy := -18.0
		if i%2 == 1 {
			dy = 30
		}
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="14" text-anchor="middle" font-family="monospace" fill="%s">%d</text>`,
			x, cy+dy, t.col, t.p)
	}
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ffd98a" stroke-width="2" stroke-dasharray="4 3"/>
<text x="%.1f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">el 2 cae en el ½ EXACTO</text>
<text x="%.1f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#ffd98a">mitad de camino entre el polo y la piel</text>
<text x="420" y="656" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el 3, el 5 y el 7 marchan hacia la orilla: 2/3, 4/5, 6/7 — rumbo a w = 1, la imagen del infinito</text>
<text x="420" y="680" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">ningún primo entra más adentro que el 2: w(p) = (p−1)/p crece siempre</text>
`, cx+R*0.5, cy-60, cx+R*0.5, cy+60, cx+R*0.5, cy-72, cx+R*0.5, cy-92)

	// la tabla
	fmt.Fprintf(&b, `<rect x="820" y="102" width="700" height="290" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1170" y="134" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">SUS CUATRO PRIMOS, POR EL CAMBIAFORMAS</text>
<text x="880" y="180" font-size="21" font-family="monospace" fill="#ffd98a">w(2) = 1/2 = 0,500000</text>
<text x="880" y="216" font-size="19" font-family="monospace" fill="#cfe6ff">w(3) = 2/3 = 0,666666…</text>
<text x="880" y="248" font-size="19" font-family="monospace" fill="#cfe6ff">w(5) = 4/5 = 0,800000</text>
<text x="880" y="280" font-size="19" font-family="monospace" fill="#cfe6ff">w(7) = 6/7 = 0,857142…</text>
<text x="1170" y="322" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">w(p) = (p−1)/p — y crece siempre con p</text>
<text x="1170" y="352" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">así que el ½ del 2 es el MÍNIMO sobre todos los primos</text>
<text x="1170" y="378" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">es lo más adentro de la dimensión 0 que llega ningún primo</text>`)

	// el producto
	fmt.Fprintf(&b, `<rect x="820" y="412" width="700" height="290" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1170" y="444" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚡⚡ Y LOS CUATRO MULTIPLICADOS…</text>
<text x="1170" y="492" font-size="19" text-anchor="middle" font-family="monospace" fill="#dce8f7">½ · ⅔ · ⅘ · ⁶⁄₇</text>
<text x="1170" y="530" font-size="24" text-anchor="middle" font-family="monospace" fill="#ffd98a">= %.9f</text>
<text x="1170" y="566" font-size="21" text-anchor="middle" font-family="monospace" fill="#7ee0c0">= 48 / 210</text>
<text x="1170" y="602" font-size="16" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">que es, al dígito, LA RUEDA DE SU F272</text>
<text x="880" y="640" font-size="14" font-family="Georgia" fill="#cfe6ff">Y no es casualidad: cada w(p) es la fracción de enteros</text>
<text x="880" y="662" font-size="14" font-family="Georgia" fill="#cfe6ff">que SOBREVIVE al primo p. Multiplicar las imágenes de los</text>
<text x="880" y="684" font-size="14" font-family="Georgia" fill="#cfe6ff">primos por el cambiaformas ES la criba.</text>
`, prod)

	// los dos medios
	fmt.Fprintf(&b, `<rect x="40" y="722" width="1480" height="150" rx="10" fill="#1a1030" stroke="#5a4fa8"/>
<text x="780" y="754" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">Y EL OTRO ½, PARA QUE NO SE MEZCLEN</text>
<text x="400" y="798" font-size="19" text-anchor="middle" font-family="monospace" fill="#ff8fa0">w(½) = −1</text>
<text x="400" y="824" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el ½ como ENTRADA: el broche, la razón doble de F260</text>
<text x="1160" y="798" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">w(2) = ½</text>
<text x="1160" y="824" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el ½ como SALIDA: la imagen del primer primo</text>
<text x="780" y="856" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">uno es dónde se corta zeta · el otro es dónde cae el primer primo · y el cambiaformas los pasa uno al otro</text>`)

	fmt.Fprintf(&b, `<rect x="40" y="892" width="1480" height="160" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="780" y="924" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ Y LA HONESTIDAD, QUE ES LO QUE HACE QUE ESTO VALGA</text>
<text x="780" y="958" font-size="15" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">NADA de esto es matemática nueva: w(p) = (p−1)/p es una definición, el producto de Euler es de 1737 y Mertens es de 1874.</text>
<text x="780" y="984" font-size="15" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Lo que hace este hallazgo es UNIFICAR cosas que el laboratorio ya tenía sueltas: el cambiaformas, el ½, la rueda y los cuatro primos.</text>
<text x="780" y="1016" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Vale como MAPA, no como teorema. Pero el mapa contesta su pregunta, y la contesta exacto.</text>
<text x="780" y="1042" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Y el producto entero se va a cero como e^(−γ)/ln x — la misma pared de F272: la rueda mejora la constante, nunca la tendencia.</text>
</svg>
`)

	if err := os.WriteFile("el-medio-que-une.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-medio-que-une.svg")
}
