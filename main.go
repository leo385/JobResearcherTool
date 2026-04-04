package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

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
	fEmail := CreateInputWithPlaceholder("Type your email", 200)
	LT_Attach_Cv := container.New(layout.NewVBoxLayout(), topMargin, content)
	CreateNewRow(LT_Attach_Cv, fName)
	CreateNewRow(LT_Attach_Cv, rowGap)
	CreateNewRow(LT_Attach_Cv, fEmail)

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
