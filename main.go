package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Maze Solver")
	maze := NewMaze(40, 30, 20)
	maze.DrawCells(&w)
	go maze.RemoveWalls()
	w.ShowAndRun()
}
