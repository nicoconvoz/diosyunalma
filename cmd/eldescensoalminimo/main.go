package main

// EL DESCENSO AL MINIMO - the captain corrected me: he did not want to climb b.
// He wants the SAME operation applied to the 5 itself, recursively, until what
// is left is the minimum number that still carries sqrt2.
//
// The operation that made 23 out of 5 is n -> n^2 - 2. So going DOWN is the
// inverse: n -> sqrt(n + 2). And the 5 does split, exactly:
//
//     5 = (sqrt7 + sqrt2)(sqrt7 - sqrt2)      because 7 - 2 = 5
//     sqrt7 = (sqrt(2+sqrt7) + sqrt2)(...)    and so on, forever
//
// Descending from 5: 5, sqrt7, sqrt(2+sqrt7), ... The question is where it
// stops, and it stops at a number that reproduces itself:
//
//     (2 + sqrt2)(2 - sqrt2) = 4 - 2 = 2
//
// THE MINIMUM IS 2. It is the fixed point of the captain's operation.
//
// Verdicts fixed BEFORE running:
//   V1 (the exact split) 5 = (sqrt7+sqrt2)(sqrt7-sqrt2) must verify to machine
//      precision, and so must every rung of the descent.
//   V2 (the floor) the descent must converge, and its limit must be the fixed
//      point of x = sqrt(x+2), which algebra says is exactly 2.
//   V3 (self-reproduction) 2 is the ONLY number > 0 that the operation returns
//      unchanged: (2+sqrt2)(2-sqrt2) = 2. Checked against neighbours.
//   V4 (why, named) n -> n^2 - 2 is ANGLE DOUBLING: with n = 2cosh(t) it sends
//      t to 2t, so descending halves t and t -> 0 gives 2cosh(0) = 2. Measured
//      by recovering t at each rung.
//   V5 (the integer side) climbing instead of descending from 3 gives
//      3 -> 7 -> 47 -> 2207 -> ... which must be the Lucas numbers L(2^k).
//      Their primality is measured, not assumed.
//   V6 (trap #19, declared) the nested radical sqrt(2+sqrt(2+...)) = 2 is a
//      classical identity, not a discovery. Reported with its name on it.
// Chebyshev, Lucas, and a nested radical every student meets. Nothing new.

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("⬇️  EL DESCENSO AL MÍNIMO — hacerle al 5 lo mismo que le hicimos al 23")
	fmt.Println("   Corrección: el capitán no quería subir el b. Quería la operación aplicada AL 5 MISMO.")
	fmt.Println()

	// ---------------------------------------------------------------- V1
	fmt.Println("§1 · el 5 SÍ se parte igual — y acá está el número que hay que poner")
	fmt.Println()
	fmt.Println("   La operación que hizo el 23 a partir del 5 es n ↦ n² − 2 (porque 5² − 2 = 23).")
	fmt.Println("   Entonces, para partir el 5 igual, hace falta un a con a² − 2 = 5, o sea a = √7:")
	fmt.Println()
	r7, r2 := math.Sqrt(7), math.Sqrt2
	fmt.Printf("        5 = (√7 + √2)(√7 − √2) = (%.12f)(%.12f) = %.12f\n", r7+r2, r7-r2, (r7+r2)*(r7-r2))
	fmt.Println()
	fmt.Println("   ⚡ Funciona: 7 − 2 = 5. Ninguno de los dos factores es entero, y el 5 queda partido.")
	fmt.Println("      Y ahora al √7 se le puede hacer lo mismo, y al que salga también, y así para siempre.")
	fmt.Println()

	// ---------------------------------------------------------------- V2 + V4
	fmt.Println("§2 · EL DESCENSO: bajando escalón por escalón desde el 5")
	fmt.Println()
	fmt.Println("        paso      el número         se parte como            t (con n = 2·cosh t)     t anterior / 2")
	x := 5.0
	tPrev := 0.0
	for i := 0; i <= 10; i++ {
		t := math.Acosh(x / 2)
		mitad := ""
		if i > 0 {
			mitad = fmt.Sprintf("%.12f", tPrev/2)
		}
		sig := math.Sqrt(x + 2)
		fmt.Printf("         %2d     %.12f    (√%.6f + √2)(√%.6f − √2)     %.12f      %s\n",
			i, x, x+2, x+2, t, mitad)
		tPrev = t
		x = sig
	}
	fmt.Println()
	fmt.Printf("   ⚡⚡ El descenso CONVERGE, y converge a %.12f — o sea, al 2.\n", x)
	fmt.Println("      Y mirá la última columna: en cada escalón el parámetro t se PARTE AL MEDIO, exacto.")
	fmt.Println()

	// ---------------------------------------------------------------- V3
	fmt.Println("§3 · POR QUÉ SE DETIENE AHÍ: el 2 se reproduce a sí mismo")
	fmt.Println()
	fmt.Printf("        (2 + √2)(2 − √2) = (%.12f)(%.12f) = %.12f\n", 2+r2, 2-r2, (2+r2)*(2-r2))
	fmt.Println()
	fmt.Println("   ⚡⚡ Da 2. El mismo 2 con el que empezó. Es el ÚNICO número que la operación devuelve intacto,")
	fmt.Println("      y por eso es el fondo: aplicarle la operación de nuevo no lo mueve más.")
	fmt.Println()
	fmt.Println("   La cuenta que lo demuestra, en un renglón:")
	fmt.Println("        buscamos x con x² − 2 = x  ⟹  x² − x − 2 = 0  ⟹  (x − 2)(x + 1) = 0  ⟹  x = 2")
	fmt.Println("        (la otra raíz, −1, no sirve: los factores tienen que dar un producto positivo)")
	fmt.Println()
	fmt.Println("        vecino        (n+√2)(n−√2) = n² − 2      ¿vuelve al mismo?")
	for _, n := range []float64{1.5, 1.9, 2.0, 2.1, 2.5, 3.0} {
		v := n*n - 2
		et := "no, se va"
		if math.Abs(v-n) < 1e-12 {
			et = "SÍ ← el punto fijo"
		}
		fmt.Printf("        %.1f              %+.6f                  %s\n", n, v, et)
	}
	fmt.Println()
	fmt.Println("   Abajo del 2 la operación te tira para atrás (1,9 → 1,61 → 0,59 → −1,65: se desbarranca).")
	fmt.Println("   Arriba del 2 te dispara para adelante (2,1 → 2,41 → 3,81 → 12,5). El 2 es el filo del cuchillo.")
	fmt.Println()

	// ---------------------------------------------------------------- V4
	fmt.Println("§4 · el nombre de lo que está pasando: la operación DUPLICA UN ÁNGULO")
	fmt.Println()
	fmt.Println("   Si escribís n = 2·cosh(t), entonces n² − 2 = 2·cosh(2t). La operación del capitán no hace")
	fmt.Println("   más que DUPLICAR t. Por eso bajar es partir t al medio, y por eso el fondo es t = 0:")
	fmt.Println()
	fmt.Printf("        n = 2·cosh(0) = %.12f  ← el mínimo, donde el ángulo se acabó\n", 2*math.Cosh(0))
	fmt.Println()
	fmt.Println("   Y con ángulos de verdad (n = 2·cos θ, para n entre −2 y 2) es exactamente la fórmula del")
	fmt.Println("   ángulo doble: 4cos²θ − 2 = 2·cos(2θ). Es la recurrencia de Chebyshev, de toda la vida.")
	fmt.Println()

	// ---------------------------------------------------------------- V5
	fmt.Println("§5 · el mismo descenso, pero PARA ARRIBA y con enteros")
	fmt.Println()
	fmt.Println("   Si en vez de bajar subís, arrancando del 3 (que es el entero más chico arriba del 2):")
	fmt.Println()
	fmt.Println("        paso        número              ¿primo?      ¿es un número de Lucas?")
	n := big.NewInt(3)
	dos := big.NewInt(2)
	lucas := map[string]string{"3": "L(2)", "7": "L(4)", "47": "L(8)", "2207": "L(16)", "4870847": "L(32)"}
	for i := 1; i <= 6; i++ {
		s := n.String()
		et := "no"
		if n.ProbablyPrime(30) {
			et = "SÍ"
		}
		lu, ok := lucas[s]
		if !ok {
			lu = "L(" + fmt.Sprint(1<<uint(i)) + ")"
		}
		corto := s
		if len(corto) > 30 {
			corto = corto[:14] + "…" + corto[len(corto)-12:] + fmt.Sprintf(" (%d cifras)", len(s))
		}
		fmt.Printf("         %2d      %-46s  %-6s      %s\n", i, corto, et, lu)
		n = new(big.Int).Sub(new(big.Int).Mul(n, n), dos)
	}
	fmt.Println()
	fmt.Println("   ⚡ Son los números de Lucas en las potencias de dos: L(2), L(4), L(8), L(16), L(32)…")
	fmt.Println("      La cadena del capitán, subiendo, camina la sucesión de Lucas de dos en dos potencias.")
	fmt.Println("      (Y la primalidad se MIDE, no se supone: fijate cuál se cae.)")
	fmt.Println()

	// ---------------------------------------------------------------- V6
	fmt.Println("§6 · la trampa decimonovena, declarada")
	fmt.Println()
	fmt.Println("   Escribir el descenso entero de una sola vez da la raíz anidada:")
	fmt.Println()
	fmt.Println("        2 = √(2 + √(2 + √(2 + √(2 + …))))")
	fmt.Println()
	v := 0.0
	for i := 0; i < 60; i++ {
		v = math.Sqrt(2 + v)
	}
	fmt.Printf("        verificado numéricamente con 60 raíces: %.15f\n", v)
	fmt.Println()
	fmt.Println("   ⚠️ Es una identidad clásica que está en cualquier libro de análisis; se reporta con el nombre")
	fmt.Println("      puesto y no como hallazgo. Lo mismo la sucesión de Lucas y la recurrencia de Chebyshev.")
	fmt.Println()

	fmt.Println("§7 · LA RESPUESTA A LA PREGUNTA")
	fmt.Println()
	fmt.Println("   «Al 5 que sale ahí quiero hacerle el equivalente a ± √2 hasta que quede un mínimo número con √2.»")
	fmt.Println()
	fmt.Println("   El mínimo es el DOS. Y no por convención ni por cansancio, sino porque:")
	fmt.Println("     · el 5 se parte en (√7 + √2)(√7 − √2), y cada resultado se vuelve a partir igual")
	fmt.Println("     · el descenso 5 → 2,6458 → 2,1554 → 2,0388 → 2,0097 → … baja siempre y no se pasa")
	fmt.Println("     · el 2 SE REPRODUCE: (2 + √2)(2 − √2) = 2, así que ahí la escalera se termina sola")
	fmt.Println("     · abajo del 2 la operación se desbarranca; arriba se dispara. El 2 es el filo")
	fmt.Println("     · y el motivo de fondo: la operación duplica un ángulo, y bajar lo parte al medio hasta cero")
	fmt.Println()
	fmt.Println("   Nada de esto es nuevo (Chebyshev, Lucas, la raíz anidada clásica).")
	fmt.Println("   Estructura cerrada no es hipótesis demostrada. Todavía no.")

	dibujar()
	fmt.Println("\n🖼️  lámina escrita: galeria/laminas/10-el-telar/el-descenso-al-minimo.svg")
}

