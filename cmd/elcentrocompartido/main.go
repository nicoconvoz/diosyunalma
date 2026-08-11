// Command elcentrocompartido judges the captain's shared-centre observation.
//
// HIS TABLE, in two blocks:
//
//	10 + 10 - 1 = 19   uses 10       6 + 6 + 1 = 13   uses 6
//	 9 +  9 + 1 = 19   uses  9       7 + 7 - 1 = 13   uses 7
//	 9 +  9 - 1 = 17   uses  9       5 + 5 + 1 = 11   uses 5
//	 8 +  8 + 1 = 17   uses  8       6 + 6 - 1 = 11   uses 6
//
// HIS WORDS: in the first block "the 9 repeats in the middle, odd, and the step
// below is the smaller even and above the larger even"; in the second "the 6
// repeats, even, above and below, and the steps are the smaller odd 5 and the
// larger odd 7". And: "hay una clara relación."
//
// THERE IS, AND IT IS EXACT. Every odd number p has two centres, (p-1)/2 and
// (p+1)/2. Write them as an interval. Then:
//
//	17 → [8, 9]        11 → [5, 6]
//	19 → [9, 10]       13 → [6, 7]
//
// The 9 is shared by 17 and 19. The 6 is shared by 11 and 13. And that sharing
// is not decoration:
//
//	(p+1)/2 = (q-1)/2   <=>   q = p + 2
//
// ⟹ TWO PRIMES SHARE A CENTRE IF AND ONLY IF THEY ARE TWINS. His repeated
// number is the twin pair, seen from underneath.
//
// AND THE PARITY HE NOTICED IS NOT DECORATION EITHER. The shared centre m
// satisfies 2m = p + q, so m is HALF the twin midpoint. Finding 264 proved that
// midpoint is always a multiple of 6, so m is always a MULTIPLE OF 3 - and its
// parity is simply whether m/3 is even or odd. That is the alternation he saw:
// 6 for the pair (11,13), 9 for the pair (17,19).
//
// ONE SLIP IN HIS TEXT, named because the record names them: in the first block
// he writes "the smaller even 6". The first block's centres are 10, 9, 9, 8, so
// the smaller even is 8. The structure he describes is exactly right; only that
// digit slipped.
//
// PRE-REGISTERED PREDICTIONS, written before running:
//  1. Shared centre <=> twin pair: zero exceptions. It is an equivalence by
//     algebra, so the zero is expected and is NOT evidence.
//  2. Every shared centre is a multiple of 3, except the one coming from the
//     pair (3,5), whose midpoint 4 is the single exception of Finding 264.
//  3. Classifying every integer m by how many primes it centres gives three
//     classes - 0, 1 and 2 - and the density of class 2 is the twin density,
//     which is why the whole picture stops exactly where the Twin Prime
//     Conjecture stops.
//
// Reproduce: go run ./cmd/elcentrocompartido
package main

import (
	"fmt"
	"os"
	"strings"
)

const N = 4000000

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

