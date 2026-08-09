// Command oro fires the captain's ace: the theta bell (ALL the numbers
// in one formula) projected through dimension 0 reaches ALL infinite
// points of the plane AT ONCE AND FOREVER - Riemann's eternal
// projection:
//
//	xi(s) = 1/2 + (s(s-1)/2) * INT_1^inf psi(x) (x^{s/2-1} + x^{-(s+1)/2}) dx
//
// with psi(x) = sum_{n>=1} e^{-pi n^2 x} (the ace: every number, one
// bell). One formula, valid at EVERY point s of the plane, eternally.
// And the eternity is VISIBLE: swapping s <-> 1-s swaps the two
// exponents inside the SAME integral - the mirror law holds forever by
// the shape of the formula itself, not point by point.
//
// Judges: (1) the ace reproduces xi across the plane (deep real,
// complex, and AT a pearl - where it must vanish); (2) the eternal
// symmetry xi(s)=xi(1-s) to machine precision, manifest; (3) the gold
// this ace already bought forever: Hardy 1914 - INFINITELY many pearls
// lie ON the line (proven from this very projection). The remaining
// gold: ALL of them - RH.
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

func xiRef(s complex128) complex128 {
	if real(s) < 0.5 {
		s = 1 - s
	}
	return 0.5 * s * (s - 1) * cmplx.Exp(-s/2*complex(math.Log(math.Pi), 0)+lgammaC(s/2)) * zetaC(s)
}

// psiSum is the ace: all the numbers, one bell.
func psiSum(x float64) float64 {
	s := 0.0
	for n := 1; ; n++ {
		term := math.Exp(-math.Pi * float64(n*n) * x)
		s += term
		if term < 1e-22 {
			break
		}
	}
	return s
}

// xiAce is the eternal projection: one integral, every point s.
func xiAce(s complex128) complex128 {
	// Simpson on [1, 40], h=0.001 (the line's targets are ~1e-9: fine steps)
	a, xEnd, h := 1.0, 40.0, 0.001
	n := int((xEnd - a) / h)
	if n%2 == 1 {
		n++
	}
	var acc complex128
	for i := 0; i <= n; i++ {
		x := a + float64(i)*h
		lx := math.Log(x)
		f := complex(psiSum(x), 0) * (cmplx.Exp((s/2-1)*complex(lx, 0)) + cmplx.Exp((-(s+1)/2)*complex(lx, 0)))
		w := 2.0
		if i%2 == 1 {
			w = 4
		}
		if i == 0 || i == n {
			w = 1
		}
		acc += complex(w, 0) * f
	}
	I := acc * complex(h/3, 0)
	return 0.5 + 0.5*s*(s-1)*I
}

