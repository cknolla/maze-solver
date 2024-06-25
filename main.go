package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Maze Solver")
	maze := NewMaze()
	maze.DrawCells(&w)
	maze.HideWall(Location{X: 5, Y: 5}, bottom)
	maze.HideWall(Location{X: 9, Y: 5}, right)
	maze.DrawLine(Location{X: 3, Y: 3}, Location{X: 4, Y: 3}, false)
	maze.DrawLine(Location{X: 4, Y: 3}, Location{X: 5, Y: 3}, true)
	w.ShowAndRun()
}
