package scheduler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvent_AppendAndTrigger(t *testing.T) {
	var count int
	e := newEvent()
	e.Append(func(ctx context.Context, task Task) {
		count++
	})

	e.Trigger(context.Background(), nil)
	assert.Equal(t, 1, count)
}
