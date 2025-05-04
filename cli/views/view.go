package views

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type View interface {
	Update(msg tea.Msg) tea.Cmd
	Name() string
	View(out *strings.Builder)
	Init() []tea.Cmd
}
