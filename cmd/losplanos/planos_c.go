// planos_c.go - desk C of the museum tour: the sonar's funnel, the three
// engines that write in the same ledger, and the frozen gear the laboratory
// tripped over. Every number drawn on these three plates is measured here.
package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// ---- desk C helpers (p08/p09/p10 prefixes per house rule) ----

// p08rnd is a fixed linear congruential generator: the funnel must draw the
// SAME thousand balls on every run, so the plate is reproducible.
type p08rnd struct{ s uint64 }

func (g *p08rnd) u() float64 {
	g.s = g.s*6364136223846793005 + 1442695040888963407
	return float64(g.s>>11) / float64(uint64(1)<<53)
}

// sigma draws |Z|/sqrt(L) for a random-phase block: the sea's own Rayleigh law.
func (g *p08rnd) sigma() float64 { return math.Sqrt(-math.Log(1 - g.u())) }

// p09rueda draws a toothed wheel of radius r.
func p09rueda(l *Lienzo, cx, cy, r float64, col string, dientes int) {
	l.Circ(cx, cy, r, cPanel, col, 2)
	l.Circ(cx, cy, r*0.16, col, "", 0)
	for i := 0; i < dientes; i++ {
		a := 2 * math.Pi * float64(i) / float64(dientes)
		l.Line(cx+r*math.Cos(a), cy+r*math.Sin(a), cx+(r+7)*math.Cos(a), cy+(r+7)*math.Sin(a), col, 2.6, "")
	}
}

// p09aguja plants a needle on a 240-degree gauge.
func p09aguja(l *Lienzo, cx, cy, r, lo, hi, v float64, col string, ancho float64) {
	f := math.Max(0, math.Min(1, (v-lo)/(hi-lo)))
	a := (150 + 240*f) * math.Pi / 180
	l.Line(cx, cy, cx+(r-10)*math.Cos(a), cy+(r-10)*math.Sin(a), col, ancho, "")
}

// p09censo reads the arbiter's census live from the hunt's own book; if the
// book is out of reach it falls back to the census transcribed from it.
func p09censo() (es []float64, dir string, fuente string) {
	dir = "t=1e+40 k=2.131827179357483e+19 L=4000"
	if raw, err := os.ReadFile("luz/cazadero.log"); err == nil {
		for _, ln := range strings.Split(string(raw), "\n") {
			if !strings.HasPrefix(ln, "ARBITRO:") && !strings.HasPrefix(ln, "FANTASMA5") {
				continue
			}
			for _, f := range strings.Fields(ln) {
				if strings.HasPrefix(f, "e256=") {
					if v, e2 := strconv.ParseFloat(f[5:], 64); e2 == nil {
						es = append(es, v)
					}
				}
			}
			if strings.HasPrefix(ln, "FANTASMA5") {
				fs := strings.Fields(ln)
				if len(fs) > 5 {
					dir = fs[2] + " " + fs[3] + " " + fs[4]
				}
			}
		}
	}
	if len(es) > 0 {
		return es, dir, "luz/cazadero.log"
	}
	return []float64{4.3e-3, 1.4e-3, 1.6e-4, 1.1e-3, 7.0e-4, 1.7e-3, 3.5e-3, 2.1e-3, 1.4e-3,
		9.1e-4, 5.2e-4, 6.8e-1, 3.1e-4, 1.1e-3, 1.0e-3, 3.2e-3, 2.1e-5, 1.3e-3, 4.4e-3, 1.6e-3,
		8.3e-4, 1.6e-3, 4.4e-4, 3.5e-4, 4.3e-3, 4.3e-3, 1.4e-3, 1.6e-4, 1.1e-3, 7.0e-4, 1.7e-3, 3.5e-3}, dir, "censo transcripto del libro"
}

// p09vuelta is one turn of the circle (cmd/circulo chirpFlip): Poisson plus
// stationary phase rewrite the chirp as its shorter dual.
func p09vuelta(b, c float64, n int64) (re, im float64, dual int64) {
	m1, m2, largo := dualRango(b, c, n)
	inv, amp := 1/(4*b), 1/math.Sqrt(2*b)
	for m := m1; m <= m2; m++ {
		x := float64(m) - c
		s, co := math.Sincos(2*math.Pi*math.Mod(-x*x*inv, 1) + math.Pi/4)
		re += amp * co
		im += amp * s
	}
	q := math.Mod(b*float64(n-1)*float64(n-1)+c*float64(n-1), 1)
	re += 0.5 * (1 + math.Cos(2*math.Pi*q))
	im += 0.5 * math.Sin(2*math.Pi*q)
	return re, im, largo
}

// p10hist buckets a measured sample into nb bins over [0,hi].
func p10hist(v []float64, nb int, hi float64) []int {
	h := make([]int, nb)
	for _, x := range v {
		i := int(x / hi * float64(nb))
		if i >= nb {
			i = nb - 1
		}
		if i < 0 {
			i = 0
		}
		h[i]++
	}
	return h
}

// ---------------------------------------------------------------------------
// PLATE 08 - EL EMBUDO DEL SONAR Y EL CARDUMEN
// ---------------------------------------------------------------------------

