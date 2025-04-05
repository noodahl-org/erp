package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	mobileapp "github.com/noodahl-org/erp/fyne/internal/app"
	"github.com/noodahl-org/erp/fyne/internal/widgets"
)

func main() {

	app := mobileapp.NewMobileApp()

	app.Main.SetContent(
		container.New(
			layout.NewAdaptiveGridLayout(1),
			container.New(
				layout.NewCustomPaddedLayout(15, 15, 15, 100),
				widgets.LoginWidget(app.Login, app.Register),
			),
		),
	)
	app.Main.ShowAndRun()
}
