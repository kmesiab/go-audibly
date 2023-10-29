package main

import (
	"fmt"
	"os"
)

const WatchPath = "./pre_transcode"

var (
	s3BucketName = "myp-pre-transcode" // overridden by environment variable
)

var AllowedAudioExtensions = []string{".mp3", ".wav", ".flac", ".aac", ".ogg"}

func handleNewAudioFile(filename string) {
	fullFileName := fmt.Sprintf("%s/%s", WatchPath, filename)
	uploadToS3(fullFileName, s3BucketName, filename)
}

func main() {
	if !validConfig() {
		PrepareLogMessagef("AWS configuration not found in environment").Error()
		return
	}

	PrepareLogMessagef("\n👀 Watching folder: %s\n", WatchPath).Info()

	Watch(WatchPath, handleNewAudioFile)

}

func validConfig() bool {

	awsRegion := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	s3BucketName = os.Getenv("AWS_S3_BUCKET_NAME")

	return !(accessKey == "" || secretKey == "" || s3BucketName == "" || awsRegion == "")
}
