package utils

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// ExtractLinks extracts all links from HTML content
func ExtractLinks(content string) []string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return nil
	}

	var links []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}

// ExtractTitle extracts the title from HTML content
func ExtractTitle(content string) string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return ""
	}

	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return title
}

// ToAbsoluteURL converts a relative URL to an absolute URL
func ToAbsoluteURL(base *url.URL, href string) string {
	if href == "" {
		return ""
	}

	// Remove fragments
	if i := strings.Index(href, "#"); i > -1 {
		href = href[:i]
	}

	// Skip certain schemes
	if strings.HasPrefix(href, "mailto:") || strings.HasPrefix(href, "tel:") || strings.HasPrefix(href, "javascript:") {
		return ""
	}

	relativeURL, err := url.Parse(href)
	if err != nil {
		return ""
	}

	absoluteURL := base.ResolveReference(relativeURL)
	return absoluteURL.String()
}
