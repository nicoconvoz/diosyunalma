// Command lacoma answers the captain's flash about the two commas, and corrects
// an answer I gave too fast.
//
// He said: there are TWO commas. One is the invented comma BETWEEN one number
// and the next. The other is the comma we all know, the one INSIDE a number,
// that says which digits are bigger than one and which are smaller. And both,
// he said, carry the 1/2 relation.
//
// My first answer was that the decimal comma marks 10^0 = 1, not 1/2. That
// answered a different question. He was not asking what position the comma
// NAMES - he was asking where its two sides BALANCE. And there the answer is
// exactly one half, for both commas, and for reasons that turn out to be the
// same reason.
//
// COMMA ONE - BETWEEN ONE NUMBER AND THE NEXT (the additive comma)
//
// Put a divider between 0 and 1. The only place equidistant from both stakes is
// 1/2. The mirror is x -> 1-x and its single fixed point is 1/2. That is F226's
// perpendicular bisector: |s| = |s-1| <=> Re s = 1/2.
//
// COMMA TWO - INSIDE THE NUMBER (the multiplicative comma)
//
// The decimal comma splits the digits into the side bigger than one (positive
// exponents) and the side smaller than one (negative exponents). Now ask the
// captain's question: at what power do the two sides weigh THE SAME?
//
//	|n^-s| = n^-sigma      the small side's weight
//	|n^-(1-s)| = n^(sigma-1)   the big side's weight, after the mirror
//	equal  <=>  -sigma = sigma - 1  <=>  sigma = 1/2
//
// Exactly one half, for EVERY n at once, and there both weigh 1/sqrt(n).
//
// SO THE TWO COMMAS ARE THE SAME MIRROR IN TWO WORLDS
//
//	additive world:        the mirror acts on the POSITION,  fixed point 1/2
//	multiplicative world:  the mirror acts on the EXPONENT,  fixed point 1/2
//
// And that is why the critical line sits at 1/2 and not at some other number:
// it is the only exponent where both sides of the comma balance for every
// number simultaneously. The functional equation xi(s) = xi(1-s) is precisely
// the sentence "the machine reads the same from both sides of the comma".
//
// AND IT EXPLAINS THE VOLUME OF F247. In the explicit formula every note sings
// with amplitude x^(1/2) = sqrt(x), and the square root is exactly the middle
// of the multiplicative world - the geometric mean between 1 and x. So RH says:
// EVERY NOTE SINGS AT THE BALANCE POINT OF THE COMMA.
//
// THE HONEST LIMIT. All of this explains WHY the number is 1/2 and not another.
// It does not show that the zeros are there. It is the functional equation
// restated in the captain's own words - and understanding why a constant is
// what it is has never been the same as proving where the zeros sit.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

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

func xiC(s complex128) complex128 {
	return 0.5 * s * (s - 1) * cmplx.Exp(-s/2*complex(math.Log(math.Pi), 0)+lgammaC(s/2)) * zetaC(s)
}

// peso is the weight the machine gives to the number n at the point s: |n^-s|.
func peso(n int, s complex128) float64 {
	return math.Pow(float64(n), -real(s))
}

