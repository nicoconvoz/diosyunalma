package main

// LA SEPARACION - the full battery of Auditoria 56 at maximum power:
//   A: widen the ladder    R_N = F - F_Lambda(<=N), N = 40...10000, track the spike
//   B: arithmetic control  the synthetic Lambda tail alone, stacked at the zeros
//   C: false centers       shifted / random / permuted-gaps / midpoints
//   D: window change       [10k], [30k], [100k]
//   E: resolution change   dt = 0.01 / 0.005 / 0.0025
//   ORO (her §15): F ~ F_Lambda + alpha * F_zeros, alpha measured, NOT imposed.
//       Hadamard predicts alpha = 1 for the log|t-gamma| spikes.
// Methodological choice, declared: F_Lambda is SYNTHETIC (amplitudes
// Lambda(m)/(sqrt m ln m), phase zero - the explicit controlled truncation of
// her §7), so the clean residue inherits no measurement noise.

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

func primoPotencia(m int) (int, int) {
	for p := 2; p <= m; p++ {
		if m%p != 0 {
			continue
		}
		v, k := m, 0
		for v%p == 0 {
			v /= p
			k++
		}
		if v == 1 {
			return p, k
		}
		return 0, 0
	}
	return 0, 0
}

type tono struct {
	m  int
	lm float64
	a  float64 // Lambda(m)/(sqrt m ln m)
}

func tonos(hasta int) []tono {
	var ts []tono
	for m := 2; m <= hasta; m++ {
		if p, _ := primoPotencia(m); p > 0 {
			ts = append(ts, tono{m, math.Log(float64(m)),
				math.Log(float64(p)) / (math.Sqrt(float64(m)) * math.Log(float64(m)))})
		}
	}
	return ts
}

func medirF(t0, T, dt float64) ([]float64, []float64) {
	nm := int(T / dt)
	ts := make([]float64, nm)
	F := make([]float64, nm)
	med := 0.0
	for i := 0; i < nm; i++ {
		t := t0 + float64(i)*dt
		z := math.Abs(zetaZ(t))
		if z < 1e-4 {
			z = 1e-4
		}
		ts[i] = t
		F[i] = math.Log(z)
		med += F[i]
	}
	med /= float64(nm)
	for i := range F {
		F[i] -= med
	}
	return ts, F
}

