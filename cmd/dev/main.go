package main

import (
	"io/ioutil"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	b, err := ioutil.ReadFile("save_test.json")
	if err != nil {
		log.Fatal(err)
	}
	myApp := app.New()
	myWindow := myApp.NewWindow("Lentry")

	//lentry := widget.NewLentry("This is a line\nor two of text\neven three\n[X|y[[\n{X|[y_\nThis is a really really long line that should really wrap but may not be long enough yet so i will type a bit more till it is wider than the screen")
	lentry := widget.NewLentry(string(b))
	lentry.Wrapping = fyne.TextWrapWord

	myWindow.SetContent(lentry)
	myWindow.Resize(fyne.NewSize(200.0, 200.0))
	myWindow.ShowAndRun()
}
