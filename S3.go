package main

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func uploadToS3(filePath, bucketName, keyName string) {

	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	sess, err := session.NewSession(&aws.Config{
		Region:      &region,
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""), // Token can be empty
	})

	absPath, err := filepath.Abs(filePath)

	if err != nil {
		PrepareLogMessagef("Failed to get absolute path: %s", err.Error()).
			Add("filepath", filePath)
	}

	// #nosec // File set to absolute above
	file, err := os.Open(absPath)

	if err != nil {
		PrepareLogMessagef("Failed to open file: %s", err.Error()).
			Add("filepath", filePath)
	}
	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			PrepareLogMessagef("Failed to close file: %s", e.Error()).
				Add("filename", file.Name()).
				Error()
		}
	}(file)

	PrepareLogMessagef("Uploading file to S3: %s", bucketName).Info()

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
		Body:   file,
	})

	if err != nil {

		PrepareLogMessagef("Failed to upload filer: %s", err.Error()).
			Add("filename", file.Name()).
			Error()
		return

	} else {

		PrepareLogMessagef("File uploaded: %s", file.Name()).Info()

		err := os.Rename(filePath, filePath+".done")

		if err != nil {
			PrepareLogMessagef("Failed to close file: %s", err.Error()).
				Add("filename", file.Name()).
				Error()
			return
		}

	}
}
