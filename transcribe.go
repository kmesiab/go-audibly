package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/transcribeservice"
)

const pollTimeout = 10 * time.Second

type TranscriptCallback func(transcriptionJob *transcribeservice.TranscriptionJob)

func CreateTranscriptionJob(
	transcriptionJobName,
	fullFileNameAndPath,
	inputBucketName,
	outputBucketName string,
	callback TranscriptCallback,
) {
	filename := filepath.Base(fullFileNameAndPath)
	bucketPath := fmt.Sprintf("s3://%s/%s", inputBucketName, filename)
	sess := session.Must(session.NewSession())
	transcribeSvc := transcribeservice.New(sess)
	input := &transcribeservice.StartTranscriptionJobInput{
		LanguageCode:         aws.String("en-US"),
		Media:                &transcribeservice.Media{MediaFileUri: aws.String(bucketPath)},
		MediaFormat:          aws.String("mp3"),
		TranscriptionJobName: aws.String(transcriptionJobName),
		OutputBucketName:     aws.String(outputBucketName),
	}

	_, err := transcribeSvc.StartTranscriptionJob(input)
	if err != nil {
		PrepareLogMessagef("Failed to create transcription job: %s", err.Error()).Error()

		return
	}

	go pollForTranscript(transcribeSvc, transcriptionJobName, callback)
}

func pollForTranscript(svc *transcribeservice.TranscribeService, jobName string, callback TranscriptCallback) {
	for {
		input := &transcribeservice.GetTranscriptionJobInput{
			TranscriptionJobName: aws.String(jobName),
		}

		output, err := svc.GetTranscriptionJob(input)
		if err != nil {
			PrepareLogMessagef("Failed to get transcription job: %s", err.Error()).Error()

			return
		}

		status := *output.TranscriptionJob.TranscriptionJobStatus
		switch status {
		case "COMPLETED":
			PrepareLogMessagef("Transcription job %s completed.", jobName).Info()
			callback(output.TranscriptionJob)

			return // break from loop
		case "FAILED":
			errMsg := *output.TranscriptionJob.FailureReason
			PrepareLogMessagef("Transcription job %s failed.", jobName).
				Add("reason", errMsg).
				Error()

			return // break from loop
		default:
			PrepareLogMessagef("Waiting for transcription job %s to complete.", jobName).Info()
		}

		time.Sleep(pollTimeout) // Poll every 5 seconds
	}
}
