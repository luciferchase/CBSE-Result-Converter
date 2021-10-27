package layout

import (
	"fyne.io/fyne/v2"
	fyneLayout "fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// I have no idea of majority of its working, but I just hacked together a way to set my
// desired padding in between the components

type Horizontal struct {
	Padding float32
}

func isHorizontalSpacer(obj fyne.CanvasObject) bool {
	if spacer, ok := obj.(fyneLayout.SpacerObject); ok {
		return spacer.ExpandHorizontal()
	}

	return false
}

func (h *Horizontal) isSpacer(obj fyne.CanvasObject) bool {
	// invisible spacers don't impact layout
	if !obj.Visible() {
		return false
	}
	return isHorizontalSpacer(obj)
}

func (h *Horizontal) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := make([]fyne.CanvasObject, 0)
	total := float32(0)
	padding := theme.Padding() + h.Padding

	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if h.isSpacer(child) {
			spacers = append(spacers, child)
			continue
		}

		total += child.MinSize().Width
	}

	x, y := float32(0), float32(0)

	var extra float32
	extra = size.Width - total - (padding * float32(len(objects)-len(spacers)-1))

	extraCell := float32(0)
	if len(spacers) > 0 {
		extraCell = extra / float32(len(spacers))
	}

	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if h.isSpacer(child) {
			x += extraCell
			continue
		}
		child.Move(fyne.NewPos(x, y))
		child.Resize(fyne.NewSize(child.MinSize().Width, size.Height))

		x += padding + child.MinSize().Width
	}
}

func (h *Horizontal) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	padding := theme.Padding() + h.Padding
	addPadding := false

	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		if h.isSpacer(child) {
			continue
		}

		minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
		minSize.Width += child.MinSize().Width
		if addPadding {
			minSize.Width += padding
		}

		addPadding = true
	}
	return minSize
}
