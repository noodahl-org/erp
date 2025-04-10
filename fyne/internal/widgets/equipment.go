package widgets

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/noodahl-org/erp/api/models"
)

func EquipmentView(submit func()) *widget.Form {
	equipment := models.Equipment{}
	make := binding.BindString(&equipment.Make)
	model := binding.BindString(&equipment.Model)

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "make", Widget: widget.NewEntryWithData(make)},
			{Text: "model", Widget: widget.NewEntryWithData(model)},
		},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", equipment.Make)
			log.Println("multiline:", equipment.Model)
		},
	}

	return form
}

func EquipmentList() *fyne.Container {
	return container.New(
		layout.NewGridLayoutWithRows(2),
		widget.NewLabel("test"),
	)
}
