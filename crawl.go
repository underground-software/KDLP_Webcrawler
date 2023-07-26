package main

import (
	"log"
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
func extractLinks(content string) []string {

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

				// If the attribute key is "href" and the URL/link is valid
				if attr.Key == "href" && isValidURL(attr.Val) {

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
