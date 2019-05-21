package scheduler

import (
	"context"
	"sort"
	"sync"
)

const (
	PriorityLowest Priority = iota + 1
	PriorityLow
	PriorityNormal
	PriorityHigh
	PriorityHighest
)

// Приоритет - это приоритет задачи
type Priority int

// tasks - сортированный список задач (по приоритету)
type tasks []Task

func (tasks tasks) Len() int {
	return len(tasks)
}

func (tasks tasks) Swap(i, j int) {
	tasks[j], tasks[i] = tasks[i], tasks[j]
}

func (tasks tasks) Less(i, j int) bool {
	return tasks[i].Priority() < tasks[j].Priority()
}

// Поиск задачи в списку задач
func (tasks tasks) indexOf(task Task) int {
	for i, t := range tasks {
		if t == task {
			return i
		}
	}

	return -1
}

// Менеджер - интерфейс менеджера задач
type Manager interface {
	// Добавление новой задачи в список задач
	Append(ctx context.Context, task Task)
	// Удаление задачи из списка задач
	Remove(ctx context.Context, task Task)
	// Запуск менеджера на выполнение
	Run(ctx context.Context)
	// Останов мереджера.
	// Останов также может быть выполнен через контекст, переданный в Run
	Shutdown()
}

// Engine - конкретная реализация менеджера задач
type engine struct {
	sync.Mutex
	tasks tasks
	work  chan struct{}
	stop  chan struct{}
	done  chan struct{}
}

func (engine *engine) Append(ctx context.Context, task Task) {
	engine.Lock()
	defer engine.Unlock()

	engine.tasks = append(engine.tasks, task)
	task.OnCreate().Trigger(ctx, task)
	sort.Sort(engine.tasks)
	engine.work <- struct{}{}
}

func (engine *engine) Remove(ctx context.Context, task Task) {
	engine.Lock()
	defer engine.Unlock()

	index := engine.tasks.indexOf(task)
	if index == -1 {
		return
	}

	engine.tasks = append(
		engine.tasks[:index],
		engine.tasks[index+1:]...,
	)
}

func (engine *engine) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(engine.done)
			return
		case <-engine.stop:
			close(engine.done)
			return
		case <-engine.work:
			task := engine.fetch()
			if task != nil {
				task.Execute(ctx)
			}
		}
	}
}

func (engine *engine) Shutdown() {
	close(engine.stop)
	<-engine.done
}

// Извлечение задачи для выполнения
func (engine *engine) fetch() Task {
	engine.Lock()
	defer engine.Unlock()

	l := len(engine.tasks)
	if l == 0 {
		return nil
	}
	task := engine.tasks[l-1]
	engine.tasks = engine.tasks[:l-1]

	return task
}

// Создание экземпляра менеджера задач
func New() Manager {
	return &engine{
		tasks: make(tasks, 0, 16),
		work:  make(chan struct{}, 1024),
		done:  make(chan struct{}),
		stop:  make(chan struct{}),
	}
}
