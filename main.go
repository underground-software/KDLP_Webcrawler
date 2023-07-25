package main

import (
	"fmt"
	"os"
)

func help() {
	// TODO: Implement
	fmt.Println("Help is on the way")

}

func main() {

	switch os.Args[1] {

	case "-h":
		fallthrough

	case "--help":
		help()

		// TODO: implement
	case "--crawl":
		// Create a new instance of crawler
		crawler := newCrawler()

		// Call the crawlURL function on the KDLP home page
		crawler.crawlURL(homeURL)

	default:
		fmt.Println(os.Args[1] + " is not a valid argument.\nRunning " + os.Args[0] + " --help may help you!")

	}
}
