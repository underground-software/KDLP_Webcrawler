/*
 * This code includes portions of Colly (https://github.com/gocolly/colly)
 * Copyright (c) 2017 Adam Tauber (asciimoo) <https://github.com/asciimoo>
 * Licensed under the Apache License, Version 2.0 (http://www.apache.org/licenses/LICENSE-2.0)
 */

package main

import (
	"fmt"
	"os"
)

func help() {

	fmt.Println("Usage:")
	fmt.Println("\t$ ./Webcrawler <option>")

	fmt.Println("\nOptions:")
	fmt.Println("\t-h, --help             shows a manual")
	fmt.Println("\t--crawl                recursively crawls domain and retrieves dead links with reference URLS")
	fmt.Println("\t--crawl-colly          recursively crawls domain and retrieves dead links with reference URLS via colly")

}

func main() {

	initializeErrorLogging()

	switch os.Args[1] {

	case "-h":
		fallthrough

	case "--help":
		help()

	case "--crawl":
		// Set domain and starting URL for crawling
		domain := "https://prod-01.kdlp.underground.software/"
		homeURL := domain + "index.md"

		// Public version
		// domain := "https://kdlp.underground.software/"
		// homeURL := domain + "index.html"

		// Call the crawl process
		runCustomCrawl(domain, homeURL)

	case "--crawl-colly":
		baseURL := "prod-01.kdlp.underground.software"
		StartCollyCrawl(baseURL)

	default:

		fmt.Println(os.Args[1] + " is not a valid argument.\nRunning " + os.Args[0] + " --help may help you!")

	}
}
