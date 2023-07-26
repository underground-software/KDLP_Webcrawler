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
