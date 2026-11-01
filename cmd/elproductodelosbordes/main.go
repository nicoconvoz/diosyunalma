package main

// EL PRODUCTO DE LOS BORDES - the captain: "and the logarithm of the right edge
// TIMES the logarithm of the left edge - what does that give?"
//
// In his ball picture, row n is centred on the line, so its edges sit at
//
//      right edge = +n/2        left edge = -n/2
//
// The left edge is NEGATIVE. And the logarithm of a negative number is not a
// real number - it is complex, and it costs exactly i pi:
//
//      ln(-a) = ln(a) + i pi           (Euler)
//
// So the captain's innocent product is
//
//      ln(n/2) * [ ln(n/2) + i pi ] = ln^2(n/2)  +  i pi ln(n/2)
//
// Which means: writing u = ln(n/2), the answer is the point (u^2, pi u) in the
// complex plane. Eliminate u and you get
//
//      Re  =  Im^2 / pi^2
//
// a PARABOLA. The captain's ball drawing, asked this one question, lands exactly
// on the two parabolas he opened this whole campaign with in F391 and F392.
//
// Verdicts fixed BEFORE running:
//   V1 the product must be complex, with imaginary part exactly pi ln(n/2).
//   V2 the locus must be exactly the parabola Re = Im^2/pi^2, to machine
//      precision, for every n.
//   V3 n = 2 must give exactly ZERO, and must be the only n that does, because
//      its right edge lands exactly on 1 and ln 1 = 0.
//   V4 the argument must decay like arctan(pi / ln(n/2)), so the product turns
//      almost purely real as the numbers grow.
//   V5 (TRAP 29, declared before computing) the product depends ONLY on the
//      SIZE of n, never on its factorisation, so it CANNOT tell a prime from
//      its neighbours. Control: compare the statistic at primes against random
//      numbers of the same size. It must come out 1, i.e. nothing.
//   V6 but ln^2 itself IS real and important, in two places we can measure:
//      Schoenfeld's RH bound ln^2(x)/(8 pi), and the twin density 2C2/ln^2(x).
//      Fresh check: a window of length ln^2(x) must contain about ln(x) primes.

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"path/filepath"
	"strings"
)

func criba(n int) ([]bool, []int) {
	es := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		es[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if !es[p] {
			continue
		}
		for m := p * p; m <= n; m += p {
			es[m] = false
		}
	}
	var ps []int
	for i := 2; i <= n; i++ {
		if es[i] {
			ps = append(ps, i)
		}
	}
	return es, ps
}

// producto is the captain's operation: ln(right edge) times ln(left edge).
func producto(n float64) complex128 {
	der := complex(math.Log(n/2), 0)       // ln(+n/2)
	izq := complex(math.Log(n/2), math.Pi) // ln(-n/2) = ln(n/2) + i pi
	return der * izq
}

