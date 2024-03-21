package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
)

func listener(p *tea.Program) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := fmt.Fprintf(writer, "fuck\n")
		if err != nil {
			return
		}

		p.Send(gotRequestMsg{data: request})
	})

	err := http.ListenAndServe("0.0.0.0:9999", nil)
	if err != nil {
		panic("fuck!!!")
	}
}
