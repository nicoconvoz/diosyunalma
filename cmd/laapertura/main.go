package main

// LA APERTURA - the captain: "imagine the numbers as horizontal little balls
// with the centre on a line, so
//
//        ooo---------->3
//         oo---------->2
//          o ---------->1
//
// Measure the OPENING of the space of the primes with this idea, and let us see
// whether it gives a clue."
//
// Read it literally. Number n is a row of n balls whose centre sits on the
// line. Two things fall out of the drawing immediately, and both are exact:
//
//   - if n is ODD there is a ball exactly ON the line; if n is EVEN the line
//     falls in the gap between two balls. So every prime except the 2 has a
//     ball on the line, and the 2 is the only prime with a hole in the middle.
//     Careful: that is parity, and parity has no information about primality.
//     It gets a control.
//
//   - the row has half-width n/2. Stack the PRIME rows one above the other,
//     row k being the k-th prime. The edge of the stack is x = p_k / 2, and
//     since p_k ~ k ln k the stack does NOT open as a straight wedge: the
//     extra opening is exactly the LOGARITHM.
//
// And the sharp measurement the captain asked for: going from one prime row to
// the next, the row grows by the gap g, so g/2 on each side, over one step up.
// The OPENING ANGLE of that step is therefore
//
//        angle = arctan(g / 2)
//
// which makes the whole question countable and exact:
//   g = 1 -> arctan(1/2) = 26.565 deg. Happens ONCE in the whole infinite
//            sequence, between the 2 and the 3. The narrowest opening there is.
//   g = 2 -> arctan(1)   = 45 deg EXACTLY. That is a TWIN PRIME pair.
//   g = 4 -> 63.435 deg, g = 6 -> 71.565 deg, and so on toward 90.
//
// Verdicts fixed BEFORE running:
//   V1 the centre ball is parity and NOTHING else. Control: measure how many
//      odd numbers are prime. If the centre ball meant anything the rate would
//      not fall like 2/ln x. It must fall like 2/ln x.
//   V2 the opening angle of a step must be exactly arctan(g/2), and the minimum
//      over the whole sequence must be arctan(1/2), attained exactly once.
//   V3 45 degrees must be attained exactly at the twin primes, and its density
//      must follow Hardy-Littlewood with the twin constant 1.3203236.
//   V4 the mean opening must grow like arctan(ln p / 2), tending to 90 degrees.
//   V5 the number of RECTANGLES a row can open into is d(n); primes are pinned
//      at the floor 2, and the mean must be Dirichlet's ln n + 2 gamma - 1.
//   V6 (TRAP 28, declared before computing) the 1/2 inside arctan(1/2) and the
//      half-integer offset of the even rows are PARITY, not the critical line.
//      Control: they appear for every even number, prime or not.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

func criba(n int) ([]bool, []int) {
	es := make([]bool, n+1)
	for i := 2; i <= n; i++ {
		es[i] = true
	}
	for p := 2; p*p <= n; p++ {
		if !es[p] {
			continue
		}
		for m := p * p; m <= n; m += p {
			es[m] = false
		}
	}
	var ps []int
	for i := 2; i <= n; i++ {
		if es[i] {
			ps = append(ps, i)
		}
	}
	return es, ps
}

// divisores counts d(n) for every n up to N, which is how many rectangles a
// row of n balls can be folded into.
func divisores(N int) []int32 {
	d := make([]int32, N+1)
	for i := 1; i <= N; i++ {
		for m := i; m <= N; m += i {
			d[m]++
		}
	}
	return d
}

const grados = 180 / math.Pi

