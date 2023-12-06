package eventbus

import "github.com/werbenhu/eventbus"

type EventBus interface {
	ubscribe(topic string, handler any) error
	Publish(topic string, payload any) error
	Unsubscribe(topic string, handler any) error
}

type eventBus struct {
	e *eventbus.EventBus
}

func NewEventBus() *eventBus {
	return &eventBus{
		e: eventbus.New(),
	}
}

func Subscribe(topic string, handler any) error {
	return eventbus.Subscribe(topic, handler)
}

func Publish(topic string, payload any) error {
	return eventbus.Publish(topic, payload)
}

func Unsubscribe(topic string, handler any) error {
	return eventbus.Unsubscribe(topic, handler)
}
