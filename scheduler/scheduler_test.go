package scheduler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEngine_AppendRemove(t *testing.T) {
	m := &engine{
		tasks: make(tasks, 0, 16),
		work:  make(chan struct{}, 1024),
	}

	task := NewTask("1", PriorityLow, nil)

	m.Append(task)
	require.Equal(t, 1, len(m.tasks))

	m.Remove(task)
	require.Equal(t, 0, len(m.tasks))
}

func TestEngine_Fetch(t *testing.T) {
	type Test struct {
		tasks []Task
		res   string
	}

	tests := map[string]Test{
		"Single": {
			tasks: []Task{
				NewTask("1", PriorityLow, nil),
			},
			res: "1",
		},
		"Multiple": {
			tasks: []Task{
				NewTask("1", PriorityLow, nil),
				NewTask("2", PriorityHighest, nil),
				NewTask("3", PriorityHigh, nil),
			},
			res: "231",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &engine{
				tasks: make(tasks, 0, 16),
				work:  make(chan struct{}, 1024),
			}

			for _, tt := range test.tasks {
				m.Append(tt)
			}

			var res string
			task := m.fetch()
			for task != nil {
				res += task.Name()
				task = m.fetch()
			}

			assert.Equal(t, test.res, res)
		})
	}
}

func TestEngine_Shutdown(t *testing.T) {
	m := New()
	go m.Run(context.Background())
	m.Shutdown()
}
