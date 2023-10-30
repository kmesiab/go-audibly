package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// FileHandler defines the type of the function to be called when a new file is detected
// by the fsnotify Watcher in the Watch function. It takes the full path of the newly
// created file as its argument.
//
// Functions matching this type can be used to further process the new file, such as uploading it to S3
// or initiating a transcription job, as demonstrated in the handleEvents and Watch functions.
type FileHandler func(string)

// handleEvents listens for file system events from the fsnotify Watcher and processes them.
// It calls onNewAudioFile whenever a new audio file is created in the watched directory.
// The function runs until the provided context is canceled, an error occurs, or the watcher is closed.
//
// Parameters:
// - ctx: The context used to control the termination of the event handling.
// - watcher: The fsnotify Watcher monitoring the directory.
// - onNewAudioFile: The function to call when a new audio file is detected.
// - done: A channel that will be closed when the function returns to signal completion.
func handleEvents(ctx context.Context, watcher *fsnotify.Watcher, onNewAudioFile FileHandler, done chan struct{}) {
	defer close(done)

	for {
		select {
		case <-ctx.Done():

			return
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				if IsAudioFile(event.Name, AllowedAudioExtensions) {
					fmt.Printf("🎵 New audio file detected: %s 🎵\n", event.Name)
					onNewAudioFile(event.Name)
				} else {
					fmt.Printf("⚠️ Ignoring new file: %s ⚠️\n", event.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			fmt.Println("Error:", err)
		}
	}
}

// Watch initializes a new file watcher to monitor the directory at the specified path.
// It triggers the onNewAudioFile function for each new audio file added to the directory.
// The function runs until the provided context is canceled or an error occurs.
//
// Parameters:
// - ctx: The context used to control the termination of the watcher.
// - path: The directory path to watch for new audio files.
// - onNewAudioFile: The function to call when a new audio file is detected.
//
// Returns:
// - An error if the watcher encounters an issue, otherwise nil.
func Watch(ctx context.Context, path string, onNewAudioFile FileHandler) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}

	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			PrepareLogMessagef("Failed to close watcher: %s", err.Error()).Error()
		}
	}(watcher)

	done := make(chan struct{})

	go handleEvents(ctx, watcher, onNewAudioFile, done)

	if err := watcher.Add(path); err != nil {
		return fmt.Errorf("error watching folder: %w", err)
	}

	<-done

	return nil
}

// ProcessExistingFiles scans the directory at the given path for existing audio files that match
// the allowed file extensions. For each matching file, it triggers the provided onNewAudioFile handler.
//
// Parameters:
// - path: The directory path to scan for existing audio files.
// - allowedExtensions: A slice of string containing the allowed audio file extensions.
// - onNewAudioFile: The function to call when a new audio file is found.
func ProcessExistingFiles(path string, allowedExtensions []string, onNewAudioFile FileHandler) {
	files, err := os.ReadDir(path)
	if err != nil {
		PrepareLogMessagef("Error reading directory: %s", err.Error()).Error()

		return
	}

	for _, file := range files {
		if IsAudioFile(file.Name(), allowedExtensions) {
			PrepareLogMessagef("Existing audio file detected: %s", file.Name()).Info()
			fullFileName := fmt.Sprintf("%s/%s", path, file.Name())
			onNewAudioFile(fullFileName)
		}
	}
}

func IsAudioFile(filename string, allowedExtensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, audioExt := range allowedExtensions {
		if ext == audioExt {
			return true
		}
	}

	return false
}
