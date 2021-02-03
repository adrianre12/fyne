package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Lentry")
	data := []string{"a[_", "b[_", "c[_"}

	content := widget.NewLentry(
		func() int { return len(data) },
		func() fyne.CanvasObject { return widget.NewLabel("Template") },
		func(i int, o fyne.CanvasObject) { o.(*widget.Label).SetText(data[i]) },
	)
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(200.0, 200.0))
	myWindow.ShowAndRun()
}
