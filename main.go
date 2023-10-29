package main

import (
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const WatchPath = "./pre_transcode"

var AllowedAudioExtensions = []string{".mp3", ".wav", ".flac", ".aac", ".ogg"}

func main() {
	Watch(WatchPath)
}

func Watch(path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		PrepareLogMessagef("Error creating watcher: %s", err.Error()).Error()
	}
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
					} else {
						PrepareLogMessagef("Ignoring new file: %s", event.Name).Info()
					}
				}
			case err = <-watcher.Errors:
				PrepareLogMessagef("The fs watcher received an error: %s", err.Error()).Error()
			}
		}
	}()

	err = watcher.Add(WatchPath)
	if err != nil {
		PrepareLogMessagef("Error watching folder: %s", err.Error()).
			Add("path", path).
			Error()
	}
	<-done
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
