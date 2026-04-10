package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-resty/resty/v2"
)

/* Struct apropriate to Qwen AI model response */
type JsonAIResponse struct {
	Output []struct {
		Content string `json:"content"`
	} `json:"output"`
}

type CVData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone number"`
	City        string `json:"city"`
	BirthDate   string `json:"birth date"`
}

type CVMappedFields struct {
	FNameEntry        *widget.Entry
	FEmailEntry       *widget.Entry
	FPhoneNumberEntry *widget.Entry
	FCityBirthEntry   *widget.Entry
	FBirthDateEntry   *widget.Entry
}

type FileState struct {
	CurrentFilePath string
	StatusLabelPath *canvas.Text
	CVMappedFields  *CVMappedFields
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

	// Init http request to AI model
	json_result := &JsonAIResponse{}
	client := resty.New()

	go func() { // run in goroutine to avoid fyne being frozen
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"model":          "qwen2.5-coder-7b-instruct",
				"input":          "Parse this CV to JSON (name(Capitalize every word), email, phone number, city, birth date (dd/mm/year)), don't use _ instead of space in key attributes, and response only JSON!" + cv_content,
				"stream":         false,
				"context_length": 8000,
			}).
			SetResult(json_result).
			Post("http://192.168.100.10:3001/api/v1/chat")

		if err != nil {
			log.Fatalf("Connection failed: %v", err)
		}

		if resp.IsError() {
			fmt.Printf("Server AI returned error: %s\n", resp.Status())
			return
		}

		// Receive AI answer parsed to JSON struct
		if len(json_result.Output) > 0 {
			rawContent := json_result.Output[0].Content
			fmt.Println(rawContent)

			// Json cleaning from necessary markdowns
			cleanJSON := strings.TrimSpace(rawContent)
			cleanJSON = strings.TrimPrefix(cleanJSON, "```json")
			cleanJSON = strings.TrimSuffix(cleanJSON, "```")
			cleanJSON = strings.TrimSpace(cleanJSON)

			fmt.Printf("After Trim Operation: %s", cleanJSON)

			finalData := &CVData{}
			err := json.Unmarshal([]byte(cleanJSON), finalData)
			if err != nil {
				fmt.Println("Error unmarshaled JSON into struct: ", err)
				return
			}

			fmt.Printf("Success, received data: %s, %s\n", finalData.Name, finalData.Email)

			if finalData != nil {
				if len(fs.CVMappedFields.FNameEntry.Text) <= 0 {
					fyne.Do(func() {
						fs.CVMappedFields.FNameEntry.SetText(finalData.Name)
					})
				}
				if len(fs.CVMappedFields.FEmailEntry.Text) <= 0 {
					fyne.Do(func() {
						fs.CVMappedFields.FEmailEntry.SetText(finalData.Email)
					})
				}
				if len(fs.CVMappedFields.FPhoneNumberEntry.Text) <= 0 {
					fyne.Do(func() {
						fs.CVMappedFields.FPhoneNumberEntry.SetText(finalData.PhoneNumber)
					})
				}
				if len(fs.CVMappedFields.FCityBirthEntry.Text) <= 0 {
					fyne.Do(func() {
						fs.CVMappedFields.FCityBirthEntry.SetText(finalData.City)
					})
				}
				if len(fs.CVMappedFields.FBirthDateEntry.Text) <= 0 {
					fyne.Do(func() {
						fs.CVMappedFields.FBirthDateEntry.SetText(finalData.BirthDate)
					})
				}
			}
		}
	}()

}

type sWindow struct {
	myApp       *fyne.App
	id          string
	windowTitle string
	width       float32
	height      float32
}

func (w *sWindow) deriveWindow() fyne.Window {
	myApp := app.NewWithID(w.id)
	myWindow := myApp.NewWindow(w.windowTitle)
	myWindow.Resize(fyne.NewSize(w.width, w.height))
	return myWindow
}

func main() {
	sWindow := &sWindow{
		id:          "com.main.jobresearcher.tool",
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
		CVMappedFields:  &CVMappedFields{},
		Window:          currentWindow,
	}

	selectBtn := widget.NewButton("Browse", func() {
		OpenFilePicker(fileState.Window, fileState.HandleFileSelection)
	})

	content := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), fileState.StatusLabelPath, selectBtn, sMyLayout.rightMargin)

	/* basic cv input fields */
	fileState.CVMappedFields.FNameEntry = widget.NewEntry()
	fileState.CVMappedFields.FNameEntry.SetPlaceHolder("Type your name")
	fileState.CVMappedFields.FNameEntry.Resize(fyne.NewSize(200, fileState.CVMappedFields.FNameEntry.MinSize().Height))

	fileState.CVMappedFields.FEmailEntry = widget.NewEntry()
	fileState.CVMappedFields.FEmailEntry.SetPlaceHolder("Type your email")
	fileState.CVMappedFields.FEmailEntry.Resize(fyne.NewSize(200, fileState.CVMappedFields.FEmailEntry.MinSize().Height))

	fileState.CVMappedFields.FPhoneNumberEntry = widget.NewEntry()
	fileState.CVMappedFields.FPhoneNumberEntry.SetPlaceHolder("Type your phone number")
	fileState.CVMappedFields.FPhoneNumberEntry.Resize(fyne.NewSize(200, fileState.CVMappedFields.FPhoneNumberEntry.MinSize().Height))

	fileState.CVMappedFields.FCityBirthEntry = widget.NewEntry()
	fileState.CVMappedFields.FCityBirthEntry.SetPlaceHolder("Type your city")
	fileState.CVMappedFields.FCityBirthEntry.Resize(fyne.NewSize(200, fileState.CVMappedFields.FCityBirthEntry.MinSize().Height))

	fileState.CVMappedFields.FBirthDateEntry = widget.NewEntry()
	fileState.CVMappedFields.FBirthDateEntry.SetPlaceHolder("Type your birth date")
	fileState.CVMappedFields.FBirthDateEntry.Resize(fyne.NewSize(200, fileState.CVMappedFields.FBirthDateEntry.MinSize().Height))

	LT_Attach_Cv := container.New(layout.NewVBoxLayout(), sMyLayout.topMargin, content)

	sMyLayout.CreateNewRow(LT_Attach_Cv, container.NewWithoutLayout(fileState.CVMappedFields.FNameEntry))
	sMyLayout.CreateNewRow(LT_Attach_Cv, sMyLayout.rowGap)
	sMyLayout.CreateNewRow(LT_Attach_Cv, container.NewWithoutLayout(fileState.CVMappedFields.FEmailEntry))
	sMyLayout.CreateNewRow(LT_Attach_Cv, sMyLayout.rowGap)
	sMyLayout.CreateNewRow(LT_Attach_Cv, container.NewWithoutLayout(fileState.CVMappedFields.FPhoneNumberEntry))
	sMyLayout.CreateNewRow(LT_Attach_Cv, sMyLayout.rowGap)
	sMyLayout.CreateNewRow(LT_Attach_Cv, container.NewWithoutLayout(fileState.CVMappedFields.FCityBirthEntry))
	sMyLayout.CreateNewRow(LT_Attach_Cv, sMyLayout.rowGap)
	sMyLayout.CreateNewRow(LT_Attach_Cv, container.NewWithoutLayout(fileState.CVMappedFields.FBirthDateEntry))

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
