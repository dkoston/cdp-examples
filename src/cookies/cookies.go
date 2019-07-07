package cookies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mafredri/cdp"
)

// GetAllCookies returns all cookies from Chrome
func GetAllCookies(ctx context.Context, network cdp.Network) (cookies []*http.Cookie, err error) {
	networkCookies, err := network.GetAllCookies(ctx)
	if err != nil {
		return cookies, err
	}

	for _, cookie := range networkCookies.Cookies {
		c := http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,

			Domain: cookie.Domain,
			Path:   cookie.Path,
			// Expires
			// RawExpires

			// MaxAge
			Secure:   cookie.Secure,
			HttpOnly: cookie.HTTPOnly,
			// Do we need SameSite here?
			// Raw
			// Unparsed
		}
		cookies = append(cookies, &c)
	}

	return cookies, nil
}

// GetCookiesForDomain returns the cookies from Chrome for a specific domain name
func GetCookiesForDomain(ctx context.Context, network cdp.Network, domain string) (cookies []*http.Cookie, err error) {
	allCookies, err := GetAllCookies(ctx, network)
	if err != nil {
		return cookies, fmt.Errorf("failed getting cookies from browser: %v", err)
	}

	for i := 0; i < len(allCookies); i++ {
		if allCookies[i].Domain == domain {
			cookies = append(cookies, allCookies[i])
		}
	}

	return cookies, nil
}
