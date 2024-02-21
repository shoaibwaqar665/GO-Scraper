package main

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	// "io"
	"log"
	// "net/http"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type scrapStruct struct {
	url   string
	image string
	title string
	text  string
}

func scrapeAndWriteCSV() []scrapStruct {
	var scrapData []scrapStruct
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	)

	c.SetRequestTimeout(time.Second * 10)

	c.OnHTML(".o-listease__item", func(e *colly.HTMLElement) {
		linkData := scrapStruct{}

		linkData.url = e.ChildAttr("a", "href")
		linkData.image = e.ChildAttr("img", "src")
		linkData.title = e.ChildAttr("a", "title")
		linkData.text = e.ChildText(".m-statement__quote")

		scrapData = append(scrapData, linkData)
	})

	c.OnHTML(".m-teaser", func(e *colly.HTMLElement) {
		linkData := scrapStruct{}

		linkData.url = e.ChildAttr("a", "href")
		linkData.image = e.ChildAttr("img", "src")
		linkData.title = e.ChildAttr("a", "title")

		scrapData = append(scrapData, linkData)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\n", r.Request.URL, err)
	})

	c.Visit("https://www.politifact.com")
	c.Wait()

	file, err := os.Create("link1.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	headers := []string{
		"url",
		"image",
		"title",
		"text",
	}
	writer.Write(headers)

	for _, dataArray := range scrapData {
		record := []string{
			dataArray.url,
			dataArray.image,
			dataArray.title,
			dataArray.text,
		}

		writer.Write(record)
	}

	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatalln("Error writing CSV:", err)
	}

	downloadImages(scrapData) // Call function to download images

	return scrapData
}

// Function to download images
func downloadImages(data []scrapStruct) {
	imageCollector := colly.NewCollector()

	imageCollector.OnResponse(func(r *colly.Response) {
		fileName := createFileName(r.Request.URL.String())
		filePath := fmt.Sprintf("imagesLink1/%s", fileName) // Define your image storage directory
		err := os.WriteFile(filePath, r.Body, 0644)
		if err != nil {
			log.Printf("Failed to save image %s: %v\n", filePath, err)
		} else {
			log.Printf("Image saved as %s\n", filePath)
		}
	})

	for _, item := range data {
		if item.image != "" {
			imageCollector.Visit(item.image)
		}
	}
}

// Function to create a unique file name for each image
func createFileName(url string) string {
	hasher := md5.New()
	hasher.Write([]byte(url))
	return hex.EncodeToString(hasher.Sum(nil)) + ".jpg"
}