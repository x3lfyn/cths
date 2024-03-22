package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
	"net/http"
	"os"
	"strings"
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

	blocKStyle := lipgloss.NewStyle().Padding(0, 1).BorderStyle(lipgloss.RoundedBorder())

	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))

	res := strings.Builder{}

	block1 := blocKStyle.Render(fmt.Sprintf("%s\n"+
		" %s Listening to requests\n\n"+
		"Requests:\n"+
		"%s", strings.Repeat(" ", termWidth/2-4), m.spinner.View(), reqsStr))

	block2 := blocKStyle.Render(
		fmt.Sprintf("%s\nfuckfuckfuck\n", strings.Repeat(" ", termWidth/2-4)))

	res.WriteString(
		lipgloss.JoinHorizontal(lipgloss.Top, block1, block2),
	)

	return res.String()
}

func main() {
	p := tea.NewProgram(initialModel())

	go listener(p)

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
