package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	"github.com/noodahl-org/erp/cli/views"
)

type erpCli struct {
	Views       map[string]views.View
	currentView string
}

func (e erpCli) Init() tea.Cmd {
	log.Println("erp cli init")
	return tea.Batch(
		e.Views[e.currentView].Init()...,
	)
}

func (e erpCli) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return e, tea.Quit
		}
	}

	// Pass the command through to our current view
	if val, ok := e.Views[e.currentView]; ok {
		return e, val.Update(msg)
	}

	return e, nil
}

func (e erpCli) View() string {
	doc := strings.Builder{}
	if val, ok := e.Views[e.currentView]; ok {
		val.View(&doc)
	}
	return doc.String()
}

func NewErpCLI(conf map[string]string) erpCli {
	return erpCli{
		currentView: "system",
		Views: map[string]views.View{
			"system": views.NewSystemsView(conf),
		},
	}
}

func main() {
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
