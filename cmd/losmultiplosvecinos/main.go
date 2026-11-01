package main

// LOS MULTIPLOS VECINOS - the captain looked at the pairs and saw what I had
// missed. He wrote:
//
//     4 y 2   son 2x2 y 2x1
//     9 y 6   son 3x3 y 3x2
//     25 y 20 son 5x5 y 5x4
//     49 y 56 son 7x7 y 7x8
//
// Both numbers in every pair are MULTIPLES OF p, and their cofactors are
// CONSECUTIVE. That is not a restatement of n^2 - n(n-1) = n: it is the reason
// the identity works, and it is a better reading than the one this house
// wrote yesterday.
//
//     p*k - p*(k-1) = p*(k - (k-1)) = p*1 = p
//
// Verdicts fixed BEFORE running:
//   V1 every pair the captain listed must decompose into p times consecutive
//      cofactors, verified for all four.
//   V2 the structural reason must be stated as what it is: subtracting two
//      CONSECUTIVE MULTIPLES of p, which differ by p by the definition of
//      multiple. That is why the trap of F403 fires - it is not about primes.
//   V3 (the new measurement, and the decisive one) if the reading is right,
//      then choosing OTHER consecutive multiples must give OTHER valid
//      factorizations of the same number. Measured: how many forms does each
//      number have? The prediction fixed here is INFINITELY MANY, and the
//      program must exhibit several per number.
//   V4 (the consequence) a pattern that offers infinitely many forms for every
//      integer cannot single out primes. Measured against composites.
// The captain's reading is credited as his. Nothing here is new mathematics.

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

func main() {
	fmt.Println("🔗 LOS MÚLTIPLOS VECINOS — lo que vio el capitán y la casa no había visto")
	fmt.Println()

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · su observación, verificada")
	fmt.Println()
	fmt.Println("        par           el de arriba        el de abajo        cofactores      ¿vecinos?")
	pares := []struct{ p, a, b int }{{2, 2, 1}, {3, 3, 2}, {5, 5, 4}, {7, 7, 8}}
	for _, x := range pares {
		arriba, abajo := x.p*x.a, x.p*x.b
		vec := "SÍ"
		if int(math.Abs(float64(x.a-x.b))) != 1 {
			vec = "no"
		}
		fmt.Printf("      %3d y %3d       %d × %d = %-4d       %d × %d = %-4d       %d y %d          %s\n",
			arriba, abajo, x.p, x.a, arriba, x.p, x.b, abajo, x.a, x.b, vec)
	}
	fmt.Println()
	fmt.Println("   ⚡ Los DOS números de cada par son múltiplos del mismo primo, y sus cofactores son VECINOS:")
	fmt.Println("      2 y 1 · 3 y 2 · 5 y 4 · 7 y 8. Siempre pegados, siempre a distancia uno.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · por qué funciona — y esta lectura es mejor que la de ayer")
	fmt.Println()
	fmt.Println("   Ayer la casa lo escribió como identidad algebraica: n² − n(n−1) = n. Es cierto, pero no")
	fmt.Println("   explica NADA. La lectura del capitán sí explica, porque saca el primo afuera:")
	fmt.Println()
	fmt.Println("        p·k − p·(k−1)  =  p · [k − (k−1)]  =  p · 1  =  p")
	fmt.Println()
	fmt.Println("   O sea: restás DOS MÚLTIPLOS VECINOS del mismo número. Y dos múltiplos vecinos de p se")
	fmt.Println("   llevan exactamente p de diferencia — porque eso ES lo que significa «múltiplo».")
	fmt.Println()
	fmt.Println("        p       múltiplos vecinos        su diferencia")
	for _, p := range []int{2, 3, 5, 7, 11, 100} {
		fmt.Printf("      %4d      %5d y %-5d            %d\n", p, p*7, p*8, p*8-p*7)
	}
	fmt.Println()
	fmt.Println("   ⚡⚡ Y ACÁ SE ENTIENDE LA TRAMPA DE AYER, de una vez: el 100 % contra el 100 % no era un")
	fmt.Println("      misterio. La identidad vale para todo número porque es LA DEFINICIÓN DE MÚLTIPLO,")
	fmt.Println("      y el 4, el 6 y el 100 tienen múltiplos vecinos igual que cualquier primo.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · LA CONSECUENCIA QUE NO HABÍAMOS MEDIDO: no hay UNA forma, hay INFINITAS")
	fmt.Println()
	fmt.Println("   Si la clave son dos múltiplos vecinos, entonces se pueden elegir OTROS y tiene que")
	fmt.Println("   funcionar igual. Basta con que el de arriba sea un cuadrado perfecto: p·m = (p·k)².")
	fmt.Println()
	fmt.Println("        número      forma                                   la cuenta                 producto")
	es, _ := criba(200000)
	for _, p := range []int{2, 3, 5, 7} {
		for k := 1; k <= 4; k++ {
			a := p * k
			arriba := a * a
			abajo := arriba - p
			r := math.Sqrt(float64(abajo))
			fmt.Printf("          %d        (%3d + √%-6d)(%3d − √%-6d)      %6d − %-6d = %-3d      %.9f\n",
				p, a, abajo, a, abajo, arriba, abajo, p, (float64(a)+r)*(float64(a)-r))
		}
		fmt.Println()
	}
	fmt.Println("   ⚡⚡ CUATRO formas distintas para cada uno, y se pueden seguir sacando para siempre: una")
	fmt.Println("      por cada k. El 2 no tiene «su» factorización — tiene infinitas, todas exactas.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · y por eso no puede distinguir primos: los compuestos también tienen infinitas")
	fmt.Println()
	fmt.Println("        número    ¿primo?     tres de sus infinitas formas")
	for _, n := range []int{2, 7, 4, 6, 100} {
		et := "no"
		if es[n] {
			et = "SÍ"
		}
		var fs []string
		for k := 1; k <= 3; k++ {
			a := n * k
			fs = append(fs, fmt.Sprintf("(%d±√%d)", a, a*a-n))
		}
		fmt.Printf("         %4d       %-3s        %s\n", n, et, strings.Join(fs, "  "))
	}
	fmt.Println()
	fmt.Println("   ⚖️ Un patrón que le da INFINITAS formas a cada número entero no puede señalar a ninguno.")
	fmt.Println("      Es la confirmación estructural de la trampa vigésima, y ahora con el porqué adentro.")
	fmt.Println()

	fmt.Println("§5 · EL VEREDICTO")
	fmt.Println()
	fmt.Println("   ✅ El capitán vio la estructura que la casa había pasado por alto: los dos números del par")
	fmt.Println("      son múltiplos del mismo primo con cofactores vecinos. Su lectura EXPLICA; la nuestra de")
	fmt.Println("      ayer solo verificaba.")
	fmt.Println("   ✅ Y esa lectura reveló algo que no habíamos medido: cada número tiene INFINITAS de estas")
	fmt.Println("      factorizaciones, no una — cuatro exhibidas por cabeza, y hay tantas como se quieran.")
	fmt.Println("   ⚖️ Y cierra la trampa de ayer con su razón adentro: la identidad vale siempre porque es la")
	fmt.Println("      definición de múltiplo, no una propiedad de los primos.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(pares)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/los-multiplos-vecinos.svg")
}

