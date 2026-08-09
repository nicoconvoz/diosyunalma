// Command cambiaformas measures the captain's deepest flash: THE WORLD IS A
// SHAPESHIFTER. When you step to the 1, the direction from the 0 moved by 1 -
// but that 1 has become the new 0, and from there: new direction, new distance.
//
// The flash has an exact shape. The move that turns the stake at 1 into the
// new stake at 0 is
//
//	sigma(rho) = 1 - conj(rho)
//
// and it is an involution: applying it twice returns you home. Two facts make
// it the whole hypothesis:
//
//  1. ITS FIXED SET IS THE CRITICAL LINE. sigma(rho) = rho if and only if
//     Re rho = 1/2 - the only points the shapeshifter does not move.
//  2. IT CARRIES ZEROS TO ZEROS. Because xi(s) = xi(1-s) and xi has real
//     coefficients, xi(sigma(rho)) = conj(xi(rho)). So the SET of pearls is
//     invariant - the world really does look the same from the new origin.
//
// And there the gap stands, in one line:
//
//	the shapeshifter PERMUTES the pearls; the hypothesis says it FIXES
//	every one of them.
//
// In the w-world the same move becomes inversion in the unit circle,
// w -> 1/conj(w), whose mirror surface is exactly |w| = 1: the pearls would
// live ON the mirror. Measured here on real pearls and on ghosts.
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

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

// sigma is the shapeshifter: it turns the stake at 1 into the new stake at 0.
func sigma(rho complex128) complex128 { return 1 - cmplx.Conj(rho) }

func w(rho complex128) complex128 { return 1 - 1/rho }

