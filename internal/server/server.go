package server

import (
	"bytes"
	"cths/pkg/bus"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type HandlerType int

const (
	FileHandler HandlerType = iota
	DirectoryHandler
)

func determineHandlerType(path string) (HandlerType, error) {
	cleanPath := filepath.Clean(path)

	info, err := os.Stat(cleanPath)
	if err != nil {
		return FileHandler, err
	}

	if info.IsDir() {
		return DirectoryHandler, nil
	}
	return FileHandler, nil
}

func ParseRequestMiddleware(next http.Handler, messageBus *bus.MessageBus) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read body once so both TUI and handler can use it
		var bodyBytes []byte
		if r.Body != nil {
			bodyBytes, _ = io.ReadAll(r.Body)
			r.Body.Close()
			// Create a new reader for the handler
			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		}

		payload := bus.RequestMessagePayload{
			Method:  r.Method,
			URL:     *r.URL,
			Headers: r.Header,
			Body:    bodyBytes,
		}

		message := bus.Message{
			Type:      bus.RequestMessage,
			Timestamp: time.Now(),
			Payload:   payload,
		}

		messageBus.Publish(message)

		next.ServeHTTP(w, r)
	})
}

func NewHandler(whatToServe string, messageBus *bus.MessageBus) http.Handler {

	var lastHandler http.Handler

	handlerType, err := determineHandlerType(whatToServe)
	if err != nil {
		log.Fatalf("error while trying to determine serve type: %v", err)
	}

	switch handlerType {
	case FileHandler:
		lastHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, whatToServe) })

	case DirectoryHandler:
		lastHandler = http.FileServer(http.Dir(whatToServe))
	}

	return ParseRequestMiddleware(lastHandler, messageBus)
}

func NewServer(port int, whatToServe string, messageBus *bus.MessageBus) *http.Server {
	return &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: NewHandler(whatToServe, messageBus),
	}
}
