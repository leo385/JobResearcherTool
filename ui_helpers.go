package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func SetMargin(x float32, y float32) *canvas.Rectangle {
	margin := canvas.NewRectangle(color.Transparent)
	margin.SetMinSize(fyne.NewSize(x, y))
	return margin
}

var leftMargin = SetMargin(20, 0)
var rightMargin = SetMargin(20, 0)
var topMargin = SetMargin(0, 20)
var rowGap = SetMargin(0, 5)

func CreateInputWithPlaceholder(placeholder string, width float32) *fyne.Container {
	input := widget.NewEntry()
	input.SetPlaceHolder(placeholder)
	input.Resize(fyne.NewSize(width, input.MinSize().Height))
	return container.NewWithoutLayout(input)
}

func CreateNewRow(parentContainer *fyne.Container, widgets ...fyne.CanvasObject) {
	newRow := container.New(layout.NewHBoxLayout(), leftMargin)
	for _, w := range widgets {
		newRow.Add(w)
	}

	parentContainer.Add(newRow)
	parentContainer.Refresh()
}
