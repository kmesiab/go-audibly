package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockFileSystem implements FileSystemInterface for testing.
type MockFileSystem struct {
	mock.Mock
}

func (m *MockFileSystem) ReadDir(dirname string) ([]os.DirEntry, error) {
	var ok bool
	var dirEntries []os.DirEntry
	args := m.Called(dirname)

	// Safely assert the type of args.Get(0)
	if dirEntries, ok = args.Get(0).([]os.DirEntry); !ok {
		return []os.DirEntry{}, args.Error(1)
	}

	return dirEntries, args.Error(1)
}

func TestProcessExistingFiles(t *testing.T) {
	allowedExtensions := &[]string{".mp3", ".wav"}
	mockFS := new(MockFileSystem)

	t.Run("NoFiles", func(t *testing.T) {
		mockFS.On("ReadDir", mock.Anything).Return([]os.DirEntry{}, nil)

		err := ProcessExistingFiles(mockFS, "/test/path", allowedExtensions, func(string) {})
		require.NoError(t, err)
		mockFS.AssertExpectations(t)
	})

	// Additional test cases would go here...
}