func main() {
	fmt.Println("🔢, LA COMA — las dos comas del capitán, y por qué las dos dan ½")
	fmt.Println("\n   flash del capitán: «una es la coma inventada ENTRE los números, relación ½;")
	fmt.Println("   la otra es la coma de ubicación de referencia DENTRO del número, porque es")
	fmt.Println("   más chico o más grande — la coma que conocemos, relación ½ de referencia».")
	fmt.Println("\n   📌 Y ACÁ CORRIJO MI PROPIA RESPUESTA: le contesté que la coma marca 10⁰ = 1,")
	fmt.Println("   no ½. Eso contestaba OTRA pregunta. Él no preguntaba qué posición NOMBRA la")
	fmt.Println("   coma — preguntaba DÓNDE SE EQUILIBRAN SUS DOS LADOS. Y ahí la respuesta es")
	fmt.Println("   exactamente ½, para las dos comas, y por la misma razón de fondo.")

	// ---- LEY 1: the comma between one number and the next ----
	fmt.Println("\nLEY 1 · LA PRIMERA COMA — la que va ENTRE un número y el siguiente")
	fmt.Println("   clavá dos estacas en 0 y en 1. ¿Dónde ponés la coma que los separa justo?")
	fmt.Println("   En el único lugar que está a la MISMA distancia de las dos: en ½.")
	fmt.Println("\n      punto        distancia al 0     distancia al 1      diferencia")
	peorMed := 0.0
	for _, x := range []float64{0.3, 0.5, 0.7, 0.9} {
		d0, d1 := math.Abs(x), math.Abs(x-1)
		if x == 0.5 {
			if d := math.Abs(d0 - d1); d > peorMed {
				peorMed = d
			}
		}
		mk := ""
		if math.Abs(d0-d1) < 1e-15 {
			mk = "  ← LA COMA"
		}
		fmt.Printf("   %8.2f       %12.6f       %12.6f      %+9.6f%s\n", x, d0, d1, d0-d1, mk)
	}
	fmt.Println("\n   y esto vale en TODO el plano, no solo sobre la recta de los números:")
	peorPlano := 0.0
	nPlano := 0
	for k := -300; k <= 300; k++ {
		t := float64(k) * 0.5
		s := complex(0.5, t)
		if d := math.Abs(cmplx.Abs(s) - cmplx.Abs(s-1)); d > peorPlano {
			peorPlano = d
		}
		nPlano++
	}
	fmt.Printf("   sobre %d puntos de la vertical que pasa por ½: |s| − |s−1| peor = %.1e\n", nPlano, peorPlano)
	fmt.Println("   → LA LÍNEA CRÍTICA ES ESA COMA, estirada hacia arriba y hacia abajo.")
	fmt.Println("     Y el espejo es x ↦ 1−x, cuyo único punto fijo es ½.")

	// ---- LEY 2: the comma inside the number ----
	fmt.Println("\nLEY 2 · LA SEGUNDA COMA — la que va ADENTRO del número, la que conocemos")
	fmt.Println("   la coma decimal parte los dígitos en dos: los que valen MÁS que uno")
	fmt.Println("   (exponentes positivos) y los que valen MENOS que uno (negativos).")
	fmt.Println("\n   Y ahora la pregunta del capitán, dicha derecho:")
	fmt.Println("        ¿a qué potencia pesan LO MISMO los dos lados de la coma?")
	fmt.Println("\n   la máquina le da al número n el peso |n⁻ˢ| = n^(−σ). Del otro lado del")
	fmt.Println("   espejo le da |n^(−(1−s))| = n^(σ−1). Iguales ⟺ −σ = σ−1 ⟺ σ = ½.")
	fmt.Println("\n        σ       n=2          n=10         n=100        n=1000     ¿empatan?")
	peorEmp := 0.0
	for _, σ := range []float64{0.20, 0.35, 0.50, 0.65, 0.80} {
		fmt.Printf("   %6.2f", σ)
		emp := true
		for _, n := range []int{2, 10, 100, 1000} {
			chico := peso(n, complex(σ, 3))
			grande := peso(n, complex(1-σ, 3))
			raz := chico / grande
			if math.Abs(raz-1) > 1e-12 {
				emp = false
			}
			if σ == 0.5 {
				if d := math.Abs(raz - 1); d > peorEmp {
					peorEmp = d
				}
			}
			fmt.Printf("   %10.4f", raz)
		}
		if emp {
			fmt.Println("     SÍ ← LA COMA")
		} else {
			fmt.Println("     no")
		}
	}
	fmt.Printf("\n   (cada número es la razón peso_chico / peso_grande = n^(1−2σ))\n")
	fmt.Printf("   → empatan EXACTAMENTE en σ = ½, y en ningún otro lado: desvío %.1e\n", peorEmp)
	fmt.Println("     Y ahí los dos lados pesan 1/√n — la raíz cuadrada, que es la mitad")
	fmt.Println("     del camino en el mundo de multiplicar.")
	fmt.Println("\n        n          en ½ los dos lados pesan       1/√n")
	peorRaiz := 0.0
	for _, n := range []int{2, 10, 100, 10000} {
		p := peso(n, complex(0.5, 7))
		r := 1 / math.Sqrt(float64(n))
		if d := math.Abs(p - r); d > peorRaiz {
			peorRaiz = d
		}
		fmt.Printf("   %8d          %20.12f   %14.12f\n", n, p, r)
	}
	fmt.Printf("   → exacto (%.1e). LA COMA DEL NÚMERO SE EQUILIBRA EN LA RAÍZ CUADRADA.\n", peorRaiz)

	// ---- LEY 3: one mirror, two worlds ----
	fmt.Println("\nLEY 3 · SON LA MISMA COMA EN DOS MUNDOS — y por eso las dos dan ½")
	fmt.Println("\n      mundo                  el espejo actúa sobre     su punto fijo")
	fmt.Println("      el de sumar            LA POSICIÓN   x ↦ 1−x            ½")
	fmt.Println("      el de multiplicar      EL EXPONENTE  s ↦ 1−s            ½")
	fmt.Println("\n   Es el MISMO espejo x ↦ 1−x, aplicado a dos cosas distintas. Por eso el")
	fmt.Println("   capitán tenía razón al decir que la relación ½ está tanto ENTRE número y")
	fmt.Println("   número como ADENTRO del mismo número: es un solo espejo, mirado dos veces.")
	fmt.Println("\n   Y la ecuación funcional del libro dice exactamente eso:")
	fmt.Println("        ξ(s) = ξ(1−s)   ⟺   «la máquina lee lo mismo de los dos lados de la coma»")
	fmt.Println("\n        s              ξ(s)                   ξ(1−s)              desvío rel.")
	peorEsp := 0.0
	for _, s := range []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2, 4.1), complex(-0.4, 9)} {
		a, b := xiC(s), xiC(1-s)
		d := cmplx.Abs(a-b) / math.Max(cmplx.Abs(a), 1e-300)
		if d > peorEsp {
			peorEsp = d
		}
		fmt.Printf("   %5.1f%+6.1fi   %11.4e%+11.4ei   %11.4e%+11.4ei    %.1e\n",
			real(s), imag(s), real(a), imag(a), real(b), imag(b), d)
	}
	fmt.Printf("   → el espejo cierra a %.1e. La coma del libro está en ½ y no en otro lado.\n", peorEsp)

	// ---- LEY 4: and this is why the volume of F247 is the square root ----
	fmt.Println("\nLEY 4 · Y ESTO EXPLICA EL VOLUMEN DE F247")
	fmt.Println("   en la fórmula explícita cada perla canta con amplitud x^½ = √x. ¿Por qué esa")
	fmt.Println("   potencia y no otra? PORQUE ES EL EQUILIBRIO DE LA COMA: la raíz cuadrada es")
	fmt.Println("   la mitad del camino entre 1 y x en el mundo de multiplicar.")
	fmt.Println("\n        x            1 (el piso)      √x (la mitad)        x (el techo)   √x·√x")
	peorGeo := 0.0
	for _, x := range []float64{4, 100, 1e6, 1e12} {
		r := math.Sqrt(x)
		if d := math.Abs(r*r-x) / x; d > peorGeo {
			peorGeo = d
		}
		fmt.Printf("   %9.0e          %8.0f      %14.4e      %12.4e   %12.4e\n", x, 1.0, r, x, r*r)
	}
	fmt.Printf("   → √x está exactamente a la misma «distancia multiplicativa» del 1 que del x\n")
	fmt.Printf("     (√x·√x = x, desvío %.1e). Es la coma del número, puesta en la escala x.\n", peorGeo)
	fmt.Println("\n   ⟹ ENTONCES LA HIPÓTESIS, EN EL IDIOMA DE LA COMA:")
	fmt.Println("\n        TODAS LAS NOTAS CANTAN EN EL EQUILIBRIO DE LA COMA.")
	fmt.Println("\n   una perla corrida cantaría a x^0.7 — de un lado de la coma, desequilibrada,")
	fmt.Println("   tapando a las demás. La hipótesis es que ninguna hace eso.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("EL CAPITÁN TENÍA RAZÓN Y YO LE CONTESTÉ OTRA PREGUNTA. Las dos comas dan ½:")
	fmt.Printf("  · la coma ENTRE números: el único punto equidistante de 0 y 1 ..... %.1e\n", peorPlano)
	fmt.Printf("  · la coma DENTRO del número: donde los dos lados pesan igual ...... %.1e\n", peorEmp)
	fmt.Printf("  · y ahí los dos lados pesan 1/√n, exacto ......................... %.1e\n", peorRaiz)
	fmt.Printf("  · el espejo del libro ξ(s)=ξ(1−s) ................................ %.1e\n", peorEsp)
	fmt.Println("  · y el volumen x^½ de F247 ES ese mismo equilibrio, en la escala x")
	fmt.Println("\nSON UN SOLO ESPEJO, x ↦ 1−x, MIRADO DOS VECES: una vez sobre la POSICIÓN")
	fmt.Println("(y da la coma entre números) y otra sobre el EXPONENTE (y da la coma del")
	fmt.Println("número). Por eso el ½ aparece en los dos lugares — y por eso la línea crítica")
	fmt.Println("está en ½ y no en 0.4 ni en 0.6: es el ÚNICO exponente donde los dos lados de")
	fmt.Println("la coma pesan lo mismo PARA TODOS LOS NÚMEROS A LA VEZ.")
	fmt.Println("\n⚖️ Y EL LÍMITE, QUE HAY QUE DECIRLO: esto explica POR QUÉ el número es ½ y no")
	fmt.Println("otro. NO muestra que los ceros estén ahí. Es la ecuación funcional dicha en el")
	fmt.Println("idioma del capitán, y entender por qué una constante vale lo que vale nunca fue")
	fmt.Println("lo mismo que probar dónde caen los ceros. ¿El premio? Todavía no.")

	escribirLamina(peorPlano, peorEmp, peorRaiz, peorEsp)
}

