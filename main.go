package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"net/http"
	"os"
)

type gotRequestMsg struct{ data *http.Request }

type model struct {
	spinner  spinner.Model
	requests []*http.Request
	err      error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}

	case gotRequestMsg:
		m.requests = append(m.requests, msg.data)
		return m, nil

	case error:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	var reqsStr string
	for i, req := range m.requests {
		reqsStr += renderRequest(req)
		if i != len(m.requests)-1 {
			reqsStr += "\n\n"
		}
	}

	res := fmt.Sprintf("\n"+
		" %s Listening to requests\n\n"+
		"Requests:\n"+
		"%s", m.spinner.View(), reqsStr)

	return res
}

func main() {
	p := tea.NewProgram(initialModel())

	go listener(p)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
