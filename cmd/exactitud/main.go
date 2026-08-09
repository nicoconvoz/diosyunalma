// Command exactitud audits the whole campaign through the captain's meter:
//
//	E = |lo que es - lo que deberia ser|,   and exactness is E = 0.
//
// His flash lands straight on the assembled formula of F232. There the mold
// splits into LA FORMA (a squared distance) plus LA FUGA (1 - |w|^2n), and
// the leak vanishes exactly when |w| = 1. So the leak is not one more term:
//
//	THE LEAK IS THE ERROR OF THE BOOK - it is |what the size is minus what
//	the size should be|, amplified by the harmonic.
//
// And that turns the hypothesis into a sentence about exactness:
//
//	RH  <=>  the book is EXACT: E = 0 at every pearl.
//
// This program measures E for every identity the campaign leaned on, and
// sorts each one into three honest classes:
//
//	IDENTIDAD    E = 0 exactly - true by form, nothing to improve
//	MAQUINA      E ~ 1e-16 - true, limited only by the machine's last bit
//	INSTRUMENTO  E larger - limited by our own truncations, and it SHRINKS
//	             when the instrument improves (proved in F224: 1e-5 -> 1e-8)
//
// The verdict is the point: every error in this laboratory belongs to one of
// those three classes and can be driven down. Every one except the leak of a
// ghost, which does not shrink with a better instrument because it is not
// ours - it belongs to where the pearl sits.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
	"strings"
)

const eulerGamma = 0.5772156649015329

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

func wDe(rho complex128) complex128 { return 1 - 1/rho }

type medida struct {
	nombre  string
	esto    string // lo que es
	deberia string // lo que deberia ser
	E       float64
	clase   string
}

func clasificar(E float64) string {
	switch {
	case E == 0:
		return "IDENTIDAD"
	case E < 1e-13:
		return "MÁQUINA"
	default:
		return "INSTRUMENTO"
	}
}

