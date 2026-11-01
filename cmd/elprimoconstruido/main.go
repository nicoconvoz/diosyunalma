package main

// EL PRIMO CONSTRUIDO - the captain's order: build me a prime out of factors
// that are NOT integers. Make the attempt even if it cannot be done. We need
// to experiment.
//
// It CAN be done, and this program does it - forwards and backwards.
//
//   Forwards: pick two irrational reals FIRST, multiply them, and a prime
//   comes out. (5 + sqrt2)(5 - sqrt2) = 25 - 2 = 23. Neither factor is an
//   integer; the product is prime. The prime was CONSTRUCTED.
//
//   Backwards: take any prime and split it. p = sqrt(p) * sqrt(p) always
//   works (sqrt of a prime is irrational - classical proof), and the
//   structured splittings a^2 - d b^2 = p work in specific rings.
//
// Verdicts fixed BEFORE running:
//   V1 (what is gained) the construction must produce genuine primes from
//      genuine non-integers, verified numerically to machine precision and
//      confirmed prime by sieve. Count how many primes each ring reaches.
//   V2 (the structure, not a coincidence) which primes a ring reaches is
//      decided by a congruence: Z[i] reaches p = 1 (mod 4), Z[sqrt2] reaches
//      p = +-1 (mod 8), Z[sqrt5] reaches p = +-1 (mod 5). Measured, with
//      violations required to be ZERO.
//   V3 (coverage) no single ring reaches every prime - each reaches about
//      half. Measured: how many rings are needed so every prime up to N is
//      reached by at least one?
//   V4 (trap #17, declared) p = sqrt(p) * sqrt(p) reaches EVERY prime and is
//      therefore worthless as structure: it factors every integer equally,
//      prime or not. Reported because the captain asked for the attempt, and
//      flagged so it is not sold as more than it is.
//   V5 (what is lost, stated plainly) in any ring where p splits, p is no
//      longer prime THERE - and it stays prime in Z regardless. Nothing here
//      makes a prime composite among the integers.
// Gauss 1832, quadratic reciprocity. Nothing new, nothing about RH.

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

type anillo struct {
	d      int // the ring is Z[sqrt(d)]
	nombre string
	regla  string
}

// factores finds integers a,b with |a^2 - d b^2| = p, i.e. a factorization of p
// into the two conjugate NON-INTEGER numbers a +- b*sqrt(d).
func factores(p, d int) (int, int, bool) {
	lim := 200
	if d < 0 {
		lim = int(math.Sqrt(float64(p)/float64(-d))) + 2
	}
	for b := 1; b <= lim; b++ {
		v := p + d*b*b
		if v < 0 {
			continue
		}
		a := int(math.Round(math.Sqrt(float64(v))))
		if a*a == v && a > 0 {
			// a^2 - d b^2 = p
			return a, b, true
		}
		v2 := -p + d*b*b
		if v2 > 0 {
			a2 := int(math.Round(math.Sqrt(float64(v2))))
			if a2*a2 == v2 && a2 > 0 {
				// a^2 - d b^2 = -p
				return a2, b, true
			}
		}
	}
	return 0, 0, false
}

func txtRaiz(d int) string {
	if d < 0 {
		if d == -1 {
			return "i"
		}
		return fmt.Sprintf("i√%d", -d)
	}
	return fmt.Sprintf("√%d", d)
}

// valorConj returns the two conjugate factors as real/complex magnitudes for display.
func valorConj(a, b, d int) (string, string, float64) {
	if d < 0 {
		m := float64(a*a - d*b*b)
		return fmt.Sprintf("%d + %d%s", a, b, txtRaiz(d)), fmt.Sprintf("%d − %d%s", a, b, txtRaiz(d)), m
	}
	r := math.Sqrt(float64(d))
	fa, fb := float64(a), float64(b)
	// If a^2 - d b^2 is negative the pair (a+b√d, a-b√d) multiplies to -p.
	// Flip to (b√d+a, b√d-a), whose product is +p and whose factors are still
	// both irrational.
	if a*a-d*b*b < 0 {
		return fmt.Sprintf("%.9f", fb*r+fa), fmt.Sprintf("%.9f", fb*r-fa), (fb*r + fa) * (fb*r - fa)
	}
	return fmt.Sprintf("%.9f", fa+fb*r), fmt.Sprintf("%.9f", fa-fb*r), (fa + fb*r) * (fa - fb*r)
}

