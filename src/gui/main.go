package main

import (
	"errors"
	// "log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	// "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	assets "github.com/luciferchase/CBSE-Result-Converter/src/gui/assets"
	backend "github.com/luciferchase/CBSE-Result-Converter/src/backend"
	layout "github.com/luciferchase/CBSE-Result-Converter/src/gui/custom/layout"
)

func main() {
	app := app.New()

	window := app.NewWindow("Result Converter")
	window.Resize(fyne.NewSize(700, 500))
	window.SetFixedSize(true)

	window.SetIcon(assets.Icon)

	logo := canvas.NewImageFromResource(assets.Logo)
	logo.SetMinSize(fyne.NewSize(700, 90))


	var class string
	classRadio := widget.NewRadioGroup([]string{"Class 10", "Class 12"},
		func(value string) {
			if value == "Class 10" {
				class = "X"
			} else if value == "Class 12" {
				class = "XII"
			}
		},
	)

	fileNameLabel := widget.NewLabel("No file selected!")

	var file fyne.URI
	openFileButton := widget.NewButtonWithIcon("Browse", theme.FileIcon(), func() {
		openFileDialog := dialog.NewFileOpen(func(callback fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if callback == nil {
				dialog.ShowError(errors.New("Select a file!"), window)
			} else {
				file = callback.URI()
				fileNameLabel.SetText(file.Name())
			}
		}, window)

		openFileDialog.Resize(fyne.NewSize(600, 400))
		openFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt"}))

		cwd, _ := os.Getwd()
		dir, _ := storage.ListerForURI(storage.NewFileURI(cwd))
		openFileDialog.SetLocation(dir)

		openFileDialog.Show()
	})
	openFileButton.Resize(fileNameLabel.MinSize())

	convertButton := widget.NewButtonWithIcon("Convert to CSV", theme.StorageIcon(), func() {
		// First check if class is selected or not
		if len(class) == 0 {
			dialog.ShowError(errors.New("Select a class first!"), window)
			return
		}

		parsedData, _, missingSubjectCodes, err := backend.Parse(file.Path(), class)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		if len(missingSubjectCodes) != 0 {
			dialogBoxContent := widget.NewLabel(
				"These subject codes are missing from the database: \n[" +
					strings.Join(missingSubjectCodes, ", ") + "]\n\n" +
					"Do you want to enter their names so it can be added to database?")

			dialog.ShowCustomConfirm("Missing Subject Codes", "Enter names", "Skip",
				dialogBoxContent, func(choice bool) {
					if choice {
						var subjectNames []*widget.Entry
						var formContent []*widget.FormItem

						for _, i := range missingSubjectCodes {
							nameEntry := widget.NewEntry()
							formContent = append(formContent,
								&widget.FormItem{Text: i, Widget: nameEntry})
							subjectNames = append(subjectNames, nameEntry)
						}

						dialog.ShowForm("Missing Subject Codes", "Save", "Skip",
							formContent, func(bool) {
								subjectCodes := make(map[string]string)
								for i, code := range missingSubjectCodes {
									subjectCodes[code] = strings.ToUpper((*subjectNames[i]).Text)
								}

								err = backend.Update(subjectCodes)
								if err != nil {
									dialog.ShowError(err, window)
									return
								} else {
									dialog.ShowInformation("Missing Subject Codes",
										"Subject Codes Updated Successfully!\nResult File is now Complete!",
										window)
									// Try again
									parsedData, _, missingSubjectCodes, err = backend.Parse(file.Path(), class)
									if err != nil {
										dialog.ShowError(err, window)
										return
									}
								}
							}, window)
					} else {
						dialog.ShowInformation("Missing Subject Codes",
							"ALERT: Few records will be missing from the result file!",
							window)
					}
				}, window)
		} else {
			dialog.ShowCustomConfirm("Save File", "Save File", "Discard",
				widget.NewLabel("File Converted Successfully to .csv!"), func(choice bool) {
					if choice {
						saveFileDialog := dialog.NewFileSave(func(callback fyne.URIWriteCloser, err error) {
							if err != nil {
								dialog.ShowError(err, window)
								return
							}
							if callback == nil {
								dialog.ShowError(errors.New("Select a file!"), window)
							} else {
								file = callback.URI()
								err = backend.Write(parsedData, class, file.Path())
								if err != nil {
									dialog.ShowError(err, window)
									return
								}
							}
						}, window)

						saveFileDialog.Resize(fyne.NewSize(600, 400))
						saveFileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))

						cwd, _ := os.Getwd()
						dir, _ := storage.ListerForURI(storage.NewFileURI(cwd))
						saveFileDialog.SetLocation(dir)

						saveFileDialog.Show()
					}
				}, window)
		}
	})

	convert := container.New(
		&layout.Vertical{Padding: 10},
		container.New(
			&layout.Horizontal{Padding: 32},
			widget.NewLabel("Select Class:"),
			classRadio,
		),
		container.New(
			&layout.Horizontal{},
			widget.NewLabel("Open Result File:"),
			container.New(
				&layout.Vertical{},
				fileNameLabel,
				openFileButton,
			),
		),
		convertButton,
	)

	analyse := container.New(
		&layout.Center{},
		widget.NewLabel("Convert a file first to see Analysis!"),
	)

	content := container.New(
		&layout.Vertical{},
		logo,
		container.New(
			&layout.Horizontal{Padding: 25},
			convert,
			widget.NewSeparator(),
			analyse,
		),
	)

	window.SetContent(content)
	window.ShowAndRun()
}