func plano08() {
	const N = 1000
	pasos := []int64{1500, 6000, 24000}
	g := &p08rnd{s: 20260816}
	niv := make([]int, N)    // level where each ball came to rest
	sg := make([]float64, N) // its sigma there
	entran, quedan := [4]int{}, [4]int{}
	cond, remo, remoTot := 0, 0.0, 0.0

	const t8 = 1e33
	nTop := math.Sqrt(t8 / (2 * math.Pi))
	k0 := 0.35 * nTop
	L := bandL(t8, k0)

	for i := 0; i < N; i++ {
		j, s := 0, 0.0
		for ; j < 3; j++ {
			entran[j]++
			s = g.sigma()
			remo += float64(pasos[j])
			if !escala2(s) {
				break // calm at this range: the wave stops here
			}
		}
		if j == 3 {
			entran[3]++
			s = g.sigma()
			remo += L
			if condena(s) {
				cond++
			}
		}
		quedan[j]++
		niv[i], sg[i] = j, s
		remoTot += L
	}
	// the condemnation rate of the judgment step, on a big deterministic sample
	g2, nc := &p08rnd{s: 7771}, 0
	const M = 200000
	for i := 0; i < M; i++ {
		if condena(g2.sigma()) {
			nc++
		}
	}
	tasa, ahorro := 100*float64(nc)/M, 100*(1-remo/remoTot)
	paso := L / k0
	fmt.Printf("SS plano08 cuelgan1=%d(%.1f%%) escalan=%.1f%% cuelgan2=%d cuelgan3=%d juicio=%d condenados=%d tasaJuicio=%.3f%% ahorroRemo=%.1f%% t=%.0e k=%.17g L=%.0f L/k=%.4e\n",
		quedan[0], 100*float64(quedan[0])/N, 100*float64(N-quedan[0])/N, quedan[1], quedan[2],
		entran[3], cond, tasa, ahorro, t8, k0, L, paso)

	l := NewLienzo(1400, 950)
	l.Titulo("⑧ EL EMBUDO DEL SONAR Y EL CARDUMEN",
		fmt.Sprintf("mil bloques caen por el embudo: %d se cuelgan en el primer escalón y sólo %d llegan al juicio — simulación medida en esta corrida", quedan[0], entran[3]))

	cx, yTop, yFin := 745.0, 116.0, 546.0
	hwT, hwF := 415.0, 156.0
	hw := func(y float64) float64 { return escala(y, yTop, yFin, hwT, hwF) }
	px := func(s, y float64) float64 {
		if s > 3 {
			s = 3
		}
		return escala(s, 0, 3, cx-hw(y), cx+hw(y))
	}
	// the funnel body and the two windows running down it
	l.Line(cx-hwT, yTop, cx-hwF, yFin, cDim, 2.6, "")
	l.Line(cx+hwT, yTop, cx+hwF, yFin, cDim, 2.6, "")
	l.Line(cx-hwT, yTop, cx+hwT, yTop, cGrid, 1.6, "5 5")
	for _, s := range []float64{anchoLo, anchoHi, juicioLo, juicioHi} {
		l.Line(px(s, yTop), yTop, px(s, yFin), yFin, cGrid, 1, "4 6")
	}
	yT := []float64{236, 348, 448, 536}
	cols := []string{cGreen, cBlue, cGold, cInk}
	for i, ys := range yT {
		lo, hi, col := anchoLo, anchoHi, cGreen
		if i == 3 {
			lo, hi, col = juicioLo, juicioHi, cBlue
		}
		xl, xr := px(lo, ys), px(hi, ys)
		l.Line(cx-hw(ys), ys, xl, ys, cRose, 1.3, "3 4") // the open slots
		l.Line(xr, ys, cx+hw(ys), ys, cRose, 1.3, "3 4")
		l.Line(xl, ys, xr, ys, col, 5.5, "") // the catching plate
		l.Circ(xl, ys, 3.4, col, "", 0)
		l.Circ(xr, ys, 3.4, col, "", 0)
		l.Flecha((cx-hw(ys)+xl)/2, ys+4, (cx-hw(ys)+xl)/2-7, ys+26, cRose, 1.6)
		l.Flecha((cx+hw(ys)+xr)/2, ys+4, (cx+hw(ys)+xr)/2+7, ys+26, cRose, 1.6)
		ip := i
		if ip > 2 {
			ip = 2
		}
		et := fmt.Sprintf("PASO %d · %d remadas", i+1, pasos[ip])
		if i == 3 {
			et = fmt.Sprintf("JUICIO · banda entera L=%.0f", L)
		}
		l.Txt(cx-hw(ys)-14, ys-6, 13.5, col, "end", et)
		l.Mono(cx-hw(ys)-14, ys+11, 11.5, cDim, "end", fmt.Sprintf("entran %d · cuelgan %d", entran[i], quedan[i]))
	}
	// the thousand balls, each drawn where it came to rest
	for i := 0; i < N; i++ {
		j := niv[i]
		if j == 3 && condena(sg[i]) {
			l.Circ(px(sg[i], yFin), yFin+10+g.u()*18, 2.8, cRose, "", 0)
			continue
		}
		y := yT[j] - 5 - g.u()*38
		l.Circ(px(sg[i], y), y, 1.95, cols[j], "", 0)
	}
	// the beast bins, drawn exactly under the two open slots of the judgment plate
	for _, s := range []struct {
		lo, hi float64
		txt    string
	}{{0, juicioLo, "MUDA  σ<0,05"}, {juicioHi, 3, "OLA  σ>2,40"}} {
		x0, x1 := px(s.lo, yFin), px(s.hi, yFin)
		l.Rect(x0, yFin+6, x1-x0, 26, "none", cRose, 0.85)
		l.Txt((x0+x1)/2, yFin+50, 11, cRose, "middle", s.txt)
	}
	l.Mono(px(1.0, yTop), yTop-10, 12, cGreen, "middle", "banda ancha 0,35 — 1,70")
	l.Mono(px(2.72, yTop), yTop-10, 12, cBlue, "middle", "banda de juicio 0,05 — 2,40")
	// the census, in the free water to the right of the funnel
	l.Txt(1010, 392, 14, cRose, "start", "EL CARDUMEN")
	l.Mono(1010, 414, 11.5, cInk, "start", fmt.Sprintf("llegan al juicio:  %d de %d", entran[3], N))
	l.Mono(1010, 432, 11.5, cRose, "start", fmt.Sprintf("bestias en esta corrida: %d", cond))
	l.Mono(1010, 450, 11.5, cDim, "start", fmt.Sprintf("el mar aleatorio condena %.3f%%", tasa))
	l.Mono(1010, 468, 11.5, cGreen, "start", fmt.Sprintf("remo ahorrado: %.1f%%", ahorro))
	l.Line(1010, 478, 1330, 478, cGrid, 1.2, "")
	xb := 1010.0 // the whole funnel as one measured bar
	for i, q := range []int{quedan[0], quedan[1], quedan[2], entran[3]} {
		w := 320 * float64(q) / float64(N)
		l.Rect(xb, 488, w, 12, cols[i], cols[i], 0.75)
		xb += w
	}
	l.Rect(1010, 488, 320, 12, "none", cDim, 0.6)
	l.Mono(1010, 516, 10.5, cDim, "start", "dónde quedaron los 1.000, por color de paso")

	// the contiguous strip and the magnifier
	l.Panel(45, 604, 1310, 154, cGold, "LA FRANJA CONTIGUA — los bloques se topan, sin huecos ni solapes: k ← k + L")
	rail := 692.0
	l.Line(90, rail, 890, rail, cDim, 1.6, "")
	for i := 0; i < 8; i++ {
		x := 90 + float64(i)*100
		col := cGold
		if i%2 == 1 {
			col = cGreen
		}
		l.Rect(x, rail-26, 100, 26, col, col, 0.18)
		l.Line(x, rail-32, x, rail+8, col, 1.6, "")
		l.Flecha(x+14, rail+22, x+86, rail+22, col, 1.4)
		l.Mono(x+50, rail+18, 10.5, cDim, "middle", "+L")
	}
	l.Line(890, rail-32, 890, rail+8, cGold, 1.6, "")
	l.Mono(90, rail-40, 11.5, cInk, "start", fmt.Sprintf("k = %.17g", k0))
	l.Mono(890, rail-40, 11.5, cInk, "end", fmt.Sprintf("L = %.0f", L))
	a, bs := fmt.Sprintf("%.17g", k0), fmt.Sprintf("%.17g", k0+L)
	pref := 0
	for pref < len(a) && pref < len(bs) && a[pref] == bs[pref] {
		pref++
	}
	mx, my := 1150.0, 682.0
	l.Line(mx+52, my+52, mx+74, my+74, cGold, 5, "")
	l.Circ(mx, my, 74, cPanel, cGold, 2.4)
	l.Mono(mx-62, my-14, 11, cDim, "start", a[:pref])
	l.Mono(mx-62+float64(pref)*6.6, my-14, 11, cGold, "start", a[pref:])
	l.Mono(mx-62, my+10, 11, cDim, "start", bs[:pref])
	l.Mono(mx-62+float64(pref)*6.6, my+10, 11, cGold, "start", bs[pref:])
	l.Line(mx-62+float64(pref)*6.6, my+16, mx-55+float64(pref)*6.6, my+16, cRose, 2.5, "")
	l.Txt(mx, my+42, 11, cRose, "middle", fmt.Sprintf("dígito %d", pref+1))
	l.Mono(mx, my-40, 11, cGold, "middle", fmt.Sprintf("L/k = %.3e", paso))

	l.Formula(45, 772, 1310, "sigma = |Z| / sqrt(L)   ·   cuelga si 0,35 < sigma < 1,70   ·   BESTIA si sigma < 0,05 o sigma > 2,40   ·   L(k) = k · (0,45/t)^(1/3)")
	l.Nota(45, 814, 1310, cGreen, []string{
		"EN CRIOLLO: el sonar no escucha el bloque entero de una. Tira un pulso corto de 1.500 remadas; si el eco suena a mar comun, cuelga y pasa al bloque siguiente.",
		fmt.Sprintf("Solo lo raro se gana pulsos mas largos (6.000, 24.000, banda entera). Aca %d de %d colgaron en el primer paso y apenas %d llegaron al juicio: se ahorro %.0f%% del remo.", quedan[0], N, entran[3], ahorro),
		fmt.Sprintf("La banda ancha (0,35 a 1,70) decide si se sigue escuchando; la banda de juicio (0,05 a 2,40) decide si es bestia. El mar aleatorio condena %.3f%% de los que llegan al final.", tasa),
		fmt.Sprintf("Abajo: los bloques se topan uno con otro. El paso L recien cambia el digito %d de k — por eso la direccion de una presa se escribe con %%.17g, o no se escribe.", pref+1),
	})
	l.Pie(fmt.Sprintf("medido en esta corrida · t=%.0e · k=%.17g · L=%.0f · L/k=%.3e", t8, k0, L, paso))
	l.Guardar("laminas/plano-08-el-embudo.svg")
}