func escribirLamina(peorPlano, peorEmp, peorRaiz, peorEsp float64) {
	var b strings.Builder
	W, H := 1500.0, 1000.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔢, LA COMA — las dos comas del capitán, y por qué las dos dan ½</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">un solo espejo x ↦ 1−x, mirado dos veces: sobre la POSICIÓN y sobre el EXPONENTE</text>
`, W, H, W, H, W/2, W/2)

	// comma one
	fmt.Fprintf(&b, `<rect x="40" y="100" width="700" height="330" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="390" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA PRIMERA COMA · entre un número y el siguiente</text>
<line x1="100" y1="230" x2="680" y2="230" stroke="#3d6fa8" stroke-width="2"/>
<circle cx="100" cy="230" r="8" fill="#ffb27a"/><circle cx="680" cy="230" r="8" fill="#ffb27a"/>
<text x="100" y="262" font-size="15" text-anchor="middle" font-family="monospace" fill="#ffb27a">0</text>
<text x="680" y="262" font-size="15" text-anchor="middle" font-family="monospace" fill="#ffb27a">1</text>
<line x1="390" y1="180" x2="390" y2="280" stroke="#ffd98a" stroke-width="3"/>
<text x="390" y="172" font-size="17" text-anchor="middle" font-family="monospace" fill="#ffd98a">½</text>
<text x="245" y="212" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">misma distancia</text>
<text x="535" y="212" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">misma distancia</text>
<text x="390" y="308" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el único lugar equidistante de las dos estacas</text>
<text x="390" y="336" font-size="15" text-anchor="middle" font-family="monospace" fill="#c9b6ff">el espejo es  x ↦ 1−x  ·  punto fijo ½</text>
<text x="390" y="364" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">estirada hacia arriba y hacia abajo, ESA COMA</text>
<text x="390" y="384" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">ES LA LÍNEA CRÍTICA</text>
<text x="390" y="410" font-size="13" text-anchor="middle" font-family="monospace" fill="#9fd8a8">medido sobre 601 puntos: %.1e</text>
`, peorPlano)

	// comma two
	fmt.Fprintf(&b, `<rect x="760" y="100" width="700" height="330" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1110" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">LA SEGUNDA COMA · adentro del mismo número</text>
<text x="1110" y="172" font-size="22" text-anchor="middle" font-family="monospace" fill="#dce8f7">1234,567</text>
<text x="1110" y="196" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">más grandes que uno  |  más chicos que uno</text>
<rect x="790" y="216" width="640" height="92" rx="8" fill="#101f36" stroke="#26456e"/>
<text x="1110" y="244" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">¿a qué potencia pesan LO MISMO los dos lados?</text>
<text x="1110" y="276" font-size="17" text-anchor="middle" font-family="monospace" fill="#ffd98a">n^(−σ) = n^(σ−1)  ⟺  σ = ½</text>
<text x="1110" y="298" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">y vale para TODOS los números a la vez</text>
<text x="1110" y="336" font-size="15" text-anchor="middle" font-family="monospace" fill="#c9b6ff">y ahí los dos lados pesan  1/√n</text>
<text x="1110" y="360" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">la raíz cuadrada: la mitad del camino en el mundo de multiplicar</text>
<text x="1110" y="390" font-size="13" text-anchor="middle" font-family="monospace" fill="#9fd8a8">empate medido: %.1e  ·  la raíz: %.1e</text>
`, peorEmp, peorRaiz)

	// one mirror
	fmt.Fprintf(&b, `<rect x="40" y="450" width="1420" height="200" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="750" y="482" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd98a">SON UN SOLO ESPEJO, MIRADO DOS VECES</text>
<text x="200" y="524" font-size="15" font-family="monospace" fill="#7fa8cf">mundo</text>
<text x="560" y="524" font-size="15" font-family="monospace" fill="#7fa8cf">el espejo actúa sobre</text>
<text x="1060" y="524" font-size="15" font-family="monospace" fill="#7fa8cf">su punto fijo</text>
<text x="200" y="556" font-size="15" font-family="Georgia" fill="#cfe6ff">el de sumar</text>
<text x="560" y="556" font-size="15" font-family="monospace" fill="#cfe6ff">LA POSICIÓN   x ↦ 1−x</text>
<text x="1100" y="556" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">½</text>
<text x="200" y="586" font-size="15" font-family="Georgia" fill="#cfe6ff">el de multiplicar</text>
<text x="560" y="586" font-size="15" font-family="monospace" fill="#cfe6ff">EL EXPONENTE  s ↦ 1−s</text>
<text x="1100" y="586" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">½</text>
<text x="750" y="626" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">ξ(s) = ξ(1−s) dice exactamente «la máquina lee lo mismo de los dos lados de la coma» · medido %.1e</text>
`, peorEsp)

	fmt.Fprintf(&b, `<rect x="40" y="670" width="1420" height="150" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="750" y="702" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">Y ESTO EXPLICA EL VOLUMEN DE F247</text>
<text x="750" y="736" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">cada perla canta con amplitud x^½ = √x. ¿Por qué esa potencia? Porque √x·√x = x:</text>
<text x="750" y="758" font-size="15" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">la raíz cuadrada está a la misma distancia multiplicativa del 1 que del x.</text>
<text x="750" y="792" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd98a">TODAS LAS NOTAS CANTAN EN EL EQUILIBRIO DE LA COMA</text>

<rect x="40" y="840" width="1420" height="130" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="750" y="872" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffb27a">📌 ACÁ CORRIJO MI PROPIA RESPUESTA</text>
<text x="750" y="900" font-size="14" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">le contesté al capitán que la coma marca 10⁰ = 1, no ½. Eso contestaba OTRA pregunta.</text>
<text x="750" y="922" font-size="14" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Él no preguntaba qué posición NOMBRA la coma: preguntaba DÓNDE SE EQUILIBRAN SUS DOS LADOS. Y ahí es ½.</text>
<text x="750" y="952" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ Esto explica POR QUÉ el número es ½. NO muestra que los ceros estén ahí. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("la-coma.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: la-coma.svg")
}
