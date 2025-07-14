// Package watcher for tracking files in direction and sub-directions
package watcher

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type EventType string

const (
	EventTypeFileCreate  EventType = "file_created"
	EventTypeFileRemoved EventType = "file_removed"
)

var (
	ErrDirNotExist    = errors.New("dir does not exist")
	ErrWalkDir        = errors.New("error from walking dirs")
	ErrCreateSnapshot = errors.New("problem with creating snapshot")
)

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

	// currentSnapshot - текущий список файлов, где ключа досточно для проверки наличия или отсутствия файла
	currentSnapshot map[string]struct{}
}

func NewDirWatcher(refreshInterval time.Duration) *Watcher {
	return &Watcher{
		refreshInterval: refreshInterval,
		Events:          make(chan Event),
		currentSnapshot: make(map[string]struct{}),
	}
}

// createSnapshot создает копию текущего состояния файлов в словарь-мапу
func (w *Watcher) createSnapshot(rootPatch string) (map[string]struct{}, error) {
	snapshot := make(map[string]struct{})

	err := filepath.WalkDir(rootPatch, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("%w: %w", ErrWalkDir, err)
		}

		// Add only files
		if !d.IsDir() {
			snapshot[path] = struct{}{}
		}

		return nil
	})

	return snapshot, err
}

// compareSnapshots сравнивает два снимка: "было" (текущий) и "стало" (новый)
func (w *Watcher) compareSnapshots(oldSnapshot, newSnapshot map[string]struct{}) {
	// Ищем новые файлы (есть в новом, нет в старом)
	for filePath := range newSnapshot {
		if _, exists := oldSnapshot[filePath]; !exists {
			w.Events <- Event{
				Type: EventTypeFileCreate,
				Path: filePath,
			}
		}
	}

	// Ищем удалённые файлы (есть в старом, нет в новом)
	for filePath := range oldSnapshot {
		if _, exists := newSnapshot[filePath]; !exists {
			w.Events <- Event{
				Type: EventTypeFileRemoved,
				Path: filePath,
			}
		}
	}
}

// WatchDir является главной функцией для работы наблюдателя
func (w *Watcher) WatchDir(ctx context.Context, path string) error {
	// Проверка наличия папки
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return ErrDirNotExist
	}

	// Создаём первый снимок (при запуске программы)
	initSnapshot, err := w.createSnapshot(path)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCreateSnapshot, err)
	}
	w.currentSnapshot = initSnapshot

	// Создаём тикер для периодических проверок
	ticker := time.NewTicker(w.refreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Контекст отменён
			return ctx.Err()
		case <-ticker.C:
			// Время для новой проверки (сработал тик)
			// создаём снимок "стало"
			newSnapshot, err := w.createSnapshot(path)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrCreateSnapshot, err)
			}

			// сравниваем снимки
			w.compareSnapshots(w.currentSnapshot, newSnapshot)

			// обновление снимка "было"
			w.currentSnapshot = newSnapshot
		}
	}
}

func (w *Watcher) Close() {
	close(w.Events)
}
