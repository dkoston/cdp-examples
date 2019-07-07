# cdp-examples

[![CircleCI](https://circleci.com/gh/dkoston/cdp-examples.svg?style=svg)](https://circleci.com/gh/dkoston/cdp-examples)

Some real world examples of using github.com/mafredri/cdp for browser automation


## Development

## Pre-requisites

- go 1.12+

1. Install dependencies with `go get -u`
2. Run [./bin/git/init-hooks.bash](./bin/git/init-hooks.bash) to enable hooks for gofmt and golint 

## Testing

`go test -v -race ./...`


## Examples

### Get Cookies

Get the cookies from Chrome for a specific domain:

`go run cmd/examples/main.go getcookies -d=domain`

Get all cookies from Chrome:

`go run cmd/examples/main.go getcookies`

Related source: [./src/cookies](./src/cookies)