package components

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Spinner struct {
	Spinner spinner.Model
	Err     error
	Msg     string
}

func NewSpinner(msg string) *Spinner {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &Spinner{
		Spinner: sp,
		Msg:     msg,
	}
}

func (s *Spinner) Update(msg tea.Msg) (spinner.Model, tea.Cmd) {
	return s.Spinner.Update(msg)
}

func (s *Spinner) View() string {
	return s.Spinner.View()
}
