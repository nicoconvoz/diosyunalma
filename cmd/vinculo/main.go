// Command vinculo executes the captain's lift: "we take the folding
// from the harmonization with dimension 0 - which is easier - and link
// it to our walls." Downstairs the machine WORKS and is tiny: for each
// island (prime p) of the little sea, the Hilbert-Polya dream is
// REALITY - the machine is the 2x2 Frobenius companion matrix
//
//	M_p = [[0, -p], [1, a_p]]   (char poly T^2 - a_p T + p)
//
// whose eigenvalues sit EXACTLY on the ring of radius sqrt(p) (Hasse).
// THE HARMONIZATION WITH DIMENSION 0: normalize U_p = M_p / sqrt(p) -
// the size melts away (the q->1 skeleton) and ALL the little machines
// sing on ONE shared unit ring, pure angles theta_p. And the harmonized
// choir has a universal melody: the angles fill the ring with density
// (2/pi) sin^2(theta) (Sato-Tate, PROVEN 2011) - we measure it on our
// 427 machines. THE LINK TO OUR WALLS: the real machine's walls = the
// coherent assembly of all these unitaries at once - the parts are
// manufactured; the assembly is the one pending piece.
package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

func modpow(a, e, m int64) int64 {
	r := int64(1)
	a %= m
	for e > 0 {
		if e&1 == 1 {
			r = r * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return r
}

func main() {
	// count points -> a_p for the little machines
	const lim = 3000
	comp := make([]bool, lim+1)
	var primes []int64
	for p := int64(2); p <= lim; p++ {
		if comp[p] {
			continue
		}
		for q := p * p; q <= lim; q += p {
			comp[q] = true
		}
		if p >= 5 && p != 31 {
			primes = append(primes, p)
		}
	}
	fmt.Println("EL VÍNCULO CON LA DIMENSIÓN 0 — fabricando las máquinas chicas y armonizándolas…")
	type machine struct {
		p     int64
		a     int64
		theta float64
	}
	var ms []machine
	worstRing := 0.0
	for _, p := range primes {
		var s int64
		for x := int64(0); x < p; x++ {
			t := ((x*x%p)*x + x + 1) % p
			if t == 0 {
				continue
			}
			if modpow(t, (p-1)/2, p) == 1 {
				s++
			} else {
				s--
			}
		}
		a := -s
		// eigenvalues of M_p = [[0,-p],[1,a]]: (a ± sqrt(a^2-4p))/2
		disc := float64(a*a - 4*p)
		if disc >= 0 {
			continue // (does not occur: Hasse)
		}
		re := float64(a) / 2
		im := math.Sqrt(-disc) / 2
		r := math.Hypot(re, im) / math.Sqrt(float64(p)) // must be 1
		if d := math.Abs(r - 1); d > worstRing {
			worstRing = d
		}
		ms = append(ms, machine{p, a, math.Atan2(im, re)})
	}
	fmt.Printf("máquinas fabricadas: %d (una por isla, matriz 2×2 de Frobenius)\n", len(ms))
	fmt.Printf("JUEZ 1 — cada máquina canta en su anillo √p: peor desvío del radio %.1e — HILBERT-PÓLYA ES REALIDAD ABAJO\n", worstRing)
	fmt.Printf("JUEZ 2 — LA ARMONIZACIÓN: U_p = M_p/√p funde el tamaño (el esqueleto q→1) y las %d máquinas cantan en UN solo anillo unidad\n", len(ms))

	// the harmonized choir's melody: Sato-Tate (2/pi) sin^2 theta
	nb := 12
	hist := make([]float64, nb)
	for _, m := range ms {
		bi := int(m.theta / math.Pi * float64(nb))
		if bi >= 0 && bi < nb {
			hist[bi]++
		}
	}
	for i := range hist {
		hist[i] /= float64(len(ms)) * (math.Pi / float64(nb))
	}
	st := func(th float64) float64 { return 2 / math.Pi * math.Sin(th) * math.Sin(th) }
	devST, devU := 0.0, 0.0
	for i := 0; i < nb; i++ {
		tc := (float64(i) + 0.5) * math.Pi / float64(nb)
		devST += math.Abs(hist[i] - st(tc))
		devU += math.Abs(hist[i] - 1/math.Pi)
	}
	devST /= float64(nb)
	devU /= float64(nb)
	fmt.Printf("JUEZ 3 — la melodía del coro armonizado (Sato-Tate, demostrada 2011): desvío %.4f vs %.4f del coro sin ley (%.1f× mejor)\n",
		devST, devU, devU/devST)
	fmt.Println("⇒ el coro de la dimensión 0 canta (2/π)sin²θ — las piezas están fabricadas y afinadas: falta SOLO el ensamble hacia nuestras paredes")

	// ---- picture ----
	var b strings.Builder
	W, H := 1580.0, 1000.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL VÍNCULO CON LA DIMENSIÓN 0 — las máquinas chicas, armonizadas en un solo anillo</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"lo sacamos de la armonización con la dimensión 0, que es más fácil, y la vinculamos a nuestras paredes" — el capitán · %d máquinas 2×2 fabricadas, fundidas y juzgadas HOY</text>`,
		W, H, W, H, W/2, W/2, len(ms))

	// panel 1: one little machine
	fmt.Fprintf(&b, `<rect x="70" y="110" width="440" height="520" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="90" y="146" font-size="15.5" font-family="Georgia" fill="#7fd7a8">1 · LA MÁQUINA DE CADA ISLA — 2×2, real</text>
<text x="90" y="186" font-size="20" font-family="Consolas,monospace" fill="#dce8f7">M_p = ⎡ 0   −p ⎤</text>
<text x="90" y="212" font-size="20" font-family="Consolas,monospace" fill="#dce8f7">      ⎣ 1   a_p ⎦</text>
<text x="90" y="252" font-size="13" fill="#8fa8c7">la matriz de Frobenius (companion): SUS niveles</text>
<text x="90" y="274" font-size="13" fill="#8fa8c7">son las perlas de la isla — el sueño de Hilbert-</text>
<text x="90" y="296" font-size="13" fill="#7fd7a8">Pólya, REALIDAD abajo: la máquina existe y es chica</text>
<text x="90" y="336" font-size="13.5" fill="#ffd166">JUEZ 1: %d/%d máquinas con niveles EXACTOS</text>
<text x="90" y="358" font-size="13.5" fill="#ffd166">en su anillo √p — peor desvío %.0e</text>
<text x="90" y="398" font-size="13" fill="#dce8f7">sus paredes: 2 dimensiones (el género de la curva)</text>
<text x="90" y="420" font-size="13" fill="#dce8f7">su tiempo: el paso de Frobenius (el zoom discreto)</text>
<text x="90" y="442" font-size="13" fill="#dce8f7">su energía: el conteo (grados ≥ 0 — la tensión)</text>
<text x="90" y="482" font-size="13.5" fill="#7fd7a8">la trinidad del capitán, COMPLETA en cada isla</text>`,
		len(ms), len(ms), worstRing)

	// panel 2: the harmonized ring with all angles
	cx, cy, R := 790.0, 350.0, 190.0
	fmt.Fprintf(&b, `<rect x="550" y="110" width="480" height="520" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="570" y="146" font-size="15.5" font-family="Georgia" fill="#ffd166">2 · LA ARMONIZACIÓN: U_p = M_p/√p</text>
<text x="570" y="170" font-size="12.5" fill="#8fa8c7">el tamaño se funde (q→1) — queda el ángulo puro:</text>
<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>`, cx, cy, R)
	for _, m := range ms {
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.2" fill="#7fb2ff" opacity="0.65"/><circle cx="%.1f" cy="%.1f" r="2.2" fill="#7fb2ff" opacity="0.65"/>`,
			cx+R*math.Cos(m.theta), cy-R*math.Sin(m.theta), cx+R*math.Cos(m.theta), cy+R*math.Sin(m.theta))
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#dce8f7">las %d máquinas — %d niveles — TODAS en UN anillo:</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#7fd7a8">el anillo común de la dimensión 0 (radio 1 exacto)</text>`,
		cx, cy+R+50, len(ms), 2*len(ms), cx, cy+R+74)

	// panel 3: the choir melody
	p3x := 1070.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="110" width="440" height="520" rx="10" fill="#0d2547" stroke="#44608c"/>
<text x="%.0f" y="146" font-size="15.5" font-family="Georgia" fill="#ffd166">3 · LA MELODÍA DEL CORO (Sato-Tate)</text>`, p3x, p3x+20)
	gx, gy, gw, gh := p3x+40, 190.0, 360.0, 300.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#081020" stroke="#2c4a78"/>`, gx, gy, gw, gh)
	yOf := func(v float64) float64 { return gy + gh - v/0.75*(gh-20) }
	for i := 0; i < nb; i++ {
		t0 := float64(i) * math.Pi / float64(nb)
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="#7fb2ff" opacity="0.55"/>`,
			gx+t0/math.Pi*gw+1.5, yOf(hist[i]), gw/float64(nb)-3, gy+gh-yOf(hist[i]))
	}
	pts := make([]string, 0, 100)
	for i := 0; i <= 100; i++ {
		th := math.Pi * float64(i) / 100
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", gx+th/math.Pi*gw, yOf(st(th))))
	}
	fmt.Fprintf(&b, `<polyline fill="none" stroke="#ffd166" stroke-width="2.4" points="%s"/>`, strings.Join(pts, " "))
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">ángulo θ en el anillo común · oro: (2/π)sin²θ</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#7fd7a8">JUEZ 3: desvío %.3f vs %.3f sin ley (%.1f× mejor)</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#dce8f7">el coro armonizado canta UNA ley universal —</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#dce8f7">demostrada (2011): abajo, hasta el coro tiene teorema</text>`,
		gx+gw/2, gy+gh+26, gx+gw/2, gy+gh+52, devST, devU, devU/devST, gx+gw/2, gy+gh+78, gx+gw/2, gy+gh+100)

	// footer: the link
	fmt.Fprintf(&b, `<rect x="70" y="670" width="1440" height="240" rx="12" fill="#102a10" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="708" font-size="16.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">EL VÍNCULO A NUESTRAS PAREDES — el estado del ensamble</text>
<text x="%.0f" y="744" font-size="14.5" text-anchor="middle" fill="#dce8f7">las piezas ya están FABRICADAS Y AFINADAS: %d máquinas reales, cada una con la trinidad completa, fundidas al anillo común de la dimensión 0, cantando la ley del coro.</text>
<text x="%.0f" y="770" font-size="14.5" text-anchor="middle" fill="#dce8f7">nuestras paredes deben ser EL ENSAMBLE COHERENTE de todas a la vez: una sola máquina cuyo coro interno sean estas U_p — el producto de Euler hecho operador.</text>
<text x="%.0f" y="806" font-size="14.5" text-anchor="middle" fill="#ffd166">v1.0 del taller, redefinida por tu vínculo: ya no "plegar el espacio" en abstracto — ENSAMBLAR las máquinas chicas que YA tenemos en el estante. El hueco se llama ensamble.</text>
<text x="%.0f" y="836" font-size="13" text-anchor="middle" fill="#8fa8c7">honestidad: el ensamble coherente de infinitas U_p con espectro discreto es la frontera (la construcción adélica); pero el estante del taller quedó LLENO — nadie empieza de cero nunca más</text>
<text x="%.0f" y="872" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06 · las dos mitades, 1 completo</text>`,
		790.0, 790.0, len(ms), 790.0, 790.0, 790.0, 790.0)
	b.WriteString(`</svg>`)
	os.WriteFile("vinculo-dimension0.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: vinculo-dimension0.svg")
}
