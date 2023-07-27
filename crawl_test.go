package main

import (
	"fmt"
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

func Test_saveDeadLinksToFile(t *testing.T) {
	type args struct {
		deadLinks []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test writing multiple dead links",
			args: args{
				deadLinks: []string{
					"https://www.example.com/deadlink1",
					"https://www.example.com/deadlink2",
					"https://www.example.com/deadlink3",
				},
			},
			wantErr: false,
		},
		{
			name: "Test writing single dead link",
			args: args{
				deadLinks: []string{
					"https://www.example.com/deadlink4",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Call the function being tested
			filepath := "test_dead_links.txt"
			err := saveDeadLinksToFile(filepath, tt.args.deadLinks)
			if (err != nil) != tt.wantErr {
				t.Errorf("saveDeadLinksToFile() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestCrawler_handleDeadLink(t *testing.T) {
	type fields struct {
		visited   map[string]bool
		deadLinks []string
	}
	type args struct {
		URL        string
		statusCode int
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantLinks  []string // Expected deadLinks after calling handleDeadLink
		wantFile   string   // Expected contents of the dead links file
		wantErrMsg string   // Expected error message (empty if no error)
	}{
		{
			name: "Handle Dead Link",
			fields: fields{
				visited:   map[string]bool{},
				deadLinks: []string{},
			},
			args: args{
				URL:        "https://example.com/deadlink",
				statusCode: 404,
			},
			wantLinks:  []string{"https://example.com/deadlink"},
			wantFile:   "https://example.com/deadlink",
			wantErrMsg: "", // No error expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				visited:   tt.fields.visited,
				deadLinks: tt.fields.deadLinks,
			}

			// Call handleDeadLink
			c.handleDeadLink(tt.args.URL, tt.args.statusCode)

			// Verify statusCode
			if tt.args.statusCode != 404 {
				t.Errorf("handleDeadLink() unexpected statusCode.\nGot: %v\nWant: 404", tt.args.statusCode)
			}

			// Verify deadLinks slice
			if !equalStringSlices(c.deadLinks, tt.wantLinks) {
				t.Errorf("handleDeadLink() unexpected deadLinks.\nGot: %v\nWant: %v", c.deadLinks, tt.wantLinks)
			}

			// Verify the contents of the dead links file
			fileContent, err := os.ReadFile("dead_links.txt")
			if err != nil {
				t.Fatalf("Error reading dead links file: %v", err)
			}
			gotFileContent := strings.TrimSpace(string(fileContent))
			if gotFileContent != tt.wantFile {
				t.Errorf("handleDeadLink() unexpected file content.\nGot: %v\nWant: %v", gotFileContent, tt.wantFile)
			}

			// Print debugging output
			fmt.Println("Got deadLinks:", c.deadLinks)
			fmt.Println("Got file content:", gotFileContent)
		})
	}
}

// Helper function to compare two string slices
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
