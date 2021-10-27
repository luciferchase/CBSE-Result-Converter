package tabs

import (
	"errors"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	nativeDialog "github.com/sqweek/dialog"

	backend "github.com/luciferchase/CBSE-Result-Converter/src/backend"
)

// Logic behind the main screen
//
// ConvertTab defines all the components in the main screen which is initialised in the main.go and rendered there.

func ConvertTab(window fyne.Window) (classRadio *widget.RadioGroup, fileNameLabel *widget.Label,
	openFileButton, convertButton *widget.Button, analysisTab *container.Scroll) {
	var err error

	// Initial analysis content
	initialContent := widget.NewLabelWithStyle(
		"Convert a file first to see Analysis!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)
	analysisContent := container.NewMax(
		initialContent,
	)
	analysisTab = container.NewVScroll(analysisContent)

	// Select Class
	var class string
	classRadio = widget.NewRadioGroup([]string{"Class 10", "Class 12"},
		func(value string) {
			// Backend takes value of class as "X" and "XII" only
			// This was because of TUI
			if value == "Class 10" {
				class = "X"
			} else if value == "Class 12" {
				class = "XII"
			}
		},
	)

	fileNameLabel = widget.NewLabel("No file selected!")

	var path string
	openFileButton = widget.NewButtonWithIcon("Browse", theme.FileIcon(), func() {
		// Fyne gives api for dialog boxes but in their attempt for cross-platform uniformed
		// dialog boxes, they made it imho ugly.
		// Hence using the awesome sqweek/dialog library for native dialog boxes
		//
		// Select only .txt files
		path, err = nativeDialog.File().Filter("Text File (*.txt, *.text)", "txt").Load()
		if err != nil {
			dialog.ShowError(errors.New("Invalid File!"), window)
			return
		} else if path == "" {
			dialog.ShowError(errors.New("No file selected!"), window)
			return
		}
		// Change fileNameLabel so that user knows which file they selected
		fileNameLabel.SetText(filepath.Base(path))
	})

	convertButton = widget.NewButtonWithIcon("Convert to CSV", theme.StorageIcon(), func() {
		// First check if class is selected or not
		if len(class) == 0 {
			dialog.ShowError(errors.New("Select a class first!"), window)
			return
		}

		parsedData, schoolResult, missingSubjectCodes, err := backend.Parse(path, class)
		if err != nil {
			dialog.ShowError(err, window)
			return
		}

		if len(missingSubjectCodes) != 0 {
			parsedData = ifSubjectCodesMissing(window, missingSubjectCodes, path, class)
		} else {
			// Update Analysis tab
			// Whenever convert button is clicked, it will fetch fresh analysis from the
			// AnalysisTab function and update the analysis tab.
			accordion := AnalysisTab(schoolResult)
			accordion.Open(1)

			// Previously, if after converting a file, another file is converted then,
			// the Analysis Tab didn't used to clear and the new analysis was rendered in the
			// same tab as well. Pretty ugly.
			//
			// To fix that, I redifened analysisContent and updated the analysisTab struct.
			// I cannot redefine analysisTab because I am returning its pointer,
			// and If I do container.NewVScroll() no content will be updated in the Analysis Tab.

			analysisContent = container.NewMax(accordion)
			analysisTab.Content = analysisContent

			dialog.ShowCustomConfirm("Save File", "Save File", "Discard",
				widget.NewLabel("File Converted Successfully to .csv!"), func(choice bool) {
					if choice {
						path, err = nativeDialog.File().Filter("XLS Worksheet (*.csv)",
							"csv").Title("Save File").Save()
						if err != nil {
							dialog.ShowError(errors.New("Invalid File!"), window)
							return
						} else if path == "" {
							dialog.ShowError(errors.New("No file selected!"), window)
							return
						}

						err = backend.Write(parsedData, class, path)
						if err != nil {
							dialog.ShowError(errors.New("An Error occurred! Please try again"),
								window)
							return
						}

						// Send notification
						//
						// Currently, it doesn't send app icon in the notification.
						// It is a known issue in fyne (https://github.com/fyne-io/fyne/issues/2592)
						fyne.CurrentApp().SendNotification(&fyne.Notification{
							Title:   "Success",
							Content: "File saved Successfully!",
						})
					}
				}, window)
		}
		// Refresh everything
		classRadio.SetSelected("")
		fileNameLabel.SetText("No file selected!")
	})
	return
}

// If in edge-cases that some subject-codes are missing from the core group of subject-codes,
// Instead of panic, it gracefully handles the problem.
func ifSubjectCodesMissing(window fyne.Window, missingSubjectCodes []string,
	path, class string) (parsedData [][]string) {
	var err error

	dialogBoxContent := widget.NewLabel(
		"These subject codes are missing from the database: \n[" +
			strings.Join(missingSubjectCodes, ", ") + "]\n\n" +
			"Do you want to enter their names so it can be added to database?",
	)

	dialog.ShowCustomConfirm("Missing Subject Codes", "Enter names", "Skip",
		dialogBoxContent, func(choice bool) {
			// If they want to enter name of the missing subject codes manually
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
							parsedData, _, missingSubjectCodes, err = backend.Parse(path, class)
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
	return
}
