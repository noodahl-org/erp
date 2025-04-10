package mobileapp

import (
	"log"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"github.com/noodahl-org/erp/fyne/internal/widgets"
	"github.com/noodahl-org/erp/api/models"
)

func (a *MobileApp) ShowDashboard() {
	var err error
	dashboard := models.Dashboard{
		User:          *a.user,
		UserEquipment: []models.UserEquipment{},
	}

	if dashboard.UserEquipment, err = a.erp.FetchUserEquipment(models.UserEquipment{UserID: a.user.ID}); err != nil {
		log.Fatal(err) //todo handle errors better
	}

	if dashboard.Tasks, err = a.erp.FetchUserMaintenanceTasks(models.UserMaintenanceTask{UserID: a.user.ID}); err != nil {
		log.Fatal(err)
	}

	a.dashboard = &dashboard

	a.Main.SetContent(container.New(
		layout.NewGridLayoutWithRows(2),
		widgets.DashboardWidget(dashboard, a.ButtonHandler),
		a.Nav,
	))
}
