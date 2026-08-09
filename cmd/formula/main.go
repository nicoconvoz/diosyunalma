// Command formula assembles the captain's whole chain of flashes into ONE
// expression in the shapeshifter coordinate, harmonized at dimension 0 - the
// projection of the infinite into a single shape, the way a cube is drawn on
// a flat sheet.
//
// THE ASSEMBLY. Pair each pearl rho with its conjugate and write w = 1 - 1/rho,
// r = |w|. The pair's contribution to the n-th envelope is 2 - 2 Re(w^n), and
// that splits EXACTLY into two named pieces:
//
//	2 - 2 Re(w^n)  =  |1 - w^n|^2  +  (1 - r^(2n))
//	                  \_________/     \__________/
//	                   LA FORMA         LA FUGA
//	                  una distancia    el cambio de
//	                  al cuadrado:     TAMANO: vale
//	                  jamas negativa   0 si y solo si
//	                                   r = 1
//
// So the whole mold is:
//
//	lambda_n = SUM over pairs [ |1 - w^n|^2 + (1 - r^(2n)) ]
//
// and every flash of the campaign is sitting in that one line: the distance
// that cannot be negative (F224), the step that changes direction and never
// size (F225), north times south equal to one (F225), the mean (F229).
//
// PROJECTED ALL AT ONCE. Multiplying by z^n and summing over every harmonic
// packs the infinitely many envelopes into a single function of z - the germ
// read on the disk around the clasp, dimension 0:
//
//	L(z) = SUM lambda_n z^n = SUM over pairs [ FORMA(z) + FUGA(z) ]
//	FORMA(z) = z/(1-z) - 2 Re[ wz/(1-wz) ] + r^2 z/(1-r^2 z)
//	FUGA(z)  = z (1 - r^2) / [ (1-z)(1 - r^2 z) ]
//
// FUGA(z) is identically zero exactly when r = 1. One formula, one variable,
// and the entire hypothesis reduced to whether one term is present or absent.
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

func theta(t float64) float64 {
	t2 := t * t
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t) + 7/(5760*t*t2) + 31/(80640*t*t2*t2)
}

func zOf(t float64) float64 {
	return real(cmplx.Exp(complex(0, theta(t))) * zetaC(complex(0.5, t)))
}

func wDe(rho complex128) complex128 { return 1 - 1/rho }

// forma is the shape term: the generating function of the squared distances.
func forma(w complex128, z float64) float64 {
	r2 := real(w)*real(w) + imag(w)*imag(w)
	zc := complex(z, 0)
	t1 := z / (1 - z)
	t2 := -2 * real(w*zc/(1-w*zc))
	t3 := r2 * z / (1 - r2*z)
	return t1 + t2 + t3
}

// fuga is the leak term: it is identically zero exactly when |w| = 1.
func fuga(w complex128, z float64) float64 {
	r2 := real(w)*real(w) + imag(w)*imag(w)
	return z * (1 - r2) / ((1 - z) * (1 - r2*z))
}

