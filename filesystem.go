package main

import (
	"os"
)

func FilterNonDirFilenames(entries []os.DirEntry) []string {
	filenames := make([]string, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filenames = append(filenames, entry.Name())
	}

	return filenames
}

func GetFilenames(folder string) (*[]string, error) {
	dirEntries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0, len(dirEntries))

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		filenames = append(filenames, entry.Name())
	}

	return &filenames, nil
}
