// Package main handles the upload and processing of audio files to an AWS S3 bucket.
package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// uploadToS3 uploads a given file to an AWS S3 bucket.
func uploadToS3(keyName string, appConfig *Config) error {
	var file *os.File
	var sess *session.Session

	filePathAndName := fmt.Sprintf("%s/%s", *appConfig.WatchFolder, keyName)

	// Initialize AWS session
	sess, err := session.NewSession(&aws.Config{
		Credentials: appConfig.AwsCredentials,
		Region:      appConfig.AwsRegion,
	})
	if err != nil {
		return err
	}

	if file, err = openAudioFile(filePathAndName, appConfig); err != nil {
		return err
	}

	// Log and upload
	PrepareLogMessagef("Uploading file to S3: %s", *appConfig.AwsS3InputBucketName).Debug()
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(*appConfig.AwsS3InputBucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})

	if err != nil {
		return err
	}

	PrepareLogMessagef("File uploaded: %s", file.Name()).Debug()

	// Move to processed folder
	return moveToProcessedFolder(keyName, appConfig, filePathAndName)
}

// moveToProcessedFolder moves the processed file to a specific folder.
func moveToProcessedFolder(keyName string, appConfig *Config, filePathAndName string) error {
	completedFileName := fmt.Sprintf("%s/%s", *appConfig.ProcessedAudioFolder, keyName)

	return os.Rename(filePathAndName, completedFileName)
}

// openAudioFile opens the audio file for reading.
func openAudioFile(filePathAndName string, appConfig *Config) (*os.File, error) {
	file, err := os.Open(filePathAndName)
	if err != nil {
		return nil, err
	}

	defer func(file *os.File) {
		if e := file.Close(); e != nil {
			PrepareLogMessagef("Fatal Error: failed to close file").
				AddError(e).
				Add("filepath", *appConfig.WatchFolder).
				Add("filename", file.Name()).
				Error()
		}
	}(file)

	return file, nil
}
