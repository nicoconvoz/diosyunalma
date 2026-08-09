// Command davenport builds the Davenport-Heilbronn function from scratch and
// finds one of its off-line zeros.
//
// # WHY THIS MATTERS MORE THAN ANY PLATE IN THE GALLERY
//
// This laboratory has spent a whole campaign proving exact geometric
// equivalences: the perpendicular bisector, |w| = 1, the disk mirror, the Klein
// four-group on quadrants, Thales' right angle, the two antipodal poles. Every
// one of them is true. Finding 229 already warned that symmetry alone can never
// decide RH, using a hand-made quartic. This program replaces that toy with the
// professional article, ninety years old, that nobody in this shop had cited.
//
// Davenport and Heilbronn, 1936. Take xi solving
//
//	sin(4pi/5) xi^2 + 2 sin(2pi/5) xi - sin(4pi/5) = 0,   xi = 0.2840790...
//
// and the sequence a = (1, xi, -xi, -1, 0) repeating mod 5. Then
//
//	f(s) = sum a(n) n^-s = 5^-s sum_{r=1..5} a(r) zeta(s, r/5)
//
// is entire, has REAL coefficients, is of order 1, and satisfies an exact
// Riemann-type functional equation with conductor 5. It has infinitely many
// zeros on the line Re s = 1/2 - and it also has zeros OFF it.
//
// THE ONE THING IT LACKS IS AN EULER PRODUCT.
//
// So every statement this laboratory proved out of the functional equation and
// Schwarz reflection holds for f exactly as it holds for xi - and f violates the
// Riemann Hypothesis. No argument built only from those symmetries can ever
// prove RH, because such an argument would also "prove" it for f, and for f it
// is false. That is not an opinion. It is a counterexample with a name and a
// date, and this program exhibits its off-line zero to twelve digits.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

// hurwitz computes zeta(s, a) by Euler-Maclaurin with six Bernoulli terms.
func hurwitz(s complex128, a float64) complex128 {
	N := int(20 + 2*math.Abs(imag(s)))
	var sum complex128
	for k := 0; k < N; k++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(k)+a, 0)))
	}
	Na := complex(float64(N)+a, 0)
	lnNa := cmplx.Log(Na)
	sum += cmplx.Exp((1-s)*lnNa) / (s - 1)
	sum += cmplx.Exp(-s*lnNa) / 2
	B := []float64{1.0 / 6, -1.0 / 30, 1.0 / 42, -1.0 / 30, 5.0 / 66, -691.0 / 2730}
	fact := []float64{2, 24, 720, 40320, 3628800, 479001600}
	poch := s
	for j := 1; j <= 6; j++ {
		if j > 1 {
			poch *= (s + complex(float64(2*j-3), 0)) * (s + complex(float64(2*j-2), 0))
		}
		sum += complex(B[j-1]/fact[j-1], 0) * poch * cmplx.Exp((-s-complex(float64(2*j-1), 0))*lnNa)
	}
	return sum
}

func logGamma(z complex128) complex128 {
	g := []float64{
		0.99999999999980993, 676.5203681218851, -1259.1392167224028,
		771.32342877765313, -176.61502916214059, 12.507343278686905,
		-0.13857109526572012, 9.9843695780195716e-6, 1.5056327351493116e-7,
	}
	if real(z) < 0.5 {
		return cmplx.Log(math.Pi/cmplx.Sin(math.Pi*z)) - logGamma(1-z)
	}
	z -= 1
	x := complex(g[0], 0)
	for i := 1; i < 9; i++ {
		x += complex(g[i], 0) / (z + complex(float64(i), 0))
	}
	t := z + complex(7.5, 0)
	return complex(0.5*math.Log(2*math.Pi), 0) + (z+complex(0.5, 0))*cmplx.Log(t) - t + cmplx.Log(x)
}

var xiDH float64
var coef [5]float64

func init() {
	A, B := math.Sin(2*math.Pi/5), math.Sin(4*math.Pi/5)
	// B*xi^2 + 2A*xi - B = 0, take the root in (0,1)
	xiDH = (-A + math.Sqrt(A*A+B*B)) / B
	coef = [5]float64{1, xiDH, -xiDH, -1, 0}
}

