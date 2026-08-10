// Command elrelieve draws what the captain's ups and downs actually form.
//
// HIS ORDER: "graph the jumps, the positive ones and the negative ones, and
// build me what they form."
//
// The signed sequence he means is the one Finding 267 identified as what his
// hand-written "plus two, minus two" terms really were:
//
//	g_i = p_{i+1} - p_i        the gaps
//	d_i = g_i - g_{i+1}        the SIGNED second difference - his ups and downs
//
// FOUR PANELS, AND TWO OF THEM CARRY A WARNING ON THE FACE.
//
//  1. THE RELIEF - the signed d_i drawn as a skyline above and below zero.
//     Real: it is the actual sequence, nothing is imposed.
//  2. THE GILBREATH TRIANGLE - take absolute differences over and over. Row 1
//     is the gaps, row 2 is |d_i|, and so on down. Every row begins with 1.
//     That is Gilbreath's conjecture (1958) and it is OPEN. Caveat that must
//     travel with it: the property is not special to the primes - Proth stated
//     it in 1878 and it holds for a broad class of sequences - so the triangle
//     is beautiful and unproved, but it is not evidence about primes alone.
//  3. THE WALK OF THE SIGNS - cumulative sum of sign(d_i). This one has real
//     content, because the ORDER of the gaps decides it.
//     ⚠️ NOT drawn: the cumulative sum of the d_i themselves, because it
//     telescopes to g_1 - g_{n+1}. That walk cannot go anywhere the gaps do not
//     already send it, and drawing it would be the seventh 0.0e+00 trap.
//  4. THE DISTRIBUTION - how often each value appears, and the shuffle control.
//     ⚠️ Most of the evenness is just "gaps are even past 2", not a fact about
//     order. The control says how much survives a shuffle.
//
// Reproduce: go run ./cmd/elrelieve
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"
	"strings"
)

const N = 2000000 // primos hasta acá para la estadística

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

// baraja mezcla una copia con un generador determinista (sin Math.random ni time):
// congruencial lineal con semilla fija, para que el control sea reproducible.
func baraja(src []int, semilla uint64) []int {
	c := append([]int(nil), src...)
	s := semilla
	for i := len(c) - 1; i > 0; i-- {
		s = s*6364136223846793005 + 1442695040888963407
		j := int((s >> 33) % uint64(i+1))
		c[i], c[j] = c[j], c[i]
	}
	return c
}

func segundas(g []int) []int {
	d := make([]int, len(g)-1)
	for i := range d {
		d[i] = g[i] - g[i+1]
	}
	return d
}

