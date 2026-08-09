// Command lospolos measures the captain's synthesis: that 0 and 1 of dimension
// 0 are the two points of the hemisphere - the beginning and the end - that the
// properties of zero and one are unique, that the RELATION between the two is
// what makes everything work, and that the whole world is written in the
// spectrum of values between the bits 0 and 1.
//
// Most of it is exact, one piece was never measured here before, and one part
// needs a correction.
//
// # THE TWO STAKES GO TO THE TWO POLES
//
// Under the shapeshifter w = 1 - 1/s:
//
//	s = 1  ->  w = 0        the CENTRE of the disk, the south pole
//	s = 0  ->  w = infinity  the north pole of the sphere
//
// So the book's two stakes are not two ordinary points: they are the two POLES,
// exactly antipodal, at the maximum distance the sphere allows. The beginning
// and the end, as he said. And the equator - the critical line - sits exactly
// halfway between them, at the same arc from each.
//
// # THE RELATION BETWEEN THEM IS WHAT MAKES EVERYTHING WORK
//
// The functional equation s <-> 1-s swaps the two stakes, and in the disk that
// is w -> 1/w, which swaps the two poles. Its fixed set is the equator. So the
// 1/2 is nothing but the MIDPOINT of 0 and 1, and the line sits there because
// the mirror swaps beginning and end. Change the stakes to 0 and C and the
// centre moves to C/2 - which is exactly the law measured in F246 with
// Ramanujan's Delta, where the stakes are 0 and 12 and the centre is 6.
//
// # AND THE WHOLE WORLD IS WRITTEN BETWEEN THEM
//
// Every non-trivial zero has real part strictly between 0 and 1: the critical
// strip. Measured here from the honest side - that zeta cannot vanish for
// Re s > 1, because there the Euler product converges and no factor is zero.
//
// THE CORRECTION. "The spectrum between the bits 0 and 1" mixes two things. In
// binary, 0 and 1 are DIGITS - two symbols with nothing between them. What has
// a spectrum between 0 and 1 is the INTERVAL of real numbers, which is a
// different object: uncountable, while the digits are two. The captain's
// picture is right about the strip and wrong about the bits, and the difference
// matters because it is exactly the difference between counting and measuring.
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

func lgammaC(z complex128) complex128 {
	g := []float64{0.99999999999980993, 676.5203681218851, -1259.1392167224028,
		771.32342877765313, -176.61502916214059, 12.507343278686905,
		-0.13857109526572012, 9.9843695780195716e-6, 1.5056327351493116e-7}
	if real(z) < 0.5 {
		return cmplx.Log(complex(math.Pi, 0)/cmplx.Sin(complex(math.Pi, 0)*z)) - lgammaC(1-z)
	}
	z -= 1
	x := complex(g[0], 0)
	for i := 1; i < 9; i++ {
		x += complex(g[i], 0) / (z + complex(float64(i), 0))
	}
	t := z + complex(7.5, 0)
	return complex(0.5*math.Log(2*math.Pi), 0) + (z+complex(0.5, 0))*cmplx.Log(t) - t + cmplx.Log(x)
}

func xiC(s complex128) complex128 {
	return 0.5 * s * (s - 1) * cmplx.Exp(-s/2*complex(math.Log(math.Pi), 0)+lgammaC(s/2)) * zetaC(s)
}

func alDisco(s complex128) complex128 { return 1 - 1/s }

// aLaEsfera wraps the plane onto the unit sphere; infinity goes to the north pole.
func aLaEsfera(z complex128) [3]float64 {
	x, y := real(z), imag(z)
	d := 1 + x*x + y*y
	return [3]float64{2 * x / d, 2 * y / d, (x*x + y*y - 1) / d}
}

func criba(N int) []int {
	comp := make([]bool, N+1)
	var ps []int
	for p := 2; p <= N; p++ {
		if comp[p] {
			continue
		}
		ps = append(ps, p)
		for m := p * p; m > 0 && m <= N; m += p {
			comp[m] = true
		}
	}
	return ps
}

