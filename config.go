package main

import (
	"fmt"
	"reflect"
	"sync"

	goenv "github.com/Netflix/go-env"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

var (
	once   = sync.Once{}
	config *Config
)

type Config struct {
	AWSConfig
	WatchFolder          *string `env:"WATCH_FOLDER"`
	TranscriptFolder     *string `env:"TRANSCRIPT_FOLDER"`
	ProcessedAudioFolder *string `env:"PROCESSED_AUDIO_FOLDER"`
}

type AWSConfig struct {
	AwsRegion             *string                  `env:"AWS_REGION"`
	AwsAccessKeyID        *string                  `env:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey    *string                  `env:"AWS_SECRET_ACCESS_KEY"`
	AwsS3InputBucketName  *string                  `env:"AWS_S3_INPUT_BUCKET_NAME"`
	AwsS3OutputBucketName *string                  `env:"AWS_S3_OUTPUT_BUCKET_NAME"`
	AwsCredentials        *credentials.Credentials `env:"-"`
}

func GetConfig() (*Config, error) {
	once.Do(func() {
		config = &Config{}
		_, err := goenv.UnmarshalFromEnviron(config)
		if err != nil {
			PrepareLogMessagef("Failed to unmarshal config").
				AddError(err).
				Error()

			return
		}

		config.AwsCredentials = credentials.NewStaticCredentials(
			*config.AwsAccessKeyID,
			*config.AwsSecretAccessKey,
			"",
		)
	})

	err := ValidateConfig(config)

	return config, err
}

func ValidateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	v := reflect.ValueOf(config)

	// Check if it's a pointer and dereference it
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Check if it's a struct type
	if v.Kind() != reflect.Struct {

		return fmt.Errorf("config is not a struct")
	}

	// Now you can safely call NumField
	m := v.NumField()

	for i := 0; i < m; i++ {
		if v.Field(i).String() == "" {

			return fmt.Errorf("%s must not be empty", v.Type().Field(i).Name)
		}
	}

	return nil
}
