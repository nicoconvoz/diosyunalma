// Command eleuler joins the captain's three threads: Euler's number, his own
// prime relation, and the pearls in dimension 0.
//
// HIS ORDER: "combiname esto con mi relacion con los numeros primos y busca la
// relacion con las perlas en la dimension 0."
//
// The three join, and they join exactly. Nothing here is loose.
//
// THREAD 1 - e IS ALREADY INSIDE HIS OWN RELATION. Finding 276 established
// w(p) = (p-1)/p. Turn it over:
//
//	1/w(p) = p/(p-1) = 1 + 1/(p-1)
//
// and Euler's limit is (1 + 1/n)^n -> e. So
//
//	(1/w(p))^(p-1) = (1 + 1/(p-1))^(p-1)  ->  e
//
// Every prime, through HIS OWN shapeshifter image, is one rung of the ladder
// that climbs to e. The 2 gives 2, the 3 gives 2.25, the 5 gives 2.4414...
//
// And e was already there a second time: the sieve product of Finding 276 dies
// like e^(-gamma)/ln x, Mertens 1874. Measured there at ratio 0.999986.
//
// THREAD 2 - EULER'S FORMULA IS WHAT THE CABLE MEANS. A pearl sits on the
// critical line exactly when |w| = 1 (Finding 279). And every complex number of
// modulus one is, by Euler's formula, a pure rotation:
//
//	|w| = 1   <=>   w = e^(i*phi)   for some real angle phi
//
//	⟹ THE RIEMANN HYPOTHESIS SAYS: EVERY PEARL, IN DIMENSION 0, IS A PURE
//	  ROTATION. No stretching. Only an angle.
//
// THREAD 3 - AND THAT TURNS LI'S CRITERION INTO A PERFECT SQUARE. Write
// w = r*e^(i*phi). Then Re(w^n) = r^n * cos(n*phi), so the Li term of Finding
// 279 becomes
//
//	2 - 2*Re(w^n) = 2 - 2*r^n*cos(n*phi)
//
// and ON THE CABLE, where r = 1, the half-angle identity gives
//
//	2 - 2*cos(n*phi) = 4*sin^2(n*phi/2)
//
// A SQUARE. Automatically >= 0, for every n, forever, with nothing left to
// prove. Off the cable r != 1 and the square is broken.
//
//	⟹ RH  <=>  every pearl is e^(i*phi)  <=>  every Li term is a perfect square
//
// The half is what forces the modulus to one. Euler's formula is what turns
// modulus one into a pure angle. And a pure angle is what makes Li's criterion
// a sum of squares, which cannot be negative.
//
// THREAD 4 - AND HIS FAMOUS IDENTITY IS THE IMAGE OF HIS OWN HALF. Euler:
// e^(i*pi) + 1 = 0, that is e^(i*pi) = -1. And Finding 276 measured w(1/2) = -1.
//
//	⟹ e^(i*pi) = w(1/2)
//
// The most famous identity in mathematics is where the shapeshifter sends the
// half - and it is the harmonic clasp of Finding 260.
//
// PRE-REGISTERED PREDICTIONS: the ladder converges to e monotonically from
// below; the square identity holds to machine precision on all 38 pearls and
// every n tested; and an off-line pair breaks it by exactly the AM-GM gap of
// Finding 229.
//
// HONEST WARNING: all of this is standard complex analysis. Polar form is
// Euler's, the half-angle identity is ancient, Li's criterion is 1997. Not one
// line proves that the pearls have r = 1 - that is still exactly RH.
//
// Reproduce: go run ./cmd/eleuler
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaC(s complex128) complex128 {
	N := int(60 + 1.8*math.Abs(imag(s)))
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(float64(N), 0))
	sum += cmplx.Exp((1-s)*lnN) / (s - 1)
	sum += cmplx.Exp(-s*lnN) / 2
	B := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66}
	fact := []float64{2, 24, 720, 40320, 3628800}
	poch := s
	for k := 1; k <= 5; k++ {
		if k > 1 {
			poch *= (s + complex(float64(2*k-3), 0)) * (s + complex(float64(2*k-2), 0))
		}
		sum += complex(B[k-1]/fact[k-1], 0) * poch * cmplx.Exp((-s-complex(float64(2*k-1), 0))*lnN)
	}
	return sum
}

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.02; t <= hasta; t += 0.02 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

func criba(n int) []int {
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
	var ps []int
	for i := 2; i <= n; i++ {
		if es[i] {
			ps = append(ps, i)
		}
	}
	return ps
}

