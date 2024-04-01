package main

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"io"
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

	fieldNameStyle := lipgloss.NewStyle().Bold(true).Underline(true)

	res := fmt.Sprintf("%s\n%s %s bytes\n%s %s\n%s %s\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).Render(req.req.Method+" "+req.req.RequestURI+" "+req.req.Proto),
		fieldNameStyle.Render("Length"),
		strconv.FormatInt(req.req.ContentLength, 10),
		fieldNameStyle.Render("From"),
		req.req.RemoteAddr,
		fieldNameStyle.Render("Host"),
		req.req.Host,
		fieldNameStyle.Render("Headers"),
		strings.Join(headers, "\n"),
	)

	b, err := io.ReadAll(req.req.Body)
	if err != nil {
		panic(err)
	}
	res += fmt.Sprintf("\n\n%s", string(b))

	return res
}
