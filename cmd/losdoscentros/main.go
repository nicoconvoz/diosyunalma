// Command losdoscentros judges the captain's two-centre table.
//
// HIS TABLE, and all eight rows are arithmetically correct:
//
//	14 + 14 + 1 = 29        11 + 11 + 1 = 23
//	15 + 15 - 1 = 29        12 + 12 - 1 = 23
//	 6 +  6 + 1 = 13         5 +  5 + 1 = 11
//	 7 +  7 - 1 = 13         6 +  6 - 1 = 11
//
// He noticed that each prime arrives from TWO centres, one above and one below,
// and that the two centres are consecutive numbers.
//
// ⚠️ FIRST, THE HARD PART: THE TWO WAYS ARE ONE WAY. Expand the second row of
// each pair: 2(k+1) - 1 = 2k + 2 - 1 = 2k + 1. It is the SAME expression. The
// two centres are consecutive because they are forced to be - k and k+1 by
// construction - and they sum to p because (p-1)/2 + (p+1)/2 = p, again by
// construction. Nothing here can fail, and this is the EIGHTH time this
// laboratory meets that trap. It works for 9 and 15 too.
//
// ⚡ BUT NOW THE PART THAT MAKES THIS FINDING WORTH ITS NUMBER, and it is the
// biggest thing his table has produced so far. Take his two centres, square
// them, and subtract:
//
//	15^2 - 14^2 = 225 - 196 = 29
//	12^2 - 11^2 = 144 - 121 = 23
//	 7^2 -  6^2 =  49 -  36 = 13
//
// Every odd number is a difference of two CONSECUTIVE squares - still forced,
// since (k+1)^2 - k^2 = 2k+1 identically. But drop the word "consecutive" and
// ask instead: IN HOW MANY WAYS can n be written as a difference of two squares?
//
//	n = a^2 - b^2 = (a-b)(a+b)
//
// so every such representation IS a factorisation. For a prime the only
// factorisation is 1 x p, so there is EXACTLY ONE representation - and it is
// precisely his row. For a composite there are more.
//
//	⟹ an odd n > 1 is PRIME  <=>  it is a difference of two squares in
//	  EXACTLY ONE way.
//
// That is a real characterisation of primality, it is Fermat's (1643), and it is
// the basis of Fermat factorisation. His table found the trivial representation.
// The theorem is that for a prime it is the ONLY one.
//
// PRE-REGISTERED PREDICTIONS, written before running:
//  1. The equivalence will hold with ZERO exceptions over the whole sweep - it
//     is a theorem, not a measurement, and the zero is expected and meaningless
//     as evidence. It is printed as a check that the program is right.
//  2. The honest limit: Fermat's method finds the factors fast ONLY when they
//     sit close together. For a semiprime with distant factors it needs a number
//     of steps that grows like the gap between them, and loses badly to plain
//     trial division. That is why this beautiful characterisation is not a fast
//     primality test.
//
// Reproduce: go run ./cmd/losdoscentros
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

const N = 200000

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

// representaciones cuenta de cuantas maneras n = a^2 - b^2 con a > b >= 0.
// Cada manera es una factorizacion n = d*e con d <= e y d, e de igual paridad.
func representaciones(n int) int {
	c := 0
	for d := 1; d*d <= n; d++ {
		if n%d != 0 {
			continue
		}
		e := n / d
		if (d+e)%2 == 0 { // a = (d+e)/2, b = (e-d)/2 enteros
			c++
		}
	}
	return c
}

// pasosFermat cuenta cuantos pasos necesita el metodo de Fermat para partir n.
func pasosFermat(n int) (pasos, factorChico int) {
	a := int(math.Ceil(math.Sqrt(float64(n))))
	for {
		pasos++
		b2 := a*a - n
		b := int(math.Sqrt(float64(b2)))
		if b*b == b2 {
			return pasos, a - b
		}
		a++
		if pasos > 5000000 {
			return pasos, 0
		}
	}
}