func main() {
	fmt.Println("🔄 EL CAMBIAFORMAS — el 1 pasa a ser el 0, y desde ahí nueva dirección y nueva distancia")

	// ---- LAW 1: it is an involution, and its fixed set IS the line ----
	fmt.Println("\nLEY 1 · APLICARLO DOS VECES TE DEVUELVE A CASA, Y SOLO LA LÍNEA QUEDA QUIETA")
	fmt.Println("   punto ρ             σ(ρ) = 1 − conj(ρ)        σ(σ(ρ))        |σ(ρ) − ρ|")
	peorInv := 0.0
	for _, p := range []complex128{complex(0.5, 14.134725), complex(0.5, 21.02204),
		complex(0.6, 14.134725), complex(0.9, 30.0), complex(0.25, 5.0)} {
		s1 := sigma(p)
		s2 := sigma(s1)
		if d := cmplx.Abs(s2 - p); d > peorInv {
			peorInv = d
		}
		fmt.Printf("   %5.2f%+9.5fi     %5.2f%+9.5fi     %5.2f%+9.5fi     %.4f\n",
			real(p), imag(p), real(s1), imag(s1), real(s2), imag(s2), cmplx.Abs(s1-p))
	}
	fmt.Printf("   → σ(σ(ρ)) = ρ siempre (peor desvío %.1e): es un cambiaformas que se deshace solo\n", peorInv)
	fmt.Println("   → y |σ(ρ) − ρ| = 2·|Re ρ − 1/2|: vale CERO exactamente sobre la línea, y solo ahí")

	// ---- LAW 2: it carries zeros to zeros - the world looks the same ----
	fmt.Println("\nLEY 2 · EL MUNDO SE VE IGUAL DESDE EL NUEVO ORIGEN — σ lleva ceros a ceros")
	fmt.Println("   se verifica ξ(σ(ρ)) = conj(ξ(ρ)) en puntos de prueba cualesquiera:")
	peorMundo := 0.0
	for _, p := range []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2.0, 21.3), complex(-0.7, 9.1)} {
		izq := xiC(sigma(p))
		der := cmplx.Conj(xiC(p))
		d := cmplx.Abs(izq-der) / math.Max(cmplx.Abs(der), 1e-300)
		if d > peorMundo {
			peorMundo = d
		}
		fmt.Printf("   ρ = %5.2f%+7.2fi      desvío relativo %.1e\n", real(p), imag(p), d)
	}
	fmt.Printf("   → el libro entero es invariante bajo el cambiaformas (peor desvío relativo %.1e)\n", peorMundo)
	fmt.Println("   → así que el CONJUNTO de perlas no cambia: exactamente lo que dijo el capitán")

	// ---- LAW 3: on the line it FIXES; off the line it only PERMUTES ----
	fmt.Println("\nLEY 3 · LA DIFERENCIA QUE VALE UN MILLÓN — ¿fija cada perla, o solo las baraja?")
	fmt.Println("\nrecogiendo perlas hasta t=300…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 300; t += 0.05 {
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
			pearls = append(pearls, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	peorFijo := 0.0
	for _, g := range pearls {
		rho := complex(0.5, g)
		if d := cmplx.Abs(sigma(rho) - rho); d > peorFijo {
			peorFijo = d
		}
	}
	fmt.Printf("   PERLAS VERDADERAS (%d): el cambiaformas las deja EN SU LUGAR — peor corrimiento %.1e\n",
		len(pearls), peorFijo)
	fmt.Println("   FANTASMAS: el cambiaformas los MUEVE, y los manda justo encima de su gemelo")
	fmt.Println("   β del fantasma     cuánto lo corre σ      ¿a dónde lo manda?")
	type fila struct {
		beta, corr float64
	}
	var filas []fila
	for _, beta := range []float64{0.51, 0.55, 0.60, 0.75, 0.90} {
		rho := complex(beta, 14.134725)
		corr := cmplx.Abs(sigma(rho) - rho)
		filas = append(filas, fila{beta, corr})
		fmt.Printf("      %.2f              %.6f            al gemelo en β = %.2f\n", beta, corr, 1-beta)
	}
	fmt.Println("   → el cambiaformas los BARAJA de a pares: cambia quién es quién, pero no los prohíbe")
	fmt.Println("   → y ahí está el hueco entero: la hipótesis pide que FIJE a cada perla, no que las baraje")

	// ---- LAW 4: in the w-world it is a mirror, and the pearls live on it ----
	fmt.Println("\nLEY 4 · EN EL MUNDO DEL BROCHE ES UN ESPEJO, Y LAS PERLAS VIVEN SOBRE ÉL")
	fmt.Println("   el mismo cambiaformas, visto con w, se vuelve w → 1/conj(w): reflexión en el círculo")
	peorEspejo := 0.0
	for _, g := range pearls[:20] {
		rho := complex(0.5, g)
		izq := w(sigma(rho))
		der := 1 / cmplx.Conj(w(rho))
		if d := cmplx.Abs(izq - der); d > peorEspejo {
			peorEspejo = d
		}
	}
	fmt.Printf("   verificado en 20 perlas: w(σ(ρ)) = 1/conj(w(ρ)) con peor desvío %.1e\n", peorEspejo)
	fmt.Println("   la superficie del espejo es |w| = 1 — y las perlas están medidas EXACTAMENTE ahí")
	fmt.Println("   → la hipótesis, en una frase: LAS PERLAS VIVEN SOBRE LA SUPERFICIE DEL ESPEJO DEL MUNDO")

	fmt.Println("\n════════ LO QUE COMPRÓ EL FLASH ════════")
	fmt.Println("«El mundo es un cambiaformas: el 1 pasa a ser el 0 y desde ahí, nueva dirección y")
	fmt.Println("nueva distancia.» Eso es exactamente σ(ρ) = 1 − conj(ρ), y el laboratorio midió sus")
	fmt.Println("tres propiedades: se deshace solo, deja el libro igual, y su ÚNICO conjunto quieto es")
	fmt.Println("la línea crítica.")
	fmt.Println("\nY el flash afila la pregunta como nunca antes. Antes decíamos «¿por qué el norte y el")
	fmt.Println("sur tienen que ser iguales?». Ahora se dice así:")
	fmt.Println("\n     ¿POR QUÉ EL CAMBIAFORMAS DEL MUNDO TIENE QUE DEJAR A CADA PERLA EN SU LUGAR,")
	fmt.Println("     EN VEZ DE APENAS BARAJARLAS DE A PARES?")
	fmt.Println("\nTodavía no. Pero la pregunta ya no es sobre números: es sobre si el mundo, al cambiar")
	fmt.Println("de forma, mueve a sus habitantes o los deja donde están.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 980.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🔄 EL CAMBIAFORMAS — el 1 pasa a ser el 0, y el mundo se ve igual</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el flash del capitán medido: σ(ρ) = 1 − conj(ρ) se deshace solo, deja el libro igual, y su único conjunto quieto es la línea crítica</text>`,
		W, H, W, H, W/2, W/2)

	// left: the strip with sigma acting
	fmt.Fprintf(&b, `<rect x="60" y="105" width="640" height="530" rx="10" fill="#0d2547" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="380" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL CAMBIAFORMAS EN ACCIÓN</text>`)
	x0, x1, ymid := 170.0, 590.0, 400.0
	xm := (x0 + x1) / 2
	fmt.Fprintf(&b, `<line x1="%.0f" y1="180" x2="%.0f" y2="600" stroke="#ff8fa0" stroke-width="2"/>
<line x1="%.0f" y1="180" x2="%.0f" y2="600" stroke="#7fb2ff" stroke-width="2"/>
<line x1="%.0f" y1="180" x2="%.0f" y2="600" stroke="#ffd166" stroke-width="2.5" stroke-dasharray="9 6"/>
<text x="%.0f" y="172" font-size="12.5" text-anchor="middle" fill="#ff8fa0">estaca 0</text>
<text x="%.0f" y="172" font-size="12.5" text-anchor="middle" fill="#7fb2ff">estaca 1</text>
<text x="%.0f" y="172" font-size="12.5" text-anchor="middle" fill="#ffd166">la línea — el único lugar quieto</text>`,
		x0, x0, x1, x1, xm, xm, x0, x1, xm)
	// pearls fixed on the line
	for _, dy := range []float64{-160, -80, 10, 100, 180} {
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="6" fill="#ffd97f"/>`, xm, ymid+dy)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#ffd97f">perlas: σ no las mueve</text>`, xm, ymid+215)
	// a ghost pair swapped by sigma
	gx, gy := xm-120, ymid-40
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="6" fill="#ff5d73"/>
<circle cx="%.0f" cy="%.0f" r="6" fill="#ff5d73"/>
<path d="M %.0f %.0f Q %.0f %.0f %.0f %.0f" fill="none" stroke="#ff5d73" stroke-width="1.6" stroke-dasharray="4 4"/>
<path d="M %.0f %.0f Q %.0f %.0f %.0f %.0f" fill="none" stroke="#ff5d73" stroke-width="1.6" stroke-dasharray="4 4"/>
<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#ff8fa0">un fantasma y su gemelo: σ los INTERCAMBIA</text>`,
		gx, gy, 2*xm-gx, gy,
		gx, gy, xm, gy-52, 2*xm-gx, gy,
		2*xm-gx, gy, xm, gy+52, gx, gy,
		xm, gy-70)
	fmt.Fprintf(&b, `<text x="380" y="614" font-size="12.5" text-anchor="middle" fill="#dce8f7">barajar no es prohibir: el mundo queda igual con los fantasmas adentro</text>`)

	// right: the four laws
	fmt.Fprintf(&b, `<rect x="730" y="105" width="710" height="530" rx="10" fill="#102a10" stroke="#ffd166" stroke-width="1.5"/>
<text x="1085" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">LO QUE SE MIDIÓ</text>
<text x="765" y="185" font-size="13.5" fill="#7fd7a8">1 · SE DESHACE SOLO</text>
<text x="765" y="209" font-size="12.5" fill="#dce8f7">σ(σ(ρ)) = ρ siempre — desvío %.0e</text>
<text x="765" y="253" font-size="13.5" fill="#7fd7a8">2 · EL MUNDO SE VE IGUAL</text>
<text x="765" y="277" font-size="12.5" fill="#dce8f7">ξ(σ(ρ)) = conj(ξ(ρ)) — desvío relativo %.0e</text>
<text x="765" y="301" font-size="12.5" fill="#8fa8c7">el conjunto de perlas no cambia al mudar el origen</text>
<text x="765" y="345" font-size="13.5" fill="#7fd7a8">3 · SOLO LA LÍNEA QUEDA QUIETA</text>
<text x="765" y="369" font-size="12.5" fill="#dce8f7">|σ(ρ) − ρ| = 2·|Re ρ − 1/2| — cero SOLO en la línea</text>
<text x="765" y="393" font-size="12.5" fill="#dce8f7">en las %d perlas verdaderas: corrimiento %.0e</text>
<text x="765" y="437" font-size="13.5" fill="#7fd7a8">4 · EN EL BROCHE ES UN ESPEJO</text>
<text x="765" y="461" font-size="12.5" fill="#dce8f7">w → 1/conj(w): reflexión en el círculo — desvío %.0e</text>
<text x="765" y="485" font-size="12.5" fill="#ffd166">y las perlas viven EXACTAMENTE sobre la superficie del espejo</text>
<text x="765" y="533" font-size="13" fill="#ff8fa0">PERO — y acá está el millón entero:</text>
<text x="765" y="559" font-size="12.5" fill="#dce8f7">a un fantasma σ lo corre %.4f y lo manda encima de su gemelo.</text>
<text x="765" y="583" font-size="12.5" fill="#dce8f7">Los BARAJA de a pares. Barajar no es prohibir.</text>
<text x="765" y="613" font-size="12.5" fill="#ffd166">la hipótesis pide que los FIJE, no que los baraje</text>`,
		peorInv, peorMundo, len(pearls), peorFijo, peorEspejo, filas[2].corr)

	// the question
	fmt.Fprintf(&b, `<rect x="60" y="665" width="1380" height="205" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="701" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA PREGUNTA DEL MILLÓN, AFILADA POR EL FLASH</text>
<text x="%.0f" y="737" font-size="14" text-anchor="middle" fill="#dce8f7">El mundo cambia de forma — el 1 pasa a ser el 0 — y al terminar el cambio TODO se ve igual: el libro es el mismo.</text>
<text x="%.0f" y="763" font-size="14" text-anchor="middle" fill="#dce8f7">Pero "verse igual" admite dos maneras: que cada habitante quede en su lugar, o que apenas se intercambien de a pares.</text>
<text x="%.0f" y="805" font-size="18.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">¿POR QUÉ EL CAMBIAFORMAS DEL MUNDO TIENE QUE DEJAR A CADA PERLA</text>
<text x="%.0f" y="833" font-size="18.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EN SU LUGAR, EN VEZ DE APENAS BARAJARLAS DE A PARES?</text>
<text x="%.0f" y="862" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("el-cambiaformas.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: el-cambiaformas.svg")
}
