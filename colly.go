/*
 * This code includes portions of Colly (https://github.com/gocolly/colly)
 * Copyright (c) 2017 Adam Tauber (asciimoo) <https://github.com/asciimoo>
 * Licensed under the Apache License, Version 2.0 (http://www.apache.org/licenses/LICENSE-2.0)
 */

package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Function to check if the URL is relative and contains "/cgit"
func isRelativeURLWithCgit(href string) bool {
	return strings.Contains(href, "/cgit")
}

// Function to handle the dead link
func handleDeadLink(deadURL string, statusCode int, deadLinks *[]string) {
	// Log the dead link with the referring URL and status code
	log.Println("Dead Link found:", deadURL, "Status Code:", statusCode)

	// Append the dead link along with the referring URL to the deadLinks slice
	*deadLinks = append(*deadLinks, "Dead Link Found: "+deadURL)

	// Save the updated deadLinks slice to the dead links file
	if err := saveDeadLinksToFile("dead_links.txt", *deadLinks); err != nil {
		log.Println("Error saving dead links to file:", err)
	}
}

// Function to start crawling with colly
func StartCollyCrawl(baseURL string) {

	// Start the timer
	startTime := time.Now()

	url := []string{baseURL}

	startingURL := "https://" + baseURL

	// Declare a slice to store the dead links
	var deadLinks []string

	c := colly.NewCollector(
		colly.AllowedDomains(url...),
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {

		// Call handleDeadLink when the response status code indicates an error
		if r != nil && r.StatusCode >= 400 {
			handleDeadLink(r.Request.URL.String(), r.StatusCode, &deadLinks)
		}
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		if isRelativeURLWithCgit(href) {
			return
		}

		// Otherwise, visit the URL
		e.Request.Visit(href)
	})

	fmt.Println("Starting crawl at:", baseURL)

	if err := c.Visit(startingURL); err != nil {
		fmt.Println("Error on start of crawl:", err)
	}

	c.Wait()

	// Stop the timer
	elapsedTime := time.Since(startTime)

	// Display the elapsed time
	fmt.Println("Elapsed time:", elapsedTime)

	// Check if there are any dead links
	if len(deadLinks) > 0 {

		// Display file path for dead links
		fmt.Println("Dead links written to: dead_links.txt")

	} else {

		// Display no dead links found in terminal, no dead links file is created
		fmt.Println("No dead links found")

	}
}
