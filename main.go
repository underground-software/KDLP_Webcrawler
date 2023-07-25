package main

import (
	"fmt"
	"os"
)

func help() {

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

	default:
		fmt.Println(os.Args[1] + " is not a valid argument.\nRunning " + os.Args[0] + " --help may help you!")

	}
}
