package views

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"go.temporal.io/sdk/client"
	"resty.dev/v3"
)

type status struct {
	index int
	msg   string
}

type SystemsView struct {
	sub      chan status
	conf     map[string]string
	input    bytes.Buffer
	statuses []string
	resty    *resty.Client
	msg      string
}

func NewSystemsView(conf map[string]string) *SystemsView {
	return &SystemsView{
		conf:  conf,
		sub:   make(chan status),
		resty: resty.New(),
		msg:   "validationg resources...",
		statuses: []string{
			"erp-api/health: ",
			"temporal: ",
		},
	}
}

func (s *SystemsView) Name() string {
	return "system"
}

func (s *SystemsView) wait() tea.Cmd {
	return func() tea.Msg {
		return status(<-s.sub)
	}
}

func (s *SystemsView) Update(msg tea.Msg) tea.Cmd {
	var cmds []tea.Cmd

	switch m := msg.(type) {
	case status:
		s.statuses[m.index] += m.msg
		cmds = append(cmds, s.wait())
	case tea.KeyMsg:
		switch m.Type {
		case tea.KeyBackspace:
			s.input.Truncate(s.input.Len() - 1)
		default:
			s.input.WriteString(m.String())
		}
	}

	return tea.Batch(cmds...)
	//do stuff with input
}

func (s *SystemsView) View(out *strings.Builder) {
	out.WriteString(fmt.Sprintf("[%s]\n", s.Name()))
	out.WriteString(fmt.Sprintf(" - %s\n", s.msg))
	for _, s := range s.statuses {
		out.WriteString(fmt.Sprintf("%s\n", s))
	}
	out.WriteString(s.input.String())
}

func (s *SystemsView) Init() []tea.Cmd {
	return []tea.Cmd{
		func() tea.Msg {
			if val, ok := s.conf["API_URL"]; ok {
				url := "http://" + val + "/health"
				res, err := s.resty.R().
					Get(url)
				if err != nil {
					s.sub <- status(status{index: 0, msg: err.Error()})
					return nil
				} else if res.Err != nil {
					s.sub <- status(status{index: 0, msg: res.Err.Error()})
					return nil
				}
				s.sub <- status(status{index: 0, msg: res.Status()})
			}
			return s.wait()
		},
		func() tea.Msg {
			if val, ok := s.conf["TEMPORAL_HOST"]; ok {
				log.Println("temporal host", s.conf["TEMPORAL_HOST"])
				_, err := client.Dial(client.Options{
					HostPort: val,
				})
				if err != nil {
					s.sub <- status(status{index: 1, msg: err.Error()})
					return nil
				}
				s.sub <- status(status{index: 1, msg: "OK"})
			}
			return s.wait()
		},
		s.wait(),
		//func() tea.Msg {
		//todo http call
		//s.sub <- status("test")
		//return nil
		//},
	}
}
