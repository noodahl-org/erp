package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/noodahl-org/erp/api/models"
	"github.com/noodahl-org/erp/cli/views"
)

type erpCli struct {
	Views       map[string]tea.Model
	currentView string
}

func (e erpCli) Init() tea.Cmd {
	log.Println("erp cli init")
	return tea.Batch(
		e.Views[e.currentView].Init(),
	)
}

func (e erpCli) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return e, tea.Quit
		}
	case NavMsg:
		e.currentView = msg.Screen
		switch msg.Screen {
		case "dashboard":
			log.Println("nav:dashboard")
			if user, ok := msg.Data.(*models.User); ok {
				e.Views["dashboard"] = views.NewDashboardView(user, Navigate)
			}
		case "equipment":
			log.Println("nav:equipment")
			if opts, ok := msg.Data.(*views.ViewOpts); ok {
				e.Views["equipment"] = views.NewEquipmentView(opts, Navigate)
			}
		}

		// Initialize the new view
		return e, e.Views[e.currentView].Init()
	}

	if val, ok := e.Views[e.currentView]; ok {
		newModel, cmd := val.Update(msg)

		// Update the view in our map if it's changed
		if newViewModel, ok := newModel.(tea.Model); ok {
			e.Views[e.currentView] = newViewModel
		}

		return e, cmd
	}
	return e, nil
}

func (e erpCli) View() string {
	if val, ok := e.Views[e.currentView]; ok {
		return val.View()
	}
	return ""
}

func NewErpCLI(conf map[string]string) erpCli {
	cli := erpCli{
		currentView: "login",
	}
	cli.Views = map[string]tea.Model{
		"login": views.NewLoginView(conf, Navigate),
	}
	return cli
}

func main() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	conf := map[string]string{
		"API_URL":         os.Getenv("API_URL"),
		"KUBE_HOST":       os.Getenv("KUBE_HOST"),
		"KUBE_USERNAME":   os.Getenv("KUBE_USERNAME"),
		"KUBE_PASS":       os.Getenv("KUBE_PASS"),
		"KUBE_KEY_BASE64": os.Getenv("KUBE_KEY_BASE64"),
	}

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(NewErpCLI(conf))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type NavMsg struct {
	Screen string
	Data   interface{}
}

func Navigate(screen string, data interface{}) tea.Msg {
	return NavMsg{
		Screen: screen,
		Data:   data,
	}
}
