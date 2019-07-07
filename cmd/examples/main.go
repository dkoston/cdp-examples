package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dkoston/cdp-examples/src/browser"
	"github.com/dkoston/cdp-examples/src/cookies"
	"github.com/dkoston/cdp-examples/src/logger"
	"github.com/jessevdk/go-flags"
)

const version = "0.0.1"

// CommandLineOptions allows for passing commands to main.go
type CommandLineOptions struct {
	Version VersionCommand `command:"version"  description:"print version and exit"`
	Cookies CookiesCommand `command:"getcookies" description:"dumps out all cookies or a domain's cookies from Chrome'"`
}

var errShowVersionString = errors.New("version command invoked")

// VersionCommand is used to specify to the application to print its version
type VersionCommand struct {
	Version string `default:"version"`
}

// Execute prints when the version command is used on the command line
func (command *VersionCommand) Execute(args []string) error {
	return errShowVersionString
}

// CookiesCommand allows passing options to the "getcookies" command
type CookiesCommand struct {
	Domain string `short:"d" long:"domain" description:"limit cookies to a specific domain name" optional:"true" env:"COOKIE_DOMAIN"`
}

func main() {
	var opts CommandLineOptions
	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	handleCliParseError(parser, err)

	log := logger.Get()

	// Launch a chrome instance with debug port open
	go browser.LaunchChrome()
	time.Sleep(1 * time.Second)

	const timeout = 5 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	domain := "https://www.google.com"

	conn, c, err := browser.ConnectCDP(ctx, domain)

	var browserCookies []*http.Cookie

	if opts.Cookies.Domain != "" {
		browserCookies, err = cookies.GetCookiesForDomain(ctx, c.Network, opts.Cookies.Domain)
		if err != nil {
			log.Errorf("failed to get cookies from %s: %v", domain, err)
		}
	} else {
		browserCookies, err = cookies.GetAllCookies(ctx, c.Network)
		if err != nil {
			log.Errorf("failed to get all browser cookies: %v", err)
		}
	}

	for i := 0; i < len(browserCookies); i++ {
		log.Noticef("%v", browserCookies[i])
	}

	err = conn.Close()
	if err != nil {
		log.Errorf("failed to close connection to Chrome remote debugging: %v", err)
	}
}

// handleCliParseError handles errors resulting from the parsing of our command line options
// by either printing out our version string, or printing the errors
func handleCliParseError(parser *flags.Parser, err error) {
	if err != nil {
		if err == errShowVersionString {
			logger.Get().Infof("[VERSION] %s", version)
			os.Exit(0)
		} else {
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}
}
