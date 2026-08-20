package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase VIII plate. Its two honest edges are drawn as large as
// its positive result: the ordered-but-not-arithmetic control, and the
// localisation that shadows the node's improvement.

const (
	pBg    = "#0b1526"
	pPanel = "#0d1830"
	pInk   = "#dce8f7"
	pDim   = "#8fb4d9"
	pGold  = "#ffd98a"
	pGreen = "#7ee0c0"
	pBlue  = "#7fb2ff"
	pRose  = "#ff9aa8"
	pGrid  = "#1d3a63"
)

type lienzo struct{ b strings.Builder }

func (l *lienzo) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func nm(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *lienzo) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *lienzo) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(t)))
}

func (l *lienzo) panel(x, y, w, h float64, c, t string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, pPanel, c))
	if t != "" {
		l.txt(x+15, y+25, 14.5, c, "start", t)
	}
}

func (l *lienzo) rect(x, y, w, h float64, f, s string, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" opacity="%.2f"/>`, x, y, w, h, f, s, op))
}

func (l *lienzo) line(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *lienzo) circ(x, y, r float64, f string) {
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, x, y, r, f))
}

func (l *lienzo) camino(p [][2]float64, c string, w float64, d string) {
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
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f"%s/>`, sb.String(), c, w, dd))
}

func mp(v, v0, v1, q0, q1 float64) float64 {
	if v1 == v0 {
		return q0
	}
	return q0 + (v-v0)/(v1-v0)*(q1-q0)
}

