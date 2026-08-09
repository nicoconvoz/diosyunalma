// Command icono draws the laboratory's icon and writes it as a Windows .ico.
//
// The mark is the shop's own signature: the unit circle (the skin), its
// horizontal diameter, the two golden points +1 and -1 that carry no pearl
// (F254, F258, F260), and one green pearl sitting on the skin. It stays
// readable down to 16 pixels, which is the only test an icon has to pass.
//
// The file is assembled by hand: an ICONDIR header, one ICONDIRENTRY per size,
// and PNG payloads (Vista and later accept PNG inside .ico). Everything is
// drawn at 4x and averaged down, which is all the antialiasing this needs.
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

var (
	fondo  = color.NRGBA{0x0b, 0x15, 0x26, 0xff}
	oro    = color.NRGBA{0xff, 0xd1, 0x66, 0xff}
	verde  = color.NRGBA{0x7f, 0xd7, 0xa8, 0xff}
	borde  = color.NRGBA{0x3d, 0x6f, 0xa8, 0xff}
	fuerte = 4
)

// dibujar renders the mark at n x n pixels, supersampled by fuerte.
func dibujar(n int) *image.NRGBA {
	S := n * fuerte
	gr := image.NewNRGBA(image.Rect(0, 0, S, S))
	c := float64(S) / 2
	radio := float64(S) * 0.33
	grosor := math.Max(float64(S)*0.045, 1.2)
	esquina := float64(S) * 0.18
	puntoR := math.Max(float64(S)*0.062, 1.5)
	perlaR := math.Max(float64(S)*0.05, 1.2)

	// the pearl sits at a pleasant angle on the skin
	ang := 2.0
	px, py := c+radio*math.Cos(ang), c-radio*math.Sin(ang)

	for y := 0; y < S; y++ {
		for x := 0; x < S; x++ {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			// rounded-square background
			dx := math.Max(math.Abs(fx-c)-(c-esquina), 0)
			dy := math.Max(math.Abs(fy-c)-(c-esquina), 0)
			if math.Hypot(dx, dy) > esquina {
				continue
			}
			col := fondo
			// a faint rim so the tile reads on light desktops
			if math.Hypot(dx, dy) > esquina-float64(S)*0.02 {
				col = borde
			}
			d := math.Hypot(fx-c, fy-c)
			// the skin
			if math.Abs(d-radio) < grosor/2 {
				col = oro
			}
			// the horizontal diameter, from -1 to +1
			if math.Abs(fy-c) < grosor/2 && math.Abs(fx-c) < radio {
				col = oro
			}
			// the two points with no pearl
			if math.Hypot(fx-(c-radio), fy-c) < puntoR || math.Hypot(fx-(c+radio), fy-c) < puntoR {
				col = oro
			}
			// one pearl on the skin
			if math.Hypot(fx-px, fy-py) < perlaR {
				col = verde
			}
			gr.SetNRGBA(x, y, col)
		}
	}

	// average the supersample down
	out := image.NewNRGBA(image.Rect(0, 0, n, n))
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			var r, g, b, a int
			for j := 0; j < fuerte; j++ {
				for i := 0; i < fuerte; i++ {
					p := gr.NRGBAAt(x*fuerte+i, y*fuerte+j)
					r += int(p.R) * int(p.A) / 255
					g += int(p.G) * int(p.A) / 255
					b += int(p.B) * int(p.A) / 255
					a += int(p.A)
				}
			}
			k := fuerte * fuerte
			out.SetNRGBA(x, y, color.NRGBA{uint8(r / k), uint8(g / k), uint8(b / k), uint8(a / k)})
		}
	}
	return out
}

func main() {
	destino := "instalador/diosyunalma.ico"
	if len(os.Args) > 1 {
		destino = os.Args[1]
	}
	// con destino .png escribe una sola lámina, para poder MIRAR el ícono antes
	// de meterlo en un acceso directo; con .ico arma el archivo de Windows
	if len(destino) > 4 && destino[len(destino)-4:] == ".png" {
		f, err := os.Create(destino)
		if err != nil {
			fmt.Println("no pude crear el png:", err)
			os.Exit(1)
		}
		defer f.Close()
		if err := png.Encode(f, dibujar(256)); err != nil {
			fmt.Println("no pude codificar el png:", err)
			os.Exit(1)
		}
		fmt.Printf("vista previa escrita: %s (256x256)\n", destino)
		return
	}

	tamanos := []int{256, 128, 64, 48, 32, 16}

	var payloads [][]byte
	for _, n := range tamanos {
		var buf bytes.Buffer
		if err := png.Encode(&buf, dibujar(n)); err != nil {
			fmt.Println("no pude codificar el png:", err)
			os.Exit(1)
		}
		payloads = append(payloads, buf.Bytes())
	}

	var ico bytes.Buffer
	binary.Write(&ico, binary.LittleEndian, uint16(0)) // reservado
	binary.Write(&ico, binary.LittleEndian, uint16(1)) // tipo: icono
	binary.Write(&ico, binary.LittleEndian, uint16(len(tamanos)))
	off := 6 + 16*len(tamanos)
	for i, n := range tamanos {
		lado := byte(n)
		if n >= 256 {
			lado = 0 // 256 se escribe como 0, así manda el formato
		}
		ico.WriteByte(lado)                                 // ancho
		ico.WriteByte(lado)                                 // alto
		ico.WriteByte(0)                                    // colores de paleta
		ico.WriteByte(0)                                    // reservado
		binary.Write(&ico, binary.LittleEndian, uint16(1))  // planos
		binary.Write(&ico, binary.LittleEndian, uint16(32)) // bits por pixel
		binary.Write(&ico, binary.LittleEndian, uint32(len(payloads[i])))
		binary.Write(&ico, binary.LittleEndian, uint32(off))
		off += len(payloads[i])
	}
	for _, p := range payloads {
		ico.Write(p)
	}

	if err := os.WriteFile(destino, ico.Bytes(), 0o644); err != nil {
		fmt.Println("no pude escribir el icono:", err)
		os.Exit(1)
	}
	fmt.Printf("icono escrito: %s (%d tamaños, %d bytes)\n", destino, len(tamanos), ico.Len())
}
