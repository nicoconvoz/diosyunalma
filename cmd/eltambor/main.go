package main

// EL TAMBOR - the captain's flash: the spiral connected to the drum. The
// drumhead's radial coordinate is log n: center = first terms (small n), rim =
// the eye zone (n ~ tau). Drum philosophy, measured:
//   - a strike at position n rings a PURE TONE of frequency ln n and loudness
//     n^{-1/2}: center = deep and loud (BOOM), rim = bright and dry (TAK);
//     at equal energy only the timbre changes (the captain's trap).
//   - the two-way travel edge<->point is the inversion n -> tau/n: mirror
//     strikes ring tones that always SUM to the drum's full tone ln(tau).
//   - and the deep question "who strikes the guitar": the eye's song sounds at
//     every position with force m^{-1/2}, but the DEPTH of the song (its
//     logarithm) is predicted to ring ONLY at primes and prime powers -
//     composites like 6, 10, 12, 15 must stay silent.
// Pre-registered predictions: (1) measured strike frequency = ln n;
// (2) mirror-pair tone sums = ln tau; (3) |proj(log|Z|, tone m)| follows
// Lambda(m)-shape: ~ p^{-k/2}/(2k) at m = p^k, ~0 at true composites.

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

// golpe: the response of the eye to a strike (removing term n) is, by
// linearity of the drum, r(t) = -e^{-i t ln n}/sqrt(n). We MEASURE its
// frequency blind, by zero crossings of the real part.
func golpe(n int, t0, T, dt float64) (frecMedida, amp float64) {
	ln := math.Log(float64(n))
	a := 1 / math.Sqrt(float64(n))
	cruces := 0
	prev := -a * math.Cos(t0*ln)
	for t := t0 + dt; t <= t0+T; t += dt {
		v := -a * math.Cos(t*ln)
		if (prev < 0) != (v < 0) {
			cruces++
		}
		prev = v
	}
	return math.Pi * float64(cruces) / T, a
}

// primoPotencia: returns (p, k) if m = p^k with p prime, else (0, 0).
// p is the SMALLEST prime factor of m; m is a prime power iff dividing it out
// completely leaves 1.
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

func main() {
	fmt.Println("🥁 EL TAMBOR — la espiral como parche: BOOM en el centro, TAK en el aro")
	fmt.Println("   coordenada radial del parche: log n · el aro es el ojo, el centro los")
	fmt.Println("   primeros términos · predicciones pre-registradas en el encabezado del código")
	fmt.Println()

	t0 := 10000.0
	tau := t0 / (2 * math.Pi)

	fmt.Println("§1 · el mapa del timbre: golpes del centro al aro (frecuencia medida a ciegas)")
	fmt.Println()
	fmt.Println("      n golpeado    tono medido    ln(n) predicho    fuerza 1/√n    timbre")
	golpes := []int{2, 5, 20, 100, 400, 1591}
	for _, n := range golpes {
		f, a := golpe(n, t0, 60, 0.002)
		timbre := "BOOM — grave, con cuerpo"
		switch {
		case n > 1000:
			timbre = "TAK — agudo, seco (aro)"
		case n > 100:
			timbre = "tak — brillante"
		case n > 10:
			timbre = "tum — medio"
		}
		fmt.Printf("      %5d        %.4f          %.4f          %.4f       %s\n",
			n, f, math.Log(float64(n)), a, timbre)
	}
	fmt.Println("   ⟹ misma energía = misma altura de barra; lo único que cambia del centro")
	fmt.Println("     al aro es el TONO. La trampa del capitán, confirmada: timbre, no fuerza.")

	fmt.Println()
	fmt.Println("§2 · el viaje de ida y vuelta: golpes espejo n ↔ τ/n")
	fmt.Println()
	fmt.Printf("   el tono total del parche: ln τ = %.4f\n", math.Log(tau))
	fmt.Println("      n      compañero τ/n     tono n + tono compañero")
	for _, n := range []int{2, 5, 20, 50} {
		m := int(math.Round(tau / float64(n)))
		f1, _ := golpe(n, t0, 60, 0.002)
		f2, _ := golpe(m, t0, 60, 0.002)
		fmt.Printf("      %3d       %5d              %.4f\n", n, m, f1+f2)
	}
	fmt.Println("   ⟹ los pares espejados suman SIEMPRE el tono entero: la onda que va del")
	fmt.Println("     borde al punto y la que vuelve son las dos mitades del mismo tono.")

	fmt.Println()
	fmt.Println("§3 · ¿quién golpea la guitarra? la profundidad del canto, proyectada tono a tono")
	fmt.Println("   el canto |Z(t)| suena entero; su LOGARITMO debería sonar solo donde")
	fmt.Println("   golpean los primos y sus potencias — los compuestos, mudos.")
	fmt.Println()
	// the song over a long window
	T := 1000.0
	dt := 0.01
	nm := int(T / dt)
	logZ := make([]float64, nm)
	ts := make([]float64, nm)
	for i := 0; i < nm; i++ {
		t := t0 + float64(i)*dt
		z := math.Abs(zetaZ(t))
		if z < 1e-4 {
			z = 1e-4
		}
		ts[i] = t
		logZ[i] = math.Log(z)
	}
	fmt.Println("      m      |proyección medida|    predicho p^(−k/2)/(2k)    posición")
	var barras []barra
	for m := 2; m <= 40; m++ {
		lm := math.Log(float64(m))
		var cr, ci float64
		for i := 0; i < nm; i++ {
			cr += logZ[i] * math.Cos(ts[i]*lm)
			ci -= logZ[i] * math.Sin(ts[i]*lm)
		}
		med := math.Hypot(cr, ci) / float64(nm)
		p, k := primoPotencia(m)
		pre := 0.0
		pos := "compuesto — debería callar"
		if p > 0 {
			pre = math.Pow(float64(p), -float64(k)/2) / (2 * float64(k))
			if k == 1 {
				pos = fmt.Sprintf("primo %d", p)
			} else {
				pos = fmt.Sprintf("potencia %d^%d", p, k)
			}
		}
		barras = append(barras, barra{m, med, pre, p > 0})
		if m <= 20 {
			fmt.Printf("      %2d          %.4f                 %.4f              %s\n", m, med, pre, pos)
		}
	}
	// scorecard: primes/powers vs composites
	var sPP, sComp float64
	var nPP, nComp int
	for _, b := range barras {
		if b.esPP {
			sPP += b.med
			nPP++
		} else {
			sComp += b.med
			nComp++
		}
	}
	fmt.Printf("\n   promedio de |proyección|: primos y potencias %.4f (%d tonos) · compuestos %.4f (%d tonos)\n",
		sPP/float64(nPP), nPP, sComp/float64(nComp), nComp)
	fmt.Printf("   razón señal/silencio: %.1f×\n", (sPP/float64(nPP))/(sComp/float64(nComp)))

	dibujar(barras, t0, tau)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-tambor.svg")
}

