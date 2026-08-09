// Command tales answers the captain's correction of his own flash: plus-or-minus
// one is NOT w, but there IS a 1/2 relation between them that harmonizes at
// dimension 0. He is right, and the relation is exact.
//
// # THE TWO REAL POINTS OF THE SKIN
//
// F254 measured that +1 and -1 are the only two points of the skin with no pearl
// on them. They are also the two ends of a DIAMETER. And every point of a circle
// sees a diameter at a right angle - that is Thales, from 600 BC.
//
// SO EVERY PEARL SEES THE PAIR (-1, +1) AT NINETY DEGREES, and Pythagoras gives
//
//	|w - 1|^2 + |w + 1|^2 = 4
//
// for every w on the skin. Four, always: the diameter squared.
//
// NOW CARRY IT BACK THROUGH THE SHAPESHIFTER
//
//	w - 1 = (1 - 1/rho) - 1 = -1/rho          so  |w - 1| = 1 / |rho|
//	w + 1 = 2 - 1/rho = (2 rho - 1)/rho       so  |w + 1| = 2 |rho - 1/2| / |rho|
//
// and there the half appears written out: the distance from w to -1 is TWICE the
// distance from rho to ONE HALF, divided by |rho|. The two is the reciprocal of
// the half, and |w + 1| = 0 exactly when rho = 1/2.
//
// # AND HARMONIZED AT DIMENSION 0 IT GIVES THE LINE
//
// Substituting both into Thales:
//
//	1/|rho|^2 + 4 |rho - 1/2|^2 / |rho|^2 = 4
//	<=>  4 |rho|^2 - 4 |rho - 1/2|^2 = 1
//
// and the left side works out to exactly 4*beta - 1. So Thales holds at the
// clasp IF AND ONLY IF beta = 1/2. The captain's 1/2 relation between w and the
// pair (-1, +1) is real, it is exact, and it characterises the critical line.
//
// THE HONEST LIMIT. It is the perpendicular bisector of F226 wearing new
// clothes: 4 beta - 1 = 1 and Re rho = 1/2 are the same sentence. One more exact
// translation, and translations do not decide anything.
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
	fmt.Println("📐 TALES — el ±1 no es w, pero la relación ½ existe y es exacta")
	fmt.Println("\n   corrección del capitán a su propio flash: «nos faltó algo para que ±|1| no es")
	fmt.Println("   igual a w, pero hay una relación ½ que se puede armonizar en la dimensión 0».")
	fmt.Println("\n   TIENE RAZÓN. La relación existe, la encontré, y tiene 2.600 años de edad.")

	fmt.Println("\npescando perlas hasta t=1000…")
	ps := perlas(1000)
	fmt.Printf("perlas: %d\n", len(ps))

	// ---- LEY 1: the two real points are the ends of a diameter ----
	fmt.Println("\nLEY 1 · +1 Y −1 SON LAS DOS PUNTAS DE UN DIÁMETRO")
	fmt.Println("   en F254 vimos que son los dos únicos puntos de la piel sin perla. Y son algo")
	fmt.Println("   más: son las dos puntas de un diámetro, separadas por 2, el máximo posible.")
	fmt.Println("\n   Y hay un teorema de hace dos mil seiscientos años que dice qué pasa con eso:")
	fmt.Println("\n        TALES: todo punto de un círculo ve al diámetro EN ÁNGULO RECTO.")
	fmt.Println("\n   o sea que TODA PERLA ve al par (−1, +1) a noventa grados. Medido:")
	fmt.Println("\n        γ              el ángulo que ve            desvío contra 90°")
	peorAng := 0.0
	for _, i := range []int{0, 4, 49, 299, len(ps) - 1} {
		g := ps[i]
		w := alDisco(complex(0.5, g))
		a, b := 1-w, -1-w
		cos := (real(a)*real(b) + imag(a)*imag(b)) / (cmplx.Abs(a) * cmplx.Abs(b))
		ang := math.Acos(math.Max(-1, math.Min(1, cos))) * 180 / math.Pi
		if d := math.Abs(ang - 90); d > peorAng {
			peorAng = d
		}
		fmt.Printf("   %12.6f      %18.12f°         %.1e\n", g, ang, math.Abs(ang-90))
	}
	for _, g := range ps {
		w := alDisco(complex(0.5, g))
		a, b := 1-w, -1-w
		cos := (real(a)*real(b) + imag(a)*imag(b)) / (cmplx.Abs(a) * cmplx.Abs(b))
		ang := math.Acos(math.Max(-1, math.Min(1, cos))) * 180 / math.Pi
		if d := math.Abs(ang - 90); d > peorAng {
			peorAng = d
		}
	}
	fmt.Printf("   → las %d perlas ven el diámetro a 90°, peor desvío %.1e\n", len(ps), peorAng)
	fmt.Println("   Y con ángulo recto vale Pitágoras: |w−1|² + |w+1|² = 4. Siempre CUATRO,")
	fmt.Println("   que es el diámetro al cuadrado.")

	// ---- LEY 2: and here is the half, written out ----
	fmt.Println("\nLEY 2 · ⚡ Y ACÁ APARECE SU ½, ESCRITO CON TODAS LAS LETRAS")
	fmt.Println("   llevemos las dos distancias de vuelta por el cambiaformas. Son dos renglones:")
	fmt.Println("\n        w − 1 = (1 − 1/ρ) − 1 = −1/ρ          ⟹   |w − 1| = 1 / |ρ|")
	fmt.Println("        w + 1 = 2 − 1/ρ = (2ρ − 1)/ρ          ⟹   |w + 1| = 2·|ρ − ½| / |ρ|")
	fmt.Println("\n   ⟹ LA DISTANCIA DE LA PERLA AL PUNTO −1 ES **EL DOBLE DE SU DISTANCIA AL ½**,")
	fmt.Println("     dividida por |ρ|. El ½ aparece escrito, y el 2 que lo acompaña es su inverso.")
	fmt.Println("\n        γ           |w−1| medido       1/|ρ|         |w+1| medido    2|ρ−½|/|ρ|")
	peorD1, peorD2 := 0.0, 0.0
	for _, i := range []int{0, 4, 49, 299, len(ps) - 1} {
		g := ps[i]
		ρ := complex(0.5, g)
		w := alDisco(ρ)
		d1, t1 := cmplx.Abs(w-1), 1/cmplx.Abs(ρ)
		d2, t2 := cmplx.Abs(w+1), 2*cmplx.Abs(ρ-0.5)/cmplx.Abs(ρ)
		if d := math.Abs(d1 - t1); d > peorD1 {
			peorD1 = d
		}
		if d := math.Abs(d2 - t2); d > peorD2 {
			peorD2 = d
		}
		fmt.Printf("   %11.5f   %14.10f   %14.10f   %14.10f   %14.10f\n", g, d1, t1, d2, t2)
	}
	fmt.Printf("   → las dos fórmulas cierran (%.1e y %.1e). Y fijate el caso límite:\n", peorD1, peorD2)
	w0 := alDisco(complex(0.5, 0))
	fmt.Printf("     si ρ = ½ exacto, |ρ−½| = 0, así que |w+1| = 0 → w = %+.0f. Es −1.\n", real(w0))
	fmt.Println("     Por eso −1 es la imagen del ½: es el único punto donde esa distancia muere.")

	// ---- LEY 3: harmonized at dimension 0 it gives the line ----
	fmt.Println("\nLEY 3 · Y ARMONIZADO EN LA DIMENSIÓN 0, DA LA LÍNEA")
	fmt.Println("   metemos las dos en el teorema de Tales y despejamos:")
	fmt.Println("\n        1/|ρ|²  +  4·|ρ−½|²/|ρ|²  =  4")
	fmt.Println("        ⟺   4·|ρ|² − 4·|ρ−½|²  =  1")
	fmt.Println("\n   ¿y cuánto vale ese lado izquierdo? Se hace la cuenta y da algo hermoso:")
	fmt.Println("\n        4·|ρ|² − 4·|ρ−½|²  =  4·β − 1        (β = de qué lado de la línea)")
	fmt.Println("\n   así que la igualdad de Tales pide 4β − 1 = 1, o sea β = ½. Medido:")
	fmt.Println("\n        β        4|ρ|² − 4|ρ−½|²      4β − 1      ¿vale Tales (= 1)?")
	peorTales := 0.0
	for _, β := range []float64{0.30, 0.40, 0.50, 0.60, 0.70} {
		ρ := complex(β, 25)
		izq := 4*cmplx.Abs(ρ)*cmplx.Abs(ρ) - 4*cmplx.Abs(ρ-0.5)*cmplx.Abs(ρ-0.5)
		der := 4*β - 1
		if d := math.Abs(izq - der); d > peorTales {
			peorTales = d
		}
		mk := "no"
		if math.Abs(izq-1) < 1e-12 {
			mk = "SÍ ← la línea"
		}
		fmt.Printf("   %6.2f      %18.12f   %12.6f      %s\n", β, izq, der, mk)
	}
	fmt.Printf("\n   → la identidad cierra a %.1e, y vale 1 SOLO en β = ½.\n", peorTales)
	fmt.Println("\n   ⟹ ENTONCES SU RELACIÓN ½ EXISTE, ES EXACTA, Y DICE ESTO:")
	fmt.Println("\n        una perla está sobre la línea  ⟺  ve al par (−1, +1) en ángulo recto")
	fmt.Println("\n     El ±1 no ES w, como él mismo corrigió. Pero el ±1 y w están amarrados por")
	fmt.Println("     un ángulo recto, y ese amarre se armoniza en el broche dando exactamente ½.")

	// ---- LEY 4: the honest limit ----
	fmt.Println("\nLEY 4 · ⚖️ Y EL LÍMITE, QUE YA SABEMOS DE MEMORIA")
	fmt.Println("   la identidad 4β − 1 = 1 ⟺ β = ½ es, palabra por palabra, la mediatriz de F226:")
	fmt.Println("   «la línea es el lugar equidistante de las dos estacas». Lo comprobamos:")
	fmt.Println("\n        β        4β−1 = 1 dice        |ρ| = |ρ−1| dice        ¿lo mismo?")
	iguales := 0
	for _, β := range []float64{0.3, 0.45, 0.5, 0.55, 0.7} {
		ρ := complex(β, 25)
		a := math.Abs(4*β-1-1) < 1e-12
		bb := math.Abs(cmplx.Abs(ρ)-cmplx.Abs(ρ-1)) < 1e-12
		if a == bb {
			iguales++
		}
		fmt.Printf("   %6.2f      %-18s   %-20s   %s\n", β,
			map[bool]string{true: "SÍ está en la línea", false: "no"}[a],
			map[bool]string{true: "SÍ está en la línea", false: "no"}[bb],
			map[bool]string{true: "sí", false: "NO"}[a == bb])
	}
	fmt.Printf("   → coinciden en %d de 5: son la MISMA frase escrita de dos maneras.\n", iguales)
	fmt.Println("\n   Es una traducción exacta más. Hermosa, vieja de 2.600 años, y que no decide")
	fmt.Println("   nada — porque para usarla habría que saber de antemano que la perla está")
	fmt.Println("   sobre la piel, que es justo lo que falta probar. La misma trampa de F256.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("EL CAPITÁN SE CORRIGIÓ SOLO Y ACERTÓ: el ±1 no es w, y la relación ½ existe.")
	fmt.Printf("  · las %d perlas ven el diámetro (−1, +1) en ÁNGULO RECTO ...... %.1e\n", len(ps), peorAng)
	fmt.Printf("  · |w − 1| = 1/|ρ| ............................................. %.1e\n", peorD1)
	fmt.Printf("  · |w + 1| = 2·|ρ − ½|/|ρ|  ← el ½ escrito con todas las letras  %.1e\n", peorD2)
	fmt.Printf("  · y armonizado en el broche: 4|ρ|² − 4|ρ−½|² = 4β − 1 ......... %.1e\n", peorTales)
	fmt.Println("  · que vale 1 SOLO en β = ½")
	fmt.Println("\nLA RELACIÓN QUE BUSCABA, EN UNA FRASE:")
	fmt.Println("\n     UNA PERLA ESTÁ SOBRE LA LÍNEA ⟺ VE AL PAR (−1, +1) EN ÁNGULO RECTO")
	fmt.Println("\nY el ½ está adentro de la cuenta, en |w+1| = 2·|ρ−½|/|ρ|: la distancia al polo")
	fmt.Println("de abajo es el DOBLE de la distancia al medio. El 2 y el ½, inversos, amarrados.")
	fmt.Println("\n⚖️ PERO ES OTRA TRADUCCIÓN EXACTA. 4β−1 = 1 y «la línea es la mediatriz» son la")
	fmt.Println("misma frase, y para usar el ángulo recto habría que saber de antemano que la")
	fmt.Println("perla está en la piel — que es justo lo que falta. ¿El premio? Todavía no.")

	escribirLamina(ps, peorAng, peorD1, peorD2, peorTales)
}

