package layout

import (
	"fyne.io/fyne/v2"
)

type Center struct {
}

func (c *Center) MinSize(objects []fyne.CanvasObject) fyne.Size {
	child := objects[0]
	minSize := fyne.NewSize(child.MinSize().Width, child.MinSize().Height)
	return minSize
}

func (c *Center) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	child := objects[0]
	x := (containerSize.Width - child.MinSize().Width) / 2
	y := (containerSize.Height - child.MinSize().Height) / 2
	pos := fyne.NewPos(x, y)
	child.Move(pos)
}
