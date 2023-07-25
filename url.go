package main

import (
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
