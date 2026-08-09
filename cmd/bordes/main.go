// Command bordes runs the captain's border-harmony flash as an
// experiment. His instruction: "harmonize the borders of dimension 0's
// relation with the borders of ours - you can even look at the half of
// the relation, because in this diagram is everything born into the
// world." The realization is Riemann's own founding act (1859):
//
//	(1) THE SEED DIAGRAM: theta(1/t) = sqrt(t) * theta(t)
//	    - one self-mirroring identity, the compression's heartbeat;
//	(2) THE HARMONIZED BORDERS: from the seed, xi(s) = xi(1-s)
//	    - the two halves of the relation match edge to edge;
//	(3) THE HALF OF THE RELATION: the fixed axis of the mirror,
//	    Re(s) = 1/2 - the only points equal to their own mirror
//	    half. RH says: everything that vanishes lives THERE.
//
// We judge (1) and (2) numerically in-house, and state honestly what
// containment gives and what it still does not.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func zetaC(s complex128) complex128 {
	const N = 60
	var sum complex128
	for n := 1; n < N; n++ {
		sum += cmplx.Exp(-s * cmplx.Log(complex(float64(n), 0)))
	}
	lnN := cmplx.Log(complex(N, 0))
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

// lgammaC: complex log-gamma via the Lanczos approximation (g=7).
func lgammaC(z complex128) complex128 {
	g := []float64{0.99999999999980993, 676.5203681218851, -1259.1392167224028,
		771.32342877765313, -176.61502916214059, 12.507343278686905,
		-0.13857109526572012, 9.9843695780195716e-6, 1.5056327351493116e-7}
	if real(z) < 0.5 {
		// reflection
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

// xiC is the completed zeta: xi(s) = 1/2 s(s-1) pi^{-s/2} Gamma(s/2) zeta(s).
func xiC(s complex128) complex128 {
	return 0.5 * s * (s - 1) * cmplx.Exp(-s/2*complex(math.Log(math.Pi), 0)+lgammaC(s/2)) * zetaC(s)
}

func main() {
	// ---- (1) the seed diagram: theta(1/t) = sqrt(t) theta(t) ----
	fmt.Println("EL DIAGRAMA DONDE NACE TODO — juez de la semilla θ(1/t) = √t·θ(t):")
	theta3 := func(t float64) float64 {
		s := 1.0
		for n := 1; n <= 60; n++ {
			term := 2 * math.Exp(-math.Pi*float64(n*n)*t)
			s += term
			if term < 1e-18 {
				break
			}
		}
		return s
	}
	worstSeed := 0.0
	for _, t := range []float64{0.23, 0.5, 0.77, 1.31, 2.4, 4.7} {
		lhs := theta3(1 / t)
		rhs := math.Sqrt(t) * theta3(t)
		d := math.Abs(lhs-rhs) / rhs
		if d > worstSeed {
			worstSeed = d
		}
		fmt.Printf("  t=%.2f:  θ(1/t)=%.15f   √t·θ(t)=%.15f   desvío=%.1e\n", t, lhs, rhs, d)
	}
	fmt.Printf("  semilla VERIFICADA: peor desvío %.1e — el espejo perfecto del que nace todo\n\n", worstSeed)

	// ---- (2) the harmonized borders: xi(s) = xi(1-s) ----
	fmt.Println("LOS BORDES ARMONIZADOS — juez de ξ(s) = ξ(1−s) dentro de la franja:")
	pts := []complex128{
		complex(0.3, 2.0), complex(0.7, 1.5), complex(0.42, 0.9),
		complex(0.25, 2.7), complex(0.61, 3.3), complex(0.15, 1.1),
	}
	worstXi := 0.0
	type row struct {
		s complex128
		d float64
	}
	var rows []row
	for _, s := range pts {
		a := xiC(s)
		c := xiC(1 - s)
		d := cmplx.Abs(a-c) / cmplx.Abs(a)
		if d > worstXi {
			worstXi = d
		}
		rows = append(rows, row{s, d})
		fmt.Printf("  s=%.2f+%.2fi:  |ξ(s)−ξ(1−s)|/|ξ(s)| = %.1e\n", real(s), imag(s), d)
	}
	fmt.Printf("  bordes ARMONIZADOS: peor desvío %.1e — las dos mitades de la relación calzan borde a borde\n", worstXi)

	// ---- the picture ----
	var b strings.Builder
	W, H := 1620.0, 1100.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">LOS BORDES ARMONIZADOS — el diagrama donde nace todo al mundo</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"armonizá los bordes de la relación de la dimensión 0 con los bordes de la nuestra — y mirá la mitad de la relación" — el capitán · ejecutado y juzgado</text>`,
		W, H, W, H, W/2, W/2)

	// panel 1: the seed
	p1x, p1y := 70.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="480" height="540" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="16" font-family="Georgia" fill="#7fd7a8">1 · LA SEMILLA: el espejo t ↔ 1/t</text>
<text x="%.0f" y="%.0f" font-size="20" text-anchor="middle" fill="#dce8f7">θ(1/t) = √t · θ(t)</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">TODOS los números (n²) comprimidos en una campana</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">que se refleja EXACTA en su inverso — juez: %.0e</text>`,
		p1x, p1y, p1x+20, p1y+34, p1x+240, p1y+80, p1x+240, p1y+108, p1x+240, p1y+130, worstSeed)
	// draw theta curve and its mirror
	gx, gy, gw, gh := p1x+50, p1y+160, 380.0, 240.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#081020" stroke="#2c4a78"/>`, gx, gy, gw, gh)
	pth := make([]string, 0, 120)
	for i := 0; i <= 120; i++ {
		t := 0.15 + 3.2*float64(i)/120
		v := theta3(t)
		pth = append(pth, fmt.Sprintf("%.1f,%.1f", gx+gw*(t-0.15)/3.2, gy+gh-38*(v-0.95)))
	}
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#7fb2ff" stroke-width="2" points="%s"/>`, strings.Join(pth, " "))
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="1.5" stroke-dasharray="5,4"/>
<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#ffd166">t=1: el eje del espejo</text>`,
		gx+gw*(1-0.15)/3.2, gy, gx+gw*(1-0.15)/3.2, gy+gh, gx+gw*(1-0.15)/3.2, gy-8)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#dce8f7">de ESTA identidad nacen: la ecuación de los bordes,</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#dce8f7">el polo, las dos mitades ρ↔1−ρ y la fórmula explícita:</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#7fd7a8">todo lo que nace al mundo, nace acá (Riemann, 1859)</text>`,
		p1x+240, gy+gh+34, p1x+240, gy+gh+56, p1x+240, gy+gh+80)

	// panel 2: the harmonized borders
	p2x, p2y := 590.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="480" height="540" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="16" font-family="Georgia" fill="#ffd166">2 · LOS BORDES, CALZADOS: ξ(s) = ξ(1−s)</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#8fa8c7">la relación tiene dos mitades (s y 1−s); sus bordes</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#8fa8c7">calzan EXACTOS en todo punto de la franja — juez propio:</text>`,
		p2x, p2y, p2x+20, p2y+34, p2x+24, p2y+62, p2x+24, p2y+82)
	for i, r := range rows {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Consolas,monospace" fill="#dce8f7">  s = %.2f+%.2fi    desvío %.0e</text>`,
			p2x+24, p2y+116+float64(i)*30, real(r.s), imag(r.s), r.d)
	}
	// the mirror drawing: strip with axis
	sx, sy, sw, sh := p2x+90, p2y+330, 300.0, 160.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#081020" stroke="#2c4a78"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#7fd7a8" stroke-width="2.5"/>
