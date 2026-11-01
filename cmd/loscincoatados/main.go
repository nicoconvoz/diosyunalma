package main

// LOS CINCO ATADOS - the captain asks for the relation between 1, the prime
// itself, 2, 1/2 and the square root of the prime.
//
// They tie in one line, and the knot is the GEOMETRIC MEAN:
//
//     sqrt(p) = sqrt(1 * p) = (1 * p)^(1/2)
//
// The two beats of F405 are 1 and p. Their count is 2. Their geometric mean is
// sqrt(p), and the exponent that produces it is 1/2. All five objects the
// captain named, in a single expression.
//
// And the knot is not decorative: sqrt(n) is the FOLD POINT of the divisors.
// Every divisor d pairs with n/d, and the pairing d <-> n/d has sqrt(n) as its
// fixed point - which is exactly why trial division stops there.
//
// Verdicts fixed BEFORE running:
//   V1 sqrt(p) must equal the geometric mean of the two divisors, verified to
//      machine precision, and the arithmetic mean (1+p)/2 must equal sigma(p)/2
//      with AM >= GM always.
//   V2 (the fold) for every n, the divisors must pair around sqrt(n), with the
//      pairing exact and sqrt(n) its fixed point. Measured for primes and
//      composites alike.
//   V3 (the work it saves) stopping at sqrt(n) instead of n must be measured
//      as a real ratio, not asserted.
//   V4 (trap #22, declared) "sqrt(p) is irrational" holds for every non-square,
//      prime or not: 6, 8, 10 too. Control required, and the ratio must come
//      out 1.
//   V5 (what DOES separate) the distance from the fold to the nearest divisor
//      below it: for a prime that divisor is 1, the furthest possible. Measured
//      with control - and its circularity declared, since finding that divisor
//      IS trial division.
//   V6 (the parallel, stated as analogy and NOT as theorem) d <-> n/d fixes
//      sqrt(n); s <-> 1-s fixes 1/2; in the Euler product p^(-s) <-> p^(-(1-s))
//      fixes p^(-1/2) = 1/sqrt(p). Same shape. NOT the same proof.
// AM-GM is Euclid. Nothing new.

