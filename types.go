package main

import (
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type gotRequestMsg struct{ data HttpRequest }

type State struct {
	handler       func(http.ResponseWriter, *http.Request)
	teaProgram    *tea.Program
	listenAddress string
}

type HttpRequest struct {
	req  *http.Request
	time time.Time
	body string
}

type model struct {
	spinner              spinner.Model
	requests             []HttpRequest
	selectedRequestIndex int
	err                  error
	width                int
}
