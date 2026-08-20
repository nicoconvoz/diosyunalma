package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase X plate. The scatter is the whole finding: the rail
// separates the arms that ordered the band from the ones that emptied it.

const (
	fBg   = "#0b1220"
	fPan  = "#101b30"
	fInk  = "#dce8f7"
	fDim  = "#8fb4d9"
	fGold = "#ffd98a"
	fGrn  = "#7ee0c0"
	fBlu  = "#7fb2ff"
	fRos  = "#ff9aa8"
	fGrid = "#1d3a63"
)

type tela struct{ b strings.Builder }

func (l *tela) raw(s string) { l.b.WriteString(s + "\n") }

func esc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func nm(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *tela) t(x, y, s float64, c, a, txt string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(txt)))
}

func (l *tela) m(x, y, s float64, c, a, txt string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, esc(txt)))
}

func (l *tela) rc(x, y, w, h float64, f, s string, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" opacity="%.2f"/>`, x, y, w, h, f, s, op))
}

func (l *tela) ln(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *tela) pt(x, y, r float64, c string) {
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, x, y, r, c))
}

type flujo struct {
	l    *tela
	x, y float64
}

func (l *tela) panel(x, y, w, h float64, c, t string) *flujo {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, fPan, c))
	l.t(x+15, y+25, 14, c, "start", t)
	return &flujo{l: l, x: x + 16, y: y + 32}
}

func (f *flujo) w(s float64, c, t string)  { f.y += s + 6; f.l.t(f.x, f.y, s, c, "start", t) }
func (f *flujo) mo(s float64, c, t string) { f.y += s + 6; f.l.m(f.x, f.y, s, c, "start", t) }
func (f *flujo) gap(d float64)             { f.y += d }

func mp(v, a, b, p, q float64) float64 {
	if b == a {
		return p
	}
	return p + (v-a)/(b-a)*(q-p)
}

