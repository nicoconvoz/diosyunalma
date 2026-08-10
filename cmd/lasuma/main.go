// Command lasuma judges the captain's table and the formula he drew from it.
//
// THE TABLE HE WROTE
//
//	0 + 0 + 0 = 0 - 0 - 0 = 0        identity
//	1/2 + 1/2 = 1 + 0   = 1          unique property
//	1/2 - 1/2 = 1/4 - 0 = 1/4        <- THIS ROW IS WRONG: 1/2 - 1/2 = 0
//	1 + 1 = 2                        first prime
//	1 - 1 = 0                        unique property
//	2 + 2 - 1 = 3      2 + 2 + 1 = 5
//	3 + 3 - 1 = 5      3 + 3 + 1 = 7
//	5 + 5 - 3 = 7      5 + 5 + 3 = 13
//
// and he concluded: (X + X) - 2Y = Z.
//
// THREE CORRECTIONS AND ONE IDENTITY.
//
// The third row is arithmetically false: a half minus a half is zero, not a
// quarter. And the formula carries a factor-two slip - his own rows say
// 2X ± Y = Z, not 2X - 2Y. Substituting his numbers proves it: 2*5 - 2*3 = 4,
// not the 7 the row gives. That second correction is not cosmetic: 2X - 2Y is
// 2(X-Y), always even, so without striking the spurious 2 nothing on the left
// could ever have been prime past 2 itself.
//
// But turn the corrected formula around and something famous appears. If
// P = 2X - Y and Q = 2X + Y are both prime, then adding them kills Y:
//
//	P + Q = 4X
//
// That cancellation is an IDENTITY. It holds for any X and Y, prime or not, so
// it can never fail and proves nothing on its own. What it buys is a coordinate:
// fix the centre, walk the offset outward. That is how every Goldbach
// verification in history is actually run.
//
// So the LAST THREE rows of his table are even numbers written as sums of two
// primes: 8 = 3+5, 12 = 5+7, 20 = 7+13. And here is the THIRD correction, the
// only one of the three with mathematics in it: with X a whole number, 4X only
// ever reaches multiples of 4, so his family is Goldbach restricted to
// n = 0 (mod 4). Free the centre 2X to be ANY integer >= 2 and it becomes
// exactly the GOLDBACH CONJECTURE - stated in the Goldbach-Euler correspondence
// of 1742 (Goldbach's letter of 7 June gives the ternary form; the binary
// statement is Euler's reply of 30 June), verified by computer up to 4x10^18,
// proved by nobody. This program sweeps the FREED centre, which is why it covers
// every even number and not only the multiples of 4.
//
// F264 and F265 came out of the SAME table on the SAME day: F264 is the centre
// with the offset held at 1, F265 is the same centre with the offset set free.
package main

