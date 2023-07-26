package main

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_openErrorLogFile(t *testing.T) {
	tests := []struct {
		name    string
		want    *os.File
		wantErr bool
	}{
		{
			name: "Test opening/creating error log file",
			want: func() *os.File {
				// Get the current working directory
				currentDir, err := os.Getwd()
				if err != nil {
					t.Fatalf("Failed to get current working directory: %v", err)
				}

				// Construct the path for the error log file in the current directory
				logFilePath := filepath.Join(currentDir, "error_log.txt")

				// Open the error log file
				file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
				if err != nil {
					t.Fatalf("Failed to open error log file: %v", err)
				}

				return file
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := openErrorLogFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("openErrorLogFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer func() {
				// Close the file before removing it
				got.Close()

				// Remove the file from the filesystem
				if err := os.Remove(got.Name()); err != nil {
					t.Errorf("Failed to remove error log file: %v", err)
				}
			}()
		})
	}
}
