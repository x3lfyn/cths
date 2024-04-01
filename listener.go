package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"time"
)

func listener(p *tea.Program) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.FileServer(http.Dir(".")).ServeHTTP(writer, request)

		p.Send(gotRequestMsg{data: HttpRequest{request, time.Now()}})
	})

	err := http.ListenAndServe("0.0.0.0:9999", nil)
	if err != nil {
		panic("fuck!!!")
	}
}
