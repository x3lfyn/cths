package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func wordwrap(msg string, width int) string {
	offset := width
	for len(msg) > offset {
		msg = msg[:offset] + "\n" + msg[offset:]
		offset += width + 1
	}
	return msg
}

func renderRequest(req HttpRequest, viewWidth int) string {
	var headers []string
	for key, val := range req.req.Header {
		headers = append(headers, wordwrap(key+": "+strings.Join(val, ", "), viewWidth))
	}

	sort.Strings(headers)

	fieldNameStyle := lipgloss.NewStyle().Bold(true).Underline(true)

	res := fmt.Sprintf("%s\n%s %s bytes\n%s %s\n%s %s\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render(req.req.Method+" "+req.req.RequestURI+" "+req.req.Proto),
		fieldNameStyle.Render("Length"),
		strconv.FormatInt(req.req.ContentLength, 10),
		fieldNameStyle.Render("From"),
		req.req.RemoteAddr,
		fieldNameStyle.Render("Host"),
		wordwrap(req.req.Host, viewWidth),
		fieldNameStyle.Render("Headers"),
		strings.Join(headers, "\n"),
	)

	res += fmt.Sprintf("\n\n%s", wordwrap(req.body, viewWidth))

	return res
}