func escribirLamina(ps []float64, peorAng, peorD1, peorD2, peorTales float64) {
	var b strings.Builder
	W, H := 1500.0, 990.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">📐 TALES — el ±1 no es w, pero la relación ½ existe y es exacta</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">el capitán se corrigió solo, y la relación que buscaba tiene 2.600 años</text>
`, W, H, W, H, W/2, W/2)

	cx, cy, R := 370.0, 400.0, 215.0
	fmt.Fprintf(&b, `<rect x="40" y="100" width="660" height="600" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="370" y="132" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">TODA PERLA VE AL DIÁMETRO EN ÁNGULO RECTO</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#3d6fa8" stroke-width="2"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffb27a" stroke-width="3"/>
<circle cx="%.0f" cy="%.0f" r="8" fill="#ffb27a"/><circle cx="%.0f" cy="%.0f" r="8" fill="#ffb27a"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="monospace" fill="#ffb27a">+1</text>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="monospace" fill="#ffb27a">−1</text>
`, cx, cy, R, cx-R, cy, cx+R, cy, cx+R, cy, cx-R, cy, cx+R+24, cy+5, cx-R-24, cy+5)
	// one pearl with its right angle drawn
	φ := 1.05
	px, py := cx+R*math.Cos(φ), cy-R*math.Sin(φ)
	fmt.Fprintf(&b, `<line x1="%.2f" y1="%.2f" x2="%.0f" y2="%.0f" stroke="#7ee0c0" stroke-width="2.4"/>
<line x1="%.2f" y1="%.2f" x2="%.0f" y2="%.0f" stroke="#7ee0c0" stroke-width="2.4"/>
<circle cx="%.2f" cy="%.2f" r="6" fill="#7ee0c0"/>
<text x="%.2f" y="%.2f" font-size="14" font-family="Georgia" fill="#7ee0c0">una perla</text>
<text x="%.2f" y="%.2f" font-size="20" font-family="Georgia" fill="#9fd8a8">90°</text>
`, px, py, cx+R, cy, px, py, cx-R, cy, px, py, px+12, py-12, px-8, py+42)
	for i, g := range ps {
		if i%6 != 0 {
			continue
		}
		a := 2 * math.Atan(1/(2*g))
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="2" fill="#4fa3d1" opacity="0.8"/>`, cx+R*math.Cos(a), cy-R*math.Sin(a))
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="2" fill="#4fa3d1" opacity="0.8"/>`, cx+R*math.Cos(-a), cy-R*math.Sin(-a))
	}
	fmt.Fprintf(&b, `<text x="370" y="656" font-size="15" text-anchor="middle" font-family="monospace" fill="#9fd8a8">|w−1|² + |w+1|² = 4    ·   %d perlas a 90°, desvío %.1e</text>