func dibujar(r Res) {
	l := &lienzo{}
	l.raw(`<svg xmlns="http://www.w3.org/2000/svg" width="1400" height="950" viewBox="0 0 1400 950">`)
	l.raw(fmt.Sprintf(`<rect width="1400" height="950" fill="%s"/>`, pBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="1372" height="922" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, pGreen))
	l.txt(700, 50, 25, pInk, "middle", "🌊🧱 POROSIDAD — el medio NO es uniforme, y eso se mide")
	l.txt(700, 76, 13.5, pGold, "middle",
		"«¿y si en este sistema hay cierta porosidad: mayor o menor, materiales más blandos y otros más duros?» — Jesús Nicolás Astorga")

	// ---- A: the arms, as bars ----------------------------------------------
	l.panel(36, 96, 700, 400, pGold, "A · LOS BRAZOS — misma fuerza total, sólo cambia el MATERIAL")
	ax, ay, ah := 250.0, 150.0, 30.0
	barras := []struct {
		nom string
		v   float64
		col string
	}{
		{"homogéneo", r.Base.s10, pDim},
		{"mezclado (5 semillas)", r.MezMedia, pRose},
		{"rampa lisa", r.Rampa.s10, pRose},
		{"onda lisa", r.Onda.s10, pBlue},
		{"ARITMÉTICO 1/log p", r.Arit.s10, pGreen},
		{"ARITMÉTICO divisores", r.Div.s10, pGreen},
		{"bloques suaves (×2)", r.Locs[1][2], pGold},
		{"homogéneo + NODO", r.BaseN.s10, pGold},
		{"ARITMÉTICO + NODO", r.AritN.s10, pGold},
	}
	maxV := 32.0
	for i, b := range barras {
		y := ay + float64(i)*ah
		w := mp(b.v, 0, maxV, 0, 420)
		l.rect(ax, y, w, ah-9, b.col, b.col, 0.55)
		l.txt(ax-10, y+ah-14, 11.5, pInk, "end", b.nom)
		l.mono(ax+w+8, y+ah-14, 11.5, b.col, "start", nm(b.v, 3))
	}
	l.line(ax+mp(r.MezMedia, 0, maxV, 0, 420), ay-6, ax+mp(r.MezMedia, 0, maxV, 0, 420), ay+9*ah, pRose, 1.4, "5 5")

	// the material itself, drawn: the arithmetic permeability field and the same
	// values shuffled. Same histogram, same mean, same spread - only the ORDER differs.
	tira := func(y float64, h []float64, col, rot string) {
		w := 420.0 / float64(len(h))
		lo, hi := h[0], h[0]
		for _, v := range h {
			if v < lo {
				lo = v
			}
			if v > hi {
				hi = v
			}
		}
		for i, v := range h {
			hh := mp(v, lo, hi, 2, 20)
			l.rect(ax+float64(i)*w, y+20-hh, math.Max(w, 0.7), hh, col, col, 0.85)
		}
		l.txt(ax-10, y+15, 11, col, "end", rot)
	}
	tira(424, r.HArit, pGreen, "h aritmético")
	mez := make([]float64, len(r.HArit))
	copy(mez, r.HArit)
	dd := &dado{s: 20260817}
	for i := len(mez) - 1; i > 0; i-- {
		j := int(dd.u() * float64(i+1))
		mez[i], mez[j] = mez[j], mez[i]
	}
	tira(450, mez, pRose, "el mismo, revuelto")
	l.txt(36+16, 486, 11.5, pInk, "start",
		"Σ²(10): más chico es mejor · la línea punteada es el mezclado · mismo histograma arriba y abajo: sólo cambia el ORDEN de las durezas.")

	// ---- B: the control that decides ---------------------------------------
	l.panel(752, 96, 612, 400, pBlue, "B · EL CONTROL QUE DECIDE — ¿es la aritmética, o sólo el ORDEN?")
	bx, by := 780.0, 150.0
	l.txt(bx, by, 13, pInk, "start", "las dos reglas aritméticas dieron casi lo mismo, así que")
	l.txt(bx, by+20, 13, pInk, "start", "hubo que probar campos igual de ordenados y SIN aritmética:")
	l.mono(bx, by+52, 12.5, pGreen, "start", fmt.Sprintf("aritmético  %s", nm(r.Arit.s10, 3)))
	l.mono(bx, by+74, 12.5, pRose, "start", fmt.Sprintf("rampa lisa  %s  ← MUCHO peor: el orden monótono ARRUINA", nm(r.Rampa.s10, 3)))
	l.mono(bx, by+96, 12.5, pBlue, "start", fmt.Sprintf("onda lisa   %s  ← empata con el aritmético", nm(r.Onda.s10, 3)))
	l.mono(bx, by+118, 12.5, pDim, "start", fmt.Sprintf("mezclado    %s ± %s", nm(r.MezMedia, 3), nm(r.MezDes, 3)))
	l.txt(bx, by+156, 13, pGold, "start", "LO QUE ESO DICE, con todas las letras:")
	l.txt(bx, by+180, 12.5, pInk, "start", "· NO es «cualquier orden»: la rampa monótona es peor que el azar.")
	l.txt(bx, by+202, 12.5, pInk, "start", "· El campo aritmético SÍ le gana al mezclado, por 2,28 sigmas.")
	l.txt(bx, by+224, 12.5, pRose, "start", "· PERO una onda lisa, sin nada de aritmética, lo empata.")
	l.rect(bx, by+248, 552, 96, pPanel, pRose, 0.45)
	l.txt(bx+276, by+272, 12.5, pRose, "middle", "⟹ lo que el medio premia es que la dureza tenga estructura")
	l.txt(bx+276, by+292, 12.5, pRose, "middle", "a escala intermedia — ni plana, ni monótona, ni revuelta.")
	l.txt(bx+276, by+314, 12.5, pGold, "middle", "La aritmética ESTÁ en esa clase, pero no se distingue")
	l.txt(bx+276, by+332, 12.5, pGold, "middle", "de una onda cualquiera. No alcanza para atribuírsela.")

	// ---- C: localisation ---------------------------------------------------
	l.panel(36, 512, 700, 300, pRose, "C · LOCALIZACIÓN — y la advertencia de su §9, cumplida")
	cx, cy, cw, ch := 120.0, 566.0, 560.0, 128.0
	l.line(cx, cy+ch, cx+cw, cy+ch, pGrid, 1.4, "")
	l.line(cx, cy, cx, cy+ch, pGrid, 1.4, "")
	var pr, s2 [][2]float64
	for i, lo := range r.Locs {
		x := cx + float64(i)*cw/float64(len(r.Locs)-1)
		pr = append(pr, [2]float64{x, mp(lo[5], 0, 0.12, cy+ch, cy)})
		s2 = append(s2, [2]float64{x, mp(lo[2], 0, 45, cy+ch, cy)})
		l.mono(x, cy+ch+16, 10.5, pDim, "middle", fmt.Sprintf("×%.0f", lo[0]))
	}
	l.camino(pr, pBlue, 2.6, "")
	l.camino(s2, pGold, 2.6, "")
	for _, p := range pr {
		l.circ(p[0], p[1], 3.8, pBlue)
	}
	for _, p := range s2 {
		l.circ(p[0], p[1], 3.8, pGold)
	}
	l.txt(cx+cw/2, cy+ch+34, 11, pDim, "middle", "contraste entre la zona blanda y la dura")
	l.txt(cx+8, cy-8, 11.5, pBlue, "start", "azul: PR/N (1 = extendido, 0 = atrapado)   ·   dorado: Σ²(10)")
	l.txt(56, 758, 12.5, pInk, "start",
		fmt.Sprintf("la porosidad SUAVE (×2) es el mejor punto sin nodo de toda la hoja: Σ²(10) = %s contra %s del homogéneo.",
			nm(r.Locs[1][2], 3), nm(r.Base.s10, 3)))
	l.txt(56, 780, 12.5, pRose, "start",
		fmt.Sprintf("pero al subir el contraste los estados se ATRAPAN (PR/N cae a %s) y la varianza empeora: demasiada porosidad tapona.",
			nm(r.Locs[len(r.Locs)-1][5], 3)))
	l.txt(56, 800, 12, pGold, "start", "hay un punto justo: ni medio plano, ni medio taponado.")

	// ---- D: the node's shadow ----------------------------------------------
	l.panel(752, 512, 612, 300, pGreen, "D · LA SOMBRA DEL NODO — antes de festejar su número")
	dx, dy := 780.0, 566.0
	l.mono(dx, dy, 13, pGold, "start", fmt.Sprintf("Σ²(10):  sin nodo %s  →  con nodo %s", nm(r.Base.s10, 3), nm(r.BaseN.s10, 3)))
	l.mono(dx, dy+26, 13, pRose, "start", fmt.Sprintf("PR/N :   sin nodo %s  →  con nodo %s", nm(r.Base.pr, 3), nm(r.BaseN.pr, 3)))
	l.txt(dx, dy+58, 12.5, pInk, "start", "el nodo baja la varianza a la mitad, sí — pero al mismo")
	l.txt(dx, dy+78, 12.5, pInk, "start", "tiempo ATRAPA los estados: la participación cae 2,5 veces.")
	l.txt(dx, dy+106, 12.5, pRose, "start", "Buena parte de esa mejora puede ser ALCANCE APARENTE por")
	l.txt(dx, dy+126, 12.5, pRose, "start", "localización, no alcance real. Es la distinción exacta que")
	l.txt(dx, dy+146, 12.5, pRose, "start", "pedía su §9, y hay que declararla ANTES de festejar.")
	l.rect(dx, dy+172, 552, 78, pPanel, pGold, 0.4)
	l.txt(dx+276, dy+196, 12.5, pGold, "middle", "la banda NO se vacía en ningún brazo:")
	l.txt(dx+276, dy+216, 12.5, pGold, "middle", "entre 187 y 311 niveles vivos de 400.")
	l.txt(dx+276, dy+238, 12, pGreen, "middle", "la disciplina que dejó la Fase VII, aplicada de entrada.")

	l.raw(fmt.Sprintf(`<rect x="36" y="826" width="1328" height="42" rx="9" fill="%s" stroke="%s" opacity="0.5"/>`, pPanel, pGold))
	l.mono(700, 843, 12.5, pGold, "middle", "C(i,j) = √(h_i·h_j) · f(|i−j|)      h = permeabilidad local      h_i = 1/log(p_i), declarada antes de calcular un solo espectro")
	l.txt(700, 861, 11.5, pInk, "middle",
		fmt.Sprintf("R6 limpio · %d modos · fuerza total fija en %.0f para todos los brazos · los ceros sólo como regla: Σ²(10) = 0,3364", r.N, r.Fuerza))
	l.txt(700, 894, 13, pGold, "middle",
		"el flash de la porosidad es de Jesús Nicolás Astorga · los tres controles obligatorios, de la auditora · las mediciones, de este taller")
	l.txt(700, 922, 13, pGold, "middle", "go run ./cmd/porosidad · estructura cerrada no es hipótesis demostrada · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "porosidad.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
	_ = math.Pi
}
