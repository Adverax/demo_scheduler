package scheduler

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTask_Execute(t *testing.T) {
	type Match struct {
		execute int
		error   int
		success int
		primary int
		err     error
	}

	type Test struct {
		Match
		OnExecute Subscriber
		OnError   Subscriber
		OnSuccess Subscriber
	}

	var match Match
	var errTest = errors.New("test")

	tests := map[string]Test{
		"Success": {
			Match: Match{
				execute: 1,
				success: 1,
				primary: 1,
			},
			OnExecute: func(ctx context.Context, task Task) error {
				match.execute++
				return nil
			},
			OnSuccess: func(ctx context.Context, task Task) error {
				match.success++
				return nil
			},
			OnError: func(ctx context.Context, task Task) error {
				match.error++
				return nil
			},
		},
		"Failure": {
			Match: Match{
				execute: 1,
				error:   1,
				primary: 1,
				err:     errTest,
			},
			OnExecute: func(ctx context.Context, task Task) error {
				match.execute++
				return nil
			},
			OnSuccess: func(ctx context.Context, task Task) error {
				match.success++
				return nil
			},
			OnError: func(ctx context.Context, task Task) error {
				match.error++
				return nil
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			match.err = nil
			match.execute = 0
			match.error = 0
			match.success = 0
			match.primary = 0

			task := NewTask(
				"MyTask",
				PriorityHigh,
				func() error {
					match.primary++
					return test.Match.err
				},
			)

			task.OnExecute().Append(test.OnExecute)
			task.OnSuccess().Append(test.OnSuccess)
			task.OnError().Append(test.OnError)

			task.Execute(context.Background())
			assert.Equal(t, test.Match.primary, match.primary)
			assert.Equal(t, test.Match.error, match.error)
			assert.Equal(t, test.Match.success, match.success)
			assert.Equal(t, test.Match.execute, match.execute)
			assert.Equal(t, test.Match.err != nil, match.error != 0)
		})
	}
}
