// Command losgemelos judges the captain's mechanic:
//
//	2 + 2 = 4  ->  4 - 1 = 3  and  4 + 1 = 5     both prime
//	3 + 3 = 6  ->  6 - 1 = 5  and  6 + 1 = 7     both prime
//
// Double a number, then look one step to each side. When both sides land on a
// prime you have found a TWIN PRIME PAIR, and the doubled number is its centre.
//
// WHAT IS TRUE, WHAT IS FALSE, AND WHAT IS OPEN - the three answers this
// program separates, because they are three different things:
//
//   - FALSE as a rule: doubling does NOT always give two primes. 4+4 = 8 gives
//     7 and 9, and 9 = 3x3. The mechanic finds twins; it does not manufacture
//     them.
//   - TRUE and provable: every twin pair after (3,5) is centred on a MULTIPLE
//     OF 6. Not most of them - all of them, and the proof is three lines.
//   - OPEN since 1849: whether the twins ever run out. That is the Twin Prime
//     Conjecture, and it is still unsolved. Zhang in 2013 proved some finite gap
//     recurs infinitely often; Maynard and Tao pushed the bound to 246. Nobody
//     has reached 2.
//
// So the captain, from "double it and look at the sides", walked straight into
// one of the famous open problems of number theory. Again.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// criba returns a sieve of primality up to n.
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
	fmt.Println("👯 LOS GEMELOS — la mecánica del capitán, juzgada")
	fmt.Println("\n   lo que dijo: «si tenemos 2 y le sumamos 2 y le restamos 1 tenemos 3, y si le")
	fmt.Println("   sumamos 1 tenemos 5. Si tenemos 3 y le sumamos 3 y le sumamos 1 tenemos 7, y")
	fmt.Println("   si le restamos 1 tenemos 5. Hay una mecánica acá, la encontré».")
	fmt.Println("\n   LA ENCONTRÓ. Y desemboca en un problema abierto desde 1849.")

	const N = 2000000
	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · LA MECÁNICA, ESCRITA")
	fmt.Println("\n        tomás n  →  lo duplicás (2n)  →  mirás un paso a cada lado: 2n−1 y 2n+1")
	fmt.Println("\n   Sus dos ejemplos, verificados:")
	fmt.Println("\n        n     2n     2n−1   2n+1    ¿los dos primos?")
	for _, n := range []int{2, 3} {
		c := 2 * n
		fmt.Printf("   %6d %6d %6d %6d    %s\n", n, c, c-1, c+1,
			map[bool]string{true: "SÍ ← un par de gemelos", false: "no"}[es[c-1] && es[c+1]])
	}
	fmt.Println("\n   Cuando los dos lados caen en primo, eso se llama un PAR DE PRIMOS GEMELOS,")
	fmt.Println("   y el número duplicado es su CENTRO. Usted encontró los dos primeros pares.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚖️ PERO NO SIEMPRE FUNCIONA, Y HAY QUE DECIRLO PRIMERO")
	fmt.Println("   Si la mecánica fabricara primos, el problema estaría resuelto. No los fabrica:")
	fmt.Println("   los ENCUENTRA. Y falla seguido. Mirá qué pasa si seguís con 4, 5, 6…")
	fmt.Println("\n        n     2n     2n−1   2n+1    ¿los dos?     qué falló")
	for n := 2; n <= 12; n++ {
		c := 2 * n
		ok := es[c-1] && es[c+1]
		fallo := ""
		if !ok {
			if !es[c-1] && !es[c+1] {
				fallo = fmt.Sprintf("ninguno: %d y %d compuestos", c-1, c+1)
			} else if !es[c-1] {
				fallo = fmt.Sprintf("%d no es primo", c-1)
			} else {
				fallo = fmt.Sprintf("%d no es primo", c+1)
			}
		}
		fmt.Printf("   %6d %6d %6d %6d    %-10s %s\n", n, c, c-1, c+1,
			map[bool]string{true: "SÍ", false: "no"}[ok], fallo)
	}
	aciertos, total := 0, 0
	for n := 2; 2*n+1 <= N; n++ {
		total++
		if es[2*n-1] && es[2*n+1] {
			aciertos++
		}
	}
	fmt.Printf("\n   → sobre %d valores de n, la mecánica acierta %d veces (%.2f%%).\n",
		total, aciertos, 100*float64(aciertos)/float64(total))
	fmt.Println("     No es una fábrica de primos. Es un DETECTOR de gemelos, y eso ya es mucho.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚡ Y ACÁ ESTÁ LA LEY EXACTA QUE SU MECÁNICA DESTAPA")
	fmt.Println("   Mirá los centros de los pares que sí funcionan:")
	fmt.Println("\n        centro    el par        ¿centro ÷ 6?")
	mostrados := 0
	for c := 4; c <= 200 && mostrados < 9; c += 2 {
		if es[c-1] && es[c+1] {
			mostrados++
			fmt.Printf("   %9d    (%d, %d)%s   %s\n", c, c-1, c+1,
				strings.Repeat(" ", 8-len(fmt.Sprintf("%d, %d", c-1, c+1))),
				map[bool]string{true: "SÍ", false: "NO ← la única excepción"}[c%6 == 0])
		}
	}
	// verificacion masiva
	conCentro6, excepciones := 0, []int{}
	for c := 4; c+1 <= N; c += 2 {
		if es[c-1] && es[c+1] {
			if c%6 == 0 {
				conCentro6++
			} else {
				excepciones = append(excepciones, c)
			}
		}
	}
	fmt.Printf("\n   VERIFICACIÓN MASIVA hasta %d:\n", N)
	fmt.Printf("        pares con centro múltiplo de 6 ...... %d\n", conCentro6)
	fmt.Printf("        excepciones ......................... %d  %v\n", len(excepciones), excepciones)
	fmt.Println("\n   ⟹ **TODO PAR DE GEMELOS, SALVO EL PRIMERO, ESTÁ CENTRADO EN UN MÚLTIPLO DE 6.**")
	fmt.Println("     No la mayoría: TODOS. Y no es una casualidad medida — se demuestra en tres")
	fmt.Println("     renglones, así que vale para siempre y para todos los que faltan.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · POR QUÉ 6, Y NO OTRO NÚMERO")
	fmt.Println("   El centro tiene que ser divisible por 2 y por 3. Las dos razones son simples:")
	fmt.Println("\n     · POR 2 — los dos vecinos tienen que ser impares (si no, son pares y")
	fmt.Println("       divisibles por 2), así que el centro está en el medio de dos impares:")
	fmt.Println("       es par por obligación.")
	fmt.Println("\n     · POR 3 — de tres números seguidos, uno SIEMPRE es múltiplo de 3. Si los")
	fmt.Println("       dos de los costados son primos mayores que 3, el múltiplo de 3 no puede")
	fmt.Println("       ser ninguno de ellos: tiene que ser el del medio.")
	fmt.Println("\n   Y un número divisible por 2 y por 3 es divisible por 6. Listo.")
	fmt.Println("\n   Comprobado: de tres seguidos, ¿siempre hay uno múltiplo de 3?")
	falla3 := 0
	for c := 4; c <= 100000; c++ {
		if (c-1)%3 != 0 && c%3 != 0 && (c+1)%3 != 0 {
			falla3++
		}
	}
	fmt.Printf("        tríos revisados hasta 100000, tríos sin múltiplo de 3: %d\n", falla3)
	fmt.Println("\n   📌 Y ACÁ ESTÁ EL PUENTE CON SU PROPIA INTUICIÓN VIEJA: el 6 = 2 × 3, o sea")
	fmt.Println("     los dos primeros primos multiplicados. Los gemelos viven en los huecos que")
	fmt.Println("     dejan el 2 y el 3, y por eso el 3 y el 5 son el único par que se escapa —")
	fmt.Println("     el 3 todavía no había terminado de nacer cuando se armó ese par.")

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · ⚖️ Y ADÓNDE LO LLEVA ESTO: A UN PROBLEMA ABIERTO DESDE 1849")
	fmt.Println("   La pregunta que sigue naturalmente a su mecánica es:")
	fmt.Println("\n        ¿LOS GEMELOS SE ACABAN ALGUNA VEZ, O SIGUEN PARA SIEMPRE?")
	fmt.Println("\n   Nadie lo sabe. Se llama la CONJETURA DE LOS PRIMOS GEMELOS, la planteó")
	fmt.Println("   Alphonse de Polignac en 1849, y sigue abierta. Lo que sí se sabe:")
	fmt.Println("\n     · 2013 — Yitang Zhang probó que ALGÚN hueco menor que 70.000.000 se repite")
	fmt.Println("       infinitas veces. Fue el primer avance real en 160 años.")
	fmt.Println("     · 2014 — Maynard y Tao bajaron ese techo a 246.")
	fmt.Println("     · Para los gemelos hace falta llegar a 2. Nadie llegó.")
	fmt.Println("\n   Mientras tanto, los gemelos se van RALEANDO. Medido acá:")
	fmt.Println("\n        hasta        pares de gemelos     uno cada tantos números")
	for _, lim := range []int{1000, 10000, 100000, 1000000, N} {
		cnt := 0
		for c := 4; c+1 <= lim; c += 2 {
			if es[c-1] && es[c+1] {
				cnt++
			}
		}
		fmt.Printf("   %11d   %16d   %20.1f\n", lim, cnt, float64(lim)/float64(cnt))
	}
	fmt.Println("\n   → cada vez más raros, pero nunca se detienen del todo en lo que podemos ver.")
	fmt.Println("     Y «nunca se detienen en lo que vemos» no es lo mismo que «no se detienen».")
	fmt.Println("     Es la misma pared de F259, F261 y F263, en otro problema.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("EL CAPITÁN ENCONTRÓ LOS PRIMOS GEMELOS, Y DE PASO SU LEY DEL CENTRO.")
	fmt.Printf("  · su mecánica (duplicar y mirar los lados) acierta %.2f%% de las veces\n",
		100*float64(aciertos)/float64(total))
	fmt.Println("  · NO fabrica primos: los detecta. 4+4 = 8 da 7 y 9, y el 9 no es primo")
	fmt.Printf("  · pero destapa una ley EXACTA: %d de %d pares centrados en múltiplo de 6\n",
		conCentro6, conCentro6+len(excepciones))
	fmt.Printf("  · con %d excepción%s: %v — el par que nació antes de que el 3 terminara\n",
		len(excepciones), map[bool]string{true: "", false: "es"}[len(excepciones) == 1], excepciones)
	fmt.Println("  · y esa ley se demuestra en tres renglones: el centro debe ser divisible")
	fmt.Println("    por 2 (vecinos impares) y por 3 (de tres seguidos, uno lo es)")
	fmt.Println("\n⚖️ Y LO QUE NO SE PUEDE CONCLUIR, dicho igual de fuerte:")
	fmt.Println("  que los gemelos sean infinitos. Eso es la conjetura de 1849 y SIGUE ABIERTA.")
	fmt.Println("  Zhang llegó a 70 millones en 2013, Maynard y Tao a 246. Para los gemelos")
	fmt.Println("  hace falta 2, y nadie llegó. Su mecánica los encuentra; no prueba que sigan.")
	fmt.Println("\n¿El premio? Todavía no — y ojo, que éste es OTRO premio, no el nuestro.")

	escribirLamina(es, N, aciertos, total, conCentro6, excepciones)
}

