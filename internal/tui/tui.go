package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"cths/pkg/bus"
)

type TUI struct {
	program *tea.Program
	model   Model
}

type RequestItem struct {
	payload   bus.RequestMessagePayload
	timestamp time.Time
}

type PanelSelected int

const (
	leftPanelSelected PanelSelected = iota
	rightPanelSelected
)

type Model struct {
	messageBus      *bus.MessageBus
	requests        []RequestItem
	selectedRequest int
	selectedPanel   PanelSelected
	leftPanel       viewport.Model
	rightPanel      viewport.Model
	ready           bool
	width           int
	height          int
}

func NewModel(messageBus *bus.MessageBus) Model {
	return Model{
		messageBus:      messageBus,
		requests:        make([]RequestItem, 0),
		selectedRequest: -1,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.waitForRequests(),
		tea.EnterAltScreen,
	)
}

func (m Model) waitForRequests() tea.Cmd {
	return func() tea.Msg {
		ch := m.messageBus.Subscribe(bus.RequestMessage)
		for msg := range ch {
			if payload, ok := msg.Payload.(bus.RequestMessagePayload); ok {
				return requestReceivedMsg{
					Payload:   payload,
					Timestamp: msg.Timestamp,
				}
			}
		}
		return nil
	}
}

type requestReceivedMsg struct {
	Payload   bus.RequestMessagePayload
	Timestamp time.Time
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.ready {
			m.width = msg.Width
			m.height = msg.Height
			m.leftPanel = viewport.New(msg.Width/3, msg.Height-6)
			m.rightPanel = viewport.New(msg.Width*2/3, msg.Height-6)
			m.ready = true
		} else {
			m.leftPanel.Width = msg.Width / 3
			m.leftPanel.Height = msg.Height - 6
			m.rightPanel.Width = msg.Width * 2 / 3
			m.rightPanel.Height = msg.Height - 6
		}
		m.updateLeftPanel()
		m.updateRightPanel()

	case requestReceivedMsg:
		item := RequestItem{
			payload:   msg.Payload,
			timestamp: msg.Timestamp,
		}
		m.requests = append(m.requests, item)
		if m.selectedRequest == -1 {
			m.selectedRequest = 0
		}
		m.updateLeftPanel()
		m.updateRightPanel()
		return m, m.waitForRequests()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.selectedRequest > 0 {
				switch m.selectedPanel {
				case leftPanelSelected:
					m.selectedRequest--
					m.updateLeftPanel()
					m.updateRightPanel()
				case rightPanelSelected:
					m.updateRightPanel()
				}
			}
		case "down", "j":
			if m.selectedRequest < len(m.requests)-1 {
				switch m.selectedPanel {
				case leftPanelSelected:
					m.selectedRequest++
					m.updateLeftPanel()
					m.updateRightPanel()
				case rightPanelSelected:
					m.updateRightPanel()
				}
			}
		case "tab":
			m.selectedPanel = (m.selectedPanel + 1) % 2
		}
	}

	var cmd tea.Cmd
	switch m.selectedPanel {
	case leftPanelSelected:
		m.leftPanel, cmd = m.leftPanel.Update(msg)
	case rightPanelSelected:
		m.rightPanel, cmd = m.rightPanel.Update(msg)
	}
	return m, cmd

}

func (m *Model) updateLeftPanel() {
	if !m.ready {
		return
	}

	var content strings.Builder

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	content.WriteString(headerStyle.Render("REQUESTS"))
	content.WriteString("\n")

	if len(m.requests) == 0 {
		content.WriteString(lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Render("  No requests yet..."))
		m.leftPanel.SetContent(content.String())
		return
	}

	for i, item := range m.requests {
		req := item.payload

		indicator := " "
		if i == m.selectedRequest {
			indicator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("51")).
				Bold(true).
				Render("*")
		}

		// Method with cyan/pink accents
		methodStyle := lipgloss.NewStyle().Bold(true)
		switch req.Method {
		case "GET":
			methodStyle = methodStyle.Foreground(lipgloss.Color("51")) // Cyan
		case "POST":
			methodStyle = methodStyle.Foreground(lipgloss.Color("205")) // Pink
		case "PUT", "PATCH":
			methodStyle = methodStyle.Foreground(lipgloss.Color("51")) // Cyan
		case "DELETE":
			methodStyle = methodStyle.Foreground(lipgloss.Color("205")) // Pink
		default:
			methodStyle = methodStyle.Foreground(lipgloss.Color("15")) // White
		}
		method := methodStyle.Render(fmt.Sprintf("%-6s", req.Method))

		// URL path (truncate if too long)
		path := req.URL.Path
		if req.URL.RawQuery != "" {
			path += "?" + req.URL.RawQuery
		}

		// Truncate path if too long
		maxPathWidth := m.leftPanel.Width - 24 - len(req.Method)
		if len(path) > maxPathWidth && maxPathWidth > 3 {
			path = path[:maxPathWidth-3] + "..."
		}

		// Time (compact format)
		timeStr := item.timestamp.Format("15:04:05")

		// Build compact single line
		timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
		content.WriteString(fmt.Sprintf("%s %s %s %s\n",
			indicator,
			method,
			lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Width(maxPathWidth).Render(path),
			timeStyle.Render(timeStr)))
	}

	m.leftPanel.SetContent(content.String())
}

