# Simple Web Crawler

This is a simple web crawler that takes a URL and scrapes for every URL on the same domain in order to generate a
sitemap.

## Build

To build the program, you need to have the following installed

- Golang v1.19+

Then run the following commands in the project folder:

```shell
# Download dependencies
go mod vendor
# Build the program
go build -o web-crawler
# Run the program
./web-crawler https://gobyexample.com/
```
## Sitemap Format

The sitemap format is as follows:

- For every page visited the parent URL is written to `sitemap.txt` in the current working directory followed by every anchor URL within the same host domain found on the page
- Each group of URLs as above is delineated by an empty line
```text
[parent-url]
[child-url]
[empty-line]
[parent-url]
[child-url]
[child-url]
[empty-line]
...
```

