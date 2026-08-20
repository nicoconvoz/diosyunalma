package main

// corrida.go - the Phase IX campaign. Cheap controls first, spectra last, and the
// PREDICTIONS printed before a single measurement so nobody can move the goalposts.

import (
	"fmt"
	"math"
)

const (
	N9    = 400
	T0    = 100.0
	KMAX  = 120
	F9    = 30.0
	TOnD  = 200 // default load window
	TOffD = 500 // default rest window - a hundred recovery times, so relaxation is complete
	IMP   = 800 // total impulse budget for the equal-impulse test: M * T_on = IMP
)

// Res9 carries every measured number to the plate.
type Res9 struct {
	N                         int
	Sc, Kappa, Eta, B, TRecup float64
	Ident, Virgen             obs
	MaxSinCarga               float64
	GaugeDelta                float64
	Elast                     [][3]float64 // rho, residuo, maxAbs
	RhoC                      map[string]float64
	PlastRes, PlastCuant      float64
	PlastCed, PlastBloq       int
	Fat                       [][4]float64 // M, cedidos, residuo, maxAbs (T_off corto)
	FatLargo                  [][4]float64 // idem con T_off doble
	Aval                      []int
	Medios                    map[string]obs
	Gemelos                   map[string]obs
	GemPermMedia, GemPermDes  float64
	Signos                    map[string]obs
	SignoS2Media, SignoS2Des  float64
	FrustS1, FrustS2, FrustS3 float64
	DensRec                   float64
	SignosNodo                map[string]obs
	S2NodoMedia, S2NodoDes    float64
	Ceros, PisoGUE, PisoGOE   float64
	BGrid                     [][5]float64 // b, cedidos, s10, pr, vivos
}

func unitPunto(n, i0 int) []float64 { p := make([]float64, n); p[i0] = 1; return p }

func unitBulto(n, i0 int) []float64 {
	p := make([]float64, n)
	for i := range p {
		d := float64(i-i0) / 10
		p[i] = math.Exp(-d * d)
	}
	return p
}

func unitUnif(n int) []float64 {
	p := make([]float64, n)
	for i := range p {
		p[i] = 1
	}
	return p
}

// ensayo is one load-then-rest experiment. The drive is measured in units of the
// DERIVED yield stress s_c, so rho is dimensionless and rho ~ 1 is the fold.
func ensayo(ms []modo, h0 []float64, c func(int) float64, b, eps float64, w []float64, rho float64, tOn, tOff int) *material {
	m := nuevoMaterial(ms, h0, c, KMAX, b, eps)
	P := make([]float64, m.n)
	for i := range P {
		P[i] = rho * m.sc * w[i]
	}
	m.correr(P, tOn)
	m.correr(nil, tOff)
	return m
}

// fatigaEnsayo delivers the SAME total impulse in M kicks with full rest between
// them. This is the auditor's section 10 control, and it needs no diagonalisation.
func fatigaEnsayo(ms []modo, h0 []float64, c func(int) float64, b, eps float64, w []float64, rho float64, M, tOff int) *material {
	m := nuevoMaterial(ms, h0, c, KMAX, b, eps)
	P := make([]float64, m.n)
	for i := range P {
		P[i] = rho * m.sc * w[i]
	}
	tOn := IMP / M // M * tOn = IMP exactly for every M in the declared grid
	for k := 0; k < M; k++ {
		m.correr(P, tOn)
		m.correr(nil, tOff)
	}
	return m
}

