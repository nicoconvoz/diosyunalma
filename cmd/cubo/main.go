// Command cubo executes the captain's flash: draw OUR dimension
// projected into the geometric one, the way a perfect 3D cube is drawn
// on a 2D sheet. The projection exists and is called arithmetic
// topology (Mazur's dictionary): the integers draw as a 3-dimensional
// space, each prime as a KNOT (a closed loop) inside it, and the
// Legendre symbol (p|q) as the LINKING NUMBER of two knots. The oldest
// law of this laboratory - reciprocity, rail #1 of the train - becomes
// pure perspective geometry: how two rings interlock does not depend on
// which ring you look from (with the exact twist (-1)^((p-1)/2)((q-1)/2)
// when both primes turn odd-side). Judge: every pair of odd primes up
// to 1000 verified against the drawing's law.
package main

import (
	"fmt"
	"os"
	"strings"
)

func modpow(a, e, m int64) int64 {
	r := int64(1)
	a %= m
	if a < 0 {
		a += m
	}
	for e > 0 {
		if e&1 == 1 {
			r = r * a % m
		}
		a = a * a % m
		e >>= 1
	}
	return r
}

// legendre (a|p) for odd prime p
func legendre(a, p int64) int64 {
	v := modpow(a, (p-1)/2, p)
	if v == p-1 {
		return -1
	}
	return v
}