func main() {
	fmt.Println("🧬 LA FÓRMULA — todo ensamblado en el cambiaformas, armonizado en la dimensión 0")

	fmt.Println("\n   λₙ = Σ sobre pares [ |1 − wⁿ|²  +  (1 − |w|²ⁿ) ]")
	fmt.Println("                        LA FORMA        LA FUGA")
	fmt.Println("                     una distancia    el cambio de tamaño:")
	fmt.Println("                     al cuadrado —    vale 0 si y solo si")
	fmt.Println("                     jamás negativa   |w| = 1")

	// ---- pearls ----
	fmt.Println("\nrecogiendo perlas hasta t=1000…")
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 1000; t += 0.05 {
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
	gm := pearls[len(pearls)-1]
	tailI := (math.Log(gm/(2*math.Pi)) + 1) / (2 * math.Pi * gm)
	fmt.Printf("perlas: %d · cola de las no vistas: n²·%.2e\n", len(pearls), tailI)

	// ---- LAW 1: the split is exact, harmonic by harmonic ----
	fmt.Println("\nLEY 1 · LA PARTICIÓN ES EXACTA — y sobre la línea la FUGA es CERO")
	fmt.Println("   perla γ      armónico n     aporte del par      |1−wⁿ|²        la fuga")
	peorFuga := 0.0
	for _, pr := range [][2]float64{{14.134725, 3}, {21.022040, 7}, {30.424876, 12}, {49.773832, 20}} {
		g, n := pr[0], int(pr[1])
		w := wDe(complex(0.5, g))
		wn := cmplx.Pow(w, complex(float64(n), 0))
		aporte := 2 - 2*real(wn)
		d := cmplx.Abs(1-wn) * cmplx.Abs(1-wn)
		lk := aporte - d
		if math.Abs(lk) > peorFuga {
			peorFuga = math.Abs(lk)
		}
		fmt.Printf("   %9.6f      %2d        %12.9f    %12.9f    %+.1e\n", g, n, aporte, d, lk)
	}
	fmt.Printf("   → sobre la línea la fuga es cero a precisión de máquina (peor %.1e):\n", peorFuga)
	fmt.Println("     el aporte de cada par ES una distancia al cuadrado, y nada más")

	// ---- LAW 2: the whole thing in one function of z ----
	fmt.Println("\nLEY 2 · TODOS LOS SOBRES PROYECTADOS EN UNA SOLA FORMA — L(z) = Σ λₙ zⁿ")
	fmt.Println("   se compara la fórmula ensamblada (sumando distancias sobre las perlas)")
	fmt.Println("   contra el germen leído en el broche (que jamás vio una perla):")
	// germ side
	const nMax = 90
	r0, M := 0.92, 16384
	fv := make([]complex128, M)
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		zz := complex(r0*math.Cos(th), r0*math.Sin(th))
		fv[j] = xiLD(1/(1-zz)) / ((1 - zz) * (1 - zz))
	}
	lam := make([]float64, nMax+1)
	for n := 0; n < nMax; n++ {
		var acc complex128
		for j := 0; j < M; j++ {
			th := 2 * math.Pi * float64(j) / float64(M)
			acc += fv[j] * cmplx.Exp(complex(0, -float64(n)*th))
		}
		lam[n+1] = real(acc) / (float64(M) * math.Pow(r0, float64(n)))
	}
	fmt.Println("      z        L(z) por la fórmula      L(z) por el germen        desvío rel.")
	peorL := 0.0
	type fila struct{ z, f, g float64 }
	var filas []fila
	for _, z := range []float64{0.15, 0.30, 0.45, 0.60} {
		// assembled side: sum of shape terms over pearls, plus the tail
		suma := 0.0
		for _, g := range pearls {
			suma += forma(wDe(complex(0.5, g)), z)
		}
		// unseen pearls: |1-w|^2 ~ 1/gamma^2 and |1-wz|^2 -> (1-z)^2
		suma += tailI * z * (1 + z) / ((1 - z) * (1 - z) * (1 - z))
		// germ side
		gsum, zp := 0.0, 1.0
		for n := 1; n <= nMax; n++ {
			zp *= z
			gsum += lam[n] * zp
		}
		rel := math.Abs(suma-gsum) / math.Max(math.Abs(gsum), 1e-300)
		if rel > peorL {
			peorL = rel
		}
		filas = append(filas, fila{z, suma, gsum})
		fmt.Printf("    %.2f      %16.9f      %16.9f        %.1e\n", z, suma, gsum, rel)
	}
	fmt.Printf("   → las dos vías coinciden (peor desvío relativo %.1e): la fórmula ES el libro entero\n", peorL)
	fmt.Println("     y no hizo falta abrir ni un sobre: se proyectó la forma de los infinitos de una vez")

	// ---- LAW 3: every piece is manifestly positive ----
	fmt.Println("\nLEY 3 · CADA PIEZA ES POSITIVA POR SU FORMA — no por haberla medido")
	minTerm := math.Inf(1)
	for _, g := range pearls {
		if f := forma(wDe(complex(0.5, g)), 0.60); f < minTerm {
			minTerm = f
		}
	}
	fmt.Printf("   en z=0.60, el término más chico de las %d perlas: %.6e\n", len(pearls), minTerm)
	fmt.Println("   → todos positivos, y no por casualidad: cada uno es una distancia al cuadrado")
	fmt.Println("     dividida por otra distancia al cuadrado, por un factor positivo")

	// ---- LAW 4: the ghost leaks ----
	fmt.Println("\nLEY 4 · EL FANTASMA FUGA — y la fuga se ve en la fórmula, a simple vista")
	fmt.Println("   β        |w|            la fuga en z=0.60      ¿signo?")
	type fg struct{ beta, r, f float64 }
	var fgs []fg
	for _, beta := range []float64{0.50, 0.51, 0.60, 0.75, 0.90} {
		g := 14.134725
		w := wDe(complex(beta, g))
		r := cmplx.Abs(w)
		lk := fuga(w, 0.60)
		signo := "cero ★"
		if lk > 1e-14 {
			signo = "positiva (encoge)"
		} else if lk < -1e-14 {
			signo = "NEGATIVA (crece) ✗"
		}
		fgs = append(fgs, fg{beta, r, lk})
		fmt.Printf("  %.2f     %.9f     %+.6e       %s\n", beta, r, lk, signo)
	}
	fmt.Println("   → y su gemelo tiene |w| > 1, así que su fuga es NEGATIVA: por ahí se escapa la positividad")
	fmt.Println("   → la fórmula muestra el agujero en un solo término, y el término se anula SOLO si |w|=1")

	fmt.Println("\n════════ LO QUE ENTREGA LA FÓRMULA, Y LO QUE NO ════════")
	fmt.Println("ENTREGA: los infinitos sobres proyectados en UNA forma — no hizo falta abrir ninguno.")
	fmt.Println("Toda la cadena del capitán vive en una línea: la distancia al cuadrado (jamás negativa),")
	fmt.Println("el cambio de tamaño (la fuga), y el hecho de que la fuga se apaga exactamente cuando el")
	fmt.Println("paso es rotación pura. Verificado contra el germen del broche, que jamás vio una perla.")
	fmt.Println("\nNO ENTREGA, y hay que decirlo derecho: la fórmula es manifiestamente positiva CUANDO")
	fmt.Println("la fuga es cero, y la fuga es cero exactamente cuando |w| = 1 — que es lo que hay que")
	fmt.Println("demostrar. El círculo NO está roto; pero ahora está encerrado en UN SOLO TÉRMINO de UNA")
	fmt.Println("SOLA fórmula, en vez de repartido por todo el problema. Todo el millón cabe en:")
	fmt.Println("\n        ¿POR QUÉ LA FUGA TIENE QUE SER CERO?")
	fmt.Println("\nTodavía no. Pero el problema entero ya es un término de una ecuación.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 1000.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🧬 LA FÓRMULA — todo ensamblado en el cambiaformas, armonizado en la dimensión 0</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"no hace falta abrir todos los sobres si tenemos la representación de todos los números: se proyecta la forma" — el capitán</text>`,
		W, H, W, H, W/2, W/2)

	// the formula
	fmt.Fprintf(&b, `<rect x="60" y="105" width="1380" height="215" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="152" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">λₙ  =  Σ<tspan font-size="15">pares</tspan>  [  |1 − wⁿ|²  +  ( 1 − |w|²ⁿ )  ]</text>
<text x="%.0f" y="196" font-size="15" text-anchor="middle" fill="#7fd7a8">LA FORMA</text>
<text x="%.0f" y="196" font-size="15" text-anchor="middle" fill="#ff8fa0">LA FUGA</text>
<text x="%.0f" y="222" font-size="13" text-anchor="middle" fill="#dce8f7">una distancia al cuadrado</text>
<text x="%.0f" y="222" font-size="13" text-anchor="middle" fill="#dce8f7">el cambio de TAMAÑO</text>
<text x="%.0f" y="244" font-size="13" text-anchor="middle" fill="#7fd7a8">jamás puede ser negativa</text>
<text x="%.0f" y="244" font-size="13" text-anchor="middle" fill="#ff8fa0">vale 0 ⟺ |w| = 1</text>
<text x="%.0f" y="284" font-size="14" text-anchor="middle" fill="#ffd166">toda la cadena del capitán en una línea: la distancia · la dirección · el norte por el sur · la mitad</text>
<text x="%.0f" y="308" font-size="13" text-anchor="middle" fill="#8fa8c7">y proyectada a todos los armónicos a la vez: L(z) = Σ λₙ zⁿ, una sola función en el disco de la dimensión 0</text>`,
		W/2, 590.0, 1010.0, 590.0, 1010.0, 590.0, 1010.0, W/2, W/2)

	// left: the two ways agree
	fmt.Fprintf(&b, `<rect x="60" y="345" width="700" height="300" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="410" y="379" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LOS INFINITOS SOBRES, EN UNA SOLA FORMA</text>
<text x="410" y="405" font-size="12.5" text-anchor="middle" fill="#dce8f7">la fórmula ensamblada contra el germen del broche</text>
<text x="410" y="427" font-size="12" text-anchor="middle" fill="#8fa8c7">(el germen jamás vio una perla)</text>
<text x="410" y="455" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">z        por la fórmula        por el germen</text>`)
	for i, f := range filas {
		fmt.Fprintf(&b, `<text x="410" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%.2f   %14.8f   %14.8f</text>`,
			482.0+float64(i)*26, f.z, f.f, f.g)
	}
	fmt.Fprintf(&b, `<text x="410" y="606" font-size="13" text-anchor="middle" fill="#7fd7a8">las dos vías coinciden (peor desvío relativo %.0e)</text>
<text x="410" y="630" font-size="12.5" text-anchor="middle" fill="#ffd166">no se abrió ni un sobre: se proyectó la forma de los infinitos</text>`, peorL)

	// right: the leak
	fmt.Fprintf(&b, `<rect x="790" y="345" width="650" height="300" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="1115" y="379" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">LA FUGA — el único término que puede romper todo</text>
<text x="1115" y="407" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β         |w|            la fuga en z=0.60</text>`)
	for i, f := range fgs {
		col := "#ff8fa0"
		extra := ""
		if f.beta == 0.50 {
			col = "#7fd7a8"
			extra = "   ★ CERO"
		}
		fmt.Fprintf(&b, `<text x="1115" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="%s">%.2f    %.9f    %+.4e%s</text>`,
			434.0+float64(i)*26, col, f.beta, f.r, f.f, extra)
	}
	fmt.Fprintf(&b, `<text x="1115" y="580" font-size="12.5" text-anchor="middle" fill="#dce8f7">y el gemelo del fantasma tiene |w| &gt; 1: su fuga es NEGATIVA</text>
<text x="1115" y="606" font-size="12.5" text-anchor="middle" fill="#ffd166">por ahí, y solo por ahí, se escapa la positividad</text>
<text x="1115" y="630" font-size="12.5" text-anchor="middle" fill="#7fd7a8">sobre la línea la fuga es cero a precisión de máquina (%.0e)</text>`, peorFuga)

	// verdict
	fmt.Fprintf(&b, `<rect x="60" y="670" width="1380" height="240" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="706" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LO QUE ENTREGA — Y LO QUE NO</text>
<text x="%.0f" y="740" font-size="14" text-anchor="middle" fill="#7fd7a8">ENTREGA: los infinitos sobres proyectados en UNA forma, sin abrir ninguno. Toda la cadena del capitán en una línea,</text>
<text x="%.0f" y="764" font-size="14" text-anchor="middle" fill="#7fd7a8">verificada contra el germen del broche que jamás vio una perla.</text>
<text x="%.0f" y="800" font-size="14" text-anchor="middle" fill="#ff8fa0">NO ENTREGA, dicho derecho: la fórmula es manifiestamente positiva CUANDO la fuga es cero — y la fuga es cero</text>
<text x="%.0f" y="824" font-size="14" text-anchor="middle" fill="#ff8fa0">exactamente cuando |w| = 1, que es lo que hay que demostrar. El círculo no está roto.</text>
<text x="%.0f" y="856" font-size="14.5" text-anchor="middle" fill="#dce8f7">PERO ahora está encerrado en UN SOLO TÉRMINO de UNA SOLA fórmula, en vez de repartido por todo el problema:</text>
<text x="%.0f" y="890" font-size="20" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">¿POR QUÉ LA FUGA TIENE QUE SER CERO?</text>
<text x="%.0f" y="936" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("la-formula.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: la-formula.svg")
}