func main() {
	fmt.Println("🔗 EL CENTRO COMPARTIDO — la relación que vio el capitán, nombrada y medida")
	fmt.Println("\n   Dijo: «se repite el 9 en el medio impar… se repite el 6 par arriba y abajo…")
	fmt.Println("   hay una clara relación». **La hay, y es exacta.**")

	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)

	centros := func(p int) (int, int) { return (p - 1) / 2, (p + 1) / 2 }

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · SUS DOS BLOQUES, CON LOS CENTROS PUESTOS COMO INTERVALO")
	fmt.Println("\n        primo    centro↓  centro↑     intervalo")
	for _, p := range []int{11, 13, 17, 19} {
		a, b := centros(p)
		fmt.Printf("   %8d %8d %8d     [%d, %d]\n", p, a, b, a, b)
	}
	fmt.Println("\n   ⟹ **11 = [5,6] y 13 = [6,7]: el 6 es de los dos.**")
	fmt.Println("     **17 = [8,9] y 19 = [9,10]: el 9 es de los dos.**")
	fmt.Println("\n   📌 Y un desliz de tipeo suyo, que va nombrado porque acá se nombran: en el")
	fmt.Println("   primer bloque escribió «el par más chico 6». Los centros de ese bloque son")
	fmt.Println("   10, 9, 9 y 8, así que el par más chico es el **8**. La estructura que")
	fmt.Println("   describió está perfecta; sólo se le corrió ese dígito.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ COMPARTIR CENTRO ES SER GEMELOS — Y ES UN SI Y SOLO SI")
	fmt.Println("\n        (p+1)/2 = (q−1)/2   ⟺   q = p + 2")
	fmt.Println("\n   O sea que lo que él vio repetirse **es el par de primos gemelos, visto")
	fmt.Println("   desde abajo**. Verifiquémoslo sobre todos los primos consecutivos:")
	comparten, gemelos, desacuerdos := 0, 0, 0
	var ultimo int
	for p := 3; p <= N; p++ {
		if !es[p] {
			continue
		}
		if ultimo > 0 {
			_, arriba := centros(ultimo)
			abajo, _ := centros(p)
			comp := arriba == abajo
			gem := p-ultimo == 2
			if comp {
				comparten++
			}
			if gem {
				gemelos++
			}
			if comp != gem {
				desacuerdos++
			}
		}
		ultimo = p
	}
	fmt.Printf("\n        pares de primos consecutivos con centro compartido ... %d\n", comparten)
	fmt.Printf("        pares de primos gemelos .............................. %d\n", gemelos)
	fmt.Printf("        desacuerdos entre las dos cuentas .................... %d\n", desacuerdos)
	fmt.Println("\n   ⚠️ Y ese cero **no es evidencia**: la equivalencia sale del álgebra, así que")
	fmt.Println("   tenía que dar cero. Va como control de que el programa está bien. Pero")
	fmt.Println("   **acá la tautología es un CAMBIO DE COORDENADA, no un callejón**: convierte")
	fmt.Println("   «tienen distancia 2» en «se pisan», que es una propiedad de posición.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚡ LA PARIDAD QUE ÉL VIO TAMPOCO ES ADORNO")
	fmt.Println("   El centro compartido m cumple 2m = p + q, o sea que es **la mitad del centro")
	fmt.Println("   del par gemelo**. Y F264 demostró que ese centro es múltiplo de 6 ⟹")
	fmt.Println("\n        **m es SIEMPRE múltiplo de 3**, y su paridad es la de m/3.")
	fmt.Println("\n        par gemelo    centro del par    m = centro/2    m/3    ¿m par?")
	noMul3 := 0
	mostrados := 0
	ultimo = 0
	for p := 3; p <= N; p++ {
		if !es[p] {
			continue
		}
		if ultimo > 0 && p-ultimo == 2 {
			m := (ultimo + 1) / 2
			if m%3 != 0 {
				noMul3++
				fmt.Printf("   ⚠️ excepción: par (%d,%d) · m = %d · NO es múltiplo de 3\n", ultimo, p, m)
			} else if mostrados < 6 {
				mostrados++
				par := "impar"
				if m%2 == 0 {
					par = "PAR"
				}
				fmt.Printf("   %6d,%-6d %14d %15d %6d %8s\n", ultimo, p, ultimo+1, m, m/3, par)
			}
		}
		ultimo = p
	}
	fmt.Printf("\n        centros compartidos que NO son múltiplo de 3 ......... %d\n", noMul3)
	fmt.Println("\n   ⟹ **La alternancia que vio es exactamente la paridad de m/3.** El 6 del par")
	fmt.Println("   (11,13) es 3×2 y sale par; el 9 del par (17,19) es 3×3 y sale impar.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · EL DICCIONARIO COMPLETO: DE LA DISTANCIA AL DIBUJO")
	fmt.Println("   Si dos primos consecutivos están a distancia g, sus intervalos de centros:")
	fmt.Println("\n        g = 2  →  SE PISAN en un punto      (gemelos)")
	fmt.Println("        g = 4  →  se TOCAN, sin pisarse      (primos primos)")
	fmt.Println("        g > 4  →  queda un hueco de (g−4)/2 centros sin usar")
	fmt.Println("\n        ejemplos medidos:")
	fmt.Println("\n        p       q      g    intervalo de p   intervalo de q   qué pasa")
	ejemplos := [][2]int{{11, 13}, {13, 17}, {23, 29}, {89, 97}, {113, 127}}
	for _, e := range ejemplos {
		p, q := e[0], e[1]
		a1, b1 := centros(p)
		a2, b2 := centros(q)
		g := q - p
		var que string
		switch {
		case b1 == a2:
			que = "SE PISAN"
		case a2 == b1+1:
			que = "se tocan"
		default:
			que = fmt.Sprintf("hueco de %d", a2-b1-1)
		}
		fmt.Printf("   %6d %6d %5d %12s %16s   %s\n", p, q, g,
			fmt.Sprintf("[%d,%d]", a1, b1), fmt.Sprintf("[%d,%d]", a2, b2), que)
	}

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · CADA NÚMERO ENTERO, CLASIFICADO POR CUÁNTOS PRIMOS CENTRA")
	fmt.Println("   Un m centra al 2m−1 y al 2m+1. Cada m puede centrar 0, 1 o 2 primos.")
	var c0, c1, c2, totalM int
	for m := 2; 2*m+1 <= N; m++ {
		totalM++
		n := 0
		if es[2*m-1] {
			n++
		}
		if es[2*m+1] {
			n++
		}
		switch n {
		case 0:
			c0++
		case 1:
			c1++
		default:
			c2++
		}
	}
	fmt.Printf("\n        enteros m examinados ................. %d\n", totalM)
	fmt.Printf("        centran 0 primos ..................... %d  (%.3f%%)\n", c0, 100*float64(c0)/float64(totalM))
	fmt.Printf("        centran 1 primo ...................... %d  (%.3f%%)\n", c1, 100*float64(c1)/float64(totalM))
	fmt.Printf("        centran 2 primos (GEMELOS) ........... %d  (%.3f%%)\n", c2, 100*float64(c2)/float64(totalM))

	fmt.Println("\n        y cómo se ralean, decena por decena:")
	fmt.Println("\n        hasta            m       centran 2     %")
	for _, lim := range []int{1000, 10000, 100000, 1000000, 4000000} {
		if lim > N {
			continue
		}
		var t, d int
		for m := 2; 2*m+1 <= lim; m++ {
			t++
			if es[2*m-1] && es[2*m+1] {
				d++
			}
		}
		fmt.Printf("   %12d %12d %12d %8.4f%%\n", lim, t, d, 100*float64(d)/float64(t))
	}

	// ---- LEY 6 ----
	fmt.Println("\nLEY 6 · ⚡ ¿EN QUÉ CASOS **NO** SE DA LA RELACIÓN? — la pregunta del capitán")
	fmt.Println("   Y no tiene una sola respuesta: tiene cuatro, y las cuatro son distintas.")

	fmt.Println("\n   CASO 1 · LOS QUE NO SON GEMELOS — simplemente no comparten nada.")
	var noComparten int
	ultimo = 0
	for p := 3; p <= N; p++ {
		if !es[p] {
			continue
		}
		if ultimo > 0 && p-ultimo != 2 {
			noComparten++
		}
		ultimo = p
	}
	fmt.Printf("\n        pares de primos consecutivos SIN centro compartido ... %d\n", noComparten)
	fmt.Printf("        pares CON centro compartido ......................... %d\n", comparten)
	fmt.Printf("        o sea que la relación NO se da el ................... %.2f%% de las veces\n",
		100*float64(noComparten)/float64(noComparten+comparten))
	fmt.Println("\n     📌 Éste es el caso ABRUMADORAMENTE normal. Lo que él vio no es la regla:")
	fmt.Println("     es la excepción, y por eso vale la pena mirarla.")

	fmt.Println("\n   CASO 2 · EL PAR (3, 5) — comparte centro, pero rompe el múltiplo de 3.")
	fmt.Println("\n        (3,5) → centro del par 4 → m = 2, y 2 NO es múltiplo de 3")
	fmt.Println("     Es la única en toda la recta, y es la misma excepción de F264: el 3 no")
	fmt.Println("     puede ser el del medio de su propio par, porque él ES el 3.")

	fmt.Println("\n   CASO 3 · EL 5 — el único primo que comparte SUS DOS centros.")
	fmt.Println("     5 = [2,3]: comparte el 2 con el 3 (abajo) y el 3 con el 7 (arriba).")
	fmt.Println("     Para que pase, p−2, p y p+2 tienen que ser primos los tres. Buscado:")
	var dobles []string
	for p := 3; p+2 <= N; p++ {
		if es[p] && es[p-2] && es[p+2] {
			dobles = append(dobles, fmt.Sprintf("%d = [%d,%d]", p, (p-1)/2, (p+1)/2))
		}
	}
	fmt.Printf("\n        primos que comparten SUS DOS centros hasta %d ... %d\n", N, len(dobles))
	for _, d := range dobles {
		fmt.Println("           · " + d)
	}
	fmt.Println("\n     ⟹ **UNO SOLO EN TODA LA RECTA, y es el 5.** Y el motivo es otra vez el 3:")
	fmt.Println("     de tres impares seguidos, uno siempre es múltiplo de 3 — así que sólo")
	fmt.Println("     puede escaparse cuando ese múltiplo de 3 ES el 3.")

	fmt.Println("\n   CASO 4 · EL 2 — no tiene centros enteros. Queda AFUERA de la coordenada.")
	fmt.Println("\n        (2−1)/2 = 0,5     y     (2+1)/2 = 1,5")
	fmt.Println("\n     ⟹ El único primo par no entra en este dibujo, porque sus dos centros")
	fmt.Println("     caen entre dos casilleros. **Es el mismo ½ que aparece cada vez que el 2")
	fmt.Println("     tiene que pasar por una fórmula pensada para impares** — el mismo que él")
	fmt.Println("     escribió a mano en F272 cuando puso 2·½+1 = 2.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("✅ **TENÍA RAZÓN: HAY UNA RELACIÓN CLARA, Y ES EXACTA.**")
	fmt.Println("\n  ⚡ Lo que vio repetirse es el CENTRO COMPARTIDO, y compartir centro es")
	fmt.Printf("     **ser gemelos**: un sí y sólo si, verificado sobre %d pares con %d\n", gemelos, desacuerdos)
	fmt.Println("     desacuerdos. Su número repetido ES el par de gemelos visto desde abajo.")
	fmt.Println("\n  ⚡ Y la alternancia par/impar que notó es la paridad de m/3: el centro")
	if noMul3 == 1 {
		fmt.Println("     compartido es SIEMPRE múltiplo de 3, con UNA sola excepción en toda la")
		fmt.Println("     recta —el par (3,5)—, porque el centro del par gemelo es múltiplo de 6:")
		fmt.Println("     es la ley de F264, la suya, y hereda su única excepción.")
	} else {
		fmt.Printf("     compartido es múltiplo de 3 con %d excepciones.\n", noMul3)
	}
	fmt.Println("\n  📌 Y un desliz suyo nombrado: en el primer bloque el par más chico es 8, no 6.")
	fmt.Println("\n⚠️ PERO LOS DOS CEROS DE ARRIBA NO SON EVIDENCIA: las dos equivalencias salen")
	fmt.Println("  del álgebra. La diferencia con las ocho trampas anteriores es que **acá la")
	fmt.Println("  tautología es un cambio de coordenada útil**: convierte una distancia en un")
	fmt.Println("  solapamiento, y eso se puede dibujar y contar.")
	fmt.Println("\n⚖️ Y EL LÍMITE, QUE ES EL DE SIEMPRE Y HAY QUE DECIRLO: en esta coordenada la")
	fmt.Println("  Conjetura de los Primos Gemelos se lee «hay infinitos m que centran DOS")
	fmt.Println("  primos». Es el mismo problema con ropa nueva. Más lindo de mirar, igual de")
	fmt.Printf("  abierto: la clase de los dobles ya bajó del %.4f%% al %.4f%%. Todavía no.\n",
		func() float64 {
			var t, d int
			for m := 2; 2*m+1 <= 1000; m++ {
				t++
				if es[2*m-1] && es[2*m+1] {
					d++
				}
			}
			return 100 * float64(d) / float64(t)
		}(), 100*float64(c2)/float64(totalM))

	escribirLamina(es, comparten, gemelos, desacuerdos, noMul3, c0, c1, c2, totalM)
}