func escribirLamina(es []bool, N, aciertos, total, conCentro6 int, excepciones []int) {
	var b strings.Builder
	W, H := 1520.0, 1000.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">👯 LOS GEMELOS — la mecánica del capitán, juzgada</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">duplicar y mirar a los costados · encuentra gemelos, destapa una ley exacta, y termina en un problema abierto desde 1849</text>
<text x="%.0f" y="118" font-size="20" text-anchor="middle" font-family="monospace" fill="#ffd98a">n  →  2n  →  2n−1  y  2n+1</text>
`, W, H, W, H, W/2, W/2, W/2)

	// la recta con los gemelos marcados
	fmt.Fprintf(&b, `<rect x="40" y="146" width="1440" height="232" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="760" y="176" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LOS PRIMEROS PARES, Y SUS CENTROS</text>`)
	x0, esc := 90.0, 6.6
	y := 250.0
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="1440" y2="%.0f" stroke="#26456e" stroke-width="1.4"/>`, x0, y, y)
	for c := 4; c <= 200; c += 2 {
		if es[c-1] && es[c+1] {
			px := x0 + float64(c)*esc
			col := "#7ee0c0"
			if c%6 != 0 {
				col = "#ff8fa0"
			}
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="4" fill="%s"/>`, px-esc, y, col)
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="4" fill="%s"/>`, px+esc, y, col)
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="6" fill="none" stroke="%s" stroke-width="1.6"/>`, px, y, col)
			fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="10" text-anchor="middle" font-family="monospace" fill="#8fa8c7">%d</text>`, px, y+22, c)
		}
	}
	fmt.Fprintf(&b, `<text x="760" y="310" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">● los dos primos del par · ○ el centro, que es 2n</text>
<text x="760" y="336" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">en rosa, la única excepción de toda la recta: el par (3, 5), centrado en 4</text>
<text x="760" y="362" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">todos los demás centros son múltiplos de 6</text>`)

	// la ley
	fmt.Fprintf(&b, `<rect x="40" y="398" width="710" height="290" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="395" y="430" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚡ LA LEY DEL CENTRO</text>
<text x="395" y="472" font-size="18" text-anchor="middle" font-family="monospace" fill="#dce8f7">todo centro es múltiplo de 6</text>
<text x="395" y="506" font-size="15" text-anchor="middle" font-family="monospace" fill="#7ee0c0">%d de %d pares · %d excepción</text>
<text x="70" y="546" font-size="14.5" font-family="Georgia" fill="#cfe6ff">Y se demuestra en tres renglones:</text>
<text x="70" y="576" font-size="14" font-family="Georgia" fill="#cfe6ff">· por 2 — los vecinos tienen que ser impares, así que el</text>
<text x="86" y="598" font-size="14" font-family="Georgia" fill="#cfe6ff">centro es par por obligación</text>
<text x="70" y="626" font-size="14" font-family="Georgia" fill="#cfe6ff">· por 3 — de tres seguidos uno SIEMPRE es múltiplo de 3, y</text>
<text x="86" y="648" font-size="14" font-family="Georgia" fill="#cfe6ff">si los costados son primos, tiene que ser el del medio</text>
<text x="395" y="678" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">2 × 3 = 6 · los dos primeros primos</text>
`, conCentro6, conCentro6+len(excepciones), len(excepciones))

	// el problema abierto
	fmt.Fprintf(&b, `<rect x="770" y="398" width="710" height="290" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1125" y="430" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚖️ Y ACÁ EMPIEZA UN PROBLEMA ABIERTO</text>
<text x="1125" y="470" font-size="16" text-anchor="middle" font-family="Georgia" fill="#dce8f7">¿los gemelos se acaban alguna vez?</text>
<text x="1125" y="498" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Conjetura de los Primos Gemelos · de Polignac, 1849</text>
<text x="800" y="540" font-size="14" font-family="Georgia" fill="#cfe6ff">2013 · Zhang: algún hueco &lt; 70.000.000 se repite infinitas</text>
<text x="816" y="562" font-size="14" font-family="Georgia" fill="#cfe6ff">veces — el primer avance en 160 años</text>
<text x="800" y="590" font-size="14" font-family="Georgia" fill="#cfe6ff">2014 · Maynard y Tao bajan el techo a 246</text>
<text x="800" y="620" font-size="14" font-family="Georgia" fill="#ffd98a">para los gemelos hace falta llegar a 2 — nadie llegó</text>
<text x="1125" y="660" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">su mecánica los ENCUENTRA; no prueba que sigan</text>
`)

	fmt.Fprintf(&b, `<rect x="40" y="708" width="1440" height="252" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="760" y="740" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">LO QUE LA MECÁNICA SÍ HACE, Y LO QUE NO</text>
<text x="70" y="778" font-size="15.5" font-family="Georgia" fill="#9fd8a8">✅ ENCUENTRA gemelos, y destapa la ley exacta del centro — que vale para siempre, no solo para lo medido.</text>
<text x="70" y="806" font-size="15.5" font-family="Georgia" fill="#9fd8a8">✅ Y explica por qué el (3, 5) es la única excepción de toda la recta infinita.</text>
<text x="70" y="842" font-size="15.5" font-family="Georgia" fill="#f3d9cf">❌ NO fabrica primos. 4 + 4 = 8 da 7 y 9, y el 9 = 3 × 3. Acierta el %.2f%% de las veces.</text>
<text x="70" y="870" font-size="15.5" font-family="Georgia" fill="#f3d9cf">❌ Y NO prueba que los gemelos sean infinitos: eso sigue abierto desde 1849.</text>
<text x="760" y="912" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Los gemelos se van raleando pero no se detienen en lo que podemos ver — y «no se detienen en lo que vemos»</text>
<text x="760" y="938" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">no es lo mismo que «no se detienen». La misma pared de F259, F261 y F263, en otro problema.</text>
</svg>
`, 100*float64(aciertos)/float64(total))

	if err := os.WriteFile("los-gemelos.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: los-gemelos.svg")
	_ = math.Pi
}
