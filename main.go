package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var filePathText = canvas.NewText("Choose CV file", color.White)
var newVerticalObjects []fyne.CanvasObject

func SetMargin(x float32, y float32) *canvas.Rectangle {
	margin := canvas.NewRectangle(color.Transparent)
	margin.SetMinSize(fyne.NewSize(x, y))
	return margin
}

var leftMargin = SetMargin(20, 0)
var rightMargin = SetMargin(20, 0)
var topMargin = SetMargin(0, 20)

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

	// Filtr only .pdf files
	pdfFilter := storage.NewExtensionFileFilter([]string{".pdf"})
	fileDialog.SetFilter(pdfFilter)

	fileDialog.Show()
}

func CreateInputWithPlaceholder(placeholder string, width float32) *fyne.Container {
	input := widget.NewEntry()
	input.SetPlaceHolder(placeholder)
	input.Resize(fyne.NewSize(width, input.MinSize().Height))
	return container.NewWithoutLayout(input)
}

func CreateNewRow(widgetObj fyne.CanvasObject, parentContainer *fyne.Container) {
	newRow := container.New(layout.NewHBoxLayout(), leftMargin, widgetObj)
	parentContainer.Add(newRow)
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("JobResearcherTool v0.0.1")
	myWindow.Resize(fyne.NewSize(800, 600))

	selectBtn := widget.NewButton("Browse", func() {
		OpenFilePicker(myWindow)
	})

	content := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), filePathText, selectBtn, rightMargin)
	/* basic cv input fields */
	fName := CreateInputWithPlaceholder("Type your name", 200)
	LT_Attach_Cv := container.New(layout.NewVBoxLayout(), topMargin, content)
	CreateNewRow(fName, LT_Attach_Cv)

	tabs := container.NewAppTabs(
		container.NewTabItem("Attach your CV", LT_Attach_Cv),
		container.NewTabItem("Set board URL", widget.NewLabel("World!")),
		container.NewTabItem("Job specification", widget.NewLabel("World!")),
	)

	tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationTop)

	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
