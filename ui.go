package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	return model{spinner: s, selectedRequestIndex: -1, width: termWidth}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "up":
			if len(m.requests) != 0 && (m.selectedRequestIndex-1 >= 0) {
				m.selectedRequestIndex -= 1
			}
			return m, nil
		case "down":
			if len(m.requests) != 0 && (m.selectedRequestIndex+1) < len(m.requests) {
				m.selectedRequestIndex += 1
			}
			return m, nil
		default:
			return m, nil
		}

	case gotRequestMsg:
		m.requests = append(m.requests, msg.data)
		if len(m.requests) == 1 {
			m.selectedRequestIndex = 0
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
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
		var style lipgloss.Style
		if i == m.selectedRequestIndex {
			style = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("36"))
		} else {
			style = lipgloss.NewStyle()
		}
		reqsStr += style.Render(req.time.Format(time.TimeOnly) + " " + req.req.RequestURI)
		if i != len(m.requests)-1 {
			reqsStr += "\n"
		}
	}

	blockStyle := lipgloss.NewStyle().Padding(0, 1).BorderStyle(lipgloss.RoundedBorder())
	blockWidth := m.width/2 - 4

	var curReqStr string
	if m.selectedRequestIndex == -1 {
		curReqStr = "No requests"
	} else {
		curReqStr = renderRequest(m.requests[m.selectedRequestIndex], blockWidth)
	}

	res := strings.Builder{}

	block1 := blockStyle.Render(fmt.Sprintf("%s\n"+
		" %s Listening on %s\n\n"+
		"%s", strings.Repeat(" ", blockWidth), m.spinner.View(), GlobalState.listenAddress, reqsStr))

	block2 := blockStyle.Render(
		fmt.Sprintf("%s\n%s\n", strings.Repeat(" ", blockWidth), curReqStr))

	res.WriteString(
		lipgloss.JoinHorizontal(lipgloss.Top, block1, block2),
	)

	return res.String()
}
