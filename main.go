package main

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

type AppContext struct {
	context.Context
	Config *Config
}

type AudioFileCallback func(ctx AppContext, fullFileName *string, wg *sync.WaitGroup)

func main() {
	cfg, err := GetConfig()
	if err != nil {
		LogMessagef("Invalid configuration").Error()

		return
	}

	existingFiles, err := GetFilenames(*cfg.WatchFolder)
	if err != nil {
		LogMessagef("Failed to read existing files").Error()

		return
	}

	ctx := AppContext{context.Background(), cfg}
	mainThreadWaitGroup := &sync.WaitGroup{}

	if len(*existingFiles) > 0 {
		processExistingFiles(ctx, *existingFiles, mainThreadWaitGroup,
			func(ctx AppContext, fullFileName *string, wg *sync.WaitGroup) {
				wg.Add(1)
				go func() {
					defer wg.Done()
					go onFileDetected(ctx, fullFileName, wg, uploadAudioFileToInputBucket)
				}()
			})
	}

	LogMessagef("✅ Go-Audibly Audio Transcription Service started.\n").Info()
	LogMessagef("👀 Watching 📁 %s for audio files.\n", *cfg.WatchFolder).Info()

	mainThreadWaitGroup.Add(1)
	go watch(ctx, mainThreadWaitGroup, uploadAudioFileToInputBucket)
	mainThreadWaitGroup.Wait()
}

func watch(
	ctx AppContext,
	wg *sync.WaitGroup,
	onNewAudioFile AudioFileCallback,
) {
	defer wg.Done() // This will be called when the function exits

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		LogMessagef("⚠️Failed to create file watcher").AddError(err).Error()

		return
	}

	err = watcher.Add(*ctx.Config.WatchFolder)

	if err != nil {
		LogMessagef("⚠️Failed to add folder to watcher").
			AddError(err).Add("watch_folder", *ctx.Config.WatchFolder).Error()
	}

	for {
		select {

		case event, ok := <-watcher.Events:

			if !ok { // 'watcher.Events' has been closed
				LogMessagef("🚨Watcher channel closed").Fatal()

				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				wg.Add(1)

				go func(eventName string) {
					defer wg.Done()

					LogMessagef("🔊 Audio file detected: %s", eventName).Info()

					onNewAudioFile(ctx, &eventName, wg)
				}(event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok { // 'watcher.Errors' has been closed
				LogMessagef("🚨 Watcher has been stopped; no more errors will be received.").Info()

				return
			}

			LogMessagef("🚨 An error was received while watching").AddError(err).Error()

			continue
		}
	}
}

func processExistingFiles(ctx AppContext, filenames []string, wg *sync.WaitGroup, callback OnFileDetectedCallback) {
	for _, filename := range filenames {

		filename = *ctx.Config.WatchFolder + "/" + filename

		absFilename, err := filepath.Abs(filename)
		if err != nil {
			LogMessagef("Failed to get absolute path").
				AddError(err).
				Add("filename", filename).
				Error()

			continue
		}

		if IsAudioFile(absFilename, ctx.Config.AllowedAudioExtensions) {
			wg.Add(1)

			go func(filename string) {
				defer wg.Done()
				callback(ctx, &filename, wg)
			}(absFilename)
		}
	}
}

func moveAudioFileToProcessedFolder(filenameAndPath, processedFolder *string) error {
	filename := filepath.Base(*filenameAndPath)
	newFilename := *processedFolder + "/" + filename
	newAbsFilename, _ := filepath.Abs(newFilename)

	return os.Rename(*filenameAndPath, newAbsFilename)
}
