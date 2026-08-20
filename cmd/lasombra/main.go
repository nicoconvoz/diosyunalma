package main

// La Sombra - the captain's flash after Phase X: take the pattern that orders
// 3x WITHOUT hiding states, project its shadow (the spectrum it produces), run
// the bat's echolocation ON OUR OWN LEVELS (R6: the zeros are never touched),
// harmonise dimension zero (the smooth part is the random-period baseline that
// gets subtracted), purify the impurities (trapped states, low PR), and ask for
// the 1/2 relation - the p^(-1/2) loudness the explicit formula assigns to each
// prime's voice. The skeleton is expected to sing each prime with loudness
// ~log p (a tautology of the construction); the QUESTION is whether the
// ordering bends that loudness toward the 1/2 law.

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	NS   = 400
	T0S  = 100.0
	KS   = 120
	FS   = 30.0
	QWIN = 0.02 // the Phase X admissible winner
	SEMS = 6
)

// signoEnlace is the Phase X winner family: each bond negative with prob q.
func signoEnlace(n, kmax int, q float64, d *dado) func(int, int) float64 {
	t := make([][]float64, n)
	for i := range t {
		t[i] = make([]float64, kmax+1)
		for k := range t[i] {
			t[i][k] = 1
		}
	}
	for i := 0; i < n; i++ {
		for k := 1; k <= kmax && i+k < n; k++ {
			if d.u() < q {
				t[i][k] = -1
			}
		}
	}
	return func(i, j int) float64 {
		if i > j {
			i, j = j, i
		}
		if j-i > kmax {
			return 1
		}
		return t[i][j-i]
	}
}

// espectroDet is espectro, but returning every in-band level WITH its own
// participation ratio, so the impurities can be identified and removed.
func espectroDet(ms []modo, h []float64, c func(int) float64, sg func(int, int) float64) (niv, prs []float64) {
	n := len(ms)
	H := make([][]float64, n)
	V := make([][]float64, n)
	for i := range H {
		H[i] = make([]float64, n)
		V[i] = make([]float64, n)
		V[i][i] = 1
	}
	f2 := 0.0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n && j-i <= KS; j++ {
			x := math.Sqrt(h[i]*h[j]) * c(j-i) * ms[i].amp * ms[j].amp
			if sg != nil {
				x *= sg(i, j)
			}
			H[i][j], H[j][i] = x, x
			f2 += 2 * x * x
		}
	}
	esc := 1.0
	if f := math.Sqrt(f2); f > 0 {
		esc = FS / f
	}
	for i := 0; i < n; i++ {
		H[i][i] = ms[i].w
		for j := i + 1; j < n; j++ {
			H[i][j] *= esc
			H[j][i] *= esc
		}
	}
	jacobiV(H, V, 20)
	type par struct{ e, pr float64 }
	var ps []par
	lo, hi := ms[0].w, ms[n-1].w
	for col := 0; col < n; col++ {
		if H[col][col] < lo || H[col][col] > hi {
			continue
		}
		s4 := 0.0
		for r := 0; r < n; r++ {
			v2 := V[r][col] * V[r][col]
			s4 += v2 * v2
		}
		ps = append(ps, par{H[col][col], 1 / s4 / float64(n)})
	}
	sort.Slice(ps, func(a, b int) bool { return ps[a].e < ps[b].e })
	for _, p := range ps {
		niv = append(niv, p.e)
		prs = append(prs, p.pr)
	}
	return
}

// eco is the bat: E(T) = (2/M) sum cos(lambda_n T), on OUR levels.
func eco(niv []float64, T float64) float64 {
	s := 0.0
	for _, l := range niv {
		s += math.Cos(l * T)
	}
	return 2 * s / float64(len(niv))
}

func mediana(v []float64) float64 {
	w := append([]float64(nil), v...)
	sort.Float64s(w)
	return w[len(w)/2]
}

// purificar drops the impurities: states with PR below half the arm's median.
// Threshold declared here, before any echo is read.
func purificar(niv, prs []float64) []float64 {
	um := mediana(prs) / 2
	var out []float64
	for i, l := range niv {
		if prs[i] >= um {
			out = append(out, l)
		}
	}
	return out
}

