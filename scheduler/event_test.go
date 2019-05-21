package scheduler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEvent_AppendAndTrigger(t *testing.T) {
	var count int
	e := newEvent()
	e.Append(func(ctx context.Context, task Task) error {
		count++
		return nil
	})

	err := e.Trigger(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}