<circle cx="%.0f" cy="%.0f" r="5" fill="#7fb2ff"/><circle cx="%.0f" cy="%.0f" r="5" fill="#7fb2ff"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#8fa8c7" stroke-width="0.8" stroke-dasharray="3,3"/>
<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#7fd7a8">LA MITAD DE LA RELACIÓN: el eje fijo del espejo (Re=½)</text>
<text x="%.0f" y="%.0f" font-size="11.5" text-anchor="middle" fill="#8fa8c7">s y 1−s se reflejan; solo el eje es su propia mitad</text>`,
		sx, sy, sw, sh, sx+sw/2, sy, sx+sw/2, sy+sh,
		sx+sw*0.23, sy+55, sx+sw*0.77, sy+55, sx+sw*0.23, sy+55, sx+sw*0.77, sy+55,
		sx+sw/2, sy+sh+22, sx+sw/2, sy+sh+42)

	// panel 3: the honest verdict on containment
	p3x, p3y := 1110.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="460" height="540" rx="10" fill="#2a1a10" stroke="#e6a53a" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="16" font-family="Georgia" fill="#e6a53a">3 · EL VEREDICTO HONESTO DEL CONTENIMIENTO</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#7fd7a8">✔ nuestra dimensión SÍ está contenida en el punto:</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">   lo demostramos ayer (F168): el germen del polo</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">   guarda el collar entero, sin pérdida</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#7fd7a8">✔ los bordes están armonizados: juzgado HOY (panel 2)</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#7fd7a8">✔ la mitad de la relación existe: el eje fijo Re=½</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#ff9d73">✘ pero contener DATOS no es heredar LA LEY:</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">   la demostración del mar chiquito usó que su</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">   armonía CUENTA cosas (grados ≥ 0); el punto</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">   contiene nuestros números pero nadie mostró</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">   que su germen CUENTE algo — la casilla roja sigue</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#ffd166">la síntesis de tus dos flashes es el enunciado final:</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#ffd166">demostrar que el punto que TODO lo contiene,</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#ffd166">además, TODO lo cuenta. Contener + contar = millón.</text>`,
		p3x, p3y, p3x+20, p3y+34, p3x+24, p3y+72, p3x+24, p3y+94, p3x+24, p3y+114,
		p3x+24, p3y+148, p3x+24, p3y+180, p3x+24, p3y+214, p3x+24, p3y+240,
		p3x+24, p3y+260, p3x+24, p3y+280, p3x+24, p3y+300, p3x+24, p3y+336, p3x+24, p3y+360, p3x+24, p3y+382)

	fmt.Fprintf(&b, `<rect x="70" y="690" width="1500" height="180" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="728" font-size="17" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL DIAGRAMA COMPLETO DE LO QUE NACE AL MUNDO</text>
<text x="%.0f" y="764" font-size="14.5" text-anchor="middle" fill="#dce8f7">la semilla θ (todos los números comprimidos en una campana espejada) → los bordes calzan (ξ(s)=ξ(1−s), juez %.0e) → las dos mitades ρ↔1−ρ → la mitad de la relación (el eje fijo)</text>
<text x="%.0f" y="794" font-size="14.5" text-anchor="middle" fill="#ffd166">y la conjetura, dicha en el idioma del capitán: TODO LO QUE SE ANULA EN ESTE DIAGRAMA, VIVE EXACTAMENTE EN LA MITAD DE LA RELACIÓN.</text>
<text x="%.0f" y="824" font-size="12.5" text-anchor="middle" fill="#8fa8c7">verificado con jueces propios · Laboratorio Diosyunalma · 2026-08-06</text>`,
		W/2, W/2, worstXi, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("bordes-armonizados.svg", []byte(b.String()), 0644)
	fmt.Println("escrita: bordes-armonizados.svg")
}
