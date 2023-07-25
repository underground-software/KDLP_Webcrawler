package main

// Function to create new instance of web crawler
func newCrawler() *Crawler {
	return &Crawler{
		visited: make(map[string]bool),
	}
}
