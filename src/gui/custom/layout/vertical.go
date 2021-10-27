package layout

import (
	"fyne.io/fyne/v2"
	fyneLayout "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// I have no idea of majority of its working, but I just hacked together a way to set my
// desired padding in between the components

type Vertical struct {
	Padding float32
}

func isVerticalSpacer(obj fyne.CanvasObject) bool {
	if spacer, ok := obj.(fyneLayout.SpacerObject); ok {
		return spacer.ExpandVertical()
	}

	return false
}

func (v *Vertical) isSpacer(obj fyne.CanvasObject) bool {
	// invisible spacers don't impact layout
	if !obj.Visible() {
		return false
	}
	return isVerticalSpacer(obj)
}

func (v *Vertical) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := make([]fyne.CanvasObject, 0)
	total := float32(0)
	padding := theme.Padding() + v.Padding

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if v.isSpacer(child) {
			spacers = append(spacers, child)
			continue
		}
		total += child.MinSize().Height
	}

	x, y := float32(0), float32(0)
	var extra float32
	extra = size.Height - total - (padding * float32(len(objects)-len(spacers)-1))

	extraCell := float32(0)
	if len(spacers) > 0 {
		extraCell = extra / float32(len(spacers))
	}

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if v.isSpacer(child) {
			y += extraCell
			continue
		}
		child.Move(fyne.NewPos(x, y))
		child.Resize(fyne.NewSize(size.Width, child.MinSize().Height))

		y += padding + child.MinSize().Height
	}
}

func (v *Vertical) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	padding := theme.Padding() + v.Padding
	addPadding := false

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if v.isSpacer(child) {
			continue
		}

		minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
		minSize.Height += child.MinSize().Height
		if addPadding {
			minSize.Height += padding
		}

		addPadding = true
	}
	return minSize
}