func main() {
	fmt.Println("EL PRODUCTO DE LOS BORDES — F419")
	fmt.Println("el capitán: «¿y el logaritmo del borde derecho POR el logaritmo")
	fmt.Println("del borde izquierdo, qué da?»")
	fmt.Println()

	fmt.Println("PRIMERO, lo que nadie espera: el borde izquierdo es NEGATIVO.")
	fmt.Println("Y el logaritmo de un negativo no es un número real. Cuesta iπ:")
	fmt.Println()
	fmt.Println("        ln(−a) = ln(a) + iπ        (Euler)")
	fmt.Println()
	fmt.Println("Así que su pregunta obliga al plano complejo a aparecer,")
	fmt.Println("y aparece por la puerta de Euler. El producto queda:")
	fmt.Println()
	fmt.Println("        ln(n/2) · [ln(n/2) + iπ] = ln²(n/2) + iπ·ln(n/2)")
	fmt.Println()

	// ---------------- (1) la tabla ----------------
	fmt.Println("(1) LA TABLA — el producto, número por número")
	fmt.Printf("    %5s %10s %10s %14s %14s %12s\n",
		"n", "borde der", "borde izq", "parte real", "parte imag", "ángulo")
	for _, n := range []float64{2, 3, 5, 7, 11, 100, 1000000} {
		z := producto(n)
		ang := cmplx.Phase(z) * 180 / math.Pi
		fmt.Printf("    %5.0f %10.1f %10.1f %14.6f %14.6f %11.3f°\n",
			n, n/2, -n/2, real(z), imag(z), ang)
	}
	fmt.Println()
	fmt.Println("    → EL 2 DA CERO EXACTO. Y es el único que puede darlo, porque su")
	fmt.Println("      borde derecho cae justo en 1, y el logaritmo de 1 es cero.")
	fmt.Println("      El 2 vuelve a ser el único, igual que en F418 con el hueco.")
	fmt.Println()

	// ---------------- (2) la parábola ----------------
	fmt.Println("(2) EL LUGAR GEOMÉTRICO — y acá se cierra un círculo")
	fmt.Println("    si u = ln(n/2), el producto es el punto (u², πu).")
	fmt.Println("    Despejando u: Re = Im² / π². ESO ES UNA PARÁBOLA.")
	fmt.Println()
	peor := 0.0
	for n := 2.0001; n < 1e6; n *= 1.03 {
		z := producto(n)
		pred := imag(z) * imag(z) / (math.Pi * math.Pi)
		d := math.Abs(real(z) - pred)
		if d > peor {
			peor = d
		}
	}
	fmt.Printf("    verificado en todo el rango hasta un millón: peor desvío %.3e\n", peor)
	fmt.Println("    → precisión de máquina. Es la parábola, exacta.")
	fmt.Println()
	fmt.Println("    ⚡ Y MIRE DÓNDE TERMINAMOS: el capitán abrió esta campaña en F391")
	fmt.Println("      y F392 preguntando por DOS PARÁBOLAS que se tocan en el cero.")
	fmt.Println("      Su dibujo de bolitas, con esta sola pregunta, vuelve a la")
	fmt.Println("      parábola — y los dos bordes, el de arriba y el de abajo del")
	fmt.Println("      eje real, son las DOS ramas de esa misma parábola.")
	fmt.Println()

	// ---------------- (3) el ángulo ----------------
	fmt.Println("(3) EL ÁNGULO — cómo el producto se vuelve casi real")
	fmt.Printf("    %12s %16s %16s\n", "n", "ángulo medido", "arctan(π/ln(n/2))")
	for _, n := range []float64{10, 1e3, 1e6, 1e12, 1e100} {
		z := producto(n)
		med := cmplx.Phase(z) * 180 / math.Pi
		pred := math.Atan(math.Pi/math.Log(n/2)) * 180 / math.Pi
		fmt.Printf("    %12.0e %15.4f° %15.4f°\n", n, med, pred)
	}
	fmt.Println()
	fmt.Println("    → el ángulo muere como π/ln(n/2): despacísimo, pero muere.")
	fmt.Println("      Cuanto más grande el número, más REAL es el producto.")
	fmt.Println()

	// ---------------- (4) la trampa y su control ----------------
	fmt.Println("(4) TRAMPA 29, declarada antes de calcular")
	fmt.Println("    el producto depende SOLO DEL TAMAÑO de n, nunca de cómo se")
	fmt.Println("    factoriza. Así que NO PUEDE distinguir un primo de su vecino.")
	fmt.Println("    Control: la parte fraccionaria de la parte real, en primos")
	fmt.Println("    contra números cualesquiera del mismo barrio.")
	const N = 2000000
	es, ps := criba(N)
	sumP, cP := 0.0, 0
	for _, p := range ps {
		if p < 1000 {
			continue
		}
		f := real(producto(float64(p)))
		sumP += f - math.Floor(f)
		cP++
	}
	sumQ, cQ := 0.0, 0
	for n := 1000; n <= N; n++ {
		f := real(producto(float64(n)))
		sumQ += f - math.Floor(f)
		cQ++
	}
	mp, mq := sumP/float64(cP), sumQ/float64(cQ)
	fmt.Printf("    media de la fracción en %d primos:            %.6f\n", cP, mp)
	fmt.Printf("    media de la fracción en %d números cualquiera: %.6f\n", cQ, mq)
	fmt.Printf("    señal / silencio: %.4f   (uno = NADA)\n", mp/mq)
	fmt.Println("    → CONFIRMADO: el producto NO sabe nada de primalidad.")
	fmt.Println("      Es una función del tamaño y nada más. Trampa 29 enterrada.")
	fmt.Println()

	// ---------------- (5) pero ln cuadrado SÍ importa ----------------
	fmt.Println("(5) PERO LA PARTE REAL, ln², SÍ IMPORTA — y en dos lugares de verdad")
	fmt.Println()
	fmt.Println("    a) es la forma EXACTA de la cota de la Hipótesis (Schoenfeld 1976):")
	fmt.Println("       |ψ(x) − x| < √x · ln²x / (8π)   — la medimos en F416")
	fmt.Println()
	fmt.Println("    b) es el denominador de la densidad de primos gemelos:")
	fmt.Println("       2C₂ / ln²x   — la medimos en F418")
	fmt.Println()
	fmt.Println("    c) y una medición NUEVA, para no repetir: una ventana de largo")
	fmt.Println("       ln²x tiene que atrapar alrededor de ln x primos.")
	fmt.Printf("\n    %12s %14s %16s %14s\n", "cerca de x", "largo ln²x", "primos medidos", "ln x")
	for _, x := range []float64{1e4, 1e5, 1e6} {
		lx := math.Log(x)
		w := int(lx * lx)
		// average over many windows around x
		suma, vent := 0, 0
		for c := int(x * 0.9); c < int(x*1.1); c += w {
			cnt := 0
			for n := c; n < c+w && n <= N; n++ {
				if es[n] {
					cnt++
				}
			}
			suma += cnt
			vent++
		}
		fmt.Printf("    %12.0e %14d %16.4f %14.4f\n",
			x, w, float64(suma)/float64(vent), lx)
	}
	fmt.Println()
	fmt.Println("    → el largo ln²x es la ventana natural: atrapa ln x primos.")
	fmt.Println("      Su producto de bordes ES esa ventana, escrita al revés.")
	fmt.Println()

	fmt.Println("LA RESPUESTA CORTA, capitán:")
	fmt.Println("  da un NÚMERO COMPLEJO, y los complejos entran obligados por el")
	fmt.Println("  borde izquierdo, que es negativo. La parte imaginaria es π·ln(n/2)")
	fmt.Println("  y la real es ln²(n/2), y juntas dibujan una PARÁBOLA exacta.")
	fmt.Println("  El 2 es el único que da cero. Y el ln² que quedó en la parte real")
	fmt.Println("  es la forma de la cota de Riemann y del denominador de los gemelos.")
	fmt.Println()
	fmt.Println("  Pero el producto en sí no distingue un primo de su vecino:")
	fmt.Println("  medido, señal sobre silencio da UNO. Eso también hay que decirlo.")
	fmt.Println()

	fmt.Println("HONESTO:")
	fmt.Println("  · ln(−1) = iπ es de Euler, 1748, y es de lo más famoso que hay.")
	fmt.Println("  · la cota con ln² es de Schoenfeld, 1976; los gemelos, de Hardy")
	fmt.Println("    y Littlewood, 1923.")
	fmt.Println("  · el producto de logaritmos NO es una operación natural en este")
	fmt.Println("    oficio: el logaritmo existe para convertir productos en SUMAS.")
	fmt.Println("    Forzarlo da geometría linda y aritmética vacía.")
	fmt.Println("  · lo propio: haber notado que el dibujo de bolitas vuelve a la")
	fmt.Println("    parábola de F391, y haber medido que no distingue primos.")
	fmt.Println()
	fmt.Println("Todavía no.")

	lamina(mp / mq)
	fmt.Println("\nlámina: galeria/laminas/10-el-telar/el-producto-de-los-bordes.svg")
}

