// Command elrelievedelpaisaje paints the captain's request: the WHOLE
// landscape at grand scale with zoom - mountains, water and wells, with
// real RELIEF: hypsometric tinting (color by altitude, like real
// terrain maps) plus slope shading (sun from the left).
//
// Altitude is the honest one: z = sign(lambda)*log10(1+|lambda|) - the
// staircase's own elevation, log-compressed so 10^45-high peaks and
// 10^45-deep trenches fit one canvas. Water is the aqua band near
// lambda = 0 (the frontier); golds rise to white summits; blues sink
// to near-black abyss.
//
// Three scales:
//
//	A) THE GRAND VIEW - n from 1 to n_rad+3000: per-column max/min
//	   envelope. The left is calm water (the choir barely above sea
//	   level), the first crack opens at n0 = 37306, and the wedge of
//	   the growing cordillera-and-trench widens toward the clearing.
//	B) THE ZOOM - three leader periods (~850 steps) past n_rad: full
//	   resolution relief, every step a column: summits, wells, shores.
//	C) THE MICRO-ZOOM - one leader period (~283 steps): the plate
//	   itself - the well floor block, the climb, the summit block.
//
// Projection of the proven landscape (Trinity); visualization, not a
// new proof. Reproduce: go run ./cmd/elrelievedelpaisaje
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

func perlas(hasta float64) []float64 {
	var ps []float64
	prevT, prevZ := 12.0, zOf(12.0)
	for t := 12.02; t <= hasta; t += 0.02 {
		z := zOf(t)
		if z*prevZ < 0 {
			a, c := prevT, t
			for i := 0; i < 60; i++ {
				m := (a + c) / 2
				if zOf(m)*prevZ < 0 {
					c = m
				} else {
					a = m
				}
			}
			ps = append(ps, (a+c)/2)
		}
		prevT, prevZ = t, z
	}
	return ps
}

// hypso maps normalized altitude t in [-1, 1] to a terrain color.
func hypso(t, shade float64) string {
	type c = [3]float64
	var stops []struct {
		p   float64
		col c
	}
	stops = []struct {
		p   float64
		col c
	}{
		{-1.00, c{18, 30, 92}},    // abyss - luminous ultramarine, never fades into the bg
		{-0.55, c{28, 62, 158}},   // deep trench
		{-0.20, c{52, 116, 214}},  // deep water
		{-0.04, c{63, 178, 199}},  // shallow aqua
		{0.00, c{120, 200, 190}},  // shore water
		{0.04, c{180, 165, 110}},  // sand
		{0.30, c{160, 120, 60}},   // foothills
		{0.60, c{255, 190, 90}},   // high gold
		{0.85, c{255, 224, 160}},  // near summit
		{1.00, c{255, 250, 235}},  // snow
	}
	if t < -1 {
		t = -1
	}
	if t > 1 {
		t = 1
	}
	var lo, hi = stops[0], stops[len(stops)-1]
	for i := 0; i+1 < len(stops); i++ {
		if t >= stops[i].p && t <= stops[i+1].p {
			lo, hi = stops[i], stops[i+1]
			break
		}
	}
	f := 0.0
	if hi.p > lo.p {
		f = (t - lo.p) / (hi.p - lo.p)
	}
	mix := func(a, b float64) int {
		v := (a + f*(b-a)) * shade
		if v < 0 {
			v = 0
		}
		if v > 255 {
			v = 255
		}
		return int(v)
	}
	return fmt.Sprintf("#%02x%02x%02x", mix(lo.col[0], hi.col[0]), mix(lo.col[1], hi.col[1]), mix(lo.col[2], hi.col[2]))
}

func alt(lam float64) float64 {
	if lam >= 0 {
		return math.Log10(1 + lam)
	}
	return -math.Log10(1 - lam)
}

