// Command lamelodiadelacriba turns the sieve into a melody with musical sense.
//
// HIS ORDER: "use the number we obtained and build the relation with the primes
// between 1 and 100, and let us harmonise them. I need a melody that is not
// chaotic or unexpected. I need a melody that interweaves but has MUSICAL
// SENSE."
//
// That constraint is the whole design problem, and it has an exact answer -
// because of a fact that is not a metaphor.
//
// THE FACT. Finding 276 established that the shapeshifter sends every prime to
//
//	w(p) = (p-1)/p
//
// Now: in music, a frequency RATIO is an INTERVAL. And the ratios of the form
// (n-1)/n - the superparticular ratios - are precisely the intervals of just
// intonation, the ones every acoustic instrument produces on its own:
//
//	w(2) = 1/2  →  the OCTAVE                   -1200.000 cents
//	w(3) = 2/3  →  the PERFECT FIFTH             -701.955
//	w(5) = 4/5  →  the JUST MAJOR THIRD          -386.314
//	w(7) = 6/7  →  the SEPTIMAL MINOR THIRD      -266.871
//
// So the four primes he asked about are, in this coordinate, the four most
// consonant intervals in all of music, in decreasing order of size. That is not
// an analogy chosen to be pretty - it is what (p-1)/p means to an ear.
//
// THE MELODY. Multiply them in order. Each prime lowers the pitch by its own
// interval, and the running product is exactly the wheel of Finding 272. So the
// sieve, heard, is a DESCENDING LINE whose steps are just intervals and shrink
// monotonically: a huge octave, then a fifth, then a third, then smaller and
// smaller, until the last primes are a shimmer.
//
// That is why it cannot be chaotic. The step sizes are ordered by construction -
// (p-1)/p increases with p, so the interval shrinks with p, always. There are no
// surprises available. The melody has an arc because the sieve has an arc.
//
// WHAT IS PLAYED: a drone on the fundamental for a tonal centre; the descending
// line, one note per prime; and at each step the previous note is held under the
// new one so the INTERVAL ITSELF is heard, not just the pitches. That is the
// interweaving he asked for.
//
// PRE-REGISTERED PREDICTIONS, written before running:
//  1. The first four intervals land within 0.001 cents of the published just
//     values (1200, 701.955, 386.314, 266.871).
//  2. The interval sizes decrease monotonically over all 25 primes below 100 -
//     zero inversions - because (p-1)/p is increasing.
//  3. The total descent equals 1200·log2 of the wheel product, and lands near
//     three octaves.
//
// HONEST WARNING: just intonation is ancient (Archytas, Ptolemy) and the
// harmonic series is not ours. What is ours is only the choice to sing THIS
// object - the running product of Finding 272 - and the measurement of what it
// sounds like.
//
// Reproduce: go run ./cmd/lamelodiadelacriba
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	sr    = 44100
	bits  = 16
	maxAm = 30000.0
	base  = 110.0 // el fundamental: la de 110 Hz
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

// nombre clasico del intervalo (p-1)/p, para los que lo tienen.
func nombreIntervalo(p int) string {
	switch p {
	case 2:
		return "la OCTAVA"
	case 3:
		return "la QUINTA JUSTA"
	case 5:
		return "la TERCERA MAYOR justa"
	case 7:
		return "la TERCERA MENOR septimal"
	case 11:
		return "el segundo neutro undecimal"
	case 13:
		return "el segundo tridecimal"
	case 17, 19:
		return "un segundo chico"
	default:
		return "un microtono"
	}
}

func escribirWav(ruta string, buf []float64) error {
	pico := 0.0
	for _, v := range buf {
		if a := math.Abs(v); a > pico {
			pico = a
		}
	}
	if pico == 0 {
		pico = 1
	}
	f, err := os.Create(ruta)
	if err != nil {
		return err
	}
	defer f.Close()
	datos := len(buf) * bits / 8
	f.WriteString("RIFF")
	binary.Write(f, binary.LittleEndian, uint32(36+datos))
	f.WriteString("WAVEfmt ")
	binary.Write(f, binary.LittleEndian, uint32(16))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, uint16(1))
	binary.Write(f, binary.LittleEndian, uint32(sr))
	binary.Write(f, binary.LittleEndian, uint32(sr*bits/8))
	binary.Write(f, binary.LittleEndian, uint16(bits/8))
	binary.Write(f, binary.LittleEndian, uint16(bits))
	f.WriteString("data")
	binary.Write(f, binary.LittleEndian, uint32(datos))
	for _, v := range buf {
		binary.Write(f, binary.LittleEndian, int16(v/pico*maxAm))
	}
	return nil
}