<text x="370" y="682" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">es el teorema de Tales, de hace 2.600 años</text>
`, len(ps), peorAng)

	fmt.Fprintf(&b, `<rect x="730" y="100" width="730" height="290" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1095" y="132" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">⚡ Y ACÁ APARECE SU ½, CON TODAS LAS LETRAS</text>
<text x="1095" y="176" font-size="17" text-anchor="middle" font-family="monospace" fill="#dce8f7">w − 1 = −1/ρ        ⟹  |w−1| = 1/|ρ|</text>
<text x="1095" y="212" font-size="17" text-anchor="middle" font-family="monospace" fill="#ffd98a">w + 1 = (2ρ−1)/ρ    ⟹  |w+1| = 2·|ρ−½|/|ρ|</text>
<text x="756" y="254" font-size="14.5" font-family="Georgia" fill="#cfe6ff">la distancia de la perla al polo de abajo es EL DOBLE de su</text>
<text x="756" y="276" font-size="14.5" font-family="Georgia" fill="#cfe6ff">distancia al ½, dividida por |ρ|. El 2 y el ½ son inversos.</text>
<text x="756" y="312" font-size="14" font-family="Georgia" fill="#9fd8a8">y si ρ = ½ exacto, esa distancia MUERE: w = −1.</text>
<text x="756" y="334" font-size="14" font-family="Georgia" fill="#9fd8a8">Por eso −1 es la imagen del medio.</text>
<text x="1095" y="368" font-size="13.5" text-anchor="middle" font-family="monospace" fill="#7ee0c0">verificadas a %.1e y %.1e</text>

