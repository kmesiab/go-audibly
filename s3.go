package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadToS3(keyName string, appConfig *Config) error {
	var err error
	var sess *session.Session

	filePathAndName := fmt.Sprintf("%s/%s", *appConfig.WatchFolder, keyName)

	if sess, err = session.NewSession(&aws.Config{
		Credentials: appConfig.AwsCredentials,
		Region:      appConfig.AwsRegion,
	}); err != nil {
		PrepareLogMessagef("Failed to create session: %s", err.Error()).Error()

		return err
	}

	// #nosec // File set to absolute above
	file, err := os.Open(filePathAndName)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			PrepareLogMessagef("Fatal Error: failed to close file").
				AddError(e).
				Add("filepath", *appConfig.WatchFolder).
				Add("filename", file.Name()).
				Error()
		}
	}(file)

	PrepareLogMessagef("Uploading file to S3: %s", *appConfig.AwsS3InputBucketName).Debug()

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(*appConfig.AwsS3InputBucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})

	if err != nil {
		return err
	} else {
		PrepareLogMessagef("File uploaded: %s", file.Name()).Info()

		completedFileName := fmt.Sprintf("%s/%s", *appConfig.ProcessedAudioFolder, keyName)
		err := os.Rename(filePathAndName, completedFileName)
		if err != nil {
			return err
		}
	}

	return nil
}
