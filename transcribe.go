package main

import (
	"errors"
	"os"
	"path"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	tsvc "github.com/aws/aws-sdk-go/service/transcribeservice"
	"github.com/google/uuid"
)

const (
	transcriptionResultSuccess    = "COMPLETED"
	transcriptionResultError      = "FAILED"
	pollTimeout                   = 5 * time.Second
	maxNumberOfSpeakersToIdentify = 5
)

func onTranscriptionJobDownloadFailed(jobName *string, err error) {
	LogMessagef("Failed to download transcription file. Canceling context.").
		AddError(err).
		Add("job name", *jobName).
		Error()
}

func startTranscriptionJob(ctx AppContext, wg *sync.WaitGroup, filename, uploadLocation *string) {
	jobName := path.Base(*filename) + "_" + uuid.New().String()

	input := &tsvc.StartTranscriptionJobInput{
		LanguageCode:         aws.String("en-US"),
		Media:                &tsvc.Media{MediaFileUri: aws.String(*uploadLocation)},
		MediaFormat:          aws.String("mp3"), // TODO: Add support for other formats
		TranscriptionJobName: aws.String(jobName),
		OutputBucketName:     aws.String(*ctx.Config.AwsS3OutputBucketName),
		Settings: &tsvc.Settings{
			ShowSpeakerLabels: aws.Bool(true),
			MaxSpeakerLabels:  aws.Int64(maxNumberOfSpeakersToIdentify), // Number of speakers to identify
		},
	}

	sess := session.Must(session.NewSession())
	svc := tsvc.New(sess)

	_, err := svc.StartTranscriptionJob(input)
	if err != nil {
		onTranscriptionJobFailed(&jobName, err)

		return
	}

	wg.Add(1)

	go func() {
		defer wg.Done()
		onTranscribeJobCreated(ctx, wg, &jobName, svc)
	}()
}

func downloadTranscriptionFile(ctx AppContext, job *tsvc.TranscriptionJob) {
	filePath := *ctx.Config.TranscriptFolder + "/" + *job.TranscriptionJobName + ".json"

	file, err := os.Create(filePath)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			LogMessagef("Error closing file").
				AddError(err).
				Add("filename", filePath).
				Error()
		}
	}(file)

	if err != nil {
		onTranscriptionJobDownloadFailed(job.TranscriptionJobName, err)

		return
	}

	_, err = file.Write([]byte(job.String()))

	if err != nil {
		onTranscriptionJobDownloadFailed(job.TranscriptionJobName, err)

		return
	}

	onTranscriptionJobDownloadSuccess(job)
}

func pollForTranscriptJob(
	ctx AppContext,
	wg *sync.WaitGroup,
	jobName *string,
	svc *tsvc.TranscribeService,
	callback func(ctx AppContext, svc *tsvc.TranscriptionJob),
) {
	input := &tsvc.GetTranscriptionJobInput{
		TranscriptionJobName: aws.String(*jobName),
	}

	for {

		output, err := svc.GetTranscriptionJob(input)
		if err != nil {
			LogMessagef("Failed to get transcription job: %s", err.Error()).Error()
			wg.Done()

			return
		}

		status := *output.TranscriptionJob.TranscriptionJobStatus

		switch status {

		case transcriptionResultSuccess:

			wg.Add(1)

			go func() {
				defer wg.Done()
				callback(ctx, output.TranscriptionJob)
			}()

			return

		case transcriptionResultError:

			err := errors.New(*output.TranscriptionJob.FailureReason)
			onTranscriptionJobFailed(jobName, err)

			return

		default:
			go onTranscriptJobWaiting(jobName)
		}

		// TODO: Move this to config
		time.Sleep(pollTimeout) // Poll every 5 seconds
	}
}

func onTranscribeJobCreated(ctx AppContext, wg *sync.WaitGroup, jobName *string, svc *tsvc.TranscribeService) {
	LogMessagef("Transcribe job created: %s", *jobName).Info()

	wg.Add(1)

	defer func() {
		defer wg.Done()
		go pollForTranscriptJob(ctx, wg, jobName, svc, onTranscriptionJobCompleted)
	}()
}

func onTranscriptionJobCompleted(ctx AppContext, job *tsvc.TranscriptionJob) {
	LogMessagef("Transcribe job completed: %s", *job.TranscriptionJobName).Info()
	// TODO: use redaction url
	go downloadTranscriptionFile(ctx, job)
	// TODO: rename file on s3 bucket
}

func onTranscriptionJobFailed(jobName *string, err error) {
	LogMessagef("Transcription job failed").
		AddError(err).
		Add("job name", *jobName).
		Error()
}

func onTranscriptJobWaiting(jobName *string) {
	LogMessagef("Waiting for transcription job: %s", *jobName).Info()
}

func onTranscriptionJobDownloadSuccess(job *tsvc.TranscriptionJob) {
	LogMessagef("File successfully transcribed").
		Add("job name", *job.TranscriptionJobName).
		Add("transcript url", *job.Transcript.TranscriptFileUri).
		Info()
}
