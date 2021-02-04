package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("List Test")

	data := []string{}
	for i := 0; i < 1000; i++ {
		line := fmt.Sprintf("%d", i)
		data = append(data, line)
	}
	list := widget.NewList(
		func() int { return len(data) },
		func() fyne.CanvasObject {
			return canvas.NewText("placeholder", theme.ForegroundColor())
		},
		func(i int, o fyne.CanvasObject) {
			textCanvas := o.(*canvas.Text)
			textCanvas.Text = data[i]
		},
	)

	myWindow.SetContent(list)
	myWindow.Resize(fyne.NewSize(400.0, 400.0))
	myWindow.Show()
	myApp.Run()
	tidyUp()
}

func tidyUp() {
	fmt.Println("Exited")
}