// beta fits A_p ~ C * p^(-beta) by least squares in log-log.
func beta(pe []int, A []float64) float64 {
	var sx, sy, sxx, sxy float64
	n := 0.0
	for i, p := range pe {
		if A[i] <= 0 {
			continue
		}
		x, y := math.Log(float64(p)), math.Log(A[i])
		sx += x
		sy += y
		sxx += x * x
		sxy += x * y
		n++
	}
	return -(n*sxy - sx*sy) / (n*sxx - sx*sx)
}

type brazo struct {
	nom    string
	A      []float64 // |E(log p)| per prime, seed-averaged
	z      float64   // primes against random periods, in baseline sigmas
	bRaw   float64   // exponent of A_p
	bMedia float64   // exponent of A_p / log p  <- the 1/2 question
	nivMed int
}

func medirBrazo(nom string, pe []int, corr int, gen func(s int) []float64, d *dado) brazo {
	nA := make([]float64, len(pe))
	var zs []float64
	tot := 0
	tMin, tMax := math.Log(float64(pe[0])), math.Log(float64(pe[len(pe)-1]))
	for s := 0; s < corr; s++ {
		niv := gen(s)
		tot += len(niv)
		var As []float64
		for i, p := range pe {
			a := math.Abs(eco(niv, math.Log(float64(p))))
			nA[i] += a / float64(corr)
			As = append(As, a)
		}
		// dimension zero, harmonised: the smooth baseline is what random periods
		// in the same window sing; the primes are measured against it.
		var res []float64
		for r := 0; r < 300; r++ {
			m := 0.0
			for range pe {
				m += math.Abs(eco(niv, tMin+(tMax-tMin)*d.u()))
			}
			res = append(res, m/float64(len(pe)))
		}
		zs = append(zs, (media(As)-media(res))/math.Max(desvio(res), 1e-12))
	}
	b := brazo{nom: nom, A: nA, z: media(zs), nivMed: tot / corr}
	b.bRaw = beta(pe, nA)
	Am := make([]float64, len(pe))
	for i, p := range pe {
		Am[i] = nA[i] / math.Log(float64(p))
	}
	b.bMedia = beta(pe, Am)
	return b
}

