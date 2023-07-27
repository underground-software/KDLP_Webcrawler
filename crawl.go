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

// Function to extract links from HTML content
func extractLinks(content string, baseURL string) []string {

	// Initialise links array
	var links []string

	// Parse HTML content
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return links // Return empty links array in case of parsing error
	}

	// Recursive function to extract links from HTML tree
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {

		// If node type is element node (type 2) and node data contains attribute "a" (anchor tag)
		if n.Type == html.ElementNode && n.Data == "a" {

			// For each attribute in the "a" tag
			for _, attr := range n.Attr {

				// If the attribute key is "href"
				if attr.Key == "href" {

					// Check if the URL is relative
					if !strings.HasPrefix(attr.Val, "http://") && !strings.HasPrefix(attr.Val, "https://") {

						// Convert the relative URL to absolute URL by appending the base URL
						attr.Val = baseURL + attr.Val
					}

					// Append valid URL to links array
					links = append(links, attr.Val)
				}
			}
		}

		// Recursively call findLinks on each child node of the current node
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findLinks(c)
		}
	}

	// Start the recursive link extraction from the root node of the HTML tree
	findLinks(doc)

	// Return the extracted links array
	return links
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

// crawlInternalURL fetches content, extracts URLs, and crawls URLs for internal links
func (c *Crawler) crawlInternalURL(URL, referringURL string) {
	// Fetch the content of the URL
	content, err := retrieveHTTPContent(URL)
	if err != nil {
		log.Println("Error fetching contents of URL:", URL, "Error:", err)
		return
	}

	// Parse HTML content and extract links
	links := extractLinks(content, c.domain)
	fmt.Println("Links found in", URL, ":", links)

	// Recursively call crawlURL for each internal link found
	for _, link := range links {
		c.crawlURL(link, URL)
	}
}