func dibujar() {
	var b strings.Builder
	W, H := 1060, 680
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, W, H, W, H)
	b.WriteString(`<rect width="100%" height="100%" fill="#0d1017"/>`)
	txt := func(x, y float64, sz int, fill, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="%d" fill="%s" text-anchor="%s" font-family="Georgia,serif">%s</text>`, x, y, sz, fill, anc, s)
	}
	txt(530, 32, 20, "#e8e2d0", "middle", "EL DESCENSO AL MÍNIMO — y el mínimo es el 2, que se reproduce a sí mismo")
	txt(530, 54, 12, "#8a8676", "middle", "partirle al 5 lo mismo que se le partió al 23, una y otra vez · F400")

	// the descent
	x0, y0, w, h := 20.0, 76.0, 620.0, 330.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#131826" stroke="#2a3040"/>`, x0, y0, w, h)
	txt(x0+w/2, y0+24, 13, "#e8e2d0", "middle", "el descenso desde el 5, escalón por escalón")
	base := y0 + h - 50
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#6ee7a0" stroke-dasharray="5,4"/>`, x0+40, base, x0+w-30, base)
	txt(x0+w-34, base+18, 12, "#6ee7a0", "end", "el 2 — el fondo")
	xx := 5.0
	for i := 0; i <= 8; i++ {
		px := x0 + 70 + float64(i)*62
		py := base - (xx-2)*70
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="6" fill="#f2c14e"/>`, px, py)
		if i < 8 {
			npx := x0 + 70 + float64(i+1)*62
			npy := base - (math.Sqrt(xx+2)-2)*70
			fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#f2c14e" stroke-width="1.4" opacity="0.6"/>`, px, py, npx, npy)
		}
		if i < 5 {
			txt(px, py-14, 11, "#e8e2d0", "middle", fmt.Sprintf("%.4f", xx))
		}
		xx = math.Sqrt(xx + 2)
	}
	txt(x0+w/2, y0+h-16, 10, "#8a8676", "middle", "5 → √7 → √(2+√7) → … cada paso parte el ángulo al medio")

	// the fixed point
	rx := 660.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="380" height="%.0f" fill="#121a16" stroke="#2a5a40"/>`, x0+640, y0, h)
	txt(rx+190, y0+26, 13, "#6ee7a0", "middle", "EL MÍNIMO SE REPRODUCE")
	txt(rx+190, y0+76, 17, "#f2c14e", "middle", "(2 + √2)(2 − √2)")
	txt(rx+190, y0+108, 15, "#e8e2d0", "middle", "= 4 − 2 =")
	txt(rx+190, y0+148, 34, "#6ee7a0", "middle", "2")
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#2a5a40"/>`, rx+30, y0+180, rx+350, y0+180)
	txt(rx+190, y0+208, 12, "#c9c3b0", "middle", "el mismo 2 con el que empezó")
	txt(rx+30, y0+242, 11, "#c9c3b0", "start", "x² − 2 = x  ⟹  (x−2)(x+1) = 0")
	txt(rx+30, y0+264, 11, "#c9c3b0", "start", "⟹  x = 2, y no hay otro")
	txt(rx+30, y0+296, 11, "#8a8676", "start", "abajo del 2 se desbarranca")
	txt(rx+30, y0+314, 11, "#8a8676", "start", "arriba del 2 se dispara")

	// why
	fmt.Fprintf(&b, `<rect x="20" y="422" width="1020" height="118" fill="#131826" stroke="#2a3040"/>`)
	txt(530, 448, 13, "#e8e2d0", "middle", "el porqué: la operación DUPLICA UN ÁNGULO")
	txt(530, 480, 14, "#5ec8ff", "middle", "n = 2·cosh(t)   ⟹   n² − 2 = 2·cosh(2t)")
	txt(530, 508, 12, "#c9c3b0", "middle", "subir duplica t · bajar lo parte al medio · el fondo es t = 0, donde 2·cosh(0) = 2")
	txt(530, 530, 11, "#8a8676", "middle", "es la recurrencia de Chebyshev; con ángulos reales, 4cos²θ − 2 = 2cos(2θ)")

	// lucas + trap
	fmt.Fprintf(&b, `<rect x="20" y="556" width="500" height="106" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 580, 12, "#e8e2d0", "middle", "y para arriba, desde el 3:")
	txt(270, 608, 13, "#f2c14e", "middle", "3 → 7 → 47 → 2207 → 4870847 → …")
	txt(270, 632, 11, "#8a8676", "middle", "los números de Lucas en las potencias de dos")
	txt(270, 650, 11, "#8a8676", "middle", "L(2) · L(4) · L(8) · L(16) · L(32)")

	fmt.Fprintf(&b, `<rect x="540" y="556" width="500" height="106" fill="#1a1216" stroke="#5a2a30"/>`)
	txt(790, 580, 12, "#ff6b6b", "middle", "⚠️ trampa decimonovena, declarada")
	txt(790, 606, 13, "#c9c3b0", "middle", "2 = √(2 + √(2 + √(2 + …)))")
	txt(790, 630, 11, "#8a8676", "middle", "identidad clásica de cualquier libro de análisis.")
	txt(790, 648, 11, "#8a8676", "middle", "Chebyshev y Lucas tampoco son nuestros. Todavía no.")
	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "el-descenso-al-minimo.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
