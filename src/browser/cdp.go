package browser

import (
	"context"

	"github.com/dkoston/cdp-examples/src/logger"
	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/protocol/page"
	"github.com/mafredri/cdp/rpcc"
)

// ConnectCDP provides a connection to dev tools via remote debugging port, and a cdp.Client instance
// The cdp.Client will cause Chrome to open the URL specified
func ConnectCDP(ctx context.Context, url string) (conn *rpcc.Conn, c *cdp.Client, err error) {
	// Use the DevTools HTTP/JSON API to manage targets (e.g. pages, webworkers).
	devt := devtool.New("http://127.0.0.1:9222")
	pt, err := devt.Get(ctx, devtool.Page)
	if err != nil {
		pt, err = devt.Create(ctx)
		if err != nil {
			return conn, c, err
		}
	}

	// Initiate a new RPC connection to the Chrome DevTools Protocol target.
	conn, err = rpcc.DialContext(ctx, pt.WebSocketDebuggerURL)
	if err != nil {
		return conn, c, err
	}
	c = cdp.NewClient(conn)

	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := c.Page.DOMContentEventFired(ctx)
	if err != nil {
		return conn, c, err
	}

	// Enable events on the Page domain, it's often preferable to create
	// event clients before enabling events so that we don't miss any.
	if err = c.Page.Enable(ctx); err != nil {
		return conn, c, err
	}

	// Enable events on the Runtime domain (JavaScript)
	if err = c.Runtime.Enable(ctx); err != nil {
		return conn, c, err
	}

	// Create the Navigate arguments with the optional Referrer field set.
	navArgs := page.NewNavigateArgs(url)
	nav, err := c.Page.Navigate(ctx, navArgs)
	if err != nil {
		return conn, c, err
	}

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		return conn, c, err
	}

	logger.Get().Debugf("Page loaded (%s) with frame ID: %s\n", url, nav.FrameID)

	err = domContent.Close()
	return conn, c, err
}

// Reload is used to reload the current page via browser automation
func Reload(ctx context.Context, p cdp.Page) (err error) {
	var reloadArgs page.ReloadArgs

	return p.Reload(ctx, &reloadArgs)
}
