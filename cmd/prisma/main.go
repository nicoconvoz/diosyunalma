// Command prisma paints the sea of numbers as a sea of color.
//
// Domain coloring of the Riemann zeta function on the critical strip:
// hue is the phase of zeta(s), brightness its magnitude. A zero is not
// one color — it is the vortex where EVERY color meets, a rainbow
// pinwheel collapsing to black. Recognizing a zero becomes a visual,
// topological act: count how many rainbows wrap around a point.
//
// That count is exactly what the flagship's sphere computes on the
// boundary of a window (the argument principle): the sphere has been
// navigating by color all along — this command makes it visible.
//
// zeta(s) off the critical line is computed by Euler–Maclaurin, exact
// far beyond pixel precision at these heights. The full 2D painting is
// affordable in shallow water; in deep water the light bucket already
// paints the 1D ribbon along the line.
//
// Usage:
//
//	go run ./cmd/prisma        # writes prisma.png
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
	"os"
	"time"
)

const (
	width  = 1800
	height = 600
	tMin   = 5.0
	tMax   = 65.0
	sigMin = -0.5
	sigMax = 1.5
	nTerms = 120
)

// zeta computes zeta(s) via Euler–Maclaurin with three correction terms.
func zeta(s complex128, lnn []float64) complex128 {
	var sum complex128
	sig, t := real(s), imag(s)
	for n := 1; n < nTerms; n++ {
		amp := math.Exp(-sig * lnn[n])
		sn, cs := math.Sincos(t * lnn[n])
		sum += complex(amp*cs, -amp*sn)
	}
	nf := complex(float64(nTerms), 0)
	ns := cmplx.Exp(-s * complex(lnn[nTerms], 0)) // N^{-s}
	sum += ns * nf / (s - 1)
	sum += ns / 2
	sum += ns * s / nf / 12
	sum -= ns * s * (s + 1) * (s + 2) / (nf * nf * nf) / 720
	sum += ns * s * (s + 1) * (s + 2) * (s + 3) * (s + 4) /
		(nf * nf * nf * nf * nf) / 30240
	return sum
}

// paint maps a zeta value to a color: hue from the phase, darkness at
// zeros, washing toward white near large magnitude, with faint
// magnitude rings for texture.
func paint(z complex128) color.RGBA {
	m := cmplx.Abs(z)
	h := (cmplx.Phase(z) + math.Pi) / (2 * math.Pi)
	v := m / (m + 0.35)
	ring := math.Log(m+1e-30) / math.Log(1.6)
	v *= 0.88 + 0.12*(ring-math.Floor(ring))
	sat := 1 / (1 + 0.02*m*m)
	return hsv(h, sat, v)
}

func hsv(h, s, v float64) color.RGBA {
	i := int(h*6) % 6
	f := h*6 - math.Floor(h*6)
	p, q, t := v*(1-s), v*(1-s*f), v*(1-s*(1-f))
	var r, g, b float64
	switch i {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	default:
		r, g, b = v, p, q
	}
	return color.RGBA{uint8(r * 255), uint8(g * 255), uint8(b * 255), 255}
}

func main() {
	fmt.Println("THE PRISM — painting the sea of numbers as a sea of color")

	lnn := make([]float64, nTerms+1)
	for n := 1; n <= nTerms; n++ {
		lnn[n] = math.Log(float64(n))
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	start := time.Now()
	for py := 0; py < height; py++ {
		sig := sigMax - (sigMax-sigMin)*float64(py)/float64(height-1)
		for px := 0; px < width; px++ {
			t := tMin + (tMax-tMin)*float64(px)/float64(width-1)
			img.SetRGBA(px, py, paint(zeta(complex(sig, t), lnn)))
		}
	}

	// a faint guide along the critical line sigma = 1/2.
	yLine := int(math.Round(float64(height-1) * (sigMax - 0.5) / (sigMax - sigMin)))
	for px := 0; px < width; px++ {
		c := img.RGBAAt(px, yLine)
		c.R = uint8(float64(c.R)*0.88 + 255*0.12)
		c.G = uint8(float64(c.G)*0.88 + 255*0.12)
		c.B = uint8(float64(c.B)*0.88 + 255*0.12)
		img.SetRGBA(px, yLine, c)
	}

	out, err := os.Create("prisma.png")
	if err != nil {
		panic(err)
	}
	if err := png.Encode(out, img); err != nil {
		panic(err)
	}
	out.Close()

	fmt.Printf("\n  strip painted: sigma in [%.1f, %.1f], t in [%.0f, %.0f], %dx%d px, %.1fs\n",
		sigMin, sigMax, tMin, tMax, width, height, time.Since(start).Seconds())
	fmt.Println("  every zero is a vortex where all colors meet - and every one")
	fmt.Println("  sits on the faint line at sigma = 1/2: the Riemann Hypothesis,")
	fmt.Println("  painted. The sphere counts these rainbow turns on any window")
	fmt.Println("  boundary - the ship has been navigating by color all along.")
	fmt.Println("\n  written: prisma.png")
}
