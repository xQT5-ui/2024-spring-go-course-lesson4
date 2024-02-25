package watcher

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const TickInterval = 250 * time.Millisecond

func newDir() (string, error) {
	return os.MkdirTemp("", "watcher-test*")
}

func mkFile(path string) (string, error) {
	f, err := os.CreateTemp(path, "test-*.txt")
	return f.Name(), err
}

func TestNotifier_Notify(t *testing.T) {
	t.Run("err, path not exist", func(t *testing.T) {
		w := NewDirWatcher(TickInterval)
		defer w.Close()
		err := w.WatchDir(context.Background(), "not_exist_path")
		assert.Error(t, err)
		assert.Equal(t, ErrDirNotExist, err)
	})

	t.Run("err, ctx cancelled", func(t *testing.T) {
		d, err := newDir()
		assert.NoError(t, err)

		w := NewDirWatcher(TickInterval)
		defer w.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		err = w.WatchDir(ctx, d)
		assert.Error(t, err)
		assert.Equal(t, "context deadline exceeded", err.Error())
	})

	t.Run("ok, no files on start", func(t *testing.T) {
		d, err := newDir()
		assert.NoError(t, err)
		defer os.RemoveAll(d)

		w := NewDirWatcher(TickInterval)
		defer w.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = w.WatchDir(ctx, d)
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					assert.NoError(t, err)
				}
			}
		}()

		var result []Event
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case e := <-w.Events:
					result = append(result, e)
				case <-ctx.Done():
					return
				}
			}
		}()

		wg.Wait()

		assert.Equal(t, 0, len(result))
	})

	t.Run("ok, has files on start", func(t *testing.T) {
		d, err := newDir()
		assert.NoError(t, err)
		defer os.RemoveAll(d)

		_, err = mkFile(d)
		assert.NoError(t, err)
		_, err = mkFile(d)
		assert.NoError(t, err)
		_, err = mkFile(d)
		assert.NoError(t, err)

		w := NewDirWatcher(TickInterval)
		defer w.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = w.WatchDir(ctx, d)
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					assert.NoError(t, err)
				}
			}
		}()

		var result []Event
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case e := <-w.Events:
					result = append(result, e)
				case <-ctx.Done():
					return
				}
			}
		}()

		wg.Wait()

		assert.Equal(t, 0, len(result))
	})

	t.Run("ok, has files on start, added files, removed files", func(t *testing.T) {
		d, err := newDir()
		assert.NoError(t, err)
		defer os.RemoveAll(d)

		_, err = mkFile(d)
		assert.NoError(t, err)
		f2, err := mkFile(d)
		assert.NoError(t, err)
		_, err = mkFile(d)
		assert.NoError(t, err)

		w := NewDirWatcher(TickInterval)
		defer w.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := w.WatchDir(ctx, d)
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					assert.NoError(t, err)
				}
			}
		}()

		var result []Event
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case e := <-w.Events:
					result = append(result, e)
				case <-ctx.Done():
					return
				}
			}
		}()

		var f4 string
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Second)
			var err error
			f4, err = mkFile(d)
			assert.NoError(t, err)
			time.Sleep(1500 * time.Millisecond)
			assert.NoError(t, os.Remove(f2))
		}()

		wg.Wait()

		assert.Equal(t, 2, len(result))

		assert.True(t, slices.ContainsFunc(result, func(e Event) bool {
			return e.Type == EventTypeFileCreate && e.Path == f4
		}))
		assert.True(t, slices.ContainsFunc(result, func(e Event) bool {
			return e.Type == EventTypeFileRemoved && e.Path == f2
		}))
	})

	t.Run("ok, subdirs watching", func(t *testing.T) {
		d, err := newDir()
		assert.NoError(t, err)
		defer os.RemoveAll(d)

		w := NewDirWatcher(TickInterval)
		defer w.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := w.WatchDir(ctx, d)
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					assert.NoError(t, err)
				}
			}
		}()

		var result []Event
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case e := <-w.Events:
					result = append(result, e)
				case <-ctx.Done():
					return
				}
			}
		}()

		var f1, f2 string
		wg.Add(1)
		go func() {
			defer wg.Done()
			var err error
			time.Sleep(500 * time.Millisecond)
			f1, err = mkFile(d)
			assert.NoError(t, err)

			d2 := filepath.Join(d, "subdir")
			assert.NoError(t, os.Mkdir(d2, 0o755))

			f2, err = mkFile(d2)
			assert.NoError(t, err)

			time.Sleep(time.Second)
			assert.NoError(t, os.Remove(f2))
			time.Sleep(time.Second)
			assert.NoError(t, os.Remove(f1))
		}()

		wg.Wait()

		assert.Equal(t, 4, len(result))

		assert.True(t, slices.ContainsFunc(result, func(e Event) bool {
			return e.Type == EventTypeFileCreate && e.Path == f1
		}))
		assert.True(t, slices.ContainsFunc(result, func(e Event) bool {
			return e.Type == EventTypeFileCreate && e.Path == f2
		}))

		assert.True(t, slices.ContainsFunc(result, func(e Event) bool {
			return e.Type == EventTypeFileRemoved && e.Path == f1
		}))
		assert.True(t, slices.ContainsFunc(result, func(e Event) bool {
			return e.Type == EventTypeFileRemoved && e.Path == f2
		}))
	})
}
