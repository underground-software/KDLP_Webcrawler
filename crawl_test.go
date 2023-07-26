package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_extractLinks(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name    string
		args    args
		domain  string
		want    []string
		wantErr bool
	}{
		{
			name: "Extract links from HTML with valid URLs",
			args: args{
				content: `
					<!DOCTYPE html>
					<html>
					<body>
						<a href="https://www.example.com/page1">Link 1</a>
						<a href="https://www.example.com/page2">Link 2</a>
					</body>
					</html>
				`,
			},
			domain: "https://www.example.com/",
			want:   []string{"https://www.example.com/page1", "https://www.example.com/page2"},
		},
		{
			name: "Extract links from HTML with invalid URLs",
			args: args{
				content: `
					<!DOCTYPE html>
					<html>
					<body>
						<a href="https://www.example.com/page1">Link 1</a>
						<a href="relative.html">Relative link</a>
					</body>
					</html>
				`,
			},
			domain: "https://www.example.com/",
			want:   []string{"https://www.example.com/page1", "https://www.example.com/relative.html"},
		},
		{
			name: "HTML parsing error",
			args: args{
				content: "<html><body>Malformed HTML",
			},
			domain: "https://www.example.com/",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractLinks(tt.args.content, tt.domain)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCrawler(t *testing.T) {
	tests := []struct {
		name string
		want *Crawler
	}{
		{
			name: "Test case 1",
			want: &Crawler{
				visited: make(map[string]bool),
			},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCrawler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCrawler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrawler_writeDeadLinksToFile(t *testing.T) {
	type fields struct {
		deadLinks []string
	}
	type args struct {
		filepath string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantOutput []string // Expected output lines in the file
	}{
		{
			name: "Write dead links to file",
			fields: fields{
				deadLinks: []string{"https://www.example.com/deadlink1", "https://www.example.com/deadlink2"},
			},
			args: args{
				filepath: "test_dead_links.txt",
			},
			wantErr:    false,
			wantOutput: []string{"https://www.example.com/deadlink1", "https://www.example.com/deadlink2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				deadLinks: tt.fields.deadLinks,
			}
			defer os.Remove(tt.args.filepath) // Cleanup the temporary file

			err := c.writeDeadLinksToFile(tt.args.filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crawler.writeDeadLinksToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Read the file and check its content
				fileContent, err := os.ReadFile(tt.args.filepath)
				if err != nil {
					t.Errorf("Failed to read the file: %v", err)
					return
				}

				outputLines := strings.Split(string(fileContent), "\n")
				// Ignore the last empty line
				outputLines = outputLines[:len(outputLines)-1]

				// Check if the output matches the expected output
				if !reflect.DeepEqual(outputLines, tt.wantOutput) {
					t.Errorf("Unexpected file content. Got: %v, want: %v", outputLines, tt.wantOutput)
				}
			}
		})
	}
}
