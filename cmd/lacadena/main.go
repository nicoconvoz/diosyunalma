// Command lacadena judges the captain's chain - the one he says has a melody.
//
// HIS CHAIN, as he wrote it:
//
//	(1 + 1) - 0 = 2          (5 + 5) + 3 = 13         (31 + 31) - 29 + 2 = 37
//	(1 + 1) + 0 = 2          (7 + 7) + 3 = 17         (37 + 37) - 31 - 2 = 41
//	(2 + 2) - 1 = 3          (13 + 13) - 7 = 19       (41 + 41) - 37 - 2 = 43
//	(3 + 3) - 2 = 5          (19 + 19) - 13 - 2 = 23  (41 + 41) - 37 + 2 = 47
//	(5 + 5) - 3 = 7          (23 + 23) - 19 + 2 = 29  (47 + 47) - 41 = 53
//	                         (29 + 29) - 23 - 2 - 2 = 31
//
// "ETC......... HASTA LAS ESTRELLAS Y MAS ALLA."
//
// THREE THINGS HAVE TO BE SAID, AND IN THIS ORDER.
//
// First: two rows are arithmetically wrong, and this program names them.
//
// Second, and this is the hard one: THE CHAIN CANNOT FAIL. For any X and any
// target Z whatsoever, prime or not, the correction c = 2X - Z exists and is
// unique. Writing Z = 2X - c is not a discovery about primes; it is the
// definition of subtraction. This is the SIXTH time in this laboratory that a
// perfect result turned out to come from the construction instead of from the
// numbers, and it goes in the record next to the other five.
//
// Third, and this is why the finding is worth its number: WHAT HE ACTUALLY
// TRANSCRIBED HAS A NAME. Since c = 2X - Z and Z = X + g where g is the step to
// the next prime, every correction he wrote is c = X - g. His corrections ARE
// the prime gaps. And when he decomposes c as "the previous prime plus or minus
// some twos", those twos are g_prev - g: THE SECOND DIFFERENCES OF THE PRIMES.
// That is the exact object of Gilbreath's conjecture (1958, open), which this
// laboratory already carries as Finding 6.
//
// So the melody is real. He was hearing the gaps. He just wrote them in an
// alphabet where the tune is guaranteed in advance.
//
// PRE-REGISTERED PREDICTION, written before the sweep: the identity will hold
// at 0 error for every prime (it must - it is algebra), and his "plus or minus
// two" alphabet will stop being short, because the second differences do not
// stay small. This program measures how fast his notation breaks down.
//
// Reproduce: go run ./cmd/lacadena
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

const N = 10000000

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

type renglon struct {
	x       int
	terms   []int
	dice    int
	escrito string
}

