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
			want: []string{"https://www.example.com/page1", "https://www.example.com/page2"},
		},
		{
			name: "Extract links from HTML with invalid URLs",
			args: args{
				content: `
					<!DOCTYPE html>
					<html>
					<body>
						<a href="https://www.example.com/page1">Link 1</a>
						<a href="InvalidLink">Invalid Link</a>
					</body>
					</html>
				`,
			},
			want: []string{"https://www.example.com/page1"},
		},
		{
			name: "HTML parsing error",
			args: args{
				content: "<html><body>Malformed HTML",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractLinks(tt.args.content)
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
