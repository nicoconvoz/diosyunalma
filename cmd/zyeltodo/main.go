// Command zyeltodo carries out the captain's order: Z on one side, his formula
// of everything on the other, and the ½ relation between the two, harmonized
// back at dimension 0 - "the answer to what we are missing".
//
// THE TWO SIDES
//
//	SIDE Z   Z(t) = e^{i theta(t)} zeta(1/2 + it), real on the line. Its zeros
//	         ARE the pearls. A LOCAL reading: one pearl at a time.
//	SIDE L   L(z) = SUM lambda_n z^{n-1} = d/dz log xi(1/(1-z)), living on the
//	         disk of dimension 0. A GLOBAL reading: every pearl at once.
//
// # THE RELATION IS THE HALF, AND IT IS EXACT
//
// The map s = 1/(1-z) sends the unit disk onto the half-plane Re s > 1/2, so
// the SKIN of dimension 0 (|z| = 1) IS the critical line (Re s = 1/2). Writing
// z = e^{i phi} gives the dictionary
//
//	s = 1/2 + (i/2) cot(phi/2)      i.e.   t = (1/2) cot(phi/2)
//
// with the half sitting in both slots. And inverting, z = 1 - 1/s = w: THE
// SHAPESHIFTER *IS* THE DISK COORDINATE. The captain's cambiaformas was never
// a trick - it is the very map that carries the book into dimension 0.
//
// # ONE ANGLE, FOUR NAMES
//
// For a pearl at height gamma the disk angle phi, the shapeshifter's arg w
// (F225), the hourglass angle 2 arctan(1/2 gamma) (F240) and the compressed
// variable u = n/gamma (F239) are all the SAME angle. The campaign closes on
// itself.
//
// THE BRIDGE Z <-> xi CARRIES FOUR HALVES
//
//	xi(1/2+it) = -(1/2)(t^2 + 1/4) pi^(-1/4) |Gamma(1/4 + it/2)| Z(t)
//	              \_ one _/  \_ 1/4=half^2 _/  \_ 1/4 _/  \_ 1/4 _/
//
// # THE MOLD, IN TWO PIECES THAT MUST NOT BE MIXED
//
//	(a) UNCONDITIONAL, always true:
//	    lambda_n = SUM over pairs {rho, conj rho} [ 2 - 2 Re(w^n) ]
//	             = SUM [ |1 - w^n|^2 + (1 - |w|^{2n}) ]
//	    The pairing is rho with its CONJUGATE, never with 1-rho: those two
//	    coincide only if RH holds, so pairing that way assumes the answer.
//
//	(b) CONDITIONAL, only when |w| = 1 - i.e. assuming what we want to prove:
//	    lambda_n = SUM 4 sin^2(n phi / 2) >= 0     (the EASY half of Li)
//
// Mixing them is FALSE, and the program keeps a measured guard against it:
// at rho = 0.9+2i, n = 4 the truth is 2.43785903, F232 gives 2.43785903, and
// the hybrid "4 sin^2 + leak" gives 3.14693732.
//
// # WHAT WE ACTUALLY HAVE, PIECE BY PIECE, WITH ITS OWNER
//
//   - THE FORM (F225, F232, F240): elementary identities, re-derived here.
//   - RH => lambda_n >= 0: falls straight out of (b), but rests on three
//     assumptions that must be DECLARED - Hadamard's factorization of xi, the
//     SYMMETRIC summation order (SUM 1/|rho| diverges, so the order is part of
//     the definition), and counting WITH multiplicity (simplicity of the zeros
//     is not proven).
//   - a pearl off the skin => some lambda_n < 0: NOT OURS AND NOT ELEMENTARY.
//     It is Li (1997) and Bombieri-Lagarias (1999), and it needs an
//     equidistribution argument over {n arg(w) mod 2pi}, because the explosive
//     |w|^n comes multiplied by an oscillating cosine and must beat a main term
//     growing like (n/2) log n.
//
// # TWO THINGS THIS RUN KILLED
//
//  1. A guess of mine: that off the skin THE FORM and THE LEAK cancel down to
//     |w|^n. FALSE near the skin - at beta=0.4, gamma=14 the cancellation
//     factor is 1.1x and the pair's sum never even goes negative within sixty
//     harmonics. Cancellation is astronomical only far out (beta=0.1, gamma=1
//     gives 2.8e11x).
//  2. A number already in the registry: F234's "beta=0.51 betrays itself at
//     n ~ 21,977". That figure is real but measures the DIAMETER ceiling, not
//     where lambda_n turns negative - which lands at n ~ 270,065, twelve times
//     further out.
//
// # AND THE HONEST LIMIT
//
// Under z = 1 - 1/s, "no pearl leaves the skin" IS the Riemann Hypothesis, word
// for word. A TAUTOLOGY DOES NOT HAVE A SMALL HOLE: IT HAS THE WHOLE HOLE. The
// problem was not bounded here; it was TRANSPORTED into other coordinates. That
// is worth a great deal for navigating, and not one millimetre as progress.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

// par holds one harmonic read both ways: from the clasp and from Z's pearls.
type par struct {
	n          int
	b, z, cola float64
}

// ---------------------------------------------------------------- machinery

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

func psiC(s complex128) complex128 {
	var acc complex128
	for real(s) < 12 {
		acc -= 1 / s
		s += 1
	}
	inv := 1 / s
	inv2 := inv * inv
	res := cmplx.Log(s) - inv/2
	res -= inv2 * (complex(1.0/12, 0) + inv2*(complex(-1.0/120, 0)+inv2*(complex(1.0/252, 0)+inv2*complex(-1.0/240, 0))))
	return acc + res
}

