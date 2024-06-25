package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Maze Solver")
	canv := NewMaze()
	canv.DrawCells(&w)
	canv.DrawLine(Location{X: 3, Y: 3}, Location{X: 4, Y: 3}, false)
	canv.DrawLine(Location{X: 4, Y: 3}, Location{X: 5, Y: 3}, true)
	w.ShowAndRun()
}
