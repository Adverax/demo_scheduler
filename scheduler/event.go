package scheduler

import (
	"context"
	"sync"
)

// Подписчик на событие
type Subscriber func(ctx context.Context, task Task)

type subscribers []Subscriber

// Абстрактное событие
type Event interface {
	// Регистрация нового подписчика
	Append(subscriber Subscriber)
	// Оповещение всех подписчиков о событии
	Trigger(ctx context.Context, task Task)
}

// Реализация абстрактного события
type event struct {
	sync.Mutex
	subscribers subscribers
}

func (e *event) Append(action Subscriber) {
	if action == nil {
		return
	}

	e.Lock()
	defer e.Unlock()

	e.subscribers = append(e.subscribers, action)
}

func (e *event) Trigger(ctx context.Context, task Task) {
	e.Lock()
	defer e.Unlock()

	for _, action := range e.subscribers {
		action(ctx, task)
	}
}

func newEvent() Event {
	return &event{
		subscribers: make(subscribers, 0, 4),
	}
}
