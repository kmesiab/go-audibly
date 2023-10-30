package main

import (
	"fmt"
	"reflect"

	goenv "github.com/Netflix/go-env"
)

type Config struct {
	AWSConfig
	WatchFolder          string `env:"WATCH_FOLDER"`
	TranscriptFolder     string `env:"TRANSCRIPT_FOLDER"`
	ProcessedAudioFolder string `env:"PROCESSED_AUDIO_FOLDER"`
}

type AWSConfig struct {
	AwsRegion             string `env:"AWS_REGION"`
	AwsAccessKeyID        string `env:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey    string `env:"AWS_SECRET_ACCESS_KEY"`
	AwsS3InputBucketName  string `env:"AWS_S3_INPUT_BUCKET_NAME"`
	AwsS3OutputBucketName string `env:"AWS_S3_OUTPUT_BUCKET_NAME"`
}

func GetConfig() (*Config, error) {
	config := &Config{}
	_, err := goenv.UnmarshalFromEnviron(config)
	if err != nil {
		return nil, err
	}

	err = ValidateConfig(config)

	return config, err
}

func ValidateConfig(config *Config) error {
	v := reflect.ValueOf(*config)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			return fmt.Errorf("%s must not be empty", v.Type().Field(i).Name)
		}
	}

	return nil
}
