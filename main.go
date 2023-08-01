package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func help() {

	fmt.Println("Usage:")
	fmt.Println("\t$ ./Webcrawler <option>")

	fmt.Println("\nOptions:")
	fmt.Println("\t-h, --help             shows a manual")
	fmt.Println("\t--crawl                recursively crawls domain and retrieves dead links with reference URLS")

}

func main() {

	// Get the current working directory
	currentDir, err := getCurrentDirectory()
	if err != nil {
		log.Fatal("Failed to get current working directory:", err)
	}

	// Construct the path for the error log file in the current directory
	logFilePath := filepath.Join(currentDir, "error_log.txt")

	// Sets logged errors to print to both error log file and terminal
	_, err = openErrorLogFile(logFilePath)
	if err != nil {
		log.Fatal("Failed to open error log file:", err)
	}

	switch os.Args[1] {

	case "-h":
		fallthrough

	case "--help":
		help()

	case "--crawl":
		// Set domain and starting URL for crawling
		domain := "https://prod-01.kdlp.underground.software/"
		homeURL := domain + "index.md"

		// // Public version
		// domain := "https://kdlp.underground.software/"
		// homeURL := domain + "index.html"

		// Create a new instance of crawler
		crawler := newCrawler(domain, homeURL)

		// Call the crawlURL function on the KDLP home page
		crawler.crawlURL(homeURL, "") // Empty string for storing URLs

		// Check if there are any dead links
		if len(crawler.deadLinks) > 0 {

			// Display file path for dead links
			fmt.Println("Dead links written to: dead_links.txt")

		} else {

			// Display no dead links found in terminal, no dead links file is created
			fmt.Println("No dead links found")

		}

	default:

		fmt.Println(os.Args[1] + " is not a valid argument.\nRunning " + os.Args[0] + " --help may help you!")

	}
}
