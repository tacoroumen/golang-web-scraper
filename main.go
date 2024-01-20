package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func downloadPage(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(body), nil
	} else {
		return "", fmt.Errorf("Failed to download: %s", url)
	}
}

func saveHTMLContent(htmlContent, folder string) error {
	savePath := filepath.Join(folder, "index.html")
	err := ioutil.WriteFile(savePath, []byte(htmlContent), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Downloaded: %s\n", savePath)
	return nil
}

func downloadAssets(doc *html.Node, baseURL, folder string) {
	var downloadAsset func(*html.Node)

	downloadAsset = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "img" || n.Data == "script" || n.Data == "link" || n.Data == "a") {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					url := urlJoin(baseURL, attr.Val)
					if strings.HasPrefix(url, baseURL) {
						response, err := http.Get(url)
						if err == nil && response.StatusCode == 200 {
							assetPath := filepath.Join(folder, strings.TrimLeft(urlParse(url).Path, "/"))
							os.MkdirAll(filepath.Dir(assetPath), os.ModePerm)
							file, err := os.Create(assetPath)
							if err == nil {
								defer file.Close()
								_, err = io.Copy(file, response.Body)
								if err == nil {
									fmt.Printf("Downloaded asset: %s\n", url)
								}
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			downloadAsset(c)
		}
	}

	downloadAsset(doc)
}

func scrapeWebsite(websiteURL string) {
	outputFolder := filepath.Join("scraped", strings.Replace(urlParse(websiteURL).Host, "www.", "", 1))
	os.MkdirAll(outputFolder, os.ModePerm)

	htmlContent, err := downloadPage(websiteURL)
	if err == nil {
		err = saveHTMLContent(htmlContent, outputFolder)
		if err == nil {
			doc, err := html.Parse(strings.NewReader(htmlContent))
			if err == nil {
				downloadAssets(doc, websiteURL, outputFolder)
			}
		}
	}
	if err != nil {
		fmt.Println(err)
	}
}

func urlJoin(base, relative string) string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return relative
	}
	relativeURL, err := url.Parse(relative)
	if err != nil {
		return relative
	}
	return baseURL.ResolveReference(relativeURL).String()
}

func urlParse(input string) *url.URL {
	parsedURL, _ := url.Parse(input)
	return parsedURL
}

func main() {
	var websiteURL string
	flag.StringVar(&websiteURL, "url", "", "URL of the website to scrape")
	flag.Parse()

	if websiteURL == "" {
		fmt.Println("Please provide the URL of the website to scrape using the -url flag.")
		return
	}

	scrapeWebsite(websiteURL)
}