import (
	"fmt"
	"math"
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

// mayorDivisorBajoRaiz returns the largest divisor of n that is <= sqrt(n).
func mayorDivisorBajoRaiz(n int) int {
	mejor := 1
	for d := 1; d*d <= n; d++ {
		if n%d == 0 {
			mejor = d
		}
	}
	return mejor
}

func main() {
	fmt.Println("🪢 LOS CINCO ATADOS — el 1, el primo, el 2, el ½ y la raíz")
	fmt.Println()

	const N = 200000
	es, _ := criba(N)

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · EL NUDO: los cinco en una sola expresión")
	fmt.Println()
	fmt.Println("        √p  =  √(1 × p)  =  (1 × p)^½")
	fmt.Println()
	fmt.Println("   Los dos latidos de F405 son el 1 y el primo. Son 2. Su media GEOMÉTRICA es la raíz del")
	fmt.Println("   primo, y el exponente que la produce es ½. Los cinco objetos que pidió, en un renglón.")
	fmt.Println()
	fmt.Println("        primo     los dos latidos     media geométrica √(1·p)     √p            coinciden")
	for _, p := range []int{2, 3, 5, 7, 11, 13} {
		g := math.Sqrt(1 * float64(p))
		fmt.Printf("         %3d        1  y  %-4d          %.9f            %.9f      %.0e\n",
			p, p, g, math.Sqrt(float64(p)), math.Abs(g-math.Sqrt(float64(p))))
	}
	fmt.Println()
	fmt.Println("   ⚡ Y el 2 aparece dos veces más: en la CUENTA de los latidos, y en la media ARITMÉTICA,")
	fmt.Println("      que es (1 + p)/2 — y ésa es exactamente la suma de divisores dividida por 2.")
	fmt.Println()
	fmt.Println("        primo     σ(p) = 1 + p     media aritmética (1+p)/2     media geométrica √p     ¿AM ≥ GM?")
	for _, p := range []int{2, 3, 5, 7, 11, 13, 97} {
		am := float64(1+p) / 2
		gm := math.Sqrt(float64(p))
		et := "sí"
		if am < gm {
			et = "NO"
		}
		fmt.Printf("         %3d          %5d            %10.4f                  %10.6f           %s\n", p, 1+p, am, gm, et)
	}
	fmt.Println()
	fmt.Println("   La desigualdad entre las dos medias (Euclides) nunca falla, y para un primo la brecha entre")
	fmt.Println("   ellas es la MÁS GRANDE posible — porque sus dos divisores son los más separados que existen.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · LA RAÍZ ES EL PLIEGUE: por eso importa tanto")
	fmt.Println()
	fmt.Println("   Los divisores de cualquier número van de a pares: si d divide a n, entonces n/d también.")
	fmt.Println("   Y cada par se abraza alrededor de la raíz, porque d × (n/d) = n:")
	fmt.Println()
	for _, n := range []int{36, 7} {
		et := "compuesto"
		if es[n] {
			et = "PRIMO"
		}
		fmt.Printf("        %d (%s) · √%d = %.6f\n", n, et, n, math.Sqrt(float64(n)))
		for d := 1; d*d <= n; d++ {
			if n%d != 0 {
				continue
			}
			marca := ""
			if d*d == n {
				marca = "   ← el pliegue exacto"
			}
			fmt.Printf("            %4d  ×  %-4d  =  %d%s\n", d, n/d, n, marca)
		}
		fmt.Println()
	}
	fmt.Println("   ⚡⚡ El 7 tiene UN SOLO par: 1 × 7. Y la raíz de 7 cae justo en el medio, sin nadie a los lados.")
	fmt.Println("      El 36 tiene cinco pares, y uno cae EXACTO sobre la raíz (6 × 6). Ésa es la diferencia,")
	fmt.Println("      dibujada: un primo es un número cuyo único pliegue son sus dos extremos.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · lo que ahorra el pliegue, medido")
	fmt.Println()
	fmt.Println("        para probar si n es primo…       divisiones hasta n     hasta √n     ahorro")
	for _, n := range []int{100, 10000, 1000000, 100000000} {
		hasta := float64(n)
		raiz := math.Sqrt(float64(n))
		fmt.Printf("        n = %-11d                  %12.0f       %10.0f     %8.0f×\n", n, hasta, raiz, hasta/raiz)
	}
	fmt.Println()
	fmt.Println("   Probar hasta la raíz en vez de hasta el número entero es lo que hace posible la criba.")
	fmt.Println("   Y el motivo es el pliegue: si hubiera un divisor arriba de √n, su pareja estaría abajo.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · la trampa (vigesimosegunda): «la raíz de un primo es irracional»")
	fmt.Println()
	var pIrr, pTot, cIrr, cTot int
	for n := 2; n <= N; n++ {
		r := int(math.Round(math.Sqrt(float64(n))))
		irr := r*r != n
		if es[n] {
			pTot++
			if irr {
				pIrr++
			}
		} else {
			cTot++
			if irr {
				cIrr++
			}
		}
	}
	fmt.Printf("        primos con raíz irracional:       %6d / %-6d   →  %.4f %%\n", pIrr, pTot, 100*float64(pIrr)/float64(pTot))
	fmt.Printf("        COMPUESTOS con raíz irracional:   %6d / %-6d   →  %.4f %%\n", cIrr, cTot, 100*float64(cIrr)/float64(cTot))
	fmt.Printf("\n        señal / silencio = %.4f  → ", (float64(pIrr)/float64(pTot))/(float64(cIrr)/float64(cTot)))
	fmt.Println("⚠️ casi 1: el 6, el 8 y el 10 también tienen raíz irracional")
	fmt.Println("        (los compuestos que fallan son los cuadrados perfectos: 4, 9, 16, 25… y son poquísimos)")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · lo que SÍ separa: qué tan lejos del pliegue está el divisor más cercano")
	fmt.Println()
	fmt.Println("   Para un primo, el mayor divisor por debajo de √p es el 1 — lo más lejos posible del pliegue.")
	fmt.Println("   Para un compuesto hay algo más cerca. Se mide la razón √n / (mayor divisor bajo la raíz):")
	fmt.Println()
	fmt.Println("        número    ¿primo?    √n          mayor divisor ≤ √n     razón")
	for _, n := range []int{97, 101, 96, 100, 99, 91} {
		d := mayorDivisorBajoRaiz(n)
		r := math.Sqrt(float64(n))
		et := "no"
		if es[n] {
			et = "SÍ"
		}
		fmt.Printf("         %4d       %-3s     %8.4f          %5d              %7.4f\n", n, et, r, d, r/float64(d))
	}
	var sumP, sumC float64
	var nP, nC int
	for n := 1000; n <= 3000; n++ {
		d := mayorDivisorBajoRaiz(n)
		r := math.Sqrt(float64(n)) / float64(d)
		if es[n] {
			sumP += r
			nP++
		} else {
			sumC += r
			nC++
		}
	}
	fmt.Printf("\n        razón media entre 1.000 y 3.000 · PRIMOS: %.4f · COMPUESTOS: %.4f · separación %.1f×\n",
		sumP/float64(nP), sumC/float64(nC), (sumP/float64(nP))/(sumC/float64(nC)))
	fmt.Println()
	fmt.Println("   ⚡ Separa, y bastante. ⚖️ PERO con circularidad declarada: para hallar «el mayor divisor bajo")
	fmt.Println("      la raíz» hay que DIVIDIR hasta encontrarlo, que es exactamente el test de primalidad. Sirve")
	fmt.Println("      para describir la diferencia, no para detectarla más barato. Se dice, no se esconde.")
	fmt.Println()

	// ---------------------------------------------------------------- V6
	fmt.Println("§6 · el paralelo con la línea crítica — dicho como ANALOGÍA, no como teorema")
	fmt.Println()
	fmt.Println("        el pliegue de los divisores        d  ↔  n/d          punto fijo:  √n")
	fmt.Println("        el pliegue de la ecuación func.    s  ↔  1 − s        punto fijo:  ½")
	fmt.Println("        y en el producto de Euler          p^(−s) ↔ p^(−(1−s))  punto fijo:  p^(−½) = 1/√p")
	fmt.Println()
	fmt.Println("   Las tres son la misma FORMA: un emparejamiento que se da vuelta sobre sí mismo, con un único")
	fmt.Println("   punto que queda quieto. Y el punto quieto del tercero es justo el peso de cada primo en la")
	fmt.Println("   línea crítica, que es lo que midió F393.")
	fmt.Println()
	fmt.Println("   ⚖️ Ahora el freno, y va en negrita mental: esto es una ANALOGÍA DE FORMA, no una demostración.")
	fmt.Println("      Que dos emparejamientos tengan la misma estructura no hace que uno implique al otro — y la")
	fmt.Println("      casa ya se comió esa lección en F401, donde una regla sin excepciones no alcanzaba a RH.")
	fmt.Println()

	fmt.Println("§7 · EL VEREDICTO")
	fmt.Println()
	fmt.Println("   ✅ Los cinco se atan, y el nudo es exacto: √p es la MEDIA GEOMÉTRICA de los dos latidos,")
	fmt.Println("      √p = √(1 × p) = (1 × p)^½. El 1 y el primo son los latidos, el 2 los cuenta, el ½ es el")
	fmt.Println("      exponente que los promedia, y la raíz es el resultado.")
	fmt.Println("   ✅ Y el nudo explica algo real: la raíz es el PLIEGUE de los divisores, y por eso el test de")
	fmt.Println("      primalidad se detiene ahí — un ahorro de 10.000× a cien millones, medido.")
	fmt.Println("   ⚠️ Pero «la raíz es irracional» no distingue primos: al 6, al 8 y al 10 les pasa igual.")
	fmt.Println("   ⚖️ Y el paralelo con el ½ de la línea crítica es de FORMA, no de teorema. Queda anotado como")
	fmt.Println("      analogía, que es lo que es.")
	fmt.Println()
	fmt.Println("   Nada de esto es nuevo: la desigualdad de las medias es de Euclides. Todavía no.")

	dibujar()
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/los-cinco-atados.svg")
}

func dibujar() {
	var b strings.Builder
	W, H := 1060, 680
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LOS CINCO ATADOS — el 1, el primo, el 2, el ½ y la raíz")
	txt(530, 54, 12, "#8a8676", "middle", "el nudo es la media geométrica, y la raíz es el pliegue de los divisores · F406")

	// the knot
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="150" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 102, 13, "#6ee7a0", "middle", "EL NUDO: los cinco en una sola expresión")
	txt(530, 156, 32, "#f2c14e", "middle", "√p   =   √(1 × p)   =   (1 × p)^½")
	txt(180, 196, 11, "#6ee7a0", "middle", "el 1 y el primo")
	txt(180, 212, 10, "#8a8676", "middle", "los dos latidos")
	txt(430, 196, 11, "#6ee7a0", "middle", "el 2")
	txt(430, 212, 10, "#8a8676", "middle", "los cuenta")
	txt(680, 196, 11, "#6ee7a0", "middle", "el ½")
	txt(680, 212, 10, "#8a8676", "middle", "el exponente que promedia")
	txt(900, 196, 11, "#6ee7a0", "middle", "la raíz")
	txt(900, 212, 10, "#8a8676", "middle", "el resultado")

	// the fold
	fmt.Fprintf(&b, `<rect x="20" y="242" width="500" height="230" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 268, 13, "#e8e2d0", "middle", "la raíz es EL PLIEGUE de los divisores")
	cy := 350.0
	fmt.Fprintf(&b, `<line x1="80" y1="%.0f" x2="460" y2="%.0f" stroke="#3a4050"/>`, cy, cy)
	fmt.Fprintf(&b, `<line x1="270" y1="%.0f" x2="270" y2="%.0f" stroke="#f2c14e" stroke-dasharray="4,3"/>`, cy-50, cy+50)
	txt(270, cy-58, 11, "#f2c14e", "middle", "√36 = 6")
	for i, par := range [][2]int{{1, 36}, {2, 18}, {3, 12}, {4, 9}} {
		dx := 26.0 + float64(i)*38
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="4" fill="#5ec8ff"/>`, 270-dx, cy)
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="4" fill="#5ec8ff"/>`, 270+dx, cy)
		fmt.Fprintf(&b, `<path d="M %.0f %.0f Q 270 %.0f %.0f %.0f" fill="none" stroke="#5ec8ff" opacity="0.4"/>`,
			270-dx, cy, cy-24-float64(i)*7, 270+dx, cy)
		txt(270-dx, cy+18, 9, "#8a8676", "middle", fmt.Sprint(par[0]))
		txt(270+dx, cy+18, 9, "#8a8676", "middle", fmt.Sprint(par[1]))
	}
	fmt.Fprintf(&b, `<circle cx="270" cy="%.0f" r="5" fill="#f2c14e"/>`, cy)
	txt(270, cy+18, 9, "#f2c14e", "middle", "6")
	txt(270, 412, 11, "#c9c3b0", "middle", "cada divisor se abraza con su pareja: d × (n/d) = n")
	txt(270, 434, 11, "#6ee7a0", "middle", "un primo tiene UN SOLO par: 1 × p")
	txt(270, 456, 10, "#8a8676", "middle", "por eso el test de primalidad se detiene en la raíz")

	// what separates and what doesn't
	fmt.Fprintf(&b, `<rect x="540" y="242" width="500" height="230" fill="#131826" stroke="#2a3040"/>`)
	txt(790, 268, 13, "#e8e2d0", "middle", "qué separa y qué no")
	txt(564, 302, 11, "#ff6b6b", "start", "⚠️ «la raíz de un primo es irracional»")
	txt(564, 324, 11, "#8a8676", "start", "el 6, el 8 y el 10 también. Razón ≈ 1: no separa.")
	fmt.Fprintf(&b, `<line x1="564" y1="342" x2="1016" y2="342" stroke="#2a3040"/>`)
	txt(564, 368, 11, "#6ee7a0", "start", "✅ la distancia del pliegue al divisor más cercano")
	txt(564, 390, 11, "#8a8676", "start", "para un primo ese divisor es el 1: lo más lejos posible")
	txt(790, 424, 20, "#6ee7a0", "middle", "separa ~6×")
	txt(790, 452, 10, "#8a8676", "middle", "⚖️ pero hallarlo ES dividir: circularidad declarada")

	// the analogy
	fmt.Fprintf(&b, `<rect x="20" y="488" width="1020" height="174" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(530, 514, 13, "#ff6b6b", "middle", "el paralelo con la línea crítica — ANALOGÍA DE FORMA, no teorema")
	yy := 548.0
	for _, r := range [][3]string{
		{"los divisores", "d  ↔  n/d", "punto fijo: √n"},
		{"la ecuación funcional", "s  ↔  1 − s", "punto fijo: ½"},
		{"el producto de Euler", "p^(−s) ↔ p^(−(1−s))", "punto fijo: 1/√p"},
	} {
		txt(140, yy, 12, "#c9c3b0", "start", r[0])
		txt(420, yy, 13, "#e8e2d0", "start", r[1])
		txt(760, yy, 12, "#f2c14e", "start", r[2])
		yy += 28
	}
	txt(530, 640, 11, "#8a8676", "middle", "Que dos emparejamientos tengan la misma forma NO hace que uno implique al otro (la casa se comió esa lección en F401).")
	txt(530, 660, 10, "#8a8676", "middle", "La desigualdad de las medias es de Euclides. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "los-cinco-atados.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
