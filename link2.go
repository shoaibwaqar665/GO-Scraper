package main

import (
	"encoding/csv"
	// "fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type link2 struct {
	url   string
	image string
	video string
	text  string
}

func scrapLink2() []link2 {
	var scrapData []link2
	c := colly.NewCollector(
		
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
	)

	c.SetRequestTimeout(time.Second * 10)

	c.OnHTML(".entry-header", func(e *colly.HTMLElement) {
		linkData := link2{}

		linkData.url = e.ChildAttr("a", "href")
		linkData.image = e.ChildAttr("img", "src")
		linkData.text = e.ChildText("h4")

		scrapData = append(scrapData, linkData)
	})

	c.OnHTML(".embed-youtube", func(e *colly.HTMLElement) {
		linkData := link2{}

		linkData.video = e.ChildAttr("iframe", "src")

		scrapData = append(scrapData, linkData)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\n", r.Request.URL, err)
	})

	c.Visit("https://www.altnews.in")
	c.Wait()

	// Download images and videos
	for i, data := range scrapData {
		if data.image != "" {
			scrapData[i].image = downloadFile(data.image, "imagesLink2")
		}
		if data.video != "" {
			scrapData[i].video = downloadFile(data.video, "videos")
		}
	}

	// Write data to CSV
	writeCSV(scrapData)

	return scrapData
}

func downloadFile(url, folder string) string {
	// Create the folder if it doesn't exist
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		if err := os.MkdirAll(folder, 0755); err != nil {
			log.Printf("Error creating directory: %v\n", err)
			return ""
		}
	}

	// Get the file name
	tokens := strings.Split(url, "/")
	fileName := filepath.Join(folder, tokens[len(tokens)-1])

	// Check if file already exists
	if _, err := os.Stat(fileName); err == nil {
		log.Printf("File already exists: %s\n", fileName)
		return fileName
	}

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error downloading file from %s: %v\n", url, err)
		return ""
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating file %s: %v\n", fileName, err)
		return ""
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Printf("Error writing data to %s: %v\n", fileName, err)
		return ""
	}

	return fileName
}

func writeCSV(scrapData []link2) {
	file, err := os.Create("link2.csv")
	if err != nil {
		log.Fatalln("Failed to create output CSV file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"url", "image", "video", "text"}
	if err := writer.Write(headers); err != nil {
		log.Fatalln("Error writing headers to CSV:", err)
	}

	for _, record := range scrapData {
		if err := writer.Write([]string{record.url, record.image, record.video, record.text}); err != nil {
			log.Fatalln("Error writing record to CSV:", err)
		}
	}
}
