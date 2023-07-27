package main

type Crawler struct {

	// Stores base domain of website being crawled
	domain string

	// Starting URL for crawler
	homeURL string

	// Map to track visited URLs
	visited map[string]bool

	// Slice to store dead links
	deadLinks []string
}
