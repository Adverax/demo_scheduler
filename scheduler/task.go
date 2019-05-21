package scheduler

import (
	"context"
)

type Task interface {
	OnCreate() Event
	OnExecute() Event
	OnSuccess() Event
	OnError() Event
	Execute(ctx context.Context)
	Priority() Priority
	Name() string
	Err() error
}

type task struct {
	name     string
	priority Priority
	action   func() error
	err      error

	onCreate  Event
	onExecute Event
	onSuccess Event
	onError   Event
}

func (task *task) Name() string {
	return task.name
}

func (task *task) Priority() Priority {
	return task.priority
}

func (task *task) Err() error {
	return task.err
}

func (task *task) OnCreate() Event {
	return task.onCreate
}

func (task *task) OnExecute() Event {
	return task.onExecute
}

func (task *task) OnSuccess() Event {
	return task.onSuccess
}

func (task *task) OnError() Event {
	return task.onError
}

func (task *task) Execute(ctx context.Context) {
	task.onExecute.Trigger(ctx, task)
	task.err = task.action()
	if task.err == nil {
		task.onSuccess.Trigger(ctx, task)
	} else {
		task.onError.Trigger(ctx, task)
	}
}

func NewTask(
	name string,
	priority Priority,
	action func() error,
) Task {
	return &task{
		name:      name,
		priority:  priority,
		action:    action,
		onCreate:  newEvent(),
		onExecute: newEvent(),
		onSuccess: newEvent(),
		onError:   newEvent(),
	}
}
