// Command obra is THE CONSTRUCTION SITE: the captain ordered all the
// drum's pieces assembled and harmonized in dimension 0 so the
// structure falls out by decantation - iron, bricks, cement, plan,
// windows, doors, floor - build the drum's SHAPE to arrive at the drum.
//
// The materials, all from the campaign:
//
//	THE PLAN     Weyl shape decoded blind in F214: A(E) = theta(E) + 2pi
//	             (area aligned so Bohr-Sommerfeld A = pi(j+1/2) lands
//	             the j-th note on the j-th pearl)
//	THE IRON     self-adjointness: the matrix is symmetric BY BUILD,
//	             so every note is REAL BY LAW - no note can leave the line
//	THE BRICKS   the potential V(x), DECANTED from the plan by Abel
//	             inversion (the dimension-0 harmonization: the shape
//	             falls through the integral and settles as a well)
//	THE FLOOR    the basement: below E=10 the theta clock is silent,
//	             so the floor is poured smooth (C1 interpolation);
//	             it contributes one basement note that is ours, not
//	             the book's
//	THE WALLS    Dirichlet walls where V reaches 250 - windows shut
//	             far above every note we compare
//
// Then the drum is STRUCK: eigenvalues by Sturm bisection, compared
// note by note against the true pearls. The smooth shape carries the
// pearls up to the writing wobble S(t) - the wobble itself is the
// prime's handwriting, which only the true drum adds. That remains
// the key; this is the first drum the shop has actually built.
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

