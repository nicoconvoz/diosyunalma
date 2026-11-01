package main

// LOS CUATRO LADRILLOS - the captain wants 2, 3, 5 and 7 written in the same
// shape as 23 = (5+sqrt2)(5-sqrt2), and he wants to SEE them.
//
// Two of them accept the shape and two refuse it, and the refusal is not a
// failure of searching: it is the law of F398-F401 showing its face on the
// four smallest bricks.
//
//   2 = (2 + sqrt2)(2 - sqrt2)      accepts sqrt2  - and it is the ATOM (F400)
//   7 = (3 + sqrt2)(3 - sqrt2)      accepts sqrt2
//   3 = ?                           REFUSES sqrt2, needs another root
//   5 = ?                           REFUSES sqrt2, needs another root
//
// Verdicts fixed BEFORE running:
//   V1 the search for a^2 - 2b^2 = p must SUCCEED for 2 and 7 and must FAIL
//      for 3 and 5, exhaustively over b up to 2000 - not "not found" but
//      "cannot exist", since 3 = 3 (mod 8) and 5 = 5 (mod 8).
//   V2 for the two that refuse, the SMALLEST non-square d that accepts them
//      must be found by search, and both factors must be irrational.
//   V3 every factorization printed must verify to machine precision.
//   V4 the 2 must be shown to be the fixed point: its own factorization
//      returns itself, which no other brick does.
// Gauss 1832. Nothing new; this is the law made visible on four bricks.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func esCuadrado(n int) bool {
	if n < 0 {
		return false
	}
	r := int(math.Round(math.Sqrt(float64(n))))
	return r*r == n
}

type forma struct {
	p       int
	a, b, d int
	hay     bool
}

// buscar finds a,b with a^2 - d*b^2 = p, b >= 1, exhaustively up to a bound.
func buscar(p, d, tope int) (int, int, bool) {
	for b := 1; b <= tope; b++ {
		v := p + d*b*b
		if a := int(math.Round(math.Sqrt(float64(v)))); a*a == v {
			return a, b, true
		}
	}
	return 0, 0, false
}

