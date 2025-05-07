package mobileapp

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/x/fyne/theme"
	"github.com/noodahl-org/erp/api/models"
	"github.com/noodahl-org/erp/client"
	"github.com/noodahl-org/erp/fyne/internal/widgets"
)

type MobileApp struct {
	fyneapp   fyne.App
	user      *models.User
	dashboard *models.Dashboard
	Main      fyne.Window
	Nav       *container.AppTabs
	erp       client.ERPClient
}

func NewMobileApp(opts ...func(*MobileApp)) *MobileApp {
	app := app.New()
	app.Settings().SetTheme(theme.AdwaitaTheme())
	win := app.NewWindow("noodahl-erp")
	win.SetMaster()
	win.Resize(fyne.NewSize(400, 600))
	a := &MobileApp{
		fyneapp: app,
		Main:    win,
		erp: client.NewERPClient(
			client.WithBaseURL("http://localhost:8081"),
		),
	}
	a.Nav = widgets.NavigationTabs(a.NavigationHandler)

	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *MobileApp) ButtonHandler(index string) {
	switch index {
	case "equipment/catalog":
		a.Nav.SelectIndex(1)
		a.Main.SetContent(
			container.NewGridWithRows(2,
				a.ShowEquipmentCatalog(),
				a.Nav,
			),
		)
	case "equipment/view":
		a.Nav.SelectIndex(2)
		a.Main.SetContent(
			container.NewGridWithRows(2,
				a.UserEquipmentCatalog(),
				a.Nav,
			),
		)
	case "maintenance/view":
		a.Main.SetContent(
			container.NewGridWithRows(2),
		)
	case "maintenance/edit":
	}
}

func (a *MobileApp) NavigationHandler(t *container.TabItem) {
	switch t.Text {
	case "dashboard":
		a.ShowDashboard()
	case "equipment":
		a.Main.SetContent(
			container.NewGridWithRows(2,
				//fix this
				a.ShowEquipmentCatalog(),
				a.Nav,
			),
		)
	case "inventory":
		a.Main.SetContent(
			container.NewGridWithRows(2,
				a.UserEquipmentCatalog(),
				a.Nav,
			),
		)
	}

}
