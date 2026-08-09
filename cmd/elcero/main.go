// Command elcero answers the captain's flash: "zero is the only point that
// unites all the reference, the small and the large; it is where everything
// starts... understanding one zero we understand them all."
//
// The sentence has an exact half and a false half, and both are worth having.
//
// THE EXACT HALF
//
// Carry the book into the disk of dimension 0 with the shapeshifter. Then:
//
//   - EVERY pearl lands at distance exactly 1 from the origin. Not close to 1:
//     one, to the last bit the machine can hold. So the origin IS the only
//     point of the plane from which every pearl is equidistant - and that is
//     the captain's "the point that unites all the reference".
//   - The small and the large meet there too. A low pearl lands at a wide
//     angle, a pearl at height 10^100 lands at an angle of 1e-100, and both
//     sit on the same skin at the same distance. The origin is what they share.
//   - So in the disk EVERY PEARL IS THE SAME PEARL, ROTATED. Same modulus,
//     only the angle changes - and the angle is nothing but the height read in
//     another coordinate.
//   - And in the compressed variable u = n/gamma (F239) even the rotation goes
//     away: two pearls a thousand apart in height give the SAME number at the
//     same u. There they are not similar. They are the same.
//
// THE HONEST HALF
//
//   - "Understanding one we understand them all" is true of the ANATOMY and
//     false of the LOCATION. Every pearl is simple, and around each one the
//     book looks identical - measured here. But no pearl tells you where the
//     next one is: the gaps between consecutive pearls vary by a factor of
//     several, and that wobble is the primes' own handwriting.
//   - And the twist the captain will like: the origin of the disk is NOT a
//     zero. It is s = 1, the POLE, where the machine screams instead of going
//     quiet. The point that unites every pearl is the one place the book blows
//     up.
//
// WHAT DOES DETERMINE EVERYTHING IS ALL OF THEM AT ONCE (Hadamard, and F247's
// measured reconstruction): the zeros together rebuild the book and the primes.
// One does not. That gap between "understanding" and "knowing" is the whole
// difference between a beautiful picture and a proof.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
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

// alDisco is the shapeshifter: it carries a point of the book into the disk.
func alDisco(s complex128) complex128 { return 1 - 1/s }