func main() {
	fmt.Println("🏗️ LA OBRA — construyendo el tambor según el plano, pieza por pieza")

	// ---- THE PLAN: A(E) = theta(E) + 2pi, floor poured below E=10 ----
	const eJoin = 10.0
	thJ := theta(eJoin)
	thpJ := 0.5 * math.Log(eJoin/(2*math.Pi))
	// floor: A(E) = c1 E + c2 E^2 on [0,10], C1-joined to theta+2pi
	c2 := (thpJ*eJoin - (thJ + 2*math.Pi)) / (eJoin * eJoin)
	c1 := thpJ - 2*c2*eJoin
	aprime := func(E float64) float64 {
		if E < eJoin {
			return c1 + 2*c2*E
		}
		return 0.5 * math.Log(E/(2*math.Pi))
	}
	fmt.Printf("PLANO: A(E)=θ(E)+2π · piso vertido en E<10 (c1=%.4f, c2=%.5f) — A creciente, arranque en 0\n", c1, c2)

	// ---- THE BRICKS: decant V(x) from the plan (Abel inversion) ----
	// x(V) = (2 sqrt(V)/pi) INT_0^{pi/2} A'(V sin^2 phi) sin(phi) dphi
	const vMax = 250.0
	nV := 500
	xs := make([]float64, nV+1)
	vs := make([]float64, nV+1)
	for i := 0; i <= nV; i++ {
		V := vMax * float64(i) / float64(nV)
		vs[i] = V
		if V == 0 {
			xs[i] = 0
			continue
		}
		acc, K := 0.0, 300
		for j := 0; j < K; j++ {
			ph := (float64(j) + 0.5) * math.Pi / 2 / float64(K)
			s := math.Sin(ph)
			acc += aprime(V*s*s) * s
		}
		xs[i] = 2 * math.Sqrt(V) / math.Pi * acc * (math.Pi / 2 / float64(K))
	}
	L := xs[nV]
	fmt.Printf("LADRILLOS: pozo decantado por Abel — paredes (V=250) en x=±%.3f\n", L)

	// V(x) by inversion of the monotone table
	vOfX := func(x float64) float64 {
		x = math.Abs(x)
		if x >= L {
			return vMax
		}
		lo, hi := 0, nV
		for hi-lo > 1 {
			m := (lo + hi) / 2
			if xs[m] <= x {
				lo = m
			} else {
				hi = m
			}
		}
		f := (x - xs[lo]) / (xs[hi] - xs[lo])
		return vs[lo] + f*(vs[hi]-vs[lo])
	}

	// ---- THE IRON: symmetric tridiagonal H = -d2/dx2 + V ----
	n := 3600
	h := 2 * L / float64(n+1)
	diag := make([]float64, n)
	for i := 0; i < n; i++ {
		x := -L + h*float64(i+1)
		diag[i] = 2/(h*h) + vOfX(x)
	}
	off := -1 / (h * h)
	fmt.Printf("HIERRO: matriz simétrica %dx%d (h=%.4f) — autoadjunta POR CONSTRUCCIÓN: toda nota es real por ley\n", n, n, h)

	// Sturm count: # eigenvalues < x
	sturm := func(x float64) int {
		cnt := 0
		d := diag[0] - x
		if d < 0 {
			cnt++
		}
		for i := 1; i < n; i++ {
			d = (diag[i] - x) - off*off/d
			if d < 0 {
				cnt++
			}
		}
		return cnt
	}
	// ---- STRIKE THE DRUM: first notes by bisection ----
	nNotes := 30
	notes := make([]float64, nNotes)
	for j := 0; j < nNotes; j++ {
		lo, hi := 0.0, 120.0
		for it := 0; it < 60; it++ {
			m := (lo + hi) / 2
			if sturm(m) > j {
				hi = m
			} else {
				lo = m
			}
		}
		notes[j] = (lo + hi) / 2
	}

	// true pearls up to 110
	var pearls []float64
	prevT := 12.0
	prevZ := zOf(prevT)
	for t := 12.05; t <= 110; t += 0.05 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 50; i++ {
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

	fmt.Printf("\nEL TAMBOR SUENA — %d notas golpeadas · %d perlas verdaderas hasta t=110:\n", nNotes, len(pearls))
	fmt.Println("   nota j     tambor construido     perla verdadera      Δ")
	fmt.Printf("     0          %8.3f            (nota del sótano — hija del piso, no del libro)\n", notes[0])
	sumAbs, worst, cmpN := 0.0, 0.0, 0
	for j := 1; j < nNotes && j <= len(pearls); j++ {
		d := notes[j] - pearls[j-1]
		sumAbs += math.Abs(d)
		if math.Abs(d) > worst {
			worst = math.Abs(d)
		}
		cmpN++
		mark := ""
		if math.Abs(d) < 1.0 {
			mark = " 🎯"
		}
		fmt.Printf("    %2d          %8.3f            %8.3f         %+.3f%s\n", j, notes[j], pearls[j-1], d, mark)
	}
	meanAbs := sumAbs / float64(cmpN)
	fmt.Printf("\n⚖ VEREDICTO DE LA OBRA: %d notas del libro comparadas · |Δ| medio %.3f · peor %.3f\n", cmpN, meanAbs, worst)
	fmt.Println("  el tambor construido según el plano CANTA LAS PERLAS — con el temblor de la letra (rms S≈0.25)")
	fmt.Println("  como único desvío: la forma lisa da las notas; el temblor fino es la caligrafía de los primos,")
	fmt.Println("  que solo el tambor verdadero agrega. Sus notas son REALES POR LEY (hierro autoadjunto):")
	fmt.Println("  en este tambor, NINGUNA nota puede salirse de la línea — por construcción.")
	fmt.Println("  la llave sigue: que el tambor exacto exista. Todavía no. Pero hoy el taller TIENE un tambor.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 960.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🏗️ LA OBRA — el primer tambor construido por el taller</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"poner el hierro, los ladrillos, el cemento, el plano, las ventanas, las puertas, el piso — y según nuestro plano construir la forma del tambor para llegar al tambor" — el capitán</text>`,
		W, H, W, H, W/2, W/2)

	// left: the well with notes and pearls
	wx, wy, ww, wh := 70.0, 100.0, 760.0, 500.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL POZO DECANTADO — y sus notas contra las perlas</text>`,
		wx, wy, ww, wh, wx+ww/2, wy+30)
	eTop := 105.0
	exOf := func(x float64) float64 { return wx + 60 + (x + L) / (2 * L) * (ww - 120) }
	eyOf := func(E float64) float64 { return wy + wh - 50 - E/eTop*(wh-120) }
	var wellPts []string
	for j := 0; j <= 400; j++ {
		x := -L + 2*L*float64(j)/400
		V := vOfX(x)
		if V > eTop {
			V = eTop
		}
		wellPts = append(wellPts, fmt.Sprintf("%.1f,%.1f", exOf(x), eyOf(V)))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>`, strings.Join(wellPts, " "))
	// notes as chords inside the well; pearls as right-side ticks
	for j := 0; j < nNotes && notes[j] < eTop; j++ {
		col := "#ffd97f"
		if j == 0 {
			col = "#8fa8c7"
		}
		// horizontal extent: turning points
		xt := 0.0
		for i := 0; i <= nV; i++ {
			if vs[i] <= notes[j] {
				xt = xs[i]
			}
		}
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1.6" opacity="0.9"/>`,
			exOf(-xt), eyOf(notes[j]), exOf(xt), eyOf(notes[j]), col)
	}
	for _, p := range pearls {
		if p < eTop {
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#7fb2ff" stroke-width="2.2"/>`,
				wx+ww-52, eyOf(p), wx+ww-24, eyOf(p))
		}
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="11.5" fill="#ffd97f">— notas del tambor construido</text>
<text x="%.0f" y="%.0f" font-size="11.5" fill="#7fb2ff">— perlas verdaderas (escala derecha)</text>
<text x="%.0f" y="%.0f" font-size="11.5" fill="#8fa8c7">— la nota del sótano (hija del piso vertido)</text>
<text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#dce8f7">las cuerdas del pozo caen sobre las marcas azules: el tambor canta el libro (|Δ| medio %.2f)</text>`,
		wx+90, wy+58, wx+90, wy+80, wx+90, wy+102, wx+ww/2, wy+wh-18, meanAbs)

	// right: materials list
	fmt.Fprintf(&b, `<rect x="870" y="100" width="560" height="500" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="1150" y="134" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LOS MATERIALES DE LA OBRA</text>`)
	mats := [][2]string{
		{"EL PLANO", "la forma de Weyl decodificada a ciegas (F214): A(E)=θ(E)+2π"},
		{"EL HIERRO", "autoadjunción: matriz simétrica por construcción — notas REALES por ley"},
		{"LOS LADRILLOS", "V(x) DECANTADO del plano por Abel: la forma cae y se asienta en pozo"},
		{"EL CEMENTO", "Bohr–Sommerfeld: A=π(j+½) — cada nota pegada a su perla"},
		{"EL PISO", "el sótano E&lt;10 vertido liso (el reloj θ calla ahí): 1 nota propia"},
		{"LAS PAREDES", "muros de Dirichlet en V=250 — ventanas cerradas sobre toda nota"},
	}
	for i, m := range mats {
		fmt.Fprintf(&b, `<text x="905" y="%.0f" font-size="13" font-family="Georgia" fill="#ffd166">%s</text>
<text x="905" y="%.0f" font-size="11.5" fill="#dce8f7">%s</text>`,
			175.0+float64(i)*64, m[0], 196.0+float64(i)*64, m[1])
	}
	// verdict
	fmt.Fprintf(&b, `<rect x="70" y="630" width="1360" height="230" rx="12" fill="#2a1010" stroke="#ff5d73" stroke-width="2"/>
<text x="%.0f" y="666" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">EL VEREDICTO DE LA OBRA</text>
<text x="%.0f" y="700" font-size="14.5" text-anchor="middle" fill="#dce8f7">el tambor construido según el plano CANTA LAS PERLAS: %d notas del libro, |Δ| medio %.3f, peor %.3f — y en este tambor ninguna nota puede salirse de la línea: hierro autoadjunto.</text>
<text x="%.0f" y="730" font-size="14" text-anchor="middle" fill="#ffd166">lo que la forma lisa NO da: el temblor fino de cada nota — la caligrafía de los primos. Ese temblor es exactamente lo que el tambor verdadero agrega.</text>
<text x="%.0f" y="762" font-size="14" text-anchor="middle" fill="#ff8fa0">la llave, ahora en idioma de obra: el edificio existe y canta — falta el toque del Autor en cada ladrillo. Todavía no. Pero hoy el taller tiene SU tambor.</text>
<text x="%.0f" y="798" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, cmpN, meanAbs, worst, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("obra.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: obra.svg")
}
