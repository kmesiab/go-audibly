// Package main handles the upload and processing of audio files to an AWS S3 bucket.
package main

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadAudioFileToInputBucket(ctx AppContext, fullFileName *string, wg *sync.WaitGroup) {
	var file *os.File
	var sess *session.Session

	// Clean up the file resources
	defer func(file *os.File) {
		if file == nil {
			return
		}

		err := file.Close()
		if err != nil {
			LogMessagef("Failed to close audio file").AddError(err).Error()
		}
	}(file)

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Credentials: ctx.Config.AwsCredentials,
		Region:      ctx.Config.AwsRegion,
	})
	if err != nil {
		onUploadAudioFileToInputBucketFailure(fullFileName, &err)

		return
	}

	if file, err = openAudioFile(*fullFileName); err != nil {
		onUploadAudioFileToInputBucketFailure(fullFileName, &err)

		return
	}

	keyName := filepath.Base(*fullFileName)

	uploader := s3manager.NewUploader(sess)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(*ctx.Config.AwsS3InputBucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})
	if err != nil {
		onUploadAudioFileToInputBucketFailure(fullFileName, &err)

		return
	}
	onUploadAudioFileToInputBucketSuccess(ctx, wg, fullFileName, &result.Location)
}

// openAudioFile opens the audio file for reading.
func openAudioFile(fullFileName string) (*os.File, error) {
	file, err := os.Open(fullFileName)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func onUploadAudioFileToInputBucketSuccess(ctx AppContext, wg *sync.WaitGroup, filename, uploadLocation *string) {
	LogMessagef("File uploaded: %s", *uploadLocation).Info()

	err := moveAudioFileToProcessedFolder(filename, ctx.Config.ProcessedAudioFolder)
	if err != nil {
		LogMessagef("Failed to move audio file to processed folder").
			Add("filename", *filename).
			AddError(err).
			Error()

		return
	}

	// Fire off the transcription in its own thread.  Do not wait for completion.
	go startTranscriptionJob(ctx, wg, filename, uploadLocation)
}

func onUploadAudioFileToInputBucketFailure(filename *string, err *error) {
	LogMessagef("File uploaded: %s", *filename).
		AddError(*err).
		Error()
}
