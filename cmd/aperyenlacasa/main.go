package main

// APERY EN LA CASA - the captain: "put Apery's constant into our assembly
// formula and into the shapeshifter, in the 1/2 relation of dimension 0. I
// want to see what happens."
//
// The two formulas are the house's own:
//   the shapeshifter   w(s) = 1 - 1/s          (F276: w(2) = 1/2 exactly)
//   the master curve   R(n) = n^(-sigma) / (2|sin(t/2n)|)   (Theorem 7)
//
// And zeta(3) = 1.2020569... arrived from somewhere else entirely: it counts
// the photons of the cosmic background (F411).
//
// THE DANGER, DECLARED BEFORE ANY NUMBER IS COMPUTED: dropping a constant into
// a formula that did not derive it is EXACTLY how numerology works. Any number
// put through any formula produces some other number, and the human eye will
// find meaning in it. So this experiment is built around its control:
//
//   whatever zeta(3) produces, the SAME insertion is run with e, pi, sqrt2,
//   the golden ratio, gamma, and five random numbers of the same size. If
//   zeta(3) does not stand out from that crowd, the answer is NOTHING HAPPENS,
//   and that is the honest result.
//
// Verdicts fixed BEFORE running:
//   V1 (trap #25) the insertion is numerology unless a control decides. State
//      it first.
//   V2 compute w(zeta3) exactly and locate it on the map: the shapeshifter
//      sends the critical line to the unit circle and s=1 (zeta's pole) to 0,
//      so |w| measures how deep inside a number sits.
//   V3 run the identical computation on nine control constants.
//   V4 put zeta(3) into the master curve's waist equation x cot x = sigma and
//      into R(n), again against the same controls.
//   V5 report whether zeta(3) stood out. Predicted before running: IT WILL
//      NOT, because nothing in either formula derives zeta(3). If the
//      prediction fails, that is the interesting case.
// Apery 1979 (the irrationality). Nothing here is new.

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// w is the house's shapeshifter from F276.
func w(s float64) float64 { return 1 - 1/s }

// cintura solves x cot x = sigma for x in (0, pi), the master curve's waist.
func cintura(sigma float64) float64 {
	lo, hi := 1e-9, math.Pi-1e-9
	for i := 0; i < 200; i++ {
		m := (lo + hi) / 2
		if m*math.Cos(m)/math.Sin(m) > sigma {
			lo = m
		} else {
			hi = m
		}
	}
	return (lo + hi) / 2
}

type cte struct {
	nombre string
	v      float64
	ctrl   bool
}

type res struct {
	c  cte
	wv float64
}

type dist struct {
	nombre string
	d      float64
	cual   string
	ctrl   bool
}

