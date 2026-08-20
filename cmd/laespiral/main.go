package main

// LA ESPIRAL - the captain's flash: a cosine in spiral, converging to a point.
// The partial sums S_X = sum_{n<=X} n^{-1/2-it} walk the complex plane: each n
// is one cosine step of length 1/sqrt(n) and angle -t*log(n). The walk advances,
// then curls into an eye. Euler-Maclaurin says the eye IS zeta(1/2+it):
//   zeta(s) = S_X + X^{1-s}/(s-1) - X^{-s}/2 + O(|s| X^{-3/2})
// So at a zero the spiral winds around the ORIGIN; off a zero it winds around a
// displaced point. Two panels + the numeric check that the eye equals |Z(t)|.
//
// Exploration mode: measure first, judge after.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// ---- fresh Riemann-Siegel (same minimal construction as always) --------------

func theta(t float64) float64 {
	return t/2*math.Log(t/(2*math.Pi)) - t/2 - math.Pi/8 + 1/(48*t)
}

func zetaZ(t float64) float64 {
	th := theta(t)
	u := math.Sqrt(t / (2 * math.Pi))
	N := int(u)
	s := 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		s += math.Cos(th-t*math.Log(fn)) / math.Sqrt(fn)
	}
	s *= 2
	p := u - float64(N)
	c0 := math.Cos(2*math.Pi*(p*p-p-1.0/16)) / math.Cos(2*math.Pi*p)
	sg := 1.0
	if N%2 == 0 {
		sg = -1
	}
	return s + sg*math.Pow(2*math.Pi/t, 0.25)*c0
}

func ceros(t0 float64, cuantos int) []float64 {
	var g []float64
	a, za := t0, zetaZ(t0)
	for b := t0 + 0.02; len(g) < cuantos; b += 0.02 {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 45; i++ {
				m := (lo + hi) / 2
				if (zlo < 0) != (zetaZ(m) < 0) {
					hi = m
				} else {
					lo = m
				}
			}
			g = append(g, (lo+hi)/2)
		}
		a, za = b, zb
	}
	return g
}

// ---- the walk of cosines ------------------------------------------------------

type punto struct{ x, y float64 }

// walk returns the partial-sum path: point k = sum_{n=1..k} e^{-i t ln n}/sqrt(n),
// starting at the origin (k=0).
func walk(t float64, N int) []punto {
	ps := make([]punto, N+1)
	x, y := 0.0, 0.0
	ps[0] = punto{0, 0}
	for n := 1; n <= N; n++ {
		fn := float64(n)
		a := t * math.Log(fn)
		r := 1 / math.Sqrt(fn)
		x += r * math.Cos(a)
		y -= r * math.Sin(a)
		ps[n] = punto{x, y}
	}
	return ps
}

// ojo applies the Euler-Maclaurin correction at cutoff X: the stabilized center
// of the spiral, which should equal zeta(1/2+it).
func ojo(t float64, ps []punto, X int) punto {
	lX := math.Log(float64(X))
	sx, sy := math.Sqrt(float64(X)), 0.0
	_ = sy
	// X^{1-s}/(s-1) with s = 1/2+it:  sqrt(X) e^{-i t lnX} / (-1/2 + i t)
	nr, ni := sx*math.Cos(t*lX), -sx*math.Sin(t*lX)
	dr, di := -0.5, t
	den := dr*dr + di*di
	cr := (nr*dr + ni*di) / den
	ci := (ni*dr - nr*di) / den
	// - X^{-s}/2 = - e^{-i t lnX} / (2 sqrt(X))
	hr := -math.Cos(t*lX) / (2 * sx)
	hi := math.Sin(t*lX) / (2 * sx)
	return punto{ps[X].x + cr + hr, ps[X].y + ci + hi}
}

func abs2(p punto) float64 { return math.Hypot(p.x, p.y) }

func main() {
	fmt.Println("🌀 LA ESPIRAL — el coseno enroscado que converge a un punto")
	fmt.Println("   ¿tenemos una cuerda? · codo a codo · modo explorador")
	fmt.Println()

	gs := ceros(100, 2)
	tz := gs[0]               // exact height of a real zero
	tm := (gs[0] + gs[1]) / 2 // middle of the gap: NOT a zero
	N := 600

	wz := walk(tz, N)
	wm := walk(tm, N)

	fmt.Printf("§1 · las dos alturas\n")
	fmt.Printf("   cero real:      γ = %.6f   (|Z| = %.6f — cero de verdad)\n", tz, math.Abs(zetaZ(tz)))
	fmt.Printf("   medio del hueco: t = %.6f   (|Z| = %.4f — acá NO hay cero)\n\n", tm, math.Abs(zetaZ(tm)))

	fmt.Println("§2 · el ojo de la espiral contra zeta — la prueba de la cuerda")
	fmt.Println("   ojo(X) = suma parcial hasta X + corrección de Euler–Maclaurin.")
	fmt.Println("   Si la espiral es un instrumento de verdad, |ojo| debe clavar |Z(t)|:")
	fmt.Println()
	fmt.Println("      X    |ojo| en el cero   |ojo| en el hueco")
	for _, X := range []int{150, 300, 450, 600} {
		fmt.Printf("   %4d        %.6f           %.6f\n", X, abs2(ojo(tz, wz, X)), abs2(ojo(tm, wm, X)))
	}
	oz, om := ojo(tz, wz, N), ojo(tm, wm, N)
	zm := math.Abs(zetaZ(tm))
	fmt.Printf("\n   referencia Riemann–Siegel:  |Z(γ)| = %.6f   |Z(hueco)| = %.6f\n", math.Abs(zetaZ(tz)), zm)
	if abs2(oz) < 0.01 && math.Abs(abs2(om)-zm) < 0.02*zm+0.01 {
		fmt.Println("   ⟹ CUERDA: el ojo cae en el ORIGEN a la altura del cero, y clava |Z|")
		fmt.Println("     en el hueco. La espiral de cosenos LEE los ceros con los ojos.")
	} else {
		fmt.Println("   ⟹ el ojo no clava la referencia: la cuerda está floja — se reporta igual.")
	}

	dibujar(wz, wm, oz, om, tz, tm, N)
}

