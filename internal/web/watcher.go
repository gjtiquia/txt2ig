package web

import (
	"log"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type FileWatcher struct {
	watcher  *fsnotify.Watcher
	debounce map[string]*time.Timer
	delay    time.Duration
	onChange func(string)
	mu       sync.Mutex
}

func NewFileWatcher() (*FileWatcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &FileWatcher{
		watcher:  fw,
		debounce: make(map[string]*time.Timer),
		delay:    300 * time.Millisecond,
	}, nil
}

func (w *FileWatcher) Watch(path string, onChange func(string)) error {
	w.onChange = onChange
	return w.watcher.Add(path)
}

func (w *FileWatcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					w.handleEvent(event.Name)
				}
			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("file watcher error: %v", err)
			}
		}
	}()
}

func (w *FileWatcher) handleEvent(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if timer, exists := w.debounce[path]; exists {
		timer.Stop()
	}

	w.debounce[path] = time.AfterFunc(w.delay, func() {
		if w.onChange != nil {
			w.onChange(path)
		}
		w.mu.Lock()
		delete(w.debounce, path)
		w.mu.Unlock()
	})
}

func (w *FileWatcher) Close() error {
	return w.watcher.Close()
}