func dibujar(pares []struct{ p, a, b int }) {
	var b strings.Builder
	W, H := 1060, 660
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "LOS MÚLTIPLOS VECINOS — lo que vio el capitán")
	txt(530, 54, 12, "#8a8676", "middle", "los dos números del par son múltiplos del mismo primo, con cofactores pegados · F404")

	// his observation
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="186" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 102, 13, "#6ee7a0", "middle", "SU OBSERVACIÓN")
	for i, x := range pares {
		px := 150.0 + float64(i)*245
		txt(px, 140, 15, "#f2c14e", "middle", fmt.Sprintf("%d  y  %d", x.p*x.a, x.p*x.b))
		txt(px, 168, 13, "#e8e2d0", "middle", fmt.Sprintf("%d×%d    y    %d×%d", x.p, x.a, x.p, x.b))
		txt(px, 196, 11, "#8a8676", "middle", fmt.Sprintf("cofactores %d y %d", x.a, x.b))
		txt(px, 216, 12, "#6ee7a0", "middle", "vecinos ✓")
	}
	txt(530, 246, 12, "#c9c3b0", "middle", "los dos son múltiplos del mismo primo, y sus cofactores están pegados — siempre")

	// the reason
	fmt.Fprintf(&b, `<rect x="20" y="278" width="500" height="180" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 304, 13, "#e8e2d0", "middle", "por qué funciona: sacar el primo afuera")
	txt(270, 348, 17, "#6ee7a0", "middle", "p·k − p·(k−1)")
	txt(270, 378, 15, "#e8e2d0", "middle", "= p · [k − (k−1)]")
	txt(270, 408, 17, "#f2c14e", "middle", "= p · 1 = p")
	txt(270, 440, 11, "#8a8676", "middle", "dos múltiplos vecinos se llevan p de diferencia")

	fmt.Fprintf(&b, `<rect x="540" y="278" width="500" height="180" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(790, 304, 13, "#ff6b6b", "middle", "y por eso la trampa de ayer")
	txt(790, 338, 12, "#c9c3b0", "middle", "el 100 % contra el 100 % no era misterio:")
	txt(790, 364, 12, "#c9c3b0", "middle", "la identidad es LA DEFINICIÓN DE MÚLTIPLO")
	txt(790, 398, 12, "#c9c3b0", "middle", "y el 4, el 6 y el 100 tienen múltiplos")
	txt(790, 420, 12, "#c9c3b0", "middle", "vecinos igual que cualquier primo")
	txt(790, 446, 11, "#8a8676", "middle", "ahora la trampa tiene su porqué adentro")

	// infinitely many
	fmt.Fprintf(&b, `<rect x="20" y="474" width="1020" height="166" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 500, 13, "#e8e2d0", "middle", "LO QUE NO HABÍAMOS MEDIDO: no hay UNA forma, hay INFINITAS")
	txt(530, 528, 12, "#8a8676", "middle", "si la clave son dos múltiplos vecinos, se pueden elegir otros — y funciona igual")
	yy := 560.0
	for _, r := range [][2]string{
		{"2", "(2 ± √2)   ·   (4 ± √14)   ·   (6 ± √34)   ·   (8 ± √62)   ·   …"},
		{"7", "(7 ± √42)  ·   (14 ± √189) ·   (21 ± √434) ·   (28 ± √777) ·   …"},
	} {
		txt(120, yy, 14, "#f2c14e", "start", "el "+r[0])
		txt(220, yy, 13, "#c9c3b0", "start", r[1])
		yy += 32
	}
	txt(530, 626, 11, "#8a8676", "middle", "un patrón que le da infinitas formas a cada entero no puede señalar a ninguno. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "los-multiplos-vecinos.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
