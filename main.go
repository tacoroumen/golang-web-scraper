package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"regexp"

	"golang.org/x/net/html"
)

func downloadPage(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		body, err := io.ReadAll(response.Body)
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
	err := os.WriteFile(savePath, []byte(htmlContent), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("Downloaded: %s\n", savePath)
	return nil
}

func downloadAssets(doc *html.Node, baseURL, folder string) {
	var downloadAsset func(*html.Node)

	downloadAsset = func(n *html.Node) {
		if n.Type == html.ElementNode && (n.Data == "img" || n.Data == "script" || n.Data == "link" || n.Data == "a" || n.Data == "svg") {
			for _, attr := range n.Attr {
				if (attr.Key == "href" || attr.Key == "src") && (n.Data != "link" || n.Data != "a" || n.Data != "script") {
					url := urlJoin(baseURL, attr.Val)
					if strings.HasPrefix(url, baseURL) {
						downloadFile(url, folder)
					}
				}
			}
		} else if n.Type == html.ElementNode && n.Data == "style" {
			// Parse and download background images from CSS styles
			parseAndDownloadBackgroundImages(n.FirstChild, baseURL, folder)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			downloadAsset(c)
		}
	}

	downloadAsset(doc)
}

func parseAndDownloadBackgroundImages(styleNode *html.Node, baseURL, folder string) {
	if styleNode != nil && styleNode.Type == html.TextNode {
		// Extract CSS styles from the <style> node
		cssStyles := styleNode.Data

		// TODO: Implement a CSS parser to extract URLs from background-image properties
		// For simplicity, you may use regular expressions or a third-party CSS parser library

		// Example using regular expression (not recommended for all cases):
		// background-image: url('example.jpg');
		re := regexp.MustCompile(`background-image:[^url]*url\(['"]?([^'"\)]+)['"]?\)`)
		matches := re.FindAllStringSubmatch(cssStyles, -1)

		for _, match := range matches {
			if len(match) >= 2 {
				url := urlJoin(baseURL, match[1])
				if strings.HasPrefix(url, baseURL) {
					downloadFile(url, folder)
				}
			}
		}
	}
}

func downloadFile(url, folder string) {
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
