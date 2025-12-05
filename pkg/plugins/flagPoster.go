package plugins

import (
	"bytes"
	"cths/pkg/bus"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"
)

type FlagPosterPlugin struct {
	name           string
	bus            *bus.MessageBus
	serverURL      string
	serverPassword string
	flagRegex      *regexp.Regexp
	client         *http.Client
}

type FlagData struct {
	Flag   string `json:"flag"`
	Sploit string `json:"sploit"`
	Team   string `json:"team"`
}

func NewFlagPosterPlugin() *FlagPosterPlugin {
	flagRegexPattern := os.Getenv("FLAG_POSTER_REGEX")
	if flagRegexPattern == "" {
		flagRegexPattern = `[A-Z0-9]{31}=` // default fallback
	}

	flagRegex := regexp.MustCompile(flagRegexPattern)
	serverURL := os.Getenv("FLAG_POSTER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080" // default fallback
	}

	serverPassword := os.Getenv("FLAG_POSTER_PASSWORD")
	if serverPassword == "" {
		serverPassword = "idk"
	}

	return &FlagPosterPlugin{
		name:           "flag-poster",
		flagRegex:      flagRegex,
		serverURL:      serverURL,
		serverPassword: serverPassword,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *FlagPosterPlugin) Name() string {
	return p.name
}

func (p *FlagPosterPlugin) Init(b *bus.MessageBus) error {
	p.bus = b
	return nil
}

func (p *FlagPosterPlugin) Start() error {
	ch := p.bus.Subscribe(bus.RequestMessage)

	go func() {
		for msg := range ch {
			if msg.Type == bus.RequestMessage {
				p.processRequest(msg)
			}
		}
	}()

	log.Printf("Plugin '%s' started (server URL: %s, regex: %s)", p.name, p.serverURL, p.flagRegex.String())
	return nil
}

func (p *FlagPosterPlugin) Stop() error {
	log.Printf("Plugin '%s' stopped", p.name)
	return nil
}

func (p *FlagPosterPlugin) processRequest(msg bus.Message) {
	payload, ok := msg.Payload.(bus.RequestMessagePayload)
	if !ok {
		return
	}

	// Search for flags in all request data
	flags := p.findFlags(payload)

	if len(flags) > 0 {
		log.Printf("Found %d flag(s), posting to server...", len(flags))
		if err := p.postFlags(flags); err != nil {
			log.Printf("Error posting flags: %v", err)
		}
	}
}

func (p *FlagPosterPlugin) findFlags(payload bus.RequestMessagePayload) []string {
	flags := make(map[string]bool) // Use map to avoid duplicates

	urlStr := payload.URL.String()
	p.findFlagsInText(urlStr, flags)

	for _, values := range payload.Headers {
		for _, value := range values {
			p.findFlagsInText(value, flags)
		}
	}

	if len(payload.Body) > 0 {
		p.findFlagsInText(string(payload.Body), flags)
	}

	result := make([]string, 0, len(flags))
	for flag := range flags {
		result = append(result, flag)
	}

	return result
}

func (p *FlagPosterPlugin) findFlagsInText(text string, flags map[string]bool) {
	matches := p.flagRegex.FindAllString(text, -1)
	for _, match := range matches {
		flags[match] = true
	}
}

func (p *FlagPosterPlugin) postFlags(flags []string) error {
	data := make([]FlagData, 0, len(flags))
	for _, flag := range flags {
		data = append(data, FlagData{
			Flag:   flag,
			Sploit: "flag_poster",
			Team:   "*",
		})
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal flags: %w", err)
	}

	apiURL, err := url.JoinPath(p.serverURL, "/api/post_flags")
	if err != nil {
		return fmt.Errorf("failed to construct API URL: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", p.serverPassword)

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var prikol []byte
	resp.Body.Read(prikol)
	log.Printf(string(prikol))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned error status: %d", resp.StatusCode)
	}

	log.Printf("Successfully posted %d flag(s) to %s", len(flags), apiURL)
	return nil
}
