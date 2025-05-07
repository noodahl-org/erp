package views

import (
	"log"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/noodahl-org/erp/api/models"
	"github.com/noodahl-org/erp/client"
)

type EquipmentView struct {
	erp         client.ERPClient
	user        *models.User
	form        *huh.Form
	equip       *models.Equipment
	userequip   *models.UserEquipment
	year        string
	runWorkflow bool
	navigate    func(string, interface{}) tea.Msg
}

func NewEquipmentView(opts *ViewOpts, navFunc func(string, interface{}) tea.Msg) tea.Model {
	log.Printf("equipment:construct, %v", opts)
	view := &EquipmentView{
		erp:       opts.ERP,
		user:      opts.User,
		equip:     &models.Equipment{},
		userequip: &models.UserEquipment{},
		navigate:  navFunc,
	}
	return view
}

func (v *EquipmentView) Init() tea.Cmd {
	v.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Description("make").Value(&v.equip.Make),
			huh.NewInput().Description("model").Value(&v.equip.Model),
			huh.NewText().Description("description").Value(&v.equip.Description),
			huh.NewInput().Description("serial").Value(&v.userequip.SerialNumber),
			huh.NewInput().Description("year").Value(&v.year),
		).Title("New Equipment"),
		huh.NewGroup(
			huh.NewConfirm().Description("Run NewEquipmentWorkflow?").Value(&v.runWorkflow),
		),
	)
	v.form.SubmitCmd = func() tea.Msg {
		var err error
		v.equip, err = v.erp.UpsertEquipment(*v.equip)
		if err != nil {
			log.Fatal(err)
		}
		v.userequip.UserID = v.user.ID
		v.userequip.EquipmentID = v.equip.ID
		year, _ := strconv.Atoi(v.year)
		v.userequip.Year = year
		v.userequip, err = v.erp.UpserUserEquipment(*v.userequip)
		if err != nil {
			log.Fatal(err)
		}

		return v.navigate("dashboard", v.user)
	}
	return nil
}

func (v *EquipmentView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	var cmd tea.Cmd
	form, cmd := v.form.Update(msg)
	if form, ok := form.(*huh.Form); ok {
		v.form = form
	}
	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

func (v EquipmentView) View() string {
	return v.form.View()
}
