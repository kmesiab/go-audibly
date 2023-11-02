package main

import (
	"os"
)

// FileSystemInterface defines an interface for filesystem operations.
// This interface abstracts the filesystem operations that your application requires.
// Currently, it includes a method for reading a directory's contents.
type FileSystemInterface interface {
	ReadDir(dirname string) ([]os.DirEntry, error)
}

// OSFileSystem is a concrete implementation of the FileSystemInterface.
// It uses the actual filesystem provided by the os package.
type OSFileSystem struct{}

// ReadDir is a method that implements the ReadDir function from the FileSystemInterface.
// It takes a directory name as input and returns a slice of os.DirEntry and an error.
// This method directly uses os.ReadDir to read the directory contents.
func (OSFileSystem) ReadDir(dirname string) ([]os.DirEntry, error) {
	return os.ReadDir(dirname)
}
