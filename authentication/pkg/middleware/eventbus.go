package middleware

import (
	"context"
	"fmt"
	"log"
	"sync"
)

// EventBus defines a simple in-memory event bus for publishing domain events
type EventBus struct {
	handlers map[string][]func(ctx context.Context, event interface{})
	mu       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]func(ctx context.Context, event interface{})),
	}
}

// Publish publishes an event to all registered handlers for that event type
func (eb *EventBus) Publish(ctx context.Context, event interface{}) error {
	eventName := fmt.Sprintf("%T", event) // Gebruik de type naam als event naam
	eb.mu.RLock()
	handlers := eb.handlers[eventName]
	eb.mu.RUnlock()

	if len(handlers) == 0 {
		log.Printf("No handlers registered for event: %s", eventName)
		return nil
	}

	// Publiceer asynchroon om de caller niet te blokkeren
	for _, handler := range handlers {
		go func(h func(ctx context.Context, event interface{})) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in event handler for %s: %v", eventName, r)
				}
			}()
			h(ctx, event)
		}(handler)
	}
	return nil
}

// Subscribe registers an event handler for a specific event type
func (eb *EventBus) Subscribe(eventType interface{}, handler func(ctx context.Context, event interface{})) {
	eventName := fmt.Sprintf("%T", eventType)
	eb.mu.Lock()
	eb.handlers[eventName] = append(eb.handlers[eventName], handler)
	eb.mu.Unlock()
	log.Printf("Subscribed handler for event: %s", eventName)
}
