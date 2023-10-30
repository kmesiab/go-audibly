package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	err := os.Setenv("ENV", "test")
	if err != nil {
		PrepareLogMessagef("Error setting environment variable: %s", err.Error()).Error()
		return
	}

	err = os.Setenv("APP_NAME", "go-audibly")
	if err != nil {
		PrepareLogMessagef("Error setting environment variable: %s", err.Error()).Error()
		return
	}

	logger := GetLogger()

	assert.NotNil(t, logger)
	assert.Equal(t, "test", logger.Data["env"])
	assert.Equal(t, "go-audibly", logger.Data["app_name"])
}

func TestPrepareLogMessage(t *testing.T) {
	msg := PrepareLogMessage("Test message")

	assert.Equal(t, "Test message", msg.Message)
	assert.NotNil(t, msg.Fields)
}

func TestPrepareLogMessagef(t *testing.T) {
	msg := PrepareLogMessagef("Formatted %s", "message")

	assert.Equal(t, "Formatted message", msg.Message)
	assert.NotNil(t, msg.Fields)
}

func TestAdd(t *testing.T) {
	msg := PrepareLogMessage("Test message").Add("key", "value")

	assert.Equal(t, "value", msg.Fields["key"])
}

// For testing Info(), Debug(), Warn(), Error() you can mock the logger or capture stdout, but that's more complex.
