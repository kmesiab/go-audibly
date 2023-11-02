package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws/session"
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

	go func() {
		PrepareLogMessagef("🔄 Processing any existing files").
			Add("watch folder", *config.WatchFolder).
			Info()

		err = ProcessExistingFiles(
			OSFileSystem{},         // Wrapper around the actual filesystem.
			*config.WatchFolder,    // Path to the directory to scan for existing audio files.
			allowedAudioExtensions, // Allowed audio file extensions.
			NewAudioFileCallback,   // Function to call when a new audio file is detected.
		)

		if err != nil {
			PrepareLogMessagef("Error processing existing files").
				AddError(err).
				Error()
		}
	}()

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
			Add("filename", filePathAndName).
			AddError(err).
			Error()

		return
	}

	absPath, _ := filepath.Abs(filePathAndName)
	filename := filepath.Base(absPath)

	// Upload the audio file to S3
	if err = uploadToS3(filename, config); err != nil {
		PrepareLogMessagef("Failed to upload file to s3").
			Add("filename", filename).
			Add("abs_path", absPath).
			AddError(err).
			Error()

		return
	}

	jobName := uuid.New().String()
	sess := session.Must(session.NewSession())

	job := &AudioTranscriptionJob{
		Name:             jobName,
		Filename:         filename,
		Service:          transcribeservice.New(sess),
		InputBucketName:  *config.AwsS3InputBucketName,
		OutputBucketName: *config.AwsS3OutputBucketName,
	}

	// Start the transcription job
	err = CreateTranscriptionJob(job, TranscriptionCallback)

	if err != nil {
		PrepareLogMessagef("Failed to create transcribe job").
			AddError(err).
			Add("absolute path", absPath).
			Add("filename", filename).
			Add("job ID", jobName).
			Error()
	}
}

func TranscriptionCallback(job *AudioTranscriptionJob, transcriptionJob *transcribeservice.TranscriptionJob) {
	transcript := transcriptionJob.Transcript.String()
	SaveTranscription(job, &transcript)
}

// SaveTranscription saves the transcription text to a file.
// It takes the filename and the contents of the transcription as arguments.
// The file is saved in the folder specified by config.TranscriptFolder.
// It logs errors for file creation, writing, and closing actions.
func SaveTranscription(job *AudioTranscriptionJob, contents *string) {
	config, err := GetConfig()
	if err != nil {
		PrepareLogMessagef("Invalid configuration").
			AddTranscribeJob(job).
			AddError(err).
			Error()

		return
	}

	fileName := filepath.Base(job.Filename)
	newFilePathAndName := fmt.Sprintf("%s/%s", *config.TranscriptFolder, fileName)

	// #nosec // File set to absolute above
	file, err := os.Create(newFilePathAndName)
	if err != nil {
		PrepareLogMessagef("Failed to create file: %s", newFilePathAndName).
			AddTranscribeJob(job).
			AddError(err).
			Error()

		return
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			PrepareLogMessagef("Error closing file").
				AddTranscribeJob(job).
				AddError(err).
				Error()
		}
	}(file)

	_, err = file.WriteString(*contents)
	if err != nil {
		PrepareLogMessagef("Error writing to file").
			AddTranscribeJob(job).
			AddError(err).
			Error()

		return
	}

	err = os.Remove(job.Filename)
	if err != nil {
		PrepareLogMessagef("Error deleting audio file").
			AddTranscribeJob(job).
			AddError(err).
			Error()

		return
	}
}