func main() {
	z3 := 0.0
	for n := 1; n <= 40000000; n++ {
		f := float64(n)
		z3 += 1 / (f * f * f)
	}
	fmt.Println("🔬 APÉRY EN LA CASA — meter ζ(3) en nuestras dos fórmulas, con control")
	fmt.Println()
	fmt.Printf("   ζ(3) sumada acá hasta 4×10⁷: %.10f   (Apéry, 1979: irracional)\n\n", z3)

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · EL PELIGRO, declarado antes de calcular nada (trampa vigesimoquinta)")
	fmt.Println()
	fmt.Println("   Meter una constante en una fórmula que NO la derivó es exactamente cómo funciona la")
	fmt.Println("   numerología: cualquier número puesto en cualquier fórmula da otro número, y el ojo humano")
	fmt.Println("   le encuentra sentido. Así que todo lo que salga se corre TAMBIÉN con nueve constantes de")
	fmt.Println("   control del mismo tamaño. Si ζ(3) no se destaca del montón, la respuesta es «no pasa nada».")
	fmt.Println()
	fmt.Println("   ⚖️ Y la predicción, escrita ANTES de correr: NO se va a destacar, porque ni el cambiaformas")
	fmt.Println("      ni la curva madre derivan a ζ(3) — viene de contar fotones, que es otro mundo. Si la")
	fmt.Println("      predicción falla, ÉSE sería el caso interesante.")
	fmt.Println()

	r := rand.New(rand.NewSource(20260919))
	ctes := []cte{
		{"ζ(3)  Apéry", z3, false},
		{"φ−1/2  áureo", (1+math.Sqrt(5))/2 - 0.5, true},
		{"√2", math.Sqrt2, true},
		{"π/e", math.Pi / math.E, true},
		{"e/2 + 1/4", math.E/2 + 0.25, true},
		{"γ+2/3", 0.5772156649015329 + 2.0/3, true},
		{"ln(10)/2", math.Log(10) / 2, true},
	}
	for i := 0; i < 3; i++ {
		ctes = append(ctes, cte{fmt.Sprintf("azar %d", i+1), 1 + 0.4*r.Float64(), true})
	}

	// ---------------------------------------------------------------- V2 + V3
	fmt.Println("§2 · POR EL CAMBIAFORMAS: w(s) = 1 − 1/s")
	fmt.Println()
	fmt.Println("   Recordemos el mapa (F276): w manda la línea crítica al círculo de radio 1, y manda el 1")
	fmt.Println("   —el polo de zeta— al CENTRO. Así que |w| mide qué tan adentro está un número.")
	fmt.Println()
	fmt.Println("        el 1 (polo)  →  0,000000     ·     el 2  →  0,500000 EXACTO     ·     el ½  →  −1 (la piel)")
	fmt.Println()
	fmt.Println("        constante              su valor        w(s) = 1 − 1/s       |w|        ¿algo notable?")
	var rs []res
	for _, c := range ctes {
		wv := w(c.v)
		rs = append(rs, res{c, wv})
		nota := ""
		if !c.ctrl {
			nota = "  ← la que pidió el capitán"
		}
		fmt.Printf("        %-18s   %10.7f      %10.7f      %8.6f%s\n", c.nombre, c.v, wv, math.Abs(wv), nota)
	}
	fmt.Println()
	// is zeta3's w special? measure distance to "nice" numbers
	fmt.Println("   ¿El w de cada una cae cerca de algún número redondo? Se mide la distancia a ½, ⅓, ¼, ⅕, 1/6:")
	fmt.Println()
	fmt.Println("        constante              w(s)         el redondo más cercano      qué tan lejos")
	redondos := []struct {
		n string
		v float64
	}{{"1/2", 0.5}, {"1/3", 1.0 / 3}, {"1/4", 0.25}, {"1/5", 0.2}, {"1/6", 1.0 / 6}, {"1/7", 1.0 / 7}, {"0", 0}}
	var ds []dist
	for _, x := range rs {
		mejor, cual := math.Inf(1), ""
		for _, rr := range redondos {
			if d := math.Abs(x.wv - rr.v); d < mejor {
				mejor, cual = d, rr.n
			}
		}
		ds = append(ds, dist{x.c.nombre, mejor, cual, x.c.ctrl})
		nota := ""
		if !x.c.ctrl {
			nota = "  ←"
		}
		fmt.Printf("        %-18s   %10.7f        %-6s                  %.6f%s\n", x.c.nombre, x.wv, cual, mejor, nota)
	}
	sort.Slice(ds, func(i, j int) bool { return ds[i].d < ds[j].d })
	fmt.Println()
	fmt.Println("        ranking de cercanía (el que MÁS cerca esté de un redondo, primero):")
	for i, d := range ds {
		et := "control"
		if !d.ctrl {
			et = "★ ζ(3)"
		}
		fmt.Printf("          %d°  %-18s  a %.6f de %-4s   (%s)\n", i+1, d.nombre, d.d, d.cual, et)
	}
	pos := 0
	for i, d := range ds {
		if !d.ctrl {
			pos = i + 1
		}
	}
	fmt.Printf("\n   ⚖️ ζ(3) quedó en el puesto %d de %d. ", pos, len(ds))
	if pos <= 2 {
		fmt.Println("Está arriba — habría que mirarlo más de cerca.")
	} else {
		fmt.Println("En el montón: NO se destaca de los controles.")
	}
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§3 · EN LA FÓRMULA DE ENSAMBLE: la cintura de la Curva Madre, x·cot x = σ")
	fmt.Println()
	fmt.Println("   La curva madre tiene una cintura única para cada σ, que sale de resolver x·cot x = σ.")
	fmt.Println("   Con σ = ½ (la línea crítica) da el valor que el Teorema 7 verificó al 1 %. Probamos con ζ(3):")
	fmt.Println()
	fmt.Println("        σ usado                valor        ¿existe la cintura?")
	for _, c := range ctes {
		nota := ""
		if !c.ctrl {
			nota = "  ← ζ(3)"
		}
		et := fmt.Sprintf("%.7f", cintura(c.v))
		if c.v >= 1 {
			et = "NO EXISTE"
		}
		fmt.Printf("        %-18s   %10.7f      %s%s\n", c.nombre, c.v, et, nota)
	}
	xMedio := cintura(0.5)
	fmt.Println()
	fmt.Println("   ⛔ NO SE PUEDE, y el motivo es matemático y no un fallo del código: la función x·cot x")
	fmt.Println("      vale a lo sumo 1 en todo el intervalo (0, π) — se acerca a 1 cuando x tiende a cero y")
	fmt.Println("      baja desde ahí. Así que la ecuación de la cintura NO TIENE SOLUCIÓN para σ mayor que 1.")
	fmt.Printf("      Y ζ(3) = %.7f es mayor que 1. La fórmula la RECHAZA en la puerta.\n", z3)
	fmt.Println()
	fmt.Printf("        para comparar, σ = ½ (la línea crítica) sí tiene cintura: %.7f, con x/π = %.7f\n", xMedio, xMedio/math.Pi)
	fmt.Println("        y todas las constantes de control también quedan afuera, por el mismo motivo.")
	fmt.Println()
	fmt.Println("   ⚡ Y eso es un dato REAL, no un fracaso: la curva madre solo habla de σ dentro de la franja")
	fmt.Println("      crítica, entre 0 y 1. ζ(3) vive afuera. No es que no diga nada: es que no está invitada.")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§4 · la otra puerta: ¿y si ζ(3) fuera el RESULTADO en vez de la entrada?")
	fmt.Println()
	fmt.Println("   En vez de meter ζ(3), preguntamos qué entrada la produce. Es la pregunta inversa, y es")
	fmt.Println("   la única honesta: si algún s natural la devolviera, eso SÍ sería un vínculo.")
	fmt.Println()
	inv := 1 / (1 - z3)
	fmt.Printf("        ¿qué s cumple w(s) = ζ(3)?      s = 1/(1−ζ(3)) = %.9f\n", inv)
	fmt.Printf("        ¿es un número con nombre?       %.9f ... no coincide con ningún natural, primo,\n", inv)
	fmt.Println("                                        ni con −e (−2,718), ni con −π (−3,1416), ni con −2π/√2.")
	fmt.Println()
	inv2 := 1 / (1 - 1/z3)
	fmt.Printf("        ¿y qué s cumple w(s) = 1/ζ(3)?  s = %.9f  — tampoco.\n", inv2)
	fmt.Println()

	// ---------------------------------------------------------------- verdict
	fmt.Println("§5 · EL VEREDICTO")
	fmt.Println()
	fmt.Printf("   ζ(3) pasada por el cambiaformas da w = %.7f, que queda a 0,0014 de 1/6. A primera vista\n", w(z3))
	fmt.Println("   eso IMPRESIONA — parece que ζ(3) casi aterriza en un sexto. Y ahí está el peligro entero")
	fmt.Println("   de esta clase de experimento, porque el control lo desarma:")
	fmt.Println()
	fmt.Printf("        ζ(3) quedó %d° de %d… y quien la superó fue UN NÚMERO AL AZAR, que cayó seis veces\n", pos, len(ds))
	fmt.Println("        más cerca de un redondo que ella. Un número sin ningún significado le ganó.")
	fmt.Println()
	fmt.Println("   Y en la fórmula de ensamble ni siquiera entra: la cintura no existe para σ > 1.")
	fmt.Println()
	fmt.Println("   ⚖️ RESULTADO: NO PASA NADA. Y eso era lo predicho antes de correr, por una razón que")
	fmt.Println("      conviene tener escrita: ζ(3) llega de contar fotones del fondo cósmico, y ni el")
	fmt.Println("      cambiaformas ni la curva madre la derivan de ninguna parte. Meterla a mano es")
	fmt.Println("      ponerle una pieza de otro motor a un motor que anda: entra, gira, y no significa nada.")
	fmt.Println()
	fmt.Println("   📌 Lo que SÍ vale de este experimento es el CONTROL. Sin las nueve constantes de al lado,")
	fmt.Println("      cualquiera de estos números hubiera parecido interesante — y ahí empieza la numerología.")
	fmt.Println("      El experimento se pidió, se hizo completo, y dio negativo. Eso también se registra.")
	fmt.Println()
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(z3, rs, ds, pos, xMedio)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/apery-en-la-casa.svg")
}

