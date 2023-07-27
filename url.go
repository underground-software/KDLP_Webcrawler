package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Checks whether a given URL is valid
func isValidURL(URL string) bool {
	_, err := url.ParseRequestURI(URL)
	return err == nil
}

// Checks whether a URL is internal
func isInternalURL(URL string, domain string) bool {
	return strings.HasPrefix(URL, domain)
}

// Function to fetch HTTP response
func fetchHTTPResponse(URL string) (*http.Response, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Function to read HTTP response body
func readHTTPResponseBody(resp *http.Response) (string, error) {
	content, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// Function to check URL's HTTP status code
func checkURLStatus(URL string) (int, error) {
	resp, err := fetchHTTPResponse(URL)
	if err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

// Function to retrieve the HTTP content of a URL page
func retrieveHTTPContent(URL string) (string, error) {
	statusCode, err := checkURLStatus(URL)
	if err != nil {
		return "", err
	}

	if statusCode < 200 || statusCode >= 300 {
		return "", fmt.Errorf("failed to fetch content: received status code %d", statusCode)
	}

	resp, err := fetchHTTPResponse(URL)
	if err != nil {
		return "", err
	}

	content, err := readHTTPResponseBody(resp)
	if err != nil {
		return "", err
	}

	return content, nil
}
