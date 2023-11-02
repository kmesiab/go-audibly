package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	tsvc "github.com/aws/aws-sdk-go/service/transcribeservice"
)

const (
	transcriptionResultSuccess = "COMPLETED"
	transcriptionResultError   = "FAILED"
	pollTimeout                = 10 * time.Second
)

type TranscriptCallback func(
	job *AudioTranscriptionJob,
	transcriptionJob *tsvc.TranscriptionJob,
)

func CreateTranscriptionJob(job *AudioTranscriptionJob, callback TranscriptCallback) error {
	filename := filepath.Base(job.Filename)
	bucketPath := fmt.Sprintf("s3://%s/%s", job.InputBucketName, filename)

	input := &tsvc.StartTranscriptionJobInput{
		LanguageCode:         aws.String("en-US"),
		Media:                &tsvc.Media{MediaFileUri: aws.String(bucketPath)},
		MediaFormat:          aws.String("mp3"),
		TranscriptionJobName: aws.String(job.Name),
		OutputBucketName:     aws.String(job.OutputBucketName),
	}

	_, err := job.Service.StartTranscriptionJob(input)
	if err != nil {
		PrepareLogMessagef("Failed to create transcription job: %s", err.Error()).Error()

		return err
	}

	go pollForTranscript(job, callback)

	return nil
}

func pollForTranscript(job *AudioTranscriptionJob, callback TranscriptCallback) {
	for {
		input := &tsvc.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String(job.Name),
		}

		output, err := job.Service.GetTranscriptionJob(input)
		if err != nil {
			PrepareLogMessagef("Failed to get transcription job: %s", err.Error()).Error()

			return
		}

		status := *output.TranscriptionJob.TranscriptionJobStatus

		switch status {
		case transcriptionResultSuccess:
			PrepareLogMessagef("Transcription job %s completed.", job.Name).
				Add("file name", job.Filename).
				Add("job name", job.Name).Info()

			go callback(job, output.TranscriptionJob)

			return
		case transcriptionResultError:
			PrepareLogMessagef("Transcription job %s failed.", job.Name).
				Add("reason", *output.TranscriptionJob.FailureReason).
				Add("file name", job.Filename).
				Add("job name", job.Name).
				Error()

			return
		default:
			PrepareLogMessagef("Waiting for transcription job %s to complete.", job.Name).
				Add("file name", job.Filename).
				Add("job name", job.Name).
				Info()
		}

		time.Sleep(pollTimeout) // Poll every 5 seconds
	}
}