func main() {
	fmt.Println("⛰️  EL RELIEVE — los saltos positivos y negativos del capitán, y lo que forman")
	fmt.Printf("\ncribando hasta %d…\n", N)
	es := criba(N)
	var primos []int
	for i := 2; i <= N; i++ {
		if es[i] {
			primos = append(primos, i)
		}
	}
	gaps := make([]int, len(primos)-1)
	for i := range gaps {
		gaps[i] = primos[i+1] - primos[i]
	}
	d := segundas(gaps)
	fmt.Printf("primos: %d · saltos: %d · subidas y bajadas: %d\n", len(primos), len(gaps), len(d))

	// ---- LEY 1: el relieve ----
	fmt.Println("\nLEY 1 · EL RELIEVE: CUÁNTOS SUBEN, CUÁNTOS BAJAN, CUÁNTOS SE QUEDAN")
	pos, neg, cero := 0, 0, 0
	for _, v := range d {
		switch {
		case v > 0:
			pos++
		case v < 0:
			neg++
		default:
			cero++
		}
	}
	fmt.Printf("\n        suben  (d > 0) ..... %8d  (%.3f%%)\n", pos, 100*float64(pos)/float64(len(d)))
	fmt.Printf("        bajan  (d < 0) ..... %8d  (%.3f%%)\n", neg, 100*float64(neg)/float64(len(d)))
	fmt.Printf("        planos (d = 0) ..... %8d  (%.3f%%)\n", cero, 100*float64(cero)/float64(len(d)))
	fmt.Printf("        desbalance ......... %8d  (%.4f%% del total)\n", pos-neg, 100*float64(pos-neg)/float64(len(d)))
	fmt.Println("\n   📌 Casi empatados, y eso NO es casualidad ni es hallazgo: si subieran mucho")
	fmt.Println("   más que las bajadas, los saltos crecerían sin freno, y se sabe que crecen")
	fmt.Println("   apenas como ln(p). El empate está forzado por el tamaño, no por el orden.")

	// ---- LEY 2: la trampa que NO vamos a dibujar ----
	fmt.Println("\nLEY 2 · ⚠️ LA CAMINATA QUE NO SE DIBUJA, Y POR QUÉ")
	fmt.Println("   Lo primero que uno quiere hacer con una tira de más y menos es sumarla.")
	fmt.Println("   Pero esta suma se TELESCOPIA:")
	fmt.Println("\n        d₁ + d₂ + … + dₙ = (g₁−g₂) + (g₂−g₃) + … = g₁ − gₙ₊₁")
	suma := 0
	for i := 0; i < 50000 && i < len(d); i++ {
		suma += d[i]
	}
	n := 50000
	if n > len(d) {
		n = len(d)
	}
	fmt.Printf("\n        sumado a mano hasta n=%d ....... %d\n", n, suma)
	fmt.Printf("        g₁ − g₍ₙ₊₁₎ ..................... %d\n", gaps[0]-gaps[n])
	fmt.Println("\n   ⟹ **IDÉNTICOS, y tienen que serlo.** Esa caminata no puede ir a ningún lado")
	fmt.Println("   que los saltos no la manden. Dibujarla sería la SÉPTIMA trampa del 0.0e+00")
	fmt.Println("   de este laboratorio. Va nombrada y NO va graficada.")

	// ---- LEY 3: la caminata de los signos, que sí tiene contenido ----
	fmt.Println("\nLEY 3 · LA CAMINATA DE LOS SIGNOS — ésta sí depende del ORDEN")
	paseo := make([]int, len(d)+1)
	maxExc, minExc := 0, 0
	for i, v := range d {
		s := 0
		if v > 0 {
			s = 1
		} else if v < 0 {
			s = -1
		}
		paseo[i+1] = paseo[i] + s
		if paseo[i+1] > maxExc {
			maxExc = paseo[i+1]
		}
		if paseo[i+1] < minExc {
			minExc = paseo[i+1]
		}
	}
	fmt.Printf("\n        posición final ..................... %d\n", paseo[len(paseo)-1])
	fmt.Printf("        excursión máxima hacia arriba ...... %d\n", maxExc)
	fmt.Printf("        excursión máxima hacia abajo ....... %d\n", minExc)
	fmt.Printf("        √n de referencia (azar puro) ....... %.0f\n", math.Sqrt(float64(len(d))))

	// control barajado
	gb := baraja(gaps, 0x5EED1234)
	db := segundas(gb)
	pb, nb, cb := 0, 0, 0
	for _, v := range db {
		switch {
		case v > 0:
			pb++
		case v < 0:
			nb++
		default:
			cb++
		}
	}
	paseoB := 0
	maxB, minB := 0, 0
	for _, v := range db {
		s := 0
		if v > 0 {
			s = 1
		} else if v < 0 {
			s = -1
		}
		paseoB += s
		if paseoB > maxB {
			maxB = paseoB
		}
		if paseoB < minB {
			minB = paseoB
		}
	}
	fmt.Println("\n   EL CONTROL — los MISMOS saltos, barajados. Lo que sobrevive al barajado es")
	fmt.Println("   de la distribución; lo que se muere era del ORDEN:")
	fmt.Printf("\n        %-34s %12s %12s\n", "", "primos", "barajado")
	fmt.Printf("        %-34s %11.3f%% %11.3f%%\n", "suben", 100*float64(pos)/float64(len(d)), 100*float64(pb)/float64(len(db)))
	fmt.Printf("        %-34s %11.3f%% %11.3f%%\n", "bajan", 100*float64(neg)/float64(len(d)), 100*float64(nb)/float64(len(db)))
	fmt.Printf("        %-34s %11.3f%% %11.3f%%\n", "planos (d = 0)", 100*float64(cero)/float64(len(d)), 100*float64(cb)/float64(len(db)))
	fmt.Printf("        %-34s %12d %12d\n", "excursión máxima arriba", maxExc, maxB)
	fmt.Printf("        %-34s %12d %12d\n", "excursión máxima abajo", minExc, minB)

	// ---- LEY 4: el triángulo ----
	fmt.Println("\nLEY 4 · ⚡ Y LO QUE FORMAN DE VERDAD: EL TRIÁNGULO DE GILBREATH")
	fmt.Println("   Tomá los saltos. Restá vecino con vecino en valor absoluto. Te queda otra")
	fmt.Println("   fila. Repetí. Y otra vez. Y otra vez. Eso arma un triángulo — y la fila 2")
	fmt.Println("   de ese triángulo es exactamente el valor absoluto de SUS más-dos-menos-dos.")
	filas := 220
	anchoT := 420
	tri := construirTriangulo(gaps, filas, anchoT)
	rotas := 0
	for r := 1; r < len(tri); r++ {
		if len(tri[r]) > 0 && tri[r][0] != 1 {
			rotas++
		}
	}
	fmt.Printf("\n        filas construidas .................. %d\n", len(tri)-1)
	fmt.Printf("        filas que NO empiezan en 1 ......... %d\n", rotas)
	fmt.Println("\n   ⟹ **TODAS EMPIEZAN CON UN 1.** Eso es la conjetura de GILBREATH (1958),")
	fmt.Println("   y sigue **sin demostrarse** — este laboratorio ya la tiene como Hallazgo 6.")
	fmt.Println("\n   ⚖️ Y EL LÍMITE QUE VIAJA CON ELLA, que hay que decirlo o no vale: **la")
	fmt.Println("   propiedad NO es exclusiva de los primos.** Proth la enunció en 1878 y vale")
	fmt.Println("   para una familia ancha de sucesiones. El triángulo es hermoso y está")
	fmt.Println("   abierto, pero NO es evidencia sobre los primos en particular.")
	fmt.Println("\n        fila   valores 0    valores 2   otros    largo")
	for _, r := range []int{1, 2, 5, 10, 50, 100, 200} {
		if r >= len(tri) {
			continue
		}
		c0, c2, otros := 0, 0, 0
		for _, v := range tri[r] {
			switch v {
			case 0:
				c0++
			case 2:
				c2++
			default:
				otros++
			}
		}
		fmt.Printf("   %8d %11d %12d %7d %8d\n", r, c0, c2, otros, len(tri[r]))
	}

	// ---- LEY 5: la distribución ----
	fmt.Println("\nLEY 5 · LA DISTRIBUCIÓN, Y CUÁNTO DE ELLA ES SOLO «LOS SALTOS SON PARES»")
	hist := map[int]int{}
	histB := map[int]int{}
	for _, v := range d {
		hist[v]++
	}
	for _, v := range db {
		histB[v]++
	}
	impares := 0
	for _, v := range d {
		if v%2 != 0 {
			impares++
		}
	}
	fmt.Printf("\n        subidas/bajadas IMPARES ............ %d de %d (%.5f%%)\n",
		impares, len(d), 100*float64(impares)/float64(len(d)))
	fmt.Println("        (sólo pueden aparecer alrededor del 2, el único primo par)")
	fmt.Println("\n          d      en los primos        %      barajado        %     ¿simétrico?")
	var claves []int
	for k := range hist {
		if k >= -12 && k <= 12 {
			claves = append(claves, k)
		}
	}
	sort.Ints(claves)
	for _, k := range claves {
		if hist[k] == 0 {
			continue
		}
		sim := "—"
		if k > 0 && hist[-k] > 0 {
			sim = fmt.Sprintf("%+.2f%%", 100*float64(hist[k]-hist[-k])/float64(hist[k]+hist[-k]))
		}
		fmt.Printf("   %8d %14d %9.3f%% %13d %9.3f%%   %s\n",
			k, hist[k], 100*float64(hist[k])/float64(len(d)),
			histB[k], 100*float64(histB[k])/float64(len(db)), sim)
	}

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("⚡ LO QUE FORMAN: **EL TRIÁNGULO DE GILBREATH.** Todas sus filas empiezan con")
	fmt.Println("  un 1, nadie sabe demostrar por qué, y hace 68 años que está así.")
	fmt.Printf("\n  · %d filas construidas · %d que no empiezan en 1\n", len(tri)-1, rotas)
	fmt.Printf("  · suben %.3f%% · bajan %.3f%% · planos %.3f%%\n",
		100*float64(pos)/float64(len(d)), 100*float64(neg)/float64(len(d)), 100*float64(cero)/float64(len(d)))
	fmt.Println("\n⚠️ Y DOS COSAS QUE NO SE DIBUJAN COMO HALLAZGO, POR HONESTIDAD:")
	fmt.Println("  · la suma corrida de los d se telescopia en g₁ − gₙ₊₁: no puede sorprender.")
	fmt.Println("  · el empate entre subidas y bajadas está forzado porque los saltos crecen")
	fmt.Println("    como ln(p): un desbalance sostenido sería imposible.")
	fmt.Println("\n⚖️ Y EL LÍMITE MÁS IMPORTANTE: el triángulo NO es propiedad de los primos.")
	fmt.Println("  Proth 1878 — vale para una familia ancha de sucesiones. Es hermoso, está")
	fmt.Println("  abierto, y no dice nada que sea SOLO de los primos. Todavía no.")

	escribirLamina(gaps, d, tri, hist, histB, pos, neg, cero, maxExc, minExc, rotas, len(d), impares)
}

