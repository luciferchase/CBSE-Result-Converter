package tabs

import (
	"strconv"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	backend "github.com/luciferchase/CBSE-Result-Converter/src/backend"
	layout "github.com/luciferchase/CBSE-Result-Converter/src/gui/custom/layout"
)

const padding = 50

// Returns all the components for the Analysis Tab
func AnalysisTab(schoolResult backend.Result) (accordion *widget.Accordion) {
	// This is a very weird way to do this and needs major refactoring :/
	//
	// Basically, since School Result has 6 components, I am dividing them
	// into 2 columns with 3 each and using Forms because it lines them up neatly.
	formItems := []*widget.FormItem{
		{Text: "Total:", Widget: widget.NewLabel(schoolResult.Total)},
		{Text: "Passed:", Widget: widget.NewLabel(schoolResult.Passed)},
		{Text: "Compartment:", Widget: widget.NewLabel(schoolResult.Comptt)},
	}
	form1 := widget.NewForm(formItems...)
	formItems = []*widget.FormItem{
		{Text: "Essential Repeat:", Widget: widget.NewLabel(schoolResult.Repeat)},
		{Text: "Absent:", Widget: widget.NewLabel(schoolResult.Absent)},
		{Text: "Other:", Widget: widget.NewLabel(schoolResult.Other)},
	}
	form2 := widget.NewForm(formItems...)
	basicResult := widget.NewAccordionItem(
		"School Result",
		container.New(
			&layout.Horizontal{Padding: padding},
			form1,
			form2,
		),
	)

	formItems = []*widget.FormItem{
		{Text: "Roll:", Widget: widget.NewLabel(schoolResult.Topper.Roll)},
		{Text: "Name:", Widget: widget.NewLabel(schoolResult.Topper.Name)},
	}
	form1 = widget.NewForm(formItems...)
	formItems = []*widget.FormItem{
		{
			Text: "Percentage:",
			// 'f' is the litreal for nice decimal notation, 2 is the no. of significant digits after .
			// and 64 is the bit size of Float used.
			Widget: widget.NewLabel(strconv.FormatFloat(schoolResult.Topper.Percentage, 'f', 2, 64)),
		},
	}
	form2 = widget.NewForm(formItems...)
	topper := widget.NewAccordionItem(
		"Topper",
		container.New(
			&layout.Horizontal{Padding: padding},
			form1,
			form2,
		),
	)

	formItems = []*widget.FormItem{
		{Text: "Roll:", Widget: widget.NewLabel(schoolResult.LowestMarks.Roll)},
		{Text: "Name:", Widget: widget.NewLabel(schoolResult.LowestMarks.Name)},
	}
	form1 = widget.NewForm(formItems...)
	formItems = []*widget.FormItem{
		{
			Text:   "Percentage:",
			Widget: widget.NewLabel(strconv.FormatFloat(schoolResult.LowestMarks.Percentage, 'f', 2, 64)),
		},
	}
	form2 = widget.NewForm(formItems...)
	lowest := widget.NewAccordionItem(
		"Worst Performer",
		container.New(
			&layout.Horizontal{Padding: padding},
			form1,
			form2,
		),
	)

	accordion = widget.NewAccordion(
		basicResult,
		topper,
		lowest,
	)

	return
}