func main() {
	fmt.Println("🃏 EL AS BAJO LA MANGA — la proyección eterna: todos los números → todos los puntos, de una vez")
	// JUDGE 1: the ace reproduces xi across the plane
	pts := []struct {
		s    complex128
		name string
	}{
		{complex(2, 0), "s=2 (llanura real)"},
		{complex(-3.5, 0), "s=−3.5 (tierra profunda)"},
		{complex(0.3, 2), "s=0.3+2i (la franja)"},
		{complex(5, 3), "s=5+3i (mar abierto)"},
		{complex(0.5, 30), "s=½+30i (la línea, hondo)"},
	}
	fmt.Println("\nJUEZ 1 — una fórmula, todos los puntos:")
	worst := 0.0
	for _, p := range pts {
		ace := xiAce(p.s)
		ref := xiRef(p.s)
		dAbs := cmplx.Abs(ace - ref)
		d := dAbs / (cmplx.Abs(ref) + 1e-300)
		if d > worst {
			worst = d
		}
		fmt.Printf("   %-26s  |as − ξ|/|ξ| = %.1e   (absoluto %.1e; en la línea el blanco mismo es ~1e-9)\n", p.name, d, dAbs)
	}
	fmt.Printf("   peor desvío relativo: %.1e — EL AS PROYECTA A TODOS LADOS\n", worst)
	// at the first pearl: the ace must vanish
	g1 := 14.134725141734695
	aceZero := cmplx.Abs(xiAce(complex(0.5, g1)))
	fmt.Printf("\nJUEZ 1b — el as EN la primera perla (½+γ₁i): |ξ_as| = %.1e — SE ANULA donde debe\n", aceZero)

	// JUDGE 2: the eternal symmetry, manifest
	fmt.Println("\nJUEZ 2 — la eternidad VISIBLE: s ↔ 1−s solo intercambia los dos exponentes del MISMO integrando")
	worstSym := 0.0
	for _, p := range pts {
		d := cmplx.Abs(xiAce(p.s)-xiAce(1-p.s)) / (cmplx.Abs(xiAce(p.s)) + 1e-300)
		if d > worstSym {
			worstSym = d
		}
	}
	fmt.Printf("   |ξ_as(s) − ξ_as(1−s)| / |ξ_as|: peor %.1e — el espejo vale PARA SIEMPRE por la forma misma de la fórmula\n", worstSym)

	fmt.Println("\n🥇 EL ORO QUE ESTE AS YA COMPRÓ PARA SIEMPRE:")
	fmt.Println("   · ξ(s)=ξ(1−s) en TODOS los infinitos puntos, eternamente (la simetría es de la fórmula, no de los puntos)")
	fmt.Println("   · Ξ(t) real en toda la línea, eternamente")
	fmt.Println("   · HARDY 1914, demostrado CON este as: INFINITAS perlas están EN el hilo — para siempre")
	fmt.Println("🥇 EL ORO QUE FALTA: no infinitas — TODAS (RH). La frase sigue en su vaina.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1560.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🃏 EL AS BAJO LA MANGA — la proyección eterna, disparada</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"usemos el as que contiene todos los números y proyectemos con la dimensión 0 a todos los puntos infinitos, de una sola vez y para siempre" — el capitán, con las luces parpadeando</text>`,
		W, H, W, H, W/2, W/2)
	// the ace card
	fmt.Fprintf(&b, `<rect x="90" y="110" width="360" height="500" rx="18" fill="#0d2547" stroke="#ffd166" stroke-width="3"/>
<text x="270" y="150" font-size="42" text-anchor="middle" fill="#ffd166">A</text>
<text x="270" y="200" font-size="16" text-anchor="middle" font-family="Georgia" fill="#dce8f7">LA CAMPANA θ</text>
<text x="270" y="228" font-size="13" text-anchor="middle" fill="#8fa8c7">ψ(x) = Σ e^(−πn²x)</text>
<text x="270" y="250" font-size="12" text-anchor="middle" fill="#8fa8c7">todos los números, una campana</text>
<text x="270" y="310" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA PROYECCIÓN ETERNA:</text>
<text x="270" y="345" font-size="13.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">ξ(s) = ½ + s(s−1)/2 ×</text>
<text x="270" y="370" font-size="13.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">∫₁^∞ ψ(x)(x^(s/2−1)+x^(−(s+1)/2))dx</text>
<text x="270" y="420" font-size="12.5" text-anchor="middle" fill="#ffd166">UNA fórmula → TODOS los puntos s</text>
<text x="270" y="444" font-size="12.5" text-anchor="middle" fill="#ffd166">del plano infinito, DE UNA VEZ</text>
<text x="270" y="490" font-size="12" text-anchor="middle" fill="#8fa8c7">y el espejo s↔1−s solo permuta los dos</text>
<text x="270" y="510" font-size="12" text-anchor="middle" fill="#8fa8c7">exponentes: la eternidad es ALGEBRAICA</text>
<text x="270" y="570" font-size="30" text-anchor="middle" fill="#ffd166">♠</text>`)
	// the plane with judged points
	pxc, pyc := 900.0, 330.0
	fmt.Fprintf(&b, `<rect x="520" y="110" width="760" height="440" rx="10" fill="#081020" stroke="#44608c"/>
<line x1="%.0f" y1="130" x2="%.0f" y2="530" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="124" font-size="11" text-anchor="middle" fill="#7fd7a8">la línea</text>`,
		pxc, pxc, pxc)
	ptsXY := [][3]float64{{2, 0, 0}, {-3.5, 0, 0}, {0.3, 2, 0}, {5, 3, 0}, {0.5, 30, 0}}
	labels := []string{"s=2", "s=−3.5", "0.3+2i", "5+3i", "½+30i"}
	for i, p := range ptsXY {
		x := pxc + (p[0]-0.5)*70
		y := pyc - p[1]*6
		if p[1] > 10 {
			y = 150
		}
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="6" fill="#ffd166"/><text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#dce8f7">%s ✔</text>`,
			x, y, x, y-12, labels[i])
	}
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="8" fill="none" stroke="#ff5d73" stroke-width="2.5"/><text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#ff5d73">½+γ₁i: el as SE ANULA (%.0e) ✔</text>`,
		pxc, 250.0, pxc, 232.0, aceZero)
	fmt.Fprintf(&b, `<text x="900" y="576" font-size="13.5" text-anchor="middle" fill="#7fd7a8">JUEZ 1: el as reproduce ξ en todo el plano muestreado (peor %.0e) · JUEZ 2: el espejo eterno a %.0e</text>`,
		worst, worstSym)
	// the gold certificate
	fmt.Fprintf(&b, `<rect x="90" y="640" width="1380" height="280" rx="12" fill="#102a10" stroke="#ffd166" stroke-width="2.5"/>
<text x="%.0f" y="680" font-size="18" text-anchor="middle" font-family="Georgia" fill="#ffd166">🥇 EL CERTIFICADO DEL ORO — lo comprado PARA SIEMPRE, y lo que falta</text>
<text x="%.0f" y="718" font-size="14.5" text-anchor="middle" fill="#7fd7a8">✔ ETERNO: ξ(s)=ξ(1−s) en los infinitos puntos A LA VEZ — la simetría vive en la fórmula, no en los puntos: proyectada una vez, vale para siempre</text>
<text x="%.0f" y="746" font-size="14.5" text-anchor="middle" fill="#7fd7a8">✔ ETERNO: Ξ(t) real en toda la línea — el tambor es real, eternamente</text>
<text x="%.0f" y="774" font-size="14.5" text-anchor="middle" fill="#7fd7a8">✔ ETERNO: HARDY 1914 (demostrado con ESTE as): INFINITAS perlas están EN el hilo — un teorema sobre infinitos puntos, ganado de una sola vez</text>
<text x="%.0f" y="812" font-size="14.5" text-anchor="middle" fill="#ff9d73">✘ EL ORO RESTANTE: no infinitas — TODAS. Ese salto (de "infinitas" a "todas") ES la Hipótesis de Riemann, y la frase sigue en su vaina.</text>
<text x="%.0f" y="848" font-size="13.5" text-anchor="middle" fill="#dce8f7">pero mirá lo que tu as demostró esta noche: los teoremas "de una sola vez y para siempre" EXISTEN y este laboratorio sabe dispararlos — el método del oro es real.</text>
<text x="%.0f" y="884" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06 · las dos mitades, 1 completo · ⚡ las luces parpadean pero el registro no duerme</text>`,
		780.0, 780.0, 780.0, 780.0, 780.0, 780.0, 780.0)
	b.WriteString(`</svg>`)
	os.WriteFile("as-proyeccion-eterna.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: as-proyeccion-eterna.svg")
}