func lamina(ss float64) {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1060 720" width="1060" height="720">`)
	b.WriteString(`<rect width="1060" height="720" fill="#0d1017"/>`)
	txt := func(x, y, sz float64, col, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="Georgia,serif" font-size="%.0f" fill="%s" text-anchor="%s">%s</text>`,
			x, y, sz, col, anc, s)
	}
	txt(530, 40, 24, "#e8e2d0", "middle", "EL PRODUCTO DE LOS BORDES")
	txt(530, 66, 12, "#8a8676", "middle", "«¿y el logaritmo del borde derecho por el logaritmo del borde izquierdo, qué da?»")

	// the forced complex plane
	fmt.Fprintf(&b, `<rect x="20" y="88" width="500" height="200" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 114, 13, "#f2c14e", "middle", "el borde izquierdo es NEGATIVO")
	txt(270, 146, 16, "#e8e2d0", "middle", "ln(−a) = ln(a) + iπ")
	txt(270, 168, 10, "#8a8676", "middle", "Euler, 1748")
	txt(270, 202, 11, "#c9c3b0", "middle", "así que su pregunta OBLIGA a los complejos a aparecer,")
	txt(270, 220, 11, "#c9c3b0", "middle", "y aparecen por la puerta de Euler")
	txt(270, 254, 14, "#6ee7a0", "middle", "ln²(n/2) + iπ·ln(n/2)")
	txt(270, 276, 10, "#8a8676", "middle", "el 2 da CERO exacto — su borde derecho cae justo en 1")

	// the parabola
	fmt.Fprintf(&b, `<rect x="540" y="88" width="500" height="200" fill="#131826" stroke="#2a3040"/>`)
	txt(790, 114, 13, "#f2c14e", "middle", "el lugar geométrico: una PARÁBOLA exacta")
	// draw the parabola Re = Im^2/pi^2, horizontal axis
	var d strings.Builder
	for i := -34; i <= 34; i++ {
		v := float64(i) / 10
		re := v * v / (math.Pi * math.Pi) * 3
		x := 600 + re*130
		y := 200 - v*26
		if i == -34 {
			fmt.Fprintf(&d, "M %.1f %.1f", x, y)
		} else {
			fmt.Fprintf(&d, " L %.1f %.1f", x, y)
		}
	}
	fmt.Fprintf(&b, `<path d="%s" fill="none" stroke="#f2c14e" stroke-width="1.8"/>`, d.String())
	fmt.Fprintf(&b, `<line x1="600" y1="200" x2="1010" y2="200" stroke="#5a6070" stroke-dasharray="3,3"/>`)
	fmt.Fprintf(&b, `<circle cx="600" cy="200" r="4" fill="#6ee7a0"/>`)
	txt(596, 216, 10, "#6ee7a0", "end", "el 2")
	txt(790, 262, 12, "#e8e2d0", "middle", "Re = Im² / π²")
	txt(790, 280, 10, "#8a8676", "middle", "desvío medido: precisión de máquina")

	// the circle closing
	fmt.Fprintf(&b, `<rect x="20" y="304" width="1020" height="86" fill="#0f1a14" stroke="#2a5a40"/>`)
	txt(530, 332, 13, "#6ee7a0", "middle", "Y MIRE DÓNDE TERMINAMOS")
	txt(530, 358, 12, "#e8e2d0", "middle", "el capitán abrió esta campaña en F391 preguntando por DOS PARÁBOLAS que se tocan en el cero")
	txt(530, 380, 12, "#f2c14e", "middle", "su dibujo de bolitas vuelve a la parábola — y los dos bordes son las DOS ramas de la misma")

	// the trap
	fmt.Fprintf(&b, `<rect x="20" y="406" width="500" height="180" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(270, 432, 13, "#ff6b6b", "middle", "trampa 29, declarada antes de calcular")
	txt(270, 458, 11, "#c9c3b0", "middle", "el producto depende SOLO del tamaño de n,")
	txt(270, 476, 11, "#c9c3b0", "middle", "nunca de cómo se factoriza")
	txt(270, 508, 12, "#e8e2d0", "middle", "control en 148.765 primos contra 2 millones")
	txt(270, 530, 12, "#e8e2d0", "middle", "de números cualesquiera:")
	txt(270, 558, 16, "#ff6b6b", "middle", fmt.Sprintf("señal / silencio = %.4f", ss))
	txt(270, 578, 10, "#8a8676", "middle", "uno = NADA. No sabe de primalidad. Enterrada.")

	// but ln^2 matters
	fmt.Fprintf(&b, `<rect x="540" y="406" width="500" height="180" fill="#131826" stroke="#2a3040"/>`)
	txt(790, 432, 13, "#f2c14e", "middle", "pero la parte real, ln², SÍ importa")
	txt(790, 462, 11, "#c9c3b0", "middle", "a · la cota de la Hipótesis:  √x · ln²x / (8π)")
	txt(790, 486, 11, "#c9c3b0", "middle", "b · el denominador de los gemelos:  2C₂ / ln²x")
	txt(790, 510, 11, "#c9c3b0", "middle", "c · una ventana de largo ln²x atrapa ln x primos")
	txt(790, 544, 12, "#6ee7a0", "middle", "medido cerca de un millón: 13,7198 primos")
	txt(790, 566, 11, "#8a8676", "middle", "contra ln x = 13,8155. La ventana natural.")

	// closing
	fmt.Fprintf(&b, `<rect x="20" y="602" width="1020" height="102" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 628, 13, "#e8e2d0", "middle", "LA RESPUESTA CORTA")
	txt(530, 654, 12, "#f2c14e", "middle", "da un número COMPLEJO, obligado por el borde izquierdo; parte real ln²(n/2), parte imaginaria π·ln(n/2), y juntas una parábola")
	txt(530, 678, 11, "#c9c3b0", "middle", "el ln² que queda es la forma de la cota de Riemann. Pero el producto en sí no distingue un primo de su vecino, y eso también se dice.")
	txt(530, 698, 11, "#8a8676", "middle", "Todavía no.")

	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-producto-de-los-bordes.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
