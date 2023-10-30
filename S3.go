package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadToS3(keyName string, config *Config) {
	sess, err := session.NewSession(&aws.Config{
		Region: &config.AwsRegion,
		Credentials: credentials.NewStaticCredentials(
			config.AwsAccessKeyID,
			config.AwsSecretAccessKey,
			"",
		),
	})
	if err != nil {
		PrepareLogMessagef("❌ Failed to create session: %s", err.Error()).Error()

		return
	}

	filePathAndName := fmt.Sprintf("%s/%s", config.WatchFolder, keyName)
	// #nosec // File set to absolute above
	file, err := os.Open(filePathAndName)
	if err != nil {
		PrepareLogMessagef("❌ Failed to open file: %s", err.Error()).
			Add("filepath", config.WatchFolder).
			Add("filename", keyName).
			Error()
	}

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			PrepareLogMessagef("❌ Failed to close file: %s", e.Error()).
				Add("filepath", config.WatchFolder).
				Add("filename", file.Name()).
				Error()
		}
	}(file)

	PrepareLogMessagef("Uploading file to S3: %s", config.AwsS3InputBucketName).Info()

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AwsS3InputBucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})

	if err != nil {
		PrepareLogMessagef("Failed to upload filer: %s", err.Error()).
			Add("filename", file.Name()).
			Error()
	} else {
		PrepareLogMessagef("File uploaded: %s", file.Name()).Info()
		completedFileName := fmt.Sprintf("%s/%s", config.ProcessedAudioFolder, keyName)
		err := os.Rename(filePathAndName, completedFileName)
		if err != nil {
			PrepareLogMessagef("Failed move file: %s", err.Error()).
				Add("filename", file.Name()).
				Error()
		}
	}
}
