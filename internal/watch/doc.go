// Package watch provides file system watching with debouncing.
//
// It monitors single or multiple files using fsnotify and triggers callbacks on write events.
// The package handles graceful signal handling and cleanup.
//
// Example usage:
//
//	watch.Watch("config.yaml", func(event fsnotify.Event) {
//		// Handle file change
//	})
//	watch.WatchMultiple(files, callback)
package watch
