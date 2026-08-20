package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

// lamina.go - the Phase IX plate. The nulls are drawn as large as the one real
// gain, and the retroactive explanation of the node gets its own box, because it
// is the finding that changes two previous phases.
//
// Every line is placed by the flujo helper rather than by hand, so a title can
// never collide with its first body line again.

const (
	qBg    = "#0b1220"
	qPanel = "#101b30"
	qInk   = "#dce8f7"
	qDim   = "#8fb4d9"
	qGold  = "#ffd98a"
	qGreen = "#7ee0c0"
	qBlue  = "#7fb2ff"
	qRose  = "#ff9aa8"
	qGrid  = "#1d3a63"

	lienzoW = 1460.0
	lienzoH = 1120.0
)

type tela struct{ b strings.Builder }

func (l *tela) raw(s string) { l.b.WriteString(s + "\n") }

func xesc(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;").Replace(s)
}

func num(v float64, d int) string {
	return strings.Replace(fmt.Sprintf("%.*f", d, v), ".", ",", 1)
}

func (l *tela) txt(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, xesc(t)))
}

func (l *tela) mono(x, y, s float64, c, a, t string) {
	l.raw(fmt.Sprintf(`<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Consolas,monospace" fill="%s" text-anchor="%s">%s</text>`, x, y, s, c, a, xesc(t)))
}

func (l *tela) rect(x, y, w, h float64, f, s string, op float64) {
	l.raw(fmt.Sprintf(`<rect x="%.2f" y="%.2f" width="%.2f" height="%.2f" fill="%s" stroke="%s" opacity="%.2f"/>`, x, y, w, h, f, s, op))
}

func (l *tela) linea(x1, y1, x2, y2 float64, c string, w float64, d string) {
	dd := ""
	if d != "" {
		dd = fmt.Sprintf(` stroke-dasharray="%s"`, d)
	}
	l.raw(fmt.Sprintf(`<line x1="%.2f" y1="%.2f" x2="%.2f" y2="%.2f" stroke="%s" stroke-width="%.2f"%s/>`, x1, y1, x2, y2, c, w, dd))
}

func (l *tela) punto(x, y, r float64, c string) {
	l.raw(fmt.Sprintf(`<circle cx="%.2f" cy="%.2f" r="%.2f" fill="%s"/>`, x, y, r, c))
}

func (l *tela) via(p [][2]float64, c string, w float64) {
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
	l.raw(fmt.Sprintf(`<path d="%s" fill="none" stroke="%s" stroke-width="%.2f"/>`, sb.String(), c, w))
}

func lin(v, v0, v1, q0, q1 float64) float64 {
	if v1 == v0 {
		return q0
	}
	return q0 + (v-v0)/(v1-v0)*(q1-q0)
}

// flujo lays out one panel's lines top to bottom. Nothing is hand-placed, so a
// collision would need the pitch itself to be wrong.
type flujo struct {
	l          *tela
	x, y, r, b float64 // cursor, panel right edge, panel bottom
}

func (l *tela) panel(x, y, w, h float64, c, t string) *flujo {
	l.raw(fmt.Sprintf(`<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="10" fill="%s" stroke="%s"/>`, x, y, w, h, qPanel, c))
	l.txt(x+15, y+25, 14.5, c, "start", t)
	return &flujo{l: l, x: x + 16, y: y + 30, r: x + w, b: y + h}
}

func (f *flujo) ln(s float64, c, t string) {
	f.y += s + 6
	f.l.txt(f.x, f.y, s, c, "start", t)
}

func (f *flujo) mn(s float64, c, t string) {
	f.y += s + 6
	f.l.mono(f.x, f.y, s, c, "start", t)
}

func (f *flujo) hueco(d float64) { f.y += d }

// nota is a highlighted closing box: centred lines inside a tinted rectangle.
func (f *flujo) nota(c string, s float64, lineas []string) {
	h := float64(len(lineas))*(s+7) + 16
	f.y += 10
	w := f.r - f.x - 16
	f.l.rect(f.x, f.y, w, h, qPanel, c, 0.5)
	y := f.y + 8
	for _, t := range lineas {
		y += s + 7
		f.l.txt(f.x+w/2, y, s, c, "middle", t)
	}
	f.y += h
}

