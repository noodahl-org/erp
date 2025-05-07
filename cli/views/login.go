package views

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/noodahl-org/erp/api/models"
)

type LoginView struct {
	spinner spinner.Model
	conf    map[string]string
	user    *models.User
	form    *huh.Form
}

func (s *LoginView) Init() tea.Cmd {
	return tea.Batch([]tea.Cmd{
		s.form.Init(),
		s.spinner.Tick,
	}...)
}

func NewLoginView(conf map[string]string, navigate func(screen string, data interface{}) tea.Msg) tea.Model {
	user := &models.User{}
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	view := &LoginView{
		spinner: spin,
		user:    user,
		conf:    conf,
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewNote().Title("Login"),
				huh.NewInput().Prompt("username: ").
					Value(&user.Username),
				huh.NewInput().EchoMode(huh.EchoModePassword).
					Prompt("password: ").
					Value(&user.Password),
			),
		),
	}
	view.form.SubmitCmd = func() tea.Msg {
		return navigate("dashboard", view.user)
	}
	return view
}

func (s *LoginView) Name() string {
	return "login"
}

func (s *LoginView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var form tea.Msg
	cmds := []tea.Cmd{}
	form, cmd = s.form.Update(msg)
	if form, ok := form.(*huh.Form); ok {
		s.form = form
	}

	cmds = append(cmds, cmd)

	s.spinner, cmd = s.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

func (s *LoginView) View() string {
	if s.form.State == huh.StateCompleted {
		return s.spinner.View()
	}
	return s.form.View()
}