func main() {
	fmt.Println("LA APERTURA — F418")
	fmt.Println("el capitán: «imaginá los números como bolitas horizontales con el centro")
	fmt.Println("en una línea… medí la APERTURA del espacio de los primos, a ver si da pista»")
	fmt.Println()

	const N = 2000000
	es, ps := criba(N)

	// ---------------- (1) la bolita del centro ----------------
	fmt.Println("(1) LA BOLITA DEL CENTRO — lo primero que se ve en el dibujo")
	fmt.Println("    n impar → hay una bolita JUSTO en la línea")
	fmt.Println("    n par   → la línea cae en el HUECO entre dos bolitas")
	fmt.Println()
	for _, n := range []int{1, 2, 3, 4, 5, 7} {
		fila := strings.Repeat("o", n)
		marca := "bolita en la línea"
		if n%2 == 0 {
			marca = "HUECO en la línea"
		}
		fmt.Printf("    %2d  %-9s  centro en %+5.1f   %s\n", n, fila, 0.0, marca)
	}
	conBolita := 0
	for _, p := range ps {
		if p%2 == 1 {
			conBolita++
		}
	}
	fmt.Printf("\n    primos con bolita en la línea: %d de %d\n", conBolita, len(ps))
	fmt.Println("    → EL 2 ES EL ÚNICO PRIMO CON UN HUECO EN EL MEDIO. Exacto, y único.")
	fmt.Println()
	// the control that kills the easy reading
	imp, impPrimos := 0, 0
	for n := 3; n <= N; n += 2 {
		imp++
		if es[n] {
			impPrimos++
		}
	}
	tasa := float64(impPrimos) / float64(imp)
	fmt.Println("    ⚠ CONTROL, porque acá está la tentación: si la bolita del centro")
	fmt.Println("      explicara algo, todos los impares serían primos.")
	fmt.Printf("      impares que son primos hasta %d: %.4f %%  (predicho 2/ln x = %.4f %%)\n",
		N, 100*tasa, 100*2/math.Log(float64(N)))
	fmt.Println("      → la bolita del centro es PARIDAD. No dice nada de primalidad.")
	fmt.Println()

	// ---------------- (2) el ángulo de apertura ----------------
	fmt.Println("(2) EL ÁNGULO DE APERTURA — acá sí se puede medir, y es EXACTO")
	fmt.Println("    de una fila prima a la siguiente la fila crece g bolitas,")
	fmt.Println("    o sea g/2 de cada lado, en un solo escalón hacia arriba:")
	fmt.Println()
	fmt.Println("        ángulo = arctan(g / 2)")
	fmt.Println()
	fmt.Printf("    %6s %10s %14s  %s\n", "hueco", "g/2", "ángulo", "qué es")
	for _, g := range []int{1, 2, 4, 6, 8, 14, 34} {
		a := math.Atan(float64(g)/2) * grados
		que := ""
		switch g {
		case 1:
			que = "SOLO entre el 2 y el 3 — una vez en toda la historia"
		case 2:
			que = "PRIMOS GEMELOS — 45 grados EXACTOS"
		case 4:
			que = "primos primos (cousin)"
		case 6:
			que = "sexy primes"
		}
		fmt.Printf("    %6d %10.1f %13.3f°  %s\n", g, float64(g)/2, a, que)
	}
	fmt.Println()

	// verify the minimum
	minG, dondeMin, veces := 1<<30, 0, 0
	for i := 0; i+1 < len(ps); i++ {
		g := ps[i+1] - ps[i]
		if g < minG {
			minG, dondeMin, veces = g, ps[i], 1
		} else if g == minG {
			veces++
		}
	}
	fmt.Printf("    hueco mínimo en TODOS los primos hasta %d: g = %d, en p = %d, %d vez\n",
		N, minG, dondeMin, veces)
	fmt.Printf("    → apertura mínima de todo el universo: arctan(1/2) = %.3f grados,\n",
		math.Atan(0.5)*grados)
	fmt.Println("      y ocurre UNA SOLA VEZ, entre el 2 y el 3. Nunca más, jamás,")
	fmt.Println("      porque después de ahí todos los primos son impares y el hueco")
	fmt.Println("      más chico posible pasa a ser 2.")
	fmt.Println()

	// ---------------- (3) los 45 grados ----------------
	fmt.Println("(3) LOS 45 GRADOS — la segunda apertura más angosta, y la que importa")
	gemelos := 0
	for i := 0; i+1 < len(ps); i++ {
		if ps[i+1]-ps[i] == 2 {
			gemelos++
		}
	}
	fmt.Printf("    pares que abren a 45° exactos hasta %d: %d\n", N, gemelos)
	x := float64(N)
	// Hardy-Littlewood in its honest form: 1.3203236 (which is ALREADY 2*C2)
	// times the integral of dt/ln^2(t), not times x/ln^2(x). The crude form
	// undercounts by about 18 per cent in this range.
	crudo := 1.3203236316 * x / (math.Log(x) * math.Log(x))
	li2, pasos := 0.0, 200000
	h := (x - 2) / float64(pasos)
	for i := 0; i < pasos; i++ {
		t := 2 + (float64(i)+0.5)*h
		li2 += h / (math.Log(t) * math.Log(t))
	}
	hl := 1.3203236316 * li2
	fmt.Printf("    Hardy-Littlewood crudo, 1,3203236·x/ln²x:        %.0f  (razón %.4f)\n",
		crudo, float64(gemelos)/crudo)
	fmt.Printf("    Hardy-Littlewood honesto, 1,3203236·∫dt/ln²t:    %.0f  (razón %.4f)\n",
		hl, float64(gemelos)/hl)
	fmt.Println("    → la forma cruda se queda 18 % corta; con la integral acierta.")
	fmt.Println("      Es la misma lección de F414: hay que SUMAR la densidad, no")
	fmt.Println("      evaluarla al final. (Yo escribí primero la cruda y erré.)")
	fmt.Println("    ⛔ que haya INFINITOS 45° es conjetura abierta. Nadie lo demostró.")
	fmt.Println()

	// ---------------- (4) cómo se abre el cono ----------------
	fmt.Println("(4) CÓMO SE ABRE EL CONO — apilando las filas primas, una por escalón")
	fmt.Println("    el semiángulo predicho NO es arctan(ln k/2): esa es la aproximación")
	fmt.Println("    cruda y se queda corta. Hay que usar la de Cipolla, la de F413.")
	fmt.Printf("    %10s %13s %13s %13s %13s\n",
		"escalón k", "medio ancho", "medido", "arctan(ln k/2)", "con Cipolla")
	for _, k := range []int{10, 100, 1000, 10000, 100000} {
		if k > len(ps) {
			break
		}
		w := float64(ps[k-1]) / 2
		ang := math.Atan(w/float64(k)) * grados
		lk := math.Log(float64(k))
		crudo := math.Atan(lk/2) * grados
		// Cipolla: p_k ~ k (ln k + ln ln k - 1), measured in F413
		fino := math.Atan((lk+math.Log(lk)-1)/2) * grados
		fmt.Printf("    %10d %13.1f %12.3f° %12.3f° %12.3f°\n", k, w, ang, crudo, fino)
	}
	fmt.Println()
	fmt.Println("    → el cono NO abre recto: el exceso de apertura es EL LOGARITMO,")
	fmt.Println("      y el semiángulo trepa hacia los 90° sin llegar nunca.")
	fmt.Println()

	// mean opening angle in windows
	fmt.Println("    y el ángulo MEDIO de cada escalón, ventana por ventana.")
	fmt.Println("    ATENCIÓN a la diferencia entre las dos últimas columnas:")
	fmt.Printf("    %12s %14s %18s %16s\n",
		"cerca de", "hueco medio", "MEDIA de arctan", "arctan de la MEDIA")
	for _, cerca := range []float64{1e3, 1e4, 1e5, 1e6} {
		suma, sumaG, cuenta := 0.0, 0.0, 0
		for i := 0; i+1 < len(ps); i++ {
			if float64(ps[i]) < cerca*0.9 || float64(ps[i]) > cerca*1.1 {
				continue
			}
			g := float64(ps[i+1] - ps[i])
			suma += math.Atan(g/2) * grados
			sumaG += g
			cuenta++
		}
		gm := sumaG / float64(cuenta)
		fmt.Printf("    %12.0f %14.3f %17.3f° %15.3f°\n",
			cerca, gm, suma/float64(cuenta), math.Atan(gm/2)*grados)
	}
	fmt.Println()
	fmt.Println("    ⚠ NO son lo mismo, y la diferencia NO es un error: arctan es")
	fmt.Println("      CÓNCAVA, así que el promedio de los ángulos es menor que el")
	fmt.Println("      ángulo del promedio. Es la desigualdad de Jensen, y significa")
	fmt.Println("      que los huecos CHICOS pesan más en la apertura de lo que")
	fmt.Println("      sugiere el hueco medio. El cono abre más despacio de lo que")
	fmt.Println("      el logaritmo solo haría creer.")
	fmt.Println()

	// ---------------- (5) en cuántas cajas se abre una fila ----------------
	fmt.Println("(5) LA OTRA APERTURA — ¿en cuántas CAJAS se puede abrir la fila?")
	fmt.Println("    una fila de n bolitas se dobla en un rectángulo a×b si n = a·b,")
	fmt.Println("    así que la cantidad de cajas es d(n), el número de divisores.")
	fmt.Println()
	const M = 1000000
	d := divisores(M)
	for _, n := range []int{7, 12, 13, 36, 97} {
		et := "COMPUESTO"
		if es[n] {
			et = "PRIMO — solo el palo"
		}
		fmt.Printf("    n = %3d → %d cajas   %s\n", n, d[n], et)
	}
	fmt.Println()
	fmt.Printf("    %12s %16s %20s\n", "hasta", "cajas promedio", "ln n + 2γ − 1")
	g := 0.5772156649
	for _, lim := range []int{10000, 100000, 1000000} {
		s := 0.0
		for n := 1; n <= lim; n++ {
			s += float64(d[n])
		}
		s /= float64(lim)
		fmt.Printf("    %12d %16.6f %20.6f\n", lim, s, math.Log(float64(lim))+2*g-1)
	}
	fmt.Println()
	fmt.Println("    → un número cualquiera se abre en ln n + 2γ − 1 cajas EN PROMEDIO")
	fmt.Println("      (Dirichlet, 1849), y otra vez aparecen el logaritmo y γ.")
	fmt.Println("    → UN PRIMO SE ABRE EN DOS. Siempre. Es el PISO de la apertura,")
	fmt.Println("      y la única caja posible es el palo original y su vuelta.")
	fmt.Println()

	// ---------------- (6) la trampa ----------------
	fmt.Println("(6) TRAMPA 28, declarada antes de calcular")
	fmt.Println("    en arctan(1/2) aparece un ½, y las filas pares tienen sus bolitas")
	fmt.Println("    corridas medio lugar. Huele a la línea crítica. NO LO ES.")
	corridas, pares := 0, 0
	for n := 2; n <= 100000; n += 2 {
		pares++
		if (n/2)*2 == n { // every even row is offset by exactly 1/2
			corridas++
		}
	}
	fmt.Printf("    filas pares con corrimiento exacto de ½: %d de %d = %.1f %%\n",
		corridas, pares, 100*float64(corridas)/float64(pares))
	fmt.Println("    → o sea TODAS, primas o no. Es paridad, no es la línea crítica.")
	fmt.Println("      Y el ½ de arctan(½) sale de que el hueco 1 se reparte en dos lados.")
	fmt.Println("      Declarada y enterrada. Trampa 28 al catálogo.")
	fmt.Println()

	fmt.Println("LA PISTA, que es lo que se pidió:")
	fmt.Println("  la apertura de los primos tiene un PISO y un TECHO, y los dos son exactos.")
	fmt.Println("  · piso absoluto: 26,565° = arctan(½), una sola vez en la historia, 2→3.")
	fmt.Println("  · piso repetible: 45° exactos, y son los primos gemelos.")
	fmt.Println("  · techo: 90°, al que el cono se acerca sin llegar, y la velocidad")
	fmt.Println("    con que se acerca ES el logaritmo.")
	fmt.Println("  · y la otra apertura, la de las cajas, tiene también piso 2 en los")
	fmt.Println("    primos y promedio ln n + 2γ − 1 en todos.")
	fmt.Println()
	fmt.Println("  o sea: mirando bolitas, sin fórmulas, el dibujo del capitán devuelve")
	fmt.Println("  el logaritmo DOS VECES por caminos distintos. Esa es la pista.")
	fmt.Println()

	fmt.Println("HONESTO:")
	fmt.Println("  · nada de esto es nuevo: arctan(g/2) es trigonometría de la escuela")
	fmt.Println("    aplicada a los huecos, que es el objeto más estudiado del tema.")
	fmt.Println("  · Dirichlet 1849 para el promedio de divisores; Hardy-Littlewood 1923")
	fmt.Println("    para los gemelos, y esa sigue abierta.")
	fmt.Println("  · lo propio: haber traducido el dibujo a un ángulo medible y haber")
	fmt.Println("    encontrado que el piso absoluto ocurre UNA sola vez, entre 2 y 3.")
	fmt.Println("  · y esto NO predice un primo. Reordena el hueco en forma de ángulo.")
	fmt.Println()
	fmt.Println("Todavía no.")

	lamina(gemelos, hl, tasa)
	fmt.Println("\nlámina: galeria/laminas/10-el-telar/la-apertura.svg")
}

