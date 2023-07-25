package main

import "testing"

func Test_isValidURL(t *testing.T) {
	type args struct {
		URL string
	}

	tests := []struct {
		name    string
		args    args
		isValid bool
	}{
		{
			name:    "Valid URL",
			args:    args{URL: "https://www.example.com"},
			isValid: true,
		},
		{
			name:    "Valid URL with no scheme",
			args:    args{URL: "www.example.com"},
			isValid: false,
		},
		{
			name:    "Invalid URL with no scheme",
			args:    args{URL: "example.com"},
			isValid: false,
		},
		{
			name:    "Invalid URL",
			args:    args{URL: "invalid"},
			isValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidURL(tt.args.URL); got != tt.isValid {
				t.Errorf("Expected isValidURL(%q) to return %v, but got: %v", tt.args.URL, got, tt.isValid)
			}
		})
	}
}

func Test_isInternalURL(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name       string
		args       args
		isInternal bool
	}{
		{
			name:       "Internal URL",
			args:       args{URL: "https://kdlp.underground.software/index.html"},
			isInternal: true,
		},

		{
			name:       "External URL 1",
			args:       args{URL: "https://www.google.com/"},
			isInternal: false,
		},

		{
			name:       "External URL 2",
			args:       args{URL: "https://bssw.io/items/the-developer-certificate-of-origin"},
			isInternal: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInternalURL(tt.args.URL); got != tt.isInternal {
				t.Errorf("Expected isInternalURL(%q) to return %v, but got: %v", tt.args.URL, got, tt.isInternal)
			}
		})
	}
}