// ---------------------------------------------------------------------------
// PLATE 09 - LOS TRES MOTORES
// ---------------------------------------------------------------------------

func plano09() {
	const t9 = 1e14 // calibration water: where the truth is still cheap
	nTop := math.Sqrt(t9 / (2 * math.Pi))
	k := 0.618 * nTop
	c, b, _ := puente(t9, k)
	n := int64(120000)
	fr, fi, dual := p09vuelta(b, c, n)
	dr, di := chirpDirect(b, c, n, true)
	mV, mD := math.Hypot(fr, fi), math.Hypot(dr, di)
	e := math.Hypot(fr-dr, fi-di) / math.Max(mD, 1e-9)
	bsc, nsc := cascada(b, c, n, 2000, true)
	// The damper's real contribution on a long block. Measured twice: at this
	// water's exact curvature, and at the same curvature nudged in its last
	// bits. The pair is the point - the reservoirs fill or stay dry depending
	// on the exact bits of b, which is the coordinate doctrine (F122) showing
	// up inside the rowing itself.
	nA := int64(2000000)
	ar, ai := chirpDirect(b, c, nA, true)
	xr, xi := chirpDirect(b, c, nA, false)
	difA := math.Hypot(ar-xr, ai-xi)
	bNudge := math.Nextafter(b, 1) // one ulp away, ~5.6e-17
	for i := 0; i < 140; i++ {
		bNudge = math.Nextafter(bNudge, 1)
	}
	yr, yi := chirpDirect(bNudge, c, nA, true)
	zr, zi := chirpDirect(bNudge, c, nA, false)
	difB := math.Hypot(yr-zr, yi-zi)
	// and how deep the curvature must go before the reservoirs fill at all
	var rb, rd []float64
	for p := 1; p <= 6; p++ {
		bb := frac(math.Pow(10, -float64(p)) * 1.2360679)
		u1, v1 := chirpDirect(bb, c, 400000, true)
		u2, v2 := chirpDirect(bb, c, 400000, false)
		rb = append(rb, bb)
		rd = append(rd, math.Hypot(u1-u2, v1-v2))
	}
	cen, dirF, fuente := p09censo()
	lo, hi, ghost, limpios := 1.0, 0.0, 0.0, 0
	for _, v := range cen {
		if v > 0.05 {
			ghost = math.Max(ghost, v)
			continue
		}
		lo, hi, limpios = math.Min(lo, v), math.Max(hi, v), limpios+1
	}
	fmt.Printf("SS plano09 b=%.12f c=%.12f n=%d dual=%d |vuelta|=%.4f |directo|=%.4f e=%.6f puerta=0.05 peldanos=%d amort(n=%d) exacta=%.2e corrida=%.2e censo=%d[%s] racimo=%.1e..%.1e fantasma=%.1e\n",
		b, c, n, dual, mV, mD, e, len(bsc), nA, difA, difB, len(cen), fuente, lo, hi, ghost)

	l := NewLienzo(1400, 950)
	l.Titulo("⑨ LOS TRES MOTORES — y el 5º fantasma con domicilio",
		fmt.Sprintf("dos motores rápidos que deben coincidir dentro del 5%% y un árbitro de 256 bits que no sabe mentir · censo de %d pinchazos leído de %s", len(cen), fuente))

	// --- panel A: the three engines writing the same number ---
	l.Panel(45, 100, 700, 500, cGreen, "LOS TRES ESCRIBEN EL MISMO NÚMERO")
	l.Txt(62, 150, 14, cGold, "start", "① LA CASCADA — cae de peldaño en peldaño")
	for i := 0; i < len(bsc) && i < 6; i++ {
		r := 12 + 22*math.Log10(float64(nsc[i]))/math.Log10(float64(n))
		cxg := 105 + float64(i)*118
		p09rueda(l, cxg, 200, r, cGold, 9)
		l.Mono(cxg, 244, 10.5, cDim, "middle", fmt.Sprintf("n=%d", nsc[i]))
		l.Mono(cxg, 258, 10, cGold, "middle", fmt.Sprintf("b=%.4f", bsc[i]))
		if i+1 < len(bsc) && i < 5 {
			l.Flecha(cxg+r+7, 200, 105+float64(i+1)*118-r-9, 200, cGreen, 1.8)
		}
	}
	l.Txt(62, 294, 14, cBlue, "start", "② EL JUEZ DIRECTO — una remada por término, sin atajos")
	for i := 0; i < 130; i++ {
		x := 70 + float64(i)*4.9
		l.Line(x, 307, x, 332, cBlue, 1.5, "")
	}
	l.Mono(70, 348, 11, cDim, "start", fmt.Sprintf("n = %d remadas, una por una", n))
	l.Txt(62, 384, 14, cRose, "start", "③ EL ÁRBITRO DE 256 BITS — anda al paso del buey")
	l.Rect(78, 420, 54, 30, cRose, cRose, 0.25) // the ox
	l.Line(78, 420, 68, 408, cRose, 2, "")
	l.Line(98, 420, 92, 407, cRose, 2, "")
	for i := 0; i < 4; i++ {
		l.Line(84+float64(i)*14, 450, 84+float64(i)*14, 466, cRose, 2, "")
	}
	l.Line(132, 438, 176, 438, cRose, 2.4, "") // the shaft
	l.Circ(196, 462, 15, cPanel, cRose, 2)     // the cart
	l.Circ(252, 462, 15, cPanel, cRose, 2)
	l.Rect(176, 404, 100, 48, cPanel, cRose, 0.9)
	for i := 0; i < 5; i++ { // the huge ledger it carries
		l.Line(184, 414+float64(i)*8, 268, 414+float64(i)*8, cDim, 1, "")
	}
	l.Mono(292, 426, 11, cDim, "start", "256 bits de mantisa,")
	l.Mono(292, 444, 11, cDim, "start", "sin el cuaderno dd de los otros dos")
	// the ledger all three sign
	l.Rect(330, 486, 390, 96, cPanel, cGold, 0.95)
	for i := 0; i < 5; i++ {
		l.Line(340, 502+float64(i)*17, 710, 502+float64(i)*17, cGrid, 1, "")
	}
	l.Txt(525, 500, 12.5, cGold, "middle", "EL LIBRO DEL BLOQUE")
	l.Mono(348, 522, 11.5, cGold, "start", fmt.Sprintf("cascada  |Z| = %.4f", mV))
	l.Mono(348, 539, 11.5, cBlue, "start", fmt.Sprintf("directo  |Z| = %.4f", mD))
	l.Mono(348, 556, 11.5, cRose, "start", fmt.Sprintf("árbitro  e256 = %.1e a %.1e", lo, hi))
	l.Mono(348, 573, 11.5, cInk, "start", fmt.Sprintf("dual del primer giro = %d términos", dual))
	for i, col := range []string{cGold, cBlue, cRose} {
		y := 520 + float64(i)*20
		l.Flecha(298, y, 326, y, col, 1.8)
	}

	// --- panel B: the two needles on one dial ---
	l.Panel(760, 100, 595, 290, cGold, "LAS DOS AGUJAS — tienen que coincidir dentro del 5%")
	dcx, dcy, dr2 := 1057.0, 238.0, 100.0
	l.Circ(dcx, dcy, dr2, cPanel, cGold, 2.4)
	loD, hiD := mD*0.90, mD*1.10
	for i := 0; i <= 40; i++ { // the acceptance wedge, hatched
		v := mD*0.95 + float64(i)/40*mD*0.10
		a := (150 + 240*math.Max(0, math.Min(1, (v-loD)/(hiD-loD)))) * math.Pi / 180
		l.Line(dcx+dr2*0.56*math.Cos(a), dcy+dr2*0.56*math.Sin(a), dcx+dr2*0.92*math.Cos(a), dcy+dr2*0.92*math.Sin(a), cGreen, 1.4, "")
	}
	for i := 0; i <= 10; i++ {
		a := (150 + 240*float64(i)/10) * math.Pi / 180
		l.Line(dcx+dr2*0.93*math.Cos(a), dcy+dr2*0.93*math.Sin(a), dcx+dr2*math.Cos(a), dcy+dr2*math.Sin(a), cDim, 1.5, "")
	}
	for _, m := range []struct {
		f   float64
		txt string
	}{{0, "−10%"}, {0.25, "−5%"}, {0.75, "+5%"}, {1, "+10%"}} {
		a := (150 + 240*m.f) * math.Pi / 180
		l.Mono(dcx+(dr2+16)*math.Cos(a), dcy+(dr2+16)*math.Sin(a)+4, 10, cDim, "middle", m.txt)
	}
	p09aguja(l, dcx, dcy, dr2, loD, hiD, mD, cBlue, 4.4)
	p09aguja(l, dcx, dcy, dr2, loD, hiD, mV, cGold, 2.2)
	l.Circ(dcx, dcy, 6, cInk, "", 0)
	l.Mono(dcx, dcy+dr2+24, 12, cInk, "middle", fmt.Sprintf("|cascada| %.3f    |directo| %.3f", mV, mD))
	l.Mono(dcx, dcy+dr2+42, 12.5, cGreen, "middle", fmt.Sprintf("e = %.2e  «  puerta 0,05", e))
	l.Txt(790, 152, 12, cBlue, "start", "aguja gruesa: el directo")
	l.Txt(790, 172, 12, cGold, "start", "aguja fina: la cascada")
	l.Txt(790, 192, 12, cGreen, "start", "franja rayada: ±5% permitido")
	l.Txt(790, 212, 12, cRose, "start", "si se pelean, la presa")
	l.Txt(790, 232, 12, cRose, "start", "se descarta sin apelación")

	// --- panel C: what the F144 damper really contributes ---
	l.Panel(760, 402, 595, 198, cRose, "EL AMORTIGUADOR F144 — cuánto aporta, medido")
	bx, by, bw, bh := 800.0, 452.0, 520.0, 100.0
	l.Ejes(bx, by, bw, bh, "curvatura b", "|con − sin amortiguador|")
	for i := range rd {
		x := bx + (float64(i)+0.9)*bw/7.0
		h := escala(math.Log10(math.Max(rd[i], 1e-9)), -8, -2, 0, bh)
		l.Rect(x-14, by+bh-h, 28, h, cRose, cRose, 0.55)
		l.Mono(x, by+bh+14, 9.5, cDim, "middle", fmt.Sprintf("%.0e", rb[i]))
	}
	xz := bx + 0.12*bw/7.0
	l.Rect(xz-10, by+bh-4, 24, 4, "none", cGreen, 0.9)
	l.Mono(xz+2, by+bh+14, 9.5, cGreen, "middle", "b real")
	l.Mono(bx+bw/2, by+bh+38, 11, cGreen, "middle",
		fmt.Sprintf("con la curvatura EXACTA de esta agua (b=%.15f, n=%d): dif = %.1e", b, nA, difA))
	l.Mono(bx+bw/2, by+bh+56, 11, cRose, "middle",
		fmt.Sprintf("corrida %.0e en los últimos bits (b=%.15f): dif = %.1e", bNudge-b, bNudge, difB))
	l.Txt(bx+bw/2, by+bh+76, 11.5, cGold, "middle",
		"el reservorio se llena o queda seco según los BITS de la curvatura: la doctrina de la coordenada, adentro del remo")

	// --- panel D: the arbiter's census ---
	l.Panel(45, 610, 1310, 148, cBlue, "EL CENSO DEL ÁRBITRO — cada pinchazo de 256 bits, y el fantasma con domicilio")
	ax, ay, aw := 90.0, 700.0, 1200.0
	l.Line(ax, ay, ax+aw, ay, cDim, 1.6, "")
	lg := func(v float64) float64 { return escala(math.Log10(v), -5.4, 0.1, ax, ax+aw) }
	for d := -5; d <= 0; d++ {
		x := lg(math.Pow(10, float64(d)))
		l.Line(x, ay, x, ay+7, cDim, 1.4, "")
		l.Mono(x, ay+20, 10.5, cDim, "middle", fmt.Sprintf("1e%d", d))
	}
	l.Line(lg(0.05), ay+10, lg(0.05), 638, cGold, 2, "5 4")
	l.Txt(lg(0.05), 632, 11.5, cGold, "middle", "PUERTA 5% — juez y árbitro")
	pila := map[int]int{}
	for _, v := range cen {
		x := lg(v)
		key := int(x / 7)
		col, r := cBlue, 4.0
		if v > 0.05 {
			col, r = cRose, 6.0
		}
		y := ay - 9 - float64(pila[key])*10
		pila[key]++
		l.Line(x, ay, x, y, col, 1.2, "")
		l.Circ(x, y, r, col, "", 0)
	}
	l.Line(lg(lo), ay+32, lg(hi), ay+32, cBlue, 2, "")
	l.Line(lg(lo), ay+27, lg(lo), ay+37, cBlue, 2, "")
	l.Line(lg(hi), ay+27, lg(hi), ay+37, cBlue, 2, "")
	l.Mono((lg(lo)+lg(hi))/2, ay+48, 11, cBlue, "middle", fmt.Sprintf("%d pinchazos limpios, de %.1e a %.1e", limpios, lo, hi))
	l.Flecha(lg(ghost)-62, 668, lg(ghost)-9, 686, cRose, 1.8)
	l.Mono(lg(ghost)-66, 652, 11, cRose, "end", fmt.Sprintf("FANTASMA 5º · e256 = %.1e", ghost))
	l.Mono(lg(ghost)-66, 668, 10.5, cRose, "end", dirF)

	l.Formula(45, 770, 1310, "|Z_cascada − Z_directo| / |Z_directo| < 0,05      y      e256 = |Z_dd − Z_256| / |Z_256| ≤ 0,05")
	l.Nota(45, 812, 1310, cBlue, []string{
		"EN CRIOLLO: el mismo bloque se calcula por dos caminos distintos — el atajo (la cascada) y el camino largo (una remada por termino) — y los dos numeros tienen que darse la mano.",
		fmt.Sprintf("Aca se dieron la mano con una diferencia de %.2e, muy por debajo de la puerta del 5%%. Si se pelean, la presa se descarta: no hay bestia sin dos firmas.", e),
		"Detras espera el tercero: un arbitro de 256 bits, lento como un carro de bueyes, que no usa el cuaderno de los otros dos. Se lo llama porque dos motores parecidos pueden equivocarse igual.",
		fmt.Sprintf("Y una vez el arbitro grito: en %s el error salto a %.1e. Ese fantasma tiene domicilio exacto — es el borde donde el cuaderno de doble precision se agota.", dirF, ghost),
	})
	l.Pie(fmt.Sprintf("medido en esta corrida · e(juez)=%.2e · amortiguador en la dirección real = %.1e · censo %d pinchazos · fantasma %.1e", e, difA, len(cen), ghost))
	l.Guardar("laminas/plano-09-los-tres-motores.svg")
}

