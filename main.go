package main

import (
	"fmt"
)

func main() {
	scrapeAndWriteCSV()
	scrapLink2()
	scrapLink3()
	fmt.Println("Scraping and CSV writing completed successfully.")
}


//   go run main.go link1.go link2.go link3.go

