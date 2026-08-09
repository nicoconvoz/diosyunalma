// Command atomo draws THE ATOM OF THE PRIMES (F149): the captain saw
// that the curlicues could draw it, and they do. The partial-sum path
// of zeta(1/2+it) (Euler-Maclaurin corrected) is a curlicue converging
// to a point - and at the ZEROS that point is the ORIGIN: the orbit
// CLOSES. Off a zero it stays open. Bohr's quantization, drawn: closed
// orbits are the levels; in the quantum-chaos dictionary the primes are
// the periodic orbits (period ln p) of the atom nobody has built.
package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"strings"
)

func main() {
	var b strings.Builder
	W, H := 1240.0, 600.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="620" y="36" font-size="23" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL ÁTOMO DE LOS PRIMOS — los orbitales dibujados por los curlicues (F149)</text>
<text x="620" y="60" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el camino de sumas parciales de &#950;(&#189;+it): en los CEROS la órbita CIERRA en el origen (nivel cuantizado);</text>
<text x="620" y="78" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">fuera de un cero queda abierta — la cuantización de Bohr, vista; los primos son las órbitas periódicas (período ln p)</text>`,
		W, H, W, H)

	panels := []struct {
		t   float64
		tag string
	}{
		{14.134725142, "t = &#947;&#8321; = 14.1347... — ÓRBITA CERRADA (primer nivel del átomo)"},
		{15.0, "t = 15.0 — órbita ABIERTA (no hay nivel: el átomo no resuena)"},
		{21.022039639, "t = &#947;&#8322; = 21.0220... — ÓRBITA CERRADA (segundo nivel)"},
	}
	for pi, pn := range panels {
		ox := 210.0 + float64(pi)*410
		oy := 330.0
		s := complex(0.5, pn.t)
		var sum complex128
		type pt struct{ x, y float64 }
		var path []pt
		minx, maxx, miny, maxy := 0.0, 0.0, 0.0, 0.0
		for n := 1; n <= 700; n++ {
			sum += cmplx.Exp(-s * complex(math.Log(float64(n)), 0))
			nf := complex(float64(n), 0)
			corr := cmplx.Exp((1 - s) * cmplx.Log(nf)) / (1 - s)
			p := sum - corr
			x, y := real(p), imag(p)
			path = append(path, pt{x, y})
			if x < minx {
				minx = x
			}
			if x > maxx {
				maxx = x
			}
			if y < miny {
				miny = y
			}
			if y > maxy {
				maxy = y
			}
		}
		scale := 300.0 / math.Max(maxx-minx, math.Max(maxy-miny, 0.1))
		cx, cy := (minx+maxx)/2, (miny+maxy)/2
		pts := make([]string, 0, len(path))
		for _, p := range path {
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", ox+(p.x-cx)*scale, oy-(p.y-cy)*scale))
		}
		// the origin (the nucleus)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="5" fill="#ffd166"/>`, ox+(0-cx)*scale, oy-(0-cy)*scale)
		fmt.Fprintf(&b, `<polyline fill="none" stroke="#7fb2ff" stroke-width="0.9" points="%s"/>`, strings.Join(pts, " "))
		// the landing point
		last := path[len(path)-1]
		col := "#7fd7a8"
		if pi == 1 {
			col = "#ff5d73"
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4" fill="%s"/>`, ox+(last.x-cx)*scale, oy-(last.y-cy)*scale, col)
		fmt.Fprintf(&b, `<text x="%.0f" y="530" font-size="12" text-anchor="middle" font-family="Georgia" fill="#ffd166">%s</text>`, ox, pn.tag)
	}
	fmt.Fprintf(&b, `<text x="620" y="565" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">punto dorado = el núcleo (origen) · verde = la órbita aterriza EN el núcleo (cero) · rojo = aterriza lejos (no-cero) — el átomo de Hilbert-Pólya, esbozado por los curlicues del capitán</text>`)
	b.WriteString(`</svg>`)
	os.WriteFile("atomo-primos.svg", []byte(b.String()), 0644)
	fmt.Println("escrito: atomo-primos.svg")
}