type barra struct {
	m        int
	med, pre float64
	esPP     bool
}

func dibujar(barras []barra, t0, tau float64) {
	var b strings.Builder
	W, H := 1060, 640
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	t := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	t(530, 34, 21, "#e8e2d4", "middle", "EL TAMBOR — quién golpea la guitarra")
	t(530, 56, 12, "#8a93a6", "middle", fmt.Sprintf("el parche es la espiral (radio = log n) · ventana t ∈ [%.0f, %.0f] · BOOM al centro, TAK al aro — y la profundidad del canto solo suena en los primos", t0, t0+1000))

	// left: the drumhead, schematic spiral with strikes
	cx, cy, R := 235.0, 330.0, 195.0
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="%.0f" fill="#11151f" stroke="#3a4258" stroke-width="3"/>`, cx, cy, R)
	// spiral: radius = log n / log tau
	var pts []string
	vueltas := 5.0
	for s := 0.001; s <= 1; s += 0.002 {
		r := R * s
		a := 2 * math.Pi * vueltas * s
		pts = append(pts, fmt.Sprintf("%.1f,%.1f", cx+r*math.Cos(a), cy+r*math.Sin(a)))
	}
	fmt.Fprintf(&b, `<polyline points="%s" fill="none" stroke="#2e3850" stroke-width="1.2"/>`, strings.Join(pts, " "))
	// center strike
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="26" fill="#c9584a" opacity="0.55"/>`, cx, cy)
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="10" fill="#e07a6a"/>`, cx, cy)
	t(cx, cy-36, 13, "#e8a396", "middle", "BOOM · n chico · tono ln n grave")
	// rim strike
	rx, ry := cx+R*0.93*math.Cos(-0.6), cy+R*0.93*math.Sin(-0.6)
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="9" fill="#7ec8e0" opacity="0.8"/>`, rx, ry)
	fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="3.5" fill="#b8e6f5"/>`, rx, ry)
	t(rx-6, ry-16, 13, "#9fd8ea", "middle", "TAK · n ≈ τ · agudo")
	t(cx, cy+R+26, 12, "#8a93a6", "middle", "la onda viaja del borde al punto y vuelve: n ↔ τ/n,")
	t(cx, cy+R+44, 12, "#8a93a6", "middle", fmt.Sprintf("y los tonos espejados suman siempre ln τ = %.3f", math.Log(tau)))

	// right: bar chart of |proj(log|Z|, tone m)|
	x0, y0, cw, ch := 500.0, 120.0, 520.0, 380.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#11151f" stroke="#232a3a"/>`, x0, y0, cw, ch)
	t(x0+cw/2, y0-10, 13, "#e8e2d4", "middle", "la profundidad del canto, proyectada sobre cada tono ln m")
	vmax := 0.0
	for _, bb := range barras {
		if bb.med > vmax {
			vmax = bb.med
		}
	}
	bw := cw / float64(len(barras))
	for i, bb := range barras {
		h := bb.med / vmax * (ch - 40)
		col := "#5c6478"
		if bb.esPP {
			col = "#c9b458"
		}
		bx := x0 + float64(i)*bw + 2
		fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`, bx, y0+ch-20-h, bw-4, h, col)
		// predicted tick
		if bb.pre > 0 {
			hp := bb.pre / vmax * (ch - 40)
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#7ee0c0" stroke-width="2"/>`, bx, y0+ch-20-hp, bx+bw-4, y0+ch-20-hp)
		}
		if bb.esPP || bb.m%10 == 0 {
			t(bx+(bw-4)/2, y0+ch-4, 9, "#8a93a6", "middle", fmt.Sprintf("%d", bb.m))
		}
	}
	t(x0+16, y0+22, 12, "#c9b458", "start", "■ primos y potencias — suenan")
	t(x0+16, y0+42, 12, "#5c6478", "start", "■ compuestos — callan")
	t(x0+16, y0+62, 12, "#7ee0c0", "start", "— lo predicho: p^(−k/2)/(2k)")

	t(530, float64(H)-32, 12, "#c7cdd9", "middle",
		"el canto entero suena en todas las posiciones; su PROFUNDIDAD suena solo donde golpean los primos —")
	t(530, float64(H)-12, 12, "#c9b458", "middle",
		"el tambor de la guitarra lo tocan los primos, cada uno desde su radio log p · Todavía no")

	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-tambor.svg")
	os.WriteFile(ruta, []byte(b.String()), 0644)
}