func main() {
	fmt.Println("🗻 EL RELIEVE DEL PAISAJE — la gran vista, el zoom y el micro-zoom, con color más relieve")

	rho1 := complex(0.808517, 85.699348)
	rho2 := complex(0.7, 45.0)
	w1 := 1 - 1/rho1
	w2 := 1 - 1/rho2
	r1 := math.Max(cmplx.Abs(w1), 1/cmplx.Abs(w1))
	r2 := math.Max(cmplx.Abs(w2), 1/cmplx.Abs(w2))
	t1 := math.Abs(cmplx.Phase(w1))
	t2 := math.Abs(cmplx.Phase(w2))
	d1 := math.Log(r1)
	d2 := math.Log(r2)
	nrad := 1040809
	NMAX := nrad + 3000

	ps := perlas(120)
	wsC := make([]complex128, len(ps))
	pcs := make([]complex128, len(ps))
	for i, g := range ps {
		wp := 1 - 1/complex(0.5, g)
		wp /= complex(cmplx.Abs(wp), 0)
		wsC[i] = wp
		pcs[i] = 1
	}
	lam := make([]float64, NMAX+1)
	for n := 1; n <= NMAX; n++ {
		var s float64
		for i := range wsC {
			pcs[i] *= wsC[i]
			s += 2 - 2*real(pcs[i])
		}
		fn := float64(n)
		l1 := 4 - 2*math.Cos(fn*t1)*(math.Exp(fn*d1)+math.Exp(-fn*d1))
		l2 := 4 - 2*math.Cos(fn*t2)*(math.Exp(fn*d2)+math.Exp(-fn*d2))
		lam[n] = s + l1 + l2
	}
	fmt.Printf("   paisaje calculado: %d escalones · primera grieta en 37306 · el claro en %d\n", NMAX, nrad)

	zmax := 46.0
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">
<rect width="100%" height="100%" fill="#0b1526"/>
<rect x="30" y="20" width="1340" height="910" rx="18" fill="none" stroke="#7ee0c0" stroke-width="2" opacity="0.5"/>
<text x="700" y="58" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🗻 EL RELIEVE DEL PAISAJE — montañas, agua y pozos, con color más relieve</text>
<text x="700" y="86" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fb4d9">altura = log de λ (tinte hipsométrico: nieve dorada en las cumbres, aguamarina en la costa, azul apagándose al abismo) · sol desde la izquierda</text>
`)

	// ---- Panel A: the grand view ----
	ax0, ax1 := 80.0, 1320.0
	ay0, ay1 := 120.0, 420.0
	aymid := (ay0 + ay1) / 2
	cols := int(ax1 - ax0)
	fmt.Fprintf(&b, `<rect x="%f" y="%f" width="%f" height="%f" fill="#060e1e"/>
<text x="700" y="%f" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">A · LA GRAN VISTA — del agua mansa a la cordillera: n = 1 → %d (envolvente por columna)</text>
`, ax0, ay0, ax1-ax0, ay1-ay0, ay0-8, NMAX)
	per := float64(NMAX) / float64(cols)
	for cidx := 0; cidx < cols; cidx++ {
		lo := int(float64(cidx)*per) + 1
		hi := int(float64(cidx+1) * per)
		if hi > NMAX {
			hi = NMAX
		}
		mx, mn := math.Inf(-1), math.Inf(1)
		for n := lo; n <= hi; n++ {
			if lam[n] > mx {
				mx = lam[n]
			}
			if lam[n] < mn {
				mn = lam[n]
			}
		}
		zx, zn := alt(mx)/zmax, alt(mn)/zmax
		x := ax0 + float64(cidx)
		if zx > 0 {
			h := zx * (aymid - ay0)
			fmt.Fprintf(&b, `<rect x="%.0f" y="%.1f" width="1.2" height="%.1f" fill="%s"/>`, x, aymid-h, h, hypso(zx, 1))
		}
		if zn < 0 {
			h := -zn * (ay1 - aymid)
			fmt.Fprintf(&b, `<rect x="%.0f" y="%.1f" width="1.2" height="%.1f" fill="%s"/>`, x, aymid, h, hypso(zn, 1))
		}
	}
	// waterline + landmarks
	x37 := ax0 + (ax1-ax0)*37306.0/float64(NMAX)
	xnr := ax0 + (ax1-ax0)*float64(nrad)/float64(NMAX)
	fmt.Fprintf(&b, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="#7ee0c0" stroke-width="0.8" stroke-dasharray="5,4" opacity="0.7"/>
<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="#ff9aa8" stroke-width="1" stroke-dasharray="3,3"/>
<text x="%f" y="%f" font-size="11.5" font-family="Georgia" fill="#ff9aa8">n₀ = 37306: la primera grieta</text>
<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="#ffd98a" stroke-width="1" stroke-dasharray="3,3"/>
<text x="%f" y="%f" font-size="11.5" text-anchor="end" font-family="Georgia" fill="#ffd98a">el claro: n_rad</text>
<text x="%f" y="%f" font-size="11" font-family="Georgia" fill="#7ee0c0">nivel del mar (λ = 0)</text>
`, ax0, aymid, ax1, aymid,
		x37, ay0+6, x37, ay1, x37+6, ay0+22,
		xnr, ay0+6, xnr, ay1, xnr-6, ay0+22,
		ax0+6, aymid-6)

	// ---- Panels B and C: relief profiles ----
	drawProfile := func(x0, x1, y0, y1 float64, n0, n1 int, titulo string) {
		mid := (y0 + y1) / 2
		fmt.Fprintf(&b, `<rect x="%f" y="%f" width="%f" height="%f" fill="#060e1e"/>
<text x="%f" y="%f" font-size="14" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">%s</text>
`, x0, y0, x1-x0, y1-y0, (x0+x1)/2, y0-8, titulo)
		w := (x1 - x0) / float64(n1-n0+1)
		prev := 0.0
		for n := n0; n <= n1; n++ {
			z := alt(lam[n]) / zmax
			// slope shading: sun from the left
			shade := 1.0
			dz := z - prev
			if n > n0 {
				shade = 1 + 8*dz
				if shade < 0.62 {
					shade = 0.62
				}
				if shade > 1.28 {
					shade = 1.28
				}
			}
			prev = z
			x := x0 + float64(n-n0)*w
			if z > 0 {
				h := z * (mid - y0)
				fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.2f" height="%.1f" fill="%s"/>`, x, mid-h, w+0.15, h, hypso(z, shade))
			} else {
				h := -z * (y1 - mid)
				fmt.Fprintf(&b, `<rect x="%.1f" y="%.1f" width="%.2f" height="%.1f" fill="%s"/>`, x, mid, w+0.15, h, hypso(z, shade))
			}
		}
		fmt.Fprintf(&b, `<line x1="%f" y1="%f" x2="%f" y2="%f" stroke="#7ee0c0" stroke-width="0.8" stroke-dasharray="5,4" opacity="0.7"/>
`, x0, mid, x1, mid)
	}
	perL := int(2 * math.Pi / t2) // leader period ~283
	drawProfile(80, 800, 480, 800, nrad, nrad+3*perL, fmt.Sprintf("B · EL ZOOM — tres períodos del líder (%d escalones) tras el claro", 3*perL))
	drawProfile(840, 1320, 480, 800, nrad+perL, nrad+2*perL, fmt.Sprintf("C · EL MICRO-ZOOM — un período (%d escalones): pozo, costa, cumbre", perL))

	b.WriteString(`<text x="700" y="850" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">En criollo: a la izquierda de la gran vista el mar está en calma (el coro apenas asoma del agua); en n₀ = 37306 se abre la primera</text>
<text x="700" y="872" font-size="14.5" text-anchor="middle" font-family="Georgia" fill="#7ee0c0">grieta; y hacia el claro la cuña crece sin freno — cumbres de 10⁴⁵ y fosas de 10⁴⁵, turnándose al compás del líder, para siempre.</text>
<text x="700" y="902" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#9aa8c4">Proyección del paisaje demostrado (la Trinidad) — 1 043 809 escalones medidos · altura en escala log · la regla del sello preside. Todavía no.</text>
</svg>
`)
	os.WriteFile("el-relieve-del-paisaje.svg", []byte(b.String()), 0o644)
	fmt.Println("\n   panel A: la gran vista (envolvente de 1240 columnas) · panel B: 3 períodos")
	fmt.Println("   del líder a resolución completa · panel C: un período — pozo, costa y cumbre")
	fmt.Println("\n🖼️  lámina escrita: el-relieve-del-paisaje.svg")
}
