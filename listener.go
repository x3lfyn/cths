package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
	"time"
	"io"
)

func listener(p *tea.Program, listenAddr string) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.FileServer(http.Dir(".")).ServeHTTP(writer, request)

		b, err := io.ReadAll(request.Body);
		if err != nil {
			panic(err);
		}

		p.Send(gotRequestMsg{data: HttpRequest{request, time.Now(), string(b)}})
	})

	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		panic(err)
	}
}
