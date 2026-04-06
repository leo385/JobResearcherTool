package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FileState struct {
	CurrentFilePath string
	StatusLabelPath *canvas.Text
	Window          fyne.Window
}

func (fs *FileState) HandleFileSelection(path string) {
	fs.CurrentFilePath = path
	fs.StatusLabelPath.Text = "Selected: " + path
	fs.StatusLabelPath.Refresh()

	// pdf file processing
	cv_content, err := readPdf(path)
	if err != nil {
		panic(err)
	}
	fmt.Println(cv_content)
}

type sWindow struct {
	windowTitle string
	width       float32
	height      float32
}

func (w *sWindow) deriveWindow() fyne.Window {
	myApp := app.New()
	myWindow := myApp.NewWindow(w.windowTitle)
	myWindow.Resize(fyne.NewSize(w.width, w.height))
	return myWindow
}

func main() {
	sWindow := &sWindow{
		windowTitle: "JobResearcherTool v0.0.1",
		width:       800,
		height:      600,
	}
	currentWindow := sWindow.deriveWindow()

	sMyLayout := &sLayout{
		leftMargin:  SetMargin(20, 0),
		rightMargin: SetMargin(20, 0),
		topMargin:   SetMargin(0, 20),
		rowGap:      SetMargin(0, 5),
	}

	fileState := &FileState{
		StatusLabelPath: canvas.NewText("Choose your CV file", color.White),
		Window:          currentWindow,
	}

	selectBtn := widget.NewButton("Browse", func() {
		OpenFilePicker(fileState.Window, fileState.HandleFileSelection)
	})

	content := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), fileState.StatusLabelPath, selectBtn, sMyLayout.rightMargin)

	/* basic cv input fields */
	fName := CreateInputWithPlaceholder("Type your name", 200)
	fEmail := CreateInputWithPlaceholder("Type your email", 200)
	LT_Attach_Cv := container.New(layout.NewVBoxLayout(), sMyLayout.topMargin, content)

	sMyLayout.CreateNewRow(LT_Attach_Cv, fName)
	sMyLayout.CreateNewRow(LT_Attach_Cv, sMyLayout.rowGap)
	sMyLayout.CreateNewRow(LT_Attach_Cv, fEmail)

	tabs := container.NewAppTabs(
		container.NewTabItem("Attach your CV", LT_Attach_Cv),
		container.NewTabItem("Set board URL", widget.NewLabel("World!")),
		container.NewTabItem("Job specification", widget.NewLabel("World!")),
	)

	tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationTop)

	currentWindow.SetContent(tabs)
	currentWindow.ShowAndRun()
}
