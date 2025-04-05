package widgets

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/noodahl-org/erp/internal/models"
)

func LoginWidget(loginFunc, registerFunc func(models.User) error) *widget.Form {
	user := models.User{}

	usernamebind := binding.BindString(&user.Username)
	passwordbind := binding.BindString(&user.Password)
	password := widget.NewPasswordEntry()
	password.Bind(passwordbind)

	var errmsg string
	errbinding := binding.BindString(&errmsg)
	errlabel := widget.NewLabelWithStyle(errmsg, fyne.TextAlignCenter, fyne.TextStyle{
		Bold:      true,
		Monospace: true,
	})
	errlabel.Bind(errbinding)

	form := widget.NewForm(
		[]*widget.FormItem{
			{Text: "username", Widget: widget.NewEntryWithData(usernamebind)},
			{Text: "password", Widget: password},
			{Text: "", Widget: widget.NewButton("login", func() {
				user.Password = password.Text
				if err := loginFunc(user); err != nil {
					if err := errbinding.Set(errmsg); err != nil {
						log.Fatal(err) //todo handle errs in the app
					}
				}
			})},
			{Text: "", Widget: widget.NewButton("register", func() {
				if err := registerFunc(user); err != nil {
					//handle register error
				}
			})},
			{Text: "", Widget: errlabel},
		}...,
	)

	return form
}