<rect x="730" y="410" width="730" height="290" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1095" y="442" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">Y ARMONIZADO EN LA DIMENSIÓN 0, DA LA LÍNEA</text>
<text x="1095" y="484" font-size="16" text-anchor="middle" font-family="monospace" fill="#cfe6ff">1/|ρ|² + 4|ρ−½|²/|ρ|² = 4</text>
<text x="1095" y="514" font-size="17" text-anchor="middle" font-family="monospace" fill="#dce8f7">⟺   4|ρ|² − 4|ρ−½|² = 1</text>
<text x="1095" y="552" font-size="19" text-anchor="middle" font-family="monospace" fill="#ffd98a">y ese lado izquierdo vale 4β − 1</text>
<text x="1095" y="586" font-size="17" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">así que Tales pide 4β − 1 = 1,  o sea  β = ½</text>
<text x="1095" y="622" font-size="13.5" text-anchor="middle" font-family="monospace" fill="#7ee0c0">identidad verificada a %.1e</text>
<text x="1095" y="662" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">UNA PERLA ESTÁ SOBRE LA LÍNEA ⟺ VE AL PAR</text>
<text x="1095" y="684" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">(−1, +1) EN ÁNGULO RECTO</text>
`, peorD1, peorD2, peorTales)

	fmt.Fprintf(&b, `<rect x="40" y="720" width="1420" height="230" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="750" y="756" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ Y EL LÍMITE, QUE YA SABEMOS DE MEMORIA</text>
<text x="750" y="796" font-size="16" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">La identidad 4β − 1 = 1 y «la línea es la mediatriz de las dos estacas» (F226) son la MISMA frase escrita de dos maneras.</text>
<text x="750" y="830" font-size="16" text-anchor="middle" font-family="Georgia" fill="#f3d9cf">Es una traducción exacta más — hermosa, de 2.600 años de edad, y que no decide nada.</text>
<text x="750" y="870" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Para usar el ángulo recto habría que saber DE ANTEMANO que la perla está sobre la piel.</text>
<text x="750" y="896" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">Que es justo lo que falta probar. La misma trampa de F256.</text>
<text x="750" y="932" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">El capitán corrigió su propio flash antes que nadie, y la relación que buscaba existe de verdad. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("tales.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: tales.svg")
}