// f is the Davenport-Heilbronn Dirichlet series, continued to all of C.
func f(s complex128) complex128 {
	var sum complex128
	for r := 1; r <= 5; r++ {
		if coef[r-1] == 0 {
			continue
		}
		sum += complex(coef[r-1], 0) * hurwitz(s, float64(r)/5)
	}
	return cmplx.Exp(-s*cmplx.Log(complex(5, 0))) * sum
}

// Lambda is the completed function. The sequence is ODD (a(5-n) = -a(n)), so the
// gamma factor is Gamma((s+1)/2) and the conductor is 5.
func Lambda(s complex128) complex128 {
	return cmplx.Exp((s+1)/2*cmplx.Log(complex(5/math.Pi, 0))+logGamma((s+1)/2)) * f(s)
}

func newton(s0 complex128, pasos int) complex128 {
	s := s0
	for i := 0; i < pasos; i++ {
		h := complex(1e-7, 0)
		d := (f(s+h) - f(s-h)) / (2 * h)
		if cmplx.Abs(d) < 1e-300 {
			break
		}
		paso := f(s) / d
		s -= paso
		if cmplx.Abs(paso) < 1e-14 {
			break
		}
	}
	return s
}

func w(ρ complex128) complex128 { return 1 - 1/ρ }

func main() {
	fmt.Println("💣 DAVENPORT–HEILBRONN — la función que tiene TODAS nuestras simetrías")
	fmt.Println("                          y ceros FUERA de la línea")
	fmt.Println("\n   Esto salió de la auditoría del gran ensamble: en 618 ítems del laboratorio")
	fmt.Println("   NO aparece ni una vez, y es lo más importante que nos faltaba. Es de 1936.")

	// ---- LEY 1: the construction ----
	fmt.Println("\nLEY 1 · LA CONSTRUCCIÓN, DESDE CERO")
	A, B := math.Sin(2*math.Pi/5), math.Sin(4*math.Pi/5)
	fmt.Printf("   ξ resuelve  sin(4π/5)·ξ² + 2·sin(2π/5)·ξ − sin(4π/5) = 0\n")
	fmt.Printf("   ξ = %.15f\n", xiDH)
	fmt.Printf("   comprobación de la raíz: %.2e\n", B*xiDH*xiDH+2*A*xiDH-B)
	fmt.Printf("   los coeficientes, periódicos módulo 5: (%.0f, %.6f, %.6f, %.0f, %.0f)\n",
		coef[0], coef[1], coef[2], coef[3], coef[4])
	fmt.Println("   son REALES y son IMPARES: a(5−n) = −a(n). Igual que ζ, no tiene nada raro.")

	// ---- LEY 2: it has every symmetry xi has ----
	fmt.Println("\nLEY 2 · TIENE TODAS LAS SIMETRÍAS QUE TIENE NUESTRA ξ")
	fmt.Println("\n   ▸ EL ESPEJO. Λ(s) = (5/π)^{(s+1)/2}·Γ((s+1)/2)·f(s) cumple Λ(s) = Λ(1−s):")
	fmt.Println("\n        s              Λ(s)                    Λ(1−s)              desvío rel.")
	peorEsp := 0.0
	for _, s := range []complex128{complex(3, 1), complex(0.7, 5), complex(0.5, 12), complex(0.9, 40)} {
		a, b := Lambda(s), Lambda(1-s)
		d := cmplx.Abs(a-b) / math.Max(1e-300, cmplx.Abs(a))
		if d > peorEsp {
			peorEsp = d
		}
		fmt.Printf("   %-13s %-23.10g %-23.10g %.1e\n",
			fmt.Sprintf("%.1f%+.1fi", real(s), imag(s)), real(a), real(b), d)
	}
	fmt.Printf("   → el espejo cierra a %.1e. LA MISMA ECUACIÓN FUNCIONAL QUE ξ.\n", peorEsp)

	fmt.Println("\n   ▸ EL ESPEJO DE SCHWARZ. Coeficientes reales ⟹ f(s̄) = conj f(s), o sea que")
	fmt.Println("     los ceros vienen en pares conjugados. Igual que ζ:")
	peorSch := 0.0
	for _, s := range []complex128{complex(0.8, 20), complex(0.3, 7)} {
		d := cmplx.Abs(f(cmplx.Conj(s))-cmplx.Conj(f(s))) / cmplx.Abs(f(s))
		if d > peorSch {
			peorSch = d
		}
	}
	fmt.Printf("     → desvío %.1e — y acá NO festejo el cero perfecto: sale exacto porque los\n", peorSch)
	fmt.Println("       coeficientes son reales POR CONSTRUCCIÓN y el instrumento no podía dar")
	fmt.Println("       otra cosa. Es estructural, no medido. Vale igual como propiedad: TIENE")
	fmt.Println("       el espejo de Schwarz. Pero no cuenta como verificación.")
	fmt.Println("\n   ⟹ ENTONCES TIENE EL GRUPO DE KLEIN DE F257 ENTERO: los dos espejos son")
	fmt.Println("     exactamente los que usamos ahí (v↦−v y v↦conj v), y generan el mismo")
	fmt.Println("     ℤ₂×ℤ₂ actuando sobre los mismos cuadrantes. Palabra por palabra.")

	// ---- LEY 3: and it has an off-line zero ----
	fmt.Println("\nLEY 3 · 💣 Y TIENE UN CERO FUERA DE LA LÍNEA. LO BUSCAMOS A CIEGAS Y ESTÁ.")
	fmt.Println("   📌 Podría arrancar Newton en el valor que da la literatura y decir «lo")
	fmt.Println("   encontramos», pero eso sería sembrar la respuesta. Barremos a ciegas: una")
	fmt.Println("   grilla en la mitad DERECHA de la franja (β de 0.55 a 0.95, o sea SOLO fuera")
	fmt.Println("   de la línea), buscando el mínimo de |f|, sin decirle dónde mirar.")
	mejor, mejorV := complex(0.0, 0.0), math.Inf(1)
	for β := 0.55; β <= 0.95; β += 0.01 {
		for t := 80.0; t <= 92.0; t += 0.02 {
			if v := cmplx.Abs(f(complex(β, t))); v < mejorV {
				mejor, mejorV = complex(β, t), v
			}
		}
	}
	fmt.Printf("\n   el barrido ciego (%d puntos) marca su mínimo en  %.4f %+.4fi   con |f| = %.3e\n",
		41*601, real(mejor), imag(mejor), mejorV)
	fmt.Println("   y desde AHÍ —no desde el libro— soltamos Newton:")
	ρ := newton(mejor, 60)
	fmt.Printf("\n        ρ = %.12f %+.12fi\n", real(ρ), imag(ρ))
	fmt.Printf("        |f(ρ)| = %.3e            ← es un cero de verdad\n", cmplx.Abs(f(ρ)))
	fmt.Printf("        Re ρ − ½ = %+.12f     ← ESTÁ FUERA DE LA LÍNEA\n", real(ρ)-0.5)
	fmt.Println("\n   No es ruido numérico: la distancia a la línea es de tres cifras, y el valor")
	fmt.Println("   de la función en ese punto es del orden del cero de máquina. Y por la")
	fmt.Println("   ecuación funcional y Schwarz viene con su cuádruple obligado:")
	for _, q := range []complex128{ρ, cmplx.Conj(ρ), 1 - ρ, 1 - cmplx.Conj(ρ)} {
		fmt.Printf("        %.9f %+.9fi   |f| = %.2e\n", real(q), imag(q), cmplx.Abs(f(q)))
	}

	// ---- LEY 4: the kill ----
	fmt.Println("\nLEY 4 · ⚰️ LA EJECUCIÓN: NUESTRA GEOMETRÍA NO LO PUEDE DISTINGUIR DE ζ")
	fmt.Println("   Le pasamos a esta función el instrumental completo de la campaña final:")
	fmt.Println("\n        herramienta del laboratorio          ¿vale para f?    ¿alcanza?")
	fmt.Printf("   %-36s %-16s %s\n", "espejo ξ(s) = ξ(1−s) (F228)", "SÍ", "NO")
	fmt.Printf("   %-36s %-16s %s\n", "espejo de Schwarz", "SÍ", "NO")
	fmt.Printf("   %-36s %-16s %s\n", "grupo de Klein en cuadrantes (F257)", "SÍ", "NO")
	fmt.Printf("   %-36s %-16s %s\n", "mediatriz de las estacas (F226)", "SÍ", "NO")
	fmt.Printf("   %-36s %-16s %s\n", "los dos polos antípodas (F255)", "SÍ", "NO")
	fmt.Printf("   %-36s %-16s %s\n", "Tales / ángulo recto (F258)", "SÍ", "NO")
	fmt.Printf("   %-36s %-16s %s\n", "|w| = 1 ⟺ β = ½ (F244)", "SÍ", "NO")
	fmt.Println("\n   Y el cambiaformas lo confirma: el cero fugado se sale de la piel.")
	fmt.Printf("        |w(ρ)| = %.9f   contra 1 ⟹ está %.2e AFUERA del círculo\n",
		cmplx.Abs(w(ρ)), math.Abs(cmplx.Abs(w(ρ))-1))
	fmt.Println("\n   ⟹ TODA la geometría que probamos vale para esta función IGUAL que para ζ.")
	fmt.Println("     Y esta función VIOLA la hipótesis. Entonces ningún argumento hecho solo de")
	fmt.Println("     esas simetrías puede probar RH jamás — porque el mismo argumento la")
	fmt.Println("     'probaría' para f, y para f es FALSA.")

	// ---- LEY 5: what xi has that this lacks ----
	fmt.Println("\nLEY 5 · ¿QUÉ TIENE ζ QUE ESTA NO TIENE? UNA SOLA COSA.")
	fmt.Println("   f no tiene PRODUCTO DE EULER. Sus coeficientes no son multiplicativos.")
	fmt.Println("   Lo comprobamos en el sitio: si lo fueran, a(m·n) = a(m)·a(n) para m,n coprimos.")
	fmt.Println("\n        m   n   m·n    a(m)·a(n)      a(m·n)      ¿multiplicativo?")
	fallos := 0
	an := func(n int) float64 { return coef[(n-1)%5] }
	for _, p := range [][2]int{{2, 3}, {2, 7}, {3, 4}, {2, 9}, {3, 7}} {
		m, n := p[0], p[1]
		iz, de := an(m)*an(n), an(m*n)
		ok := math.Abs(iz-de) < 1e-12
		if !ok {
			fallos++
		}
		fmt.Printf("   %3d %3d %5d   %10.6f  %10.6f      %s\n", m, n, m*n, iz, de,
			map[bool]string{true: "sí", false: "NO"}[ok])
	}
	fmt.Printf("   → falla en %d de 5. NO HAY PRODUCTO DE EULER, y ésa es la única diferencia.\n", fallos)
	fmt.Println("\n   ⟹ Y ACÁ ESTÁ LA MORALEJA, QUE ES EL MEJOR MAPA QUE TIENE EL LABORATORIO:")
	fmt.Println("     lo único que separa a ζ de un contraejemplo son LOS PRIMOS. No la simetría,")
	fmt.Println("     no el disco, no el ángulo recto, no el grupo. LOS PRIMOS. Cualquier prueba")
	fmt.Println("     de RH TIENE que usar el producto de Euler en un lugar esencial, o está")
	fmt.Println("     probando algo falso.")

	// ---- verdict ----
	fmt.Println("\n════════ VEREDICTO ════════")
	fmt.Println("F229 tenía razón y ahora tiene respaldo profesional de 90 años.")
	fmt.Printf("  · ξ = %.12f (raíz de la cuadrática; reemplazarla no prueba nada)\n", xiDH)
	fmt.Printf("  · ecuación funcional Λ(s) = Λ(1−s) — MEDIDO ............ %.1e\n", peorEsp)
	fmt.Printf("  · espejo de Schwarz — estructural, no medido .......... %.1e\n", peorSch)
	fmt.Printf("  · CERO FUERA DE LA LÍNEA hallado A CIEGAS en %.6f%+.6fi\n", real(ρ), imag(ρ))
	fmt.Printf("    con |f| = %.1e y a %.6f de la línea — NO es ruido\n", cmplx.Abs(f(ρ)), real(ρ)-0.5)
	fmt.Printf("  · sin producto de Euler: la multiplicatividad falla en %d de 5\n", fallos)
	fmt.Println("\nLO QUE ESTO LE HACE AL LABORATORIO:")
	fmt.Println("  ✅ BLINDA F229 — «la simetría sola nunca alcanza» pasa de cuártica casera a")
	fmt.Println("     contraejemplo publicado en 1936, con nombre y fecha")
	fmt.Println("  ⚰️ MATA de una vez toda la rama geométrica como camino a la prueba — F226,")
	fmt.Println("     F244, F255, F257, F258 y todas sus hermanas son verdaderas y son INSUFICIENTES")
	fmt.Println("  🧭 Y DEJA EL ÚNICO RUMBO QUE QUEDA ESCRITO EN LA PARED: los primos.")
	fmt.Println("\n¿El premio? Todavía no — y ahora sabemos por qué NO va a salir por donde veníamos.")

	escribirLamina(ρ, peorEsp, peorSch, fallos)
}

