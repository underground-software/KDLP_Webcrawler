# KDLP_Webcrawler

Golang application for identifying dead links on KDLP website

## Setup

1. Go to the project directory and compile the source code.

   ```bash
   go build -o webcrawler  
   ```

## Usage

- Usage:

    ```bash
    ./webcrawler <option>
    ```

- Options:

   ```bash
   -h, --help       shows a manual
   --crawl          recursively crawls domain and retrieves dead links with reference URLS
   --crawl-colly    crawls domain for dead links via colly
   ```
