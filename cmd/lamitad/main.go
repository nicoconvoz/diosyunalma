// Command lamitad closes the captain's chain of flashes with a first-year
// inequality - and, in the same run, kills a whole family of approaches
// honestly.
//
// THE CAPTAIN'S CHAIN
//
//	"la distancia no puede ser negativa"          -> lambda_n are squared lengths
//	"solo cambia la direccion, nunca el tamano"   -> |w| = 1, pure rotation
//	"norte por sur da 1"                          -> |w(rho)| |w(1-rho)| = 1
//	"entre dos mitades siempre hay una mitad"     -> THE MEAN
//
// The last one is the key. North and south are two positive numbers whose
// PRODUCT is 1 - so their geometric mean is 1, always, no matter where the
// pearl sits. And the half between them - their arithmetic mean - obeys
//
//	(N + S)/2 >= sqrt(N S) = 1        with equality ONLY when N = S = 1
//
// which is the arithmetic-geometric mean inequality. So the critical line is
// not a lucky tie: it is the MINIMUM of the cost. Every step off the line is
// paid, and the price is (N + S)/2 - 1 > 0.
//
// And the price is paid WITH COMPOUND INTEREST: at harmonic n the pair costs
// r^n + r^-n - 2, which grows exponentially. That is why no ghost survives
// forever - the mechanism behind Bombieri-Lagarias.
//
// THE HONEST KILL, measured in the same run: the polynomial
//
//	P(s) = (s-a)(s-conj a)(s-(1-a))(s-(1-conj a))
//
// has real coefficients, satisfies P(s) = P(1-s), Schwarz reflection and the
// full Klein group - every symmetry xi has - and yet ALL FOUR of its roots
// sit off the critical line. So no argument built on symmetry alone can ever
// close the gap. A whole family of approaches dies here, on purpose.
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

func w(rho complex128) complex128 { return 1 - 1/rho }