func main() {
	fmt.Println("🔨 EL PRIMO CONSTRUIDO — fabricar un primo con factores que NO son enteros")
	fmt.Println("   Orden del capitán: «hacé el intento aunque digas que no se pueda». Se puede. Acá está.")
	fmt.Println()

	es, ps := criba(2000000)

	// ------------------------------------------------------------ FORWARDS
	fmt.Println("§1 · HACIA ADELANTE: elegir primero dos números irracionales y que salga un primo")
	fmt.Println()
	fmt.Println("   Se eligen a y b enteros, se arman los DOS FACTORES IRRACIONALES a ± b√2, se multiplican,")
	fmt.Println("   y recién después se pregunta si lo que salió es primo. El primo no se busca: se FABRICA.")
	fmt.Println()
	fmt.Println("        factor 1            factor 2            producto      ¿primo?")
	hechos := 0
	for b := 1; b <= 3 && hechos < 9; b++ {
		for a := 2; a <= 14 && hechos < 9; a++ {
			n := a*a - 2*b*b
			if n <= 1 || !es[n] {
				continue
			}
			f1 := float64(a) + float64(b)*math.Sqrt2
			f2 := float64(a) - float64(b)*math.Sqrt2
			fmt.Printf("     %16.12f   ×  %16.12f  =  %8.6f      %s   (%d ± %d√2)\n",
				f1, f2, f1*f2, "SÍ", a, b)
			hechos++
		}
	}
	fmt.Println()
	fmt.Println("   ⚡ Ninguno de los dos factores es entero — son irracionales, con infinitos decimales que no se")
	fmt.Println("      repiten. Su producto da un entero exacto, y ese entero es PRIMO. El capitán pidió esto")
	fmt.Println("      y esto es, literal: 7 = 3,414213562373… × 1,585786437627…")
	fmt.Println()

	// ------------------------------------------------------------ BACKWARDS
	fmt.Println("§2 · HACIA ATRÁS: agarrar un primo dado y partirlo en dos irracionales")
	fmt.Println()
	anillos := []anillo{
		{2, "ℤ[√2]", "p ≡ ±1 (mod 8)"},
		{3, "ℤ[√3]", "p ≡ ±1 (mod 12)"},
		{5, "ℤ[√5]", "p ≡ ±1 (mod 5)"},
		{-1, "ℤ[i]", "p ≡ 1 (mod 4)"},
		{-2, "ℤ[√−2]", "p ≡ 1,3 (mod 8)"},
	}
	fmt.Println("        p        anillo      los dos factores                                    producto verificado")
	for _, p := range []int{7, 11, 17, 19, 23, 29, 31, 41} {
		puesto := false
		for _, an := range anillos {
			a, b, ok := factores(p, an.d)
			if !ok {
				continue
			}
			s1, s2, prod := valorConj(a, b, an.d)
			if an.d < 0 {
				fmt.Printf("      %3d      %-8s   (%s)(%s)%s= %.0f\n", p, an.nombre, s1, s2,
					strings.Repeat(" ", max(2, 36-len(s1)-len(s2))), prod)
			} else {
				fmt.Printf("      %3d      %-8s   %s × %s   =  %.9f\n", p, an.nombre, s1, s2, prod)
			}
			puesto = true
			break
		}
		if !puesto {
			fmt.Printf("      %3d      —          ninguno de los anillos probados lo alcanza\n", p)
		}
	}
	fmt.Println()

	// ------------------------------------------------------------ V2
	fmt.Println("§3 · no es capricho: qué primos alcanza cada anillo lo decide una congruencia")
	fmt.Println()
	fmt.Println("        anillo      regla                  alcanzados hasta 10.000     fracción     violaciones")
	const NV = 10000
	for _, an := range anillos {
		al, tot, viol := 0, 0, 0
		for _, p := range ps {
			if p > NV {
				break
			}
			if p == 2 {
				continue
			}
			tot++
			_, _, ok := factores(p, an.d)
			if ok {
				al++
			}
			var pred bool
			switch an.d {
			case 2:
				pred = p%8 == 1 || p%8 == 7
			case 3:
				pred = p%12 == 1 || p%12 == 11
			case 5:
				pred = p%5 == 1 || p%5 == 4
			case -1:
				pred = p%4 == 1
			case -2:
				pred = p%8 == 1 || p%8 == 3
			}
			// A prime dividing 2d RAMIFIES: it is reachable but obeys no
			// splitting congruence. The first version of this rule omitted
			// that case and reported one false violation per real ring.
			if an.d%p == 0 || (2*an.d)%p == 0 {
				pred = true
			}
			if ok != pred {
				viol++
			}
		}
		fmt.Printf("        %-8s    %-18s        %6d / %-6d          %.4f        %d\n",
			an.nombre, an.regla, al, tot, float64(al)/float64(tot), viol)
	}
	fmt.Println()
	fmt.Println("   ⚡ Cada anillo alcanza ~la MITAD de los primos, y cuál alcanza lo dice el resto de una división.")
	fmt.Println("      Cero violaciones: la ley de reciprocidad cuadrática de Gauss, verificada sobre 1.228 primos.")
	fmt.Println()

	// ------------------------------------------------------------ V3
	fmt.Println("§4 · ¿cuántos anillos hacen falta para alcanzarlos a TODOS?")
	fmt.Println()
	fmt.Println("        anillos usados                        primos NO alcanzados hasta 10.000")
	acum := []anillo{}
	for _, an := range anillos {
		acum = append(acum, an)
		falta := 0
		var ejemplo []int
		for _, p := range ps {
			if p > NV {
				break
			}
			if p == 2 {
				continue
			}
			hall := false
			for _, a2 := range acum {
				if _, _, ok := factores(p, a2.d); ok {
					hall = true
					break
				}
			}
			if !hall {
				falta++
				if len(ejemplo) < 5 {
					ejemplo = append(ejemplo, p)
				}
			}
		}
		nombres := []string{}
		for _, a2 := range acum {
			nombres = append(nombres, a2.nombre)
		}
		ej := ""
		if len(ejemplo) > 0 {
			ej = fmt.Sprintf("  (ej: %v)", ejemplo)
		}
		fmt.Printf("        %-42s  %5d%s\n", strings.Join(nombres, " "), falta, ej)
	}
	fmt.Println()
	fmt.Println("   Ningún anillo solo alcanza a todos, pero juntando unos pocos casi no queda ninguno afuera.")
	fmt.Println()

	// ------------------------------------------------------------ V4
	fmt.Println("§5 · la ruta que alcanza a TODOS — y por qué no vale gran cosa (trampa decimoséptima)")
	fmt.Println()
	fmt.Println("        p        √p × √p, los dos irracionales            producto")
	for _, p := range []int{2, 3, 5, 7, 11, 13} {
		r := math.Sqrt(float64(p))
		fmt.Printf("      %3d      %.12f × %.12f      =  %.9f\n", p, r, r, r*r)
	}
	fmt.Println()
	fmt.Println("   ⚠️ Esto funciona para TODO primo — la raíz de un primo es irracional siempre, demostrado desde")
	fmt.Println("      los griegos. Pero por eso mismo NO DICE NADA: también funciona para el 4, el 6 y el 9.")
	fmt.Println("      Factoriza a todos por igual, primos y compuestos. Es la trampa del 0.0e+00 número diecisiete,")
	fmt.Println("      y se reporta porque el capitán pidió el intento completo, con el cartel puesto.")
	fmt.Println()

	// ------------------------------------------------------------ V5
	fmt.Println("§6 · LO QUE SE GANÓ Y LO QUE NO")
	fmt.Println()
	fmt.Println("   SE GANÓ, y es literal lo que se pidió:")
	fmt.Println("     · un primo fabricado eligiendo primero dos irracionales: (5+√2)(5−√2) = 23")
	fmt.Println("     · primos partidos en dos factores no enteros, verificados a doce decimales")
	fmt.Println("     · la ley exacta de qué primo se deja partir en qué anillo, con cero violaciones")
	fmt.Println()
	fmt.Println("   NO SE GANÓ, y hay que decirlo con todas las letras:")
	fmt.Println("     · el primo SIGUE SIENDO PRIMO entre los enteros. Nada de esto lo vuelve compuesto ahí.")
	fmt.Println("     · dentro de cada anillo, el número que se parte DEJA de ser primo — pero los factores")
	fmt.Println("       nuevos son primos ELLOS, así que el problema no desaparece: se muda de casa.")
	fmt.Println("     · ningún anillo alcanza a todos los primos; siempre hay que elegir cuál.")
	fmt.Println("     · nada de esto es nuevo: es Gauss, 1832, y la reciprocidad cuadrática.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(ps, es, anillos)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-primo-construido.svg")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dibujar(ps []int, es []bool, anillos []anillo) {
	var b strings.Builder
	W, H := 1060, 660
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL PRIMO CONSTRUIDO — fabricado con dos números que no son enteros")
	txt(530, 54, 12, "#8a8676", "middle", "se eligen primero los factores irracionales; el primo sale después · F398")

	// the headline construction
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="180" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 106, 13, "#6ee7a0", "middle", "HACIA ADELANTE — el primo no se busca, se FABRICA")
	txt(190, 156, 15, "#f2c14e", "middle", "5 + √2")
	txt(190, 182, 12, "#8a8676", "middle", "6,414213562373…")
	txt(300, 160, 18, "#e8e2d0", "middle", "×")
	txt(410, 156, 15, "#f2c14e", "middle", "5 − √2")
	txt(410, 182, 12, "#8a8676", "middle", "3,585786437627…")
	txt(520, 160, 18, "#e8e2d0", "middle", "=")
	txt(620, 160, 26, "#6ee7a0", "middle", "23")
	txt(620, 186, 11, "#6ee7a0", "middle", "PRIMO")
	txt(760, 140, 11, "#c9c3b0", "start", "Ninguno de los dos factores es entero:")
	txt(760, 158, 11, "#c9c3b0", "start", "los dos tienen infinitos decimales que")
	txt(760, 176, 11, "#c9c3b0", "start", "no se repiten nunca. Su producto es")
	txt(760, 194, 11, "#c9c3b0", "start", "un entero exacto, y es primo.")
	txt(530, 228, 11, "#8a8676", "middle", "y también: 7 = (3+√2)(3−√2) · 17 = (1+4i)(1−4i) · 31 = (7+3√2)(7−3√2)")

	// coverage by ring
	fmt.Fprintf(&b, `<rect x="20" y="272" width="640" height="250" fill="#131826" stroke="#2a3040"/>`)
	txt(340, 300, 13, "#e8e2d0", "middle", "qué primos alcanza cada anillo (y la regla que lo decide)")
	yy := 332.0
	reglas := []struct{ n, r string }{
		{"ℤ[i]", "p ≡ 1 (mod 4)"},
		{"ℤ[√2]", "p ≡ ±1 (mod 8)"},
		{"ℤ[√3]", "p ≡ ±1 (mod 12)"},
		{"ℤ[√5]", "p ≡ ±1 (mod 5)"},
		{"ℤ[√−2]", "p ≡ 1,3 (mod 8)"},
	}
	for _, r := range reglas {
		txt(50, yy, 13, "#f2c14e", "start", r.n)
		txt(150, yy, 12, "#c9c3b0", "start", r.r)
		fmt.Fprintf(&b, `<rect x="330" y="%.0f" width="150" height="12" fill="#2a3040"/>`, yy-10)
		fmt.Fprintf(&b, `<rect x="330" y="%.0f" width="75" height="12" fill="#6ee7a0" opacity="0.8"/>`, yy-10)
		txt(500, yy, 12, "#8a8676", "start", "≈ la mitad")
		yy += 34
	}
	txt(340, 505, 11, "#8a8676", "middle", "cero violaciones en 1.228 primos — reciprocidad cuadrática, Gauss 1832")

	// what was lost
	fmt.Fprintf(&b, `<rect x="680" y="272" width="360" height="250" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(860, 300, 13, "#ff6b6b", "middle", "lo que NO se ganó")
	txt(700, 332, 11, "#c9c3b0", "start", "· el primo SIGUE SIENDO PRIMO entre")
	txt(700, 350, 11, "#c9c3b0", "start", "  los enteros. Nada lo vuelve compuesto.")
	txt(700, 378, 11, "#c9c3b0", "start", "· los factores nuevos son primos ELLOS:")
	txt(700, 396, 11, "#c9c3b0", "start", "  el problema no desaparece, se muda.")
	txt(700, 424, 11, "#c9c3b0", "start", "· ningún anillo alcanza a todos.")
	txt(700, 452, 11, "#c9c3b0", "start", "· nada de esto es nuevo: Gauss, 1832.")
	txt(700, 490, 11, "#8a8676", "start", "⚠️ trampa 17: √p × √p alcanza a todos")
	txt(700, 506, 11, "#8a8676", "start", "   pero factoriza al 4 y al 6 igual.")

	// the primes painted by reachability
	fmt.Fprintf(&b, `<rect x="20" y="538" width="1020" height="104" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 562, 12, "#e8e2d0", "middle", "los primeros primos, y el primer anillo que los parte")
	n := 0
	for _, p := range ps {
		if p > 130 || n > 28 {
			break
		}
		if p == 2 {
			n++
			continue
		}
		px := 60.0 + float64(n)*33
		col := "#3a4050"
		for _, an := range anillos {
			if _, _, ok := factores(p, an.d); ok {
				switch an.d {
				case -1:
					col = "#5ec8ff"
				case 2:
					col = "#f2c14e"
				case 3:
					col = "#6ee7a0"
				case 5:
					col = "#ff9f6b"
				default:
					col = "#c58bff"
				}
				break
			}
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="596" r="11" fill="%s" opacity="0.9"/>`, px, col)
		txt(px, 600, 10, "#0d1017", "middle", fmt.Sprintf("%d", p))
		n++
	}
	txt(60, 630, 10, "#8a8676", "start", "● celeste ℤ[i]  ● amarillo ℤ[√2]  ● verde ℤ[√3]  ● naranja ℤ[√5]  ● violeta ℤ[√−2]")
	txt(1020, 630, 10, "#8a8676", "end", "Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-primo-construido.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
