package main

/*
  https://docs.fyne.io/
  https://github.com/fyne-io/fyne
  go get fyne.io/fyne/v2@latest
  go install fyne.io/tools/cmd/fyne@latest
*/
import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")

	w.SetContent(widget.NewLabel("Hello World!"))
	w.ShowAndRun()
}
