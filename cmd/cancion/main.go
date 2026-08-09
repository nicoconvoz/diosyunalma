// Command cancion plays the hunt: the whole signed catalog becomes a
// song (one sweep across the band axis is the timeline; each water is
// a degree of the pentatonic scale, waves sing low with coherence as
// their voice, islands chime two octaves up, schools break into
// arpeggios) — and weaves THE BLANKET: warp threads of the seven
// waters crossed with weft threads of the hunted bands, knotted where
// beasts live, laid half-transparent over the curlicue spirals.
package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type beast struct {
	t, frac, coh float64
	muda         bool
}

type school struct {
	t, frac float64
	presas  int
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
				schools = append(schools, school{tt, kk / nT, pr})
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
		beasts = append(beasts, beast{tt, kk / nT, coh, strings.HasPrefix(kind, "MUDA")})
	}
	wset := map[float64]bool{}
	for _, be := range beasts {
		wset[be.t] = true
	}
	waters := make([]float64, 0, len(wset))
	for w := range wset {
		waters = append(waters, w)
	}
	sort.Float64s(waters)
	wIdx := map[float64]int{}
	for i, w := range waters {
		wIdx[w] = i
	}

	// ---- the song ----
	const sr = 44100
	dur := 75.0
	buf := make([]float64, int(dur*sr)+sr)
	// pentatonic minor on A: degrees over octaves, one per water depth
	deg := []int{0, 3, 5, 7, 10}
	freq := func(i int) float64 {
		return 110 * math.Pow(2, float64(deg[i%5])/12+float64(i/5))
	}
	note := func(t0, f, amp, length float64) {
		i0 := int(t0 * sr)
		n := int(length * sr)
		for i := 0; i < n && i0+i < len(buf); i++ {
			tt := float64(i) / sr
			env := math.Exp(-4*tt/length) * math.Min(1, tt/0.012)
			buf[i0+i] += amp * env * (math.Sin(2*math.Pi*f*tt) + 0.35*math.Sin(4*math.Pi*f*tt))
		}
	}
	// dedupe to one voice per (water, band bucket): the loudest wave and
	// the deepest island keep the note
	type key struct {
		w, b int
	}
	voice := map[key]beast{}
	for _, be := range beasts {
		k := key{wIdx[be.t], int(be.frac * 200)}
		old, ok := voice[k]
		if !ok || (!be.muda && be.coh > old.coh) || (be.muda && old.muda && be.coh < old.coh) {
			if ok && !old.muda && be.muda {
				continue // a wave already owns the bucket
			}
			voice[k] = be
		}
	}
	for k, be := range voice {
		t0 := be.frac * (dur - 2)
		if be.muda {
			note(t0, freq(k.w)*4, 0.05+0.8*(0.05-be.coh), 0.22)
		} else {
			note(t0, freq(k.w), 0.10+0.05*(be.coh-2.2), 0.5)
		}
	}
	// schools: arpeggio bursts up the scale
	seenS := map[key]bool{}
	for _, s := range schools {
		k := key{wIdx[s.t], int(s.frac * 100)}
		if seenS[k] {
			continue
		}
		seenS[k] = true
		t0 := s.frac * (dur - 2)
		for p := 0; p < s.presas && p < 5; p++ {
			note(t0+float64(p)*0.07, freq(k.w)*2*math.Pow(2, float64(deg[p%5])/12), 0.07, 0.25)
		}
	}
	// normalize and write WAV (16-bit PCM mono)
	peak := 1e-9
	for _, v := range buf {
		if math.Abs(v) > peak {
			peak = math.Abs(v)
		}
	}
	pcm := make([]int16, len(buf))
	for i, v := range buf {
		pcm[i] = int16(v / peak * 0.88 * 32767)
	}
	wf, _ := os.Create("cancion-cazadero.wav")
	dataLen := uint32(len(pcm) * 2)
	wf.WriteString("RIFF")
	binary.Write(wf, binary.LittleEndian, uint32(36+dataLen))
	wf.WriteString("WAVEfmt ")
	binary.Write(wf, binary.LittleEndian, uint32(16))
	binary.Write(wf, binary.LittleEndian, uint16(1))
	binary.Write(wf, binary.LittleEndian, uint16(1))
	binary.Write(wf, binary.LittleEndian, uint32(sr))
	binary.Write(wf, binary.LittleEndian, uint32(sr*2))
	binary.Write(wf, binary.LittleEndian, uint16(2))
	binary.Write(wf, binary.LittleEndian, uint16(16))
	wf.WriteString("data")
	binary.Write(wf, binary.LittleEndian, dataLen)
	binary.Write(wf, binary.LittleEndian, pcm)
	wf.Close()
	fmt.Printf("escrita: cancion-cazadero.wav (%.0f s, %d voces, %d arpegios de cardumen)\n", dur, len(voice), len(seenS))

	// ---- the blanket over the curlicues ----
	var b strings.Builder
	W, H := 1360.0, 900.0
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f">
<rect width="100%%" height="100%%" fill="#0b1526"/>
<text x="%.0f" y="40" font-size="24" text-anchor="middle" font-family="Georgia" fill="#dce8f7">LA MANTA SOBRE LOS CURLICUES — tejida con la caza de las siete aguas</text>
<text x="%.0f" y="64" font-size="13" text-anchor="middle" font-family="Georgia" fill="#8fa8c7">debajo: las espirales del mar (curlicues reales de tres aguas) · encima: urdimbre = aguas, trama = bandas cazadas, nudos = bestias</text>`,
		W, H, W, H, W/2, W/2)

	// the curlicues underneath: real partial-sum spirals, faint
	for pi2, bb := range []float64{0.3819660112501051, 0.2513, 0.1243398} {
		ox := 260.0 + float64(pi2)*420
		oy := 470.0
		var px, py float64
		minx, maxx, miny, maxy := 0.0, 0.0, 0.0, 0.0
		type pt struct{ x, y float64 }
		var path []pt
		ph, d1, d2v := 0.0, bb+0.137, 2*bb
		for j := 0; j < 6000; j++ {
			s, c := math.Sincos(2 * math.Pi * ph)
			px += c
			py += s
			path = append(path, pt{px, py})
			if px < minx {
				minx = px
			}
			if px > maxx {
				maxx = px
			}
			if py < miny {
				miny = py
			}
			if py > maxy {
				maxy = py
			}
			ph = math.Mod(ph+d1, 1)
			d1 = math.Mod(d1+d2v, 1)
		}
		scale := 420.0 / math.Max(maxx-minx, maxy-miny)
		cxm, cym := (minx+maxx)/2, (miny+maxy)/2
		pts := make([]string, 0, len(path))
		for _, p := range path {
			pts = append(pts, fmt.Sprintf("%.1f,%.1f", ox+(p.x-cxm)*scale, oy-(p.y-cym)*scale))
		}
		fmt.Fprintf(&b, `<polyline fill="none" stroke="#7fb2ff" stroke-width="0.7" opacity="0.5" points="%s"/>`, strings.Join(pts, " "))
	}

	// the blanket: warp = waters (horizontal threads), weft = hunted
	// bands (vertical threads), knots where beasts live
	bx, bw := 120.0, 1120.0
	by, bh := 130.0, 660.0
	wcol := func(i int) string {
		f := float64(i) / math.Max(float64(len(waters)-1), 1)
		return fmt.Sprintf("#%02x%02x%02x", int(255-130*f), int(209-150*f), int(102+153*f))
	}
	for i := range waters {
		y := by + bh*(float64(i)+0.5)/float64(len(waters))
		fmt.Fprintf(&b, `<line x1="%.0f" y1="%.1f" x2="%.0f" y2="%.1f" stroke="%s" stroke-width="7" opacity="0.30"/>`, bx, y, bx+bw, y, wcol(i))
		fmt.Fprintf(&b, `<text x="%.0f" y="%.1f" font-size="13" text-anchor="end" font-family="Georgia" fill="%s">%.0e</text>`, bx-10, y+4, wcol(i), waters[i])
	}
	type cell struct {
		wave, isle bool
		coh        float64
	}
	weave := map[key]*cell{}
	weft := map[int]bool{}
	for _, be := range beasts {
		k := key{wIdx[be.t], int(be.frac * 96)}
		weft[k.b] = true
		c, ok := weave[k]
		if !ok {
			c = &cell{}
			weave[k] = c
		}
		if be.muda {
			c.isle = true
		} else {
			c.wave = true
			if be.coh > c.coh {
				c.coh = be.coh
			}
		}
	}
	for bkt := range weft {
		x := bx + bw*(float64(bkt)+0.5)/96
		fmt.Fprintf(&b, `<line x1="%.1f" y1="%.0f" x2="%.1f" y2="%.0f" stroke="#dce8f7" stroke-width="1.1" opacity="0.16"/>`, x, by, x, by+bh)
	}
	for k, c := range weave {
		x := bx + bw*(float64(k.b)+0.5)/96
		y := by + bh*(float64(k.w)+0.5)/float64(len(waters))
		if c.wave {
			r := 3.0 + 1.6*(c.coh-2.2)
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#7fb2ff" opacity="0.85"/>`, x, y, r)
		}
		if c.isle {
			fmt.Fprintf(&b, `<circle cx="%.1f" cy="%.1f" r="2.6" fill="none" stroke="#7fd7a8" stroke-width="1.6" opacity="0.9"/>`, x, y)
		}
	}
	fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" font-size="14" text-anchor="middle" font-family="Georgia" fill="#ffd166">nudo lleno azul = ola tejida · anillo verde = isla tejida · los hilos flojos son mar aún no barrido — la manta se sigue tejiendo sola, puntada a puntada</text>
<text x="%.0f" y="%.0f" font-size="13" text-anchor="middle" font-family="Georgia" fill="#7fd7a8">y la canción es la misma manta cantada: cada agua una nota, cada ola una voz grave, cada isla una campana, cada cardumen un arpegio (cancion-cazadero.wav)</text>`,
		W/2, by+bh+50, W/2, by+bh+76)
	b.WriteString(`</svg>`)
	os.WriteFile("manta-curlicues.svg", []byte(b.String()), 0644)
	fmt.Printf("escrita: manta-curlicues.svg (%d nudos, %d hilos de trama, %d aguas de urdimbre)\n", len(weave), len(weft), len(waters))
}