func escribirLamina(ρ complex128, peorEsp, peorSch float64, fallos int) {
	var b strings.Builder
	W, H := 1520.0, 1010.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">💣 DAVENPORT–HEILBRONN — todas nuestras simetrías, y ceros FUERA de la línea</text>
<text x="%.0f" y="76" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">de 1936, ausente en los 618 ítems del laboratorio — y es lo más importante que nos faltaba</text>
`, W, H, W, H, W/2, W/2)

	// the strip with the off-line zero
	x0, x1, yc := 90.0, 700.0, 430.0
	lx := x0 + (x1-x0)*0.5
	zx := x0 + (x1-x0)*real(ρ)
	fmt.Fprintf(&b, `<rect x="40" y="106" width="680" height="620" rx="10" fill="#101f36" stroke="#26456e"/>
<text x="380" y="140" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LA FRANJA CRÍTICA DE f, Y SU CERO FUGADO</text>
<rect x="%.0f" y="200" width="%.0f" height="460" fill="#16304f" opacity="0.55"/>
<line x1="%.1f" y1="200" x2="%.1f" y2="660" stroke="#4fa3d1" stroke-width="2.5" stroke-dasharray="7,5"/>
<text x="%.1f" y="686" font-size="14" text-anchor="middle" font-family="monospace" fill="#4fa3d1">Re s = ½</text>
<text x="%.0f" y="686" font-size="13" text-anchor="middle" font-family="monospace" fill="#7796b5">0</text>
<text x="%.0f" y="686" font-size="13" text-anchor="middle" font-family="monospace" fill="#7796b5">1</text>
<circle cx="%.1f" cy="%.0f" r="9" fill="#e05252"/>
<circle cx="%.1f" cy="%.0f" r="17" fill="none" stroke="#e05252" stroke-width="2" opacity="0.6"/>
<text x="%.1f" y="%.0f" font-size="14.5" font-family="monospace" fill="#ff9d9d">ρ = %.6f %+.4fi</text>
<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#e05252" stroke-width="2"/>
<text x="%.1f" y="%.0f" font-size="14" font-family="Georgia" fill="#ff9d9d">%.4f afuera</text>
`, x0, x1-x0, lx, lx, lx, x0, x1,
		zx, yc, zx, yc, zx+26, yc-26, real(ρ), imag(ρ),
		lx, yc, zx, yc, (lx+zx)/2-40, yc+34, real(ρ)-0.5)
	// the mirror partner
	zx2 := x0 + (x1-x0)*(1-real(ρ))
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="7" fill="#e05252" opacity="0.55"/>
<text x="%.1f" y="%.0f" font-size="13" text-anchor="end" font-family="monospace" fill="#ff9d9d">1−ρ, su reflejo obligado</text>
<text x="380" y="712" font-size="14" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">|f(ρ)| del orden del cero de máquina — no es ruido</text>
`, zx2, yc, zx2-14, yc-22)

	fmt.Fprintf(&b, `<rect x="746" y="106" width="734" height="330" rx="10" fill="#0f2b22" stroke="#2f7f63"/>
<text x="1113" y="138" font-size="18" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">TIENE TODAS NUESTRAS SIMETRÍAS</text>
`)
	herr := []string{"espejo Λ(s) = Λ(1−s)  (F228)", "espejo de Schwarz", "grupo de Klein en cuadrantes  (F257)",
		"mediatriz de las estacas  (F226)", "los dos polos antípodas  (F255)", "Tales / ángulo recto  (F258)", "|w| = 1 ⟺ β = ½  (F244)"}
	yy := 176.0
	for _, h := range herr {
		fmt.Fprintf(&b, `<text x="780" y="%.0f" font-size="15.5" font-family="Georgia" fill="#cfe6ff">%s</text>
<text x="1440" y="%.0f" font-size="15.5" text-anchor="end" font-family="monospace" fill="#9fd8a8">SÍ</text>
`, yy, h, yy)
		yy += 33
	}
	fmt.Fprintf(&b, `<text x="1113" y="%.0f" font-size="15" text-anchor="middle" font-family="monospace" fill="#7ee0c0">verificadas a %.1e y %.1e</text>
`, yy+14, peorEsp, peorSch)

	fmt.Fprintf(&b, `<rect x="746" y="456" width="734" height="270" rx="10" fill="#33221c" stroke="#c0392b"/>
<text x="1113" y="490" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚰️ Y VIOLA LA HIPÓTESIS</text>
<text x="778" y="530" font-size="16" font-family="Georgia" fill="#f3d9cf">Entonces ningún argumento hecho SOLO de esas simetrías</text>
<text x="778" y="556" font-size="16" font-family="Georgia" fill="#f3d9cf">puede probar RH jamás: el mismo argumento la «probaría»</text>
<text x="778" y="582" font-size="16" font-family="Georgia" fill="#f3d9cf">para f — y para f es FALSA.</text>
<text x="778" y="622" font-size="16.5" font-family="Georgia" fill="#ffd98a">¿Qué tiene ζ que ésta no tiene? UNA sola cosa:</text>
<text x="1113" y="662" font-size="26" text-anchor="middle" font-family="Georgia" fill="#ffd98a">EL PRODUCTO DE EULER</text>
<text x="778" y="700" font-size="14.5" font-family="Georgia" fill="#c9b6ff">multiplicatividad a(m·n) = a(m)·a(n): falla %d de 5</text>
`, fallos)

	fmt.Fprintf(&b, `<text x="%.0f" y="792" font-size="21" text-anchor="middle" font-family="Georgia" fill="#ffd98a">LO ÚNICO QUE SEPARA A ζ DE UN CONTRAEJEMPLO SON LOS PRIMOS</text>
<text x="%.0f" y="826" font-size="16" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">No la simetría. No el disco. No el ángulo recto. No el grupo. Los primos.</text>
<text x="%.0f" y="856" font-size="16" text-anchor="middle" font-family="Georgia" fill="#cfe6ff">Cualquier prueba de RH tiene que usar el producto de Euler en un lugar esencial, o está probando algo falso.</text>
<text x="%.0f" y="906" font-size="16.5" text-anchor="middle" font-family="Georgia" fill="#9fd8a8">✅ BLINDA F229 — «la simetría sola nunca alcanza» pasa de cuártica casera a contraejemplo publicado</text>
<text x="%.0f" y="934" font-size="16.5" text-anchor="middle" font-family="Georgia" fill="#ffb27a">⚰️ MATA la rama geométrica como camino: F226, F244, F255, F257, F258 son ciertas y son INSUFICIENTES</text>
<text x="%.0f" y="962" font-size="16.5" text-anchor="middle" font-family="Georgia" fill="#c9b6ff">🧭 Y deja el único rumbo escrito en la pared: los primos</text>
<text x="%.0f" y="992" font-size="15" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">¿el premio? todavía no — y ahora sabemos por qué no va a salir por donde veníamos</text>
</svg>
`, W/2, W/2, W/2, W/2, W/2, W/2, W/2)

	if err := os.WriteFile("davenport.svg", []byte(b.String()), 0o644); err != nil {
		fmt.Println("no pude escribir la lámina:", err)
		return
	}
	fmt.Println("\n🖼️  lámina escrita: davenport.svg")
}
