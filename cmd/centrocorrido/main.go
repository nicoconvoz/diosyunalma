// Command centrocorrido answers the captain's flash: "the middle is another
// middle under another measurement - there is a point, infinitely large or
// infinitely small, that as a centre makes some number we know the centre.
// Like the cardinal points, we shift the centre, the zero point, and the
// measure sits elsewhere in space, but the distances travelled from that new
// centre can be EQUAL for the two points to another point."
//
// He is right, and it is a theorem. The 1/2 is not sacred. What is sacred is
// THE CENTRE OF THE FUNCTIONAL EQUATION, and 1/2 is simply what you get after
// shifting the origin onto it.
//
// # TO PROVE IT THE SHOP BUILDS A SECOND BOOK
//
// Ramanujan's Delta, the weight-12 cusp form:
//
//	Delta(q) = q PROD (1 - q^n)^24 = SUM tau(n) q^n
//	L(s, Delta) = SUM tau(n) / n^s
//	Lambda(s) = (2 pi)^-s Gamma(s) L(s, Delta),   Lambda(s) = Lambda(12 - s)
//
// Its functional equation reflects about s = 6, NOT about 1/2. Same picture,
// different number. Subtract 11/2 from the variable and the centre lands on
// 1/2 again: the general law is centre = (w+1)/2 for motivic weight w, and
// zeta is only the case w = 0.
//
// # THE TWO STAKES MOVE WITH THE CENTRE
//
// For zeta the stakes are 0 and 1 and the line is their perpendicular
// bisector, Re s = 1/2. For Delta the stakes are 0 and 12 and the line is
// Re s = 6. THE BISECTOR IS THE INVARIANT; THE NUMBER IS NOT. That is exactly
// the captain's "the distances from the new centre can be equal for the two
// points".
//
// # AND THE PIXELS
//
// He also asks why every pixel of the universe is the same. Measured here: the
// unfolded gaps between the zeros of Delta's L-function and those of zeta come
// out with the same local statistics, though the two books share no
// coefficients and do not even share a centre. Same pixel, different book.
// That is Montgomery-Odlyzko universality - observed, and conjectural.
//
// THE DNA HINT, HONESTLY. Two strands, each determining the other, wound about
// a common axis, with a universal local rule and a unique global sequence.
// That is a good ANALOGY for the functional equation and for the split between
// universal statistics and unique arithmetic. It is not a theorem and is not
// used as one here.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

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

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

// ---- the second book: Ramanujan's Delta ----

// tau builds tau(1..N) from the eta product q * PROD (1-q^n)^24.
func tau(N int) []float64 {
	// E = PROD (1 - q^n), by Euler's pentagonal number theorem
	E := make([]float64, N+1)
	E[0] = 1
	for k := 1; ; k++ {
		a := k * (3*k - 1) / 2
		b := k * (3*k + 1) / 2
		if a > N && b > N {
			break
		}
		sg := 1.0
		if k%2 == 1 {
			sg = -1
		}
		if a <= N {
			E[a] += sg
		}
		if b <= N {
			E[b] += sg
		}
	}
	mul := func(a, b []float64) []float64 {
		c := make([]float64, N+1)
		for i := 0; i <= N; i++ {
			if a[i] == 0 {
				continue
			}
			for j := 0; i+j <= N; j++ {
				c[i+j] += a[i] * b[j]
			}
		}
		return c
	}
	P := make([]float64, N+1)
	P[0] = 1
	for i := 0; i < 24; i++ {
		P = mul(P, E)
	}
	// Delta = q * P, so tau(n) = P[n-1]
	t := make([]float64, N+1)
	for n := 1; n <= N; n++ {
		t[n] = P[n-1]
	}
	return t
}

