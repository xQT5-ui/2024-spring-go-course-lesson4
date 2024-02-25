package watcher

import (
	"context"
	"errors"
	"time"
)

type EventType string

const (
	EventTypeFileCreate  EventType = "file_created"
	EventTypeFileRemoved EventType = "file_removed"
)

var ErrDirNotExist = errors.New("dir does not exist")

// Event - событие, формируемое при изменении в файловой системе.
type Event struct {
	// Type - тип события.
	Type EventType
	// Path - путь к объекту изменения.
	Path string
}

type Watcher struct {
	// Events - канал событий
	Events chan Event
	// refreshInterval - интервал обновления списка файлов
	refreshInterval time.Duration

	// можете добавить свои поля
}

func NewDirWatcher(refreshInterval time.Duration) *Watcher {
	return &Watcher{
		refreshInterval: refreshInterval,
		Events:          make(chan Event),
	}
}

func (w *Watcher) WatchDir(ctx context.Context, path string) error {
	// TODO реализовать функцию
	return nil
}

func (w *Watcher) Close() {
	close(w.Events)
}
