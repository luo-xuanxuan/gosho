package gosho

import (
	"strings"

	"golang.org/x/net/html"
)

func parseHTMLNode(node *html.Node, tag string, id string, class string) []*html.Node {
	var nodes []*html.Node

	// Recursively search for matching nodes
	var searchNodes func(n *html.Node)
	searchNodes = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag {
			var hasID, hasClass bool
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == id {
					hasID = true
				}
				if attr.Key == "class" && attr.Val == class {
					hasClass = true
				}
			}
			// Add node if it matches both id and class (or if id/class are empty)
			if (id == "" || hasID) && (class == "" || hasClass) {
				nodes = append(nodes, n)
			}
		}
		// Recursively search the child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			searchNodes(c)
		}
	}

	// Start searching from the root of the document
	searchNodes(node)
	return nodes
}

func getAttr(n *html.Node, attribute string) string {
	for _, attr := range n.Attr {
		if attr.Key == attribute {
			return attr.Val
		}
	}
	return ""
}

func extractText(n *html.Node) string {
	var sb strings.Builder

	// Define a recursive function to traverse child nodes
	var eText func(*html.Node)
	eText = func(n *html.Node) {
		// If the node is a TextNode, append its content
		if n.Type == html.TextNode {
			sb.WriteString(n.Data)
		}
		// Recursively traverse child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			eText(c)
		}
	}

	// Start the text extraction from the root node
	eText(n)

	// Get the concatenated text
	text := sb.String()

	// Strip all whitespace (spaces, newlines, tabs, etc.)
	// First, remove extra newlines, tabs, and replace all spaces with no space
	text = strings.ReplaceAll(text, "\n", "")
	text = strings.ReplaceAll(text, "\t", "")
	text = strings.ReplaceAll(text, " ", "")
	return text
}
