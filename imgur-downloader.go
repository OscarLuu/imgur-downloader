package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/html"
)

func createFile(imageURL string, image string) {
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(image + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("File size: ", b)
}

func getVal(tag html.Token) {
	// Finds that id value of the imgur links
	for _, div := range tag.Attr {
		if div.Key == "id" {
			if len(div.Val) == 7 {
				image := div.Val
				encodeURL(image)
			}
		}
	}
}

func encodeURL(image string) {
	imageURL := "https://i.imgur.com/" + image + ".jpg"
	createFile(imageURL, image)
}

func tokenizer(resp *http.Response) {
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// Error Token is end of the document
			return
		case tt == html.StartTagToken:
			tag := z.Token()

			// Check to see if the tag has <div> if not move to the next line.
			div := tag.Data == "div"
			if !div {
				continue
			}
			getVal(tag)
		}
	}
}

func main() {
	// Create directory
	currentTime := time.Now()
	dir := os.Args[2] + currentTime.Format("01-02-2006")

	os.Mkdir(dir, 0700)
	os.Chdir(dir)

	// Fetch URL
	url := os.Args[1]
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	tokenizer(resp)
}
