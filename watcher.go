package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var once sync.Once

type FileHandler func(string) // define a FileHandler type that takes a string argument (filename)

func Watch(path string, onNewAudioFile FileHandler) {

	// Process existing files on initialization
	once.Do(func() {
		ProcessExistingFiles(path, AllowedAudioExtensions, s3BucketName)
	})

	// Start a new watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		PrepareLogMessagef("Error creating watcher: %s", err.Error()).Error()
	}

	// Defer the cleanup
	defer func(watcher *fsnotify.Watcher) {
		e := watcher.Close()
		if err != nil {
			PrepareLogMessagef("Error closing watcher: %s", e.Error()).Error()
		}
	}(watcher)

	done := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Create == fsnotify.Create {
					if IsAudioFile(event.Name, AllowedAudioExtensions) {
						PrepareLogMessagef("New audio file detected: %s", event.Name).Info()
						// Call the callback function here
						onNewAudioFile(event.Name)
					} else {
						PrepareLogMessagef("Ignoring new file: %s", event.Name).Info()
					}
				}
			case err = <-watcher.Errors:
				PrepareLogMessagef("The fs watcher received an error: %s", err.Error()).Error()
			}
		}
	}()

	// Start watching
	err = watcher.Add(WatchPath)

	if err != nil {
		PrepareLogMessagef("Error watching folder: %s", err.Error()).
			Add("path", path).
			Error()
	}
	<-done
}

// ProcessExistingFiles Handle existing files on initialization
func ProcessExistingFiles(path string, allowedExtensions []string, s3BucketName string) {
	files, err := os.ReadDir(path)
	if err != nil {
		PrepareLogMessagef("Error reading directory: %s", err.Error()).Error()
		return
	}

	for _, file := range files {
		if IsAudioFile(file.Name(), allowedExtensions) {
			PrepareLogMessagef("Existing audio file detected: %s", file.Name()).Info()
			uploadToS3(path+"/"+file.Name(), s3BucketName, file.Name())
		}
	}
}

// IsAudioFile checks if a file is an audio file based on its extension
func IsAudioFile(filename string, allowedExtensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, audioExt := range allowedExtensions {
		if ext == audioExt {
			return true
		}
	}
	return false
}