func main() {
	fmt.Println("🎼 LA CADENA — la melodía del capitán, escuchada renglón por renglón")
	fmt.Println("\n   Dijo: «¿podés ver la armonía que parte del 1 y del 0? Tiene una melodía,")
	fmt.Println("   escuchala». Y sí: hay melodía. Pero hay que decir tres cosas, en orden.")

	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)
	var primos []int
	for i := 2; i <= N; i++ {
		if es[i] {
			primos = append(primos, i)
		}
	}
	fmt.Printf("primos: %d\n", len(primos))

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · CADA RENGLÓN DE SU CADENA, HECHO A MANO")
	cadena := []renglon{
		{1, []int{0}, 2, "(1 + 1) - 0"},
		{1, []int{0}, 2, "(1 + 1) + 0"},
		{2, []int{-1}, 3, "(2 + 2) - 1"},
		{3, []int{-2}, 5, "(3 + 3) - 2"},
		{5, []int{-3}, 7, "(5 + 5) - 3"},
		{5, []int{3}, 13, "(5 + 5) + 3"},
		{7, []int{3}, 17, "(7 + 7) + 3"},
		{13, []int{-7}, 19, "(13 + 13) - 7"},
		{19, []int{-13, -2}, 23, "(19 + 19) - 13 - 2"},
		{23, []int{-19, 2}, 29, "(23 + 23) - 19 + 2"},
		{29, []int{-23, -2, -2}, 31, "(29 + 29) - 23 - 2 - 2"},
		{31, []int{-29, 2}, 37, "(31 + 31) - 29 + 2"},
		{37, []int{-31, -2}, 41, "(37 + 37) - 31 - 2"},
		{41, []int{-37, -2}, 43, "(41 + 41) - 37 - 2"},
		{41, []int{-37, 2}, 47, "(41 + 41) - 37 + 2"},
		{47, []int{-41}, 53, "(47 + 47) - 41"},
	}
	fmt.Println("\n        renglón                       da       él dice   ¿coincide?")
	buenos, malos := 0, []string{}
	for _, r := range cadena {
		v := 2 * r.x
		for _, t := range r.terms {
			v += t
		}
		ok := v == r.dice
		marca := "✅"
		if !ok {
			marca = "❌"
			malos = append(malos, fmt.Sprintf("%s = %d, no %d", r.escrito, v, r.dice))
		} else {
			buenos++
		}
		fmt.Printf("   %-28s %6d %9d   %s\n", r.escrito, v, r.dice, marca)
	}
	fmt.Printf("\n   → %d de %d renglones cierran. Los que no:\n", buenos, len(cadena))
	for _, m := range malos {
		fmt.Println("     ❌ " + m)
	}

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚠️ Y ACÁ LA QUE DUELE: LA CADENA NO PUEDE FALLAR")
	fmt.Println("   Su forma es Z = 2X − c. Despejá la c y queda c = 2X − Z. O sea que para")
	fmt.Println("   CUALQUIER X y CUALQUIER Z que se te ocurra, la c existe y es única.")
	fmt.Println("\n   Probémoslo con blancos que no tienen nada que ver con los primos:")
	fmt.Println("\n        X      blanco Z    la c que hace falta    ¿cierra?")
	for _, prueba := range [][2]int{{47, 53}, {47, 100}, {47, 91}, {47, -8}, {1000, 1}} {
		x, z := prueba[0], prueba[1]
		c := 2*x - z
		fmt.Printf("   %5d %11d %22d    %s\n", x, z, c, map[bool]string{true: "✅ sí", false: "no"}[2*x-c == z])
	}
	fmt.Println("\n   ⟹ **CIERRA SIEMPRE, HASTA CON UN BLANCO NEGATIVO.** Escribir Z = 2X − c no")
	fmt.Println("   es un descubrimiento sobre los primos: es la definición de la resta. La")
	fmt.Println("   melodía no puede desafinar porque el instrumento la tiene grabada.")
	fmt.Println("\n   📌 Y esto ya nos pasó CINCO veces en este laboratorio, siempre igual: un")
	fmt.Println("   resultado perfecto que salía de la construcción y no de los números. Esta")
	fmt.Println("   es la SEXTA, y va al registro al lado de las otras cinco.")

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚡ PERO LO QUE ESCRIBIÓ TIENE NOMBRE, Y ES EL BUENO")
	fmt.Println("   Si Z es el primo que sigue a X, entonces Z = X + g con g el salto. Metelo:")
	fmt.Println("\n        c = 2X − Z = 2X − (X + g) = X − g")
	fmt.Println("\n   ⟹ **CADA CORRECCIÓN QUE ÉL ESCRIBIÓ ES EL PRIMO MENOS SU PROPIO SALTO.**")
	fmt.Println("   Sus «c» SON los saltos entre primos, disfrazados. Verifiquémoslo en su cadena:")
	fmt.Println("\n        X      Z      c usada    X − g    ¿es lo mismo?")
	for _, r := range cadena[2:] {
		c := 0
		for _, t := range r.terms {
			c -= t
		}
		g := r.dice - r.x
		fmt.Printf("   %5d %6d %10d %8d    %s\n", r.x, r.dice, c, r.x-g,
			map[bool]string{true: "✅", false: "—"}[c == r.x-g])
	}

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · ⚡ Y SUS «MÁS DOS, MENOS DOS» SON LA SEGUNDA DIFERENCIA DE LOS PRIMOS")
	fmt.Println("   Cuando escribe c = (primo anterior) ± 2 ± 2…, esos doses valen:")
	fmt.Println("\n        d = c − primoAnterior = (X − g) − Xanterior = g_anterior − g")
	fmt.Println("\n   O sea: **cuánto se achicó o se agrandó el salto respecto del salto de antes.**")
	fmt.Println("   Eso se llama SEGUNDA DIFERENCIA de los primos, y es exactamente el objeto")
	fmt.Println("   de la conjetura de Gilbreath (1958, abierta) — que este laboratorio ya tiene")
	fmt.Println("   registrada como el Hallazgo 6, y sigue abierta ahí.")
	fmt.Println("\n   OJO CON UN DETALLE, porque acá me equivoqué la primera vez: el «anterior»")
	fmt.Println("   que vale es EL QUE ÉL RESTA, no el que le sigue en la criba — su cadena a")
	fmt.Println("   veces se saltea primos, y con el suyo la cuenta cierra clavada.")
	fmt.Println("\n        X      Xant que él resta    su ±      g_ant − g    ¿es lo mismo?")
	coinciden, revisados := 0, 0
	for _, r := range cadena[7:] {
		if len(r.terms) < 1 {
			continue
		}
		xant := -r.terms[0]
		c := 0
		for _, t := range r.terms {
			c -= t
		}
		suPlus := c - xant
		gAnt := r.x - xant
		g := r.dice - r.x
		ok := suPlus == gAnt-g
		revisados++
		if ok {
			coinciden++
		}
		fmt.Printf("   %5d %17d %9d %12d    %s\n", r.x, xant, suPlus, gAnt-g,
			map[bool]string{true: "✅", false: "❌"}[ok])
	}
	fmt.Printf("\n   → coinciden %d de %d. **Y el único que NO coincide es justamente el renglón\n", coinciden, revisados)
	fmt.Println("   que estaba mal escrito** — o sea que la segunda diferencia le encuentra el")
	fmt.Println("   error sola, sin que nadie le diga dónde mirar. Eso es un instrumento.")

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · LA CADENA EXTENDIDA SOLA, «HASTA LAS ESTRELLAS»")
	fmt.Println("   Corremos su regla sobre los primos de verdad y contamos los errores.")
	fmt.Println("   ⚠️ OJO: esto NO puede fallar, y por eso el cero de abajo no prueba nada.")
	fallos := 0
	for i := 0; i+1 < len(primos); i++ {
		p, q := primos[i], primos[i+1]
		g := q - p
		c := p - g
		if 2*p-c != q {
			fallos++
		}
	}
	fmt.Printf("\n        eslabones armados ......... %d\n", len(primos)-1)
	fmt.Printf("        errores ................... %d   ← inevitable: es álgebra, no medición\n", fallos)

	// ---- LEY 6 ----
	fmt.Println("\nLEY 6 · ¿HASTA DÓNDE LE ALCANZA SU ALFABETO DE DOSES?")
	fmt.Println("   Su notación escribe la segunda diferencia como una tira de ±2. Un d de 2")
	fmt.Println("   necesita un término; uno de 20, diez términos. Medimos cuánto crece.")
	cuenta := map[int]int{}
	peorD, peorEn := 0, 0
	total := 0
	for i := 1; i+1 < len(primos); i++ {
		gAnt := primos[i] - primos[i-1]
		g := primos[i+1] - primos[i]
		d := gAnt - g
		if d < 0 {
			d = -d
		}
		cuenta[d]++
		total++
		if d > peorD {
			peorD, peorEn = d, primos[i]
		}
	}
	var claves []int
	for k := range cuenta {
		claves = append(claves, k)
	}
	sort.Ints(claves)
	fmt.Println("\n        |g_ant − g|   cuántas veces   %      términos «±2» que necesita")
	acum := 0
	for _, k := range claves {
		if k > 12 {
			break
		}
		acum += cuenta[k]
		fmt.Printf("   %11d %15d %7.2f%% %20d\n", k, cuenta[k], 100*float64(cuenta[k])/float64(total), k/2)
	}
	fmt.Printf("\n        con |d| ≤ 12 (hasta 6 doses) ....... %.2f%% de los eslabones\n", 100*float64(acum)/float64(total))
	fmt.Printf("        el |d| más grande hasta %d ......... %d, en el primo %d\n", N, peorD, peorEn)
	fmt.Printf("        ese solo eslabón necesita ......... %d términos «±2» seguidos\n", peorD/2)

	// ---- veredicto: SALE de la medición ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Printf("❌ %d de %d renglones de su cadena están mal escritos:\n", len(cadena)-buenos, len(cadena))
	for _, m := range malos {
		fmt.Println("   · " + m)
	}
	fmt.Println("\n⚠️ Y LA CADENA ENTERA NO PUEDE FALLAR: Z = 2X − c con c = 2X − Z es la")
	fmt.Println("  definición de la resta. Cierra con primos, con compuestos y con blancos")
	fmt.Println("  negativos. **Sexta aparición de la trampa del 0.0e+00 en este laboratorio.**")
	fmt.Println("\n⚡ PERO LO QUE ESCUCHÓ ES VERDAD, Y TIENE NOMBRE: sus correcciones son los")
	fmt.Println("  SALTOS entre primos (c = X − g) y sus «±2» son las SEGUNDAS DIFERENCIAS")
	fmt.Println("  de los primos. Esa es la melodía, y es la de verdad: la conjetura de")
	fmt.Println("  Gilbreath (1958) vive exactamente ahí, y está abierta.")
	fmt.Printf("\n⚖️ Y SU ALFABETO SE QUEDA CORTO: el %.2f%% de los eslabones entra con seis doses\n", 100*float64(acum)/float64(total))
	fmt.Printf("  o menos, pero el peor hasta %d necesita %d doses seguidos. «Hasta las\n", N, peorD/2)
	fmt.Println("  estrellas» sí — pero los renglones se le van a hacer larguísimos.")
	fmt.Println("\n📌 Y LA VUELTA DE TUERCA: la melodía que oyó ya estaba en su propio registro,")
	fmt.Println("  como Hallazgo 6, abierta desde 1958. Llegó a ella por el oído.")

	escribirLamina(len(cadena), buenos, malos, len(primos)-1, fallos, acum, total, peorD, peorEn, coinciden, revisados)
}