// ---- the plate ----------------------------------------------------------------

func dibujar(wz, wm []punto, oz, om punto, tz, tm float64, N int) {
	// shared scale over both walks
	minx, maxx := math.Inf(1), math.Inf(-1)
	miny, maxy := math.Inf(1), math.Inf(-1)
	for _, w := range [][]punto{wz, wm} {
		for _, p := range w {
			minx, maxx = math.Min(minx, p.x), math.Max(maxx, p.x)
			miny, maxy = math.Min(miny, p.y), math.Max(maxy, p.y)
		}
	}
	pad := 0.35
	minx, maxx, miny, maxy = minx-pad, maxx+pad, miny-pad, maxy+pad
	lado := 470.0
	esc := lado / math.Max(maxx-minx, maxy-miny)
	cxm := (minx + maxx) / 2
	cym := (miny + maxy) / 2

	var b strings.Builder
	W, H := 1060, 640
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	t := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	t(530, 34, 21, "#e8e2d4", "middle", "LA ESPIRAL — el coseno enroscado que converge a un punto")
	t(530, 56, 12, "#8a93a6", "middle", fmt.Sprintf("la caminata de las sumas parciales Σ n^(−1/2−it), un paso por cada n hasta N = %d · misma escala en los dos paneles", N))

	panel := func(x0 float64, w []punto, o punto, titulo, sub, colOjo string) {
		cx := x0 + lado/2
		cy := 95 + lado/2
		fmt.Fprintf(&b, `<rect x="%.0f" y="95" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x0, lado, lado)
		px := func(p punto) (float64, float64) { return cx + (p.x-cxm)*esc, cy + (p.y-cym)*esc }
		// the walk, colored early->late
		for i := 1; i < len(w); i++ {
			f := float64(i) / float64(len(w)-1)
			r := int(90 + 165*f)
			g := int(140 - 60*f)
			bb := int(220 - 60*f)
			x1, y1 := px(w[i-1])
			x2, y2 := px(w[i])
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="rgb(%d,%d,%d)" stroke-width="1.1" opacity="0.85"/>`, x1, y1, x2, y2, r, g, bb)
		}
		// origin cross
		ox, oy := px(punto{0, 0})
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#e8e2d4" stroke-width="1"/>`, ox-7, oy, ox+7, oy)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#e8e2d4" stroke-width="1"/>`, ox, oy-7, ox, oy+7)
		t(ox+10, oy-6, 11, "#e8e2d4", "start", "origen (el cero)")
		// the eye
		ex, ey := px(o)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4.5" fill="none" stroke="%s" stroke-width="1.6"/>`, ex, ey, colOjo)
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="1.6" fill="%s"/>`, ex, ey, colOjo)
		t(x0+lado/2, 88, 14, "#e8e2d4", "middle", titulo)
		t(x0+lado/2, 95+lado+22, 12, "#8a93a6", "middle", sub)
	}

	panel(30, wz, oz, fmt.Sprintf("a la altura del CERO real  γ = %.4f", tz),
		fmt.Sprintf("el ojo cae en el ORIGEN: |ojo| = %.4f — la espiral se muerde la cola en el cero exacto", abs2(oz)), "#7ee0c0")
	panel(560, wm, om, fmt.Sprintf("en el MEDIO del hueco  t = %.4f", tm),
		fmt.Sprintf("mismo coseno, otra altura: el ojo queda DESPLAZADO, |ojo| = %.4f = |ζ(1/2+it)|", abs2(om)), "#f2a6c0")

	t(530, float64(H)-12, 12, "#c9b458", "middle",
		"la cuerda: el punto donde la espiral converge ES ζ(1/2+it) — un cero de Riemann es la altura donde el ojo pisa el origen · Todavía no")

	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "la-espiral.svg")
	os.WriteFile(ruta, []byte(b.String()), 0644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", ruta)
}
