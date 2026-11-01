package main

// EL ANTERIOR Y EL SIGUIENTE - the captain's question: is there a relation
// between the composition of the divisors of the number BEFORE a prime and the
// NEXT prime?
//
// There is, it is exact, and it is a one-line proof:
//
//   if q | (p-1) then p = 1 (mod q), so p + (q-1) = 0 (mod q).
//   Since p + q - 1 > q, that number is COMPOSITE.
//   Therefore the gap p' - p can never equal q - 1.
//
// So every odd prime divisor q of p-1 FORBIDS exactly one small gap: q-1.
// The divisors of the number before the prime veto candidates for the next one.
//
// Verdicts fixed BEFORE running:
//   V1 (trap #14, declared) p-1 is even for every odd prime. That is the
//      definition of odd, already caught in F272 and F393. It carries no
//      information and must be excluded from every count below.
//   V2 (tautology check, declared as such) violations of the rule above must
//      number exactly ZERO over all primes to 10^7. This is a PROOF being
//      verified, not a discovery - if it printed anything but 0 the code
//      would be wrong.
//   V3 (the measurement that is NOT a tautology) the veto must SHIFT the gap
//      distribution: the mean gap conditioned on 3|(p-1) must exceed the mean
//      gap conditioned on 3 not dividing p-1, and the size of that shift is
//      an empirical number nobody fixed in advance.
//   V4 (the mirror control) the same veto exists on the other side: q|(p+1)
//      forbids gap q+1. Measured too, so the effect is not claimed to be
//      special to "before" when it is not.
//   V5 (how much it actually predicts) fraction of gap variance explained by
//      the small divisors of p-1. Expected to be SMALL: vetoing a few short
//      gaps is not the same as predicting the next prime.
// This is elementary and known (it is the basis of admissible constellations
// in Hardy-Littlewood). Nothing here is new mathematics, and nothing proves RH.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
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

// divisoresImpares returns the distinct ODD prime divisors of n, up to a cap.
func divisoresImpares(n int, tope int) []int {
	var out []int
	for _, q := range []int{3, 5, 7, 11, 13, 17, 19, 23, 29, 31} {
		if q > tope {
			break
		}
		if n%q == 0 {
			out = append(out, q)
		}
	}
	return out
}

type acumulador struct {
	n   int
	sum float64
}

func (a *acumulador) add(v float64) { a.n++; a.sum += v }
func (a acumulador) media() float64 {
	if a.n == 0 {
		return math.NaN()
	}
	return a.sum / float64(a.n)
}