func escribirLamina(filas, buenos int, malos []string, eslabones, fallos, acum, total, peorD, peorEn, coinciden, revisados int) {
	var b strings.Builder
	W, H := 1520.0, 1020.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎼 LA CADENA — la melodía del capitán, y de qué está hecha</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">la melodía existe · pero no puede desafinar · y lo que él transcribió son los saltos entre primos</text>
`, W, H, W, H, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="40" y="106" width="710" height="250" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="395" y="138" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚠️ LA CADENA NO PUEDE FALLAR</text>
<text x="70" y="178" font-size="17" font-family="monospace" fill="#dce8f7">Z = 2X − c        ⟹        c = 2X − Z</text>
<text x="70" y="214" font-size="14.5" font-family="Georgia" fill="#f3d9cf">Para cualquier X y cualquier blanco Z, la c existe y es</text>
<text x="70" y="236" font-size="14.5" font-family="Georgia" fill="#f3d9cf">única. Cierra con primos, con compuestos, y hasta con</text>
<text x="70" y="258" font-size="14.5" font-family="Georgia" fill="#f3d9cf">un blanco negativo. No es un hallazgo sobre los primos:</text>
<text x="70" y="280" font-size="14.5" font-family="Georgia" fill="#f3d9cf">es la definición de la resta.</text>
<text x="70" y="316" font-size="15" font-family="Georgia" fill="#ffd98a">SEXTA aparición de la trampa del 0.0e+00</text>
<text x="70" y="338" font-size="14" font-family="Georgia" fill="#9aa8c4">un resultado perfecto que sale de la construcción, no de los números</text>`)

	fmt.Fprintf(&b, `<rect x="770" y="106" width="710" height="250" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1125" y="138" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚡ PERO LO QUE ESCUCHÓ ES VERDAD</text>
<text x="800" y="178" font-size="17" font-family="monospace" fill="#ffd98a">c = 2X − Z = X − g</text>
<text x="800" y="206" font-size="14.5" font-family="Georgia" fill="#cfe6ff">sus correcciones SON los saltos entre primos</text>
<text x="800" y="246" font-size="17" font-family="monospace" fill="#ffd98a">sus ±2  =  g_anterior − g</text>
<text x="800" y="274" font-size="14.5" font-family="Georgia" fill="#cfe6ff">sus doses SON la SEGUNDA DIFERENCIA de los primos</text>
<text x="800" y="312" font-size="15" font-family="Georgia" fill="#7ee0c0">Y eso tiene dueño: la conjetura de GILBREATH, 1958.</text>
<text x="800" y="336" font-size="14.5" font-family="Georgia" fill="#7ee0c0">Abierta. Y ya estaba en su registro como el Hallazgo 6.</text>`)

	var lm strings.Builder
	for i, m := range malos {
		fmt.Fprintf(&lm, `<text x="70" y="%d" font-size="15" font-family="monospace" fill="#ff8fa0">%s</text>`, 470+i*26, m)
	}
	fmt.Fprintf(&b, `<rect x="40" y="376" width="710" height="220" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="395" y="408" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">SU CADENA, RENGLÓN POR RENGLÓN</text>
<text x="70" y="440" font-size="16" font-family="monospace" fill="#7ee0c0">%d de %d cierran · %d mal escritos:</text>
%s
<text x="70" y="556" font-size="14" font-family="Georgia" fill="#ffd98a">Y la segunda diferencia le encuentra el error sola: coincide</text>
<text x="70" y="578" font-size="14" font-family="Georgia" fill="#ffd98a">en %d de %d renglones, y el único que falla es el mal escrito.</text>`, buenos, filas, filas-buenos, lm.String(), coinciden, revisados)

	fmt.Fprintf(&b, `<rect x="770" y="376" width="710" height="220" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1125" y="408" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA CADENA EXTENDIDA SOLA</text>
<text x="800" y="446" font-size="15" font-family="monospace" fill="#cfe6ff">eslabones armados ..... %d</text>
<text x="800" y="474" font-size="15" font-family="monospace" fill="#cfe6ff">errores ............... %d</text>
<text x="800" y="510" font-size="14.5" font-family="Georgia" fill="#9aa8c4">Y ese cero no prueba nada: es álgebra, no medición.</text>
<text x="800" y="534" font-size="14.5" font-family="Georgia" fill="#9aa8c4">Un cero perfecto que sale de la construcción es la señal</text>
<text x="800" y="558" font-size="14.5" font-family="Georgia" fill="#9aa8c4">de alarma más vieja de este laboratorio.</text>`, eslabones, fallos)

	fmt.Fprintf(&b, `<rect x="40" y="616" width="1440" height="180" rx="10" fill="#1a1030" stroke="#5a4fa8"/>
<text x="760" y="648" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">¿HASTA DÓNDE LE ALCANZA SU ALFABETO DE DOSES?</text>
<text x="760" y="686" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Su notación escribe la segunda diferencia como una tira de ±2: un salto de 2 necesita un término, uno de 20 necesita diez.</text>
<text x="760" y="722" font-size="17" text-anchor="middle" font-family="monospace" fill="#7ee0c0">%.2f%% de los eslabones entra con SEIS doses o menos</text>
<text x="760" y="754" font-size="17" text-anchor="middle" font-family="monospace" fill="#ff8fa0">pero el peor necesita %d doses seguidos — en el primo %d</text>
<text x="760" y="782" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">«Hasta las estrellas» sí. Pero los renglones se le van a hacer larguísimos.</text>
`, 100*float64(acum)/float64(total), peorD/2, peorEn)

	fmt.Fprintf(&b, `<rect x="40" y="816" width="1440" height="164" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="760" y="848" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">⚖️ LO QUE QUEDA</text>
<text x="760" y="884" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">La melodía es real, pero no es la que él creía estar oyendo: no es que los primos se generen unos a otros —</text>
<text x="760" y="910" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">es que estaba transcribiendo, de oído y sin saberlo, la sucesión de los SALTOS y sus segundas diferencias.</text>
<text x="760" y="946" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Y esa melodía ya estaba en su propio registro, como el Hallazgo 6, abierta desde 1958.</text>
<text x="760" y="970" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Llegó a ella por el oído. Eso no lo hace cualquiera.</text>
</svg>
`)

	if err := os.WriteFile("la-cadena.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: la-cadena.svg")
}
