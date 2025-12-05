package app

import (
	"cths/internal/server"
	"cths/internal/tui"
	"cths/pkg/bus"
	"cths/pkg/plugins"
	"log"
	"net/http"
)

type App struct {
	server   *http.Server
	plugins  []plugins.Plugin
	tui      *tui.TUI
	headless bool
}

func NewApp(port int, whatToServe string, headless bool) *App {
	mb := bus.NewMessageBus()

	pluginList := []plugins.Plugin{
		plugins.NewFlagPosterPlugin(),
	}

	for _, plugin := range pluginList {
		if err := plugin.Init(mb); err != nil {
			log.Printf("Failed to initialize plugin '%s': %v", plugin.Name(), err)
		}
	}

	app := &App{
		server:   server.NewServer(port, whatToServe, mb),
		plugins:  pluginList,
		headless: headless,
	}

	if !headless {
		app.tui = tui.NewTUI(mb)
	}

	return app
}

func (a *App) Run() error {
	for _, plugin := range a.plugins {
		if err := plugin.Start(); err != nil {
			log.Printf("Failed to start plugin '%s': %v", plugin.Name(), err)
		}
	}

	go a.startServer()
	if !a.headless {
		go a.startTUI()
	}

	select {}
}

func (a *App) startServer() {
	log.Printf("starting server on %v", a.server.Addr)
	if err := a.server.ListenAndServe(); err != nil {
		log.Fatalf("error when running server: %v", err)
	}
}

func (a *App) startTUI() {
	if err := a.tui.Run(); err != nil {
		log.Fatalf("error when running tui: %v", err)
	}
}