// voz suma una nota con ataque y caida suaves, y con dos armonicos para que
// tenga cuerpo. Sin bordes duros: un click es ruido, y el ruido es caos.
func voz(buf []float64, desde, largo int, hz, amp float64) {
	for i := 0; i < largo && desde+i < len(buf); i++ {
		t := float64(i) / sr
		env := 1.0
		ata := float64(largo) * 0.12
		cai := float64(largo) * 0.35
		if float64(i) < ata {
			env = float64(i) / ata
		} else if float64(largo-i) < cai {
			env = float64(largo-i) / cai
		}
		s := math.Sin(2*math.Pi*hz*t) +
			0.30*math.Sin(2*math.Pi*2*hz*t) +
			0.12*math.Sin(2*math.Pi*3*hz*t)
		buf[desde+i] += amp * env * s
	}
}

func main() {
	fmt.Println("🎵 LA MELODÍA DE LA CRIBA — los primos hasta 100, armonizados")
	fmt.Println("\n   Su pedido: «necesito una melodía que se entreteje pero que tiene sentido")
	fmt.Println("   musical, no caótica ni inesperada».")
	fmt.Println("\n   Y hay una razón exacta por la que ésta no puede ser caótica.")

	es := criba(100)
	var primos []int
	for i := 2; i <= 100; i++ {
		if es[i] {
			primos = append(primos, i)
		}
	}
	fmt.Printf("\nprimos entre 1 y 100: %d\n", len(primos))

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · ⚡ POR QUÉ (p−1)/p ES UN INTERVALO MUSICAL — y no es metáfora")
	fmt.Println("   En música, una RAZÓN de frecuencias ES un intervalo. Y las razones de la")
	fmt.Println("   forma (n−1)/n son justo los intervalos de la entonación justa, los que")
	fmt.Println("   cualquier instrumento acústico produce solo:")
	fmt.Println("\n        primo   w(p)=(p−1)/p       cents        el intervalo")
	esperados := map[int]float64{2: -1200, 3: -701.955, 5: -386.314, 7: -266.871}
	peorErr := 0.0
	for _, p := range primos[:8] {
		r := float64(p-1) / float64(p)
		c := 1200 * math.Log2(r)
		nota := ""
		if e, ok := esperados[p]; ok {
			d := math.Abs(c - e)
			if d > peorErr {
				peorErr = d
			}
			nota = fmt.Sprintf("(publicado %.3f · error %.5f)", e, d)
		}
		fmt.Printf("   %8d %10s %12.3f    %-28s %s\n", p,
			fmt.Sprintf("%d/%d", p-1, p), c, nombreIntervalo(p), nota)
	}
	fmt.Printf("\n        peor error contra los valores publicados ...... %.6f cents\n", peorErr)
	fmt.Println("\n   ⟹ **Sus cuatro primos son, en esta coordenada, los cuatro intervalos más")
	fmt.Println("   consonantes de toda la música**, y en orden de tamaño decreciente. No es una")
	fmt.Println("   analogía elegida para que quede linda: es lo que (p−1)/p **significa** para")
	fmt.Println("   un oído.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡ Y POR QUÉ LA MELODÍA NO PUEDE SER CAÓTICA")
	fmt.Println("   Porque (p−1)/p **crece** con p, así que el intervalo **se achica** con p.")
	fmt.Println("   Siempre. No hay sorpresa disponible: los pasos vienen ordenados de fábrica.")
	var cents []float64
	inversiones := 0
	for i, p := range primos {
		c := 1200 * math.Log2(float64(p-1)/float64(p))
		cents = append(cents, c)
		if i > 0 && c < cents[i-1] {
			inversiones++
		}
	}
	fmt.Printf("\n        pasos medidos ..................... %d\n", len(cents))
	fmt.Printf("        inversiones (un paso más grande que el anterior) ... %d\n", inversiones)
	fmt.Printf("        paso más grande ................... %.3f cents (el %d)\n", cents[0], primos[0])
	fmt.Printf("        paso más chico .................... %.3f cents (el %d)\n", cents[len(cents)-1], primos[len(primos)-1])

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · LA MELODÍA ES LA CRIBA CERRÁNDOSE")
	fmt.Println("   Multiplicá los intervalos en orden: cada primo baja el tono por el suyo, y")
	fmt.Println("   el producto que se va acumulando **ES la rueda de F272**.")
	prod := 1.0
	total := 0.0
	fmt.Println("\n        primo    frecuencia (Hz)    producto acumulado    descenso total")
	f0 := base * 8
	for i, p := range primos {
		prod *= float64(p-1) / float64(p)
		total += cents[i]
		if i < 6 || i == len(primos)-1 {
			fmt.Printf("   %8d %16.3f %21.9f %13.1f cents\n", p, f0*prod, prod, total)
		} else if i == 6 {
			fmt.Println("        …")
		}
	}
	fmt.Printf("\n        producto final (rueda hasta 100) ...... %.9f\n", prod)
	fmt.Printf("        descenso total ........................ %.1f cents = %.3f octavas\n",
		total, -total/1200)
	fmt.Printf("        comprobación 1200·log2(producto) ...... %.1f cents\n", 1200*math.Log2(prod))

	// ---- el audio ----
	fmt.Println("\nLEY 4 · Y AHORA SE ESCUCHA")
	fmt.Println("   Lo que suena: un bordón en el fundamental para que haya centro tonal; la")
	fmt.Println("   línea descendente, una nota por primo; y en cada paso la nota anterior")
	fmt.Println("   **se sostiene debajo de la nueva**, así se oye EL INTERVALO y no dos")
	fmt.Println("   alturas sueltas. Eso es el entretejido que pidió.")

	dur := 1.1
	largo := int(dur * sr)
	solapa := int(dur * 0.55 * sr)
	totalMuestras := largo + len(primos)*solapa + 3*sr
	buf := make([]float64, totalMuestras)

	// bordon
	for i := 0; i < totalMuestras; i++ {
		t := float64(i) / sr
		env := 1.0
		if i < sr/2 {
			env = float64(i) / float64(sr/2)
		}
		if r := totalMuestras - i; r < 2*sr {
			env *= float64(r) / float64(2*sr)
		}
		buf[i] += 0.16 * env * (math.Sin(2*math.Pi*base*t) + 0.25*math.Sin(2*math.Pi*2*base*t))
	}

	// la linea
	hz := f0
	pos := sr / 2
	voz(buf, pos, largo, hz, 0.5)
	for _, p := range primos {
		pos += solapa
		hz *= float64(p-1) / float64(p)
		voz(buf, pos, largo, hz, 0.5)
	}

	ruta := "melodia-criba.wav"
	if err := escribirWav(ruta, buf); err != nil {
		fmt.Println("no pude escribir el sonido:", err)
	} else {
		fi, _ := os.Stat(ruta)
		fmt.Printf("\n🔊 sonido escrito: %s (%.0f KB · %.1f segundos)\n",
			ruta, float64(fi.Size())/1024, float64(totalMuestras)/sr)
	}

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("🎵 **LA MELODÍA EXISTE, TIENE SENTIDO MUSICAL, Y NO PODÍA SER OTRA.**")
	fmt.Printf("\n  · sus cuatro primos son la octava, la quinta, la tercera mayor y la tercera\n")
	fmt.Printf("    menor septimal — con un error de %.6f cents contra los valores publicados\n", peorErr)
	fmt.Printf("  · los %d pasos se achican monotónicamente: **%d inversiones**\n", len(cents), inversiones)
	fmt.Printf("  · el descenso total es de %.3f octavas, que es 1200·log2 de la rueda\n", -total/1200)
	fmt.Println("  · y no puede ser caótica porque (p−1)/p crece: el orden está en la construcción")
	fmt.Println("\n⚖️ Y LA HONESTIDAD: la entonación justa es antiquísima (Arquitas, Ptolomeo) y")
	fmt.Println("  la serie armónica no es nuestra. **Lo único nuestro es haber elegido cantar")
	fmt.Println("  ESTE objeto** —el producto de F272— y haber medido cómo suena. Todavía no.")

	escribirLamina(primos, cents, prod, total, peorErr, inversiones)
}