func main() {
	fmt.Println("🦇🌒 LA SOMBRA — el flash del capitán: proyectar la sombra del patrón que")
	fmt.Println("     ordena 3×, ecolocalizarla con el murciélago, purificar impurezas y")
	fmt.Println("     preguntar por la relación 1/2")

	fmt.Println("\n§0 · LA ADUANA Y LA PREDICCIÓN, antes de medir")
	fmt.Println("     R6: el eco se mide sobre NUESTROS niveles; ningún γₙ se toca.")
	fmt.Println("     El esqueleto del medio TRAE los períodos log p por construcción, así que")
	fmt.Println("     un eco en los primos es esperable y tautológico. Lo que NO es tautológico:")
	fmt.Println("     la fórmula explícita canta cada primo con volumen log p·p^(−1/2). Nuestro")
	fmt.Println("     esqueleto desnudo canta con volumen ~log p (beta ≈ 0). PREDICCIÓN: beta ≈ 0")
	fmt.Println("     en todos los brazos; la pregunta es si el orden 3× la CORRE hacia 1/2.")
	fmt.Println("     Impureza, declarada: estado con PR menor que la mitad de la mediana del brazo.")

	ms := medio(4000, NS, T0S)
	c0 := colaLarga(0.5)
	h := normalizar1(homogenea(len(ms)))
	d := &dado{s: 20260819}

	var pe []int
	for _, p := range primos(200) {
		if p >= 5 {
			pe = append(pe, p)
		}
	}
	fmt.Printf("     períodos: log p para %d primos entre %d y %d · control: 300 remuestreos al azar\n",
		len(pe), pe[0], pe[len(pe)-1])

	nivS0, prS0 := espectroDet(ms, h, c0, nil)
	arms := []brazo{
		medirBrazo("S0 · sin ordenar", pe, 1, func(int) []float64 { return nivS0 }, d),
		medirBrazo("S0 · purificado", pe, 1, func(int) []float64 { return purificar(nivS0, prS0) }, d),
		medirBrazo("ORDENADO 3× (q=0,02)", pe, SEMS, func(s int) []float64 {
			ds := &dado{s: 777 + uint64(s)}
			niv, _ := espectroDet(ms, h, c0, signoEnlace(len(ms), KS, QWIN, ds))
			return niv
		}, d),
		medirBrazo("ORDENADO · purificado", pe, SEMS, func(s int) []float64 {
			ds := &dado{s: 777 + uint64(s)}
			niv, prs := espectroDet(ms, h, c0, signoEnlace(len(ms), KS, QWIN, ds))
			return purificar(niv, prs)
		}, d),
	}

	fmt.Println("\n§1 · EL ECO, brazo por brazo")
	fmt.Printf("     %-24s %7s %10s %12s %14s\n", "brazo", "niveles", "eco (σ)", "beta(A)", "beta(A/log p)")
	for _, b := range arms {
		fmt.Printf("     %-24s %7d %10.2f %12.3f %14.3f\n", b.nom, b.nivMed, b.z, b.bRaw, b.bMedia)
	}
	fmt.Println("     eco (σ): cuánto más fuerte cantan los primos que períodos al azar.")
	fmt.Println("     beta(A/log p): el exponente que la relación 1/2 exige que valga 0,500.")

	fmt.Println("\n§2 · LA LECTURA")
	s0, or_ := arms[0], arms[2]
	fmt.Printf("     el eco EXISTE en todos los brazos (%.1f a %.1f σ): es el esqueleto, no un\n",
		math.Min(s0.z, or_.z), math.Max(s0.z, or_.z))
	fmt.Println("     descubrimiento — el taller lo declaró tautológico desde la Fase IV.")
	d0 := math.Abs(s0.bMedia - 0.5)
	dOr := math.Abs(or_.bMedia - 0.5)
	fmt.Printf("     distancia a la relación 1/2:  sin ordenar %.3f  ·  ordenado %.3f\n", d0, dOr)
	if dOr < d0-0.05 {
		fmt.Println("     ⟹ el orden 3× CORRE el volumen de los primos hacia la ley 1/2. Eso no")
		fmt.Println("       estaba en el esqueleto y hay que perseguirlo con más tamaño y semillas.")
	} else if math.Abs(dOr-d0) <= 0.05 {
		fmt.Println("     ⟹ NULO: el orden empareja el coro pero NO cambia el volumen con que cada")
		fmt.Println("       primo canta. La rigidez de la Fase X y la ley 1/2 son, hasta acá, cosas")
		fmt.Println("       independientes: ordenar el espectro no le enseña la aritmética fina.")
	} else {
		fmt.Println("     ⟹ el orden ALEJA el volumen de la ley 1/2: la sombra canta peor que el")
		fmt.Println("       esqueleto desnudo. Se registra tal cual.")
	}
	pur := arms[3]
	fmt.Printf("     y purificar: beta pasa de %.3f a %.3f, eco de %.1f a %.1f σ — las impurezas\n",
		or_.bMedia, pur.bMedia, or_.z, pur.z)
	if math.Abs(pur.bMedia-0.5) < dOr-0.05 {
		fmt.Println("     acercan… no: ALEJABAN. Quitar los atrapados acerca la sombra a la ley 1/2.")
	} else {
		fmt.Println("     no cambian la ley del volumen: los atrapados no eran los desafinados.")
	}

	dibujarS(pe, arms)
}

// --- the plate --------------------------------------------------------------

