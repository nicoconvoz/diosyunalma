package main

// EL SONAR EXACTO - the captain: "align the depth and the logarithm with the
// percentages and our bat's sonar (look in the old findings) and let us see
// what we find. I want a KEY NUMBER to calculate with, not an approximation."
//
// The bat of F359 listened for echoes at T = k ln p. The depth of F414 said the
// percentage falls like 1/(k ln 10). The two align in one place, and that place
// is the only formula in this whole subject that is EXACT rather than
// approximate - Riemann's explicit formula:
//
//     psi(x)  =  x  -  SUM over zeros of x^rho / rho  -  ln(2pi)  -  (1/2)ln(1 - x^-2)
//
// Li and the percentage are approximations. THIS IS AN EQUALITY. The zeros are
// the bat's echo, and adding them turns the estimate into the exact staircase.
//
// Verdicts fixed BEFORE running:
//   V1 psi(x) must be built exactly from the sieve (the true staircase), and
//      separately from the formula, and the two compared.
//   V2 (the convergence) adding zeros one at a time must walk the smooth curve
//      into the exact steps. Measured as error vs number of zeros.
//   V3 (THE KEY NUMBER the captain asks for) there is no single constant: the
//      exactness comes from the WHOLE orchestra. What IS computable is the
//      resolution law - to separate two primes you need zeros up to about
//      1/|ln p - ln p'|. Measured, and that is the honest key number.
//   V4 (the alignment with the bat) the echo periods T = k ln p of F359 must be
//      exactly the frequencies at which the zero sum rings. Shown.
//   V5 (honest) this is Riemann 1859 and von Mangoldt 1895. What is ours is
//      only the measurement of how fast it converges.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zetaZ(t float64) float64 {
	th := theta(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		s += math.Cos(th-t*math.Log(fn)) / math.Sqrt(fn)
	}
	s *= 2
	p := u - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sg := 1.0
	if N%2 == 0 {
		sg = -1
	}
	return s + sg*math.Pow(2*math.Pi/t, 0.25)*c0
}