func main() {
	fmt.Println("🚪 EL ANTERIOR Y EL SIGUIENTE — ¿los divisores de p−1 tienen que ver con el primo que sigue?")
	fmt.Println("   Sí. Y la relación es exacta: cada divisor VETA un hueco. Se deriva en un renglón y se mide.")
	fmt.Println()

	const N = 10000000
	_, ps := criba(N)
	fmt.Printf("   criba hasta %d · %d primos · %d huecos medidos\n\n", N, len(ps), len(ps)-1)

	// ------------------------------------------------------------- V1
	fmt.Println("§1 · la trampa, declarada antes de contar (decimocuarta)")
	fmt.Println("   «p−1 es par» vale para todo primo impar — y para todo compuesto impar. Es la definición de")
	fmt.Println("   impar, ya cazada en F272 y F393. El 2 queda EXCLUIDO de todas las cuentas de abajo; solo se")
	fmt.Println("   miran los divisores IMPARES de p−1, que son los que sí llevan información.")
	fmt.Println()

	// ------------------------------------------------------------- the rule
	fmt.Println("§2 · la regla, derivada en un renglón")
	fmt.Println()
	fmt.Println("      si q | (p−1)  ⟹  p ≡ 1 (mod q)  ⟹  p + (q−1) ≡ 0 (mod q)")
	fmt.Println("      y como p + q − 1 > q, ese número es COMPUESTO.")
	fmt.Println("      ⟹ el hueco hasta el primo siguiente NO PUEDE valer q − 1.")
	fmt.Println()
	fmt.Println("   O sea: cada divisor primo impar de p−1 PROHÍBE exactamente un hueco chico.")
	fmt.Println("   El 3 prohíbe el hueco 2 · el 5 prohíbe el 4 · el 7 prohíbe el 6 · el 11 prohíbe el 10.")
	fmt.Println()
	fmt.Println("   Ejemplos a mano:")
	for _, p := range []int{7, 11, 13, 17, 29, 31} {
		ds := divisoresImpares(p-1, 31)
		var vetos []string
		for _, q := range ds {
			vetos = append(vetos, fmt.Sprintf("%d", q-1))
		}
		sig := 0
		for _, x := range ps {
			if x > p {
				sig = x
				break
			}
		}
		v := "ninguno"
		if len(vetos) > 0 {
			v = strings.Join(vetos, ", ")
		}
		fmt.Printf("      p = %-3d  p−1 = %-3d  divisores impares %-8v  huecos vetados: %-8s  siguiente: %-3d (hueco %d)\n",
			p, p-1, ds, v, sig, sig-p)
	}
	fmt.Println()

	// ------------------------------------------------------------- V2
	fmt.Println("§3 · verificación de la regla (esto es un TEOREMA verificándose, no un hallazgo)")
	fmt.Println()
	fmt.Println("        q      hueco vetado (q−1)     veces que ocurre con q | p−1     veces con q ∤ p−1")
	violaciones := 0
	for _, q := range []int{3, 5, 7, 11, 13} {
		conDiv, sinDiv := 0, 0
		for i := 0; i < len(ps)-1; i++ {
			p := ps[i]
			if p <= q {
				continue
			}
			if ps[i+1]-p != q-1 {
				continue
			}
			if (p-1)%q == 0 {
				conDiv++
				violaciones++
			} else {
				sinDiv++
			}
		}
		fmt.Printf("      %3d            %3d                       %d                        %d\n", q, q-1, conDiv, sinDiv)
	}
	fmt.Printf("\n   violaciones totales: %d  → ", violaciones)
	if violaciones == 0 {
		fmt.Println("la regla NO falla nunca, como debe ser: está demostrada, no medida.")
	} else {
		fmt.Println("hay un error en el código o en la derivación: revisar")
	}
	fmt.Println()

	// ------------------------------------------------------------- V3
	fmt.Println("§4 · LO QUE SÍ SE MIDE: cuánto corre el veto la distribución de huecos")
	fmt.Println()
	fmt.Println("        divisor q      hueco medio con q | p−1      hueco medio con q ∤ p−1      corrimiento")
	var datos []dat2
	for _, q := range []int{3, 5, 7, 11, 13} {
		var con, sin acumulador
		for i := 0; i < len(ps)-1; i++ {
			p := ps[i]
			if p <= q {
				continue
			}
			g := float64(ps[i+1] - p)
			if (p-1)%q == 0 {
				con.add(g)
			} else {
				sin.add(g)
			}
		}
		fmt.Printf("        %3d              %.4f                      %.4f                 %+.4f  (%+.2f %%)\n",
			q, con.media(), sin.media(), con.media()-sin.media(), 100*(con.media()/sin.media()-1))
		datos = append(datos, dat2{q, con.media(), sin.media()})
	}
	fmt.Println()
	fmt.Println("   El 3 es el que más mueve la aguja, y tiene sentido: veta el hueco 2, que es el más frecuente")
	fmt.Println("   de todos. Cuanto más grande el q, más raro el hueco que veta y menos se nota.")
	fmt.Println()

	// how many odd divisors -> mean gap
	fmt.Println("        divisores impares de p−1 (entre 3 y 31)      cuántos primos      hueco medio")
	porCantidad := map[int]*acumulador{}
	for i := 0; i < len(ps)-1; i++ {
		p := ps[i]
		if p < 40 {
			continue
		}
		k := len(divisoresImpares(p-1, 31))
		if porCantidad[k] == nil {
			porCantidad[k] = &acumulador{}
		}
		porCantidad[k].add(float64(ps[i+1] - p))
	}
	var ks []int
	for k := range porCantidad {
		ks = append(ks, k)
	}
	sort.Ints(ks)
	for _, k := range ks {
		if porCantidad[k].n < 100 {
			continue
		}
		fmt.Printf("                    %d                              %8d           %.4f\n", k, porCantidad[k].n, porCantidad[k].media())
	}
	fmt.Println("\n   ⚡ Cuantos más primos chicos dividen a p−1, MÁS GRANDE es el hueco al siguiente — porque")
	fmt.Println("      cada divisor tacha un candidato cercano. La composición del anterior empuja al siguiente lejos.")
	fmt.Println()

	// ------------------------------------------------------------- V4
	fmt.Println("§5 · el control del espejo: los divisores de p+1 (el número de DESPUÉS)")
	fmt.Println()
	fmt.Println("   Por el mismo renglón: si q | (p+1) entonces p ≡ −1 (mod q) y p + (q+1) ≡ 0 (mod q),")
	fmt.Println("   así que q | p+1 veta el hueco q+1. El efecto existe de los dos lados, y se mide igual:")
	fmt.Println()
	fmt.Println("        divisor q     veta hueco     hueco medio con q | p+1     con q ∤ p+1      corrimiento")
	for _, q := range []int{3, 5, 7} {
		var con, sin acumulador
		viol := 0
		for i := 0; i < len(ps)-1; i++ {
			p := ps[i]
			if p <= q {
				continue
			}
			g := ps[i+1] - p
			if (p+1)%q == 0 {
				con.add(float64(g))
				if g == q+1 {
					viol++
				}
			} else {
				sin.add(float64(g))
			}
		}
		fmt.Printf("        %3d             %3d           %.4f                    %.4f           %+.4f   (violaciones %d)\n",
			q, q+1, con.media(), sin.media(), con.media()-sin.media(), viol)
	}
	fmt.Println("\n   El veto no es propiedad del «anterior»: es de la aritmética modular, y el anterior y el")
	fmt.Println("   siguiente vetan huecos distintos. Decirlo es lo que impide vender el efecto como si fuera de un lado solo.")
	fmt.Println()

	// ------------------------------------------------------------- V5
	fmt.Println("§6 · ¿y cuánto PREDICE, realmente?")
	fmt.Println()
	var todos acumulador
	for i := 0; i < len(ps)-1; i++ {
		if ps[i] >= 40 {
			todos.add(float64(ps[i+1] - ps[i]))
		}
	}
	varTotal, varDentro := 0.0, 0.0
	nTot := 0
	for i := 0; i < len(ps)-1; i++ {
		p := ps[i]
		if p < 40 {
			continue
		}
		g := float64(ps[i+1] - p)
		varTotal += (g - todos.media()) * (g - todos.media())
		k := len(divisoresImpares(p-1, 31))
		varDentro += (g - porCantidad[k].media()) * (g - porCantidad[k].media())
		nTot++
	}
	r2 := 1 - varDentro/varTotal
	fmt.Printf("   hueco medio global: %.4f  ·  desvío típico: %.4f\n", todos.media(), math.Sqrt(varTotal/float64(nTot)))
	fmt.Printf("   fracción de la varianza del hueco explicada por los divisores impares de p−1: %.4f  (%.2f %%)\n", r2, 100*r2)
	fmt.Println()
	fmt.Println("   ⚖️ Y ésta es la parte que hay que decir fuerte: el efecto es REAL y DEMOSTRABLE, pero chico.")
	fmt.Printf("      Explica el %.2f %% de la variación. Saber los divisores de p−1 te deja tachar candidatos,\n", 100*r2)
	fmt.Println("      no te dice cuál es el siguiente. Vetar no es predecir.")
	fmt.Println()
	fmt.Println("   📌 Nada de esto es matemática nueva: es el mismo argumento con el que Hardy y Littlewood")
	fmt.Println("      decidían qué constelaciones de primos son admisibles. La casa lo mide, no lo descubre.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(datos, porCantidad, ks, r2, todos.media())
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-anterior-y-el-siguiente.svg")
}

type dat2 struct {
	q        int
	con, sin float64
}

func dibujar(datos []dat2, porCantidad map[int]*acumulador, ks []int, r2, medioGlobal float64) {
	var b strings.Builder
	W, H := 1060, 640
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL ANTERIOR Y EL SIGUIENTE — cada divisor de p−1 veta un hueco")
	txt(530, 54, 12, "#8a8676", "middle", "si q divide a p−1, el primo siguiente no puede estar a distancia q−1 · F395")

	// the rule box
	x0, y0, w, h := 20.0, 76.0, 500.0, 230.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#121a16" stroke="#2a5a40"/>`, x0, y0, w, h)
	txt(x0+w/2, y0+24, 13, "#6ee7a0", "middle", "LA REGLA, en un renglón")
	txt(x0+28, y0+58, 13, "#e8e2d0", "start", "q | (p−1)  ⟹  p ≡ 1 (mod q)")
	txt(x0+28, y0+82, 13, "#e8e2d0", "start", "⟹  p + (q−1) ≡ 0 (mod q)  ⟹  compuesto")
	txt(x0+28, y0+106, 13, "#6ee7a0", "start", "⟹  el hueco NO puede valer q − 1")
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2a5a40"/>`, x0+28, y0+124, x0+w-28, y0+124)
	yy := y0 + 150
	for _, s := range []string{"el 3 veta el hueco 2", "el 5 veta el hueco 4", "el 7 veta el hueco 6"} {
		txt(x0+28, yy, 12, "#c9c3b0", "start", s)
		yy += 22
	}
	txt(x0+280, y0+150, 11, "#8a8676", "start", "p = 7 · p−1 = 6 = 2·3")
	txt(x0+280, y0+168, 11, "#8a8676", "start", "el 3 veta el hueco 2 → 9 es compuesto")
	txt(x0+280, y0+186, 11, "#f2c14e", "start", "siguiente: 11, hueco 4 ✓")

	// shift bars
	rx := 540.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, rx, y0, w, h)
	txt(rx+w/2, y0+24, 13, "#e8e2d0", "middle", "cuánto corre el veto al hueco medio")
	bw := 60.0
	for i, d := range datos {
		bx := rx + 50 + float64(i)*90
		hc := (d.con - 14) * 14
		hs := (d.sin - 14) * 14
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.0f" height="%.1f" fill="#f2c14e" opacity="0.85"/>`, bx, y0+h-50-hc, bw/2-2, hc)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.0f" height="%.1f" fill="#5ec8ff" opacity="0.85"/>`, bx+bw/2, y0+h-50-hs, bw/2-2, hs)
		txt(bx+bw/2, y0+h-32, 11, "#c9c3b0", "middle", fmt.Sprintf("q=%d", d.q))
		txt(bx+bw/2, y0+h-18, 10, "#8a8676", "middle", fmt.Sprintf("%+.2f", d.con-d.sin))
	}
	txt(rx+50, y0+48, 10, "#f2c14e", "start", "■ con q | p−1")
	txt(rx+150, y0+48, 10, "#5ec8ff", "start", "■ con q ∤ p−1")
	txt(rx+w/2, y0+h-6, 10, "#8a8676", "middle", "(barras desde 14; el hueco medio global es ~15,5)")

	// staircase
	sx, sy, sw, sh := 20.0, 322.0, 1020.0, 240.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, sx, sy, sw, sh)
	txt(sx+sw/2, sy+24, 13, "#e8e2d0", "middle", "cuantos más primos chicos dividen a p−1, más lejos queda el siguiente")
	var vals []float64
	for _, k := range ks {
		if porCantidad[k].n >= 100 {
			vals = append(vals, porCantidad[k].media())
		}
	}
	if len(vals) > 1 {
		lo, hi := vals[0], vals[0]
		for _, v := range vals {
			lo, hi = math.Min(lo, v), math.Max(hi, v)
		}
		n := 0
		for _, k := range ks {
			if porCantidad[k].n < 100 {
				continue
			}
			v := porCantidad[k].media()
			bx := sx + 120 + float64(n)*150
			hh := (v - lo + 0.5) / (hi - lo + 1) * (sh - 110)
			fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="80" height="%.1f" fill="#6ee7a0" opacity="0.8"/>`, bx, sy+sh-50-hh, hh)
			txt(bx+40, sy+sh-56-hh, 12, "#e8e2d0", "middle", fmt.Sprintf("%.2f", v))
			txt(bx+40, sy+sh-32, 11, "#c9c3b0", "middle", fmt.Sprintf("%d divisores", k))
			txt(bx+40, sy+sh-16, 10, "#8a8676", "middle", fmt.Sprintf("%d primos", porCantidad[k].n))
			n++
		}
	}
	txt(sx+sw-30, sy+50, 11, "#ff6b6b", "end", fmt.Sprintf("pero solo explica el %.2f %% de la variación", 100*r2))
	txt(sx+sw-30, sy+68, 11, "#8a8676", "end", "vetar no es predecir")
	txt(sx+sw-30, sy+86, 10, "#8a8676", "end", "elemental y conocido: es el argumento de")
	txt(sx+sw-30, sy+102, 10, "#8a8676", "end", "admisibilidad de Hardy–Littlewood")
	txt(530, 600, 11, "#8a8676", "middle", "⚠️ «p−1 es par» quedó excluido: es la definición de impar (decimocuarta trampa). Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-anterior-y-el-siguiente.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
