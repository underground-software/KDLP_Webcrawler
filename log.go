package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func openErrorLogFile() (*os.File, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct the path for the error log file in the current directory
	logFilePath := filepath.Join(currentDir, "error_log.txt")

	// Open (or create then open) the error log file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// Create a multi-writer to print log to both the file and terminal simultaneously
	logOutput := io.MultiWriter(file, os.Stderr)

	// Set the log output to the multi-writer
	log.SetOutput(logOutput)

	return file, nil
}