func main() {
	// primes up to 1000
	const lim = 1000
	comp := make([]bool, lim+1)
	var primes []int64
	for p := int64(2); p <= lim; p++ {
		if comp[p] {
			continue
		}
		for q := p * p; q <= lim; q += p {
			comp[q] = true
		}
		if p >= 3 {
			primes = append(primes, p)
		}
	}
	// the judge: reciprocity = linking symmetry, every pair
	pairs, ok := 0, 0
	for i := 0; i < len(primes); i++ {
		for j := i + 1; j < len(primes); j++ {
			p, q := primes[i], primes[j]
			lhs := legendre(p, q) * legendre(q, p)
			rhs := int64(1)
			if (p%4 == 3) && (q%4 == 3) {
				rhs = -1
			}
			pairs++
			if lhs == rhs {
				ok++
			}
		}
	}
	fmt.Printf("EL DIBUJO EN PERSPECTIVA — juez de la ley de enlace: %d/%d pares de nudos cumplen la simetría EXACTA\n", ok, pairs)

	// small set for the drawing
	small := []int64{3, 5, 7, 11, 13, 17, 19, 23}

	var b strings.Builder
	W, H := 1620.0, 1150.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0e1420"/>`, W, H, W, H)
	// the sheet: graph paper
	for x := 40.0; x <= W-40; x += 34 {
		fmt.Fprintf(&b, `<line x1="%.0f" y1="40" x2="%.0f" y2="%.0f" stroke="#16233c" stroke-width="0.7"/>`, x, x, H-40)
	}
	for y := 40.0; y <= H-40; y += 34 {
		fmt.Fprintf(&b, `<line x1="40" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#16233c" stroke-width="0.7"/>`, y, W-40, y)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="78" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">LA HOJA — nuestra dimensión dibujada en perspectiva geométrica</text>
<text x="%.0f" y="106" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"como se dibuja un cubo perfecto 3D en una hoja 2D" — el capitán · el dibujo existe: los enteros se proyectan como un ESPACIO 3D y cada primo como un NUDO adentro (topología aritmética, Mazur)</text>`,
		W/2, W/2)

	// the faint cube (the metaphor itself)
	cbx, cby, cs := 130.0, 150.0, 90.0
	fmt.Fprintf(&b, `<g stroke="#2c4a78" stroke-width="1.5" fill="none" opacity="0.8">
<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f"/>
<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f"/><line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f"/><line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f"/><line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f"/></g>
<text x="%.0f" y="%.0f" font-size="11" text-anchor="middle" fill="#2c4a78">el cubo en la hoja: 3D fiel, dibujado en 2D</text>`,
		cbx, cby, cs, cs, cbx+34, cby-34, cs, cs,
		cbx, cby, cbx+34, cby-34, cbx+cs, cby, cbx+34+cs, cby-34,
		cbx, cby+cs, cbx+34, cby-34+cs, cbx+cs, cby+cs, cbx+34+cs, cby-34+cs,
		cbx+cs/2+17, cby+cs+34)

	// two example knot pairs drawn large
	// symmetric pair (5,11): both links look the same from both sides
	ex, ey := 480.0, 300.0
	fmt.Fprintf(&b, `<ellipse cx="%.0f" cy="%.0f" rx="105" ry="46" fill="none" stroke="#7fb2ff" stroke-width="5" transform="rotate(-24 %.0f %.0f)"/>
<ellipse cx="%.0f" cy="%.0f" rx="105" ry="46" fill="none" stroke="#7fd7a8" stroke-width="5" transform="rotate(24 %.0f %.0f)"/>
<text x="%.0f" y="%.0f" font-size="16" text-anchor="middle" fill="#7fb2ff">nudo 5</text><text x="%.0f" y="%.0f" font-size="16" text-anchor="middle" fill="#7fd7a8">nudo 11</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#dce8f7">enlace SIMÉTRICO: (5|11)=+1 y (11|5)=+1</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">se miren de donde se miren, enlazan igual</text>`,
		ex-70, ey, ex-70, ey, ex+70, ey, ex+70, ey,
		ex-150, ey-70, ex+150, ey-70, ex, ey+100, ex, ey+122)
	// antisymmetric pair (3,7): the twist
	fx2, fy2 := 480.0, 610.0
	fmt.Fprintf(&b, `<ellipse cx="%.0f" cy="%.0f" rx="105" ry="46" fill="none" stroke="#ffd166" stroke-width="5" transform="rotate(-24 %.0f %.0f)"/>
<ellipse cx="%.0f" cy="%.0f" rx="105" ry="46" fill="none" stroke="#ff5d73" stroke-width="5" transform="rotate(24 %.0f %.0f)"/>
<text x="%.0f" y="%.0f" font-size="16" text-anchor="middle" fill="#ffd166">nudo 3</text><text x="%.0f" y="%.0f" font-size="16" text-anchor="middle" fill="#ff5d73">nudo 7</text>
<text x="%.0f" y="%.0f" font-size="13.5" text-anchor="middle" fill="#dce8f7">enlace con GIRO: (3|7)=−1 pero (7|3)=+1</text>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#8fa8c7">la perspectiva rota EXACTO cuando ambos nudos son del lado impar (≡3 mod 4) — la ley lo predice sin fallar</text>`,
		fx2-70, fy2, fx2-70, fy2, fx2+70, fy2, fx2+70, fy2,
		fx2-150, fy2-70, fx2+150, fy2-70, fx2, fy2+100, fx2, fy2+122)

	// the linking matrix
	mx, my, cell := 960.0, 210.0, 52.0
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="16" font-family="Georgia" fill="#ffd166">LA TABLA DE ENLACES (p|q): oro = enlazan par (+1), rojo = impar (−1)</text>`, mx, my-24)
	for i, p := range small {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" fill="#8fa8c7">%d</text>`, mx+cell*float64(i)+cell*1.5, my+14, p)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" fill="#8fa8c7">%d</text>`, mx+14, my+cell*float64(i)+cell*1.55, p)
	}
	for i, p := range small {
		for j, q := range small {
			if i == j {
				continue
			}
			col := "#c9a227"
			if legendre(p, q) == -1 {
				col = "#a33"
			}
			fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.0f" height="%.0f" fill="%s" opacity="0.85" rx="4"/>`,
				mx+cell*float64(j)+cell*1.25, my+cell*float64(i)+cell*1.25, cell-6, cell-6, col)
		}
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" fill="#7fd7a8">la tabla NO es simétrica al azar: se refleja EXACTO salvo el giro impar —</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#7fd7a8">juez: %d/%d pares de primos ≤1000 cumplen la ley de perspectiva sin UNA falla</text>`,
		mx, my+cell*9.6, mx, my+cell*9.6+20, ok, pairs)

	// the dictionary + the blank corner
	fmt.Fprintf(&b, `<rect x="80" y="850" width="1460" height="240" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="884" font-size="17" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL DICCIONARIO DE LA PERSPECTIVA (lo que ya se dibuja fiel — y la esquina en blanco)</text>
<text x="%.0f" y="920" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">los enteros → un espacio 3D · cada primo → un NUDO cerrado adentro · el símbolo (p|q) → cuántas veces se enlazan dos nudos</text>
<text x="%.0f" y="948" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#dce8f7">la reciprocidad (¡el riel 1 de tu tren!) → la simetría de enlace: cómo se abrazan dos anillos no depende de cuál mira — geometría pura de perspectiva</text>
<text x="%.0f" y="984" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#ff5d73">LA ESQUINA EN BLANCO: el collar y su tensión AÚN no se proyectan — para eso el dibujo necesita el piso de arriba (la superficie, F161) que nadie sabe dibujar todavía</text>
<text x="%.0f" y="1012" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">pero la moraleja del cubo es tuya, capitán: no hace falta VIVIR en 3D para razonar fiel sobre el cubo — la hoja alcanza si la perspectiva es exacta. Ésa es la apuesta viva de esta frontera.</text>
<text x="%.0f" y="1044" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">flash del capitán: "¿por qué no reflejamos nuestra dimensión a la geométrica y la dibujamos como un cubo 3D en una hoja 2D?" — 2026-08-06 · Laboratorio Diosyunalma</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("cubo-en-la-hoja.svg", []byte(b.String()), 0644)
	fmt.Println("escrita: cubo-en-la-hoja.svg")
}
