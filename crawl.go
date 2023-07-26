package main

import (
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Function to create new instance of web crawler
func newCrawler() *Crawler {
	return &Crawler{
		visited: make(map[string]bool),
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

// Function to write dead links to a file
func (c *Crawler) writeDeadLinksToFile(filepath string) error {
	// Open the file in append mode
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write dead links to the file
	for _, link := range c.deadLinks {
		if _, err := file.WriteString(link + "\n"); err != nil {
			return err
		}
	}

	return nil
}