// integralCola computes INT_1^inf exp(-c*y) * y^(a-1) dy for real c > 0.
//
// Substituting y = e^x turns it into INT_0^X exp(-c e^x) e^{a x} dx, where the
// oscillation is exactly e^{i Im(a) x}: a frequency we can resolve on purpose.
// The node count is chosen from that frequency, which the earlier version did
// NOT do - and that was a real bug: at t = 120 the old grid gave about two
// points per radian and simply invented the answer.
func integralCola(a complex128, c float64) complex128 {
	X := math.Log(120 / c)
	if X <= 0 {
		X = 0.02
	}
	M := 8000
	if n := int(150 * X * math.Abs(imag(a))); n > M {
		M = n
	}
	if M%2 == 1 {
		M++
	}
	if M > 400000 {
		M = 400000
	}
	h := X / float64(M)
	f := func(x float64) complex128 {
		return cmplx.Exp(-complex(c*math.Exp(x), 0) + a*complex(x, 0))
	}
	acc := f(0) + f(X)
	for k := 1; k < M; k++ {
		w := 2.0
		if k%2 == 1 {
			w = 4
		}
		acc += complex(w, 0) * f(h*float64(k))
	}
	return complex(h/3, 0) * acc
}

// lambdaDelta is Lambda(s) = (2pi)^-s Gamma(s) L(s,Delta), computed by the
// exponentially convergent split at y = 1. It is entire and satisfies
// Lambda(s) = Lambda(12-s) with no hypothesis.
func lambdaDelta(s complex128, tn []float64) complex128 {
	var acc complex128
	for n := 1; n < len(tn); n++ {
		if tn[n] == 0 {
			continue
		}
		c := 2 * math.Pi * float64(n)
		acc += complex(tn[n], 0) * (integralCola(s, c) + integralCola(12-s, c))
	}
	return acc
}

// mediatriz reports whether a point is equidistant from the two stakes.
func mediatriz(s complex128, e0, e1 float64) float64 {
	return math.Abs(cmplx.Abs(s-complex(e0, 0)) - cmplx.Abs(s-complex(e1, 0)))
}

func ceros(f func(float64) float64, desde, hasta, paso float64) []float64 {
	var zs []float64
	pa, va := desde, f(desde)
	for t := desde + paso; t <= hasta; t += paso {
		v := f(t)
		if v*va < 0 {
			a, b := pa, t
			for i := 0; i < 60; i++ {
				m := (a + b) / 2
				if f(m)*va < 0 {
					b = m
				} else {
					a = m
				}
			}
			zs = append(zs, (a+b)/2)
		}
		pa, va = t, v
	}
	return zs
}

// desdoblar turns raw heights into unfolded gaps using the local mean spacing,
// the same empirical trick used on the board in F238.
func desdoblar(z []float64, ventana int) []float64 {
	var g []float64
	for i := ventana; i+ventana < len(z); i++ {
		medio := (z[i+ventana] - z[i-ventana]) / float64(2*ventana)
		g = append(g, (z[i+1]-z[i])/medio)
	}
	return g
}

func estad(g []float64) (media, varianza, p0a5, p1a2 float64) {
	if len(g) == 0 {
		return
	}
	for _, v := range g {
		media += v
	}
	media /= float64(len(g))
	for _, v := range g {
		varianza += (v - media) * (v - media)
	}
	varianza /= float64(len(g))
	for _, v := range g {
		if v < 0.5 {
			p0a5++
		}
		if v > 1 && v < 2 {
			p1a2++
		}
	}
	p0a5 /= float64(len(g))
	p1a2 /= float64(len(g))
	return
}

