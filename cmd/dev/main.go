package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Lentry")

	content := widget.NewLentry("This is a line\nor two of text\neven three\n[X|y[[\n{X|[y_\nThis is a really really long line that should really wrap but may not be long enough yet so i will type a bit more till it is wider than the screen")
	content.Wrapping = fyne.TextWrapWord

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(200.0, 200.0))
	myWindow.ShowAndRun()
}
