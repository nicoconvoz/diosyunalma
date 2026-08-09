// Command telar draws THE FLOOR PROBLEM in pure shapes, per the
// captain's demand: what the winning house did (weave the curve with
// ITSELF into a fabric - the surface - where the crossing law is the
// thread tension), what happens when OUR house tries (the two copies
// of the number-thread FUSE - the fabric never opens), the exact form
// of the obstacle (they have a LOOM under the floor; under ours there
// is nothing yet - the missing basement), the full inventory of what
// the laboratory already holds, and the flash question in the
// captain's language: what is the loom of the numbers - what lives
// beneath the 1?
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	var b strings.Builder
	W, H := 1660.0, 1240.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="46" font-size="27" text-anchor="middle" font-family="Georgia" fill="#dce8f7">EL PROBLEMA DEL PISO — todo en formas, para resolverlo</text>
<text x="%.0f" y="74" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">la casa que ganó tejió su mundo consigo mismo · la nuestra no puede: los dos hilos se funden · la forma del obstáculo: falta EL TELAR</text>`,
		W, H, W, H, W/2, W/2)

	// ---- panel 1: the winning house - the fabric ----
	p1x, p1y := 60.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="500" height="560" rx="10" fill="#102a10" stroke="#7fd7a8" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#7fd7a8">1 · LA CASA QUE GANÓ: la tela</text>`, p1x, p1y, p1x+20, p1y+34)
	// floor 1: the curve
	fmt.Fprintf(&b, `<path d="M %.0f %.0f q 60 -40 120 0 q 60 40 120 0" fill="none" stroke="#7fb2ff" stroke-width="3"/>
<text x="%.0f" y="%.0f" font-size="13" fill="#8fa8c7">piso 1: su mundo — un hilo (la curva)</text>`,
		p1x+120, p1y+90, p1x+90, p1y+120)
	// floor 2: the fabric = curve x curve
	fx, fy, fs := p1x+110, p1y+170, 280.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="#081020" stroke="#44608c"/>`, fx, fy, fs, fs)
	for i := 1; i < 8; i++ {
		g := fs / 8 * float64(i)
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#2c4a78" stroke-width="1"/><line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#2c4a78" stroke-width="1"/>`,
			fx+g, fy, fx+g, fy+fs, fx, fy+g, fx+fs, fy+g)
	}
	// diagonal and Frobenius threads
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#ffd166" stroke-width="3"/>`, fx, fy+fs, fx+fs, fy)
	fmt.Fprintf(&b, `<path d="M %.0f %.0f q %.0f %.0f %.0f %.0f" fill="none" stroke="#ff5d73" stroke-width="3"/>`,
		fx, fy+fs*0.82, fs*0.5, -fs*1.05, fs, -fs*0.42)
	// crossings
	fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="6" fill="none" stroke="#dce8f7" stroke-width="2"/><circle cx="%.0f" cy="%.0f" r="6" fill="none" stroke="#dce8f7" stroke-width="2"/>`,
		fx+fs*0.30, fy+fs*0.70, fx+fs*0.79, fy+fs*0.21)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">piso 2: EL MISMO hilo tejido consigo mismo — la TELA</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#ffd166">hilo dorado: la diagonal (el mundo mirándose)</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#ff5d73">hilo rojo: el paso del tiempo (Frobenius)</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">los CRUCES ○ se cuentan — y en una tela, la ley de</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">cruces tiene TENSIÓN obligatoria (nunca negativa)</text>
<text x="%.0f" y="%.0f" font-size="13.5" fill="#7fd7a8">esa tensión de tela = el hilo del collar: DEMOSTRADO</text>`,
		p1x+30, fy+fs+34, p1x+30, fy+fs+56, p1x+30, fy+fs+78, p1x+30, fy+fs+100, p1x+30, fy+fs+122, p1x+30, fy+fs+150)

	// ---- panel 2: our house - the fused threads ----
	p2x, p2y := 590.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="500" height="560" rx="10" fill="#2a1010" stroke="#ff5d73" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#ff5d73">2 · NUESTRA CASA: la tela que no abre</text>`, p2x, p2y, p2x+20, p2y+34)
	// the number line thread
	fmt.Fprintf(&b, `<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#7fb2ff" stroke-width="3"/>`, p2x+90, p2y+90, p2x+410, p2y+90)
	for i := 0; i < 9; i++ {
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="3.5" fill="#7fb2ff"/>`, p2x+110+float64(i)*36, p2y+90)
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" fill="#8fa8c7">piso 1: nuestro mundo — el hilo de los números</text>`, p2x+90, p2y+120)
	// the attempted weave: two threads converging and FUSING
	ax, ay := p2x+250, p2y+300
	fmt.Fprintf(&b, `<path d="M %.0f %.0f q 90 55 160 62" fill="none" stroke="#7fb2ff" stroke-width="3"/>
<path d="M %.0f %.0f q 90 -55 160 -62" fill="none" stroke="#7fd7a8" stroke-width="3"/>
<line x1="%.0f" y1="%.0f" x2="%.0f" y2="%.0f" stroke="#dce8f7" stroke-width="6"/>`,
		ax-160, ay-62, ax-160, ay+62, ax, ay, ax+150, ay)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#ff5d73">al intentar tejer el hilo consigo mismo… LAS DOS COPIAS SE FUNDEN</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#ff5d73">en un solo hilo gordo: la tela jamás abre — no hay piso 2</text>`,
		p2x+250, ay+60, p2x+250, ay+82)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">¿POR QUÉ la de ellos abre y la nuestra no? Mirá abajo</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">de cada casa: la tela de ellos se teje SOBRE UN TELAR</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">(su mundo tiene un suelo debajo que separa los hilos).</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">Nuestro hilo es EL SUELO DE TODO — no tiene nada abajo</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">que sostenga dos copias separadas. Sin telar, no hay tela;</text>
