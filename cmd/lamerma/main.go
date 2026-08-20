package main

// LA MERMA - strict order: stop hunting anything. Investigate eps_k = 1 - q_k.
// Decisions fixed BEFORE looking:
//   contraction phase: cycles with eps > 0 before the first sustained
//     crossing into expansion (3 consecutive eps < 0), capped at k = 60
//   three candidate laws, all fit, mechanical winner by R^2 of the linearized
//     forms: POWER (ln eps vs ln k), EXPONENTIAL (ln eps vs k),
//     SHIFTED HYPERBOLA (1/eps vs k linear <=> eps = c/(k+a))
//   geometric decomposition: if r ~ n^{-p} locally and a cycle spans
//     [n_k, n_k + dn_k], then eps_k ~ p * dn_k / n_k, so the local exponent
//     p_k = eps_k * n_k / dn_k is measured per cycle - constant p would NAME
//     the law's origin
//   landmarks: k* = first sustained crossing, and n(k*)/tau
// The banned numbers appear nowhere.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

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

func cerosDesde(t0 float64, cuantos int) []float64 {
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

type punto struct{ x, y float64 }

func walk(t, sigma float64, N int) []punto {
	ps := make([]punto, N+1)
	x, y := 0.0, 0.0
	for n := 1; n <= N; n++ {
		fn := float64(n)
		a := t * math.Log(fn)
		r := math.Pow(fn, -sigma)
		x += r * math.Cos(a)
		y -= r * math.Sin(a)
		ps[n] = punto{x, y}
	}
	return ps
}

func ojoEM(t, sigma float64, ps []punto, X int) punto {
	lX := math.Log(float64(X))
	mod := math.Pow(float64(X), 1-sigma)
	nr, ni := mod*math.Cos(t*lX), -mod*math.Sin(t*lX)
	dr, di := sigma-1, t
	den := dr*dr + di*di
	cr := (nr*dr + ni*di) / den
	ci := (ni*dr - nr*di) / den
	m2 := math.Pow(float64(X), -sigma) / 2
	hr := -m2 * math.Cos(t*lX)
	hi := m2 * math.Sin(t*lX)
	return punto{ps[X].x + cr + hr, ps[X].y + ci + hi}
}

// ciclos: angular cycle boundaries plus per-cycle mean radius and span.
type ciclo struct {
	nIni, nFin int
	rMedio     float64
}

func ciclos(ps []punto, C punto) []ciclo {
	prev, th0, th := 0.0, 0.0, 0.0
	base := math.NaN()
	kPrev := 0
	ini := 1
	var cs []ciclo
	var suma float64
	var cnt int
	for n := 1; n < len(ps); n++ {
		dx, dy := ps[n].x-C.x, ps[n].y-C.y
		r := math.Hypot(dx, dy)
		a := math.Atan2(dy, dx)
		if n == 1 {
			th = a
			base = a
			th0 = a
			_ = th0
		} else {
			d := a - prev
			for d > math.Pi {
				d -= 2 * math.Pi
			}
			for d < -math.Pi {
				d += 2 * math.Pi
			}
			th += d
		}
		prev = a
		k := int(math.Floor(math.Abs(th-base) / (2 * math.Pi)))
		if k > kPrev {
			if cnt > 0 {
				cs = append(cs, ciclo{ini, n, suma / float64(cnt)})
			}
			ini, suma, cnt = n, 0, 0
			kPrev = k
		}
		suma += r
		cnt++
	}
	return cs
}

// ajusteLineal returns slope, intercept, R2 of y vs x.
func ajusteLineal(x, y []float64) (m, b, r2 float64) {
	var sx, sy, sxx, sxy, syy, n float64
	for i := range x {
		sx += x[i]
		sy += y[i]
		sxx += x[i] * x[i]
		sxy += x[i] * y[i]
		syy += y[i] * y[i]
		n++
	}
	m = (n*sxy - sx*sy) / (n*sxx - sx*sx)
	b = (sy - m*sx) / n
	ssr := 0.0
	for i := range x {
		d := y[i] - (m*x[i] + b)
		ssr += d * d
	}
	sst := syy - sy*sy/n
	return m, b, 1 - ssr/sst
}

type resultado struct {
	nom               string
	eps               []float64
	nMed, dN          []float64
	alfaPot, r2Pot    float64
	tasaExp, r2Exp    float64
	cHip, aHip, r2Hip float64
	pLocalMed         float64
	kCruce            int
	nCruceSobreTau    float64
}

func mediana(v []float64) float64 {
	w := append([]float64(nil), v...)
	for i := 1; i < len(w); i++ {
		for j := i; j > 0 && w[j] < w[j-1]; j-- {
			w[j], w[j-1] = w[j-1], w[j]
		}
	}
	if len(w) == 0 {
		return math.NaN()
	}
	return w[len(w)/2]
}

func main() {
	fmt.Println("📉 LA MERMA — ε_k = 1 − q_k bajo el microscopio (orden nueva)")
	fmt.Println("   tres leyes candidatas, ganador por R² · descomposición ε = p·Δn/n")
	fmt.Println()

	gs := cerosDesde(1000, 2)
	g1, g2 := gs[0], gs[1]
	g5 := cerosDesde(5000, 1)[0]
	casos := []struct {
		nom      string
		t, sigma float64
	}{
		{"cero γ≈1000", g1, 0.5},
		{"cerca del cero (+0,15)", g1 + 0.15, 0.5},
		{"control: medio del hueco", (g1 + g2) / 2, 0.5},
		{"cero γ≈5000", g5, 0.5},
		{"instrumento σ=0,3", g1, 0.3},
		{"instrumento σ=0,7", g1, 0.7},
	}

	var res []resultado
	for _, cs := range casos {
		tau := cs.t / (2 * math.Pi)
		N := int(8 * tau)
		ps := walk(cs.t, cs.sigma, N)
		C := ojoEM(cs.t, cs.sigma, ps, N)
		cls := ciclos(ps, C)

		var eps, nMed, dN []float64
		for i := 1; i < len(cls); i++ {
			q := cls[i].rMedio / cls[i-1].rMedio
			eps = append(eps, 1-q)
			nMed = append(nMed, float64(cls[i].nIni+cls[i].nFin)/2)
			dN = append(dN, float64(cls[i].nFin-cls[i].nIni))
		}
		// k* = first sustained crossing (3 consecutive eps < 0)
		kC := len(eps)
		for k := 0; k+2 < len(eps); k++ {
			if eps[k] < 0 && eps[k+1] < 0 && eps[k+2] < 0 {
				kC = k
				break
			}
		}
		// fit set: eps > 0, k < kC, capped at 60
		var lk, le, kk, invE []float64
		var pLoc []float64
		for k := 0; k < kC && k < 60; k++ {
			if eps[k] <= 0 {
				continue
			}
			lk = append(lk, math.Log(float64(k+1)))
			le = append(le, math.Log(eps[k]))
			kk = append(kk, float64(k+1))
			invE = append(invE, 1/eps[k])
			if dN[k] > 0 {
				pLoc = append(pLoc, eps[k]*nMed[k]/dN[k])
			}
		}
		r := resultado{nom: cs.nom, eps: eps, nMed: nMed, dN: dN, kCruce: kC}
		if kC < len(eps) {
			r.nCruceSobreTau = nMed[kC] / tau
		}
		r.alfaPot, _, r.r2Pot = ajusteLineal(lk, le)
		r.tasaExp, _, r.r2Exp = ajusteLineal(kk, le)
		mH, bH, r2H := ajusteLineal(kk, invE)
		r.cHip, r.aHip, r.r2Hip = 1/mH, bH/mH, r2H
		r.pLocalMed = mediana(pLoc)
		res = append(res, r)
	}

	fmt.Println("§1 · la merma, primeros 10 ciclos (caso base, cero γ≈1000)")
	fmt.Println()
	fmt.Println("      k     ε_k        n medio del ciclo    Δn    exponente local p_k = ε·n/Δn")
	b := res[0]
	for k := 0; k < 10 && k < len(b.eps); k++ {
		fmt.Printf("     %2d   %+.4f        %7.1f          %4.0f          %.3f\n",
			k+1, b.eps[k], b.nMed[k], b.dN[k], b.eps[k]*b.nMed[k]/b.dN[k])
	}

	fmt.Println()
	fmt.Println("§2 · las tres leyes, ajustadas en la fase de contracción (ganador por R²)")
	fmt.Println()
	fmt.Println("   caso                        POTENCIA α (R²)      EXPONENCIAL tasa (R²)     HIPÉRBOLA c/(k+a): c, a (R²)")
	for _, r := range res {
		fmt.Printf("   %-26s   %+.3f (%.4f)      %+.4f (%.4f)         %.3f, %+.2f (%.4f)\n",
			r.nom, r.alfaPot, r.r2Pot, r.tasaExp, r.r2Exp, r.cHip, r.aHip, r.r2Hip)
	}

	fmt.Println()
	fmt.Println("§3 · el origen geométrico: exponente local p_k (mediana por caso) y el cruce")
	fmt.Println()
	fmt.Println("   caso                        p̃ local     k* (cruce)    n(k*)/τ")
	for _, r := range res {
		fmt.Printf("   %-26s    %.3f        %4d         %.3f\n", r.nom, r.pLocalMed, r.kCruce, r.nCruceSobreTau)
	}

	fmt.Println()
	fmt.Println("§4 · ¿la ley es idéntica entre casos? (puntos 5 y 6 de la orden)")
	// spread of the power exponent across sigma=0.5 cases vs zero-control differences
	var difs []float64
	nBase := len(res[0].eps)
	if len(res[2].eps) < nBase {
		nBase = len(res[2].eps)
	}
	if nBase > 30 {
		nBase = 30
	}
	for k := 0; k < nBase; k++ {
		difs = append(difs, math.Abs(res[0].eps[k]-res[2].eps[k]))
	}
	maxD, sumD := 0.0, 0.0
	for _, d := range difs {
		sumD += d
		if d > maxD {
			maxD = d
		}
	}
	fmt.Printf("   exponentes α de los cuatro casos σ=0,5: %+.3f / %+.3f / %+.3f / %+.3f\n",
		res[0].alfaPot, res[1].alfaPot, res[2].alfaPot, res[3].alfaPot)
	fmt.Printf("   y de los instrumentos σ=0,3 / σ=0,7:   %+.3f / %+.3f\n", res[4].alfaPot, res[5].alfaPot)
	fmt.Printf("   diferencia ε_k cero↔control, primeros %d ciclos: media %.4f · máx %.4f\n",
		nBase, sumD/float64(len(difs)), maxD)
	fmt.Printf("   (la merma del ciclo 1 del caso base es %.4f: la diferencia cero↔control es %.0f× menor)\n",
		res[0].eps[0], res[0].eps[0]/(sumD/float64(len(difs))))

	dibujar(res)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/la-merma.svg")
}

func dibujar(res []resultado) {
	var b strings.Builder
	W, H := 1060, 620
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	t := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	t(530, 34, 21, "#e8e2d4", "middle", "LA MERMA — cómo muere la contracción, ciclo a ciclo")
	t(530, 56, 12, "#8a93a6", "middle", "izquierda: ε_k en log-log (¿recta = ley de potencia?) · derecha: 1/ε_k contra k (¿recta = hipérbola?)")

	cols := []string{"#7ee0c0", "#6aa9ff", "#f2a6c0", "#c9b458", "#f0a05a", "#e06a6a"}

	// left: log-log
	x0, y0, cw, ch := 60.0, 110.0, 440.0, 400.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x0, y0, cw, ch)
	lxmin, lxmax := 0.0, math.Log10(60)
	lymin, lymax := -3.0, math.Log10(0.5)
	px := func(k float64) float64 { return x0 + 20 + (math.Log10(k)-lxmin)/(lxmax-lxmin)*(cw-40) }
	py := func(e float64) float64 { return y0 + ch - 20 - (math.Log10(e)-lymin)/(lymax-lymin)*(ch-40) }
	for ci, r := range res {
		for k, e := range r.eps {
			if e <= 0 || k >= 60 {
				continue
			}
			if math.Log10(e) < lymin {
				continue
			}
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.4" fill="%s" opacity="0.8"/>`, px(float64(k+1)), py(e), cols[ci])
		}
	}
	t(x0+cw/2, y0+ch+18, 11, "#8a93a6", "middle", "k (log) → ε_k (log), los seis casos superpuestos")

	// right: 1/eps vs k
	x1, y1 := 560.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x1, y1, cw, ch)
	px2 := func(k float64) float64 { return x1 + 20 + k/60*(cw-40) }
	py2 := func(v float64) float64 { return y1 + ch - 20 - v/80*(ch-40) }
	for ci, r := range res {
		for k, e := range r.eps {
			if e <= 0 || k >= 60 {
				continue
			}
			v := 1 / e
			if v > 80 {
				continue
			}
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.4" fill="%s" opacity="0.8"/>`, px2(float64(k+1)), py2(v), cols[ci])
		}
	}
	t(x1+cw/2, y1+ch+18, 11, "#8a93a6", "middle", "k → 1/ε_k · una recta acá significa ε = c/(k+a)")

	for ci, r := range res {
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.1f" r="4" fill="%s"/>`, 80.0, 84+float64((ci%3))*0-4, cols[ci])
		_ = r
	}
	// legend inline
	lx := 70.0
	for ci, r := range res {
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="3.5" fill="%s"/>`, lx, 86.0, cols[ci])
		t(lx+8, 90, 10, "#c7cdd9", "start", r.nom)
		lx += float64(9*len(r.nom)) + 26
	}

	t(530, float64(H)-12, 12, "#c9b458", "middle",
		"la ley y su origen geométrico ε = p·Δn/n, medidos a ciegas · Todavía no")

	b.WriteString(`</svg>`)
	os.WriteFile(filepath.Join("galeria", "laminas", "10-el-telar", "la-merma.svg"), []byte(b.String()), 0644)
}