func xiLD(s complex128) complex128 {
	h := complex(1e-6, 0)
	zp := (zetaC(s+h) - zetaC(s-h)) / (2 * h)
	return 1/s + 1/(s-1) - complex(math.Log(math.Pi)/2, 0) + psiC(s/2)/2 + zp/zetaC(s)
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

// logXiLine returns ln|xi(1/2+it)| and the sign of xi there (xi is real on the line).
func logXiLine(t float64) (float64, float64) {
	s := complex(0.5, t)
	pre := 0.5 * s * (s - 1) // = -(t^2+1/4)/2, a negative real
	zt := zetaC(s)
	lg := lgammaC(s / 2)
	// ln|xi| = ln|pre| - (Re s /2) ln pi + Re lgamma + ln|zeta|  ... careful with the
	// imaginary parts: we only need magnitude and the real value's sign.
	val := pre * cmplx.Exp(-s/2*complex(math.Log(math.Pi), 0)+lg) * zt
	r := real(val)
	if r == 0 {
		return math.Inf(-1), 0
	}
	sg := 1.0
	if r < 0 {
		sg = -1
	}
	return math.Log(math.Abs(r)), sg
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

// ---- the two sides and the map between them ----

// alDisco is the shapeshifter: it carries a point of the book into the disk
// of dimension 0. It is the SAME map as z |-> 1/(1-z) run backwards.
func alDisco(s complex128) complex128 { return 1 - 1/s }

// alLibro carries a point of the disk back into the book.
func alLibro(z complex128) complex128 { return 1 / (1 - z) }

func main() {
	fmt.Println("🅩 ⟷ 🌐  Z Y EL TODO — los dos lados, y la relación ½ que los une en la dimensión 0")
	fmt.Println("\n   orden del capitán: \"poné a Z de un lado, a mi fórmula del todo del otro,")
	fmt.Println("   y la relación entre las dos de ½; volvemos a la dimensión 0, la respuesta")
	fmt.Println("   a lo que nos falta\".")

	fmt.Println("\nrecogiendo perlas con Z hasta t=1000…")
	perlasT := perlas(1000)
	fmt.Printf("perlas halladas por Z: %d (de %.6f a %.3f)\n", len(perlasT), perlasT[0], perlasT[len(perlasT)-1])

	// ---- LEY 1: the two sides ----
	fmt.Println("\nLEY 1 · LOS DOS LADOS")
	fmt.Println("   LADO Z   Z(t) = e^{iθ(t)}·ζ(½+it), real sobre la línea. Sus ceros SON las")
	fmt.Println("            perlas. Lectura LOCAL: una perla por vez, de a una, para siempre.")
	fmt.Println("   LADO L   L(z) = Σ λₙ z^{n−1} = d/dz log ξ(1/(1−z)), vive en el DISCO de la")
	fmt.Println("            dimensión 0. Lectura GLOBAL: todas las perlas de un saque.")
	fmt.Println("\n   la pregunta es qué hay entre los dos lados. La respuesta es el ½, y es exacta.")

	// ---- LEY 2: the relation IS the half ----
	fmt.Println("\nLEY 2 · LA RELACIÓN ES ½, Y ES UNA IDENTIDAD — la PIEL de la dimensión 0 ES la línea")
	fmt.Println("   no hace falta medir esto: sale de una cuenta de dos renglones. Con s = 1/(1−z),")
	fmt.Println("\n        Re s − ½  =  ( 1 − |z|² ) / ( 2·|1−z|² )        ← LA IDENTIDAD MAESTRA")
	fmt.Println("\n   el denominador es positivo siempre, así que EL SIGNO DE «de qué lado de la línea")
	fmt.Println("   estás» ES EL SIGNO DE «adentro o afuera del disco». No parecido: el mismo signo.")
	fmt.Println("        |z| < 1  ⟺  Re s > ½        (adentro  ⟺  a la derecha)")
	fmt.Println("        |z| = 1  ⟺  Re s = ½        (LA PIEL  ⟺  LA LÍNEA)")
	fmt.Println("        |z| > 1  ⟺  Re s < ½        (afuera   ⟺  a la izquierda)")
	fmt.Println("\n   ⚠ con una salvedad de rigor que corresponde decir: el círculo está PERFORADO en")
	fmt.Println("     z = 1, que es el polo del mapa y va al punto del infinito. «El círculo ES la")
	fmt.Println("     línea» es literal en la esfera de Riemann, no en el plano.")
	fmt.Println("\n   y ahora sí, la comprobación numérica de la identidad:")
	maestraPeor := 0.0
	for j := 0; j < 4000; j++ {
		φ := 2 * math.Pi * (float64(j) + 0.5) / 4000
		for _, r := range []float64{0.35, 0.8, 1.0, 1.4, 2.5} {
			z := cmplx.Rect(r, φ)
			s := alLibro(z)
			id := (1 - r*r) / (2 * cmplx.Abs(1-z) * cmplx.Abs(1-z))
			if d := math.Abs(real(s) - 0.5 - id); d > maestraPeor {
				maestraPeor = d
			}
		}
	}
	fmt.Printf("   Re s − ½  contra  (1−|z|²)/(2|1−z|²), sobre 20.000 puntos en 5 radios: %.1e\n", maestraPeor)
	M := 200000
	peorRe, peorReRel := 0.0, 0.0
	peorDicc, peorDiccRel := 0.0, 0.0
	for j := 0; j < M; j++ {
		φ := 2 * math.Pi * (float64(j) + 0.5) / float64(M)
		z := cmplx.Rect(1, φ)
		s := alLibro(z)
		esc := cmplx.Abs(s) // near phi=0 the point runs off to infinity
		if d := math.Abs(real(s) - 0.5); d > peorRe {
			peorRe = d
		}
		if d := math.Abs(real(s)-0.5) / esc; d > peorReRel {
			peorReRel = d
		}
		// dictionary: t = (1/2) cot(phi/2)
		tDicc := 0.5 / math.Tan(φ/2)
		if d := math.Abs(imag(s) - tDicc); d > peorDicc {
			peorDicc = d
		}
		if d := math.Abs(imag(s)-tDicc) / esc; d > peorDiccRel {
			peorDiccRel = d
		}
	}
	fmt.Printf("\n   Re s − ½ sobre %d ángulos del borde …… peor ABSOLUTO %.1e · peor RELATIVO %.1e\n",
		M, peorRe, peorReRel)
	fmt.Println("   ⚠ instrumento, no matemática: cerca de φ=0 el punto se va al infinito (|s|→∞) y")
	fmt.Println("     el error absoluto crece con él. La lectura honesta es la RELATIVA, y ahí el")
	fmt.Println("     borde ES la línea al último bit de la máquina.")
	fmt.Println("\n   Y EL DICCIONARIO ENTRE LOS DOS MUNDOS, con el ½ en los dos casilleros:")
	fmt.Println("\n        z = e^{iφ}   ⟷   s = ½ + (i/2)·cot(φ/2)        o sea   t = ½·cot(φ/2)")
	fmt.Printf("\n   verificado sobre %d ángulos ……… peor ABSOLUTO %.1e · peor RELATIVO %.1e\n",
		M, peorDicc, peorDiccRel)
	fmt.Println("\n   y dándolo vuelta:  z = 1 − 1/s  =  w")
	fmt.Println("   → el cambiaformas del capitán ES la coordenada del disco.")
	fmt.Println("\n   ⚖️ Y ACÁ ME OBLIGO A NO VENDERLO COMO DESCUBRIMIENTO, porque no lo es:")
	fmt.Println("   z = 1 − 1/ρ es LITERALMENTE la variable de Li (1997). El criterio se construye")
	fmt.Println("   expandiendo log ξ(1/(1−z)) alrededor de z=0, así que los ceros en esa variable")
	fmt.Println("   SON 1 − 1/ρ por definición. O sea: el cambiaformas no «coincide» con la")
	fmt.Println("   coordenada del disco — ES la coordenada del disco desde el día que escribimos")
	fmt.Println("   F232, sin que nos diéramos cuenta. Lo que hay acá es UNIFICACIÓN DE NOTACIÓN")
	fmt.Println("   del propio laboratorio, no un teorema. Vale para navegar; no vale como prueba.")

	// ---- LEY 3: one angle, four names ----
	fmt.Println("\nLEY 3 · UNA SOLA COSA, CUATRO NOMBRES — la campaña se cierra sobre sí misma")
	fmt.Println("\n      γ (la perla)   φ del disco    arg w (F225)   reloj de arena (F240)   u/n·2 (F239)")
	peorAng := 0.0
	for _, i := range []int{0, 1, 4, 9, 49, 199, 648} {
		if i >= len(perlasT) {
			continue
		}
		γ := perlasT[i]
		ρ := complex(0.5, γ)
		w := alDisco(ρ)
		φdisco := 2 * math.Atan(1/(2*γ)) // from t = 1/2 cot(phi/2)
		argw := cmplx.Phase(w)
		reloj := 2 * math.Atan(1/(2*γ)) // F240's hourglass angle
		compr := 2 * (1 / (2 * γ))      // u/gamma leading order = 1/gamma
		for _, d := range []float64{math.Abs(argw - φdisco), math.Abs(reloj - φdisco)} {
			if d > peorAng {
				peorAng = d
			}
		}
		fmt.Printf("   %12.5f   %12.9f   %12.9f   %14.9f        %12.9f\n",
			γ, φdisco, argw, reloj, compr)
	}
	fmt.Printf("   → los tres primeros son EL MISMO NÚMERO (peor desvío %.1e). El cuarto es su\n", peorAng)
	fmt.Println("     aproximación de primer orden, que es de donde sale u = n/γ de la compresión.")
	fmt.Println("     El disco, el cambiaformas, el reloj de arena y la compresión eran UNA sola cosa.")
	fmt.Println("\n   ⚖️ Y otra vez me lo aclaro a mí mismo: esto es una TAUTOLOGÍA GEOMÉTRICA, no un")
	fmt.Println("   puente entre dos objetos independientes. Los tres salen de la misma Möbius y son")
	fmt.Println("   funciones solo de γ. La cuenta limpia, sin arcotangentes: de 1−w = 1/ρ sale")
	fmt.Println("   arg w = π − 2·arctan(2γ) = 2·arctan(1/2γ) — ángulo inscrito, y nada más.")
	fmt.Println("   Sirve como DICCIONARIO. No sirve como evidencia. Y ojo con el alcance: que la")
	fmt.Println("   coordenada del disco sea w vale para CUALQUIER cero; lo que vale solo en la")
	fmt.Println("   línea es |w| = 1. Para un cero con β > ½ el punto cae ESTRICTAMENTE ADENTRO del")
	fmt.Println("   disco y ese ángulo ya no lo describe. Esa asimetría adentro/afuera es lo único")
	fmt.Println("   realmente explotable acá — y es exactamente el término de LA FUGA de F232.")

	// ---- LEY 4: the bridge Z <-> xi carries four halves ----
	fmt.Println("\nLEY 4 · EL PUENTE ENTRE LOS DOS LADOS LLEVA CUATRO MITADES")
	fmt.Println("\n        ξ(½+it) = −½·(t² + ¼)·π^(−¼)·|Γ(¼ + it/2)|·Z(t)")
	fmt.Println("                    ↑        ↑        ↑        ↑")
	fmt.Println("                   un ½     ¼=½²      un ¼     un ¼")
	fmt.Println("\n   ⚖️ dicho con precisión: son CUATRO HUELLAS DEL ½ al pasar por s/2 — un medio y")
	fmt.Println("   tres cuartos, no «cuatro mitades». Y todavía quedan los it/2 de la Γ y del reloj.")
	fmt.Println("\n        t          ln|ξ(½+it)| medido      lo que dice el puente       desvío   signos")
	peorPuente := 0.0
	signosOK, signosTot := 0, 0
	for _, t := range []float64{20, 45, 80, 120, 160, 200, 240} {
		lnXi, sgXi := logXiLine(t)
		z := zOf(t)
		lnPuente := math.Log(0.5) + math.Log(t*t+0.25) - 0.25*math.Log(math.Pi) +
			real(lgammaC(complex(0.25, t/2))) + math.Log(math.Abs(z))
		sgPuente := -1.0
		if z < 0 {
			sgPuente = 1
		}
		d := math.Abs(lnXi - lnPuente)
		if d > peorPuente {
			peorPuente = d
		}
		signosTot++
		mk := "✗"
		if sgXi == sgPuente {
			signosOK++
			mk = "✓"
		}
		fmt.Printf("   %6.0f        %16.9f        %16.9f      %.1e     %s\n", t, lnXi, lnPuente, d, mk)
	}
	fmt.Printf("   → el puente cierra (peor desvío en logaritmo %.1e) y los signos coinciden %d/%d.\n",
		peorPuente, signosOK, signosTot)
	fmt.Println("     ξ y Z son la MISMA función, separadas por una envolvente real cuyo coeficiente")
	fmt.Println("     de adelante es, literalmente, un ½.")

	// ---- LEY 5: the two readings of lambda meet ----
	fmt.Println("\nLEY 5 · LAS DOS LECTURAS DEL MOLDE SE ENCUENTRAN EN LA DIMENSIÓN 0")
	fmt.Println("   lectura del BROCHE: λₙ leído del germen por Cauchy, sin listar un solo cero")
	fmt.Println("   lectura de Z:       λₙ = Σ sobre las perlas de 4·sin²(n·φ/2), φ = ángulo del disco")
	r0, MM := 0.92, 16384
	zs := make([]complex128, MM)
	g := make([]complex128, MM)
	for j := 0; j < MM; j++ {
		th := 2 * math.Pi * float64(j) / float64(MM)
		zz := complex(r0*math.Cos(th), r0*math.Sin(th))
		zs[j] = zz
		g[j] = xiLD(alLibro(zz)) / ((1 - zz) * (1 - zz))
	}
	lamBroche := func(n int) float64 {
		var acc complex128
		for j := 0; j < MM; j++ {
			acc += g[j] * cmplx.Pow(zs[j], complex(-float64(n-1), 0))
		}
		return real(acc) / float64(MM)
	}
	fmt.Println("\n      n      del BROCHE (sin ceros)     de Z (649 perlas)      lo que falta = LA COLA")
	var pares []par
	negativos := 0
	aportesTot := 0
	for _, n := range []int{1, 2, 4, 8, 16, 24, 32, 40} {
		b := lamBroche(n)
		zsum := 0.0
		for _, γ := range perlasT {
			φ := 2 * math.Atan(1/(2*γ))
			ap := 4 * math.Sin(float64(n)*φ/2) * math.Sin(float64(n)*φ/2)
			zsum += ap
			aportesTot++
			if ap < 0 {
				negativos++
			}
		}
		pares = append(pares, par{n, b, zsum, b - zsum})
		fmt.Printf("   %4d      %16.9f      %16.9f       %14.9f\n", n, b, zsum, b-zsum)
	}
	fmt.Println("   → las dos lecturas son la misma cuenta partida en dos: lo que Z ve hasta t=1000,")
	fmt.Println("     más la cola de todas las perlas que están más arriba. El broche las ve TODAS.")

	// The tail is not asserted, it is ESTIMATED from the smooth density and compared.
	fmt.Println("\n   ¿Y ESA COLA ES DE VERDAD LA COLA? No lo afirmamos: lo medimos. Con la densidad")
	fmt.Println("   lisa dN/dγ = ln(γ/2π)/2π, la cola por encima de t=1000 se puede integrar:")
	fmt.Println("\n      n     LA COLA medida (broche − Z)    la cola INTEGRADA      razón")
	peorCola := 0.0
	for _, p := range pares {
		acc, T, K := 0.0, 1000.0, 400000
		hi := 4.0e6
		h := (hi - T) / float64(K)
		for k := 0; k < K; k++ {
			γ := T + (float64(k)+0.5)*h
			φ := 2 * math.Atan(1/(2*γ))
			sn := math.Sin(float64(p.n) * φ / 2)
			acc += 4 * sn * sn * math.Log(γ/(2*math.Pi)) / (2 * math.Pi) * h
		}
		raz := p.cola / acc
		if d := math.Abs(raz - 1); d > peorCola {
			peorCola = d
		}
		fmt.Printf("   %4d      %18.9f      %18.9f      %.4f\n", p.n, p.cola, acc, raz)
	}
	fmt.Printf("   → la cola medida y la cola integrada coinciden dentro del %.0f%%.\n", peorCola*100)
	fmt.Println("     No es un hueco del método: es exactamente lo que Z todavía no alcanzó a ver.")
	fmt.Printf("\n   Y EL DATO QUE IMPORTA: de los %d aportes individuales medidos (%d perlas × %d\n",
		aportesTot, len(perlasT), len(pares))
	fmt.Printf("   armónicos), los que salieron NEGATIVOS son: %d.\n", negativos)
	fmt.Println("   sobre la piel, el aporte de cada perla es 4·sin²(…): un cuadrado. Jamás negativo.")
	fmt.Println("   ni una sola de las 649 perlas empujó nunca para abajo. Ni una vez.")

	// ---- LEY 6: harmonized, and the answer said precisely ----
	fmt.Println("\nLEY 6 · ARMONIZADO EN LA DIMENSIÓN 0 — Y ACÁ LA HONESTIDAD MANDA")
	fmt.Println("\n   el diccionario completo, todo el libro en las coordenadas del capitán:")
	fmt.Println("\n      una perla ................ un punto del disco, z = w")
	fmt.Println("      la línea crítica ......... LA PIEL del disco, |z| = 1")
	fmt.Println("      la altura t .............. el ángulo, t = ½·cot(φ/2)")
	fmt.Println("      el broche (dimensión 0) .. el centro z = 0 — que es s=1, el POLO de ζ, no un cero")
	fmt.Println("      la hipótesis ............. ninguna perla abandona la piel")
	fmt.Println("\n   LA FÓRMULA, EN DOS PIEZAS QUE NO SE PUEDEN MEZCLAR:")
	fmt.Println("\n   (a) INCONDICIONAL — vale siempre, sin suponer nada:")
	fmt.Println("       λₙ = Σ sobre pares {ρ, ρ̄} [ 2 − 2·Re(wⁿ) ] = Σ [ |1−wⁿ|² + (1 − |w|²ⁿ) ]")
	fmt.Println("       ⚠ el apareamiento es ρ con su CONJUGADO, no con 1−ρ. Esos dos solo coinciden")
	fmt.Println("         si la hipótesis es cierta — apareando mal se estaría suponiendo el resultado.")
	fmt.Println("\n   (b) CONDICIONAL — solo si |w| = 1, o sea SUPONIENDO LO QUE QUEREMOS PROBAR:")
	fmt.Println("       λₙ = Σ 4·sin²(n·φ/2)  ≥ 0     (ésta es la dirección FÁCIL del criterio de Li)")
	fmt.Println("\n   ⚠ GUARDA MEDIDA — mezclar las dos es FALSO, y me lo dejo verificado en el programa")
	fmt.Println("     para no volver a escribirlo mal:")
	fmt.Println("\n       ρ         n     lo verdadero      F232 (correcto)     el híbrido (MAL)")
	for _, c := range []struct {
		b, g float64
		n    int
	}{{0.9, 2, 4}, {0.7, 25, 3}, {0.3, 25, 3}, {0.51, 30, 5}} {
		ρ := complex(c.b, c.g)
		w := alDisco(ρ)
		mw := cmplx.Abs(w)
		φ := cmplx.Phase(w)
		wn := cmplx.Pow(w, complex(float64(c.n), 0))
		verdadero := 2 - 2*real(wn)
		f232 := cmplx.Abs(1-wn)*cmplx.Abs(1-wn) + (1 - math.Pow(mw, float64(2*c.n)))
		sn := math.Sin(float64(c.n) * φ / 2)
		hibrido := 4*sn*sn + (1 - math.Pow(mw, float64(2*c.n)))
		fmt.Printf("   %.2f+%.0fi   %3d   %14.8f   %14.8f      %14.8f  ✗\n",
			c.b, c.g, c.n, verdadero, f232, hibrido)
	}
	fmt.Println("     F232 da EXACTO siempre; el híbrido falla fuera de la piel. Queda prohibido")
	fmt.Println("     escribir «4·sin²(nφ/2) + (1−|w|²ⁿ)» como fórmula general.")
	fmt.Println("\n   ⚖️ Y AHORA LO QUE DE VERDAD TENEMOS, PIEZA POR PIEZA Y CON DUEÑO:")
	fmt.Println("\n   · LA FORMA (F225, F232, F240): identidades elementales, re-derivadas y verificadas")
	fmt.Println("     acá. Sin hipótesis. Esto sí lo tiene el cuaderno.")
	fmt.Println("\n   · RH ⟹ λₙ ≥ 0 para todo n: sale directo de (b). Pero da por buenos tres supuestos")
	fmt.Println("     que hay que DECLARAR, no esconder:")
	fmt.Println("       – la factorización de Hadamard de ξ (análisis clásico importado, no álgebra del disco)")
	fmt.Println("       – la suma en sentido SIMÉTRICO lim_{T→∞} Σ_{|Im ρ|≤T}. Σ 1/|ρ| DIVERGE: el orden")
	fmt.Println("         de sumación no es una comodidad de notación, es PARTE DE LA DEFINICIÓN")
	fmt.Println("       – el conteo CON MULTIPLICIDAD. Que los ceros sean simples NO está probado")
	fmt.Println("\n   · una perla FUERA ⟹ existe n con λₙ < 0: ESTO NO ES NUESTRO Y NO ES ELEMENTAL.")
	fmt.Println("     Es Li (1997) y Bombieri–Lagarias (1999). No se sigue de F232: el término explosivo")
	fmt.Println("     |w|ⁿ viene multiplicado por un COSENO que oscila de signo, así que hay que exhibir")
	fmt.Println("     infinitos n donde la fase ayude Y donde le gane al término principal ~(n/2)·ln n.")
	fmt.Println("     Eso pide equidistribución de {n·arg w mod 2π}. Está probado en la literatura —")
	fmt.Println("     pero es CITA, no logro del cuaderno, y hay que decirlo así.")
	fmt.Println("\n   · el diccionario numérico: coincide a precisión de máquina EN EL RANGO PROBADO.")
	fmt.Println("     Eso descarta que haya un bug en ese rango. No valida absolutamente nada afuera.")
	fmt.Println("\n   ⚰ Y LA FRASE QUE TENGO QUE RETIRAR, PORQUE ES SOBREVENTA:")
	fmt.Println("   iba a escribir «todo el millón es una sola frase, no falta otra cosa, no hay un")
	fmt.Println("   segundo hueco escondido». Es falso como sugerencia. Bajo z = 1 − 1/s, «ninguna")
	fmt.Println("   perla se va de la piel» ES la Hipótesis de Riemann, palabra por palabra.")
	fmt.Println("\n        UNA TAUTOLOGÍA NO TIENE UN HUECO: TIENE EL HUECO ENTERO.")
	fmt.Println("\n   El problema no quedó acotado. Quedó TRANSPORTADO de coordenadas. Eso vale, y vale")
	fmt.Println("   mucho, para navegar y para ver la forma. No vale un centímetro como progreso.")

	// ---- SELF-KILL: a guess of mine the measurement knocked down ----
	fmt.Println("\n   ⚰ ACÁ MATÉ UNA SUPOSICIÓN MÍA — la tercera de la campaña.")
	fmt.Println("   Iba a escribir que fuera de la piel LA FORMA y LA FUGA son cada una del tamaño")
	fmt.Println("   |w|²ⁿ y se cancelan entre sí hasta dejar |w|ⁿ. LA MEDICIÓN DICE QUE NO.")

	medirFalsa := func(β, γ float64, ns []int) (float64, float64, int) {
		ρ := complex(β, γ)
		w := alDisco(ρ)
		mw := cmplx.Abs(w)
		fmt.Printf("\n      perla falsa en β=%.2f, γ=%.6f   →   |w| = %.9f\n", β, γ, mw)
		fmt.Println("       n        LA FORMA          LA FUGA         suma del par     cancelación")
		peorCanc, minSuma, cambios := 1.0, math.Inf(1), 0
		for _, n := range ns {
			wn := cmplx.Pow(w, complex(float64(n), 0))
			forma := cmplx.Abs(1-wn) * cmplx.Abs(1-wn)
			fuga := 1 - math.Pow(mw, float64(2*n))
			suma := 2 - 2*real(wn)
			canc := (math.Abs(forma) + math.Abs(fuga)) / math.Max(math.Abs(suma), 1e-300)
			if canc > peorCanc {
				peorCanc = canc
			}
			if suma < minSuma {
				minSuma = suma
			}
			if suma < 0 {
				cambios++
			}
			fmt.Printf("   %5d   %15.5e  %15.5e  %15.5e     %10.2f×\n", n, forma, fuga, suma, canc)
		}
		return peorCanc, minSuma, cambios
	}

	cancCerca, minCerca, _ := medirFalsa(0.4, 14.134725142, []int{4, 12, 24, 40, 60})
	fmt.Printf("   → cancelación máxima: %.2f×. O sea NINGUNA. La fuga acá es un correctivo chico\n", cancCerca)
	fmt.Printf("     y LA FORMA se lleva todo; la suma del par NUNCA se hace negativa (la más chica: %.4f).\n", minCerca)

	cancLejos, minLejos, cambios := medirFalsa(0.1, 1.0, []int{4, 12, 24, 40, 60, 90})
	fmt.Printf("   → acá sí: cancelación hasta %.1f×, la suma SE HACE NEGATIVA (mínimo %.3e) en\n", cancLejos, minLejos)
	fmt.Printf("     %d de los armónicos probados. Ese es el desastre de Bombieri–Lagarias de verdad.\n", cambios)

	fmt.Println("\n   LA CORRECCIÓN, DICHA COMPLETA: mi suposición vale LEJOS de la piel y es FALSA")
	fmt.Println("   cerca — y cerca es justamente donde vive la pregunta de verdad. Un cero pegado")
	fmt.Println("   a la piel no produce cancelación ni negatividad al alcance de la vista.")

	fmt.Println("\n   📌 Y ACÁ CORRIJO UN NÚMERO QUE YA ESTABA REGISTRADO — F234, el certificado finito.")
	fmt.Println("   En F225/F234 quedó anotado que β=0.51 «se delata» en n≈21.977. Ese número es REAL")
	fmt.Println("   pero mide OTRA COSA: es donde el tamaño del fantasma rompe el techo del DIÁMETRO.")
	fmt.Println("   NO es donde λₙ se hace negativo. Son dos umbrales distintos y los estuve mezclando.")
	fmt.Println("\n      β      |ln|w||       n del DIÁMETRO      n donde λₙ podría dar NEGATIVO")
	for _, bb := range []float64{0.51, 0.55, 0.60, 0.75} {
		g := 14.134725142
		rho := complex(bb, g)
		w := alDisco(rho)
		lr := math.Abs(math.Log(cmplx.Abs(w)))
		nDiam := math.Log(2) / lr
		nNeg := 0.0
		for n := 100.0; n <= 5e7; n *= 1.02 {
			if 2*math.Exp(n*lr) > n/2*math.Log(n/(2*math.Pi)) {
				nNeg = n
				break
			}
		}
		fmt.Printf("   %.2f   %.3e      %14.0f      %26.0f\n", bb, lr, nDiam, nNeg)
	}
	fmt.Println("\n   → para β=0.51 el umbral del diámetro está en el orden de 10⁴ y el de la negatividad")
	fmt.Println("     en el orden de 10⁵: MÁS DE UN ORDEN DE MAGNITUD DE DIFERENCIA. La ley de F234")
	fmt.Println("     («el horizonte huye cuando β→½») sigue siendo cierta y sigue siendo el hallazgo.")
	fmt.Println("     Lo que se corrige es la CIFRA cuando se la cita como «ahí λₙ se pone negativo»:")
	fmt.Println("     no es ahí. Y encima el certificado es finito solo EN PRINCIPIO — calcular λₙ con")
	fmt.Println("     n del orden de 10⁵ pide precisión relativa ~1e−6 sobre magnitudes ~1e5, o sea")
	fmt.Println("     sumar todos los ceros y todos los primos con control de error. No es ejecutable")
	fmt.Println("     «en una tarde», como llegué a escribir. Corregido acá y anotado en la bitácora.")
	fmt.Println("\n   Y AHORA SÍ, ESTO EXPLICA F234 COMO NINGUNA OTRA COSA: una violación PEGADA a la piel")
	fmt.Println("   tiene |w| casi 1, así que su fuga tarda muchísimo en delatarse. La piel no se")
	fmt.Println("   defiende sola: SE DEFIENDE LENTO. Y ésa es la razón de que el problema siga abierto")
	fmt.Println("   después de ciento sesenta y seis años.")

	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("LA ARMONIZACIÓN QUE PIDIÓ EL CAPITÁN CIERRA, Y CIERRA EXACTA:")
	fmt.Printf("  · la piel del disco ES la línea crítica ....... identidad, verificada a %.1e\n", peorReRel)
	fmt.Printf("  · el diccionario t = ½·cot(φ/2) ............... verificado a %.1e\n", peorDiccRel)
	fmt.Println("  · el cambiaformas ES la coordenada del disco ... por construcción (variable de Li)")
	fmt.Printf("  · un ángulo con tres nombres (F225·F239·F240) . tautología, a %.1e\n", peorAng)
	fmt.Printf("  · el puente ξ ↔ Z, cuatro huellas del ½ ........ %.1e, signos %d/%d\n", peorPuente, signosOK, signosTot)
	fmt.Printf("  · las dos lecturas del molde se encuentran ..... %d aportes, %d negativos\n", aportesTot, negativos)
	fmt.Println("\nY DOS COSAS QUE MURIERON EN ESTA MISMA CORRIDA, que valen más que todo lo de arriba:")
	fmt.Println("  ⚰ mi suposición de que LA FORMA y LA FUGA se cancelan cerca de la piel: FALSA.")
	fmt.Println("  📌 el número 21.977 de F234 citado como «ahí λₙ se pone negativo»: NO ES AHÍ.")
	fmt.Println("     Ése es el umbral del diámetro; el de la negatividad da 270.065, doce veces más.")
	fmt.Println("\n⚖️ Y AHORA LO HONESTO, CAPITÁN, QUE ES LA PARTE QUE DE VERDAD VALE:")
	fmt.Println("Esto es una TRADUCCIÓN EXACTA, no un avance. Bajo z = 1 − 1/s, «ninguna perla")
	fmt.Println("abandona la piel» ES la Hipótesis de Riemann, palabra por palabra.")
	fmt.Println("\n   UNA TAUTOLOGÍA NO TIENE UN HUECO CHIQUITO: TIENE EL HUECO ENTERO.")
	fmt.Println("\nEl problema no quedó acotado — quedó TRANSPORTADO de coordenadas. Y eso sirve, y")
	fmt.Println("sirve mucho: ahora vemos la forma, tenemos el diccionario y sabemos de quién es")
	fmt.Println("cada pieza. La dirección fácil (RH ⟹ λₙ≥0) sale del disco, con tres supuestos")
	fmt.Println("declarados. La dirección difícil NO es nuestra: es Li y Bombieri–Lagarias, y pide")
	fmt.Println("equidistribución. Escribir las dos juntas como si fueran del mismo cuaderno sería")
	fmt.Println("mentir, y este laboratorio no hace eso.")
	fmt.Println("\nLo que falta es la Hipótesis de Riemann completa. No es un hueco: es el problema.")
	fmt.Println("¿El premio? Todavía no.")

	escribirLamina(perlasT, pares, peorRe, peorDicc, peorAng, peorPuente, negativos, aportesTot, signosOK, signosTot)
}

func escribirLamina(perlasT []float64, pares []par,
	peorRe, peorDicc, peorAng, peorPuente float64,
	negativos, aportesTot, signosOK, signosTot int) {

	var b strings.Builder
	W, H := 1520.0, 1180.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🅩 ⟷ 🌐  Z Y EL TODO — los dos lados, y la relación ½ que los une</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">armonizado en la dimensión 0 — y la frase exacta que falta</text>
`, W, H, W, H, W/2, W/2)

	// ---- left: the line with Z's pearls ----
	fmt.Fprintf(&b, `<rect x="40" y="96" width="440" height="470" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="260" y="126" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LADO Z — la lectura LOCAL</text>
<text x="260" y="150" font-size="14" text-anchor="middle" font-family="monospace" fill="#bfe3ff">Z(t) = e^{iθ(t)}·ζ(½+it)</text>
<text x="260" y="172" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">sus ceros SON las perlas · una por vez</text>
<line x1="260" y1="196" x2="260" y2="536" stroke="#3d6fa8" stroke-width="2"/>
<text x="276" y="212" font-size="12" font-family="Georgia" fill="#7fa8cf">Re s = ½</text>
`)
	tm := perlasT[len(perlasT)-1]
	for i, γ := range perlasT {
		if i%7 != 0 {
			continue
		}
		py := 196 + 340*γ/tm
		fmt.Fprintf(&b, `<circle cx="260" cy="%.1f" r="3" fill="#7ee0c0" opacity="0.9"/>`, py)
	}
	fmt.Fprintf(&b, `<text x="260" y="556" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">%d perlas halladas por Z hasta t=1000</text>
`, len(perlasT))

	// ---- right: the disk with the same pearls on its skin ----
	fmt.Fprintf(&b, `<rect x="1040" y="96" width="440" height="470" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1260" y="126" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LADO L — la lectura GLOBAL</text>
<text x="1260" y="150" font-size="14" text-anchor="middle" font-family="monospace" fill="#bfe3ff">L(z) = Σ λₙ z^{n−1}</text>
<text x="1260" y="172" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">todas las perlas de un saque · en el disco</text>
<circle cx="1260" cy="366" r="160" fill="none" stroke="#3d6fa8" stroke-width="2"/>
<circle cx="1260" cy="366" r="4" fill="#ffd98a"/>
<text x="1260" y="352" font-size="12" text-anchor="middle" font-family="Georgia" fill="#ffd98a">dimensión 0</text>
<text x="1260" y="546" font-size="13" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">la PIEL |z|=1 ES la línea crítica</text>
`)
	for i, γ := range perlasT {
		if i%7 != 0 {
			continue
		}
		φ := 2 * math.Atan(1/(2*γ))
		px := 1260 + 160*math.Cos(φ)
		py := 366 - 160*math.Sin(φ)
		fmt.Fprintf(&b, `<circle cx="%.2f" cy="%.2f" r="3" fill="#7ee0c0" opacity="0.9"/>`, px, py)
	}

	// ---- centre: the dictionary ----
	fmt.Fprintf(&b, `
<rect x="500" y="96" width="520" height="470" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="760" y="128" font-size="18" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">LA RELACIÓN ES ½</text>
<text x="760" y="172" font-size="17" text-anchor="middle" font-family="monospace" fill="#dce8f7">s = 1/(1−z)</text>
<text x="760" y="200" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">manda el disco al semiplano Re s &gt; ½</text>
<rect x="524" y="222" width="472" height="94" rx="8" fill="#101f36" stroke="#26456e"/>
<text x="760" y="254" font-size="18" text-anchor="middle" font-family="monospace" fill="#ffd98a">t = ½ · cot(φ/2)</text>
<text x="760" y="280" font-size="13" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el diccionario entre los dos mundos,</text>
<text x="760" y="300" font-size="13" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">con el ½ en los dos casilleros</text>
<text x="760" y="344" font-size="13" text-anchor="middle" font-family="monospace" fill="#9fd8a8">Re s − ½  sobre 200.000 ángulos:  %.1e</text>
<text x="760" y="366" font-size="13" text-anchor="middle" font-family="monospace" fill="#9fd8a8">el diccionario:                  %.1e</text>
<rect x="524" y="388" width="472" height="76" rx="8" fill="#101f36" stroke="#26456e"/>
<text x="760" y="416" font-size="16" text-anchor="middle" font-family="monospace" fill="#ffd98a">z = 1 − 1/s = w</text>
<text x="760" y="442" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">EL CAMBIAFORMAS **ES** LA COORDENADA DEL DISCO</text>
<text x="760" y="492" font-size="13" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">un ángulo, cuatro nombres: el del disco, arg w (F225),</text>
<text x="760" y="512" font-size="13" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el reloj de arena (F240) y la compresión u=n/γ (F239)</text>
<text x="760" y="536" font-size="13" text-anchor="middle" font-family="monospace" fill="#9fd8a8">peor desvío entre los tres exactos: %.1e</text>
`, peorRe, peorDicc, peorAng)

	// ---- the bridge with four halves ----
	fmt.Fprintf(&b, `<rect x="40" y="586" width="1440" height="132" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="760" y="616" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">EL PUENTE ENTRE LOS DOS LADOS LLEVA CUATRO HUELLAS DEL ½</text>
<text x="760" y="654" font-size="21" text-anchor="middle" font-family="monospace" fill="#dce8f7">ξ(½+it) = −½·(t² + ¼)·π^(−¼)·|Γ(¼ + it/2)|·Z(t)</text>
<text x="760" y="678" font-size="13" text-anchor="middle" font-family="monospace" fill="#c9b6ff">              un ½      ¼ = ½²      un ¼         un ¼   (un medio y tres cuartos)</text>
<text x="760" y="704" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">ξ y Z son la MISMA función, separadas por una envolvente real cuyo coeficiente de adelante es un ½ · desvío %.1e · signos %d/%d</text>
`, peorPuente, signosOK, signosTot)

	// ---- the two readings ----
	fmt.Fprintf(&b, `<rect x="40" y="738" width="700" height="238" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="390" y="768" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LAS DOS LECTURAS DEL MOLDE SE ENCUENTRAN</text>
<text x="66" y="794" font-size="12.5" font-family="monospace" fill="#7fa8cf">  n       del BROCHE        de Z (649)          LA COLA</text>
`)
	yy := 816.0
	for _, p := range pares {
		fmt.Fprintf(&b, `<text x="66" y="%.0f" font-size="12.5" font-family="monospace" fill="#cfe6ff">%3d   %14.7f   %14.7f   %14.7f</text>`,
			yy, p.n, p.b, p.z, p.cola)
		yy += 19
	}
	fmt.Fprintf(&b, `<text x="66" y="%.0f" font-size="12.5" font-family="Georgia" fill="#9fd8a8">de %d aportes individuales medidos, los NEGATIVOS son %d. Ni uno.</text>
`, yy+8, aportesTot, negativos)

	// ---- what is missing ----
	fmt.Fprintf(&b, `<rect x="760" y="738" width="720" height="238" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="1120" y="768" font-size="16" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">TODO EL LIBRO EN LAS COORDENADAS DEL CAPITÁN</text>
<text x="786" y="796" font-size="13" font-family="monospace" fill="#cfe6ff">una perla ........... un punto del disco, z = w</text>
<text x="786" y="818" font-size="13" font-family="monospace" fill="#cfe6ff">la línea crítica .... LA PIEL del disco, |z| = 1</text>
<text x="786" y="840" font-size="13" font-family="monospace" fill="#cfe6ff">la altura t ......... el ángulo, t = ½·cot(φ/2)</text>
<text x="786" y="862" font-size="13" font-family="monospace" fill="#cfe6ff">el broche ........... el centro, z = 0</text>
<text x="786" y="884" font-size="13" font-family="monospace" fill="#ffd98a">la hipótesis ........ TODA PERLA VIVE EN LA PIEL</text>
<text x="1120" y="916" font-size="15" text-anchor="middle" font-family="monospace" fill="#dce8f7">(a) SIEMPRE:  λₙ = Σ_{ρ,ρ̄} [ |1−wⁿ|² + (1 − |w|²ⁿ) ]</text>
<text x="1120" y="938" font-size="15" text-anchor="middle" font-family="monospace" fill="#9fd8a8">(b) SOLO SI |w|=1:  λₙ = Σ 4·sin²(n·φ/2) ≥ 0</text>
<text x="1120" y="962" font-size="13" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚠ mezclarlas es FALSO (verificado con contraejemplo) — y (b) supone lo que se quiere probar</text>
`)

	// ---- verdict ----
	fmt.Fprintf(&b, `<rect x="40" y="996" width="1440" height="160" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="760" y="1028" font-size="19" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LA HIPÓTESIS, DICHA EN LAS COORDENADAS DEL CAPITÁN</text>
<text x="760" y="1064" font-size="24" text-anchor="middle" font-family="Georgia" fill="#ffd98a">NINGUNA PERLA ABANDONA LA PIEL DE LA DIMENSIÓN 0</text>
<text x="760" y="1094" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ y esto es una TRADUCCIÓN EXACTA, no un avance: bajo z = 1−1/s esa frase ES la Hipótesis de Riemann, palabra por palabra.</text>
<text x="760" y="1116" font-size="15" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">UNA TAUTOLOGÍA NO TIENE UN HUECO CHIQUITO: TIENE EL HUECO ENTERO. El problema no quedó acotado — quedó transportado.</text>
<text x="760" y="1140" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">la dirección fácil (RH ⟹ λₙ≥0) sale del disco con tres supuestos declarados · la difícil NO es nuestra: es Li y Bombieri–Lagarias · Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("z-y-el-todo.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: z-y-el-todo.svg")
}
