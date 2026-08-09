// Command arbolmitades measures the captain's bisection flash: halve 7, halve
// it again, halve it forever - you get a deep and unreachable sea. And the
// midpoint relation is inherited: 6 - 1/2 - 7 and 7 - 1/2 - 8 give 6 - 1/2 - 8,
// so +inf - 1/2 - -inf, and in dimension 0: 0 - 1/2 - 1.
//
// The flash has a sharp consequence nobody had written down here yet. Bisect
// the segment between the two stakes forever and you generate the dyadic
// tree, which is dense: every possible place a pearl could sit is reached in
// the limit. Now flip the world with the shapeshifter, beta -> 1 - beta:
//
//	EVERY node of the infinite tree is paired with another - EXCEPT ONE.
//	The very first cut, 1/2, is its own pair.
//
// So the captain's infinite sea has exactly one self-paired point, and the
// hypothesis says the book puts every single pearl on it. Measured here:
// the pairing at each depth, and the PRICE (the AM-GM cost of F229) at every
// node of the tree - zero at the root and strictly positive everywhere else.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
	"strings"
)

func w(rho complex128) complex128 { return 1 - 1/rho }

// precio is the AM-GM cost of sitting at beta instead of the midpoint:
// (N + S)/2 - 1, where N and S are the north and south sizes.
func precio(beta, gamma float64) float64 {
	N := cmplx.Abs(w(complex(beta, gamma)))
	S := cmplx.Abs(w(complex(1-beta, gamma)))
	return (N+S)/2 - 1
}

