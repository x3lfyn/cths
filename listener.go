package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func responder(writer http.ResponseWriter, request *http.Request) {
	if state.responder == "fileserver" {
		http.FileServer(http.Dir(".")).ServeHTTP(writer, request)
	} else if state.responder == "file" {
		file, err := os.ReadFile(state.file)
		if err != nil {
			panic(err)
		}

		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Write(file)

	} else {
		fmt.Fprint(writer, state.string)
	}

	b, err := io.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	state.program.Send(gotRequestMsg{data: HttpRequest{request, time.Now(), string(b)}})
}

func listener() {
	if state.responder == "FileServer" {

	}

	http.HandleFunc("/", responder)

	err := http.ListenAndServe(state.address, nil)
	if err != nil {
		panic(err)
	}
}
