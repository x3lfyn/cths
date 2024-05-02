package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"net/http"
	"os"
	"time"
)

func RequestCatcher(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		handler(writer, request)

		b, err := io.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}

		GlobalState.teaProgram.Send(gotRequestMsg{data: HttpRequest{request, time.Now(), string(b)}})
	}
}

func RunServerAndTui(handler func(http.ResponseWriter, *http.Request)) {
	p := tea.NewProgram(initialModel())
	GlobalState.teaProgram = p

	go func() {
		http.HandleFunc("/", RequestCatcher(handler))

		err := http.ListenAndServe(GlobalState.listenAddress, nil)
		if err != nil {
			panic(err)
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