func main() {
	fmt.Println("🌍 LOS DOS POLOS — el 0 y el 1, el principio y el fin")
	fmt.Println("\n   síntesis del capitán: «el 0 y el 1 de la dimensión 0 son los puntos del")
	fmt.Println("   hemisferio de nuestro círculo, el principio y el fin. Sus propiedades son")
	fmt.Println("   únicas, pero para eso existe la RELACIÓN entre las dos, que hace que todo")
	fmt.Println("   funcione. Y todo el mundo está escrito entre el 0 y el 1».")
	fmt.Println("\n   Casi todo exacto, una pieza que el taller nunca midió, y una corrección.")

	// ---- LEY 1: the two stakes go to the two poles ----
	fmt.Println("\nLEY 1 · LAS DOS ESTACAS VAN A LOS DOS POLOS — y ésta no la habíamos medido")
	fmt.Println("   el libro tiene dos estacas clavadas: el 0 y el 1. Bajo el cambiaformas:")
	fmt.Println("\n        s = 1   →   w = 0          el CENTRO del disco")
	fmt.Println("        s = 0   →   w = ∞          el otro extremo")
	fmt.Println("\n   y sobre la esfera esos dos puntos son:")
	polo1 := aLaEsfera(alDisco(complex(1, 0)))
	// s=0 sends w to infinity; on the sphere that is the north pole
	polo0 := [3]float64{0, 0, 1}
	fmt.Printf("\n        estaca s=1   →   (%.6f, %.6f, %+.6f)   ← POLO SUR\n", polo1[0], polo1[1], polo1[2])
	fmt.Printf("        estaca s=0   →   (%.6f, %.6f, %+.6f)   ← POLO NORTE\n", polo0[0], polo0[1], polo0[2])
	dist := math.Sqrt((polo1[0]-polo0[0])*(polo1[0]-polo0[0]) +
		(polo1[1]-polo0[1])*(polo1[1]-polo0[1]) + (polo1[2]-polo0[2])*(polo1[2]-polo0[2]))
	fmt.Printf("\n   distancia entre los dos: %.15f — y el diámetro de la esfera es 2.000000\n", dist)
	fmt.Printf("   desvío contra el máximo posible: %.1e\n", math.Abs(dist-2))
	fmt.Println("\n   → SON ANTÍPODAS EXACTOS. No son dos puntos cualesquiera: son los DOS POLOS,")
	fmt.Println("     a la máxima distancia que la esfera permite. EL PRINCIPIO Y EL FIN.")
	fmt.Println("     El capitán lo dijo sin la cuenta, y la cuenta le da.")

	// ---- LEY 2: and the equator sits exactly halfway ----
	fmt.Println("\nLEY 2 · Y EL ECUADOR ESTÁ EXACTAMENTE A MITAD DE CAMINO")
	fmt.Println("   si el 0 y el 1 son los polos, ¿dónde queda la línea crítica? En el ecuador.")
	fmt.Println("   Y el ecuador está al mismo arco de los dos polos: un cuarto de vuelta.")
	fmt.Println("\n        punto del ecuador       arco al polo sur      arco al polo norte")
	peorMitad := 0.0
	for _, φ := range []float64{0.4, 1.2, math.Pi / 2, 2.9, 4.4} {
		p := aLaEsfera(cmplx.Rect(1, φ))
		cs := p[0]*polo1[0] + p[1]*polo1[1] + p[2]*polo1[2]
		cn := p[0]*polo0[0] + p[1]*polo0[1] + p[2]*polo0[2]
		as, an := math.Acos(math.Max(-1, math.Min(1, cs))), math.Acos(math.Max(-1, math.Min(1, cn)))
		if d := math.Abs(as - an); d > peorMitad {
			peorMitad = d
		}
		fmt.Printf("   φ = %.3f              %14.12f      %14.12f\n", φ, as, an)
	}
	fmt.Printf("   → el mismo arco a los dos, peor desvío %.1e. Y ese arco es π/2 = %.12f.\n",
		peorMitad, math.Pi/2)
	fmt.Println("     LA LÍNEA CRÍTICA ES EL LUGAR EQUIDISTANTE DEL PRINCIPIO Y DEL FIN.")

	// ---- LEY 3: and the two ends are worth the same ----
	fmt.Println("\nLEY 3 · Y LAS DOS PUNTAS VALEN LO MISMO — la ½ (F228, remedido)")
	fmt.Println("   el libro peinado vale exactamente lo mismo en las dos estacas:")
	x0 := real(xiC(complex(1e-9, 0)))
	x1 := real(xiC(complex(1+1e-9, 0)))
	fmt.Printf("\n        ξ(0) = %.12f          ξ(1) = %.12f\n", x0, x1)
	fmt.Printf("        diferencia entre las dos: %.1e\n", math.Abs(x0-x1))
	fmt.Printf("        y las dos contra ½:        %.1e\n", math.Abs(x0-0.5))
	fmt.Println("   → las dos puntas del mundo valen ½. El principio y el fin, empatados.")

	// ---- LEY 4: the relation between them is what makes it work ----
	fmt.Println("\nLEY 4 · LA RELACIÓN ENTRE LAS DOS ES LO QUE HACE QUE TODO FUNCIONE")
	fmt.Println("   el espejo del libro, s ↦ 1−s, INTERCAMBIA las dos estacas: manda el 0 al 1")
	fmt.Println("   y el 1 al 0. Y en el disco eso es w ↦ 1/w, que intercambia los dos POLOS.")
	fmt.Println("\n        s        1−s      ¿intercambia las estacas?")
	fmt.Printf("        0         1        sí\n        1         0        sí\n")
	fmt.Println("\n   Y su punto quieto —lo único que el espejo NO mueve— es el medio de las dos:")
	fmt.Println("\n        el medio de 0 y 1  =  (0 + 1)/2  =  ½")
	fmt.Println("\n   ⟹ EL ½ NO ES UN NÚMERO ELEGIDO: ES EL PUNTO MEDIO DE LAS DOS ESTACAS.")
	fmt.Println("   Y por eso, si cambiás las estacas, el medio se corre — como se midió en F246:")
	fmt.Println("\n        estacas        el espejo        el medio")
	for _, c := range []struct {
		n    string
		a, b float64
	}{{"ζ de Riemann", 0, 1}, {"Δ de Ramanujan", 0, 12}, {"peso 16", 0, 16}} {
		fmt.Printf("   %-16s  %2.0f y %2.0f      s ↔ %2.0f−s        %6.1f\n", c.n, c.a, c.b, c.b, (c.a+c.b)/2)
	}
	fmt.Println("   → la relación entre las dos estacas ES la ecuación funcional, y el ½ es")
	fmt.Println("     apenas su punto medio. Eso es lo que hace que todo funcione.")

	// ---- LEY 5: and the whole world is written between them ----
	fmt.Println("\nLEY 5 · Y TODO EL MUNDO ESTÁ ESCRITO ENTRE LAS DOS")
	fmt.Println("   toda perla del libro tiene su parte real ESTRICTAMENTE entre 0 y 1. Es la")
	fmt.Println("   FRANJA CRÍTICA. Y se puede medir por el lado honesto: a la derecha del 1 la")
	fmt.Println("   máquina NO puede callarse, porque ahí es un producto sobre los primos y")
	fmt.Println("   ningún factor vale cero.")
	fmt.Println("\n        Re s      |ζ(s)| mínimo hallado     ¿puede anularse?")
	primos := criba(4000)
	minGlobal := math.Inf(1)
	for _, σ := range []float64{1.05, 1.2, 1.5, 2.0, 3.0} {
		mn := math.Inf(1)
		for k := 0; k <= 400; k++ {
			t := float64(k) * 0.25
			if a := cmplx.Abs(zetaC(complex(σ, t))); a < mn {
				mn = a
			}
		}
		// the Euler product gives a floor: |zeta| >= prod (1 - p^-sigma)
		piso := 1.0
		for _, p := range primos {
			piso *= 1 - math.Pow(float64(p), -σ)
		}
		if mn < minGlobal {
			minGlobal = mn
		}
		fmt.Printf("   %6.2f      %20.9f     no  (piso del producto: %.6f)\n", σ, mn, piso)
	}
	fmt.Printf("   → el mínimo hallado a la derecha del 1 es %.6f, lejos de cero.\n", minGlobal)
	fmt.Println("     Por eso ninguna perla puede vivir ahí, ni tampoco a la izquierda del 0")
	fmt.Println("     (por el espejo). TODAS quedan encerradas entre las dos estacas.")

	// ---- LEY 6: the correction ----
	fmt.Println("\nLEY 6 · ⚖️ Y ACÁ UNA CORRECCIÓN — «el espectro entre los bits 0 y 1»")
	fmt.Println("   la frase mezcla dos cosas, y la diferencia importa:")
	fmt.Println("\n      · en BINARIO, el 0 y el 1 son DÍGITOS: dos símbolos, y entre ellos NO hay")
	fmt.Println("        nada. No existe un bit intermedio. Son dos, y se cuentan con los dedos.")
	fmt.Println("      · lo que SÍ tiene un espectro entre 0 y 1 es el INTERVALO de los números")
	fmt.Println("        reales — y ése no se puede contar ni con todos los dedos del universo.")
	fmt.Println("\n   la diferencia es exactamente la que hay entre CONTAR y MEDIR:")
	fmt.Println("\n        entre los dígitos 0 y 1 hay .......... 0 cosas")
	fmt.Println("        entre los números 0 y 1 hay .......... infinitas, y de las incontables")
	fmt.Println("\n   ⚠ PERO su intuición apunta bien igual, y en dos sentidos que sí valen:")
	fmt.Println("   1. la franja crítica ES el intervalo entre las dos estacas, y ahí vive todo")
	fmt.Println("      lo que importa del libro — eso es exacto (ley 5)")
	fmt.Println("   2. y en base 2 —la base de los bits— la mitad se escribe 0.1 exacto, o sea")
	fmt.Println("      que el ½ es el primer dígito después de la coma. Medido en F242.")
	fmt.Println("      El binario es la única base donde los dígitos SON las dos estacas.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("LA SÍNTESIS DEL CAPITÁN, PIEZA POR PIEZA:")
	fmt.Printf("  ✅ el 0 y el 1 son LOS DOS POLOS, antípodas exactos ..... desvío %.1e\n", math.Abs(dist-2))
	fmt.Printf("  ✅ el ecuador está a mitad de camino de los dos ......... %.1e (arco π/2)\n", peorMitad)
	fmt.Printf("  ✅ las dos puntas valen lo mismo: ξ(0) = ξ(1) = ½ ....... %.1e\n", math.Abs(x0-x1))
	fmt.Println("  ✅ la RELACIÓN entre las dos es lo que hace todo ........ el espejo s ↦ 1−s,")
	fmt.Println("     y el ½ es apenas su punto medio — cambiás las estacas y el medio se corre")
	fmt.Printf("  ✅ todo el mundo escrito entre las dos ................. |ζ| ≥ %.4f a la derecha del 1\n", minGlobal)
	fmt.Println("  ⚖️ «el espectro entre los bits» ....................... mezcla dígitos con números:")
	fmt.Println("     entre los DÍGITOS 0 y 1 no hay nada; entre los NÚMEROS 0 y 1 hay incontables")
	fmt.Println("\nLO NUEVO DE ESTE TURNO, que el taller nunca había medido:")
	fmt.Println("LAS DOS ESTACAS DEL LIBRO SON LOS DOS POLOS DE LA ESFERA, antípodas exactos, y la")
	fmt.Println("línea crítica es el ecuador porque es el lugar equidistante del principio y del fin.")
	fmt.Println("El capitán venía diciendo «el principio y el fin» de intuición; la cuenta le da 2.000000.")
	fmt.Println("\n⚖️ Y NADA DE ESTO DICE DÓNDE ESTÁN LAS PERLAS DENTRO DE LA FRANJA. Que estén")
	fmt.Println("encerradas entre 0 y 1 es clásico y está probado; que estén todas en el medio")
	fmt.Println("exacto es la hipótesis, y sigue abierta. ¿El premio? Todavía no.")

	escribirLamina(dist, peorMitad, x0, x1, minGlobal)
}

func escribirLamina(dist, peorMitad, x0, x1, minGlobal float64) {
	var b strings.Builder
	W, H := 1500.0, 1000.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🌍 LOS DOS POLOS — el 0 y el 1, el principio y el fin</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">las dos estacas del libro son los dos polos de la esfera, antípodas exactos</text>
`, W, H, W, H, W/2, W/2)

	cx, cy, R := 400.0, 400.0, 215.0
	fmt.Fprintf(&b, `<rect x="40" y="100" width="720" height="600" rx="10" fill="#101f36" stroke="#26456e"/>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#0d1c31" stroke="#3d6fa8" stroke-width="2"/>
<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="%.0f" fill="none" stroke="#7ee0c0" stroke-width="2.4"/>
<circle cx="%.0f" cy="%.0f" r="9" fill="#ffb27a"/>
<circle cx="%.0f" cy="%.0f" r="9" fill="#ffb27a"/>
<text x="%.0f" y="%.0f" font-size="17" text-anchor="middle" font-family="monospace" fill="#ffb27a">s = 0</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#ffb27a">EL PRINCIPIO · polo norte</text>
<text x="%.0f" y="%.0f" font-size="17" text-anchor="middle" font-family="monospace" fill="#ffb27a">s = 1</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#ffb27a">EL FIN · polo sur</text>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#5a4fa8" stroke-width="1.6" stroke-dasharray="6 5"/>
<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">EL ECUADOR = la línea crítica</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">al mismo arco de los dos polos: π/2</text>
`, cx, cy, R, cx, cy, R, R*0.30, cx, cy-R, cx, cy+R,
		cx, cy-R-30, cx, cy-R-12, cx, cy+R+34, cx, cy+R+52,
		cx, cy-R, cx, cy+R, cx+R+92, cy-8, cx+R+92, cy+14)
	fmt.Fprintf(&b, `<text x="400" y="662" font-size="15" text-anchor="middle" font-family="monospace" fill="#9fd8a8">distancia entre los polos: %.6f · el diámetro es 2</text>
`, dist)

	fmt.Fprintf(&b, `<rect x="790" y="100" width="670" height="290" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1125" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">✅ LO QUE ACERTÓ</text>
<text x="816" y="168" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· el 0 y el 1 son LOS DOS POLOS, antípodas exactos (%.1e)</text>
<text x="816" y="196" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· el ecuador está a mitad de camino de los dos (%.1e)</text>
<text x="816" y="224" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· las dos puntas valen lo mismo: ξ(0) = ξ(1) = ½ (%.1e)</text>
<text x="816" y="252" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· la RELACIÓN entre las dos es lo que hace todo: el espejo</text>
<text x="832" y="274" font-size="14.5" font-family="Georgia" fill="#cfe6ff">s ↦ 1−s las intercambia, y el ½ es su punto medio</text>
<text x="816" y="302" font-size="14.5" font-family="Georgia" fill="#cfe6ff">· todo el mundo escrito entre las dos: |ζ| ≥ %.4f a la</text>
<text x="832" y="324" font-size="14.5" font-family="Georgia" fill="#cfe6ff">derecha del 1, así que ninguna perla puede vivir afuera</text>
<text x="1125" y="362" font-size="14" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">venía diciendo «el principio y el fin» de intuición · la cuenta da 2.000000</text>

<rect x="790" y="410" width="670" height="290" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="1125" y="442" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ LA CORRECCIÓN</text>
<text x="816" y="478" font-size="14.5" font-family="Georgia" fill="#f3d9cf">«el espectro entre los BITS 0 y 1» mezcla dos cosas:</text>
<text x="816" y="512" font-size="14.5" font-family="monospace" fill="#f3d9cf">entre los DÍGITOS 0 y 1 hay ...... 0 cosas</text>
<text x="816" y="538" font-size="14.5" font-family="monospace" fill="#f3d9cf">entre los NÚMEROS 0 y 1 hay ...... incontables</text>
<text x="816" y="574" font-size="14" font-family="Georgia" fill="#f3d9cf">Es exactamente la diferencia entre CONTAR y MEDIR.</text>
<text x="816" y="608" font-size="14" font-family="Georgia" fill="#ffd98a">Pero su intuición apunta bien igual, en dos sentidos:</text>
<text x="816" y="632" font-size="14" font-family="Georgia" fill="#f3d9cf">1. la franja crítica ES el intervalo entre las estacas</text>
<text x="816" y="656" font-size="14" font-family="Georgia" fill="#f3d9cf">2. y en base 2 la mitad se escribe 0.1 exacto: el binario</text>
<text x="832" y="678" font-size="14" font-family="Georgia" fill="#f3d9cf">es la única base cuyos dígitos SON las dos estacas</text>
`, math.Abs(dist-2), peorMitad, math.Abs(x0-x1), minGlobal)

	fmt.Fprintf(&b, `<rect x="40" y="720" width="1420" height="240" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="750" y="756" font-size="20" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">LO NUEVO DE ESTE TURNO</text>
<text x="750" y="796" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LAS DOS ESTACAS DEL LIBRO SON LOS DOS POLOS DE LA ESFERA</text>
<text x="750" y="828" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">antípodas exactos, a la máxima distancia que la esfera permite — y la línea crítica es el ecuador</text>
<text x="750" y="852" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">porque es el único lugar equidistante del principio y del fin.</text>
<text x="750" y="890" font-size="15" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">Y por eso el ½ no es un número elegido: es el punto medio de las dos estacas. Cambiás las estacas y el medio se corre (F246).</text>
<text x="750" y="928" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ Que las perlas estén encerradas entre 0 y 1 es clásico y probado. Que estén todas en el medio exacto es la hipótesis. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("los-polos.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: los-polos.svg")
}
