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
}

type task struct {
	name     string
	priority Priority
	action   func() error

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
	err := task.execute(ctx)
	if err == nil {
		err = task.onSuccess.Trigger(ctx, task)
	}
	if err != nil {
		_ = task.onError.Trigger(ctx, task)
	}
}

func (task *task) execute(ctx context.Context) error {
	err := task.onExecute.Trigger(ctx, task)
	if err != nil {
		return err
	}
	return task.action()
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