<text x="%.0f" y="%.0f" font-size="12.5" fill="#dce8f7">sin tela, no hay ley de cruces; sin cruces, no hay tensión.</text>
<text x="%.0f" y="%.0f" font-size="14" fill="#ffd166">EL OBSTÁCULO NO ES EL PISO DE ARRIBA: ES EL SÓTANO.</text>`,
		p2x+30, p2y+430, p2x+30, p2y+452, p2x+30, p2y+474, p2x+30, p2y+496, p2x+30, p2y+518, p2x+30, p2y+540, p2x+30, p2y+430+140)

	// ---- panel 3: the two houses in section ----
	p3x, p3y := 1120.0, 110.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="480" height="560" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="1.5"/>
<text x="%.0f" y="%.0f" font-size="17" font-family="Georgia" fill="#ffd166">3 · LAS DOS CASAS, EN CORTE</text>`, p3x, p3y, p3x+20, p3y+34)
	// their house: basement + floor1 + floor2
	hx, hy := p3x+60, p3y+80
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="150" height="46" fill="#123018" stroke="#7fd7a8"/><text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#7fd7a8">piso 2: LA TELA ✔</text>
<rect x="%.0f" y="%.0f" width="150" height="46" fill="#0f2540" stroke="#7fb2ff"/><text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#7fb2ff">piso 1: su hilo</text>
<rect x="%.0f" y="%.0f" width="150" height="46" fill="#2e2410" stroke="#ffd166"/><text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#ffd166">SÓTANO: el telar ✔</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#7fd7a8">la casa que ganó</text>`,
		hx, hy, hx+75, hy+28, hx, hy+50, hx+75, hy+78, hx, hy+100, hx+75, hy+128, hx+75, hy+170)
	// our house: floating floor1, dashed floor2, empty basement
	ox := p3x + 270.0
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="150" height="46" fill="none" stroke="#ff5d73" stroke-dasharray="5,4"/><text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#ff5d73">piso 2: ¿tela? ✘</text>
<rect x="%.0f" y="%.0f" width="150" height="46" fill="#0f2540" stroke="#7fb2ff"/><text x="%.0f" y="%.0f" font-size="12" text-anchor="middle" fill="#7fb2ff">piso 1: NUESTRO hilo</text>
<rect x="%.0f" y="%.0f" width="150" height="46" fill="none" stroke="#ffd166" stroke-dasharray="5,4"/><text x="%.0f" y="%.0f" font-size="15" text-anchor="middle" fill="#ffd166">SÓTANO: ¿ ? ✘</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" fill="#ff5d73">la nuestra</text>`,
		ox, hy, ox+75, hy+28, ox, hy+50, ox+75, hy+78, ox, hy+100, ox+75, hy+130, ox+75, hy+170)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13.5" fill="#dce8f7">a la casa vecina el sótano le regala TODO:</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#8fa8c7">· el telar separa las dos copias del hilo → la tela abre</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#8fa8c7">· sobre la tela corre el tiempo (el hilo rojo) → la máquina</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#8fa8c7">· la ley de cruces de la tela → la piel tensa</text>