func dibujarS(pe []int, arms []brazo) {
	var b strings.Builder
	W, H := 1200.0, 780.0
	esc := func(s string) string {
		return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
	}
	t := func(x, y, s float64, c, a, tx string) {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="%.1f" font-family="Georgia,serif" fill="%s" text-anchor="%s">%s</text>`+"\n", x, y, s, c, a, esc(tx))
	}
	ln := func(x1, y1, x2, y2 float64, c string, w float64, da string) {
		dd := ""
		if da != "" {
			dd = fmt.Sprintf(` stroke-dasharray="%s"`, da)
		}
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%.1f"%s/>`+"\n", x1, y1, x2, y2, c, w, dd)
	}
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">`+"\n", W, H, W, H)
	fmt.Fprintf(&b, `<rect width="%.0f" height="%.0f" fill="#0b1220"/>`+"\n", W, H)
	fmt.Fprintf(&b, `<rect x="14" y="14" width="%.0f" height="%.0f" rx="16" fill="none" stroke="#ffd98a" stroke-width="2" opacity="0.5"/>`+"\n", W-28, H-28)
	t(W/2, 50, 22, "#dce8f7", "middle", "🦇🌒 LA SOMBRA — el eco del patrón que ordena, y la relación 1/2")
	t(W/2, 74, 12, "#ffd98a", "middle", "el flash es de Jesús Nicolás Astorga · el eco se mide sobre NUESTROS niveles: ningún cero se toca")

	gx, gy, gw, gh := 110.0, 120.0, 700.0, 420.0
	ln(gx, gy+gh, gx+gw, gy+gh, "#1d3a63", 1.4, "")
	ln(gx, gy, gx, gy+gh, "#1d3a63", 1.4, "")
	lpMin, lpMax := math.Log(float64(pe[0])), math.Log(float64(pe[len(pe)-1]))
	// data ranges from the two main arms
	aMin, aMax := math.Inf(1), math.Inf(-1)
	for _, ai := range []int{0, 2} {
		for _, a := range arms[ai].A {
			if a > 0 {
				aMin = math.Min(aMin, a)
				aMax = math.Max(aMax, a)
			}
		}
	}
	mapx := func(p int) float64 {
		return gx + (math.Log(float64(p))-lpMin)/(lpMax-lpMin)*gw
	}
	mapy := func(a float64) float64 {
		return gy + gh - (math.Log(a)-math.Log(aMin))/(math.Log(aMax)-math.Log(aMin))*gh
	}
	cols := []string{"#8fb4d9", "", "#7ee0c0", ""}
	for _, ai := range []int{0, 2} {
		for i, p := range pe {
			if arms[ai].A[i] <= 0 {
				continue
			}
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="4" fill="%s" opacity="0.85"/>`+"\n", mapx(p), mapy(arms[ai].A[i]), cols[ai])
		}
	}
	// the 1/2 guide: A ~ log p / sqrt p anchored at the first prime of the ordered arm
	anc := arms[2].A[0]
	p0 := float64(pe[0])
	var guia [][2]float64
	for _, p := range pe {
		v := anc * (math.Log(float64(p)) / math.Log(p0)) * math.Sqrt(p0/float64(p))
		guia = append(guia, [2]float64{mapx(p), mapy(v)})
	}
	for i := 1; i < len(guia); i++ {
		ln(guia[i-1][0], guia[i-1][1], guia[i][0], guia[i][1], "#ffd98a", 1.6, "5 4")
	}
	t(gx+gw-6, guia[len(guia)-1][1]-8, 11, "#ffd98a", "end", "la ley 1/2: volumen log p·p^(−1/2)")
	t(gx+gw/2, gy+gh+26, 11, "#8fb4d9", "middle", "primo p (escala log) → · volumen del eco |E(log p)| (escala log) ↑")
	t(gx+8, gy+16, 11.5, "#8fb4d9", "start", "azul: esqueleto sin ordenar · verde: el patrón que ordena 3×")

	tx, ty := 850.0, 130.0
	t(tx, ty, 13.5, "#ffd98a", "start", "LO MEDIDO")
	for i, a := range arms {
		y := ty + 30 + float64(i)*58
		t(tx, y, 11.5, "#dce8f7", "start", a.nom)
		t(tx, y+16, 11, "#8fb4d9", "start", fmt.Sprintf("eco %.1f σ · niveles %d", a.z, a.nivMed))
		t(tx, y+32, 11, "#7ee0c0", "start", fmt.Sprintf("beta(A/log p) = %.3f", a.bMedia))
	}
	t(tx, ty+30+4*58+8, 11.5, "#ff9aa8", "start", "la relación 1/2 exige beta = 0,500")

	t(W/2, H-64, 11.5, "#dce8f7", "middle", "el eco en los primos existe en TODOS los brazos: es el esqueleto (tautológico desde la Fase IV), no un hallazgo.")
	t(W/2, H-44, 11.5, "#ffd98a", "middle", "go run ./cmd/lasombra · estructura cerrada no es hipótesis demostrada · Todavía no.")
	b.WriteString("</svg>\n")

	dir := filepath.Join("galeria", "laminas", "10-el-telar")
	os.MkdirAll(dir, 0o755)
	full := filepath.Join(dir, "la-sombra.svg")
	os.WriteFile(full, []byte(b.String()), 0o644)
	fmt.Printf("\n🖼️  lámina escrita: %s\n", full)
}