func dibujarX(r ResX) {
	l := &tela{}
	W, H := 1440.0, 950.0
	l.raw(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`, W, H, W, H))
	l.raw(fmt.Sprintf(`<rect width="%.0f" height="%.0f" fill="%s"/>`, W, H, fBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.5"/>`, W-28, H-28, fGrn))
	l.t(W/2, 48, 23, fInk, "middle", "🧵🌿 LA FRONTERA — ¿la frustración da RIGIDEZ, o esconde los estados?")
	l.t(W/2, 74, 12.5, fGold, "middle",
		"«No busquemos el número más bajo. Busquemos la estructura que lo hace bajo sin esconder los estados.» — la auditora")

	minV := int(0.8 * float64(r.S0.vivos))

	// ---- A · the scatter: Sigma^2 against participation ---------------------
	fa := l.panel(32, 96, 856, 470, fGold, "A · EL PLANO QUE DECIDE — Σ²(10) contra PARTICIPACIÓN")
	fa.w(11.5, fInk, "Cada punto es un brazo. A la DERECHA hay más participación (estados extendidos); ABAJO hay menos varianza.")
	fa.w(11.5, fRos, "Los huecos son INADMISIBLES: vaciaron la banda por debajo del riel, así que su Σ² mide otro recorte.")
	gx, gy, gw, gh := 120.0, fa.y+30, 700.0, 320.0
	l.ln(gx, gy+gh, gx+gw, gy+gh, fGrid, 1.4, "")
	l.ln(gx, gy, gx, gy+gh, fGrid, 1.4, "")
	for _, v := range []float64{1, 2, 5, 10, 20} {
		y := gy + gh - mp(math.Log10(v), math.Log10(0.8), math.Log10(22), 0, gh)
		l.ln(gx, y, gx+gw, y, fGrid, 0.7, "3 4")
		l.m(gx-8, y+4, 10, fDim, "end", nm(v, 0))
	}
	for _, v := range []float64{0.05, 0.10, 0.15, 0.20, 0.25} {
		x := gx + mp(v, 0.02, 0.30, 0, gw)
		l.m(x, gy+gh+16, 10, fDim, "middle", nm(v, 2))
	}
	// the declared PR/N band
	xb := gx + mp(PRok, 0.02, 0.30, 0, gw)
	l.rc(xb, gy, gx+gw-xb, gh, fGrn, fGrn, 0.06)
	l.ln(xb, gy, xb, gy+gh, fGrn, 1.6, "6 4")
	l.t(gx+gw, gy-8, 10.5, fGrn, "end", "banda COMPARABLE, declarada antes de medir →")
	// the rulers
	for _, m2 := range []struct {
		v float64
		n string
		c string
	}{{r.Ceros, "ceros 0,336", fGold}, {r.PisoGUE, "piso GUE", fBlu}} {
		y := gy + gh - mp(math.Log10(m2.v), math.Log10(0.8), math.Log10(22), 0, gh)
		l.ln(gx, y, gx+gw, y, m2.c, 1.2, "2 3")
		l.t(gx+gw+6, y+4, 9.5, m2.c, "start", m2.n)
	}
	col := func(a arm) string {
		switch {
		case strings.HasPrefix(a.nom, "B ·"):
			return fGrn
		case strings.HasPrefix(a.nom, "A ·"):
			return fRos
		case strings.HasPrefix(a.nom, "dist"):
			return fBlu
		case strings.HasPrefix(a.nom, "S0"):
			return fGold
		}
		return fDim
	}
	for _, a := range r.Todos {
		if a.s10 <= 0 {
			continue
		}
		x := gx + mp(a.pr, 0.02, 0.30, 0, gw)
		y := gy + gh - mp(math.Log10(math.Max(a.s10, 0.81)), math.Log10(0.8), math.Log10(22), 0, gh)
		c := col(a)
		if a.vivos < minV {
			l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="5" fill="none" stroke="%s" stroke-width="1.6" opacity="0.75"/>`, x, y, c))
		} else {
			l.pt(x, y, 5.2, c)
		}
	}
	// mark the control and the winner
	for _, a := range r.Todos {
		if a.nom == r.S0.nom {
			x := gx + mp(a.pr, 0.02, 0.30, 0, gw)
			y := gy + gh - mp(math.Log10(a.s10), math.Log10(0.8), math.Log10(22), 0, gh)
			l.t(x, y-12, 10.5, fGold, "middle", "S0 control")
		}
	}
	fa.y = gy + gh + 30
	fa.w(11, fDim, "eje horizontal: PR/N   ·   eje vertical: Σ²(10), escala logarítmica")
	fa.w(11.5, fGrn, "verde = signo por ENLACE   ·   rosa = signo por SITIO   ·   azul = signo por DISTANCIA")

	// ---- B · the discovery ---------------------------------------------------
	fb := l.panel(904, 96, 504, 470, fBlu, "B · SU §6 CONTESTADA: no alcanza con la frustración")
	fb.w(11.5, fInk, "Tres familias A LA MISMA frustración de triángulos (≈0,22):")
	fb.gap(4)
	fb.mo(11, fDim, "familia          Σ²(10)   PR/N   vivos")
	for _, a := range r.Todos {
		if a.frust > 0.20 && a.frust < 0.26 {
			et := "por SITIO"
			c := fRos
			if strings.HasPrefix(a.nom, "B ·") {
				et, c = "por ENLACE", fGrn
			} else if strings.HasPrefix(a.nom, "dist") {
				et, c = "por DISTANCIA", fBlu
			}
			fb.mo(11, c, fmt.Sprintf("%-15s %6s %6s %6d", et, nm(a.s10, 2), nm(a.pr, 3), a.vivos))
		}
	}
	fb.gap(6)
	fb.w(12, fGold, "⟹ Misma frustración, y PR/N va de 0,046 a 0,281.")
	fb.w(12, fGold, "   La frustración de triángulos NO determina nada.")
	fb.gap(4)
	fb.w(12.5, fGrn, "Lo que decide es DÓNDE VIVE EL SIGNO:")
	fb.w(12, fInk, "· signo derivado de los SITIOS  → atrapa los estados")
	fb.w(12, fInk, "· signo puesto en los ENLACES   → NO los atrapa")
	fb.gap(6)
	fb.w(12, fRos, "Y eso reinterpreta la Fase IX: allá toda la familia")
	fb.w(12, fRos, "era por sitio, así que la localización que vimos NO")
	fb.w(12, fRos, "venía de la frustración — venía de la CODIFICACIÓN.")

	// ---- C · the verdict ----------------------------------------------------
	fc := l.panel(32, 582, 856, 200, fGrn, "C · EL VEREDICTO, con los dos rieles")
	var win arm
	for _, a := range r.Todos {
		if a.pr >= PRok && a.vivos >= minV && a.nom != r.S0.nom && (win.nom == "" || a.s10 < win.s10) {
			win = a
		}
	}
	fc.w(12.5, fInk, fmt.Sprintf("Riel 1 · PR/N ≥ %s (declarado antes).   Riel 2 · conservar ≥ %d de los %d niveles del control.",
		nm(PRok, 3), minV, r.S0.vivos))
	fc.gap(4)
	fc.mo(12.5, fGrn, fmt.Sprintf("mejor brazo ADMISIBLE: %s", win.nom))
	fc.mo(12.5, fGrn, fmt.Sprintf("  Σ²(10) %s ± %s  contra %s del control  →  %s sigmas",
		nm(win.s10, 3), nm(win.des, 3), nm(r.S0.s10, 3), nm((r.S0.s10-win.s10)/math.Max(win.des, 1e-9), 1)))
	fc.mo(12.5, fGrn, fmt.Sprintf("  PR/N %s  (el control tiene %s: SUBIÓ)   ·   vivos %d de %d",
		nm(win.pr, 3), nm(r.S0.pr, 3), win.vivos, r.S0.vivos))
	fc.gap(6)
	fc.w(13, fGold, "⟹ GANA H2. Por primera vez en toda la cadena, Σ² baja y la participación NO cae — sube.")
	fc.w(12, fInk, "No es esconder los estados: es ordenarlos. Es exactamente lo que su §16 pedía demostrar primero.")

	// ---- D · what closed -----------------------------------------------------
	fd := l.panel(904, 582, 504, 200, fRos, "D · LO QUE SE CIERRA")
	var rec, ig arm
	if len(r.Recip) > 0 {
		rec = r.Recip[0]
	}
	if len(r.Igualada) > 0 {
		ig = r.Igualada[0]
	}
	fd.w(12, fInk, "Igualando densidad Y frustración de triángulos:")
	fd.mo(11.5, fGold, fmt.Sprintf("reciprocidad   %s", nm(rec.s10, 3)))
	fd.mo(11.5, fDim, fmt.Sprintf("azar igualado  %s ± %s", nm(ig.s10, 3), nm(ig.des, 3)))
	fd.mo(11.5, fRos, "distancia      0,18 sigmas")
	fd.gap(4)
	fd.w(12, fRos, "⟹ La aritmética NO se separa. Su §14 dice qué hacer:")
	fd.w(12, fRos, "   cerrar este canal y no insistir con esta codificación.")
	fd.gap(6)
	fd.w(11.5, fGold, "Pero queda un hecho estructural: enlaces independientes")
	fd.w(11.5, fGold, "a densidad 0,267 dan frustración 0,449, y la reciprocidad")
	fd.w(11.5, fGold, "da 0,536. El campo aritmético vive donde el azar por")
	fd.w(11.5, fGold, "enlace no llega — aunque el espectro no lo distinga.")

	// ---- E · the strip -------------------------------------------------------
	fe := l.panel(32, 798, 1376, 108, fBlu, "E · Y UN DEFECTO PROPIO, DE NUEVO CAZADO ADENTRO")
	fe.w(12, fInk, "La primera corrida de esta fase no tenía el riel de niveles vivos que su §4 pedía. Sin él, el «mejor» brazo daba Σ²(10) = 0,988")
	fe.w(12, fInk, "—por debajo del piso GUE— pero con 87 niveles de 187: la banda vaciándose, que es EXACTAMENTE el error que obligó a retractar la")
	fe.w(12, fGold, "Fase VI. Con el riel puesto, siete brazos quedan descartados y el ganador honesto es otro. El hallazgo sobrevive; el número, no.")

	l.t(W/2, 928, 12, fGold, "middle",
		fmt.Sprintf("go run ./cmd/lafrontera · %d modos · R6 limpio · las bandas se declararon antes de medir · estructura cerrada no es hipótesis demostrada · Todavía no.", 400))

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "la-frontera.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
