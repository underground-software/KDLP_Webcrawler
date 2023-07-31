package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_extractValidLinks(t *testing.T) {
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
			name: "Extract links from HTML with relative URLs",
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
			name: "Extract links from HTML with http scheme",
			args: args{
				content: `
				<!DOCTYPE html>
				<html>
				<body>
					<a href="http://www.example.com/page1">Link 1</a>
					<a href="relative.html">Relative link</a>
				</body>
				</html>
			`,
			},
			domain: "https://www.example.com/",
			want:   []string{"http://www.example.com/page1", "https://www.example.com/relative.html"},
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
			got := extractValidLinks(tt.args.content, tt.domain)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: Improve?
func Test_newCrawler(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		homeURL string
		want    *Crawler
	}{
		{
			name:    "Test case 1",
			domain:  "https://website.test/",
			homeURL: "https://website.test/index.html",
			want: &Crawler{
				visited:   make(map[string]bool),
				deadLinks: []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newCrawler(tt.domain, tt.homeURL)

			// Check the domain and homeURL fields directly
			if got.domain != tt.domain || got.homeURL != tt.homeURL {
				t.Errorf("newCrawler() = %v, want %v", got, tt.want)
			}

			// Compare the visited and deadLinks fields using reflect.DeepEqual
			if !reflect.DeepEqual(got.visited, tt.want.visited) || !reflect.DeepEqual(got.deadLinks, tt.want.deadLinks) {
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

	t.Cleanup(func() {
		if err := os.Remove("test_dead_links.txt"); err != nil {
			t.Fatalf("Error removing test_dead_links.txt: %v", err)
		}
	})
}

func TestCrawler_handleDeadLink(t *testing.T) {
	type fields struct {
		visited   map[string]bool
		deadLinks []string
	}
	type args struct {
		referringURL string
		URL          string
		statusCode   int
	}

	tests := []struct {
		name       string
		fields     fields
		args       args
		wantLinks  []string // Expected deadLinks after calling handleDeadLink
		wantFile   string   // Expected contents of the dead links file
		wantErrMsg string   // Expected error message
	}{
		{
			name: "Handle Dead Link",
			fields: fields{
				visited:   map[string]bool{},
				deadLinks: []string{},
			},
			args: args{
				referringURL: "https://example.com",
				URL:          "https://example.com/deadlink",
				statusCode:   404,
			},
			wantLinks: []string{"dead link https://example.com/deadlink found at: https://example.com"},
			wantFile:  "dead link https://example.com/deadlink found at: https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				visited:   tt.fields.visited,
				deadLinks: tt.fields.deadLinks,
			}

			// Call handleDeadLink with the given arguments
			c.handleDeadLink(tt.args.referringURL, tt.args.URL, tt.args.statusCode)

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

func TestCrawler_crawlURL(t *testing.T) {
	type fields struct {
		visited   map[string]bool
		deadLinks []string
	}
	type args struct {
		URL          string
		referenceURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Handle Dead Link",
			fields: fields{
				visited:   map[string]bool{},
				deadLinks: []string{},
			},
			args: args{
				URL:          "https://example.com/deadlink",
				referenceURL: "https://example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				visited:   tt.fields.visited,
				deadLinks: tt.fields.deadLinks,
			}
			c.crawlURL(tt.args.URL, tt.args.referenceURL)
		})
	}
}

func TestCrawler_crawlInternalURL(t *testing.T) {
	type fields struct {
		domain    string
		homeURL   string
		visited   map[string]bool
		deadLinks []string
	}
	type args struct {
		URL          string
		referringURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Valid internal URL with links",
			fields: fields{
				domain:    "https://website.I.Am.Testing/",
				homeURL:   "https://website.I.Am.Testing/index.html",
				visited:   make(map[string]bool),
				deadLinks: []string{},
			},
			args: args{
				URL:          "https://website.I.Am.Testing/page1",
				referringURL: "https://website.I.Am.Testing/index.html",
			},
		},
		{
			name: "Invalid URL",
			fields: fields{
				domain:    "https://website.I.Am.Testing/",
				homeURL:   "https://website.I.Am.Testing/index.html",
				visited:   make(map[string]bool),
				deadLinks: []string{},
			},
			args: args{
				URL:          "invalid",
				referringURL: "https://website.I.Am.Testing/index.html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crawler{
				domain:    tt.fields.domain,
				homeURL:   tt.fields.homeURL,
				visited:   tt.fields.visited,
				deadLinks: tt.fields.deadLinks,
			}
			c.crawlInternalURL(tt.args.URL, tt.args.referringURL)

		})
	}
}