func main() {
	fmt.Println("⭕ EL CERO — el único punto que une toda la referencia")
	fmt.Println("\n   flash del capitán: «el cero es el único punto que une toda la referencia,")
	fmt.Println("   y lo pequeño y lo grande; es de donde parte todo… entendiendo un cero los")
	fmt.Println("   entendemos a todos».")
	fmt.Println("\n   Tiene una mitad EXACTA y una mitad FALSA, y las dos valen la pena.")

	fmt.Println("\npescando perlas hasta t=1000…")
	ps := perlas(1000)
	fmt.Printf("perlas: %d (de %.6f a %.3f)\n", len(ps), ps[0], ps[len(ps)-1])

	// ---- LEY 1: the origin is the only point equidistant from every pearl ----
	fmt.Println("\nLEY 1 · EL CERO ES EL ÚNICO PUNTO DEL QUE TODAS LAS PERLAS ESTÁN A LA MISMA DISTANCIA")
	fmt.Println("   llevá cada perla al disco con el cambiaformas y medí su distancia a distintos")
	fmt.Println("   candidatos a centro. El único que las ve a TODAS igual de lejos es el cero.")
	fmt.Println("\n      candidato a centro      distancia media    dispersión (máx − mín)")
	type cand struct {
		n string
		z complex128
	}
	mejorD := math.Inf(1)
	for _, c := range []cand{
		{"EL CERO  (0, 0)", 0},
		{"(0.1, 0)", complex(0.1, 0)},
		{"(0, 0.1)", complex(0, 0.1)},
		{"(−0.3, 0)", complex(-0.3, 0)},
		{"(0.5, 0.5)", complex(0.5, 0.5)},
		{"(0.9, 0)", complex(0.9, 0)},
	} {
		mn, mx, suma := math.Inf(1), 0.0, 0.0
		for _, g := range ps {
			d := cmplx.Abs(alDisco(complex(0.5, g)) - c.z)
			if d < mn {
				mn = d
			}
			if d > mx {
				mx = d
			}
			suma += d
		}
		disp := mx - mn
		if disp < mejorD {
			mejorD = disp
		}
		fmt.Printf("   %-22s  %14.9f     %.3e\n", c.n, suma/float64(len(ps)), disp)
	}
	fmt.Printf("\n   → gana el CERO con dispersión %.1e: las %d perlas están TODAS a distancia 1.\n", mejorD, len(ps))
	fmt.Println("     Ningún otro punto del disco las ve parejas. Es literal lo que dijo el capitán:")
	fmt.Println("     EL CERO ES EL ÚNICO PUNTO QUE UNE TODA LA REFERENCIA.")

	// ---- LEY 2: and that zero is the pole, not a zero ----
	fmt.Println("\nLEY 2 · Y ESE CERO NO ES UN CERO: ES EL POLO")
	fmt.Println("   el centro del disco, z = 0, corresponde a s = 1 en el libro. Y en s = 1 la")
	fmt.Println("   máquina no se calla: EXPLOTA. Es su único polo.")
	fmt.Println("\n        s              |ζ(s)|            qué pasa ahí")
	for _, e := range []float64{0.1, 0.01, 1e-4, 1e-6} {
		v := cmplx.Abs(zetaC(complex(1+e, 0)))
		fmt.Printf("   1 + %.0e        %14.2f      se dispara como 1/(s−1)\n", e, v)
	}
	fmt.Println("   → el punto que une a todas las perlas es justo el lugar donde el libro grita.")
	fmt.Println("     Lo pequeño y lo grande se juntan en el único lugar que no es silencio.")

	// ---- LEY 3: every pearl is the same pearl, rotated ----
	fmt.Println("\nLEY 3 · TODAS LAS PERLAS SON LA MISMA PERLA, ROTADA")
	fmt.Println("   en el disco cada perla es un punto de la piel: mismo tamaño, distinto ángulo.")
	fmt.Println("   Y el ángulo no es información nueva: es la altura, leída en otra coordenada.")
	fmt.Println("\n        γ (la altura)        |w| (el tamaño)          el ángulo")
	peorTam := 0.0
	for _, i := range []int{0, 1, 9, 99, 399, len(ps) - 1} {
		g := ps[i]
		w := alDisco(complex(0.5, g))
		if d := math.Abs(cmplx.Abs(w) - 1); d > peorTam {
			peorTam = d
		}
		fmt.Printf("   %14.6f       %18.16f     %14.9f\n", g, cmplx.Abs(w), cmplx.Phase(w))
	}
	for _, g := range ps {
		if d := math.Abs(cmplx.Abs(alDisco(complex(0.5, g))) - 1); d > peorTam {
			peorTam = d
		}
	}
	fmt.Printf("   → el tamaño es 1 en las %d, peor desvío %.1e. Lo ÚNICO que cambia es el ángulo.\n",
		len(ps), peorTam)
	fmt.Println("     Entender una es entender la forma de todas: son la misma, giradas.")

	// ---- LEY 4: in the compressed variable even the rotation goes ----
	fmt.Println("\nLEY 4 · Y EN LA VARIABLE COMPRIMIDA NI SIQUIERA GIRAN — SON LA MISMA (F239)")
	fmt.Println("   el aporte de una perla depende solo de u = n/γ. Dos perlas separadas por mil")
	fmt.Println("   de altura dan EL MISMO número en el mismo u:")
	fmt.Println("\n        u        γ=20          γ=200         γ=2000        γ=20000     exacto")
	peorU := 0.0
	for _, u := range []float64{0.25, 0.50, 1.00, 2.00} {
		fmt.Printf("   %6.2f", u)
		exacto := 4 * math.Sin(u/2) * math.Sin(u/2) / (u * u)
		for _, g := range []float64{20, 200, 2000, 20000} {
			n := u * g
			φ := 2 * math.Atan(1/(2*g))
			v := 4 * math.Sin(n*φ/2) * math.Sin(n*φ/2) / (n * n / (g * g))
			v /= g * g / (g * g)
			ap := 4 * math.Sin(n*φ/2) * math.Sin(n*φ/2) / (u * u)
			_ = v
			if d := math.Abs(ap - exacto); d > peorU {
				peorU = d
			}
			fmt.Printf("   %12.9f", ap)
		}
		fmt.Printf("   %12.9f\n", exacto)
	}
	fmt.Printf("   → una sola forma para todas las alturas (peor desvío %.1e). En esa coordenada\n", peorU)
	fmt.Println("     las perlas no se parecen: SON LA MISMA. Ahí el capitán tiene razón entera.")

	// ---- LEY 5: the honest half ----
	fmt.Println("\nLEY 5 · ⚖️ LA MITAD HONESTA — «entendiendo una las entendemos a todas» es")
	fmt.Println("   cierto de la ANATOMÍA y falso de la UBICACIÓN. Las dos cosas, medidas:")
	fmt.Println("\n   (a) LA ANATOMÍA ES LA MISMA: cada perla es SIMPLE — el libro la cruza con")
	fmt.Println("       pendiente distinta de cero, nunca la roza. Si una fuera doble, la pendiente")
	fmt.Println("       daría cero ahí.")
	fmt.Println("\n        γ            pendiente |Z'(γ)|        ¿simple?")
	minPend := math.Inf(1)
	simples := 0
	for _, i := range []int{0, 1, 4, 49, 199, len(ps) - 1} {
		g := ps[i]
		h := 1e-5
		d := math.Abs((zOf(g+h) - zOf(g-h)) / (2 * h))
		fmt.Printf("   %12.6f     %18.9f        %s\n", g, d, map[bool]string{true: "sí", false: "NO"}[d > 1e-3])
	}
	for _, g := range ps {
		h := 1e-5
		d := math.Abs((zOf(g+h) - zOf(g-h)) / (2 * h))
		if d < minPend {
			minPend = d
		}
		if d > 1e-3 {
			simples++
		}
	}
	fmt.Printf("   → %d de %d perlas son simples; la pendiente más chica de todas es %.4f.\n",
		simples, len(ps), minPend)
	fmt.Println("     Ninguna se rozó. La anatomía local es idéntica en todas: entendés una,")
	fmt.Println("     entendés cómo es cualquiera por dentro.")

	fmt.Println("\n   (b) LA UBICACIÓN NO: ninguna perla te dice dónde está la siguiente.")
	var huecos []float64
	for i := 1; i < len(ps); i++ {
		huecos = append(huecos, ps[i]-ps[i-1])
	}
	ordenados := append([]float64(nil), huecos...)
	sort.Float64s(ordenados)
	suma := 0.0
	for _, h := range huecos {
		suma += h
	}
	medio := suma / float64(len(huecos))
	fmt.Printf("\n        huecos medidos: %d · el más chico %.4f · el más grande %.4f · medio %.4f\n",
		len(huecos), ordenados[0], ordenados[len(ordenados)-1], medio)
	fmt.Printf("        el más grande es %.1f veces el más chico\n", ordenados[len(ordenados)-1]/ordenados[0])
	fmt.Println("   → si una perla contara dónde está la próxima, todos los huecos serían iguales.")
	fmt.Println("     No lo son. Ese temblor es la caligrafía de los números primos, y es")
	fmt.Println("     EXACTAMENTE la información que una sola perla no tiene.")

	// ---- LEY 6: all of them together do determine everything ----
	fmt.Println("\nLEY 6 · PERO TODAS JUNTAS SÍ DETERMINAN TODO")
	fmt.Println("   Hadamard probó que el libro se reconstruye entero multiplicando un factor por")
	fmt.Println("   cada perla. Y en F247 el taller lo midió: con 269 perlas y NADA MÁS salió la")
	fmt.Println("   escalera de los primos, a desvío medio 0.0945.")
	fmt.Println("\n        una perla      → su forma, sí. Su vecina, no.")
	fmt.Println("        todas juntas   → el libro entero, y con él los primos.")
	fmt.Println("\n   Esa diferencia entre ENTENDER y SABER es toda la distancia que falta.")

	fmt.Println("\nLEY 7 · «ES LA QUIETUD» — y el capitán le puso el nombre exacto")
	fmt.Println("   un espejo mueve todo… menos un punto. Ese punto es el centro, y es el único")
	fmt.Println("   que se queda quieto. El capitán lo escribió solo:")
	fmt.Println("\n        0 = ( x + (−x) ) / 2      el cero es el medio entre un número y su opuesto")
	fmt.Println("\n   Y ésa es exactamente la ley que el taller midió en F246, en el caso C = 0:")
	fmt.Println("\n        el espejo x ↦ C − x  tiene un solo punto quieto, y está en C/2")
	fmt.Println("\n        espejo            qué mueve            el punto QUIETO")
	for _, c := range []struct {
		n string
		C float64
	}{{"x ↦ −x        (el del capitán)", 0}, {"x ↦ 1−x       (el del libro)", 1},
		{"s ↦ 12−s      (la Δ de Ramanujan)", 12}, {"x ↦ 7−x       (uno inventado)", 7}} {
		quieto := c.C / 2
		peor, mueve := 0.0, 0.0
		for k := -200; k <= 200; k++ {
			x := quieto + float64(k)*0.037
			d := math.Abs((c.C - x) - x)
			if math.Abs(x-quieto) < 1e-12 {
				if d > peor {
					peor = d
				}
			} else if d > mueve {
				mueve = d
			}
		}
		fmt.Printf("   %-34s  hasta %7.2f        %6.2f  (desvío %.1e)\n", c.n, mueve, quieto, peor)
	}
	fmt.Println("   → cada espejo tiene UN solo punto quieto y está en la mitad. El del capitán da 0;")
	fmt.Println("     el del libro da ½; el de Ramanujan da 6. La misma ley, distinto espejo.")
	fmt.Println("\n   ⚡ Y ACÁ VIENE LO QUE NO ESPERABA. Hay DOS quietudes distintas:")
	fmt.Println("\n        la del ESPEJO      : cuánto te mueve  =  | (1−s) − s |")
	fmt.Println("        la de la DERIVADA  : d/ds [ s(1−s) ]  =  | 1 − 2s |")
	fmt.Println("\n   Son LA MISMA EXPRESIÓN. No parecidas: la misma, letra por letra.")
	fmt.Println("\n        s              |(1−s) − s|          |d/ds s(1−s)|        desvío")
	peorQ := 0.0
	for _, s := range []complex128{complex(0.5, 0), complex(0.3, 2), complex(0.5, 14.134725), complex(2, -3), complex(-1, 0.7)} {
		a := cmplx.Abs((1 - s) - s)
		d := cmplx.Abs(1 - 2*s)
		if v := math.Abs(a - d); v > peorQ {
			peorQ = v
		}
		fmt.Printf("   %5.2f%+7.2fi     %16.9f     %16.9f     %.1e\n", real(s), imag(s), a, d, math.Abs(a-d))
	}
	fmt.Printf("   → idénticas (%.1e). Entonces «el punto que el espejo no mueve» y «el punto\n", peorQ)
	fmt.Println("     donde muere la derivada» son EL MISMO PUNTO, y no por casualidad: s(1−s) es")
	fmt.Println("     la cuadrática que respeta el espejo, así que su derivada MIDE cuánto te mueve.")
	fmt.Println("\n   ⟹ LA QUIETUD TIENE DOS CARAS Y SON UNA:")
	fmt.Println("        el espejo no te mueve  ⟺  la derivada se muere  ⟺  estás en el centro")
	fmt.Println("\n   Por eso el ½ del libro es el punto quieto, y por eso el 0 del disco es el punto")
	fmt.Println("   quieto de su propio espejo. El capitán dijo «es la quietud» y le acertó al nombre")
	fmt.Println("   técnico: en matemática se llama PUNTO FIJO, pero quietud es mejor palabra.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("EL FLASH DEL CAPITÁN, PARTIDO EN DOS Y MEDIDO:")
	fmt.Printf("  ✅ el cero une toda la referencia ................ dispersión %.1e contra el resto\n", mejorD)
	fmt.Printf("  ✅ todas las perlas son la misma, rotada ......... |w|=1 en %d, desvío %.1e\n", len(ps), peorTam)
	fmt.Printf("  ✅ y en la variable comprimida, la misma a secas .. %.1e\n", peorU)
	fmt.Printf("  ✅ la anatomía es idéntica: %d/%d simples ....... pendiente mínima %.4f\n", simples, len(ps), minPend)
	fmt.Printf("  ❌ la ubicación NO se deduce de una .............. el hueco mayor es %.1f× el menor\n",
		ordenados[len(ordenados)-1]/ordenados[0])
	fmt.Println("  🌀 y el punto que las une no es un cero: ES EL POLO, donde el libro grita")
	fmt.Println("\nASÍ QUE SÍ, CAPITÁN, PERO CON UNA COMA:")
	fmt.Println("Entendiendo una perla entendés CÓMO SON todas — su tamaño, su forma, su anatomía.")
	fmt.Println("Lo que no entendés es DÓNDE están las otras. Y el premio no paga por entender")
	fmt.Println("cómo son: paga por demostrar dónde están.")
	fmt.Println("\n⚖️ Todo lo de arriba está medido y es exacto. Nada de eso demuestra la hipótesis:")
	fmt.Println("que todas tengan tamaño 1 en el disco ES estar sobre la línea, así que medirlo en")
	fmt.Println("las que ya pescamos sobre la línea no prueba nada sobre las que no vimos.")
	fmt.Println("¿El premio? Todavía no.")

	escribirLamina(ps, mejorD, peorTam, peorU, simples, minPend, huecos, ordenados, medio)
}

func escribirLamina(ps []float64, mejorD, peorTam, peorU float64, simples int,
	minPend float64, huecos, ordenados []float64, medio float64) {

	var b strings.Builder
	W, H := 1500.0, 1030.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⭕ EL CERO — el único punto que une toda la referencia</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">el flash del capitán tiene una mitad exacta y una mitad falsa, y las dos valen</text>
`, W, H, W, H, W/2, W/2)

	// the disk with every pearl on the skin
	cx, cy, R := 340.0, 400.0, 240.0
	fmt.Fprintf(&b, `<rect x="40" y="100" width="600" height="600" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="340" y="132" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">TODAS A LA MISMA DISTANCIA DEL CERO</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#3d6fa8" stroke-width="2"/>
`, cx, cy, R)
	for i, g := range ps {
		if i%3 != 0 {
			continue
		}
		φ := 2 * math.Atan(1/(2*g))
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.2f" y2="%.2f" stroke="#1f4470" stroke-width="0.6" opacity="0.55"/>`,
			cx, cy, cx+R*math.Cos(φ), cy-R*math.Sin(φ))
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.2f" y2="%.2f" stroke="#1f4470" stroke-width="0.6" opacity="0.55"/>`,
			cx, cy, cx+R*math.Cos(-φ), cy-R*math.Sin(-φ))
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="2.6" fill="#7ee0c0"/>`, cx+R*math.Cos(φ), cy-R*math.Sin(φ))
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="2.6" fill="#7ee0c0"/>`, cx+R*math.Cos(-φ), cy-R*math.Sin(-φ))
	}
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="7" fill="#ffb27a"/>
<text x="%.0f" y="%.0f" font-size="15" font-family="Georgia" fill="#ffb27a">el cero</text>
<text x="340" y="672" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">%d perlas, todas a distancia 1 · dispersión %.1e</text>
<text x="340" y="694" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">ningún otro punto del disco las ve parejas</text>
`, cx, cy, cx+14, cy+24, len(ps), mejorD)

	// the two halves
	fmt.Fprintf(&b, `<rect x="670" y="100" width="790" height="284" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1065" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">✅ LA MITAD EXACTA</text>
