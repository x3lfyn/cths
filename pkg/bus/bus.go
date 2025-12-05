package bus

import (
	"net/http"
	"net/url"
	"time"
)

type MessageType int

const (
	RequestMessage MessageType = iota
	PluginMessage
	// add whatever message types u want your plugins to handle
)

type Message struct {
	Type      MessageType
	Timestamp time.Time
	Payload   any
}

type RequestMessagePayload struct {
	Method  string
	URL     url.URL
	Headers http.Header
	Body    []byte
}

type PluginMessagePayload struct {
	Content string
}

type MessageBus struct {
	subscribers map[MessageType][]chan Message
}

func NewMessageBus() *MessageBus {
	return &MessageBus{
		subscribers: make(map[MessageType][]chan Message),
	}
}

func (b *MessageBus) Subscribe(messageType MessageType) <-chan Message {
	ch := make(chan Message, 50)
	b.subscribers[messageType] = append(b.subscribers[messageType], ch)
	return ch
}

func (b *MessageBus) Publish(message Message) {
	for _, sub := range b.subscribers[message.Type] {
		select {
		case sub <- message:
		default:
		}
	}
}

func (b *MessageBus) Close() {
	for msgType, subs := range b.subscribers {
		for _, ch := range subs {
			close(ch)
		}
		delete(b.subscribers, msgType)
	}
}
