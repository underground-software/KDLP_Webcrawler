package main

import (
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
func isInternalURL(URL string) bool {
	return strings.HasPrefix(URL, domain)
}

// Function to fetch URL's HTTP status code
func checkURLStatus(URL string) (int, error) {

	// Fetch HTTP contents of URL
	resp, err := http.Get(URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Return status code
	return resp.StatusCode, nil

}
