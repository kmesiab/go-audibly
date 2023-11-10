package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	t.Setenv("ENV", "test")
	t.Setenv("APP_NAME", "go-audibly")

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

func TestLogMessagef(t *testing.T) {
	msg := LogMessagef("Formatted %s", "message")

	assert.Equal(t, "Formatted message", msg.Message)
	assert.NotNil(t, msg.Fields)
}

func TestAdd(t *testing.T) {
	msg := PrepareLogMessage("Test message").Add("key", "value")

	assert.Equal(t, "value", msg.Fields["key"])
}