func construirTriangulo(seq []int, filas, ancho int) [][]int {
	if len(seq) > ancho {
		seq = seq[:ancho]
	}
	tri := [][]int{append([]int(nil), seq...)}
	cur := seq
	for r := 1; r <= filas && len(cur) > 1; r++ {
		nx := make([]int, len(cur)-1)
		for i := range nx {
			v := cur[i+1] - cur[i]
			if v < 0 {
				v = -v
			}
			nx[i] = v
		}
		tri = append(tri, nx)
		cur = nx
	}
	return tri
}

// trianguloPNG dibuja el triángulo como PNG y lo devuelve como data URI.
func trianguloPNG(tri [][]int, ancho, alto int) string {
	img := image.NewRGBA(image.Rect(0, 0, ancho, alto))
	for y := 0; y < alto; y++ {
		fila := []int(nil)
		if y+1 < len(tri) {
			fila = tri[y+1]
		}
		for x := 0; x < ancho; x++ {
			v := 0
			if x < len(fila) {
				v = fila[x]
			}
			img.Set(x, y, rgbTri(v))
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return ""
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func rgbTri(v int) color.RGBA {
	switch {
	case v == 0:
		return color.RGBA{0x13, 0x20, 0x38, 0xff}
	case v == 1:
		return color.RGBA{0xff, 0xd9, 0x8a, 0xff}
	case v == 2:
		return color.RGBA{0x2f, 0x7f, 0x63, 0xff}
	case v == 4:
		return color.RGBA{0x7e, 0xe0, 0xc0, 0xff}
	case v == 6:
		return color.RGBA{0x7f, 0xb2, 0xff, 0xff}
	case v <= 12:
		return color.RGBA{0xc9, 0xb6, 0xff, 0xff}
	default:
		return color.RGBA{0xff, 0x8f, 0xa0, 0xff}
	}
}

func escribirLamina(gaps, d []int, tri [][]int, hist, histB map[int]int,
	pos, neg, cero, maxExc, minExc, rotas, total, impares int) {

	var b strings.Builder
	W, H := 1640.0, 1780.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="48" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⛰️ EL RELIEVE — las subidas y bajadas del capitán, y lo que forman</text>
<text x="%.0f" y="78" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">d = salto anterior − salto siguiente · sus «más dos, menos dos» de F267 · %d subidas y bajadas medidas</text>
`, W, H, W, H, W/2, W/2, total)

	// ---------- PANEL 1: el relieve ----------
	base := 300.0
	fmt.Fprintf(&b, `<rect x="40" y="102" width="1560" height="330" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="820" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">1 · EL RELIEVE — los primeros 480 pasos, tal cual son</text>
<line x1="70" y1="%.0f" x2="1570" y2="%.0f" stroke="#5b7ba6" stroke-width="1"/>
<text x="60" y="%.0f" font-size="12" text-anchor="end" font-family="monospace" fill="#8fa8c7">0</text>
`, base, base, base+4)
	esc := 7.0
	nB := 480
	if nB > len(d) {
		nB = len(d)
	}
	anchoB := 1500.0 / float64(nB)
	for i := 0; i < nB; i++ {
		v := float64(d[i]) * esc
		if v > 110 {
			v = 110
		}
		if v < -110 {
			v = -110
		}
		x := 70 + float64(i)*anchoB
		col := "#7ee0c0"
		if d[i] < 0 {
			col = "#ff8fa0"
		} else if d[i] == 0 {
			col = "#5b7ba6"
		}
		y, h := base-v, v
		if v < 0 {
			y, h = base, -v
		}
		if h < 1 {
			h = 1
			y = base - 0.5
		}
		fmt.Fprintf(&b, `<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s"/>`, x, y, anchoB*0.82, h, col)
	}
	fmt.Fprintf(&b, `
<text x="70" y="176" font-size="14" font-family="Georgia" fill="#7ee0c0">verde: el salto se achicó (subida del relieve)</text>
<text x="70" y="196" font-size="14" font-family="Georgia" fill="#ff8fa0">rojo: el salto se agrandó (bajada)</text>
<text x="1570" y="176" font-size="14" text-anchor="end" font-family="monospace" fill="#cfe6ff">suben %d (%.2f%%) · bajan %d (%.2f%%) · planos %d (%.2f%%)</text>
<text x="1570" y="196" font-size="13" text-anchor="end" font-family="Georgia" fill="#9aa8c4">el empate está FORZADO: si no, los saltos crecerían sin freno, y crecen como ln(p)</text>
<text x="820" y="418" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">sólo %d de %d son impares (%.5f%%): sólo pueden aparecer cerca del 2, el único primo par</text>
`, pos, 100*float64(pos)/float64(total), neg, 100*float64(neg)/float64(total), cero, 100*float64(cero)/float64(total),
		impares, total, 100*float64(impares)/float64(total))

	// ---------- PANEL 2: el triangulo ----------
	fmt.Fprintf(&b, `<rect x="40" y="452" width="1560" height="880" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="820" y="484" font-size="19" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">2 · ⚡ LO QUE FORMAN: EL TRIÁNGULO DE GILBREATH</text>
<text x="820" y="510" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">fila 1 = los saltos · fila 2 = el valor absoluto de SUS más-dos-menos-dos · y después restá vecinos, otra vez, y otra vez</text>
`)
	cx, cy, celda := 96.0, 528.0, 3.6
	// se dibuja por TRAMOS de color igual, no celda por celda: la lámina pasa de
	// 86.000 rectángulos a unos pocos miles y de 4,8 MB a un archivo manejable.
	// El color de fondo del panel es el mismo que el del 0, así que el 0 no se dibuja.
	// El triángulo va como PNG embebido, un píxel por celda. Dibujarlo con un
	// rectángulo por celda daba 86.000 nodos y 4,8 MB; en tramos daba 18.000 y
	// 1,3 MB. Como imagen queda en unos pocos KB y se ve igual de nítido con
	// image-rendering:pixelated.
	fmt.Fprintf(&b, `<image x="%.2f" y="%.2f" width="%.2f" height="%.2f" image-rendering="pixelated" href="%s"/>`,
		cx, cy, 400*celda, 215*celda, trianguloPNG(tri, 400, 215))
	fmt.Fprintf(&b, `
<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ffd98a" stroke-width="2"/>
<text x="%.1f" y="%.1f" font-size="14" text-anchor="end" font-family="Georgia" fill="#ffd98a">esta columna</text>
<text x="%.1f" y="%.1f" font-size="14" text-anchor="end" font-family="Georgia" fill="#ffd98a">es TODA de unos</text>
<text x="%.1f" y="%.1f" font-size="13" text-anchor="end" font-family="Georgia" fill="#9aa8c4">y nadie sabe</text>
<text x="%.1f" y="%.1f" font-size="13" text-anchor="end" font-family="Georgia" fill="#9aa8c4">demostrar por qué</text>
`, cx-4, cy, cx-4, cy+215*celda, cx-14, cy+40, cx-14, cy+60, cx-14, cy+86, cx-14, cy+104)

	lx := cx + 400*celda + 40
	fmt.Fprintf(&b, `
<text x="%.0f" y="%.0f" font-size="15" font-family="Georgia" fill="#dce8f7">LOS COLORES</text>
<rect x="%.0f" y="%.0f" width="14" height="14" fill="#ffd98a"/><text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#cfe6ff">1</text>
<rect x="%.0f" y="%.0f" width="14" height="14" fill="#132038" stroke="#26456e"/><text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#cfe6ff">0</text>
<rect x="%.0f" y="%.0f" width="14" height="14" fill="#2f7f63"/><text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#cfe6ff">2</text>
<rect x="%.0f" y="%.0f" width="14" height="14" fill="#7ee0c0"/><text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#cfe6ff">4</text>
<rect x="%.0f" y="%.0f" width="14" height="14" fill="#7fb2ff"/><text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#cfe6ff">6</text>
<rect x="%.0f" y="%.0f" width="14" height="14" fill="#c9b6ff"/><text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#cfe6ff">8 a 12</text>
<rect x="%.0f" y="%.0f" width="14" height="14" fill="#ff8fa0"/><text x="%.0f" y="%.0f" font-size="13" font-family="monospace" fill="#cfe6ff">más de 12</text>
<text x="%.0f" y="%.0f" font-size="15" font-family="Georgia" fill="#ffd98a">%d filas</text>
<text x="%.0f" y="%.0f" font-size="15" font-family="Georgia" fill="#7ee0c0">%d no empiezan en 1</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#cfe6ff">Conjetura de GILBREATH,</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#cfe6ff">1958. Verificada hasta</text>
<text x="%.0f" y="%.0f" font-size="14" font-family="Georgia" fill="#cfe6ff">profundidades enormes.</text>
<text x="%.0f" y="%.0f" font-size="15" font-family="Georgia" fill="#ffd98a">ABIERTA.</text>
<text x="%.0f" y="%.0f" font-size="13.5" font-family="Georgia" fill="#9aa8c4">Ya estaba en su propio</text>
<text x="%.0f" y="%.0f" font-size="13.5" font-family="Georgia" fill="#9aa8c4">registro: Hallazgo 6.</text>
`,
		lx, cy+14,
		lx, cy+30, lx+22, cy+42, lx, cy+54, lx+22, cy+66, lx, cy+78, lx+22, cy+90,
		lx, cy+102, lx+22, cy+114, lx, cy+126, lx+22, cy+138, lx, cy+150, lx+22, cy+162,
		lx, cy+174, lx+22, cy+186,
		lx, cy+230, len(tri)-1,
		lx, cy+256, rotas,
		lx, cy+300, lx, cy+320, lx, cy+340,
		lx, cy+366,
		lx, cy+398, lx, cy+416)

	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="360" height="150" rx="8" fill="#33221c" stroke="#c0392b"/>
<text x="%.0f" y="%.0f" font-size="14.5" font-family="Georgia" fill="#ffb27a">⚖️ EL LÍMITE QUE VIAJA CON ESTO</text>
<text x="%.0f" y="%.0f" font-size="13.5" font-family="Georgia" fill="#f3d9cf">La propiedad NO es exclusiva de los</text>
<text x="%.0f" y="%.0f" font-size="13.5" font-family="Georgia" fill="#f3d9cf">primos: Proth la enunció en 1878 y</text>
<text x="%.0f" y="%.0f" font-size="13.5" font-family="Georgia" fill="#f3d9cf">vale para una familia ancha de</text>
<text x="%.0f" y="%.0f" font-size="13.5" font-family="Georgia" fill="#f3d9cf">sucesiones. Es hermoso y está abierto,</text>
<text x="%.0f" y="%.0f" font-size="13.5" font-family="Georgia" fill="#f3d9cf">pero no dice nada que sea SÓLO</text>
<text x="%.0f" y="%.0f" font-size="13.5" font-family="Georgia" fill="#f3d9cf">de los primos.</text>
`, lx, cy+450, lx+12, cy+476, lx+12, cy+500, lx+12, cy+520, lx+12, cy+540, lx+12, cy+560, lx+12, cy+580, lx+12, cy+600)

	// ---------- PANEL 3: la caminata de los signos ----------
	fmt.Fprintf(&b, `<rect x="40" y="1352" width="770" height="250" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="425" y="1384" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">3 · LA CAMINATA DE LOS SIGNOS</text>
<text x="425" y="1408" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">sumás +1 por cada subida y −1 por cada bajada. Ésta SÍ depende del orden.</text>
`)
	// mini-paseo dibujado sobre los primeros 6000 pasos
	px, py, pw, ph := 70.0, 1500.0, 710.0, 70.0
	pasos := 6000
	if pasos > len(d) {
		pasos = len(d)
	}
	acum, mx, mn := 0, 0, 0
	serie := make([]int, pasos+1)
	for i := 0; i < pasos; i++ {
		if d[i] > 0 {
			acum++
		} else if d[i] < 0 {
			acum--
		}
		serie[i+1] = acum
		if acum > mx {
			mx = acum
		}
		if acum < mn {
			mn = acum
		}
	}
	rango := float64(mx - mn)
	if rango < 1 {
		rango = 1
	}
	var pl strings.Builder
	for i := 0; i <= pasos; i += 4 {
		x := px + pw*float64(i)/float64(pasos)
		y := py - ph*(float64(serie[i]-mn)/rango-0.5)*2
		fmt.Fprintf(&pl, "%.1f,%.1f ", x, y)
	}
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#5b7ba6" stroke-width="1"/>
<polyline points="%s" fill="none" stroke="#c9b6ff" stroke-width="1.6"/>
<text x="70" y="1440" font-size="13.5" font-family="monospace" fill="#cfe6ff">excursión máxima arriba: %d · abajo: %d</text>
<text x="70" y="1462" font-size="13" font-family="Georgia" fill="#9aa8c4">primeros %d pasos dibujados · sobre %d medidos</text>
<text x="425" y="1586" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">⚠️ y la OTRA caminata, la de los valores, NO se dibuja: se telescopia en g₁ − gₙ₊₁ y no puede sorprender</text>
`, px, py, px+pw, py, pl.String(), maxExc, minExc, pasos, total)

	// ---------- PANEL 4: la distribucion y el control ----------
	fmt.Fprintf(&b, `<rect x="830" y="1352" width="770" height="250" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1215" y="1384" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">4 · LA DISTRIBUCIÓN, CONTRA SU CONTROL BARAJADO</text>
<text x="1215" y="1406" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">barras llenas: los primos · contorno: los MISMOS saltos, barajados</text>
`)
	hx, hy, hw := 880.0, 1560.0, 660.0
	maxH := 0
	for k := -12; k <= 12; k += 2 {
		if hist[k] > maxH {
			maxH = hist[k]
		}
		if histB[k] > maxH {
			maxH = histB[k]
		}
	}
	if maxH == 0 {
		maxH = 1
	}
	bw := hw / 13.0
	for j, k := 0, -12; k <= 12; k, j = k+2, j+1 {
		h := 120.0 * float64(hist[k]) / float64(maxH)
		hb := 120.0 * float64(histB[k]) / float64(maxH)
		x := hx + float64(j)*bw
		col := "#7ee0c0"
		if k < 0 {
			col = "#ff8fa0"
		} else if k == 0 {
			col = "#5b7ba6"
		}
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`, x, hy-h, bw*0.72, h, col)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="none" stroke="#ffd98a" stroke-width="1.2"/>`, x, hy-hb, bw*0.72, hb)
		fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="11" text-anchor="middle" font-family="monospace" fill="#8fa8c7">%d</text>`, x+bw*0.36, hy+16, k)
	}
	fmt.Fprintf(&b, `
<text x="1215" y="1594" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">⚠️ lo que sobrevive al barajado es de la DISTRIBUCIÓN, no del orden: casi toda esta figura lo es</text>
`)

	// ---------- cierre ----------
	fmt.Fprintf(&b, `<rect x="40" y="1622" width="1560" height="130" rx="10" fill="#1a1030" stroke="#5a4fa8"/>
<text x="820" y="1654" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚖️ LO QUE FORMAN, DICHO SIN ADORNOS</text>
<text x="820" y="1686" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Sus subidas y bajadas son la segunda fila de un triángulo que baja para siempre, y cuya primera columna es toda de unos.</text>
<text x="820" y="1712" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Eso se llama conjetura de Gilbreath, es de 1958, sigue abierta — y no es exclusiva de los primos (Proth, 1878).</text>
<text x="820" y="1740" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd98a">El capitán la escuchó de oído antes de saber que existía. Eso no lo hace cualquiera. Pero todavía no.</text>
</svg>
`)

	if err := os.WriteFile("el-relieve.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-relieve.svg")
}
