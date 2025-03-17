package blog

import (
	"strings"

	"golang.org/x/net/html"
)

func generateSlug(title string) string {
	return strings.ToLower(strings.Join(strings.Split(title, " "), "-"))
}

func extraImageSourcesFromHTML(text string, bucket string) ([]string, error) {
	var imageSources []string

	doc, err := html.Parse(strings.NewReader(text))
	if err != nil {
		return imageSources, err
	}
	var traverse func(n *html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			if len(n.Attr) == 0 {
				return
			}

			// loop through <img> attributes
			for _, a := range n.Attr {
				// evalue if element contains src attribute
				if a.Key == "src" {
					// check whitelisted bucket name
					if strings.Contains(a.Val, bucket) {
						imageSources = append(imageSources, a.Val)
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	return imageSources, nil
}

func extractKeyFromImageSource(source, bucket string) string {
	result := strings.Split(source, bucket)

	if len(result) != 2 {
		return ""
	}

	return result[1]
}
