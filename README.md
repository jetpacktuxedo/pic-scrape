# pic-scape

Scrape all posted media for a bluesky repo (user).

## Status
This is quite a bit jank right now, but it technically works. 
There is very little filtering at all, so it will also scrape the embeded media included every time the user has shared a link. 
It should add proper file extensions based on the type of content that was downloaded, but if it finds a type that doesn't match it will download the file with no extension and log the discovered mime type.
Failing to download or save a file isn't fatal, just logged.
There is no auth mechanism at all right now, so this will probably only work on public repos.

## Usage
`pic-scrape bsky.app`
`pic-scrape.exe jay.bsky.team`

## Building
Windows: `GOOS=windows GOARCH=amd64 go build -o pic-scrape.exe`
Linux: `GOOS=linux GOARCH=amd64 go build -o pic-scrape`
MacOS: `GOOS=darwin GOARCH=amd64 go build -o pic-scrape-macos-intel`
MacOS: `GOOS=darwin GOARCH=arm64 go build -o pic-scrape-macos-silicon`