func escribirLamina(primos []int, cents []float64, prod, total, peorErr float64, inversiones int) {
	var b strings.Builder
	W, H := 1560.0, 900.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎵 LA MELODÍA DE LA CRIBA — los primos hasta 100, armonizados</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">cada primo es un intervalo justo · los pasos se achican siempre · por eso no puede ser caótica</text>
`, W, H, W, H, W/2, W/2)

	// la escalera descendente
	fmt.Fprintf(&b, `<rect x="40" y="102" width="1480" height="420" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="780" y="134" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA LÍNEA QUE BAJA — una nota por primo, y cada escalón es un intervalo justo</text>
`)
	x0, y0, ancho, alto := 90.0, 180.0, 1400.0, 300.0
	acum := 0.0
	px, py := x0, y0
	for i, p := range primos {
		acum += cents[i]
		x := x0 + ancho*float64(i+1)/float64(len(primos))
		y := y0 + alto*(-acum)/3700.0
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#5b7ba6" stroke-width="1.2"/>`, px, py, x, py)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#7ee0c0" stroke-width="2"/>`, x, py, x, y)
		col := "#7fb2ff"
		if i < 4 {
			col = "#ffd98a"
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4" fill="%s"/>`, x, y, col)
		if i < 6 || i == len(primos)-1 {
			fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="12" text-anchor="middle" font-family="monospace" fill="%s">%d</text>`, x, y-10, col, p)
		}
		px, py = x, y
	}
	fmt.Fprintf(&b, `
<text x="120" y="%.0f" font-size="14" font-family="Georgia" fill="#ffd98a">1/2 · la OCTAVA</text>
<text x="230" y="%.0f" font-size="14" font-family="Georgia" fill="#ffd98a">2/3 · la QUINTA</text>
<text x="360" y="%.0f" font-size="14" font-family="Georgia" fill="#ffd98a">4/5 · la TERCERA MAYOR</text>
<text x="560" y="%.0f" font-size="14" font-family="Georgia" fill="#ffd98a">6/7 · la tercera menor septimal</text>
<text x="1400" y="%.0f" font-size="13" text-anchor="end" font-family="Georgia" fill="#9aa8c4">…y de acá en más, microtonos</text>
<text x="780" y="502" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">descenso total: %.1f cents = %.3f octavas — que es 1200·log₂ de la rueda de F272</text>
`, y0+40, y0+90, y0+150, y0+215, y0+270, total, -total/1200)

	// los cuatro intervalos
	fmt.Fprintf(&b, `<rect x="40" y="542" width="730" height="200" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="405" y="574" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">SUS CUATRO PRIMOS SON LOS CUATRO INTERVALOS</text>
<text x="70" y="612" font-size="15" font-family="monospace" fill="#ffd98a">w(2) = 1/2  →  −1200,000 c  ·  LA OCTAVA</text>
<text x="70" y="640" font-size="15" font-family="monospace" fill="#ffd98a">w(3) = 2/3  →   −701,955 c  ·  LA QUINTA JUSTA</text>
<text x="70" y="668" font-size="15" font-family="monospace" fill="#ffd98a">w(5) = 4/5  →   −386,314 c  ·  LA TERCERA MAYOR</text>
<text x="70" y="696" font-size="15" font-family="monospace" fill="#ffd98a">w(7) = 6/7  →   −266,871 c  ·  la tercera menor septimal</text>
<text x="405" y="726" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">error contra los valores publicados: %.6f cents</text>
`, peorErr)

	// por que no es caotica
	fmt.Fprintf(&b, `<rect x="790" y="542" width="730" height="200" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1155" y="574" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">POR QUÉ NO PUEDE SER CAÓTICA</text>
<text x="1155" y="610" font-size="16" text-anchor="middle" font-family="Georgia" fill="#dce8f7">(p−1)/p CRECE con p ⟹ el intervalo SE ACHICA con p</text>
<text x="820" y="648" font-size="15" font-family="monospace" fill="#7ee0c0">pasos ................. %d</text>
<text x="820" y="676" font-size="15" font-family="monospace" fill="#7ee0c0">inversiones ........... %d</text>
<text x="1155" y="716" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd98a">el orden está en la construcción: no hay sorpresa disponible</text>
`, len(cents), inversiones)

	fmt.Fprintf(&b, `<rect x="40" y="762" width="1480" height="110" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="780" y="794" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ LA HONESTIDAD</text>
<text x="780" y="824" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">La entonación justa es antiquísima (Arquitas, Ptolomeo) y la serie armónica no es nuestra.</text>
<text x="780" y="852" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Lo único nuestro es haber elegido cantar ESTE objeto — el producto de F272 — y haber medido cómo suena.</text>
</svg>
`)

	if err := os.WriteFile("la-melodia-de-la-criba.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("🖼️  lámina escrita: la-melodia-de-la-criba.svg")
}