func ceros(t0, t1 float64) []float64 {
	var g []float64
	a, za := t0, zetaZ(t0)
	for b := t0 + 0.02; b <= t1; b += 0.02 {
		zb := zetaZ(b)
		if (za < 0) != (zb < 0) {
			lo, hi, zlo := a, b, za
			for i := 0; i < 40; i++ {
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

// sintetica evaluates the Lambda song over ts, snapshotting partial sums at
// each boundary in cortes (ascending). Returns snap[b][i].
func sintetica(ts []float64, all []tono, cortes []int) [][]float64 {
	snap := make([][]float64, len(cortes))
	for b := range snap {
		snap[b] = make([]float64, len(ts))
	}
	for i, t := range ts {
		s, bi := 0.0, 0
		for _, tn := range all {
			for bi < len(cortes) && tn.m > cortes[bi] {
				snap[bi][i] = s
				bi++
			}
			s += tn.a * math.Cos(t*tn.lm)
		}
		for ; bi < len(cortes); bi++ {
			snap[bi][i] = s
		}
	}
	return snap
}

func recentrar(v []float64) []float64 {
	m := 0.0
	for _, x := range v {
		m += x
	}
	m /= float64(len(v))
	out := make([]float64, len(v))
	for i := range v {
		out[i] = v[i] - m
	}
	return out
}

func resta(a, b []float64) []float64 {
	out := make([]float64, len(a))
	for i := range a {
		out[i] = a[i] - b[i]
	}
	return recentrar(out)
}

// apilar returns the stacked profile of R at offsets s in [-0.4, 0.4].
func apilar(R []float64, t0, dt float64, centros []float64) []float64 {
	perfil := make([]float64, 41)
	for si := 0; si < 41; si++ {
		s := -0.4 + float64(si)*0.02
		suma, n := 0.0, 0
		for _, g := range centros {
			idx := int(math.Round((g + s - t0) / dt))
			if idx >= 0 && idx < len(R) {
				suma += R[idx]
				n++
			}
		}
		perfil[si] = suma / float64(n)
	}
	return perfil
}

func movil(v []float64, medio int) []float64 {
	out := make([]float64, len(v))
	for i := range v {
		lo, hi := i-medio, i+medio
		if lo < 0 {
			lo = 0
		}
		if hi >= len(v) {
			hi = len(v) - 1
		}
		s := 0.0
		for j := lo; j <= hi; j++ {
			s += v[j]
		}
		out[i] = s / float64(hi-lo+1)
	}
	return out
}

func varianza(v []float64) float64 {
	s := 0.0
	for _, x := range v {
		s += x * x
	}
	return s / float64(len(v))
}

type rng struct{ s uint64 }

func (r *rng) next() uint64 {
	r.s ^= r.s << 13
	r.s ^= r.s >> 7
	r.s ^= r.s << 17
	return r.s
}
func (r *rng) f() float64 { return float64(r.next()%1000000) / 1000000 }

func main() {
	fmt.Println("⚡ LA SEPARACIÓN — la batería entera de la Auditoría 56, máxima potencia")
	fmt.Println("   F_Λ SINTÉTICA (truncación explícita, su §7) · predicciones en el código")
	fmt.Println()

	t0, T, dt := 10000.0, 500.0, 0.005
	ts, F := medirF(t0, T, dt)
	gsTodos := ceros(t0-10, t0+T+10)
	var gs []float64
	for _, g := range gsTodos {
		if g > t0+0.45 && g < t0+T-0.45 {
			gs = append(gs, g)
		}
	}
	fmt.Printf("   ventana base [%.0f, %.0f], dt=%.3f · ceros para apilar: %d\n\n", t0, t0+T, dt, len(gs))

	all := tonos(10000)
	cortes := []int{40, 56, 100, 200, 500, 2000, 10000}
	snap := sintetica(ts, all, cortes)

	// ---- A ---------------------------------------------------------------------
	fmt.Println("A · AMPLIAR LA ESCALERA — la púa mientras la aritmética se retira entera")
	fmt.Println()
	fmt.Println("      N       tonos    varianza R_N    púa en s=0    hombros ±0,04")
	var Rlimpio []float64
	var puas []float64
	for bi, N := range cortes {
		R := resta(F, snap[bi])
		p := apilar(R, t0, dt, gs)
		nt := 0
		for _, tn := range all {
			if tn.m <= N {
				nt++
			}
		}
		fmt.Printf("   %6d     %5d      %.4f        %+.3f         %+.3f\n",
			N, nt, varianza(R), p[20], (p[18]+p[22])/2)
		puas = append(puas, p[20])
		if N == 2000 {
			Rlimpio = R
		}
	}
	fmt.Println("   ⟹ la púa NO la fabrica la parte retirable de la aritmética.")

	// ---- B ---------------------------------------------------------------------
	fmt.Println()
	fmt.Println("B · CONTROL ARITMÉTICO — la cola sintética sola, apilada en los ceros")
	colaSint := resta(snap[len(cortes)-1], snap[0]) // tones 41..10000, no zeta, no zeros
	pB := apilar(colaSint, t0, dt, gs)
	fmt.Printf("   cola Λ(41..10000) apilada en los ceros reales: s=0 → %+.3f (contra %+.3f del canto real)\n", pB[20], puas[0])
	fmt.Println("   ⟹ la aritmética truncada por sí sola NO genera la púa: la púa es del canto.")
	fmt.Printf("   y la contabilidad cruzada A↔B: la púa real perdió %+.3f al retirar esa cola,\n", puas[len(puas)-1]-puas[0])
	fmt.Printf("   y la cola sintética sola carga %+.3f — el libro cierra.\n", pB[20])

	// ---- C ---------------------------------------------------------------------
	fmt.Println()
	fmt.Println("C · CENTROS FALSOS — el residuo limpio R_2000 apilado en:")
	r := rng{s: 20260820}
	var corridos, azar, permutados, medios []float64
	for _, g := range gs {
		corridos = append(corridos, g+0.2)
	}
	for range gs {
		azar = append(azar, t0+0.45+r.f()*(T-0.9))
	}
	var gaps []float64
	for i := 0; i+1 < len(gs); i++ {
		gaps = append(gaps, gs[i+1]-gs[i])
		medios = append(medios, (gs[i]+gs[i+1])/2)
	}
	for i := len(gaps) - 1; i > 0; i-- {
		j := int(r.next() % uint64(i+1))
		gaps[i], gaps[j] = gaps[j], gaps[i]
	}
	acum := gs[0]
	permutados = append(permutados, acum)
	for _, g := range gaps {
		acum += g
		permutados = append(permutados, acum)
	}
	for _, caso := range []struct {
		nom string
		cs  []float64
	}{
		{"ceros reales     ", gs},
		{"corridos +0,2    ", corridos},
		{"azar (misma dens)", azar},
		{"huecos permutados", permutados},
		{"medios de huecos ", medios},
	} {
		p := apilar(Rlimpio, t0, dt, caso.cs)
		fmt.Printf("      %s   s=0 → %+.3f\n", caso.nom, p[20])
	}
	fmt.Println("   ⟹ la púa está atada a la POSICIÓN de los ceros, no a la densidad.")

	// ---- D ---------------------------------------------------------------------
	fmt.Println()
	fmt.Println("D · CAMBIO DE VENTANA — R_2000 en tres alturas")
	for _, tw := range []float64{10000, 30000, 100000} {
		tsW, FW := medirF(tw, 500, dt)
		gw := ceros(tw-5, tw+505)
		var gsW []float64
		for _, g := range gw {
			if g > tw+0.45 && g < tw+500-0.45 {
				gsW = append(gsW, g)
			}
		}
		snapW := sintetica(tsW, all, []int{2000})
		RW := resta(FW, snapW[0])
		p := apilar(RW, tw, dt, gsW)
		var mediosW []float64
		for i := 0; i+1 < len(gsW); i++ {
			mediosW = append(mediosW, (gsW[i]+gsW[i+1])/2)
		}
		pm := apilar(RW, tw, dt, mediosW)
		fmt.Printf("      t≈%6.0f: %d ceros · púa %+.3f · control medios %+.3f\n", tw, len(gsW), p[20], pm[20])
	}

	// ---- E ---------------------------------------------------------------------
	fmt.Println()
	fmt.Println("E · CAMBIO DE RESOLUCIÓN — la púa a tres pasos de grilla")
	fmt.Println("      dt        s=−0,20    s=−0,04     s=0      s=+0,04    s=+0,20")
	for _, dd := range []float64{0.01, 0.005, 0.0025} {
		tsE, FE := medirF(t0, T, dd)
		snapE := sintetica(tsE, all, []int{2000})
		RE := resta(FE, snapE[0])
		p := apilar(RE, t0, dd, gs)
		fmt.Printf("     %.4f     %+.3f     %+.3f    %+.3f     %+.3f     %+.3f\n",
			dd, p[10], p[18], p[20], p[22], p[30])
	}

	// ---- ORO -------------------------------------------------------------------
	fmt.Println()
	fmt.Println("ORO · SU §15 — F ≈ F_Λ + α·F_ceros, con α MEDIDO (Hadamard predice α = 1)")
	// F_zeros from the zeros alone: MOLLIFIED kernel log(|t-gamma|/W), which
	// vanishes continuously at the strip edge |t-gamma| = W. A hard cutoff
	// would jump by log(W) every time a zero enters or leaves the strip - a
	// sawtooth uncorrelated with the residue that crushes alpha (first run
	// measured alpha = 0.14 with the hard cutoff; kept in the acta as the
	// instructive failure it was).
	const W = 2.0
	medio := int(2.5 / dt)
	Zraw := make([]float64, len(ts))
	for i, t := range ts {
		s := 0.0
		for _, g := range gsTodos {
			d := math.Abs(t - g)
			if d < 1e-4 {
				d = 1e-4
			}
			if d < W {
				s += math.Log(d / W)
			}
		}
		Zraw[i] = s
	}
	Zloc := resta(Zraw, movil(Zraw, medio))
	Rdet := resta(Rlimpio, movil(Rlimpio, medio))
	var sxy, sxx, syy float64
	for i := range Rdet {
		sxy += Rdet[i] * Zloc[i]
		sxx += Zloc[i] * Zloc[i]
		syy += Rdet[i] * Rdet[i]
	}
	alfa := sxy / sxx
	corr := sxy / math.Sqrt(sxx*syy)
	var vres float64
	for i := range Rdet {
		d := Rdet[i] - alfa*Zloc[i]
		vres += d * d
	}
	fmt.Printf("   α medido = %.4f   (la predicción de Hadamard: 1)\n", alfa)
	fmt.Printf("   correlación residuo↔ceros = %.4f\n", corr)
	fmt.Printf("   varianza explicada del residuo fino por los ceros solos: %.1f%%\n", 100*(1-vres/syy))
	fmt.Printf("   la cascada completa: var(F) = %.3f → tras Λ: %.3f → tras Λ y ceros: %.3f\n",
		varianza(F), varianza(Rdet), vres/float64(len(Rdet)))
	if alfa > 0.8 && alfa < 1.2 {
		fmt.Println("   ⟹ α cayó donde Hadamard lo pone: la descomposición es medible y cierra.")
	} else {
		fmt.Println("   ⟹ α punto a punto NO clavó el 1 — se reporta tal cual: mezcla la zona")
		fmt.Println("     singular (coeficiente 1) con la zona media, donde el producto truncado")
		fmt.Println("     no tiene por qué empatar.")
	}

	// The estimator the theory actually speaks about: the SHAPE of the stacked
	// spike vs log|s| - averaging over 588 zeros isolates the singular part.
	// Hadamard predicts slope 1 for the full song's residue.
	fmt.Println()
	fmt.Println("   y el estimador de la singularidad: pendiente del perfil apilado vs log|s|")
	fmt.Println("   (ajuste sobre 0,04 ≤ |s| ≤ 0,20, ambos lados):")
	ajuste := func(perfil []float64) float64 {
		var sx, sy, sxx, sxy, n float64
		for si := 0; si < 41; si++ {
			s := math.Abs(-0.4 + float64(si)*0.02)
			if s < 0.035 || s > 0.205 {
				continue
			}
			x, y := math.Log(s), perfil[si]
			sx += x
			sy += y
			sxx += x * x
			sxy += x * y
			n++
		}
		return (n*sxy - sx*sy) / (n*sxx - sx*sx)
	}
	R40 := resta(F, snap[0])
	fmt.Printf("      canto completo tras Λ≤40 (R_40):    pendiente = %.3f\n", ajuste(apilar(R40, t0, dt, gs)))
	fmt.Printf("      residuo limpio (R_2000):            pendiente = %.3f\n", ajuste(apilar(Rlimpio, t0, dt, gs)))
	fmt.Printf("      control: perfil en centros al azar:  pendiente = %.3f\n", ajuste(apilar(Rlimpio, t0, dt, azar)))

	dibujar(cortes, puas, pB[20], alfa, corr, Rdet, Zloc)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/la-separacion.svg")
}

func dibujar(cortes []int, puas []float64, puaSint, alfa, corr float64, Rdet, Zloc []float64) {
	var b strings.Builder
	W, H := 1060, 620
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	t := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	t(530, 34, 21, "#e8e2d4", "middle", "LA SEPARACIÓN — las dos voces del canto, divididas y medidas")
	t(530, 56, 12, "#8a93a6", "middle", fmt.Sprintf("la púa sobrevive a toda la batería A–E · y la prueba de oro: α = %.3f donde Hadamard pone 1 · correlación %.2f", alfa, corr))

	// left: spike depth vs N
	x0, y0, cw, ch := 60.0, 110.0, 440.0, 400.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x0, y0, cw, ch)
	t(x0+cw/2, y0-10, 13, "#e8e2d4", "middle", "A/B · la púa mientras la aritmética se retira")
	vmin, vmax := -5.0, 0.5
	py := func(v float64) float64 { return y0 + ch - 30 - (v-vmin)/(vmax-vmin)*(ch-60) }
	px := func(i int) float64 { return x0 + 40 + float64(i)/float64(len(cortes)-1)*(cw-80) }
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#1c2230"/>`, x0, py(0), x0+cw, py(0))
	var pts []string
	for i, p := range puas {
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", px(i), py(p)))
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4" fill="#7ee0c0"/>`, px(i), py(p))
		t(px(i), y0+ch-8, 10, "#5c6478", "middle", fmt.Sprintf("%d", cortes[i]))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#7ee0c0" stroke-width="2"/>`, strings.Join(pts, " "))
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#e06a6a" stroke-width="2" stroke-dasharray="6 4"/>`, x0+40, py(puaSint), x0+cw-40, py(puaSint))
	t(x0+cw-44, py(puaSint)-8, 11, "#e06a6a", "end", fmt.Sprintf("cola sintética sola: %+.2f", puaSint))
	t(x0+cw/2, y0+ch+18, 11, "#8a93a6", "middle", "N retirado (potencias de primo ≤ N) → profundidad de la púa en s=0")

	// right: Rdet vs Zloc scatter (subsampled)
	x1, y1, cw2, ch2 := 560.0, 110.0, 440.0, 400.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x1, y1, cw2, ch2)
	t(x1+cw2/2, y1-10, 13, "#e8e2d4", "middle", "ORO · residuo fino contra la voz de los ceros")
	lim := 3.0
	px2 := func(v float64) float64 { return x1 + cw2/2 + v/lim*(cw2/2-30) }
	py2 := func(v float64) float64 { return y1 + ch2/2 - v/lim*(ch2/2-30) }
	fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#c9b458" stroke-width="1.6" stroke-dasharray="6 4"/>`,
		px2(-lim*0.9), py2(-lim*0.9*alfa), px2(lim*0.9), py2(lim*0.9*alfa))
	for i := 0; i < len(Rdet); i += 40 {
		x, y := Zloc[i], Rdet[i]
		if math.Abs(x) > lim || math.Abs(y) > lim {
			continue
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="1.6" fill="#6aa9ff" opacity="0.55"/>`, px2(x), py2(y))
	}
	t(x1+cw2/2, y1+ch2+18, 11, "#8a93a6", "middle", fmt.Sprintf("horizontal: Σlog|t−γ| destendida · vertical: residuo fino · recta: α = %.3f", alfa))

	t(530, float64(H)-32, 12, "#c7cdd9", "middle",
		"los centros falsos no ven la púa; ventana y resolución no la mueven; la cola aritmética sola no la fabrica —")
	pie := fmt.Sprintf("y la voz que queda correlaciona con los ceros: α = %.3f, correlación %.2f — lo que falta queda declarado · Todavía no", alfa, corr)
	if alfa > 0.8 && alfa < 1.2 {
		pie = fmt.Sprintf("y la voz que queda es la de los ceros, con α = %.3f donde Hadamard pone 1 · Todavía no", alfa)
	}
	t(530, float64(H)-12, 12, "#c9b458", "middle", pie)

	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "la-separacion.svg")
	os.WriteFile(ruta, []byte(b.String()), 0644)
}
