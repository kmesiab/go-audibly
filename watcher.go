package main

import (
	"path/filepath"
	"strings"
	"sync"
)

type OnFileDetectedCallback func(ctx AppContext, fullFileName *string, wg *sync.WaitGroup)

// FileHandler defines the type of the function to be called when a new file is detected
// by the fsnotify Watcher in the Watch function. It takes the full path of the newly
// created file as its argument.
//
// Functions matching this type can be used to further process the new file, such as uploading it to S3
// or initiating a transcription job, as demonstrated in the HandleFileEvents and Watch functions.
type FileHandler func(string)

// ProcessExistingFiles scans the directory at the given path for existing audio files that match
// the allowed file extensions. For each matching file, it triggers the provided onNewAudioFile handler.
//
// Parameters:
// - path: The directory path to scan for existing audio files.
// - allowedExtensions: A slice of string containing the allowed audio file extensions.
// - onNewAudioFile: The function to call when a new audio file is found.
func ProcessExistingFiles(
	filePaths *[]string,
	allowedAudioExtensions *[]string,
	onNewAudioFile FileHandler,
) {
	for _, file := range *filePaths {
		if IsAudioFile(file, allowedAudioExtensions) {
			LogMessagef("Existing audio file detected: %s", file).Info()

			onNewAudioFile(file)
		}
	}
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

func onFileDetected(ctx AppContext, filename *string, wg *sync.WaitGroup, callback OnFileDetectedCallback) {
	if IsAudioFile(*filename, ctx.Config.AllowedAudioExtensions) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			callback(ctx, filename, wg)
		}()

	}
}
