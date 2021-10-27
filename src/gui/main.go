package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	assets "github.com/luciferchase/CBSE-Result-Converter/src/gui/assets"
	content "github.com/luciferchase/CBSE-Result-Converter/src/gui/content"
	layout "github.com/luciferchase/CBSE-Result-Converter/src/gui/custom/layout"
)

func main() {
	// I tried "Result Converter" but it crashes and I don't like result.converter
	// hence using "Result" for now.
	// Not sure if this is a bug or a feature.
	app := app.NewWithID("Result")

	window := app.NewWindow("Result Converter")
	window.CenterOnScreen()
	// I have opted instead of letting the app decide the dimension, fixing it beforehand
	// this way I could ensure every thing looks nice and easy considering the end users
	// are going to be mostly window owners with majority of them having 1366x768 resolution.
	// I know this is not the most optimum way, but I like this.
	window.Resize(fyne.NewSize(580, 350))
	window.SetFixedSize(true)
	// Used Flaticon (https://www.flaticon.com) for the icon
	window.SetIcon(assets.Icon)

	// Used Logojoy (https://logojoy.com/dashboard) for the inspiration for the logo
	logo := canvas.NewImageFromResource(assets.Logo)
	logo.SetMinSize(fyne.NewSize(570, 100))

	appTabs := mainScreen(window)

	// Design for the main screen is really simple having only 2 components
	// 1. Logo - panning horizontally at the top
	// 2. Tabs - I. Convert Tab
	//					a) Class Radio (Required)
	//					b) Browse File (Required)
	// 					c) Convert Button
	// 			 II. Analysis Tab
	// 					a) School Result (Everything related to the school performance)
	// 					b) Topper Performance
	//					c) Worst Performance
	content := container.NewVBox(
		logo,
		appTabs,
	)

	window.SetContent(content)
	window.ShowAndRun()
}

func mainScreen(window fyne.Window) fyne.CanvasObject {
	// Initialise components
	classRadio, fileNameLabel, openFileButton,
		convertButton, analyse := content.ConvertTab(window)
	classRadio.Required = true

	// Convert Tab
	formItems := []*widget.FormItem{
		{
			Text:   "Class",
			Widget: classRadio,
		},
		{
			Text:     "Input File",
			Widget:   container.NewHBox(openFileButton, fileNameLabel),
			HintText: "Select a valid CBSE issued result text file",
		},
	}
	form := widget.NewForm(formItems...)

	// Used my own layout for better control over the padding between the components
	convert := container.New(
		&layout.Vertical{Padding: 20},
		form,
		convertButton,
	)

	// Workaround to set tabs to the middle of the screen
	// \t produces square like things only
	// One of the reason for fixing the size of the window.
	appTabs := container.NewAppTabs(
		container.NewTabItem("						    Convert						    ", convert),
		container.NewTabItem("						    Analyse						    ", analyse),
	)
	appTabs.SetTabLocation(container.TabLocationTop)

	return appTabs
}