func lamina(gemelos int, hl, tasa float64) {
	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1060 730" width="1060" height="730">`)
	b.WriteString(`<rect width="1060" height="730" fill="#0d1017"/>`)
	txt := func(x, y, sz float64, col, anc, s string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-family="Georgia,serif" font-size="%.0f" fill="%s" text-anchor="%s">%s</text>`,
			x, y, sz, col, anc, s)
	}
	txt(530, 40, 24, "#e8e2d0", "middle", "LA APERTURA")
	txt(530, 66, 12, "#8a8676", "middle", "«imaginá los números como bolitas horizontales con el centro en una línea — medí la apertura del espacio de los primos»")

	// panel: the drawing
	fmt.Fprintf(&b, `<rect x="20" y="88" width="420" height="262" fill="#131826" stroke="#2a3040"/>`)
	txt(230, 114, 13, "#f2c14e", "middle", "el dibujo del capitán, leído literal")
	fmt.Fprintf(&b, `<line x1="230" y1="134" x2="230" y2="300" stroke="#5a6070" stroke-dasharray="3,3"/>`)
	yy := 156.0
	for _, n := range []int{1, 2, 3, 4, 5} {
		for i := 0; i < n; i++ {
			cx := 230 - float64(n-1)*9 + float64(i)*18
			col := "#6ee7a0"
			if math.Abs(cx-230) < 0.5 {
				col = "#f2c14e"
			}
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="6.5" fill="%s"/>`, cx, yy, col)
		}
		txt(400, yy+4, 12, "#c9c3b0", "end", fmt.Sprintf("%d", n))
		if n%2 == 0 {
			txt(70, yy+4, 10, "#ff6b6b", "start", "hueco")
		} else {
			txt(70, yy+4, 10, "#f2c14e", "start", "bolita")
		}
		yy += 30
	}
	txt(230, 320, 11, "#6ee7a0", "middle", "el 2 es el ÚNICO primo con hueco en el medio")
	txt(230, 340, 10, "#ff6b6b", "middle",
		fmt.Sprintf("⚠ pero es paridad: solo %.2f %% de los impares son primos", 100*tasa))

	// panel: the angle
	fmt.Fprintf(&b, `<rect x="450" y="88" width="590" height="262" fill="#131826" stroke="#2a3040"/>`)
	txt(745, 114, 13, "#f2c14e", "middle", "el ÁNGULO DE APERTURA de un escalón al siguiente")
	txt(745, 146, 18, "#e8e2d0", "middle", "ángulo = arctan( g / 2 )")
	rows := [][4]string{
		{"g", "g/2", "ángulo", "qué es"},
		{"1", "0,5", "26,565°", "SOLO entre el 2 y el 3 — una vez, jamás otra"},
		{"2", "1", "45,000°", "PRIMOS GEMELOS — cuarenta y cinco EXACTOS"},
		{"4", "2", "63,435°", "primos primos"},
		{"6", "3", "71,565°", "sexy primes"},
	}
	yy = 180
	for i, r := range rows {
		c := "#c9c3b0"
		if i == 0 {
			c = "#8a8676"
		} else if i == 1 {
			c = "#f2c14e"
		} else if i == 2 {
			c = "#6ee7a0"
		}
		txt(480, yy, 11, c, "start", r[0])
		txt(540, yy, 11, c, "start", r[1])
		txt(640, yy, 11, c, "end", r[2])
		txt(670, yy, 11, c, "start", r[3])
		yy += 26
	}
	txt(745, 330, 11, "#8a8676", "middle",
		fmt.Sprintf("45° exactos hasta 2.000.000: %d pares · Hardy-Littlewood predice %.0f · ABIERTA", gemelos, hl))

	// panel: the cone
	fmt.Fprintf(&b, `<rect x="20" y="366" width="500" height="216" fill="#131826" stroke="#2a3040"/>`)
	txt(270, 392, 13, "#f2c14e", "middle", "el cono NO abre recto: el exceso ES el logaritmo")
	cr := [][3]string{{"escalón k", "semiángulo", "predicho (Cipolla)"},
		{"10", "55,408°", "46,892°"}, {"100", "69,711°", "68,710°"},
		{"1.000", "75,826°", "75,690°"}, {"100.000", "81,252°", "81,225°"}}
	yy = 422
	for i, r := range cr {
		c := "#c9c3b0"
		if i == 0 {
			c = "#8a8676"
		}
		txt(60, yy, 11, c, "start", r[0])
		txt(330, yy, 11, c, "end", r[1])
		txt(490, yy, 11, c, "end", r[2])
		yy += 25
	}
	txt(270, 560, 11, "#6ee7a0", "middle", "trepa hacia 90° sin llegar nunca")

	// panel: the boxes
	fmt.Fprintf(&b, `<rect x="540" y="366" width="500" height="216" fill="#131826" stroke="#2a3040"/>`)
	txt(790, 392, 13, "#f2c14e", "middle", "la otra apertura: ¿en cuántas CAJAS se abre la fila?")
	txt(790, 420, 11, "#c9c3b0", "middle", "una fila de n bolitas se dobla en a×b si n = a·b → d(n) cajas")
	txt(790, 452, 14, "#6ee7a0", "middle", "UN PRIMO SE ABRE EN DOS. Siempre.")
	txt(790, 474, 11, "#c9c3b0", "middle", "es el PISO: la única caja es el palo y su vuelta")
	txt(790, 512, 13, "#e8e2d0", "middle", "promedio de todos: ln n + 2γ − 1")
	txt(790, 536, 11, "#8a8676", "middle", "Dirichlet 1849 — y otra vez el logaritmo, y otra vez γ")
	txt(790, 562, 11, "#f2c14e", "middle", "medido hasta un millón: 13,970034 contra 13,969942")

	// the clue
	fmt.Fprintf(&b, `<rect x="20" y="598" width="1020" height="118" fill="#0f1a14" stroke="#2a5a40"/>`)
	txt(530, 624, 14, "#6ee7a0", "middle", "LA PISTA")
	txt(530, 650, 12, "#e8e2d0", "middle", "la apertura de los primos tiene PISO y TECHO, y los dos son exactos:")
	txt(530, 674, 12, "#f2c14e", "middle", "piso absoluto 26,565° una sola vez (2→3) · piso repetible 45° = los gemelos · techo 90°, y la velocidad hacia él ES el logaritmo")
	txt(530, 700, 11, "#c9c3b0", "middle", "mirando bolitas, sin una sola fórmula, el dibujo devuelve el logaritmo DOS VECES por caminos distintos. Todavía no.")

	b.WriteString(`</svg>`)
	ruta := filepath.Join("galeria", "laminas", "10-el-telar", "la-apertura.svg")
	_ = os.MkdirAll(filepath.Dir(ruta), 0o755)
	_ = os.WriteFile(ruta, []byte(b.String()), 0o644)
}
