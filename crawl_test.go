package main

import (
	"reflect"
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
