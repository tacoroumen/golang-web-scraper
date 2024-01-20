# Website Scraper

This is a simple command-line tool written in Go for scraping a website and downloading its HTML content along with associated assets (images, scripts, stylesheets, etc.). The tool uses Go's standard library and the "golang.org/x/net/html" package for parsing HTML.

## Usage

```bash
go run main.go -url <website_url>
```

Replace `<website_url>` with the URL of the website you want to scrape.

## Features

- Downloads the HTML content of the specified website.
- Downloads images, scripts, stylesheets, and other assets linked in the HTML.
- Parses CSS styles and extracts background images.

## Dependencies

The tool uses only the standard library and the "golang.org/x/net/html" package for HTML parsing.

## How It Works

1. **Download HTML Content**: The tool fetches the HTML content of the specified website using an HTTP GET request.

2. **Save HTML Content**: The HTML content is saved to a local file named "index.html" in a folder named "scraped," which is created in the current working directory.

3. **Parse HTML for Assets**: The HTML content is parsed to identify assets such as images, scripts, stylesheets, links, and SVGs.

4. **Download Assets**: The identified assets are downloaded to the "scraped" folder. The tool recognizes and handles various HTML elements like `img`, `script`, `link`, `a`, and `svg`.

5. **Parse CSS Styles**: The tool extracts background images from CSS styles within the HTML content.

6. **Download Background Images**: If background images are found in the CSS styles, they are downloaded and saved to the "scraped" folder.

## Note

- The CSS parsing currently uses a simple regular expression approach for demonstration purposes. For a more robust solution, consider using a dedicated CSS parsing library.
- This tool may not handle all edge cases, and adjustments might be needed for specific websites.

## License

This tool is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.