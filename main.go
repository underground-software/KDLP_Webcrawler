package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Println("Usage:")
	fmt.Println("\t$ ./KDLP_Webcrawler.git <option>") // TODO make easier

	fmt.Println("\nOptions:")
	fmt.Println("\t-h, --help             shows a manual")
	fmt.Println("\t--crawl                recursively crawls domain and retrieves dead links")

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
