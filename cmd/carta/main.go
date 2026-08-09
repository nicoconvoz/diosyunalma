// Command carta draws the full atlas of the standing hunt: every beast
// ever signed by the judge (luz/cazadero.log), laid out water by water,
// plus THE SHAPE — the coherence spectrum of the sea, where the catch
// reveals its anatomy: two branches (waves above, islands below) split
// by the calm water the sonar never wakes for.
package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type beast struct {
	t, k, ola, coh, juez float64
	L                    int64
	frac                 float64
	muda                 bool
}

type school struct {
	t, k, frac float64
	presas     int
}

func main() {
	data, err := os.ReadFile("luz/cazadero.log")
	if err != nil {
		fmt.Println("no hay libro de caza")
		return
	}
	var beasts []beast
	var schools []school
	for _, ln := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(ln, "CARDUMEN") {
			var tt, kk float64
			var pr int
			if n, _ := fmt.Sscanf(ln, "CARDUMEN: t=%g k=%g presas=%d", &tt, &kk, &pr); n == 3 {
				nT := math.Sqrt(tt / (2 * math.Pi))
				schools = append(schools, school{tt, kk, kk / nT, pr})
			}
			continue
		}
		var kind string
		var tt, kk, ola, coh, juez float64
		var LL int64
		n, _ := fmt.Sscanf(ln, "BESTIA %s t=%g k=%g L=%d |ola|=%g coh=%gσ juez=%g", &kind, &tt, &kk, &LL, &ola, &coh, &juez)
		if n < 6 {
			continue
		}
		nT := math.Sqrt(tt / (2 * math.Pi))
		beasts = append(beasts, beast{tt, kk, ola, coh, juez, LL, kk / nT, strings.HasPrefix(kind, "MUDA")})
	}
	// waters present, sorted
	wset := map[float64]bool{}
	for _, be := range beasts {
		wset[be.t] = true
	}
	waters := make([]float64, 0, len(wset))
	for w := range wset {
		waters = append(waters, w)
	}
	sort.Float64s(waters)

	var b strings.Builder
	laneH, laneGap := 64.0, 40.0
	base := 120.0
	shapeY := base + float64(len(waters))*(laneH+laneGap) + 30
	shapeH := 300.0
	H := shapeY + shapeH + 150
	W := 1460.0
	x0, plotW := 130.0, 1240.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<defs><linearGradient id="mar" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#123564"/><stop offset="1" stop-color="#0a1f42"/></linearGradient></defs>
