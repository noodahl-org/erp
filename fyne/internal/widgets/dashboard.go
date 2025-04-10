package widgets

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/noodahl-org/erp/api/models"
)

func DashboardWidget(dashboard models.Dashboard, handler func(string)) *fyne.Container {
	return container.New(
		layout.NewAdaptiveGridLayout(2),
		widget.NewCard("Equipment", fmt.Sprintf("%v", len(dashboard.UserEquipment)), container.New(
			layout.NewGridLayoutWithColumns(2),
			widget.NewButton("Catalog", func() {
				handler("equipment/catalog")
			}),
			widget.NewButton("View", func() {
				handler("equipment/view")
			}),
		)),
		widget.NewCard("Maintenance", fmt.Sprintf("%v Open Task(s)", len(dashboard.Tasks)), container.New(
			layout.NewGridLayoutWithColumns(2),
			widget.NewButton("View", func() {
				handler("maintenance/view")
			}),
			widget.NewButton("Edit", func() {
				handler("maintenance/edit")
			}),
		)),

		widget.NewCard("Documents", "", widget.NewButton("", func() {
			handler("test/test")
		})),
	)
}