// ---------------------------------------------------------------------------
// PLATE 10 - EL ENGRANAJE CONGELADO (la lámina honesta)
// ---------------------------------------------------------------------------

func plano10() {
	const t10 = 1e34
	nTop := math.Sqrt(t10 / (2 * math.Pi))
	kk := 0.35 * nTop
	_, b0, _ := puente(t10, kk)
	franja := make([]float64, 0, 1500)
	deriva := 0.0
	for i := 0; i < 1500; i++ {
		_, bb, _ := puente(t10, kk)
		franja = append(franja, bb)
		if d := math.Abs(bb - b0); d > deriva {
			deriva = d
		}
		kk += bandL(t10, kk)
	}
	const nBlk = 49664
	bRac, bIrr := 4.0/49.0, frac(math.Sqrt(2)/3)
	muR := make([]float64, 0, 200)
	muI := make([]float64, 0, 200)
	for i := 0; i < 200; i++ {
		cc := frac(0.6180339887498949 * float64(i+1))
		x1, y1 := chirpDirect(bRac, cc, nBlk, true)
		x2, y2 := chirpDirect(bIrr, cc, nBlk, true)
		muR = append(muR, math.Hypot(x1, y1))
		muI = append(muI, math.Hypot(x2, y2))
	}
	sR := append([]float64(nil), muR...)
	sI := append([]float64(nil), muI...)
	sort.Float64s(sR)
	sort.Float64s(sI)
	medR, medI := sR[100], sI[100]
	rq, rL := math.Sqrt(49), math.Sqrt(nBlk)
	fmt.Printf("SS plano10 b0=%.12f (4/49=%.12f) deriva=%.3e bloques=%d bRac=%.10f medianaRac=%.3f sqrt(49)=%.1f bIrr=%.10f medianaIrr=%.3f sqrt(L)=%.2f L=%d\n",
		b0, 4.0/49, deriva, len(franja), bRac, medR, rq, bIrr, medI, rL, nBlk)

	l := NewLienzo(1400, 950)
	l.Titulo("⑩ EL ENGRANAJE CONGELADO — la lámina honesta",
		fmt.Sprintf("1.500 bloques distintos, un solo engranaje: b se quedó clavado en 4/49 con una deriva máxima de %.1e en toda la franja", deriva))

	// --- the strip and the needle that never moved ---
	l.Panel(45, 100, 1310, 222, cRose, "LA FRANJA DE 1.500 BLOQUES CONTIGUOS — la aguja de curvatura no se mueve")
	p09rueda(l, 140, 210, 46, cRose, 12)
	p09aguja(l, 140, 210, 46, 0, 0.5, b0, cGold, 3.4)
	l.Mono(140, 274, 11.5, cGold, "middle", fmt.Sprintf("b = %.9f", b0))
	l.Mono(140, 292, 11.5, cRose, "middle", "= 4/49, de punta a punta")
	sx0, sx1, sy := 250.0, 1330.0, 296.0
	l.Line(sx0, sy, sx1, sy, cDim, 1.4, "")
	for i := 0; i < len(franja); i += 10 {
		x := escala(float64(i), 0, 1499, sx0, sx1)
		l.Line(x, sy, x, sy-escala(franja[i], 0, 0.12, 0, 84), cGold, 2.6, "")
	}
	for _, i := range []int{0, 749, 1499} {
		x := escala(float64(i), 0, 1499, sx0, sx1)
		l.Line(x, 212, x, sy-escala(franja[i], 0, 0.12, 0, 84)-4, cGrid, 1, "4 4")
		p09rueda(l, x, 166, 22, cGold, 8)
		p09aguja(l, x, 166, 22, 0, 0.5, franja[i], cRose, 2)
		l.Mono(x, 204, 10, cDim, "middle", fmt.Sprintf("bloque %d", i+1))
	}
	l.Mono(sx0, sy+16, 10.5, cRose, "start", fmt.Sprintf("deriva máxima: %.1e", deriva))
	l.Mono(sx1, sy+16, 10.5, cDim, "end", "una barra cada diez bloques · altura = b medido en ese bloque")

	// --- the law that froze it ---
	l.Panel(45, 334, 470, 356, cGold, "LA LEY QUE LO CONGELÓ:  b_raw = 1/(2·frac0²)")
	lx, ly, lw, lh := 100.0, 386.0, 380.0, 222.0
	l.Ejes(lx, ly, lw, lh, "frac0 = k/nTop", "b_raw")
	pts := make([][2]float64, 0, 121)
	for i := 0; i <= 120; i++ {
		f0 := 0.3 + float64(i)/120*0.7
		v := math.Min(1/(2*f0*f0), 6)
		pts = append(pts, [2]float64{escala(f0, 0.3, 1.0, lx, lx+lw), escala(v, 0, 6, ly+lh, ly)})
	}
	l.Camino(pts, cGold, 2.4, 1)
	for d := 1; d <= 5; d++ {
		y := escala(float64(d), 0, 6, ly+lh, ly)
		l.Line(lx, y, lx+lw, y, cGrid, 1, "3 5")
		l.Mono(lx-6, y+4, 9.5, cDim, "end", fmt.Sprintf("%d", d))
	}
	for _, m := range []struct {
		f0  float64
		txt string
	}{{0.35, "200/49"}, {0.70, "50/49"}} {
		v := 1 / (2 * m.f0 * m.f0)
		x := escala(m.f0, 0.3, 1.0, lx, lx+lw)
		y := escala(math.Min(v, 6), 0, 6, ly+lh, ly)
		l.Line(x, ly+lh, x, y, cRose, 1.4, "4 4")
		l.Circ(x, y, 5.5, cRose, "", 0)
		l.Mono(x+9, y-5, 11, cRose, "start", m.txt)
		l.Mono(x, ly+lh+16, 10.5, cRose, "middle", fmt.Sprintf("frac0 = %.2f", m.f0))
	}
	for i, m := range []struct{ f0 float64 }{{0.35}, {0.70}} {
		v := 1 / (2 * m.f0 * m.f0)
		l.Mono(lx+lw, ly+34+float64(i)*18, 10.5, cGold, "end",
			fmt.Sprintf("frac0=%.2f → b_raw=%.4f → b=frac= %.6f", m.f0, v, frac(v)))
	}
	l.Mono(lx, ly+lh+42, 10.5, cDim, "start", "el t se cancela solo: la curvatura del bloque")
	l.Mono(lx, ly+lh+58, 10.5, cDim, "start", "depende de DÓNDE está dentro de su propio mar,")
	l.Mono(lx, ly+lh+74, 10.5, cDim, "start", "no de qué tan profundo es el mar.")

	// --- the two histograms ---
	l.Panel(527, 334, 828, 356, cBlue, "DOS HISTOGRAMAS DE |ola| — 200 direcciones c, bloque de 49.664 términos")
	hx, hy, hw2, hh := 580.0, 390.0, 720.0, 212.0
	l.Ejes(hx, hy, hw2, hh, "|ola|", "")
	const nb, top = 36, 460.0
	ha, hb := p10hist(muR, nb, top), p10hist(muI, nb, top)
	mxc := 1
	for i := 0; i < nb; i++ {
		if ha[i] > mxc {
			mxc = ha[i]
		}
		if hb[i] > mxc {
			mxc = hb[i]
		}
	}
	l.Mono(hx-8, hy+12, 9.5, cDim, "end", fmt.Sprintf("%d", mxc))
	l.Mono(hx-8, hy+hh, 9.5, cDim, "end", "0")
	pasoB := hw2 / nb
	for i := 0; i < nb; i++ {
		x := hx + float64(i)*pasoB
		if ha[i] > 0 {
			h := float64(ha[i]) / float64(mxc) * hh
			l.Rect(x, hy+hh-h, pasoB-1, h, cGold, cGold, 0.55)
			if ha[i] == mxc {
				l.Mono(x+pasoB+8, hy+hh-h+14, 11, cGold, "start", fmt.Sprintf("%d de 200 acá", ha[i]))
			}
		}
		if hb[i] > 0 {
			h := float64(hb[i]) / float64(mxc) * hh
			l.Rect(x+pasoB*0.25, hy+hh-h, pasoB*0.5, h, cBlue, cBlue, 0.8)
		}
	}
	for _, m := range []struct {
		v   float64
		col string
		txt string
	}{{rq, cGold, fmt.Sprintf("√49 = %.0f", rq)}, {rL, cBlue, fmt.Sprintf("√L = %.0f", rL)}} {
		x := escala(m.v, 0, top, hx, hx+hw2)
		l.Line(x, hy, x, hy+hh, m.col, 1.6, "6 5")
		l.Txt(x, hy-8, 11.5, m.col, "middle", m.txt)
	}
	for _, m := range []struct {
		v   float64
		col string
		txt string
	}{{medR, cGold, fmt.Sprintf("mediana %.2f", medR)}, {medI, cBlue, fmt.Sprintf("mediana %.1f", medI)}} {
		x := escala(m.v, 0, top, hx, hx+hw2)
		l.Flecha(x, hy+hh+36, x, hy+hh+8, m.col, 2)
		l.Mono(x, hy+hh+50, 11, m.col, "middle", m.txt)
	}
	l.Mono(hx+hw2, hy+16, 11.5, cGold, "end", fmt.Sprintf("b racional 4/49 = %.7f", bRac))
	l.Mono(hx+hw2, hy+34, 11.5, cBlue, "end", fmt.Sprintf("b irracional frac(√2/3) = %.7f", bIrr))
	l.Mono(hx, hy+hh+70, 10.5, cDim, "start", "eje vertical: cuántas de las 200 direcciones caen en cada casilla")

	l.Formula(45, 704, 1310, "b = frac( 1/(2·frac0²) )   ·   frac0 = 0,35 → 200/49 → b = 4/49   ·   |ola| ≈ sqrt(q) con b = p/q racional,   ≈ sqrt(L) con b irracional")
	l.Nota(45, 746, 1310, cRose, []string{
		"EN CRIOLLO: el laboratorio barrio 1.500 bloques pegados uno al otro creyendo que medía mil quinientos pedazos distintos del mar.",
		fmt.Sprintf("Pero la curvatura de un bloque solo depende de la fraccion frac0 = k/nTop, y esa fraccion estaba fija: la aguja quedo clavada en b = 4/49 (deriva %.1e en toda la franja).", deriva),
		fmt.Sprintf("Un b racional de denominador 49 apila la ola cerca de raiz de 49 = %.0f. Medido: mediana %.2f. Con un b irracional la misma cuenta se desparrama: mediana %.1f, cerca de raiz de L = %.0f.", rq, medR, medI, rL),
		"O sea: el censo no estaba escuchando mil quinientos mares. Estaba escuchando mil quinientas veces el ruido de UNA sola pieza.",
		"Esta lamina no celebra nada. Queda colgada en el museo porque una maquina que no muestra sus trampas no es una maquina: es un cuento.",
	})
	l.Pie(fmt.Sprintf("medido en esta corrida · b0=%.9f · deriva=%.1e · mediana racional %.2f vs √49=%.0f · mediana irracional %.1f vs √L=%.0f", b0, deriva, medR, rq, medI, rL))
	l.Guardar("laminas/plano-10-el-engranaje-congelado.svg")
}