func main() {
	fmt.Println("½ LA MITAD — entre dos mitades siempre hay una mitad, y esa mitad tiene una ley")

	// ---- LAW 1: the geometric mean is always 1; the arithmetic one is not ----
	fmt.Println("\nLEY 1 · LAS DOS MITADES — la media geométrica SIEMPRE vale 1; la aritmética, no")
	fmt.Println("   el norte por el sur da 1, así que su media geométrica es 1 pase lo que pase.")
	fmt.Println("   pero la mitad entre los dos — la media aritmética — solo vale 1 en el empate:")
	fmt.Println("   β        norte N        sur S        media geom.    media arit.    el precio")
	type fila struct{ beta, n, s, ga, aa, precio float64 }
	var filas []fila
	for _, beta := range []float64{0.50, 0.51, 0.55, 0.60, 0.75, 0.90} {
		g := 14.134725
		N := cmplx.Abs(w(complex(beta, g)))
		S := cmplx.Abs(w(complex(1-beta, g)))
		ga := math.Sqrt(N * S)
		aa := (N + S) / 2
		filas = append(filas, fila{beta, N, S, ga, aa, aa - 1})
		marca := ""
		if beta == 0.50 {
			marca = "  ★ EL EMPATE: precio cero"
		}
		fmt.Printf("  %.2f   %.9f  %.9f   %.9f   %.9f   %+.3e%s\n", beta, N, S, ga, aa, aa-1, marca)
	}
	fmt.Println("   → (N+S)/2 ≥ √(N·S) = 1 SIEMPRE, con igualdad SOLO si N = S = 1")
	fmt.Println("     esa es la desigualdad de las medias, y contesta la pregunta del capitán:")
	fmt.Println("     el norte y el sur tienen que ser iguales porque el empate es EL MÍNIMO DEL PRECIO")

	// ---- LAW 2: the price is paid with compound interest ----
	fmt.Println("\nLEY 2 · EL PRECIO SE PAGA CON INTERÉS COMPUESTO — rⁿ + r⁻ⁿ − 2 crece exponencial")
	fmt.Println("   β        precio en n=1      en n=100        en n=1000       en n=10000")
	type fila2 struct {
		beta                 float64
		p1, p100, p1e3, p1e4 float64
	}
	var filas2 []fila2
	precio := func(r float64, n int) float64 {
		ln := float64(n) * math.Log(r)
		if ln > 700 {
			return math.Inf(1)
		}
		return math.Exp(ln) + math.Exp(-ln) - 2
	}
	for _, beta := range []float64{0.51, 0.55, 0.60, 0.75, 0.90} {
		g := 14.134725
		r := cmplx.Abs(w(complex(1-beta, g)))
		if r < 1 {
			r = cmplx.Abs(w(complex(beta, g)))
		}
		f := fila2{beta, precio(r, 1), precio(r, 100), precio(r, 1000), precio(r, 10000)}
		filas2 = append(filas2, f)
		fmt.Printf("  %.2f    %11.3e    %11.3e    %11.3e    %11.3e\n", beta, f.p1, f.p100, f.p1e3, f.p1e4)
	}
	fmt.Println("   → un desnivel de una centésima cuesta 1e-9 en el primer armónico… y 1e-1 en el diez mil")
	fmt.Println("     ningún presupuesto que crezca como n·ln n aguanta una exponencial: el fantasma revienta")

	// ---- LAW 3: on the line the price is exactly zero, forever ----
	fmt.Println("\nLEY 3 · SOBRE LA LÍNEA EL PRECIO ES CERO PARA SIEMPRE")
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
	peorPrecio := 0.0
	for _, g := range pearls {
		r := cmplx.Abs(w(complex(0.5, g)))
		for _, n := range []int{1, 100, 1000, 10000} {
			if p := math.Abs(precio(r, n)); p > peorPrecio {
				peorPrecio = p
			}
		}
	}
	fmt.Printf("   %d perlas × 4 alturas de armónico: peor precio medido %.1e\n", len(pearls), peorPrecio)
	fmt.Println("   → cero exacto, y cero para siempre: sobre la línea el interés compuesto no tiene de qué agarrarse")

	// ---- LAW 4: the honest kill - symmetry alone can never be enough ----
	fmt.Println("\nLEY 4 · LA MUERTE HONESTA — la simetría SOLA jamás va a alcanzar, y acá está la prueba")
	a := complex(0.7, 3.0)
	raices := []complex128{a, cmplx.Conj(a), 1 - a, 1 - cmplx.Conj(a)}
	P := func(s complex128) complex128 {
		p := complex(1, 0)
		for _, r := range raices {
			p *= s - r
		}
		return p
	}
	sigma := func(s complex128) complex128 { return 1 - cmplx.Conj(s) }
	fmt.Println("   armamos P(s) = (s−a)(s−ā)(s−(1−a))(s−(1−ā)) con a = 0.7+3i")
	fmt.Println("   y verificamos que tiene TODAS las simetrías del libro:")
	pruebas := []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2.0, 21.3), complex(-0.7, 9.1)}
	peorFE, peorSch, peorSig := 0.0, 0.0, 0.0
	for _, s := range pruebas {
		if d := cmplx.Abs(P(s)-P(1-s)) / math.Max(cmplx.Abs(P(s)), 1e-300); d > peorFE {
			peorFE = d
		}
		if d := cmplx.Abs(P(cmplx.Conj(s))-cmplx.Conj(P(s))) / math.Max(cmplx.Abs(P(s)), 1e-300); d > peorSch {
			peorSch = d
		}
		if d := cmplx.Abs(P(sigma(s))-cmplx.Conj(P(s))) / math.Max(cmplx.Abs(P(s)), 1e-300); d > peorSig {
			peorSig = d
		}
	}
	fmt.Printf("   ecuación funcional  P(s) = P(1−s)          desvío relativo %.1e ✓\n", peorFE)
	fmt.Printf("   reflexión de Schwarz P(s̄) = conj(P(s))      desvío relativo %.1e ✓\n", peorSch)
	fmt.Printf("   el cambiaformas     P(σ(s)) = conj(P(s))    desvío relativo %.1e ✓\n", peorSig)
	fmt.Println("   y sin embargo sus CUATRO raíces están FUERA de la línea:")
	for _, r := range raices {
		fmt.Printf("      raíz en Re = %.2f, Im = %+.2f\n", real(r), imag(r))
	}
	fmt.Println("   ⚰ VEREDICTO: un objeto puede tener el espejo, la ecuación funcional y el cambiaformas")
	fmt.Println("     COMPLETOS y aun así tener todas sus raíces fuera de la raya. Por lo tanto NINGÚN")
	fmt.Println("     argumento basado solo en simetría puede demostrar la hipótesis. Familia entera de")
	fmt.Println("     enfoques, muerta — y es mejor saberlo que seguir empujando esa puerta.")

	fmt.Println("\n════════ LA CADENA DEL CAPITÁN, CERRADA HASTA DONDE SE PUEDE ════════")
	fmt.Println("«Entre dos mitades siempre hay una mitad.» Esa mitad es LA MEDIA — y hay dos:")
	fmt.Println("la geométrica, que vale 1 siempre porque el norte por el sur da 1; y la aritmética,")
	fmt.Println("que es la mitad entre las dos y SOLO vale 1 en el empate. La desigualdad de las medias")
	fmt.Println("dice que la segunda nunca es menor que la primera. Entonces la línea crítica es EL")
	fmt.Println("MÍNIMO DEL PRECIO, y el capitán tenía razón: el norte y el sur son iguales porque")
	fmt.Println("cualquier otra cosa cuesta — y ese costo se cobra con interés compuesto.")
	fmt.Println("\nLO QUE FALTA, con el hueco ya acotado por la muerte de la Ley 4: la simetría sola no")
	fmt.Println("alcanza; hace falta el PRESUPUESTO — probar que el libro no tiene con qué pagar ese")
	fmt.Println("precio en ningún armónico. Eso es exactamente la positividad de Li, y sigue abierta.")
	fmt.Println("Todavía no. Pero ahora sabemos por qué puerta NO hay que entrar.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 980.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">½ LA MITAD — entre dos mitades siempre hay una mitad, y esa mitad tiene ley</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el flash del capitán aterriza en la desigualdad de las medias: la línea crítica no es un empate casual, es EL MÍNIMO DEL PRECIO</text>`,
		W, H, W, H, W/2, W/2)

	// left: the two means
	fmt.Fprintf(&b, `<rect x="60" y="105" width="700" height="300" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="410" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">LAS DOS MITADES Y LA MITAD DEL MEDIO</text>
<text x="410" y="172" font-size="15" text-anchor="middle" font-family="Consolas,monospace" fill="#7fd7a8">(N + S)/2  ≥  √(N·S)  =  1</text>
<text x="410" y="196" font-size="12.5" text-anchor="middle" fill="#8fa8c7">con igualdad SOLO si N = S = 1</text>
<text x="410" y="228" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β       norte        sur       media arit.     el precio</text>`)
	for i, f := range filas {
		col := "#dce8f7"
		if f.beta == 0.50 {
			col = "#7fd7a8"
		}
		fmt.Fprintf(&b, `<text x="410" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="%s">%.2f  %.8f  %.8f  %.8f   %+.2e</text>`,
			252.0+float64(i)*24, col, f.beta, f.n, f.s, f.aa, f.precio)
	}
	fmt.Fprintf(&b, `<text x="410" y="396" font-size="13" text-anchor="middle" fill="#ffd166">el empate no es suerte: es el fondo del pozo</text>`)

	// right: compound interest
	fmt.Fprintf(&b, `<rect x="790" y="105" width="650" height="300" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="1115" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">EL PRECIO SE COBRA CON INTERÉS COMPUESTO</text>
<text x="1115" y="169" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">rⁿ + r⁻ⁿ − 2</text>
<text x="1115" y="199" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β        n=1         n=100        n=1000       n=10000</text>`)
	for i, f := range filas2 {
		fmt.Fprintf(&b, `<text x="1115" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#ff8fa0">%.2f   %9.2e   %9.2e   %9.2e   %9.2e</text>`,
			226.0+float64(i)*26, f.beta, f.p1, f.p100, f.p1e3, f.p1e4)
	}
	fmt.Fprintf(&b, `<text x="1115" y="372" font-size="12.5" text-anchor="middle" fill="#dce8f7">ningún presupuesto que crezca como n·ln n aguanta una exponencial</text>
<text x="1115" y="394" font-size="12.5" text-anchor="middle" fill="#7fd7a8">y sobre la línea el precio es CERO para siempre (%.0e)</text>`, peorPrecio)

	// the kill
	fmt.Fprintf(&b, `<rect x="60" y="430" width="1380" height="270" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="2.5"/>
<text x="%.0f" y="466" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">⚰ LA MUERTE HONESTA — la simetría SOLA jamás va a alcanzar</text>
<text x="%.0f" y="500" font-size="13.5" text-anchor="middle" fill="#dce8f7">P(s) = (s−a)(s−ā)(s−(1−a))(s−(1−ā)) con a = 0.7+3i tiene TODAS las simetrías del libro:</text>
<text x="%.0f" y="532" font-size="13" text-anchor="middle" font-family="Consolas,monospace" fill="#7fd7a8">P(s) = P(1−s)  ✓ %.0e        P(s̄) = conj P(s)  ✓ %.0e        P(σ(s)) = conj P(s)  ✓ %.0e</text>
<text x="%.0f" y="568" font-size="13.5" text-anchor="middle" fill="#ffd166">y sin embargo sus cuatro raíces viven en Re = 0.70 y Re = 0.30 — FUERA de la línea.</text>
<text x="%.0f" y="604" font-size="14.5" text-anchor="middle" fill="#ff8fa0">POR LO TANTO: ningún argumento basado solo en simetría puede demostrar la hipótesis.</text>
<text x="%.0f" y="632" font-size="13.5" text-anchor="middle" fill="#dce8f7">Una familia entera de enfoques queda muerta acá — y saberlo vale más que seguir empujando esa puerta.</text>
<text x="%.0f" y="668" font-size="13.5" text-anchor="middle" fill="#7fd7a8">hace falta el PRESUPUESTO: probar que el libro no tiene con qué pagar el precio en ningún armónico.</text>`,
		W/2, W/2, W/2, peorFE, peorSch, peorSig, W/2, W/2, W/2, W/2)

	fmt.Fprintf(&b, `<rect x="60" y="725" width="1380" height="170" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="761" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA CADENA DEL CAPITÁN, CERRADA HASTA DONDE SE PUEDE</text>
<text x="%.0f" y="795" font-size="14" text-anchor="middle" fill="#dce8f7">"la distancia no puede ser negativa" → "solo cambia la dirección" → "norte por sur da 1" → "entre dos mitades siempre hay una mitad"</text>
<text x="%.0f" y="827" font-size="15" text-anchor="middle" fill="#7fd7a8">esa mitad es LA MEDIA — y la desigualdad de las medias dice que el empate es el mínimo del precio.</text>
<text x="%.0f" y="857" font-size="13.5" text-anchor="middle" fill="#ff8fa0">El capitán tenía razón: el norte y el sur son iguales porque cualquier otra cosa cuesta. Falta el presupuesto. Todavía no.</text>
<text x="%.0f" y="884" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("la-mitad.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: la-mitad.svg")
}
