package main

import (
	"flag"

	"cths/internal/app"
)

func main() {

	port := flag.Int("port", 6969, "port to listen on")
	whatToServe := flag.String("serve", ".", "file or directory to serve")
	headless := flag.Bool("headless", false, "run without TUI")
	flag.Parse()

	a := app.NewApp(*port, *whatToServe, *headless)
	a.Run()
}
