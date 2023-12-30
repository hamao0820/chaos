package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 640
	plotCount    = int(1e6)
)

var (
	a, b, c, d float32
	countMap   [screenHeight + 1][screenWidth + 1]int
)

func init() {
	a = rand.Float32()*6 - 3
	b = rand.Float32()*6 - 3
	c = rand.Float32()*6 - 3
	d = rand.Float32()*6 - 3
}

type Game struct {
	offscreen    *ebiten.Image
	offscreenPix []byte
}

func NewGame() *Game {
	g := &Game{
		offscreen:    ebiten.NewImage(screenWidth, screenHeight),
		offscreenPix: make([]byte, screenWidth*screenHeight*4),
	}

	g.updateOffscreen()
	return g
}

func (g *Game) updateOffscreen() {
	var x, y float32
	for i := 0; i < plotCount; i++ {
		nx := math.Sin(float64(a)*float64(y)) - math.Cos(float64(b)*float64(x))
		ny := math.Sin(float64(c)*float64(x)) - math.Cos(float64(d)*float64(y))

		px := math.Floor(float64((nx + 2) / 4 * screenWidth))
		py := math.Floor(float64((ny + 2) / 4 * screenHeight))

		if px < 0 || px >= screenWidth || py < 0 || py >= screenHeight {
			continue
		}

		countMap[int(py)][int(px)]++

		x = float32(nx)
		y = float32(ny)
	}

	for y := 0; y <= screenHeight; y++ {
		for x := 0; x <= screenWidth; x++ {
			if countMap[y][x] == 0 {
				continue
			}

			ratio := float64(countMap[y][x]) / float64(plotCount) * screenWidth
			if ratio > 1 {
				ratio = 0.99
			}

			k := 1 - math.Pow(1-ratio, 50)

			off := y*screenWidth*4 + x*4
			g.offscreenPix[off] = 0xff
			g.offscreenPix[off+1] = byte(0xff * k)
			g.offscreenPix[off+2] = 0xff
			g.offscreenPix[off+3] = byte(0xff * k)

		}
	}

	g.offscreen.WritePixels(g.offscreenPix)
}

func (g *Game) Update() error {
	switch {
	case inpututil.IsKeyJustPressed(ebiten.KeyEnter):
		a = rand.Float32()*6 - 3
		b = rand.Float32()*6 - 3
		c = rand.Float32()*6 - 3
		d = rand.Float32()*6 - 3

		g.offscreen.Dispose()
		g.offscreen.Fill(color.Black)
		g.offscreen = ebiten.NewImage(screenWidth, screenHeight)
		g.offscreenPix = make([]byte, screenWidth*screenHeight*4)
		countMap = [screenHeight + 1][screenWidth + 1]int{}
		g.updateOffscreen()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.offscreen, nil)
	ebitenutil.DebugPrint(screen, "Press Enter to generate a new image")
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("a: %f, b: %f, c: %f, d: %f", a, b, c, d), 0, 20)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Mandelbrot (Ebitengine Demo)")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