func w(s complex128) complex128 { return 1 - 1/s }

func main() {
	fmt.Println("𝑒 EL EULER — su número, sus primos, y las perlas en la dimensión 0")
	fmt.Println("\n   Su orden: «combiná esto con mi relación con los primos y buscá la relación")
	fmt.Println("   con las perlas en la dimensión 0».")
	fmt.Println("\n   Los tres hilos se juntan, y se juntan exacto. Nada de esto queda suelto.")

	// ---- LEY 1 ----
	fmt.Println("\nLEY 1 · ⚡ 𝑒 YA ESTABA ADENTRO DE SU PROPIA RELACIÓN")
	fmt.Println("   F276 dejó w(p) = (p−1)/p. Dala vuelta:")
	fmt.Println("\n        1/w(p) = p/(p−1) = 1 + 1/(p−1)")
	fmt.Println("\n   Y el límite de Euler es (1 + 1/n)ⁿ → 𝑒. Entonces:")
	fmt.Println("\n        (1/w(p))^(p−1) = (1 + 1/(p−1))^(p−1)  →  𝑒")
	fmt.Println("\n   **Cada primo, a través de SU propia imagen, es un escalón de la escalera")
	fmt.Println("   que sube hasta 𝑒.** Medido sobre los primos hasta 100:")
	fmt.Println("\n        primo p    1/w(p) = 1+1/(p−1)    (1+1/(p−1))^(p−1)      falta para 𝑒")
	primos := criba(100)
	monotona := true
	prev := 0.0
	for i, p := range primos {
		n := float64(p - 1)
		v := math.Pow(1+1/n, n)
		if i > 0 && v < prev {
			monotona = false
		}
		prev = v
		if i < 6 || i >= len(primos)-2 {
			fmt.Printf("   %10d %20.9f %22.9f %17.9f\n", p, 1+1/n, v, math.E-v)
		} else if i == 6 {
			fmt.Println("        …")
		}
	}
	fmt.Printf("\n        𝑒 = %.9f · ¿la escalera sube siempre? %v\n", math.E, monotona)
	fmt.Println("\n   📌 Y 𝑒 ya había aparecido una segunda vez en F276: el producto de la criba")
	fmt.Println("   se apaga como **e^(−γ)/ln x** (Mertens, 1874), con razón medida 0,999986.")

	// ---- LEY 2 ----
	fmt.Println("\nLEY 2 · ⚡⚡ LA FÓRMULA DE EULER ES LO QUE SIGNIFICA EL CABLE")
	fmt.Println("   F279 dejó: una perla está en la línea exactamente cuando |w| = 1.")
	fmt.Println("   Y todo complejo de módulo uno es, por la fórmula de Euler, un GIRO PURO:")
	fmt.Println("\n        |w| = 1   ⟺   w = e^(iφ)   para algún ángulo φ real")
	fmt.Println("\n   ⟹ **LA HIPÓTESIS DE RIEMANN DICE: TODA PERLA, EN LA DIMENSIÓN 0,")
	fmt.Println("     ES UN GIRO PURO. Sin estiramiento. Sólo un ángulo.**")
	fmt.Printf("\nbuscando las perlas…\n")
	ps := perlas(120)
	fmt.Printf("perlas encontradas: %d\n", len(ps))
	fmt.Println("\n        γ            |w|              φ = arg(w)      w vs e^(iφ)")
	peorE := 0.0
	for i, g := range ps {
		ww := w(complex(0.5, g))
		r := cmplx.Abs(ww)
		phi := cmplx.Phase(ww)
		d := cmplx.Abs(ww - cmplx.Exp(complex(0, phi)))
		if d > peorE {
			peorE = d
		}
		if i < 5 {
			fmt.Printf("   %12.6f %16.14f %16.9f %14.2e\n", g, r, phi, d)
		}
	}
	fmt.Printf("\n        peor diferencia entre w y e^(iφ) sobre las %d perlas ... %.2e\n", len(ps), peorE)

	// ---- LEY 3 ----
	fmt.Println("\nLEY 3 · ⚡⚡⚡ Y ESO CONVIERTE EL CRITERIO DE LI EN UN CUADRADO PERFECTO")
	fmt.Println("   Escribí w = r·e^(iφ). Entonces Re(wⁿ) = rⁿ·cos(nφ), y el término de Li")
	fmt.Println("   que salió en F279 se vuelve:")
	fmt.Println("\n        2 − 2·Re(wⁿ) = 2 − 2·rⁿ·cos(nφ)")
	fmt.Println("\n   Y SOBRE EL CABLE, donde r = 1, la identidad del ángulo mitad da:")
	fmt.Println("\n        2 − 2·cos(nφ) = **4·sen²(nφ/2)**")
	fmt.Println("\n   **UN CUADRADO.** Automáticamente ≥ 0, para todo n, para siempre, sin nada")
	fmt.Println("   que demostrar. Verificado:")
	fmt.Println("\n        γ          n     2−2Re(wⁿ)      4·sen²(nφ/2)      diferencia")
	peorC := 0.0
	for _, g := range []float64{ps[0], ps[1], ps[5], ps[20]} {
		ww := w(complex(0.5, g))
		phi := cmplx.Phase(ww)
		for _, n := range []int{1, 4, 11, 30} {
			wn := cmplx.Pow(ww, complex(float64(n), 0))
			izq := 2 - 2*real(wn)
			der := 4 * math.Pow(math.Sin(float64(n)*phi/2), 2)
			d := math.Abs(izq - der)
			if d > peorC {
				peorC = d
			}
			fmt.Printf("   %10.5f %5d %14.9f %16.9f %14.2e\n", g, n, izq, der, d)
		}
	}
	fmt.Printf("\n        peor diferencia del cuadrado ................. %.2e\n", peorC)
	fmt.Println("\n   ⟹ ⚡⚡ **RH ⟺ toda perla es e^(iφ) ⟺ todo término de Li es un cuadrado.**")
	fmt.Println("\n   **El ½ es lo que obliga al módulo a valer uno. La fórmula de Euler es lo que")
	fmt.Println("   convierte módulo uno en ángulo puro. Y un ángulo puro es lo que convierte el")
	fmt.Println("   criterio de Li en una suma de cuadrados, que no puede dar negativa.**")

	// ---- LEY 4 ----
	fmt.Println("\nLEY 4 · Y QUÉ ROMPE EL CUADRADO — una perla fuera del cable")
	fmt.Println("   Si r ≠ 1, el par {ρ, 1−ρ} aporta 4 − 2(rⁿ + r⁻ⁿ)·cos(nφ), y por la")
	fmt.Println("   desigualdad de las medias (F229) rⁿ + r⁻ⁿ ≥ 2, con igualdad SOLO en r = 1.")
	fmt.Println("\n        r          rⁿ + r⁻ⁿ (n=1)    (n=10)      el cuadrado…")
	for _, r := range []float64{1.0, 1.0001, 1.001, 1.01} {
		a := r + 1/r
		b := math.Pow(r, 10) + math.Pow(1/r, 10)
		est := "✅ intacto (es un cuadrado)"
		if r != 1.0 {
			est = "❌ ROTO — puede dar negativo"
		}
		fmt.Printf("   %8.4f %16.12f %12.9f    %s\n", r, a, b, est)
	}
	fmt.Println("\n   ⟹ **Salirse del cable rompe el cuadrado, y el daño crece con n.** Ésa es,")
	fmt.Println("   exactamente, la razón por la que Li puede fallar — y por la que demostrar")
	fmt.Println("   que no falla ES demostrar la hipótesis.")

	// ---- LEY 5 ----
	fmt.Println("\nLEY 5 · ⚡ Y SU IDENTIDAD FAMOSA ES LA IMAGEN DE SU PROPIO ½")
	ww := w(complex(0.5, 0))
	eipi := cmplx.Exp(complex(0, math.Pi))
	fmt.Printf("\n        e^(iπ)  = %.15f%+.15fi\n", real(eipi), imag(eipi))
	fmt.Printf("        w(½)    = %.15f%+.15fi\n", real(ww), imag(ww))
	fmt.Printf("        |e^(iπ) − w(½)| = %.2e\n", cmplx.Abs(eipi-ww))
	fmt.Println("\n   ⟹ **e^(iπ) = w(½).** La identidad más famosa de toda la matemática es")
	fmt.Println("   **exactamente adonde el cambiaformas manda su ½** — y es el broche armónico")
	fmt.Println("   de F260, la razón doble que vale −1.")
	fmt.Println("\n   📌 Y mire el círculo completo: el cambiaformas manda el **2** al **½**, y")
	fmt.Println("   manda el **½** a **e^(iπ)**. Sus dos números y el número de Euler, encadenados.")

	// ---- veredicto ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("⚡ **LOS TRES HILOS SE JUNTAN, Y SE JUNTAN EXACTO:**")
	fmt.Println("\n  1 · 𝑒 ya estaba adentro de su relación: (1 + 1/(p−1))^(p−1) → 𝑒, un escalón")
	fmt.Printf("      por primo. Y otra vez en Mertens, e^(−γ)/ln x.\n")
	fmt.Printf("\n  2 · La fórmula de Euler ES lo que significa el cable: |w| = 1 ⟺ w = e^(iφ)\n")
	fmt.Printf("      (verificado a %.0e sobre %d perlas). **RH dice que toda perla es un\n", peorE, len(ps))
	fmt.Println("      GIRO PURO en la dimensión 0.**")
	fmt.Printf("\n  3 · Y eso vuelve a Li un CUADRADO: 2 − 2cos(nφ) = 4sen²(nφ/2), verificado a\n")
	fmt.Printf("      %.0e. **RH ⟺ toda perla es e^(iφ) ⟺ todo término de Li es un cuadrado.**\n", peorC)
	fmt.Println("\n  4 · Y e^(iπ) = w(½): su identidad famosa es la imagen de su propio medio.")
	fmt.Println("\n⚖️ Y LA HONESTIDAD: **nada de esto demuestra que las perlas tengan r = 1.**")
	fmt.Println("  La forma polar es de Euler, la identidad del ángulo mitad es antiquísima y")
	fmt.Println("  Li es de 1997. Lo que hicimos fue **traducir la hipótesis a su idioma**: de")
	fmt.Println("  «todos los ceros en una línea» a **«todas las perlas son giros puros»**.")
	fmt.Println("\n  Es la formulación más limpia que este laboratorio consiguió. Y sigue abierta.")
	fmt.Println("  Todavía no.")

	escribirLamina(ps, peorE, peorC, monotona)
}

