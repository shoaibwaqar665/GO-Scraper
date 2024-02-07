package main

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type link3 struct {
	url   string
	image string
	title string
	text  string
}

func scrapLink3() []link3 {
	var scrapData []link3

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	)

	// Set a delay between requests to avoid being blocked
	c.SetRequestTimeout(time.Second * 10)

	c.OnHTML(".o-listease__item", func(e *colly.HTMLElement) {
		linkData := link3{}

		linkData.url = e.ChildAttr("a", "href")
		linkData.image = e.ChildAttr("img", "src")
		linkData.title = e.ChildAttr("a", "title")
		linkData.text = e.ChildText(".m-statement__quote")

		scrapData = append(scrapData, linkData)
	})

	c.OnHTML(".m-teaser", func(e *colly.HTMLElement) {
		linkData := link3{}

		linkData.url = e.ChildAttr("a", "href")
		linkData.image = e.ChildAttr("img", "src")
		linkData.title = e.ChildAttr("a", "title")

		scrapData = append(scrapData, linkData)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\n", r.Request.URL, err)
	})

	c.Visit("https://www.politifact.com")

	// Wait for the collector to finish
	c.Wait()

	file, err := os.Create("link1.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	// Initializing a file writer
	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"image",
		"title",
		"text",
	}
	writer.Write(headers)

	for _, dataArray := range scrapData {
		// Converting a data to an array of strings
		record := []string{
			dataArray.url,
			dataArray.image,
			dataArray.title,
			dataArray.text,
		}

		writer.Write(record)
	}

	writer.Flush()

	// Check for any errors in writing the CSV file
	if err := writer.Error(); err != nil {
		log.Fatalln("Error writing CSV:", err)
	}

	return scrapData
}
