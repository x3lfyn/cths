package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"sort"
	"strconv"
	"strings"
)

func renderRequest(req HttpRequest) string {
	var headers []string
	for key, val := range req.req.Header {
		headers = append(headers, key+": "+strings.Join(val, ", "))
	}

	sort.Strings(headers)

	return fmt.Sprintf("%s\n%s %s bytes\n%s %s\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render(req.req.Method+" "+req.req.RequestURI),
		lipgloss.NewStyle().Bold(true).Underline(true).Render("Length"),
		strconv.FormatInt(req.req.ContentLength, 10),
		lipgloss.NewStyle().Bold(true).Underline(true).Render("From"),
		req.req.RemoteAddr,
		lipgloss.NewStyle().Bold(true).Underline(true).Render("Headers"),
		strings.Join(headers, "\n"),
	)
}
