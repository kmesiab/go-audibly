package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/google/uuid"
)

func main() {
	var (
		err                    error
		config                 *Config
		allowedAudioExtensions = &[]string{".mp3", ".wav", ".flac", ".aac", ".ogg"}
	)

	if config, err = GetConfig(); err != nil {
		PrepareLogMessagef("Invalid configuration").
			AddError(err).
			Error()

		return
	}

	PrepareLogMessagef("🔄 Processing any existing files in folder: %s", *config.WatchFolder).Info()

	err = ProcessExistingFiles(*config.WatchFolder, allowedAudioExtensions, NewAudioFileCallback)

	if err != nil {
		PrepareLogMessagef("Error processing existing files").
			AddError(err).
			Error()

		return
	}

	PrepareLogMessagef("👀 Watching folder: %s", *config.WatchFolder).Info()

	err = Watch(context.Background(), *config.WatchFolder, allowedAudioExtensions, NewAudioFileCallback)

	if err != nil {
		PrepareLogMessagef("Error watching folder").
			AddError(err).
			Error()

		return
	}

	PrepareLogMessage("✅ Done").Info()
}

// NewAudioFileCallback is called whenever a new audio file is detected in the watched directory.
// It uploads the newly detected audio file to S3, and starts a transcription job.
func NewAudioFileCallback(filePathAndName string) {
	config, err := GetConfig()
	if err != nil {
		PrepareLogMessagef("Invalid configuration").
			AddError(err).
			Error()

		return
	}

	absPath, _ := filepath.Abs(filePathAndName)
	filename := filepath.Base(absPath)

	// Upload the audio file to S3
	if err = uploadToS3(filename, config); err != nil {
		PrepareLogMessagef("Failed to upload file to s3").
			AddError(err).
			Error()

		return
	}

	jobID := uuid.New().String()
	// Start the transcription job
	CreateTranscriptionJob(
		jobID,
		absPath,
		*config.AwsS3InputBucketName,
		*config.AwsS3OutputBucketName,
		TranscriptionCallback,
	)
}

func TranscriptionCallback(transcriptionJob *transcribeservice.TranscriptionJob) {
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
		PrepareLogMessagef("Invalid configuration").
			AddError(err).
			Error()

		return
	}

	absPath, err := filepath.Abs(*fullFilePathAndName)
	fileName := filepath.Base(absPath)

	if err != nil {
		PrepareLogMessagef("Failed to get absolute path").
			AddError(err).
			Error()

		return
	}

	newFilePathAndName := fmt.Sprintf("%s/%s", *config.TranscriptFolder, fileName)
	// #nosec // File set to absolute above
	file, err := os.Create(newFilePathAndName)
	if err != nil {
		PrepareLogMessagef("Failed to create file: %s", newFilePathAndName).Error()

		return
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			PrepareLogMessagef("Error closing file").
				AddError(err).
				Error()
		}
	}(file)

	_, err = file.WriteString(*contents)
	if err != nil {
		PrepareLogMessagef("Error writing to file").
			AddError(err).
			Error()

		return
	}

	err = os.Remove(*fullFilePathAndName)
	if err != nil {
		PrepareLogMessagef("Error deleting audio file").
			AddError(err).
			Error()

		return
	}
}
