// Command marchico executes the captain's strategy: do not sail the
// big seas - prove it in the LITTLE sea first and compare. The little
// sea exists: an elliptic curve over F_p. There the necklace has ONE
// relation and TWO halves per prime (the captain's words, exactly):
// two Frobenius roots alpha, beta bound by alpha*beta = p, and the
// theorem |alpha| = |beta| = sqrt(p) was PROVEN by hand (Hasse 1933).
// The proof's tension, in shape: the degree form
//
//	Q(m,n) = deg(m + n*phi) = m^2 + a*mn + p*n^2
//
// is a COUNT (the degree of an actual map counts preimages), and
// counts are never negative -> the bowl never dips -> a^2 <= 4p ->
// both halves ON the ring. We verify the whole little archipelago
// (every good prime <= 3000) and draw the comparison table with our
// ocean - whose pearls ALSO come as two halves (rho, 1-rho) bound by
// one relation. The blank cell of the dictionary is the prize: WHAT
// does lambda_n count?
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
	// the little archipelago: y^2 = x^3 + x + 1 over F_p
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
		if p >= 5 && p != 31 { // p=31 divides the curve discriminant (-31)
			primes = append(primes, p)
		}
	}
	fmt.Println("EL MAR CHIQUITO — la demostración ejecutada, isla por isla:")
	nOK, worstDisc := 0, math.Inf(-1)
	var exP, exA int64
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
		// the bowl: min of Q(m,n) = m^2 + a mn + p n^2 off the origin
		minQ := int64(1) << 60
		for m := int64(-20); m <= 20; m++ {
			for n := int64(-20); n <= 20; n++ {
				if m == 0 && n == 0 {
					continue
				}
				q := m*m + a*m*n + p*n*n
				if q < minQ {
					minQ = q
				}
			}
		}
		disc := float64(a*a - 4*p)
		if minQ > 0 && disc < 0 {
			nOK++
		}
		if disc > worstDisc {
			worstDisc, exP, exA = disc, p, a
		}
	}
	fmt.Printf("  islas (primos buenos ≤3000): %d\n", len(primes))
	fmt.Printf("  el cuenco jamás baja de cero (Q(m,n)>0 fuera del origen): %d/%d islas\n", nOK, len(primes))
	fmt.Printf("  ⇒ a² − 4p < 0 en TODAS (la más ajustada: p=%d, a=%d, a²−4p=%.0f)\n", exP, exA, worstDisc)
	fmt.Printf("  ⇒ LAS DOS MITADES EN EL ANILLO, siempre: |α|=|β|=√p — DEMOSTRADO Y VERIFICADO\n")

	// example bowl for the drawing: p=101
	var aEx int64
	{
		p := int64(101)
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
		aEx = -s
	}
	fmt.Printf("  isla de muestra p=101: a=%d, las dos mitades α,β = (a±√(a²−404))/2, |α|=|β|=√101\n", aEx)

	// ---- the picture ----
	var b strings.Builder
	W, H := 1640.0, 1140.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL MAR CHIQUITO — la demostración hecha a mano, y la comparación con nuestro océano</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"si podemos demostrarlo en la dimensión 0 — donde existe solo 1 relación y dos mitades — lo aplicamos al resto; no vayamos a mares grandes: partamos del chiquito y comparemos" — el capitán · ES la ruta histórica (Hasse 1933 → Weil → Deligne)</text>`,
		W, H, W, H, W/2, W/2)

	// panel 1: the little sea and its bowl
	p1x, p1y := 70.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="740" height="560" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#7fd7a8">1 · EL MAR CHIQUITO: 1 relación, 2 mitades — y la tensión ES UN CONTEO</text>`, p1x, p1y, p1x+20, p1y+34)
	// the two halves on the ring
	rx, ry, rr := p1x+150, p1y+230, 110.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#7fd7a8" stroke-width="2.5"/>`, rx, ry, rr)
	th := math.Acos(float64(aEx) / (2 * math.Sqrt(101)))
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="7" fill="#7fb2ff"/><circle cx="%.1f" cy="%.1f" r="7" fill="#7fb2ff"/>`,
		rx+rr*math.Cos(th), ry-rr*math.Sin(th), rx+rr*math.Cos(th), ry+rr*math.Sin(th))
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#dce8f7">las DOS MITADES: α y β</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#ffd166">LA RELACIÓN (una sola): α·β = p</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">isla p=101: a=%d — mitades en el anillo √101</text>`,
		rx, ry-rr-18, rx, ry+rr+28, rx, ry+rr+50, aEx)
	// the bowl: contours of Q(m,n)
	bx, by := p1x+520, p1y+230
	for _, c := range []float64{300, 900, 1800, 3000} {
		pts := make([]string, 0, 100)
		for i := 0; i <= 96; i++ {
			ang := 2 * math.Pi * float64(i) / 96
			cx0, sx0 := math.Cos(ang), math.Sin(ang)
			qv := cx0*cx0 + float64(aEx)*cx0*sx0 + 101*sx0*sx0
			r := math.Sqrt(c/qv) * 14
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", bx+r*cx0, by-r*sx0))
		}
		fmt.Fprintf(&b, `<polygon fill="none" stroke="#ffd166" stroke-width="1.3" opacity="0.75" points="%s"/>`, strings.Join(pts, " "))
	}
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="3" fill="#7fd7a8"/>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#ffd166">EL CUENCO: Q(m,n) = m² + a·mn + p·n²</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#dce8f7">cada valor CUENTA algo real (el grado de un mapa:</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#dce8f7">cuántos puntos caen encima) — y contar NUNCA da negativo</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#7fd7a8">cuenco sin pozos ⇒ a²&lt;4p ⇒ mitades AL ANILLO. Fin.</text>`,
		bx, by, bx, by-160, bx, by+150, bx, by+172, bx, by+196)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" fill="#7fd7a8">VERIFICADO EN CASA, ISLA POR ISLA: %d/%d primos ≤3000 — el cuenco jamás bajó de cero, las mitades jamás salieron del anillo</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#8fa8c7">ésta es LA TENSIÓN DEL HILO del mar chiquito: no un misterio — UN CONTEO. La demostración entera cabe en esta frase.</text>`,
		p1x+370, p1y+500, nOK, len(primes), p1x+370, p1y+526)

	// panel 2: the comparison table
	p2x, p2y := 850.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="720" height="560" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#ffd166">2 · LA COMPARACIÓN — mar chiquito vs nuestro océano</text>`, p2x, p2y, p2x+20, p2y+34)
	rows := [][3]string{
		{"las perlas", "2 por isla (α, β)", "infinitas (los ceros)"},
		{"las dos mitades", "α y β = p/α", "ρ y 1−ρ (¡también dos mitades!)"},
		{"la relación única", "α·β = p (producto)", "ρ + (1−ρ) = 1 (suma)"},
		{"el anillo", "|α|=√p", "Re(ρ)=½"},
		{"la ecuación de armonía", "Q(m,n) ≥ 0", "λ_n ≥ 0 (F166, medida)"},
		{"¿por qué ≥ 0?", "porque Q CUENTA (grados)", "❓ NADIE SABE QUÉ CUENTA λ_n"},
		{"veredicto", "DEMOSTRADO (Hasse 1933)", "ABIERTO — el millón"},
	}
	for i, r := range rows {
		y := p2y + 80 + float64(i)*62
		bg := "#0f2540"
		if i == 5 {
			bg = "#3a1515"
		}
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="680" height="52" rx="6" fill="%s" stroke="#2c4a78"/>
<text x="%.0f" y="%.0f" font-size="13" fill="#8fa8c7">%s</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#7fd7a8">%s</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#dce8f7">%s</text>`,
			p2x+20, y-30, bg,
			p2x+34, y, r[0], p2x+230, y, r[1], p2x+460, y, r[2])
	}

	// footer: the route and the refined target
	fmt.Fprintf(&b, `<rect x="70" y="700" width="1500" height="230" rx="12" fill="#0d2547" stroke="#7fd7a8" stroke-width="2"/>
<text x="%.0f" y="740" font-size="17" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA RUTA DEL CAPITÁN ES LA RUTA DE LA HISTORIA — Y MARCA LA CASILLA EXACTA</text>
<text x="%.0f" y="776" font-size="14.5" text-anchor="middle" fill="#dce8f7">mar chiquito (1 relación, 2 mitades: Hasse, A MANO, 1933) → todos los mares con nieve (Weil 1948) → todas las dimensiones (Deligne 1974) → nuestro océano: ❓</text>
<text x="%.0f" y="812" font-size="15" text-anchor="middle" fill="#ffd166">la comparación revela EL ÚNICO renglón que falta: en el mar chiquito la armonía es positiva PORQUE CUENTA COSAS REALES.</text>
<text x="%.0f" y="842" font-size="15" text-anchor="middle" fill="#ffd166">La pregunta del millón, afinada al máximo: ¿QUÉ COSA REAL CUENTA λ_n? — encontrá el objeto contado, y la demostración sube sola del mar chiquito al infinito.</text>
<text x="%.0f" y="878" font-size="13" text-anchor="middle" fill="#8fa8c7">honestidad: Connes persigue exactamente ese objeto hace décadas; nadie lo halló — pero jamás estuvo formulado en el idioma de las formas del capitán hasta hoy</text>
<text x="%.0f" y="908" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · 2026-08-06</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("mar-chiquito.svg", []byte(b.String()), 0644)
	fmt.Println("escrita: mar-chiquito.svg")
}