<text x="%.0f" y="44" font-size="26" text-anchor="middle" font-family="Georgia" fill="#dce8f7">🧭 LA CARTA DEL ABISMO — atlas completo de la cacería</text>
<text x="%.0f" y="70" font-size="14" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">%d bestias firmadas por el juez · %d cardúmenes · %d aguas anexadas · eje horizontal: banda k/nTop (0 → 1)</text>`,
		W, H, W, H, W/2, W/2, len(beasts), len(schools), len(waters))

	// color per water: shallow gold -> deep violet
	wcol := func(t float64) string {
		i := 0
		for j, w := range waters {
			if w == t {
				i = j
			}
		}
		f := float64(i) / math.Max(float64(len(waters)-1), 1)
		r := int(255 - 130*f)
		g := int(209 - 150*f)
		bl := int(102 + 153*f)
		return fmt.Sprintf("#%02x%02x%02x", r, g, bl)
	}

	// ---- the lanes: one per water, every catch plotted ----
	for i, w := range waters {
		y := base + float64(i)*(laneH+laneGap)
		col := wcol(w)
		fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="8" fill="url(#mar)" stroke="#44608c"/>`, x0, y, plotW, laneH)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="16" text-anchor="end" font-family="Georgia" fill="%s">t=%.0e</text>`, x0-12, y+laneH/2+5, col, w)
		nW, nI := 0, 0
		for _, be := range beasts {
			if be.t != w {
				continue
			}
			x := x0 + plotW*be.frac
			if be.muda {
				nI++
				fmt.Fprintf(&b, `<path d="M %.1f %.1f l 4 -9 l 4 9 z" fill="#7fd7a8" opacity="0.85"/>`, x-4, y+laneH-8)
			} else {
				nW++
				r := math.Max(2, 1.5+2.2*(be.coh-2.2))
				fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#7fb2ff" opacity="0.75"/>`, x, y+laneH*0.35, r)
			}
		}
		// dedupe schools by band bucket so the pod stays readable
		seen := map[int]bool{}
		for _, s := range schools {
			if s.t == w {
				bkt := int(s.frac * 200)
				if seen[bkt] {
					continue
				}
				seen[bkt] = true
				fmt.Fprintf(&b, `<text x="%.1f" y="%.0f" font-size="15" text-anchor="middle">🐬</text>`, x0+plotW*s.frac, y-4)
			}
		}
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" font-family="Georgia" fill="#8fa8c7">🌊 %d · 🏝 %d</text>`, x0+plotW+14, y+laneH/2+5, nW, nI)
	}

	// ---- THE SHAPE: coherence spectrum (log sigma vs band) ----
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="19" text-anchor="middle" font-family="Georgia" fill="#ffd166">LA FORMA DEL MAR — el espectro de coherencia: dos ramas y el desierto en calma</text>`, W/2, shapeY-6)
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="8" fill="#081020" stroke="#44608c"/>`, x0, shapeY, plotW, shapeH)
	// sigma log scale from 0.003 to 6
	yOf := func(sig float64) float64 {
		lo, hi := math.Log(0.003), math.Log(6.0)
		f := (math.Log(sig) - lo) / (hi - lo)
		return shapeY + shapeH - f*shapeH
	}
	// calm-sea band (0.35..1.8): the sonar's dismissal zone
	fmt.Fprintf(&b, `<rect x="%.0f" y="%.1f" width="%.0f" height="%.1f" fill="#8fa8c7" opacity="0.10"/>`, x0, yOf(1.8), plotW, yOf(0.35)-yOf(1.8))
	fmt.Fprintf(&b, `<text x="%.0f" y="%.1f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">el mar en calma (σ≈1): el sonar cuelga — aquí no vive nada que valga la pólvora</text>`, x0+plotW/2, yOf(0.8))
	for _, gl := range []float64{0.003, 0.01, 0.05, 0.35, 1, 1.8, 2.4, 4} {
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="#2c4a78" stroke-dasharray="3,5" opacity="0.7"/><text x="%.0f" y="%.1f" font-size="11" text-anchor="end" fill="#8fa8c7">%g</text>`,
			x0, yOf(gl), x0+plotW, yOf(gl), x0-6, yOf(gl)+4, gl)
	}
	for _, be := range beasts {
		x := x0 + plotW*be.frac
		fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.1" fill="%s" opacity="0.65"/>`, x, yOf(math.Max(be.coh, 0.0031)), wcol(be.t))
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.1f" font-size="14" font-family="Georgia" fill="#7fb2ff">↑ LA RAMA DE LAS OLAS (conspiración de fase)</text>`, x0+16, yOf(3.4))
	fmt.Fprintf(&b, `<text x="%.0f" y="%.1f" font-size="14" font-family="Georgia" fill="#7fd7a8">↓ LA RAMA DE LAS ISLAS (cancelación casi perfecta)</text>`, x0+16, yOf(0.0125))

	// ---- footer: records and the water colour key ----
	fy := shapeY + shapeH + 40
	var maxW, minI beast
	minI.coh = 99
	for _, be := range beasts {
		if !be.muda && be.coh > maxW.coh {
			maxW = be
		}
		if be.muda && be.coh < minI.coh {
			minI = be
		}
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" font-family="Georgia" fill="#dce8f7">degradé de aguas: `, W/2, fy)
	for _, w := range waters {
		fmt.Fprintf(&b, `<tspan fill="%s">%.0e</tspan>  `, wcol(w), w)
	}
	fmt.Fprintf(&b, `</text><text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd166">👑 ola récord: %.3fσ (|ola|=%.0f, t=%.0e) · 💰 silencio récord: %.4fσ (t=%.0e) · círculo azul/color = ola · triángulo verde = isla · 🐬 = cardumen</text>`,
		W/2, fy+28, maxW.coh, maxW.ola, maxW.t, minI.coh, minI.t)
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">la forma que dibuja la caza ES la anatomía del mar: dos ramas pobladas y un desierto en el medio — el tren solo pisa donde hay vida</text>`, W/2, fy+56)
	b.WriteString(`</svg>`)
	os.WriteFile("carta-abismo.svg", []byte(b.String()), 0644)
	fmt.Printf("escrita: carta-abismo.svg (%d bestias, %d cardúmenes, %d aguas)\n", len(beasts), len(schools), len(waters))
}