func dibujar(z3 float64, rs []res, ds []dist, pos int, xMedio float64) {
	var b strings.Builder
	W, H := 1060, 680
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "APÉRY EN LA CASA — el experimento que se pidió, y dio negativo")
	txt(530, 54, 12, "#8a8676", "middle", "meter ζ(3) en nuestras fórmulas, con nueve controles al lado · F412")

	// the map
	fmt.Fprintf(&b, `<rect x="20" y="76" width="1020" height="230" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 102, 13, "#e8e2d0", "middle", "el cambiaformas w(s) = 1 − 1/s, y dónde cae cada constante")
	cy := 200.0
	fmt.Fprintf(&b, `<line x1="80" y1="%.0f" x2="980" y2="%.0f" stroke="#3a4050" stroke-width="1.5"/>`, cy, cy)
	marca := func(v float64, lab, col string, arriba bool) {
		px := 530 + v*880
		dy := 22.0
		if arriba {
			dy = -16
		}
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.0f" r="5" fill="%s"/>`, px, cy, col)
		txt(px, cy+dy, 11, col, "middle", lab)
	}
	marca(-1, "½ → −1 (la piel)", "#8a8676", true)
	marca(0, "1 → 0 (el polo)", "#8a8676", true)
	marca(0.5, "2 → ½ EXACTO", "#f2c14e", true)
	for _, x := range rs {
		col := "#5ec8ff"
		if !x.c.ctrl {
			col = "#ff6b6b"
		}
		marca(x.wv, "", col, false)
	}
	txt(530+w(z3)*880, cy+42, 12, "#ff6b6b", "middle", fmt.Sprintf("ζ(3) → %.4f", w(z3)))
	txt(530, cy+80, 11, "#5ec8ff", "middle", "● los puntos celestes son las nueve constantes de control — ζ(3) cae entre ellas, sin distinguirse")

	// the ranking
	fmt.Fprintf(&b, `<rect x="20" y="322" width="500" height="250" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 348, 13, "#e8e2d0", "middle", "ranking de cercanía a un número redondo")
	yy := 380.0
	for i, d := range ds {
		if i >= 6 {
			break
		}
		col := "#8a8676"
		if !d.ctrl {
			col = "#ff6b6b"
		}
		txt(50, yy, 11, col, "start", fmt.Sprintf("%d°  %s", i+1, d.nombre))
		txt(490, yy, 11, col, "end", fmt.Sprintf("a %.5f de %s", d.d, d.cual))
		yy += 28
	}
	txt(270, 552, 11, "#ff6b6b", "middle", fmt.Sprintf("ζ(3) quedó %d° de %d: en el montón", pos, len(ds)))

	// the verdict
	fmt.Fprintf(&b, `<rect x="540" y="322" width="500" height="250" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(790, 348, 13, "#ff6b6b", "middle", "RESULTADO: NO PASA NADA")
	txt(564, 384, 11, "#c9c3b0", "start", "Y estaba predicho antes de correr, porque")
	txt(564, 404, 11, "#c9c3b0", "start", "ζ(3) viene de contar fotones del fondo")
	txt(564, 424, 11, "#c9c3b0", "start", "cósmico, y ninguna de nuestras dos")
	txt(564, 444, 11, "#c9c3b0", "start", "fórmulas la deriva de ningún lado.")
	fmt.Fprintf(&b, `<line x1="564" y1="466" x2="1016" y2="466" stroke="#5a2a30"/>`)
	txt(564, 492, 11, "#f2c14e", "start", "Meterla a mano es ponerle una pieza de")
	txt(564, 512, 11, "#f2c14e", "start", "otro motor a un motor que anda: entra,")
	txt(564, 532, 11, "#f2c14e", "start", "gira, y no significa nada.")
	txt(790, 558, 10, "#8a8676", "middle", "trampa vigesimoquinta, declarada antes de calcular")

	// why it matters
	fmt.Fprintf(&b, `<rect x="20" y="588" width="1020" height="76" fill="#121a16" stroke="#2a5a40"/>`)
	txt(530, 614, 12, "#6ee7a0", "middle", "Lo que SÍ vale: el CONTROL. Sin las nueve constantes al lado, cualquiera de estos números habría parecido interesante.")
	txt(530, 640, 11, "#8a8676", "middle", "El experimento se pidió, se hizo completo, y dio negativo. Eso también se registra. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "apery-en-la-casa.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
