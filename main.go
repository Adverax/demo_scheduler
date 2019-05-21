package main

import (
	"context"
	"scheduler/scheduler"
)

func main() {
	m := scheduler.New()
	// Настроить шедулер

	// ...
	// m.Append(scheduler.NewTask(...))

	// Запустить его на выполнение
	go m.Run(context.Background())

	// Потокобезопасен.
	// Можно запускать в отельной горутине и добавлять задачи или обработчики событий.
	m.Shutdown()
}
