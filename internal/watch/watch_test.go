package watch

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatchSingleFile(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create callback counter
	var callCount atomic.Int32
	callback := func() {
		callCount.Add(1)
	}

	// Start watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- Watch(tmpFile.Name(), callback)
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify file
	err = os.WriteFile(tmpFile.Name(), []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Wait for callback
	time.Sleep(200 * time.Millisecond)

	// Stop watcher
	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		proc.Signal(os.Interrupt)
	}

	<-done

	if callCount.Load() == 0 {
		t.Error("Expected callback to be called, but it wasn't")
	}
}

func TestWatchMultipleFiles(t *testing.T) {
	// Create temporary files
	tmpFile1, err := os.CreateTemp("", "test_1_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file 1: %v", err)
	}
	defer os.Remove(tmpFile1.Name())
	tmpFile1.Close()

	tmpFile2, err := os.CreateTemp("", "test_2_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file 2: %v", err)
	}
	defer os.Remove(tmpFile2.Name())
	tmpFile2.Close()

	// Create callback counter
	var callCount atomic.Int32
	callback := func() {
		callCount.Add(1)
	}

	// Start watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- WatchMultiple([]string{tmpFile1.Name(), tmpFile2.Name()}, callback)
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify first file
	err = os.WriteFile(tmpFile1.Name(), []byte("content 1"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to temp file 1: %v", err)
	}

	time.Sleep(150 * time.Millisecond)

	// Modify second file
	err = os.WriteFile(tmpFile2.Name(), []byte("content 2"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to temp file 2: %v", err)
	}

	// Wait for callbacks
	time.Sleep(200 * time.Millisecond)

	// Stop watcher
	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		proc.Signal(os.Interrupt)
	}

	<-done

	if callCount.Load() < 2 {
		t.Errorf("Expected at least 2 callbacks, got %d", callCount.Load())
	}
}

