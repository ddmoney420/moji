package watch

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	debounce time.Duration
	mu       sync.Mutex
}

var DefaultWatcher = &Watcher{
	debounce: 100 * time.Millisecond,
}

// SetDebounce sets the debounce duration for the default watcher
func SetDebounce(duration time.Duration) {
	DefaultWatcher.mu.Lock()
	defer DefaultWatcher.mu.Unlock()
	DefaultWatcher.debounce = duration
}

// Watch watches a single file and calls callback on Write events
func Watch(path string, callback func()) error {
	return DefaultWatcher.Watch(path, callback)
}

// WatchMultiple watches multiple files and calls callback on Write events
func WatchMultiple(paths []string, callback func()) error {
	return DefaultWatcher.WatchMultiple(paths, callback)
}

// Watch (Watcher method) watches a single file and calls callback on Write events
func (w *Watcher) Watch(path string, callback func()) error {
	return w.WatchMultiple([]string{path}, callback)
}

// WatchMultiple (Watcher method) watches multiple files and calls callback on Write events
func (w *Watcher) WatchMultiple(paths []string, callback func()) error {
	if len(paths) == 0 {
		return fmt.Errorf("no paths provided to watch")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	// Add paths to watcher
	for _, path := range paths {
		// Check if path exists
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat %s: %w", path, err)
		}

		if info.IsDir() {
			if err := watcher.Add(path); err != nil {
				return fmt.Errorf("failed to watch directory %s: %w", path, err)
			}
		} else {
			// Watch the directory containing the file
			dir := path[:len(path)-len(info.Name())]
			if dir == "" {
				dir = "."
			}
			if err := watcher.Add(dir); err != nil {
				return fmt.Errorf("failed to watch directory %s: %w", dir, err)
			}
		}
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Debounce timer
	var debounceTimer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}

			// Only trigger on Write events
			if event.Op&fsnotify.Write != 0 {
				// Cancel previous timer if it exists
				if debounceTimer != nil {
					debounceTimer.Stop()
				}

				// Set new debounce timer
				w.mu.Lock()
				duration := w.debounce
				w.mu.Unlock()

				debounceTimer = time.AfterFunc(duration, func() {
					callback()
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			return fmt.Errorf("watcher error: %w", err)

		case <-sigChan:
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return nil
		}
	}
}