import (
	"fmt"
	"os"
	"strings"
)

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
	fmt.Println("➕ LA SUMA DE DOS — la tabla del capitán, revisada renglón por renglón")
	fmt.Println("\n   Me pasó una tabla y sacó de ella una fórmula: (X + X) − 2Y = Z.")
	fmt.Println("   Hay DOS correcciones que hacerle… y abajo de las dos, un hallazgo grande.")

	// el índice más grande que se lee es c+y < 2c ≤ 400000, así que la criba
	// se dimensiona exactamente para eso y no se infla el número que se imprime.
	const N = 400000
	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · LA TABLA, RENGLÓN POR RENGLÓN")
	fmt.Println("\n        renglón                        ¿da?      qué pasa")
	fmt.Printf("   %-30s %-9s %s\n", "0 + 0 + 0 = 0 − 0 − 0 = 0", "✅ SÍ", "el cero es el único que no se mueve")
	fmt.Printf("   %-30s %-9s %s\n", "½ + ½ = 1 + 0 = 1", "✅ SÍ", "las dos mitades hacen el uno")
	fmt.Printf("   %-30s %-9s %s\n", "½ − ½ = ¼ − 0 = ¼", "❌ NO", "½ − ½ = 0, no ¼")
	fmt.Printf("   %-30s %-9s %s\n", "1 + 1 = 2", "✅ SÍ", "y sí: el primer primo")
	fmt.Printf("   %-30s %-9s %s\n", "1 − 1 = 0", "✅ SÍ", "vuelve al cero")
	for _, r := range []struct{ x, y, m, p int }{{2, 1, 3, 5}, {3, 1, 5, 7}, {5, 3, 7, 13}} {
		ok := 2*r.x-r.y == r.m && 2*r.x+r.y == r.p && es[r.m] && es[r.p]
		fmt.Printf("   %-30s %-9s %s\n",
			fmt.Sprintf("%d+%d−%d = %d  y  +%d = %d", r.x, r.x, r.y, r.m, r.y, r.p),
			map[bool]string{true: "✅ SÍ", false: "no"}[ok], "los dos resultados son primos")
	}
	fmt.Println("\n   📌 EL ÚNICO RENGLÓN FALSO: **½ − ½ = 0**, no ¼. Ojo que el ¼ sí aparece en")
	fmt.Println("   otro lado — es ½ × ½, la mitad DE la mitad. Restar y partir no son lo mismo:")
	fmt.Println("   restar te devuelve al cero, partir te lleva más adentro. Su tabla trae ONCE")
	fmt.Println("   igualdades — acá los tres últimos pares van de a dos por línea — así que los")
	fmt.Println("   DIEZ renglones restantes están perfectos.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚖️ Y LA FÓRMULA TIENE UN DOS DE MÁS")
	fmt.Println("   Escribió  (X + X) − 2Y = Z.  Probémosla con sus propios renglones:")
	fmt.Println("\n        renglón          su fórmula: 2X − 2Y     lo que el renglón dice")
	for _, r := range []struct{ x, y, z int }{{2, 1, 3}, {3, 1, 5}, {5, 3, 7}} {
		fmt.Printf("   %-16s %-23d %d\n",
			fmt.Sprintf("%d+%d−%d = %d", r.x, r.x, r.y, r.z), 2*r.x-2*r.y, r.z)
	}
	fmt.Println("\n   → no coincide. Lo que sus renglones dicen de verdad es, sin el dos:")
	fmt.Println("\n        (X + X) − Y = Z        y        (X + X) + Y = Z′")
	fmt.Println("\n   Es un desliz de escritura, no de idea — la idea está bien y ahora lo vemos.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚡ Y AHORA DÉ VUELTA LA FÓRMULA, QUE ACÁ ESTÁ LO GRANDE")
	fmt.Println("   Tomemos sus dos resultados juntos y sumémoslos:")
	fmt.Println("\n        P = 2X − Y")
	fmt.Println("        Q = 2X + Y")
	fmt.Println("        ─────────────")
	fmt.Println("        P + Q = 4X        ← la Y SE BORRA")
	fmt.Println("\n   Y como P y Q son primos, cada renglón suyo es UN NÚMERO PAR ESCRITO COMO")
	fmt.Println("   SUMA DE DOS PRIMOS. Sus tres renglones, dados vuelta:")
	fmt.Println("\n        X     Y     P      Q      P + Q = 4X")
	for _, r := range []struct{ x, y int }{{2, 1}, {3, 1}, {5, 3}} {
		p, q := 2*r.x-r.y, 2*r.x+r.y
		fmt.Printf("   %5d %5d %5d %6d      %d + %d = %d\n", r.x, r.y, p, q, p, q, p+q)
	}
	fmt.Println("\n   ⟹ 8 = 3+5 · 12 = 5+7 · 20 = 7+13. **Usted escribió tres descomposiciones**")
	fmt.Println("     **de Goldbach sin saber que se llamaban así.**")
	fmt.Println("\n   📌 PERO OJO, Y ESTO HAY QUE DECIRLO: P + Q = 4X es una IDENTIDAD. Vale para")
	fmt.Println("   cualquier X y cualquier Y, primos o no — no puede fallar, y por lo tanto sola")
	fmt.Println("   no demuestra NADA. Lo que compra no es verdad: es una COORDENADA. Fijás el")
	fmt.Println("   centro y caminás el desvío hacia afuera. Así se corre toda verificación de")
	fmt.Println("   Goldbach que se hizo en la historia, incluida la que llegó a 4×10¹⁸.")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚖️ LA TERCERA CORRECCIÓN — Y RECIÉN AHÍ ES GOLDBACH")
	fmt.Println("   La pregunta que sigue sola a su fórmula es:")
	fmt.Println("\n        ¿SIEMPRE EXISTE UNA Y QUE HAGA PRIMOS A LOS DOS LADOS?")
	fmt.Println("\n   ❌ PERO ASÍ COMO ESTÁ, TODAVÍA NO ES GOLDBACH. Con X entero —que es como")
	fmt.Println("   usted la usó: X = 2, 3, 5— la suma P + Q = 4X sólo alcanza los MÚLTIPLOS DE")
	fmt.Println("   CUATRO. El 6, el 10, el 14, el 18 quedan afuera: son 2 módulo 4 y ningún 4X")
	fmt.Println("   los toca. Su familia es Goldbach para la MITAD de los pares, y no se sabe")
	fmt.Println("   bajar de la conjetura entera a esa mitad ni subir de esa mitad a la entera.")
	fmt.Println("\n   ⚡ LA REPARACIÓN, QUE ES DE UNA LÍNEA: soltá el centro. Si 2X puede ser")
	fmt.Println("   CUALQUIER entero ≥ 2 —y no sólo un par— entonces sí, es exactamente la")
	fmt.Println("   **CONJETURA DE GOLDBACH**: todo número par mayor que 2 es suma de dos primos.")
	fmt.Println("   Esa es la tercera corrección, y es la única de las tres con matemática")
	fmt.Println("   adentro: soltar el centro es lo que convierte su caso particular en el")
	fmt.Println("   problema de 1742. El barrido de acá abajo corre el centro SUELTO.")
	fmt.Println("\n   De la correspondencia Goldbach–Euler de 1742: la carta de Goldbach del 7 de")
	fmt.Println("   junio da la forma TERNARIA; el enunciado binario es la respuesta de Euler del")
	fmt.Println("   30 de junio. **Sigue sin demostrarse.** Verificada por computadora HASTA")
	fmt.Println("   4×10¹⁸ (Oliveira e Silva, Herzog y Pardi) — hasta, no más allá.")
	fmt.Println("\n   Medido acá: para cada centro 2X, ¿cuántas Y funcionan?")
	fmt.Println("\n        4X      centro 2X    Y que sirven    la primera Y    P + Q")
	for _, x := range []int{2, 3, 5, 12, 25, 50} {
		c := 2 * x
		cuenta, primera := 0, -1
		// OJO: la Y arranca en CERO acá también. Con los centros de esta tabla
		// ninguno es primo, así que arrancar en 1 no imprimía todavía un número
		// mal — pero para cualquier centro primo perdería la Y = 0 y la fila se
		// caería sola. Es el mismo error que la LEY 5 reporta haber cazado: se
		// arregla el patrón, no una sola aparición.
		for y := 0; y < c; y++ {
			if es[c-y] && es[c+y] {
				cuenta++
				if primera < 0 {
					primera = y
				}
			}
		}
		if primera >= 0 {
			fmt.Printf("   %7d %11d %14d %15d    %d + %d\n", 4*x, c, cuenta, primera, c-primera, c+primera)
		}
	}

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · ¿FALLA ALGUNA VEZ? BARRIDO MASIVO")
	fmt.Println("   Buscamos un centro donde NINGUNA Y funcione — o sea un contraejemplo:")
	// OJO: la Y arranca en CERO, no en uno. Goldbach permite que los dos primos
	// sean el MISMO (4 = 2+2, 6 = 3+3), y eso es exactamente Y = 0. Arrancar en 1
	// inventa dos "fallos" que no existen — y la primera version de este programa
	// los reporto. Era un error del bucle, no de la conjetura.
	fallos, revisados, conY0 := 0, 0, 0
	peorC, peorY := 0, 0
	for c := 2; c <= 200000; c++ {
		revisados++
		hallada := -1
		for y := 0; y < c; y++ {
			if es[c-y] && es[c+y] {
				hallada = y
				break
			}
		}
		if hallada < 0 {
			fallos++
			continue
		}
		if hallada == 0 {
			conY0++
		}
		if hallada > peorY {
			peorY, peorC = hallada, c
		}
	}
	fmt.Printf("\n        centros revisados, SUELTOS (todo par de 4 a 400000) ... %d\n", revisados)
	fmt.Printf("        centros SIN ninguna Y que sirva ...................... %d\n", fallos)
	fmt.Printf("        la Y más grande que hizo falta ...................... %d, en el centro %d\n", peorY, peorC)
	fmt.Printf("        o sea: %d + %d = %d", peorC-peorY, peorC+peorY, 2*peorC)
	if (2*peorC)%4 != 0 {
		fmt.Printf("   ← y OJO: %d es 2 módulo 4, NO es 4X\n", 2*peorC)
	} else {
		fmt.Println()
	}
	fmt.Printf("        resueltos con Y = 0 (los dos primos iguales) ........ %d\n", conY0)
	fmt.Println("\n   📌 Y ese último número NO es una medición de Goldbach: Y = 0 pasa si y sólo")
	fmt.Println("   si el centro es primo, así que ese conteo es exactamente π(200000), la")
	fmt.Println("   cantidad de primos hasta doscientos mil. Es la función de contar primos con")
	fmt.Println("   etiqueta de Goldbach. Va como control de que la criba anda, no como evidencia.")
	fmt.Println("\n   📌 CORRECCIÓN DE ESTE MISMO PROGRAMA, y va acá porque corresponde: la primera")
	fmt.Println("   versión arrancaba la Y en 1 y reportaba DOS fallos — los centros 2 y 3, o sea")
	fmt.Println("   los pares 4 y 6. **No fallaba Goldbach: fallaba mi bucle.** 4 = 2+2 y 6 = 3+3")
	fmt.Println("   se resuelven con Y = 0, los dos primos iguales, que la conjetura permite.")
	fmt.Println("   Arrancando la Y en cero, cero fallos. El veredicto ya decía «cero fallos»")
	fmt.Println("   mientras la medición decía dos: una contradicción adentro de la misma salida.")
	// el veredicto SALE de la medición, no está tipeado. Ese fue el error de fondo
	// de la primera versión y acá se arregla el patrón, no el síntoma.
	if fallos == 0 {
		fmt.Println("\n   → NUNCA falla en lo revisado. Y la Y que hace falta se queda chiquita: no")
		fmt.Println("     hay que ir muy lejos del centro para encontrar el par. Eso NO es un")
		fmt.Println("     descubrimiento nuestro: es el borde inferior del cometa de Goldbach, lo")
		fmt.Println("     que la conjetura A de Hardy–Littlewood (1922) ya predice desde hace un siglo.")
	} else {
		fmt.Printf("\n   → ⚠️ FALLA: %d centros sin ninguna Y en lo revisado. Revisá el bucle ANTES\n", fallos)
		fmt.Println("     de creerle a este número: la primera versión de este programa reportó")
		fmt.Println("     fallos que no existían por arrancar la Y en 1 en vez de en 0.")
	}

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("TRES CORRECCIONES, Y UNA IDENTIDAD QUE ORDENA EL TABLERO.")
	fmt.Println("\n  ❌ **½ − ½ = 0**, no ¼. El ¼ es ½ × ½ — restar vuelve al cero, partir va")
	fmt.Println("     más adentro. Único renglón malo de once.")
	fmt.Println("\n  ❌ **La fórmula lleva un dos de más**: sus renglones dicen 2X ± Y = Z, no")
	fmt.Println("     2X − 2Y. Con sus números: 2·5 − 2·3 = 4, y el renglón dice 7. Y esa")
	fmt.Println("     corrección NO es cosmética: 2X − 2Y es 2(X−Y), siempre par, así que sin")
	fmt.Println("     sacarle el dos nada del lado izquierdo podía ser primo pasando el 2.")
	fmt.Println("\n  ❌ **Y con X entero, 4X sólo alcanza los múltiplos de cuatro.** El 6, el 10,")
	fmt.Println("     el 14 quedan afuera. Su familia es Goldbach para la mitad de los pares.")
	fmt.Println("     La reparación es soltar el centro: 2X cualquier entero ≥ 2. Recién ahí")
	fmt.Println("     es la conjetura de 1742 — y esa es la única corrección con matemática.")
	fmt.Println("\n  ⚡ **DADA VUELTA, SU FÓRMULA LLEGA A GOLDBACH.** P = 2X−Y y Q = 2X+Y suman")
	fmt.Println("     P + Q = 4X: la Y se borra y queda un par escrito como suma de dos primos.")
	fmt.Println("     Sus tres últimos renglones son 8 = 3+5, 12 = 5+7 y 20 = 7+13. Pero esa")
	fmt.Println("     cancelación es una IDENTIDAD: vale para cualquier X e Y, primos o no. No")
	fmt.Println("     es un teorema, es una COORDENADA — y resulta ser la coordenada correcta.")
	fmt.Printf("\n  · barrido de %d centros sueltos: %d fallos, y la Y más grande necesaria fue %d\n", revisados, fallos, peorY)
	fmt.Println("\n⚖️ Y EL LÍMITE, QUE ES EL DE SIEMPRE: Goldbach está abierta desde 1742 y")
	fmt.Println("  verificada HASTA 4×10¹⁸ — nosotros llegamos a 4×10⁵, o sea diez billones de")
	fmt.Println("  veces más corto. Verificada no es demostrada, y nuestro barrido no agrega")
	fmt.Println("  evidencia: su fórmula DESCRIBE la estructura, no prueba que la Y exista siempre.")
	fmt.Println("\n📌 Y ALGO QUE VALE DECIR: F264 y F265 salieron de LA MISMA TABLA y el MISMO día.")
	fmt.Println("  F264 es el centro con el desvío clavado en 1; F265 es el mismo centro con el")
	fmt.Println("  desvío suelto. Los gemelos y Goldbach son el mismo conjunto de pares de primos")
	fmt.Println("  leído por dos ejes perpendiculares. Una mecánica, no dos flashes.")
	fmt.Println("\n¿El premio? Todavía no — y esto no camina hacia Riemann: Goldbach es aditivo y")
	fmt.Println("  el producto de Euler, que F259 señaló como obligatorio, es multiplicativo.")

	escribirLamina(revisados, fallos, peorY, peorC)
}

