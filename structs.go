package main

type Crawler struct {

	// Map to track visited URLs
	visited map[string]bool

	// Slice to store dead links
	deadLinks []string
}
