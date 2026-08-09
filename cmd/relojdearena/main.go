// Command relojdearena measures the captain's last image: an hourglass stacked
// at a point in every direction. Dimension 0 is the point 0 - there are no
// other dimensions, only one with a zone so dense; and the same thing that
// exists there exists at every other scale, because every point, infinitely
// small or infinitely large, carries the same 1/2 relation.
//
// The picture is exact. Under the shapeshifter w = 1 - 1/rho a pearl at height
// gamma lands on the unit circle at angle
//
//	theta(gamma) = 2 arctan(1/(2 gamma))  ~  1/gamma
//
// so as the height grows the angle SHRINKS: pearl after pearl piles up against
// the clasp at w = +1. Every scale of the world - the tens, the millions, the
// 10^42 the train hunted - is stacked at one single point, and none of them
// ever reaches it. That is the hourglass, and its waist is dimension 0.
//
// And the 1/2 relation holds at every one of those scales, exactly:
//
//	|w| = 1  <=>  |rho| = |rho - 1|  <=>  Re rho = 1/2
//
// measured here from gamma = 14 up to gamma = 10^100 - eighty-six orders of
// magnitude - with the same answer to the machine's last bit.
//
// THE HONEST PART, and it matters: on the line that relation is an IDENTITY.
// It costs nothing and proves nothing new, because "carrying the 1/2 relation
// at every scale" is exactly what being on the line MEANS. The hourglass is
// the portrait of the hypothesis, not its proof. What no measurement settles
// is whether every pearl is on the line to begin with.
package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
)

