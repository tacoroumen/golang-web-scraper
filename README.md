# Web Scraper

This is a simple web scraper written in Go that downloads the HTML content of a given website and saves it along with its assets (images, scripts, links) to a local directory.

## Usage

### Prerequisites
Make sure you have Go installed on your machine.

### Installation
Clone this repository and navigate to the project directory:

```bash
git clone https://github.com/tacoroumen/golang-web-scraper
cd golang-web-scraper
```

Build the project:

```bash
go build
```

### Running the Scraper

To run the web scraper, use the following command:

```bash
./web-scraper -url <website-url>
```

Replace `<website-url>` with the URL of the website you want to scrape.

### Example

```bash
./web-scraper -url https://example.com
```

This will create a folder named `scraped` in the current directory. Inside the `scraped` folder, a subfolder with the name of the website's domain (without 'www.') will be created. The HTML content of the website will be saved as `index.html` in this folder, and the assets (images, scripts, links) will be downloaded and saved in their respective directories.

## Flags

- `-url`: Specifies the URL of the website to scrape.

## Dependencies

The following external packages are used in this project:

- `golang.org/x/net/html`: HTML parsing package.

## License

This web scraper is released under the [MIT License](LICENSE). Feel free to use and modify it according to your needs.