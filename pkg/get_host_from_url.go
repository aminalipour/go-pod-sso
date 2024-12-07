package pkg

import (
	"net"
	"net/url"
)

func GetHostFromURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host, _, err := net.SplitHostPort(parsedURL.Host)
	if err != nil {
		// If there is no port, parsedURL.Host contains only the hostname
		host = parsedURL.Host
	}
	return host
}
