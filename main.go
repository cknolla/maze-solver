package main

import (
	"fyne.io/fyne/v2/app"
	"math/rand/v2"
)

func main() {
	colCount := 80
	rowCount := 50
	myApp := app.New()
	w := myApp.NewWindow("Maze Solver")
	maze := NewMaze(colCount, rowCount, 20, rand.NewPCG(rand.Uint64(), rand.Uint64()))
	maze.animationDelay = 0.002
	maze.DrawCells(&w)
	go maze.Animate()
	w.ShowAndRun()
}