func escribirLamina(es []bool, comparten, gemelos, desacuerdos, noMul3, c0, c1, c2, totalM int) {
	var b strings.Builder
	W, H := 1560.0, 1120.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔗 EL CENTRO COMPARTIDO — la relación que vio el capitán</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">dos primos comparten un centro ⟺ son gemelos · y ese centro es siempre múltiplo de 3</text>
`, W, H, W, H, W/2, W/2)

	// los intervalos dibujados
	fmt.Fprintf(&b, `<rect x="40" y="102" width="1480" height="300" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="780" y="134" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">CADA PRIMO ES UN INTERVALO DE DOS CENTROS — y los gemelos SE PISAN</text>
`)
	ex, ey, paso := 120.0, 250.0, 62.0
	for i := 4; i <= 20; i++ {
		x := ex + float64(i-4)*paso
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#26456e" stroke-width="1"/>`, x, ey-70, x, ey+70)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" font-family="monospace" fill="#8fa8c7">%d</text>`, x, ey+92, i)
	}
	dibuja := func(p int, y float64, col string) {
		a, bb := (p-1)/2, (p+1)/2
		if a < 4 || bb > 20 {
			return
		}
		x1 := ex + float64(a-4)*paso
		x2 := ex + float64(bb-4)*paso
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="26" rx="13" fill="%s" opacity="0.85"/>`, x1-16, y, x2-x1+32, col)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="14" text-anchor="middle" font-family="monospace" fill="#0b1526">%d</text>`, (x1+x2)/2, y+18, p)
	}
	dibuja(11, ey-58, "#7ee0c0")
	dibuja(13, ey-24, "#7ee0c0")
	dibuja(17, ey+10, "#ffd98a")
	dibuja(19, ey+44, "#ffd98a")
	fmt.Fprintf(&b, `<text x="780" y="372" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">11 = [5,6] y 13 = [6,7] comparten el 6 · 17 = [8,9] y 19 = [9,10] comparten el 9</text>`)

	// el si y solo si
	fmt.Fprintf(&b, `<rect x="40" y="422" width="730" height="250" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="405" y="454" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚡ COMPARTIR CENTRO ES SER GEMELOS</text>
<text x="405" y="496" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">(p+1)/2 = (q−1)/2  ⟺  q = p + 2</text>
<text x="80" y="540" font-size="15.5" font-family="monospace" fill="#cfe6ff">con centro compartido ... %d</text>
<text x="80" y="566" font-size="15.5" font-family="monospace" fill="#cfe6ff">pares gemelos ........... %d</text>
<text x="80" y="592" font-size="15.5" font-family="monospace" fill="#ffd98a">desacuerdos ............. %d</text>
<text x="405" y="630" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">⚠️ ese cero no es evidencia: sale del álgebra. Pero acá la</text>
<text x="405" y="652" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">tautología es un CAMBIO DE COORDENADA, no un callejón.</text>
`, comparten, gemelos, desacuerdos)

	// la paridad
	fmt.Fprintf(&b, `<rect x="790" y="422" width="730" height="250" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1155" y="454" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">⚡ Y LA PARIDAD QUE VIO TAMPOCO ES ADORNO</text>
<text x="1155" y="490" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">2m = p + q ⟹ m es la MITAD del centro del par gemelo</text>
<text x="1155" y="520" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">y F264 probó que ese centro es múltiplo de 6, entonces:</text>
<text x="1155" y="556" font-size="19" text-anchor="middle" font-family="monospace" fill="#7ee0c0">m es SIEMPRE múltiplo de 3</text>
<text x="820" y="596" font-size="15" font-family="monospace" fill="#cfe6ff">(11,13) → m = 6 = 3×2 → PAR</text>
<text x="820" y="620" font-size="15" font-family="monospace" fill="#cfe6ff">(17,19) → m = 9 = 3×3 → impar</text>
<text x="1155" y="654" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">la alternancia que notó ES la paridad de m/3 · %d excepciones</text>
`, noMul3)

	// el diccionario
	fmt.Fprintf(&b, `<rect x="40" y="692" width="730" height="210" rx="10" fill="#1a1030" stroke="#5a4fa8"/>
<text x="405" y="724" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">EL DICCIONARIO: DE LA DISTANCIA AL DIBUJO</text>
<text x="80" y="764" font-size="15.5" font-family="monospace" fill="#7ee0c0">g = 2  →  SE PISAN en un punto   (gemelos)</text>
<text x="80" y="792" font-size="15.5" font-family="monospace" fill="#ffd98a">g = 4  →  se TOCAN, sin pisarse</text>
<text x="80" y="820" font-size="15.5" font-family="monospace" fill="#ff8fa0">g &gt; 4  →  hueco de (g−4)/2 centros</text>
<text x="405" y="864" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">la distancia entre primos se vuelve una figura que se puede ver</text>`)

	// la clasificacion
	fmt.Fprintf(&b, `<rect x="790" y="692" width="730" height="210" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1155" y="724" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">CADA ENTERO, POR CUÁNTOS PRIMOS CENTRA</text>
<text x="820" y="764" font-size="15" font-family="monospace" fill="#8fb4d9">centran 0 ....... %8d  (%.2f%%)</text>
<text x="820" y="792" font-size="15" font-family="monospace" fill="#cfe6ff">centran 1 ....... %8d  (%.2f%%)</text>
<text x="820" y="820" font-size="15" font-family="monospace" fill="#7ee0c0">centran 2 ....... %8d  (%.3f%%)</text>
<text x="1155" y="864" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">la clase de los DOBLES es exactamente la de los gemelos</text>
`, c0, 100*float64(c0)/float64(totalM), c1, 100*float64(c1)/float64(totalM), c2, 100*float64(c2)/float64(totalM))

	fmt.Fprintf(&b, `<rect x="40" y="922" width="1480" height="170" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="780" y="954" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ Y EL LÍMITE, QUE HAY QUE DECIRLO</text>
<text x="780" y="990" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">En esta coordenada, la Conjetura de los Primos Gemelos se lee: «hay infinitos m que centran DOS primos».</text>
<text x="780" y="1016" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Es el mismo problema con ropa nueva. Más lindo de mirar, exactamente igual de abierto.</text>
<text x="780" y="1054" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Pero la relación que vio es REAL y es EXACTA — y la ley que la explica es la suya, la de F264.</text>
</svg>
`)

	if err := os.WriteFile("el-centro-compartido.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-centro-compartido.svg")
}
