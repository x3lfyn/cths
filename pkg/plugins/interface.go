// pkg/plugins/interface.go
package plugins

import "cths/pkg/bus"

type Plugin interface {
	Name() string
	Init(bus *bus.MessageBus) error
	Start() error
	Stop() error
}