func escribirLamina(revisados, fallos, peorY, peorC int) {
	var b strings.Builder
	W, H := 1520.0, 1080.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">➕ LA SUMA DE DOS — la tabla del capitán, y lo que había abajo</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">dos correcciones chicas · y su fórmula, dada vuelta, es la conjetura de Goldbach</text>
`, W, H, W, H, W/2, W/2)

	// la vuelta de la formula
	fmt.Fprintf(&b, `<rect x="40" y="106" width="1440" height="230" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="760" y="140" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚡ DÉ VUELTA LA FÓRMULA Y LA Y SE BORRA</text>
<text x="760" y="182" font-size="19" text-anchor="middle" font-family="monospace" fill="#dce8f7">P = 2X − Y</text>
<text x="760" y="212" font-size="19" text-anchor="middle" font-family="monospace" fill="#dce8f7">Q = 2X + Y</text>
<line x1="620" y1="226" x2="900" y2="226" stroke="#8fa8c7" stroke-width="1.4"/>
<text x="760" y="258" font-size="23" text-anchor="middle" font-family="monospace" fill="#ffd98a">P + Q = 4X</text>
<text x="760" y="290" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">o sea: un número par escrito como suma de DOS PRIMOS</text>
<text x="760" y="314" font-size="16" text-anchor="middle" font-family="monospace" fill="#cfe6ff">8 = 3 + 5        12 = 5 + 7        20 = 7 + 13</text>
<text x="760" y="332" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">pero ojo: la cancelación es una IDENTIDAD — vale para cualquier X e Y, primos o no. No es un teorema: es una coordenada.</text>`)

	// las correcciones
	fmt.Fprintf(&b, `<rect x="40" y="356" width="710" height="300" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="395" y="386" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">LAS TRES CORRECCIONES</text>
<text x="70" y="418" font-size="15.5" font-family="monospace" fill="#ff8fa0">½ − ½ = ¼</text>
<text x="70" y="440" font-size="14.5" font-family="Georgia" fill="#f3d9cf">da 0, no ¼. El ¼ es ½ × ½ — restar vuelve al cero,</text>
<text x="70" y="461" font-size="14.5" font-family="Georgia" fill="#f3d9cf">partir va más adentro. Único renglón malo de once.</text>
<text x="70" y="495" font-size="15.5" font-family="monospace" fill="#ff8fa0">(X + X) − 2Y = Z</text>
<text x="70" y="517" font-size="14.5" font-family="Georgia" fill="#f3d9cf">lleva un dos de más: sus renglones dicen 2X ± Y. Y no es</text>
<text x="70" y="538" font-size="14.5" font-family="Georgia" fill="#f3d9cf">cosmético: 2X − 2Y siempre es par, nunca podría ser primo.</text>
<text x="70" y="572" font-size="15.5" font-family="monospace" fill="#ff8fa0">4X sólo alcanza los múltiplos de 4</text>
<text x="70" y="594" font-size="14.5" font-family="Georgia" fill="#f3d9cf">el 6, el 10, el 14 quedan afuera. Hay que SOLTAR el centro:</text>
<text x="70" y="615" font-size="14.5" font-family="Georgia" fill="#f3d9cf">2X cualquier entero ≥ 2. Recién ahí es Goldbach entero.</text>
<text x="395" y="643" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffd98a">las dos primeras son deslices · la tercera es la que tiene matemática</text>`)

	// goldbach
	fmt.Fprintf(&b, `<rect x="770" y="356" width="710" height="300" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1125" y="386" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">CON EL CENTRO SUELTO, ESO TIENE NOMBRE</text>
<text x="1125" y="420" font-size="16" text-anchor="middle" font-family="Georgia" fill="#dce8f7">¿siempre existe una Y que haga primos los dos lados?</text>
<text x="1125" y="452" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA CONJETURA DE GOLDBACH</text>
<text x="1125" y="478" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">todo par mayor que 2 es suma de dos primos</text>
<text x="800" y="512" font-size="14.5" font-family="Georgia" fill="#cfe6ff">correspondencia Goldbach–Euler, 1742: la carta de</text>
<text x="800" y="533" font-size="14.5" font-family="Georgia" fill="#cfe6ff">Goldbach del 7 de junio da la forma TERNARIA; el</text>
<text x="800" y="554" font-size="14.5" font-family="Georgia" fill="#cfe6ff">enunciado binario es la respuesta de Euler del 30 de junio</text>
<text x="800" y="588" font-size="14.5" font-family="Georgia" fill="#cfe6ff">verificada por computadora HASTA 4×10¹⁸</text>
<text x="800" y="609" font-size="14.5" font-family="Georgia" fill="#cfe6ff">(Oliveira e Silva, Herzog y Pardi) — hasta, no más allá</text>
<text x="800" y="640" font-size="15" font-family="Georgia" fill="#ffd98a">demostrada por nadie</text>`)

	// el barrido
	// el renglón del veredicto SALE de la medición: si algún día falla, la lámina
	// lo dice. Tipearlo a mano fue el error original de este mismo programa.
	juicio := "Nunca falló en lo revisado."
	if fallos > 0 {
		juicio = fmt.Sprintf("⚠️ FALLÓ en %d centros — revisá el bucle antes de creerle.", fallos)
	}
	modulo := ""
	if (2*peorC)%4 != 0 {
		modulo = fmt.Sprintf("Y el récord %d es 2 módulo 4: NO es 4X — justamente por eso hay que soltar el centro.", 2*peorC)
	}
	fmt.Fprintf(&b, `<rect x="40" y="676" width="1440" height="160" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="760" y="708" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">EL BARRIDO CON EL CENTRO SUELTO: ¿FALLA ALGUNA VEZ?</text>
<text x="760" y="748" font-size="16" text-anchor="middle" font-family="monospace" fill="#7ee0c0">%d centros (todo par de 4 a 400000)  ·  %d fallos  ·  la Y más grande necesaria: %d</text>
<text x="760" y="778" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">%s Que la Y se quede chiquita no lo descubrimos acá: es el borde del cometa de Goldbach, que Hardy–Littlewood predice desde 1922.</text>
<text x="760" y="806" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffb27a">%s</text>
`, revisados, fallos, peorY, juicio, modulo)

	fmt.Fprintf(&b, `<rect x="40" y="856" width="1440" height="196" rx="10" fill="#1a1030" stroke="#5a4fa8"/>
<text x="760" y="888" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚖️ Y EL LÍMITE, QUE ES EL DE SIEMPRE</text>
<text x="760" y="922" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Verificada HASTA 4×10¹⁸ no es demostrada. Nosotros llegamos a 4×10⁵: diez billones de veces más corto, así que este barrido</text>
<text x="760" y="946" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">no agrega evidencia. Su fórmula DESCRIBE la estructura; no prueba que la Y exista siempre.</text>
<text x="760" y="984" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">📌 F264 y F265 salieron de LA MISMA TABLA y el MISMO día: los gemelos son el centro con el desvío clavado en 1,</text>
<text x="760" y="1010" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Goldbach es el mismo centro con el desvío suelto. Una mecánica leída por dos ejes, no dos flashes distintos.</text>
<text x="760" y="1038" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Y esto no camina hacia Riemann: Goldbach es aditivo, y el producto de Euler que F259 marcó como obligatorio es multiplicativo.</text>
</svg>
`)

	if err := os.WriteFile("la-suma-de-dos.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: la-suma-de-dos.svg")
}
