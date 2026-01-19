package validation

import "fmt"

// EventBus handles event emission and subscription
type EventBus struct {
	subscribers map[string][]EventHandler
}

// EventHandler represents an event handler function
type EventHandler func(event Event)

// Event represents a validation event
type Event struct {
	Name    string
	Data    interface{}
	Emitter string
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]EventHandler),
	}
}

// Subscribe registers a handler for a specific event
func (eb *EventBus) Subscribe(eventName string, handler EventHandler) {
	eb.subscribers[eventName] = append(eb.subscribers[eventName], handler)
}

// Emit triggers all handlers for a specific event
func (eb *EventBus) Emit(event Event) {
	if handlers, ok := eb.subscribers[event.Name]; ok {
		for _, handler := range handlers {
			handler(event)
		}
	}
}

// EmitValidateEvents emits before and after validation events
func EmitValidateEvents(eventBus *EventBus, inputPath string, result *ValidationResult) {
	if eventBus == nil {
		return
	}

	eventBus.Emit(Event{
		Name:    "beforeValidate",
		Data:    map[string]string{"inputPath": inputPath},
		Emitter: "validator",
	})

	eventBus.Emit(Event{
		Name:    "afterValidate",
		Data:    result,
		Emitter: "validator",
	})
}

// String returns string representation of event
func (e Event) String() string {
	return fmt.Sprintf("Event{Name: %s, Emitter: %s}", e.Name, e.Emitter)
}