func (m *Model) updateRightPanel() {
	if !m.ready {
		return
	}

	if m.selectedRequest < 0 || m.selectedRequest >= len(m.requests) {
		m.rightPanel.SetContent(lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Render("No request selected"))
		return
	}

	item := m.requests[m.selectedRequest]
	req := item.payload

	var content strings.Builder

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205"))

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")).
		Bold(true)

	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))

	methodStyle := lipgloss.NewStyle().Bold(true)
	switch req.Method {
	case "GET":
		methodStyle = methodStyle.Foreground(lipgloss.Color("51"))
	case "POST":
		methodStyle = methodStyle.Foreground(lipgloss.Color("205"))
	case "PUT", "PATCH":
		methodStyle = methodStyle.Foreground(lipgloss.Color("51"))
	case "DELETE":
		methodStyle = methodStyle.Foreground(lipgloss.Color("205"))
	default:
		methodStyle = methodStyle.Foreground(lipgloss.Color("15"))
	}

	content.WriteString(methodStyle.Render(req.Method))

	fullURL := " " + req.URL.Path
	if req.URL.RawQuery != "" {
		fullURL += "?" + req.URL.RawQuery
	}
	if req.URL.Scheme != "" && req.URL.Host != "" {
		fullURL = fmt.Sprintf("%s://%s%s", req.URL.Scheme, req.URL.Host, fullURL)
	} else if req.URL.Host != "" {
		fullURL = fmt.Sprintf("%s%s", req.URL.Host, fullURL)
	}
	content.WriteString(valueStyle.Width(m.rightPanel.Width - 4 - len(req.Method)).Render(fullURL))
	content.WriteString("\n\n")

	if req.URL.RawQuery != "" {
		content.WriteString(fmt.Sprintf("%s\n", sectionStyle.Render("QUERY")))
		kv := req.URL.Query()

		for k, v := range kv {
			content.WriteString(fmt.Sprintf("%s: %s\n",
				keyStyle.Padding(0, 0, 0, 2).Render(k),
				valueStyle.Padding(0, 2, 0, 0).Width(m.rightPanel.Width-4).Render(strings.Join(v, ",")), // may have multiple values with the same key
			))
		}

		content.WriteString("\n\n")
	}

	content.WriteString(fmt.Sprintf("%s\n", sectionStyle.Render("HEADERS")))
	if len(req.Headers) == 0 {
		content.WriteString(valueStyle.Render("  [none]"))
		content.WriteString("\n\n")
	} else {
		for key, values := range req.Headers {
			headerValue := strings.Join(values, ", ")
			content.WriteString(fmt.Sprintf("%s: %s\n",
				keyStyle.Padding(0, 0, 0, 2).Render(key),
				valueStyle.Padding(0, 2, 0, 0).Width(m.rightPanel.Width-4).Render(headerValue),
			))
		}
		content.WriteString("\n\n")
	}

	content.WriteString(fmt.Sprintf("%s ", sectionStyle.Render("BODY")))
	if len(req.Body) == 0 {
		content.WriteString(valueStyle.Render("[empty]"))
		content.WriteString("\n")
	} else {
		content.WriteString(fmt.Sprintf("(%d bytes)\n", len(req.Body)))
		bodyStr := strings.TrimSpace(string(req.Body))
		content.WriteString(fmt.Sprintf("%s\n", valueStyle.Padding(0, 2, 0).Width(m.rightPanel.Width-4).Render(bodyStr)))
	}

	m.rightPanel.SetContent(content.String())
}

func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}
	borderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("15"))
	inUseBorderStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("51"))

	leftPanelContent := m.leftPanel.View()
	rightPanelContent := m.rightPanel.View()

	left := ""
	right := ""

	switch m.selectedPanel {
	case leftPanelSelected:
		left = inUseBorderStyle.
			Width(m.width/3 - 4).
			Height(m.height - 4).
			Render(leftPanelContent)
		right = borderStyle.
			Width(m.width*2/3 - 4).
			Height(m.height - 4).
			Render(rightPanelContent)

	case rightPanelSelected:
		left = borderStyle.
			Width(m.width/3 - 4).
			Height(m.height - 4).
			Render(leftPanelContent)
		right = inUseBorderStyle.
			Width(m.width*2/3 - 4).
			Height(m.height - 4).
			Render(rightPanelContent)
	}

	layout := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	requestCount := fmt.Sprintf("%d", len(m.requests))
	if m.selectedRequest >= 0 && m.selectedRequest < len(m.requests) {
		requestCount = fmt.Sprintf("%d/%d", m.selectedRequest+1, len(m.requests))
	}

	footer := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Render(fmt.Sprintf("↑↓ nav • %s • q quit", requestCount))

	return lipgloss.JoinVertical(lipgloss.Left, layout, footer)
}

func NewTUI(mb *bus.MessageBus) *TUI {
	model := NewModel(mb)

	t := TUI{
		program: tea.NewProgram(model, tea.WithAltScreen()),
		model:   model,
	}

	return &t
}

func (t TUI) Run() error {
	_, err := t.program.Run()
	return err
}
