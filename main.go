// Package main launches the life app
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

const (
	minXCount = 50
	minYCount = 40
)

func show(app fyne.App) {
	board := newBoard(minXCount, minYCount)
	board.load()
	game := newGame(board)

	window := app.NewWindow("Life")
	window.SetIcon(resourceIconPng)

	window.SetContent(game.buildUI())
	window.Canvas().SetOnTypedRune(game.typedRune)

	// start the board animation before we show the window - it will block
	game.animate()

	window.Show()
}

func main() {
	app := app.New()
	app.SetIcon(resourceIconPng)

	show(app)
	app.Run()
}