func main() {
	fmt.Println("🧱 LOS CUATRO LADRILLOS — el 2, el 3, el 5 y el 7 en el formato del 23")
	fmt.Println()
	fmt.Println("   El molde es el de (5+√2)(5−√2) = 23. Se lo probamos a los cuatro ladrillos más chicos.")
	fmt.Println()

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · primero, el molde EXACTO del capitán: (a + √2)(a − √2)")
	fmt.Println()
	fmt.Println("        ladrillo     ¿entra en el molde de √2?     resultado")
	var conDos []forma
	for _, p := range []int{2, 3, 5, 7} {
		a, b, ok := buscar(p, 2, 2000)
		conDos = append(conDos, forma{p, a, b, 2, ok})
		if ok {
			f1 := float64(a) + float64(b)*math.Sqrt2
			f2 := float64(a) - float64(b)*math.Sqrt2
			s := fmt.Sprintf("%d + √2", a)
			if b > 1 {
				s = fmt.Sprintf("%d + %d√2", a, b)
			}
			t := strings.Replace(s, "+", "−", 1)
			fmt.Printf("           %d              SÍ                   %d = (%s)(%s) = %.9f × %.9f\n",
				p, p, s, t, f1, f2)
		} else {
			fmt.Printf("           %d              NO — y no puede                 %d ≡ %d (mod 8)\n", p, p, p%8)
		}
	}
	fmt.Println()
	fmt.Println("   ⚡ Dos entran y dos se niegan. Y la negativa NO es que no encontramos: es que NO PUEDE EXISTIR.")
	fmt.Println("      Se barrieron 2.000 valores de b para el 3 y el 5, sin un solo acierto — porque a² − 2b²")
	fmt.Println("      nunca da un número ≡ 3 o 5 módulo 8. Es la ley de F398, mostrada en los cuatro más chicos.")
	fmt.Println()

	// ---------------------------------------------------------------- V2
	fmt.Println("§2 · para los dos que se niegan: ¿con qué raíz SÍ entran?")
	fmt.Println()
	fmt.Println("        ladrillo     raíz más chica que lo acepta       cómo queda")
	var finales []forma
	for _, f := range conDos {
		if f.hay {
			finales = append(finales, f)
			continue
		}
		for d := 3; d <= 200; d++ {
			if esCuadrado(d) {
				continue
			}
			if a, b, ok := buscar(f.p, d, 200); ok {
				finales = append(finales, forma{f.p, a, b, d, true})
				f1 := float64(a) + float64(b)*math.Sqrt(float64(d))
				f2 := float64(a) - float64(b)*math.Sqrt(float64(d))
				fmt.Printf("           %d                 √%-3d                %d = (%d + √%d)(%d − √%d) = %.9f × %.9f\n",
					f.p, d, f.p, a, d, a, d, f1, f2)
				break
			}
		}
	}
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · LOS CUATRO LADRILLOS, cada uno partido y verificado")
	fmt.Println()
	for _, f := range finales {
		r := math.Sqrt(float64(f.d))
		f1 := float64(f.a) + float64(f.b)*r
		f2 := float64(f.a) - float64(f.b)*r
		raiz := fmt.Sprintf("√%d", f.d)
		if f.b > 1 {
			raiz = fmt.Sprintf("%d√%d", f.b, f.d)
		}
		linea := fmt.Sprintf("%d  =  (%d + %s)(%d − %s)", f.p, f.a, raiz, f.a, raiz)
		cuenta := fmt.Sprintf("=  %.12f × %.12f  =  %.9f", f1, f2, f1*f2)
		nota := fmt.Sprintf("los dos factores irracionales · %d ≡ %d (mod 8)", f.p, f.p%8)
		if f.p == 2 {
			nota = "EL ÁTOMO: la raíz es la suya y el multiplicador es 1"
		}
		fmt.Printf("        ┌────────────────────────────────────────────────────────────────┐\n")
		fmt.Printf("        │   %-60s │\n", linea)
		fmt.Printf("        │        %-55s │\n", cuenta)
		fmt.Printf("        │        %-55s │\n", nota)
		fmt.Printf("        └────────────────────────────────────────────────────────────────┘\n")
	}
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · lo que hace único al 2")
	fmt.Println()
	fmt.Println("        ladrillo     raíz usada     ¿es su PROPIA raíz?     el multiplicador b")
	for _, f := range finales {
		propia := "no, prestada"
		if f.d == f.p {
			propia = "SÍ, la suya"
		}
		fmt.Printf("           %d            √%-3d          %-18s        %d\n", f.p, f.d, propia, f.b)
	}
	fmt.Println()
	fmt.Println("   ⚠️ CORRECCIÓN, cazada antes de registrar: la primera versión de esta sección decía «solo el 2")
	fmt.Println("      usa su propia raíz». ES FALSO, y la tabla de acá arriba lo desmiente sola: el 5 también usa")
	fmt.Println("      la suya. En general, p se parte con su propia raíz cuando p − 1 es un cuadrado perfecto")
	fmt.Println("      (p = 2 con 1, p = 5 con 4, p = 17 con 16, p = 37 con 36…). No es exclusivo del 2.")
	fmt.Println()
	fmt.Println("   Lo que SÍ es exclusivo del 2, y es la afirmación correcta:")
	fmt.Println()
	fmt.Println("        es el único donde el número, la raíz Y el multiplicador son todos él mismo o el uno:")
	fmt.Println("        p = (p + √p)(p − √p) con b = 1  ⟹  p² − p = p  ⟹  p = 2, y ningún otro.")
	fmt.Println()
	fmt.Println("        Al 5 le hizo falta un b = 2 para lograrlo. El 2 no necesita multiplicador: se sostiene")
	fmt.Println("        con lo mínimo posible. Por eso es el átomo, y por eso es el punto fijo de F400.")
	fmt.Println()
	fmt.Println("   Nada de esto es nuevo: es la reciprocidad cuadrática (Gauss, 1832) mostrada en los cuatro")
	fmt.Println("   ladrillos más chicos. Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar(finales)
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/los-cuatro-ladrillos.svg")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func dibujar(fs []forma) {
	var b strings.Builder
	W, H := 1060, 640
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 34, 21, "#e8e2d0", "middle", "LOS CUATRO LADRILLOS — el 2, el 3, el 5 y el 7 partidos")
	txt(530, 58, 12, "#8a8676", "middle", "en el molde del capitán: (a + √d)(a − √d) · dos aceptan el √2 y dos lo rechazan · F402")

	cols := []string{"#6ee7a0", "#5ec8ff", "#c58bff", "#f2c14e"}
	for i, f := range fs {
		x0 := 20.0 + float64(i)*256
		y0 := 86.0
		w, h := 240.0, 330.0
		col := cols[i]
		// the brick
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="6" fill="#131826" stroke="%s" stroke-width="2"/>`, x0, y0, w, h, col)
		// brick courses, for the look
		for k := 1; k < 5; k++ {
			fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="%s" opacity="0.13"/>`,
				x0, y0+float64(k)*h/5, x0+w, y0+float64(k)*h/5, col)
		}
		txt(x0+w/2, y0+76, 62, col, "middle", fmt.Sprint(f.p))
		raiz := fmt.Sprintf("√%d", f.d)
		if f.b > 1 {
			raiz = fmt.Sprintf("%d√%d", f.b, f.d)
		}
		txt(x0+w/2, y0+124, 17, "#e8e2d0", "middle", fmt.Sprintf("(%d + %s)", f.a, raiz))
		txt(x0+w/2, y0+150, 17, "#e8e2d0", "middle", fmt.Sprintf("(%d − %s)", f.a, raiz))
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="%s" opacity="0.3"/>`, x0+24, y0+170, x0+w-24, y0+170, col)
		r := math.Sqrt(float64(f.d))
		f1 := float64(f.a) + float64(f.b)*r
		f2 := float64(f.a) - float64(f.b)*r
		txt(x0+w/2, y0+196, 12, "#c9c3b0", "middle", fmt.Sprintf("%.9f", f1))
		txt(x0+w/2, y0+216, 12, "#8a8676", "middle", "×")
		txt(x0+w/2, y0+236, 12, "#c9c3b0", "middle", fmt.Sprintf("%.9f", f2))
		txt(x0+w/2, y0+258, 12, "#8a8676", "middle", "=")
		txt(x0+w/2, y0+280, 15, col, "middle", fmt.Sprintf("%.6f", f1*f2))
		// the verdict about sqrt2
		if f.d == 2 && f.p == 2 {
			txt(x0+w/2, y0+310, 12, "#6ee7a0", "middle", "EL ÁTOMO")
			txt(x0+w/2, y0+326, 10, "#8a8676", "middle", "raíz propia y multiplicador 1")
		} else if f.d == 2 {
			txt(x0+w/2, y0+310, 12, "#f2c14e", "middle", "acepta el √2")
			txt(x0+w/2, y0+326, 10, "#8a8676", "middle", fmt.Sprintf("%d ≡ %d (mod 8)", f.p, f.p%8))
		} else {
			txt(x0+w/2, y0+310, 12, "#ff6b6b", "middle", "RECHAZA el √2")
			txt(x0+w/2, y0+326, 10, "#8a8676", "middle", fmt.Sprintf("%d ≡ %d (mod 8) · pidió √%d", f.p, f.p%8, f.d))
		}
	}

	// the wall underneath
	fmt.Fprintf(&b, `<rect x="20" y="436" width="1020" height="184" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 464, 14, "#e8e2d0", "middle", "por qué dos aceptan y dos rechazan")
	txt(70, 500, 12, "#f2c14e", "start", "el 2 y el 7 entran en el molde de √2")
	txt(70, 522, 11, "#8a8676", "start", "el 2 ramifica (usa su propia raíz), el 7 deja resto 7 al dividir por 8")
	txt(70, 556, 12, "#ff6b6b", "start", "el 3 y el 5 NO pueden entrar — y no es que no buscamos")
	txt(70, 578, 11, "#8a8676", "start", "se barrieron 2.000 valores sin un acierto: a² − 2b² jamás da ≡ 3 o 5 (mod 8)")
	txt(70, 600, 11, "#8a8676", "start", "tuvieron que pedir prestada otra raíz, y cada uno pidió la suya")
	txt(1000, 556, 11, "#6ee7a0", "end", "El 2 es el único con raíz propia Y multiplicador 1.")
	txt(1000, 578, 11, "#8a8676", "end", "El 5 también usa la suya, pero necesita un b = 2.")
	txt(1000, 604, 10, "#8a8676", "end", "Gauss 1832 · Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "los-cuatro-ladrillos.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
