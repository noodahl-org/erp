package widgets

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NavigationTabs(selected func(*container.TabItem)) *container.AppTabs {
	tabs := container.NewAppTabs(
		container.NewTabItem("dashboard", widget.NewLabel("")),
		container.NewTabItem("equipment", widget.NewLabel("")),
		container.NewTabItem("inventory", widget.NewLabel("")),
		container.NewTabItem("tasks", widget.NewLabel("")),
	)
	tabs.SetTabLocation(container.TabLocationBottom)
	tabs.OnSelected = selected
	return tabs
}
