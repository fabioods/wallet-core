package events

import "errors"

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (e *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := e.handlers[eventName]; ok {
		for _, h := range e.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}

	e.handlers[eventName] = append(e.handlers[eventName], handler)
	return nil
}

func (e *EventDispatcher) Dispatch(event EventInterface) error {
	for _, handler := range e.handlers[event.GetName()] {
		err := handler.Handle(event)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := e.handlers[eventName]; ok {
		for i, h := range e.handlers[eventName] {
			if h == handler {
				e.handlers[eventName] = append(e.handlers[eventName][:i], e.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (e *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := e.handlers[eventName]; ok {
		for _, h := range e.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (e *EventDispatcher) Clear() error {
	e.handlers = make(map[string][]EventHandlerInterface)
	return nil
}