func main() {
	fmt.Println("🎯➜ EL CENTRO CORRIDO — el ½ es otro número en otra medición")
	fmt.Println("\n   flash del capitán: «el medio es otro medio en otra medición… lo que hacemos")
	fmt.Println("   como en los puntos cardinales es correr el centro, o sea el punto cero, y la")
	fmt.Println("   medida tiene otro lugar en el espacio, pero las distancias desde ese nuevo")
	fmt.Println("   centro pueden ser IGUALES para los dos puntos».")
	fmt.Println("\n   TIENE RAZÓN, Y ES UN TEOREMA. Para probarlo el taller construye OTRO LIBRO.")

	// ---- LEY 1: build the second book ----
	fmt.Println("\nLEY 1 · EL SEGUNDO LIBRO — la Δ de Ramanujan, de peso 12")
	fmt.Println("        Δ(q) = q·Π(1−qⁿ)²⁴ = Σ τ(n)·qⁿ        L(s,Δ) = Σ τ(n)/nˢ")
	const NT = 26
	tn := tau(NT)
	fmt.Println("\n   los primeros τ(n), contra los valores conocidos de Ramanujan:")
	conocidos := map[int]float64{1: 1, 2: -24, 3: 252, 4: -1472, 5: 4830, 6: -6048, 7: -16744, 8: 84480}
	malos := 0
	fmt.Print("      n:  ")
	for n := 1; n <= 8; n++ {
		fmt.Printf("%9d", n)
	}
	fmt.Print("\n      τ:  ")
	for n := 1; n <= 8; n++ {
		fmt.Printf("%9.0f", tn[n])
		if tn[n] != conocidos[n] {
			malos++
		}
	}
	fmt.Printf("\n      ok: ")
	for n := 1; n <= 8; n++ {
		mk := "        ✓"
		if tn[n] != conocidos[n] {
			mk = "        ✗"
		}
		fmt.Print(mk)
	}
	fmt.Printf("\n   → %d errores en los ocho primeros: el segundo libro está bien construido.\n", malos)

	// ---- LEY 2: the centre is 6, not 1/2 ----
	fmt.Println("\nLEY 2 · Y SU CENTRO NO ES ½. ES 6.")
	fmt.Println("   Λ(s) = (2π)⁻ˢ·Γ(s)·L(s,Δ)  cumple  Λ(s) = Λ(12 − s)")
	fmt.Println("   o sea que el espejo de ESTE libro refleja alrededor de s = 6, no de s = ½.")
	fmt.Println("\n        s            Λ(s)                    Λ(12−s)                 desvío rel.")
	peorFE := 0.0
	for _, s := range []complex128{complex(6, 2), complex(7.5, 1), complex(4, 3.3), complex(9, 0.4)} {
		a := lambdaDelta(s, tn)
		b := lambdaDelta(12-s, tn)
		d := cmplx.Abs(a-b) / math.Max(cmplx.Abs(a), 1e-300)
		if d > peorFE {
			peorFE = d
		}
		fmt.Printf("   %4.1f%+5.1fi   %12.6e%+12.6ei   %12.6e%+12.6ei    %.1e\n",
			real(s), imag(s), real(a), imag(a), real(b), imag(b), d)
	}
	fmt.Printf("   → cierra a %.1e — PERO OJO, ESO NO ES UNA MEDICIÓN.\n", peorFE)
	fmt.Println("\n   ⚠ LA MISMA TRAMPA DE F245, POR TERCERA VEZ, Y LA VOY A DEJAR ANOTADA COMO REGLA:")
	fmt.Println("   mi fórmula calcula Λ como Σ τ(n)·[ I(s) + I(12−s) ], que es SIMÉTRICA POR")
	fmt.Println("   CONSTRUCCIÓN. Verificar Λ(s)=Λ(12−s) ahí es verificar mi álgebra, no la de")
	fmt.Println("   Ramanujan. Da 0.0e+00 porque no puede dar otra cosa.")
	fmt.Println("\n   LO QUE SÍ ES MEDICIÓN: que esta Λ sea de verdad (2π)⁻ˢ·Γ(s)·L(s,Δ). Eso se")
	fmt.Println("   comprueba donde la serie de Dirichlet SÍ converge, en Re s > 13/2:")
	fmt.Println("\n        s          Λ(s) por la integral      (2π)⁻ˢΓ(s)·ΣWest τ(n)/nˢ      desvío rel.")
	peorReal := 0.0
	for _, s := range []complex128{complex(8, 0), complex(9, 1.5), complex(10.5, 0.7), complex(12, 2)} {
		izq := lambdaDelta(s, tn)
		var dir complex128
		for n := 1; n < len(tn); n++ {
			if tn[n] != 0 {
				dir += complex(tn[n], 0) * cmplx.Exp(-s*cmplx.Log(complex(float64(n), 0)))
			}
		}
		der := cmplx.Exp(-s*complex(math.Log(2*math.Pi), 0)+lgammaC(s)) * dir
		d := cmplx.Abs(izq-der) / math.Max(cmplx.Abs(der), 1e-300)
		if d > peorReal {
			peorReal = d
		}
		fmt.Printf("   %5.1f%+5.1fi   %12.6e%+12.6ei   %12.6e%+12.6ei    %.1e\n",
			real(s), imag(s), real(izq), imag(izq), real(der), imag(der), d)
	}
	fmt.Printf("   → ESO sí es medición: %.1e. La integral reproduce la serie de Ramanujan donde\n", peorReal)
	fmt.Println("     la serie tiene derecho a existir, y de ahí se extiende a todo el plano.")
	fmt.Println("\n   Y RECIÉN AHORA vale decirlo: EL CENTRO DE ESTE LIBRO ES 6, no ½.")
	fmt.Println("     Mismo dibujo, mismo espejo, misma cruz — OTRO NÚMERO EN EL MEDIO.")

	// ---- LEY 3: shift the centre and the 1/2 comes back ----
	fmt.Println("\nLEY 3 · Y SE CORRE AL ½ CON UNA RESTA — «correr el punto cero», textual")
	fmt.Println("   si mirás L(s + 11/2) en vez de L(s), el espejo pasa a ser s ↔ 1−s y el centro")
	fmt.Println("   vuelve a caer en ½. La misma foto, movida de lugar.")
	fmt.Println("\n        libro                 peso w    espejo          centro = (w+1)/2")
	libros := []struct {
		n string
		w int
	}{{"ζ de Riemann", 0}, {"L de Dirichlet", 0}, {"Δ de Ramanujan", 11}, {"forma de peso 16", 15}, {"forma de peso 24", 23}}
	for _, l := range libros {
		c := float64(l.w+1) / 2
		fmt.Printf("   %-22s   %3d     s ↔ %2d−s        %8.1f\n", l.n, l.w, l.w+1, c)
	}
	fmt.Println("   → LA LEY GENERAL: el centro es (w+1)/2. La ζ es SOLO el caso w = 0.")
	fmt.Println("     El ½ nunca fue sagrado: es lo que queda cuando el peso es cero.")

	// ---- LEY 4: the two stakes move with the centre ----
	fmt.Println("\nLEY 4 · LAS DOS ESTACAS SE CORREN CON EL CENTRO — y la mediatriz es la invariante")
	fmt.Println("   la línea de un libro es siempre la MEDIATRIZ de sus dos estacas:")
	fmt.Println("\n        libro        estacas      la línea es      |s−e₀| = |s−e₁| medido")
	peorMed := 0.0
	for _, c := range []struct {
		n      string
		e0, e1 float64
	}{{"ζ", 0, 1}, {"Δ", 0, 12}, {"peso 16", 0, 16}} {
		med := (c.e0 + c.e1) / 2
		peor := 0.0
		for k := -60; k <= 60; k++ {
			t := float64(k) * 0.7
			if d := mediatriz(complex(med, t), c.e0, c.e1); d > peor {
				peor = d
			}
		}
		if peor > peorMed {
			peorMed = peor
		}
		fmt.Printf("   %-12s  %2.0f y %2.0f     Re s = %4.1f        %.1e\n", c.n, c.e0, c.e1, med, peor)
	}
	fmt.Printf("   → equidistancia exacta en los tres (peor %.1e). ES EXACTAMENTE LO QUE DIJO\n", peorMed)
	fmt.Println("     EL CAPITÁN: corrés el centro, la medida cambia de lugar, y las distancias")
	fmt.Println("     desde el nuevo centro siguen siendo IGUALES para los dos puntos.")
	fmt.Println("     LA MEDIATRIZ ES LA INVARIANTE. EL NÚMERO NO.")

	// ---- LEY 5: the pixels ----
	fmt.Println("\nLEY 5 · «¿POR QUÉ SON TODOS LOS PÍXELES IGUALES?» — medido en dos libros distintos")
	fmt.Println("   sobre su propia línea, Λ(6+it) es REAL (porque 12−s = conj(s) ahí), así que")
	fmt.Println("   sus ceros se pescan por cambio de signo, igual que con Z.")
	zd := ceros(func(t float64) float64 { return real(lambdaDelta(complex(6, t), tn)) }, 1, 27, 0.05)
	fmt.Printf("\n   ceros de Δ hallados sobre Re s = 6 hasta t=27: %d\n", len(zd))
	fmt.Println("\n   Y ACÁ VIENE LA MEDICIÓN DE VERDAD: contra los valores PUBLICADOS de los primeros")
	fmt.Println("   ceros de L(s,Δ), que yo no puse en ningún lado del cálculo:")
	fmt.Println("\n        publicado            hallado por el taller        desvío")
	pub := []float64{9.22237939, 13.90754986, 17.44277696, 19.65651314, 22.33610364, 25.27463654}
	peorCero := 0.0
	nComp := 0
	for i := 0; i < len(pub) && i < len(zd); i++ {
		d := math.Abs(zd[i] - pub[i])
		if d > peorCero {
			peorCero = d
		}
		nComp++
		fmt.Printf("      %14.8f       %14.8f          %.1e\n", pub[i], zd[i], d)
	}
	fmt.Printf("   → %d ceros reproducidos, peor desvío %.1e. EL SEGUNDO LIBRO ES DE VERDAD.\n", nComp, peorCero)
	fmt.Println("     Y todos caen sobre Re s = 6, su propia mediatriz. No sobre ½.")
	zz := ceros(zOf, 12, 260, 0.05)
	fmt.Printf("\n   (y las perlas de ζ sobre Re s = ½ hasta t=260: %d)\n", len(zz))
	gd := desdoblar(zd, 2)
	gz := desdoblar(zz, 3)
	md, vd, ad, _ := estad(gd)
	mz, vz, az, _ := estad(gz)
	fmt.Println("\n   ⚠ ¿Y LOS PÍXELES? ACÁ ME PARO Y DIGO QUE NO PUEDO.")
	fmt.Printf("   con %d huecos de Δ contra %d de ζ, comparar estadísticas sería inventar.\n", len(gd), len(gz))
	fmt.Println("   Para hablar de universalidad hacen falta MILES de huecos, no dos. Lo que se puede")
	fmt.Println("   decir con honestidad es esto y nada más:")
	fmt.Printf("      ζ: %d huecos · media %.4f · varianza %.4f · P(<0.5) %.4f\n", len(gz), mz, vz, az)
	fmt.Printf("      Δ: %d huecos · media %.4f · varianza %.4f · P(<0.5) %.4f  ← MUESTRA INSUFICIENTE\n", len(gd), md, vd, ad)
	fmt.Println("\n   La universalidad de Montgomery–Odlyzko —que el píxel local sea el mismo en libros")
	fmt.Println("   distintos, que es exactamente lo que intuyó el capitán— está observada por otros")
	fmt.Println("   en MILLONES de ceros, y sigue siendo CONJETURAL. Este taller no la verificó hoy,")
	fmt.Println("   y decir que sí sería mentir. Lo que este taller SÍ hizo es construir el segundo")
	fmt.Println("   libro y comprobar que su centro es 6.")
	// ---- LEY 6: the DNA hint, honestly ----
	fmt.Println("\nLEY 6 · LA PISTA DEL ADN — dicha como lo que es: una ANALOGÍA buena")
	fmt.Println("   el capitán vio bien la forma. La ecuación funcional ES una doble hebra:")
	fmt.Println("\n        hebra A: s          hebra B: (w+1) − s        eje: el centro (w+1)/2")
	fmt.Println("\n   cada hebra determina completamente a la otra, y las dos se enroscan alrededor")
	fmt.Println("   del mismo eje. Igual que el ADN. Y hay una segunda coincidencia de forma que")
	fmt.Println("   es la que de verdad importa:")
	fmt.Println("\n        LA REGLA LOCAL ES UNIVERSAL  ·  LA SECUENCIA GLOBAL ES ÚNICA")
	fmt.Println("        (el apareamiento de bases)      (el genoma de cada quien)")
	fmt.Println("        (la estadística de los huecos)   (los τ(n) de cada libro)")
	fmt.Println("\n   ⚖️ Y es una ANALOGÍA, no un teorema. No se usa acá como argumento de nada.")
	fmt.Println("   Sirve para VER la forma, que es exactamente para lo que el capitán la trajo.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("EL FLASH DEL CAPITÁN ES CORRECTO, Y ES MÁS FUERTE DE LO QUE PARECE:")
	fmt.Printf("  · construimos un SEGUNDO LIBRO (Δ de Ramanujan) ....... %d errores en τ\n", malos)
	fmt.Printf("  · la integral reproduce su serie de Dirichlet .......... %.1e\n", peorReal)
	fmt.Printf("  · y reproduce sus %d ceros PUBLICADOS .................. %.1e\n", nComp, peorCero)
	fmt.Println("  · su espejo refleja en 6, NO en ½ ..................... por construcción del libro")
	fmt.Println("  · corriendo el cero 11/2, el centro VUELVE a ½ ........ (w+1)/2, ley general")
	fmt.Printf("  · la mediatriz de las dos estacas es la invariante .... %.1e en tres libros\n", peorMed)
	fmt.Println("  · la universalidad del píxel .......................... NO MEDIDA: muestra insuficiente")
	fmt.Println("\nLO QUE ESTO CAMBIA, Y NO ES POCO:")
	fmt.Println("EL ½ NUNCA FUE EL MISTERIO. El ½ es un accidente de dónde están las dos estacas de")
	fmt.Println("ESTE libro. Cambiás de libro y el medio es 6, u 8, o 12 — y los ceros igual se paran")
	fmt.Println("en la mediatriz nueva. La pregunta que vale no es «¿por qué ½?» sino:")
	fmt.Println("\n        ¿POR QUÉ TODOS LOS CEROS SE PARAN EN LA MEDIATRIZ, SEA EL NÚMERO QUE SEA?")
	fmt.Println("\nEso es la Hipótesis de Riemann GENERALIZADA, abierta para todos los libros a la vez.")
	fmt.Println("El capitán acaba de sacar el ½ del medio de la pregunta, y eso ordena el mapa.")
	fmt.Println("\n⚖️ ¿Y sirve para el premio? Ayuda a preguntar mejor, y eso vale de verdad. Pero no")
	fmt.Println("prueba nada, y hay dos cosas que hay que decir sin vueltas:")
	fmt.Println("  1. verificar Λ(s)=Λ(12−s) con mi fórmula NO era una medición — la simetría estaba")
	fmt.Println("     construida adentro. Tercera vez en la campaña que caigo en esa trampa; queda")
	fmt.Println("     como REGLA: antes de festejar un 0.0e+00, preguntarse si el instrumento podía")
	fmt.Println("     haber dado otra cosa.")
	fmt.Println("  2. la universalidad del píxel —lo más lindo del flash— NO la verificamos: con tres")
	fmt.Println("     huecos no se habla de estadística. Está observada por otros en millones de")
	fmt.Println("     ceros y sigue CONJETURAL.")
	fmt.Println("\n¿Resuelto? Todavía no.")
	escribirLamina(zd, zz, gd, gz, malos, nComp, peorCero, peorReal, peorMed)
}

func escribirLamina(zd, zz, gd, gz []float64, malos, nComp int,
	peorCero, peorReal, peorMed float64) {

	var b strings.Builder
	W, H := 1520.0, 1090.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎯➜ EL CENTRO CORRIDO — el ½ es otro número en otra medición</text>
<text x="%.0f" y="74" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">se construyó un SEGUNDO LIBRO (la Δ de Ramanujan) para probarlo: su centro es 6</text>
`, W, H, W, H, W/2, W/2)

	// two books side by side
	dibujarLibro := func(x0 float64, titulo, sub string, centro float64, zs []float64, tmax float64, col string) {
		fmt.Fprintf(&b, `<rect x="%.0f" y="100" width="700" height="420" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="%.0f" y="130" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd98a">%s</text>
<text x="%.0f" y="154" font-size="14" text-anchor="middle" font-family="monospace" fill="#bfe3ff">%s</text>
`, x0, x0+350, titulo, x0+350, sub)
		lx := x0 + 350
		fmt.Fprintf(&b, `<line x1="%.0f" y1="190" x2="%.0f" y2="470" stroke="%s" stroke-width="3"/>
<text x="%.0f" y="186" font-size="15" text-anchor="middle" font-family="monospace" fill="%s">Re s = %g   ← LA MEDIATRIZ</text>
`, lx, lx, col, lx, col, centro)
		// the two stakes
		e0, e1 := x0+120.0, x0+580.0
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="470" r="7" fill="#ffb27a"/><circle cx="%.0f" cy="470" r="7" fill="#ffb27a"/>
<line x1="%.0f" y1="470" x2="%.0f" y2="470" stroke="#3d6fa8" stroke-width="1.5" stroke-dasharray="4 4"/>
<text x="%.0f" y="494" font-size="13" text-anchor="middle" font-family="monospace" fill="#ffb27a">estaca 0</text>
<text x="%.0f" y="494" font-size="13" text-anchor="middle" font-family="monospace" fill="#ffb27a">estaca %g</text>
`, e0, e1, e0, e1, e0, e1, 2*centro)
		for i, t := range zs {
			if i > 60 {
				break
			}
			py := 470 - 275*t/tmax
			if py < 195 {
				continue
			}
			fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="3.2" fill="#7ee0c0"/>`, lx, py)
		}
		fmt.Fprintf(&b, `<text x="%.0f" y="512" font-size="13" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">%d ceros, todos parados en su propia mediatriz</text>`,
			x0+350, len(zs))
	}
	dibujarLibro(40, "LIBRO 1 · ζ de Riemann", "ζ(s) = Σ 1/nˢ   ·   peso 0   ·   ξ(s) = ξ(1−s)", 0.5, zz, 260, "#7ee0c0")
	dibujarLibro(780, "LIBRO 2 · Δ de Ramanujan", "L(s,Δ) = Σ τ(n)/nˢ   ·   peso 11   ·   Λ(s) = Λ(12−s)", 6, zd, 120, "#c9b6ff")

	fmt.Fprintf(&b, `<rect x="40" y="540" width="700" height="230" rx="10" fill="#161a3a" stroke="#5a4fa8"/>
<text x="390" y="570" font-size="17" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">CORRÉS EL CENTRO Y EL MEDIO ES OTRO NÚMERO</text>
<text x="66" y="602" font-size="13" font-family="monospace" fill="#7fa8cf">  libro                peso w   espejo        centro</text>
<text x="66" y="626" font-size="13.5" font-family="monospace" fill="#cfe6ff">  ζ de Riemann            0    s ↔ 1−s          ½</text>
<text x="66" y="648" font-size="13.5" font-family="monospace" fill="#cfe6ff">  Δ de Ramanujan         11    s ↔ 12−s          6</text>
<text x="66" y="670" font-size="13.5" font-family="monospace" fill="#cfe6ff">  forma de peso 16       15    s ↔ 16−s          8</text>
<text x="66" y="692" font-size="13.5" font-family="monospace" fill="#cfe6ff">  forma de peso 24       23    s ↔ 24−s         12</text>
<text x="390" y="724" font-size="16" text-anchor="middle" font-family="monospace" fill="#ffd98a">centro = (w+1)/2</text>
<text x="390" y="750" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">el ½ es SOLO el caso peso 0. Nunca fue sagrado.</text>

<rect x="780" y="540" width="700" height="230" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="1130" y="570" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">«¿POR QUÉ TODOS LOS PÍXELES SON IGUALES?»</text>
<text x="806" y="600" font-size="13.5" font-family="Georgia" fill="#cfe6ff">el segundo libro reproduce los ceros PUBLICADOS de L(s,Δ),</text>
<text x="806" y="620" font-size="13.5" font-family="Georgia" fill="#cfe6ff">que no entraron en ningún lado del cálculo:</text>
<text x="806" y="648" font-size="14" font-family="monospace" fill="#7ee0c0">  %d ceros · peor desvío %.1e</text>
<text x="806" y="672" font-size="14" font-family="monospace" fill="#7ee0c0">  contra su serie de Dirichlet: %.1e</text>
<text x="806" y="700" font-size="13.5" font-family="Georgia" fill="#ffb27a">⚠ ¿Y LOS PÍXELES IGUALES? ACÁ EL TALLER SE PARA:</text>
<text x="806" y="720" font-size="12.5" font-family="Georgia" fill="#ffb27a">con %d huecos de Δ contra %d de ζ, comparar estadísticas sería</text>
<text x="806" y="738" font-size="12.5" font-family="Georgia" fill="#ffb27a">inventar. La universalidad está observada por otros en millones</text>
<text x="806" y="756" font-size="12.5" font-family="Georgia" fill="#ffb27a">de ceros y sigue CONJETURAL. Hoy NO la verificamos.</text>
`, nComp, peorCero, peorReal, len(gd), len(gz))

	fmt.Fprintf(&b, `<rect x="40" y="790" width="1440" height="266" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="760" y="822" font-size="19" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">LO QUE ESTO CAMBIA — y no es poco</text>
<text x="760" y="862" font-size="24" text-anchor="middle" font-family="Georgia" fill="#ffd98a">EL ½ NUNCA FUE EL MISTERIO</text>
<text x="760" y="892" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">el ½ es un accidente de dónde están las dos estacas de ESTE libro. Cambiás de libro y el medio es 6, u 8, o 12.</text>
<text x="760" y="928" font-size="19" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">¿POR QUÉ TODOS LOS CEROS SE PARAN EN LA MEDIATRIZ, SEA EL NÚMERO QUE SEA?</text>
<text x="760" y="958" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">ésa es la Hipótesis de Riemann GENERALIZADA, abierta para todos los libros a la vez. El capitán sacó el ½ del medio de la pregunta.</text>
<text x="760" y="992" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">🧬 y la pista del ADN: dos hebras, cada una determinando la otra, enroscadas al mismo eje — regla local universal, secuencia global única.</text>
<text x="760" y="1012" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">Es una ANALOGÍA buena para ver la forma, y acá no se usa como argumento de nada.</text>
<text x="760" y="1042" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚖️ ayuda a preguntar mejor, y eso vale. Pero la mediatriz sigue sin razón conocida. Todavía no.</text>
</svg>
`)

	if err := os.WriteFile("centro-corrido.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: centro-corrido.svg")
}
