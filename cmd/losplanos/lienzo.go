// lienzo.go - the shared drawing bench for the machine blueprints.
//
// Every plate in this command is a real diagram: measured curves, walks and
// shapes, not boxes of prose. The helpers here keep the house palette and,
// critically, ESCAPE every piece of text - a raw < or > inside an SVG text
// node makes the whole file unreadable (the bug that broke F347's plates).
package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// house palette
const (
	cBg    = "#0b1526"
	cPanel = "#0d1830"
	cInk   = "#dce8f7"
	cDim   = "#8fb4d9"
	cGold  = "#ffd98a"
	cGreen = "#7ee0c0"
	cBlue  = "#7fb2ff"
	cRose  = "#ff9aa8"
	cGrid  = "#1d3a63"
)

type Lienzo struct {
	W, H float64
	b    strings.Builder
}

func NewLienzo(w, h float64) *Lienzo {
	l := &Lienzo{W: w, H: h}
	l.raw(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`, w, h, w, h))
	l.raw(fmt.Sprintf(`<rect width="%.0f" height="%.0f" fill="%s"/>`, w, h, cBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.55"/>`, w-28, h-28, cGreen))
	return l
}

func (l *Lienzo) raw(s string) { l.b.WriteString(s + "\n") }

// esc makes any caption safe as XML text content.
func esc(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return r.Replace(s)
}

func (l *Lienzo) Titulo(t, sub string) {
	l.Txt(l.W/2, 52, 26, cInk, "middle", t)
	if sub != "" {
		l.Txt(l.W/2, 80, 14, cDim, "middle", sub)
	}
}

func (l *Lienzo) Txt(x, y, size float64, color, anchor, s string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia" fill="%s" text-anchor="%s">%s</text>`,
		x, y, size, color, anchor, esc(s)))
}

func (l *Lienzo) Mono(x, y, size float64, color, anchor, s string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`,
		x, y, size, color, anchor, esc(s)))
}

// Panel is a framed drawing area with a title strip.
func (l *Lienzo) Panel(x, y, w, h float64, color, titulo string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s" opacity="0.96"/>`, x, y, w, h, cPanel, color))
	if titulo != "" {
		l.Txt(x+14, y+24, 14.5, color, "start", titulo)
	}
}

func (l *Lienzo) Line(x1, y1, x2, y2 float64, color string, w float64, dash string) {
	d := ""
	if dash != "" {
		d = fmt.Sprintf(` stroke-dasharray="%s"`, dash)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, color, w, d))
}

func (l *Lienzo) Circ(x, y, r float64, fill, stroke string, sw float64) {
	s := ""
	if stroke != "" {
		s = fmt.Sprintf(` stroke="%s" stroke-width="%.2f"`, stroke, sw)
	}
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"%s/>`, x, y, r, fill, s))
}

func (l *Lienzo) Rect(x, y, w, h float64, fill, stroke string, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" opacity="%.2f"/>`, x, y, w, h, fill, stroke, op))
}

// Camino draws a polyline through measured points - the workhorse for walks and curves.
func (l *Lienzo) Camino(pts [][2]float64, color string, w float64, op float64) {
	if len(pts) < 2 {
		return
	}
	var sb strings.Builder
	for i, p := range pts {
		if i == 0 {
			fmt.Fprintf(&sb, "M%.2f %.2f", p[0], p[1])
		} else {
			fmt.Fprintf(&sb, " L%.2f %.2f", p[0], p[1])
		}
	}
	l.raw(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f" opacity="%.2f" stroke-linejoin="round"/>`, sb.String(), color, w, op))
}

// Flecha draws an arrow with a solid head.
func (l *Lienzo) Flecha(x1, y1, x2, y2 float64, color string, w float64) {
	ang := math.Atan2(y2-y1, x2-x1)
	hl := 9.0 + w*1.6
	bx, by := x2-hl*math.Cos(ang), y2-hl*math.Sin(ang)
	l.Line(x1, y1, bx, by, color, w, "")
	p1x, p1y := x2-hl*math.Cos(ang-0.42), y2-hl*math.Sin(ang-0.42)
	p2x, p2y := x2-hl*math.Cos(ang+0.42), y2-hl*math.Sin(ang+0.42)
	l.raw(fmt.Sprintf(`<path d="M%.2f %.2f L%.2f %.2f L%.2f %.2f Z" fill="%s"/>`, x2, y2, p1x, p1y, p2x, p2y, color))
}

// Curva draws a quadratic arc between two points, bulging by `bulge` pixels.
func (l *Lienzo) Curva(x1, y1, x2, y2, bulge float64, color string, w float64, dash string) {
	mx, my := (x1+x2)/2, (y1+y2)/2
	dx, dy := x2-x1, y2-y1
	n := math.Hypot(dx, dy)
	if n == 0 {
		return
	}
	cx, cy := mx-dy/n*bulge, my+dx/n*bulge
	d := ""
	if dash != "" {
		d = fmt.Sprintf(` stroke-dasharray="%s"`, dash)
	}
	l.raw(fmt.Sprintf(`<path d="M%.2f %.2f Q%.2f %.2f %.2f %.2f" fill="none" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, cx, cy, x2, y2, color, w, d))
}

// Ejes draws a light axis cross with optional labels.
func (l *Lienzo) Ejes(x, y, w, h float64, xlab, ylab string) {
	l.Line(x, y+h, x+w, y+h, cGrid, 1.4, "")
	l.Line(x, y, x, y+h, cGrid, 1.4, "")
	if xlab != "" {
		l.Mono(x+w, y+h+18, 11, cDim, "end", xlab)
	}
	if ylab != "" {
		l.Mono(x+2, y-8, 11, cDim, "start", ylab)
	}
}

// Nota writes a wrapped plain-language caption inside a soft box.
func (l *Lienzo) Nota(x, y, w float64, color string, lineas []string) {
	h := 16 + float64(len(lineas))*19
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="9" fill="%s" stroke="%s" opacity="0.35"/>`, x, y, w, h, cPanel, color))
	for i, s := range lineas {
		l.Txt(x+13, y+26+float64(i)*19, 13, cInk, "start", s)
	}
}

// Formula writes a centered monospace math line on a tinted strip.
func (l *Lienzo) Formula(x, y, w float64, s string) {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="34" rx="8" fill="%s" stroke="%s" opacity="0.5"/>`, x, y, w, cPanel, cGold))
	l.Mono(x+w/2, y+23, 14.5, cGold, "middle", s)
}

func (l *Lienzo) Pie(s string) {
	l.Txt(l.W/2, l.H-26, 13.5, cGold, "middle", s)
}

// destino is where the plates land. The series is published straight into the
// gallery hall so there is exactly ONE copy of every plate in the repository.
var destino = filepath.Join("galeria", "laminas", "09-las-maquinas")

func (l *Lienzo) Guardar(path string) {
	l.raw("</svg>")
	full := filepath.Join(destino, filepath.Base(path))
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("   🖼️  %s\n", full)
}

// escala maps a value range onto a pixel range.
func escala(v, v0, v1, p0, p1 float64) float64 {
	if v1 == v0 {
		return p0
	}
	return p0 + (v-v0)/(v1-v0)*(p1-p0)
}