func main() {
	fmt.Println("🌳 EL ÁRBOL DE LAS MITADES — partir y partir: el mar profundo e inalcanzable")

	// ---- LAW 1: halve 7 forever ----
	fmt.Println("\nLEY 1 · PARTIR EL 7 PARA SIEMPRE — se acerca al fondo y nunca lo toca")
	fmt.Println("   corte        lo que queda de 7             ¿llegó al 0?")
	v := 7.0
	for _, k := range []int{1, 2, 3, 10, 50, 200, 1000} {
		v = 7.0 / math.Pow(2, float64(k))
		llego := "no"
		if v == 0 {
			llego = "SÍ (se acabó la máquina, no el número)"
		}
		fmt.Printf("   %4d         %.6e                  %s\n", k, v, llego)
	}
	fmt.Println("   → cada corte lo acerca a la mitad de lo que faltaba y JAMÁS lo entrega:")
	fmt.Println("     eso es el mar profundo e inalcanzable del capitán — y es la dimensión 0")

	// ---- LAW 2: the midpoint relation is inherited ----
	fmt.Println("\nLEY 2 · LA RELACIÓN SE HEREDA — 6 —½— 7 —½— 8  entonces  6 —½— 8")
	fmt.Println("   a      b      c       ¿b es el medio de a y c?     el medio de los extremos")
	for _, tr := range [][3]float64{{6, 7, 8}, {0, 0.5, 1}, {-100, 0, 100}, {2, 5, 8}} {
		a, b, c := tr[0], tr[1], tr[2]
		medio := (a + c) / 2
		ok := "sí"
		if math.Abs(medio-b) > 1e-12 {
			ok = "no"
		}
		fmt.Printf("  %7.2f %6.2f %7.2f        %s                        %.4f\n", a, b, c, ok, medio)
	}
	fmt.Println("   → y llevado al límite: +∞ —½— −∞ se encuentran en la dimensión 0,")
	fmt.Println("     y en el mundo del libro esa cadena es exactamente  0 —½— 1")

	// ---- LAW 3: the infinite tree has exactly ONE self-paired node ----
	fmt.Println("\nLEY 3 · EL ÁRBOL INFINITO TIENE UN SOLO NODO QUE ES SU PROPIA PAREJA")
	fmt.Println("   se parte [0,1] una y otra vez y se da vuelta el mundo con β → 1−β:")
	fmt.Println("   profundidad   nodos del árbol   parejas   nodos que son su propia pareja")
	type nivel struct {
		d, nodos, parejas, solos int
	}
	var niveles []nivel
	for d := 1; d <= 12; d++ {
		den := 1 << d
		vistos := map[float64]bool{}
		for k := 1; k < den; k++ {
			vistos[float64(k)/float64(den)] = true
		}
		nodos := len(vistos)
		solos := 0
		for b := range vistos {
			if math.Abs(b-(1-b)) < 1e-15 {
				solos++
			}
		}
		niveles = append(niveles, nivel{d, nodos, (nodos - solos) / 2, solos})
		if d <= 4 || d == 8 || d == 12 {
			fmt.Printf("      %2d          %8d       %7d              %d\n", d, nodos, (nodos-solos)/2, solos)
		}
	}
	fmt.Println("   → a cualquier profundidad, por más nodos que haya, SIEMPRE hay exactamente UNO")
	fmt.Println("     que se mira al espejo y se ve a sí mismo: el 1/2. Todos los demás vienen de a dos")

	// ---- LAW 4: the price on every node of the tree ----
	fmt.Println("\nLEY 4 · EL PRECIO EN CADA NODO DEL ÁRBOL — cero en la raíz, positivo en todo lo demás")
	fmt.Println("   (el precio es el de F229: (N+S)/2 − 1, lo que cuesta no estar en el medio)")
	g := 14.134725
	den := 16
	var betas []float64
	for k := 1; k < den; k++ {
		betas = append(betas, float64(k)/float64(den))
	}
	sort.Float64s(betas)
	fmt.Println("   β = k/16      precio          ¿su pareja?")
	peorRaiz := 0.0
	minFuera := math.Inf(1)
	for _, b := range betas {
		p := precio(b, g)
		par := fmt.Sprintf("%.4f", 1-b)
		if math.Abs(b-0.5) < 1e-15 {
			par = "★ ELLA MISMA"
			if math.Abs(p) > peorRaiz {
				peorRaiz = math.Abs(p)
			}
		} else if p < minFuera {
			minFuera = p
		}
		fmt.Printf("   %.4f      %+.6e     %s\n", b, p, par)
	}
	fmt.Printf("   → el precio en la raíz es %.1e (cero) y el más barato de todos los demás es %.2e\n",
		peorRaiz, minFuera)
	fmt.Println("   → el árbol entero paga; solo la raíz entra gratis")

	fmt.Println("\n════════ LO QUE DIJO EL CAPITÁN, MEDIDO ════════")
	fmt.Println("«Partí el 7 infinitas veces: te queda un mar profundo e inalcanzable. Y la relación se")
	fmt.Println("hereda: 6 —½— 7 —½— 8, entonces 6 —½— 8; entonces +∞ —½— −∞; y en la dimensión 0,")
	fmt.Println("0 —½— 1.» Todo eso está medido acá y es exacto.")
	fmt.Println("\nY el flash entrega un remate que el laboratorio no había escrito: el mar infinito de")
	fmt.Println("cortes alcanza CUALQUIER lugar donde una perla podría vivir, y al dar vuelta el mundo")
	fmt.Println("todos esos lugares se emparejan de a dos… menos UNO. El primer corte, el 1/2, es el")
	fmt.Println("único que se mira al espejo y se ve a sí mismo. Y ES GRATIS: el único nodo del árbol")
	fmt.Println("infinito donde el precio es cero.")
	fmt.Println("\nLA HIPÓTESIS, DICHA CON EL ÁRBOL DEL CAPITÁN: de todos los infinitos lugares que")
	fmt.Println("nacen de partir y partir, el libro elige SIEMPRE el único que es su propia pareja")
	fmt.Println("y el único que no cobra. Falta demostrar que no puede elegir otro. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🌳 EL ÁRBOL DE LAS MITADES — un solo nodo es su propia pareja, y es gratis</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"partí el 7 infinitas veces: un mar profundo e inalcanzable… 0 —½— 1" — el capitán · el árbol infinito y su único punto fijo</text>`,
		W, H, W, H, W/2, W/2)

	// the tree drawing
	tx, ty, tw := 90.0, 130.0, 1320.0
	fmt.Fprintf(&b, `<rect x="60" y="105" width="1380" height="330" rx="10" fill="#0d2547" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="%.0f" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">PARTIR Y PARTIR — cada corte duplica los lugares, y el espejo los empareja de a dos</text>`, W/2)
	for d := 1; d <= 5; d++ {
		y := ty + 40 + float64(d)*48
		den := 1 << d
		for k := 1; k < den; k++ {
			bt := float64(k) / float64(den)
			x := tx + 30 + bt*(tw-60)
			col, r := "#7fb2ff", 4.0
			if math.Abs(bt-0.5) < 1e-15 {
				col, r = "#ffd166", 7.0
			}
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="%s"/>`, x, y, r, col)
		}
		fmt.Fprintf(&b, `<text x="%.0f" y="%.1f" font-size="11" fill="#8fa8c7">corte %d</text>`, tx-20, y+4, d)
	}
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#ffd166" stroke-width="1.5" stroke-dasharray="6 5"/>
<text x="%.1f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#ffd166">el 1/2 — el único que se mira al espejo y se ve a sí mismo</text>
<text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" fill="#ff8fa0">estaca 0</text>
<text x="%.1f" y="%.0f" font-size="12" text-anchor="middle" fill="#7fb2ff">estaca 1</text>`,
		tx+30+0.5*(tw-60), ty+50, tx+30+0.5*(tw-60), ty+300,
		tx+30+0.5*(tw-60), ty+320,
		tx+30, ty+50, tx+tw-30, ty+50)

	// left bottom: the unreachable sea
	fmt.Fprintf(&b, `<rect x="60" y="465" width="660" height="240" rx="10" fill="#102a10" stroke="#ffd166" stroke-width="1.5"/>
<text x="390" y="499" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL MAR INALCANZABLE — partir el 7</text>
<text x="390" y="527" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">corte          lo que queda de 7</text>`)
	for i, k := range []int{1, 3, 10, 50, 200, 1000} {
		fmt.Fprintf(&b, `<text x="390" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%4d           %.4e</text>`,
			552.0+float64(i)*24, k, 7.0/math.Pow(2, float64(k)))
	}
	fmt.Fprintf(&b, `<text x="390" y="700" font-size="12.5" text-anchor="middle" fill="#7fd7a8">se acerca siempre y no llega nunca: la dimensión 0</text>`)

	// right bottom: the price
	fmt.Fprintf(&b, `<rect x="760" y="465" width="680" height="240" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="1100" y="499" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">EL PRECIO EN CADA NODO — solo la raíz entra gratis</text>
<text x="1100" y="527" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β = k/16          precio (N+S)/2 − 1</text>`)
	for i, k := range []int{2, 4, 6, 8, 10, 12} {
		bt := float64(k) / 16
		p := precio(bt, g)
		col := "#ff8fa0"
		extra := ""
		if k == 8 {
			col = "#7fd7a8"
			extra = "   ★ GRATIS"
		}
		fmt.Fprintf(&b, `<text x="1100" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="%s">%.4f            %+.4e%s</text>`,
			552.0+float64(i)*24, col, bt, p, extra)
	}
	fmt.Fprintf(&b, `<text x="1100" y="700" font-size="12.5" text-anchor="middle" fill="#dce8f7">el árbol entero paga; el 1/2 es el único nodo con precio cero</text>`)

	fmt.Fprintf(&b, `<rect x="60" y="735" width="1380" height="160" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="771" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA HIPÓTESIS, DICHA CON EL ÁRBOL DEL CAPITÁN</text>
<text x="%.0f" y="805" font-size="14.5" text-anchor="middle" fill="#dce8f7">De partir y partir nacen infinitos lugares donde una perla podría vivir. Al dar vuelta el mundo, todos se emparejan de a dos… menos uno.</text>
<text x="%.0f" y="837" font-size="15.5" text-anchor="middle" fill="#7fd7a8">El libro elige SIEMPRE ese único: el que es su propia pareja, y el único que no cobra.</text>
<text x="%.0f" y="866" font-size="13.5" text-anchor="middle" fill="#ff8fa0">Falta demostrar que no puede elegir otro. Todavía no.</text>
<text x="%.0f" y="890" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("arbol-de-mitades.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: arbol-de-mitades.svg")
}