func dibujar9(r Res9) {
	l := &tela{}
	W, H := lienzoW, lienzoH
	l.raw(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`, W, H, W, H))
	l.raw(fmt.Sprintf(`<rect width="%.0f" height="%.0f" fill="%s"/>`, W, H, qBg))
	l.raw(fmt.Sprintf(`<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="%s" stroke-width="2" opacity="0.5"/>`, W-28, H-28, qGold))
	l.txt(W/2, 50, 24, qInk, "middle", "🔧🧵 LA MECÁNICA DEL MEDIO — once palabras, un solo dial")
	l.txt(W/2, 76, 12.5, qGold, "middle",
		"«¿y si no sólo hay porosidad, sino resistencia, tensión, compresión, cizallamiento, torsión, elasticidad, plasticidad,")
	l.txt(W/2, 94, 12.5, qGold, "middle",
		"dureza, fragilidad, tenacidad y fatiga?» — Jesús Nicolás Astorga")

	// ================= A · the collapse =====================================
	fa := l.panel(32, 112, 452, 432, qGold, "A · EL COLAPSO — 11 palabras → 4 objetos")
	fa.ln(12, qInk, "La auditora prohibió once perillas. No hubo once:")
	fa.mn(12.5, qGold, "U(x) = (s_c·b/π)·(1 − cos(πx/b))")
	fa.ln(11.5, qDim, "UNA función escalar, leída de cuatro maneras.")

	// the potential with its four readings marked
	uy, uh, uw := fa.y+16, 112.0, 386.0
	ux := fa.x + 4
	l.linea(ux, uy+uh, ux+uw, uy+uh, qGrid, 1.2, "")
	var pot [][2]float64
	for k := 0; k <= 240; k++ {
		xx := -2.0 + 4.0*float64(k)/240
		yy := 1 - math.Cos(math.Pi*xx)
		pot = append(pot, [2]float64{ux + lin(xx, -2, 2, 0, uw), uy + uh - lin(yy, 0, 2, 4, uh-6)})
	}
	l.via(pot, qGold, 2.4)
	for _, w := range []float64{-2, 0, 2} {
		l.punto(ux+lin(w, -2, 2, 0, uw), uy+uh-4, 4.2, qGreen)
	}
	for _, br := range []float64{-1, 1} {
		bx := ux + lin(br, -2, 2, 0, uw)
		l.punto(bx, uy+6, 4.2, qRose)
		l.linea(bx, uy+6, bx, uy+uh, qRose, 1, "3 3")
	}
	fa.y = uy + uh
	fa.ln(10.5, qDim, "x en unidades de b · los pozos verdes son los estados plásticos")
	fa.hueco(4)
	fa.mn(11, qGreen, "curvatura del fondo  →  DUREZA")
	fa.mn(11, qRose, "el máximo, s_c       →  RESISTENCIA")
	fa.mn(11, qGold, "el área de la barrera →  TENACIDAD")
	fa.mn(11, qRose, "el salto al cruzarlo  →  FRAGILIDAD")
	fa.mn(11, qGreen, "en qué pozo queda    →  ELÁSTICO / PLÁSTICO")
	fa.mn(11, qBlue, "el signo de x        →  TRACCIÓN / COMPRESIÓN")
	fa.hueco(4)
	fa.ln(11.5, qGold, "Estado: UN real por sitio, igual que Fase VIII.")
	fa.ln(11.5, qGold, fmt.Sprintf("s_c y η son DERIVADOS: %s y %.1e.", num(r.Sc, 4), r.Eta))

	// ================= B · the mechanics measured ===========================
	fb := l.panel(500, 112, 452, 432, qBlue, "B · LA MECÁNICA MEDIDA — los nulos a la vista")
	rcT := r.RhoC["homogéneo · eps=+1 · tracción"]
	rcC := r.RhoC["homogéneo · eps=+1 · compresión"]
	fb.ln(12, qGreen, "✓ ELASTICIDAD: residuo 8e−25 · reversible")
	fb.ln(12, qGreen, fmt.Sprintf("✓ recuperación DERIVADA ≈ %s pasos — su §5", num(r.TRecup, 1)))
	fb.ln(12, qGreen, fmt.Sprintf("✓ RESISTENCIA medida: rho_c = %s (tracción)", num(rcT, 4)))
	fb.ln(12, qGold, fmt.Sprintf("   compresión %s ⟹ asimetría %s, que NO", num(rcC, 4), num((rcT-rcC)/(rcT+rcC), 3)))
	fb.ln(12, qGold, "   se puso: sale de h = h⁰·exp(x). Pero el SIGNO")
	fb.ln(12, qRose, "   salió al REVÉS de la predicción pre-registrada.")
	fb.ln(12, qGreen, fmt.Sprintf("✓ PLASTICIDAD: %d sitio en un pozo entero", r.PlastCed))
	fb.hueco(4)
	fb.ln(12, qRose, "✗ FATIGA: NULO. El mismo impulso en 1 golpe o en")
	fb.ln(12, qInk, "   32 difiere 2,7e−25 · y al doblar el reposo cae")
	fb.ln(12, qInk, "   a 9e−48: es relajación, no memoria. Era la ÚNICA")
	fb.ln(12, qInk, "   palabra dejada como predicción falsable en vez")
	fb.ln(12, qInk, "   de construirla — por eso el nulo vale algo.")
	fb.ln(12, qRose, fmt.Sprintf("✗ GRIETAS: avalancha %v — cada golpe tumba UN", r.Aval))
	fb.ln(12, qInk, "   sitio y nada se propaga. Fallas aisladas.")
	fb.nota(qRose, 12, []string{
		"✗ Y EL ESPECTRO NO SE MUEVE",
		fmt.Sprintf("memoria sola %s contra %s del virgen", num(r.Medios["C · homogéneo dinámico"].s10, 3), num(r.Ident.s10, 3)),
		"La mecánica es real y espectralmente MUDA.",
	})

	// ================= C · the sign channel =================================
	fc := l.panel(968, 112, 460, 432, qGreen, "C · EL CANAL DE SIGNOS — la única ganancia real")
	fc.ln(11.5, qInk, "Fase VIII tenía el núcleo TODO positivo: no había")
	fc.ln(11.5, qInk, "tensión. Un signo de ENLACE la da a costo CERO, y")
	fc.ln(11.5, qInk, "como signo² = 1 la fuerza total queda intacta.")
	bars := []struct {
		n string
		v float64
		c string
	}{
		{"S0 · todos +1", r.Ident.s10, qDim},
		{"S1 · gauge (test)", r.Signos["S1 gauge"].s10, qDim},
		{"S2 · azar frustrado", r.SignoS2Media, qBlue},
		{"S3 · reciprocidad", r.Signos["S3 reciprocidad"].s10, qGreen},
		{"nodo solo", r.SignosNodo["S0 llano"].s10, qGold},
		{"nodo + azar", r.S2NodoMedia, qGold},
		{"nodo + reciprocidad", r.SignosNodo["S3 reciprocidad"].s10, qGold},
	}
	bx3, bh := 1148.0, 25.0
	fc.hueco(10)
	for _, b := range bars {
		w := lin(b.v, 0, 19, 0, 216)
		l.rect(bx3, fc.y, w, bh-8, b.c, b.c, 0.6)
		l.txt(bx3-8, fc.y+bh-12, 11, qInk, "end", b.n)
		l.mono(bx3+w+6, fc.y+bh-12, 11, b.c, "start", num(b.v, 3))
		fc.y += bh
	}
	fc.ln(11.5, qBlue, "S1 devuelve S0 EXACTO (Δ = 0): el teorema de gauge,")
	fc.ln(11.5, qBlue, "verificado bit a bit. Queda como test unitario.")
	fc.nota(qGold, 11.5, []string{
		"⚡ EL NODO ESTABA HACIENDO FRUSTRACIÓN",
		fmt.Sprintf("nodo solo %s  ·  azar SIN nodo %s ± %s",
			num(r.SignosNodo["S0 llano"].s10, 3), num(r.SignoS2Media, 3), num(r.SignoS2Des, 3)),
		"a 0,34 σ: toda la mejora del reloj de arena sale con",
		"signos al azar, sin nodo y sin un solo primo adentro.",
	})

	// ================= D · the arithmetic blade =============================
	fd := l.panel(32, 560, 690, 254, qRose, "D · LA HOJA ARITMÉTICA, TERCERA APLICACIÓN — y da lo mismo")
	fd.ln(12.5, qInk, "La reciprocidad (p ≡ 3 mod 4) contra signos AL AZAR a la MISMA densidad medida:")
	fd.mn(12, qGreen, fmt.Sprintf("sin nodo   reciprocidad %s   azar %s ± %s   →  0,45 σ",
		num(r.Signos["S3 reciprocidad"].s10, 3), num(r.SignoS2Media, 3), num(r.SignoS2Des, 3)))
	fd.mn(12, qGold, fmt.Sprintf("con nodo   reciprocidad %s   azar %s ± %s   →  1,17 σ",
		num(r.SignosNodo["S3 reciprocidad"].s10, 3), num(r.S2NodoMedia, 3), num(r.S2NodoDes, 3)))
	fd.hueco(4)
	fd.ln(12.5, qRose, "⟹ La aritmética NO se separa del azar. Lo que mueve el espectro es la")
	fd.ln(12.5, qRose, "   FRUSTRACIÓN: que la mitad de los triángulos dé producto de signos negativo.")
	fd.ln(11.5, qInk, fmt.Sprintf("   medida: gauge %s · azar %s · reciprocidad %s — iguales",
		num(r.FrustS1, 3), num(r.FrustS2, 3), num(r.FrustS3, 3)))
	fd.hueco(4)
	fd.ln(12, qGold, "Es la MISMA hoja que mató la atribución aritmética en Fase VIII (una onda lisa a")
	fd.ln(12, qGold, "0,81 σ) y la misma respuesta. Tres mecanismos distintos, tres veces lo mismo.")
	fd.ln(11.5, qDim, fmt.Sprintf("Y el precio: PR/N cae de %s a %s ⟹ buena parte es alcance APARENTE.",
		num(r.Ident.pr, 3), num(r.Signos["S2 azar frustrado"].pr, 3)))

	// ================= E · the ceiling ======================================
	fe := l.panel(738, 560, 690, 254, qGold, "E · EL TECHO — y por qué reencuadra el taller entero")
	fe.hueco(30)
	ex, ew := fe.x+40, 570.0
	ey := fe.y
	l.linea(ex, ey, ex+ew, ey, qGrid, 1.4, "")
	marks := []struct {
		v float64
		n string
		c string
	}{
		{r.Ceros, "ceros", qGold},
		{r.PisoGUE, "piso GUE", qBlue},
		{r.PisoGOE, "piso GOE", qDim},
		{r.SignosNodo["S3 reciprocidad"].s10, "mejor de la fase", qGreen},
		{r.Ident.s10, "virgen", qRose},
	}
	for _, m := range marks {
		x := ex + lin(math.Log10(m.v), math.Log10(0.30), math.Log10(20), 0, ew)
		l.linea(x, ey-9, x, ey+9, m.c, 2.2, "")
		l.punto(x, ey, 4, m.c)
		l.mono(x, ey-16, 10, m.c, "middle", num(m.v, 3))
		l.txt(x, ey+24, 10, m.c, "middle", m.n)
	}
	fe.y = ey + 30
	fe.ln(10.5, qDim, "                                Σ²(10), escala logarítmica")
	fe.hueco(4)
	fe.ln(12.5, qGold, "⟹ EL OBJETIVO ESTÁ DEBAJO DE LOS DOS PISOS DE MATRIZ ALEATORIA.")
	fe.ln(12, qInk, "Los ceros son MÁS rígidos que GUE a esta distancia. Ninguna ingeniería de clase")
	fe.ln(12, qInk, "de simetría puede llegar — y un modelo que se acerca al piso GOE se volvió una")
	fe.ln(12, qInk, "matriz aleatoria sin estructura: eso es PERDER contenido aritmético disfrazado")
	fe.ln(12, qInk, "de ganarlo. El hueco que falta no es un déficit de simetría.")

	// ================= F · the two theorems =================================
	ff := l.panel(32, 830, 1396, 186, qBlue, "F · DOS TEOREMAS QUE NO NECESITARON CORRER NADA")
	ff.ln(12.5, qGreen, "1 · UNA FASE DE SITIO ES GAUGE PURA.  H' = D H D† con D = diag(e^{−iθ}) ⟹ mismo espectro y misma PR/N, exactamente.")
	ff.ln(12, qInk, "    ⟹ Resuelve el PENDIENTE 3 de la Fase VIII —«¿entra la aritmética como FASE?»— en NEGATIVO para la fase de sitio: no")
	ff.ln(12, qInk, "    mueve ni un observable. No hacía falta correr, hacía falta mirar. Se verificó igual, y dio Δ = 0 exacto.")
	ff.hueco(6)
	ff.ln(12.5, qGold, "2 · CORRECCIÓN A SU §8: LA TORSIÓN NO NECESITA UNA RED 2D.  Un signo o fase de ENLACE de la forma u_i − u_j es un coborde,")
	ff.ln(12, qInk, "    y tampoco mueve nada. Lo que falta NO es dimensión: con kmax = 120 el grafo ya está lleno de triángulos (todo i<j<k con")
	ff.ln(12, qInk, "    k−i ≤ 120). Falta que la fase de enlace NO sea un coborde. Una cadena de vecinos próximos sería un árbol y ahí sí la")
	ff.ln(12, qInk, "    torsión sería imposible: la cola larga que eligió la Fase VII quitó esa excusa.")

	l.txt(W/2, 1046, 11.5, qDim, "middle",
		fmt.Sprintf("%d modos · fuerza total fija · R6 limpio: la ley se declaró sin mirar un solo cero · Σσ = 0 y la invariancia de escala, verificadas a 1e−13", r.N))
	l.txt(W/2, 1068, 12, qGold, "middle",
		"el flash de las once palabras es de Jesús Nicolás Astorga · los dieciséis controles, de la auditoría · el defecto del brazo uniforme lo cazó nuestro propio teorema")
	l.txt(W/2, 1092, 12.5, qGold, "middle", "go run ./cmd/lamecanica · estructura cerrada no es hipótesis demostrada · Todavía no.")

	l.raw("</svg>")
	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		panic(err)
	}
	full := filepath.Join(dir, "la-mecanica.svg")
	if err := os.WriteFile(full, []byte(l.b.String()), 0o644); err != nil {
		panic(err)
	}
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
