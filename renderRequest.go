package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"net/http"
	"strconv"
	"strings"
)

func renderRequest(req *http.Request) string {
	var headers string
	for key, val := range req.Header {
		headers += key + ": " + strings.Join(val, ", ") + "\n"
	}

	return fmt.Sprintf("%s\n%s %s bytes\n%s %s\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render(req.Method+" "+req.RequestURI),
		lipgloss.NewStyle().Bold(true).Underline(true).Render("Length"),
		strconv.FormatInt(req.ContentLength, 10),
		lipgloss.NewStyle().Bold(true).Underline(true).Render("From"),
		req.RemoteAddr,
		lipgloss.NewStyle().Bold(true).Underline(true).Render("Headers"),
		headers,
	)
}
