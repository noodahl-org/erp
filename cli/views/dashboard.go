package views

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/noodahl-org/erp/api/models"
	"github.com/noodahl-org/erp/cli/views/components"
	"github.com/noodahl-org/erp/client"
)

type DashboardView struct {
	spinner  spinner.Model
	user     *models.User
	form     *huh.Form
	help     help.Model
	erp      client.ERPClient
	dash     *models.Dashboard
	navigate func(screen string, data interface{}) tea.Msg
	selected string
}

func NewDashboardView(user *models.User, navFunc func(string, interface{}) tea.Msg) tea.Model {
	log.Printf("dashboard:construct, %v, %v", user, navFunc)
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("205"))

	help := help.New()

	view := &DashboardView{
		spinner:  spin,
		user:     user,
		help:     help,
		navigate: navFunc,
	}
	return view
}

func (v *DashboardView) Init() tea.Cmd {
	log.Printf("dashboard:init, %v", v)
	//setup client
	v.erp = client.NewERPClient(
		client.WithBaseURL(os.Getenv("API_URL")),
	)
	//hydrate user
	var err error
	v.user, err = v.erp.FetchUser(*v.user)
	if err != nil {
		log.Fatal(err)
	}

	v.dash, err = v.erp.FetchDashboard(v.user.ID)
	if err != nil {
		log.Fatal(err)
	}

	//todo: template
	desc := fmt.Sprintf(`user: %s [%s]
	api url: %s
	user equipment: %v
	active tasks: %v
		`, components.Header.Render(v.user.Username), v.user.ID[:4],
		os.Getenv("API_URL"),
		len(v.dash.UserEquipment),
		len(v.dash.Tasks),
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Dashboard").
				Description(desc).Next(true),
		), //.Next(true),
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("nav").
				Options(
					huh.NewOption("Create Equipment", "equipment"),
					huh.NewOption("View Equipment", "user_equipment"),
					huh.NewOption("Tasks", "tasks"),
				).Value(&v.selected),
		),
	)

	form.SubmitCmd = func() tea.Msg {
		return v.navigate(v.selected, &ViewOpts{
			User: v.user,
			ERP:  v.erp,
		})
	}
	v.form = form
	return tea.Batch([]tea.Cmd{v.spinner.Tick, v.form.Init()}...)
}

func (v *DashboardView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	var cmd tea.Cmd
	v.spinner, cmd = v.spinner.Update(msg)
	cmds = append(cmds, cmd)

	form, cmd := v.form.Update(msg)
	if form, ok := form.(*huh.Form); ok {
		v.form = form
	}

	cmds = append(cmds, cmd)

	return v, tea.Batch(cmds...)
}

func (v *DashboardView) View() string {
	return v.form.View()
}