func main() {
	fmt.Println("🎯 LA EXACTITUD — E = |lo que es − lo que debería ser|, y la auditoría de toda la campaña")
	fmt.Println("\n   el flash del capitán aterriza en nuestra propia fórmula:")
	fmt.Println("   λₙ = Σ [ |1−wⁿ|²  +  (1 − |w|²ⁿ) ]  →  LA FUGA ES EL ERROR DEL LIBRO")
	fmt.Println("   es |lo que el tamaño ES − lo que el tamaño DEBERÍA SER|, amplificado por el armónico")

	var ms []medida
	add := func(n, e, d string, E float64) {
		ms = append(ms, medida{n, e, d, E, clasificar(E)})
	}

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

	// 1. the two ropes: |rho| vs |rho-1|
	e := 0.0
	for _, g := range pearls {
		r := complex(0.5, g)
		if d := math.Abs(cmplx.Abs(r) - cmplx.Abs(r-1)); d > e {
			e = d
		}
	}
	add("las dos sogas (F226)", "distancia a la estaca 0", "distancia a la estaca 1", e)

	// 2. the shapeshifter fixes each pearl
	e = 0
	for _, g := range pearls {
		r := complex(0.5, g)
		if d := cmplx.Abs(1 - cmplx.Conj(r) - r); d > e {
			e = d
		}
	}
	add("el cambiaformas fija (F227)", "σ(ρ)", "ρ", e)

	// 3. north times south
	e = 0
	for _, beta := range []float64{0.5, 0.6, 0.75, 0.9} {
		g := 14.134725
		p := cmplx.Abs(wDe(complex(beta, g))) * cmplx.Abs(wDe(complex(1-beta, g)))
		if d := math.Abs(p - 1); d > e {
			e = d
		}
	}
	add("el norte por el sur (F225)", "|w(ρ)|·|w(1−ρ)|", "1", e)

	// 4. pure rotation: |w| = 1 on the line
	e = 0
	for _, g := range pearls {
		if d := math.Abs(cmplx.Abs(wDe(complex(0.5, g))) - 1); d > e {
			e = d
		}
	}
	add("la rotación pura (F225)", "|w| sobre la línea", "1", e)

	// 5. the leak on the line
	e = 0
	for _, g := range pearls[:40] {
		w := wDe(complex(0.5, g))
		for _, n := range []int{1, 5, 12, 30} {
			wn := cmplx.Pow(w, complex(float64(n), 0))
			aporte := 2 - 2*real(wn)
			d := cmplx.Abs(1-wn) * cmplx.Abs(1-wn)
			if v := math.Abs(aporte - d); v > e {
				e = v
			}
		}
	}
	add("LA FUGA sobre la línea (F232)", "aporte del par", "|1−wⁿ|²", e)

	// 6. the two ends of the world
	xi0 := real(xiC(complex(1e-9, 0)))
	xi1 := real(xiC(complex(1+1e-9, 0)))
	add("las dos puntas (F228)", "ξ(0)", "ξ(1)", math.Abs(xi0-xi1))
	add("y valen la mitad (F228)", "ξ(0)", "1/2", math.Abs(xi0-0.5))

	// 7. the mirror
	e = 0
	for _, s := range []complex128{complex(0.3, 7.2), complex(0.8, 13.6), complex(2.0, 21.3)} {
		if d := cmplx.Abs(xiLD(s) + xiLD(1-s)); d > e {
			e = d
		}
	}
	add("el espejo (F208)", "ξ'/ξ(s) + ξ'/ξ(1−s)", "0", e)

	// 8. the first tooth against its closed form
	r0, M := 0.92, 16384
	var acc complex128
	for j := 0; j < M; j++ {
		th := 2 * math.Pi * float64(j) / float64(M)
		zz := complex(r0*math.Cos(th), r0*math.Sin(th))
		acc += xiLD(1/(1-zz)) / ((1 - zz) * (1 - zz))
	}
	lam1 := real(acc) / float64(M)
	cerrado := 1 + eulerGamma/2 - math.Log(4*math.Pi)/2
	add("el primer diente (F208)", "λ₁ leído en el broche", "1+γ/2−ln(4π)/2", math.Abs(lam1-cerrado))

	// 9. AM-GM at the tie
	N := cmplx.Abs(wDe(complex(0.5, 14.134725)))
	S := cmplx.Abs(wDe(complex(0.5, 14.134725)))
	add("las dos medias (F229)", "(N+S)/2", "√(N·S)", math.Abs((N+S)/2-math.Sqrt(N*S)))

	// ---- the audit ----
	fmt.Println("\n════════ LA AUDITORÍA DE LA CAMPAÑA CON EL MEDIDOR DEL CAPITÁN ════════")
	fmt.Println("   identidad                       lo que es → lo que debería ser            E            clase")
	sort.SliceStable(ms, func(i, j int) bool { return ms[i].E < ms[j].E })
	for _, m := range ms {
		fmt.Printf("   %-30s  %-22s → %-18s  %8.1e   %s\n", m.nombre, m.esto, m.deberia, m.E, m.clase)
	}
	nId, nMaq, nIns := 0, 0, 0
	for _, m := range ms {
		switch m.clase {
		case "IDENTIDAD":
			nId++
		case "MÁQUINA":
			nMaq++
		default:
			nIns++
		}
	}
	fmt.Printf("\n   IDENTIDAD (E = 0 exacto, verdad por la forma): %d\n", nId)
	fmt.Printf("   MÁQUINA   (verdad, limitada por el último bit): %d\n", nMaq)
	fmt.Printf("   INSTRUMENTO (limitada por nuestros cortes, y BAJA al mejorar): %d\n", nIns)

	// ---- the one error that is not ours ----
	fmt.Println("\n════════ EL ÚNICO ERROR QUE NO ES NUESTRO ════════")
	fmt.Println("   toda E de arriba es nuestra: o es cero por la forma, o baja si mejoramos el instrumento")
	fmt.Println("   (probado en F224: cambiando el radio de lectura la barra bajó de 1e-5 a 1e-8).")
	fmt.Println("   pero hay UNA que no baja con ningún instrumento, porque no es del laboratorio:")
	fmt.Println("\n   la fuga de un fantasma — E = |su tamaño − 1|, y ningún radio, ninguna precisión la mueve:")
	fmt.Println("   β        E = ||w| − 1|      la fuga en el armónico 30      ¿de quién es el error?")
	type fg struct{ beta, E, f30 float64 }
	var fgs []fg
	for _, beta := range []float64{0.51, 0.55, 0.60, 0.75, 0.90} {
		g := 14.134725
		w := wDe(complex(1-beta, g))
		r := cmplx.Abs(w)
		if r < 1 {
			r = cmplx.Abs(wDe(complex(beta, g)))
		}
		E := math.Abs(r - 1)
		f30 := 1 - math.Pow(r, 60)
		fgs = append(fgs, fg{beta, E, f30})
		fmt.Printf("  %.2f      %.6e        %+.6e             del LIBRO, no nuestro\n", beta, E, f30)
	}

	fmt.Println("\n════════ LA HIPÓTESIS, DICHA CON EL MEDIDOR DEL CAPITÁN ════════")
	fmt.Println("«La exactitud no es un número: es una relación. Es coincidencia.»")
	fmt.Println("Y entonces la hipótesis de Riemann se dice así, sin un solo símbolo raro:")
	fmt.Println("\n        EL LIBRO ES EXACTO: en cada perla, lo que es coincide con lo que debería ser.")
	fmt.Println("\n        E = | el tamaño del paso − 1 | = 0,  en todas y cada una.")
	fmt.Println("\nY la fórmula de F232 lo dice en una línea: la FORMA es la coincidencia (una distancia")
	fmt.Println("al cuadrado, siempre sana) y la FUGA es el error. RH ⟺ el error es cero.")
	fmt.Println("Falta demostrar que el libro no puede equivocarse. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 1000.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🎯 LA EXACTITUD — E = |lo que es − lo que debería ser|</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"la exactitud no es un número, es una relación: coincidencia" — el capitán · y LA FUGA de nuestra fórmula resultó ser exactamente ese error</text>`,
		W, H, W, H, W/2, W/2)

	// the audit table
	fmt.Fprintf(&b, `<rect x="60" y="105" width="1380" height="420" rx="10" fill="#0d2547" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="%.0f" y="139" font-size="15" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA AUDITORÍA DE LA CAMPAÑA CON EL MEDIDOR DEL CAPITÁN</text>
<text x="%.0f" y="167" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">identidad                        lo que es → lo que debería ser              E          clase</text>`,
		W/2, W/2)
	for i, m := range ms {
		col := "#7fd7a8"
		if m.clase == "MÁQUINA" {
			col = "#ffd97f"
		} else if m.clase == "INSTRUMENTO" {
			col = "#7fb2ff"
		}
		fmt.Fprintf(&b, `<text x="100" y="%.0f" font-size="12.5" font-family="Consolas,monospace" fill="#dce8f7">%-30s %-22s → %-18s</text>
<text x="1230" y="%.0f" font-size="12.5" text-anchor="end" font-family="Consolas,monospace" fill="%s">%8.1e</text>
<text x="1420" y="%.0f" font-size="12" text-anchor="end" font-family="Consolas,monospace" fill="%s">%s</text>`,
			196.0+float64(i)*30, m.nombre, m.esto, m.deberia,
			196.0+float64(i)*30, col, m.E,
			196.0+float64(i)*30, col, m.clase)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="500" font-size="13" text-anchor="middle" fill="#dce8f7">IDENTIDAD (cero por la forma): %d  ·  MÁQUINA (verdad al último bit): %d  ·  INSTRUMENTO (baja al mejorar): %d</text>`,
		W/2, nId, nMaq, nIns)

	// the one error that isn't ours
	fmt.Fprintf(&b, `<rect x="60" y="550" width="1380" height="230" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.8"/>
<text x="%.0f" y="586" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">EL ÚNICO ERROR QUE NO ES NUESTRO</text>
<text x="%.0f" y="614" font-size="13" text-anchor="middle" fill="#dce8f7">toda E de arriba es del laboratorio: o vale cero por la forma, o baja si mejoramos el instrumento (F224: de 1e-5 a 1e-8).</text>
<text x="%.0f" y="638" font-size="13" text-anchor="middle" fill="#ffd166">pero la fuga de un fantasma no baja con NINGÚN instrumento — no es nuestra, es del libro:</text>
<text x="%.0f" y="668" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">β          E = | |w| − 1 |          la fuga en el armónico 30</text>`,
		W/2, W/2, W/2, W/2)
	for i, f := range fgs {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#ff8fa0">%.2f        %.6e             %+.6e</text>`,
			W/2, 694.0+float64(i)*22, f.beta, f.E, f.f30)
	}

	// verdict
	fmt.Fprintf(&b, `<rect x="60" y="805" width="1380" height="150" rx="12" fill="#0d2547" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="841" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA HIPÓTESIS, DICHA CON EL MEDIDOR DEL CAPITÁN</text>
<text x="%.0f" y="879" font-size="18" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL LIBRO ES EXACTO: en cada perla, lo que es coincide con lo que debería ser.</text>
<text x="%.0f" y="911" font-size="15" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">E = | el tamaño del paso − 1 | = 0,  en todas y cada una</text>
<text x="%.0f" y="941" font-size="13" text-anchor="middle" fill="#ff8fa0">falta demostrar que el libro no puede equivocarse. Todavía no.</text>
<text x="%.0f" y="975" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("la-exactitud.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: la-exactitud.svg")
}