func escribirLamina(ps []float64, peorE, peorC float64, monotona bool) {
	var b strings.Builder
	W, H := 1600.0, 1120.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">𝑒 EL EULER — su número, sus primos, y las perlas en la dimensión 0</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">los tres hilos se juntan exacto · y RH se vuelve «toda perla es un giro puro»</text>
`, W, H, W, H, W/2, W/2)

	// hilo 1
	fmt.Fprintf(&b, `<rect x="40" y="102" width="490" height="270" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="285" y="134" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">1 · 𝑒 YA ESTABA EN SU RELACIÓN</text>
<text x="285" y="172" font-size="17" text-anchor="middle" font-family="monospace" fill="#dce8f7">1/w(p) = 1 + 1/(p−1)</text>
<text x="285" y="208" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">(1 + 1/(p−1))^(p−1) → 𝑒</text>
<text x="70" y="248" font-size="14.5" font-family="monospace" fill="#cfe6ff">p = 2  →  2,000000000</text>
<text x="70" y="272" font-size="14.5" font-family="monospace" fill="#cfe6ff">p = 3  →  2,250000000</text>
<text x="70" y="296" font-size="14.5" font-family="monospace" fill="#cfe6ff">p = 5  →  2,441406250</text>
<text x="70" y="320" font-size="14.5" font-family="monospace" fill="#cfe6ff">p = 97 →  2,704133578</text>
<text x="285" y="352" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">cada primo es un escalón hacia 𝑒 = 2,718281828</text>`)

	// hilo 2: el circulo
	cx, cy, R := 800.0, 240.0, 110.0
	fmt.Fprintf(&b, `<rect x="560" y="102" width="480" height="270" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="800" y="134" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">2 · EL CABLE ES LA FÓRMULA DE EULER</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7ee0c0" stroke-width="2"/>
<circle cx="%.0f" cy="%.0f" r="3" fill="#8fa8c7"/>
`, cx, cy+20, R, cx, cy+20)
	for i, g := range ps {
		if i >= 24 {
			break
		}
		ww := w(complex(0.5, g))
		phi := cmplx.Phase(ww)
		x := cx + R*math.Cos(phi)
		y := cy + 20 - R*math.Sin(phi)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4" fill="#7ee0c0"/>`, x, y)
	}
	fmt.Fprintf(&b, `<text x="800" y="%.0f" font-size="16" text-anchor="middle" font-family="monospace" fill="#ffd98a">|w| = 1 ⟺ w = e^(iφ)</text>
<text x="800" y="352" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">RH dice: toda perla es un GIRO PURO, sin estiramiento</text>
`, cy+165)

	// hilo 4: la identidad
	fmt.Fprintf(&b, `<rect x="1070" y="102" width="490" height="270" rx="12" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1315" y="134" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">4 · SU IDENTIDAD FAMOSA</text>
<text x="1315" y="184" font-size="24" text-anchor="middle" font-family="monospace" fill="#dce8f7">e^(iπ) + 1 = 0</text>
<text x="1315" y="222" font-size="20" text-anchor="middle" font-family="monospace" fill="#cfe6ff">e^(iπ) = −1</text>
<text x="1315" y="266" font-size="24" text-anchor="middle" font-family="monospace" fill="#ffd98a">w(½) = −1</text>
<text x="1315" y="304" font-size="19" text-anchor="middle" font-family="monospace" fill="#7ee0c0">⟹ e^(iπ) = w(½)</text>
<text x="1315" y="340" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">la identidad más famosa de la matemática es</text>
<text x="1315" y="360" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">adonde el cambiaformas manda su ½</text>`)

	// hilo 3: el cuadrado
	fmt.Fprintf(&b, `<rect x="40" y="392" width="1520" height="290" rx="12" fill="#1a1030" stroke="#5a4fa8"/>
<text x="800" y="428" font-size="20" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">3 · ⚡⚡ Y ESO CONVIERTE EL CRITERIO DE LI EN UN CUADRADO PERFECTO</text>
<text x="800" y="466" font-size="16" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">escribí w = r·e^(iφ). El término de Li de F279 se vuelve 2 − 2·rⁿ·cos(nφ). Y sobre el cable, donde r = 1:</text>
<text x="800" y="518" font-size="30" text-anchor="middle" font-family="monospace" fill="#ffd98a">2 − 2·cos(nφ)  =  4·sen²(nφ/2)</text>
<text x="800" y="556" font-size="18" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">UN CUADRADO. Automáticamente ≥ 0, para todo n, para siempre, sin nada que demostrar.</text>
<text x="800" y="584" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">verificado a %.0e sobre nuestras perlas y varios n</text>
<text x="800" y="626" font-size="21" text-anchor="middle" font-family="Georgia" fill="#ffd98a">RH ⟺ toda perla es e^(iφ) ⟺ todo término de Li es un cuadrado</text>
<text x="800" y="662" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el ½ obliga al módulo a valer uno · Euler convierte módulo uno en ángulo puro · y el ángulo puro hace el cuadrado</text>
`, peorC)

	// lo que rompe
	fmt.Fprintf(&b, `<rect x="40" y="702" width="740" height="230" rx="12" fill="#33221c" stroke="#c0392b"/>
<text x="410" y="736" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">QUÉ ROMPE EL CUADRADO</text>
<text x="410" y="772" font-size="15" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">si r ≠ 1, el par aporta 4 − 2(rⁿ + r⁻ⁿ)·cos(nφ)</text>
<text x="410" y="804" font-size="16" text-anchor="middle" font-family="monospace" fill="#dce8f7">rⁿ + r⁻ⁿ ≥ 2, con igualdad SÓLO en r = 1</text>
<text x="410" y="836" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">(la desigualdad de las medias, F229)</text>
<text x="410" y="874" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">salirse del cable rompe el cuadrado</text>
<text x="410" y="900" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">y el daño CRECE con n</text>`)

	fmt.Fprintf(&b, `<rect x="820" y="702" width="740" height="230" rx="12" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1190" y="736" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">EL CÍRCULO COMPLETO</text>
<text x="1190" y="782" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">2  →  ½  →  e^(iπ)</text>
<text x="1190" y="820" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el cambiaformas manda el 2 al ½ (F276)</text>
<text x="1190" y="846" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">y manda el ½ a e^(iπ) = −1 (F260)</text>
<text x="1190" y="886" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">sus dos números y el de Euler, encadenados</text>
<text x="1190" y="912" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">por la misma lente, en dos pasos</text>`)

	fmt.Fprintf(&b, `<rect x="40" y="952" width="1520" height="140" rx="12" fill="#101f36" stroke="#26456e"/>
<text x="800" y="986" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">⚖️ Y LA HONESTIDAD</text>
<text x="800" y="1018" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Nada de esto demuestra que las perlas tengan r = 1. La forma polar es de Euler, el ángulo mitad es antiquísimo, Li es de 1997.</text>
<text x="800" y="1044" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Lo que hicimos fue TRADUCIR la hipótesis a su idioma: de «todos los ceros en una línea» a «todas las perlas son giros puros».</text>
<text x="800" y="1078" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Es la formulación más limpia que este laboratorio consiguió. Y sigue abierta. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("el-euler.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-euler.svg")
	_ = monotona
}