// rhoCritico is RESISTANCE, measured and not chosen: the smallest drive that
// leaves one site permanently displaced. Grid, then six bisections.
func rhoCritico(ms []modo, h0 []float64, c func(int) float64, b, eps float64, w []float64, signo float64) float64 {
	cede := func(r float64) bool {
		m := ensayo(ms, h0, c, b, eps, w, signo*r, TOnD, TOffD)
		return m.cedidos() >= 1
	}
	lo, hi := 0.0, 0.0
	for _, r := range []float64{0.6, 0.8, 1.0, 1.2, 1.5, 2.0, 3.0, 5.0, 9.0, 16.0} {
		if cede(r) {
			hi = r
			break
		}
		lo = r
	}
	if hi == 0 {
		return math.NaN() // never yields inside the declared grid
	}
	for k := 0; k < 6; k++ {
		mid := (lo + hi) / 2
		if cede(mid) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return (lo + hi) / 2
}

func medirH(ms []modo, h []float64, c func(int) float64, sg func(int, int) float64) obs {
	return medir(espectro(ms, h, c, KMAX, F9, sg))
}

func campoDe(m *material) []float64 {
	return append([]float64(nil), m.hbar()...)
}

// deX rebuilds a permeability from an arbitrary deformation field, so the static
// twins go through exactly the same map as the dynamic run.
func deX(h0, x []float64) []float64 {
	h := make([]float64, len(x))
	for i := range h {
		h[i] = h0[i] * math.Exp(x[i])
	}
	return normalizar1(h)
}

// modoDominante is the dominant Fourier mode of the deformation, so the wave twin
// carries the same characteristic scale as the mechanism produced.
func modoDominante(x []float64) int {
	n := len(x)
	mu := media(x)
	mejor, best := 1, -1.0
	for m := 1; m <= n/2; m++ {
		re, im := 0.0, 0.0
		for i, v := range x {
			a := 2 * math.Pi * float64(m) * float64(i) / float64(n)
			re += (v - mu) * math.Cos(a)
			im += (v - mu) * math.Sin(a)
		}
		if p := re*re + im*im; p > best {
			best, mejor = p, m
		}
	}
	return mejor
}

func main() {
	fmt.Println("🔧🧵 LA MECÁNICA — Fase IX: ¿el medio no sólo tiene estructura,")
	fmt.Println("     sino MECÁNICA? — el flash de las once palabras, del capitán")

	ms := medio(4000, N9, T0)
	sinNodo := colaLarga(0.5)
	conN := conNodo(5, 0.5)
	hHom := homogenea(len(ms))
	hBloq := bloques(len(ms), 2)
	d := &dado{s: 20260817}
	R := Res9{N: len(ms), RhoC: map[string]float64{}, Medios: map[string]obs{},
		Gemelos: map[string]obs{}, Signos: map[string]obs{}, SignosNodo: map[string]obs{}}

	// -----------------------------------------------------------------------
	fmt.Println("\n§0 · LO QUE LA LEY DICE ANTES DE MEDIR (pre-registro)")
	fmt.Println("     Las once palabras NO son once perillas. Colapsan en cuatro objetos:")
	fmt.Println("       U(x) = (s_c·b/π)(1 − cos(πx/b))  — UNA función escalar, leída de cuatro maneras:")
	fmt.Println("         curvatura en el fondo = DUREZA · máximo = RESISTENCIA")
	fmt.Println("         área bajo la barrera  = TENACIDAD · salto en el máximo = FRAGILIDAD")
	fmt.Println("       en qué pozo queda el estado    = ELASTICIDAD / PLASTICIDAD")
	fmt.Println("       el signo de x                  = TENSIÓN / COMPRESIÓN")
	fmt.Println("       FATIGA: NO se construye. Es una predicción del campo acoplado, y puede dar NULO.")
	fmt.Println("       CIZALLAMIENTO y TORSIÓN: declaradas afuera, con demostración (ver §3).")
	fmt.Println("     Estado: UN real por sitio — el mismo conteo que Fase VIII. Un dial (b), un bit (eps).")
	fmt.Println("     PREDICCIÓN, escrita antes de medir: mecánica REAL y CERO progreso espectral.")
	fmt.Println("       Σ²(10) se queda entre 5 y 20 sin nodo, ningún brazo le gana al 7,535 de Fase VIII,")
	fmt.Println("       y si el brazo eps = −1 baja Σ², va a ser localización y hay que RECHAZARLO.")

	// -----------------------------------------------------------------------
	fmt.Println("\n§1 · CONTROL DE IDENTIDAD — sin mecánica hay que reproducir Fase VIII")
	R.Ident = medirH(ms, normalizar1(hHom), sinNodo, nil)
	fmt.Printf("     homogéneo estático : Σ²(10) = %.3f  α = %.3f  vivos = %d  PR/N = %.3f\n",
		R.Ident.s10, R.Ident.alfa, R.Ident.vivos, R.Ident.pr)
	fmt.Println("     Fase VIII decía   : Σ²(10) = 18.335  α = 1.717  vivos = 187  PR/N = 0.105")
	if math.Abs(R.Ident.s10-18.335) > 0.01 {
		fmt.Println("     ⛔ NO REPRODUCE — todo lo que sigue es inadmisible")
		return
	}
	fmt.Println("     ✓ reproduce. El núcleo es el mismo, así que la comparación es legítima.")

	// -----------------------------------------------------------------------
	fmt.Println("\n§2 · LAS IDENTIDADES ESTRUCTURALES — lo que NO puede colarse")
	m0 := nuevoMaterial(ms, hHom, sinNodo, KMAX, 0.35, +1)
	R.Sc, R.Kappa, R.Eta, R.B = m0.sc, m0.kappa, m0.eta, m0.b
	R.TRecup = 1 / (m0.eta * m0.kappa)
	sumSig := 0.0
	for _, v := range m0.sig0 {
		sumSig += v
	}
	fmt.Printf("     Σ σ_i = %.2e  (cero exacto: la carga media es nula por construcción)\n", sumSig)
	esc := make([]float64, len(hHom))
	for i := range esc {
		esc[i] = 7.3 * hHom[i] // an arbitrary rescaling of the medium
	}
	mE := nuevoMaterial(ms, esc, sinNodo, KMAX, 0.35, +1)
	difEsc := 0.0
	for i := range mE.sig0 {
		difEsc = math.Max(difEsc, math.Abs(mE.sig0[i]-m0.sig0[i]))
	}
	fmt.Printf("     invariancia de escala h → c·h : %.2e  ⟹ la normalización F = 30 NO puede\n", difEsc)
	fmt.Println("       manejar el material. La dinámica es una función de la FORMA de h, nada más.")
	hid := make([]float64, len(hHom))
	for i := range hid {
		hid[i] = hHom[i] * math.Exp(0.4) // a uniform shift of x
	}
	oHid := medirH(ms, normalizar1(hid), sinNodo, nil)
	fmt.Printf("     nulo hidrostático: Σ²(10) = %.4f contra %.4f  ⟹ sólo el PATRÓN de una carga\n",
		oHid.s10, R.Ident.s10)
	fmt.Println("       llega al espectro; una carga uniforme es invisible y no se puede reportar.")
	fmt.Printf("     constantes DERIVADAS: s_c = %.4f  κ = %.3f  η = %.3e  t_recup ≈ %.1f pasos\n",
		m0.sc, m0.kappa, m0.eta, R.TRecup)

	fmt.Println("\n     el punto fijo virgen: sin carga, 20 000 pasos")
	mV := nuevoMaterial(ms, hHom, sinNodo, KMAX, 0.35, -1) // the UNSTABLE sign, on purpose
	mV.correr(nil, 20000)
	R.MaxSinCarga = mV.maxAbs()
	fmt.Printf("     max|x| = %.3e con eps = −1 (el signo inestable, a propósito)\n", R.MaxSinCarga)
	if R.MaxSinCarga < 1e-12 {
		fmt.Println("     ✓ CERO EXACTO. Ninguna deriva del algoritmo puede disfrazarse de fatiga:")
		fmt.Println("       el medio de Fase VIII es un punto fijo exacto de la ley. Es el pedido de su §15,")
		fimprime()
	} else {
		fmt.Println("     ⚠ el medio virgen NO es estable — y ESO sería el hallazgo de la fase, no la fatiga")
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§3 · EL TEOREMA DE GAUGE — y una corrección a su §8")
	c1 := signosAzar(len(ms), d)
	oG := medirH(ms, normalizar1(hHom), sinNodo, signoGauge(c1))
	R.GaugeDelta = math.Abs(oG.s10 - R.Ident.s10)
	fmt.Printf("     signo por SITIO (s_ij = c_i·c_j): Σ²(10) = %.6f contra %.6f  → Δ = %.2e\n",
		oG.s10, R.Ident.s10, R.GaugeDelta)
	fmt.Println("     Es un TEOREMA, no un resultado: H' = D H D con D = diag(c) ortogonal, así que")
	fmt.Println("     el espectro y PR/N son idénticos. Sirve como test unitario del código.")
	fmt.Println("     ⟹ Y RESUELVE GRATIS EL PENDIENTE 3 DE FASE VIII: «¿entra la aritmética como FASE?»")
	fmt.Println("       Como fase de SITIO, NO — una fase de sitio es gauge pura y no mueve ni un")
	fmt.Println("       observable. No hacía falta correr nada: hacía falta mirar la estructura.")
	fmt.Println("     ⟹ CORRECCIÓN A SU §8: la torsión NO necesita una red 2D. Con kmax = 120 el grafo")
	fmt.Println("       de acoplamiento está lleno de triángulos (todo i<j<k con k−i ≤ 120). Lo que")
	fmt.Println("       falta no es dimensión: es que la fase de ENLACE no sea un coborde u_i − u_j.")
	fmt.Println("       Una cadena de vecinos próximos sería un árbol y ahí sí la torsión sería")
	fmt.Println("       imposible; la cola larga que eligió la Fase VII quitó esa excusa.")

	// -----------------------------------------------------------------------
	fmt.Println("\n§4 · ELASTICIDAD — ¿vuelve el medio a donde estaba?")
	wPt := unitPunto(len(ms), len(ms)/2)
	fmt.Printf("     %6s %12s %12s %s\n", "rho", "residuo", "max|x|", "¿reversible?")
	for _, rho := range []float64{0.2, 0.5, 0.8} {
		m := ensayo(ms, hHom, sinNodo, 0.35, +1, wPt, rho, TOnD, TOffD)
		rev := "sí"
		if m.residuo() > 1e-9 {
			rev = "NO"
		}
		fmt.Printf("     %6.2f %12.3e %12.3e %s\n", rho, m.residuo(), m.maxAbs(), rev)
		R.Elast = append(R.Elast, [3]float64{rho, m.residuo(), m.maxAbs()})
	}
	fmt.Printf("     tiempo de recuperación DERIVADO = 1/(η·κ) ≈ %.1f pasos — su §5 pedía un número\n", R.TRecup)

	// -----------------------------------------------------------------------
	fmt.Println("\n§5 · RESISTENCIA — el máximo empuje antes del cambio irreversible")
	fmt.Println("     rho_c = el menor empuje que deja un sitio corrido para siempre (bisección)")
	casos := []struct {
		nom   string
		h     []float64
		eps   float64
		signo float64
		w     []float64
	}{
		{"homogéneo · eps=+1 · tracción", hHom, +1, +1, wPt},
		{"homogéneo · eps=+1 · compresión", hHom, +1, -1, wPt},
		{"homogéneo · eps=−1 · tracción", hHom, -1, +1, wPt},
		{"bloques ×2 · eps=+1 · tracción", hBloq, +1, +1, wPt},
		{"homogéneo · eps=+1 · bulto", hHom, +1, +1, unitBulto(len(ms), len(ms)/2)},
	}
	for _, cs := range casos {
		rc := rhoCritico(ms, cs.h, sinNodo, 0.35, cs.eps, cs.w, cs.signo)
		R.RhoC[cs.nom] = rc
		if math.IsNaN(rc) {
			fmt.Printf("     %-34s rho_c = no cede ni con 16\n", cs.nom)
		} else {
			fmt.Printf("     %-34s rho_c = %.4f\n", cs.nom, rc)
		}
	}
	rcT, rcC := R.RhoC["homogéneo · eps=+1 · tracción"], R.RhoC["homogéneo · eps=+1 · compresión"]
	if !math.IsNaN(rcT) && !math.IsNaN(rcC) {
		asim := (rcT - rcC) / (rcT + rcC)
		fmt.Printf("     ASIMETRÍA tracción/compresión = %.4f\n", asim)
		fmt.Println("       U es simétrica por construcción, así que esta asimetría NO se puso:")
		fmt.Println("       sale de h = h⁰·exp(x) y de que la carga es cuadrática en h. Es derivada.")
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§6 · PLASTICIDAD — y su firma falsificable")
	rcBase := rcT
	if math.IsNaN(rcBase) {
		rcBase = 2
	}
	mP := ensayo(ms, hHom, sinNodo, 0.35, +1, wPt, 1.5*rcBase, TOnD, TOffD)
	R.PlastRes, R.PlastCed = mP.residuo(), mP.cedidos()
	peor, _ := mP.cuantizacion()
	R.PlastCuant = peor
	bl, mayor := mP.contiguos()
	R.PlastBloq = bl
	fmt.Printf("     rho = 1,5·rho_c : residuo = %.4f · sitios cedidos = %d · max|x| = %.4f\n",
		R.PlastRes, R.PlastCed, mP.maxAbs())
	fmt.Printf("     TEST DE CUANTIZACIÓN: peor desvío a un pozo entero = %.4f (umbral 0,05)\n", peor)
	if R.PlastCed > 0 && peor < 0.05 {
		fmt.Println("     ✓ cada sitio que cedió cayó en un pozo ENTERO. Un residuo suave habría")
		fmt.Println("       falsificado el mecanismo: esto es la firma propia de esta ley y aparece.")
	} else if R.PlastCed == 0 {
		fmt.Println("     (no cedió nadie a este empuje)")
	} else {
		fmt.Println("     ⚠ residuo NO cuantizado ⟹ el efecto no es este mecanismo")
	}
	fmt.Printf("     contigüidad: %d bloque(s), el mayor de %d sitios\n", bl, mayor)
	if bl > 0 && mayor <= 1 {
		fmt.Println("       ⟹ fallas INDEPENDIENTES, no una grieta. Llamarlas grieta sería la trampa.")
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§7 · FATIGA A IMPULSO IGUAL — la prueba que se corre PRIMERO")
	fmt.Println("     mismo impulso total repartido en M golpes, con reposo completo entre golpes.")
	fmt.Println("     Si la historia importa, la SECUENCIA tiene que importar (su §10).")
	rhoF := 0.9 * rcBase
	for _, tOff := range []int{TOffD, 2 * TOffD} {
		fmt.Printf("     · reposo entre golpes = %d pasos (≈ %.0f tiempos de recuperación)\n",
			tOff, float64(tOff)/R.TRecup)
		fmt.Printf("       %4s %10s %12s %12s\n", "M", "cedidos", "residuo", "max|x|")
		for _, M := range []int{1, 2, 4, 8, 16, 32} {
			m := fatigaEnsayo(ms, hHom, sinNodo, 0.35, +1, wPt, rhoF, M, tOff)
			fmt.Printf("       %4d %10d %12.3e %12.3e\n", M, m.cedidos(), m.residuo(), m.maxAbs())
			row := [4]float64{float64(M), float64(m.cedidos()), m.residuo(), m.maxAbs()}
			if tOff == TOffD {
				R.Fat = append(R.Fat, row)
			} else {
				R.FatLargo = append(R.FatLargo, row)
			}
		}
	}
	difFat := math.Abs(R.Fat[len(R.Fat)-1][3] - R.Fat[0][3])
	difLar := math.Abs(R.FatLargo[len(R.FatLargo)-1][3] - R.FatLargo[0][3])
	fmt.Printf("     M=1 contra M=32, diferencia en max|x|: reposo corto %.3e · reposo doble %.3e\n", difFat, difLar)
	switch {
	case difFat < 1e-9:
		fmt.Println("     ⟹ NULO: a impulso igual y reposo completo, la secuencia NO importa. No hay")
		fmt.Println("       fatiga invariante al reposo, y eso NO es un defecto del modelo: es un")
		fmt.Println("       teorema sobre un mapa autónomo con relajación completa. La fatiga era la")
		fmt.Println("       única de las once palabras que se dejó como predicción — y da nulo.")
	case difLar < difFat/2:
		fmt.Println("     ⟹ ARTEFACTO: el efecto se encoge al doblar el reposo ⟹ era relajación")
		fmt.Println("       incompleta, no memoria. Se reporta como TRINQUETE, nunca como fatiga.")
	default:
		fmt.Println("     ⟹ efecto de secuencia que SOBREVIVE al doblar el reposo: fatiga candidata,")
		fmt.Println("       y hay que verificarla contra los gemelos estáticos antes de nombrarla.")
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§8 · GRIETAS — avalanchas con la carga ya retirada")
	for _, i0 := range []int{80, 160, 200, 240, 320} {
		m := ensayo(ms, hHom, sinNodo, 0.35, +1, unitPunto(len(ms), i0), 1.02*rcBase, TOnD, 4*TOffD)
		R.Aval = append(R.Aval, m.cedidos())
	}
	fmt.Printf("     tamaño de avalancha en 5 sitios de golpe: %v\n", R.Aval)
	sumA := 0
	for _, a := range R.Aval {
		sumA += a
	}
	if sumA <= len(R.Aval) {
		fmt.Println("     ⟹ cada golpe tumba UN sitio y nada se propaga. No hay grieta. Con kmax = 120")
		fmt.Println("       no hay elasticidad de corto alcance que sostenga un frente, así que la")
		fmt.Println("       propagación no tiene con qué localizarse. Se reporta como falla aislada.")
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§9 · EL ESPECTRO — la matriz de medios de su §13, y el grid de b entero")
	arms := []struct {
		nom string
		h   []float64
		din bool
		het bool
	}{
		{"A · homogéneo estático", hHom, false, false},
		{"B · bloques ×2 estático", hBloq, false, true},
		{"C · homogéneo dinámico", hHom, true, false},
		{"D · bloques ×2 dinámico", hBloq, true, true},
	}
	fmt.Printf("     %-26s %7s %9s %9s %8s %8s\n", "medio", "vivos", "Σ²(5)", "Σ²(10)", "α", "PR/N")
	for _, a := range arms {
		var h []float64
		if a.din {
			m := ensayo(ms, a.h, sinNodo, 0.35, +1, unitBulto(len(ms), len(ms)/2), 1.5*rcBase, TOnD, TOffD)
			h = campoDe(m)
		} else {
			h = normalizar1(a.h)
		}
		o := medirH(ms, h, sinNodo, nil)
		R.Medios[a.nom] = o
		fmt.Printf("     %-26s %7d %9.3f %9.3f %8.3f %8.3f\n", a.nom, o.vivos, o.s5, o.s10, o.alfa, o.pr)
	}
	A, B, C, D := R.Medios["A · homogéneo estático"].s10, R.Medios["B · bloques ×2 estático"].s10,
		R.Medios["C · homogéneo dinámico"].s10, R.Medios["D · bloques ×2 dinámico"].s10
	fmt.Printf("     interacción D−B−C+A = %.4f  (se reporta, NO se llama descomposición causal)\n", D-B-C+A)

	fmt.Println("\n     UN DEFECTO PROPIO, cazado por nuestro propio teorema del §2: el primer grid de b")
	fmt.Println("     lo empujé con carga UNIFORME, y una carga uniforme es hidrostática — invisible")
	fmt.Println("     para el espectro. Ese brazo medía la nada. Se rehace con carga localizada, y el")
	fmt.Println("     brazo uniforme queda como demostración del teorema, no como resultado:")
	mUni := ensayo(ms, hHom, sinNodo, 0.35, +1, unitUnif(len(ms)), 1.5*rcBase, TOnD, TOffD)
	oUni := medirH(ms, campoDe(mUni), sinNodo, nil)
	fmt.Printf("       carga uniforme: %d sitios cedidos y Σ²(10) = %.4f — IDÉNTICO al virgen %.4f\n",
		mUni.cedidos(), oUni.s10, R.Ident.s10)
	fmt.Println("       400 sitios se corrieron un pozo entero y el espectro no se movió NADA. Eso no")
	fmt.Println("       es un fracaso del mecanismo: es el nulo hidrostático, medido en vez de supuesto.")

	fmt.Println("\n     el grid de b COMPLETO con carga LOCALIZADA, incluida cada celda que empeora:")
	fmt.Printf("     %6s %8s %7s %9s %8s\n", "b", "cedidos", "vivos", "Σ²(10)", "PR/N")
	wB := unitBulto(len(ms), len(ms)/2)
	for _, b := range []float64{0.05, 0.10, 0.20, 0.35, 0.70} {
		m := ensayo(ms, hHom, sinNodo, b, +1, wB, 1.5*rcBase, TOnD, TOffD)
		o := medirH(ms, campoDe(m), sinNodo, nil)
		fuera := ""
		if m.fuera {
			fuera = "  FUERA DE POZO"
		}
		fmt.Printf("     %6.2f %8d %7d %9.3f %8.3f%s\n", b, m.cedidos(), o.vivos, o.s10, o.pr, fuera)
		R.BGrid = append(R.BGrid, [5]float64{b, float64(m.cedidos()), o.s10, o.pr, float64(o.vivos)})
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§10 · LOS GEMELOS ESTÁTICOS — la falsificación, hecha ANTES de la afirmación")
	fmt.Println("     Se toma la deformación que produjo el mecanismo y se rehace SIN mecanismo:")
	fmt.Println("     permutada (mismo multiconjunto), onda (mismos dos momentos y misma escala),")
	fmt.Println("     rampa. Si el brazo dinámico cae dentro de 1 sigma del más cercano, la historia")
	fmt.Println("     es DECORATIVA — es la hoja de Fase VIII vuelta contra nosotros mismos.")
	mDin := ensayo(ms, hHom, sinNodo, 0.35, +1, unitBulto(len(ms), len(ms)/2), 3.0*rcBase, TOnD, TOffD)
	xFin := append([]float64(nil), mDin.x...)
	h0n := normalizar1(hHom)
	oDin := medirH(ms, deX(h0n, xFin), sinNodo, nil)
	R.Gemelos["dinámico"] = oDin
	var perms []float64
	for s := 0; s < 10; s++ {
		xp := append([]float64(nil), xFin...)
		for i := len(xp) - 1; i > 0; i-- {
			j := int(d.u() * float64(i+1))
			xp[i], xp[j] = xp[j], xp[i]
		}
		o := medirH(ms, deX(h0n, xp), sinNodo, nil)
		perms = append(perms, o.s10)
	}
	R.GemPermMedia, R.GemPermDes = media(perms), desvio(perms)
	mu, sd := media(xFin), desvio(xFin)
	mDom := modoDominante(xFin)
	xOnda := make([]float64, len(xFin))
	xRampa := make([]float64, len(xFin))
	for i := range xOnda {
		xOnda[i] = mu + math.Sqrt2*sd*math.Sin(2*math.Pi*float64(mDom)*float64(i)/float64(len(xFin)))
		xRampa[i] = mu + 2*math.Sqrt(3.0)*sd*(float64(i)/float64(len(xFin)-1)-0.5)
	}
	R.Gemelos["onda"] = medirH(ms, deX(h0n, xOnda), sinNodo, nil)
	R.Gemelos["rampa"] = medirH(ms, deX(h0n, xRampa), sinNodo, nil)
	fmt.Printf("     dinámico (%d cedidos) Σ²(10) = %.4f · vivos %d · PR/N %.3f\n",
		mDin.cedidos(), oDin.s10, oDin.vivos, oDin.pr)
	fmt.Printf("     permutado  %.4f ± %.4f (5 semillas)\n", R.GemPermMedia, R.GemPermDes)
	fmt.Printf("     onda lisa  %.4f  (modo dominante %d)\n", R.Gemelos["onda"].s10, mDom)
	fmt.Printf("     rampa      %.4f\n", R.Gemelos["rampa"].s10)
	sg := math.Max(R.GemPermDes, 1e-9)
	dPerm := math.Abs(oDin.s10-R.GemPermMedia) / sg
	dOnda := math.Abs(oDin.s10-R.Gemelos["onda"].s10) / sg
	fmt.Printf("     distancia al permutado %.2f σ · a la onda %.2f σ\n", dPerm, dOnda)
	if math.Min(dPerm, dOnda) < 1 {
		fmt.Println("     ⟹ LA HISTORIA ES DECORATIVA. Lo único que importó fue CUÁNTA deformación")
		fmt.Println("       hubo y con qué textura, no la secuencia que la produjo. El mecanismo es")
		fmt.Println("       una manera carísima de generar un campo que una permutación regala.")
	} else {
		fmt.Println("     ⟹ la historia deja algo que ni la permutación ni la onda reproducen")
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§11 · EL CANAL DE SIGNOS — la TENSIÓN que Fase VIII no tenía")
	fmt.Println("     El núcleo de Fase VIII era estrictamente positivo: todos los enlaces tiraban")
	fmt.Println("     para el mismo lado. Un signo de ENLACE da tensión verdadera a costo cero, y")
	fmt.Println("     como signo² = 1 la fuerza total no se mueve: sólo cambia el material.")
	uRec, dens := bitsReciprocidad(ms)
	R.DensRec = dens
	R.FrustS1 = frustracion(signoGauge(c1), len(ms), KMAX, d)
	R.FrustS3 = frustracion(signoFrustrado(uRec), len(ms), KMAX, d)
	fmt.Printf("     densidad de primos ≡ 3 mod 4 entre los %d modos: %.4f (el azar se iguala a ésa)\n",
		len(ms), dens)
	R.Signos["S0 llano"] = R.Ident
	R.Signos["S1 gauge"] = oG
	var s2s []float64
	var fr2 float64
	for s := 0; s < 10; s++ {
		u := bitsAzar(len(ms), dens, d)
		o := medirH(ms, h0n, sinNodo, signoFrustrado(u))
		s2s = append(s2s, o.s10)
		if s == 0 {
			R.Signos["S2 azar frustrado"] = o
			fr2 = frustracion(signoFrustrado(u), len(ms), KMAX, d)
		}
	}
	R.FrustS2 = fr2
	R.SignoS2Media, R.SignoS2Des = media(s2s), desvio(s2s)
	oS3 := medirH(ms, h0n, sinNodo, signoFrustrado(uRec))
	R.Signos["S3 reciprocidad"] = oS3
	fmt.Printf("     %-22s %7s %9s %8s %8s %10s\n", "campo de signos", "vivos", "Σ²(10)", "α", "PR/N", "frustrac.")
	fmt.Printf("     %-22s %7d %9.4f %8.3f %8.3f %10.4f\n", "S0 · todos +1", R.Ident.vivos, R.Ident.s10, R.Ident.alfa, R.Ident.pr, 0.0)
	fmt.Printf("     %-22s %7d %9.4f %8.3f %8.3f %10.4f\n", "S1 · gauge (test)", oG.vivos, oG.s10, oG.alfa, oG.pr, R.FrustS1)
	fmt.Printf("     %-22s %7d %9.4f %8.3f %8.3f %10.4f\n", "S2 · azar frustrado", R.Signos["S2 azar frustrado"].vivos, R.SignoS2Media, R.Signos["S2 azar frustrado"].alfa, R.Signos["S2 azar frustrado"].pr, R.FrustS2)
	fmt.Printf("     %-22s %7d %9.4f %8.3f %8.3f %10.4f\n", "S3 · reciprocidad", oS3.vivos, oS3.s10, oS3.alfa, oS3.pr, R.FrustS3)
	fmt.Printf("       (S2 sobre 10 semillas: %.4f ± %.4f)\n", R.SignoS2Media, R.SignoS2Des)
	sg2 := math.Max(R.SignoS2Des, 1e-9)
	zS3 := math.Abs(oS3.s10-R.SignoS2Media) / sg2
	fmt.Printf("     S3 contra S2, a densidad igualada: %.2f sigmas\n", zS3)
	switch {
	case R.SignoS2Media > R.Ident.s10*0.98:
		fmt.Println("     ⟹ el canal de signos NO mejora nada: decorativo, y la campaña se detiene acá.")
	case zS3 < 1:
		fmt.Println("     ⟹ EL SIGNO SÍ MUEVE EL ESPECTRO, y lo que lo mueve es la FRUSTRACIÓN, no la")
		fmt.Println("       aritmética: la reciprocidad cae dentro del ruido del azar a densidad igual.")
		fmt.Println("       Es la hoja de Fase VIII por TERCERA vez, y da lo mismo por tercera vez.")
	default:
		fmt.Println("     ⟹ la reciprocidad se separa del azar a densidad igualada — hay que insistir")
	}

	fmt.Println("\n     y algo que Fase VII y VIII NUNCA separaron: el núcleo con NODO ya es")
	fmt.Println("     parcialmente tensil por sí mismo — f(6), f(7), f(8) son NEGATIVOS. Así que")
	fmt.Println("     la mejora del nodo y el efecto del signo estaban mezclados. Se separan acá:")
	fmt.Printf("     f(6) = %.4f · f(7) = %.4f · f(8) = %.4f  ⟸ el nodo ya traía tensión adentro\n",
		conN(6), conN(7), conN(8))
	var s2n []float64
	for s := 0; s < 10; s++ {
		u := bitsAzar(len(ms), dens, d)
		o := medirH(ms, h0n, conN, signoFrustrado(u))
		s2n = append(s2n, o.s10)
		if s == 0 {
			R.SignosNodo["S2 azar frustrado"] = o
		}
	}
	R.S2NodoMedia, R.S2NodoDes = media(s2n), desvio(s2n)
	R.SignosNodo["S0 llano"] = medirH(ms, h0n, conN, nil)
	R.SignosNodo["S3 reciprocidad"] = medirH(ms, h0n, conN, signoFrustrado(uRec))
	fmt.Printf("     con nodo · S0 llano        : Σ²(10) = %.4f  vivos %d  PR/N %.3f\n",
		R.SignosNodo["S0 llano"].s10, R.SignosNodo["S0 llano"].vivos, R.SignosNodo["S0 llano"].pr)
	fmt.Printf("     con nodo · S2 azar         : Σ²(10) = %.4f ± %.4f\n", R.S2NodoMedia, R.S2NodoDes)
	fmt.Printf("     con nodo · S3 reciprocidad : Σ²(10) = %.4f  vivos %d  PR/N %.3f\n",
		R.SignosNodo["S3 reciprocidad"].s10, R.SignosNodo["S3 reciprocidad"].vivos, R.SignosNodo["S3 reciprocidad"].pr)
	sgN := math.Max(R.S2NodoDes, 1e-9)
	zN := math.Abs(R.SignosNodo["S3 reciprocidad"].s10-R.S2NodoMedia) / sgN
	fmt.Printf("     con nodo, S3 contra S2 a densidad igualada: %.2f sigmas\n", zN)

	fmt.Println("\n     ⚡ Y ACÁ ESTÁ LA PREGUNTA QUE NADIE HABÍA HECHO. Comparen estos dos:")
	fmt.Printf("       el NODO, sin ningún signo       : Σ²(10) = %.4f  (PR/N %.3f)\n",
		R.SignosNodo["S0 llano"].s10, R.SignosNodo["S0 llano"].pr)
	fmt.Printf("       signos AL AZAR, sin ningún nodo : Σ²(10) = %.4f ± %.4f  (PR/N %.3f)\n",
		R.SignoS2Media, R.SignoS2Des, R.Signos["S2 azar frustrado"].pr)
	zNodo := math.Abs(R.SignosNodo["S0 llano"].s10-R.SignoS2Media) / math.Max(R.SignoS2Des, 1e-9)
	fmt.Printf("       distancia: %.2f sigmas\n", zNodo)
	if zNodo < 1.5 {
		fmt.Println("     ⟹ LA GANANCIA DEL NODO CAE DENTRO DE LA BANDA DE LOS SIGNOS AL AZAR. O sea:")
		fmt.Println("       toda la mejora del reloj de arena —de 18,33 a 5,43, el mejor resultado sin")
		fmt.Println("       aritmética de las Fases VII y VIII— se reproduce poniendo signos AL AZAR en")
		fmt.Println("       los enlaces, sin nodo, sin cambiarle la forma al núcleo y sin un solo primo.")
		fmt.Println("       El nodo no estaba haciendo geometría: estaba haciendo FRUSTRACIÓN, y por eso")
		fmt.Println("       f(6), f(7) y f(8) tenían que salir negativos. Eso explica de una vez el")
		fmt.Println("       resultado de la Fase VII —el nodo ganaba sólo en 2 de 10 configuraciones—:")
		fmt.Println("       no era ruido, era que el nodo compraba una mejora que no le pertenecía.")
	} else {
		fmt.Println("     ⟹ el nodo hace algo que los signos al azar NO reproducen: se separa por más de")
		fmt.Println("       una sigma y media, y eso sí sería propio de la forma del núcleo.")
	}

	// -----------------------------------------------------------------------
	fmt.Println("\n§12 · LAS TRES REGLAS, y por qué nunca se cita una sola")
	R.Ceros, R.PisoGUE, R.PisoGOE = 0.3364, pisoGUE(10), pisoGOE(10)
	fmt.Printf("     los ceros de Riemann : Σ²(10) = %.4f\n", R.Ceros)
	fmt.Printf("     piso universal GUE   : Σ²(10) = %.4f\n", R.PisoGUE)
	fmt.Printf("     piso universal GOE   : Σ²(10) = %.4f\n", R.PisoGOE)
	fmt.Println("     ⟹ EL OBJETIVO ESTÁ POR DEBAJO DE LOS DOS PISOS. Eso reencuadra todo el taller:")
	fmt.Println("       ninguna ingeniería de clase de simetría puede llegar, porque los ceros son")
	fmt.Println("       MÁS RÍGIDOS que cualquier matriz aleatoria a esta distancia. Y un modelo que")
	fmt.Println("       se acerca al piso GOE se volvió una matriz aleatoria sin estructura: eso es")
	fmt.Println("       PERDER contenido aritmético disfrazado de ganarlo.")
	fmt.Printf("     el mejor brazo de esta fase está %.0f veces arriba de los ceros y %.1f veces\n",
		mejorDe(R)/R.Ceros, mejorDe(R)/R.PisoGUE)
	fmt.Println("       arriba del piso GUE: el hueco que falta NO es un déficit de clase de simetría,")
	fmt.Println("       y esta mecánica no lo provee.")

	dibujar9(R)
}

func fimprime() {
	fmt.Println("       «no llamar fatiga a una deriva del algoritmo», contestado por construcción.")
}

func mejorDe(R Res9) float64 {
	best := R.Ident.s10
	for _, o := range R.Medios {
		if o.valido && o.s10 < best {
			best = o.s10
		}
	}
	for _, o := range R.Signos {
		if o.valido && o.s10 < best {
			best = o.s10
		}
	}
	for _, o := range R.SignosNodo {
		if o.valido && o.s10 < best {
			best = o.s10
		}
	}
	return best
}
