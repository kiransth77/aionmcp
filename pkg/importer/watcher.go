package importer

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

// FileWatcher watches specification files for changes and triggers reloads
type FileWatcher struct {
	watcher  *fsnotify.Watcher
	manager  *ImporterManager
	logger   *zap.Logger
	mu       sync.RWMutex
	watching map[string]string      // file path -> source ID
	debounce map[string]*time.Timer // debounce timers for file changes
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewFileWatcher creates a new file watcher
func NewFileWatcher(manager *ImporterManager, logger *zap.Logger) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create file watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	fw := &FileWatcher{
		watcher:  watcher,
		manager:  manager,
		logger:   logger,
		watching: make(map[string]string),
		debounce: make(map[string]*time.Timer),
		ctx:      ctx,
		cancel:   cancel,
	}

	// Start watching in a goroutine
	go fw.watch()

	return fw, nil
}

// WatchSpec starts watching a specification file for changes
func (w *FileWatcher) WatchSpec(source SpecSource) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Only watch local files, not URLs
	if source.Path == "" || filepath.IsAbs(source.Path) == false {
		return fmt.Errorf("can only watch local file paths")
	}

	// Get absolute path
	absPath, err := filepath.Abs(source.Path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Add to watcher
	if err := w.watcher.Add(absPath); err != nil {
		return fmt.Errorf("failed to add file to watcher: %w", err)
	}

	// Track the mapping
	w.watching[absPath] = source.ID

	w.logger.Info("Started watching specification file",
		zap.String("source_id", source.ID),
		zap.String("path", absPath),
		zap.String("type", string(source.Type)))

	return nil
}

// UnwatchSpec stops watching a specification file
func (w *FileWatcher) UnwatchSpec(sourceID string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Find the file path for this source ID
	var pathToRemove string
	for path, id := range w.watching {
		if id == sourceID {
			pathToRemove = path
			break
		}
	}

	if pathToRemove == "" {
		return fmt.Errorf("source ID not found in watched files: %s", sourceID)
	}

	// Remove from watcher
	if err := w.watcher.Remove(pathToRemove); err != nil {
		w.logger.Warn("Failed to remove file from watcher",
			zap.String("path", pathToRemove),
			zap.Error(err))
	}

	// Clean up tracking
	delete(w.watching, pathToRemove)

	// Cancel any pending debounce timer
	if timer, exists := w.debounce[pathToRemove]; exists {
		timer.Stop()
		delete(w.debounce, pathToRemove)
	}

	w.logger.Info("Stopped watching specification file",
		zap.String("source_id", sourceID),
		zap.String("path", pathToRemove))

	return nil
}

// watch runs the file watching loop
func (w *FileWatcher) watch() {
	defer w.watcher.Close()

	for {
		select {
		case <-w.ctx.Done():
			w.logger.Info("File watcher stopped")
			return

		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			w.handleFileEvent(event)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

			w.logger.Error("File watcher error", zap.Error(err))
		}
	}
}

// handleFileEvent processes file system events
func (w *FileWatcher) handleFileEvent(event fsnotify.Event) {
	w.mu.RLock()
	sourceID, exists := w.watching[event.Name]
	w.mu.RUnlock()

	if !exists {
		return // Not watching this file
	}

	// Only handle write and create events
	if event.Op&fsnotify.Write == 0 && event.Op&fsnotify.Create == 0 {
		return
	}

	w.logger.Debug("File change detected",
		zap.String("path", event.Name),
		zap.String("source_id", sourceID),
		zap.String("operation", event.Op.String()))

	// Debounce rapid file changes (common with editors that save frequently)
	w.debounceReload(event.Name, sourceID)
}

// debounceReload debounces rapid file changes to avoid excessive reloads
func (w *FileWatcher) debounceReload(path, sourceID string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Cancel existing timer if any
	if timer, exists := w.debounce[path]; exists {
		timer.Stop()
	}

	// Create new debounce timer
	w.debounce[path] = time.AfterFunc(500*time.Millisecond, func() {
		w.performReload(path, sourceID)

		// Clean up timer
		w.mu.Lock()
		delete(w.debounce, path)
		w.mu.Unlock()
	})
}

// performReload actually performs the specification reload
func (w *FileWatcher) performReload(path, sourceID string) {
	w.logger.Info("Reloading specification due to file change",
		zap.String("source_id", sourceID),
		zap.String("path", path))

	start := time.Now()

	// Reload the specification
	result, err := w.manager.ReloadSpec(w.ctx, sourceID)
	if err != nil {
		w.logger.Error("Failed to reload specification",
			zap.String("source_id", sourceID),
			zap.String("path", path),
			zap.Error(err))
		return
	}

	// Log reload results
	w.logger.Info("Specification reloaded successfully",
		zap.String("source_id", sourceID),
		zap.String("path", path),
		zap.Int("tools_count", len(result.Tools)),
		zap.Int("errors_count", len(result.Errors)),
		zap.Int("warnings_count", len(result.Warnings)),
		zap.Duration("reload_duration", time.Since(start)))

	// Log any errors or warnings
	for _, err := range result.Errors {
		w.logger.Warn("Reload error", zap.Error(err))
	}
	for _, warning := range result.Warnings {
		w.logger.Warn("Reload warning", zap.String("warning", warning))
	}
}

// Stop stops the file watcher
func (w *FileWatcher) Stop() {
	w.cancel()
}

// GetWatchedFiles returns a list of currently watched files
func (w *FileWatcher) GetWatchedFiles() map[string]string {
	w.mu.RLock()
	defer w.mu.RUnlock()

	// Return a copy to avoid race conditions
	result := make(map[string]string)
	for path, sourceID := range w.watching {
		result[path] = sourceID
	}
	return result
}

// IsWatching checks if a source is currently being watched
func (w *FileWatcher) IsWatching(sourceID string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, id := range w.watching {
		if id == sourceID {
			return true
		}
	}
	return false
}
