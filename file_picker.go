package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

var filePathText = canvas.NewText("Choose CV file", color.White)

func OpenFilePicker(myWindow fyne.Window) {
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, myWindow)
			return
		}
		if reader == nil {
			fmt.Println("Canceled choosed file")
			return
		}

		filePath := reader.URI().Path()
		fmt.Println("Attached file:", filePath)
		filePathText.Text = filePath
		filePathText.Refresh()

		defer reader.Close()

	}, myWindow)

	// Filter only .pdf files
	pdfFilter := storage.NewExtensionFileFilter([]string{".pdf"})
	fileDialog.SetFilter(pdfFilter)

	fileDialog.Show()
}
