package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func OpenFilePicker(myWindow fyne.Window, onSuccess func(string)) {
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

		// delegate to external file processing
		onSuccess(filePath)

		defer reader.Close()

	}, myWindow)

	// Filter only .pdf files
	pdfFilter := storage.NewExtensionFileFilter([]string{".pdf"})
	fileDialog.SetFilter(pdfFilter)

	fileDialog.Show()
}