<text x="696" y="168" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· el cero es el ÚNICO punto que une toda la referencia</text>
<text x="696" y="196" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· todas las perlas son LA MISMA PERLA, rotada</text>
<text x="716" y="218" font-size="13" font-family="monospace" fill="#7ee0c0">|w| = 1 en %d perlas · desvío %.1e</text>
<text x="696" y="248" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· en la variable comprimida ni siquiera giran: son la misma</text>
<text x="716" y="270" font-size="13" font-family="monospace" fill="#7ee0c0">una sola forma para toda altura · %.1e</text>
<text x="696" y="300" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· la anatomía es idéntica: todas simples, ninguna se roza</text>
<text x="716" y="322" font-size="13" font-family="monospace" fill="#7ee0c0">%d/%d simples · pendiente mínima %.4f</text>
<text x="696" y="356" font-size="13.5" font-family="Georgia" fill="#9fd8a8">entendiendo una entendés CÓMO SON todas</text>

<rect x="670" y="400" width="790" height="300" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="1065" y="432" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">❌ LA MITAD QUE NO</text>
<text x="696" y="468" font-size="14.5" font-family="Georgia" fill="#f3d9cf">ninguna perla te dice dónde está la siguiente.</text>
<text x="696" y="496" font-size="13.5" font-family="monospace" fill="#f3d9cf">hueco más chico %.4f · más grande %.4f · medio %.4f</text>
<text x="696" y="518" font-size="13.5" font-family="monospace" fill="#ffb27a">el mayor es %.1f veces el menor</text>
<text x="696" y="550" font-size="14" font-family="Georgia" fill="#f3d9cf">si una contara dónde está la próxima, todos los huecos</text>
<text x="696" y="570" font-size="14" font-family="Georgia" fill="#f3d9cf">serían iguales. No lo son. Ese temblor es la caligrafía</text>
<text x="696" y="590" font-size="14" font-family="Georgia" fill="#f3d9cf">de los primos — justo lo que una sola perla no tiene.</text>
<text x="696" y="626" font-size="14.5" font-family="Georgia" fill="#ffd98a">🌀 y el punto que las une NO es un cero: es el POLO,</text>
<text x="696" y="648" font-size="14.5" font-family="Georgia" fill="#ffd98a">   el único lugar donde el libro grita en vez de callarse.</text>
<text x="696" y="682" font-size="13.5" font-family="Georgia" fill="#f3d9cf">lo pequeño y lo grande se juntan en el único punto que no es silencio.</text>
`, len(ps), peorTam, peorU, simples, len(ps), minPend,
		ordenados[0], ordenados[len(ordenados)-1], medio, ordenados[len(ordenados)-1]/ordenados[0])

	fmt.Fprintf(&b, `<rect x="40" y="720" width="1420" height="270" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="750" y="754" font-size="20" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">SÍ, CAPITÁN — PERO CON UNA COMA</text>
<text x="750" y="794" font-size="17" text-anchor="middle" font-family="Georgia" fill="#dce8f7">Entendiendo una perla entendés CÓMO SON todas: su tamaño, su forma, su anatomía.</text>
<text x="750" y="822" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Lo que no entendés es DÓNDE están las otras.</text>
<text x="750" y="856" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Y todas juntas sí determinan todo: Hadamard reconstruye el libro entero con un factor por perla,</text>
<text x="750" y="878" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">y en F247 el taller lo midió — 269 perlas y nada más devolvieron la escalera de los primos.</text>
<text x="750" y="912" font-size="16" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">Esa diferencia entre ENTENDER y SABER es toda la distancia que falta.</text>
<text x="750" y="950" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ que todas tengan tamaño 1 en el disco ES estar sobre la línea: medirlo en las que ya pescamos</text>
<text x="750" y="972" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffb27a">sobre la línea no prueba nada sobre las que no vimos. El premio no paga por entender. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("el-cero.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: el-cero.svg")
}
