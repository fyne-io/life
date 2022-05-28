package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const cellSize = 8

type gameRenderer struct {
	render   *canvas.Raster
	objects  []fyne.CanvasObject
	imgCache *image.RGBA

	aliveColor color.Color
	deadColor  color.Color

	game *game
}

func (g *gameRenderer) MinSize() fyne.Size {
	pixDensity := g.game.pixelDensity()
	return fyne.NewSize(float32(minXCount*cellSize)/pixDensity, float32(minYCount*cellSize)/pixDensity)
}

func (g *gameRenderer) Layout(size fyne.Size) {
	g.render.Resize(size)
}

func (g *gameRenderer) ApplyTheme() {
	g.aliveColor = theme.ForegroundColor()
	g.deadColor = theme.BackgroundColor()
}

func (g *gameRenderer) Refresh() {
	canvas.Refresh(g.render)
}

func (g *gameRenderer) Objects() []fyne.CanvasObject {
	return g.objects
}

func (g *gameRenderer) Destroy() {
}

func (g *gameRenderer) draw(w, h int) image.Image {
	pixDensity := g.game.pixelDensity()
	pixW, pixH := g.game.cellForCoord(w, h, pixDensity)

	img := g.imgCache
	if img == nil || img.Bounds().Size().X != pixW || img.Bounds().Size().Y != pixH {
		img = image.NewRGBA(image.Rect(0, 0, pixW, pixH))
		g.imgCache = img
	}
	g.game.board.ensureGridSize(pixW, pixH)

	for y := 0; y < pixH; y++ {
		for x := 0; x < pixW; x++ {
			if x < g.game.board.width && y < g.game.board.height && g.game.board.currentGenCells[y][x] {
				img.Set(x, y, g.aliveColor)
			} else {
				img.Set(x, y, g.deadColor)
			}
		}
	}

	return img
}

type game struct {
	widget.BaseWidget

	genText *widget.Label
	board   *board
	paused  bool
}

func (g *game) CreateRenderer() fyne.WidgetRenderer {
	renderer := &gameRenderer{game: g}

	render := canvas.NewRaster(renderer.draw)
	render.ScaleMode = canvas.ImageScalePixels
	renderer.render = render
	renderer.objects = []fyne.CanvasObject{render}
	renderer.ApplyTheme()

	return renderer
}

func (g *game) cellForCoord(x, y int, density float32) (int, int) {
	xpos := int(float32(x) / float32(cellSize) / density)
	ypos := int(float32(y) / float32(cellSize) / density)

	return xpos, ypos
}

func (g *game) toggleRun() {
	g.paused = !g.paused
}

func (g *game) animate() {
	go func() {
		tick := time.NewTicker(time.Second / 6)

		for range tick.C {
			if g.paused {
				continue
			}

			g.board.nextGen()
			g.updateGeneration()
			g.Refresh()
		}
	}()
}

func (g *game) typedRune(r rune) {
	if r == ' ' {
		g.toggleRun()
	}
}

func (g *game) Tapped(ev *fyne.PointEvent) {
	pixDensity := g.pixelDensity()
	xpos, ypos := g.cellForCoord(int(ev.Position.X*pixDensity), int(ev.Position.Y*pixDensity), pixDensity)

	if ev.Position.X < 0 || ev.Position.Y < 0 || xpos >= g.board.width || ypos >= g.board.height {
		return
	}

	g.board.currentGenCells[ypos][xpos] = !g.board.currentGenCells[ypos][xpos]

	g.Refresh()
}

func (g *game) TappedSecondary(ev *fyne.PointEvent) {
}

func (g *game) buildUI() fyne.CanvasObject {
	var pause *widget.Button
	pause = widget.NewButton("Pause", func() {
		g.paused = !g.paused

		if g.paused {
			pause.SetText("Resume")
		} else {
			pause.SetText("Pause")
		}
	})

	title := container.NewGridWithColumns(2, g.genText, pause)
	return container.NewBorder(title, nil, nil, nil, g)
}

func (g *game) updateGeneration() {
	g.genText.SetText(fmt.Sprintf("Generation %d", g.board.generation))
}

func (g *game) pixelDensity() float32 {
	c := fyne.CurrentApp().Driver().CanvasForObject(g)
	if c == nil {
		return 1.0
	}

	pixW, _ := c.PixelCoordinateForPosition(fyne.NewPos(cellSize, cellSize))
	return float32(pixW) / float32(cellSize)
}

func newGame(b *board) *game {
	g := &game{board: b, genText: widget.NewLabel("Generation 0")}
	g.ExtendBaseWidget(g)

	return g
}