func main() {
	fmt.Println("⏳ EL RELOJ DE ARENA — todas las escalas apiladas en un punto, y el mismo ½ en cada una")

	// ---- LAW 1: every scale piles at the clasp ----
	fmt.Println("\nLEY 1 · TODAS LAS ESCALAS SE APILAN CONTRA EL BROCHE")
	fmt.Println("   bajo el cambiaformas, una perla a la altura γ cae en el círculo al ángulo")
	fmt.Println("   θ(γ) = 2·arctan(1/2γ) ≈ 1/γ — cuanto más honda, MÁS CERCA del broche:")
	fmt.Println("\n       altura γ            el ángulo θ          ¿llegó al broche?")
	type fila struct {
		exp   int
		g, th float64
	}
	var filas []fila
	for _, e := range []int{1, 3, 6, 12, 21, 42, 80, 200} {
		g := math.Pow(10, float64(e))
		th := 2 * math.Atan2(1, 2*g)
		filas = append(filas, fila{e, g, th})
		llego := "no"
		if th == 0 {
			llego = "SÍ (se acabó la máquina, no el mundo)"
		}
		fmt.Printf("      10^%-3d          %.6e            %s\n", e, th, llego)
	}
	fmt.Println("   → el ángulo se achica sin fin y JAMÁS toca el broche: la cintura del reloj de arena")
	fmt.Println("     es la dimensión 0, y ninguna escala llega — todas se apilan contra ella")

	// ---- LAW 2: the 1/2 relation at every scale, to the last bit ----
	fmt.Println("\nLEY 2 · LA RELACIÓN ½ EN TODAS LAS ESCALAS — de γ=14 a γ=10^100, con 200 dígitos")
	fmt.Println("   se mide en precisión extendida (el árbitro), para que no quede duda de máquina:")
	fmt.Println("\n       altura γ         | |ρ| − |ρ−1| |          la relación ½")
	prec := uint(700)
	medio := new(big.Float).SetPrec(prec).SetFloat64(0.5)
	type fb struct {
		et  string
		dev *big.Float
	}
	var fbs []fb
	for _, et := range []string{"14.134725", "1e6", "1e12", "1e21", "1e42", "1e100"} {
		g, _, _ := big.ParseFloat(et, 10, prec, big.ToNearestEven)
		g2 := new(big.Float).SetPrec(prec).Mul(g, g)
		// |rho|^2 = 1/4 + g^2   and   |rho-1|^2 = 1/4 + g^2  (identical on the line)
		a := new(big.Float).SetPrec(prec).Add(new(big.Float).SetPrec(prec).Mul(medio, medio), g2)
		b := new(big.Float).SetPrec(prec).Add(new(big.Float).SetPrec(prec).Mul(medio, medio), g2)
		dev := new(big.Float).SetPrec(prec).Sub(a, b)
		dev.Abs(dev)
		fbs = append(fbs, fb{et, dev})
		fmt.Printf("      %-12s        %s                    se cumple ✓\n", et, dev.Text('e', 3))
	}
	fmt.Println("   → EXACTAMENTE CERO en todas las escalas, con 700 bits de precisión:")
	fmt.Println("     de la decena al 10^100 la relación ½ es la misma, sin degradarse jamás")

	// ---- LAW 3: the waist is infinitely dense ----
	fmt.Println("\nLEY 3 · LA CINTURA ES INFINITAMENTE DENSA — cuántas perlas caben en cada tajada")
	fmt.Println("   la densidad de perlas hasta γ es (γ/2π)·ln(γ/2π): se apiñan sin límite contra el broche")
	fmt.Println("\n       hasta la altura     perlas que caben        el ángulo que ocupan")
	for _, e := range []int{2, 6, 12, 21, 42} {
		g := math.Pow(10, float64(e))
		n := g / (2 * math.Pi) * math.Log(g/(2*math.Pi))
		th := 2 * math.Atan2(1, 2*g)
		fmt.Printf("      10^%-3d            %.4e            desde θ=%.2e hacia el broche\n", e, n, th)
	}
	fmt.Println("   → cada vez MÁS perlas en un ángulo cada vez MÁS chico: eso es la zona densa")
	fmt.Println("     que nombró el capitán. Una sola dimensión, con una cintura infinitamente apretada")

	// ---- LAW 4: the honest part ----
	fmt.Println("\n⚖️ LEY 4 · LA PARTE HONESTA — el reloj de arena es el RETRATO, no la demostración")
	fmt.Println("   sobre la línea, la relación ½ a toda escala es una IDENTIDAD: se cumple por álgebra,")
	fmt.Println("   no por medición, y por eso da cero exacto a 700 bits en cualquier altura.")
	fmt.Println("   Pero justamente por ser identidad NO PRUEBA NADA NUEVO: «llevar la relación ½ en")
	fmt.Println("   todas las escalas» es exactamente lo que SIGNIFICA estar sobre la línea.")
	fmt.Println("   El reloj de arena dibuja la hipótesis con una fidelidad hermosa — y lo que ninguna")
	fmt.Println("   medición resuelve es si todas las perlas están sobre la línea para empezar.")

	fmt.Println("\n════════ LO QUE VIO EL CAPITÁN ════════")
	fmt.Println("«Un reloj de arena apilado en un punto en todas las direcciones; la dimensión 0 es el")
	fmt.Println("punto 0, no hay más dimensiones: es una sola con una zona tan densa. Y lo mismo que")
	fmt.Println("existe ahí existe en las demás escalas: cada punto infinitamente chico o infinitamente")
	fmt.Println("grande tiene esta relación, el ½.»")
	fmt.Println("\nMedido: es exacto. Todas las escalas del mundo se apilan contra un solo punto sin")
	fmt.Println("alcanzarlo jamás; la cintura se aprieta sin fin; y la relación ½ vale idéntica de la")
	fmt.Println("decena al 10^100, con cero exacto a 700 bits. Una sola dimensión, un solo punto,")
	fmt.Println("una sola relación repetida en cada escala.")
	fmt.Println("\nY la honestidad al pie: eso es el retrato de la hipótesis, no su llave. Todavía no.")

	// ---- picture ----
	var b strings.Builder
	W, H := 1500.0, 980.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="25" text-anchor="middle" font-family="Georgia" fill="#dce8f7">⏳ EL RELOJ DE ARENA — todas las escalas apiladas en un punto</text>
<text x="%.0f" y="74" font-size="13.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"la dimensión 0 es el punto 0; no hay más dimensiones, es una sola con una zona tan densa — y cada punto tiene la relación ½" — el capitán</text>`,
		W, H, W, H, W/2, W/2)

	// the hourglass
	hx, hy := 400.0, 400.0
	fmt.Fprintf(&b, `<rect x="60" y="105" width="680" height="580" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.8"/>
<text x="400" y="141" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA CINTURA ES LA DIMENSIÓN 0</text>
<path d="M %.0f %.0f L %.0f %.0f L %.0f %.0f L %.0f %.0f Z" fill="#12305c" opacity="0.45" stroke="#7fd7a8" stroke-width="2"/>
<path d="M %.0f %.0f L %.0f %.0f L %.0f %.0f L %.0f %.0f Z" fill="#12305c" opacity="0.45" stroke="#7fd7a8" stroke-width="2"/>
<circle cx="%.0f" cy="%.0f" r="7" fill="#ffd166"/>
<text x="%.0f" y="%.0f" font-size="12.5" text-anchor="middle" fill="#ffd166">el broche · dimensión 0</text>`,
		hx-230, hy-230, hx+230, hy-230, hx+8, hy, hx-8, hy,
		hx-8, hy, hx+8, hy, hx+230, hy+230, hx-230, hy+230,
		hx, hy, hx+150, hy+4)
	// scales stacking toward the waist
	for i, e := range []int{1, 3, 6, 12, 21, 42} {
		d := 220.0 / math.Pow(1.55, float64(i))
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ffd97f" stroke-width="1.6" opacity="%.2f"/>
<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="#ffd97f" stroke-width="1.6" opacity="%.2f"/>
<text x="%.1f" y="%.1f" font-size="10.5" text-anchor="end" fill="#8fa8c7">10^%d</text>`,
			hx-d, hy-d, hx+d, hy-d, 0.35+0.1*float64(i),
			hx-d, hy+d, hx+d, hy+d, 0.35+0.1*float64(i),
			hx-d-8, hy-d+4, e)
	}
	fmt.Fprintf(&b, `<text x="400" y="660" font-size="12.5" text-anchor="middle" fill="#dce8f7">cada escala se apila más cerca del punto — y ninguna lo alcanza jamás</text>`)

	// the measurements
	fmt.Fprintf(&b, `<rect x="770" y="105" width="670" height="290" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.8"/>
<text x="1105" y="141" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LA RELACIÓN ½ EN TODAS LAS ESCALAS</text>
<text x="1105" y="169" font-size="12.5" text-anchor="middle" fill="#dce8f7">medida con 700 bits, para que no quede duda de máquina:</text>
<text x="1105" y="197" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">altura γ            | |ρ| − |ρ−1| |</text>`)
	for i, f := range fbs {
		fmt.Fprintf(&b, `<text x="1105" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">%-12s          %s</text>`,
			222.0+float64(i)*24, f.et, f.dev.Text('e', 3))
	}
	fmt.Fprintf(&b, `<text x="1105" y="380" font-size="13.5" text-anchor="middle" fill="#ffd166">CERO EXACTO de la decena al 10^100 — la relación no se degrada nunca</text>`)

	// the angles
	fmt.Fprintf(&b, `<rect x="770" y="420" width="670" height="265" rx="10" fill="#0d2547" stroke="#7fb2ff" stroke-width="1.8"/>
<text x="1105" y="456" font-size="15.5" text-anchor="middle" font-family="Georgia" fill="#7fb2ff">EL ÁNGULO SE ACHICA SIN FIN</text>
<text x="1105" y="484" font-size="12" text-anchor="middle" font-family="Consolas,monospace" fill="#8fa8c7">altura γ              el ángulo θ hasta el broche</text>`)
	for i, f := range filas {
		if i%2 == 1 {
			continue
		}
		fmt.Fprintf(&b, `<text x="1105" y="%.0f" font-size="12.5" text-anchor="middle" font-family="Consolas,monospace" fill="#dce8f7">10^%-4d                %.4e</text>`,
			508.0+float64(i/2)*24, f.exp, f.th)
	}
	fmt.Fprintf(&b, `<text x="1105" y="638" font-size="13" text-anchor="middle" fill="#7fd7a8">se acerca siempre, no llega nunca — la cintura del reloj</text>
<text x="1105" y="664" font-size="12.5" text-anchor="middle" fill="#8fa8c7">y cada vez MÁS perlas caben en un ángulo cada vez MÁS chico</text>`)

	// honest footer
	fmt.Fprintf(&b, `<rect x="60" y="715" width="1380" height="185" rx="12" fill="#2a1010" stroke="#ff5d73" stroke-width="2"/>
<text x="%.0f" y="751" font-size="16" text-anchor="middle" font-family="Georgia" fill="#ff8fa0">⚖️ LA PARTE HONESTA — el reloj de arena es el RETRATO, no la llave</text>
<text x="%.0f" y="787" font-size="14" text-anchor="middle" fill="#dce8f7">Sobre la línea, la relación ½ a toda escala es una IDENTIDAD: se cumple por álgebra, y por eso da cero exacto a 700 bits.</text>
<text x="%.0f" y="815" font-size="14" text-anchor="middle" fill="#ffd166">Pero justamente por ser identidad no prueba nada nuevo: llevar el ½ en todas las escalas ES lo que significa estar en la línea.</text>
<text x="%.0f" y="849" font-size="14" text-anchor="middle" fill="#ff8fa0">El reloj de arena dibuja la hipótesis con fidelidad hermosa. Lo que ninguna medición resuelve es si TODAS las perlas están en ella.</text>
<text x="%.0f" y="879" font-size="14" text-anchor="middle" fill="#7fd7a8">Todavía no.</text>
<text x="%.0f" y="940" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">Laboratorio Diosyunalma · las dos mitades, 1 completo ⚓ · y sobre todos los libros, el Otro Libro</text>`,
		W/2, W/2, W/2, W/2, W/2, W/2)
	b.WriteString(`</svg>`)
	os.WriteFile("reloj-de-arena.svg", []byte(b.String()), 0644)
	fmt.Println("\nescrita: reloj-de-arena.svg")
}
