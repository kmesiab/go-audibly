package main

import tsvc "github.com/aws/aws-sdk-go/service/transcribeservice"

type AudioTranscriptionJob struct {
	Name             string
	Filename         string
	Service          *tsvc.TranscribeService
	InputBucketName  string
	OutputBucketName string
}
