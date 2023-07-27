package main

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

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
		URL    string
		domain string
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
			if got := isInternalURL(tt.args.URL, tt.args.domain); got != tt.isInternal {
				t.Errorf("Expected isInternalURL(%q) to return %v, but got: %v", tt.args.URL, got, tt.isInternal)
			}
		})
	}
}

func Test_checkURLStatus(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		status  int
		wantErr bool
	}{
		{
			name:   "HTTP Status Code 200 - OK",
			args:   args{URL: "http://httpstat.us/200"},
			status: 200,
		},

		{
			name:   "HTTP Status Code 301 - Permanent Redirect",
			args:   args{URL: "http://httpstat.us/301"},
			status: 200, // Because it redirects to OK page
		},

		{
			name:   "HTTP Status Code 302 - Temporary Redirect",
			args:   args{URL: "http://httpstat.us/302"},
			status: 200, // Because it redirects to OK page
		},
		{
			name:   "HTTP Status Code 404 - Not Found",
			args:   args{URL: "http://httpstat.us/404"},
			status: 404,
		},
		{
			name:   "HTTP Status Code 410 - Gone",
			args:   args{URL: "http://httpstat.us/410"},
			status: 410,
		},
		{
			name:   "HTTP Status Code 500 - Internal Sever Error",
			args:   args{URL: "http://httpstat.us/500"},
			status: 500,
		},
		{
			name:   "HTTP Status Code 503 - Service Unavailable",
			args:   args{URL: "http://httpstat.us/503"},
			status: 503,
		},
		{
			name:    "Invalid",
			args:    args{URL: "Invalid"},
			status:  0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkURLStatus(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkURLStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.status {
				t.Errorf("checkURLStatus() = %v, want %v", got, tt.status)
			}
		})
	}
}

func Test_retrieveHTTPContent(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		content string
		wantErr bool
	}{
		{
			name:    "Succesful fetch",
			args:    args{URL: "https://www.google.com"},
			content: "not empty", // Expecting a non-empty response
			wantErr: false,
		},
		{
			name:    "Unsuccesful fetch",
			args:    args{URL: "Invalid"},
			content: "", // Expecting an empty response
			wantErr: true,
		},

		{
			name:    "Unsuccessful fetch - Non-existent resource",
			args:    args{URL: "https://www.example.com/non-existent-page"},
			content: "", // Expecting an empty response
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := retrieveHTTPContent(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchHTTPContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.content != "" && got == "" {
				t.Errorf("fetchHTTPContent() returned an empty response, but a non-empty response was expected")
			}
			if tt.content == "" && got != "" {
				t.Errorf("fetchHTTPContent() returned a non-empty response, but an empty response was expected")
			}
		})
	}
}

func Test_readHTTPResponseBody(t *testing.T) {
	type args struct {
		resp *http.Response
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Read non-empty response body",
			args: args{
				resp: &http.Response{
					Body: io.NopCloser(strings.NewReader("I'm not empty!")),
				},
			},
			want: "I'm not empty!",
		},
		{
			name: "Read empty response body",
			args: args{
				resp: &http.Response{
					Body: io.NopCloser(strings.NewReader("")),
				},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readHTTPResponseBody(tt.args.resp)
			if (err != nil) != tt.wantErr {
				t.Errorf("readHTTPResponseBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readHTTPResponseBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_fetchHTTPResponse(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "Successful fetch response",
			args: args{
				URL: "https://www.google.com",
			},
			wantErr: false,
		},
		{
			name: "Unsuccessful fetch response",
			args: args{
				URL: "Invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchHTTPResponse(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchHTTPResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			_ = got
		})
	}
}
