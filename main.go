package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/google/uuid"
)

var AllowedAudioExtensions = []string{".mp3"} // ".wav", ".flac", ".aac", ".ogg"

func main() {
	config, err := GetConfig()
	if err != nil {
		PrepareLogMessagef("❌ Invalid configuration: %s", err.Error()).Error()

		return
	}

	PrepareLogMessagef("🔄 Processing any existing files in folder: %s", config.WatchFolder).Info()
	ProcessExistingFiles(config.WatchFolder, AllowedAudioExtensions, handleNewAudioFile)

	PrepareLogMessagef("👀 Watching folder: %s", config.WatchFolder).Info()
	err = Watch(context.Background(), config.WatchFolder, handleNewAudioFile)

	if err != nil {
		PrepareLogMessagef("❌ Error watching folder: %s", err.Error()).Error()

		return
	}

	PrepareLogMessage("✅ Done").Info()
}

// handleNewAudioFile is called whenever a new audio file is detected in the watched directory.
// It uploads the newly detected audio file to S3, and starts a transcription job.
func handleNewAudioFile(filePathAndName string) {
	config, err := GetConfig()
	if err != nil {
		PrepareLogMessagef("❌ Invalid configuration: %s", err.Error()).Error()

		return
	}

	absPath, _ := filepath.Abs(filePathAndName)
	filename := filepath.Base(absPath)

	uploadToS3(filename, config)

	jobID := uuid.New().String()
	// Start the transcription job
	CreateTranscriptionJob(
		jobID,
		absPath,
		config.AwsS3InputBucketName,
		config.AwsS3OutputBucketName,
		handleTranscriptionCallback,
	)
}

func handleTranscriptionCallback(transcriptionJob *transcribeservice.TranscriptionJob) {
	transcript := transcriptionJob.Transcript.String()
	SaveTranscription(transcriptionJob.TranscriptionJobName, &transcript)
}

// SaveTranscription saves the transcription text to a file.
// It takes the filename and the contents of the transcription as arguments.
// The file is saved in the folder specified by config.TranscriptFolder.
// It logs errors for file creation, writing, and closing actions.
func SaveTranscription(fullFilePathAndName, contents *string) {
	config, err := GetConfig()
	if err != nil {
		PrepareLogMessagef("❌ Invalid configuration: %s", err.Error()).Error()

		return
	}

	absPath, err := filepath.Abs(*fullFilePathAndName)

	fileName := filepath.Base(absPath)

	if err != nil {
		PrepareLogMessagef("❌ Failed to get absolute path: %s", err.Error()).
			Add("filepath", absPath).
			Error()
	}

	newFilePathAndName := fmt.Sprintf("%s/%s", config.TranscriptFolder, fileName)
	// #nosec // File set to absolute above
	file, err := os.Create(newFilePathAndName)
	if err != nil {
		PrepareLogMessagef("Failed to create file: %s", newFilePathAndName).Error()
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			PrepareLogMessagef("Error closing file: %s\n", err.Error()).Error()
		}
	}(file)

	_, err = file.WriteString(*contents)
	if err != nil {
		PrepareLogMessagef("Error writing to file: %s\n", err.Error()).Error()
	}

	err = os.Remove(*fullFilePathAndName)
	if err != nil {
		PrepareLogMessagef("Error deleting audio file: %s\n", err.Error()).Error()
	}
}
