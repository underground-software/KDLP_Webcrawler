package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Function to create a new instance of the web crawler
func newCrawler(domain, homeURL string) *Crawler {
	return &Crawler{
		domain:    domain,
		homeURL:   homeURL,
		visited:   make(map[string]bool),
		deadLinks: []string{},
	}
}

func findLinks(n *html.Node, baseURL string) []string {
	// Initialize a slice to store the final unique valid links in order of discovery
	var links []string

	// If node type is element node (type 2) and node data contains attribute "a" (anchor tag)
	if n.Type == html.ElementNode && n.Data == "a" {

		// For each attribute in the "a" tag
		for _, attr := range n.Attr {

			// Resolve the URL to handle relative URLs correctly
			if attr.Key == "href" {
				absoluteURL, err := resolveURL(baseURL, attr.Val)
				if err != nil {
					log.Println("Error resolving URL:", err)
					continue
				}

				// Check if the resolved URL is valid and store it in the links slice
				if isValidURL(absoluteURL) {
					links = append(links, absoluteURL)
				} else {
					log.Println("Invalid URL found:", absoluteURL)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		// Recursively call findLinks only for the child nodes
		links = append(links, findLinks(c, baseURL)...)
	}

	return links
}

// Function to extract valid links from HTML content
func extractValidLinks(content string, baseURL string) []string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return nil // Return nil in case of parsing error
	}

	// Start the recursive link extraction from the root node of the HTML tree
	return findLinks(doc, baseURL)
}

func saveDeadLinksToFile(filepath string, deadLinks []string) error {

	// Open the file in write-only mode
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each dead link to the file, one link per line
	for _, link := range deadLinks {
		if _, err := fmt.Fprintln(file, link); err != nil {
			return err
		}
	}

	return nil
}

// HandleDeadLink handles the case when the URL is a dead link.
func (c *Crawler) handleDeadLink(referringURL, deadURL string, statusCode int) {

	// Log the dead link with the referring URL and status code
	log.Println("Dead Link:", deadURL, "found on:", referringURL, "Status Code:", statusCode)

	// Append the dead link along with the referring URL to the deadLinks slice
	c.deadLinks = append(c.deadLinks, "dead link "+deadURL+" found at: "+referringURL)

	// Save the updated deadLinks slice to the dead links file
	if err := saveDeadLinksToFile("dead_links.txt", c.deadLinks); err != nil {
		log.Println("Error saving dead links to file:", err)
	}
}

// Initiates crawl process for a URL
func (c *Crawler) crawlURL(URL, referenceURL string) {
	// Check if the URL has already been visited
	if c.visited[URL] {
		fmt.Println("Already visited:", URL)
		return
	}

	// Mark the URL as visited
	c.visited[URL] = true
	fmt.Println("Added", URL, "to visited map")

	// Check if the URL is valid
	if !isValidURL(URL) {
		log.Println("Invalid URL:", URL)
		return
	}

	// Check if the URL should be treated as a fake URL and skip processing it
	if isFakeURL(URL) {
		fmt.Println("Fake URL:", URL)
		return
	}

	// Fetch the status of the URL
	statusCode, err := checkURLStatus(URL)
	if err != nil {
		log.Println("Error checking status for URL:", URL, "Error:", err)
		return
	}

	if statusCode == 404 {
		c.handleDeadLink(referenceURL, URL, statusCode)
		return
	}

	// If internal link: Fetch content, extract URLs, and crawl URLs
	if isInternalURL(URL, c.domain) {
		c.crawlInternalURL(URL, referenceURL)
	}
}

// fetches content, extracts URLs, and crawls URLs for internal links
func (c *Crawler) crawlInternalURL(URL, referringURL string) {
	// Fetch the content of the URL
	content, err := retrieveHTTPContent(URL)
	if err != nil {
		log.Println("Error fetching contents of URL:", URL, "Error:", err)
		return
	}

	// Parse HTML content and extract links
	links := extractValidLinks(content, URL)

	// Iterate through the links and only crawl unvisited internal links
	for _, link := range links {
		if !c.visited[link] {
			c.crawlURL(link, URL)
		}
	}
}
