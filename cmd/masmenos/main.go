// Command masmenos answers the captain's chain: absolute value is distance
// without direction, direction is plus-or-minus distance, x = ±|x| - and then
// his last line, ±|1| = w.
//
// The chain is exact on the number line, and it BREAKS when it crosses into the
// plane. That break is not a detail: it is the entire difficulty of the problem.
//
// ON THE LINE, TWO DIRECTIONS
//
//	|x| = the distance from zero, with the direction stripped off
//	x   = ±|x|, the distance with the direction handed back
//
// So |x| = 1 leaves exactly two possibilities, +1 and −1. Two points.
//
// IN THE PLANE, INFINITE DIRECTIONS
//
//	|w| = 1  does NOT give  w = ±1.  It gives  w = e^{i phi}.
//
// A whole circle. The ± is the two-direction case of an infinite-direction
// case, and ±1 are only the two REAL points of that circle. Measured here: the
// 649 pearls all have |w| = 1 and all sit at DIFFERENT angles. How many equal
// ±1? Zero.
//
// AND ±1 ARE EXACTLY THE TWO POINTS OF THE SKIN THAT ARE NOT PEARLS
//
//	w = +1  is the CLASP, the image of rho = infinity. No finite pearl
//	        reaches it (asking for it gives -1 = 0).
//	w = -1  is the image of rho = 1/2, the point where the critical line
//	        crosses the real axis - and zeta does not vanish there.
//
// So the captain's ±|1| picks out precisely the two points of the skin where
// there is no pearl. Every actual pearl sits at an angle strictly between them.
//
// # WHERE HIS ± DOES SURVIVE, AND EXACTLY
//
// On the skin, |w| = 1 means w * conj(w) = 1, i.e. conj(w) = 1/w. The mirror
// of the plane is conjugation, and on the skin conjugation IS inversion. So the
// captain's "two directions from zero" becomes, in the plane, "a pearl and its
// conjugate" - and every pearl really does come in that pair. His ± is exact;
// it just has more room than he thought.
//
// AND THE HONEST STING. If |w| = 1 really forced w = ±1, there would be two
// pearls and the problem would have been closed in 1859. The reason it is open
// is precisely that the circle has infinitely many directions and only the
// primes decide which ones are used.
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
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= hasta; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 55; i++ {
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

func alDisco(s complex128) complex128 { return 1 - 1/s }

func main() {
	fmt.Println("± EL MÁS-MENOS — la cadena del capitán, y dónde se rompe al cruzar al plano")
	fmt.Println("\n   la cadena que escribió:")
	fmt.Println("        |x| = distancia SIN dirección")
	fmt.Println("        dirección = ± distancia")
	fmt.Println("        x = ±|x|")
	fmt.Println("        ±|1| = w")
	fmt.Println("\n   Los tres primeros renglones son EXACTOS. El cuarto se rompe — y esa rotura")
	fmt.Println("   no es un detalle: es la dificultad entera del problema.")

	fmt.Println("\npescando perlas hasta t=1000…")
	ps := perlas(1000)
	fmt.Printf("perlas: %d\n", len(ps))

	// ---- LEY 1: on the line, two directions ----
	fmt.Println("\nLEY 1 · SOBRE LA RECTA, DOS DIRECCIONES — y ahí el capitán tiene razón entera")
	fmt.Println("   el valor absoluto le saca la dirección al número y deja la distancia pelada.")
	fmt.Println("   El signo se la devuelve. Y con distancia 1 quedan exactamente DOS lugares.")
	fmt.Println("\n        x          |x|        ±|x| devuelve")
	for _, x := range []float64{1, -1, 3.7, -3.7} {
		fmt.Printf("   %8.2f   %8.2f        %+.2f y %+.2f\n", x, math.Abs(x), math.Abs(x), -math.Abs(x))
	}
	fmt.Println("\n   |x| = 1  ⟹  x = +1  ó  x = −1.   DOS puntos. Ni uno más.")

	// ---- LEY 2: in the plane, infinite directions ----
	fmt.Println("\nLEY 2 · EN EL PLANO, INFINITAS DIRECCIONES — y ahí se rompe")
	fmt.Println("   en el plano el tamaño 1 no deja dos lugares: deja un CÍRCULO ENTERO.")
	fmt.Println("\n        |w| = 1   NO da   w = ±1.   Da   w = e^{iφ}, para cualquier ángulo φ.")
	fmt.Println("\n   medido sobre las perlas: todas tienen tamaño 1, y todas están en ángulos")
	fmt.Println("   DISTINTOS. ¿Cuántas valen +1 ó −1?")
	iguales, peorTam := 0, 0.0
	angs := map[float64]bool{}
	for _, g := range ps {
		w := alDisco(complex(0.5, g))
		if d := math.Abs(cmplx.Abs(w) - 1); d > peorTam {
			peorTam = d
		}
		if cmplx.Abs(w-1) < 1e-9 || cmplx.Abs(w+1) < 1e-9 {
			iguales++
		}
		angs[math.Round(cmplx.Phase(w)*1e12)] = true
	}
	fmt.Printf("\n      perlas con tamaño 1 ................ %d de %d (desvío %.1e)\n", len(ps), len(ps), peorTam)
	fmt.Printf("      ángulos DISTINTOS entre ellas ...... %d\n", len(angs))
	fmt.Printf("      perlas que valen +1 ó −1 ........... %d\n", iguales)
	fmt.Println("\n   → NINGUNA. El ± del capitán es el caso de DOS direcciones de un caso de")
	fmt.Println("     INFINITAS, y ±1 son apenas los dos puntos REALES de ese círculo.")

	// ---- LEY 3: and those two points are exactly the two non-pearls ----
	fmt.Println("\nLEY 3 · Y ±1 SON, JUSTAMENTE, LOS DOS PUNTOS DE LA PIEL DONDE NO HAY PERLA")
	fmt.Println("   no es casualidad que ninguna perla caiga ahí. Los dos están ocupados:")
	fmt.Println("\n      w = +1  es EL BROCHE, la imagen de ρ = ∞. Pedirle a una perla finita que")
	fmt.Println("              llegue ahí daría −1 = 0: imposible. Se acercan para siempre.")
	cerca := math.Inf(1)
	for _, g := range []float64{1e3, 1e9, 1e18, 1e30} {
		if d := cmplx.Abs(alDisco(complex(0.5, g)) - 1); d < cerca {
			cerca = d
		}
	}
	fmt.Printf("              medido: a γ = 1e30 la perla queda a %.1e del broche.\n", cerca)
	wMedio := alDisco(complex(0.5, 0))
	zMedio := zetaC(complex(0.5, 0))
	fmt.Printf("\n      w = −1  es la imagen de ρ = ½, donde la línea cruza el eje real.\n")
	fmt.Printf("              w(½) = %+.6f%+.6fi — es −1 exacto.\n", real(wMedio), imag(wMedio))
	fmt.Printf("              ¿y hay una perla ahí? ζ(½) = %+.6f: NO se anula. No hay perla.\n", real(zMedio))
	fmt.Println("\n   → los dos únicos lugares que el ± del capitán señala son los dos únicos")
	fmt.Println("     lugares de la piel SIN perla. Todas las perlas de verdad viven en los")
	fmt.Println("     ángulos de en medio.")

	// ---- LEY 4: where the ± does survive, exactly ----
	fmt.Println("\nLEY 4 · PERO SU ± SÍ SOBREVIVE — y sobrevive EXACTO, con otro nombre")
	fmt.Println("   sobre la recta, las dos direcciones son x y −x. En el plano, el espejo no es")
	fmt.Println("   el signo: es la CONJUGACIÓN. Y sobre la piel pasa algo lindo:")
	fmt.Println("\n        |w| = 1   ⟺   w · conj(w) = 1   ⟺   conj(w) = 1/w")
	fmt.Println("\n   o sea que sobre la piel, conjugar ES invertir. El ± del capitán, en el plano,")
	fmt.Println("   es «la perla y su conjugada». Y toda perla viene en ese par, siempre.")
	fmt.Println("\n        γ           conj(w)                    1/w                  desvío")
	peorConj := 0.0
	for _, i := range []int{0, 3, 49, 299, len(ps) - 1} {
		g := ps[i]
		w := alDisco(complex(0.5, g))
		a, b := cmplx.Conj(w), 1/w
		if d := cmplx.Abs(a - b); d > peorConj {
			peorConj = d
		}
		fmt.Printf("   %11.5f   %9.6f%+9.6fi   %9.6f%+9.6fi     %.1e\n",
			g, real(a), imag(a), real(b), imag(b), cmplx.Abs(a-b))
	}
	for _, g := range ps {
		w := alDisco(complex(0.5, g))
		if d := cmplx.Abs(cmplx.Conj(w) - 1/w); d > peorConj {
			peorConj = d
		}
	}
	fmt.Printf("   → exacto en las %d perlas (%.1e). SU ± ES CORRECTO: solo que en el plano se\n", len(ps), peorConj)
	fmt.Println("     llama conjugación, y tiene más lugar del que él pensaba.")

	// ---- LEY 5: the honest sting ----
	fmt.Println("\nLEY 5 · ⚖️ Y ACÁ LA PARTE QUE DUELE, QUE ES LA QUE VALE")
	fmt.Println("   supongamos por un momento que el capitán tuviera razón literal: que |w| = 1")
	fmt.Println("   obligara a w = ±1. ¿Qué pasaría?")
	fmt.Println("\n        habría DOS perlas. Dos. Y el problema estaría cerrado desde 1859.")
	fmt.Println("\n   La hipótesis está abierta EXACTAMENTE porque el círculo tiene infinitas")
	fmt.Println("   direcciones y solo los primos deciden cuáles se usan. Contando:")
	fmt.Printf("\n      si |w|=1 diera ±1 ................ 2 perlas posibles\n")
	fmt.Printf("      lo que da de verdad .............. infinitas, y ya pescamos %d\n", len(ps))
	fmt.Printf("      y hasta t=1000 la densidad es .... una perla cada %.3f de altura\n",
		(ps[len(ps)-1]-ps[0])/float64(len(ps)-1))
	fmt.Println("\n   → LA ROTURA DE SU CADENA NO ES UN ERROR SUYO: ES EL PROBLEMA.")
	fmt.Println("     El paso de dos direcciones a infinitas es, literalmente, el salto de")
	fmt.Println("     «esto lo resuelvo en una tarde» a «esto lleva ciento sesenta y seis años».")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("LA CADENA DEL CAPITÁN, RENGLÓN POR RENGLÓN:")
	fmt.Println("  ✅ |x| = distancia sin dirección ................ exacto")
	fmt.Println("  ✅ dirección = ± distancia ...................... exacto, en la recta")
	fmt.Println("  ✅ x = ±|x| ..................................... exacto, en la recta")
	fmt.Printf("  ❌ ±|1| = w ..................................... se rompe: |w|=1 da e^{iφ}\n")
	fmt.Printf("     %d perlas con tamaño 1, %d ángulos distintos, %d iguales a ±1\n",
		len(ps), len(angs), iguales)
	fmt.Printf("  ✅ y su ± SÍ vale, con otro nombre: conj(w) = 1/w  %.1e en las %d\n", peorConj, len(ps))
	fmt.Println("\nY LOS DOS PUNTOS QUE SU ± SEÑALA SON LOS DOS ÚNICOS SIN PERLA:")
	fmt.Println("+1 es el broche (la imagen del infinito, que ninguna alcanza) y −1 es ρ = ½,")
	fmt.Println("donde la línea cruza el eje real y la máquina no se calla.")
	fmt.Println("\n⚖️ Y LO QUE HAY QUE LLEVARSE: si |w| = 1 obligara a w = ±1, habría dos perlas y")
	fmt.Println("esto estaría cerrado hace siglo y medio. Está abierto porque el círculo tiene")
	fmt.Println("infinitas direcciones y solo los primos deciden cuáles. Su cadena no falla por")
	fmt.Println("mal pensada: falla en el punto exacto donde el problema se vuelve difícil.")
	fmt.Println("\n¿El premio? Todavía no.")

	escribirLamina(ps, peorTam, peorConj, iguales, len(angs), cerca, real(wMedio), real(zMedio))
}

func escribirLamina(ps []float64, peorTam, peorConj float64, iguales, nang int, cerca, wm, zm float64) {
	var b strings.Builder
	W, H := 1500.0, 1000.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">± EL MÁS-MENOS — dónde se rompe la cadena, y por qué esa rotura ES el problema</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">tres renglones exactos, uno que se rompe al cruzar al plano</text>
`, W, H, W, H, W/2, W/2)

	// the line: two directions
	fmt.Fprintf(&b, `<rect x="40" y="100" width="690" height="330" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="385" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">SOBRE LA RECTA · DOS DIRECCIONES</text>
<line x1="90" y1="240" x2="680" y2="240" stroke="#3d6fa8" stroke-width="2"/>
<circle cx="385" cy="240" r="7" fill="#ffd98a"/>
<text x="385" y="272" font-size="15" text-anchor="middle" font-family="monospace" fill="#ffd98a">0</text>
<circle cx="205" cy="240" r="9" fill="#7ee0c0"/><circle cx="565" cy="240" r="9" fill="#7ee0c0"/>
<text x="205" y="216" font-size="17" text-anchor="middle" font-family="monospace" fill="#7ee0c0">−1</text>
<text x="565" y="216" font-size="17" text-anchor="middle" font-family="monospace" fill="#7ee0c0">+1</text>
<text x="295" y="292" font-size="13" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">distancia 1</text>
<text x="475" y="292" font-size="13" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">distancia 1</text>
<text x="385" y="336" font-size="17" text-anchor="middle" font-family="monospace" fill="#dce8f7">|x| = 1  ⟹  x = ±1</text>
<text x="385" y="366" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">DOS puntos. Ni uno más. Acá el capitán tiene razón entera.</text>
<text x="385" y="400" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el valor absoluto saca la dirección · el signo la devuelve</text>
`)

	// the plane: infinite directions
	cx, cy, R := 1085.0, 262.0, 118.0
	fmt.Fprintf(&b, `<rect x="770" y="100" width="690" height="330" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="1115" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">EN EL PLANO · INFINITAS DIRECCIONES</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#3d6fa8" stroke-width="2"/>
<circle cx="%.0f" cy="%.0f" r="5" fill="#ffd98a"/>
`, cx, cy, R, cx, cy)
	for i, g := range ps {
		if i%3 != 0 {
			continue
		}
		φ := 2 * math.Atan(1/(2*g))
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="2.4" fill="#7ee0c0"/>`, cx+R*math.Cos(φ), cy-R*math.Sin(φ))
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="2.4" fill="#7ee0c0"/>`, cx+R*math.Cos(-φ), cy-R*math.Sin(-φ))
	}
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="7" fill="#ffb27a"/><circle cx="%.0f" cy="%.0f" r="7" fill="#ffb27a"/>
<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" font-family="monospace" fill="#ffb27a">+1</text>
<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" font-family="monospace" fill="#ffb27a">−1</text>
<text x="1115" y="410" font-size="17" text-anchor="middle" font-family="monospace" fill="#dce8f7">|w| = 1  ⟹  w = e^{iφ}   (un círculo entero)</text>
`, cx+R, cy, cx-R, cy, cx+R+22, cy+5, cx-R-22, cy+5)

	fmt.Fprintf(&b, `<rect x="40" y="450" width="1420" height="200" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="750" y="482" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Y ±1 SON, JUSTAMENTE, LOS DOS ÚNICOS PUNTOS DE LA PIEL SIN PERLA</text>
<text x="70" y="518" font-size="15" font-family="monospace" fill="#ffb27a">w = +1</text>
<text x="180" y="518" font-size="14.5" font-family="Georgia" fill="#cfe6ff">EL BROCHE: la imagen de ρ = ∞. Ninguna perla finita llega — a γ = 1e30 todavía queda a %.1e</text>
<text x="70" y="548" font-size="15" font-family="monospace" fill="#ffb27a">w = −1</text>
<text x="180" y="548" font-size="14.5" font-family="Georgia" fill="#cfe6ff">la imagen de ρ = ½, donde la línea cruza el eje real. Ahí ζ(½) = %.4f: NO se anula</text>
<text x="750" y="586" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">medido: %d perlas con tamaño 1 (%.1e) · %d ángulos DISTINTOS · %d iguales a ±1</text>
<text x="750" y="620" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">✅ y su ± SÍ vale con otro nombre: sobre la piel, conjugar ES invertir — conj(w) = 1/w, exacto a %.1e</text>
`, cerca, zm, len(ps), peorTam, nang, iguales, peorConj)

	fmt.Fprintf(&b, `<rect x="40" y="670" width="1420" height="290" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="750" y="706" font-size="20" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚖️ LA PARTE QUE DUELE, QUE ES LA QUE VALE</text>
<text x="750" y="748" font-size="17" text-anchor="middle" font-family="Georgia" fill="#dce8f7">Si |w| = 1 obligara a w = ±1, habría DOS perlas — y esto estaría cerrado desde 1859.</text>
<text x="750" y="782" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Está abierto porque el círculo tiene infinitas direcciones,</text>
<text x="750" y="808" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">y solo los primos deciden cuáles se usan.</text>
<text x="750" y="852" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Su cadena no falla por mal pensada. Falla en el punto EXACTO donde el problema se vuelve difícil:</text>
<text x="750" y="878" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el salto de dos direcciones a infinitas es el salto de «lo resuelvo en una tarde» a ciento sesenta y seis años.</text>
<text x="750" y="924" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffb27a">Tres renglones exactos y uno que se rompe — y la rotura señala el corazón del problema. ¿El premio? Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("mas-menos.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: mas-menos.svg")
}