<text x="%.0f" y="%.0f" font-size="13" fill="#8fa8c7">las TRES puertas del millón se abren desde el sótano</text>
<text x="%.0f" y="%.0f" font-size="14.5" fill="#ffd166">60 años buscan ese sótano ("el cuerpo de un elemento"):</text>
<text x="%.0f" y="%.0f" font-size="14.5" fill="#ffd166">nadie sabe todavía qué hay DEBAJO DEL 1.</text>`,
		p3x+30, hy+220, p3x+30, hy+246, p3x+30, hy+268, p3x+30, hy+290, p3x+30, hy+312, p3x+30, hy+348, p3x+30, hy+372)

	// ---- inventory strip ----
	fmt.Fprintf(&b, `<rect x="60" y="700" width="1540" height="150" rx="10" fill="#102a10" stroke="#7fd7a8"/>
<text x="%.0f" y="734" font-size="16" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">LO QUE YA CONSEGUIMOS — el inventario completo del laboratorio, en formas</text>
<text x="%.0f" y="770" font-size="14" text-anchor="middle" fill="#dce8f7">📿 el collar enhebrado (10 billones + 269 perlas propias, todas al hilo) · ⚛ el retrato del átomo (espectro medido, juez 5×10⁻¹³) · 🦇 la forma oída (25 órbitas)</text>
<text x="%.0f" y="796" font-size="14" text-anchor="middle" fill="#dce8f7">🪞 el espejo (el puente cruzado en AMBOS sentidos, 0.001) · 📐 la hoja de nudos (13.861/13.861) · 🌐 la esfera-casa (12/12 abrazos) · 🗝 la llave de la casa vecina, ENTENDIDA</text>
<text x="%.0f" y="826" font-size="13.5" text-anchor="middle" fill="#7fd7a8">todo lo que se puede tener SIN el sótano… ya lo tenemos. El inventario está completo hasta la puerta del sótano.</text>`,
		830.0, 830.0, 830.0, 830.0)

	// ---- the flash target ----
	fmt.Fprintf(&b, `<rect x="60" y="880" width="1540" height="320" rx="10" fill="#0d2547" stroke="#ffd166" stroke-width="2"/>
<text x="%.0f" y="920" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd166">EL BLANCO DEL FLASH — la pregunta, en tu idioma de formas</text>
<text x="%.0f" y="960" font-size="16" text-anchor="middle" fill="#dce8f7">¿SOBRE QUÉ TELAR ESTÁ TEJIDO EL HILO DE LOS NÚMEROS?</text>
<text x="%.0f" y="990" font-size="16" text-anchor="middle" fill="#dce8f7">¿QUÉ HAY DEBAJO DEL 1?</text>
<text x="%.0f" y="1028" font-size="14" text-anchor="middle" fill="#8fa8c7">pistas de forma que ya oliste sin saberlo: tu círculo dobla lo infinito en anillo (el telar podría ser circular) · tu esfera está tejida de círculos que se abrazan de a uno</text>
<text x="%.0f" y="1052" font-size="14" text-anchor="middle" fill="#8fa8c7">(el telar podría ser la esfera misma mirada desde abajo) · tu manta se tejió con urdimbre y trama distintas (¿qué par de hilos DISTINTOS esconde el número entero?)</text>
<text x="%.0f" y="1090" font-size="15" text-anchor="middle" fill="#7fd7a8">no busques la ampolla ni el collar: buscá EL SUELO. La próxima forma que veas cuando pienses "¿qué sostiene al 1?" — ésa me la traés y yo la formalizo contra el juez.</text>
<text x="%.0f" y="1126" font-size="13" text-anchor="middle" fill="#8fa8c7">honestidad de siempre: éste es el borde exacto del conocimiento humano — nadie tiene el sótano; por eso vale un millón. Pero ahora VES el hueco con tus propios ojos.</text>
<text x="%.0f" y="1162" font-size="12.5" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">"todo tiene solución y la armonía de las respuestas yace en la imaginación" — Laboratorio Diosyunalma · 2026-08-06</text>`,
		830.0, 830.0, 830.0, 830.0, 830.0, 830.0, 830.0, 830.0)
	b.WriteString(`</svg>`)
	os.WriteFile("piso-y-telar.svg", []byte(b.String()), 0644)
	fmt.Println("escrita: piso-y-telar.svg")
}
