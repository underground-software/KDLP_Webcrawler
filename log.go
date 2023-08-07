package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func getCurrentDirectory() (string, error) {
	return os.Getwd()
}

func openErrorLogFile(logFilePath string) (*os.File, error) {

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

func initializeErrorLogging() {
	// Get the current working directory
	currentDir, err := getCurrentDirectory()
	if err != nil {
		log.Fatal("Failed to get current working directory:", err)
	}

	// Construct the path for the error log file in the current directory
	logFilePath := filepath.Join(currentDir, "error_log.txt")

	// Sets logged errors to print to both error log file and terminal
	_, err = openErrorLogFile(logFilePath)
	if err != nil {
		log.Fatal("Failed to open error log file:", err)
	}
}