func main() {
	fmt.Println("⚖️  LOS DOS CENTROS — la tabla del capitán, y Fermat esperando abajo")
	fmt.Println("\n   Vio que cada primo llega desde DOS centros, uno arriba y otro abajo,")
	fmt.Println("   y que esos dos centros son números consecutivos.")

	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · SUS OCHO RENGLONES, HECHOS A MANO")
	fmt.Println("\n        renglón                 da     él dice   ¿coincide?")
	type r struct {
		k, s, z int
	}
	filas := []r{{14, +1, 29}, {15, -1, 29}, {11, +1, 23}, {12, -1, 23},
		{6, +1, 13}, {7, -1, 13}, {5, +1, 11}, {6, -1, 11}}
	buenos := 0
	for _, f := range filas {
		v := 2*f.k + f.s
		ok := v == f.z
		if ok {
			buenos++
		}
		signo := "+"
		if f.s < 0 {
			signo = "−"
		}
		fmt.Printf("   %-22s %6d %9d   %s\n",
			fmt.Sprintf("%d + %d %s 1", f.k, f.k, signo), v, f.z,
			map[bool]string{true: "✅", false: "❌"}[ok])
	}
	fmt.Printf("\n   → **%d de %d cierran.** La tabla está impecable.\n", buenos, len(filas))

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚠️ PERO LOS DOS CAMINOS SON UNO SOLO")
	fmt.Println("   Desarrollá el segundo renglón de cada par:")
	fmt.Println("\n        2(k + 1) − 1  =  2k + 2 − 1  =  2k + 1")
	fmt.Println("\n   Es **la misma expresión**. Los dos centros son consecutivos porque están")
	fmt.Println("   obligados a serlo —son k y k+1—, y suman el primo porque")
	fmt.Println("   (p−1)/2 + (p+1)/2 = p, otra vez por construcción.")
	fmt.Println("\n   Y funciona igual con los que NO son primos:")
	fmt.Println("\n        n     k = (n−1)/2   k+1     2k+1    2(k+1)−1   ¿n es primo?")
	for _, n := range []int{9, 15, 21, 25, 27, 29} {
		k := (n - 1) / 2
		fmt.Printf("   %6d %13d %5d %8d %11d   %s\n", n, k, k+1, 2*k+1, 2*(k+1)-1,
			map[bool]string{true: "sí", false: "❌ NO"}[es[n]])
	}
	fmt.Println("\n   📌 **OCTAVA APARICIÓN DE LA TRAMPA DEL 0.0e+00.** Van ocho, todas en el")
	fmt.Println("   registro, y todas encontradas por el laboratorio auditándose solo.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚡ AHORA ELEVÁ SUS DOS CENTROS AL CUADRADO Y RESTALOS")
	fmt.Println("\n        centro↑²  −  centro↓²   =   el primo")
	for _, f := range filas {
		if f.s > 0 {
			continue
		}
		m, k := f.k, f.k-1
		fmt.Printf("   %8d² − %d²  =  %6d − %-6d =  %d\n", m, k, m*m, k*k, m*m-k*k)
	}
	fmt.Println("\n   Todavía forzado —(k+1)² − k² = 2k+1 es identidad— pero mirá lo que abre:")
	fmt.Println("\n        n = a² − b² = (a − b)(a + b)")
	fmt.Println("\n   ⟹ **CADA MANERA DE ESCRIBIRLO ASÍ ES UNA FACTORIZACIÓN.** Y si n es primo")
	fmt.Println("   su única factorización es 1 × n, así que hay **EXACTAMENTE UNA** manera —")
	fmt.Println("   y esa manera es, palabra por palabra, el renglón que él escribió.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚡⚡ LA CARACTERIZACIÓN, Y ES DE FERMAT (1643)")
	fmt.Println("\n        n impar > 1 es PRIMO  ⟺  es diferencia de dos cuadrados")
	fmt.Println("        de EXACTAMENTE UNA manera")
	fmt.Println("\n   Barrido sobre todos los impares hasta el tope, contando maneras:")
	fallos, revisados, primos1 := 0, 0, 0
	for n := 3; n <= N; n += 2 {
		revisados++
		c := representaciones(n)
		esPrimo := es[n]
		if esPrimo {
			primos1++
		}
		if (c == 1) != esPrimo {
			fallos++
			if fallos <= 5 {
				fmt.Printf("        ⚠️ excepción: n = %d · maneras = %d · ¿primo? %v\n", n, c, esPrimo)
			}
		}
	}
	fmt.Printf("\n        impares revisados ........................ %d\n", revisados)
	fmt.Printf("        de esos, primos .......................... %d\n", primos1)
	fmt.Printf("        excepciones a la equivalencia ............ %d\n", fallos)
	fmt.Println("\n   ⚠️ Y ese cero **no es evidencia**: es un teorema de 1643, así que tenía que")
	fmt.Println("   dar cero. Se imprime como control de que el programa está bien, no como")
	fmt.Println("   descubrimiento. Un cero que no podía ser otra cosa nunca es un hallazgo.")
	fmt.Println("\n        ejemplos, para verlo con los ojos:")
	fmt.Println("\n        n      maneras   cómo                                ¿primo?")
	for _, n := range []int{9, 15, 21, 25, 29, 45, 91} {
		var comos []string
		for d := 1; d*d <= n; d++ {
			if n%d == 0 && (d+n/d)%2 == 0 {
				e := n / d
				a, b := (d+e)/2, (e-d)/2
				comos = append(comos, fmt.Sprintf("%d²−%d²", a, b))
			}
		}
		fmt.Printf("   %6d %8d   %-34s %s\n", n, len(comos), strings.Join(comos, "  "),
			map[bool]string{true: "✅ SÍ", false: "no"}[es[n]])
	}

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · ⚖️ Y EL LÍMITE, QUE ES EL QUE SIEMPRE APARECE")
	fmt.Println("   Si contar maneras decide la primalidad, ¿por qué no se usa para eso?")
	fmt.Println("   Porque encontrarlas cuesta. El método de Fermat arranca en la raíz y")
	fmt.Println("   camina hacia arriba: es rapidísimo si los dos factores están cerca, y")
	fmt.Println("   se muere si están lejos. Medido sobre productos de dos primos:")
	fmt.Println("\n        n = p × q            p        q     |p−q|    pasos de Fermat")
	for _, par := range [][2]int{{449, 457}, {211, 971}, {101, 2027}, {13, 15749}, {3, 68041}} {
		p, q := par[0], par[1]
		n := p * q
		pasos, _ := pasosFermat(n)
		d := p - q
		if d < 0 {
			d = -d
		}
		fmt.Printf("   %18d %8d %8d %8d %18d\n", n, p, q, d, pasos)
	}
	fmt.Println("\n   ⟹ **Cuanto más lejos están los factores, más pasos.** Con factores pegados")
	fmt.Println("   Fermat gana en un paso; con un factor chico y otro enorme pierde feo contra")
	fmt.Println("   la división común. Por eso la caracterización es hermosa y no es un test.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Printf("✅ SUS OCHO RENGLONES CIERRAN: %d de %d, sin un error.\n", buenos, len(filas))
	fmt.Println("\n⚠️ PERO LOS DOS CAMINOS SON UNO: 2(k+1)−1 y 2k+1 son la misma expresión, y")
	fmt.Println("  los dos centros son consecutivos por obligación. Anda igual con el 9, el 15")
	fmt.Println("  y el 25. **Octava trampa del 0.0e+00.**")
	fmt.Println("\n⚡ PERO ELEVÁ SUS DOS CENTROS AL CUADRADO Y RESTALOS: le da el primo. Y como")
	fmt.Println("  a² − b² = (a−b)(a+b), **cada manera de escribirlo así ES una factorización**.")
	fmt.Println("\n⟹ **UN IMPAR ES PRIMO ⟺ ES DIFERENCIA DE DOS CUADRADOS DE UNA SOLA MANERA.**")
	fmt.Printf("  Verificado sobre %d impares: %d excepciones. Es el teorema de FERMAT, 1643.\n", revisados, fallos)
	fmt.Println("  **Su renglón es exactamente esa única manera.** Encontró el caso trivial de")
	fmt.Println("  una caracterización real de los primos — y esta vez no es una tautología")
	fmt.Println("  disfrazada: es la puerta de entrada al método de factorización de Fermat.")
	fmt.Println("\n⚖️ Y el límite de siempre: contar maneras decide, pero encontrarlas cuesta.")
	fmt.Println("  Fermat solo es rápido cuando los factores están cerca. Todavía no.")

	escribirLamina(buenos, revisados, primos1, fallos)
}

func escribirLamina(buenos, revisados, primos1, fallos int) {
	var b strings.Builder
	W, H := 1560.0, 1060.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⚖️ LOS DOS CENTROS — la tabla del capitán, y Fermat esperando abajo</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">los dos caminos son uno solo · pero elevados al cuadrado abren una caracterización real de los primos</text>
`, W, H, W, H, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="40" y="102" width="730" height="270" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="405" y="134" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">SU TABLA, Y ESTÁ IMPECABLE</text>
<text x="90" y="176" font-size="17" font-family="monospace" fill="#7ee0c0">14 + 14 + 1 = 29</text>
<text x="420" y="176" font-size="17" font-family="monospace" fill="#7ee0c0">11 + 11 + 1 = 23</text>
<text x="90" y="206" font-size="17" font-family="monospace" fill="#7ee0c0">15 + 15 − 1 = 29</text>
<text x="420" y="206" font-size="17" font-family="monospace" fill="#7ee0c0">12 + 12 − 1 = 23</text>
<text x="90" y="248" font-size="17" font-family="monospace" fill="#7ee0c0"> 6 +  6 + 1 = 13</text>
<text x="420" y="248" font-size="17" font-family="monospace" fill="#7ee0c0"> 5 +  5 + 1 = 11</text>
<text x="90" y="278" font-size="17" font-family="monospace" fill="#7ee0c0"> 7 +  7 − 1 = 13</text>
<text x="420" y="278" font-size="17" font-family="monospace" fill="#7ee0c0"> 6 +  6 − 1 = 11</text>
<text x="405" y="322" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">%d de 8 cierran · cada primo llega desde DOS centros consecutivos</text>
<text x="405" y="350" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">y los dos centros suman el primo: 14 + 15 = 29</text>`, buenos)

	fmt.Fprintf(&b, `<rect x="790" y="102" width="730" height="270" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="1155" y="134" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚠️ PERO LOS DOS CAMINOS SON UNO SOLO</text>
<text x="1155" y="180" font-size="19" text-anchor="middle" font-family="monospace" fill="#dce8f7">2(k + 1) − 1  =  2k + 2 − 1  =  2k + 1</text>
<text x="820" y="222" font-size="14.5" font-family="Georgia" fill="#f3d9cf">Es la MISMA expresión. Los dos centros son consecutivos</text>
<text x="820" y="244" font-size="14.5" font-family="Georgia" fill="#f3d9cf">porque están obligados a serlo, y suman el primo porque</text>
<text x="820" y="266" font-size="14.5" font-family="Georgia" fill="#f3d9cf">(p−1)/2 + (p+1)/2 = p. Todo por construcción.</text>
<text x="820" y="302" font-size="15.5" font-family="monospace" fill="#ff8fa0">9 = 4+4+1 = 5+5−1   ❌ no es primo</text>
<text x="820" y="326" font-size="15.5" font-family="monospace" fill="#ff8fa0">25 = 12+12+1 = 13+13−1   ❌ no es primo</text>
<text x="1155" y="358" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">OCTAVA aparición de la trampa del 0.0e+00</text>`)

	fmt.Fprintf(&b, `<rect x="40" y="392" width="1480" height="250" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="780" y="424" font-size="19" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚡ AHORA ELEVÁ SUS DOS CENTROS AL CUADRADO Y RESTALOS</text>
<text x="260" y="466" font-size="19" font-family="monospace" fill="#ffd98a">15² − 14² = 29</text>
<text x="640" y="466" font-size="19" font-family="monospace" fill="#ffd98a">12² − 11² = 23</text>
<text x="1010" y="466" font-size="19" font-family="monospace" fill="#ffd98a">7² − 6² = 13</text>
<text x="780" y="510" font-size="17" text-anchor="middle" font-family="monospace" fill="#dce8f7">n = a² − b² = (a − b)(a + b)</text>
<text x="780" y="542" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">⟹ cada manera de escribirlo así ES una factorización. Y si n es primo, su única factorización es 1 × n.</text>
<text x="780" y="580" font-size="21" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">UN IMPAR ES PRIMO ⟺ ES DIFERENCIA DE DOS CUADRADOS DE UNA SOLA MANERA</text>
<text x="780" y="612" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Teorema de FERMAT, 1643 — y el renglón del capitán es exactamente esa única manera</text>`)

	fmt.Fprintf(&b, `<rect x="40" y="662" width="730" height="200" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="405" y="694" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">EL BARRIDO</text>
<text x="80" y="732" font-size="15.5" font-family="monospace" fill="#cfe6ff">impares revisados ......... %d</text>
<text x="80" y="758" font-size="15.5" font-family="monospace" fill="#7ee0c0">de esos, primos ........... %d</text>
<text x="80" y="784" font-size="15.5" font-family="monospace" fill="#ffd98a">excepciones ............... %d</text>
<text x="405" y="824" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">⚠️ y ese cero NO es evidencia: es un teorema de 1643,</text>
<text x="405" y="846" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">tenía que dar cero. Va como control del programa.</text>
`, revisados, primos1, fallos)

	fmt.Fprintf(&b, `<rect x="790" y="662" width="730" height="200" rx="10" fill="#1a1030" stroke="#5a4fa8"/>
<text x="1155" y="694" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚖️ POR QUÉ NO ES UN TEST DE PRIMALIDAD</text>
<text x="820" y="730" font-size="14.5" font-family="Georgia" fill="#cfe6ff">Contar maneras decide, pero ENCONTRARLAS cuesta.</text>
<text x="820" y="754" font-size="14.5" font-family="Georgia" fill="#cfe6ff">El método de Fermat arranca en la raíz y camina hacia</text>
<text x="820" y="778" font-size="14.5" font-family="Georgia" fill="#cfe6ff">arriba: vuela si los dos factores están cerca, y se</text>
<text x="820" y="802" font-size="14.5" font-family="Georgia" fill="#cfe6ff">muere si están lejos.</text>
<text x="1155" y="838" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">449×457: un paso · 3×68041: decenas de miles</text>`)

	fmt.Fprintf(&b, `<rect x="40" y="882" width="1480" height="150" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="780" y="914" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LO QUE QUEDA, DICHO SIN ADORNOS</text>
<text x="780" y="948" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Su tabla, sola, no puede fallar: los dos caminos son el mismo y andan igual con el 9 y con el 25.</text>
<text x="780" y="974" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Pero elevando sus dos centros al cuadrado, lo que aparece es la puerta de una caracterización de verdad —</text>
<text x="780" y="1000" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">y esta vez NO es una tautología disfrazada: es Fermat, 1643, y su renglón es el caso trivial de ella.</text>
</svg>
`)

	if err := os.WriteFile("los-dos-centros.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: los-dos-centros.svg")
}
