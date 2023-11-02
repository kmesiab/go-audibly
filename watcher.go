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
// or initiating a transcription job, as demonstrated in the HandleFileEvents and Watch functions.
type FileHandler func(string)

// HandleFileEvents listens for file system events from the fsnotify Watcher and processes them.
// It calls onNewAudioFile whenever a new audio file is created in the watched directory.
// The function runs until the provided context is canceled, an error occurs, or the watcher is closed.
//
// Parameters:
// - ctx: The context to be used for cancelling the function.
// - watcher: The fsnotify Watcher monitoring the directory.
// - onNewAudioFile: The function to call when a new audio file is detected.
// - done: A channel that will be closed when the function returns to signal completion.
func HandleFileEvents(
	ctx context.Context,
	watcher *fsnotify.Watcher,
	allowedAudioExtensions *[]string,
	onNewAudioFile FileHandler,
	done chan struct{},
) {
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
				if IsAudioFile(event.Name, allowedAudioExtensions) {
					PrepareLogMessagef("🎵 New audio file detected: %s", event.Name).
						Info()

					onNewAudioFile(event.Name)

				} else {
					PrepareLogMessagef("⚠️ Ignoring new file: %s ⚠️\n", event.Name).
						Info()
				}
			}
		case err, ok := <-watcher.Errors:

			PrepareLogMessagef("An error occurred watching for file changes").
				AddError(err).
				Error()

			if !ok {
				return
			}
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
func Watch(
	ctx context.Context,
	watchPath string,
	allowedExtension *[]string,
	onNewAudioFile FileHandler,
) error {

	var err error
	var watcher *fsnotify.Watcher
	doneWatching := make(chan struct{})

	if watcher, err = fsnotify.NewWatcher(); err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}

	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			PrepareLogMessagef("Failed to close watcher").
				AddError(err).
				Error()
		}
	}(watcher)

	go HandleFileEvents(ctx, watcher, allowedExtension, onNewAudioFile, doneWatching)

	if err := watcher.Add(watchPath); err != nil {
		return fmt.Errorf("error watching folder: %w", err)
	}

	<-doneWatching

	return nil
}

// ProcessExistingFiles scans the directory at the given path for existing audio files that match
// the allowed file extensions. For each matching file, it triggers the provided onNewAudioFile handler.
//
// Parameters:
// - path: The directory path to scan for existing audio files.
// - allowedExtensions: A slice of string containing the allowed audio file extensions.
// - onNewAudioFile: The function to call when a new audio file is found.
func ProcessExistingFiles(
	fso FileSystemInterface,
	path string,
	allowedAudioExtensions *[]string,
	onNewAudioFile FileHandler,
) error {
	var err error
	var files []os.DirEntry

	if files, err = fso.ReadDir(path); err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}

	for _, file := range files {
		if IsAudioFile(file.Name(), allowedAudioExtensions) {

			PrepareLogMessagef("Existing audio file detected: %s", file.Name()).Info()

			fullFileName := fmt.Sprintf("%s/%s", path, file.Name())
			onNewAudioFile(fullFileName)
		}
	}

	return nil
}

func IsAudioFile(filename string, allowedExtensions *[]string) bool {
	ext := strings.ToLower(filepath.Ext(filename))

	for _, audioExt := range *allowedExtensions {
		if ext == audioExt {
			return true
		}
	}

	return false
}
