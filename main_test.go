package main

import (
	"fyne.io/fyne/v2/app"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// prevent "app not running" errors
	app.New()
	code := m.Run()
	os.Exit(code)
}