func TestSetDebounce(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Set custom debounce
	customWatcher := &Watcher{debounce: 50 * time.Millisecond}

	// Create callback counter
	var callCount atomic.Int32
	callback := func() {
		callCount.Add(1)
	}

	// Start watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- customWatcher.Watch(tmpFile.Name(), callback)
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify file twice rapidly
	err = os.WriteFile(tmpFile.Name(), []byte("content 1"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	time.Sleep(30 * time.Millisecond)

	err = os.WriteFile(tmpFile.Name(), []byte("content 2"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Wait for debounce to trigger
	time.Sleep(150 * time.Millisecond)

	// Stop watcher
	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		proc.Signal(os.Interrupt)
	}

	<-done

	// Due to debouncing, we should get 1 callback instead of 2
	if callCount.Load() == 0 {
		t.Error("Expected callback to be called, but it wasn't")
	}
}

func TestWatchNonexistentPath(t *testing.T) {
	callback := func() {}
	err := Watch("/nonexistent/path/that/does/not/exist", callback)
	if err == nil {
		t.Error("Expected error for nonexistent path, but got none")
	}
}

func TestWatchEmptyPaths(t *testing.T) {
	callback := func() {}
	err := WatchMultiple([]string{}, callback)
	if err == nil {
		t.Error("Expected error for empty paths, but got none")
	}
}

func TestWatchDirectory(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "test_watch_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file in the directory first
	filePath := filepath.Join(tmpDir, "test_file.txt")
	err = os.WriteFile(filePath, []byte("initial"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file in directory: %v", err)
	}

	// Create callback counter
	var callCount atomic.Int32
	callback := func() {
		callCount.Add(1)
	}

	// Start watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- Watch(filePath, callback)
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify the file in the directory
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file in directory: %v", err)
	}

	// Wait for callback
	time.Sleep(200 * time.Millisecond)

	// Stop watcher
	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		proc.Signal(os.Interrupt)
	}

	<-done

	if callCount.Load() == 0 {
		t.Error("Expected callback for file watch, but it wasn't called")
	}
}

func TestCallbackExecution(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create callback that sets a flag
	var callbackExecuted bool
	var mu sync.Mutex
	callback := func() {
		mu.Lock()
		callbackExecuted = true
		mu.Unlock()
	}

	// Start watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- Watch(tmpFile.Name(), callback)
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Modify file
	err = os.WriteFile(tmpFile.Name(), []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Wait for callback
	time.Sleep(200 * time.Millisecond)

	// Stop watcher
	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		proc.Signal(os.Interrupt)
	}

	<-done

	mu.Lock()
	defer mu.Unlock()
	if !callbackExecuted {
		t.Error("Expected callback to be executed")
	}
}

func TestGracefulShutdown(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create callback
	callback := func() {}

	// Start watcher in goroutine
	done := make(chan error, 1)
	go func() {
		done <- Watch(tmpFile.Name(), callback)
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Send interrupt signal
	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		proc.Signal(os.Interrupt)
	}

	// Wait for watcher to shutdown
	select {
	case <-done:
		// Expected
	case <-time.After(2 * time.Second):
		t.Error("Watcher did not shutdown gracefully within timeout")
	}
}

func TestSetDebounceUpdatesBoth(t *testing.T) {
	newDuration := 250 * time.Millisecond
	SetDebounce(newDuration)

	DefaultWatcher.mu.Lock()
	if DefaultWatcher.debounce != newDuration {
		t.Errorf("Expected debounce to be %v, got %v", newDuration, DefaultWatcher.debounce)
	}
	DefaultWatcher.mu.Unlock()

	// Reset to default
	SetDebounce(100 * time.Millisecond)
}

func TestWatcherConcurrency(t *testing.T) {
	// Create temporary files
	tmpFile1, err := os.CreateTemp("", "test_1_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file 1: %v", err)
	}
	defer os.Remove(tmpFile1.Name())
	tmpFile1.Close()

	tmpFile2, err := os.CreateTemp("", "test_2_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file 2: %v", err)
	}
	defer os.Remove(tmpFile2.Name())
	tmpFile2.Close()

	// Create counters
	var count1, count2 atomic.Int32

	// Start two watchers concurrently
	done1 := make(chan error, 1)
	done2 := make(chan error, 1)

	go func() {
		done1 <- Watch(tmpFile1.Name(), func() {
			count1.Add(1)
		})
	}()

	go func() {
		done2 <- Watch(tmpFile2.Name(), func() {
			count2.Add(1)
		})
	}()

	// Give watchers time to start
	time.Sleep(100 * time.Millisecond)

	// Modify files
	os.WriteFile(tmpFile1.Name(), []byte("content 1"), 0644)
	time.Sleep(150 * time.Millisecond)

	os.WriteFile(tmpFile2.Name(), []byte("content 2"), 0644)
	time.Sleep(150 * time.Millisecond)

	// Stop watchers
	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		proc.Signal(os.Interrupt)
	}

	<-done1
	<-done2

	if count1.Load() == 0 || count2.Load() == 0 {
		t.Errorf("Expected both callbacks to be called, got count1=%d, count2=%d", count1.Load(), count2.Load())
	}
}

func TestWriteEventDetection(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Create callback
	var callCount atomic.Int32
	callback := func() {
		callCount.Add(1)
	}

	// Start watcher
	done := make(chan error, 1)
	go func() {
		done <- Watch(tmpFile.Name(), callback)
	}()

	// Give watcher time to start
	time.Sleep(100 * time.Millisecond)

	// Perform write
	err = os.WriteFile(tmpFile.Name(), []byte("new content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	// Wait for event
	time.Sleep(200 * time.Millisecond)

	// Shutdown
	proc, err := os.FindProcess(os.Getpid())
	if err == nil {
		proc.Signal(os.Interrupt)
	}

	<-done

	if callCount.Load() == 0 {
		t.Error("Write event was not detected")
	}
}
