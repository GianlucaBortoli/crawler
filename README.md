# crawler
Simple web crawler in Go

## Description
This is a simple website crawler. You can pass it a URL and it will find all the links it can
reach as long as they're in the same sub-domain.

A textual sitemap is printed in the end to show where the URL visit reached.

## Build
The project uses go modules and requires go `1.12`.

By default the binary will be placed in `./bin` and has the race detector enabled.
```bash
make build
```

## Run the crawler
After building the tool, you can run it as follows
```bash
./bin/crawler
```

The only requires parameter is the URL where it needs to start from.
For example
```bash
./bin/crawler -url http://example.com
```

Other parameters can be configu

## Run unit tests
This target runs all the unit tests with the race detector enabled.

It shows also the overall coverage
```bash
make test
```

## Design decisions
TODO

## Limitations
TODO

## Possible enhancements
TODO
