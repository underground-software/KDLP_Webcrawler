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

	case "--crawl":

	}
}
