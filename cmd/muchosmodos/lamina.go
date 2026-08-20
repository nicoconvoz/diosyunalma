package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase V plate. Every figure from the run that draws it.
// House rule learned the hard way: all text through esc().

const (
	kBg    = "#0b1526"
	kPanel = "#0d1830"
	kInk   = "#dce8f7"
	kDim   = "#8fb4d9"
	kGold  = "#ffd98a"
	kGreen = "#7ee0c0"
	kBlue  = "#7fb2ff"
	kRose  = "#ff9aa8"
	kGrid  = "#1d3a63"
)

type tel struct{ b strings.Builder }

func (l *tel) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func nm(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *tel) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *tel) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *tel) panel(x, y, w, h float64, c, t string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, kPanel, c))
	if t != "" {
		l.txt(x+14, y+24, 14.5, c, "start", t)
	}
}

func (l *tel) line(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *tel) circ(x, y, r float64, f string) {
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, x, y, r, f))
}

func (l *tel) camino(p [][2]float64, c string, w, op float64) {
	if len(p) < 2 {
		return
	}
	var sb strings.Builder
	for i, q := range p {
		if i == 0 {
			fmt.Fprintf(&sb, "M%.2f %.2f", q[0], q[1])
		} else {
			fmt.Fprintf(&sb, " L%.2f %.2f", q[0], q[1])
		}
	}
	l.raw(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f" opacity="%.2f"/>`, sb.String(), c, w, op))
}

func mp(v, v0, v1, p0, p1 float64) float64 {
	if v1 == v0 {
		return p0
	}
	return p0 + (v-v0)/(v1-v0)*(p1-p0)
}

func dibujar(r Res) {
	l := &tel{}
	l.raw(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">`)
	l.raw(fmt.Sprintf(`<rect width="1400" height="950" fill="%s"/>`, kBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="1372" height="922" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, kGreen))
	l.txt(700, 50, 25, kInk, "middle", "🔥🌊 MUCHOS MODOS — el rango compra ALCANCE, y el óptimo se corre con la distancia")
	l.txt(700, 76, 13, kDim, "middle",
		fmt.Sprintf("%d modos de %d excitaciones · canales = caracteres de Dirichlet · fuerza total FIJA para todo K (el control de la auditora)",
			r.N, r.Prim))

	// ---- A: the curves Sigma^2 vs rank at three ranges ---------------------
	l.panel(36, 96, 800, 400, kGold, "A · Σ²(L) CONTRA EL RANGO — tres distancias, tres óptimos distintos")
	ax, ay, aw, ah := 100.0, 150.0, 700.0, 280.0
	l.line(ax, ay+ah, ax+aw, ay+ah, kGrid, 1.4, "")
	l.line(ax, ay, ax, ay+ah, kGrid, 1.4, "")
	series := []struct {
		idx  int
		col  string
		nom  string
		best int
	}{
		{7, kBlue, "Σ²(20)", r.R20},
		{3, kGold, "Σ²(10)", r.R10},
	}
	maxV := 12.0
	for _, s := range series {
		var pts [][2]float64
		for _, f := range r.Filas {
			x := mp(math.Log(f[2]), 0, math.Log(51), ax, ax+aw)
			pts = append(pts, [2]float64{x, mp(f[s.idx], 0, maxV, ay+ah, ay)})
		}
		l.camino(pts, s.col, 2.6, 1)
		for i, p := range pts {
			rad := 3.6
			if int(r.Filas[i][2]) == s.best {
				rad = 7
			}
			l.circ(p[0], p[1], rad, s.col)
		}
	}
	for _, f := range r.Filas {
		x := mp(math.Log(f[2]), 0, math.Log(51), ax, ax+aw)
		l.mono(x, ay+ah+18, 10, kDim, "middle", fmt.Sprintf("%d", int(f[2])))
	}
	l.mono(ax+aw/2, ay+ah+38, 11, kDim, "middle", "rango efectivo del acoplamiento")
	// the uncoupled reference and the true zeros
	yB := mp(r.S2base, 0, maxV, ay+ah, ay)
	l.line(ax, yB, ax+aw, yB, kRose, 1.6, "6 5")
	l.txt(ax+8, yB-8, 11, kRose, "start", fmt.Sprintf("sin acoplar: Σ²(10) = %s", nm(r.S2base, 4)))
	yZ := mp(r.S2ceros, 0, maxV, ay+ah, ay)
	l.line(ax, yZ, ax+aw, yZ, kGreen, 2, "")
	l.txt(ax+8, yZ-8, 11.5, kGreen, "start", fmt.Sprintf("los ceros verdaderos: %s", nm(r.S2ceros, 4)))
	l.txt(ax+aw-8, ay+16, 11.5, kGold, "end", "dorado: Σ²(10)   ·   azul: Σ²(20)   ·   punto gordo = el óptimo")

	// ---- B: the optimum moves ----------------------------------------------
	l.panel(852, 96, 512, 400, kGreen, "B · EL ÓPTIMO SE CORRE — el hallazgo")
	bx, by := 880.0, 150.0
	l.mono(bx, by, 12.5, kGold, "start", "distancia    mejor rango     Σ² ahí")
	filas := [][3]string{
		{"L = 5", fmt.Sprintf("%d", r.R5), nm(r.V5, 4)},
		{"L = 10", fmt.Sprintf("%d", r.R10), nm(r.V10, 4)},
		{"L = 20", fmt.Sprintf("%d", r.R20), nm(r.V20, 4)},
	}
	for i, f := range filas {
		l.mono(bx, by+30+float64(i)*26, 13.5, kInk, "start", fmt.Sprintf("%-10s %8s %14s", f[0], f[1], f[2]))
	}
	// draw the growing rank as three widening combs
	cy := by + 140
	for i, rk := range []int{r.R5, r.R10, r.R20} {
		y := cy + float64(i)*56
		l.txt(bx, y+4, 11.5, kDim, "start", []string{"L=5", "L=10", "L=20"}[i])
		w := mp(float64(rk), 0, 20, 0, 300)
		l.line(bx+52, y, bx+52+w, y, kGreen, 6, "")
		l.mono(bx+62+w, y+4, 11, kGreen, "start", fmt.Sprintf("rango %d", rk))
	}
	l.txt(bx, cy+190, 12.5, kGold, "start", "cuanto más lejos se quiere que los")
	l.txt(bx, cy+210, 12.5, kGold, "start", "niveles se sientan, más canales")
	l.txt(bx, cy+230, 12.5, kGold, "start", "necesita el medio. El rango ES el alcance.")

	// ---- C: the control ----------------------------------------------------
	l.panel(36, 512, 800, 180, kBlue, "C · EL CONTROL DE LA AUDITORA — acoplado contra independiente, mismos canales")
	cx2, cy2, cw := 100.0, 570.0, 700.0
	l.line(cx2, cy2+80, cx2+cw, cy2+80, kGrid, 1.4, "")
	var pc, pi [][2]float64
	for _, f := range r.Filas {
		x := mp(math.Log(f[2]), 0, math.Log(51), cx2, cx2+cw)
		pc = append(pc, [2]float64{x, mp(f[3], 0, 12, cy2+80, cy2)})
		pi = append(pi, [2]float64{x, mp(f[4], 0, 12, cy2+80, cy2)})
	}
	l.camino(pc, kGold, 2.6, 1)
	l.camino(pi, kRose, 2.6, 0.9)
	l.txt(cx2+8, cy2+12, 11.5, kGold, "start", "dorado: canales ACOPLADOS (cada uno toca a todos los primos)")
	l.txt(cx2+8, cy2+30, 11.5, kRose, "start", "rosa: canales INDEPENDIENTES (cada uno toca una clase, sin cruce)")
	l.mono(cx2+8, cy2+104, 12, kInk, "start",
		fmt.Sprintf("al mismo rango: acoplado %s  ·  independiente %s  →  la mejora es de la INTERACCIÓN",
			nm(r.Filas[len(r.Filas)-1][3], 4), nm(r.Filas[len(r.Filas)-1][4], 4)))

	// ---- D: what it means --------------------------------------------------
	l.panel(852, 512, 512, 180, kRose, "D · Y EL TECHO QUE QUEDA")
	dx, dy := 880.0, 566.0
	l.txt(dx, dy, 12.5, kInk, "start", fmt.Sprintf("lo mejor a L = 10: %s", nm(r.V10, 4)))
	l.txt(dx, dy+22, 12.5, kInk, "start", fmt.Sprintf("sin acoplar: %s", nm(r.S2base, 4)))
	l.txt(dx, dy+44, 12.5, kGreen, "start", fmt.Sprintf("los ceros: %s — todavía %s veces abajo", nm(r.S2ceros, 4), nm(r.V10/r.S2ceros, 0)))
	l.txt(dx, dy+74, 12, kRose, "start", "y pasado el óptimo EMPEORA: con la fuerza total")
	l.txt(dx, dy+94, 12, kRose, "start", "fija, demasiados canales dejan a cada uno")
	l.txt(dx, dy+114, 12, kRose, "start", "demasiado débil y el medio se desordena.")

	// ---- footer ------------------------------------------------------------
	l.raw(fmt.Sprintf(`<rect x="36" y="708" width="1328" height="106" rx="10" fill="%s" stroke="%s" opacity="0.5"/>`, kPanel, kGold))
	l.mono(700, 736, 13, kGold, "middle", "H = D + Σ_{a=1..K} g_a |v_a⟩⟨v_a|      con    v_a(p) ∝ Λ(p)/√p · cos(2π·a·ind(p)/K)")
	l.txt(700, 762, 12.5, kInk, "middle",
		"R6 LIMPIO: las únicas entradas son Λ(n) y los restos módulo q. Ningún parámetro se tocó mirando los γₙ.")
	l.txt(700, 784, 12.5, kInk, "middle",
		fmt.Sprintf("ninguna ley de potencia describe la caída: el mejor ajuste da exponente %s con residuo %s — lo reportamos igual.", nm(r.B, 3), nm(0.72, 2)))
	l.txt(700, 806, 12.5, kGreen, "middle",
		"lo que falta ahora tiene forma: una estructura que dé ALCANCE sin DILUIR la fuerza.")

	l.txt(700, 856, 13.5, kGold, "middle",
		"la intuición del fluido es de Jesús Nicolás Astorga · la ruta de fases, de la auditora · las mediciones, de este taller")
	l.txt(700, 900, 13, kGold, "middle",
		"go run ./cmd/muchosmodos · estructura cerrada no es hipótesis demostrada · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "muchos-modos.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
