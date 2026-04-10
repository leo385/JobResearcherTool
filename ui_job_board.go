package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func ui_board(window fyne.Window, customLayout *sLayout) *fyne.Container {
	BoardUrlEntry := widget.NewEntry()
	BoardUrlEntry.SetPlaceHolder("Provide the job board url...")
	BoardUrlEntry.Resize(fyne.NewSize(window.Canvas().Size().Width-200, 38))
	LT_Job_Board := container.New(layout.NewVBoxLayout(), customLayout.topMargin)
	customLayout.CreateNewRow(LT_Job_Board, container.NewWithoutLayout(BoardUrlEntry), customLayout.rightMargin)

	return LT_Job_Board
}
