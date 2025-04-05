package mobileapp

import (
	"github.com/noodahl-org/erp/internal/models"
)

func (a *MobileApp) Login(user models.User) error {
	//todo actual auth
	var err error
	a.user, err = a.erp.FetchUser(user)
	if err != nil {
		return err
	}

	a.ShowDashboard()

	return nil
}

func (a *MobileApp) Register(user models.User) error {
	return nil
}
