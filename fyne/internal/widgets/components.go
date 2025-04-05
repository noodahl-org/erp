package widgets

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type Listable interface {
	ListLabel() string
}

func ToListWidget[T Listable](data []T) *widget.List {
	list := widget.NewList(
		//length of the list
		func() int {
			return len(data)
		},
		//create item
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(index widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(data[index].ListLabel())
		},
	)
	return list
}