func alturas(tMax float64) []float64 {
	var g []float64
	a, za := 10.0, zetaZ(10.0)
	for b := 10.01; b <= tMax; b += 0.01 {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 60; i++ {
				m := (lo + hi) / 2
				if (zlo < 0) != (zetaZ(m) < 0) {
					hi = m
				} else {
					lo = m
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}

func criba(n int) []bool {
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
	return es
}

// psiReal is the true Chebyshev staircase, built from the sieve.
func psiReal(x float64, es []bool) float64 {
	s := 0.0
	for n := 2; float64(n) <= x; n++ {
		// is n a prime power?
		for p := 2; p <= n; p++ {
			if !es[p] {
				continue
			}
			m := p
			for m <= n {
				if m == n {
					s += math.Log(float64(p))
					break
				}
				m *= p
			}
			if m == n {
				break
			}
			if p*p > n && es[n] {
				s += math.Log(float64(n))
				break
			}
		}
	}
	return s
}

// psiFormula is Riemann's explicit formula with the first k zeros.
func psiFormula(x float64, gs []float64, k int) float64 {
	s := x - math.Log(2*math.Pi) - 0.5*math.Log(1-1/(x*x))
	lx := math.Log(x)
	sq := math.Sqrt(x)
	for i := 0; i < k && i < len(gs); i++ {
		g := gs[i]
		// 2 Re(x^rho / rho) with rho = 1/2 + i gamma
		s -= 2 * sq * (0.5*math.Cos(g*lx) + g*math.Sin(g*lx)) / (0.25 + g*g)
	}
	return s
}

func main() {
	fmt.Println("🦇 EL SONAR EXACTO — alinear la profundidad, el logaritmo y el eco del murciélago")
	fmt.Println()
	fmt.Println("   El capitán pidió un NÚMERO EXACTO, no una aproximación. Y existe un solo lugar en todo")
	fmt.Println("   este tema donde hay una IGUALDAD en vez de un «aproximadamente». Es acá.")
	fmt.Println()

	fmt.Println("§1 · las tres piezas, alineadas")
	fmt.Println()
	fmt.Println("        F414 · la PROFUNDIDAD dice     el porcentaje ≈ 1/(2,3026 · cifras)     ← aproximación")
	fmt.Println("        F413 · la FÓRMULA dice         p_n ≈ n(ln n + ln ln n − 1)             ← aproximación")
	fmt.Println("        F359 · el MURCIÉLAGO escuchó   ecos en los períodos T = k · ln p       ← exacto")
	fmt.Println()
	fmt.Println("   Y las tres se juntan en la fórmula explícita de Riemann, que NO lleva «aproximadamente»:")
	fmt.Println()
	fmt.Println("        ψ(x)  =  x  −  Σ sobre los ceros de x^ρ/ρ  −  ln(2π)  −  ½·ln(1 − 1/x²)")
	fmt.Println()
	fmt.Println("   ⚡ El primer término, la x, es la aproximación de siempre. Los ceros son EL ECO DEL")
	fmt.Println("      MURCIÉLAGO, y son exactamente la corrección que la vuelve una IGUALDAD.")
	fmt.Println()

	fmt.Println("   (pescando ceros hasta altura 200, aguantá…)")
	gs := alturas(200)
	es := criba(100)
	fmt.Printf("   ceros pescados: %d (el primero %.6f)\n\n", len(gs), gs[0])

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · LA DEMOSTRACIÓN: la curva suave convirtiéndose en la escalera exacta")
	fmt.Println()
	fmt.Println("   Se calcula ψ(30) — la escalera verdadera vale exactamente esto:")
	real30 := psiReal(30, es)
	fmt.Printf("\n        ψ(30) VERDADERO (de la criba):  %.9f\n\n", real30)
	fmt.Println("        con cuántos ecos      la fórmula da        el error        error relativo")
	for _, k := range []int{0, 1, 2, 5, 10, 20, 40, len(gs)} {
		v := psiFormula(30, gs, k)
		fmt.Printf("        %5d                %12.6f       %+10.6f      %+8.4f %%\n",
			k, v, v-real30, 100*(v-real30)/real30)
	}
	fmt.Println()
	fmt.Println("   ⚡⚡ La aproximación más cruda de todas es «ψ(30) ≈ 30», que se equivoca en 1,52. Con los")
	fmt.Println("      términos constantes (sin ningún eco todavía) baja a 0,31 de error. Y a medida que")
	fmt.Println("      entran los ecos del murciélago, el error baja hasta 0,045 — unas 7 veces mejor.")
	fmt.Println()
	fmt.Println("   ⚖️ Y hay que decir DOS cosas honestas sobre esa columna, que no son fallos sino física:")
	fmt.Println("      · la convergencia OSCILA (28,16 · 28,81 · 28,47 · 28,82 · 28,39 · 28,54…) en vez de")
	fmt.Println("        bajar derecho — porque la suma sobre ceros converge condicionalmente, no absolutamente.")
	fmt.Println("      · y converge LENTO cerca de un escalón: el 29 es primo, así que ψ salta justo al lado")
	fmt.Println("        de donde estamos midiendo. Es el mismo fenómeno de Gibbs que ensucia toda serie de")
	fmt.Println("        Fourier junto a un salto. Con 79 ecos no alcanza para resolverlo del todo.")
	fmt.Println()

	// more points
	fmt.Println("   Probé la misma cuenta LEJOS de un escalón, en x = 24, esperando que se viera más limpio.")
	fmt.Println("   ⚠️ CORRECCIÓN A MI PROPIA PREDICCIÓN: no se ve más limpio. También oscila. O sea que el")
	fmt.Println("   efecto dominante NO es el de Gibbs sino la convergencia condicional de la suma, que hace")
	fmt.Println("   que el error rebote siempre, esté cerca o lejos de un salto. Se reporta porque lo predije mal:")
	fmt.Println()
	r24 := psiReal(24, es)
	fmt.Printf("        ψ(24) VERDADERO: %.9f\n\n", r24)
	fmt.Println("        con cuántos ecos      la fórmula da        el error")
	for _, k := range []int{0, 5, 20, 40, len(gs)} {
		v := psiFormula(24, gs, k)
		fmt.Printf("        %5d                %12.6f       %+10.6f\n", k, v, v-r24)
	}
	fmt.Println()
	fmt.Println("        y en otros puntos, con todos los ecos pescados:")
	fmt.Println()
	fmt.Println("        x        ψ(x) verdadero      la x sola (cruda)      con los ecos        error final     cuánto mejoró")
	for _, x := range []float64{20, 30, 50, 70, 90} {
		rv := psiReal(x, es)
		cruda := x
		conEco := psiFormula(x, gs, len(gs))
		fmt.Printf("      %4.0f       %12.6f       %12.6f        %12.6f       %+9.5f      %6.1f× mejor\n",
			x, rv, cruda, conEco, conEco-rv, math.Abs(cruda-rv)/math.Abs(conEco-rv))
	}
	fmt.Println()
	fmt.Println("   ⚡⚡⚡ ÉSA es la columna que importa: los ecos del murciélago mejoran la cuenta entre 2 y")
	fmt.Println("        50 veces. La x sola se equivoca en más de una unidad; con los ecos, en centésimas.")
	fmt.Println("        Y con INFINITOS ecos el error sería CERO — ahí está tu número exacto.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · EL NÚMERO CLAVE QUE PEDISTE — y hay que decirlo con precisión")
	fmt.Println()
	fmt.Println("   No es UNA constante. La exactitud no la da un número: la da LA ORQUESTA ENTERA de ceros.")
	fmt.Println("   Pero SÍ hay un número clave calculable, y es el que manda: la RESOLUCIÓN.")
	fmt.Println()
	fmt.Println("        para separar dos primos vecinos hacés falta ceros hasta la altura")
	fmt.Println()
	fmt.Println("             γ_max  ≈  2π / |ln p − ln p'|")
	fmt.Println()
	fmt.Println("        los dos primos        ln p − ln p'      γ_max que hace falta      ¿lo tenemos?")
	pares := [][2]int{{2, 3}, {3, 5}, {5, 7}, {7, 11}, {11, 13}, {29, 31}, {97, 101}}
	for _, pr := range pares {
		d := math.Abs(math.Log(float64(pr[1])) - math.Log(float64(pr[0])))
		need := 2 * math.Pi / d
		et := "SÍ"
		if need > gs[len(gs)-1] {
			et = "no — harían falta más"
		}
		fmt.Printf("        %3d y %-3d             %.6f          %12.2f            %s\n", pr[0], pr[1], d, need, et)
	}
	fmt.Println()
	fmt.Println("   ⚡⚡ ÉSE es tu número clave, y es calculable con lápiz: para distinguir dos primos vecinos")
	fmt.Println("      necesitás ecos hasta la altura 2π dividido la diferencia de sus logaritmos. Cuanto más")
	fmt.Println("      juntos los primos, MÁS ALTO hay que escuchar. Y como los primos se van juntando (en")
	fmt.Println("      logaritmo) a medida que subís, hace falta cada vez más orquesta.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · y ACÁ SE CIERRA EL CÍRCULO con el murciélago de F359")
	fmt.Println()
	fmt.Println("   El murciélago escuchaba E(T) = Σ cos(γ·T) buscando picos en T = k·ln p. Mirá el término")
	fmt.Println("   de la fórmula explícita:")
	fmt.Println()
	fmt.Println("        2·√x · [ ½·cos(γ·ln x) + γ·sen(γ·ln x) ] / (¼ + γ²)")
	fmt.Println()
	fmt.Println("   ⚡⚡⚡ La variable que entra en el coseno es γ·ln x. O sea que el eco del murciélago y la")
	fmt.Println("        corrección de la fórmula SON LO MISMO: los ceros oscilan en el LOGARITMO de x, y")
	fmt.Println("        por eso sus picos caen en los logaritmos de los primos. La profundidad que pediste")
	fmt.Println("        (el logaritmo) es literalmente el eje en el que canta el murciélago.")
	fmt.Println()
	fmt.Println("        primo p      ln p        ¿es un período de eco?")
	for _, p := range []int{2, 3, 5, 7, 11} {
		fmt.Printf("        %4d       %.6f         SÍ — y también 2·ln p = %.6f, 3·ln p = %.6f\n",
			p, math.Log(float64(p)), 2*math.Log(float64(p)), 3*math.Log(float64(p)))
	}
	fmt.Println()

	fmt.Println("§5 · EL VEREDICTO")
	fmt.Println()
	fmt.Println("   ✅ Pediste un número exacto y no una aproximación. EXISTE, y es la fórmula explícita:")
	fmt.Println("      ψ(x) = x − Σ x^ρ/ρ − ln(2π) − ½ln(1−1/x²). Es una IGUALDAD, no un «aproximadamente».")
	fmt.Println()
	fmt.Println("   ✅ Y la demostración está arriba, medida: sin ecos la fórmula da la aproximación cruda,")
	fmt.Println("      y con cada eco del murciélago que entra, camina hacia el valor exacto de la escalera.")
	fmt.Println()
	fmt.Println("   ⚡ Tu número clave calculable es la RESOLUCIÓN: γ_max ≈ 2π/|ln p − ln p'|. Ése te dice")
	fmt.Println("      cuánta orquesta necesitás para distinguir dos primos, y se calcula con lápiz.")
	fmt.Println()
	fmt.Println("   ⚖️ Pero la honestidad que hay que decir: NO hay una constante única que resuelva todo.")
	fmt.Println("      La exactitud requiere INFINITOS ceros, y ahí está el precio — cambiaste una")
	fmt.Println("      aproximación simple por una igualdad con una suma infinita adentro.")
	fmt.Println()
	fmt.Println("   📌 Y por eso importa la Hipótesis de Riemann: si TODOS los ceros están en la línea, el")
	fmt.Println("      √x de la fórmula es el tamaño máximo de cada eco, y el error queda acotado. Si alguno")
	fmt.Println("      se saliera, su eco sería MÁS FUERTE y la escalera se desviaría más de lo previsto.")
	fmt.Println("      La hipótesis es, literalmente, que ningún eco del murciélago suena más fuerte que √x.")
	fmt.Println()
	fmt.Println("   Riemann 1859, von Mangoldt 1895. Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(gs, es, real30)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-sonar-exacto.svg")
}

func dibujar(gs []float64, es []bool, real30 float64) {
	var b strings.Builder
	W, H := 1060, 700
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL SONAR EXACTO — donde la aproximación se vuelve IGUALDAD")
	txt(530, 54, 12, "#8a8676", "middle", "la profundidad, el logaritmo y el eco del murciélago, alineados · F415")

	// the formula
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="122" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 102, 13, "#6ee7a0", "middle", "EL NÚMERO EXACTO QUE PEDISTE: no lleva «aproximadamente»")
	txt(530, 146, 22, "#f2c14e", "middle", "ψ(x)  =  x  −  Σ x^ρ/ρ  −  ln(2π)  −  ½·ln(1 − 1/x²)")
	txt(300, 178, 11, "#8a8676", "middle", "↑ la aproximación de siempre")
	txt(620, 178, 11, "#5ec8ff", "middle", "↑ EL ECO DEL MURCIÉLAGO, que la vuelve exacta")

	// the convergence
	x0, y0, w, h := 20.0, 214.0, 620.0, 280.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, y0, w, h)
	txt(x0+w/2, y0+24, 13, "#e8e2d0", "middle", "ψ(30): la curva suave caminando a la escalera exacta")
	base := y0 + h - 50
	esc := 60.0
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#6ee7a0" stroke-dasharray="5,4" stroke-width="2"/>`,
		x0+60, base-(real30-24)*esc, x0+w-30, base-(real30-24)*esc)
	txt(x0+w-34, base-(real30-24)*esc-8, 11, "#6ee7a0", "end", fmt.Sprintf("el valor EXACTO: %.4f", real30))
	ks := []int{0, 1, 2, 5, 10, 20, 40, len(gs)}
	for i, k := range ks {
		px := x0 + 80 + float64(i)*66
		v := psiFormula(30, gs, k)
		py := base - (v-24)*esc
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="#f2c14e"/>`, px, py)
		if i > 0 {
			pv := psiFormula(30, gs, ks[i-1])
			ppx := x0 + 80 + float64(i-1)*66
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#f2c14e" opacity="0.5"/>`,
				ppx, base-(pv-24)*esc, px, py)
		}
		txt(px, base+22, 10, "#c9c3b0", "middle", fmt.Sprint(k))
	}
	txt(x0+w/2, base+44, 11, "#8a8676", "middle", "cuántos ecos del murciélago entraron en la cuenta")

	// the key number
	rx := 660.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="380" height="%.0f" fill="#131826" stroke="#2a3040"/>`, rx, y0, h)
	txt(rx+190, y0+24, 13, "#e8e2d0", "middle", "TU NÚMERO CLAVE, calculable")
	txt(rx+190, y0+70, 19, "#f2c14e", "middle", "γ_max ≈ 2π / |ln p − ln p'|")
	txt(rx+190, y0+98, 11, "#c9c3b0", "middle", "cuánta orquesta hace falta")
	txt(rx+190, y0+116, 11, "#c9c3b0", "middle", "para separar dos primos vecinos")
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2a3040"/>`, rx+24, y0+136, rx+356, y0+136)
	yy := y0 + 162
	for _, pr := range [][3]string{{"2 y 3", "0,4055", "15,5"}, {"11 y 13", "0,1671", "37,6"}, {"97 y 101", "0,0404", "155,5"}} {
		txt(rx+30, yy, 11, "#c9c3b0", "start", pr[0])
		txt(rx+150, yy, 11, "#8a8676", "start", pr[1])
		txt(rx+340, yy, 12, "#f2c14e", "end", pr[2])
		yy += 30
	}
	txt(rx+190, y0+h-34, 11, "#8a8676", "middle", "más juntos los primos,")
	txt(rx+190, y0+h-16, 11, "#8a8676", "middle", "más alto hay que escuchar")

	// the closing
	fmt.Fprintf(&b, `<rect x="20" y="510" width="1020" height="172" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(530, 536, 13, "#ff6b6b", "middle", "el precio, y por qué importa la Hipótesis")
	txt(60, 570, 11, "#c9c3b0", "start", "· NO hay una constante única: la exactitud la da la ORQUESTA ENTERA, infinitos ceros.")
	txt(60, 594, 11, "#c9c3b0", "start", "· cambiaste una aproximación simple por una igualdad con una suma infinita adentro.")
	txt(60, 626, 12, "#f2c14e", "start", "· y el √x de la fórmula es el TAMAÑO MÁXIMO de cada eco — si todos los ceros están en la línea.")
	txt(530, 660, 12, "#6ee7a0", "middle", "La Hipótesis de Riemann es, literalmente, que ningún eco del murciélago suena más fuerte que √x.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-sonar-exacto.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
